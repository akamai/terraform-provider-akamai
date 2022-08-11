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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceEval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalCreate,
		ReadContext:   resourceEvalRead,
		UpdateContext: resourceEvalUpdate,
		DeleteContext: resourceEvalDelete,
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
			"eval_operation": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					Start,
					Stop,
					Restart,
					Update,
					Complete,
				}, false)),
				Description: "Evaluation mode operation (START, STOP, RESTART, UPDATE or COMPLETE)",
			},
			"eval_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ASE_MANUAL",
					"ASE_AUTO",
				}, false),
				Description: "Evaluation mode (ASE_AUTO or ASE_MANUAL)",
			},
			"current_ruleset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Versioning information for the Kona Rule Set currently in use in production",
			},
			"evaluating_ruleset": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Versioning information for the Kona Rule Set being evaluated",
			},
			"eval_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether an evaluation is currently in progress",
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the evaluation period ends",
			},
		},
	}
}

func resourceEvalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalCreate")
	logger.Debugf(" in resourceEvalCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ruleevaluation", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	evaloperation, err := tools.GetStringValue("eval_operation", d)
	if err != nil {
		return diag.FromErr(err)
	}

	evalmode, err := tools.GetStringValue("eval_mode", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createEval := appsec.UpdateEvalRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Eval:     evaloperation,
		Mode:     evalmode,
	}

	_, err = client.UpdateEval(ctx, createEval)
	if err != nil {
		logger.Errorf("calling 'createEval': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createEval.ConfigID, createEval.PolicyID))

	return resourceEvalRead(ctx, d, m)
}

func resourceEvalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRead")
	logger.Debugf(" in resourceEvalRead")

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

	getEval := appsec.GetEvalRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	eval, err := client.GetEval(ctx, getEval)
	if err != nil {
		logger.Errorf("calling 'getEval': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getEval.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getEval.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("current_ruleset", eval.Current); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("eval_status", eval.Eval); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("evaluating_ruleset", eval.Evaluating); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("expiration_date", eval.Expires); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceEvalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalUpdate")
	logger.Debugf(" in resourceEvalUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ruleevaluation", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	evaloperation, err := tools.GetStringValue("eval_operation", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	evalmode, err := tools.GetStringValue("eval_mode", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateEval := appsec.UpdateEvalRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Eval:     evaloperation,
		Mode:     evalmode,
	}

	_, err = client.UpdateEval(ctx, updateEval)
	if err != nil {
		logger.Errorf("calling 'updateEval': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceEvalRead(ctx, d, m)
}

func resourceEvalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalDelete")
	logger.Debugf(" in resourceEvalDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ruleevaluation", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	removeEval := appsec.RemoveEvalRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Eval:     "STOP",
	}

	_, err = client.RemoveEval(ctx, removeEval)
	if err != nil {
		logger.Errorf("calling 'removeEval': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}

// Definition of constant variables
const (
	Start    = "START"
	Stop     = "STOP"
	Restart  = "RESTART"
	Update   = "UPDATE"
	Complete = "COMPLETE"
)
