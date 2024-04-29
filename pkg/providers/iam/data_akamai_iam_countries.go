package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMCountries() *schema.Resource {
	return &schema.Resource{
		Description: "List all the possible countries that Akamai supports.",
		ReadContext: dataIAMCountriesRead,
		Schema: map[string]*schema.Schema{
			"countries": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported countries.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataIAMCountriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "dataIAMCountriesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching supported countries")

	res, err := client.SupportedCountries(ctx)
	if err != nil {
		logger.Error("Could not get supported countries", "error", err)
		return diag.FromErr(err)
	}

	countries := []interface{}{}
	for _, country := range res {
		countries = append(countries, country)
	}

	if err := d.Set("countries", countries); err != nil {
		logger.Error("Could not set countries in state", "error", err)
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_countries")
	return nil
}
