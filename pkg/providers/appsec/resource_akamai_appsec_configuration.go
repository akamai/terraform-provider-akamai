package appsec

import (
	"context"
	"errors"
	"fmt"
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
func resourceConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationCreate,
		ReadContext:   resourceConfigurationRead,
		UpdateContext: resourceConfigurationUpdate,
		DeleteContext: resourceConfigurationDelete,
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
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationCreate")

	createConfiguration := appsec.CreateConfigurationRequest{}

	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfiguration.Name = name

	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfiguration.Description = description

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfiguration.ContractID = contractID

	groupID, err := tools.GetIntValue("group_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfiguration.GroupID = groupID

	hostnamelist := d.Get("host_names").(*schema.Set)
	hnl := make([]string, 0, len(hostnamelist.List()))

	for _, h := range hostnamelist.List() {
		hnl = append(hnl, h.(string))

	}
	createConfiguration.Hostnames = hnl

	postresp, errc := client.CreateConfiguration(ctx, createConfiguration)
	if errc != nil {
		logger.Errorf("calling 'createConfiguration': %s", errc.Error())
		return diag.FromErr(errc)
	}

	if err := d.Set("config_id", postresp.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", postresp.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(postresp.ConfigID))

	return resourceConfigurationRead(ctx, d, m)
}

func resourceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationUpdate")

	updateConfiguration := appsec.UpdateConfigurationRequest{}

	ID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	updateConfiguration.ConfigID = ID

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

	return resourceConfigurationRead(ctx, d, m)
}

func resourceConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRemove")

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

func resourceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRead")

	getConfiguration := appsec.GetConfigurationsRequest{}

	ID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getConfiguration.ConfigID = ID

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
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(configidfound))

	return nil
}
