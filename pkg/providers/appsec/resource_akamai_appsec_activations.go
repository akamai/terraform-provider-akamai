package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
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
				DiffSuppressFunc: suppressNoteFieldForAppSecActivation,
			},
			"notification_emails": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of email addresses to be notified with the results of the activation",
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

	activationResp, err := createActivation(ctx, client, createActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(activationResp.ActivationID))

	if err := d.Set("status", activationResp.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	getActivationRequest := appsec.GetActivationsRequest{
		ActivationID: activationResp.ActivationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = pollActivation(ctx, client, activation.Status, getActivationRequest); err != nil {
		return diag.FromErr(err)
	}
	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRead")
	logger.Debug("in resourceActivationsRead")

	activationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getActivations := appsec.GetActivationsRequest{
		ActivationID: activationID,
	}

	activations, err := client.GetActivations(ctx, getActivations)
	if err != nil {
		logger.Errorf("calling 'getActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	network, err := tf.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	if activations.Action == string(appsec.ActivationTypeActivate) &&
		activations.Status == appsec.StatusDeactivated &&
		(string(activations.Network) == network) {
		d.SetId("")
		return nil
	}

	if err := d.Set("status", activations.Status); err != nil {
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

	activationResp, err := createActivation(ctx, client, createActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(activationResp.ActivationID))

	if err := d.Set("status", activationResp.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	getActivationRequest := appsec.GetActivationsRequest{
		ActivationID: activationResp.ActivationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = pollActivation(ctx, client, activation.Status, getActivationRequest); err != nil {
		return diag.FromErr(err)
	}

	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRemove")
	logger.Debug("in resourceActivationsDelete")

	activationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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

	removeActivationRequest := appsec.RemoveActivationsRequest{
		ActivationID:       activationID,
		Action:             string(appsec.ActivationTypeDeactivate),
		Network:            network,
		Note:               note,
		NotificationEmails: notificationEmails,
	}
	removeActivationRequest.ActivationConfigs = append(removeActivationRequest.ActivationConfigs, appsec.ActivationConfigs{
		ConfigID:      configID,
		ConfigVersion: version,
	})

	postresp, err := client.RemoveActivations(ctx, removeActivationRequest)
	if err != nil {
		logger.Errorf("calling 'removeActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	getActivationRequest := appsec.GetActivationsRequest{
		ActivationID: postresp.ActivationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	for activation.Status != appsec.StatusDeactivated && activation.Status != appsec.StatusAborted && activation.Status != appsec.StatusFailed {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationRequest)

			if err != nil {
				return diag.FromErr(err)
			}
			activation = act

		case <-ctx.Done():
			return diag.Errorf("activation context terminated: %s", ctx.Err())
		}
	}

	if err := d.Set("status", activation.Status); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsImport")
	logger.Debug("in appsec_activation resource's resourceImporter")

	iDParts, err := splitID(d.Id(), 3, "configID:version:network")
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

func suppressNoteFieldForAppSecActivation(_, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue != newValue && d.HasChanges("config_id", "version", "network") {
		return false
	}
	return true
}

func createActivation(ctx context.Context, client appsec.APPSEC, request appsec.CreateActivationsRequest) (*appsec.CreateActivationsResponse, error) {
	log := hclog.FromContext(ctx)

	errMsg := "create failed"
	switch request.Action {
	case string(appsec.ActivationTypeActivate):
		errMsg = "create activation failed"
	case string(appsec.ActivationTypeDeactivate):
		errMsg = "create deactivation failed"
	}

	createActivationRetry := CreateActivationRetry

	for {
		log.Debug("creating activation")
		create, err := client.CreateActivations(ctx, request, true)

		if err == nil {
			return create, nil
		}
		log.Debug("%s: retrying: %w", errMsg, err)

		if !isCreateActivationErrorRetryable(err) {
			return nil, fmt.Errorf("%s: %s", errMsg, err)
		}

		select {
		case <-time.After(createActivationRetry):
			createActivationRetry = capDuration(createActivationRetry*2, 5*time.Minute)
			continue

		case <-ctx.Done():
			return nil, fmt.Errorf("activation context terminated: %w", ctx.Err())
		}
	}

}

func pollActivation(ctx context.Context, client appsec.APPSEC, activationStatus appsec.StatusValue, getActivationRequest appsec.GetActivationsRequest) error {
	retriesMax := 5
	retries5xx := 0

	for activationStatus != appsec.StatusActive && activationStatus != appsec.StatusAborted && activationStatus != appsec.StatusFailed {
		select {
		case <-time.After(tf.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationRequest)
			if err != nil {
				var target = &appsec.Error{}
				if !errors.As(err, &target) {
					return fmt.Errorf("error has unexpected type: %T", err)
				}
				if isCreateActivationErrorRetryable(target) {
					retries5xx = retries5xx + 1
					if retries5xx > retriesMax {
						return fmt.Errorf("reached max number of 5xx retries: %d", retries5xx)
					}
					continue
				}
				return err
			}
			retries5xx = 0
			activationStatus = act.Status

		case <-ctx.Done():
			return fmt.Errorf("activation context terminated: %s", ctx.Err())
		}
	}
	return nil
}

func isCreateActivationErrorRetryable(err error) bool {
	var responseErr = &appsec.Error{}
	if !errors.As(err, &responseErr) {
		return false
	}
	if responseErr.StatusCode < 500 &&
		responseErr.StatusCode != 422 &&
		responseErr.StatusCode != 409 {
		return false
	}
	return true
}

func capDuration(t time.Duration, tMax time.Duration) time.Duration {
	if t > tMax {
		return tMax
	}
	return t
}
