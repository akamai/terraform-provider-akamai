package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceContentProtectionJavaScriptInjectionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentProtectionJavaScriptInjectionRuleCreate,
		ReadContext:   resourceContentProtectionJavaScriptInjectionRuleRead,
		UpdateContext: resourceContentProtectionJavaScriptInjectionRuleUpdate,
		DeleteContext: resourceContentProtectionJavaScriptInjectionRuleDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
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
			"content_protection_javascript_injection_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of a content protection JavaScript injection rule",
			},
			"content_protection_javascript_injection_rule": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
				Description:      "The content protection JavaScript injection rule",
			},
		},
	}
}

func resourceContentProtectionJavaScriptInjectionRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionJavaScriptInjectionRuleCreateAction")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "ContentProtectionJavaScriptInjectionRule", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("content_protection_javascript_injection_rule", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateContentProtectionJavaScriptInjectionRuleRequest{
		ConfigID:         configID,
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateContentProtectionJavaScriptInjectionRule(ctx, request)
	if err != nil {
		logger.Errorf("calling 'CreateContentProtectionJavaScriptInjectionRule': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, response["contentProtectionJavaScriptInjectionRuleId"]))
	return ContentProtectionJavaScriptInjectionRuleRead(ctx, d, m, false)
}

func resourceContentProtectionJavaScriptInjectionRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return ContentProtectionJavaScriptInjectionRuleRead(ctx, d, m, true)
}

// ContentProtectionJavaScriptInjectionRuleRead read JavaScript injection rule for content protector
func ContentProtectionJavaScriptInjectionRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}, readFromCache bool) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionJavaScriptInjectionRuleRead")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:contentProtectionJavaScriptInjectionRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, int(configID), m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	contentProtectionJavaScriptInjectionRuleID := iDParts[2]

	request := botman.GetContentProtectionJavaScriptInjectionRuleRequest{
		ConfigID:         configID,
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		ContentProtectionJavaScriptInjectionRuleID: contentProtectionJavaScriptInjectionRuleID,
	}
	var response map[string]interface{}
	if readFromCache {
		response, err = getContentProtectionJavaScriptInjectionRule(ctx, request, m)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		response, err = client.GetContentProtectionJavaScriptInjectionRule(ctx, request)
		if err != nil {
			logger.Errorf("calling 'GetContentProtectionJavaScriptInjectionRule': %s", err.Error())
			return diag.FromErr(err)
		}
	}

	// Removing contentProtectionJavaScriptInjectionRuleId from response to suppress diff
	delete(response, "contentProtectionJavaScriptInjectionRuleId")

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", securityPolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("content_protection_javascript_injection_rule_id", contentProtectionJavaScriptInjectionRuleID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":          configID,
		"security_policy_id": securityPolicyID,
		"content_protection_javascript_injection_rule_id": contentProtectionJavaScriptInjectionRuleID,
		"content_protection_javascript_injection_rule":    string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceContentProtectionJavaScriptInjectionRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionJavaScriptInjectionRuleUpdateAction")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:contentProtectionJavaScriptInjectionRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "ContentProtectionJavaScriptInjectionRule", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	contentProtectionJavaScriptInjectionRuleID := iDParts[2]

	jsonPayload, err := getJSONPayload(d, "content_protection_javascript_injection_rule", "contentProtectionJavaScriptInjectionRuleId", contentProtectionJavaScriptInjectionRuleID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateContentProtectionJavaScriptInjectionRuleRequest{
		ConfigID: configID,
		Version:  int64(version),
		ContentProtectionJavaScriptInjectionRuleID: contentProtectionJavaScriptInjectionRuleID,
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateContentProtectionJavaScriptInjectionRule(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateContentProtectionJavaScriptInjectionRule': %s", err.Error())
		return diag.FromErr(err)
	}
	return ContentProtectionJavaScriptInjectionRuleRead(ctx, d, m, false)
}

func resourceContentProtectionJavaScriptInjectionRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionJavaScriptInjectionRuleDeleteAction")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:contentProtectionJavaScriptInjectionRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "ContentProtectionJavaScriptInjectionRule", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	contentProtectionJavaScriptInjectionRuleID := iDParts[2]

	request := botman.RemoveContentProtectionJavaScriptInjectionRuleRequest{
		ConfigID:         configID,
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		ContentProtectionJavaScriptInjectionRuleID: contentProtectionJavaScriptInjectionRuleID,
	}

	err = client.RemoveContentProtectionJavaScriptInjectionRule(ctx, request)
	if err != nil {
		logger.Errorf("calling 'RemoveContentProtectionJavaScriptInjectionRule': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
