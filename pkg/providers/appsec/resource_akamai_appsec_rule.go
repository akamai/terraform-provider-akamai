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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleCreate,
		ReadContext:   resourceRuleRead,
		UpdateContext: resourceRuleUpdate,
		DeleteContext: resourceRuleDelete,
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
			"rule_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rule_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: ValidateActions,
			},
			"condition_exception": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffsConditionException,
			},
		},
	}
}

func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleCreate")
	logger.Debugf("in resourceRuleCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "rule", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ruleid, err := tools.GetIntValue("rule_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	conditionexception, err := tools.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	jsonPayloadRaw := []byte(conditionexception)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	getWAFMode := appsec.GetWAFModeRequest{}

	getWAFMode.ConfigID = configid
	getWAFMode.Version = version
	getWAFMode.PolicyID = policyid

	wafmode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if wafmode.Mode == AseAuto { // action is read only, only exception is writable
		createRule := appsec.UpdateConditionExceptionRequest{}
		createRule.ConfigID = configid
		createRule.Version = version
		createRule.PolicyID = policyid
		createRule.RuleID = ruleid
		var ruleConditionException appsec.RuleConditionException
		err = json.Unmarshal([]byte(rawJSON), &ruleConditionException)
		if err != nil {
			return diag.FromErr(err)
		}
		createRule.Conditions = ruleConditionException.Conditions
		createRule.Exception = ruleConditionException.Exception
		resp, err := client.UpdateRuleConditionException(ctx, createRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
		logger.Debugf("calling 'UpdateRule Response': %s", resp)
		d.SetId(fmt.Sprintf("%d:%s:%d", createRule.ConfigID, createRule.PolicyID, createRule.RuleID))
	} else {
		action, err := tools.GetStringValue("rule_action", d)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := validateActionAndConditionException(action, conditionexception); err != nil {
			return diag.FromErr(err)
		}

		createRule := appsec.UpdateRuleRequest{}
		createRule.ConfigID = configid
		createRule.Version = version
		createRule.PolicyID = policyid
		createRule.RuleID = ruleid
		createRule.Action = action
		createRule.JsonPayloadRaw = rawJSON

		resp, err := client.UpdateRule(ctx, createRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
		logger.Debugf("calling 'UpdateRule Response': %s", resp)
		d.SetId(fmt.Sprintf("%d:%s:%d", createRule.ConfigID, createRule.PolicyID, createRule.RuleID))

	}

	return resourceRuleRead(ctx, d, m)
}

func resourceRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleRead")
	logger.Debugf("in resourceRuleRead")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:ruleid")
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

	getRule := appsec.GetRuleRequest{}
	getRule.ConfigID = configid
	getRule.Version = version
	getRule.PolicyID = policyid
	getRule.RuleID = ruleid

	rule, err := client.GetRule(ctx, getRule)
	if err != nil {
		logger.Errorf("calling 'GetRule': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getRule.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getRule.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("rule_id", getRule.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("rule_action", rule.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if !rule.IsEmptyConditionException() {
		jsonBody, err := json.Marshal(rule.ConditionException)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("condition_exception", string(jsonBody)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleUpdate")
	logger.Debugf("in resourceRuleUpdate")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:ruleid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	policyid := idParts[1]
	version := getModifiableConfigVersion(ctx, configid, "threatIntel", m)
	ruleid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}
	conditionexception, err := tools.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	jsonPayloadRaw := []byte(conditionexception)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	getWAFMode := appsec.GetWAFModeRequest{}

	getWAFMode.ConfigID = configid
	getWAFMode.Version = version
	getWAFMode.PolicyID = policyid

	wafmode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if wafmode.Mode == AseAuto { // action is read only, only exception is writable
		updateRule := appsec.UpdateConditionExceptionRequest{}
		updateRule.ConfigID = configid
		updateRule.Version = version
		updateRule.PolicyID = policyid
		updateRule.RuleID = ruleid
		var ruleConditionException appsec.RuleConditionException
		err = json.Unmarshal([]byte(rawJSON), &ruleConditionException)
		if err != nil {
			return diag.FromErr(err)
		}
		updateRule.Conditions = ruleConditionException.Conditions
		updateRule.Exception = ruleConditionException.Exception
		resp, err := client.UpdateRuleConditionException(ctx, updateRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
		logger.Debugf("calling 'UpdateRule Response': %s", resp)
	} else {

		action, err := tools.GetStringValue("rule_action", d)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := validateActionAndConditionException(action, conditionexception); err != nil {
			return diag.FromErr(err)
		}

		updateRule := appsec.UpdateRuleRequest{}
		updateRule.ConfigID = configid
		updateRule.Version = version
		updateRule.PolicyID = policyid
		updateRule.RuleID = ruleid
		updateRule.Action = action
		updateRule.JsonPayloadRaw = rawJSON

		_, err = client.UpdateRule(ctx, updateRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
	}

	return resourceRuleRead(ctx, d, m)
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleRemove")
	logger.Debugf("in resourceEvalRuleDelete")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:ruleid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "rule", m)
	policyid := idParts[1]
	ruleid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	getWAFMode := appsec.GetWAFModeRequest{}

	getWAFMode.ConfigID = configid
	getWAFMode.Version = version
	getWAFMode.PolicyID = policyid

	wafmode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if wafmode.Mode == AseAuto {
		updateRule := appsec.UpdateConditionExceptionRequest{}
		updateRule.ConfigID = configid
		updateRule.Version = version
		updateRule.PolicyID = policyid
		updateRule.RuleID = ruleid

		_, err = client.UpdateRuleConditionException(ctx, updateRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
	} else {
		updateRule := appsec.UpdateRuleRequest{}
		updateRule.ConfigID = configid
		updateRule.Version = version
		updateRule.PolicyID = policyid
		updateRule.RuleID = ruleid
		updateRule.Action = "none"

		_, err = client.UpdateRule(ctx, updateRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return nil
}
