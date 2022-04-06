package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *providerOld) dsStates() *schema.Resource {
	return &schema.Resource{
		Description: "List US states or Canadian provinces",
		ReadContext: p.tfCRUD("ds:States:Read", p.dsStatesRead),
		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a US state or Canadian province",
			},
			"states": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported states",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func (p *providerOld) dsStatesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	Country := d.Get("country").(string)

	logger.Debug("Fetching states")
	res, err := p.client.ListStates(ctx, iam.ListStatesRequest{Country: Country})
	if err != nil {
		logger.WithError(err).Error("Could not get states")
		return diag.FromErr(err)
	}

	states := []interface{}{}
	for _, state := range res {
		states = append(states, state)
	}

	if err := d.Set("states", states); err != nil {
		logger.WithError(err).Error("Could not set states in resource state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_states")
	return nil
}
