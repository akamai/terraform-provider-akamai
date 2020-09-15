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
			"config_list": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceConfigurationRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read Configuration")

	configuration := appsec.NewConfigurationResponse()
	configurationversion := appsec.NewConfigurationVersionResponse()
	configName := d.Get("name").(string)
	version := d.Get("version").(int)

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

		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CONFIG value  %v\n", configval.ID))
		configurationversion.ConfigID = configval.ID
		err := configurationversion.GetConfigurationVersion(CorrelationID)
		if err != nil {
			edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
			return nil
		}

		//configlist = configlist + " ConfigID/Name=" + strconv.Itoa(configval.ID) + "/" + configval.Name + " VersionList="
		configlist = configlist + "\n" + strconv.Itoa(configval.ID) + " " + configval.Name
		for _, configversionval := range configurationversion.VersionList {
			configlist = configlist + " " + strconv.Itoa(configversionval.Version)
		}
		configlist = configlist + "\n"

		if configval.Name == configName {

			for _, configversionval := range configurationversion.VersionList {
				edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CONFIG Version value  %v\n", configversionval.Version))
				if configversionval.Version == version {
					d.Set("version", strconv.Itoa(version))
				}
				//	configlist = configlist + " " + strconv.Itoa(configversionval.Version)
			}

			d.Set("config_id", configval.ID)
			d.Set("latest_version", configval.LatestVersion)
			configidfound = configval.ID
		}
		d.Set("config_list", configlist)
	}
	d.SetId(strconv.Itoa(configidfound))

	return nil
}
