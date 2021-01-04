package appsec

import (
	"context"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceAdvancedSettingsPolicyLogging() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsPolicyLoggingUpdate,
		ReadContext:   resourceAdvancedSettingsPolicyLoggingRead,
		UpdateContext: resourceAdvancedSettingsPolicyLoggingUpdate,
		DeleteContext: resourceAdvancedSettingsPolicyLoggingDelete,
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
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logging": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceAdvancedSettingsPolicyLoggingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPolicyLoggingRead")

	getAdvancedSettingsPolicyLogging := appsec.GetAdvancedSettingsPolicyLoggingRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsPolicyLogging.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsPolicyLogging.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsPolicyLogging.PolicyID = policyid

	advancedsettingspolicylogging, err := client.GetAdvancedSettingsPolicyLogging(ctx, getAdvancedSettingsPolicyLogging)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsPolicyLogging': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "advancedSettingsPolicyLoggingDS", advancedsettingspolicylogging)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getAdvancedSettingsPolicyLogging.ConfigID))

	return nil
}

func resourceAdvancedSettingsPolicyLoggingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceAdvancedSettingsPolicyLoggingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPolicyLoggingUpdate")

	updateAdvancedSettingsPolicyLogging := appsec.UpdateAdvancedSettingsPolicyLoggingRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsPolicyLogging.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsPolicyLogging.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsPolicyLogging.PolicyID = policyid

	_, erru := client.UpdateAdvancedSettingsPolicyLogging(ctx, updateAdvancedSettingsPolicyLogging)
	if erru != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPolicyLogging': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceAdvancedSettingsPolicyLoggingRead(ctx, d, m)
}
