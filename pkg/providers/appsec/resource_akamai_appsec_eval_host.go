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
func resourceEvalHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalHostCreate,
		ReadContext:   resourceEvalHostRead,
		UpdateContext: resourceEvalHostUpdate,
		DeleteContext: resourceEvalHostDelete,
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

func resourceEvalHostCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalHostCreate")
	logger.Debug("in resourceEvalHostCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnameset, err := tools.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}

	hostnamelist := make([]string, 0, len(hostnameset.List()))
	for _, hostname := range hostnameset.List() {
		hostnamelist = append(hostnamelist, hostname.(string))
	}

	updateEvalHost := appsec.UpdateEvalHostRequest{
		ConfigID:  configID,
		Version:   getModifiableConfigVersion(ctx, configID, "evalhost", m),
		Hostnames: hostnamelist,
	}

	_, erru := client.UpdateEvalHost(ctx, updateEvalHost)
	if erru != nil {
		logger.Errorf("calling 'updateEvalHost': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d", configID))

	return resourceEvalHostRead(ctx, d, m)
}

func resourceEvalHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalHostRead")
	logger.Debug("in resourceEvalHostRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getEvalHostsRequest := appsec.GetEvalHostsRequest{
		ConfigID: configID,
		Version:  getLatestConfigVersion(ctx, configID, m),
	}

	evalHostResponse, err := client.GetEvalHosts(ctx, getEvalHostsRequest)
	if err != nil {
		logger.Errorf("calling 'getEvalHost': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getEvalHostsRequest.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	evalhostnameset := schema.Set{F: schema.HashString}
	for _, hostname := range evalHostResponse.Hostnames {
		evalhostnameset.Add(hostname)
	}
	if err := d.Set("hostnames", evalhostnameset.List()); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceEvalHostUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalHostUpdate")
	logger.Debug("in resourceEvalHostUpdate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	hostnames, err := tools.GetSetValue("hostnames", d)
	if err != nil {
		return diag.FromErr(err)
	}

	hostnamelist := make([]string, 0, len(hostnames.List()))
	for _, hostname := range hostnames.List() {
		hostnamelist = append(hostnamelist, hostname.(string))
	}

	updateEvalHost := appsec.UpdateEvalHostRequest{
		ConfigID:  configID,
		Version:   getModifiableConfigVersion(ctx, configID, "evalhost", m),
		Hostnames: hostnamelist,
	}

	_, erru := client.UpdateEvalHost(ctx, updateEvalHost)
	if erru != nil {
		logger.Errorf("calling 'updateEvalHost': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceEvalHostRead(ctx, d, m)
}

func resourceEvalHostDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalHostDelete")
	logger.Debug("in resourceEvalHostDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	hostnamelist := make([]string, 0)

	removeEvalHost := appsec.RemoveEvalHostRequest{
		ConfigID:  configID,
		Version:   getModifiableConfigVersion(ctx, configID, "evalhost", m),
		Hostnames: hostnamelist,
	}

	_, erru := client.RemoveEvalHost(ctx, removeEvalHost)
	if erru != nil {
		logger.Errorf("calling 'updateEvalHost': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}
