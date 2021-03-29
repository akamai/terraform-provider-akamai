package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceRuleUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleUpgradeUpdate,
		ReadContext:   resourceRuleUpgradeRead,
		UpdateContext: resourceRuleUpgradeUpdate,
		DeleteContext: resourceRuleUpgradeDelete,
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
			"current_ruleset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eval_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRuleUpgradeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleUpgradeRead")

	getRuleUpgrade := appsec.GetRuleUpgradeRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRuleUpgrade.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRuleUpgrade.Version = version

	if d.HasChange("version") {
		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRuleUpgrade.Version = version
	}

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRuleUpgrade.PolicyID = policyid

	_, errr := client.GetRuleUpgrade(ctx, getRuleUpgrade)
	if errr != nil {
		logger.Errorf("calling 'getRuleUpgrade': %s", errr.Error())
		return diag.FromErr(errr)
	}

	d.SetId(strconv.Itoa(getRuleUpgrade.ConfigID))

	return nil
}

func resourceRuleUpgradeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceRuleUpgradeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleUpgradeUpdate")

	updateRuleUpgrade := appsec.UpdateRuleUpgradeRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleUpgrade.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleUpgrade.Version = version

	if d.HasChange("version") {
		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRuleUpgrade.Version = version
	}

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRuleUpgrade.PolicyID = policyid

	updateRuleUpgrade.Upgrade = true

	ruleupgrade, erru := client.UpdateRuleUpgrade(ctx, updateRuleUpgrade)
	if erru != nil {
		logger.Errorf("calling 'updateRuleUpgrade': %s", erru.Error())
		return diag.FromErr(erru)
	}

	if err := d.Set("current_ruleset", ruleupgrade.Current); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("mode", ruleupgrade.Mode); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("eval_status", ruleupgrade.Eval); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(updateRuleUpgrade.ConfigID))

	return resourceRuleUpgradeRead(ctx, d, m)
}
