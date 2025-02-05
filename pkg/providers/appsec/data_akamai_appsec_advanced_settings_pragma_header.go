package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdvancedSettingsPragmaHeader() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedSettingsPragmaHeaderRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of the security policy",
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

func dataSourceAdvancedSettingsPragmaHeaderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAdvancedSettingsPragmaHeaderRead")

	getAdvancedSettingsPragma := appsec.GetAdvancedSettingsPragmaRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAdvancedSettingsPragma.ConfigID = configID

	if getAdvancedSettingsPragma.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsPragma.PolicyID = policyID

	advancedsettingspragma, err := client.GetAdvancedSettingsPragma(ctx, getAdvancedSettingsPragma)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsPragmaHeader': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "advancedSettingsPragmaHeaderDS", advancedsettingspragma)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(advancedsettingspragma)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getAdvancedSettingsPragma.ConfigID))

	return nil
}
