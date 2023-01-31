package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	// getWAFModeMutex enforces single-thread access to the GetWAFMode call
	getWAFModeMutex sync.Mutex
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Required:    true,
				Description: "Unique identifier of the rule",
			},
			"rule_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: ValidateActions,
				Description:      "Action to be taken when the rule is triggered",
			},
			"condition_exception": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsConditionException,
				Description:      "JSON-formatted condition and exception information for the rule",
			},
		},
	}
}

func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleCreate")
	logger.Debugf("in resourceRuleCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "rule", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ruleID, err := tools.GetIntValue("rule_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	conditionexception, err := tools.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	jsonPayloadRaw := []byte(conditionexception)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	wafMode, err := getWAFMode(ctx, m, configID, version, policyID)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if wafMode == AseAuto { // action is read only, only condition exception is writable
		ruleConditionException := appsec.RuleConditionException{}
		if conditionexception != "" {
			err = json.Unmarshal([]byte(rawJSON), &ruleConditionException)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		createRule := appsec.UpdateConditionExceptionRequest{
			ConfigID:               configID,
			Version:                version,
			PolicyID:               policyID,
			RuleID:                 ruleID,
			Conditions:             ruleConditionException.Conditions,
			Exception:              ruleConditionException.Exception,
			AdvancedExceptionsList: ruleConditionException.AdvancedExceptionsList,
		}

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

		createRule := appsec.UpdateRuleRequest{
			ConfigID:       configID,
			Version:        version,
			PolicyID:       policyID,
			RuleID:         ruleID,
			Action:         action,
			JsonPayloadRaw: rawJSON,
		}

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

func getWAFMode(ctx context.Context, m interface{}, configID int, version int, policyID string) (string, error) {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "getWAFMode")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getWAFMode", configID, version, policyID)
	getWAFModeResponse := &appsec.GetWAFModeResponse{}
	if err := meta.CacheGet(inst, cacheKey, getWAFModeResponse); err == nil {
		logger.Debugf("returning wafMode %s for config/version/policy %d/%d/%s",
			getWAFModeResponse.Mode, configID, version, policyID)
		return getWAFModeResponse.Mode, nil
	}

	logger.Debugf("requesting getWAFMode mutex lock")
	getWAFModeMutex.Lock()
	defer func() {
		logger.Debugf("releasing getWAFMode mutex lock")
		getWAFModeMutex.Unlock()
	}()

	err := meta.CacheGet(inst, cacheKey, getWAFModeResponse)
	if err == nil {
		logger.Debugf("returning wafMode %s for config/version/policy %d/%d/%s",
			getWAFModeResponse.Mode, configID, version, policyID)
		return getWAFModeResponse.Mode, nil
	}
	// Any error response other than 'not found' or 'cache disabled' is a problem.
	if !akamai.IsNotFoundError(err) && !errors.Is(err, akamai.ErrCacheDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return "", err
	}

	getWAFModeRequest := appsec.GetWAFModeRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	wafMode, err := client.GetWAFMode(ctx, getWAFModeRequest)
	if err != nil {
		logger.Errorf("calling 'GetWAFMode': %s", err.Error())
		return "", err
	}
	if err := meta.CacheSet(inst, cacheKey, wafMode); err != nil {
		if !errors.Is(err, akamai.ErrCacheDisabled) {
			logger.Errorf("error caching WAFMode: %s", err.Error())
		}
	}
	return wafMode.Mode, nil
}

func resourceRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleRead")
	logger.Debugf("in resourceRuleRead")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:ruleID")
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

	getRule := appsec.GetRuleRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		RuleID:   ruleID,
	}

	rule, err := client.GetRule(ctx, getRule)
	if err != nil {
		logger.Errorf("calling 'GetRule': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getRule.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getRule.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("rule_id", getRule.RuleID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if err := d.Set("rule_action", rule.Action); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if !rule.IsEmptyConditionException() {
		jsonBody, err := json.Marshal(rule.ConditionException)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("condition_exception", string(jsonBody)); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	return nil
}

func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleUpdate")
	logger.Debugf("in resourceRuleUpdate")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:ruleID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	version, err := getModifiableConfigVersion(ctx, configID, "rule", m)
	if err != nil {
		return diag.FromErr(err)
	}
	ruleID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}
	conditionexception, err := tools.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	jsonPayloadRaw := []byte(conditionexception)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	wafMode, err := getWAFMode(ctx, m, configID, version, policyID)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if wafMode == AseAuto { // action is read only, only exception is writable
		ruleConditionException := appsec.RuleConditionException{}
		if conditionexception != "" {
			err = json.Unmarshal([]byte(rawJSON), &ruleConditionException)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		updateRule := appsec.UpdateConditionExceptionRequest{
			ConfigID:               configID,
			Version:                version,
			PolicyID:               policyID,
			RuleID:                 ruleID,
			Conditions:             ruleConditionException.Conditions,
			Exception:              ruleConditionException.Exception,
			AdvancedExceptionsList: ruleConditionException.AdvancedExceptionsList,
		}

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

		updateRule := appsec.UpdateRuleRequest{
			ConfigID:       configID,
			Version:        version,
			PolicyID:       policyID,
			RuleID:         ruleID,
			Action:         action,
			JsonPayloadRaw: rawJSON,
		}

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
	logger := meta.Log("APPSEC", "resourceRuleDelete")
	logger.Debugf("in resourceRuleDelete")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:ruleID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "rule", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	ruleID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	wafMode, err := getWAFMode(ctx, m, configID, version, policyID)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if wafMode == AseAuto {
		updateRule := appsec.UpdateConditionExceptionRequest{
			ConfigID: configID,
			Version:  version,
			PolicyID: policyID,
			RuleID:   ruleID,
		}

		_, err = client.UpdateRuleConditionException(ctx, updateRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
	} else {
		updateRule := appsec.UpdateRuleRequest{
			ConfigID: configID,
			Version:  version,
			PolicyID: policyID,
			RuleID:   ruleID,
			Action:   "none",
		}
		_, err = client.UpdateRule(ctx, updateRule)
		if err != nil {
			logger.Errorf("calling 'UpdateRule': %s", err.Error())
			return diag.FromErr(err)
		}
	}
	return nil
}
