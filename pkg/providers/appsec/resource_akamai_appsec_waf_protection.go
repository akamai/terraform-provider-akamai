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

func resourceWAFProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionCreate")
	logger.Debugf("in resourceWAFProtectionCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "wafProtection", m)
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
		ApplyApplicationLayerControls: enabled,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     policyProtections.ApplyNetworkLayerControls,
		ApplyRateControls:             policyProtections.ApplyRateControls,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         policyProtections.ApplySlowPostControls,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))
	return resourceWAFProtectionRead(ctx, d, m)
}

func resourceWAFProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionRead")
	logger.Debugf("in resourceWAFProtectionRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	policyID := iDParts[1]

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
	enabled := policyProtections.ApplyApplicationLayerControls

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
	outputtext, err := RenderTemplates(ots, "wafProtectionDS", policyProtections)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	return nil
}

func resourceWAFProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionUpdate")
	logger.Debugf("in resourceWAFProtectionUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "wafProtection", m)
	policyID := iDParts[1]
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
		ApplyApplicationLayerControls: enabled,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     policyProtections.ApplyNetworkLayerControls,
		ApplyRateControls:             policyProtections.ApplyRateControls,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         policyProtections.ApplySlowPostControls,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceWAFProtectionRead(ctx, d, m)
}

func resourceWAFProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFProtectionDelete")
	logger.Debugf("in resourceWAFProtectionDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "wafProtection", m)
	policyID := iDParts[1]

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
		ApplyApplicationLayerControls: false,
		ApplyBotmanControls:           policyProtections.ApplyBotmanControls,
		ApplyNetworkLayerControls:     policyProtections.ApplyNetworkLayerControls,
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
