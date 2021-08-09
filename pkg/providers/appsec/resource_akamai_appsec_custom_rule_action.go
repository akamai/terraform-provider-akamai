package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceCustomRuleAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRuleActionCreate,
		ReadContext:   resourceCustomRuleActionRead,
		UpdateContext: resourceCustomRuleActionUpdate,
		DeleteContext: resourceCustomRuleActionDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_rule_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
		},
	}
}

func resourceCustomRuleActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionCreate")
	logger.Debugf("!!! in resourceCustomRuleActionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "customRuleAction", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	ruleid, err := tools.GetIntValue("custom_rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	customruleaction, err := tools.GetStringValue("custom_rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createCustomRuleAction := appsec.UpdateCustomRuleActionRequest{}
	createCustomRuleAction.ConfigID = configid
	createCustomRuleAction.Version = version
	createCustomRuleAction.PolicyID = policyid
	createCustomRuleAction.RuleID = ruleid
	createCustomRuleAction.Action = customruleaction

	_, erru := client.UpdateCustomRuleAction(ctx, createCustomRuleAction)
	if erru != nil {
		logger.Errorf("calling 'createCustomRuleAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s:%d", createCustomRuleAction.ConfigID, createCustomRuleAction.PolicyID, createCustomRuleAction.RuleID))

	return resourceCustomRuleActionRead(ctx, d, m)
}

func resourceCustomRuleActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionRead")
	logger.Debugf("!!! in resourceCustomRuleActionRead")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:customruleid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]
	ruleid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	getCustomRuleAction := appsec.GetCustomRuleActionRequest{}
	getCustomRuleAction.ConfigID = configid
	getCustomRuleAction.Version = version
	getCustomRuleAction.PolicyID = policyid
	getCustomRuleAction.RuleID = ruleid

	customruleaction, err := client.GetCustomRuleAction(ctx, getCustomRuleAction)
	if err != nil {
		logger.Errorf("calling 'getCustomRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}
	if err := d.Set("config_id", getCustomRuleAction.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getCustomRuleAction.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("custom_rule_id", getCustomRuleAction.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("custom_rule_action", customruleaction.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	return nil
}

func resourceCustomRuleActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionUpdate")
	logger.Debugf("!!! in resourceCustomRuleActionUpdate")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:customruleid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "customRuleAction", m)
	policyid := idParts[1]
	ruleid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}
	customruleaction, err := tools.GetStringValue("custom_rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateCustomRuleAction := appsec.UpdateCustomRuleActionRequest{}
	updateCustomRuleAction.ConfigID = configid
	updateCustomRuleAction.Version = version
	updateCustomRuleAction.PolicyID = policyid
	updateCustomRuleAction.RuleID = ruleid
	updateCustomRuleAction.Action = customruleaction

	_, erru := client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if erru != nil {
		logger.Errorf("calling 'updateCustomRuleAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceCustomRuleActionRead(ctx, d, m)
}

func resourceCustomRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionDelete")
	logger.Debugf("!!! in resourceCustomRuleActionDelete")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:customruleid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "customRuleAction", m)
	policyid := idParts[1]
	ruleid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	updateCustomRuleAction := appsec.UpdateCustomRuleActionRequest{}
	updateCustomRuleAction.ConfigID = configid
	updateCustomRuleAction.Version = version
	updateCustomRuleAction.PolicyID = policyid
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

// Constant values
const (
	Alert = "alert"
	Deny  = "deny"
	None  = "none"
)
