package dns

import (
	"context"
	"sort"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func dataSourceDNSRecordSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSRecordSetRead,
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

func dataSourceDNSRecordSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	akactx := akamai.ContextGet(inst.Name())
	zone := d.Get("zone").(string)
	host := d.Get("host").(string)
	recordtype := d.Get("record_type").(string)

	log := akactx.Log("DNS", "dataSourceDNSRecordSetRead")
	//CorrelationID := "[DNS][dataSourceDNSRecordSetRead-" + akactx.OperationID() + "]"

	log.Debug("[Akamai DNSv2] Start Searching for records %s %s %s ", zone, host, recordtype)
	// Warning or Errors can be collected in a slice type
	var diags diag.Diagnostics

	rdata, err := dnsv2.GetRdata(zone, host, recordtype)
	if err != nil {
		diags = append(diags, diag.Errorf("error looking up A records for %q: %s", host, err)...)
	} else {
		log.Debug("[Akamai DNSv2] Searching for records [%v]", rdata)
		sort.Strings(rdata)

		d.Set("rdata", rdata)
		d.SetId(host)
	}
	return diags
}
