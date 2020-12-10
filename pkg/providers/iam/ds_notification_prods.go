package iam

import (
	"context"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsNotificationProds() *schema.Resource {
	return &schema.Resource{
		Description: "List all products a user can subscribe to and receive notifications for on the account",
		ReadContext: p.tfCRUD("ds:NotificationProducts:Read", p.dsNotificationProductsRead),
		Schema: map[string]*schema.Schema{
			"products": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Products a user can subscribe to and receive notifications for on the account",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func (p *provider) dsNotificationProductsRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := log.FromContext(ctx)

	logger.Debug("Fetching notification products")
	res, err := p.client.ListProducts(ctx)
	if err != nil {
		logger.WithError(err).Error("Could not get notification products")
		return diag.FromErr(err)
	}

	products := []interface{}{}
	for _, ct := range res {
		products = append(products, ct)
	}

	if err := d.Set("products", products); err != nil {
		logger.WithError(err).Error("Could not set notification products in state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_notification_prods")
	return nil
}
