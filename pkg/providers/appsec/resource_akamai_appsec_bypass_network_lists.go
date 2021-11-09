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
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
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
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsCreate")
	logger.Debug("in resourceBypassNetworkListsCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	networkListIDSet, err := tools.GetSetValue("bypass_network_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	networkListIDList := make([]string, 0, len(networkListIDSet.List()))
	for _, networkListID := range networkListIDSet.List() {
		networkListIDList = append(networkListIDList, networkListID.(string))
	}

	updateBypassNetworkLists := appsec.UpdateBypassNetworkListsRequest{
		ConfigID:     configid,
		Version:      getModifiableConfigVersion(ctx, configid, "bypassnetworklists", m),
		PolicyID:     policyid,
		NetworkLists: networkListIDList,
	}

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
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getBypassNetworkLists := appsec.GetBypassNetworkListsRequest{
		ConfigID: configid,
		Version:  getLatestConfigVersion(ctx, configid, m),
		PolicyID: policyid,
	}

	bypassNetworkListsResponse, err := client.GetBypassNetworkLists(ctx, getBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'GetBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	networkListIDSet := schema.Set{F: schema.HashString}
	for _, networkList := range bypassNetworkListsResponse.NetworkLists {
		networkListIDSet.Add(networkList.ID)
	}
	if err := d.Set("bypass_network_list", networkListIDSet.List()); err != nil {
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
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	networkListIDSet, err := tools.GetSetValue("bypass_network_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	networkListIDList := make([]string, 0, len(networkListIDSet.List()))
	for _, networkListID := range networkListIDSet.List() {
		networkListIDList = append(networkListIDList, networkListID.(string))
	}

	updateBypassNetworkLists := appsec.UpdateBypassNetworkListsRequest{
		ConfigID:     configid,
		Version:      getModifiableConfigVersion(ctx, configid, "bypassnetworklists", m),
		PolicyID:     policyid,
		NetworkLists: networkListIDList,
	}

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
	logger.Debug("in resourceBypassNetworkListsDelete")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	// Send an empty list to remove the entire current list.
	networkListIDList := make([]string, 0)

	removeBypassNetworkLists := appsec.RemoveBypassNetworkListsRequest{
		ConfigID:     configid,
		Version:      getModifiableConfigVersion(ctx, configid, "bypassnetworklists", m),
		PolicyID:     policyid,
		NetworkLists: networkListIDList,
	}

	_, erru := client.RemoveBypassNetworkLists(ctx, removeBypassNetworkLists)
	if erru != nil {
		logger.Errorf("calling 'RemoveBypassNetworkLists': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}
