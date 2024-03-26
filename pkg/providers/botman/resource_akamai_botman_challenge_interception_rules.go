package botman

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceChallengeInterceptionRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChallengeInterceptionRulesCreate,
		ReadContext:   resourceChallengeInterceptionRulesRead,
		UpdateContext: resourceChallengeInterceptionRulesUpdate,
		DeleteContext: resourceChallengeInterceptionRulesDelete,
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
			"challenge_interception_rules": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceChallengeInterceptionRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeInterceptionRulesCreate")
	logger.Debugf("in resourceChallengeInterceptionRulesCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "challengeInterceptionRules", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("challenge_interception_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateChallengeInterceptionRulesRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateChallengeInterceptionRules(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return resourceChallengeInterceptionRulesRead(ctx, d, m)
}

func resourceChallengeInterceptionRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeInterceptionRulesRead")
	logger.Debugf("in resourceChallengeInterceptionRulesRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetChallengeInterceptionRulesRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetChallengeInterceptionRules(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                    configID,
		"challenge_interception_rules": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(configID))
	return nil
}

func resourceChallengeInterceptionRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeInterceptionRulesUpdate")
	logger.Debugf("in resourceChallengeInterceptionRulesUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "challengeInterceptionRules", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("challenge_interception_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateChallengeInterceptionRulesRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateChallengeInterceptionRules(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateChallengeInterceptionRules': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceChallengeInterceptionRulesRead(ctx, d, m)
}

func resourceChallengeInterceptionRulesDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceChallengeInterceptionRulesDelete")
	logger.Debugf("in resourceChallengeInterceptionRulesDelete")
	logger.Info("Botman API does not support challenge interception rules deletion - resource will only be removed from state")

	return nil
}
