package botman

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomBotCategorySequence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomBotCategorySequenceCreate,
		ReadContext:   resourceCustomBotCategorySequenceRead,
		UpdateContext: resourceCustomBotCategorySequenceUpdate,
		DeleteContext: resourceCustomBotCategorySequenceDelete,
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
			"category_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCustomBotCategorySequenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategorySequenceCreate")
	logger.Debugf("in resourceCustomBotCategorySequenceCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "customBotCategorySequence", m)
	if err != nil {
		return diag.FromErr(err)
	}

	sequence, err := tf.GetListValue("category_ids", d)
	if err != nil {
		return diag.FromErr(err)
	}
	var stringSequence []string
	for _, val := range sequence {
		stringSequence = append(stringSequence, val.(string))
	}

	request := botman.UpdateCustomBotCategorySequenceRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		Sequence: stringSequence,
	}

	_, err = client.UpdateCustomBotCategorySequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateCustomBotCategorySequence': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return resourceCustomBotCategorySequenceRead(ctx, d, m)
}

func resourceCustomBotCategorySequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategorySequenceUpdate")
	logger.Debugf("in resourceCustomBotCategorySequenceUpdate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "customBotCategorySequence", m)
	if err != nil {
		return diag.FromErr(err)
	}

	sequence, err := tf.GetListValue("category_ids", d)
	if err != nil {
		return diag.FromErr(err)
	}
	var stringSequence []string
	for _, val := range sequence {
		stringSequence = append(stringSequence, val.(string))
	}

	request := botman.UpdateCustomBotCategorySequenceRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		Sequence: stringSequence,
	}

	_, err = client.UpdateCustomBotCategorySequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateCustomBotCategorySequence': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceCustomBotCategorySequenceRead(ctx, d, m)
}

func resourceCustomBotCategorySequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomBotCategorySequenceRead")
	logger.Debugf("in resourceCustomBotCategorySequenceRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetCustomBotCategorySequenceRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetCustomBotCategorySequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomBotCategorySequence': %s", err.Error())
		return diag.FromErr(err)
	}

	fields := map[string]interface{}{
		"config_id":    configID,
		"category_ids": response.Sequence,
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceCustomBotCategorySequenceDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceCustomBotCategorySequenceDelete")
	logger.Debugf("in resourceCustomBotCategorySequenceDelete")
	logger.Info("Botman API does not support custom bot category sequence deletion - resource will only be removed from state")

	return nil
}
