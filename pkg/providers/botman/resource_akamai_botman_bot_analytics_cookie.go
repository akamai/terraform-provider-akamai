package botman

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBotAnalyticsCookie() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBotAnalyticsCookieCreate,
		ReadContext:   resourceBotAnalyticsCookieRead,
		UpdateContext: resourceBotAnalyticsCookieUpdate,
		DeleteContext: resourceBotAnalyticsDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bot_analytics_cookie": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceBotAnalyticsCookieCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotAnalyticsCookieCreate")
	logger.Debugf("in resourceBotAnalyticsCookieCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botAnalyticsCookie", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("bot_analytics_cookie", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotAnalyticsCookieRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: (json.RawMessage)(jsonPayloadString),
	}

	_, err = client.UpdateBotAnalyticsCookie(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return resourceBotAnalyticsCookieRead(ctx, d, m)
}

func resourceBotAnalyticsCookieRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotAnalyticsCookieRead")
	logger.Debugf("in resourceBotAnalyticsCookieRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetBotAnalyticsCookieRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetBotAnalyticsCookie(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":            configID,
		"bot_analytics_cookie": string(jsonBody),
	}
	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceBotAnalyticsCookieUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotAnalyticsCookieUpdate")
	logger.Debugf("in resourceBotAnalyticsCookieUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botAnalyticsCookie", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("bot_analytics_cookie", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotAnalyticsCookieRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateBotAnalyticsCookie(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceBotAnalyticsCookieRead(ctx, d, m)
}

func resourceBotAnalyticsDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("botman", "resourceBotAnalyticsDelete")
	logger.Debugf("in resourceBotAnalyticsDelete")
	logger.Info("Botman API does not support bot analytics cookie settings deletion - resource will only be removed from state")

	return nil
}
