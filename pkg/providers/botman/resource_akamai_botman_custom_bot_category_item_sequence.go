package botman

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomBotCategoryItemSequence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomBotCategoryItemSequenceCreate,
		ReadContext:   resourceCustomBotCategoryItemSequenceRead,
		UpdateContext: resourceCustomBotCategoryItemSequenceUpdate,
		DeleteContext: resourceCustomBotCategoryItemSequenceDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifyCategoryIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"category_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the bot category",
			},
			"bot_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Unique identifiers of bots in this category, sorted in preferred order",
			},
		},
	}
}
func resourceCustomBotCategoryItemSequenceUpsert(ctx context.Context, d *schema.ResourceData, m interface{}) (int64, diag.Diagnostics) {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryItemSequenceUpsert")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, int(configID), "CustomBotCategoryItemSequence", m)
	if err != nil {
		return configID, diag.FromErr(err)
	}

	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}

	sequence, err := tf.GetTypedListValue[string]("bot_ids", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}
	var stringSequence []string
	stringSequence = append(stringSequence, sequence...)

	request := botman.UpdateCustomBotCategoryItemSequenceRequest{
		ConfigID:   configID,
		Version:    int64(version),
		CategoryID: categoryID,
		Sequence:   botman.UUIDSequence{Sequence: stringSequence},
	}

	_, err = client.UpdateCustomBotCategoryItemSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpsertCustomBotCategoryItemSequence': %s", err.Error())
		return configID, diag.FromErr(err)
	}
	return configID, nil
}

func resourceCustomBotCategoryItemSequenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configID, diagnostics := resourceCustomBotCategoryItemSequenceUpsert(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}
	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d:%s", configID, categoryID))
	return resourceCustomBotCategoryItemSequenceRead(ctx, d, m)
}

func resourceCustomBotCategoryItemSequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, diagnostics := resourceCustomBotCategoryItemSequenceUpsert(ctx, d, m)
	if diagnostics != nil {
		return diagnostics
	}
	return resourceCustomBotCategoryItemSequenceRead(ctx, d, m)
}

func resourceCustomBotCategoryItemSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategoryItemSequenceRead")

	idParts, err := id.Split(d.Id(), 2, "configID:categoryID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID := idParts[1]
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetCustomBotCategoryItemSequenceRequest{
		ConfigID:   int64(configID),
		Version:    int64(version),
		CategoryID: categoryID,
	}

	response, err := client.GetCustomBotCategoryItemSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomBotCategoryItemSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	fields := map[string]interface{}{
		"config_id":   configID,
		"category_id": categoryID,
		"bot_ids":     response.Sequence,
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceCustomBotCategoryItemSequenceDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceCustomBotCategoryItemSequenceDelete")
	logger.Info("Botman API does not support custom bot category item sequence deletion - resource will only be removed from state")
	return nil
}
