package dns

import (
	"context"
	"fmt"
	"sort"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func dataSourceDNSRecordSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)

	zone := d.Get("zone").(string)
	host := d.Get("host").(string)
	recordtype := d.Get("record_type").(string)

	logger := meta.Log("[Akamai DNS]", "dataSourceDNSRecordSetRead")

	logger.WithFields(log.Fields{
		"zone":       zone,
		"host":       host,
		"recordtype": recordtype,
	}).Debug("Start Searching for records")

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
	logger.WithField("rdata", rdata).Debug("Recordset found.")
	sort.Strings(rdata)

	d.Set("rdata", rdata)
	d.SetId(host)

	return diags
}
