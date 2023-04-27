package property

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
)

func dataSourcePropertyMultipleGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyMultipleGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of groups",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contract_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataPropertyMultipleGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "dataPropertyMultipleGroupsRead")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("[Akamai Property Groups] Start Searching for group records")

	groups, err := getGroups(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	grpList := make([]map[string]interface{}, 0, len(groups.Groups.Items))
	for _, g := range groups.Groups.Items {
		contractIDs := make([]string, 0, len(g.ContractIDs))
		for _, contractID := range g.ContractIDs {
			contractIDs = append(contractIDs, contractID)
		}
		grpList = append(grpList, map[string]interface{}{
			"group_id":        g.GroupID,
			"group_name":      g.GroupName,
			"parent_group_id": g.ParentGroupID,
			"contract_ids":    contractIDs,
		})
	}

	if err := d.Set("groups", grpList); err != nil {
		return diag.FromErr(fmt.Errorf("%w:%q", tf.ErrValueSet, err.Error()))
	}

	jsonBody, err := json.Marshal(grpList)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tools.GetSHAString(string(jsonBody)))

	logger.Debugf("[Akamai Property Groups] Done")

	return nil
}
