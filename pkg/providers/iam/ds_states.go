package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsStates() *schema.Resource {
	return &schema.Resource{
		Description: "TODO",
		ReadContext: p.dsStatesRead,
		Schema:      map[string]*schema.Schema{},
	}
}

func (p *provider) dsStatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
