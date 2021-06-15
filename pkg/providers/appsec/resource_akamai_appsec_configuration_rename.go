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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceConfigurationRename() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationRenameCreate,
		ReadContext:   resourceConfigurationRenameRead,
		UpdateContext: resourceConfigurationRenameUpdate,
		DeleteContext: resourceConfigurationRenameDelete,
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConfigurationRenameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameUpdate")
	logger.Debugf("!!! in resourceConfigurationRenameCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateConfiguration := appsec.UpdateConfigurationRequest{}
	updateConfiguration.ConfigID = configid
	updateConfiguration.Name = name
	updateConfiguration.Description = description

	_, erru := client.UpdateConfiguration(ctx, updateConfiguration)
	if erru != nil {
		logger.Errorf("calling 'updateConfiguration': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(strconv.Itoa(updateConfiguration.ConfigID))

	return resourceConfigurationRenameRead(ctx, d, m)
}

func resourceConfigurationRenameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameRead")
	logger.Debugf("!!! in resourceConfigurationRenameRead")

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

	if err := d.Set("config_id", getConfiguration.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("name", configuration.Name); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("description", configuration.Description); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceConfigurationRenameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameUpdate")
	logger.Debugf("!!! in resourceConfigurationRenameRead")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateConfiguration := appsec.UpdateConfigurationRequest{}
	updateConfiguration.ConfigID = configid
	updateConfiguration.Name = name
	updateConfiguration.Description = description

	_, erru := client.UpdateConfiguration(ctx, updateConfiguration)
	if erru != nil {
		logger.Errorf("calling 'updateConfiguration': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceConfigurationRenameRead(ctx, d, m)
}

func resourceConfigurationRenameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(context.TODO(), d, m)
}
