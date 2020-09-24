package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigurationVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConfigurationVersionRead,
		Schema: map[string]*schema.Schema{
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"latest_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"staging_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"production_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceConfigurationVersionRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceConfigurationVersionRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read ConfigurationVersion")

	configurationversion := appsec.NewConfigurationVersionResponse()
	configurationversion.ConfigID = d.Get("config_id").(int)
	configVersion := d.Get("version").(int)

	err := configurationversion.GetConfigurationVersion(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ConfigurationVersion   %v\n", configurationversion))
	d.Set("latest_version", configurationversion.LastCreatedVersion)

	for _, configval := range configurationversion.VersionList {

		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CONFIG value  %v\n", configval.Version))
		if configval.Version == configVersion {
			d.Set("config_id", configval.ConfigID)
			d.Set("staging_status", configval.Staging.Status)
			d.Set("production_status", configval.Production.Status)
			d.SetId(strconv.Itoa(configval.ConfigID))
		}
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "configurationVersion", configurationversion)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ConfigurationVesion outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}
	d.SetId(strconv.Itoa(configurationversion.ConfigID))

	return nil
}
