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
func resourceConfigurationVersionClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationVersionCloneCreate,
		ReadContext:   resourceConfigurationVersionCloneRead,
		UpdateContext: resourceConfigurationVersionCloneUpdate,
		DeleteContext: resourceConfigurationVersionCloneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceConfigurationVersionCloneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationVersionCloneCreate")

	createConfigurationVersionClone := appsec.CreateConfigurationVersionCloneRequest{}

	createConfigurationVersionClone.ConfigID = d.Get("config_id").(int)
	createConfigurationVersionClone.CreateFromVersion = d.Get("create_from_version").(int)

	ccr, err := client.CreateConfigurationVersionClone(ctx, createConfigurationVersionClone)
	if err != nil {
		logger.Errorf("calling 'createConfigurationVersionClone': %s", err.Error())
		return diag.FromErr(err)
	}

	d.Set("version", ccr.Version)
	d.SetId(strconv.Itoa(ccr.Version))

	return resourceConfigurationVersionCloneRead(ctx, d, m)
}

func resourceConfigurationVersionCloneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationVersionCloneRead")

	getConfigurationVersionClone := appsec.GetConfigurationVersionCloneRequest{}

	getConfigurationVersionClone.ConfigID = d.Get("config_id").(int)
	//getConfigurationVersionClone.Version = d.Get("create_from_version").(int)
	version, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getConfigurationVersionClone.Version = version

	_, err := client.GetConfigurationVersionClone(ctx, getConfigurationVersionClone)
	if err != nil {
		logger.Errorf("calling 'getConfigurationVersionClone': %s", err.Error())
		return diag.FromErr(err)
	}

	//d.SetId(strconv.Itoa(configurationversionclone.ConfigID))
	d.SetId(strconv.Itoa(version))
	return nil
}

func resourceConfigurationVersionCloneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceConfigurationVersionCloneRead")

	removeConfigurationVersionClone := appsec.RemoveConfigurationVersionCloneRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeConfigurationVersionClone.ConfigID = configid

	version, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	removeConfigurationVersionClone.Version = version

	_, errd := client.RemoveConfigurationVersionClone(ctx, removeConfigurationVersionClone)
	if errd != nil {
		logger.Errorf("calling 'getConfigurationVersionClone': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")
	return nil
}

func resourceConfigurationVersionCloneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}
