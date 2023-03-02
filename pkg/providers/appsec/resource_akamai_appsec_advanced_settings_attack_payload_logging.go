package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAdvancedSettingsAttackPayloadLogging() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsAttackPayloadLoggingCreate,
		ReadContext:   resourceAdvancedSettingsAttackPayloadLoggingRead,
		UpdateContext: resourceAdvancedSettingsAttackPayloadLoggingUpdate,
		DeleteContext: resourceAdvancedSettingsAttackPayloadLoggingDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
			"attack_payload_logging": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentAttackPayloadLoggingSettingsDiffs,
				Description:      "Whether to enable, disable, or update attack payload logging settings",
			},
		},
	}
}

func resourceAdvancedSettingsAttackPayloadLoggingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAttackPayloadLoggingCreate")
	logger.Debugf("in resourceAdvancedSettingsAttackPayloadLoggingCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "attackPayloadLoggingSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	jsonPostPayload := d.Get("attack_payload_logging")
	jsonPayloadRaw := []byte(jsonPostPayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createAdvancedSettingsAttackPayloadLogging := appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{
		ConfigID:       configID,
		Version:        version,
		PolicyID:       policyID,
		JSONPayloadRaw: rawJSON,
	}

	_, err = client.UpdateAdvancedSettingsAttackPayloadLogging(ctx, createAdvancedSettingsAttackPayloadLogging)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(createAdvancedSettingsAttackPayloadLogging.PolicyID) > 0 {
		d.SetId(fmt.Sprintf("%d:%s", createAdvancedSettingsAttackPayloadLogging.ConfigID, createAdvancedSettingsAttackPayloadLogging.PolicyID))
	} else {
		d.SetId(fmt.Sprintf("%d", createAdvancedSettingsAttackPayloadLogging.ConfigID))
	}

	return resourceAdvancedSettingsAttackPayloadLoggingRead(ctx, d, m)
}

func resourceAdvancedSettingsAttackPayloadLoggingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAttackPayloadLoggingRead")
	logger.Debugf("in resourceAdvancedSettingsLoggingRead")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "attackPayloadLoggingSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	advancedSettingsAttackPayloadLogging, err := client.GetAdvancedSettingsAttackPayloadLogging(ctx, appsec.GetAdvancedSettingsAttackPayloadLoggingRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	})
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsAttackPayloadLogging': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", policyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(advancedSettingsAttackPayloadLogging)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("attack_payload_logging", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAdvancedSettingsAttackPayloadLoggingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAttackPayloadLoggingUpdate")
	logger.Debugf("in resourceAdvancedSettingsAttackPayloadLoggingUpdate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "attackPayloadLoggingSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	jsonPostPayload := d.Get("attack_payload_logging")
	jsonPayloadRaw := []byte(jsonPostPayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	_, err = client.UpdateAdvancedSettingsAttackPayloadLogging(ctx, appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{
		ConfigID:       configID,
		Version:        version,
		PolicyID:       policyID,
		JSONPayloadRaw: rawJSON,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAdvancedSettingsAttackPayloadLoggingRead(ctx, d, m)
}

func resourceAdvancedSettingsAttackPayloadLoggingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAttackPayloadLoggingDelete")
	logger.Debugf("in resourceAdvancedSettingsAttackPayloadLoggingDelete")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "attackPayloadLoggingSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	removeAdvancedSettingsAttackPayloadLogging := appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	removeAdvancedSettingsAttackPayloadLogging.RequestBody.Type = appsec.AttackPayload
	removeAdvancedSettingsAttackPayloadLogging.ResponseBody.Type = appsec.AttackPayload
	removeAdvancedSettingsAttackPayloadLogging.Enabled = true

	if policyID != "" {
		removeAdvancedSettingsAttackPayloadLogging.Override = false
	}

	_, err = client.RemoveAdvancedSettingsAttackPayloadLogging(ctx, removeAdvancedSettingsAttackPayloadLogging)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
