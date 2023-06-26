package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEvalGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvalGroupsRead,
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
			"attack_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique idenfier of the evaluation group for which to retrieve information",
			},
			"attack_group_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action to be taken for the evaluation group if one was specified",
			},
			"condition_exception": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON-formatted condition and exception information for the evaluation group if one was specified",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceEvalGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceEvalGroupsRead")

	getAttackGroups := appsec.GetAttackGroupsRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAttackGroups.ConfigID = configID

	if getAttackGroups.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAttackGroups.PolicyID = policyID

	attackgroup, err := tf.GetStringValue("attack_group", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAttackGroups.Group = attackgroup

	attackgroups, err := client.GetEvalGroups(ctx, getAttackGroups)
	if err != nil {
		logger.Errorf("calling 'getEvalGroups': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "EvalGroupDS", attackgroups)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if len(attackgroups.AttackGroups) == 1 {
		if err := d.Set("attack_group_action", attackgroups.AttackGroups[0].Action); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		conditionException, err := json.Marshal(attackgroups.AttackGroups[0].ConditionException)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("condition_exception", string(conditionException)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		jsonBody, err := json.Marshal(attackgroups)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("json", string(jsonBody)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	d.SetId(strconv.Itoa(getAttackGroups.ConfigID))

	return nil
}
