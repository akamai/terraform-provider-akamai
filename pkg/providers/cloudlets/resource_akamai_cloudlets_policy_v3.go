package cloudlets

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type v3PolicyStrategy struct {
	client v3.Cloudlets
}

func (strategy v3PolicyStrategy) createPolicy(ctx context.Context, cloudletName, cloudletCode string, groupID int64) (int64, error) {
	createPolicyReq := v3.CreatePolicyRequest{
		Name:         cloudletName,
		CloudletType: v3.CloudletType(cloudletCode),
		GroupID:      groupID,
	}
	createPolicyResp, err := strategy.client.CreatePolicy(ctx, createPolicyReq)
	if err != nil {
		return 0, err
	}

	return createPolicyResp.ID, nil
}

func (strategy v3PolicyStrategy) updatePolicyVersion(ctx context.Context, d *schema.ResourceData, policyID, version int64, description, matchRulesJSON string, newVersionRequired bool) (error, error) {
	matchRules := make(v3.MatchRules, 0)
	if matchRulesJSON != "" {
		if err := json.Unmarshal([]byte(matchRulesJSON), &matchRules); err != nil {
			return fmt.Errorf("unmarshalling match rules JSON: %s", err), nil
		}
	}

	if newVersionRequired {
		createPolicyReq := v3.CreatePolicyVersionRequest{
			CreatePolicyVersion: v3.CreatePolicyVersion{
				MatchRules:  matchRules,
				Description: ptr.To(description),
			},
			PolicyID: policyID,
		}

		createPolicyRes, err := strategy.client.CreatePolicyVersion(ctx, createPolicyReq)
		if err != nil {
			return nil, err
		}
		if err = d.Set("version", createPolicyRes.PolicyVersion); err != nil {
			return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()), nil
		}
		return setWarnings(d, createPolicyRes.MatchRulesWarnings), nil
	}

	updatePolicyVersionReq := v3.UpdatePolicyVersionRequest{
		UpdatePolicyVersion: v3.UpdatePolicyVersion{
			MatchRules:  matchRules,
			Description: ptr.To(description),
		},
		PolicyID:      policyID,
		PolicyVersion: version,
	}
	updatePolicyRes, err := strategy.client.UpdatePolicyVersion(ctx, updatePolicyVersionReq)
	if err != nil {
		return nil, err
	}
	return setWarnings(d, updatePolicyRes.MatchRulesWarnings), nil
}

func (strategy v3PolicyStrategy) updatePolicy(ctx context.Context, policyID, groupID int64, _ string) error {
	updatePolicyReq := v3.UpdatePolicyRequest{
		PolicyID: policyID,
		Body: v3.UpdatePolicyRequestBody{
			GroupID: groupID,
		},
	}
	_, err := strategy.client.UpdatePolicy(ctx, updatePolicyReq)
	return err
}

func (strategy v3PolicyStrategy) newPolicyVersionIsNeeded(ctx context.Context, policyID, version int64) (bool, error) {
	policyVersion, err := strategy.client.GetPolicyVersion(ctx, v3.GetPolicyVersionRequest{
		PolicyID:      policyID,
		PolicyVersion: version,
	})
	if err != nil {
		return false, err
	}

	return policyVersion.Immutable, nil

}

func (strategy v3PolicyStrategy) readPolicy(ctx context.Context, policyID int64, version *int64) (map[string]any, error) {
	policy, err := strategy.client.GetPolicy(ctx, v3.GetPolicyRequest{PolicyID: policyID})
	if err != nil {
		return nil, err
	}

	attrs := make(map[string]any)
	attrs["name"] = policy.Name
	attrs["group_id"] = strconv.FormatInt(policy.GroupID, 10)
	attrs["cloudlet_code"] = policy.CloudletType
	attrs["is_shared"] = true

	if version == nil {
		attrs["description"] = ""
		attrs["match_rules"] = ""
		return attrs, nil
	}

	policyVersion, err := strategy.client.GetPolicyVersion(ctx, v3.GetPolicyVersionRequest{
		PolicyID:      policyID,
		PolicyVersion: *version,
	})
	if err != nil {
		return nil, err
	}

	attrs["description"] = policyVersion.Description
	var matchRulesJSON []byte
	if len(policyVersion.MatchRules) > 0 {
		matchRulesJSON, err = json.MarshalIndent(policyVersion.MatchRules, "", "  ")
		if err != nil {
			return nil, err
		}
	}
	attrs["match_rules"] = string(matchRulesJSON)
	attrs["version"] = policyVersion.PolicyVersion

	return attrs, nil
}

func (strategy v3PolicyStrategy) deletePolicy(ctx context.Context, policyID int64) error {
	err := deactivatePolicyVersions(ctx, policyID, strategy.client)
	if err != nil {
		return err
	}

	err = strategy.client.DeletePolicy(ctx, v3.DeletePolicyRequest{PolicyID: policyID})

	return err
}

func (strategy v3PolicyStrategy) getVersionStrategy(meta meta.Meta) versionStrategy {
	return v3VersionStrategy{ClientV3(meta)}
}

func (strategy v3PolicyStrategy) setPolicyType(d *schema.ResourceData) error {
	return d.Set("is_shared", true)
}

func deactivatePolicyVersions(ctx context.Context, policyID int64, client v3.Cloudlets) error {
	policy, err := client.GetPolicy(ctx, v3.GetPolicyRequest{
		PolicyID: policyID,
	})
	if err != nil {
		return err
	}

	if policyHasOngoingActivations(policy) {
		policy, err = waitForOngoingActivations(ctx, policyID, client)
		if err != nil {
			return err
		}
	}

	anyDeactivationTriggered := false
	if policy.CurrentActivations.Staging.Effective != nil && policy.CurrentActivations.Staging.Effective.Operation == v3.OperationActivation {
		_, err := client.DeactivatePolicy(ctx, v3.DeactivatePolicyRequest{
			PolicyID:      policyID,
			Network:       policy.CurrentActivations.Staging.Effective.Network,
			PolicyVersion: policy.CurrentActivations.Staging.Effective.PolicyVersion,
		})
		if err != nil {
			return err
		}
		anyDeactivationTriggered = true
	}
	if policy.CurrentActivations.Production.Effective != nil && policy.CurrentActivations.Production.Effective.Operation == v3.OperationActivation {
		_, err := client.DeactivatePolicy(ctx, v3.DeactivatePolicyRequest{
			PolicyID:      policyID,
			Network:       policy.CurrentActivations.Production.Effective.Network,
			PolicyVersion: policy.CurrentActivations.Production.Effective.PolicyVersion,
		})
		if err != nil {
			return err
		}
		anyDeactivationTriggered = true
	}

	if !anyDeactivationTriggered {
		return nil
	}

	for {
		select {
		case <-time.After(DeletionPolicyPollInterval):
			everythingDeactivated, err := verifyVersionDeactivated(ctx, policyID, client)
			if err != nil {
				return err
			}
			if everythingDeactivated {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("retry timeout reached: %s", ctx.Err())
		}
	}
}

func policyHasOngoingActivations(policy *v3.Policy) bool {
	return (policy.CurrentActivations.Staging.Latest != nil && policy.CurrentActivations.Staging.Latest.Status == v3.ActivationStatusInProgress) ||
		(policy.CurrentActivations.Production.Latest != nil && policy.CurrentActivations.Production.Latest.Status == v3.ActivationStatusInProgress)
}

func waitForOngoingActivations(ctx context.Context, policyID int64, client v3.Cloudlets) (*v3.Policy, error) {
	for {
		select {
		case <-time.After(DeletionPolicyPollInterval):
			policy, err := client.GetPolicy(ctx, v3.GetPolicyRequest{PolicyID: policyID})
			if err != nil {
				return nil, err
			}
			if !policyHasOngoingActivations(policy) {
				return policy, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("retry timeout reached: %s", ctx.Err())
		}
	}
}

func verifyVersionDeactivated(ctx context.Context, policyID int64, client v3.Cloudlets) (bool, error) {
	policy, err := client.GetPolicy(ctx, v3.GetPolicyRequest{
		PolicyID: policyID,
	})
	if err != nil {
		return false, err
	}

	return (policy.CurrentActivations.Staging.Effective == nil || policy.CurrentActivations.Staging.Effective.Operation == v3.OperationDeactivation) &&
		(policy.CurrentActivations.Production.Effective == nil || policy.CurrentActivations.Production.Effective.Operation == v3.OperationDeactivation), nil
}

func (strategy v3PolicyStrategy) isFirstVersionCreated() bool {
	return false
}
