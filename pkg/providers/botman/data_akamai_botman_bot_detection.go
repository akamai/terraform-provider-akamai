package botman

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBotDetection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBotDetectionRead,
		Schema: map[string]*schema.Schema{
			"detection_name": {
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

func dataSourceBotDetectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceBotDetectionRead")

	detectionName, err := tf.GetStringValue("detection_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := botman.GetBotDetectionListRequest{
		DetectionName: detectionName,
	}

	response, err := client.GetBotDetectionList(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetBotDetectionList': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(tools.GetSHAString(string(jsonBody)))
	return nil
}
