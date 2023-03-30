package dns

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAuthoritiesSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthoritiesSetRead,
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

func dataSourceAuthoritiesSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("AkamaiDNS", "dataSourceDNSAuthoritiesRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	contract, err := tools.GetStringValue("contract", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID := strings.TrimPrefix(contract, "ctr_")
	// Warning or Errors can be collected in a slice type
	var diags diag.Diagnostics

	logger.WithField("contractid", contractID).Debug("Start Searching for authority records")

	ns, err := inst.Client(meta).GetNameServerRecordList(ctx, contractID)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("error looking up ns records for %s", contractID),
			Detail:   err.Error(),
		})
	}
	logger.WithField("records", ns).Debug("Searching for records")

	sort.Strings(ns)
	if err := d.Set("authorities", ns); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(contractID)
	return diags
}
