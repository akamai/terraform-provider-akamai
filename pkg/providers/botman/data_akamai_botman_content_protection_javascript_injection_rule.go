package botman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContentProtectionJavaScriptInjectionRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContentProtectionJavaScriptInjectionRuleRead,
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
				Optional:    true,
				Description: "Unique identifier of a content protection JavaScript injection rule",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
		},
	}
}

func dataSourceContentProtectionJavaScriptInjectionRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceContentProtectionJavaScriptInjectionRuleRead")
	logger.Debugf("in dataSourceContentProtectionJavaScriptInjectionRuleRead")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, int(configID), m)
	if err != nil {
		return diag.FromErr(err)
	}

	contentProtectionJavaScriptInjectionRuleID, err := tf.GetStringValue("content_protection_javascript_injection_rule_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := botman.GetContentProtectionJavaScriptInjectionRuleListRequest{
		ConfigID:         configID,
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		ContentProtectionJavaScriptInjectionRuleID: contentProtectionJavaScriptInjectionRuleID,
	}

	response, err := client.GetContentProtectionJavaScriptInjectionRuleList(ctx, request)

	if err != nil {
		logger.Errorf("calling 'GetContentProtectionJavaScriptInjectionRuleList': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))
	return nil
}
