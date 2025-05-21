package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of a specific rule for which to retrieve information",
			},
			"rule_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action to take when the specified rule is triggered",
			},
			"condition_exception": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The condition and exception information for the specific rule if one was specified",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceRulesRead")

	getRules := appsec.GetRulesRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getRules.ConfigID = configID

	if getRules.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getRules.PolicyID = policyID

	ruleID, err := tf.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRules.RuleID = ruleID

	rules, err := client.GetRules(ctx, getRules)
	if err != nil {
		logger.Errorf("calling 'getRules': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	getWAFMode := appsec.GetWAFModeRequest{
		ConfigID: configID,
		Version:  getRules.Version,
		PolicyID: policyID,
	}

	wafMode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	templateName := "RulesWithConditionExceptionDS"
	if wafMode.Mode == AseAuto || wafMode.Mode == AseManual {
		templateName = "ASERulesWithConditionExceptionDS"
	}

	outputtext, err := RenderTemplates(ots, templateName, rules)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if len(rules.Rules) == 1 {
		if err := d.Set("rule_action", rules.Rules[0].Action); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		conditionException, err := json.Marshal(rules.Rules[0].ConditionException)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("condition_exception", string(conditionException)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		jsonBody, err := json.Marshal(rules)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("json", string(jsonBody)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	d.SetId(strconv.Itoa(getRules.ConfigID))

	return nil
}
