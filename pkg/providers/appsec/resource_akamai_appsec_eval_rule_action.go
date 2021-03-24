package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceEvalRuleAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalRuleActionUpdate,
		ReadContext:   resourceEvalRuleActionRead,
		UpdateContext: resourceEvalRuleActionUpdate,
		DeleteContext: resourceEvalRuleActionDelete,
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
			"rule_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceEvalRuleActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleActionRead")

	getEvalRuleAction := appsec.GetEvalRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEvalRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEvalRuleAction.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			getEvalRuleAction.Version = version
		}

		policyid := s[2]
		getEvalRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEvalRuleAction.RuleID = ruleid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEvalRuleAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEvalRuleAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEvalRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEvalRuleAction.RuleID = ruleid
	}
	evalruleaction, err := client.GetEvalRuleAction(ctx, getEvalRuleAction)
	if err != nil {
		logger.Warnf("calling 'getEvalRuleAction': %s", err.Error())
	}

	if err := d.Set("rule_id", getEvalRuleAction.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", getEvalRuleAction.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getEvalRuleAction.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getEvalRuleAction.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("rule_id", getEvalRuleAction.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("rule_action", evalruleaction.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s:%d", getEvalRuleAction.ConfigID, getEvalRuleAction.Version, getEvalRuleAction.PolicyID, getEvalRuleAction.RuleID))

	return nil
}

func resourceEvalRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleActionRemove")

	removeEvalRuleAction := appsec.UpdateEvalRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEvalRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEvalRuleAction.Version = version

		policyid := s[2]
		removeEvalRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEvalRuleAction.RuleID = ruleid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEvalRuleAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEvalRuleAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEvalRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEvalRuleAction.RuleID = ruleid
	}
	removeEvalRuleAction.Action = "none"

	_, erru := client.UpdateEvalRuleAction(ctx, removeEvalRuleAction)
	if erru != nil {
		logger.Warnf("calling 'removeEvalRuleAction': %s", erru.Error())
	}
	d.SetId("")
	return nil
}

func resourceEvalRuleActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRuleActionUpdate")

	updateEvalRuleAction := appsec.UpdateEvalRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEvalRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEvalRuleAction.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			updateEvalRuleAction.Version = version
		}

		policyid := s[2]
		updateEvalRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEvalRuleAction.RuleID = ruleid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEvalRuleAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEvalRuleAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEvalRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEvalRuleAction.RuleID = ruleid
	}
	ruleaction, err := tools.GetStringValue("rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEvalRuleAction.Action = ruleaction

	_, erru := client.UpdateEvalRuleAction(ctx, updateEvalRuleAction)
	if erru != nil {
		logger.Warnf("calling 'updateEvalRuleAction': %s", erru.Error())
	}

	return resourceEvalRuleActionRead(ctx, d, m)
}
