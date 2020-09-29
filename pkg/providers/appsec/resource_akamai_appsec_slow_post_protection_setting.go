package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSlowPostProtectionSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceSlowPostProtectionSettingUpdate,
		Read:   resourceSlowPostProtectionSettingRead,
		Update: resourceSlowPostProtectionSettingUpdate,
		Delete: resourceSlowPostProtectionSettingDelete,
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

func resourceSlowPostProtectionSettingRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSlowPostProtectionSettingRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SlowPostProtectionSetting")

	slowpostprotectionsetting := appsec.NewSlowPostProtectionSettingResponse()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)

	err := slowpostprotectionsetting.GetSlowPostProtectionSetting(configid, version, policyid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.SetId(strconv.Itoa(configid))

	return nil
}

func resourceSlowPostProtectionSettingDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSlowPostProtectionSettingDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting SlowPostProtectionSetting")

	return schema.Noop(d, meta)
}

func resourceSlowPostProtectionSettingUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSlowPostProtectionSettingUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating SlowPostProtectionSetting")

	slowpostprotectionsetting := appsec.NewSlowPostProtectionSettingResponse()

	slowpostprotectionsettingspost := appsec.NewSlowPostProtectionSettingsPost()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	slowpostprotectionsettingspost.Action = d.Get("slow_rate_action").(string)
	slowpostprotectionsettingspost.SlowRateThreshold.Rate = d.Get("slow_rate_threshold_rate").(int)
	slowpostprotectionsettingspost.SlowRateThreshold.Period = d.Get("slow_rate_threshold_period").(int)
	slowpostprotectionsettingspost.DurationThreshold.Timeout = d.Get("duration_threshold_timeout").(int)

	err := slowpostprotectionsetting.UpdateSlowPostProtectionSetting(configid, version, policyid, slowpostprotectionsettingspost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	return resourceSlowPostProtectionSettingRead(d, meta)

}
