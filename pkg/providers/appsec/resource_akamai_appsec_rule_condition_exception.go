package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceRuleConditionException() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleConditionExceptionUpdate,
		ReadContext:   resourceRuleConditionExceptionRead,
		UpdateContext: resourceRuleConditionExceptionUpdate,
		DeleteContext: resourceRuleConditionExceptionDelete,
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
			"rule_id": {
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

func resourceRuleConditionExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleConditionExceptionRead")

	getRuleConditionException := appsec.GetRuleConditionExceptionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRuleConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRuleConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRuleConditionException.PolicyID = policyid

	ruleid, err := tools.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRuleConditionException.RuleID = ruleid

	ruleconditionexception, err := client.GetRuleConditionException(ctx, getRuleConditionException)
	if err != nil {
		logger.Errorf("calling 'getRuleConditionException': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "RuleConditionExceptions", ruleconditionexception)

	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(ruleconditionexception)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("condition_exception", string(jsonBody))

	d.SetId(strconv.Itoa(getRuleConditionException.ConfigID))

	return nil
}

func resourceRuleConditionExceptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleConditionExceptionRemove")

	removeRuleConditionException := appsec.RemoveRuleConditionExceptionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeRuleConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeRuleConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeRuleConditionException.PolicyID = policyid

	ruleid, err := tools.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeRuleConditionException.RuleID = ruleid

	logger.Errorf("calling 'RemoveRuleConditionException': %v", removeRuleConditionException)

	_, errd := client.RemoveRuleConditionException(ctx, removeRuleConditionException)
	if errd != nil {
		logger.Errorf("calling 'RemoveRuleConditionException': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceRuleConditionExceptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleConditionExceptionUpdate")

	updateRuleConditionException := appsec.UpdateRuleConditionExceptionRequest{}

	jsonpostpayload := d.Get("condition_exception")

	json.Unmarshal([]byte(jsonpostpayload.(string)), &updateRuleConditionException)

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleConditionException.PolicyID = policyid

	ruleid, err := tools.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleConditionException.RuleID = ruleid
	logger.Errorf("calling 'updateRuleConditionException': %s", updateRuleConditionException)
	_, erru := client.UpdateRuleConditionException(ctx, updateRuleConditionException)
	if erru != nil {
		logger.Errorf("calling 'updateRuleConditionException': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceRuleConditionExceptionRead(ctx, d, m)
}
