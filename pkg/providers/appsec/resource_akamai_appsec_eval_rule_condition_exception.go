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

func resourceEvalRuleConditionExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleConditionExceptionRead")

	getEvalRuleConditionException := appsec.GetEvalRuleConditionExceptionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEvalRuleConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEvalRuleConditionException.Version = version

		policyid := s[2]
		getEvalRuleConditionException.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEvalRuleConditionException.RuleID = ruleid
	} else {
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

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEvalRuleConditionException.RuleID = ruleid
	}
	evalruleconditionexception, err := client.GetEvalRuleConditionException(ctx, getEvalRuleConditionException)
	if err != nil {
		logger.Warnf("calling 'getEvalRuleConditionException': %s", err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "EvalRuleConditionExceptions", evalruleconditionexception)

	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	jsonBody, err := json.Marshal(evalruleconditionexception)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("condition_exception", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("rule_id", getEvalRuleConditionException.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", getEvalRuleConditionException.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getEvalRuleConditionException.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getEvalRuleConditionException.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s:%d", getEvalRuleConditionException.ConfigID, getEvalRuleConditionException.Version, getEvalRuleConditionException.PolicyID, getEvalRuleConditionException.RuleID))

	d.SetId(strconv.Itoa(getEvalRuleConditionException.ConfigID))

	return nil
}

func resourceEvalRuleConditionExceptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleConditionExceptionRemove")

	removeEvalRuleConditionException := appsec.RemoveEvalRuleConditionExceptionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEvalRuleConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEvalRuleConditionException.Version = version

		policyid := s[2]
		removeEvalRuleConditionException.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEvalRuleConditionException.RuleID = ruleid
	} else {
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

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEvalRuleConditionException.RuleID = ruleid
	}

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

	updateEvalRuleConditionException := appsec.UpdateEvalRuleConditionExceptionRequest{}

	jsonpostpayload := d.Get("condition_exception")

	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)
	updateEvalRuleConditionException.JsonPayloadRaw = rawJSON
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEvalRuleConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEvalRuleConditionException.Version = version

		policyid := s[2]
		updateEvalRuleConditionException.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEvalRuleConditionException.RuleID = ruleid
	} else {
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

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEvalRuleConditionException.RuleID = ruleid
	}
	_, erru := client.UpdateEvalRuleConditionException(ctx, updateEvalRuleConditionException)
	if erru != nil {
		logger.Warnf("calling 'updateEvalRuleConditionException': %s", erru.Error())
	}

	return resourceEvalRuleConditionExceptionRead(ctx, d, m)
}
