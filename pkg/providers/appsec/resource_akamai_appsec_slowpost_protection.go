package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSlowPostProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlowPostProtectionUpdate,
		ReadContext:   resourceSlowPostProtectionRead,
		UpdateContext: resourceSlowPostProtectionUpdate,
		DeleteContext: resourceSlowPostProtectionDelete,
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
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceSlowPostProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionRead")

	getSlowPostProtection := appsec.GetSlowPostProtectionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtection.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtection.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtection.PolicyID = policyid

	slowpostprotection, err := client.GetSlowPostProtection(ctx, getSlowPostProtection)
	if err != nil {
		logger.Errorf("calling 'getSlowPostProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "slowpostProtectionDS", slowpostprotection)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	if err := d.Set("enabled", slowpostprotection.ApplySlowPostControls); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(getSlowPostProtection.ConfigID))

	return nil
}

func resourceSlowPostProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionRemove")

	removeSlowPostProtection := appsec.UpdateSlowPostProtectionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSlowPostProtection.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSlowPostProtection.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSlowPostProtection.PolicyID = policyid

	removeSlowPostProtection.ApplySlowPostControls = false

	_, errd := client.UpdateSlowPostProtection(ctx, removeSlowPostProtection)
	if errd != nil {
		logger.Errorf("calling 'removeSlowPostProtection': %s", errd.Error())
		return diag.FromErr(errd)
	}
	d.SetId("")
	return nil
}

func resourceSlowPostProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionUpdate")

	updateSlowPostProtection := appsec.UpdateSlowPostProtectionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtection.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtection.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtection.PolicyID = policyid

	applyslowpostcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtection.ApplySlowPostControls = applyslowpostcontrols

	_, erru := client.UpdateSlowPostProtection(ctx, updateSlowPostProtection)
	if erru != nil {
		logger.Errorf("calling 'updateSlowPostProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSlowPostProtectionRead(ctx, d, m)
}
