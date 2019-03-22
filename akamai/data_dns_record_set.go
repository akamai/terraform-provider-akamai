package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"sort"
)

func dataSourceDnsRecordSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsRecordSetRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"record_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rdata": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceDnsRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	zone := d.Get("zone").(string)
	host := d.Get("host").(string)
	record_type := d.Get("record_type").(string)

	log.Printf("[DEBUG] [Akamai DNSv2] Start Searching for records %s %s %s ", zone, host, record_type)

	rdata, err := dnsv2.GetRdata(zone, host, record_type)
	if err != nil {
		return fmt.Errorf("error looking up A records for %q: %s", host, err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%v]", rdata)

	sort.Strings(rdata)

	d.Set("rdata", rdata)
	d.SetId(host)

	return nil
}
