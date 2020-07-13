package akamai

import (
	"errors"
	"fmt"
	"log"
	"strings"

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
	var name string
	_, ok := d.GetOk("name")
	getDefault := false
	if !ok {
		name = "default"
		getDefault = true
	} else {
		name = d.Get("name").(string)
	}

	log.Printf("[DEBUG] [Akamai Property Groups] Start Searching for property group records %s ", name)

	groups := papi.NewGroups()
	err := groups.GetGroups()
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
		var foundGroups []*papi.Group
		foundGroups, err := groups.FindGroupsByName(name)

		// Make sure the group belongs to the specified contract
		if err == nil && contractOk {
			for _, foundGroup := range foundGroups {
				for _, c := range foundGroup.ContractIDs {
					if c == contract.(string) || c == "ctr_"+contract.(string) {
						group = foundGroup
						goto groupFound
					}
				}
			}

			err = fmt.Errorf("group does not belong to contract %s", contract)
		}
	}

	if err != nil {
		return fmt.Errorf("error looking up Group for %q: %s", name, err)
	}

	log.Printf("[DEBUG] [Akamai Property Groups] Searching for records [%v]", group)

groupFound:
	d.Set("id", group.GroupID)
	d.SetId(group.GroupID)

	return nil
}
