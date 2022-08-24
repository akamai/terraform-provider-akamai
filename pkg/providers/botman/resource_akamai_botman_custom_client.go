package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCustomClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomClientCreate,
		ReadContext:   resourceCustomClientRead,
		UpdateContext: resourceCustomClientUpdate,
		DeleteContext: resourceCustomClientDelete,
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
			"custom_client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_client": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceCustomClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomClientCreateAction")
	logger.Debugf("in resourceCustomClientCreateAction")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomClient", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tools.GetStringValue("custom_client", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateCustomClientRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateCustomClient(ctx, request)
	if err != nil {
		logger.Errorf("calling 'CreateCustomClient': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, (response)["customClientId"].(string)))

	return resourceCustomClientRead(ctx, d, m)
}

func resourceCustomClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomClientReadAction")
	logger.Debugf("in resourceCustomClientReadAction")

	iDParts, err := splitID(d.Id(), 2, "configID:customClientID")
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

	customClientID := iDParts[1]

	request := botman.GetCustomClientRequest{
		ConfigID:       int64(configID),
		Version:        int64(version),
		CustomClientID: customClientID,
	}

	response, err := client.GetCustomClient(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomClient': %s", err.Error())
		return diag.FromErr(err)
	}

	// Removing customClientId from response to suppress diff
	delete(response, "customClientId")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":        configID,
		"custom_client_id": customClientID,
		"custom_client":    string(jsonBody),
	}
	if err := tools.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceCustomClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomClientUpdateAction")
	logger.Debugf("in resourceCustomClientUpdateAction")

	iDParts, err := splitID(d.Id(), 2, "configID:customClientID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomClient", m)
	if err != nil {
		return diag.FromErr(err)
	}

	customClientID := iDParts[1]

	jsonPayload, err := getJSONPayload(d, "custom_client", "customClientId", customClientID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateCustomClientRequest{
		ConfigID:       int64(configID),
		Version:        int64(version),
		CustomClientID: customClientID,
		JsonPayload:    jsonPayload,
	}

	_, err = client.UpdateCustomClient(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateCustomClient': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceCustomClientRead(ctx, d, m)
}

func resourceCustomClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomClientDeleteAction")
	logger.Debugf("in resourceCustomClientDeleteAction")

	iDParts, err := splitID(d.Id(), 2, "configID:customClientID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomClient", m)
	if err != nil {
		return diag.FromErr(err)
	}

	customClientID := iDParts[1]

	request := botman.RemoveCustomClientRequest{
		ConfigID:       int64(configID),
		Version:        int64(version),
		CustomClientID: customClientID,
	}

	err = client.RemoveCustomClient(ctx, request)
	if err != nil {
		logger.Errorf("calling 'RemoveCustomClient': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
