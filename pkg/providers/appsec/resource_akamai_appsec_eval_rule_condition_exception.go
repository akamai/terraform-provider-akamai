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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceEvalRuleConditionException() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalRuleConditionExceptionUpdate,
		ReadContext:   resourceEvalRuleConditionExceptionRead,
		UpdateContext: resourceEvalRuleConditionExceptionUpdate,
		DeleteContext: resourceEvalRuleConditionExceptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Required: true,
			},
			"condition_exception": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffsConditionException,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceEvalRuleConditionExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleConditionExceptionRead")

	getEvalRuleConditionException := v2.GetEvalRuleConditionExceptionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleConditionException.PolicyID = policyid

	evalruleid, err := tools.GetIntValue("eval_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEvalRuleConditionException.RuleID = evalruleid

	evalruleconditionexception, err := client.GetEvalRuleConditionException(ctx, getEvalRuleConditionException)
	if err != nil {
		logger.Warnf("calling 'getEvalRuleConditionException': %s", err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "EvalRuleConditionExceptions", evalruleconditionexception)

	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getEvalRuleConditionException.ConfigID))

	return nil
}

func resourceEvalRuleConditionExceptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleConditionExceptionRemove")

	removeEvalRuleConditionException := v2.RemoveEvalRuleConditionExceptionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeEvalRuleConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeEvalRuleConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeEvalRuleConditionException.PolicyID = policyid

	evalruleid, err := tools.GetIntValue("eval_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeEvalRuleConditionException.RuleID = evalruleid

	logger.Errorf("calling 'RemoveEvalRuleConditionException': %v", removeEvalRuleConditionException)

	_, errd := client.RemoveEvalRuleConditionException(ctx, removeEvalRuleConditionException)
	if errd != nil {
		logger.Errorf("calling 'RemoveEvalRuleConditionException': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceEvalRuleConditionExceptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleConditionExceptionUpdate")

	updateEvalRuleConditionException := v2.UpdateEvalRuleConditionExceptionRequest{}

	jsonpostpayload := d.Get("condition_exception")

	json.Unmarshal([]byte(jsonpostpayload.(string)), &updateEvalRuleConditionException)

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalRuleConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalRuleConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalRuleConditionException.PolicyID = policyid

	evalruleid, err := tools.GetIntValue("eval_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalRuleConditionException.RuleID = evalruleid

	_, erru := client.UpdateEvalRuleConditionException(ctx, updateEvalRuleConditionException)
	if erru != nil {
		logger.Warnf("calling 'updateEvalRuleConditionException': %s", erru.Error())
	}

	return resourceEvalRuleConditionExceptionRead(ctx, d, m)
}
