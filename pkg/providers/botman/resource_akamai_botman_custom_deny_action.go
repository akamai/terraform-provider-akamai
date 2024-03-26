package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCustomDenyAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomDenyActionCreate,
		ReadContext:   resourceCustomDenyActionRead,
		UpdateContext: resourceCustomDenyActionUpdate,
		DeleteContext: resourceCustomDenyActionDelete,
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
			"custom_deny_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceCustomDenyActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomDenyActionCreateAction")
	logger.Debugf("in resourceCustomDenyActionCreateAction")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomDenyAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("custom_deny_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateCustomDenyActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateCustomDenyAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, str.From((response)["actionId"])))

	return resourceCustomDenyActionRead(ctx, d, m)
}

func resourceCustomDenyActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomDenyActionRead")
	logger.Debugf("in resourceCustomDenyActionRead")

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

	request := botman.GetCustomDenyActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	response, err := client.GetCustomDenyAction(ctx, request)

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
		"config_id":          configID,
		"action_id":          actionID,
		"custom_deny_action": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceCustomDenyActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomDenyActionUpdate")
	logger.Debugf("in resourceCustomDenyActionUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomDenyAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	jsonPayload, err := getJSONPayload(d, "custom_deny_action", "actionId", actionID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateCustomDenyActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		ActionID:    actionID,
		JsonPayload: jsonPayload,
	}

	_, err = client.UpdateCustomDenyAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceCustomDenyActionRead(ctx, d, m)
}

func resourceCustomDenyActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomDenyActionDelete")
	logger.Debugf("in resourceCustomDenyActionDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomDenyAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	removeCustomDenyAction := botman.RemoveCustomDenyActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	err = client.RemoveCustomDenyAction(ctx, removeCustomDenyAction)
	if err != nil {
		logger.Errorf("calling 'removeCustomDenyAction': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
