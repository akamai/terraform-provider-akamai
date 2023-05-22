package property

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
)

func dataSourcePropertyContract() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyContractRead,
		Schema: map[string]*schema.Schema{
			"group": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"group", "group_id", "group_name"},
				Deprecated:   akamai.NoticeDeprecatedUseAlias("group"),
			},
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group", "group_id", "group_name"},
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group", "group_id", "group_name"},
			},
		},
	}
}

func dataPropertyContractRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)

	log := meta.Log("PAPI", "dataPropertyContractRead")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)

	// check if one of group_id/group_name exists.
	group, err := tf.ResolveKeyStringState(d, "group_id", "group_name")
	if err != nil {
		// if both group_id/group_name not present in the state, check for group.
		group, err = tf.GetStringValue("group", d)
		// If no group found, just return the first contract
		if err != nil {
			if !errors.Is(err, tf.ErrNotFound) {
				return diag.FromErr(err)
			}
			contracts, err := inst.Client(meta).GetContracts(ctx)
			if err != nil {
				return diag.Errorf("error looking up Contracts for group %v: %s", group, err)
			}
			if len(contracts.Contracts.Items) == 0 {
				return diag.Errorf("%v", ErrNoContractsFound)
			}
			d.SetId(contracts.Contracts.Items[0].ContractID)
			return nil
		}
	}

	// Otherwise find the group and return it's first contract
	log.Debug("[Akamai Property Contract] Start Searching for property contract by group")

	groups, err := getGroups(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, g := range groups.Groups.Items {
		if g.GroupID != group && g.GroupID != tools.AddPrefix(group, "grp_") && g.GroupName != group {
			continue
		}
		if len(g.ContractIDs) == 0 {
			return diag.Errorf("%v: %v", ErrLookingUpContract, group)
		}

		// set group_id/group_name/group in state.
		if err := d.Set("group_id", tools.AddPrefix(g.GroupID, "grp_")); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
		if err := d.Set("group_name", g.GroupName); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
		d.SetId(g.ContractIDs[0])
		return nil
	}

	return diag.Errorf("%v; groupID: %v", ErrNoContractsFound, group)
}
