package cloudlets

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets/v3"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudletsPolicyActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyActivationCreate,
		ReadContext:   resourcePolicyActivationRead,
		UpdateContext: resourcePolicyActivationUpdate,
		DeleteContext: resourcePolicyActivationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyActivationImport,
		},
		Schema: resourceCloudletsPolicyActivationSchema(),
		Timeouts: &schema.ResourceTimeout{
			Default: &PolicyActivationResourceTimeout,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourceCloudletsPolicyActivationV0().CoreConfigSchema().ImpliedType(),
			Upgrade: timeouts.MigrateToExplicit(),
		}},
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
			ValidateDiagFunc: tf.ValidateNetwork,
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
			MinItems:    1,
			Description: "Set of property IDs to link to this Cloudlets policy. It is required for non-shared policies",
		},
		"is_shared": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates if policy that is being activated is a shared policy",
		},
		"timeouts": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Enables to set timeout for processing",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default": {
						Type:             schema.TypeString,
						Optional:         true,
						ValidateDiagFunc: timeouts.ValidateDurationFormat,
					},
				},
			},
		},
	}
}

var (
	// ActivationPollMinimum is the minimum polling interval for activation creation
	ActivationPollMinimum = time.Minute

	// ActivationPollInterval is the interval for polling an activation status on creation
	ActivationPollInterval = ActivationPollMinimum

	// MaxListActivationsPollRetries is the maximum number of retries for calling ListActivations request in case of returning empty list
	MaxListActivationsPollRetries = 5

	// PolicyActivationResourceTimeout is the default timeout for the resource operations
	PolicyActivationResourceTimeout = time.Minute * 90

	// PolicyActivationRetryPollMinimum is the minimum polling interval for retrying policy activation
	PolicyActivationRetryPollMinimum = time.Second * 15

	// PolicyActivationRetryTimeout is the default timeout for the policy activation retries
	PolicyActivationRetryTimeout = time.Minute * 10

	// ErrNetworkName is used when the user inputs an invalid network name
	ErrNetworkName = errors.New("invalid network name")

	policyActivationRetryRegexp = regexp.MustCompile(`requested propertyname \\"[A-Za-z0-9.\-_]+\\" does not exist`)
)

func resourcePolicyActivationDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationDelete")
	logger.Debug("Deleting cloudlets policy activation")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))

	strategy := getActivationStrategy(rd, meta, logger)

	policyID, err := tf.GetIntValueAsInt64("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tf.GetIntValueAsInt64("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = strategy.deactivatePolicy(ctx, policyID, version, network); err != nil {
		return diag.FromErr(err)
	}

	logger.Debugf("All properties have been removed from policy ID %d", policyID)
	rd.SetId("")
	return nil
}

func resourcePolicyActivationUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationUpdate")

	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	strategy := getActivationStrategy(rd, meta, logger)

	if !rd.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	// 2. In such case, create a new version to activate (for creation, look into resource policy)
	policyID, err := tf.GetIntValueAsInt64("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tf.GetIntValueAsInt64("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = strategy.setupCloudletSpecificData(rd, network); err != nil {
		return diag.FromErr(err)
	}

	isAlreadyActive, id, err := strategy.isReactivationNotNeeded(ctx, policyID, version, rd.HasChange("version"))
	if err != nil {
		if restoreDiags := diag.FromErr(tf.RestoreOldValues(rd, []string{"version", "associated_properties"})); len(restoreDiags) > 0 {
			return append(restoreDiags, diag.FromErr(err)...)
		}
		return diag.FromErr(err)
	}

	if isAlreadyActive {
		// all is active for the given version, policyID and network, proceed to read stage
		logger.Debugf("This policy (ID=%d, version=%d) is already active.", policyID, version)
		rd.SetId(id)
		return resourcePolicyActivationRead(ctx, rd, m)
	}

	// something has changed, we need to reactivate it
	if err = strategy.reactivateVersion(ctx, policyID, version); err != nil {
		if restoreDiags := diag.FromErr(tf.RestoreOldValues(rd, []string{"version", "associated_properties"})); restoreDiags != nil {
			return append(restoreDiags, diag.FromErr(err)...)
		}
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}

	// poll until active
	id, err = strategy.waitForActivation(ctx, policyID, version)
	if err != nil {
		return diag.Errorf("%v update: %s", ErrPolicyActivation, err.Error())
	}
	rd.SetId(id)

	return resourcePolicyActivationRead(ctx, rd, m)
}

func resourcePolicyActivationCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))

	logger.Debug("Creating policy activation")

	policyID, err := tf.GetIntValueAsInt64("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tf.GetIntValueAsInt64("version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	strategy, isShared, err := discoverActivationStrategy(ctx, policyID, meta, logger)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = strategy.setupCloudletSpecificData(rd, network); err != nil {
		return diag.FromErr(err)
	}

	isActive, id, err := strategy.isVersionAlreadyActive(ctx, policyID, version)
	if err != nil {
		return diag.FromErr(err)
	}

	if isActive {
		// if the given version is active, just refresh status and quit
		rd.SetId(id)
		if err = rd.Set("is_shared", isShared); err != nil {
			return diag.Errorf("was not able to set `is_shared` computed field: %s", err)
		}
		return resourcePolicyActivationRead(ctx, rd, m)
	}

	// at this point, we are sure that the given version is not active
	pollingActivationTries := PolicyActivationRetryPollMinimum

	for {
		err = strategy.activateVersion(ctx, policyID, version)
		if err == nil {
			break
		}

		select {
		case <-time.After(pollingActivationTries):
			logger.Debugf("retrying policy activation after %s", pollingActivationTries)
			if pollingActivationTries > PolicyActivationRetryTimeout || !strategy.shouldRetryActivation(err) {
				return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
			}

			pollingActivationTries = 2 * pollingActivationTries
			continue
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return diag.Errorf("timeout waiting for retrying policy activation: last error: %s", err)
			}
			if errors.Is(ctx.Err(), context.Canceled) {
				return diag.Errorf("operation canceled while waiting for retrying policy activation, last error: %s", err)
			}

		}
	}

	// wait until policy activation is done
	id, err = strategy.waitForActivation(ctx, policyID, version)
	if err != nil {
		return diag.Errorf("%v create: %s", ErrPolicyActivation, err.Error())
	}

	rd.SetId(id)

	if err = rd.Set("is_shared", isShared); err != nil {
		return diag.Errorf("was not able to set `is_shared` computed field: %s", err)
	}

	return resourcePolicyActivationRead(ctx, rd, m)
}

func resourcePolicyActivationRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	strategy := getActivationStrategy(rd, meta, logger)

	logger.Debug("Reading policy activations")

	policyID, err := tf.GetIntValueAsInt64("policy_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	attrs, err := strategy.readActivationFromServer(ctx, policyID, network)
	if err != nil {
		return diag.Errorf("policy activation read: %s", err.Error())
	}

	if attrs == nil {
		rd.SetId("")
		return nil
	}

	if err = tf.SetAttrs(rd, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getActivationStrategy(rd *schema.ResourceData, m meta.Meta, logger log.Interface) activationStrategy {
	if rd.Get("is_shared").(bool) {
		return &v3ActivationStrategy{client: ClientV3(m), logger: logger}
	}
	return &v2ActivationStrategy{client: Client(m), logger: logger}
}

func resourcePolicyActivationImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyActivationImport")
	logger.Debugf("Import Policy Activation")

	resID := d.Id()
	parts := strings.Split(resID, ":")

	if len(parts) != 2 {
		return nil, fmt.Errorf("import id should be of format: <policy_id>:<network>, for example: 1234:staging")
	}

	policyID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	network := parts[1]

	strategy, _, err := discoverActivationStrategy(ctx, policyID, meta, logger)
	if err != nil {
		return nil, err
	}

	attrs, id, err := strategy.fetchValuesForImport(ctx, policyID, network)
	if err != nil {
		return nil, err
	}

	err = tf.SetAttrs(d, attrs)
	if err != nil {
		return nil, err
	}

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
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
	activations, err := waitForListPolicyActivations(ctx, client, cloudlets.ListPolicyActivationsRequest{
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
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activations, err = waitForListPolicyActivations(ctx, client, cloudlets.ListPolicyActivationsRequest{
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

	net = tf.StateNetwork(net)

	switch net {
	case "production":
		return cloudlets.PolicyActivationNetworkProduction, nil
	case "staging":
		return cloudlets.PolicyActivationNetworkStaging, nil
	}

	return "", ErrNetworkName
}

func statePolicyActivationNetwork(i interface{}) string {

	net := tf.StateNetwork(i)

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
	activations, err := waitForListPolicyActivations(ctx, client, cloudlets.ListPolicyActivationsRequest{PolicyID: policyID})
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
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activations, err = waitForListPolicyActivations(ctx, client, cloudlets.ListPolicyActivationsRequest{
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

// waitForListPolicyActivations polls server until the ListPolicyActivations returns non-empty list
func waitForListPolicyActivations(ctx context.Context, client cloudlets.Cloudlets, listPolicyActivationsRequest cloudlets.ListPolicyActivationsRequest) ([]cloudlets.PolicyActivation, error) {
	listActivationsPollRetries := MaxListActivationsPollRetries
	activations, err := client.ListPolicyActivations(ctx, listPolicyActivationsRequest)
	if err != nil {
		return nil, err
	}

	for len(activations) == 0 && listActivationsPollRetries > 0 {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			activations, err = client.ListPolicyActivations(ctx, listPolicyActivationsRequest)
			if err != nil {
				return nil, err
			}
			listActivationsPollRetries--

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

func discoverActivationStrategy(ctx context.Context, policyID int64, meta meta.Meta, logger log.Interface) (activationStrategy, bool, error) {
	v2Client := Client(meta)
	_, v2Err := v2Client.GetPolicy(ctx, cloudlets.GetPolicyRequest{PolicyID: policyID})
	if v2Err == nil {
		return &v2ActivationStrategy{client: v2Client, logger: logger}, false, nil
	}

	v3Client := ClientV3(meta)
	_, V3err := v3Client.GetPolicy(ctx, v3.GetPolicyRequest{PolicyID: policyID})
	if V3err == nil {
		return &v3ActivationStrategy{client: v3Client, logger: logger}, true, nil
	}

	return nil, false, fmt.Errorf("could not get policy %d: neither as V2 (%s) nor as V3 (%s)", policyID, v2Err, V3err)

}

type activationStrategy interface {
	isVersionAlreadyActive(ctx context.Context, policyID, version int64) (bool, string, error)
	setupCloudletSpecificData(rd *schema.ResourceData, network string) error
	activateVersion(ctx context.Context, policyID, version int64) error
	shouldRetryActivation(err error) bool
	reactivateVersion(ctx context.Context, policyID, version int64) error
	waitForActivation(ctx context.Context, policyID, version int64) (string, error)
	readActivationFromServer(ctx context.Context, policyID int64, network string) (map[string]any, error)
	isReactivationNotNeeded(ctx context.Context, policyID, version int64, hasVersionChange bool) (bool, string, error)
	deactivatePolicy(ctx context.Context, policyID, version int64, network string) error
	getPolicyActivation(ctx context.Context, policyID int64, network string) (*policyActivationDataSourceModel, error)
	fetchValuesForImport(ctx context.Context, policyID int64, network string) (map[string]any, string, error)
}
