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
			"uniqueid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "uniqueId",
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
	logger := meta.Log("NETWORKLIST", "resourceNetworkListRead")

	getNetworkList := network.GetNetworkListsRequest{}

	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getNetworkList.Name = name

	networkListType, err := tools.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getNetworkList.Type = networkListType

	networklist, err := client.GetNetworkLists(ctx, getNetworkList)
	if err != nil {
		logger.Errorf("calling 'getNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	if len(networklist.NetworkLists) > 0 {
		d.SetId(networklist.NetworkLists[0].UniqueID)
		if err := d.Set("uniqueid", networklist.NetworkLists[0].UniqueID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	jsonBody, err := json.Marshal(networklist)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	ids := make([]string, 0, len(networklist.NetworkLists))
	for _, networkList := range networklist.NetworkLists {
		ids = append(ids, networkList.UniqueID)
	}
	if err := d.Set("list", ids); err != nil {
		logger.Errorf("error setting 'list': %s", err.Error())
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "networkListsDS", networklist)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}
