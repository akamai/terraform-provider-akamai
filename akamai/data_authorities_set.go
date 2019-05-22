package akamai

import (
	"fmt"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"sort"
)

func dataSourceAuthoritiesSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAuthoritiesSetRead,
		Schema: map[string]*schema.Schema{
			"contract": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authorities": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceAuthoritiesSetRead(d *schema.ResourceData, meta interface{}) error {
	contractid := d.Get("contract").(string)

	log.Printf("[DEBUG] [Akamai DNSv2] Start Searching for authority records %s ", contractid)

	ns, err := dnsv2.GetNameServerRecordList(contractid)
	if err != nil {
		return fmt.Errorf("error looking up A records for %q: %s", contractid, err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%v]", ns)

	sort.Strings(ns)

	d.Set("authorities", ns)
	d.SetId(contractid)

	return nil
}
