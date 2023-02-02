package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceAdvancedSettingsEvasivePathMatch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsEvasivePathMatchCreate,
		ReadContext:   resourceAdvancedSettingsEvasivePathMatchRead,
		UpdateContext: resourceAdvancedSettingsEvasivePathMatchUpdate,
		DeleteContext: resourceAdvancedSettingsEvasivePathMatchDelete,
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
			"enable_path_match": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable the evasive path match setting",
			},
		},
	}
}

func resourceAdvancedSettingsEvasivePathMatchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsEvasivePathMatchCreate")
	logger.Debugf("in resourceAdvancedSettingsEvasivePathMatchCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evasivePathMatchSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enablePathMatch, err := tools.GetBoolValue("enable_path_match", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createAdvancedSettingsEvasivePathMatch := appsec.UpdateAdvancedSettingsEvasivePathMatchRequest{
		ConfigID:        configID,
		Version:         version,
		PolicyID:        policyID,
		EnablePathMatch: enablePathMatch,
	}

	_, err = client.UpdateAdvancedSettingsEvasivePathMatch(ctx, createAdvancedSettingsEvasivePathMatch)
	if err != nil {
		logger.Errorf("calling 'createAdvancedSettingsEvasivePathMatch': %s", err.Error())
		return diag.FromErr(err)
	}

	if len(createAdvancedSettingsEvasivePathMatch.PolicyID) > 0 {
		d.SetId(fmt.Sprintf("%d:%s", createAdvancedSettingsEvasivePathMatch.ConfigID, createAdvancedSettingsEvasivePathMatch.PolicyID))
	} else {
		d.SetId(fmt.Sprintf("%d", createAdvancedSettingsEvasivePathMatch.ConfigID))
	}

	return resourceAdvancedSettingsEvasivePathMatchRead(ctx, d, m)
}

func resourceAdvancedSettingsEvasivePathMatchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsEvasivePathMatchRead")
	logger.Debugf("in resourceAdvancedSettingsEvasivePathMatchRead")

	getAdvancedSettingsEvasivePathMatch := appsec.GetAdvancedSettingsEvasivePathMatchRequest{}
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

		getAdvancedSettingsEvasivePathMatch.ConfigID = configID
		getAdvancedSettingsEvasivePathMatch.Version = version
		getAdvancedSettingsEvasivePathMatch.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getLatestConfigVersion(ctx, configID, m)
		if err != nil {
			return diag.FromErr(err)
		}

		getAdvancedSettingsEvasivePathMatch.ConfigID = configID
		getAdvancedSettingsEvasivePathMatch.Version = version
	}

	advancedsettingsevasivepathmatch, err := client.GetAdvancedSettingsEvasivePathMatch(ctx, getAdvancedSettingsEvasivePathMatch)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsEvasivePathMatch': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAdvancedSettingsEvasivePathMatch.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getAdvancedSettingsEvasivePathMatch.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("enable_path_match", advancedsettingsevasivepathmatch.EnablePathMatch); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAdvancedSettingsEvasivePathMatchUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsEvasivePathMatchUpdate")
	logger.Debugf("in resourceAdvancedSettingsEvasivePathMatchUpdate")

	updateAdvancedSettingsEvasivePathMatch := appsec.UpdateAdvancedSettingsEvasivePathMatchRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
		if err != nil {
			return diag.FromErr(err)
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "evasivePathMatchSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}
		policyID := iDParts[1]

		updateAdvancedSettingsEvasivePathMatch.ConfigID = configID
		updateAdvancedSettingsEvasivePathMatch.Version = version
		updateAdvancedSettingsEvasivePathMatch.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "evasivePathMatchSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}

		updateAdvancedSettingsEvasivePathMatch.ConfigID = configID
		updateAdvancedSettingsEvasivePathMatch.Version = version
	}
	enablePathMatch, err := tools.GetBoolValue("enable_path_match", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAdvancedSettingsEvasivePathMatch.EnablePathMatch = enablePathMatch

	_, err = client.UpdateAdvancedSettingsEvasivePathMatch(ctx, updateAdvancedSettingsEvasivePathMatch)
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsEvasivePathMatch': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAdvancedSettingsEvasivePathMatchRead(ctx, d, m)
}

func resourceAdvancedSettingsEvasivePathMatchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsEvasivePathMatchDelete")
	logger.Debugf("in resourceAdvancedSettingsEvasivePathMatchDelete")

	removeAdvancedSettingsEvasivePathMatch := appsec.RemoveAdvancedSettingsEvasivePathMatchRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		iDParts, err := splitID(d.Id(), 2, "configID:policyID")
		if err != nil {
			return diag.FromErr(err)
		}
		configID, err := strconv.Atoi(iDParts[0])
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "evasivePathMatchSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}
		policyID := iDParts[1]

		removeAdvancedSettingsEvasivePathMatch.ConfigID = configID
		removeAdvancedSettingsEvasivePathMatch.Version = version
		removeAdvancedSettingsEvasivePathMatch.PolicyID = policyID
	} else {
		configID, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		version, err := getModifiableConfigVersion(ctx, configID, "evasivePathMatchSetting", m)
		if err != nil {
			return diag.FromErr(err)
		}

		removeAdvancedSettingsEvasivePathMatch.ConfigID = configID
		removeAdvancedSettingsEvasivePathMatch.Version = version
	}

	removeAdvancedSettingsEvasivePathMatch.EnablePathMatch = false

	_, err := client.RemoveAdvancedSettingsEvasivePathMatch(ctx, removeAdvancedSettingsEvasivePathMatch)
	if err != nil {
		logger.Errorf("calling 'removeAdvancedSettingsEvasivePathMatch': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
