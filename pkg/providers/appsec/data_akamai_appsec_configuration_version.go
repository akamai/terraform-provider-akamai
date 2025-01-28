package appsec

import (
	"context"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfigurationVersion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationVersionRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Version of the security configuration for which to return information",
			},
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Latest version of the security configuration",
			},
			"staging_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the specified version in staging",
			},
			"production_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the specified version in production",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceConfigurationVersionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceConfigurationVersionRead")

	getConfigurationVersion := appsec.GetConfigurationVersionsRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getConfigurationVersion.ConfigID = configID

	version, err := tf.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getConfigurationVersion.ConfigVersion = version

	configurationversion, err := client.GetConfigurationVersions(ctx, getConfigurationVersion)
	if err != nil {
		logger.Errorf("calling 'getConfigurationVersion': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("latest_version", configurationversion.LastCreatedVersion); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	for _, configval := range configurationversion.VersionList {

		if configval.Version == version {

			if err := d.Set("config_id", configval.ConfigID); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}

			if err := d.Set("staging_status", configval.Staging.Status); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}

			if err := d.Set("production_status", configval.Production.Status); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}

			d.SetId(strconv.Itoa(configval.ConfigID))
		}
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "configurationVersion", configurationversion)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	d.SetId(strconv.Itoa(configurationversion.ConfigID))

	return nil
}
