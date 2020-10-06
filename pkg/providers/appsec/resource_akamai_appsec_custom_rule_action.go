package appsec

import (
	"context"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
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
			"rule_id": {
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

	getCustomRuleAction.ConfigID = d.Get("config_id").(int)
	getCustomRuleAction.Version = d.Get("version").(int)
	getCustomRuleAction.PolicyID = d.Get("policy_id").(string)
	getCustomRuleAction.ID = d.Get("rule_id").(int)

	customruleaction, err := client.GetCustomRuleAction(ctx, getCustomRuleAction)
	if err != nil {
		logger.Warnf("calling 'getCustomRuleAction': %s", err.Error())
	}

	d.Set("rule_id", getCustomRuleAction.ID)
	d.SetId(strconv.Itoa(customruleaction.RuleID))

	return nil
}

func resourceCustomRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionRemove")

	updateCustomRuleAction := v2.UpdateCustomRuleActionRequest{}

	updateCustomRuleAction.ConfigID = d.Get("config_id").(int)
	updateCustomRuleAction.Version = d.Get("version").(int)
	updateCustomRuleAction.PolicyID = d.Get("policy_id").(string)
	updateCustomRuleAction.ID = d.Get("rule_id").(int)
	updateCustomRuleAction.Action = "none"

	_, err := client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if err != nil {
		logger.Warnf("calling 'removeCustomRuleAction': %s", err.Error())
	}

	d.SetId("")

	return nil
}

func resourceCustomRuleActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionUpdate")

	updateCustomRuleAction := v2.UpdateCustomRuleActionRequest{}

	updateCustomRuleAction.ConfigID = d.Get("config_id").(int)
	updateCustomRuleAction.Version = d.Get("version").(int)
	updateCustomRuleAction.PolicyID = d.Get("policy_id").(string)
	updateCustomRuleAction.ID = d.Get("rule_id").(int)
	updateCustomRuleAction.Action = d.Get("custom_rule_action").(string)

	_, erru := client.UpdateCustomRuleAction(ctx, updateCustomRuleAction)
	if erru != nil {
		logger.Warnf("calling 'updateCustomRuleAction': %s", erru.Error())
	}

	return resourceCustomRuleActionRead(ctx, d, m)
}
