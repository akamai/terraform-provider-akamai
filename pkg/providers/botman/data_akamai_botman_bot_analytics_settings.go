package botman

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBotAnalyticsSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBotAnalyticsSettingsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON-formatted bot analytics settings.",
			},
		},
	}
}

func dataSourceBotAnalyticsSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "dataSourceBotAnalyticsSettingsRead")
	logger.Debugf("in dataSourceBotAnalyticsSettingsRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetBotAnalyticsSettingsRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetBotAnalyticsSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetBotAnalyticsSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(configID))
	return nil
}
