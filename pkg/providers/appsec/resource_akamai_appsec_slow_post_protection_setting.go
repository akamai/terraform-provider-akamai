package appsec

import (
	"context"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
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

	getSlowPostProtectionSetting.ConfigID = d.Get("config_id").(int)
	getSlowPostProtectionSetting.Version = d.Get("version").(int)
	getSlowPostProtectionSetting.PolicyID = d.Get("policy_id").(string)

	_, err := client.GetSlowPostProtectionSetting(ctx, getSlowPostProtectionSetting)
	if err != nil {
		logger.Warnf("calling 'getSlowPostProtectionSetting': %s", err.Error())
	}

	d.SetId(strconv.Itoa(getSlowPostProtectionSetting.ConfigID))

	return nil
}

func resourceSlowPostProtectionSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//meta := akamai.Meta(m)
	//client := inst.Client(meta)
	//logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingRemove")

	return schema.NoopContext(nil, d, m)
}

func resourceSlowPostProtectionSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSlowPostProtectionSettingUpdate")

	updateSlowPostProtectionSetting := v2.UpdateSlowPostProtectionSettingRequest{}

	//slowpostprotectionsettingspost := appsec.NewSlowPostProtectionSettingsPost()

	updateSlowPostProtectionSetting.ConfigID = d.Get("config_id").(int)
	updateSlowPostProtectionSetting.Version = d.Get("version").(int)
	updateSlowPostProtectionSetting.PolicyID = d.Get("policy_id").(string)
	updateSlowPostProtectionSetting.Action = d.Get("slow_rate_action").(string)
	updateSlowPostProtectionSetting.SlowRateThreshold.Rate = d.Get("slow_rate_threshold_rate").(int)
	updateSlowPostProtectionSetting.SlowRateThreshold.Period = d.Get("slow_rate_threshold_period").(int)
	updateSlowPostProtectionSetting.DurationThreshold.Timeout = d.Get("duration_threshold_timeout").(int)

	_, err := client.UpdateSlowPostProtectionSetting(ctx, updateSlowPostProtectionSetting)
	if err != nil {
		logger.Warnf("calling 'updateSlowPostProtectionSetting': %s", err.Error())
	}

	return resourceSlowPostProtectionSettingRead(ctx, d, m)
}
