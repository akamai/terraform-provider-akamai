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
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"create_from_version": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"rule_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.CreateFrom.ConfigID = configid

	version, err := tools.GetIntValue("create_from_version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createConfigurationClone.CreateFrom.Version = version

	ccr, err := client.CreateConfigurationClone(ctx, createConfigurationClone)
	if err != nil {
		logger.Errorf("calling 'createConfigurationClone': %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Errorf("calling 'createConfigurationClone CCR ': %v", ccr)

	if err := d.Set("version", ccr.Version); err != nil {
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

	configurationclone, err := client.GetConfigurationClone(ctx, getConfigurationClone)
	if err != nil {
		logger.Errorf("calling 'getConfigurationClone': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configurationclone.ConfigID))

	return nil
}

func resourceConfigurationCloneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}

func resourceConfigurationCloneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}
