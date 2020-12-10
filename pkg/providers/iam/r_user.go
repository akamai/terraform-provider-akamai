package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) resUser() *schema.Resource {
	return &schema.Resource{
		Description:   "TODO",
		CreateContext: p.resUserCreate,
		ReadContext:   p.resUserRead,
		UpdateContext: p.resUserUpdate,
		DeleteContext: p.resUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: p.resUserImport,
		},
		Schema: map[string]*schema.Schema{
			"identity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's unique identity ID",
			},
		},
	}
}

func (p *provider) resUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
