package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBotCategoryException() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBotCategoryExceptionCreate,
		ReadContext:   resourceBotCategoryExceptionRead,
		UpdateContext: resourceBotCategoryExceptionUpdate,
		DeleteContext: resourceBotCategoryExceptionDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
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
			"bot_category_exception": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceBotCategoryExceptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotCategoryExceptionCreate")
	logger.Debugf("in resourceBotCategoryExceptionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botCategoryException", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("bot_category_exception", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotCategoryExceptionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateBotCategoryException(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))

	return resourceBotCategoryExceptionRead(ctx, d, m)
}

func resourceBotCategoryExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotCategoryExceptionRead")
	logger.Debugf("in resourceBotCategoryExceptionRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
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

	request := botman.GetBotCategoryExceptionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
	}

	response, err := client.GetBotCategoryException(ctx, request)

	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":              configID,
		"security_policy_id":     securityPolicyID,
		"bot_category_exception": string(jsonBody),
	}
	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceBotCategoryExceptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceBotCategoryExceptionUpdate")
	logger.Debugf("in resourceBotCategoryExceptionUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "botCategoryException", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	jsonPayloadString, err := tf.GetStringValue("bot_category_exception", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateBotCategoryExceptionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpdateBotCategoryException(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceBotCategoryExceptionRead(ctx, d, m)
}

func resourceBotCategoryExceptionDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceBotCategoryExceptionDelete")
	logger.Debugf("in resourceBotCategoryExceptionDelete")
	logger.Info("Botman API does not support bot category exception deletion - resource will only be removed from state")

	return nil
}
