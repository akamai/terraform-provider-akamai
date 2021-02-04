package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRuleConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRuleConditionException.Version = version

		policyid := s[2]
		getRuleConditionException.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRuleConditionException.RuleID = ruleid
	} else {
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
	}
	ruleconditionexception, err := client.GetRuleConditionException(ctx, getRuleConditionException)
	if err != nil {
		logger.Errorf("calling 'getRuleConditionException': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "RuleConditionExceptions", ruleconditionexception)

	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	jsonBody, err := json.Marshal(ruleconditionexception)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("condition_exception", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", getRuleConditionException.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getRuleConditionException.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getRuleConditionException.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s:%d", getRuleConditionException.ConfigID, getRuleConditionException.Version, getRuleConditionException.PolicyID, getRuleConditionException.RuleID))

	return nil
}

func resourceRuleConditionExceptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleConditionExceptionRemove")

	removeRuleConditionException := appsec.RemoveRuleConditionExceptionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeRuleConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeRuleConditionException.Version = version

		policyid := s[2]
		removeRuleConditionException.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeRuleConditionException.RuleID = ruleid
	} else {
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

	}
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

	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)
	updateRuleConditionException.JsonPayloadRaw = rawJSON
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleConditionException.Version = version

		policyid := s[2]
		updateRuleConditionException.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleConditionException.RuleID = ruleid
	} else {
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
	}
	_, erru := client.UpdateRuleConditionException(ctx, updateRuleConditionException)
	if erru != nil {
		logger.Errorf("calling 'updateRuleConditionException': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceRuleConditionExceptionRead(ctx, d, m)
}
