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
func resourceWAFProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWAFProtectionCreate,
		ReadContext:   resourceWAFProtectionRead,
		UpdateContext: resourceWAFProtectionUpdate,
		DeleteContext: resourceWAFProtectionDelete,
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceWAFProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionCreate")
	logger.Debugf("!!! in resourceWAFProtectionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "wafProtection", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	applyapplicationlayercontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createWAFProtection := appsec.UpdateWAFProtectionRequest{}
	createWAFProtection.ConfigID = configid
	createWAFProtection.Version = version
	createWAFProtection.PolicyID = policyid
	createWAFProtection.ApplyApplicationLayerControls = applyapplicationlayercontrols

	_, erru := client.UpdateWAFProtection(ctx, createWAFProtection)
	if erru != nil {
		logger.Errorf("calling 'createWAFProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createWAFProtection.ConfigID, createWAFProtection.PolicyID))

	return resourceWAFProtectionRead(ctx, d, m)
}

func resourceWAFProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionRead")
	logger.Debugf("!!! in resourceSlowPostProtectionSettingRead")

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

	getWAFProtection := appsec.GetWAFProtectionRequest{}
	getWAFProtection.ConfigID = configid
	getWAFProtection.Version = version
	getWAFProtection.PolicyID = policyid

	wafprotection, err := client.GetWAFProtection(ctx, getWAFProtection)
	if err != nil {
		logger.Errorf("calling 'getWAFProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getWAFProtection.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getWAFProtection.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enabled", wafprotection.ApplyApplicationLayerControls); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "wafProtectionDS", wafprotection)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceWAFProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionUpdate")
	logger.Debugf("!!! in resourceSlowPostProtectionSettingUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "wafProtection", m)
	policyid := idParts[1]
	applyapplicationlayercontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateWAFProtection := appsec.UpdateWAFProtectionRequest{}
	updateWAFProtection.ConfigID = configid
	updateWAFProtection.Version = version
	updateWAFProtection.PolicyID = policyid
	updateWAFProtection.ApplyApplicationLayerControls = applyapplicationlayercontrols

	_, erru := client.UpdateWAFProtection(ctx, updateWAFProtection)
	if erru != nil {
		logger.Errorf("calling 'updateWAFProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceWAFProtectionRead(ctx, d, m)
}

func resourceWAFProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionDelete")
	logger.Debugf("!!! in resourceWAFProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "wafProtection", m)
	policyid := idParts[1]

	removeWAFProtection := appsec.UpdateWAFProtectionRequest{}
	removeWAFProtection.ConfigID = configid
	removeWAFProtection.Version = version
	removeWAFProtection.PolicyID = policyid
	removeWAFProtection.ApplyApplicationLayerControls = false

	_, errd := client.UpdateWAFProtection(ctx, removeWAFProtection)
	if errd != nil {
		logger.Errorf("calling 'removeWAFProtection': %s", errd.Error())
		return diag.FromErr(errd)
	}
	d.SetId("")
	return nil
}
