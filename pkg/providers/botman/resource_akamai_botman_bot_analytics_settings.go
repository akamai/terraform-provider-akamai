package botman

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBotAnalyticsSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBotAnalyticsSettingsCreate,
		ReadContext:   resourceBotAnalyticsSettingsRead,
		UpdateContext: resourceBotAnalyticsSettingsUpdate,
		DeleteContext: resourceBotAnalyticsSettingsDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"bot_analytics_settings": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
				Description:      "JSON-formatted bot analytics settings.",
			},
		},
	}
}

func resourceBotAnalyticsSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "resourceBotAnalyticsSettingsCreate")
	logger.Debugf("in resourceBotAnalyticsSettingsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botAnalyticsSettings", meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("bot_analytics_settings", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotAnalyticsSettingsRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JSONPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateBotAnalyticsSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateBotAnalyticsSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return resourceBotAnalyticsSettingsRead(ctx, d, m)
}

func resourceBotAnalyticsSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "resourceBotAnalyticsSettingsRead")
	logger.Debugf("in resourceBotAnalyticsSettingsRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetBotAnalyticsSettingsRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetBotAnalyticsSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetBotAnalyticsSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":              configID,
		"bot_analytics_settings": string(jsonBody),
	}
	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceBotAnalyticsSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "resourceBotAnalyticsSettingsUpdate")
	logger.Debugf("in resourceBotAnalyticsSettingsUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botAnalyticsSettings", meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("bot_analytics_settings", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotAnalyticsSettingsRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JSONPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateBotAnalyticsSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateBotAnalyticsSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceBotAnalyticsSettingsRead(ctx, d, m)
}

func resourceBotAnalyticsSettingsDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceBotAnalyticsSettingsDelete")
	logger.Debugf("in resourceBotAnalyticsSettingsDelete")
	logger.Info("Botman API does not support bot analytics settings deletion - resource will only be removed from state")

	return nil
}
