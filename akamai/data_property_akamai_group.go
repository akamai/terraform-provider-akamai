package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourcePropertyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePropertyGroupsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourcePropertyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	log.Printf("[DEBUG] [Akamai Property Groups] Start Searching for property group records %s ", name)

	groups := papi.NewGroups()
	err := groups.GetGroups()
	if err != nil {
		return fmt.Errorf("error looking up Groups  for %q: %s", name, err)
	}

	group, err := groups.FindGroupId(name)

	if err != nil {
		return fmt.Errorf("error looking up Group  for %q: %s", name, err)
	}

	log.Printf("[DEBUG] [Akamai Property Groups] Searching for records [%v]", group)

	d.Set("id", group.GroupID)
	d.SetId(group.GroupID)

	return nil
}
