package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceAdvancedSettingsPrefetch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsPrefetchCreate,
		ReadContext:   resourceAdvancedSettingsPrefetchRead,
		UpdateContext: resourceAdvancedSettingsPrefetchUpdate,
		DeleteContext: resourceAdvancedSettingsPrefetchDelete,
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enable_app_layer": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"all_extensions": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"extensions": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of extensions",
			},
			"enable_rate_controls": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceAdvancedSettingsPrefetchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchCreate")
	logger.Debugf("!!! in resourceAdvancedSettingsPrefetchCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "prefetchSetting", m)
	enableAppLayer, err := tools.GetBoolValue("enable_app_layer", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	allExtensions, err := tools.GetBoolValue("all_extensions", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	extensions := d.Get("extensions").(*schema.Set)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	exts := make([]string, 0, len(extensions.List()))
	for _, h := range extensions.List() {
		exts = append(exts, h.(string))

	}
	enableRateControls, err := tools.GetBoolValue("enable_rate_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{}
	createAdvancedSettingsPrefetch.ConfigID = configid
	createAdvancedSettingsPrefetch.Version = version
	createAdvancedSettingsPrefetch.EnableAppLayer = enableAppLayer
	createAdvancedSettingsPrefetch.AllExtensions = allExtensions
	createAdvancedSettingsPrefetch.Extensions = exts
	//logger.Errorf("calling 'getAdvancedSettingsPrefetch': Extensions %v", createAdvancedSettingsPrefetch.Extensions)
	createAdvancedSettingsPrefetch.EnableRateControls = enableRateControls

	_, erru := client.UpdateAdvancedSettingsPrefetch(ctx, createAdvancedSettingsPrefetch)
	if erru != nil {
		logger.Errorf("calling 'createAdvancedSettingsPrefetch': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d", createAdvancedSettingsPrefetch.ConfigID))

	return resourceAdvancedSettingsPrefetchRead(ctx, d, m)
}

func resourceAdvancedSettingsPrefetchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchRead")
	logger.Debugf("!!! resourceAdvancedSettingsPrefetchRead")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)

	getAdvancedSettingsPrefetch := appsec.GetAdvancedSettingsPrefetchRequest{}
	getAdvancedSettingsPrefetch.ConfigID = configid
	getAdvancedSettingsPrefetch.Version = version

	prefetchget, err := client.GetAdvancedSettingsPrefetch(ctx, getAdvancedSettingsPrefetch)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsPrefetch': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAdvancedSettingsPrefetch.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enable_app_layer", prefetchget.EnableAppLayer); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("all_extensions", prefetchget.AllExtensions); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("extensions", prefetchget.Extensions); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enable_rate_controls", prefetchget.EnableRateControls); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceAdvancedSettingsPrefetchUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchUpdate")
	logger.Debugf("!!! resourceAdvancedSettingsPrefetchUpdate")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "prefetchSetting", m)
	enableAppLayer, err := tools.GetBoolValue("enable_app_layer", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	allExtensions, err := tools.GetBoolValue("all_extensions", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	extensions := d.Get("extensions").(*schema.Set)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	exts := make([]string, 0, len(extensions.List()))
	for _, h := range extensions.List() {
		exts = append(exts, h.(string))

	}
	enableRateControls, err := tools.GetBoolValue("enable_rate_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{}
	updateAdvancedSettingsPrefetch.ConfigID = configid
	updateAdvancedSettingsPrefetch.Version = version
	updateAdvancedSettingsPrefetch.EnableAppLayer = enableAppLayer
	updateAdvancedSettingsPrefetch.AllExtensions = allExtensions
	updateAdvancedSettingsPrefetch.Extensions = exts
	logger.Errorf("calling 'getAdvancedSettingsPrefetch': Extensions %v", updateAdvancedSettingsPrefetch.Extensions)
	updateAdvancedSettingsPrefetch.EnableRateControls = enableRateControls

	_, erru := client.UpdateAdvancedSettingsPrefetch(ctx, updateAdvancedSettingsPrefetch)
	if erru != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPrefetch': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceAdvancedSettingsPrefetchRead(ctx, d, m)
}

func resourceAdvancedSettingsPrefetchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchDelete")
	logger.Debugf("!!! resourceAdvancedSettingsPrefetchDelete")

	configid, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "prefetchSetting", m)

	removeAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{}
	removeAdvancedSettingsPrefetch.ConfigID = configid
	removeAdvancedSettingsPrefetch.Version = version
	removeAdvancedSettingsPrefetch.EnableAppLayer = false
	removeAdvancedSettingsPrefetch.EnableRateControls = false

	_, erru := client.UpdateAdvancedSettingsPrefetch(ctx, removeAdvancedSettingsPrefetch)
	if erru != nil {
		logger.Errorf("calling 'removeAdvancedSettingsPrefetch': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}
