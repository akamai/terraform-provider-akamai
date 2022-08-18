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
// https://techdocs.akamai.com/application-security/reference/api
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
	logger.Debugf("in resourceAdvancedSettingsPragmaHeaderCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
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

	_, err = client.UpdateAdvancedSettingsPragma(ctx, createAdvancedSettingsPragma)
	if err != nil {
		logger.Errorf("calling 'createAdvancedSettingsPragma': %s", err.Error())
		return diag.FromErr(err)
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
	logger.Debug("in resourceAdvancedSettingsPragmaHeaderRead")

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
		version, err := getLatestConfigVersion(ctx, configID, m)
		if err != nil {
			return diag.FromErr(err)
		}
		policyID := iDParts[1]

		getAdvancedSettingsPragma.ConfigID = configID
		getAdvancedSettingsPragma.Version = version
		getAdvancedSettingsPragma.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getLatestConfigVersion(ctx, configID, m)
		if err != nil {
			return diag.FromErr(err)
		}

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
	logger.Debugf("in resourceAdvancedSettingsPragmaHeaderDelete")

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
		version, err := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}
		policyID := iDParts[1]

		removeAdvancedSettingsPragma.ConfigID = configID
		removeAdvancedSettingsPragma.Version = version
		removeAdvancedSettingsPragma.PolicyID = policyID

	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}

		removeAdvancedSettingsPragma.ConfigID = configID
		removeAdvancedSettingsPragma.Version = version
	}

	_, err := client.UpdateAdvancedSettingsPragma(ctx, removeAdvancedSettingsPragma)
	if err != nil {
		logger.Errorf("calling 'removeAdvancedSettingsLogging': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}

func resourceAdvancedSettingsPragmaHeaderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPragmaHeaderUpdate")
	logger.Debugf("in resourceAdvancedSettingsPragmaHeaderUpdate")

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
		version, err := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}

		policyID := iDParts[1]

		updateAdvancedSettingsPragma.ConfigID = configID
		updateAdvancedSettingsPragma.Version = version
		updateAdvancedSettingsPragma.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "pragmaSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}
		updateAdvancedSettingsPragma.ConfigID = configID
		updateAdvancedSettingsPragma.Version = version
	}

	jsonpostpayload := d.Get("pragma_header")

	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateAdvancedSettingsPragma.JsonPayloadRaw = rawJSON
	_, err := client.UpdateAdvancedSettingsPragma(ctx, updateAdvancedSettingsPragma)
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPragma': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAdvancedSettingsPragmaHeaderRead(ctx, d, m)
}
