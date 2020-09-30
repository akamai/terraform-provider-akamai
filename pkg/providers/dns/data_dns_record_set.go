package dns

import (
	"context"
	"fmt"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/apex/log"
	"sort"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
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
	logger := meta.Log("[Akamai DNS]", "dataSourceDNSRecordSetRead")
	zone, err := tools.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	host, err := tools.GetStringValue("host", d)
	if err != nil {
		return diag.FromErr(err)
	}
	recordType, err := tools.GetStringValue("record_type", d)
	if err != nil {
		return diag.FromErr(err)
	}

	logger.WithFields(log.Fields{
		"zone":       zone,
		"host":       host,
		"recordtype": recordType,
	}).Debug("Start Searching for records")
	// Warning or Errors can be collected in a slice type
	var diags diag.Diagnostics
	rdata, err := dnsv2.GetRdata(zone, host, recordType)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed retriving recordset: %s", host),
			Detail:   err.Error(),
		})
	}
	logger.WithField("rdata", rdata).Debug("Recordset found.")
	sort.Strings(rdata)

	if err := d.Set("rdata", rdata); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(host)
	return nil
}
