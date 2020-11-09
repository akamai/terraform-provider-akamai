package appsec

import (
	"context"
	"errors"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSlowPostProtectionSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlowPostProtectionSettingUpdate,
		ReadContext:   resourceSlowPostProtectionSettingRead,
		UpdateContext: resourceSlowPostProtectionSettingUpdate,
		DeleteContext: resourceSlowPostProtectionSettingDelete,
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
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slow_rate_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Alert,
					Deny,
					None,
				}, false),
			},
			"slow_rate_threshold_rate": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"slow_rate_threshold_period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"duration_threshold_timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceSlowPostProtectionSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingRead")

	getSlowPostProtectionSetting := v2.GetSlowPostProtectionSettingRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtectionSetting.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtectionSetting.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtectionSetting.PolicyID = policyid

	_, errg := client.GetSlowPostProtectionSetting(ctx, getSlowPostProtectionSetting)
	if errg != nil {
		logger.Errorf("calling 'getSlowPostProtectionSetting': %s", errg.Error())
		return diag.FromErr(errg)
	}

	d.SetId(strconv.Itoa(getSlowPostProtectionSetting.ConfigID))

	return nil
}

func resourceSlowPostProtectionSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceSlowPostProtectionSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingUpdate")

	updateSlowPostProtectionSetting := v2.UpdateSlowPostProtectionSettingRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtectionSetting.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtectionSetting.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtectionSetting.PolicyID = policyid

	slowrateaction, err := tools.GetStringValue("slow_rate_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtectionSetting.Action = slowrateaction

	slowratethresholdrate, err := tools.GetIntValue("slow_rate_threshold_rate", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtectionSetting.SlowRateThreshold.Rate = slowratethresholdrate

	slowratethresholdperiod, err := tools.GetIntValue("slow_rate_threshold_period", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtectionSetting.SlowRateThreshold.Period = slowratethresholdperiod

	durationthresholdtimeout, err := tools.GetIntValue("duration_threshold_timeout", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSlowPostProtectionSetting.DurationThreshold.Timeout = durationthresholdtimeout

	_, erru := client.UpdateSlowPostProtectionSetting(ctx, updateSlowPostProtectionSetting)
	if erru != nil {
		logger.Errorf("calling 'updateSlowPostProtectionSetting': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSlowPostProtectionSettingRead(ctx, d, m)
}
