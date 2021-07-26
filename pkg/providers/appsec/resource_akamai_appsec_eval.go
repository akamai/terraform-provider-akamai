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
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
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
			"eval_operation": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Start,
					Stop,
					Restart,
					Update,
					Complete,
				}, false),
			},
			"current_ruleset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"evaluating_ruleset": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eval_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEvalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalCreate")
	logger.Debugf("!!! in resourceEvalCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ruleevaluation", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	evaloperation, err := tools.GetStringValue("eval_operation", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createEval := appsec.UpdateEvalRequest{}
	createEval.ConfigID = configid
	createEval.Version = version
	createEval.PolicyID = policyid
	createEval.Eval = evaloperation

	_, erru := client.UpdateEval(ctx, createEval)
	if erru != nil {
		logger.Errorf("calling 'createEval': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createEval.ConfigID, createEval.PolicyID))

	return resourceEvalRead(ctx, d, m)
}

func resourceEvalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRead")
	logger.Debugf("!!! in resourceEvalRead")

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

	getEval := appsec.GetEvalRequest{}
	getEval.ConfigID = configid
	getEval.Version = version
	getEval.PolicyID = policyid

	eval, err := client.GetEval(ctx, getEval)
	if err != nil {
		logger.Errorf("calling 'getEval': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getEval.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getEval.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("current_ruleset", eval.Current); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("eval_status", eval.Eval); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("evaluating_ruleset", eval.Evaluating); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("expiration_date", eval.Expires); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceEvalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalUpdate")
	logger.Debugf("!!! in resourceEvalUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ruleevaluation", m)
	policyid := idParts[1]
	evaloperation, err := tools.GetStringValue("eval_operation", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateEval := appsec.UpdateEvalRequest{}
	updateEval.ConfigID = configid
	updateEval.Version = version
	updateEval.PolicyID = policyid
	updateEval.Eval = evaloperation

	_, erru := client.UpdateEval(ctx, updateEval)
	if erru != nil {
		logger.Errorf("calling 'updateEval': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceEvalRead(ctx, d, m)
}

func resourceEvalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalDelete")
	logger.Debugf("!!! in resourceEvalDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ruleevaluation", m)
	policyid := idParts[1]

	removeEval := appsec.RemoveEvalRequest{}
	removeEval.ConfigID = configid
	removeEval.Version = version
	removeEval.PolicyID = policyid
	removeEval.Eval = "STOP"

	_, erru := client.RemoveEval(ctx, removeEval)
	if erru != nil {
		logger.Errorf("calling 'removeEval': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
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
