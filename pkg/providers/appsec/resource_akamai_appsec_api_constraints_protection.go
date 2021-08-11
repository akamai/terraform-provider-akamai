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
func resourceAPIConstraintsProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIConstraintsProtectionCreate,
		ReadContext:   resourceAPIConstraintsProtectionRead,
		UpdateContext: resourceAPIConstraintsProtectionUpdate,
		DeleteContext: resourceAPIConstraintsProtectionDelete,
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

func resourceAPIConstraintsProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionCreate")
	logger.Debugf("in resourceAPIConstraintsProtectionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "apiConstraintsProtection", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	applyapiconstraintscontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createAPIConstraintsProtection := appsec.UpdateAPIConstraintsProtectionRequest{
		ConfigID:            configid,
		Version:             version,
		PolicyID:            policyid,
		ApplyAPIConstraints: applyapiconstraintscontrols,
	}

	_, erru := client.UpdateAPIConstraintsProtection(ctx, createAPIConstraintsProtection)
	if erru != nil {
		logger.Errorf("calling 'createAPIConstraintsProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createAPIConstraintsProtection.ConfigID, createAPIConstraintsProtection.PolicyID))

	return resourceAPIConstraintsProtectionRead(ctx, d, m)
}

func resourceAPIConstraintsProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionRead")
	logger.Debugf("in resourceAPIConstraintsProtectionRead")

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

	getAPIConstraintsProtection := appsec.GetAPIConstraintsProtectionRequest{
		ConfigID: configid,
		Version:  version,
		PolicyID: policyid,
	}

	protection, err := client.GetAPIConstraintsProtection(ctx, getAPIConstraintsProtection)
	if err != nil {
		logger.Errorf("calling 'getAPIConstraintsProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAPIConstraintsProtection.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getAPIConstraintsProtection.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enabled", protection.ApplyAPIConstraints); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "rateProtectionDS", protection)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceAPIConstraintsProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionUpdate")
	logger.Debugf("in resourceAPIConstraintsProtectionUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "apiConstraintsProtection", m)
	policyid := idParts[1]
	applyapiconstraintscontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateAPIConstraintsProtection := appsec.UpdateAPIConstraintsProtectionRequest{
		ConfigID:            configid,
		Version:             version,
		PolicyID:            policyid,
		ApplyAPIConstraints: applyapiconstraintscontrols,
	}

	_, erru := client.UpdateAPIConstraintsProtection(ctx, updateAPIConstraintsProtection)
	if erru != nil {
		logger.Errorf("calling 'updateAPIConstraintsProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceAPIConstraintsProtectionRead(ctx, d, m)
}

func resourceAPIConstraintsProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionDelete")
	logger.Debugf("in resourceAPIConstraintsProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "apiConstraintsProtection", m)
	policyid := idParts[1]

	removeAPIConstraintsProtection := appsec.UpdateAPIConstraintsProtectionRequest{
		ConfigID:            configid,
		Version:             version,
		PolicyID:            policyid,
		ApplyAPIConstraints: false,
	}

	_, errd := client.UpdateAPIConstraintsProtection(ctx, removeAPIConstraintsProtection)
	if errd != nil {
		logger.Errorf("calling 'updateAPIConstraintsProtection': %s", errd.Error())
		return diag.FromErr(errd)
	}
	d.SetId("")
	return nil
}
