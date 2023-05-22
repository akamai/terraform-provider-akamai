package botman

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTransactionalEndpointProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTransactionalEndpointProtectionCreate,
		ReadContext:   resourceTransactionalEndpointProtectionRead,
		UpdateContext: resourceTransactionalEndpointProtectionUpdate,
		DeleteContext: resourceTransactionEndpointProtectionDelete,
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
			"transactional_endpoint_protection": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceTransactionalEndpointProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceTransactionalEndpointProtectionCreate")
	logger.Debugf("in resourceTransactionalEndpointProtectionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "transactionalEndpointProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("transactional_endpoint_protection", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateTransactionalEndpointProtectionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateTransactionalEndpointProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateTransactionalEndpointProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return resourceTransactionalEndpointProtectionRead(ctx, d, m)
}

func resourceTransactionalEndpointProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceTransactionalEndpointProtectionRead")
	logger.Debugf("in resourceTransactionalEndpointProtectionRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetTransactionalEndpointProtectionRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetTransactionalEndpointProtection(ctx, request)

	if err != nil {
		logger.Errorf("calling 'GetTransactionalEndpointProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", request.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                         configID,
		"transactional_endpoint_protection": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceTransactionalEndpointProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceTransactionalEndpointProtectionUpdate")
	logger.Debugf("in resourceTransactionalEndpointProtectionUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "transactionalEndpointProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("transactional_endpoint_protection", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateTransactionalEndpointProtectionRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateTransactionalEndpointProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateTransactionalEndpointProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceTransactionalEndpointProtectionRead(ctx, d, m)
}

func resourceTransactionEndpointProtectionDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("botman", "resourceTransactionEndpointProtectionDelete")
	logger.Debugf("in resourceTransactionEndpointProtectionDelete")
	logger.Info("Botman API does not support transactional endpoint protection deletion - resource will only be removed from state")

	return nil
}
