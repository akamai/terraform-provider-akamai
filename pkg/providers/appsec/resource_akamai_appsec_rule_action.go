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
func resourceRuleAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleActionUpdate,
		ReadContext:   resourceRuleActionRead,
		UpdateContext: resourceRuleActionUpdate,
		DeleteContext: resourceRuleActionDelete,
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

func resourceRuleActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleActionRead")

	getRuleAction := appsec.GetRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRuleAction.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			getRuleAction.Version = version
		}

		policyid := s[2]
		getRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRuleAction.RuleID = ruleid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRuleAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRuleAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRuleAction.RuleID = ruleid
	}
	ruleaction, err := client.GetRuleAction(ctx, getRuleAction)
	if err != nil {
		logger.Errorf("calling 'getRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getRuleAction.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getRuleAction.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getRuleAction.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("rule_id", getRuleAction.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("rule_action", ruleaction.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(fmt.Sprintf("%d:%d:%s:%d", getRuleAction.ConfigID, getRuleAction.Version, getRuleAction.PolicyID, getRuleAction.RuleID))

	return nil
}

func resourceRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleActionRemove")

	updateRuleAction := appsec.UpdateRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleAction.Version = version

		policyid := s[2]
		updateRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleAction.RuleID = ruleid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.RuleID = ruleid
	}
	updateRuleAction.Action = "none"

	_, erru := client.UpdateRuleAction(ctx, updateRuleAction)
	if erru != nil {
		logger.Errorf("calling 'removeRuleAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}

func resourceRuleActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleActionUpdate")

	updateRuleAction := appsec.UpdateRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleAction.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			updateRuleAction.Version = version
		}
		policyid := s[2]
		updateRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRuleAction.RuleID = ruleid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleAction.RuleID = ruleid
	}
	ruleaction, err := tools.GetStringValue("rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleAction.Action = ruleaction

	_, erru := client.UpdateRuleAction(ctx, updateRuleAction)
	if erru != nil {
		logger.Errorf("calling 'updateRuleAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceRuleActionRead(ctx, d, m)
}
