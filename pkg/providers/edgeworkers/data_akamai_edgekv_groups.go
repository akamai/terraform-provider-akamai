package edgeworkers

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceEdgeKVGroups() *schema.Resource {
	return &schema.Resource{
		Description: "List edgeKV groups for given namespace and network",
		ReadContext: dataEdgeKVGroupsRead,
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
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of groups within the specified namespace that contain EdgeKV items.",
			},
		},
	}
}

func dataEdgeKVGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("EdgeKV", "dataEdgeKVGroupsRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Reading EdgeKV namespace groups")

	namespaceName, err := tf.GetStringValue("namespace_name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	network, err := tf.GetStringValue("network", d)
	if err != nil {
		return diag.FromErr(err)
	}

	groups, err := client.ListGroupsWithinNamespace(ctx, edgeworkers.ListGroupsWithinNamespaceRequest{
		Network:     edgeworkers.NamespaceNetwork(network),
		NamespaceID: namespaceName,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s", namespaceName, network))
	return nil
}
