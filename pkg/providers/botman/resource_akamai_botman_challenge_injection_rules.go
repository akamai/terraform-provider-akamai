package botman

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceChallengeInjectionRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChallengeInjectionRulesCreate,
		ReadContext:   resourceChallengeInjectionRulesRead,
		UpdateContext: resourceChallengeInjectionRulesUpdate,
		DeleteContext: resourceChallengeInjectionRulesDelete,
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
			"challenge_injection_rules": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceChallengeInjectionRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeInjectionRulesCreate")
	logger.Debugf("in resourceChallengeInjectionRulesCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "challengeInjectionRules", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("challenge_injection_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateChallengeInjectionRulesRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateChallengeInjectionRules(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return resourceChallengeInjectionRulesRead(ctx, d, m)
}

func resourceChallengeInjectionRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeInjectionRulesRead")
	logger.Debugf("in resourceChallengeInjectionRulesRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetChallengeInjectionRulesRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetChallengeInjectionRules(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                 configID,
		"challenge_injection_rules": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(configID))
	return nil
}

func resourceChallengeInjectionRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeInjectionRulesUpdate")
	logger.Debugf("in resourceChallengeInjectionRulesUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "challengeInjectionRules", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("challenge_injection_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateChallengeInjectionRulesRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateChallengeInjectionRules(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateChallengeInjectionRules': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceChallengeInjectionRulesRead(ctx, d, m)
}

func resourceChallengeInjectionRulesDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceChallengeInjectionRulesDelete")
	logger.Debugf("in resourceChallengeInjectionRulesDelete")
	logger.Info("Botman API does not support challenge injection rules deletion - resource will only be removed from state")

	return nil
}
