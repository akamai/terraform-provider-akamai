package iam

import (
	"context"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsContactTypes() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all contact types that Akamai supports",
		ReadContext: p.tfCRUD("ds:ContactTypes:Read", p.dsContactTypesRead),
		Schema: map[string]*schema.Schema{
			"contact_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported contact types",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func (p *provider) dsContactTypesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := log.FromContext(ctx)

	logger.Debug("Fetching supported contact types")
	res, err := p.client.SupportedContactTypes(ctx)
	if err != nil {
		logger.WithError(err).Error("Could not get supported contact types")
		return diag.FromErr(err)
	}

	types := []interface{}{}
	for _, ct := range res {
		types = append(types, ct)
	}

	if err := d.Set("contact_types", types); err != nil {
		logger.WithError(err).Error("Could not set contact types in state")
		return diag.FromErr(err)
	}

	d.SetId("contact types")
	return nil
}
