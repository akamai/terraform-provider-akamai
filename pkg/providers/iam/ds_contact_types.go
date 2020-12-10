package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsContactTypes() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all contact types that Akamai supports",
		ReadContext: p.tfCRUD("ds:ContractTypes:Read", p.dsContactTypesRead),
		Schema: map[string]*schema.Schema{
			"contact_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported contact types",
				Elem:        &schema.Schema{Type: schema.TypeString, Computed: true},
			},
		},
	}
}

func (p *provider) dsContactTypesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
