package appsec

import (
	"context"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
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
			"host_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Hostnames to be protected by the new configuration",
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
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRead")

	configName, err := tf.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	configurations, err := client.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	outputConfigurations := appsec.GetConfigurationsResponse{}
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	if configName != "" {
		found := false
		for _, config := range configurations.Configurations {
			if config.Name == configName {
				found = true
				outputConfigurations.Configurations = append(outputConfigurations.Configurations, config)
				if err := d.Set("config_id", config.ID); err != nil {
					return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
				}
				if err := d.Set("latest_version", config.LatestVersion); err != nil {
					return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
				}
				if err := d.Set("staging_version", config.StagingVersion); err != nil {
					return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
				}
				if err := d.Set("production_version", config.ProductionVersion); err != nil {
					return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
				}
				getSelectedHostnamesRequest := appsec.GetSelectedHostnamesRequest{
					ConfigID: config.ID,
					Version:  config.LatestVersion,
				}

				// Fetch selected hostnames for the config version
				selectedHostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnamesRequest)
				if err != nil {
					logger.Errorf("calling 'GetSelectedHostnames': %s", err.Error())
					return diag.FromErr(err)
				}
				selectedHostnameList := make([]string, 0, len(selectedHostnames.HostnameList))
				for _, hostname := range selectedHostnames.HostnameList {
					selectedHostnameList = append(selectedHostnameList, hostname.Hostname)
				}

				// Set host_names for the config version
				if err = d.Set("host_names", selectedHostnameList); err != nil {
					return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
				}
				d.SetId(strconv.Itoa(config.ID))
				break
			}
		}
		if !found {
			return diag.Errorf("configuration '%s' not found", configName)
		}
	} else {
		if len(configurations.Configurations) > 0 {
			outputConfigurations = *configurations
			d.SetId(strconv.Itoa(configurations.Configurations[0].ID))
		} else {
			d.SetId(strconv.Itoa(0))
		}
	}

	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "configuration", outputConfigurations)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}
