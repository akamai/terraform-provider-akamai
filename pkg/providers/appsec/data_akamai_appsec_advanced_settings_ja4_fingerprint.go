package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdvancedSettingsJA4Fingerprint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedSettingsJA4FingerprintRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation of JA4 Fingerprint settings",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation of JA4 Fingerprint settings",
			},
		},
	}
}

func dataSourceAdvancedSettingsJA4FingerprintRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAdvancedSettingsJA4FingerprintRead")

	getAdvancedSettingsJA4FingerprintReq := appsec.GetAdvancedSettingsJA4FingerprintRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAdvancedSettingsJA4FingerprintReq.ConfigID = configID

	if getAdvancedSettingsJA4FingerprintReq.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	advancedSettingsJA4Fingerprint, err := client.GetAdvancedSettingsJA4Fingerprint(ctx, getAdvancedSettingsJA4FingerprintReq)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsJA4Fingerprint': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputText, err := RenderTemplates(ots, "advancedSettingsJA4FingerprintDS", advancedSettingsJA4Fingerprint)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputText); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(advancedSettingsJA4Fingerprint)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getAdvancedSettingsJA4FingerprintReq.ConfigID))

	return nil
}
