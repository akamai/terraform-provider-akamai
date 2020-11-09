package appsec

import (
	"context"
	"errors"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceKRSRuleAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKRSRuleActionUpdate,
		ReadContext:   resourceKRSRuleActionRead,
		UpdateContext: resourceKRSRuleActionUpdate,
		DeleteContext: resourceKRSRuleActionDelete,
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
			"krs_rule_action": {
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
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceKRSRuleActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceKRSRuleActionRead")

	getKRSRuleAction := v2.GetKRSRuleActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getKRSRuleAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getKRSRuleAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getKRSRuleAction.PolicyID = policyid

	ruleid, err := tools.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getKRSRuleAction.RuleID = ruleid

	krsruleaction, err := client.GetKRSRuleAction(ctx, getKRSRuleAction)
	if err != nil {
		logger.Errorf("calling 'getKRSRuleAction': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "KRSRuleAction", krsruleaction)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getKRSRuleAction.ConfigID))

	return nil
}

func resourceKRSRuleActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceKRSRuleActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceKRSRuleActionUpdate")

	updateKRSRuleAction := v2.UpdateKRSRuleActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateKRSRuleAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateKRSRuleAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateKRSRuleAction.PolicyID = policyid

	ruleid, err := tools.GetIntValue("rule_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateKRSRuleAction.RuleID = ruleid

	_, erru := client.UpdateKRSRuleAction(ctx, updateKRSRuleAction)
	if erru != nil {
		logger.Errorf("calling 'updateKRSRuleAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceKRSRuleActionRead(ctx, d, m)
}
