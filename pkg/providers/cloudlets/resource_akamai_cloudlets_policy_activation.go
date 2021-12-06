package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/apex/log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudletsPolicyActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyActivationCreate,
		ReadContext:   resourcePolicyActivationRead,
		UpdateContext: resourcePolicyActivationUpdate,
		DeleteContext: resourcePolicyActivationDelete,
		Schema:        resourceCloudletsPolicyActivationSchema(),
		Timeouts: &schema.ResourceTimeout{
			Default: &PolicyActivationResourceTimeout,
		},
	}
}

func resourceCloudletsPolicyActivationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Activation status for this Cloudlets policy",
		},
		"policy_id": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "ID of the Cloudlets policy you want to activate",
			ForceNew:    true,
		},
		"network": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: tools.ValidateNetwork,
			StateFunc:        statePolicyActivationNetwork,
			Description:      "The network you want to activate the policy version on (options are Staging and Production)",
		},
		"version": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Cloudlets policy version you want to activate",
		},
		"associated_properties": {
			Type:        schema.TypeSet,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Set of property IDs to link to this Cloudlets policy",
		},
	}
}

var (
	// ActivationPollMinimum is the minimum polling interval for activation creation
	ActivationPollMinimum = time.Minute

	// ActivationPollInterval is the interval for polling an activation status on creation
	ActivationPollInterval = ActivationPollMinimum

	// PolicyActivationResourceTimeout is the default timeout for the resource operations
	PolicyActivationResourceTimeout = time.Minute * 90

	// ErrNetworkName is used when the user inputs an invalid network name
	ErrNetworkName = errors.New("invalid network name")
)

func resourcePolicyActivationDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationDelete")
	logger.Debug("Deleting cloudlets policy activation")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	pID, err := tools.GetIntValue("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := int64(pID)

	network, err := getPolicyActivationNetwork(strings.Split(rd.Id(), ":")[1])
	if err != nil {
		return diag.FromErr(err)
	}

	policyProperties, err := client.GetPolicyProperties(ctx, cloudlets.GetPolicyPropertiesRequest{PolicyID: policyID})
	if err != nil {
		return diag.Errorf("%s: cannot find policy %d properties: %s", ErrPolicyActivation.Error(), policyID, err.Error())
	}
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  network,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	logger.Debugf("Removing all policy (ID=%d) properties", policyID)
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
		if err = waitForNotPendingPolicyActivation(ctx, logger, client, policyID, network); err != nil {
			return diag.FromErr(err)
		}

		// proceed to delete property from policy
		err = client.DeletePolicyProperty(ctx, cloudlets.DeletePolicyPropertyRequest{
			PolicyID:   policyID,
			PropertyID: policyProperty.ID,
			Network:    network,
		})
		if err != nil {
			return diag.Errorf("%s: cannot delete property '%s' from policy ID %d and network '%s'. Please, try once again later.\n%s", ErrPolicyActivation.Error(), propertyName, policyID, network, err.Error())
		}
	}
	logger.Debugf("All properties have been removed from policy ID %d", policyID)
	rd.SetId("")
	return nil
}

func resourcePolicyActivationUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationUpdate")

	// 1. check if version has changed.
	if !rd.HasChanges("version", "associated_properties") {
		logger.Debugf("nothing to update")
		return nil
	}

	logger.Debugf("proceeding to create and activate a new policy activation version")

	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	// 2. In such case, create a new version to activate (for creation, look into resource policy)
	policyID, err := tools.GetIntValue("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	activationNetwork, err := getPolicyActivationNetwork(network)
	if err != nil {
		return diag.FromErr(err)
	}

	v, err := tools.GetIntValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	version := int64(v)
	// policy version validation
	_, err = client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID:  int64(policyID),
		Version:   version,
		OmitRules: true,
	})
	if err != nil {
		if diagnostics := tools.RestoreOldValues(rd, []string{"version", "associated_properties"}); diagnostics != nil {
			return diagnostics
		}
		return diag.Errorf("%s: cannot find the given policy version (%d): %s", ErrPolicyActivation.Error(), version, err.Error())
	}

	associatedProps, err := tools.GetSetValue("associated_properties", rd)
	if errors.Is(err, tools.ErrNotFound) || len(associatedProps.List()) == 0 {
		if diagnostics := tools.RestoreOldValues(rd, []string{"version", "associated_properties"}); diagnostics != nil {
			return diagnostics
		}
		return diag.Errorf("Field associated_properties should not be empty. If you want to remove all policy associated properties, please run `terraform destroy` instead.")
	}
	if err != nil {
		return diag.FromErr(err)
	}
	var newPolicyProperties []string
	for _, prop := range associatedProps.List() {
		newPolicyProperties = append(newPolicyProperties, prop.(string))
	}
	sort.Strings(newPolicyProperties)

	// 3. look for activations with this version which is active in the given network
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: int64(policyID),
		Network:  activationNetwork,
	})
	if err != nil {
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}
	// activations, at this point, contains old and new activations

	// sort by activation date, reverse. To find out the state of the latest activations
	activations = sortPolicyActivationsByDate(activations)

	// find out which properties are activated in those activations
	// version does not matter at this point
	activeProps := getActiveProperties(activations)

	// 4. all "additional_properties" are active for the given version, policyID and network, proceed to read stage
	if reflect.DeepEqual(activeProps, newPolicyProperties) && !rd.HasChanges("version") {
		// in such case, return
		logger.Debugf("This policy (ID=%d, version=%d) is already active.", policyID, version)
		return resourcePolicyActivationRead(ctx, rd, m)
	}

	// 5. Activate policy version. This will include new associated_properties + the ones which need to be removed
	// it will fail if any of the associated_properties are not valid
	logger.Debugf("Proceeding to activate the policy ID=%d (version=%d, properties=[%s], network='%s') is not active.",
		policyID, version, strings.Join(newPolicyProperties, ", "), activationNetwork)

	err = client.ActivatePolicyVersion(ctx, cloudlets.ActivatePolicyVersionRequest{
		PolicyID: int64(policyID),
		Async:    true,
		Version:  version,
		PolicyVersionActivation: cloudlets.PolicyVersionActivation{
			Network:                 activationNetwork,
			AdditionalPropertyNames: newPolicyProperties,
		},
	})
	if err != nil {
		if diagnostics := tools.RestoreOldValues(rd, []string{"version", "associated_properties"}); diagnostics != nil {
			return diagnostics
		}
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}

	// 6. remove from the server all unnecessary policy associated_properties
	removedProperties, err := syncToServerRemovedProperties(ctx, logger, client, int64(policyID), activationNetwork, activeProps, newPolicyProperties)
	if err != nil {
		return diag.FromErr(err)
	}

	// 7. poll until active
	_, err = waitForPolicyActivation(ctx, client, int64(policyID), version, activationNetwork, newPolicyProperties, removedProperties)
	if err != nil {
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}

	return resourcePolicyActivationRead(ctx, rd, m)
}

func resourcePolicyActivationCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Creating policy activation")

	policyID, err := tools.GetIntValue("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	versionActivationNetwork, err := getPolicyActivationNetwork(network)
	if err != nil {
		return diag.FromErr(err)
	}
	associatedProps, err := tools.GetSetValue("associated_properties", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	var associatedProperties []string
	for _, prop := range associatedProps.List() {
		associatedProperties = append(associatedProperties, prop.(string))
	}
	sort.Strings(associatedProperties)

	v, err := tools.GetIntValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	version := int64(v)

	logger.Debugf("checking if policy version %d is active", version)
	policyVersion, err := client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		Version:   version,
		PolicyID:  int64(policyID),
		OmitRules: true,
	})
	if err != nil {
		return diag.Errorf("%s: cannot find the given policy version (%d): %s", ErrPolicyActivation.Error(), version, err.Error())
	}
	policyActivations := sortPolicyActivationsByDate(policyVersion.Activations)

	// just the first activations must correspond to the given properties
	var activeProperties []string
	for _, act := range policyActivations {
		if act.Network == versionActivationNetwork &&
			act.PolicyInfo.Status == cloudlets.PolicyActivationStatusActive {
			activeProperties = append(activeProperties, act.PropertyInfo.Name)
		}
	}
	sort.Strings(activeProperties)
	if reflect.DeepEqual(activeProperties, associatedProperties) {
		// if the given version is active, just refresh status and quit
		logger.Debugf("policy %d, with version %d and properties [%s], is already active in %s. Fetching all details from server", policyID, version, strings.Join(associatedProperties, ", "), string(versionActivationNetwork))
		rd.SetId(formatPolicyActivationID(int64(policyID), cloudlets.PolicyActivationNetwork(network)))
		return resourcePolicyActivationRead(ctx, rd, m)
	}

	// at this point, we are sure that the given version is not active
	logger.Debugf("activating policy %d version %d, network %s and properties [%s]", policyID, version, string(versionActivationNetwork), strings.Join(associatedProperties, ", "))
	err = client.ActivatePolicyVersion(ctx, cloudlets.ActivatePolicyVersionRequest{
		PolicyID: int64(policyID),
		Version:  version,
		Async:    true,
		PolicyVersionActivation: cloudlets.PolicyVersionActivation{
			Network:                 versionActivationNetwork,
			AdditionalPropertyNames: associatedProperties,
		},
	})
	if err != nil {
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}

	// wait until policy activation is done
	act, err := waitForPolicyActivation(ctx, client, int64(policyID), version, versionActivationNetwork, associatedProperties, nil)
	if err != nil {
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}
	rd.SetId(formatPolicyActivationID(act[0].PolicyInfo.PolicyID, act[0].Network))

	return resourcePolicyActivationRead(ctx, rd, m)
}

func resourcePolicyActivationRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading policy activations")

	policyID, err := tools.GetIntValue("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	net, err := getPolicyActivationNetwork(network)
	if err != nil {
		return diag.FromErr(err)
	}

	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: int64(policyID),
		Network:  net,
	})
	if err != nil {
		return diag.Errorf("%v read: %s", ErrPolicyActivation, err.Error())
	}

	if len(activations) == 0 {
		return diag.Errorf("%v read: cannot find any activation for the given policy (%d) and network ('%s')", ErrPolicyActivation, policyID, net)
	}

	activations = sortPolicyActivationsByDate(activations)

	if err := rd.Set("status", activations[0].PolicyInfo.Status); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := rd.Set("version", activations[0].PolicyInfo.Version); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	associatedProperties := getActiveProperties(activations)
	if err := rd.Set("associated_properties", associatedProperties); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func formatPolicyActivationID(policyID int64, network cloudlets.PolicyActivationNetwork) string {
	return fmt.Sprintf("%d:%s", policyID, network)
}

func getActiveProperties(policyActivations []cloudlets.PolicyActivation) []string {
	var activeProps []string
	for _, act := range policyActivations {
		if act.PolicyInfo.Status == cloudlets.PolicyActivationStatusActive {
			activeProps = append(activeProps, act.PropertyInfo.Name)
		}
	}
	sort.Strings(activeProps)
	return activeProps
}

// waitForPolicyActivation polls server until the activation has active status or until context is closed (because of timeout, cancellation or context termination)
func waitForPolicyActivation(ctx context.Context, client cloudlets.Cloudlets, policyID, version int64, network cloudlets.PolicyActivationNetwork, additionalProps, removedProperties []string) ([]cloudlets.PolicyActivation, error) {
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  network,
	})
	if err != nil {
		return nil, err
	}
	activations = filterActivations(activations, version, additionalProps)
	for len(activations) > 0 {
		allActive, allRemoved := true, true
	activations:
		for _, act := range activations {
			if act.PolicyInfo.Version == version {
				if act.PolicyInfo.Status == cloudlets.PolicyActivationStatusFailed ||
					strings.Contains(act.PolicyInfo.StatusDetail, "fail") {
					return nil, fmt.Errorf("%v: policyID %d activation failure: %s", ErrPolicyActivation, act.PolicyInfo.PolicyID, act.PolicyInfo.StatusDetail)
				}
				if act.PolicyInfo.Status != cloudlets.PolicyActivationStatusActive {
					allActive = false
					break
				}
			}
			for _, property := range removedProperties {
				if property == act.PropertyInfo.Name {
					allRemoved = false
					break activations
				}
			}
		}
		if allActive && allRemoved {
			return activations, nil
		}
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activations, err = client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
				PolicyID: policyID,
				Network:  network,
			})
			if err != nil {
				return nil, err
			}
			activations = filterActivations(activations, version, additionalProps)

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, ErrPolicyActivationTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil, ErrPolicyActivationCanceled
			}
			return nil, fmt.Errorf("%v: %w", ErrPolicyActivationContextTerminated, ctx.Err())
		}
	}

	if len(activations) == 0 {
		return nil, fmt.Errorf("%v: policyID %d: not all properties are active", ErrPolicyActivation, policyID)
	}

	return activations, nil
}

// filterActivations filters the latest activation for the given properties and version. In case of length mismatch (not all
// properties present in the last activation): it returns nil.
func filterActivations(activations []cloudlets.PolicyActivation, version int64, properties []string) []cloudlets.PolicyActivation {
	// inverse sorting by activation date -> first activations will be the most recent
	activations = sortPolicyActivationsByDate(activations)
	var lastActivationBlock []cloudlets.PolicyActivation
	var lastActivationDate int64
	// collect lastActivationBlock slice, with all activations sharing the latest activation date
	for _, act := range activations {
		// Each call to cloudlets.ActivatePolicyVersion() will result in a different activation date, and each activated
		// property will have the same activation date.
		if lastActivationDate != 0 && lastActivationDate != act.PolicyInfo.ActivationDate {
			break
		}
		lastActivationDate = act.PolicyInfo.ActivationDate
		lastActivationBlock = append(lastActivationBlock, act)
	}
	// find out if the all given properties were activated with the given policy version in last activation date
	allPropertiesActive := true
	for _, name := range properties {
		propertyPresent := false
		for _, act := range lastActivationBlock {
			if act.PropertyInfo.Name == name && act.PolicyInfo.Version == version {
				propertyPresent = true
				break
			}
		}
		if !propertyPresent {
			allPropertiesActive = false
			break
		}
	}
	if !allPropertiesActive {
		return nil
	}
	return lastActivationBlock
}

func sortPolicyActivationsByDate(activations []cloudlets.PolicyActivation) []cloudlets.PolicyActivation {
	sort.Slice(activations, func(i, j int) bool {
		return activations[i].PolicyInfo.ActivationDate > activations[j].PolicyInfo.ActivationDate
	})
	return activations
}

func getPolicyActivationNetwork(net string) (cloudlets.PolicyActivationNetwork, error) {

	net = tools.StateNetwork(net)

	switch net {
	case "production":
		return cloudlets.PolicyActivationNetworkProduction, nil
	case "staging":
		return cloudlets.PolicyActivationNetworkStaging, nil
	}

	return "", ErrNetworkName
}

func statePolicyActivationNetwork(i interface{}) string {

	net := tools.StateNetwork(i)

	switch net {
	case "production":
		return string(cloudlets.PolicyActivationNetworkProduction)
	case "staging":
		return string(cloudlets.PolicyActivationNetworkStaging)
	}

	// this should never happen :-)
	return net
}

func syncToServerRemovedProperties(ctx context.Context, logger log.Interface, client cloudlets.Cloudlets, policyID int64, network cloudlets.PolicyActivationNetwork, activeProps, newPolicyProperties []string) ([]string, error) {
	policyProperties, err := client.GetPolicyProperties(ctx, cloudlets.GetPolicyPropertiesRequest{PolicyID: policyID})
	if err != nil {
		return nil, fmt.Errorf("%w: cannot find policy %d properties: %s", ErrPolicyActivation, policyID, err.Error())
	}
	removedProperties := make([]string, 0)
activePropertiesLoop:
	for _, activeProp := range activeProps {
		for _, newProp := range newPolicyProperties {
			if activeProp == newProp {
				continue activePropertiesLoop
			}
		}
		// find out property id
		associateProperty, ok := policyProperties[activeProp]
		if !ok {
			logger.Warnf("Policy %d server side discrepancies: '%s' is not present in GetPolicyProperties response", policyID, activeProp)
			continue activePropertiesLoop
		}
		propertyID := associateProperty.ID

		// wait for removal until there aren't any pending activations
		if err = waitForNotPendingPolicyActivation(ctx, logger, client, policyID, network); err != nil {
			return nil, err
		}

		// remove property from policy
		logger.Debugf("proceeding to delete property '%s' from policy (ID=%d)", activeProp, policyID)
		if err := client.DeletePolicyProperty(ctx, cloudlets.DeletePolicyPropertyRequest{PolicyID: policyID, PropertyID: propertyID, Network: network}); err != nil {
			return nil, fmt.Errorf("%w: cannot remove policy %d property %d and network '%s'. Please, try once again later.\n%s", ErrPolicyActivation, policyID, propertyID, network, err.Error())
		}
		removedProperties = append(removedProperties, activeProp)
	}

	// wait for removal until there aren't any pending activations
	if err = waitForNotPendingPolicyActivation(ctx, logger, client, policyID, network); err != nil {
		return nil, err
	}

	// at this point, there are no activations in pending state
	return removedProperties, nil
}

func waitForNotPendingPolicyActivation(ctx context.Context, logger log.Interface, client cloudlets.Cloudlets, policyID int64, network cloudlets.PolicyActivationNetwork) error {
	logger.Debugf("waiting until there none of the policy (ID=%d) activations are in pending state", policyID)
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{PolicyID: policyID})
	if err != nil {
		return fmt.Errorf("%w: failed to list policy activations for policy %d: %s", ErrPolicyActivation, policyID, err.Error())
	}
	for len(activations) > 0 {
		pending := false
		for _, act := range activations {
			if act.PolicyInfo.Status == cloudlets.PolicyActivationStatusFailed {
				return fmt.Errorf("%v: policyID %d: %s", ErrPolicyActivation, act.PolicyInfo.PolicyID, act.PolicyInfo.StatusDetail)
			}
			if act.PolicyInfo.Status == cloudlets.PolicyActivationStatusPending {
				pending = true
				break
			}
		}
		if !pending {
			break
		}
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activations, err = client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
				PolicyID: policyID,
				Network:  network,
			})
			if err != nil {
				return fmt.Errorf("%w: failed to list policy activations for policy %d: %s", ErrPolicyActivation, policyID, err.Error())
			}

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return ErrPolicyActivationTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return ErrPolicyActivationCanceled
			}
			return fmt.Errorf("%v: %w", ErrPolicyActivationContextTerminated, ctx.Err())
		}
	}

	return nil
}
