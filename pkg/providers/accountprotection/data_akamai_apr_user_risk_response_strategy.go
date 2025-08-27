package accountprotection

import (
	"context"
	"encoding/json"
	"strconv"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserRiskResponseStrategy() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSourceUserRiskResponseStrategy,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifies a security configuration.",
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readDataSourceUserRiskResponseStrategy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "readDataSourceUserRiskResponseStrategy")
	logger.Debugf("in readDataSourceUserRiskResponseStrategy")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.GetUserRiskResponseStrategyRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetUserRiskResponseStrategy(ctx, request)

	if err != nil {
		logger.Errorf("calling 'GetUserRiskResponseStrategy': %s", err.Error())
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
