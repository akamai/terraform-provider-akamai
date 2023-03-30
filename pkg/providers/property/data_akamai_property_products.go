package property

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyProductsRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tools.IsNotBlank,
			},
			"products": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of products",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product_name": {Type: schema.TypeString, Computed: true},
						"product_id":   {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataPropertyProductsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "dataPropertyProductsRead")

	// create context with logging
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))

	client := inst.Client(meta)

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err) // fixme kind of error
	}

	logger.Debugf("[Akamai Property Products] Start searching for product records")

	prdResp, err := client.GetProducts(ctx, papi.GetProductsRequest{ContractID: contractID})
	if err != nil {
		return diag.FromErr(err) // fixme kind of error
	}

	products := make([]map[string]string, 0, len(prdResp.Products.Items))
	for _, prd := range prdResp.Products.Items {
		product := map[string]string{
			"product_name": prd.ProductName, "product_id": prd.ProductID}
		products = append(products, product)
	}

	if err := d.Set("products", products); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %q", tools.ErrValueSet, err.Error()))
	}

	jsonBody, err := json.Marshal(products)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tools.GetSHAString(string(jsonBody)))

	logger.Debugf("[Akamai Property Products] Start searching for product records")

	return nil
}
