package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
func resourceAdvancedSettingsPragmaHeader() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsPragmaHeaderCreate,
		ReadContext:   resourceAdvancedSettingsPragmaHeaderRead,
		UpdateContext: resourceAdvancedSettingsPragmaHeaderUpdate,
		DeleteContext: resourceAdvancedSettingsPragmaHeaderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pragma_header": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceAdvancedSettingsPragmaHeaderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPragmaHeaderCreate")
	logger.Debug("in resourceAdvancedSettingsPragmaHeaderCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)

	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	jsonpostpayload := d.Get("pragma_header")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createAdvancedSettingsPragma := appsec.UpdateAdvancedSettingsPragmaRequest{
		ConfigID:       configID,
		Version:        version,
		PolicyID:       policyID,
		JsonPayloadRaw: rawJSON,
	}

	_, erru := client.UpdateAdvancedSettingsPragma(ctx, createAdvancedSettingsPragma)
	if erru != nil {
		logger.Errorf("calling 'createAdvancedSettingsPragma': %s", erru.Error())
		return diag.FromErr(erru)
	}

	if len(createAdvancedSettingsPragma.PolicyID) > 0 {
		d.SetId(fmt.Sprintf("%d:%s", createAdvancedSettingsPragma.ConfigID, createAdvancedSettingsPragma.PolicyID))
	} else {
		d.SetId(fmt.Sprintf("%d", createAdvancedSettingsPragma.ConfigID))
	}

	return resourceAdvancedSettingsPragmaHeaderRead(ctx, d, m)
}

func resourceAdvancedSettingsPragmaHeaderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPragmaHeaderRead")
	logger.Debug("resourceAdvancedSettingsPragmaHeaderRead")

	getAdvancedSettingsPragma := appsec.GetAdvancedSettingsPragmaRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
		if err != nil {
			return diag.FromErr(err)
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return diag.FromErr(err)
		}
		version := getLatestConfigVersion(ctx, configID, m)
		policyID := iDParts[1]

		getAdvancedSettingsPragma.ConfigID = configID
		getAdvancedSettingsPragma.Version = version
		getAdvancedSettingsPragma.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version := getLatestConfigVersion(ctx, configID, m)

		getAdvancedSettingsPragma.ConfigID = configID
		getAdvancedSettingsPragma.Version = version
	}

	advancedsettingspragma, err := client.GetAdvancedSettingsPragma(ctx, getAdvancedSettingsPragma)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsPragmaRead': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAdvancedSettingsPragma.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getAdvancedSettingsPragma.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(advancedsettingspragma)

	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("pragma_header", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAdvancedSettingsPragmaHeaderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPragmaHeaderDelete")
	logger.Debug("resourceAdvancedSettingsPragmaHeaderDelete")

	jsonPayloadRaw := []byte("{}")
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	removeAdvancedSettingsPragma := appsec.UpdateAdvancedSettingsPragmaRequest{
		JsonPayloadRaw: rawJSON,
	}

	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
		if err != nil {
			return diag.FromErr(err)
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return diag.FromErr(err)
		}
		version := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)
		policyID := iDParts[1]

		removeAdvancedSettingsPragma.ConfigID = configID
		removeAdvancedSettingsPragma.Version = version
		removeAdvancedSettingsPragma.PolicyID = policyID

	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)

		removeAdvancedSettingsPragma.ConfigID = configID
		removeAdvancedSettingsPragma.Version = version
	}

	_, erru := client.UpdateAdvancedSettingsPragma(ctx, removeAdvancedSettingsPragma)
	if erru != nil {
		logger.Errorf("calling 'removeAdvancedSettingsLogging': %s", erru.Error())
		return diag.FromErr(erru)
	}
	d.SetId("")
	return nil
}

func resourceAdvancedSettingsPragmaHeaderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPragmaHeaderUpdate")
	logger.Debugf("resourceAdvancedSettingsPragmaHeaderUpdate")

	updateAdvancedSettingsPragma := appsec.UpdateAdvancedSettingsPragmaRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
		if err != nil {
			return diag.FromErr(err)
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return diag.FromErr(err)
		}
		version := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)

		policyID := iDParts[1]

		updateAdvancedSettingsPragma.ConfigID = configID
		updateAdvancedSettingsPragma.Version = version
		updateAdvancedSettingsPragma.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)

		updateAdvancedSettingsPragma.ConfigID = configID
		updateAdvancedSettingsPragma.Version = version
	}

	jsonpostpayload := d.Get("pragma_header")

	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateAdvancedSettingsPragma.JsonPayloadRaw = rawJSON
	_, erru := client.UpdateAdvancedSettingsPragma(ctx, updateAdvancedSettingsPragma)
	if erru != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPragma': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceAdvancedSettingsPragmaHeaderRead(ctx, d, m)
}
