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
func resourceSlowPostProtectionSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceSlowPostProtectionSettingsUpdate,
		Read:   resourceSlowPostProtectionSettingsRead,
		Update: resourceSlowPostProtectionSettingsUpdate,
		Delete: resourceSlowPostProtectionSettingsDelete,
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

func resourceSlowPostProtectionSettingsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSlowPostProtectionSettingsRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SlowPostProtectionSettings")

	slowpostprotectionsettings := appsec.NewSlowPostProtectionSettingsResponse()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)

	err := slowpostprotectionsettings.GetSlowPostProtectionSettings(configid, version, policyid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.SetId(strconv.Itoa(configid))

	return nil
}

func resourceSlowPostProtectionSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSlowPostProtectionSettingsDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting SlowPostProtectionSettings")

	return schema.Noop(d, meta)
}

func resourceSlowPostProtectionSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceSlowPostProtectionSettingsUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating SlowPostProtectionSettings")

	slowpostprotectionsettings := appsec.NewSlowPostProtectionSettingsResponse()

	slowpostprotectionsettingspost := appsec.NewSlowPostProtectionSettingsPost()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	slowpostprotectionsettingspost.Action = d.Get("slow_rate_action").(string)
	slowpostprotectionsettingspost.SlowRateThreshold.Rate = d.Get("slow_rate_threshold_rate").(int)
	slowpostprotectionsettingspost.SlowRateThreshold.Period = d.Get("slow_rate_threshold_period").(int)
	slowpostprotectionsettingspost.DurationThreshold.Timeout = d.Get("duration_threshold_timeout").(int)

	err := slowpostprotectionsettings.UpdateSlowPostProtectionSettings(configid, version, policyid, slowpostprotectionsettingspost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	return resourceSlowPostProtectionSettingsRead(d, meta)

}
