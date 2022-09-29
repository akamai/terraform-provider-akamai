package networklists

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	network "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetworkList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkListRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					IP,
					Geo,
				}, false)),
			},
			"network_list_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The ID of a specific network list to retrieve. If not supplied, information about all network lists will be returned.",
			},
			"contract_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Contract ID",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Group ID",
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
			"list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "dataSourceNetworkListRead")

	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	networkListType, err := tools.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	networkListID, err := tools.GetStringValue("network_list_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if networkListID != "" {
		networkList, err := client.GetNetworkList(ctx, network.GetNetworkListRequest{
			UniqueID: networkListID,
		})
		if err != nil {
			logger.Errorf("calling 'GetNetworkList': %s", err.Error())
			return diag.FromErr(err)
		}
		d.SetId(networkList.UniqueID)

		if err := d.Set("network_list_id", networkList.UniqueID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		if err := d.Set("group_id", networkList.GroupID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		if err := d.Set("contract_id", networkList.ContractID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		jsonBody, err := json.MarshalIndent(networkList, "", "  ")
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("json", string(jsonBody)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		if err := d.Set("list", []string{networkList.UniqueID}); err != nil {
			logger.Errorf("error setting 'list': %s", err.Error())
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		getNetworkListsResponse := network.GetNetworkListsResponse{
			NetworkLists: []network.GetNetworkListsResponseListElement{{
				ElementCount:    networkList.ElementCount,
				Name:            networkList.Name,
				NetworkListType: networkList.NetworkListType,
				ReadOnly:        networkList.ReadOnly,
				Shared:          networkList.Shared,
				SyncPoint:       networkList.SyncPoint,
				Type:            networkList.Type,
				UniqueID:        networkList.UniqueID,
				Description:     networkList.Description,
			}},
		}
		ots := OutputTemplates{}
		InitTemplates(ots)
		outputText, err := RenderTemplates(ots, "networkListsDS", getNetworkListsResponse)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		if err := d.Set("output_text", outputText); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	} else {
		networkLists, err := client.GetNetworkLists(ctx, network.GetNetworkListsRequest{
			Name: name,
			Type: networkListType,
		})
		if err != nil {
			logger.Errorf("calling 'GetNetworkLists': %s", err.Error())
			return diag.FromErr(err)
		}
		if len(networkLists.NetworkLists) > 0 {
			d.SetId(networkLists.NetworkLists[0].UniqueID)
		}
		jsonBody, err := json.MarshalIndent(networkLists, "", "  ")
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("json", string(jsonBody)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		IDs := make([]string, 0, len(networkLists.NetworkLists))
		for _, networkList := range networkLists.NetworkLists {
			IDs = append(IDs, networkList.UniqueID)
		}
		if err := d.Set("list", IDs); err != nil {
			logger.Errorf("error setting 'list': %s", err.Error())
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		ots := OutputTemplates{}
		InitTemplates(ots)

		outputText, err := RenderTemplates(ots, "networkListsDS", networkLists)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("output_text", outputText); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}
