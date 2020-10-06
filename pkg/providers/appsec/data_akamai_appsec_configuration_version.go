package appsec

import (
	"context"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigurationVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationVersionRead,
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

func dataSourceConfigurationVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationVersionRead")

	getConfigurationVersion := v2.GetConfigurationVersionsRequest{}

	getConfigurationVersion.ConfigID = d.Get("config_id").(int)
	configVersion := d.Get("version").(int)

	configurationversion, err := client.GetConfigurationVersions(ctx, getConfigurationVersion)
	if err != nil {
		logger.Warnf("calling 'getConfigurationVersion': %s", err.Error())
	}

	d.Set("latest_version", configurationversion.LastCreatedVersion)

	for _, configval := range configurationversion.VersionList {

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
	if err == nil {
		d.Set("output_text", outputtext)
	}
	d.SetId(strconv.Itoa(configurationversion.ConfigID))

	return nil
}
