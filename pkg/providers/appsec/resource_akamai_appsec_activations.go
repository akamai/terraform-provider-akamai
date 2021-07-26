package appsec

import (
	"context"
	"errors"
	"fmt"

	// "log"
	"strconv"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STAGING",
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Activation Notes",
			},
			"activate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"notification_emails": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
	logger.Debug("!!! in resourceActivationsCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	note, err := tools.GetStringValue("notes", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	activate, err := tools.GetBoolValue("activate", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	notificationEmailsSet, err := tools.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tools.SetToStringSlice(notificationEmailsSet)

	if !activate {
		d.SetId("none")
		return nil
	}

	activationConfig := appsec.ActivationConfigs{}
	activationConfig.ConfigID = configid
	activationConfig.ConfigVersion = version

	createActivationRequest := appsec.CreateActivationsRequest{}
	createActivationRequest.Action = "ACTIVATE"
	createActivationRequest.Network = network
	createActivationRequest.Note = note
	createActivationRequest.ActivationConfigs = append(createActivationRequest.ActivationConfigs, activationConfig)
	createActivationRequest.NotificationEmails = notificationEmails

	postresp, err := client.CreateActivations(ctx, createActivationRequest, true)
	if err != nil {
		logger.Errorf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	getActivationRequest := appsec.GetActivationsRequest{}
	getActivationRequest.ActivationID = postresp.ActivationID
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
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRead")
	logger.Debug("!!! in resourceActivationsRead")

	activationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getActivations := appsec.GetActivationsRequest{}
	getActivations.ActivationID = activationID

	activations, err := client.GetActivations(ctx, getActivations)
	if err != nil {
		logger.Errorf("calling 'getActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("status", activations.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceActivationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsUpdate")
	logger.Debug("!!! in resourceActivationsUpdate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	note, err := tools.GetStringValue("notes", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	activate, err := tools.GetBoolValue("activate", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	notificationEmailsSet, err := tools.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tools.SetToStringSlice(notificationEmailsSet)

	if !activate {
		d.SetId("none")
		return nil
	}

	activationConfig := appsec.ActivationConfigs{}
	activationConfig.ConfigID = configid
	activationConfig.ConfigVersion = version

	createActivationRequest := appsec.CreateActivationsRequest{}
	createActivationRequest.Action = "ACTIVATE"
	createActivationRequest.Network = network
	createActivationRequest.Note = note
	createActivationRequest.ActivationConfigs = append(createActivationRequest.ActivationConfigs, activationConfig)
	createActivationRequest.NotificationEmails = notificationEmails

	postresp, err := client.CreateActivations(ctx, createActivationRequest, true)
	if err != nil {
		logger.Errorf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	getActivationRequest := appsec.GetActivationsRequest{}
	getActivationRequest.ActivationID = postresp.ActivationID
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
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	return resourceActivationsRead(ctx, d, m)
}

func resourceActivationsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRemove")
	logger.Debug("!!! in resourceActivationsDelete")

	activationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	activate, err := tools.GetBoolValue("activate", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	notificationEmailsSet, err := tools.GetSetValue("notification_emails", d)
	if err != nil {
		return diag.FromErr(err)
	}
	notificationEmails := tools.SetToStringSlice(notificationEmailsSet)

	if !activate {
		d.SetId("none")
		return nil
	}

	activationConfig := appsec.ActivationConfigs{}
	activationConfig.ConfigID = configid
	activationConfig.ConfigVersion = version

	removeActivationRequest := appsec.RemoveActivationsRequest{}
	removeActivationRequest.ActivationID = activationID
	removeActivationRequest.Action = "DEACTIVATE"
	removeActivationRequest.Network = network
	removeActivationRequest.ActivationConfigs = append(removeActivationRequest.ActivationConfigs, activationConfig)
	removeActivationRequest.NotificationEmails = notificationEmails

	postresp, err := client.RemoveActivations(ctx, removeActivationRequest)
	if err != nil {
		logger.Errorf("calling 'removeActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	getActivationRequest := appsec.GetActivationsRequest{}
	getActivationRequest.ActivationID = activationID

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
			return diag.FromErr(fmt.Errorf("activation context terminated: %w", ctx.Err()))
		}
	}

	if err := d.Set("status", activation.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId("")

	return nil
}

func lookupActivation(ctx context.Context, client appsec.APPSEC, query appsec.GetActivationsRequest) (*appsec.GetActivationsResponse, error) {
	activations, err := client.GetActivations(ctx, query)
	if err != nil {
		return nil, err
	}

	return activations, nil
}
