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
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"type"},
			},
			"uniqueid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "uniqueId",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					IP,
					Geo,
				}, false),
				RequiredWith: []string{"name"},
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

	_type, err := tools.GetStringValue("type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	networklist, err := client.GetNetworkLists(ctx, getNetworkList)
	if err != nil {
		logger.Errorf("calling 'getNetworkLists': %s", err.Error())
		return diag.FromErr(err)
	}

	if len(networklist.NetworkLists) > 0 {
		if len(name) > 0 {
			for _, networkList := range networklist.NetworkLists {
				if networkList.UniqueID == name && networkList.Type == _type {
					d.SetId(networkList.UniqueID)
					if err := d.Set("uniqueid", networkList.UniqueID); err != nil {
						return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
					}

					jsonBody, err := json.Marshal(networkList)
					if err != nil {
						return diag.FromErr(err)
					}
					if err := d.Set("json", string(jsonBody)); err != nil {
						return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
					}

					ids := make([]string, 0, 1)
					ids = append(ids, networkList.UniqueID)
					if err := d.Set("list", ids); err != nil {
						logger.Errorf("error setting 'list': %s", err.Error())
						return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
					}
					break // can there be more than one list with same name & type?
				}
			}
		} else {
			d.SetId(networklist.NetworkLists[0].UniqueID)
			if err := d.Set("uniqueid", networklist.NetworkLists[0].UniqueID); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
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
		}
	}

	return nil
}
