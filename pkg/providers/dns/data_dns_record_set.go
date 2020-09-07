package dns

import (
	"context"
	"fmt"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sort"
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

	logger := akactx.Log("[Akamai DNS]", "dataSourceDNSRecordSetRead")

	logger.Debug("Start Searching for records", "zone", zone, "host", host, "recordtype", recordtype)
	// Warning or Errors can be collected in a slice type
	var diags diag.Diagnostics

	rdata, err := dnsv2.GetRdata(zone, host, recordtype)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed retriving recordset: %s", host),
			Detail:   err.Error(),
		})
	}
	logger.Debug("Recordset found.", "rdata", rdata)
	sort.Strings(rdata)

	d.Set("rdata", rdata)
	d.SetId(host)

	return diags
}
