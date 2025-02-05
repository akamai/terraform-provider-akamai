package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBypassNetworkLists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBypassNetworkListsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"bypass_network_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Names and unique identifiers of network lists used in the bypass network lists settings",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceBypassNetworkListsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceBypassNetworkListsRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	getBypassNetworkLists := appsec.GetWAPBypassNetworkListsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	bypassnetworklists, err := client.GetWAPBypassNetworkLists(ctx, getBypassNetworkLists)
	if err != nil {
		logger.Errorf("calling 'getBypassNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "bypassNetworkListsDS", bypassnetworklists)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(bypassnetworklists)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	nldata := make([]string, 0, len(bypassnetworklists.NetworkLists))

	for _, hosts := range bypassnetworklists.NetworkLists {
		nldata = append(nldata, hosts.ID)
	}

	if err := d.Set("bypass_network_list", nldata); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getBypassNetworkLists.ConfigID))

	return nil
}
