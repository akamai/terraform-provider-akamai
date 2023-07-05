package edgeworkers

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceEdgeKVGroupItems() *schema.Resource {
	return &schema.Resource{
		Description: "This data source can be used to retrieve all items which belong to the selected group.",
		ReadContext: dataEdgeKVGroupItems,
		Schema: map[string]*schema.Schema{
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of an EdgeKV namespace.",
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(edgeworkers.NamespaceStagingNetwork), string(edgeworkers.NamespaceProductionNetwork),
				}, false)),
				Description: "The network against which to execute the API request.",
			},
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the EdgeKV group.",
			},
			"items": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A map of items within the specified group. Each item consists of an item key and a value.",
			},
		},
	}
}

func dataEdgeKVGroupItems(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "dataEdgeKVGroupItems")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Reading EdgeKV group items")

	namespace, err := tf.GetStringValue("namespace_name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupName, err := tf.GetStringValue("group_name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	items, err := client.ListItems(ctx, edgeworkers.ListItemsRequest{
		ItemsRequestParams: edgeworkers.ItemsRequestParams{
			NamespaceID: namespace,
			GroupID:     groupName,
			Network:     edgeworkers.ItemNetwork(network),
		},
	})
	if err != nil {
		return diag.Errorf("could not list items: %s", err)
	}

	itemsMap, err := getItems(ctx, items, client, network, namespace, groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("items", itemsMap); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", namespace, network, groupName))
	return nil
}
