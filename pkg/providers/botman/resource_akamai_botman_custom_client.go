package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/meta"
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
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "resourceCustomClientCreateAction")
	logger.Debugf("in resourceCustomClientCreateAction")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomClient", meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("custom_client", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateCustomClientRequest{
		ConfigID:    configID,
		Version:     version,
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
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "resourceCustomClientReadAction")
	logger.Debugf("in resourceCustomClientReadAction")

	iDParts, err := id.Split(d.Id(), 2, "configID:customClientID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	customClientID := iDParts[1]

	request := botman.GetCustomClientRequest{
		ConfigID:       configID,
		Version:        version,
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
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceCustomClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "resourceCustomClientUpdateAction")
	logger.Debugf("in resourceCustomClientUpdateAction")

	iDParts, err := id.Split(d.Id(), 2, "configID:customClientID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomClient", meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	customClientID := iDParts[1]

	jsonPayload, err := getJSONPayload(d, "custom_client", "customClientId", customClientID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateCustomClientRequest{
		ConfigID:       configID,
		Version:        version,
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
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "resourceCustomClientDeleteAction")
	logger.Debugf("in resourceCustomClientDeleteAction")

	iDParts, err := id.Split(d.Id(), 2, "configID:customClientID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomClient", meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	customClientID := iDParts[1]

	request := botman.RemoveCustomClientRequest{
		ConfigID:       configID,
		Version:        version,
		CustomClientID: customClientID,
	}

	err = client.RemoveCustomClient(ctx, request)
	if err != nil {
		logger.Errorf("calling 'RemoveCustomClient': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
