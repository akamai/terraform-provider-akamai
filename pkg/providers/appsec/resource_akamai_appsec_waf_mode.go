package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceWAFMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWAFModeCreate,
		ReadContext:   resourceWAFModeRead,
		UpdateContext: resourceWAFModeUpdate,
		DeleteContext: resourceWAFModeDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					AAG,
					KRS,
					AseAuto,
					AseManual,
				}, false)),
				Description: "How Kona Rule Set rules should be upgraded (KRS, AAG, ASE_MANUAL or ASE_AUTO)",
			},
			"current_ruleset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Versioning information for the current Kona Rule Set",
			},
			"eval_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether an evaluation is currently in progress",
			},
			"eval_ruleset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Versioning information for the Kona Rule Set being evaluated, if applicable",
			},
			"eval_expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date on which the evaluation period ends, if applicable",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func resourceWAFModeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFModeCreate")
	logger.Debugf(" in resourceWAFModeCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "wafMode", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	mode, err := tools.GetStringValue("mode", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createWAFMode := appsec.UpdateWAFModeRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Mode:     mode,
	}

	_, err = client.UpdateWAFMode(ctx, createWAFMode)
	if err != nil {
		logger.Errorf("calling 'createWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createWAFMode.ConfigID, createWAFMode.PolicyID))

	return resourceWAFModeRead(ctx, d, m)
}

func resourceWAFModeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFModeRead")
	logger.Debugf(" in resourceWAFModeRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	getWAFMode := appsec.GetWAFModeRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	wafMode, err := client.GetWAFMode(ctx, getWAFMode)
	if err != nil {
		logger.Errorf("calling 'getWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getWAFMode.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getWAFMode.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("mode", wafMode.Mode); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("current_ruleset", wafMode.Current); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("eval_status", wafMode.Eval); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("eval_ruleset", wafMode.Evaluating); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("eval_expiration_date", wafMode.Expires); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "wafModesDS", wafMode)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceWAFModeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFModeUpdate")
	logger.Debugf(" in resourceWAFModeUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "wafMode", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	mode, err := tools.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateWAFMode := appsec.UpdateWAFModeRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Mode:     mode,
	}

	_, err = client.UpdateWAFMode(ctx, updateWAFMode)
	if err != nil {
		logger.Errorf("calling 'updateWAFMode': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceWAFModeRead(ctx, d, m)
}

func resourceWAFModeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(ctx, d, m)
}

// Definition of constant variables
const (
	// AAG = Automated Attack Groups
	AAG = "AAG"

	// KRS = Kona Rule Sets
	KRS = "KRS"

	// AseAuto = Adaptive Security Engine - Auto
	AseAuto = "ASE_AUTO"

	// AseManual = Adaptive Security Engine - Manual
	AseManual = "ASE_MANUAL"
)
