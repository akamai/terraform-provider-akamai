package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSiemSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSiemSettingsUpdate,
		ReadContext:   resourceSiemSettingsRead,
		UpdateContext: resourceSiemSettingsUpdate,
		DeleteContext: resourceSiemSettingsDelete,
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
			"enable_siem": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_for_all_policies": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_botman_siem": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"siem_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceSiemSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemSettingsRead")

	getSiemSettings := v2.GetSiemSettingsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getSiemSettings.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getSiemSettings.Version = version
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSiemSettings.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSiemSettings.Version = version
	}
	siemsettings, err := client.GetSiemSettings(ctx, getSiemSettings)
	if err != nil {
		logger.Errorf("calling 'getSiemSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "siemsettingsDS", siemsettings)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	if err := d.Set("config_id", getSiemSettings.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getSiemSettings.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("enable_siem", siemsettings.EnableSiem); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("enable_for_all_policies", siemsettings.EnableForAllPolicies); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("enable_botman_siem", siemsettings.EnabledBotmanSiemEvents); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("siem_id", siemsettings.SiemDefinitionID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_ids", siemsettings.FirewallPolicyIds); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d", getSiemSettings.ConfigID, getSiemSettings.Version))

	return nil
}

func resourceSiemSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemSettingsUpdate")

	removeSiemSettings := v2.RemoveSiemSettingsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeSiemSettings.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeSiemSettings.Version = version
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeSiemSettings.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeSiemSettings.Version = version
	}
	removeSiemSettings.EnableSiem = false

	siemID, err := tools.GetIntValue("siem_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSiemSettings.SiemDefinitionID = siemID

	_, erru := client.RemoveSiemSettings(ctx, removeSiemSettings)
	if erru != nil {
		logger.Errorf("calling 'removeSiemSettings': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}

func resourceSiemSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemSettingsUpdate")

	updateSiemSettings := v2.UpdateSiemSettingsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSiemSettings.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSiemSettings.Version = version
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSiemSettings.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSiemSettings.Version = version
	}
	enableSiem, err := tools.GetBoolValue("enable_siem", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSiemSettings.EnableSiem = enableSiem

	enableForAllPolicies, err := tools.GetBoolValue("enable_for_all_policies", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSiemSettings.EnableForAllPolicies = enableForAllPolicies

	enableBotmanSiem, err := tools.GetBoolValue("enable_botman_siem", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSiemSettings.EnabledBotmanSiemEvents = enableBotmanSiem

	siemID, err := tools.GetIntValue("siem_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSiemSettings.SiemDefinitionID = siemID

	security_policy_ids := d.Get("security_policy_ids").([]interface{})
	spids := make([]string, 0, len(security_policy_ids))

	for _, h := range security_policy_ids {
		spids = append(spids, h.(string))

	}

	updateSiemSettings.FirewallPolicyIds = spids

	_, erru := client.UpdateSiemSettings(ctx, updateSiemSettings)
	if erru != nil {
		logger.Errorf("calling 'updateSiemSettings': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSiemSettingsRead(ctx, d, m)
}
