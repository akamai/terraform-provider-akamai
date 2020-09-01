package dns

import (
	"context"
	"sort"
	"strings"

        "github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

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

func dataSourceAuthoritiesSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

        akactx := akamai.ContextGet(inst.Name())
        log := akactx.Log("DNS", "dataSourceDNSAuthoritiesRead")

        contractid := strings.TrimPrefix(d.Get("contract").(string), "ctr_")
        // Warning or Errors can be collected in a slice type
        var diags diag.Diagnostics

	log.Debug("[Akamai DNSv2] Start Searching for authority records %s ", contractid)

	ns, err := dnsv2.GetNameServerRecordList(contractid)
	if err != nil {
		diags = append(diags, diag.Errorf("error looking up A records for %q: %s", contractid, err)...)
	} else {
		log.Debug("[Akamai DNSv2] Searching for records [%v]", ns)

		sort.Strings(ns)
		d.Set("authorities", ns)
		d.SetId(contractid)
	}
	return diags
}
