package botman

import (
	"context"
	"encoding/json"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/hash"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBotAnalyticsSettingsValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBotAnalyticsSettingsValuesRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON-formatted list of bot analytics settings values.",
			},
		},
	}
}

func dataSourceBotAnalyticsSettingsValuesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "dataSourceBotAnalyticsSettingsValuesRead")

	response, err := client.GetBotAnalyticsSettingsValues(ctx)
	if err != nil {
		logger.Errorf("calling 'GetBotAnalyticsSettingsValues': %s", err.Error())
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
