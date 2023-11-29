package property

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	akameta "github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyContract() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyContractRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group_id", "group_name"},
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group_id", "group_name"},
			},
		},
	}
}

func dataPropertyContractRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akameta.Must(m)
	log := meta.Log("PAPI", "dataPropertyContractRead")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)

	// check if one of group_id/group_name exists.
	group, err := tf.ResolveKeyStringState(d, "group_id", "group_name")
	if err != nil {
		// If no group found, just return the first contract if only one exists
		if !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		contracts, err := Client(meta).GetContracts(ctx)
		if err != nil {
			return diag.Errorf("error looking up Contracts for group %v: %s", group, err)
		}
		if len(contracts.Contracts.Items) == 0 {
			return diag.Errorf("%v", ErrNoContractsFound)
		}
		if len(contracts.Contracts.Items) > 1 {
			return ErrMultipleContractsFound
		}
		d.SetId(contracts.Contracts.Items[0].ContractID)
		return nil
	}

	// Otherwise find the group and return it's first contract
	log.Debug("[Akamai Property Contract] Start Searching for property contract by group")

	groups, err := getGroups(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	var foundGroups []*papi.Group
	for _, g := range groups.Groups.Items {
		if isGroupEqual(g, group) {
			foundGroups = append(foundGroups, g)
		}
	}

	if len(foundGroups) > 1 {
		return diag.Errorf("there is more than 1 group with the same name. Based on provided data, it is impossible to determine which one should be returned. Please use group_id attribute")
	} else if len(foundGroups) == 0 {
		return diag.Errorf("%v; groupID: %v", ErrNoContractsFound, group)
	}
	if len(foundGroups[0].ContractIDs) == 0 {
		return diag.Errorf("%v: %v", ErrLookingUpContract, group)
	}
	if len(foundGroups[0].ContractIDs) > 1 {
		return ErrMultipleContractsInGroup
	}

	if err = d.Set("group_id", str.AddPrefix(foundGroups[0].GroupID, "grp_")); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err = d.Set("group_name", foundGroups[0].GroupName); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(foundGroups[0].ContractIDs[0])

	return nil
}

func isGroupEqual(group *papi.Group, target string) bool {
	return group.GroupID == target || group.GroupID == str.AddPrefix(target, "grp_") || group.GroupName == target
}
