package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomRuleActions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCustomRuleActionsRead,
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

func dataSourceCustomRuleActionsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceCustomRuleActionsRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read CustomRuleActions")

	customruleactions := appsec.NewCustomRuleActionsResponse()
	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	//ruleid := d.Get("rule_id").(int)

	err := customruleactions.GetCustomRuleActions(configid, version, policyid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CustomRuleActions   %v\n", customruleactions))

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "customRuleAction", customruleactions)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("customRuleAction outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	//d.Set("rule_id", ruleid)
	d.SetId(strconv.Itoa(configid))

	return nil
}
