package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsCountries() *schema.Resource {
	return &schema.Resource{
		Description: "TODO",
		ReadContext: p.dsCountriesRead,
		Schema:      map[string]*schema.Schema{},
	}
}

func (p *provider) dsCountriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
