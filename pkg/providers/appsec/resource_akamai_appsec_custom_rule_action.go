package appsec

import (
	"context"
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
func resourceCustomRuleAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRuleActionUpdate,
		ReadContext:   resourceCustomRuleActionRead,
		UpdateContext: resourceCustomRuleActionUpdate,
		DeleteContext: resourceCustomRuleActionDelete,
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
			"custom_rule_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Alert,
					Deny,
					None,
				}, false),
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_rule_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceCustomRuleActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionRead")

	getCustomRuleAction := v2.GetCustomRuleActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleAction.PolicyID = policyid

	ruleid, err := tools.GetIntValue("custom_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleAction.RuleID = ruleid

	customruleaction, err := client.GetCustomRuleAction(ctx, getCustomRuleAction)
	if err != nil {
		logger.Errorf("calling 'getCustomRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Warnf("calling 'getCustomRuleAction': %s", customruleaction)

	d.Set("custom_rule_id", ruleid)
	d.SetId(strconv.Itoa(ruleid))

	return nil
}

func resourceCustomRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionRemove")

	updateCustomRuleAction := v2.UpdateCustomRuleActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.PolicyID = policyid

	ruleid, err := tools.GetIntValue("custom_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.RuleID = ruleid

	updateCustomRuleAction.Action = "none"

	_, errd := client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if errd != nil {
		logger.Errorf("calling 'removeCustomRuleAction': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceCustomRuleActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionUpdate")

	updateCustomRuleAction := v2.UpdateCustomRuleActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.PolicyID = policyid

	ruleid, err := tools.GetIntValue("custom_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.RuleID = ruleid

	customruleaction, err := tools.GetStringValue("custom_rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
		return diag.FromErr(err)
	}
	updateCustomRuleAction.Action = customruleaction

	_, erru := client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if erru != nil {
		logger.Errorf("calling 'updateCustomRuleAction': %s", erru.Error())
	}

	return resourceCustomRuleActionRead(ctx, d, m)
}

const (
	Alert = "alert"
	Deny  = "deny"
	None  = "none"
)
