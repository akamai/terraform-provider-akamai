package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfiguration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConfigurationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"latest_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"staging_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"production_version": {
				Type:     schema.TypeInt,
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

func dataSourceConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceConfigurationRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read Configuration")

	configuration := appsec.NewConfigurationResponse()
	configName := d.Get("name").(string)

	err := configuration.GetConfiguration(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Configuration   %v\n", configuration))

	var configlist string
	var configidfound int
	configlist = configlist + " ConfigID Name  VersionList" + "\n"

	for _, configval := range configuration.Configurations {

		if configval.Name == configName {
			d.Set("config_id", configval.ID)
			d.Set("latest_version", configval.LatestVersion)
			d.Set("staging_version", configval.StagingVersion)
			d.Set("production_version", configval.ProductionVersion)
			configidfound = configval.ID
		}
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "configuration", configuration)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Configuration outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(configidfound))

	return nil
}
