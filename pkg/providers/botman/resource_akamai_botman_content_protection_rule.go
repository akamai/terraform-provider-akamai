package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceContentProtectionRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentProtectionRuleCreate,
		ReadContext:   resourceContentProtectionRuleRead,
		UpdateContext: resourceContentProtectionRuleUpdate,
		DeleteContext: resourceContentProtectionRuleDelete,
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
			"content_protection_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of a content protection rule",
			},
			"content_protection_rule": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
				Description:      "The content protection rule",
			},
		},
	}
}

func resourceContentProtectionRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionRuleCreateAction")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "ContentProtectionRule", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("content_protection_rule", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateContentProtectionRuleRequest{
		ConfigID:         configID,
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateContentProtectionRule(ctx, request)
	if err != nil {
		logger.Errorf("calling 'CreateContentProtectionRule': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, response["contentProtectionRuleId"]))
	return ContentProtectionRuleRead(ctx, d, m, false)
}

func resourceContentProtectionRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return ContentProtectionRuleRead(ctx, d, m, true)
}

// ContentProtectionRuleRead read content protector rule
func ContentProtectionRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}, readFromCache bool) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionRuleRead")

	idParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:contentProtectionRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(idParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, int(configID), m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := idParts[1]

	contentProtectionRuleID := idParts[2]

	request := botman.GetContentProtectionRuleRequest{
		ConfigID:                configID,
		Version:                 int64(version),
		SecurityPolicyID:        securityPolicyID,
		ContentProtectionRuleID: contentProtectionRuleID,
	}
	var response map[string]interface{}
	if readFromCache {
		response, err = getContentProtectionRule(ctx, request, m)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		response, err = client.GetContentProtectionRule(ctx, request)
		if err != nil {
			logger.Errorf("calling 'GetContentProtectionRule': %s", err.Error())
			return diag.FromErr(err)
		}
	}

	// Removing contentProtectionRuleId from response to suppress diff
	delete(response, "contentProtectionRuleId")

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", securityPolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("content_protection_rule_id", contentProtectionRuleID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                  configID,
		"security_policy_id":         securityPolicyID,
		"content_protection_rule_id": contentProtectionRuleID,
		"content_protection_rule":    string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceContentProtectionRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionRuleUpdateAction")

	idParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:contentProtectionRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(idParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "ContentProtectionRule", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := idParts[1]

	contentProtectionRuleID := idParts[2]

	jsonPayload, err := getJSONPayload(d, "content_protection_rule", "contentProtectionRuleId", contentProtectionRuleID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateContentProtectionRuleRequest{
		ConfigID:                configID,
		Version:                 int64(version),
		ContentProtectionRuleID: contentProtectionRuleID,
		SecurityPolicyID:        securityPolicyID,
		JsonPayload:             jsonPayload,
	}

	_, err = client.UpdateContentProtectionRule(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateContentProtectionRule': %s", err.Error())
		return diag.FromErr(err)
	}
	return ContentProtectionRuleRead(ctx, d, m, false)
}

func resourceContentProtectionRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionRuleDeleteAction")

	idParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:contentProtectionRuleID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(idParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "ContentProtectionRule", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := idParts[1]

	contentProtectionRuleID := idParts[2]

	request := botman.RemoveContentProtectionRuleRequest{
		ConfigID:                configID,
		Version:                 int64(version),
		SecurityPolicyID:        securityPolicyID,
		ContentProtectionRuleID: contentProtectionRuleID,
	}

	err = client.RemoveContentProtectionRule(ctx, request)
	if err != nil {
		logger.Errorf("calling 'RemoveContentProtectionRule': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
