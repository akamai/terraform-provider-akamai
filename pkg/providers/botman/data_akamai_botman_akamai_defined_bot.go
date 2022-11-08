package botman

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAkamaiDefinedBot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAkamaiDefinedBotRead,
		Schema: map[string]*schema.Schema{
			"bot_name": {
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

func dataSourceAkamaiDefinedBotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceAkamaiDefinedBotRead")

	botName, err := tools.GetStringValue("bot_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := botman.GetAkamaiDefinedBotListRequest{
		BotName: botName,
	}

	response, err := client.GetAkamaiDefinedBotList(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetAkamaiDefinedBotList': %s", err.Error())
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
