package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAdvancedSettingsRequestBody() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsRequestBodyCreate,
		ReadContext:   resourceAdvancedSettingsRequestBodyRead,
		UpdateContext: resourceAdvancedSettingsRequestBodyUpdate,
		DeleteContext: resourceAdvancedSettingsRequestBodyDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: resourceAdvancedSettingsRequestBodyImport,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of the security policy",
			},
			"request_body_inspection_limit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Request body inspection size limit in KB allowed values are 'default', 8, 16, 32",
			},
			"request_body_inspection_limit_override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the Request body inspection size should be overridden at policy",
			},
		},
	}
}

func resourceAdvancedSettingsRequestBodyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsRequestBodyImport")
	logger.Debugf("Import AdvancedSettingsRequestBody")

	client := inst.Client(meta)

	getAdvancedSettingsRequestBody := appsec.GetAdvancedSettingsRequestBodyRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := id.Split(d.Id(), 2, "configID:policyID")
		if err != nil {
			return nil, err
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return nil, err
		}
		version, err := getLatestConfigVersion(ctx, configID, m)
		if err != nil {
			return nil, err
		}
		policyID := iDParts[1]

		getAdvancedSettingsRequestBody.ConfigID = configID
		getAdvancedSettingsRequestBody.Version = version
		getAdvancedSettingsRequestBody.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return nil, err
		}
		version, err := getLatestConfigVersion(ctx, configID, m)
		if err != nil {
			return nil, err
		}

		getAdvancedSettingsRequestBody.ConfigID = configID
		getAdvancedSettingsRequestBody.Version = version
	}
	d.SetId(fmt.Sprintf("%d:%s", getAdvancedSettingsRequestBody.ConfigID, getAdvancedSettingsRequestBody.PolicyID))

	advancedSettingsRequestBody, err := client.GetAdvancedSettingsRequestBody(ctx, getAdvancedSettingsRequestBody)

	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsRequestBody': %s", err.Error())
		return nil, err
	}
	if err := d.Set("config_id", getAdvancedSettingsRequestBody.ConfigID); err != nil {
		return nil, err
	}
	if err := d.Set("security_policy_id", getAdvancedSettingsRequestBody.PolicyID); err != nil {
		return nil, err
	}
	if err := d.Set("request_body_inspection_limit", advancedSettingsRequestBody.RequestBodyInspectionLimitInKB); err != nil {
		return nil, err
	}
	if getAdvancedSettingsRequestBody.PolicyID != "" {
		if err := d.Set("request_body_inspection_limit_override", advancedSettingsRequestBody.RequestBodyInspectionLimitOverride); err != nil {
			return nil, err
		}
	}
	return []*schema.ResourceData{d}, nil

}

func resourceAdvancedSettingsRequestBodyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsRequestBodyCreate")
	logger.Debugf("in resourceAdvancedSettingsRequestBodyCreate")

	return upsertAdvancedSettingsRequestBody(ctx, d, m)
}

func upsertAdvancedSettingsRequestBody(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "requestBodySetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	requestBodyInspectionLimitInKB, err := tf.GetStringValue("request_body_inspection_limit", d)
	if err != nil {
		return diag.FromErr(err)
	}

	req := appsec.UpdateAdvancedSettingsRequestBodyRequest{
		ConfigID:                       configID,
		Version:                        version,
		PolicyID:                       policyID,
		RequestBodyInspectionLimitInKB: appsec.RequestBodySizeLimit(requestBodyInspectionLimitInKB),
	}
	if policyID != "" {
		override, err := tf.GetBoolValue("request_body_inspection_limit_override", d)
		if err != nil {
			return diag.FromErr(err)
		}
		req.RequestBodyInspectionLimitOverride = override
	}
	_, err = client.UpdateAdvancedSettingsRequestBody(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))

	return resourceAdvancedSettingsRequestBodyRead(ctx, d, m)
}
func resourceAdvancedSettingsRequestBodyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsRequestBodyRead")
	logger.Debugf("in resourceAdvancedSettingsRequestBodyRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	advancedSettingsRequestBody, err := client.GetAdvancedSettingsRequestBody(ctx, appsec.GetAdvancedSettingsRequestBodyRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	})
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsRequestBody': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", policyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("request_body_inspection_limit", advancedSettingsRequestBody.RequestBodyInspectionLimitInKB); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if policyID != "" {
		if err := d.Set("request_body_inspection_limit_override", advancedSettingsRequestBody.RequestBodyInspectionLimitOverride); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	return nil
}

func resourceAdvancedSettingsRequestBodyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsRequestBodyUpdate")
	logger.Debugf("in resourceAdvancedSettingsRequestBodyUpdate")

	return upsertAdvancedSettingsRequestBody(ctx, d, m)
}

func resourceAdvancedSettingsRequestBodyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsRequestBodyDelete")
	logger.Debugf("in resourceAdvancedSettingsRequestBodyDelete")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "requestBodySetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	req := appsec.RemoveAdvancedSettingsRequestBodyRequest{
		ConfigID:                       configID,
		Version:                        version,
		PolicyID:                       policyID,
		RequestBodyInspectionLimitInKB: appsec.Default,
	}

	if policyID != "" {
		req.RequestBodyInspectionLimitOverride = false
	}

	_, err = client.RemoveAdvancedSettingsRequestBody(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
