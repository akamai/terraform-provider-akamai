package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBotDetectionAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBotDetectionActionCreate,
		ReadContext:   resourceBotDetectionActionRead,
		UpdateContext: resourceBotDetectionActionUpdate,
		DeleteContext: resourceBotDetectionActionDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
			verifyDetectionIDUnchanged,
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
			"detection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bot_detection_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceBotDetectionActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotDetectionActionCreate")
	logger.Debugf("in resourceBotDetectionActionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botDetectionAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	detectionID, err := tf.GetStringValue("detection_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayload, err := getJSONPayload(d, "bot_detection_action", "detectionId", detectionID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotDetectionActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		DetectionID:      detectionID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateBotDetectionAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, detectionID))

	return botDetectionActionRead(ctx, d, m, false)
}

func resourceBotDetectionActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return botDetectionActionRead(ctx, d, m, true)
}
func botDetectionActionRead(ctx context.Context, d *schema.ResourceData, m interface{}, readFromCache bool) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotDetectionActionRead")
	logger.Debugf("in resourceBotDetectionActionRead")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:botDetectionID")
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

	detectionID := iDParts[2]

	getBotDetectionActionRequest := botman.GetBotDetectionActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		DetectionID:      detectionID,
	}

	var response map[string]interface{}
	if readFromCache {
		response, err = getBotDetectionAction(ctx, getBotDetectionActionRequest, m)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		response, err = client.GetBotDetectionAction(ctx, getBotDetectionActionRequest)
		if err != nil {
			logger.Errorf("calling 'GetBotDetectionAction': %s", err.Error())
			return diag.FromErr(err)
		}
	}

	// Removing detectionId from response to suppress diff
	delete(response, "detectionId")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":            configID,
		"security_policy_id":   securityPolicyID,
		"detection_id":         detectionID,
		"bot_detection_action": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceBotDetectionActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotDetectionActionUpdate")
	logger.Debugf("in resourceBotDetectionActionUpdate")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:botDetectionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botDetectionAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	detectionID := iDParts[2]

	jsonPayload, err := getJSONPayload(d, "bot_detection_action", "detectionId", detectionID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotDetectionActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		DetectionID:      detectionID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateBotDetectionAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return botDetectionActionRead(ctx, d, m, false)
}

func resourceBotDetectionActionDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceBotDetectionActionDelete")
	logger.Debugf("in resourceBotDetectionActionDelete")
	logger.Info("Botman API does not support bot detection category action deletion - resource will only be removed from state")

	return nil
}
