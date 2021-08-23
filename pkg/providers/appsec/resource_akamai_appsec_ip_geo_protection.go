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
func resourceIPGeoProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPGeoProtectionCreate,
		ReadContext:   resourceIPGeoProtectionRead,
		UpdateContext: resourceIPGeoProtectionUpdate,
		DeleteContext: resourceIPGeoProtectionDelete,
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

func resourceIPGeoProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionCreate")
	logger.Debugf("in resourceIPGeoProtectionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ipgeoProtection", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enabled, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	getPolicyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configid,
		Version:  version,
		PolicyID: policyid,
	}

	policyProtections, err := client.GetPolicyProtections(ctx, getPolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	updatePolicyProtectionsRequest := appsec.UpdatePolicyProtectionsRequest{
		ConfigID:                      configid,
		Version:                       version,
		PolicyID:                      policyid,
		ApplyAPIConstraints:           policyProtections.ApplyAPIConstraints,
		ApplyApplicationLayerControls: policyProtections.ApplyApplicationLayerControls,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     enabled,
		ApplyRateControls:             policyProtections.ApplyRateControls,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         policyProtections.ApplySlowPostControls,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configid, policyid))
	return resourceIPGeoProtectionRead(ctx, d, m)
}

func resourceIPGeoProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionRead")
	logger.Debugf("in resourceIPGeoProtectionRead")

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

	policyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configid,
		Version:  version,
		PolicyID: policyid,
	}

	policyProtections, err := client.GetPolicyProtections(ctx, policyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}
	enabled := policyProtections.ApplyNetworkLayerControls

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", policyid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enabled", enabled); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "networkProtectionDS", enabled)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceIPGeoProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionUpdate")
	logger.Debugf("in resourceIPGeoProtectionUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "networkProtection", m)
	policyid := idParts[1]
	enabled, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	getPolicyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configid,
		Version:  version,
		PolicyID: policyid,
	}
	policyProtections, err := client.GetPolicyProtections(ctx, getPolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	updatePolicyProtectionsRequest := appsec.UpdatePolicyProtectionsRequest{
		ConfigID:                      configid,
		Version:                       version,
		PolicyID:                      policyid,
		ApplyAPIConstraints:           policyProtections.ApplyAPIConstraints,
		ApplyApplicationLayerControls: policyProtections.ApplyApplicationLayerControls,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     enabled,
		ApplyRateControls:             policyProtections.ApplyRateControls,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         policyProtections.ApplySlowPostControls,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceIPGeoProtectionRead(ctx, d, m)
}

func resourceIPGeoProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionDelete")
	logger.Debugf("in resourceIPGeoProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ipgeoProtection", m)
	policyid := idParts[1]

	getPolicyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configid,
		Version:  version,
		PolicyID: policyid,
	}
	policyProtections, err := client.GetPolicyProtections(ctx, getPolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	updatePolicyProtectionsRequest := appsec.UpdatePolicyProtectionsRequest{
		ConfigID:                      configid,
		Version:                       version,
		PolicyID:                      policyid,
		ApplyAPIConstraints:           policyProtections.ApplyAPIConstraints,
		ApplyApplicationLayerControls: policyProtections.ApplyApplicationLayerControls,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     false,
		ApplyRateControls:             policyProtections.ApplyRateControls,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         policyProtections.ApplySlowPostControls,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
