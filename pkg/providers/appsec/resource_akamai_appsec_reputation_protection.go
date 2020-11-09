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
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceReputationProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProtectionUpdate,
		ReadContext:   resourceReputationProtectionRead,
		UpdateContext: resourceReputationProtectionUpdate,
		DeleteContext: resourceReputationProtectionDelete,
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
			"policy_id": {
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

func resourceReputationProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionRead")

	getReputationProtection := v2.GetReputationProtectionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProtection.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProtection.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProtection.PolicyID = policyid

	reputationprotection, err := client.GetReputationProtection(ctx, getReputationProtection)
	if err != nil {
		logger.Errorf("calling 'getReputationProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "reputationProtectionDS", reputationprotection)
	if err == nil {
		d.Set("output_text", outputtext)
	}
	d.Set("enabled", reputationprotection.ApplyReputationControls)
	d.SetId(strconv.Itoa(getReputationProtection.ConfigID))

	return nil
}

func resourceReputationProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceReputationProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionUpdate")

	updateReputationProtection := v2.UpdateReputationProtectionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationProtection.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationProtection.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationProtection.PolicyID = policyid

	applyreputationcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationProtection.ApplyReputationControls = applyreputationcontrols

	_, erru := client.UpdateReputationProtection(ctx, updateReputationProtection)
	if erru != nil {
		logger.Errorf("calling 'updateReputationProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceReputationProtectionRead(ctx, d, m)
}
