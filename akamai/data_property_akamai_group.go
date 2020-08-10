package akamai

import (
	"errors"
	"fmt"
	"strings"

	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePropertyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePropertyGroupsRead,
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

func dataSourcePropertyGroupsRead(d *schema.ResourceData, _ interface{}) error {
	CorrelationID := "[PAPI][dataSourcePropertyGroupsRead-" + CreateNonce() + "]"
	var name string
	_, ok := d.GetOk("name")
	var getDefault bool
	if !ok {
		name = "default"
		getDefault = true
	} else {
		name = d.Get("name").(string)
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  [Akamai Property Groups] Start Searching for property group records %s ", name))
	groups := papi.NewGroups()
	err := groups.GetGroups(CorrelationID)
	if err != nil {
		return fmt.Errorf("looking up Groups for %q: %s", name, err)
	}
	contract, contractOk := d.GetOk("contract")
	var contractStr string
	if contractOk {
		contractStr, ok = contract.(string)
	}

	group, err := findGroupByName(name, contractStr, groups, getDefault)
	if err != nil {
		return fmt.Errorf("looking up Group for %q: %s", name, err)
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Searching for records [%v]", group))
	d.SetId(group.GroupID)
	return nil
}

var (
	ErrNoGroupsFound       = errors.New("no groups found")
	ErrGroupNotFound       = errors.New("could not find group for given group ID")
	ErrFindingGroupsByName = errors.New("could not find groups for given name")
	ErrNoContractProvided  = errors.New("contract ID is required for non-default name")
	ErrGroupNotInContract  = errors.New("group does not belong to contract")
)

/*
findGroupByName returns Group struct based on provided name, contract and default name provided
for default name, either a group is returned based on provided contract, or in case of empty contract, first group is returned
TODO: we should decide whether returning first group from slice of groups is proper business behaviour

for non-default name, if contract was provided, a group with matching contract ID should be returned
in case of non-default name, contract is mandatory
*/
func findGroupByName(name string, contract string, groups *papi.Groups, isDefault bool) (*papi.Group, error) {
	var group *papi.Group
	var err error
	if isDefault {
		name = groups.AccountName
		if contract != "" {
			name += "-" + strings.TrimPrefix(contract, "ctr_")
			group, err = groups.FindGroup(name)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", err.Error(), ErrGroupNotFound)
			}
		} else {
			// Find the first one
			if len(groups.Groups.Items) == 0 {
				return nil, ErrNoGroupsFound
			}
			group = groups.Groups.Items[0]
		}
	} else {
		// for non-default name, contract is required
		if contract == "" {
			return nil, fmt.Errorf("%w: %s", ErrNoContractProvided, name)
		}
		var foundGroups []*papi.Group
		foundGroups, err := groups.FindGroupsByName(name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), ErrFindingGroupsByName)
		}
		// Make sure the group belongs to the specified contract
	FoundGroupsLoop:
		for _, foundGroup := range foundGroups {
			for _, c := range foundGroup.ContractIDs {
				if c == contract || c == "ctr_"+contract {
					group = foundGroup
					break FoundGroupsLoop
				}
			}
		}
		if group == nil {
			return nil, fmt.Errorf("%w: %s", ErrGroupNotInContract, contract)
		}
	}
	return group, nil
}
