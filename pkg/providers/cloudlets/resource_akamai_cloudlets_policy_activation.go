package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

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
			Optional:    true,
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

func resourcePolicyActivationDelete(_ context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationDelete")
	logger.Debug("Deleting cloudlets policy activation from local schema only")
	logger.Info("Cloudlets API does not support policy activation version deletion - resource will only be removed from state")
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

	// 3. look for activation with this version which is active
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: int64(policyID),
		Network:  activationNetwork,
	})
	if err != nil {
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}
	activeProps := getAssociatedProperties(activations, version)

	associatedProps, err := tools.GetSetValue("associated_properties", rd)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	additionalProps := []string{}
	for _, prop := range associatedProps.List() {
		additionalProps = append(additionalProps, prop.(string))
	}
	sort.Strings(additionalProps)
	// find out if there are activations with status==active and network==activationNetwork
	var active bool
	for _, act := range activations {
		if act.PolicyInfo.Status == cloudlets.StatusActive {
			if act.PolicyInfo.Version == version && act.Network == activationNetwork {
				active = true
				break
			}
		}
	}
	if active && reflect.DeepEqual(activeProps, additionalProps) {
		// in such case, return
		logger.Debugf("This policy (ID=%d, version=%d) is already active.", policyID, version)
		return resourcePolicyActivationRead(ctx, rd, m)
	}

	logger.Debugf("This policy (ID=%d, version=%d) is not active in '%s' network. Proceeding to activation.", policyID, version, activationNetwork)

	err = client.ActivatePolicyVersion(ctx, cloudlets.ActivatePolicyVersionRequest{
		PolicyID: int64(policyID),
		Async:    true,
		Version:  version,
		RequestBody: cloudlets.ActivatePolicyVersionRequestBody{
			Network:                 activationNetwork,
			AdditionalPropertyNames: additionalProps,
		},
	})
	if err != nil {
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}

	// 4. poll until active
	_, err = waitForPolicyActivation(ctx, client, int64(policyID), version, activationNetwork)
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
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	additionalProps := []string{}
	for _, prop := range associatedProps.List() {
		additionalProps = append(additionalProps, prop.(string))
	}
	sort.Strings(additionalProps)

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
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}
	activeProperties := []string{}
	var policyVersionActivation *cloudlets.Activation
	for _, act := range policyVersion.Activations {
		if cloudlets.VersionActivationNetwork(act.Network) == versionActivationNetwork &&
			act.PolicyInfo.Status == cloudlets.StatusActive {
			activeProperties = append(activeProperties, act.PropertyInfo.Name)
			policyVersionActivation = act
		}
	}
	sort.Strings(activeProperties)
	if reflect.DeepEqual(activeProperties, additionalProps) {
		// if the given version is active, just refresh status and quit
		logger.Debugf("policy version %d is already active in %s, fetching all details from server", version, string(versionActivationNetwork))
		rd.SetId(formatPolicyActivationID(policyVersionActivation.PolicyInfo))
		return resourcePolicyActivationRead(ctx, rd, m)
	}

	// at this point, we are sure that the given version is not active
	logger.Debugf("activating policy version %d for policy %d", version, policyID)
	err = client.ActivatePolicyVersion(ctx, cloudlets.ActivatePolicyVersionRequest{
		PolicyID: int64(policyID),
		Version:  version,
		Async:    true,
		RequestBody: cloudlets.ActivatePolicyVersionRequestBody{
			Network:                 versionActivationNetwork,
			AdditionalPropertyNames: additionalProps,
		},
	})
	if err != nil {
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}

	// wait until policy activation is done
	act, err := waitForPolicyActivation(ctx, client, int64(policyID), version, versionActivationNetwork)
	if err != nil {
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}
	rd.SetId(formatPolicyActivationID(act[0].PolicyInfo))

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
	v, err := tools.GetIntValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version := int64(v)

	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: int64(policyID),
		Network:  net,
	})
	if err != nil {
		return diag.Errorf("%v read: %s", ErrPolicyActivation, err.Error())
	}

	associatedProperties := getAssociatedProperties(activations, version)

	for _, act := range activations {
		if act.PolicyInfo.Version == version {
			if err := rd.Set("status", act.PolicyInfo.Status); err != nil {
				return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
			}
			if err := rd.Set("version", version); err != nil {
				return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
			}
			if err := rd.Set("associated_properties", associatedProperties); err != nil {
				return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
			}

			return nil
		}
	}

	return diag.Errorf("%v read: cannot find the given policy activation version (%d)", ErrPolicyActivation, version)
}

func formatPolicyActivationID(policyInfo cloudlets.PolicyInfo) string {
	return fmt.Sprintf("%d:%d", policyInfo.PolicyID, policyInfo.ActivationDate)
}

func getAssociatedProperties(policyActivations []cloudlets.PolicyActivation, version int64) []string {
	activeProps := []string{}
	for _, act := range policyActivations {
		if act.PolicyInfo.Status == cloudlets.StatusActive && act.PolicyInfo.Version == version {
			activeProps = append(activeProps, act.PropertyInfo.Name)
		}
	}
	sort.Strings(activeProps)
	return activeProps
}

// getActivePolicyActivations gets active policy activations for given policy id, version and network from the server
func getActivePolicyActivations(ctx context.Context, client cloudlets.Cloudlets, policyID, version int64, network cloudlets.VersionActivationNetwork) ([]cloudlets.PolicyActivation, error) {
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  network,
	})
	if err != nil {
		return nil, err
	}

	policyActivations := []cloudlets.PolicyActivation{}

	for _, act := range activations {
		if act.PolicyInfo.Version == version && act.PolicyInfo.Status == cloudlets.StatusActive {
			policyActivations = append(policyActivations, act)
		}
	}

	if len(policyActivations) > 0 {
		return policyActivations, nil
	}

	return nil, fmt.Errorf("%v: no activations found for given version", ErrPolicyActivation)
}

// waitForPolicyActivation polls server until the activation has active status or until context is closed (because of timeout, cancellation or context termination)
func waitForPolicyActivation(ctx context.Context, client cloudlets.Cloudlets, policyID, version int64, network cloudlets.VersionActivationNetwork) ([]cloudlets.PolicyActivation, error) {
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: policyID,
		Network:  network,
	})
	if err != nil {
		return nil, err
	}
	for len(activations) > 0 {
		allActive := true
		for _, act := range activations {
			if act.PolicyInfo.Version == version {
				if act.PolicyInfo.Status != cloudlets.StatusActive {
					allActive = false
				}
				if act.PolicyInfo.Status == cloudlets.StatusFailed {
					return nil, fmt.Errorf("%v: policyID %d: %s", ErrPolicyActivation, act.PolicyInfo.PolicyID, act.PolicyInfo.StatusDetail)
				}
			}
		}
		if allActive {
			return activations, nil
		}
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activations, err = getActivePolicyActivations(ctx, client, policyID, version, network)
			if err != nil {
				return nil, err
			}

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

	return activations, nil
}

func getPolicyActivationNetwork(net string) (cloudlets.VersionActivationNetwork, error) {

	net = tools.StateNetwork(net)

	switch net {
	case "production":
		return cloudlets.VersionActivationNetworkProduction, nil
	case "staging":
		return cloudlets.VersionActivationNetworkStaging, nil
	}

	return "", ErrNetworkName
}

func statePolicyActivationNetwork(i interface{}) string {

	net := tools.StateNetwork(i)

	switch net {
	case "production":
		return string(cloudlets.VersionActivationNetworkProduction)
	case "staging":
		return string(cloudlets.VersionActivationNetworkStaging)
	}

	// this should never happen :-)
	return net
}
