package appsec

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of a specific security information for which to retrieve information",
			},
			"config_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier of the security configuration",
			},
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Latest version of the security configuration",
			},
			"staging_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of the security configuration currently deployed in staging",
			},
			"production_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of the security configuration currently deployed in production",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRead")

	getConfiguration := appsec.GetConfigurationsRequest{}

	configName := d.Get("name").(string)

	configuration, err := client.GetConfigurations(ctx, getConfiguration)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	var configID int

	for _, configval := range configuration.Configurations {

		if configval.Name == configName {
			if err := d.Set("config_id", configval.ID); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}

			if err := d.Set("latest_version", configval.LatestVersion); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}

			if err := d.Set("staging_version", configval.StagingVersion); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}

			if err := d.Set("production_version", configval.ProductionVersion); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
			configID = configval.ID
		}
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "configuration", configuration)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	d.SetId(strconv.Itoa(configID))

	return nil
}
