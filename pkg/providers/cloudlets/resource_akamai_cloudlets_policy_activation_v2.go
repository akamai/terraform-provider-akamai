package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type v2ActivationStrategy struct {
	client               cloudlets.Cloudlets
	logger               log.Interface
	network              cloudlets.PolicyActivationNetwork
	associatedProperties []string
	activeProps          []string
	removedProps         []string
}

func (strategy *v2ActivationStrategy) setupCloudletSpecificData(rd *schema.ResourceData, network string) error {
	versionActivationNetwork, err := getPolicyActivationNetwork(network)
	if err != nil {
		return err
	}
	associatedProps, err := tf.GetSetValue("associated_properties", rd)
	if err != nil {
		if errors.Is(err, tf.ErrNotFound) {
			return fmt.Errorf("'associated_properties' is required for non-shared policies")
		}
		return err
	}
	var associatedProperties []string
	for _, prop := range associatedProps.List() {
		associatedProperties = append(associatedProperties, prop.(string))
	}
	sort.Strings(associatedProperties)

	strategy.network = versionActivationNetwork
	strategy.associatedProperties = associatedProperties
	return nil
}

func (strategy *v2ActivationStrategy) isVersionAlreadyActive(ctx context.Context, policyID, version int64) (bool, string, error) {

	strategy.logger.Debugf("checking if policy version %d is active", version)
	policyVersion, err := strategy.client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		Version:   version,
		PolicyID:  policyID,
		OmitRules: true,
	})
	if err != nil {
		return false, "", fmt.Errorf("%s: cannot find the given policy version (%d): %s", ErrPolicyActivation.Error(), version, err.Error())
	}
	policyActivations := sortPolicyActivationsByDate(policyVersion.Activations)

	// just the first activations must correspond to the given properties
	var activeProperties []string
	for _, act := range policyActivations {
		if act.Network == strategy.network &&
			act.PolicyInfo.Status == cloudlets.PolicyActivationStatusActive {
			activeProperties = append(activeProperties, act.PropertyInfo.Name)
		}
	}
	sort.Strings(activeProperties)

	isActive := reflect.DeepEqual(activeProperties, strategy.associatedProperties)
	if isActive {
		strategy.logger.Debugf("policy %d, with version %d and properties [%s], is already active in %s. Fetching all details from server", policyID, version, strings.Join(strategy.associatedProperties, ", "), string(strategy.network))
	}
	return isActive, formatPolicyActivationID(policyID, strategy.network), nil
}

func (strategy *v2ActivationStrategy) activateVersion(ctx context.Context, policyID, version int64) error {
	strategy.logger.Debugf("activating policy %d version %d, network %s and properties [%s]", policyID, version, string(strategy.network), strings.Join(strategy.associatedProperties, ", "))
	_, err := strategy.client.ActivatePolicyVersion(ctx, cloudlets.ActivatePolicyVersionRequest{
		PolicyID: policyID,
		Version:  version,
		Async:    true,
		PolicyVersionActivation: cloudlets.PolicyVersionActivation{
			Network:                 strategy.network,
			AdditionalPropertyNames: strategy.associatedProperties,
		},
	})
	return err
}

func (strategy *v2ActivationStrategy) reactivateVersion(ctx context.Context, policyID, version int64) error {
	// Activate policy version. This will include new associated_properties + the ones which need to be removed
	// it will fail if any of the associated_properties are not valid
	if err := strategy.activateVersion(ctx, policyID, version); err != nil {
		return err
	}

	// 6. remove from the server all unnecessary policy associated_properties
	removedProperties, err := syncToServerRemovedProperties(ctx, strategy.logger, strategy.client, policyID, strategy.network, strategy.activeProps, strategy.associatedProperties)
	strategy.removedProps = removedProperties
	return err
}

func (strategy *v2ActivationStrategy) waitForActivation(ctx context.Context, policyID, version int64) (string, error) {
	act, err := waitForPolicyActivation(ctx, strategy.client, policyID, version, strategy.network, strategy.associatedProperties, strategy.removedProps)
	if err != nil {
		return "", err
	}
	return formatPolicyActivationID(act[0].PolicyInfo.PolicyID, act[0].Network), nil
}

func (strategy *v2ActivationStrategy) readActivationFromServer(ctx context.Context, policyID int64, network string) (map[string]any, error) {
	net, err := getPolicyActivationNetwork(network)
	if err != nil {
		return nil, err
	}

	activations, err := waitForListPolicyActivations(ctx, strategy.client, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  net,
	})
	if err != nil {
		return nil, fmt.Errorf("%v read: %s", ErrPolicyActivation, err.Error())
	}

	if len(activations) == 0 {
		return nil, fmt.Errorf("%v read: cannot find any activation for the given policy (%d) and network ('%s')", ErrPolicyActivation, policyID, net)
	}

	activations = sortPolicyActivationsByDate(activations)
	associatedProperties := getActiveProperties(activations)

	attrs := map[string]any{
		"status":                activations[0].PolicyInfo.Status,
		"version":               activations[0].PolicyInfo.Version,
		"associated_properties": associatedProperties,
	}

	return attrs, nil
}

func (strategy *v2ActivationStrategy) isReactivationNotNeeded(ctx context.Context, policyID, version int64, hasVersionChange bool) (bool, string, error) {
	// policy version validation
	_, err := strategy.client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID:  policyID,
		Version:   version,
		OmitRules: true,
	})
	if err != nil {
		return false, "", fmt.Errorf("%s: cannot find the given policy version (%d): %s", ErrPolicyActivation.Error(), version, err.Error())
	}

	// look for activations with this version which is active in the given network
	activations, err := waitForListPolicyActivations(ctx, strategy.client, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  strategy.network,
	})
	if err != nil {
		return false, "", fmt.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}
	// activations, at this point, contains old and new activations

	// sort by activation date, reverse. To find out the state of the latest activations
	activations = sortPolicyActivationsByDate(activations)

	// find out which properties are activated in those activations
	// version does not matter at this point
	activeProps := getActiveProperties(activations)
	strategy.activeProps = activeProps

	isAlreadyActive := reflect.DeepEqual(activeProps, strategy.associatedProperties) && !hasVersionChange && activations[0].PolicyInfo.Version == version
	return isAlreadyActive, formatPolicyActivationID(policyID, strategy.network), nil
}

func (strategy *v2ActivationStrategy) deactivatePolicy(ctx context.Context, policyID, _ int64, net string) error {
	network, err := getPolicyActivationNetwork(net)
	if err != nil {
		return err
	}

	policyProperties, err := strategy.client.GetPolicyProperties(ctx, cloudlets.GetPolicyPropertiesRequest{PolicyID: policyID})
	if err != nil {
		return fmt.Errorf("%s: cannot find policy %d properties: %s", ErrPolicyActivation.Error(), policyID, err.Error())
	}
	activations, err := waitForListPolicyActivations(ctx, strategy.client, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  network,
	})
	if err != nil {
		return err
	}

	strategy.logger.Debugf("Removing all policy (ID=%d) properties", policyID)
	for propertyName, policyProperty := range policyProperties {
		// filter out property by network
		validProperty := false
		for _, act := range activations {
			if act.PropertyInfo.Name == propertyName {
				validProperty = true
				break
			}
		}
		if !validProperty {
			continue
		}
		// wait for removal until there aren't any pending activations
		if err = waitForNotPendingPolicyActivation(ctx, strategy.logger, strategy.client, policyID, network); err != nil {
			return err
		}

		// proceed to delete property from policy
		err = strategy.client.DeletePolicyProperty(ctx, cloudlets.DeletePolicyPropertyRequest{
			PolicyID:   policyID,
			PropertyID: policyProperty.ID,
			Network:    network,
		})
		if err != nil {
			return fmt.Errorf("%s: cannot delete property '%s' from policy ID %d and network '%s'. Please, try once again later.\n%s", ErrPolicyActivation.Error(), propertyName, policyID, network, err.Error())
		}
	}
	return nil
}

func (strategy *v2ActivationStrategy) shouldRetryActivation(err error) bool {
	return err != nil && policyActivationRetryRegexp.MatchString(strings.ToLower(err.Error()))
}

func (strategy *v2ActivationStrategy) fetchValuesForImport(ctx context.Context, policyID int64, network string) (map[string]any, string, error) {
	activations, err := strategy.client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  cloudlets.PolicyActivationNetwork(network),
	})
	if err != nil {
		return nil, "", err
	}

	var activation *cloudlets.PolicyActivation
	for _, act := range activations {
		if string(act.Network) == network && act.PolicyInfo.Status == cloudlets.PolicyActivationStatusActive {
			activation = &act
			break
		}
	}
	if activation == nil || len(activations) == 0 {
		return nil, "", fmt.Errorf("no active activation has been found for policy_id: '%d' and network: '%s'", policyID, network)
	}

	return map[string]any{
		"network":   activation.Network,
		"policy_id": activation.PolicyInfo.PolicyID,
		"is_shared": false,
	}, fmt.Sprintf("%d:%s", policyID, activation.Network), nil
}
