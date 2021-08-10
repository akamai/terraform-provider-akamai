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
			StateContext: schema.ImportStatePassthroughContext,
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
			"create_from_config_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"create_from_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_id": {
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
	logger.Debug("in resourceConfigurationCreate")

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tools.GetIntValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnameset, err := tools.GetSetValue("host_names", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames := make([]string, 0, len(hostnameset.List()))
	for _, h := range hostnameset.List() {
		hostnames = append(hostnames, h.(string))
	}
	createFromConfigID, err := tools.GetIntValue("create_from_config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createFromVersion, err := tools.GetIntValue("create_from_version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if createFromVersion > 0 && createFromConfigID > 0 {
		createConfigurationClone := appsec.CreateConfigurationCloneRequest{}
		createConfigurationClone.CreateFrom.ConfigID = createFromConfigID
		createConfigurationClone.CreateFrom.Version = createFromVersion
		createConfigurationClone.Name = name
		createConfigurationClone.Description = description
		createConfigurationClone.ContractID = contractID
		createConfigurationClone.GroupID = groupID
		createConfigurationClone.Hostnames = hostnames

		response, err := client.CreateConfigurationClone(ctx, createConfigurationClone)
		if err != nil {
			logger.Errorf("calling 'createConfigurationClone': %s", err.Error())
			return diag.FromErr(err)
		}
		if err := d.Set("config_id", response.ConfigID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		d.SetId(fmt.Sprintf("%d", response.ConfigID))

	} else {
		createConfiguration := appsec.CreateConfigurationRequest{}
		createConfiguration.Name = name
		createConfiguration.Description = description
		createConfiguration.ContractID = contractID
		createConfiguration.GroupID = groupID
		createConfiguration.Hostnames = hostnames

		response, err := client.CreateConfiguration(ctx, createConfiguration)
		if err != nil {
			logger.Errorf("calling 'createConfiguration': %s", err.Error())
			return diag.FromErr(err)
		}
		if err := d.Set("config_id", response.ConfigID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		d.SetId(fmt.Sprintf("%d", response.ConfigID))
	}

	return resourceConfigurationRead(ctx, d, m)
}

func resourceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRead")
	logger.Debug("in resourceConfigurationRead")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getConfiguration := appsec.GetConfigurationRequest{}
	getConfiguration.ConfigID = configid

	configuration, err := client.GetConfiguration(ctx, getConfiguration)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	d.Set("name", configuration.Name)
	d.Set("description", configuration.Description)
	d.Set("config_id", configuration.ID)

	getSelectedHostname := appsec.GetSelectedHostnameRequest{}
	getSelectedHostname.ConfigID = configid
	getSelectedHostname.Version = getLatestConfigVersion(ctx, configid, m)

	selectedhostnames, err := client.GetSelectedHostname(ctx, getSelectedHostname)
	if err != nil {
		logger.Errorf("calling 'getSelectedHostname': %s", err.Error())
		return diag.FromErr(err)
	}
	selectedhostnameset := schema.Set{F: schema.HashString}
	for _, hostname := range selectedhostnames.HostnameList {
		selectedhostnameset.Add(hostname.Hostname)
	}

	d.Set("host_names", selectedhostnameset.List())

	return nil
}

func resourceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationUpdate")
	logger.Debug("in resourceConfigurationUpdate")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateConfiguration := appsec.UpdateConfigurationRequest{}
	updateConfiguration.ConfigID = configid
	updateConfiguration.Name = name
	updateConfiguration.Description = description

	resp, erru := client.UpdateConfiguration(ctx, updateConfiguration)
	if erru != nil {
		logger.Errorf("calling 'updateConfiguration': %s", erru.Error())
		logger.Debugf("response is %w", resp)
		return diag.FromErr(erru)
	}

	if d.HasChange("host_names") {
		hostnameset, err := tools.GetSetValue("host_names", d)
		if err != nil {
			return diag.FromErr(err)
		}
		hostnamelist := tools.SetToStringSlice(hostnameset)
		hostnames := make([]appsec.Hostname, 0, len(hostnamelist))
		for _, name := range hostnamelist {
			hostname := appsec.Hostname{Hostname: name}
			hostnames = append(hostnames, hostname)
		}

		updateSelectedHostname := appsec.UpdateSelectedHostnameRequest{}
		updateSelectedHostname.ConfigID = configid
		updateSelectedHostname.Version = getModifiableConfigVersion(ctx, configid, "configuration", m)
		updateSelectedHostname.HostnameList = hostnames

		_, err = client.UpdateSelectedHostname(ctx, updateSelectedHostname)
		if err != nil {
			logger.Errorf("calling 'UpdateSelectedHostname': %s", err.Error())
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return resourceConfigurationRead(ctx, d, m)
}

func resourceConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRemove")
	logger.Debug("in resourceConfigurationDelete")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Check whether any versions of this config have ever been activated
	getConfigVersionsRequest := appsec.GetConfigurationVersionsRequest{}
	getConfigVersionsRequest.ConfigID = configid
	configurationVersions, err := client.GetConfigurationVersions(ctx, getConfigVersionsRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, configVersion := range configurationVersions.VersionList {
		if configVersion.Production.Status != "Inactive" || configVersion.Staging.Status != "Inactive" {
			return diag.Errorf("cannot delete configuration '%s' as version %d has been active in staging or production",
				configurationVersions.ConfigName, configVersion.Version)
		}
	}

	removeConfiguration := appsec.RemoveConfigurationRequest{}
	removeConfiguration.ConfigID = configid

	_, errd := client.RemoveConfiguration(ctx, removeConfiguration)
	if errd != nil {
		logger.Errorf("calling 'removeConfiguration': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")
	return nil
}
