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
	"github.com/apex/log"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudletsApplicationLoadBalancerActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationLoadBalancerActivationCreate,
		ReadContext:   resourceApplicationLoadBalancerActivationRead,
		UpdateContext: resourceApplicationLoadBalancerActivationUpdate,
		DeleteContext: resourceApplicationLoadBalancerActivationDelete,
		Schema:        resourceCloudletsApplicationLoadBalancerActivationSchema(),
		Timeouts: &schema.ResourceTimeout{
			Default: &ApplicationLoadBalancerActivationResourceTimeout,
		},
	}
}

func resourceCloudletsApplicationLoadBalancerActivationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"origin_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The conditional origin’s unique identifier",
		},
		"network": {
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validateNetwork,
			StateFunc:        stateALBActivationNetwork,
			ForceNew:         true,
			Description:      "The network you want to activate the application load balancer version on (options are Staging and Production)",
		},
		"version": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Cloudlets application load balancer version you want to activate",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Activation status for this application load balancer",
		},
	}
}

var (
	// ALBActivationPollMinimum is the minimum polling interval for activation creation
	ALBActivationPollMinimum = time.Second * 15
	// ALBActivationPollInterval is the interval for polling an activation status on creation
	ALBActivationPollInterval = ALBActivationPollMinimum

	// ApplicationLoadBalancerActivationResourceTimeout is the default timeout for the resource operations
	ApplicationLoadBalancerActivationResourceTimeout = time.Minute * 20
)

func resourceApplicationLoadBalancerActivationDelete(_ context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourceApplicationLoadBalancerActivationDelete")
	logger.Debug("Deleting cloudlets application load balancer activation from local schema only")
	logger.Info("Cloudlets API does not support application load balancer activation version deletion - resource will only be removed from state")
	rd.SetId("")
	return nil
}

func resourceApplicationLoadBalancerActivationUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourceApplicationLoadBalancerActivationUpdate")

	if !rd.HasChanges("version") {
		logger.Debugf("nothing has changed, nothing to update")
		return resourceApplicationLoadBalancerActivationRead(ctx, rd, m)
	}

	logger.Debugf("version number has changed: proceed to create and activate a new application load balancer activation version")

	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	_, err := resourceApplicationLoadBalancerActivationChange(ctx, rd, logger, client)
	if err != nil {
		return diag.Errorf("%v update: %s", ErrApplicationLoadBalancerActivation, err.Error())
	}
	return resourceApplicationLoadBalancerActivationRead(ctx, rd, m)
}

func resourceApplicationLoadBalancerActivationCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourceApplicationLoadBalancerActivationCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Creating application load balancer activation")

	activation, err := resourceApplicationLoadBalancerActivationChange(ctx, rd, logger, client)
	if err != nil {
		return diag.Errorf("%v create: %s", ErrApplicationLoadBalancerActivation, err.Error())
	}
	rd.SetId(fmt.Sprintf("%s:%s", activation.OriginID, activation.Network))
	return resourceApplicationLoadBalancerActivationRead(ctx, rd, m)
}

func resourceApplicationLoadBalancerActivationChange(ctx context.Context, rd *schema.ResourceData, logger log.Interface, client cloudlets.Cloudlets) (*cloudlets.ActivationResponse, error) {
	originID, err := tools.GetStringValue("origin_id", rd)
	if err != nil {
		return nil, err
	}
	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return nil, err
	}
	activationNetwork, err := getALBActivationNetwork(network)
	if err != nil {
		return nil, err
	}
	v, err := tools.GetIntValue("version", rd)
	if err != nil {
		return nil, err
	}
	version := int64(v)

	logger.Debugf("checking if application load balancer version %d is active", version)
	activations, err := client.GetLoadBalancerActivations(ctx, originID)
	if err != nil {
		return nil, err
	}

	for _, act := range activations {
		if act.Network == activationNetwork && act.Version == version {
			if act.Status == cloudlets.ActivationStatusActive {
				// if the given version is active, just refresh status and quit
				logger.Debugf("application load balancer version %d is already active in %s, fetching all details from server", version, string(activationNetwork))
				return &act, nil
			}
			break
		}
	}

	// at this point, we are sure that the given version is not active
	logger.Debugf("activating application load balancer version %d", version)

	activation, err := client.ActivateLoadBalancerVersion(ctx, cloudlets.ActivateLoadBalancerVersionRequest{
		OriginID: originID,
		Async:    true,
		ActivationRequest: cloudlets.ActivationRequestParams{
			Network: activationNetwork,
			Version: version,
		},
	})
	if err != nil {
		return activation, err
	}

	// wait until application load balancer activation is done
	activation, err = waitForLoadBalancerActivation(ctx, client, originID, version, activationNetwork)
	if err != nil {
		return nil, err
	}

	if err := rd.Set("status", activation.Status); err != nil {
		return nil, err
	}
	if err := rd.Set("version", activation.Version); err != nil {
		return nil, err
	}
	return activation, nil
}

func resourceApplicationLoadBalancerActivationRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourceApplicationLoadBalancerActivationRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading application load balancer activations")

	originID, err := tools.GetStringValue("origin_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	net, err := getALBActivationNetwork(network)
	if err != nil {
		return diag.FromErr(err)
	}
	var version int64
	v, err := tools.GetIntValue("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	version = int64(v)

	activations, err := client.GetLoadBalancerActivations(ctx, originID)
	if err != nil {
		return diag.Errorf("%v read: %s", ErrApplicationLoadBalancerActivation, err.Error())
	}

	for _, act := range activations {
		if act.Version == version && act.Network == net && act.Status == cloudlets.ActivationStatusActive {
			if err := rd.Set("status", act.Status); err != nil {
				return diag.FromErr(err)
			}
			if err := rd.Set("version", act.Version); err != nil {
				return diag.FromErr(err)
			}

			return nil
		}
	}

	return diag.Errorf("%v: cannot find the given application load balancer activation version '%d' for network '%s'", ErrApplicationLoadBalancerActivation, version, net)
}

func getApplicationLoadBalancerActivation(ctx context.Context, client cloudlets.Cloudlets, originID string, version int64, network cloudlets.ActivationNetwork) ([]cloudlets.ActivationResponse, error) {
	activations, err := client.GetLoadBalancerActivations(ctx, originID)
	filteredActivations := make([]cloudlets.ActivationResponse, 0, len(activations))
	if err != nil {
		return nil, err
	}

	for _, act := range activations {
		if act.Version == version && act.Network == network {
			filteredActivations = append(filteredActivations, act)
		}
	}

	if len(filteredActivations) > 0 {
		return filteredActivations, nil
	}
	return nil, fmt.Errorf("%v: application load balancer activation version not found", ErrApplicationLoadBalancerActivation)
}

// waitForLoadBalancerActivation polls server until the activation has active status or until context is closed (because of timeout, cancellation or context termination)
func waitForLoadBalancerActivation(ctx context.Context, client cloudlets.Cloudlets, originID string, version int64, network cloudlets.ActivationNetwork) (*cloudlets.ActivationResponse, error) {
	activations, err := getApplicationLoadBalancerActivation(ctx, client, originID, version, network)
	if err != nil {
		return nil, err
	}
	for !hasStatus(activations, cloudlets.ActivationStatusActive) {
		if !hasStatus(activations, cloudlets.ActivationStatusPending) {
			return nil, fmt.Errorf("%v: originID %s", ErrApplicationLoadBalancerActivation, activations[0].OriginID)
		}
		select {
		case <-time.After(tools.MaxDuration(ALBActivationPollInterval, ALBActivationPollMinimum)):
			activations, err = getApplicationLoadBalancerActivation(ctx, client, originID, version, network)
			if err != nil {
				return nil, err
			}

		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, ErrApplicationLoadBalancerActivationTimeout
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil, ErrApplicationLoadBalancerActivationCanceled
			}
			return nil, fmt.Errorf("%v: %w", ErrApplicationLoadBalancerActivationContextTerminated, ctx.Err())
		}
	}
	for _, activation := range activations {
		if activation.Status == cloudlets.ActivationStatusActive {
			return &activation, nil
		}
	}
	// should not reach here
	return nil, nil
}

func hasStatus(activations []cloudlets.ActivationResponse, status cloudlets.ActivationStatus) bool {
	for _, activation := range activations {
		if activation.Status == status {
			return true
		}
	}
	return false
}

func getALBActivationNetwork(net string) (cloudlets.ActivationNetwork, error) {

	switch net {
	case "PRODUCTION", "prod", "production":
		return cloudlets.ActivationNetworkProduction, nil
	case "STAGING", "staging":
		return cloudlets.ActivationNetworkStaging, nil
	}

	return "", ErrNetworkName
}

func stateALBActivationNetwork(i interface{}) string {

	val, ok := i.(string)
	if !ok {
		panic(fmt.Sprintf("value type is not a string: %T", i))
	}

	net, err := getALBActivationNetwork(val)
	if err != nil {
		return ""
	}
	return string(net)
}

// validateNetwork defines network validation logic
func validateNetwork(i interface{}, _ cty.Path) diag.Diagnostics {
	val, ok := i.(string)
	if !ok {
		return diag.Errorf("'network' value is not a string: %v", i)
	}
	switch val {
	case "PRODUCTION", "STAGING", "prod", "production", "staging":
		return nil
	}
	return diag.Errorf("'%s' is an invalid network value: should be 'PRODUCTION', 'STAGING', 'prod', 'production', 'staging'", val)
}
