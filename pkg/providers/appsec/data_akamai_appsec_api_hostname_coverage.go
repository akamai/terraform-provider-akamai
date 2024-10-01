package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAPIHostnameCoverage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAPIHostnameCoverageRead,
		Schema: map[string]*schema.Schema{
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

func dataSourceAPIHostnameCoverageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAPIHostnameCoverageRead")

	getAPIHostnameCoverage := appsec.GetApiHostnameCoverageRequest{}

	apihostnamecoverage, err := client.GetApiHostnameCoverage(ctx, getAPIHostnameCoverage)
	if err != nil {
		logger.Errorf("calling 'getApiHostnameCoverage': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "apiHostnameCoverageDS", apihostnamecoverage)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(apihostnamecoverage)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if len(apihostnamecoverage.HostnameCoverage) > 0 {
		for _, configval := range apihostnamecoverage.HostnameCoverage {
			if configval.Configuration != nil && configval.Configuration.ID != 0 {
				d.SetId(strconv.Itoa(configval.Configuration.ID))
			}
		}
	}

	return nil
}
