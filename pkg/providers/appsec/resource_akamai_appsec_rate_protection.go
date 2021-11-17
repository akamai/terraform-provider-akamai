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
	logger.Debugf("in resourceRateProtectionCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "rateProtection", m)
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enabled, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	getPolicyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	policyProtections, err := client.GetPolicyProtections(ctx, getPolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	updatePolicyProtectionsRequest := appsec.UpdatePolicyProtectionsRequest{
		ConfigID:                      configID,
		Version:                       version,
		PolicyID:                      policyID,
		ApplyAPIConstraints:           policyProtections.ApplyAPIConstraints,
		ApplyApplicationLayerControls: policyProtections.ApplyApplicationLayerControls,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     policyProtections.ApplyNetworkLayerControls,
		ApplyRateControls:             enabled,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         policyProtections.ApplySlowPostControls,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))

	return resourceRateProtectionRead(ctx, d, m)
}

func resourceRateProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionRead")
	logger.Debugf("in resourceRateProtectionRead")

	idParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	policyID := idParts[1]

	policyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	policyProtections, err := client.GetPolicyProtections(ctx, policyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}
	enabled := policyProtections.ApplyRateControls

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", policyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("enabled", enabled); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "rateProtectionDS", policyProtections)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	return nil
}

func resourceRateProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionUpdate")
	logger.Debugf("in resourceRateProtectionUpdate")

	idParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "rateProtection", m)
	policyID := idParts[1]
	enabled, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	getPolicyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	policyProtections, err := client.GetPolicyProtections(ctx, getPolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	updatePolicyProtectionsRequest := appsec.UpdatePolicyProtectionsRequest{
		ConfigID:                      configID,
		Version:                       version,
		PolicyID:                      policyID,
		ApplyAPIConstraints:           policyProtections.ApplyAPIConstraints,
		ApplyApplicationLayerControls: policyProtections.ApplyApplicationLayerControls,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     policyProtections.ApplyNetworkLayerControls,
		ApplyRateControls:             enabled,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         policyProtections.ApplySlowPostControls,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceRateProtectionRead(ctx, d, m)
}

func resourceRateProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionDelete")
	logger.Debugf("in resourceRateProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "rateProtection", m)
	policyID := idParts[1]

	getPolicyProtectionsRequest := appsec.GetPolicyProtectionsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	policyProtections, err := client.GetPolicyProtections(ctx, getPolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling GetPolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	updatePolicyProtectionsRequest := appsec.UpdatePolicyProtectionsRequest{
		ConfigID:                      configID,
		Version:                       version,
		PolicyID:                      policyID,
		ApplyAPIConstraints:           policyProtections.ApplyAPIConstraints,
		ApplyApplicationLayerControls: policyProtections.ApplyApplicationLayerControls,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     policyProtections.ApplyNetworkLayerControls,
		ApplyRateControls:             false,
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
