package appsec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the security policy governing the bypass network lists",
			},
			"bypass_network_list": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of network list IDs that compose the bypass list",
			},
		},
	}
}

func resourceBypassNetworkListsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsCreate")
	logger.Debug("in resourceBypassNetworkListsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	networkListIDSet, err := tf.GetSetValue("bypass_network_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	networkListIDList := make([]string, 0, len(networkListIDSet.List()))
	for _, networkListID := range networkListIDSet.List() {
		networkListIDList = append(networkListIDList, networkListID.(string))
	}

	version, err := getModifiableConfigVersion(ctx, configID, "bypassnetworklists", m)
	if err != nil {
		return diag.FromErr(err)
	}
	updateBypassNetworkLists := appsec.UpdateWAPBypassNetworkListsRequest{
		ConfigID:     configID,
		Version:      version,
		PolicyID:     policyID,
		NetworkLists: networkListIDList,
	}

	_, err = client.UpdateWAPBypassNetworkLists(ctx, updateBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'UpdateWAPBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d", configID))

	return resourceBypassNetworkListsRead(ctx, d, m)
}

func resourceBypassNetworkListsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsRead")
	logger.Debug("in resourceBypassNetworkListsRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	getBypassNetworkLists := appsec.GetWAPBypassNetworkListsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	bypassNetworkListsResponse, err := client.GetWAPBypassNetworkLists(ctx, getBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'GetWAPBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	networkListIDSet := schema.Set{F: schema.HashString}
	for _, networkList := range bypassNetworkListsResponse.NetworkLists {
		networkListIDSet.Add(networkList.ID)
	}
	if err := d.Set("bypass_network_list", networkListIDSet.List()); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceBypassNetworkListsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsUpdate")
	logger.Debug("in resourceBypassNetworkListsUpdate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	networkListIDSet, err := tf.GetSetValue("bypass_network_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	networkListIDList := make([]string, 0, len(networkListIDSet.List()))
	for _, networkListID := range networkListIDSet.List() {
		networkListIDList = append(networkListIDList, networkListID.(string))
	}

	version, err := getModifiableConfigVersion(ctx, configID, "bypassnetworklists", m)
	if err != nil {
		return diag.FromErr(err)
	}
	updateBypassNetworkLists := appsec.UpdateWAPBypassNetworkListsRequest{
		ConfigID:     configID,
		Version:      version,
		PolicyID:     policyID,
		NetworkLists: networkListIDList,
	}

	_, err = client.UpdateWAPBypassNetworkLists(ctx, updateBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'UpdateWAPBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceBypassNetworkListsRead(ctx, d, m)

}

func resourceBypassNetworkListsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceBypassNetworkListsDelete")
	logger.Debug("in resourceBypassNetworkListsDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Send an empty list to remove the entire current list.
	networkListIDList := make([]string, 0)

	version, err := getModifiableConfigVersion(ctx, configID, "bypassnetworklists", m)
	if err != nil {
		return diag.FromErr(err)
	}
	removeBypassNetworkLists := appsec.RemoveWAPBypassNetworkListsRequest{
		ConfigID:     configID,
		Version:      version,
		PolicyID:     policyID,
		NetworkLists: networkListIDList,
	}

	_, err = client.RemoveWAPBypassNetworkLists(ctx, removeBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'RemoveWAPBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
