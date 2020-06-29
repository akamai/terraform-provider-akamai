package akamai

import (
	"fmt"

	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourcePropertyContract() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePropertyContractRead,
		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePropertyContractRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][dataSourcePropertyContractRead-" + CreateNonce() + "]"
	_, groupOk := d.GetOk("group")
	group := d.Get("group").(string)
	contracts := papi.NewContracts()
	// If no group, just return the first contract
	if !groupOk {
		err := contracts.GetContracts(CorrelationID)
		if err != nil {
			return fmt.Errorf("error looking up Contracts for group %q: %s", group, err)
		}

		d.SetId(contracts.Contracts.Items[0].ContractID)
		return nil
	}

	// Otherwise find the group and return it's first contract
	///log.Print("[DEBUG] [Akamai Property Contract] Start Searching for property contract by group")
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " [Akamai Property Contract] Start Searching for property contract by group")
	groups, err := papi.GetGroups()
	if err != nil {
		return fmt.Errorf("error looking up Group Contract for %q: %s", group, err)
	}

	for _, g := range groups.Groups.Items {
		if g.GroupID == group || g.GroupID == "grp_"+group || g.GroupName == group {
			d.SetId(g.ContractIDs[0])
			return nil
		}
	}

	return fmt.Errorf("error looking up Group Contract for %q: %s", group, err)
}
