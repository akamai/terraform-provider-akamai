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

func dataSourceSiemDefinitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiemDefinitionsRead,
		Schema: map[string]*schema.Schema{
			"siem_definition_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON Siem Definition",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceSiemDefinitionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemDefinitionsRead")

	getSiemDefinitions := appsec.GetSiemDefinitionsRequest{}

	siem_definition_name, err := tools.GetStringValue("siem_definition_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSiemDefinitions.SiemDefinitionName = siem_definition_name

	siemdefinitions, err := client.GetSiemDefinitions(ctx, getSiemDefinitions)
	if err != nil {
		logger.Errorf("calling 'getSiemDefinitions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "siemDefinitionsDS", siemdefinitions)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(siemdefinitions)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	//d.Set("enabled", GetSlowPostProtection.ApplySlowPostControls)
	//d.SetId(strconv.Itoa(getRateProtection.ConfigID))

	if len(siemdefinitions.SiemDefinitions) > 0 {
		d.SetId(strconv.Itoa(siemdefinitions.SiemDefinitions[0].ID))
	}

	return nil
}
