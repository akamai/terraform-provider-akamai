package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
)

func dataSourcePropertyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "group_name"},
				Deprecated:   akamai.NoticeDeprecatedUseAlias("name"),
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "group_name"},
			},
			"contract": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract", "contract_id"},
				Deprecated:   akamai.NoticeDeprecatedUseAlias("contract"),
			},
			"contract_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract", "contract_id"},
			},
		},
	}
}

func dataPropertyGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataPropertyGroupRead")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)

	var (
		name       string
		getDefault bool
	)

	// check and load group_name, if not exists then check group.
	name, err := tools.ResolveKeyStringState(d, "group_name", "name")
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		name = "default"
		getDefault = true
	}

	log.Debugf("[Akamai Property Group] Start Searching for property group records %s ", name)

	groups, err := getGroups(ctx, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// check and load contract_id, if not exists then check contract.
	contractID, err := tools.ResolveKeyStringState(d, "contract_id", "contract")
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	group, err := findGroupByName(name, contractID, groups, getDefault)
	if err != nil {
		return diag.Errorf("%v: %v: %v", ErrLookingUpGroupByName, name, err)
	}

	// set group_name/name in state.
	if err := d.Set("group_name", group.GroupName); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("name", group.GroupName); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	if len(group.ContractIDs) != 0 {
		contractID = group.ContractIDs[0]
	}
	// set contract/contract_id in state.
	if err := d.Set("contract", contractID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
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
	for _, group := range groups.Groups.Items {
		if group.GroupName == name {
			foundGroups = append(foundGroups, group)
		}
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

func getGroups(ctx context.Context, meta akamai.OperationMeta) (*papi.GetGroupsResponse, error) {
	groups := &papi.GetGroupsResponse{}
	if err := meta.CacheGet(inst, "groups", groups); err != nil {
		if !akamai.IsNotFoundError(err) && !errors.Is(err, akamai.ErrCacheDisabled) {
			return nil, err
		}
		groups, err = inst.Client(meta).GetGroups(ctx)
		if err != nil {
			return nil, err
		}
		if err := meta.CacheSet(inst, "groups", groups); err != nil {
			if !errors.Is(err, akamai.ErrCacheDisabled) {
				return nil, err
			}
		}
	}

	return groups, nil
}
