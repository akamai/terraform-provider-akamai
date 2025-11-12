package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceActivations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationsCreate,
		ReadContext:   resourceActivationsRead,
		UpdateContext: resourceActivationsUpdate,
		DeleteContext: resourceActivationsDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: resourceImporter,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration to be activated",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Version of the security configuration to be activated",
			},
			"network": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STAGING",
				Description: "Network on which to activate the configuration version (STAGING or PRODUCTION)",
			},
			"note": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Note describing the activation. Will use timestamp if omitted.",
				DiffSuppressFunc: suppressFieldsForAppSecActivation,
			},
			"notification_emails": {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Description:      "List of email addresses to be notified with the results of the activation",
				DiffSuppressFunc: suppressActivationEmailFieldForAppSecActivation,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The results of the activation",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &AppsecResourceTimeout,
		},
	}
}

const (
	// ActivationPollMinimum is the minimum polling interval for activation creation
	ActivationPollMinimum = time.Minute
)

// activationParams contains the parameters for activation and deactivation operations, used by helper methods
type activationParams struct {
	ConfigID           int
	Version            int
	Network            string
	Note               string
	NotificationEmails []string
	ResourceData       *schema.ResourceData
	Logger             akalog.Interface
}

var (
	// ActivationPollInterval is the interval for polling an activation status on creation
	ActivationPollInterval = ActivationPollMinimum

	// AppsecResourceTimeout is the default timeout for the resource operations
	AppsecResourceTimeout = time.Minute * 90

	// CreateActivationRetry poll wait time code waits between retries for activation creation
	CreateActivationRetry = 10 * time.Second
)

func resourceActivationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsCreate")
	logger.Debug("in resourceActivationsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tf.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	var note string
	note, err = tf.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if note == "" {
		note, err = defaultActivationNote(false)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	notificationEmailsSet, err := tf.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tf.SetToStringSlice(notificationEmailsSet)

	// Create activation params
	params := activationParams{
		ConfigID:           configID,
		Version:            version,
		Network:            network,
		Note:               note,
		NotificationEmails: notificationEmails,
		ResourceData:       d,
		Logger:             logger,
	}

	return activateVersion(ctx, client, params)
}

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRead")
	logger.Debug("in resourceActivationsRead")

	// Get config ID and network from resource state
	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the currently active version for this config and network
	currentActiveVersion, err := findCurrentActiveVersion(ctx, client, configID, network)
	if err != nil {
		logger.Errorf("calling 'findCurrentActiveVersion': %s", err.Error())
		return diag.FromErr(err)
	}

	// If no active version found, mark as destroyed
	if currentActiveVersion == nil {
		d.SetId("")
		return nil
	}

	// Set the state values from the current active activation
	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("network", network); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("version", currentActiveVersion.Version); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("note", currentActiveVersion.Notes); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("notification_emails", currentActiveVersion.NotificationEmails); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("status", string(currentActiveVersion.Status)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceActivationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsUpdate")
	logger.Debug("in resourceActivationsUpdate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tf.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	var note string
	note, err = tf.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if note == "" {
		note, err = defaultActivationNote(false)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	notificationEmailsSet, err := tf.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tf.SetToStringSlice(notificationEmailsSet)

	// Create activation params
	params := activationParams{
		ConfigID:           configID,
		Version:            version,
		Network:            network,
		Note:               note,
		NotificationEmails: notificationEmails,
		ResourceData:       d,
		Logger:             logger,
	}

	return activateVersion(ctx, client, params)
}

func resourceActivationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRemove")
	logger.Debug("in resourceActivationsDelete")

	// Get the config values from state
	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := tf.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}
	var note string
	note, err = tf.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if note == "" {
		note, err = defaultActivationNote(true)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	notificationEmailsSet, err := tf.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tf.SetToStringSlice(notificationEmailsSet)

	// Create deactivation params
	params := activationParams{
		ConfigID:           configID,
		Version:            version,
		Network:            network,
		Note:               note,
		NotificationEmails: notificationEmails,
		ResourceData:       d,
		Logger:             logger,
	}

	return deactivateVersion(ctx, client, params)
}

func resourceImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsImport")
	logger.Debug("in appsec_activation resource's resourceImporter")

	iDParts, err := id.Split(d.Id(), 3, "configID:version:network")
	if err != nil {
		return nil, err
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return nil, err
	}
	version, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return nil, err
	}
	network := iDParts[2]
	if !(network == "STAGING" || network == "PRODUCTION") {
		return nil, fmt.Errorf("bad network value %s; must be either STAGING or PRODUCTION", network)

	}

	request := appsec.GetActivationHistoryRequest{
		ConfigID: configID,
	}
	response, err := client.GetActivationHistory(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, activation := range response.ActivationHistory {
		if activation.Version == version && activation.Network == network {
			d.SetId(strconv.Itoa(activation.ActivationID))
			if err = d.Set("config_id", configID); err != nil {
				return nil, err
			}
			if err = d.Set("network", network); err != nil {
				return nil, err
			}
			if err = d.Set("note", activation.Notes); err != nil {
				return nil, err
			}
			if err = d.Set("version", version); err != nil {
				return nil, err
			}
			if err = d.Set("notification_emails", activation.NotificationEmails); err != nil {
				return nil, err
			}

			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("no activation found for configId %d, version %d, network %s", configID, version, network)

}

func lookupActivation(ctx context.Context, client appsec.APPSEC, query appsec.GetActivationsRequest) (*appsec.GetActivationsResponse, error) {
	activations, err := client.GetActivations(ctx, query)
	if err != nil {
		return nil, err
	}

	return activations, nil
}

func defaultActivationNote(deactivating bool) (string, error) {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return "", err
	}

	formattedTime := time.Now().In(location).Format(time.RFC850)
	if deactivating {
		return fmt.Sprintf("Deactivation request %s", formattedTime), nil
	}
	return fmt.Sprintf("Activation request %s", formattedTime), nil
}

// retryActivationRequest executes an activation API call with exponential backoff retry logic
func retryActivationRequest(
	ctx context.Context,
	errMsg string,
	apiCall func() error,
) error {
	log := hclog.FromContext(ctx)
	retryDelay := CreateActivationRetry

	for {
		log.Debug("attempting activation request")
		err := apiCall()

		if err == nil {
			return nil
		}
		log.Debug("%s: retrying: %w", errMsg, err)

		if !isCreateActivationErrorRetryable(err) {
			return fmt.Errorf("%s: %s", errMsg, err)
		}

		select {
		case <-time.After(retryDelay):
			retryDelay = date.CapDuration(retryDelay*2, 5*time.Minute)
			continue

		case <-ctx.Done():
			return fmt.Errorf("activation context terminated: %w", ctx.Err())
		}
	}
}

func createActivation(ctx context.Context, client appsec.APPSEC, request appsec.CreateActivationsRequest) (*appsec.CreateActivationsResponse, error) {
	errMsg := "create failed"
	switch request.Action {
	case string(appsec.ActivationTypeActivate):
		errMsg = "create activation failed"
	case string(appsec.ActivationTypeDeactivate):
		errMsg = "create deactivation failed"
	}

	var result *appsec.CreateActivationsResponse
	err := retryActivationRequest(ctx, errMsg, func() error {
		create, err := client.CreateActivations(ctx, request, true)
		if err == nil {
			result = create
		}
		return err
	})

	return result, err
}

func pollActivation(ctx context.Context, client appsec.APPSEC, activationStatus appsec.StatusValue, getActivationsRequest appsec.GetActivationsRequest) (appsec.StatusValue, error) {
	retriesMax := 5
	retries5xx := 0

	for activationStatus != appsec.StatusActive && activationStatus != appsec.StatusAborted && activationStatus != appsec.StatusFailed {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationsRequest)
			if err != nil {
				var target = &appsec.Error{}
				if !errors.As(err, &target) {
					return "", fmt.Errorf("error has unexpected type: %T", err)
				}
				if isCreateActivationErrorRetryable(target) {
					retries5xx = retries5xx + 1
					if retries5xx > retriesMax {
						return "", fmt.Errorf("reached max number of 5xx retries: %d", retries5xx)
					}
					continue
				}
				return "", err
			}
			retries5xx = 0
			activationStatus = act.Status

		case <-ctx.Done():
			return "", fmt.Errorf("activation context terminated: %s", ctx.Err())
		}
	}
	return activationStatus, nil
}

func isCreateActivationErrorRetryable(err error) bool {
	var responseErr = &appsec.Error{}
	if !errors.As(err, &responseErr) {
		return false
	}
	if responseErr.StatusCode < 500 &&
		responseErr.StatusCode != 422 {
		return false
	}
	return true
}

// generateHostMoveWarning creates a warning message for host move operations
func generateHostMoveWarning(hostsToMove []appsec.HostToMove, currentConfigID int) diag.Diagnostic {
	if len(hostsToMove) == 0 {
		return diag.Diagnostic{}
	}

	var hostnames []string
	var sourceConfigID int

	for _, host := range hostsToMove {
		hostnames = append(hostnames, host.Host)
		sourceConfigID = host.FromConfig.ConfigID // All hosts are from same config due to validation
	}

	var hostnamesStr string
	if len(hostnames) == 1 {
		hostnamesStr = fmt.Sprintf("Hostname %s was", hostnames[0])
	} else {
		hostnamesStr = fmt.Sprintf("Hostnames %v were", hostnames)
	}

	message := fmt.Sprintf("%s moved from config %d to config %d. Refresh the source config %d resource(s) to update state.",
		hostnamesStr, sourceConfigID, currentConfigID, sourceConfigID)

	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Host Move Detected",
		Detail:   message,
	}
}

// createActivationWithValidation validates and creates activation with host move if needed
// Returns the activation response, host move validation result, and any diagnostics
func createActivationWithValidation(
	ctx context.Context,
	client appsec.APPSEC,
	configID int,
	version int,
	network string,
	note string,
	notificationEmails []string,
) (*appsec.CreateActivationsResponse, *appsec.GetHostMoveValidationResponse, diag.Diagnostics) {
	// Check for host move validation before creating activation
	hostMoveValidationRequest := appsec.GetHostMoveValidationRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		Network:       appsec.NetworkValue(network),
	}

	hostMoveValidation, err := client.GetHostMoveValidation(ctx, hostMoveValidationRequest)
	if err != nil {
		return nil, nil, diag.Errorf("calling 'GetHostMoveValidation': %s", err.Error())
	}

	// Validate that all hosts are from the same source config
	if hostMoveValidation != nil && len(hostMoveValidation.HostsToMove) > 0 {
		if err := validateSingleSourceConfig(hostMoveValidation.HostsToMove); err != nil {
			return nil, hostMoveValidation, diag.FromErr(err)
		}
	}

	// Choose activation path based on whether host move is needed
	var activationResp *appsec.CreateActivationsResponse

	if hostMoveValidation != nil && len(hostMoveValidation.HostsToMove) > 0 {
		// Use host move activation path
		activationResp, err = activateWithHostMove(ctx, client, configID, version, network, note, notificationEmails, hostMoveValidation.HostsToMove)
		if err != nil {
			return nil, hostMoveValidation, diag.FromErr(err)
		}
	} else {
		// Use regular activation path
		activationResp, err = activate(ctx, client, configID, version, network, note, notificationEmails)
		if err != nil {
			return nil, hostMoveValidation, diag.FromErr(err)
		}
	}

	return activationResp, hostMoveValidation, nil
}

// activateWithHostMove creates an activation using the host move API
func activateWithHostMove(
	ctx context.Context,
	client appsec.APPSEC,
	configID int,
	version int,
	network string,
	note string,
	notificationEmails []string,
	hostsToMove []appsec.HostToMove,
) (*appsec.CreateActivationsResponse, error) {
	createActivationWithHostMoveRequest := appsec.CreateActivationsWithHostMoveRequest{
		ConfigID:           configID,
		ConfigVersion:      version,
		Action:             string(appsec.ActivationTypeActivate),
		Network:            appsec.NetworkValue(network),
		Note:               note,
		NotificationEmails: notificationEmails,
		HostsToMove:        hostsToMove,
	}

	errMsg := "create activation with host move failed"
	var result *appsec.CreateActivationsResponse

	err := retryActivationRequest(ctx, errMsg, func() error {
		hostMoveActivationResp, err := client.CreateActivationsWithHostMove(ctx, createActivationWithHostMoveRequest)
		if err == nil {
			// Convert to standard response format
			result = &appsec.CreateActivationsResponse{
				ActivationID: hostMoveActivationResp.ActivationID,
				Status:       appsec.StatusValue(hostMoveActivationResp.Status),
				Network:      appsec.NetworkValue(hostMoveActivationResp.Network),
			}
		}
		return err
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

// activate creates a regular activation (without host move)
func activate(
	ctx context.Context,
	client appsec.APPSEC,
	configID int,
	version int,
	network string,
	note string,
	notificationEmails []string,
) (*appsec.CreateActivationsResponse, error) {
	createActivationRequest := appsec.CreateActivationsRequest{
		Action:             string(appsec.ActivationTypeActivate),
		Network:            network,
		Note:               note,
		NotificationEmails: notificationEmails,
	}
	createActivationRequest.ActivationConfigs = append(createActivationRequest.ActivationConfigs, appsec.ActivationConfigs{
		ConfigID:      configID,
		ConfigVersion: version,
	})

	return createActivation(ctx, client, createActivationRequest)
}

// validateSingleSourceConfig ensures all hosts to move come from the same source config
func validateSingleSourceConfig(hostsToMove []appsec.HostToMove) error {
	if len(hostsToMove) == 0 {
		return nil
	}

	firstConfigID := hostsToMove[0].FromConfig.ConfigID
	for _, host := range hostsToMove {
		if host.FromConfig.ConfigID != firstConfigID {
			return fmt.Errorf("you can't move hostnames from more than one security configuration at a time. Instead, make successive updates, one for each source security configuration")
		}
	}
	return nil
}

// findCurrentActiveVersion finds the currently active version for a config and network
func findCurrentActiveVersion(ctx context.Context, client appsec.APPSEC, configID int, network string) (*appsec.Activation, error) {
	request := appsec.GetActivationHistoryRequest{
		ConfigID: configID,
	}

	response, err := client.GetActivationHistory(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get activation history for config %d: %w", configID, err)
	}

	// Find the most recent activation with active status for the specified network
	// The activation history should be ordered, but we'll look through all to find
	// the most recent one with active status
	var latestActiveActivation *appsec.Activation

	for _, activation := range response.ActivationHistory {
		if activation.Network == network &&
			string(activation.Status) == string(appsec.StatusActive) {
			// Since we want the most recent, we'll take the first active one we find
			// as the history should be in reverse chronological order
			latestActiveActivation = &activation
			break
		}
	}

	return latestActiveActivation, nil
}

// findCurrentActiveOrPendingVersion finds the currently active or pending version for a config and network
func findCurrentActiveOrPendingVersion(ctx context.Context, client appsec.APPSEC, configID int, network string) (*appsec.Activation, error) {
	request := appsec.GetActivationHistoryRequest{
		ConfigID: configID,
	}

	response, err := client.GetActivationHistory(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get activation history for config %d: %w", configID, err)
	}

	// Look for active OR pending activations for the specified network
	for _, activation := range response.ActivationHistory {
		if activation.Network == network {
			status := string(activation.Status)

			// Return active versions immediately (highest priority)
			if status == string(appsec.StatusActive) {
				return &activation, nil
			}

			// Also return pending activations for handling
			if status == string(appsec.StatusPending) ||
				status == string(appsec.StatusInProgress) ||
				status == string(appsec.StatusNew) {
				return &activation, nil
			}
		}
	}

	return nil, nil
}

// findCurrentActiveOrPendingDeactivation finds active versions or pending deactivations for a config and network
func findCurrentActiveOrPendingDeactivation(ctx context.Context, client appsec.APPSEC, configID int, network string, logger akalog.Interface) (*appsec.Activation, error) {
	request := appsec.GetActivationHistoryRequest{
		ConfigID: configID,
	}

	response, err := client.GetActivationHistory(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to get activation history for config %d: %w", configID, err)
	}

	// Look for active versions OR pending deactivations for the specified network
	for _, activation := range response.ActivationHistory {
		if activation.Network == network {
			status := string(activation.Status)

			// Return pending deactivations first (highest priority for delete operation)
			if isPendingDeactivation(status) {
				logger.Debugf("found pending deactivation for version %d on %s (status: %s)", activation.Version, network, status)
				return &activation, nil
			}

			// Return active versions
			if status == string(appsec.StatusActive) {
				logger.Debugf("found active version %d on %s", activation.Version, network)
				return &activation, nil
			}
		}
	}

	return nil, nil
}

// isPendingDeactivation checks if the status represents a pending/in-progress deactivation
func isPendingDeactivation(status string) bool {
	return status == "PENDING_DEACTIVATION" ||
		status == "DEACTIVATION_IN_PROGRESS" ||
		status == "DEACTIVATION_PENDING"
}

// activateVersion orchestrates the activation of a configuration version
func activateVersion(ctx context.Context, client appsec.APPSEC, params activationParams) diag.Diagnostics {
	// Check if there's already an active or pending version for this config and network
	currentVersion, err := findCurrentActiveOrPendingVersion(ctx, client, params.ConfigID, params.Network)
	if err != nil {
		params.Logger.Warnf("unable to check current version: %s", err.Error())
		// Continue with activation since this is not a critical error
	} else if currentVersion != nil {
		diags := handleCurrentVersion(ctx, client, currentVersion, params)
		if diags != nil {
			return diags
		}
		// If handleCurrentVersion returns nil, it means we should proceed with activation
		// (this happens when current version is different and active)
		if currentVersion.Version == params.Version {
			// Same version handling was completed, return
			return nil
		}
	}

	// Proceed with creating a new activation
	return performActivation(ctx, client, params)
}

// handleCurrentVersion determines how to handle an existing activation
func handleCurrentVersion(ctx context.Context, client appsec.APPSEC, currentVersion *appsec.Activation, params activationParams) diag.Diagnostics {
	if currentVersion.Version == params.Version {
		return handleSameVersion(ctx, client, currentVersion, params)
	}
	return handleDifferentVersion(currentVersion, params)
}

// handleSameVersion handles the case where the requested version is already active or pending
func handleSameVersion(ctx context.Context, client appsec.APPSEC, currentVersion *appsec.Activation, params activationParams) diag.Diagnostics {
	status := string(currentVersion.Status)
	params.ResourceData.SetId(strconv.Itoa(currentVersion.ActivationID))

	if status == string(appsec.StatusActive) {
		params.Logger.Infof("version %d is already active on %s for config %d, using existing activation", params.Version, params.Network, params.ConfigID)

		// Set the status field since the version is already active
		if err := params.ResourceData.Set("status", status); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		return nil
	}

	// Same version is already being activated - wait for completion
	params.Logger.Infof("version %d is already being activated on %s for config %d (status: %s), waiting for completion", params.Version, params.Network, params.ConfigID, status)

	getActivationsRequest := appsec.GetActivationsRequest{
		ActivationID: currentVersion.ActivationID,
	}
	finalStatus, err := pollActivation(ctx, client, appsec.StatusValue(currentVersion.Status), getActivationsRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the final status after successful polling
	if err := params.ResourceData.Set("status", string(finalStatus)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

// handleDifferentVersion handles the case where a different version is currently active or pending
func handleDifferentVersion(currentVersion *appsec.Activation, params activationParams) diag.Diagnostics {
	status := string(currentVersion.Status)

	if status == string(appsec.StatusActive) {
		params.Logger.Infof("current active version on %s for config %d is %d, will activate version %d", params.Network, params.ConfigID, currentVersion.Version, params.Version)
		return nil
	}

	// Different version is pending - this is problematic
	return diag.Errorf("cannot activate version %d while version %d is %s on %s for config %d",
		params.Version, currentVersion.Version, status, params.Network, params.ConfigID)
}

// performActivation creates and polls a new activation with host move support
func performActivation(ctx context.Context, client appsec.APPSEC, params activationParams) diag.Diagnostics {
	// Handle host move validation and activation
	activationResp, hostMoveValidation, diags := createActivationWithValidation(ctx, client,
		params.ConfigID, params.Version, params.Network, params.Note, params.NotificationEmails)
	if diags != nil {
		return diags
	}

	params.ResourceData.SetId(strconv.Itoa(activationResp.ActivationID))

	if err := params.ResourceData.Set("status", activationResp.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	getActivationsRequest := appsec.GetActivationsRequest{
		ActivationID: activationResp.ActivationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationsRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	finalStatus, err := pollActivation(ctx, client, activation.Status, getActivationsRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the final status after successful polling
	if err := params.ResourceData.Set("status", string(finalStatus)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	// Collect warnings for host move operations
	var warnings diag.Diagnostics
	if hostMoveValidation != nil && len(hostMoveValidation.HostsToMove) > 0 {
		if warning := generateHostMoveWarning(hostMoveValidation.HostsToMove, params.ConfigID); warning.Summary != "" {
			warnings = append(warnings, warning)
		}
	}

	return warnings
}

// deactivateVersion orchestrates the deactivation of a configuration version
func deactivateVersion(ctx context.Context, client appsec.APPSEC, params activationParams) diag.Diagnostics {
	// Check if there's already a pending deactivation for this version
	currentVersion, err := findCurrentActiveOrPendingDeactivation(ctx, client, params.ConfigID, params.Network, params.Logger)
	if err != nil {
		params.Logger.Errorf("unable to check current version status: %s", err.Error())
		return diag.FromErr(err)
	}

	// If deactivation is already in progress, wait for it
	if currentVersion != nil && isPendingDeactivation(currentVersion.Status) {
		return waitForDeactivation(ctx, client, currentVersion, params)
	}

	// If no active version exists, nothing to deactivate
	if currentVersion == nil {
		params.Logger.Infof("no active version found for config %d on %s, nothing to deactivate", params.ConfigID, params.Network)
		return nil
	}

	// Proceed with creating a new deactivation request
	return performDeactivation(ctx, client, currentVersion.ActivationID, params)
}

// waitForDeactivation waits for an existing pending deactivation to complete
func waitForDeactivation(ctx context.Context, client appsec.APPSEC, currentVersion *appsec.Activation, params activationParams) diag.Diagnostics {
	params.Logger.Infof("deactivation already in progress for version %d on %s (status: %s), waiting for completion",
		currentVersion.Version, params.Network, currentVersion.Status)
	params.ResourceData.SetId(strconv.Itoa(currentVersion.ActivationID))

	getActivationsRequest := appsec.GetActivationsRequest{
		ActivationID: currentVersion.ActivationID,
	}

	// Poll until deactivation completes
	activation, err := lookupActivation(ctx, client, getActivationsRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	for activation.Status != appsec.StatusDeactivated && activation.Status != appsec.StatusAborted && activation.Status != appsec.StatusFailed {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationsRequest)
			if err != nil {
				return diag.FromErr(err)
			}
			activation = act

		case <-ctx.Done():
			return diag.Errorf("deactivation context terminated: %s", ctx.Err())
		}
	}

	if err := params.ResourceData.Set("status", activation.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

// performDeactivation creates and polls a new deactivation request
func performDeactivation(ctx context.Context, client appsec.APPSEC, activationID int, params activationParams) diag.Diagnostics {
	removeActivationRequest := appsec.RemoveActivationsRequest{
		ActivationID:       activationID,
		Action:             string(appsec.ActivationTypeDeactivate),
		Network:            params.Network,
		Note:               params.Note,
		NotificationEmails: params.NotificationEmails,
	}
	removeActivationRequest.ActivationConfigs = append(removeActivationRequest.ActivationConfigs, appsec.ActivationConfigs{
		ConfigID:      params.ConfigID,
		ConfigVersion: params.Version,
	})

	postResp, err := client.RemoveActivations(ctx, removeActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	params.ResourceData.SetId(strconv.Itoa(postResp.ActivationID))

	if err := params.ResourceData.Set("status", postResp.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	getActivationsRequest := appsec.GetActivationsRequest{
		ActivationID: postResp.ActivationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationsRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	for activation.Status != appsec.StatusDeactivated && activation.Status != appsec.StatusAborted && activation.Status != appsec.StatusFailed {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationsRequest)
			if err != nil {
				return diag.FromErr(err)
			}
			activation = act

		case <-ctx.Done():
			return diag.Errorf("activation context terminated: %s", ctx.Err())
		}
	}

	if err := params.ResourceData.Set("status", activation.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}
