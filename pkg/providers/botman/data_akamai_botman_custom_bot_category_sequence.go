package botman

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomBotCategorySequence() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomBotCategorySequenceRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"category_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCustomBotCategorySequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceCustomBotCategorySequenceRead")
	logger.Debugf("in dataSourceCustomBotCategorySequenceRead")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetCustomBotCategorySequenceRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetCustomBotCategorySequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomBotCategorySequence': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("category_ids", response.Sequence); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(configID))
	return nil
}
