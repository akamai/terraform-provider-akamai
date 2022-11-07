package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceServeAlternateAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServeAlternateActionCreate,
		ReadContext:   resourceServeAlternateActionRead,
		UpdateContext: resourceServeAlternateActionUpdate,
		DeleteContext: resourceServeAlternateActionDelete,
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
			"serve_alternate_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceServeAlternateActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceServeAlternateActionCreateAction")
	logger.Debugf("in resourceServeAlternateActionCreateAction")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ServeAlternateAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tools.GetStringValue("serve_alternate_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateServeAlternateActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateServeAlternateAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, tools.ConvertToString((response)["actionId"])))

	return resourceServeAlternateActionRead(ctx, d, m)
}

func resourceServeAlternateActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceServeAlternateActionRead")
	logger.Debugf("in resourceServeAlternateActionRead")

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

	request := botman.GetServeAlternateActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	response, err := client.GetServeAlternateAction(ctx, request)

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
		"config_id":              configID,
		"action_id":              actionID,
		"serve_alternate_action": string(jsonBody),
	}
	if err := tools.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceServeAlternateActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceServeAlternateActionUpdate")
	logger.Debugf("in resourceServeAlternateActionUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ServeAlternateAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	jsonPayload, err := getJSONPayload(d, "serve_alternate_action", "actionId", actionID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateServeAlternateActionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		ActionID:    actionID,
		JsonPayload: jsonPayload,
	}

	_, err = client.UpdateServeAlternateAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceServeAlternateActionRead(ctx, d, m)
}

func resourceServeAlternateActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceServeAlternateActionDelete")
	logger.Debugf("in resourceServeAlternateActionDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:actionID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ServeAlternateAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	actionID := iDParts[1]

	removeServeAlternateAction := botman.RemoveServeAlternateActionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		ActionID: actionID,
	}

	err = client.RemoveServeAlternateAction(ctx, removeServeAlternateAction)
	if err != nil {
		logger.Errorf("calling 'removeServeAlternateAction': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
