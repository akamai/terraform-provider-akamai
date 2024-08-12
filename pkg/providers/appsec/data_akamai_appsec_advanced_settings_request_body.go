package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdvancedSettingsRequestBody() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedSettingsRequestBodyRead,
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

func dataSourceAdvancedSettingsRequestBodyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAdvancedSettingsRequestBodyRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	advancedSettingsRequestBody, err := client.GetAdvancedSettingsRequestBody(ctx, appsec.GetAdvancedSettingsRequestBodyRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	})
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsRequestBody': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	if policyID != "" {
		outputText, err := RenderTemplates(ots, "advancedSettingsRequestBodyPolicyDS", advancedSettingsRequestBody)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("output_text", outputText); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

	} else {
		outputText, err := RenderTemplates(ots, "advancedSettingsRequestBodyDS", advancedSettingsRequestBody)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("output_text", outputText); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	jsonBody, err := json.Marshal(advancedSettingsRequestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))

	return nil
}
