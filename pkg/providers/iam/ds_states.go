package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsStates() *schema.Resource {
	return &schema.Resource{
		Description: "List US states or Canadian provinces",
		ReadContext: p.tfCRUD("ds:States:Read", p.dsStatesRead),
		Schema:      map[string]*schema.Schema{},
	}
}

func (p *provider) dsStatesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
