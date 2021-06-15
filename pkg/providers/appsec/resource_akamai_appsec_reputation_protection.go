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
func resourceReputationProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProtectionCreate,
		ReadContext:   resourceReputationProtectionRead,
		UpdateContext: resourceReputationProtectionUpdate,
		DeleteContext: resourceReputationProtectionDelete,
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"enabled": {
				Type:     schema.TypeBool,
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

func resourceReputationProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionCreate")
	logger.Debugf("!!! in resourceReputationProtectionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProtection", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	applyreputationcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createReputationProtection := appsec.UpdateReputationProtectionRequest{}
	createReputationProtection.ConfigID = configid
	createReputationProtection.Version = version
	createReputationProtection.PolicyID = policyid
	createReputationProtection.ApplyReputationControls = applyreputationcontrols

	_, erru := client.UpdateReputationProtection(ctx, createReputationProtection)
	if erru != nil {
		logger.Errorf("calling 'updateReputationProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createReputationProtection.ConfigID, createReputationProtection.PolicyID))

	return resourceReputationProtectionRead(ctx, d, m)
}

func resourceReputationProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionRead")
	logger.Debugf("!!! in resourceReputationProtectionRead")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]

	getReputationProtection := appsec.GetReputationProtectionRequest{}
	getReputationProtection.ConfigID = configid
	getReputationProtection.Version = version
	getReputationProtection.PolicyID = policyid

	reputationprotection, err := client.GetReputationProtection(ctx, getReputationProtection)
	if err != nil {
		logger.Errorf("calling 'getReputationProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getReputationProtection.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getReputationProtection.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enabled", reputationprotection.ApplyReputationControls); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "reputationProtectionDS", reputationprotection)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	return nil
}

func resourceReputationProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionUpdate")
	logger.Debugf("!!! in resourceSlowPostProtectionUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProtection", m)
	policyid := idParts[1]
	applyreputationcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateReputationProtection := appsec.UpdateReputationProtectionRequest{}
	updateReputationProtection.ConfigID = configid
	updateReputationProtection.Version = version
	updateReputationProtection.PolicyID = policyid
	updateReputationProtection.ApplyReputationControls = applyreputationcontrols

	_, erru := client.UpdateReputationProtection(ctx, updateReputationProtection)
	if erru != nil {
		logger.Errorf("calling 'updateReputationProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceReputationProtectionRead(ctx, d, m)
}

func resourceReputationProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionDelete")
	logger.Debugf("!!! in resourceReputationProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProtection", m)
	policyid := idParts[1]

	removeReputationProtection := appsec.RemoveReputationProtectionRequest{}
	removeReputationProtection.ConfigID = configid
	removeReputationProtection.Version = version
	removeReputationProtection.PolicyID = policyid
	removeReputationProtection.ApplyReputationControls = false

	_, errd := client.RemoveReputationProtection(ctx, removeReputationProtection)
	if errd != nil {
		logger.Errorf("calling 'removeReputationProtection': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")
	return nil
}
