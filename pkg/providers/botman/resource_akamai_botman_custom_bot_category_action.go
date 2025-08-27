package botman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomBotCategoryAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomBotCategoryActionCreate,
		ReadContext:   resourceCustomBotCategoryActionRead,
		UpdateContext: resourceCustomBotCategoryActionUpdate,
		DeleteContext: resourceCustomBotCategoryActionDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
			verifyCategoryIDUnchanged,
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
			"category_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_bot_category_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceCustomBotCategoryActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryActionCreate")
	logger.Debugf("in resourceCustomBotCategoryActionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "customBotCategoryAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	jsonPayload, err := getJSONPayload(d, "custom_bot_category_action", "categoryId", categoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateCustomBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateCustomBotCategoryAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, categoryID))

	return customBotCategoryActionRead(ctx, d, m, false)
}

func resourceCustomBotCategoryActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return customBotCategoryActionRead(ctx, d, m, true)
}

func customBotCategoryActionRead(ctx context.Context, d *schema.ResourceData, m interface{}, readFromCache bool) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryActionRead")
	logger.Debugf("in resourceCustomBotCategoryActionRead")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:categoryID")
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

	categoryID := iDParts[2]

	request := botman.GetCustomBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
	}

	var response map[string]interface{}
	if readFromCache {
		response, err = getCustomBotCategoryAction(ctx, request, m)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		response, err = client.GetCustomBotCategoryAction(ctx, request)
		if err != nil {
			logger.Errorf("calling 'GetCustomBotCategoryAction': %s", err.Error())
			return diag.FromErr(err)
		}
	}

	// Removing categoryId from response to suppress diff
	delete(response, "categoryId")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                  configID,
		"security_policy_id":         securityPolicyID,
		"category_id":                categoryID,
		"custom_bot_category_action": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceCustomBotCategoryActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryActionUpdate")
	logger.Debugf("in resourceCustomBotCategoryActionUpdate")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:categoryID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "customBotCategoryAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	categoryID := iDParts[2]

	jsonPayload, err := getJSONPayload(d, "custom_bot_category_action", "categoryId", categoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateCustomBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateCustomBotCategoryAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return customBotCategoryActionRead(ctx, d, m, false)
}

func resourceCustomBotCategoryActionDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceCustomBotCategoryActionDelete")
	logger.Debugf("in resourceCustomBotCategoryActionDelete")
	logger.Info("Botman API does not support custom bot category action deletion - resource will only be removed from state")

	return nil
}
