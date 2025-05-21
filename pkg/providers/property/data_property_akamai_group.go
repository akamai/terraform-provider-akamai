package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	akameta "github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func dataSourcePropertyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyGroupRead,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataPropertyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akameta.Must(m)
	log := meta.Log("PAPI", "dataPropertyGroupRead")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)

	var getDefault bool

	groupName, err := tf.GetStringValue("group_name", d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		groupName = "default"
		getDefault = true
	}

	log.Debugf("[Akamai Property Group] Start Searching for property group records %s ", groupName)

	groups, err := getGroups(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	group, err := findGroupByName(groupName, contractID, groups, getDefault)
	if err != nil {
		return diag.Errorf("%v: %v: %v", ErrLookingUpGroupByName, groupName, err)
	}

	if err = d.Set("group_name", group.GroupName); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	if err = d.Set("contract_id", contractID); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(group.GroupID)
	return nil
}

/*
findGroupByName returns Group struct based on provided name, contract and default name provided
for default name, either a group is returned based on provided contract, or in case of empty contract, first group is returned
TODO: we should decide whether returning first group from slice of groups is proper business behaviour

for non-default name, if contract was provided, a group with matching contract ID should be returned
in case of non-default name, contract is mandatory
*/
func findGroupByName(name, contract string, groups *papi.GetGroupsResponse, isDefault bool) (*papi.Group, error) {
	var group *papi.Group

	if isDefault {
		name = groups.AccountName
		if contract != "" {
			var found bool

			name += "-" + strings.TrimPrefix(contract, "ctr_")
			for _, group = range groups.Groups.Items {
				if group.GroupID == name {
					found = true
					break
				}
			}

			if !found {
				return nil, fmt.Errorf("group with id %q not found: %w", name, ErrLookingUpGroupByName)
			}
			return group, nil
		}
		// Find the first one
		if len(groups.Groups.Items) == 0 {
			return nil, ErrNoGroupsFound
		}
		return groups.Groups.Items[0], nil
	}

	// for non-default name, contract is required
	if contract == "" {
		return nil, fmt.Errorf("%v: %s", ErrNoContractProvided, name)
	}

	var foundGroups []*papi.Group
	for _, group = range groups.Groups.Items {
		if group.GroupName == name {
			foundGroups = append(foundGroups, group)
		}
	}

	contractMap := make(map[string][]string)
	for _, grp := range foundGroups {
		if ctrs, ok := contractMap[grp.GroupName]; ok {
			if slices.Equal(ctrs, grp.ContractIDs) {
				return nil, fmt.Errorf("there is more than 1 group with the same name and contract combination. Based on provided data, it is impossible to determine which one should be returned")
			}
		}
		contractMap[grp.GroupName] = grp.ContractIDs
	}

	// Make sure the group belongs to the specified contract
	for _, foundGroup := range foundGroups {
		for _, c := range foundGroup.ContractIDs {
			if c == contract || c == "ctr_"+contract {
				return foundGroup, nil
			}
		}
	}

	return nil, fmt.Errorf("%v: %s", ErrGroupNotInContract, contract)
}

func getGroups(ctx context.Context, meta akameta.Meta) (*papi.GetGroupsResponse, error) {
	groups, err := Client(meta).GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
