package gtm

import (
	"context"
	"fmt"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/configgtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"

	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGTMDefaultDatacenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGTMDefaultDatacenterRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  gtm.MapDefaultDC,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{
					gtm.MapDefaultDC,
					gtm.Ipv4DefaultDC,
					gtm.Ipv6DefaultDC,
				})),
			},
			"datacenter_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGTMDefaultDatacenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "dataSourceGTMDefaultDatacenterRead")

	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		logger.Errorf("[Error] GTM dataSourceGTMDefaultDatacenterRead: Domain not initialized")

		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	// get or create default dc
	dcid, ok := d.Get("datacenter").(int)
	if !ok {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("[Error] GTM dataSourceGTMDefaultDatacenterRead: invalid Default Datacenter %d in configuration", dcid),
		})
	}
	logger.WithFields(log.Fields{
		"domain": domain,
		"dcid":   dcid,
	}).Debug("Start Default Datacenter Retrieval")

	var defaultDC = inst.Client(meta).NewDatacenter(ctx)
	switch dcid {
	case gtm.MapDefaultDC:
		defaultDC, err = inst.Client(meta).CreateMapsDefaultDatacenter(ctx, domain)
	case gtm.Ipv4DefaultDC:
		defaultDC, err = inst.Client(meta).CreateIPv4DefaultDatacenter(ctx, domain)
	case gtm.Ipv6DefaultDC:
		defaultDC, err = inst.Client(meta).CreateIPv6DefaultDatacenter(ctx, domain)
	default:
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("[Error] GTM dataSourceGTMDefaultDatacenterRead: invalid Default Datacenter %d in configuration", dcid),
		})
	}

	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("[Error] GTM dataSourceGTMDefaultDatacenterRead: invalid Default Datacenter %d in configuration", dcid),
			Detail:   err.Error(),
		})
	}
	if err := d.Set("nickname", defaultDC.Nickname); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "GTM dataSourceGTMDefaultDatacenterRead: setting nickname failed.",
			Detail:   err.Error(),
		})
	}
	if err := d.Set("datacenter_id", defaultDC.DatacenterId); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "GTM dataSourceGTMDefaultDatacenterRead: setting datacenter id failed.",
			Detail:   err.Error(),
		})
	}
	defaultDatacenterID := fmt.Sprintf("%s:%s:%d", domain, "default_datcenter", defaultDC.DatacenterId)
	logger.Debugf("DataSourceGTMDefaultDatacenterRead: generated Default DC Resource Id: %s", defaultDatacenterID)
	d.SetId(defaultDatacenterID)

	return nil
}
