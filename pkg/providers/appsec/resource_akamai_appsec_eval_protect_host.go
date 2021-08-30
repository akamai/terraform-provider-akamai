package appsec

import (
	"context"
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
func resourceEvalProtectHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalProtectHostCreate,
		ReadContext:   resourceEvalProtectHostRead,
		UpdateContext: resourceEvalProtectHostUpdate,
		DeleteContext: resourceEvalProtectHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hostnames": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceEvalProtectHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalProtectHostCreate")
	logger.Debug("in resourceEvalProtectHostCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnameset, err := tools.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateEvalProtectHost := appsec.UpdateEvalProtectHostRequest{}
	updateEvalProtectHost.ConfigID = configid
	updateEvalProtectHost.Version = getModifiableConfigVersion(ctx, configid, "evalprotecthostnames", m)
	hostnamelist := make([]string, 0, len(hostnameset.List()))
	for _, hostname := range hostnameset.List() {
		hostnamelist = append(hostnamelist, hostname.(string))
	}
	updateEvalProtectHost.Hostnames = hostnamelist

	_, err = client.UpdateEvalProtectHost(ctx, updateEvalProtectHost)
	if err != nil {
		logger.Errorf("calling 'updateEvalProtectHost': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", configid))

	return resourceEvalProtectHostRead(ctx, d, m)
}

func resourceEvalProtectHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalProtectHostRead")
	logger.Debug("in resourceEvalProtectHostRead")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getEvalProtectHostsRequest := appsec.GetEvalProtectHostsRequest{
		ConfigID: configid,
		Version:  getLatestConfigVersion(ctx, configid, m),
	}

	evalprotecthostnames, err := client.GetEvalProtectHosts(ctx, getEvalProtectHostsRequest)
	if err != nil {
		logger.Errorf("calling 'updateEvalProtectHost': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	evalprotecthostnameset := schema.Set{F: schema.HashString}
	for _, hostname := range evalprotecthostnames.Hostnames {
		evalprotecthostnameset.Add(hostname)
	}
	if err := d.Set("hostnames", evalprotecthostnameset.List()); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceEvalProtectHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalProtectHostUpdate")
	logger.Debug("in resourceEvalProtectHostUpdate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames, err := tools.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateEvalProtectHost := appsec.UpdateEvalProtectHostRequest{}
	updateEvalProtectHost.ConfigID = configid
	updateEvalProtectHost.Version = getModifiableConfigVersion(ctx, configid, "evalprotecthostnames", m)
	hostnamelist := make([]string, 0, len(hostnames.List()))
	for _, hostname := range hostnames.List() {
		hostnamelist = append(hostnamelist, hostname.(string))
	}
	updateEvalProtectHost.Hostnames = hostnamelist

	_, erru := client.UpdateEvalProtectHost(ctx, updateEvalProtectHost)
	if erru != nil {
		logger.Errorf("calling 'updateEvalProtectHost': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceEvalProtectHostRead(ctx, d, m)
}

func resourceEvalProtectHostDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("APPSEC", "resourceEvalProtectHostDelete")
	logger.Debug("in resourceEvalProtectHostDelete")

	return schema.NoopContext(context.TODO(), d, m)
}
