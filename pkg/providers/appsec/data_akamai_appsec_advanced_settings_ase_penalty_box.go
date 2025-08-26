package appsec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdvancedSettingsAsePenaltyBox() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedSettingsAsePenaltyBoxRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation for the ASE Penalty Box settings.",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation for the ASE Penalty Box settings.",
			},
		},
	}
}

func dataSourceAdvancedSettingsAsePenaltyBoxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAdvancedSettingsAsePenaltyBoxRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	advancedSettingsAsePenaltyBox, err := client.GetAdvancedSettingsAsePenaltyBox(ctx, appsec.GetAdvancedSettingsAsePenaltyBoxRequest{
		ConfigID: configID,
		Version:  version,
	})
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsAsePenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputText, err := RenderTemplates(ots, "advancedSettingsAsePenaltyBoxDS", advancedSettingsAsePenaltyBox)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputText); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(advancedSettingsAsePenaltyBox)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d", configID))

	return nil
}
