package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"slow_rate_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					Alert,
					Abort,
				}, false)),
				Description: "Action to be taken when slow POST protection is triggered",
			},
			"slow_rate_threshold_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Average rate (in bytes per second over the specified time period) allowed before the specified action is triggered",
			},
			"slow_rate_threshold_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Amount of time (in seconds) that the server should allow a request before marking the request as being too slow",
			},
			"duration_threshold_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum amount of time (in seconds) within which the first 8KB of the POST body must be received to avoid triggering the specified action",
			},
		},
	}
}

func resourceSlowPostProtectionSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingCreate")
	logger.Debugf("in resourceSlowPostProtectionSettingCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "slowpostSettings", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	slowrateaction, err := tf.GetStringValue("slow_rate_action", d)
	if err != nil {
		return diag.FromErr(err)
	}
	slowratethresholdrate, err := tf.GetIntValue("slow_rate_threshold_rate", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowratethresholdperiod, err := tf.GetIntValue("slow_rate_threshold_period", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	durationthresholdtimeout, err := tf.GetIntValue("duration_threshold_timeout", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	createSlowPostProtectionSetting := appsec.UpdateSlowPostProtectionSettingRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Action:   slowrateaction,
	}
	createSlowPostProtectionSetting.SlowRateThreshold.Rate = slowratethresholdrate
	createSlowPostProtectionSetting.SlowRateThreshold.Period = slowratethresholdperiod
	createSlowPostProtectionSetting.DurationThreshold.Timeout = durationthresholdtimeout

	_, err = client.UpdateSlowPostProtectionSetting(ctx, createSlowPostProtectionSetting)
	if err != nil {
		logger.Errorf("calling 'updateSlowPostProtectionSetting': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createSlowPostProtectionSetting.ConfigID, createSlowPostProtectionSetting.PolicyID))

	return resourceSlowPostProtectionSettingRead(ctx, d, m)
}

func resourceSlowPostProtectionSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingRead")
	logger.Debugf("in resourceSlowPostProtectionSettingRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	getSlowPostProtectionSettingsRequest := appsec.GetSlowPostProtectionSettingsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	slowPostProtectionSettings, err := client.GetSlowPostProtectionSettings(ctx, getSlowPostProtectionSettingsRequest)
	if err != nil {
		logger.Errorf("calling 'getSlowPostProtectionSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err = d.Set("security_policy_id", policyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err = d.Set("slow_rate_action", slowPostProtectionSettings.Action); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if slowPostProtectionSettings.SlowRateThreshold != nil {
		if err := d.Set("slow_rate_threshold_rate", slowPostProtectionSettings.SlowRateThreshold.Rate); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
		if err := d.Set("slow_rate_threshold_period", slowPostProtectionSettings.SlowRateThreshold.Period); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	if slowPostProtectionSettings.DurationThreshold != nil {
		if err := d.Set("duration_threshold_timeout", slowPostProtectionSettings.DurationThreshold.Timeout); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	return nil
}

func resourceSlowPostProtectionSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingUpdate")
	logger.Debugf("in resourceSlowPostProtectionSettingUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "slowpostSettings", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	slowrateaction, err := tf.GetStringValue("slow_rate_action", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowratethresholdrate, err := tf.GetIntValue("slow_rate_threshold_rate", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	slowratethresholdperiod, err := tf.GetIntValue("slow_rate_threshold_period", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	durationthresholdtimeout, err := tf.GetIntValue("duration_threshold_timeout", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSlowPostProtectionSetting := appsec.UpdateSlowPostProtectionSettingRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Action:   slowrateaction,
	}
	updateSlowPostProtectionSetting.SlowRateThreshold.Rate = slowratethresholdrate
	updateSlowPostProtectionSetting.SlowRateThreshold.Period = slowratethresholdperiod
	updateSlowPostProtectionSetting.DurationThreshold.Timeout = durationthresholdtimeout

	_, err = client.UpdateSlowPostProtectionSetting(ctx, updateSlowPostProtectionSetting)
	if err != nil {
		logger.Errorf("calling 'updateSlowPostProtectionSetting': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceSlowPostProtectionSettingRead(ctx, d, m)
}

func resourceSlowPostProtectionSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingDelete")
	logger.Debugf("in resourceSlowPostProtectionSettingDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "slowpostSettings", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	request := appsec.UpdateSlowPostProtectionRequest{
		ConfigID:              configID,
		Version:               version,
		PolicyID:              policyID,
		ApplySlowPostControls: false,
	}
	_, err = client.UpdateSlowPostProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateSlowPostProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}

// Definition of constant variables
const (
	Abort = "abort"
)
