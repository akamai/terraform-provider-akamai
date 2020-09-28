package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMatchTargets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMatchTargetsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
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

func dataSourceMatchTargetsRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceMatchTargetsRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read MatchTargets")

	matchtargets := appsec.NewMatchTargetsResponse()
	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)

	err := matchtargets.GetMatchTargets(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("MatchTargets   %v\n", matchtargets))
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("MatchTargets API  %v\n", matchtargets.MatchTargets.APITargets))
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("MatchTargets WEB  %v\n", matchtargets.MatchTargets.WebsiteTargets[0].SecurityPolicy.PolicyID))

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "DSmatchTarget", matchtargets)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("matchTarget outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(configid))

	return nil
}
