package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTuningRecommendations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTuningRecommendationsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attack_group": {
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

func dataSourceTuningRecommendationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceTuningRecommendationsRead")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := tools.GetStringValue("attack_group", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	var jsonBody []byte

	if group != "" {

		getAttackGroupRecommendationsRequest := appsec.GetAttackGroupRecommendationsRequest{
			ConfigID: configid,
			Version:  getLatestConfigVersion(ctx, configid, m),
			PolicyID: policyid,
			Group:    group,
		}

		response, err := client.GetAttackGroupRecommendations(ctx, getAttackGroupRecommendationsRequest)
		if err != nil {
			logger.Errorf("calling 'GetAttackGroupRecommendations': %s", err.Error())
			return diag.FromErr(err)
		}

		jsonBody, err = json.Marshal(response)
		if err != nil {
			return diag.FromErr(err)
		}

	} else {

		getTuningRecommendationsRequest := appsec.GetTuningRecommendationsRequest{
			ConfigID: configid,
			Version:  getLatestConfigVersion(ctx, configid, m),
			PolicyID: policyid,
		}

		response, err := client.GetTuningRecommendations(ctx, getTuningRecommendationsRequest)
		if err != nil {
			logger.Errorf("calling 'GetTuningRecommendations': %s", err.Error())
			return diag.FromErr(err)
		}

		jsonBody, err = json.Marshal(response)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(configid))

	return nil
}
