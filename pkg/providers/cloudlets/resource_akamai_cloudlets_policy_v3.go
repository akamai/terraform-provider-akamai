package cloudlets

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cloudlets "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type v3Strategy struct {
	client cloudlets.Cloudlets
}

func (v3 v3Strategy) createPolicy(ctx context.Context, cloudletName string, cloudletCode string, groupID int64) (int64, error, error) {
	createPolicyReq := cloudlets.CreatePolicyRequest{
		Name:         cloudletName,
		CloudletType: cloudlets.CloudletType(cloudletCode),
		GroupID:      groupID,
	}
	createPolicyResp, err := v3.client.CreatePolicy(ctx, createPolicyReq)
	if err != nil {
		return 0, err, nil
	}

	// v2 is creating version implicitly, but v3 does it explicitly
	_, err = v3.client.CreatePolicyVersion(ctx, cloudlets.CreatePolicyVersionRequest{
		PolicyID: createPolicyResp.ID,
		CreatePolicyVersion: cloudlets.CreatePolicyVersion{
			MatchRules: make(cloudlets.MatchRules, 0),
		},
	})
	return createPolicyResp.ID, nil, err
}

func (v3 v3Strategy) updatePolicyVersion(ctx context.Context, d *schema.ResourceData, policyID int64, description string, matchRulesJSON string, version int64, newVersionRequired bool) (error, error) {
	matchRules := make(cloudlets.MatchRules, 0)
	if matchRulesJSON != "" {
		if err := json.Unmarshal([]byte(matchRulesJSON), &matchRules); err != nil {
			return fmt.Errorf("unmarshalling match rules JSON: %s", err), nil
		}
	}

	if newVersionRequired {
		createPolicyReq := cloudlets.CreatePolicyVersionRequest{
			CreatePolicyVersion: cloudlets.CreatePolicyVersion{
				MatchRules:  matchRules,
				Description: tools.StringPtr(description),
			},
			PolicyID: policyID,
		}

		createPolicyRes, err := v3.client.CreatePolicyVersion(ctx, createPolicyReq)
		if err != nil {
			return err, nil
		}
		if err = d.Set("version", createPolicyRes.Version); err != nil {
			return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()), nil
		}
		return setWarnings(d, createPolicyRes.MatchRulesWarnings), nil
	}

	updatePolicyVersionReq := cloudlets.UpdatePolicyVersionRequest{
		UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
			MatchRules:  matchRules,
			Description: tools.StringPtr(description),
		},
		PolicyID: policyID,
		Version:  version,
	}
	updatePolicyRes, err := v3.client.UpdatePolicyVersion(ctx, updatePolicyVersionReq)
	if err != nil {
		return nil, err
	}
	return setWarnings(d, updatePolicyRes.MatchRulesWarnings), nil
}

func (v3 v3Strategy) updatePolicy(ctx context.Context, policyID int64, _ string, groupID int64) error {
	updatePolicyReq := cloudlets.UpdatePolicyRequest{
		PolicyID: policyID,
		BodyParams: cloudlets.UpdatePolicyBodyParams{
			GroupID: groupID,
		},
	}
	_, err := v3.client.UpdatePolicy(ctx, updatePolicyReq)
	return err
}

func (v3 v3Strategy) newPolicyVersionIsNeeded(ctx context.Context, policyID, version int64) (bool, error) {
	policyVersion, err := v3.client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID: policyID,
		Version:  version,
	})
	if err != nil {
		return false, err
	}

	return policyVersion.Immutable, nil

}

func (v3 v3Strategy) readPolicy(ctx context.Context, policyID int64, version *int64) (map[string]any, error) {
	policy, err := v3.client.GetPolicy(ctx, cloudlets.GetPolicyRequest{PolicyID: policyID})
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

	policyVersion, err := v3.client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID: policyID,
		Version:  *version,
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
	attrs["version"] = policyVersion.Version

	return attrs, nil
}

func (v3 v3Strategy) deletePolicy(ctx context.Context, policyID int64) error {
	err := deactivatePolicyVersions(ctx, policyID, v3.client)
	if err != nil {
		return err
	}

	err = v3.client.DeletePolicy(ctx, cloudlets.DeletePolicyRequest{PolicyID: policyID})

	return err
}

func deactivatePolicyVersions(ctx context.Context, policyID int64, client cloudlets.Cloudlets) error {
	policy, err := client.GetPolicy(ctx, cloudlets.GetPolicyRequest{
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
	if policy.CurrentActivations.Staging.Effective != nil && policy.CurrentActivations.Staging.Effective.Operation == cloudlets.OperationActivation {
		_, err := client.DeactivatePolicy(ctx, cloudlets.DeactivatePolicyRequest{
			PolicyID:      policyID,
			Network:       policy.CurrentActivations.Staging.Effective.Network,
			PolicyVersion: int(policy.CurrentActivations.Staging.Effective.PolicyVersion),
		})
		if err != nil {
			return err
		}
		anyDeactivationTriggered = true
	}
	if policy.CurrentActivations.Production.Effective != nil && policy.CurrentActivations.Production.Effective.Operation == cloudlets.OperationActivation {
		_, err := client.DeactivatePolicy(ctx, cloudlets.DeactivatePolicyRequest{
			PolicyID:      policyID,
			Network:       policy.CurrentActivations.Production.Effective.Network,
			PolicyVersion: int(policy.CurrentActivations.Production.Effective.PolicyVersion),
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

func policyHasOngoingActivations(policy *cloudlets.Policy) bool {
	return (policy.CurrentActivations.Staging.Latest != nil && policy.CurrentActivations.Staging.Latest.Status == cloudlets.ActivationStatusInProgress) ||
		(policy.CurrentActivations.Production.Latest != nil && policy.CurrentActivations.Production.Latest.Status == cloudlets.ActivationStatusInProgress)
}

func waitForOngoingActivations(ctx context.Context, policyID int64, client cloudlets.Cloudlets) (*cloudlets.Policy, error) {
	for {
		select {
		case <-time.After(DeletionPolicyPollInterval):
			policy, err := client.GetPolicy(ctx, cloudlets.GetPolicyRequest{PolicyID: policyID})
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

func verifyVersionDeactivated(ctx context.Context, policyID int64, client cloudlets.Cloudlets) (bool, error) {
	policy, err := client.GetPolicy(ctx, cloudlets.GetPolicyRequest{
		PolicyID: policyID,
	})
	if err != nil {
		return false, err
	}

	return (policy.CurrentActivations.Staging.Effective == nil || policy.CurrentActivations.Staging.Effective.Operation == cloudlets.OperationDeactivation) &&
		(policy.CurrentActivations.Production.Effective == nil || policy.CurrentActivations.Production.Effective.Operation == cloudlets.OperationDeactivation), nil
}
