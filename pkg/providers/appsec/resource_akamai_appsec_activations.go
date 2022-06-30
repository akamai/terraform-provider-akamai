package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
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
				Description: "The ID of the security configuration to be activated",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The version of the security configuration to be activated",
			},
			"network": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STAGING",
				Description: "The network on which to activate the configuration version",
			},
			"activate": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    true,
				Deprecated: `The setting activate has been deprecated; "terraform apply" will always perform activation. (Use "terraform destroy" for deactivation.)`,
			},
			"note": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "A note describing the activation. Will use timestamp if omitted.",
				ConflictsWith: []string{"notes"},
			},
			"notes": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "A note describing the activation",
				Deprecated:    `The setting notes has been deprecated. Use "note" instead.`,
				ConflictsWith: []string{"note"},
			},
			"notification_emails": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of email addresses to be notified with the results of the activation",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The results of the activation",
			},
		},
	}
}

const (
	// ActivationPollMinimum is the minumum polling interval for activation creation
	ActivationPollMinimum = time.Minute
)

var (
	// ActivationPollInterval is the interval for polling an activation status on creation
	ActivationPollInterval = ActivationPollMinimum
)

func resourceActivationsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsCreate")
	logger.Debug("in resourceActivationsCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	var note string
	note, err = tools.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if note == "" {
		note, err = tools.GetStringValue("notes", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
	}
	if note == "" {
		note, err = defaultActivationNote(false)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	notificationEmailsSet, err := tools.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tools.SetToStringSlice(notificationEmailsSet)

	createActivationRequest := appsec.CreateActivationsRequest{
		Action:             "ACTIVATE",
		Network:            network,
		Note:               note,
		NotificationEmails: notificationEmails,
	}
	createActivationRequest.ActivationConfigs = append(createActivationRequest.ActivationConfigs, appsec.ActivationConfigs{
		ConfigID:      configID,
		ConfigVersion: version,
	})

	postresp, err := client.CreateActivations(ctx, createActivationRequest, true)
	if err != nil {
		logger.Errorf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	getActivationRequest := appsec.GetActivationsRequest{
		ActivationID: postresp.ActivationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	for activation.Status != appsec.StatusActive {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationRequest)
			if err != nil {
				return diag.FromErr(err)
			}
			activation = act

		case <-ctx.Done():
			return diag.Errorf("activation context terminated: %s", ctx.Err())
		}
	}

	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
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

	if err := d.Set("status", activations.Status); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceActivationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsUpdate")
	logger.Debug("in resourceActivationsUpdate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	var note string
	note, err = tools.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if errors.Is(err, tools.ErrNotFound) {
		note, err = tools.GetStringValue("notes", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
	}
	if note == "" {
		note, err = defaultActivationNote(false)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	notificationEmailsSet, err := tools.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tools.SetToStringSlice(notificationEmailsSet)

	createActivationRequest := appsec.CreateActivationsRequest{
		Action:             "ACTIVATE",
		Network:            network,
		Note:               note,
		NotificationEmails: notificationEmails,
	}
	createActivationRequest.ActivationConfigs = append(createActivationRequest.ActivationConfigs, appsec.ActivationConfigs{
		ConfigID:      configID,
		ConfigVersion: version,
	})

	postresp, err := client.CreateActivations(ctx, createActivationRequest, true)
	if err != nil {
		logger.Errorf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	getActivationRequest := appsec.GetActivationsRequest{
		ActivationID: postresp.ActivationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	for activation.Status != appsec.StatusActive {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationRequest)

			if err != nil {
				return diag.FromErr(err)
			}
			activation = act

		case <-ctx.Done():
			return diag.Errorf("activation context terminated: %s", ctx.Err())
		}
	}

	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRemove")
	logger.Debug("in resourceActivationsDelete")

	activationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	var note string
	note, err = tools.GetStringValue("note", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if note == "" {
		note, err = tools.GetStringValue("notes", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
	}
	if note == "" {
		note, err = defaultActivationNote(true)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	notificationEmailsSet, err := tools.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tools.SetToStringSlice(notificationEmailsSet)

	removeActivationRequest := appsec.RemoveActivationsRequest{
		ActivationID:       activationID,
		Action:             "DEACTIVATE",
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
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	getActivationRequest := appsec.GetActivationsRequest{
		ActivationID: activationID,
	}

	activation, err := lookupActivation(ctx, client, getActivationRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	for activation.Status != appsec.StatusDeactivated {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
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
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	return nil
}

func resourceImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsImport")
	logger.Debug("in resourceActivationsCreate")

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
			if err = d.Set("activate", true); err != nil {
				return nil, err
			}
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
