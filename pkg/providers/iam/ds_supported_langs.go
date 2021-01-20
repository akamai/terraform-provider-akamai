package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsLanguages() *schema.Resource {
	return &schema.Resource{
		Description: "List all the possible languages Akamai supports",
		ReadContext: p.tfCRUD("ds:Languages:Read", p.dsLanguagesRead),
		Schema: map[string]*schema.Schema{
			"languages": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Languages supported by Akamai",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func (p *provider) dsLanguagesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	logger.Debug("Fetching supported supported languages")
	res, err := p.client.SupportedLanguages(ctx)
	if err != nil {
		logger.WithError(err).Error("Could not get supported supported languages")
		return diag.FromErr(err)
	}

	languages := []interface{}{}
	for _, language := range res {
		languages = append(languages, language)
	}

	if err := d.Set("languages", languages); err != nil {
		logger.WithError(err).Error("Could not set supported languages in state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_supported_langs")
	return nil
}
