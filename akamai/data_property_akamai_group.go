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

func dataSourcePropertyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][dataSourcePropertyGroupsRead-" + CreateNonce() + "]"
	var name string
	_, ok := d.GetOk("name")
	getDefault := false
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
		return fmt.Errorf("error looking up Groups for %q: %s", name, err)
	}

	var group *papi.Group
	contract, contractOk := d.GetOk("contract")

	if getDefault {
		name = groups.AccountName
		if contractOk {
			name += "-" + strings.TrimPrefix(contract.(string), "ctr_")
			group, err = groups.FindGroup(name)
		} else {
			// Find the first one
			if len(groups.Groups.Items) > 0 {
				group = groups.Groups.Items[0]
				goto groupFound
			} else {
				err = errors.New("no groups found")
			}
		}
	} else {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Searching for records NON DEFAULT [%s]", name))
		var foundGroups []*papi.Group
		foundGroups, err := groups.FindGroupsByName(name)
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Searching for records NON DEFAULT ERR=[%v]", err))
		if err == nil {
			if contractOk {
				edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Searching for records NON DEFAULT CONTRACT SUPPLIED=[%s]", contract.(string)))
				// if contract is specified make sure the group belongs to the specified contract
				for _, foundGroup := range foundGroups {
					for _, c := range foundGroup.ContractIDs {
						if c == contract.(string) || c == "ctr_"+contract.(string) {
							group = foundGroup
							goto groupFound
						}
					}
				}

				err = fmt.Errorf("group does not belong to contract %s", contract)
			} else {
				// contract is unspecified. if there is only one group return it, if more return an error
				var groupCount = len(foundGroups)
				edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Searching for records NON DEFAULT CONTRACT BLANK=[%d]", groupCount))
				if groupCount == 1 {
					group = foundGroups[0]
					goto groupFound
				} else {
					//err = fmt.Errorf("There is more then one group with name %s, please add contractId to data source definition to select one.", name)
					return fmt.Errorf("There is more then one group with name %s, please add contractId to data source definition to select one.", name)
				}
			}
		}
	}

	if err != nil {
		return fmt.Errorf("error looking up Group for %q: %s", name, err)
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Searching for records [%v]", group))

groupFound:
	d.Set("id", group.GroupID)
	d.SetId(group.GroupID)

	return nil
}
