package property

import (
	"context"
	"errors"
	"fmt"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePropertyGroupsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contract": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePropertyGroupsRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	akactx := akamai.ContextGet(inst.Name())
	log := akactx.Log("PAPI", "dataSourcePropertyGroupsRead")
	CorrelationID := "[PAPI][dataSourcePropertyGroupsRead-" + akactx.OperationID() + "]"
	var name string
	name, err := tools.GetStringValue("name", d)
	var getDefault bool
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		name = "default"
		getDefault = true
	}

	log.Debug("[Akamai Property Groups] Start Searching for property group records %s ", name)
	groups := papi.NewGroups()
	err = groups.GetGroups(CorrelationID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %q: %s", ErrLookingUpGroupByName, name, err))
	}
	contract, err := tools.GetStringValue("contract", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	group, err := findGroupByName(name, contract, groups, getDefault)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%w: %q: %s", ErrLookingUpGroupByName, name, err))
	}

	log.Debug("Searching for records [%v]", group)
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
func findGroupByName(name, contract string, groups *papi.Groups, isDefault bool) (*papi.Group, error) {
	var group *papi.Group
	var err error
	if isDefault {
		name = groups.AccountName
		if contract != "" {
			name += "-" + strings.TrimPrefix(contract, "ctr_")
			group, err = groups.FindGroup(name)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", err.Error(), ErrLookingUpGroupByName)
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
		return nil, fmt.Errorf("%w: %s", ErrNoContractProvided, name)
	}
	var foundGroups []*papi.Group
	foundGroups, err = groups.FindGroupsByName(name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrLookingUpGroupByName)
	}
	// Make sure the group belongs to the specified contract
	for _, foundGroup := range foundGroups {
		for _, c := range foundGroup.ContractIDs {
			if c == contract || c == "ctr_"+contract {
				return foundGroup, nil
			}
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrGroupNotInContract, contract)
}
