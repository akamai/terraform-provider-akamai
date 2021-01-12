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
func resourceBypassNetworkLists() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBypassNetworkListsUpdate,
		ReadContext:   resourceBypassNetworkListsRead,
		UpdateContext: resourceBypassNetworkListsUpdate,
		DeleteContext: resourceBypassNetworkListsDelete,
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
			"bypass_network_list": {
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

func resourceBypassNetworkListsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsRead")

	getBypassNetworkLists := appsec.GetBypassNetworkListsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getBypassNetworkLists.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getBypassNetworkLists.Version = version

	bypassnetworklists, err := client.GetBypassNetworkLists(ctx, getBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'getBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "bypassNetworkListsDS", bypassnetworklists)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getBypassNetworkLists.ConfigID))

	return nil
}

func resourceBypassNetworkListsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceBypassNetworkListsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsUpdate")

	updateBypassNetworkLists := appsec.UpdateBypassNetworkListsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateBypassNetworkLists.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateBypassNetworkLists.Version = version

	netlist := d.Get("bypass_network_list").([]interface{})
	nru := make([]string, 0, len(netlist))

	for _, h := range netlist {
		nru = append(nru, h.(string))

	}
	updateBypassNetworkLists.NetworkLists = nru

	_, erru := client.UpdateBypassNetworkLists(ctx, updateBypassNetworkLists)
	if erru != nil {
		logger.Errorf("calling 'updateBypassNetworkLists': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceBypassNetworkListsRead(ctx, d, m)
}
