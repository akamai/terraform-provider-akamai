package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudlets/v3"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type v3ActivationStrategy struct {
	client       v3.Cloudlets
	logger       log.Interface
	network      v3.Network
	activationID int64
}

func (strategy *v3ActivationStrategy) setupCloudletSpecificData(rd *schema.ResourceData, network string) error {
	if rd.HasChange("associated_properties") {
		return fmt.Errorf("cannot provide 'associated_properties' for shared policy")
	}
	net, err := strategy.parseNetwork(network)
	if err != nil {
		return err
	}
	strategy.network = net
	return nil
}

func (strategy *v3ActivationStrategy) parseNetwork(network string) (v3.Network, error) {
	switch tf.StateNetwork(strings.ToLower(network)) {
	case "staging":
		return v3.StagingNetwork, nil
	case "production":
		return v3.ProductionNetwork, nil
	}

	return v3.StagingNetwork, fmt.Errorf("'%s' is an invalid network value: should be 'production', 'prod', 'p', 'staging', 'stag' or 's'", network)
}
func (strategy *v3ActivationStrategy) isVersionAlreadyActive(ctx context.Context, policyID, version int64) (bool, string, error) {
	policy, err := strategy.client.GetPolicy(ctx, v3.GetPolicyRequest{
		PolicyID: policyID,
	})
	if err != nil {
		return false, "", err
	}
	if strategy.network == v3.StagingNetwork {
		return policy.CurrentActivations.Staging.Effective != nil &&
			policy.CurrentActivations.Staging.Effective.PolicyVersion == version &&
			policy.CurrentActivations.Staging.Effective.Status == v3.ActivationStatusSuccess &&
			policy.CurrentActivations.Staging.Effective.Operation == v3.OperationActivation, strategy.getID(policyID, strategy.network), nil
	}
	return policy.CurrentActivations.Production.Effective != nil &&
		policy.CurrentActivations.Production.Effective.PolicyVersion == version &&
		policy.CurrentActivations.Production.Effective.Status == v3.ActivationStatusSuccess &&
		policy.CurrentActivations.Production.Effective.Operation == v3.OperationActivation, strategy.getID(policyID, strategy.network), nil
}

func (strategy *v3ActivationStrategy) activateVersion(ctx context.Context, policyID, version int64) error {
	activation, err := strategy.client.ActivatePolicy(ctx, v3.ActivatePolicyRequest{
		PolicyID:      policyID,
		PolicyVersion: version,
		Network:       strategy.network,
	})
	if activation != nil {
		strategy.activationID = activation.ID
	}
	return err
}

func (strategy *v3ActivationStrategy) reactivateVersion(ctx context.Context, policyID, version int64) error {
	return strategy.activateVersion(ctx, policyID, version)
}

func (strategy *v3ActivationStrategy) waitForActivation(ctx context.Context, policyID, _ int64) (string, error) {
	for {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activation, err := strategy.client.GetPolicyActivation(ctx, v3.GetPolicyActivationRequest{
				PolicyID:     policyID,
				ActivationID: strategy.activationID,
			})
			if err != nil {
				return "", err
			}
			if activation != nil {
				switch activation.Status {
				case v3.ActivationStatusSuccess:
					return strategy.getID(policyID, strategy.network), nil
				case v3.ActivationStatusFailed:
					return "", fmt.Errorf("activation failed for policy %d", policyID)
				}
			}
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return "", ErrPolicyActivationTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return "", ErrPolicyActivationCanceled
			}
			return "", fmt.Errorf("%v: %w", ErrPolicyActivationContextTerminated, ctx.Err())
		}
	}
}

func (strategy *v3ActivationStrategy) getID(policyID int64, network v3.Network) string {
	return fmt.Sprintf("%d:%s", policyID, network)
}

func (strategy *v3ActivationStrategy) readActivationFromServer(ctx context.Context, policyID int64, network string) (map[string]any, error) {
	policy, err := strategy.client.GetPolicy(ctx, v3.GetPolicyRequest{
		PolicyID: policyID,
	})
	if err != nil {
		return nil, err
	}

	net, err := strategy.parseNetwork(network)
	if err != nil {
		return nil, err
	}

	switch net {
	case v3.StagingNetwork:
		if policy.CurrentActivations.Staging.Effective != nil {
			return extractAttrsForActivation(policy.CurrentActivations.Staging.Effective), nil
		}
	case v3.ProductionNetwork:
		if policy.CurrentActivations.Production.Effective != nil {
			return extractAttrsForActivation(policy.CurrentActivations.Production.Effective), nil
		}
	}

	return nil, nil
}

func extractAttrsForActivation(effective *v3.PolicyActivation) map[string]any {
	return map[string]any{
		"policy_id": effective.PolicyID,
		"network":   mapNetworkToV2(effective.Network),
		"version":   effective.PolicyVersion,
		"status":    effective.Status,
	}
}

func mapNetworkToV2(network v3.Network) cloudlets.PolicyActivationNetwork {
	if network == v3.StagingNetwork {
		return cloudlets.PolicyActivationNetworkStaging
	}
	return cloudlets.PolicyActivationNetworkProduction
}

func (strategy *v3ActivationStrategy) isReactivationNotNeeded(ctx context.Context, policyID, version int64, _ bool) (bool, string, error) {
	isActive, id, err := strategy.isVersionAlreadyActive(ctx, policyID, version)
	if err != nil {
		return false, "", fmt.Errorf("policy activation update: %w", err)
	}
	return isActive, id, nil
}

func (strategy *v3ActivationStrategy) deactivatePolicy(ctx context.Context, policyID, version int64, network string) error {
	net, err := strategy.parseNetwork(network)
	if err != nil {
		return err
	}
	deactivation, err := strategy.client.DeactivatePolicy(ctx, v3.DeactivatePolicyRequest{
		PolicyID:      policyID,
		PolicyVersion: version,
		Network:       net,
	})
	if err != nil {
		return err
	}
	if deactivation != nil {
		strategy.activationID = deactivation.ID
	}

	_, err = strategy.waitForActivation(ctx, policyID, -1)
	return err
}

func (strategy *v3ActivationStrategy) shouldRetryActivation(err error) bool {
	if err == nil {
		return false
	}
	v3Error := new(v3.Error)
	if errors.As(err, &v3Error) && v3Error.Status >= 500 {
		return true
	}
	return false
}

func (strategy *v3ActivationStrategy) fetchValuesForImport(ctx context.Context, policyID int64, network string) (map[string]any, string, error) {
	net, err := strategy.parseNetwork(network)
	if err != nil {
		return nil, "", err
	}
	policy, err := strategy.client.GetPolicy(ctx, v3.GetPolicyRequest{PolicyID: policyID})
	if err != nil {
		return nil, "", err
	}
	var activationToCheck *v3.PolicyActivation

	switch net {
	case v3.StagingNetwork:
		activationToCheck = policy.CurrentActivations.Staging.Effective
	case v3.ProductionNetwork:
		activationToCheck = policy.CurrentActivations.Production.Effective
	}

	if activationToCheck == nil || activationToCheck.Operation != v3.OperationActivation || activationToCheck.Status != v3.ActivationStatusSuccess {
		return nil, "", fmt.Errorf("no active activation has been found for policy_id: '%d' and network: '%s'", policyID, network)
	}

	return map[string]any{
		"network":   activationToCheck.Network,
		"policy_id": policy.ID,
		"is_shared": true,
	}, strategy.getID(policyID, net), nil
}
