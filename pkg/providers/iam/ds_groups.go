package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsGroups() *schema.Resource {
	return &schema.Resource{
		Description: "TODO",
		ReadContext: p.dsGroupsRead,
		Schema:      map[string]*schema.Schema{},
	}
}

func (p *provider) dsGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
