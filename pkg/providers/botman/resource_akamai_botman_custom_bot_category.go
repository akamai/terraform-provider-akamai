package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCustomBotCategory() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomBotCategoryCreate,
		ReadContext:   resourceCustomBotCategoryRead,
		UpdateContext: resourceCustomBotCategoryUpdate,
		DeleteContext: resourceCustomBotCategoryDelete,
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
			"category_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_bot_category": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceCustomBotCategoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryCreateAction")
	logger.Debugf("in resourceCustomBotCategoryCreateAction")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomBotCategory", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("custom_bot_category", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateCustomBotCategoryRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	response, err := client.CreateCustomBotCategory(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, str.From((response)["categoryId"])))

	return resourceCustomBotCategoryRead(ctx, d, m)
}

func resourceCustomBotCategoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryRead")
	logger.Debugf("in resourceCustomBotCategoryRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:categoryID")
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

	categoryID := iDParts[1]

	request := botman.GetCustomBotCategoryRequest{
		ConfigID:   int64(configID),
		Version:    int64(version),
		CategoryID: categoryID,
	}

	response, err := client.GetCustomBotCategory(ctx, request)

	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	// Removing categoryId from response to suppress diff
	delete(response, "categoryId")
	// Removing read-only fields
	delete(response, "metadata")
	delete(response, "ruleId")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":           configID,
		"category_id":         categoryID,
		"custom_bot_category": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceCustomBotCategoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryUpdate")
	logger.Debugf("in resourceCustomBotCategoryUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:categoryID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomBotCategory", m)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID := iDParts[1]

	jsonPayload, err := getJSONPayload(d, "custom_bot_category", "categoryId", categoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateCustomBotCategoryRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		CategoryID:  categoryID,
		JsonPayload: jsonPayload,
	}

	_, err = client.UpdateCustomBotCategory(ctx, request)
	if err != nil {
		logger.Errorf("calling 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceCustomBotCategoryRead(ctx, d, m)
}

func resourceCustomBotCategoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryDelete")
	logger.Debugf("in resourceCustomBotCategoryDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:categoryID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "CustomBotCategory", m)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID := iDParts[1]

	removeCustomBotCategory := botman.RemoveCustomBotCategoryRequest{
		ConfigID:   int64(configID),
		Version:    int64(version),
		CategoryID: categoryID,
	}

	err = client.RemoveCustomBotCategory(ctx, removeCustomBotCategory)
	if err != nil {
		logger.Errorf("calling 'removeCustomBotCategory': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
