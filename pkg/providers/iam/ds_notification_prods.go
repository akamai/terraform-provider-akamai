package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsNotificationProds() *schema.Resource {
	return &schema.Resource{
		Description: "Return all products a user can subscribe to and receive notifications for on the account",
		ReadContext: p.tfCRUD("ds:NotificationProducts:Read", p.dsNotificationProductsRead),
		Schema:      map[string]*schema.Schema{},
	}
}

func (p *provider) dsNotificationProductsRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
