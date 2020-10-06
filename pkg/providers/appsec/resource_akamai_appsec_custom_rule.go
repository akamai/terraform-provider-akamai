package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
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
			"rules": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"rule_id": {
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

	createCustomRule := v2.CreateCustomRuleRequest{}

	createCustomRule.ConfigID = d.Get("config_id").(int)
	jsonpostpayload := d.Get("rules").(string)
	json.Unmarshal([]byte(jsonpostpayload), &createCustomRule)

	customrule, err := client.CreateCustomRule(ctx, createCustomRule)
	if err != nil {
		logger.Warnf("calling 'createCustomRule': %s", err.Error())
	}

	d.Set("rule_id", customrule.ID)
	d.SetId(strconv.Itoa(customrule.ID))

	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleUpdate")

	updateCustomRule := v2.UpdateCustomRuleRequest{}

	updateCustomRule.ConfigID = d.Get("config_id").(int)
	updateCustomRule.ID, _ = strconv.Atoi(d.Id())
	jsonpostpayload := d.Get("rules").(string)
	json.Unmarshal([]byte(jsonpostpayload), &updateCustomRule)

	_, err := client.UpdateCustomRule(ctx, updateCustomRule)
	if err != nil {
		logger.Warnf("calling 'updateCustomRule': %s", err.Error())
	}

	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleRemove")

	removeCustomRule := v2.RemoveCustomRuleRequest{}

	removeCustomRule.ConfigID = d.Get("config_id").(int)
	removeCustomRule.ID, _ = strconv.Atoi(d.Id())

	_, err := client.RemoveCustomRule(ctx, removeCustomRule)
	if err != nil {
		logger.Warnf("calling 'removeCustomRule': %s", err.Error())
	}

	d.SetId("")

	return nil
}

func resourceCustomRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleRead")

	getCustomRule := v2.GetCustomRuleRequest{}

	getCustomRule.ConfigID = d.Get("config_id").(int)
	getCustomRule.ID, _ = strconv.Atoi(d.Id())

	customrule, err := client.GetCustomRule(ctx, getCustomRule)
	if err != nil {
		logger.Warnf("calling 'getCustomRule': %s", err.Error())
	}

	d.Set("rule_id", customrule.ID)
	d.SetId(strconv.Itoa(customrule.ID))

	return nil
}
