package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourcePropertyContract() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePropertyContractRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePropertyContractRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	log.Printf("[DEBUG] [Akamai Property Contract] Start Searching for property group Contract records %s ", name)

	groups := papi.NewGroups()
	err := groups.GetGroups()
	if err != nil {
		return fmt.Errorf("error looking up Groups  for %q: %s", name, err)
	}

	group, err := groups.FindGroupId(name)

	if err != nil {
		return fmt.Errorf("error looking up Group Contract for %q: %s", name, err)
	}

	log.Printf("[DEBUG] [Akamai Property Contract] Searching for records [%v]", group)

	d.Set("id", group.ContractIDs)
	d.SetId(group.ContractIDs[0])

	return nil
}
