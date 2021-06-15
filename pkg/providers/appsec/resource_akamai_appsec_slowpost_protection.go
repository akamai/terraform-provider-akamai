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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSlowPostProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlowPostProtectionCreate,
		ReadContext:   resourceSlowPostProtectionRead,
		UpdateContext: resourceSlowPostProtectionUpdate,
		DeleteContext: resourceSlowPostProtectionDelete,
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
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

func resourceSlowPostProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionCreate")
	logger.Debugf("!!! in resourceCustomRuleActionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "slowpostProtection", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	applyslowpostcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createSlowPostProtection := appsec.UpdateSlowPostProtectionRequest{}
	createSlowPostProtection.ConfigID = configid
	createSlowPostProtection.Version = version
	createSlowPostProtection.PolicyID = policyid
	createSlowPostProtection.ApplySlowPostControls = applyslowpostcontrols

	_, erru := client.UpdateSlowPostProtection(ctx, createSlowPostProtection)
	if erru != nil {
		logger.Errorf("calling 'createSlowPostProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createSlowPostProtection.ConfigID, createSlowPostProtection.PolicyID))

	return resourceSlowPostProtectionRead(ctx, d, m)
}

func resourceSlowPostProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionRead")
	logger.Debugf("!!! in resourceSlowPostProtectionRead")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]

	getSlowPostProtection := appsec.GetSlowPostProtectionRequest{}
	getSlowPostProtection.ConfigID = configid
	getSlowPostProtection.Version = version
	getSlowPostProtection.PolicyID = policyid

	slowpostprotection, err := client.GetSlowPostProtection(ctx, getSlowPostProtection)
	if err != nil {
		logger.Errorf("calling 'getSlowPostProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getSlowPostProtection.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getSlowPostProtection.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enabled", slowpostprotection.ApplySlowPostControls); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "slowpostProtectionDS", slowpostprotection)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceSlowPostProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionUpdate")
	logger.Debugf("!!! in resourceSlowPostProtectionUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "slowpostProtection", m)
	policyid := idParts[1]
	applyslowpostcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSlowPostProtection := appsec.UpdateSlowPostProtectionRequest{}
	updateSlowPostProtection.ConfigID = configid
	updateSlowPostProtection.Version = version
	updateSlowPostProtection.PolicyID = policyid
	updateSlowPostProtection.ApplySlowPostControls = applyslowpostcontrols

	_, erru := client.UpdateSlowPostProtection(ctx, updateSlowPostProtection)
	if erru != nil {
		logger.Errorf("calling 'updateSlowPostProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSlowPostProtectionRead(ctx, d, m)
}

func resourceSlowPostProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionRemove")
	logger.Debugf("!!! in resourceSlowPostProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "slowpostProtection", m)
	policyid := idParts[1]

	removeSlowPostProtection := appsec.UpdateSlowPostProtectionRequest{}
	removeSlowPostProtection.ConfigID = configid
	removeSlowPostProtection.Version = version
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
