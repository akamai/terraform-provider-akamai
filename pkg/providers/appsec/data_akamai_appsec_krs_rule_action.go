package appsec

import (
	"context"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKRSRuleActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKRSRuleActionsRead,
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

func dataSourceKRSRuleActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceKRSRuleActionsRead")

	getKRSRuleActions := v2.GetKRSRuleActionsRequest{}

	getKRSRuleActions.ConfigID = d.Get("config_id").(int)
	getKRSRuleActions.Version = d.Get("version").(int)
	getKRSRuleActions.PolicyID = d.Get("policy_id").(string)

	krsruleactions, err := client.GetKRSRuleActions(ctx, getKRSRuleActions)
	if err != nil {
		logger.Errorf("calling 'getKRSRuleActions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "KRSRulesDS", krsruleactions)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getKRSRuleActions.ConfigID))

	return nil
}
