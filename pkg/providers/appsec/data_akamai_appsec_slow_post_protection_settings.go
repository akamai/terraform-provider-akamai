package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSlowPostProtectionSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSlowPostProtectionSettingsRead,
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
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceSlowPostProtectionSettingsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceSlowPostProtectionSettingsRead-" + tools.CreateNonce() + "]"

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

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("SlowPostProtectionSettings   %v\n", slowpostprotectionsettings))

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "slowPostDS", slowpostprotectionsettings)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("slowPost outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(configid))

	return nil
}
