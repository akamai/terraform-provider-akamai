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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceRuleUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleUpgradeCreate,
		ReadContext:   resourceRuleUpgradeRead,
		UpdateContext: resourceRuleUpgradeUpdate,
		DeleteContext: resourceRuleUpgradeDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_ruleset": {
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

func resourceRuleUpgradeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleUpgradeCreate")
	logger.Debugf("!!! in resourceRuleUpgradeCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "krsRuleUgrade", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createRuleUpgrade := appsec.UpdateRuleUpgradeRequest{}
	createRuleUpgrade.ConfigID = configid
	createRuleUpgrade.Version = version
	createRuleUpgrade.PolicyID = policyid
	createRuleUpgrade.Upgrade = true

	_, err = client.UpdateRuleUpgrade(ctx, createRuleUpgrade)
	if err != nil {
		logger.Errorf("calling 'createRuleUpgrade': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createRuleUpgrade.ConfigID, createRuleUpgrade.PolicyID))

	return resourceRuleUpgradeRead(ctx, d, m)
}

func resourceRuleUpgradeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleUpgradeRead")
	logger.Debugf("!!! in resourceRuleUpgradeRead")

	idParts, err := splitID(d.Id(), 2, "configid:policyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]

	getWAFMode := appsec.GetWAFModeRequest{}
	getWAFMode.ConfigID = configid
	getWAFMode.Version = version
	getWAFMode.PolicyID = policyid

	wafmode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getWAFMode.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getWAFMode.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("mode", wafmode.Mode); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("current_ruleset", wafmode.Current); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("eval_status", wafmode.Eval); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceRuleUpgradeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRuleUpgradeUpdate")
	logger.Debugf("!!! in resourceRuleUpgradeUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:policyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "securityPolicyRename", m)
	policyid := idParts[1]

	updateRuleUpgrade := appsec.UpdateRuleUpgradeRequest{}
	updateRuleUpgrade.ConfigID = configid
	updateRuleUpgrade.Version = version
	updateRuleUpgrade.PolicyID = policyid
	updateRuleUpgrade.Upgrade = true

	_, err = client.UpdateRuleUpgrade(ctx, updateRuleUpgrade)
	if err != nil {
		logger.Errorf("calling 'updateRuleUpgrade': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceRuleUpgradeRead(ctx, d, m)
}

func resourceRuleUpgradeDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}
