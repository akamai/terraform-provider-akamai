package botman

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomClientSequence() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomClientSequenceCreate,
		ReadContext:   resourceCustomClientSequenceRead,
		UpdateContext: resourceCustomClientSequenceUpdate,
		DeleteContext: resourceCustomClientSequenceDelete,
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
			"custom_client_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCustomClientSequenceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	configID, diagnostics := resourceCustomClientSequenceUpsert(ctx, d, m, "resourceCustomClientSequenceCreate")
	if diagnostics != nil {
		return diagnostics
	}
	d.SetId(strconv.Itoa(configID))
	return resourceCustomClientSequenceRead(ctx, d, m)
}

func resourceCustomClientSequenceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, diagnostics := resourceCustomClientSequenceUpsert(ctx, d, m, "resourceCustomClientSequenceUpdate")
	if diagnostics != nil {
		return diagnostics
	}
	return resourceCustomClientSequenceRead(ctx, d, m)
}

func resourceCustomClientSequenceUpsert(ctx context.Context, d *schema.ResourceData, m interface{}, operation string) (int, diag.Diagnostics) {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", operation)
	logger.Debugf("in %s", operation)

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "customClientSequence", m)
	if err != nil {
		return configID, diag.FromErr(err)
	}
	sequence, err := tf.GetListValue("custom_client_ids", d)
	if err != nil {
		return configID, diag.FromErr(err)
	}
	var stringSequence []string
	for _, val := range sequence {
		stringSequence = append(stringSequence, val.(string))
	}

	request := botman.UpdateCustomClientSequenceRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		Sequence: stringSequence,
	}

	_, err = client.UpdateCustomClientSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateCustomClientSequence': %s", err.Error())
		return configID, diag.FromErr(err)
	}

	return configID, nil
}

func resourceCustomClientSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceCustomClientSequenceRead")
	logger.Debugf("in resourceCustomClientSequenceRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetCustomClientSequenceRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetCustomClientSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomClientSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	fields := map[string]interface{}{
		"config_id":         configID,
		"custom_client_ids": response.Sequence,
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceCustomClientSequenceDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceCustomClientSequenceDelete")
	logger.Debugf("in resourceCustomClientSequenceDelete")
	logger.Info("Botman API does not support custom client sequence deletion - resource will only be removed from state")

	return nil
}
