package botman

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceContentProtectionRuleSequence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContentProtectionRuleSequenceCreate,
		ReadContext:   resourceContentProtectionRuleSequenceRead,
		UpdateContext: resourceContentProtectionRuleSequenceUpdate,
		DeleteContext: resourceContentProtectionRuleSequenceDelete,
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
			"content_protection_rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Unique identifiers of content protection rules, listed in the order of their evaluation",
			},
		},
	}
}

func resourceContentProtectionRuleSequenceUpsert(ctx context.Context, d *schema.ResourceData, m interface{}) (int64, diag.Diagnostics) {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionRuleSequenceUpsert")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "ContentProtectionRuleSequence", m)
	if err != nil {
		return configID, diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}

	sequence, err := tf.GetTypedListValue[string]("content_protection_rule_ids", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}
	var stringSequence []string
	stringSequence = append(stringSequence, sequence...)

	request := botman.UpdateContentProtectionRuleSequenceRequest{
		ConfigID:                      configID,
		Version:                       int64(version),
		SecurityPolicyID:              securityPolicyID,
		ContentProtectionRuleSequence: botman.ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: stringSequence},
	}

	_, err = client.UpdateContentProtectionRuleSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateContentProtectionRuleSequence': %s", err.Error())
		return configID, diag.FromErr(err)
	}
	return configID, nil
}

func resourceContentProtectionRuleSequenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configID, diagnostics := resourceContentProtectionRuleSequenceUpsert(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}
	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))
	return resourceContentProtectionRuleSequenceRead(ctx, d, m)
}

func resourceContentProtectionRuleSequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, diagnostics := resourceContentProtectionRuleSequenceUpsert(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}
	return resourceContentProtectionRuleSequenceRead(ctx, d, m)
}

func resourceContentProtectionRuleSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceContentProtectionRuleSequenceRead")

	idParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
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
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetContentProtectionRuleSequenceRequest{
		ConfigID:         configID,
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
	}

	response, err := client.GetContentProtectionRuleSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetContentProtectionRuleSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	fields := map[string]interface{}{
		"config_id":                   configID,
		"security_policy_id":          securityPolicyID,
		"content_protection_rule_ids": response.ContentProtectionRuleSequence,
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceContentProtectionRuleSequenceDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceContentProtectionRuleSequenceDelete")
	logger.Info("Botman API does not support content protection rule sequence deletion - resource will only be removed from state")
	return nil
}
