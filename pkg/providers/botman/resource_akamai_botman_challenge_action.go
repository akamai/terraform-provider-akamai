package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceChallengeAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChallengeActionCreate,
		ReadContext:   resourceChallengeActionRead,
		UpdateContext: resourceChallengeActionUpdate,
		DeleteContext: resourceChallengeActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"action_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"challenge_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceChallengeActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeActionCreateAction")
	logger.Debugf("in resourceChallengeActionCreateAction")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ChallengeAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("challenge_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateChallengeActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateChallengeAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, tools.ConvertToString((response)["actionId"])))

	return resourceChallengeActionRead(ctx, d, m)
}

func resourceChallengeActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeActionRead")
	logger.Debugf("in resourceChallengeActionRead")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
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

	actionID := iDParts[1]

	request := botman.GetChallengeActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	response, err := client.GetChallengeAction(ctx, request)

	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	// Removing actionId from response to suppress diff
	delete(response, "actionId")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":        configID,
		"action_id":        actionID,
		"challenge_action": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceChallengeActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeActionUpdate")
	logger.Debugf("in resourceChallengeActionUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ChallengeAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	jsonPayload, err := getJSONPayload(d, "challenge_action", "actionId", actionID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateChallengeActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		ActionID:    actionID,
		JsonPayload: jsonPayload,
	}

	_, err = client.UpdateChallengeAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceChallengeActionRead(ctx, d, m)
}

func resourceChallengeActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceChallengeActionDelete")
	logger.Debugf("in resourceChallengeActionDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ChallengeAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	removeChallengeAction := botman.RemoveChallengeActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	err = client.RemoveChallengeAction(ctx, removeChallengeAction)
	if err != nil {
		logger.Errorf("calling 'removeChallengeAction': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
