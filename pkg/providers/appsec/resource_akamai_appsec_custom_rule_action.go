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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
			"security_policy_id": {
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

	getCustomRuleAction := appsec.GetCustomRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getCustomRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getCustomRuleAction.Version = version

		policyid := s[2]
		getCustomRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getCustomRuleAction.RuleID = ruleid

	} else {
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

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getCustomRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("custom_rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getCustomRuleAction.RuleID = ruleid
	}
	customruleaction, err := client.GetCustomRuleAction(ctx, getCustomRuleAction)
	if err != nil {
		logger.Errorf("calling 'getCustomRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("custom_rule_id", getCustomRuleAction.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", getCustomRuleAction.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getCustomRuleAction.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getCustomRuleAction.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("custom_rule_action", customruleaction.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s:%d", getCustomRuleAction.ConfigID, getCustomRuleAction.Version, getCustomRuleAction.PolicyID, getCustomRuleAction.RuleID))

	return nil
}

func resourceCustomRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionRemove")

	updateCustomRuleAction := appsec.UpdateCustomRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateCustomRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateCustomRuleAction.Version = version

		policyid := s[2]
		updateCustomRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateCustomRuleAction.RuleID = ruleid

	} else {
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

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateCustomRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("custom_rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateCustomRuleAction.RuleID = ruleid
	}
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

	updateCustomRuleAction := appsec.UpdateCustomRuleActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateCustomRuleAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateCustomRuleAction.Version = version

		policyid := s[2]
		updateCustomRuleAction.PolicyID = policyid

		ruleid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateCustomRuleAction.RuleID = ruleid

	} else {
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

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateCustomRuleAction.PolicyID = policyid

		ruleid, err := tools.GetIntValue("custom_rule_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateCustomRuleAction.RuleID = ruleid
	}
	customruleaction, err := tools.GetStringValue("custom_rule_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRuleAction.Action = customruleaction

	_, erru := client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if erru != nil {
		logger.Errorf("calling 'updateCustomRuleAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceCustomRuleActionRead(ctx, d, m)
}

const (
	Alert = "alert"
	Deny  = "deny"
	None  = "none"
)
