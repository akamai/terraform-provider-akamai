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
func resourceRateProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRateProtectionCreate,
		ReadContext:   resourceRateProtectionRead,
		UpdateContext: resourceRateProtectionUpdate,
		DeleteContext: resourceRateProtectionDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
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

func resourceRateProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionCreate")
	logger.Debugf("!!! in resourceRateProtectionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "rateProtection", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	applyratecontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createRateProtection := appsec.UpdateRateProtectionRequest{}
	createRateProtection.ConfigID = configid
	createRateProtection.Version = version
	createRateProtection.PolicyID = policyid
	createRateProtection.ApplyRateControls = applyratecontrols

	_, erru := client.UpdateRateProtection(ctx, createRateProtection)
	if erru != nil {
		logger.Errorf("calling 'createRateProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createRateProtection.ConfigID, createRateProtection.PolicyID))

	return resourceRateProtectionRead(ctx, d, m)
}

func resourceRateProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionRead")
	logger.Debugf("!!! in resourceReputationProtectionRead")

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

	getRateProtection := appsec.GetRateProtectionRequest{}
	getRateProtection.ConfigID = configid
	getRateProtection.Version = version
	getRateProtection.PolicyID = policyid

	rateprotection, err := client.GetRateProtection(ctx, getRateProtection)
	if err != nil {
		logger.Errorf("calling 'getRateProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getRateProtection.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getRateProtection.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enabled", rateprotection.ApplyRateControls); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "rateProtectionDS", rateprotection)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceRateProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionUpdate")
	logger.Debugf("!!! in resourceRateProtectionUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "rateProtection", m)
	policyid := idParts[1]
	applyratecontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateRateProtection := appsec.UpdateRateProtectionRequest{}
	updateRateProtection.ConfigID = configid
	updateRateProtection.Version = version
	updateRateProtection.PolicyID = policyid
	updateRateProtection.ApplyRateControls = applyratecontrols

	_, erru := client.UpdateRateProtection(ctx, updateRateProtection)
	if erru != nil {
		logger.Errorf("calling 'updateRateProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceRateProtectionRead(ctx, d, m)
}

func resourceRateProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionDelete")

	logger.Debugf("!!! in resourceRateProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "rateProtection", m)
	policyid := idParts[1]

	removeRateProtection := appsec.UpdateRateProtectionRequest{}
	removeRateProtection.ConfigID = configid

	removeRateProtection.Version = version
	removeRateProtection.PolicyID = policyid
	removeRateProtection.ApplyRateControls = false

	_, errd := client.UpdateRateProtection(ctx, removeRateProtection)
	if errd != nil {
		logger.Errorf("calling 'updateRateProtection': %s", errd.Error())
		return diag.FromErr(errd)
	}
	d.SetId("")
	return nil
}
