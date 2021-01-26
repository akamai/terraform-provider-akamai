package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApiHostnameCoverage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiHostnameCoverageRead,
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

func dataSourceApiHostnameCoverageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiHostnameCoverageRead")

	getApiHostnameCoverage := appsec.GetApiHostnameCoverageRequest{}

	apihostnamecoverage, err := client.GetApiHostnameCoverage(ctx, getApiHostnameCoverage)
	if err != nil {
		logger.Errorf("calling 'getApiHostnameCoverage': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "apiHostnameCoverageDS", apihostnamecoverage)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(apihostnamecoverage)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if len(apihostnamecoverage.HostnameCoverage) > 0 {
		for _, configval := range apihostnamecoverage.HostnameCoverage {

			if configval.Configuration.ID != 0 {
				d.SetId(strconv.Itoa(configval.Configuration.ID))
			}
		}
	}

	return nil
}
