package botman

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRecategorizedAkamaiDefinedBot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecategorizedAkamaiDefinedBotRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bot_id": {
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

func dataSourceRecategorizedAkamaiDefinedBotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceRecategorizedAkamaiDefinedBotRead")
	logger.Debugf("in dataSourceRecategorizedAkamaiDefinedBotRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	botID, err := tf.GetStringValue("bot_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := botman.GetRecategorizedAkamaiDefinedBotListRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		BotID:    botID,
	}

	response, err := client.GetRecategorizedAkamaiDefinedBotList(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetRecategorizedAkamaiDefinedBotList': %s", err.Error())
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
