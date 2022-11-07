package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSiemSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiemSettingsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
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

func dataSourceSiemSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceSiemSettingsRead")

	getSiemSettings := v2.GetSiemSettingsRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getSiemSettings.ConfigID = configID

	if getSiemSettings.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	siemsettings, err := client.GetSiemSettings(ctx, getSiemSettings)
	if err != nil {
		logger.Errorf("calling 'getSiemSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext := ""
	settingstext, err := RenderTemplates(ots, "siemsettingsDS", siemsettings)
	if err != nil {
		return diag.FromErr(err)
	}
	outputtext = outputtext + settingstext
	policiestext, err := RenderTemplates(ots, "siempoliciesDS", siemsettings)
	if err != nil {
		return diag.FromErr(err)
	}
	outputtext = outputtext + policiestext
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(siemsettings)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(strconv.Itoa(getSiemSettings.ConfigID))

	return nil
}
