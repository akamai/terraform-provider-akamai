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
		Schema:      map[string]*schema.Schema{},
	}
}

func (p *provider) dsLanguagesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
