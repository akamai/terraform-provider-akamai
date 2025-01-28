package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEvalRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvalRulesRead,
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
				Description: "Unique identifier of the evaluation rule for which to retrieve information",
			},
			"eval_rule_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action to be taken for the evaluation rule if one was specified",
			},
			"condition_exception": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON-formatted condition and exception information for the rule if one was specified",
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

func dataSourceEvalRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceEvalRuleActionsRead")

	getEvalRules := appsec.GetEvalRulesRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getEvalRules.ConfigID = configID

	if getEvalRules.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getEvalRules.PolicyID = policyID

	ruleID, err := tf.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRules.RuleID = ruleID

	evalrules, err := client.GetEvalRules(ctx, getEvalRules)
	if err != nil {
		logger.Errorf("calling 'getEvalRules': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "RulesWithConditionExceptionDS", evalrules)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if len(evalrules.Rules) == 1 {
		if err := d.Set("eval_rule_action", evalrules.Rules[0].Action); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		conditionException, err := json.Marshal(evalrules.Rules[0].ConditionException)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("condition_exception", string(conditionException)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		jsonBody, err := json.Marshal(evalrules)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("json", string(jsonBody)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	d.SetId(strconv.Itoa(getEvalRules.ConfigID))

	return nil
}
