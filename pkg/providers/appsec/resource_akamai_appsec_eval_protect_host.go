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
func resourceEvalProtectHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalProtectHostUpdate,
		ReadContext:   resourceEvalProtectHostRead,
		UpdateContext: resourceEvalProtectHostUpdate,
		DeleteContext: resourceEvalProtectHostDelete,
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

func resourceEvalProtectHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalProtectHostRead")

	getEvalProtectHost := appsec.GetEvalProtectHostRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalProtectHost.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalProtectHost.Version = version

	evalprotecthost, err := client.GetEvalProtectHost(ctx, getEvalProtectHost)
	if err != nil {
		logger.Errorf("calling 'getEvalProtectHost': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "evalProtectHostDS", evalprotecthost)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getEvalProtectHost.ConfigID))

	return nil
}

func resourceEvalProtectHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceEvalProtectHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalProtectHostUpdate")

	updateEvalProtectHost := appsec.UpdateEvalProtectHostRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalProtectHost.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalProtectHost.Version = version

	hostnames := d.Get("hostnames").([]interface{})
	hn := make([]string, 0, len(hostnames))

	for _, h := range hostnames {
		hn = append(hn, h.(string))

	}
	updateEvalProtectHost.Hostnames = hn

	_, erru := client.UpdateEvalProtectHost(ctx, updateEvalProtectHost)
	if erru != nil {
		logger.Errorf("calling 'updateEvalProtectHost': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceEvalProtectHostRead(ctx, d, m)
}
