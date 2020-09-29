package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCustomRulesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
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

func dataSourceCustomRulesRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceCustomRulesRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read CustomRules")

	customrules := appsec.NewCustomRulesResponse()
	configid := d.Get("config_id").(int)

	err := customrules.GetCustomRules(configid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CustomRules   %v\n", customrules))

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "customRules", customrules)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("customrules outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(configid))

	return nil
}
