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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceEval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalUpdate,
		ReadContext:   resourceEvalRead,
		UpdateContext: resourceEvalUpdate,
		DeleteContext: resourceEvalDelete,
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
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceEvalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalRead")

	getEval := appsec.GetEvalRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEval.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getEval.Version = version

		policyid := s[2]
		getEval.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEval.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEval.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEval.PolicyID = policyid

		evaloperation, err := tools.GetStringValue("eval_operation", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getEval.Eval = evaloperation
	}
	eval, err := client.GetEval(ctx, getEval)
	if err != nil {
		logger.Errorf("calling 'getEval': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "EvalDS", eval)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
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

	if err := d.Set("config_id", getEval.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getEval.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getEval.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", getEval.ConfigID, getEval.Version, getEval.PolicyID))

	return nil
}

func resourceEvalDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalDelete")

	removeEval := appsec.RemoveEvalRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEval.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeEval.Version = version

		policyid := s[2]
		removeEval.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEval.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEval.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeEval.PolicyID = policyid
	}
	removeEval.Eval = "STOP"

	_, erru := client.RemoveEval(ctx, removeEval)
	if erru != nil {
		logger.Errorf("calling 'removeEval': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}

func resourceEvalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalUpdate")

	updateEval := appsec.UpdateEvalRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEval.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateEval.Version = version

		policyid := s[2]
		updateEval.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEval.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEval.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateEval.PolicyID = policyid
	}
	evaloperation, err := tools.GetStringValue("eval_operation", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateEval.Eval = evaloperation

	_, erru := client.UpdateEval(ctx, updateEval)
	if erru != nil {
		logger.Errorf("calling 'updateEval': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceEvalRead(ctx, d, m)
}

const (
	Start    = "START"
	Stop     = "STOP"
	Restart  = "RESTART"
	Update   = "UPDATE"
	Complete = "COMPLETE"
)
