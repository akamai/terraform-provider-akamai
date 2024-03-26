package botman

import (
	"context"
	"encoding/json"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/hash"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBotAnalyticsCookieValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBotAnalyticsCookieValuesRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBotAnalyticsCookieValuesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceBotAnalyticsCookieValuesRead")

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
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(hash.GetSHAString(string(jsonBody)))
	return nil
}
