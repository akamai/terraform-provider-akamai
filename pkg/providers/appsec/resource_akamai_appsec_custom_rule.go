package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
func resourceCustomRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRuleCreate,
		ReadContext:   resourceCustomRuleRead,
		UpdateContext: resourceCustomRuleUpdate,
		DeleteContext: resourceCustomRuleDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"custom_rule": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceCustomRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleCreate")
	logger.Debugf("in resourceCustomRuleCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	jsonpostpayload := d.Get("custom_rule")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createCustomRule := appsec.CreateCustomRuleRequest{}
	createCustomRule.ConfigID = configid
	createCustomRule.JsonPayloadRaw = rawJSON

	customrule, err := client.CreateCustomRule(ctx, createCustomRule)
	if err != nil {
		logger.Errorf("calling 'createCustomRule': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("custom_rule_id", customrule.ID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d", createCustomRule.ConfigID, customrule.ID))

	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleRead")
	logger.Debugf("in resourceCustomRuleRead")

	idParts, err := splitID(d.Id(), 2, "configid:custom_rule_id")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customruleid, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	getCustomRule := appsec.GetCustomRuleRequest{}
	getCustomRule.ConfigID = configid
	getCustomRule.ID = customruleid
	customrule, err := client.GetCustomRule(ctx, getCustomRule)
	if err != nil {
		logger.Errorf("calling 'getCustomRule': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getCustomRule.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("custom_rule_id", getCustomRule.ID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	jsonBody, err := json.Marshal(customrule)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("custom_rule", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	return nil
}

func resourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleUpdate")
	logger.Debugf("in resourceCustomRuleUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:custom_rule_id")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	customruleid, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	jsonpostpayload := d.Get("custom_rule")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateCustomRule := appsec.UpdateCustomRuleRequest{}
	updateCustomRule.ConfigID = configid
	updateCustomRule.ID = customruleid
	updateCustomRule.JsonPayloadRaw = rawJSON

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
	logger := meta.Log("APPSEC", "resourceCustomRuleDelete")
	logger.Debugf("in resourceCustomRuleDelete")

	idParts, err := splitID(d.Id(), 2, "configid:custom_rule_id")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	customruleid, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	getCustomRules := appsec.GetCustomRulesRequest{}
	getCustomRules.ConfigID = configid
	getCustomRules.ID = customruleid

	customrules, err := client.GetCustomRules(ctx, getCustomRules)
	if err != nil {
		logger.Errorf("calling 'getCustomRules': %s", err.Error())
		return diag.FromErr(err)
	}

	var status string = customrules.CustomRules[0].Status
	if strings.Compare(status, "unused") == 0 {

		removeCustomRule := appsec.RemoveCustomRuleRequest{}
		removeCustomRule.ConfigID = configid
		removeCustomRule.ID = customruleid

		_, errd := client.RemoveCustomRule(ctx, removeCustomRule)
		if errd != nil {
			logger.Errorf("calling 'removeCustomRule': %s", errd.Error())
			return diag.FromErr(errd)
		}
		d.SetId("")
	} else {
		return diag.Errorf("custom rule %d cannot be deleted, it is either active or in use", customruleid)
	}
	return nil
}
