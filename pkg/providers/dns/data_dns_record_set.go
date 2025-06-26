package dns

import (
	"context"
	"fmt"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNSRecordSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSRecordSetRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The domain zone, including any nested subdomains.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A domain name, including the parent zone.",
			},
			"record_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DNS record type.",
			},
			"rdata": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "An array of data strings that represent multiple records within a set.",
			},
		},
	}
}

func dataSourceDNSRecordSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("AkamaiDNS", "dataSourceDNSRecordSetRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	zone, err := tf.GetStringValue("zone", d)
	if err != nil {
		return diag.FromErr(err)
	}
	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	recordType, err := tf.GetStringValue("record_type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debug("Start Searching for records", log.Fields{
		"zone":       zone,
		"name":       name,
		"recordtype": recordType})

	// Warning or Errors can be collected in a slice type
	var diags diag.Diagnostics
	rdata, err := inst.Client(meta).GetRdata(ctx, dns.GetRdataRequest{
		Name:       name,
		Zone:       zone,
		RecordType: recordType,
	})
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Failed retrieving recordset: %s", name),
			Detail:   err.Error(),
		})
	}
	logger.Debug("Recordset found.", "rdata", rdata)
	if recordType != RRTypeTxt {
		sort.Strings(rdata)
	}

	if err := d.Set("rdata", rdata); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}
	d.SetId(name)
	return nil
}
