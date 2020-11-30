package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAttackGroupConditionException() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttackGroupConditionExceptionRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attack_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceAttackGroupConditionExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupConditionExceptionsRead")

	getAttackGroupConditionException := v2.GetAttackGroupConditionExceptionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.PolicyID = policyid

	attackgroup, err := tools.GetStringValue("attack_group", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAttackGroupConditionException.Group = attackgroup

	attackgroupconditionexception, err := client.GetAttackGroupConditionException(ctx, getAttackGroupConditionException)
	if err != nil {
		logger.Errorf("calling 'getAttackGroupConditionException': %s", err.Error())
		return diag.FromErr(err)
	}

	logger.Errorf("calling 'getAttackGroupConditionException FOUND : %s", attackgroupconditionexception)

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "AttackGroupConditionException", attackgroupconditionexception)

	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(attackgroupconditionexception)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("json", string(jsonBody))

	d.SetId(strconv.Itoa(getAttackGroupConditionException.ConfigID))

	return nil
}
