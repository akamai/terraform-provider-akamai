package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"custom_rule_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the custom rule whose action is being modified",
			},
			"custom_rule_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions,
				Description:      "Action to be taken when the custom rule is invoked",
			},
		},
	}
}

func resourceCustomRuleActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionCreate")
	logger.Debugf("in resourceCustomRuleActionCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "customRuleAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ruleID, err := tools.GetIntValue("custom_rule_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	customruleaction, err := tools.GetStringValue("custom_rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createCustomRuleAction := appsec.UpdateCustomRuleActionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		RuleID:   ruleID,
		Action:   customruleaction,
	}

	_, err = client.UpdateCustomRuleAction(ctx, createCustomRuleAction)
	if err != nil {
		logger.Errorf("calling 'createCustomRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%d", configID, policyID, ruleID))

	return resourceCustomRuleActionRead(ctx, d, m)
}

func resourceCustomRuleActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionRead")
	logger.Debugf("in resourceCustomRuleActionRead")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:customRuleID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	ruleID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	getCustomRuleAction := appsec.GetCustomRuleActionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		RuleID:   ruleID,
	}

	customruleaction, err := client.GetCustomRuleAction(ctx, getCustomRuleAction)
	if err != nil {
		logger.Errorf("calling 'getCustomRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getCustomRuleAction.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getCustomRuleAction.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("custom_rule_id", getCustomRuleAction.RuleID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("custom_rule_action", customruleaction.Action); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	return nil
}

func resourceCustomRuleActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionUpdate")
	logger.Debugf("in resourceCustomRuleActionUpdate")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:customRuleID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "customRuleAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	ruleID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}
	customruleaction, err := tools.GetStringValue("custom_rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateCustomRuleAction := appsec.UpdateCustomRuleActionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		RuleID:   ruleID,
		Action:   customruleaction,
	}

	_, err = client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if err != nil {
		logger.Errorf("calling 'updateCustomRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceCustomRuleActionRead(ctx, d, m)
}

func resourceCustomRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionDelete")
	logger.Debugf("in resourceCustomRuleActionDelete")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:customRuleID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "customRuleAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	ruleID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	updateCustomRuleAction := appsec.UpdateCustomRuleActionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		RuleID:   ruleID,
		Action:   "none",
	}

	_, err = client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if err != nil {
		logger.Errorf("calling 'removeCustomRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}

// Constant values
const (
	Alert = "alert"
	Deny  = "deny"
	None  = "none"
)
