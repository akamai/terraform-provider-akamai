package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEvalRuleActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvalRuleActionsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"eval_rule_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceEvalRuleActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleActionsRead")

	getEvalRuleActions := v2.GetEvalRuleActionsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleActions.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleActions.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleActions.PolicyID = policyid

	evalruleid, err := tools.GetIntValue("eval_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleActions.RuleID = evalruleid

	evalruleactions, err := client.GetEvalRuleActions(ctx, getEvalRuleActions)
	if err != nil {
		logger.Errorf("calling 'getEvalRuleActions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "EvalRulesActionsDS", evalruleactions)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(evalruleactions)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("json", string(jsonBody))

	if len(evalruleactions.RuleActions) > 0 {
		d.Set("action", evalruleactions.RuleActions[0].Action)
	}

	d.SetId(strconv.Itoa(getEvalRuleActions.ConfigID))

	return nil
}
