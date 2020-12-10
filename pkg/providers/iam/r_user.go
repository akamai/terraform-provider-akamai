package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) resUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a user in your account",
		CreateContext: p.tfCRUD("res:User:Create", p.resUserCreate),
		ReadContext:   p.tfCRUD("res:User:Read", p.resUserRead),
		UpdateContext: p.tfCRUD("res:User:Update", p.resUserUpdate),
		DeleteContext: p.tfCRUD("res:User:Delete", p.resUserDelete),
		Importer:      p.tfImporter("res:User:Import", p.resUserImport),
		Schema: map[string]*schema.Schema{
			"identity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's unique identity ID",
			},
		},
	}
}

func (p *provider) resUserCreate(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserUpdate(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserDelete(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func (p *provider) resUserImport(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
