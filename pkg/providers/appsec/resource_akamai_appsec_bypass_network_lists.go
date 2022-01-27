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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the security policy governing the bypass network lists",
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

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
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
		ConfigID:     configID,
		Version:      getModifiableConfigVersion(ctx, configID, "bypassnetworklists", m),
		PolicyID:     policyID,
		NetworkLists: networkListIDList,
	}

	_, err = client.UpdateBypassNetworkLists(ctx, updateBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'UpdateBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d", configID))

	return resourceBypassNetworkListsRead(ctx, d, m)
}

func resourceBypassNetworkListsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsRead")
	logger.Debug("in resourceBypassNetworkListsRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getBypassNetworkLists := appsec.GetBypassNetworkListsRequest{
		ConfigID: configID,
		Version:  getLatestConfigVersion(ctx, configID, m),
		PolicyID: policyID,
	}

	bypassNetworkListsResponse, err := client.GetBypassNetworkLists(ctx, getBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'GetBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	networkListIDSet := schema.Set{F: schema.HashString}
	for _, networkList := range bypassNetworkListsResponse.NetworkLists {
		networkListIDSet.Add(networkList.ID)
	}
	if err := d.Set("bypass_network_list", networkListIDSet.List()); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceBypassNetworkListsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsUpdate")
	logger.Debug("in resourceBypassNetworkListsUpdate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
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
		ConfigID:     configID,
		Version:      getModifiableConfigVersion(ctx, configID, "bypassnetworklists", m),
		PolicyID:     policyID,
		NetworkLists: networkListIDList,
	}

	_, err = client.UpdateBypassNetworkLists(ctx, updateBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'UpdateBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceBypassNetworkListsRead(ctx, d, m)

}

func resourceBypassNetworkListsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsDelete")
	logger.Debug("in resourceBypassNetworkListsDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	// Send an empty list to remove the entire current list.
	networkListIDList := make([]string, 0)

	removeBypassNetworkLists := appsec.RemoveBypassNetworkListsRequest{
		ConfigID:     configID,
		Version:      getModifiableConfigVersion(ctx, configID, "bypassnetworklists", m),
		PolicyID:     policyID,
		NetworkLists: networkListIDList,
	}

	_, err = client.RemoveBypassNetworkLists(ctx, removeBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'RemoveBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
