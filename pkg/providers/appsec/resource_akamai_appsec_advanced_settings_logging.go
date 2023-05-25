package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceAdvancedSettingsLogging() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsLoggingCreate,
		ReadContext:   resourceAdvancedSettingsLoggingRead,
		UpdateContext: resourceAdvancedSettingsLoggingUpdate,
		DeleteContext: resourceAdvancedSettingsLoggingDelete,
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
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of the security policy",
			},
			"logging": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentLoggingSettingsDiffs,
				Description:      "Whether to enable, disable, or update HTTP header logging settings",
			},
		},
	}
}

func resourceAdvancedSettingsLoggingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsLoggingCreate")
	logger.Debugf("in resourceAdvancedSettingsLoggingCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "loggingSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	jsonpostpayload := d.Get("logging")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createAdvancedSettingsLogging := appsec.UpdateAdvancedSettingsLoggingRequest{
		ConfigID:       configID,
		Version:        version,
		PolicyID:       policyID,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateAdvancedSettingsLogging(ctx, createAdvancedSettingsLogging)
	if err != nil {
		logger.Errorf("calling 'UpdateAdvancedSettingsLogging': %s", err.Error())
		return diag.FromErr(err)
	}

	if len(createAdvancedSettingsLogging.PolicyID) > 0 {
		d.SetId(fmt.Sprintf("%d:%s", createAdvancedSettingsLogging.ConfigID, createAdvancedSettingsLogging.PolicyID))
	} else {
		d.SetId(fmt.Sprintf("%d", createAdvancedSettingsLogging.ConfigID))
	}

	return resourceAdvancedSettingsLoggingRead(ctx, d, m)
}

func resourceAdvancedSettingsLoggingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsLoggingRead")
	logger.Debugf("in resourceAdvancedSettingsLoggingRead")

	getAdvancedSettingsLogging := appsec.GetAdvancedSettingsLoggingRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
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
		policyID := iDParts[1]

		getAdvancedSettingsLogging.ConfigID = configID
		getAdvancedSettingsLogging.Version = version
		getAdvancedSettingsLogging.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getLatestConfigVersion(ctx, configID, m)
		if err != nil {
			return diag.FromErr(err)
		}

		getAdvancedSettingsLogging.ConfigID = configID
		getAdvancedSettingsLogging.Version = version
	}

	advancedsettingslogging, err := client.GetAdvancedSettingsLogging(ctx, getAdvancedSettingsLogging)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsLogging': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAdvancedSettingsLogging.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getAdvancedSettingsLogging.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(advancedsettingslogging)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("logging", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAdvancedSettingsLoggingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsLoggingUpdate")
	logger.Debugf("in resourceAdvancedSettingsLoggingUpdate")

	updateAdvancedSettingsLogging := appsec.UpdateAdvancedSettingsLoggingRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
		if err != nil {
			return diag.FromErr(err)
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "loggingSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}
		policyID := iDParts[1]

		updateAdvancedSettingsLogging.ConfigID = configID
		updateAdvancedSettingsLogging.Version = version
		updateAdvancedSettingsLogging.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "loggingSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}

		updateAdvancedSettingsLogging.ConfigID = configID
		updateAdvancedSettingsLogging.Version = version
	}

	jsonpostpayload := d.Get("logging")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateAdvancedSettingsLogging.JsonPayloadRaw = rawJSON
	_, err := client.UpdateAdvancedSettingsLogging(ctx, updateAdvancedSettingsLogging)
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsLogging': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAdvancedSettingsLoggingRead(ctx, d, m)
}

func resourceAdvancedSettingsLoggingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsLoggingDelete")
	logger.Debugf("in resourceAdvancedSettingsLoggingDelete")

	removeAdvancedSettingsLogging := appsec.RemoveAdvancedSettingsLoggingRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
		if err != nil {
			return diag.FromErr(err)
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "loggingSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}
		policyID := iDParts[1]

		removeAdvancedSettingsLogging.ConfigID = configID
		removeAdvancedSettingsLogging.Version = version
		removeAdvancedSettingsLogging.PolicyID = policyID
		removeAdvancedSettingsLogging.Override = false
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "loggingSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}

		removeAdvancedSettingsLogging.ConfigID = configID
		removeAdvancedSettingsLogging.Version = version
		removeAdvancedSettingsLogging.AllowSampling = false
	}

	_, err := client.RemoveAdvancedSettingsLogging(ctx, removeAdvancedSettingsLogging)
	if err != nil {
		logger.Errorf("calling 'removeAdvancedSettingsLogging': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
