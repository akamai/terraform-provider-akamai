package akamai

import (
	"fmt"
	"log"
	"sort"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDNSRecordSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSRecordSetRead,
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

func dataSourceDNSRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	zone := d.Get("zone").(string)
	host := d.Get("host").(string)
	recordtype := d.Get("record_type").(string)

	log.Printf("[DEBUG] [Akamai DNSv2] Start Searching for records %s %s %s ", zone, host, recordtype)

	rdata, err := dnsv2.GetRdata(zone, host, recordtype)
	if err != nil {
		return fmt.Errorf("error looking up A records for %q: %s", host, err)
	}

	log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%v]", rdata)

	sort.Strings(rdata)

	d.Set("rdata", rdata)
	d.SetId(host)

	return nil
}
