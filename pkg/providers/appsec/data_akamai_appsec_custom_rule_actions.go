package appsec

import (
	"context"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomRuleActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomRuleActionsRead,
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
			"custom_rule_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of the custom rule for which to return information",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceCustomRuleActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceCustomRuleActionsRead")

	getCustomRuleActions := appsec.GetCustomRuleActionsRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getCustomRuleActions.ConfigID = configID

	if getCustomRuleActions.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getCustomRuleActions.PolicyID = policyID

	customRuleID, err := tf.GetIntValue("custom_rule_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomRuleActions.RuleID = customRuleID

	customruleactions, err := client.GetCustomRuleActions(ctx, getCustomRuleActions)
	if err != nil {
		logger.Errorf("calling 'getCustomRuleActions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "customRuleAction", customruleactions)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getCustomRuleActions.ConfigID))

	return nil
}
