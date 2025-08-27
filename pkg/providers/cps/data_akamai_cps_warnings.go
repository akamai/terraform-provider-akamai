package cps

import (
	"context"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPSWarnings() *schema.Resource {
	return &schema.Resource{
		Description: "Returns a map of pre- and post-verification warnings alongside with identifiers to be used in acknowledging warnings lists",
		ReadContext: dataCPSWarningsRead,
		Schema: map[string]*schema.Schema{
			"warnings": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Map of pre- and post-verification warnings consisting of the warning id and description",
			},
		},
	}
}

func dataCPSWarningsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "dataCPSWarningsRead")

	if err := d.Set("warnings", warningMap); err != nil {
		logger.Error("could not set cps warnings", "error", err)
		return diag.FromErr(err)
	}

	d.SetId("akamai_cps_warnings")

	return nil
}
