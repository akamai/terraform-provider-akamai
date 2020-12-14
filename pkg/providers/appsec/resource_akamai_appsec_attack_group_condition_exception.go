package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

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
func resourceAttackGroupConditionException() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttackGroupConditionExceptionUpdate,
		ReadContext:   resourceAttackGroupConditionExceptionRead,
		UpdateContext: resourceAttackGroupConditionExceptionUpdate,
		DeleteContext: resourceAttackGroupConditionExceptionDelete,
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
			"attack_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition_exception": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceAttackGroupConditionExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupConditionExceptionRead")

	getAttackGroupConditionException := appsec.GetAttackGroupConditionExceptionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.PolicyID = policyid

	attackgroup, err := tools.GetStringValue("attack_group", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.Group = attackgroup

	attackgroupconditionexception, err := client.GetAttackGroupConditionException(ctx, getAttackGroupConditionException)
	if err != nil {
		logger.Errorf("calling 'getAttackGroupConditionException': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "AttackGroupConditionExceptions", attackgroupconditionexception)

	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getAttackGroupConditionException.ConfigID))

	return nil
}

func resourceAttackGroupConditionExceptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupConditionExceptionRemove")

	removeAttackGroupConditionException := appsec.RemoveAttackGroupConditionExceptionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeAttackGroupConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeAttackGroupConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeAttackGroupConditionException.PolicyID = policyid

	attackgroup, err := tools.GetStringValue("attack_group", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeAttackGroupConditionException.Group = attackgroup

	logger.Errorf("calling 'RemoveAttackGroupConditionException': %v", removeAttackGroupConditionException)

	_, errd := client.RemoveAttackGroupConditionException(ctx, removeAttackGroupConditionException)
	if errd != nil {
		logger.Errorf("calling 'RemoveAttackGroupConditionException': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")
	return nil
}

func resourceAttackGroupConditionExceptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupConditionExceptionUpdate")

	updateAttackGroupConditionException := appsec.UpdateAttackGroupConditionExceptionRequest{}

	jsonpostpayload := d.Get("condition_exception")

	json.Unmarshal([]byte(jsonpostpayload.(string)), &updateAttackGroupConditionException)

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAttackGroupConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAttackGroupConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAttackGroupConditionException.PolicyID = policyid

	attackgroup, err := tools.GetStringValue("attack_group", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAttackGroupConditionException.Group = attackgroup

	resp, erru := client.UpdateAttackGroupConditionException(ctx, updateAttackGroupConditionException)
	if erru != nil {
		logger.Errorf("calling 'updateAttackGroupConditionException': %s", erru.Error())
		return diag.FromErr(erru)
	}
	logger.Warnf("calling 'updateAttackGroupConditionException': %s", resp)
	return resourceAttackGroupConditionExceptionRead(ctx, d, m)
}
