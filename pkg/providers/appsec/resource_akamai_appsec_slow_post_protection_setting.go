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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSlowPostProtectionSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlowPostProtectionSettingCreate,
		ReadContext:   resourceSlowPostProtectionSettingRead,
		UpdateContext: resourceSlowPostProtectionSettingUpdate,
		DeleteContext: resourceSlowPostProtectionSettingDelete,
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
			"slow_rate_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Alert,
					Abort,
				}, false),
			},
			"slow_rate_threshold_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"slow_rate_threshold_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"duration_threshold_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

func resourceSlowPostProtectionSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingCreate")
	logger.Debugf("in resourceSlowPostProtectionSettingCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "slowpostSettings", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowrateaction, err := tools.GetStringValue("slow_rate_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowratethresholdrate, err := tools.GetIntValue("slow_rate_threshold_rate", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowratethresholdperiod, err := tools.GetIntValue("slow_rate_threshold_period", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	durationthresholdtimeout, err := tools.GetIntValue("duration_threshold_timeout", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createSlowPostProtectionSetting := appsec.UpdateSlowPostProtectionSettingRequest{}
	createSlowPostProtectionSetting.ConfigID = configid
	createSlowPostProtectionSetting.Version = version
	createSlowPostProtectionSetting.PolicyID = policyid
	createSlowPostProtectionSetting.Action = slowrateaction
	createSlowPostProtectionSetting.SlowRateThreshold.Rate = slowratethresholdrate
	createSlowPostProtectionSetting.SlowRateThreshold.Period = slowratethresholdperiod
	createSlowPostProtectionSetting.DurationThreshold.Timeout = durationthresholdtimeout

	_, erru := client.UpdateSlowPostProtectionSetting(ctx, createSlowPostProtectionSetting)
	if erru != nil {
		logger.Errorf("calling 'updateSlowPostProtectionSetting': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createSlowPostProtectionSetting.ConfigID, createSlowPostProtectionSetting.PolicyID))

	return resourceSlowPostProtectionSettingRead(ctx, d, m)
}

func resourceSlowPostProtectionSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingRead")
	logger.Debugf("in resourceSlowPostProtectionSettingRead")

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

	getSlowPostProtectionSettingsRequest := appsec.GetSlowPostProtectionSettingsRequest{
		ConfigID: configid,
		Version:  version,
		PolicyID: policyid,
	}

	slowPostProtectionSettings, errg := client.GetSlowPostProtectionSettings(ctx, getSlowPostProtectionSettingsRequest)
	if errg != nil {
		logger.Errorf("calling 'getSlowPostProtectionSettings': %s", errg.Error())
		return diag.FromErr(errg)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", policyid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("slow_rate_action", slowPostProtectionSettings.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if slowPostProtectionSettings.SlowRateThreshold != nil {
		if err := d.Set("slow_rate_threshold_rate", slowPostProtectionSettings.SlowRateThreshold.Rate); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		if err := d.Set("slow_rate_threshold_period", slowPostProtectionSettings.SlowRateThreshold.Period); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	if slowPostProtectionSettings.DurationThreshold != nil {
		if err := d.Set("duration_threshold_timeout", slowPostProtectionSettings.DurationThreshold.Timeout); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	return nil
}

func resourceSlowPostProtectionSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingUpdate")
	logger.Debugf("in resourceSlowPostProtectionSettingUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "slowpostSettings", m)
	policyid := idParts[1]
	slowrateaction, err := tools.GetStringValue("slow_rate_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowratethresholdrate, err := tools.GetIntValue("slow_rate_threshold_rate", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowratethresholdperiod, err := tools.GetIntValue("slow_rate_threshold_period", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	durationthresholdtimeout, err := tools.GetIntValue("duration_threshold_timeout", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSlowPostProtectionSetting := appsec.UpdateSlowPostProtectionSettingRequest{}
	updateSlowPostProtectionSetting.ConfigID = configid
	updateSlowPostProtectionSetting.Version = version
	updateSlowPostProtectionSetting.PolicyID = policyid
	updateSlowPostProtectionSetting.Action = slowrateaction
	updateSlowPostProtectionSetting.SlowRateThreshold.Rate = slowratethresholdrate
	updateSlowPostProtectionSetting.SlowRateThreshold.Period = slowratethresholdperiod
	updateSlowPostProtectionSetting.DurationThreshold.Timeout = durationthresholdtimeout

	_, erru := client.UpdateSlowPostProtectionSetting(ctx, updateSlowPostProtectionSetting)
	if erru != nil {
		logger.Errorf("calling 'updateSlowPostProtectionSetting': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSlowPostProtectionSettingRead(ctx, d, m)
}

func resourceSlowPostProtectionSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingDelete")
	logger.Debugf("in resourceSlowPostProtectionSettingDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "slowpostSettings", m)
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
		ApplyNetworkLayerControls:     policyProtections.ApplyNetworkLayerControls,
		ApplyRateControls:             policyProtections.ApplyRateControls,
		ApplyReputationControls:       policyProtections.ApplyReputationControls,
		ApplySlowPostControls:         false,
	}
	_, err = client.UpdatePolicyProtections(ctx, updatePolicyProtectionsRequest)
	if err != nil {
		logger.Errorf("calling UpdatePolicyProtections: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// Definition of constant variables
const (
	Abort = "abort"
)
