package cloudlets

import (
	"context"
	"errors"
	"fmt"
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
		},
		"network": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: tools.ValidateNetwork,
			StateFunc:        tools.StateNetwork,
			Description:      "The network you want to activate the policy version on (options are Staging and Production)",
		},
		"version": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Cloudlets policy version you want to activate",
		},
		"associated_properties": {
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of property IDs to link to this Cloudlets policy",
		},
	}
}

const (
	// ActivationPollMinimum is the minimum polling interval for activation creation
	ActivationPollMinimum = time.Minute
)

var (
	// ActivationPollInterval is the interval for polling an activation status on creation
	ActivationPollInterval = ActivationPollMinimum

	// PolicyActivationResourceTimeout is the default timeout for the resource operations
	PolicyActivationResourceTimeout = time.Minute * 90
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
	if !rd.HasChanges("version") {
		logger.Debugf("version number has not changed, nothing to update")
		return nil
	}

	logger.Debugf("version number has changed: proceed to create and activate a new policy activation version")

	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	// 2. In such case, create a new version to activate (for creation, look into resource policy)
	policyID, err := tools.GetIntValue("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	policy, err := client.GetPolicy(ctx, int64(policyID))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%v update: %s", ErrPolicyActivation, err.Error()))
	}

	createVersionResponse, err := client.CreatePolicyVersion(ctx, cloudlets.CreatePolicyVersionRequest{
		PolicyID:            int64(policyID),
		CreatePolicyVersion: cloudlets.CreatePolicyVersion{},
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("%v update: %s", ErrPolicyActivation, err.Error()))
	}

	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	activationNetwork := cloudlets.VersionActivationNetwork(network)

	// 3. look for activation with this version which is active
	if len(policy.Activations) > 0 {
		propertyName := policy.Activations[0].PropertyInfo.Name
		activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
			PolicyID:     int64(policyID),
			Network:      activationNetwork,
			PropertyName: propertyName,
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("%v update: %s", ErrPolicyActivation, err.Error()))
		}

		for _, act := range activations {
			if act.PolicyInfo.Version == createVersionResponse.Version {
				if act.PolicyInfo.Status == cloudlets.StatusActive {
					// in such case, return
					logger.Debugf("This policy (ID=%d, version=%d) is already active.", policyID, createVersionResponse.Version)
					return resourcePolicyActivationRead(ctx, rd, m)
				}
			}
		}
	}
	logger.Debugf("This policy (ID=%d, version=%d) is not active. Proceeding to activation.", policyID, createVersionResponse.Version)

	// otherwise, create the activation for version and network
	associatedProps, err := tools.GetListValue("associated_properties", rd)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	additionalProps := make([]string, 0, len(associatedProps))
	for _, prop := range associatedProps {
		additionalProps = append(additionalProps, prop.(string))
	}

	err = client.ActivatePolicyVersion(ctx, cloudlets.ActivatePolicyVersionRequest{
		PolicyID: int64(policyID),
		Async:    true,
		Version:  createVersionResponse.Version,
		RequestBody: cloudlets.ActivatePolicyVersionRequestBody{
			Network:                 activationNetwork,
			AdditionalPropertyNames: additionalProps,
		},
	})
	if err != nil {
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}

	// 4. poll until active
	if len(policy.Activations) == 0 {
		policy, err = client.GetPolicy(ctx, int64(policyID))
		if err != nil {
			return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
		}
	}
	propertyName := policy.Activations[0].PropertyInfo.Name
	_, err = getActivePolicyActivation(ctx, client, propertyName, int64(policyID), createVersionResponse.Version, activationNetwork)
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
	versionActivationNetwork := cloudlets.VersionActivationNetwork(tools.StateNetwork(network))
	associatedProps, err := tools.GetListValue("associated_properties", rd)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	additionalProps := make([]string, 0, len(associatedProps))
	for _, prop := range associatedProps {
		additionalProps = append(additionalProps, prop.(string))
	}

	version, err := tools.GetIntValue("version", rd)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version64 := int64(version)
	// if no version provided, we shall fetch latest instead
	if version64 == 0 {
		logger.Debugf("finding latest activation for the policy with ID==%d", int64(policyID))
		v, err := getLatestPolicyActivationVersion(ctx, client, int64(policyID))
		if err != nil {
			return diag.FromErr(err)
		}
		version64 = *v
	}

	logger.Debugf("checking if policy version %d is active", version64)
	policyVersion, err := client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		Version:   version64,
		PolicyID:  int64(policyID),
		OmitRules: true,
	})
	if err != nil {
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}
	var propertyName string
	for _, act := range policyVersion.Activations {
		if cloudlets.VersionActivationNetwork(act.Network) == versionActivationNetwork {
			if act.PolicyInfo.Status == cloudlets.StatusActive {
				// if the given version is active, just refresh status and quit
				logger.Debugf("policy version %d is already active in %s, fetching all details from server", version64, string(versionActivationNetwork))
				return resourcePolicyActivationRead(ctx, rd, m)
			}
			propertyName = act.PropertyInfo.Name
			break
		}
	}

	// at this point, we are sure that the given version is not active
	logger.Debugf("activating policy version %d", version64)
	err = client.ActivatePolicyVersion(ctx, cloudlets.ActivatePolicyVersionRequest{
		PolicyID: int64(policyID),
		Version:  version64,
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
	activation, err := getActivePolicyActivation(ctx, client, propertyName, int64(policyID), version64, versionActivationNetwork)
	if err != nil {
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}

	if err := rd.Set("status", activation.PolicyInfo.Status); err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(formatPolicyActivationID(activation))

	return nil
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
	net := cloudlets.VersionActivationNetwork(network)
	version, err := tools.GetIntValue("version", rd)
	var version64 int64
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		v, err := getLatestPolicyActivationVersion(ctx, client, int64(policyID))
		if err != nil {
			return diag.FromErr(err)
		}
		version64 = *v
	} else {
		version64 = int64(version)
	}

	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID: int64(policyID),
		Network:  net,
	})
	if err != nil {
		return diag.Errorf("%v read: %s", ErrPolicyActivation, err.Error())
	}

	for _, act := range activations {
		if act.PolicyInfo.Version == version64 {
			if err := rd.Set("status", act.PolicyInfo.Status); err != nil {
				return diag.FromErr(err)
			}
			rd.SetId(formatPolicyActivationID(&act))

			return nil
		}
	}

	return diag.FromErr(fmt.Errorf("%v: cannot find the given policy activation version (%d)", ErrPolicyActivation, version64))
}

func formatPolicyActivationID(activation *cloudlets.PolicyActivation) string {
	return fmt.Sprintf("%d:%d", activation.PolicyInfo.PolicyID, activation.PolicyInfo.ActivationDate)
}

// getPolicyActivation gets a policy activation from server
func getPolicyActivation(ctx context.Context, client cloudlets.Cloudlets, propertyName string, policyID, version int64, network cloudlets.VersionActivationNetwork) (*cloudlets.PolicyActivation, error) {
	activations, err := client.ListPolicyActivations(ctx, cloudlets.ListPolicyActivationsRequest{
		PolicyID:     policyID,
		Network:      network,
		PropertyName: propertyName,
	})
	if err != nil {
		return nil, err
	}

	for _, act := range activations {
		if act.PolicyInfo.Version == version {
			return &act, nil
		}
	}

	return nil, fmt.Errorf("%v: policy activation version not found", ErrPolicyActivation)
}

// getActivePolicyActivation polls server until the activation has active status or until context is closed (because of timeout, cancellation or context termination)
func getActivePolicyActivation(ctx context.Context, client cloudlets.Cloudlets, propertyName string, policyID, version int64, network cloudlets.VersionActivationNetwork) (*cloudlets.PolicyActivation, error) {
	activation, err := getPolicyActivation(ctx, client, propertyName, policyID, version, network)
	if err != nil {
		return nil, err
	}
	for activation != nil && activation.PolicyInfo.Status != cloudlets.StatusActive {
		if activation.PolicyInfo.Status == cloudlets.StatusFailed {
			return nil, fmt.Errorf("%v: policyID %d: %s", ErrPolicyActivation, activation.PolicyInfo.PolicyID, activation.PolicyInfo.StatusDetail)
		}
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activation, err = getPolicyActivation(ctx, client, propertyName, policyID, version, network)
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

	return activation, nil
}

func getLatestPolicyActivationVersion(ctx context.Context, client cloudlets.Cloudlets, policyID int64) (*int64, error) {
	versions, err := client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
		PolicyID: int64(policyID),
	})
	if err != nil {
		return nil, fmt.Errorf("%v: %s", ErrPolicyActivation, err.Error())
	}
	var version int64
	for _, v := range versions {
		if version < v.Version {
			version = v.Version
		}
	}
	return tools.Int64Ptr(version), nil
}
