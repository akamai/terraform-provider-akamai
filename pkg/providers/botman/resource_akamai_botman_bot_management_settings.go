package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBotManagementSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBotManagementSettingsCreate,
		ReadContext:   resourceBotManagementSettingsRead,
		UpdateContext: resourceBotManagementSettingsUpdate,
		DeleteContext: resourceBotManagementSettingsDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bot_management_settings": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceBotManagementSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotManagementSettingsCreate")
	logger.Debugf("in resourceBotManagementSettingsCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botManagementSettings", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tools.GetStringValue("bot_management_settings", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotManagementSettingRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateBotManagementSetting(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))

	return resourceBotManagementSettingsRead(ctx, d, m)
}

func resourceBotManagementSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotManagementSettingsRead")
	logger.Debugf("in resourceBotManagementSettingsRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	request := botman.GetBotManagementSettingRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
	}

	response, err := client.GetBotManagementSetting(ctx, request)

	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":               configID,
		"security_policy_id":      securityPolicyID,
		"bot_management_settings": string(jsonBody),
	}
	if err = tools.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceBotManagementSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotManagementSettingsUpdate")
	logger.Debugf("in resourceBotManagementSettingsUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botManagementSettings", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	jsonPayloadString, err := tools.GetStringValue("bot_management_settings", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotManagementSettingRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateBotManagementSetting(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceBotManagementSettingsRead(ctx, d, m)
}

func resourceBotManagementSettingsDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("botman", "resourceBotManagementSettingsDelete")
	logger.Debugf("in resourceBotManagementSettingsDelete")
	logger.Info("Botman API does not support bot management settings deletion - resource will only be removed from state")

	return nil
}
