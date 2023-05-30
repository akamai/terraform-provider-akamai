package botman

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRecategorizedAkamaiDefinedBot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecategorizedAkamaiDefinedBotCreate,
		ReadContext:   resourceRecategorizedAkamaiDefinedBotRead,
		UpdateContext: resourceRecategorizedAkamaiDefinedBotUpdate,
		DeleteContext: resourceRecategorizedAkamaiDefinedBotDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifyBotIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bot_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRecategorizedAkamaiDefinedBotCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceRecategorizedAkamaiDefinedBotCreateAction")
	logger.Debugf("in resourceRecategorizedAkamaiDefinedBotCreateAction")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "RecategorizedAkamaiDefinedBot", m)
	if err != nil {
		return diag.FromErr(err)
	}

	botID, err := tf.GetStringValue("bot_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateRecategorizedAkamaiDefinedBotRequest{
		ConfigID:   int64(configID),
		Version:    int64(version),
		BotID:      botID,
		CategoryID: categoryID,
	}

	response, err := client.CreateRecategorizedAkamaiDefinedBot(ctx, request)
	if err != nil {
		logger.Errorf("calling 'CreateRecategorizedAkamaiDefinedBot': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, response.BotID))

	return resourceRecategorizedAkamaiDefinedBotRead(ctx, d, m)
}

func resourceRecategorizedAkamaiDefinedBotRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceRecategorizedAkamaiDefinedBotReadAction")
	logger.Debugf("in resourceRecategorizedAkamaiDefinedBotReadAction")

	iDParts, err := splitID(d.Id(), 2, "configID:botID")
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

	botID := iDParts[1]

	request := botman.GetRecategorizedAkamaiDefinedBotRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		BotID:    botID,
	}

	response, err := client.GetRecategorizedAkamaiDefinedBot(ctx, request)

	if err != nil {
		logger.Errorf("calling 'GetRecategorizedAkamaiDefinedBot': %s", err.Error())
		return diag.FromErr(err)
	}

	fields := map[string]interface{}{
		"config_id":   configID,
		"bot_id":      botID,
		"category_id": response.CategoryID,
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceRecategorizedAkamaiDefinedBotUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceRecategorizedAkamaiDefinedBotUpdateAction")
	logger.Debugf("in resourceRecategorizedAkamaiDefinedBotUpdateAction")

	iDParts, err := splitID(d.Id(), 2, "configID:botID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "RecategorizedAkamaiDefinedBot", m)
	if err != nil {
		return diag.FromErr(err)
	}

	botID := iDParts[1]

	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateRecategorizedAkamaiDefinedBotRequest{
		ConfigID:   int64(configID),
		Version:    int64(version),
		BotID:      botID,
		CategoryID: categoryID,
	}

	_, err = client.UpdateRecategorizedAkamaiDefinedBot(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateRecategorizedAkamaiDefinedBot': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceRecategorizedAkamaiDefinedBotRead(ctx, d, m)
}

func resourceRecategorizedAkamaiDefinedBotDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceRecategorizedAkamaiDefinedBotDeleteAction")
	logger.Debugf("in resourceRecategorizedAkamaiDefinedBotDeleteAction")

	iDParts, err := splitID(d.Id(), 2, "configID:botID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "RecategorizedAkamaiDefinedBot", m)
	if err != nil {
		return diag.FromErr(err)
	}

	botID := iDParts[1]

	request := botman.RemoveRecategorizedAkamaiDefinedBotRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
		BotID:    botID,
	}

	err = client.RemoveRecategorizedAkamaiDefinedBot(ctx, request)
	if err != nil {
		logger.Errorf("calling 'RemoveRecategorizedAkamaiDefinedBot': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
