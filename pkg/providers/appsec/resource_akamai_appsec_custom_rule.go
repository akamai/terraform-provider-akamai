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
func resourceCustomRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRuleCreate,
		ReadContext:   resourceCustomRuleRead,
		UpdateContext: resourceCustomRuleUpdate,
		DeleteContext: resourceCustomRuleDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_rule": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"custom_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceCustomRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleCreate")

	createCustomRule := appsec.CreateCustomRuleRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createCustomRule.ConfigID = configid

	jsonpostpayload := d.Get("custom_rule").(string)
	json.Unmarshal([]byte(jsonpostpayload), &createCustomRule)

	customrule, err := client.CreateCustomRule(ctx, createCustomRule)
	if err != nil {
		logger.Errorf("calling 'createCustomRule': %s", err.Error())
		return diag.FromErr(err)
	}

	d.Set("custom_rule_id", customrule.ID)
	d.SetId(strconv.Itoa(customrule.ID))

	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleUpdate")

	updateCustomRule := appsec.UpdateCustomRuleRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateCustomRule.ConfigID = configid

	ID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	updateCustomRule.ID = ID

	jsonpostpayload := d.Get("custom_rule").(string)
	json.Unmarshal([]byte(jsonpostpayload), &updateCustomRule)

	_, erru := client.UpdateCustomRule(ctx, updateCustomRule)
	if erru != nil {
		logger.Errorf("calling 'updateCustomRule': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleRemove")

	removeCustomRule := appsec.RemoveCustomRuleRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeCustomRule.ConfigID = configid

	ID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	removeCustomRule.ID = ID

	_, errd := client.RemoveCustomRule(ctx, removeCustomRule)
	if errd != nil {
		logger.Errorf("calling 'removeCustomRule': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceCustomRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleRead")

	getCustomRule := appsec.GetCustomRuleRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRule.ConfigID = configid

	ID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getCustomRule.ID = ID

	customrule, err := client.GetCustomRule(ctx, getCustomRule)
	if err != nil {
		logger.Errorf("calling 'getCustomRule': %s", err.Error())
		return diag.FromErr(err)
	}

	d.Set("custom_rule_id", customrule.ID)
	d.SetId(strconv.Itoa(customrule.ID))

	return nil
}
