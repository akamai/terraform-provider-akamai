package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceAdvancedSettingsPrefetch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsPrefetchUpdate,
		ReadContext:   resourceAdvancedSettingsPrefetchRead,
		UpdateContext: resourceAdvancedSettingsPrefetchUpdate,
		DeleteContext: resourceAdvancedSettingsPrefetchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
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
			"enable_rate_controls": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"extensions": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of extensions",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceAdvancedSettingsPrefetchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchRead")

	getAdvancedSettingsPrefetch := appsec.GetAdvancedSettingsPrefetchRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getAdvancedSettingsPrefetch.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getAdvancedSettingsPrefetch.Version = version

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAdvancedSettingsPrefetch.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAdvancedSettingsPrefetch.Version = version
	}
	advancedsettingsprefetch, err := client.GetAdvancedSettingsPrefetch(ctx, getAdvancedSettingsPrefetch)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsPrefetch': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "advancedSettingsPrefetchDS", advancedsettingsprefetch)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	if err := d.Set("config_id", getAdvancedSettingsPrefetch.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getAdvancedSettingsPrefetch.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d", getAdvancedSettingsPrefetch.ConfigID, getAdvancedSettingsPrefetch.Version))

	return nil
}

func resourceAdvancedSettingsPrefetchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchRemove")

	updateAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAdvancedSettingsPrefetch.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAdvancedSettingsPrefetch.Version = version

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAdvancedSettingsPrefetch.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAdvancedSettingsPrefetch.Version = version
	}
	updateAdvancedSettingsPrefetch.EnableAppLayer = false

	updateAdvancedSettingsPrefetch.EnableRateControls = false

	_, erru := client.UpdateAdvancedSettingsPrefetch(ctx, updateAdvancedSettingsPrefetch)
	if erru != nil {
		logger.Errorf("calling 'removeAdvancedSettingsPrefetch': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}

func resourceAdvancedSettingsPrefetchUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPrefetchUpdate")

	updateAdvancedSettingsPrefetch := appsec.UpdateAdvancedSettingsPrefetchRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAdvancedSettingsPrefetch.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAdvancedSettingsPrefetch.Version = version

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAdvancedSettingsPrefetch.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAdvancedSettingsPrefetch.Version = version
	}
	enableAppLayer, err := tools.GetBoolValue("enable_app_layer", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsPrefetch.EnableAppLayer = enableAppLayer

	allExtensions, err := tools.GetBoolValue("all_extensions", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsPrefetch.AllExtensions = allExtensions

	enableRateControls, err := tools.GetBoolValue("enable_rate_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsPrefetch.EnableRateControls = enableRateControls

	extensions := d.Get("extensions").([]interface{})
	exts := make([]string, 0, len(extensions))

	for _, h := range extensions {
		exts = append(exts, h.(string))

	}

	updateAdvancedSettingsPrefetch.Extensions = exts

	logger.Errorf("calling 'getAdvancedSettingsPrefetch': Extensions %v", updateAdvancedSettingsPrefetch.Extensions)

	_, erru := client.UpdateAdvancedSettingsPrefetch(ctx, updateAdvancedSettingsPrefetch)
	if erru != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPrefetch': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceAdvancedSettingsPrefetchRead(ctx, d, m)
}
