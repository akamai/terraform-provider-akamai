package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSiemDefinitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiemDefinitionsRead,
		Schema: map[string]*schema.Schema{
			"siem_definition_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of a specific SIEM definition for which to retrieve information",
			},
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

func dataSourceSiemDefinitionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceSiemDefinitionsRead")

	getSiemDefinitions := appsec.GetSiemDefinitionsRequest{}

	siemDdefinitionName, err := tf.GetStringValue("siem_definition_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSiemDefinitions.SiemDefinitionName = siemDdefinitionName

	siemdefinitions, err := client.GetSiemDefinitions(ctx, getSiemDefinitions)
	if err != nil {
		logger.Errorf("calling 'getSiemDefinitions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "siemDefinitionsDS", siemdefinitions)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(siemdefinitions)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if len(siemdefinitions.SiemDefinitions) > 0 {
		d.SetId(strconv.Itoa(siemdefinitions.SiemDefinitions[0].ID))
	}

	return nil
}
