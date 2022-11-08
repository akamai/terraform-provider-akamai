package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceAdvancedSettingsPrefetch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsPrefetchCreate,
		ReadContext:   resourceAdvancedSettingsPrefetchRead,
		UpdateContext: resourceAdvancedSettingsPrefetchUpdate,
		DeleteContext: resourceAdvancedSettingsPrefetchDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
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
			"enable_app_layer": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable or disable prefetch requests",
			},
			"all_extensions": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable prefetch requests for all file extensions",
			},
			"extensions": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of file extensions",
			},
			"enable_rate_controls": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable prefetch requests for rate controls",
			},
		},
	}
}

func resourceAdvancedSettingsPrefetchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchCreate")
	logger.Debugf("in resourceAdvancedSettingsPrefetchCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "prefetchSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	enableAppLayer, err := tools.GetBoolValue("enable_app_layer", d)
	if err != nil {
		return diag.FromErr(err)
	}
	allExtensions, err := tools.GetBoolValue("all_extensions", d)
	if err != nil {
		return diag.FromErr(err)
	}
	extensions := d.Get("extensions").(*schema.Set)
	if err != nil {
		return diag.FromErr(err)
	}
	exts := make([]string, 0, len(extensions.List()))
	for _, h := range extensions.List() {
		exts = append(exts, h.(string))

	}
	enableRateControls, err := tools.GetBoolValue("enable_rate_controls", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{
		ConfigID:           configID,
		Version:            version,
		EnableAppLayer:     enableAppLayer,
		AllExtensions:      allExtensions,
		Extensions:         exts,
		EnableRateControls: enableRateControls,
	}

	_, err = client.UpdateAdvancedSettingsPrefetch(ctx, createAdvancedSettingsPrefetch)
	if err != nil {
		logger.Errorf("calling 'createAdvancedSettingsPrefetch': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", createAdvancedSettingsPrefetch.ConfigID))

	return resourceAdvancedSettingsPrefetchRead(ctx, d, m)
}

func resourceAdvancedSettingsPrefetchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchRead")
	logger.Debugf("in resourceAdvancedSettingsPrefetchRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	getAdvancedSettingsPrefetch := appsec.GetAdvancedSettingsPrefetchRequest{
		ConfigID: configID,
		Version:  version,
	}

	prefetchget, err := client.GetAdvancedSettingsPrefetch(ctx, getAdvancedSettingsPrefetch)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsPrefetch': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAdvancedSettingsPrefetch.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("enable_app_layer", prefetchget.EnableAppLayer); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("all_extensions", prefetchget.AllExtensions); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("extensions", prefetchget.Extensions); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("enable_rate_controls", prefetchget.EnableRateControls); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAdvancedSettingsPrefetchUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchUpdate")
	logger.Debugf("in resourceAdvancedSettingsPrefetchUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "prefetchSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
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

	updateAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{
		ConfigID:           configID,
		Version:            version,
		EnableAppLayer:     enableAppLayer,
		AllExtensions:      allExtensions,
		Extensions:         exts,
		EnableRateControls: enableRateControls,
	}

	_, err = client.UpdateAdvancedSettingsPrefetch(ctx, updateAdvancedSettingsPrefetch)
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPrefetch': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAdvancedSettingsPrefetchRead(ctx, d, m)
}

func resourceAdvancedSettingsPrefetchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchDelete")
	logger.Debugf("in resourceAdvancedSettingsPrefetchDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "prefetchSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	removeAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{
		ConfigID:           configID,
		Version:            version,
		EnableAppLayer:     false,
		EnableRateControls: false,
	}

	_, err = client.UpdateAdvancedSettingsPrefetch(ctx, removeAdvancedSettingsPrefetch)
	if err != nil {
		logger.Errorf("calling 'removeAdvancedSettingsPrefetch': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
