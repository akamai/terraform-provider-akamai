package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"custom_rule": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
				Description:      "JSON-formatted definition of the custom rule",
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
	logger.Debugf("in resourceCustomRuleCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonpostpayload, err := tools.GetStringValue("custom_rule", d)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonPayloadRaw := []byte(jsonpostpayload)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createCustomRule := appsec.CreateCustomRuleRequest{
		ConfigID:       configID,
		JsonPayloadRaw: rawJSON,
	}

	customrule, err := client.CreateCustomRule(ctx, createCustomRule)
	if err != nil {
		logger.Errorf("calling 'createCustomRule': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("custom_rule_id", customrule.ID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%d", createCustomRule.ConfigID, customrule.ID))

	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleRead")
	logger.Debugf("in resourceCustomRuleRead")

	iDParts, err := splitID(d.Id(), 2, "configID:customRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customRuleID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	getCustomRule := appsec.GetCustomRuleRequest{
		ConfigID: configID,
		ID:       customRuleID,
	}

	customrule, err := client.GetCustomRule(ctx, getCustomRule)
	if err != nil {
		logger.Errorf("calling 'getCustomRule': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getCustomRule.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if err := d.Set("custom_rule_id", getCustomRule.ID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(customrule)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("custom_rule", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	return nil
}

func resourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleUpdate")
	logger.Debugf("in resourceCustomRuleUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:customRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	customRuleID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	jsonpostpayload := d.Get("custom_rule")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateCustomRule := appsec.UpdateCustomRuleRequest{
		ConfigID:       configID,
		ID:             customRuleID,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateCustomRule(ctx, updateCustomRule)
	if err != nil {
		logger.Errorf("calling 'updateCustomRule': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceCustomRuleRead(ctx, d, m)
}

func resourceCustomRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleDelete")
	logger.Debugf("in resourceCustomRuleDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:customRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	customRuleID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	getCustomRules := appsec.GetCustomRulesRequest{
		ConfigID: configID,
		ID:       customRuleID,
	}

	customrules, err := client.GetCustomRules(ctx, getCustomRules)
	if err != nil {
		logger.Errorf("calling 'getCustomRules': %s", err.Error())
		return diag.FromErr(err)
	}

	var status string = customrules.CustomRules[0].Status
	if strings.Compare(status, "unused") == 0 {

		removeCustomRule := appsec.RemoveCustomRuleRequest{
			ConfigID: configID,
			ID:       customRuleID,
		}

		_, err = client.RemoveCustomRule(ctx, removeCustomRule)
		if err != nil {
			logger.Errorf("calling 'removeCustomRule': %s", err.Error())
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("custom rule %d cannot be deleted, it is either active or in use", customRuleID)
	}
	return nil
}
