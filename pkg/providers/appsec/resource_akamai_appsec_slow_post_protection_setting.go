package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
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

	getSlowPostProtectionSetting := appsec.GetSlowPostProtectionSettingRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getSlowPostProtectionSetting.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getSlowPostProtectionSetting.Version = version

		policyid := s[2]
		getSlowPostProtectionSetting.PolicyID = policyid

	} else {
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

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSlowPostProtectionSetting.PolicyID = policyid
	}
	getslowpost, errg := client.GetSlowPostProtectionSetting(ctx, getSlowPostProtectionSetting)
	if errg != nil {
		logger.Errorf("calling 'getSlowPostProtectionSetting': %s", errg.Error())
		return diag.FromErr(errg)
	}

	if err := d.Set("config_id", getSlowPostProtectionSetting.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getSlowPostProtectionSetting.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getSlowPostProtectionSetting.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("slow_rate_action", getslowpost.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("slow_rate_threshold_rate", getslowpost.SlowRateThreshold.Rate); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("slow_rate_threshold_period", getslowpost.SlowRateThreshold.Period); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("duration_threshold_timeout", getslowpost.DurationThreshold.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", getSlowPostProtectionSetting.ConfigID, getSlowPostProtectionSetting.Version, getSlowPostProtectionSetting.PolicyID))

	return nil
}

func resourceSlowPostProtectionSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingDelete")

	updateSlowPostProtection := appsec.UpdateSlowPostProtectionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSlowPostProtection.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSlowPostProtection.Version = version

		policyid := s[2]
		updateSlowPostProtection.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSlowPostProtection.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSlowPostProtection.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSlowPostProtection.PolicyID = policyid
	}
	updateSlowPostProtection.ApplySlowPostControls = false

	_, erru := client.UpdateSlowPostProtection(ctx, updateSlowPostProtection)
	if erru != nil {
		logger.Errorf("calling 'updateSlowPostProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}
	d.SetId("")
	return nil
}

func resourceSlowPostProtectionSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingUpdate")

	updateSlowPostProtectionSetting := appsec.UpdateSlowPostProtectionSettingRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSlowPostProtectionSetting.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSlowPostProtectionSetting.Version = version

		policyid := s[2]
		updateSlowPostProtectionSetting.PolicyID = policyid

	} else {
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

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSlowPostProtectionSetting.PolicyID = policyid
	}
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

const (
	Abort = "abort"
)
