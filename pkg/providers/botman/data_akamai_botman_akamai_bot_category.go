package botman

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAkamaiBotCategory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAkamaiBotCategoryRead,
		Schema: map[string]*schema.Schema{
			"category_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAkamaiBotCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceAkamaiBotCategoryRead")

	categoryName, err := tools.GetStringValue("category_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := botman.GetAkamaiBotCategoryListRequest{
		CategoryName: categoryName,
	}

	response, err := client.GetAkamaiBotCategoryList(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetAkamaiBotCategoryList': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(tools.GetSHAString(string(jsonBody)))
	return nil
}
