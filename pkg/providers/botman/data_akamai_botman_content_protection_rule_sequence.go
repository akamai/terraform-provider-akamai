package botman

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContentProtectionRuleSequence() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContentProtectionRuleSequenceRead,
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
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Unique identifiers of content protection rules, listed in the order of their evaluation",
			},
		},
	}
}

func dataSourceContentProtectionRuleSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceContentProtectionRuleSequenceRead")
	logger.Debugf("in dataSourceContentProtectionRuleSequenceRead")

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

	request := botman.GetContentProtectionRuleSequenceRequest{
		ConfigID:         configID,
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
	}

	response, err := client.GetContentProtectionRuleSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GeContentProtectionRuleSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("content_protection_rule_ids", response.ContentProtectionRuleSequence); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.FormatInt(configID, 10))
	return nil
}
