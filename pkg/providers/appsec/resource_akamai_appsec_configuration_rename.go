package appsec

import (
	"context"
	"errors"
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
// https://techdocs.akamai.com/application-security/reference/api
func resourceConfigurationRename() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationRenameCreate,
		ReadContext:   resourceConfigurationRenameRead,
		UpdateContext: resourceConfigurationRenameUpdate,
		DeleteContext: resourceConfigurationRenameDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
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
	logger := meta.Log("APPSEC", "resourceConfigurationRenameCreate")
	logger.Debugf("in resourceConfigurationRenameCreate")

	configID, err := tools.GetIntValue("config_id", d)
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

	updateConfiguration := appsec.UpdateConfigurationRequest{
		ConfigID:    configID,
		Name:        name,
		Description: description,
	}

	_, err = client.UpdateConfiguration(ctx, updateConfiguration)
	if err != nil {
		logger.Errorf("calling 'updateConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(updateConfiguration.ConfigID))

	return resourceConfigurationRenameRead(ctx, d, m)
}

func resourceConfigurationRenameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameRead")
	logger.Debugf("in resourceConfigurationRenameRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getConfiguration := appsec.GetConfigurationRequest{
		ConfigID: configID,
	}

	configuration, err := client.GetConfiguration(ctx, getConfiguration)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getConfiguration.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("name", configuration.Name); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("description", configuration.Description); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceConfigurationRenameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationRenameUpdate")
	logger.Debugf("in resourceConfigurationRenameUpdate")

	configID, err := strconv.Atoi(d.Id())
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

	updateConfiguration := appsec.UpdateConfigurationRequest{
		ConfigID:    configID,
		Name:        name,
		Description: description,
	}

	_, err = client.UpdateConfiguration(ctx, updateConfiguration)
	if err != nil {
		logger.Errorf("calling 'updateConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceConfigurationRenameRead(ctx, d, m)
}

func resourceConfigurationRenameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(ctx, d, m)
}
