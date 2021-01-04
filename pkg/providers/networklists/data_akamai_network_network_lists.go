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
)

func dataSourceNetworkList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkListRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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

	/*ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "securityPoliciesDS", securitypolicy)
	if err == nil {
		d.Set("output_text", outputtext)
	}
	*/

	return nil
}
