package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTuningRecommendations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTuningRecommendationsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration for which to return tuning recommendations",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of the security policy for which to return tuning recommendations",
			},
			"ruleset_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(appsec.RulesetTypeActive),
				ValidateFunc: validation.StringInSlice([]string{
					string(appsec.RulesetTypeActive),
					string(appsec.RulesetTypeEvaluation),
				}, false),
				Description: "Type of the ruleset of the security configuration for which to return tuning recommendations",
			},
			"attack_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique name of a specific attack group for which to return tuning recommendations",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique name of a specific rule id for which to return tuning recommendations",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON-formatted list of the tuning recommendations for the security policy, attack group or rule",
			},
		},
	}
}

func dataSourceTuningRecommendationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceTuningRecommendationsRead")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := tools.GetStringValue("attack_group", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	ruleID, err := tools.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if group != "" && ruleID != 0 {
		var diags diag.Diagnostics
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "wrong field values specified",
			Detail:   "attack_group and rule_id fields can not be used together",
		}
		diags = append(diags, diag)
		return diags
	}

	rulesetType, err := tools.GetStringValue("ruleset_type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	var jsonBody []byte

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	if group != "" {
		getAttackGroupRecommendationsRequest := appsec.GetAttackGroupRecommendationsRequest{
			ConfigID:    configID,
			Version:     version,
			PolicyID:    policyID,
			Group:       group,
			RulesetType: appsec.RulesetType(rulesetType),
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
	} else if ruleID != 0 {
		getRuleRecommendationsRequest := appsec.GetRuleRecommendationsRequest{
			ConfigID:    configID,
			Version:     version,
			PolicyID:    policyID,
			RuleID:      ruleID,
			RulesetType: appsec.RulesetType(rulesetType),
		}
		response, err := client.GetRuleRecommendations(ctx, getRuleRecommendationsRequest)
		if err != nil {
			logger.Errorf("calling 'GetRuleRecommendations': %s", err.Error())
			return diag.FromErr(err)
		}

		jsonBody, err = json.Marshal(response)
		if err != nil {
			return diag.FromErr(err)
		}

	} else {
		getTuningRecommendationsRequest := appsec.GetTuningRecommendationsRequest{
			ConfigID:    configID,
			Version:     version,
			PolicyID:    policyID,
			RulesetType: appsec.RulesetType(rulesetType),
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
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(configID))

	return nil
}
