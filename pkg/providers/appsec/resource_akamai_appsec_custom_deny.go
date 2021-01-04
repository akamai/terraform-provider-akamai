package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceCustomDeny() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomDenyCreate,
		ReadContext:   resourceCustomDenyRead,
		UpdateContext: resourceCustomDenyUpdate,
		DeleteContext: resourceCustomDenyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_deny": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"custom_deny_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "custom_deny_id",
			},
		},
	}
}

func resourceCustomDenyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyCreate")

	createCustomDeny := appsec.CreateCustomDenyRequest{}

	jsonpostpayload := d.Get("custom_deny")

	if err := json.Unmarshal([]byte(jsonpostpayload.(string)), &createCustomDeny); err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createCustomDeny.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createCustomDeny.Version = version

	postresp, errc := client.CreateCustomDeny(ctx, createCustomDeny)
	if errc != nil {
		logger.Errorf("calling 'createCustomDeny': %s", errc.Error())
		return diag.FromErr(errc)
	}

	d.SetId(string(postresp.ID))

	return resourceCustomDenyRead(ctx, d, m)
}

func resourceCustomDenyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyUpdate")

	updateCustomDeny := appsec.UpdateCustomDenyRequest{}

	jsonpostpayload := d.Get("custom_deny")

	if err := json.Unmarshal([]byte(jsonpostpayload.(string)), &updateCustomDeny); err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomDeny.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomDeny.Version = version

	updateCustomDeny.ID = d.Id()

	_, erru := client.UpdateCustomDeny(ctx, updateCustomDeny)
	if erru != nil {
		logger.Errorf("calling 'updateCustomDeny': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceCustomDenyRead(ctx, d, m)
}

func resourceCustomDenyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyRemove")

	removeCustomDeny := appsec.RemoveCustomDenyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeCustomDeny.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeCustomDeny.Version = version

	removeCustomDeny.ID = d.Id()

	_, errd := client.RemoveCustomDeny(ctx, removeCustomDeny)
	if errd != nil {
		logger.Errorf("calling 'removeCustomDeny': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceCustomDenyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyRead")

	getCustomDeny := appsec.GetCustomDenyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomDeny.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomDeny.Version = version

	getCustomDeny.ID = d.Id()

	customdeny, err := client.GetCustomDeny(ctx, getCustomDeny)
	if err != nil {
		logger.Errorf("calling 'getCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "CustomDenyDS", customdeny)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	if err := d.Set("custom_deny_id", customdeny.ID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(string(customdeny.ID))

	return nil
}
