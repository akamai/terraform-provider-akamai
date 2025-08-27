package accountprotection

import (
	"context"
	"encoding/json"
	"strconv"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUserRiskResponseStrategy() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceUserRiskResponseStrategy,
		ReadContext:   readResourceUserRiskResponseStrategy,
		UpdateContext: updateResourceUserRiskResponseStrategy,
		DeleteContext: deleteResourceUserRiskResponseStrategy,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifies a security configuration.",
			},
			"user_risk_response_strategy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func createResourceUserRiskResponseStrategy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "createResourceUserRiskResponseStrategy")
	logger.Debugf("in createResourceUserRiskResponseStrategy")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "userRiskResponseStrategy", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("user_risk_response_strategy", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.UpsertUserRiskResponseStrategyRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpsertUserRiskResponseStrategy(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpsertUserRiskResponseStrategy': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return readResourceUserRiskResponseStrategy(ctx, d, m)
}

func readResourceUserRiskResponseStrategy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "readResourceUserRiskResponseStrategy")
	logger.Debugf("in readResourceUserRiskResponseStrategy")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.GetUserRiskResponseStrategyRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetUserRiskResponseStrategy(ctx, request)

	if err != nil {
		logger.Errorf("error calling GetUserRiskResponseStrategy: %v", err)
		return diag.FromErr(err)
	}

	delete(response, "metadata")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                   configID,
		"user_risk_response_strategy": string(jsonBody),
	}
	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func updateResourceUserRiskResponseStrategy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "updateResourceUserRiskResponseStrategy")
	logger.Debugf("in updateResourceUserRiskResponseStrategy")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "userRiskResponseStrategy", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("user_risk_response_strategy", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.UpsertUserRiskResponseStrategyRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpsertUserRiskResponseStrategy(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpsertUserRiskResponseStrategy': %s", err.Error())
		return diag.FromErr(err)
	}

	return readResourceUserRiskResponseStrategy(ctx, d, m)
}

func deleteResourceUserRiskResponseStrategy(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("accountprotection", "deleteResourceUserRiskResponseStrategy")
	logger.Debugf("in deleteResourceUserRiskResponseStrategy")
	logger.Info("APR API does not support user risk response strategy deletion - resource will only be removed from state")

	d.SetId("")
	return nil
}
