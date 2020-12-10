package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsCountries() *schema.Resource {
	return &schema.Resource{
		Description: "List all the possible countries that Akamai supports",
		ReadContext: p.tfCRUD("ds:Countries:Read", p.dsCountriesRead),
		Schema:      map[string]*schema.Schema{},
	}
}

func (p *provider) dsCountriesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
