package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceSiemSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSiemSettingsCreate,
		ReadContext:   resourceSiemSettingsRead,
		UpdateContext: resourceSiemSettingsUpdate,
		DeleteContext: resourceSiemSettingsDelete,
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
			"enable_siem": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable SIEM",
			},
			"enable_for_all_policies": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable SIEM on all security policies in the security configuration",
			},
			"security_policy_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of IDs of security policy for which SIEM integration is to be enabled",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enable_botman_siem": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether Bot Manager events should be included in SIEM events",
			},
			"siem_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the SIEM settings being modified",
			},
			"exceptions": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem:        getExceptionsResource(),
				Description: "Describes all the protections and actions to be excluded from SIEM events",
			},
		},
	}
}

func getExceptionsResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip_geo": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include ip geo events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"bot_management": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include bot management events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"rate": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include rate events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"url_protection": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include url protection events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"slow_post": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include slow post events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"custom_rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include custom rules events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"waf": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include waf events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"api_request_constraints": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include api request constraints events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"client_rep": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include client reputation events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"malware_protection": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include malware protection events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
			"apr_protection": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Whether there should be an exception to include apr protection events in SIEM",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
			},
		},
	}
}

func resourceSiemSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemSettingsCreate")
	logger.Debugf("in resourceSiemSettingsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "siemSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	enableSiem, err := tf.GetBoolValue("enable_siem", d)
	if err != nil {
		return diag.FromErr(err)
	}
	enableForAllPolicies, err := tf.GetBoolValue("enable_for_all_policies", d)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyIDs, err := tf.GetSetValue("security_policy_ids", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	spIDs := make([]string, 0, len(securityPolicyIDs.List()))
	for _, h := range securityPolicyIDs.List() {
		spIDs = append(spIDs, h.(string))
	}
	siemID, err := tf.GetIntValue("siem_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	exceptions, err := getAllExceptions(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createSiemSettings := appsec.UpdateSiemSettingsRequest{
		ConfigID:             configID,
		Version:              version,
		EnableSiem:           enableSiem,
		EnableForAllPolicies: enableForAllPolicies,
		FirewallPolicyIDs:    spIDs,
		SiemDefinitionID:     siemID,
		Exceptions:           exceptions,
	}

	enableBotmanSiem, err := tf.GetBoolValue("enable_botman_siem", tf.NewRawConfig(d))
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if !errors.Is(err, tf.ErrNotFound) {
		createSiemSettings.EnabledBotmanSiemEvents = ptr.To(enableBotmanSiem)
	}

	_, err = client.UpdateSiemSettings(ctx, createSiemSettings)
	if err != nil {
		logger.Errorf("calling 'createSiemSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", createSiemSettings.ConfigID))

	return resourceSiemSettingsRead(ctx, d, m)
}

func resourceSiemSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemSettingsRead")
	logger.Debugf("in resourceSiemSettingsRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	getSiemSettings := appsec.GetSiemSettingsRequest{
		ConfigID: configID,
		Version:  version,
	}

	siemsettings, err := client.GetSiemSettings(ctx, getSiemSettings)
	if err != nil {
		logger.Errorf("calling 'getSiemSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getSiemSettings.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("enable_siem", siemsettings.EnableSiem); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("enable_for_all_policies", siemsettings.EnableForAllPolicies); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_ids", siemsettings.FirewallPolicyIDs); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("enable_botman_siem", siemsettings.EnabledBotmanSiemEvents); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("siem_id", siemsettings.SiemDefinitionID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := setActionsFromExceptions(d, siemsettings.Exceptions); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceSiemSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemSettingsUpdate")
	logger.Debugf("in resourceSiemSettingsUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "siemSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	enableSiem, err := tf.GetBoolValue("enable_siem", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	enableForAllPolicies, err := tf.GetBoolValue("enable_for_all_policies", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	securityPolicyIDs, err := tf.GetSetValue("security_policy_ids", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	spIDs := make([]string, 0, len(securityPolicyIDs.List()))
	for _, h := range securityPolicyIDs.List() {
		spIDs = append(spIDs, h.(string))

	}
	siemID, err := tf.GetIntValue("siem_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	exceptions, err := getAllExceptions(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateSiemSettings := appsec.UpdateSiemSettingsRequest{
		ConfigID:             configID,
		Version:              version,
		EnableSiem:           enableSiem,
		EnableForAllPolicies: enableForAllPolicies,
		FirewallPolicyIDs:    spIDs,
		SiemDefinitionID:     siemID,
		Exceptions:           exceptions,
	}

	enableBotmanSiem, err := tf.GetBoolValue("enable_botman_siem", tf.NewRawConfig(d))
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if !errors.Is(err, tf.ErrNotFound) {
		updateSiemSettings.EnabledBotmanSiemEvents = ptr.To(enableBotmanSiem)
	}

	_, err = client.UpdateSiemSettings(ctx, updateSiemSettings)
	if err != nil {
		logger.Errorf("calling 'updateSiemSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceSiemSettingsRead(ctx, d, m)
}

func resourceSiemSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSiemSettingsDelete")
	logger.Debugf("in resourceSiemSettingsDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "siemSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}

	removeSiemSettings := appsec.RemoveSiemSettingsRequest{
		ConfigID:   configID,
		Version:    version,
		EnableSiem: false,
	}

	_, err = client.RemoveSiemSettings(ctx, removeSiemSettings)
	if err != nil {
		logger.Errorf("calling 'updateSiemSettings': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}

func getAllExceptions(d *schema.ResourceData) ([]appsec.Exception, error) {
	exceptions := make([]appsec.Exception, 0)
	exceptionsMap := getConfigParamProtectionMapping()
	exceptionsConfig, err := tf.GetListValue("exceptions", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	err = validateExceptions(exceptionsConfig)
	if err != nil {
		return nil, err
	}
	for _, exception := range exceptionsConfig {
		exceptionMap := exception.(map[string]interface{})
		for key := range exceptionMap {
			actions := make([]string, 0)
			if _, ok := exceptionsMap[key]; ok {
				actionsSet, ok := exceptionMap[key].(*schema.Set)
				if !ok {
					return nil, fmt.Errorf("wrong type conversion: expected *schema.Set, got %T", actionsSet)
				}
				for _, action := range actionsSet.List() {
					actions = append(actions, action.(string))
				}
				if len(actions) > 0 {
					exceptions = append(exceptions, appsec.Exception{Protection: exceptionsMap[key], ActionTypes: actions})
				}
			}
		}
	}
	return exceptions, nil
}

func validateExceptions(exceptionsConfig []interface{}) error {
	if len(exceptionsConfig) > 0 {
		_, ok := exceptionsConfig[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("Invalid exceptions configuration")
		}
	}
	return nil
}

func setActionsFromExceptions(d *schema.ResourceData, exceptions []appsec.Exception) error {
	if err := d.Set("exceptions", exceptionsToState(exceptions)); err != nil {
		return fmt.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func getConfigParamProtectionMapping() map[string]string {
	exceptionsMap := map[string]string{
		"bot_management":          "botmanagement",
		"rate":                    "rate",
		"ip_geo":                  "ipgeo",
		"url_protection":          "urlProtection",
		"slow_post":               "slowpost",
		"custom_rules":            "customrules",
		"waf":                     "waf",
		"api_request_constraints": "apirequestconstraints",
		"client_rep":              "clientrep",
		"malware_protection":      "malwareprotection",
		"apr_protection":          "aprProtection",
	}
	return exceptionsMap
}

func exceptionsToState(exceptions []appsec.Exception) []interface{} {
	exceptionsMap := getConfigParamProtectionMapping()
	out := make([]interface{}, 0, len(exceptions))
	exceptionMap := make(map[string]interface{})

	for configParamName, configParamVal := range exceptionsMap {
		for _, t := range exceptions {
			if configParamVal == t.Protection {
				exceptionMap[configParamName] = t.ActionTypes
				break
			}
		}
	}
	if len(exceptionMap) > 0 {
		out = append(out, exceptionMap)
	}

	return out
}
