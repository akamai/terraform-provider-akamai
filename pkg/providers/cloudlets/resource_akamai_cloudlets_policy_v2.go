package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type v2PolicyStrategy struct {
	client cloudlets.Cloudlets
}

var cloudletIDs = map[string]int{
	"ER":  0,
	"VP":  1,
	"FR":  3,
	"IG":  4,
	"AP":  5,
	"AS":  6,
	"CD":  7,
	"IV":  8,
	"ALB": 9,
	"MMB": 10,
	"MMA": 11,
}

func (strategy v2PolicyStrategy) createPolicy(ctx context.Context, cloudletName, cloudletCode string, groupID int64) (int64, error) {
	createPolicyReq := cloudlets.CreatePolicyRequest{
		Name:       cloudletName,
		CloudletID: int64(cloudletIDs[cloudletCode]),
		GroupID:    groupID,
	}
	createPolicyResp, err := strategy.client.CreatePolicy(ctx, createPolicyReq)
	if err != nil {
		return 0, err
	}
	return createPolicyResp.PolicyID, nil
}

func (strategy v2PolicyStrategy) updatePolicyVersion(ctx context.Context, d *schema.ResourceData, policyID, version int64, description, matchRulesJSON string, newVersionRequired bool) (error, error) {
	matchRules := make(cloudlets.MatchRules, 0)
	if matchRulesJSON != "" {
		if err := json.Unmarshal([]byte(matchRulesJSON), &matchRules); err != nil {
			return fmt.Errorf("unmarshalling match rules JSON: %s", err), nil
		}
	}

	matchRuleFormat, err := tf.GetStringValue("match_rule_format", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err, nil
	}

	if newVersionRequired {
		createVersionRequest := cloudlets.CreatePolicyVersionRequest{
			CreatePolicyVersion: cloudlets.CreatePolicyVersion{
				MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
				MatchRules:      matchRules,
				Description:     description,
			},
			PolicyID: policyID,
		}
		createVersionResp, err := strategy.client.CreatePolicyVersion(ctx, createVersionRequest)
		if err != nil {
			return err, nil
		}
		if err := d.Set("version", createVersionResp.Version); err != nil {
			return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()), nil
		}
		return setWarnings(d, createVersionResp.Warnings), nil
	}

	updateVersionRequest := cloudlets.UpdatePolicyVersionRequest{
		UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
			MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
			MatchRules:      matchRules,
			Description:     description,
		},
		PolicyID: policyID,
		Version:  version,
	}

	updateVersionResp, err := strategy.client.UpdatePolicyVersion(ctx, updateVersionRequest)
	if err != nil {
		return nil, err
	}
	return setWarnings(d, updateVersionResp.Warnings), nil
}

func (strategy v2PolicyStrategy) updatePolicy(ctx context.Context, policyID, groupID int64, cloudletName string) error {
	updatePolicyReq := cloudlets.UpdatePolicyRequest{
		UpdatePolicy: cloudlets.UpdatePolicy{
			Name:    cloudletName,
			GroupID: groupID,
		},
		PolicyID: policyID,
	}
	_, err := strategy.client.UpdatePolicy(ctx, updatePolicyReq)
	return err
}

func (strategy v2PolicyStrategy) newPolicyVersionIsNeeded(ctx context.Context, policyID, version int64) (bool, error) {
	versionResp, err := strategy.client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID:  policyID,
		Version:   version,
		OmitRules: true,
	})
	if err != nil {
		return false, err
	}

	return len(versionResp.Activations) > 0, nil
}

func (strategy v2PolicyStrategy) readPolicy(ctx context.Context, policyID int64, version *int64) (map[string]any, error) {

	policy, err := strategy.client.GetPolicy(ctx, cloudlets.GetPolicyRequest{PolicyID: policyID})
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]interface{})
	attrs["name"] = policy.Name
	attrs["group_id"] = strconv.FormatInt(policy.GroupID, 10)
	attrs["cloudlet_code"] = policy.CloudletCode
	attrs["cloudlet_id"] = policy.CloudletID
	attrs["is_shared"] = false

	if version == nil {
		return attrs, nil
	}

	policyVersion, err := strategy.client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID: policyID,
		Version:  *version,
	})
	if err != nil {
		return nil, err
	}

	attrs["description"] = policyVersion.Description
	attrs["match_rule_format"] = policyVersion.MatchRuleFormat
	var matchRulesJSON []byte
	if len(policyVersion.MatchRules) > 0 {
		matchRulesJSON, err = json.MarshalIndent(policyVersion.MatchRules, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	attrs["match_rules"] = string(matchRulesJSON)
	attrs["version"] = policyVersion.Version

	return attrs, nil

}

func (strategy v2PolicyStrategy) deletePolicy(ctx context.Context, policyID int64) error {
	policyVersions, err := getAllV2PolicyVersions(ctx, policyID, strategy.client)
	if err != nil {
		return err
	}
	for _, ver := range policyVersions {
		if err := strategy.client.DeletePolicyVersion(ctx, cloudlets.DeletePolicyVersionRequest{
			PolicyID: policyID,
			Version:  ver.Version,
		}); err != nil {
			return err
		}
	}

	activationPending := true
	for activationPending {
		select {
		case <-time.After(DeletionPolicyPollInterval):
			if err = strategy.client.RemovePolicy(ctx, cloudlets.RemovePolicyRequest{PolicyID: policyID}); err != nil {
				statusErr := new(cloudlets.Error)
				// if error does not contain information about pending activations, return it as it is not expected
				if errors.As(err, &statusErr) && !strings.Contains(statusErr.Detail, "Unable to delete policy because an activation for this policy is still pending") {
					return fmt.Errorf("remove policy error: %s", err)
				}
				continue
			}
			activationPending = false
		case <-ctx.Done():
			return fmt.Errorf("retry timeout reached: %s", ctx.Err())
		}
	}
	return err
}

func (strategy v2PolicyStrategy) getVersionStrategy(meta meta.Meta) versionStrategy {
	return v2VersionStrategy{Client(meta)}
}

func (strategy v2PolicyStrategy) setPolicyType(d *schema.ResourceData) error {
	return d.Set("is_shared", false)
}

func (strategy v2PolicyStrategy) isFirstVersionCreated() bool {
	return true
}
