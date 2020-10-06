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

func dataSourceRatePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRatePoliciesRead,
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

func dataSourceRatePoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePoliciesRead")
	CorrelationID := "[APPSEC][resourceRatePolicies-" + meta.OperationID() + "]"

	getRatePolicies := v2.GetRatePoliciesRequest{}

	getRatePolicies.ConfigID = d.Get("config_id").(int)
	getRatePolicies.ConfigVersion = d.Get("version").(int)

	ratepolicies, err := client.GetRatePolicies(ctx, getRatePolicies)
	if err != nil {
		logger.Warnf("calling 'getRatePolicies': %s", err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "ratePolicies", ratepolicies)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ratePolicies outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getRatePolicies.ConfigID))

	return nil
}
