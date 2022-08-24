package botman

import (
	"context"
	"encoding/json"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBotAnalyticsCookieValue() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBotAnalyticsCookieValueRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBotAnalyticsCookieValueRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceBotAnalyticsCookieValueRead")

	response, err := client.GetBotAnalyticsCookieValues(ctx)
	if err != nil {
		logger.Errorf("calling 'GetBotAnalyticsCookieValues': %s", err.Error())
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
