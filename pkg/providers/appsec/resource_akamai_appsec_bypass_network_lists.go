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
func resourceBypassNetworkLists() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBypassNetworkListsCreate,
		ReadContext:   resourceBypassNetworkListsRead,
		UpdateContext: resourceBypassNetworkListsUpdate,
		DeleteContext: resourceBypassNetworkListsDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bypass_network_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceBypassNetworkListsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsUpdate")
	logger.Debug("in resourceBypassNetworkListsCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	networklistidset, err := tools.GetSetValue("bypass_network_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateBypassNetworkLists := appsec.UpdateBypassNetworkListsRequest{}
	updateBypassNetworkLists.ConfigID = configid
	updateBypassNetworkLists.Version = getModifiableConfigVersion(ctx, configid, "bypassnetworklists", m)
	networklistidlist := make([]string, 0, len(networklistidset.List()))
	for _, networklistid := range networklistidset.List() {
		networklistidlist = append(networklistidlist, networklistid.(string))
	}
	updateBypassNetworkLists.NetworkLists = networklistidlist

	_, err = client.UpdateBypassNetworkLists(ctx, updateBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'UpdateBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d", configid))

	return resourceBypassNetworkListsRead(ctx, d, m)
}

func resourceBypassNetworkListsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsRead")
	logger.Debug("in resourceBypassNetworkListsRead")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getBypassNetworkLists := appsec.GetBypassNetworkListsRequest{}
	getBypassNetworkLists.ConfigID = configid
	getBypassNetworkLists.Version = getLatestConfigVersion(ctx, configid, m)

	bypassnetworklistsresponse, err := client.GetBypassNetworkLists(ctx, getBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'GetBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	networklistidset := schema.Set{F: schema.HashString}
	for _, networklist := range bypassnetworklistsresponse.NetworkLists {
		networklistidset.Add(networklist.ID)
	}
	if err := d.Set("bypass_network_list", networklistidset.List()); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceBypassNetworkListsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsUpdate")
	logger.Debug("in resourceBypassNetworkListsUpdate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	networklistidset, err := tools.GetSetValue("bypass_network_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateBypassNetworkLists := appsec.UpdateBypassNetworkListsRequest{}
	updateBypassNetworkLists.ConfigID = configid
	updateBypassNetworkLists.Version = getModifiableConfigVersion(ctx, configid, "bypassnetworklists", m)
	networklistidlist := make([]string, 0, len(networklistidset.List()))
	for _, networklistid := range networklistidset.List() {
		networklistidlist = append(networklistidlist, networklistid.(string))
	}
	updateBypassNetworkLists.NetworkLists = networklistidlist

	_, erru := client.UpdateBypassNetworkLists(ctx, updateBypassNetworkLists)
	if erru != nil {
		logger.Errorf("calling 'UpdateBypassNetworkLists': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceBypassNetworkListsRead(ctx, d, m)

}

func resourceBypassNetworkListsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsDelete")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	bypassnetworkidset, err := tools.GetSetValue("bypass_network_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	removeBypassNetworkLists := appsec.RemoveBypassNetworkListsRequest{}
	removeBypassNetworkLists.ConfigID = configid
	removeBypassNetworkLists.Version = getModifiableConfigVersion(ctx, configid, "bypassnetworklists", m)
	networklistidlist := make([]string, 0, len(bypassnetworkidset.List()))
	for _, networklistid := range bypassnetworkidset.List() {
		networklistidlist = append(networklistidlist, networklistid.(string))
	}
	removeBypassNetworkLists.NetworkLists = networklistidlist

	_, erru := client.RemoveBypassNetworkLists(ctx, removeBypassNetworkLists)
	if erru != nil {
		logger.Errorf("calling 'RemoveBypassNetworkLists': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}
