package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceConditionalAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConditionalActionCreate,
		ReadContext:   resourceConditionalActionRead,
		UpdateContext: resourceConditionalActionUpdate,
		DeleteContext: resourceConditionalActionDelete,
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
			"conditional_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceConditionalActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceConditionalActionCreateAction")
	logger.Debugf("in resourceConditionalActionCreateAction")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ConditionalAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("conditional_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateConditionalActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateConditionalAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, tools.ConvertToString((response)["actionId"])))

	return resourceConditionalActionRead(ctx, d, m)
}

func resourceConditionalActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceConditionalActionRead")
	logger.Debugf("in resourceConditionalActionRead")

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

	request := botman.GetConditionalActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	response, err := client.GetConditionalAction(ctx, request)

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
		"conditional_action": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceConditionalActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceConditionalActionUpdate")
	logger.Debugf("in resourceConditionalActionUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ConditionalAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	jsonPayload, err := getJSONPayload(d, "conditional_action", "actionId", actionID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateConditionalActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		ActionID:    actionID,
		JsonPayload: jsonPayload,
	}

	_, err = client.UpdateConditionalAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceConditionalActionRead(ctx, d, m)
}

func resourceConditionalActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceConditionalActionDelete")
	logger.Debugf("in resourceConditionalActionDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ConditionalAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	removeConditionalAction := botman.RemoveConditionalActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	err = client.RemoveConditionalAction(ctx, removeConditionalAction)
	if err != nil {
		logger.Errorf("calling 'removeConditionalAction': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
