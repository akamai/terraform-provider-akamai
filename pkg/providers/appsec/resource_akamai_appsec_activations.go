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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceActivations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActivationsCreate,
		ReadContext:   resourceActivationsRead,
		DeleteContext: resourceActivationsDelete,
		UpdateContext: resourceActivationsUpdate,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,

				Default: "STAGING",
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,

				Default: "Activation Notes",
			},
			"activate": {
				Type:     schema.TypeBool,
				Optional: true,

				Default: true,
			},
			"notification_emails": {
				Type:     schema.TypeSet,
				Required: true,

				Elem: &schema.Schema{Type: schema.TypeString},
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

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		logger.Debugf("Done")
		return nil
	}

	createActivations := appsec.CreateActivationsRequest{}

	createActivationsreq := appsec.GetActivationsRequest{}

	ap := appsec.ActivationConfigs{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	ap.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	ap.ConfigVersion = version

	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createActivations.Network = network

	createActivations.Action = "ACTIVATE"
	createActivations.ActivationConfigs = append(createActivations.ActivationConfigs, ap)
	createActivations.NotificationEmails = tools.SetToStringSlice(d.Get("notification_emails").(*schema.Set))

	postresp, err := client.CreateActivations(ctx, createActivations, true)
	if err != nil {
		logger.Errorf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	createActivationsreq.ActivationID = postresp.ActivationID
	activation, err := lookupActivation(ctx, client, createActivationsreq)
	for activation.Status != appsec.StatusActive {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, createActivationsreq)

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

func resourceActivationsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsUpdate")

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		logger.Debugf("Done")
		return nil
	}

	createActivations := appsec.CreateActivationsRequest{}

	getActivationsreq := appsec.GetActivationsRequest{}

	ap := appsec.ActivationConfigs{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	ap.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	ap.ConfigVersion = version

	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createActivations.Network = network

	createActivations.Action = "ACTIVATE"
	createActivations.ActivationConfigs = append(createActivations.ActivationConfigs, ap)
	createActivations.NotificationEmails = tools.SetToStringSlice(d.Get("notification_emails").(*schema.Set))

	activations, err := client.CreateActivations(ctx, createActivations, true)
	if err != nil {
		logger.Errorf("calling 'createActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(activations.ActivationID))

	if err := d.Set("status", activations.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	getActivationsreq.ActivationID = activations.ActivationID
	activation, err := lookupActivation(ctx, client, getActivationsreq)
	for activation.Status != appsec.StatusActive {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, getActivationsreq)

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

	activate, err := tools.GetBoolValue("activate", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if !activate {
		d.SetId("none")
		logger.Debugf("Done")
		return nil
	}

	removeActivations := appsec.RemoveActivationsRequest{}
	createActivationsreq := appsec.GetActivationsRequest{}

	if d.Id() == "" || d.Id() == "none" {
		return nil
	}

	ActivationID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(err)
	}
	removeActivations.ActivationID = ActivationID

	ap := appsec.ActivationConfigs{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	ap.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	ap.ConfigVersion = version

	network, err := tools.GetStringValue("network", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeActivations.Network = network

	removeActivations.NotificationEmails = tools.SetToStringSlice(d.Get("notification_emails").(*schema.Set))

	removeActivations.Action = "DEACTIVATE"

	removeActivations.ActivationConfigs = append(removeActivations.ActivationConfigs, ap)

	postresp, err := client.RemoveActivations(ctx, removeActivations)

	if err != nil {
		logger.Errorf("calling 'removeActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(postresp.ActivationID))

	if err := d.Set("status", postresp.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	createActivationsreq.ActivationID = postresp.ActivationID

	activation, err := lookupActivation(ctx, client, createActivationsreq)
	for activation.Status != appsec.StatusDeactivated {
		select {
		case <-time.After(tools.MaxDuration(ActivationPollInterval, ActivationPollMinimum)):
			act, err := client.GetActivations(ctx, createActivationsreq)

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

func resourceActivationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceActivationsRead")

	getActivations := appsec.GetActivationsRequest{}

	if d.Id() == "" || d.Id() == "none" {
		return nil
	}

	activationID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getActivations.ActivationID = activationID

	activations, err := client.GetActivations(ctx, getActivations)
	if err != nil {
		logger.Errorf("calling 'getActivations': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("status", activations.Status); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(strconv.Itoa(activations.ActivationID))

	return nil
}

func lookupActivation(ctx context.Context, client appsec.APPSEC, query appsec.GetActivationsRequest) (*appsec.GetActivationsResponse, error) {
	activations, err := client.GetActivations(ctx, query)
	if err != nil {
		return nil, err
	}

	return activations, nil

	return nil, nil
}
