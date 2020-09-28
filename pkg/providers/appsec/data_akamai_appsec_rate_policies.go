package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRatePolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRatePoliciesRead,
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

func dataSourceRatePoliciesRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceRatePoliciesRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read RatePolicies")

	ratepolicies := appsec.NewRatePoliciesResponse()
	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)

	err := ratepolicies.GetRatePolicies(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("RatePolicies   %v\n", ratepolicies))

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "ratePolicies", ratepolicies)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ratePolicies outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(configid))

	return nil
}
