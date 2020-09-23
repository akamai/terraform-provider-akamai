package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRatePolicyActions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRatePolicyActionsRead,
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

func dataSourceRatePolicyActionsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceRatePolicyActionsRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read RatePolicyActions")

	ratepolicyactions := appsec.NewRatePolicyActionsResponse()
	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)

	err := ratepolicyactions.GetRatePolicyActions(configid, version, policyid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("RatePolicyActions   %v\n", ratepolicyactions))

	for _, configval := range ratepolicyactions.RatePolicyActions {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ratepolicyaction  configval %v\n", configval.ID))
		d.SetId(strconv.Itoa(configval.ID))
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "ratePolicyActions", ratepolicyactions)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ratePolicyActions outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	return nil
}
