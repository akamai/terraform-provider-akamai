package appsec

import (
	"context"
	"encoding/json"
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
func resourceAdvancedSettingsLogging() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsLoggingUpdate,
		ReadContext:   resourceAdvancedSettingsLoggingRead,
		UpdateContext: resourceAdvancedSettingsLoggingUpdate,
		DeleteContext: resourceAdvancedSettingsLoggingDelete,
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
				Optional: true,
			},
			"logging": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
		},
	}
}

func resourceAdvancedSettingsLoggingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsLoggingRead")

	getAdvancedSettingsLogging := appsec.GetAdvancedSettingsLoggingRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsLogging.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsLogging.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsLogging.PolicyID = policyid

	advancedsettingslogging, err := client.GetAdvancedSettingsLogging(ctx, getAdvancedSettingsLogging)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsLogging': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "advancedSettingsLoggingDS", advancedsettingslogging)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getAdvancedSettingsLogging.ConfigID))

	return nil
}

func resourceAdvancedSettingsLoggingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceAdvancedSettingsLoggingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsLoggingUpdate")

	updateAdvancedSettingsLogging := appsec.UpdateAdvancedSettingsLoggingRequest{}

	jsonpostpayload := d.Get("logging").(string)

	if err := json.Unmarshal([]byte(jsonpostpayload), &updateAdvancedSettingsLogging); err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsLogging.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsLogging.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsLogging.PolicyID = policyid

	_, erru := client.UpdateAdvancedSettingsLogging(ctx, updateAdvancedSettingsLogging)
	if erru != nil {
		logger.Errorf("calling 'updateAdvancedSettingsLogging': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceAdvancedSettingsLoggingRead(ctx, d, m)
}
