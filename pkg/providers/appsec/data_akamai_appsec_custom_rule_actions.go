package appsec

import (
	"context"
	"fmt"
	"strconv"

	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomRuleActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomRuleActionsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceCustomRuleActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomRuleActionsRead")
	CorrelationID := "[APPSEC][resourceCustomRuleActions-" + meta.OperationID() + "]"

	getCustomRuleActions := v2.GetCustomRuleActionsRequest{}

	getCustomRuleActions.ConfigID = d.Get("config_id").(int)
	getCustomRuleActions.Version = d.Get("version").(int)
	getCustomRuleActions.PolicyID = d.Get("policy_id").(string)

	customruleactions, err := client.GetCustomRuleActions(ctx, getCustomRuleActions)
	if err != nil {
		logger.Warnf("calling 'getCustomRuleActions': %s", err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "customRuleAction", customruleactions)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("customRuleAction outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	//d.Set("rule_id", ruleid)
	d.SetId(strconv.Itoa(getCustomRuleActions.ConfigID))

	return nil
}
