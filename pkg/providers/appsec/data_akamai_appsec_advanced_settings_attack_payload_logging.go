package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdvancedSettingsAttackPayloadLogging() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedSettingsAttackPayloadLoggingRead,
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

func dataSourceAdvancedSettingsAttackPayloadLoggingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAdvancedSettingsAttackPayloadLoggingRead")

	getAdvancedSettingsAttackPayloadLogging := appsec.GetAdvancedSettingsAttackPayloadLoggingRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAdvancedSettingsAttackPayloadLogging.ConfigID = configID

	if getAdvancedSettingsAttackPayloadLogging.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsAttackPayloadLogging.PolicyID = policyID

	advancedSettingsAttackPayloadLogging, err := client.GetAdvancedSettingsAttackPayloadLogging(ctx, getAdvancedSettingsAttackPayloadLogging)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsAttackPayloadLogging': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputText, err := RenderTemplates(ots, "advancedSettingsAttackPayloadLoggingDS", advancedSettingsAttackPayloadLogging)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputText); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(advancedSettingsAttackPayloadLogging)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))

	return nil
}
