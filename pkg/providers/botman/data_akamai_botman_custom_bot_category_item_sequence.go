package botman

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomBotCategoryItemSequence() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomBotCategoryItemSequenceRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"category_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the bot category",
			},
			"bot_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Unique identifiers of bots in this category, sorted in preferred order",
			},
		},
	}
}

func dataSourceCustomBotCategoryItemSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceCustomBotCategoryItemSequenceRead")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, int(configID), m)
	if err != nil {
		return diag.FromErr(err)
	}
	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetCustomBotCategoryItemSequenceRequest{
		ConfigID:   configID,
		Version:    int64(version),
		CategoryID: categoryID,
	}

	response, err := client.GetCustomBotCategoryItemSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomBotCategoryItemSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("bot_ids", response.Sequence); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, categoryID))
	return nil
}
