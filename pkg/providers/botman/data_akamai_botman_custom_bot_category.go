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

func dataSourceCustomBotCategory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomBotCategoryRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"category_id": {
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

func dataSourceCustomBotCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceCustomBotCategoryRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := botman.GetCustomBotCategoryListRequest{
		ConfigID:   int64(configID),
		Version:    int64(version),
		CategoryID: categoryID,
	}

	response, err := client.GetCustomBotCategoryList(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomBotCategoryList': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(configID))
	return nil
}
