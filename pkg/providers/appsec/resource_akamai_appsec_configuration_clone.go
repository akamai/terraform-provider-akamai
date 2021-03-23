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
func resourceConfigurationClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationCloneCreate,
		ReadContext:   resourceConfigurationCloneRead,
		UpdateContext: resourceConfigurationCloneUpdate,
		DeleteContext: resourceConfigurationCloneDelete,
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
				Required: true,
			},
			"create_from_version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"config_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Config id  of cloned configuration",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of cloned configuration",
			},
		},
	}
}

func resourceConfigurationCloneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationCloneCreate")

	createConfigurationClone := appsec.CreateConfigurationCloneRequest{}

	createFromConfigID, err := tools.GetIntValue("create_from_config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.CreateFrom.ConfigID = createFromConfigID

	version, err := tools.GetIntValue("create_from_version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.CreateFrom.Version = version

	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.Name = name

	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.Description = description

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.ContractID = contractID

	groupID, err := tools.GetIntValue("group_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.GroupID = groupID

	hostnamelist := d.Get("host_names").(*schema.Set)
	hnl := make([]string, 0, len(hostnamelist.List()))

	for _, h := range hostnamelist.List() {
		hnl = append(hnl, h.(string))

	}
	createConfigurationClone.Hostnames = hnl

	ccr, err := client.CreateConfigurationClone(ctx, createConfigurationClone)
	if err != nil {
		logger.Errorf("calling 'createConfigurationClone': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("version", ccr.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", ccr.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(ccr.ConfigID))

	return resourceConfigurationCloneRead(ctx, d, m)
}

func resourceConfigurationCloneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationCloneRead")

	getConfigurationClone := appsec.GetConfigurationCloneRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getConfigurationClone.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getConfigurationClone.Version = version

	Configurationclone, err := client.GetConfigurationClone(ctx, getConfigurationClone)
	if err != nil {
		logger.Errorf("calling 'getConfigurationClone': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(Configurationclone.ConfigID))

	return nil
}

func resourceConfigurationCloneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRemove")

	removeConfiguration := appsec.RemoveConfigurationRequest{}

	configid, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
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

func resourceConfigurationCloneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceConfigurationCloneRead(ctx, d, m)
}
