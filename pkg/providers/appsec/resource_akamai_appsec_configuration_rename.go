package appsec

import (
	"context"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceConfigurationRename() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationRenameUpdate,
		ReadContext:   resourceConfigurationRenameRead,
		UpdateContext: resourceConfigurationRenameUpdate,
		DeleteContext: resourceConfigurationRenameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceConfigurationRenameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameUpdate")

	updateConfiguration := appsec.UpdateConfigurationRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateConfiguration.ConfigID = configid

	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateConfiguration.Name = name

	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateConfiguration.Description = description

	_, erru := client.UpdateConfiguration(ctx, updateConfiguration)
	if erru != nil {
		logger.Errorf("calling 'updateConfiguration': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceConfigurationRenameRead(ctx, d, m)
}

func resourceConfigurationRenameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameRemove")

	removeConfiguration := appsec.RemoveConfigurationRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeConfiguration.ConfigID = configid

	_, errd := client.RemoveConfiguration(ctx, removeConfiguration)
	if errd != nil {
		logger.Errorf("calling 'removeConfiguration': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceConfigurationRenameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameRead")

	getConfiguration := appsec.GetConfigurationsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getConfiguration.ConfigID = configid

	configName, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getConfiguration.Name = configName

	configuration, err := client.GetConfigurations(ctx, getConfiguration)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	var configlist string
	var configidfound int
	configlist = configlist + " ConfigID Name  VersionList" + "\n"

	for _, configval := range configuration.Configurations {

		if configval.ID == configid {
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
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(configidfound))

	return nil
}
