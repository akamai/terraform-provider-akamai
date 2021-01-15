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
func resourceEvalHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalHostUpdate,
		ReadContext:   resourceEvalHostRead,
		UpdateContext: resourceEvalHostUpdate,
		DeleteContext: resourceEvalHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hostnames": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceEvalHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalHostRead")

	getEvalHost := appsec.GetEvalHostRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalHost.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalHost.Version = version

	evalhost, err := client.GetEvalHost(ctx, getEvalHost)
	if err != nil {
		logger.Errorf("calling 'getEvalHost': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "evalHostDS", evalhost)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getEvalHost.ConfigID))

	return nil
}

func resourceEvalHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalHostRemove")

	removeEvalHost := appsec.RemoveEvalHostRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeEvalHost.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeEvalHost.Version = version

	removeEvalHost.Hostnames = nil

	_, erru := client.RemoveEvalHost(ctx, removeEvalHost)
	if erru != nil {
		logger.Errorf("calling 'updateEvalHost': %s", erru.Error())
		return diag.FromErr(erru)
	}
	d.SetId("")
	return nil
}

func resourceEvalHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalHostUpdate")

	updateEvalHost := appsec.UpdateEvalHostRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalHost.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalHost.Version = version

	hostnames := d.Get("hostnames").([]interface{})
	hn := make([]string, 0, len(hostnames))

	for _, h := range hostnames {
		hn = append(hn, h.(string))

	}
	updateEvalHost.Hostnames = hn

	_, erru := client.UpdateEvalHost(ctx, updateEvalHost)
	if erru != nil {
		logger.Errorf("calling 'updateEvalHost': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceEvalHostRead(ctx, d, m)
}
