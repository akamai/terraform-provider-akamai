package appsec

import (
	"context"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
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

	createConfigurationClone := v2.CreateConfigurationCloneRequest{}

	createConfigurationClone.ConfigID = d.Get("config_id").(int)
	createConfigurationClone.CreateFromVersion = d.Get("create_from_version").(int)

	ccr, err := client.CreateConfigurationClone(ctx, createConfigurationClone)
	if err != nil {
		logger.Warnf("calling 'createConfigurationClone': %s", err.Error())
	}
	logger.Warnf("calling 'createConfigurationClone CCR ': %v", ccr)
	//	d.Set("version", ccr.Version)
	d.SetId(strconv.Itoa(ccr.ConfigID))

	return resourceConfigurationCloneRead(ctx, d, m)
}

func resourceConfigurationCloneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationCloneRead")

	getConfigurationClone := v2.GetConfigurationCloneRequest{}

	getConfigurationClone.ConfigID = d.Get("config_id").(int)
	getConfigurationClone.Version = d.Get("create_from_version").(int)

	configurationclone, err := client.GetConfigurationClone(ctx, getConfigurationClone)
	if err != nil {
		logger.Warnf("calling 'getConfigurationClone': %s", err.Error())
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
