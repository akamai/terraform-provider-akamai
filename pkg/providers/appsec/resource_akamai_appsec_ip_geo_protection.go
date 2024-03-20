package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceIPGeoProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPGeoProtectionCreate,
		ReadContext:   resourceIPGeoProtectionRead,
		UpdateContext: resourceIPGeoProtectionUpdate,
		DeleteContext: resourceIPGeoProtectionDelete,
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
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable IP/Geo protection",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func resourceIPGeoProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionCreate")
	logger.Debugf("in resourceIPGeoProtectionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ipgeoProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	enabled, err := tf.GetBoolValue("enabled", d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateIPGeoProtection(ctx, appsec.UpdateIPGeoProtectionRequest{
		ConfigID:                  configID,
		Version:                   version,
		PolicyID:                  policyID,
		ApplyNetworkLayerControls: enabled,
	})
	if err != nil {
		logger.Errorf("calling UpdateIPGeoProtection: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))
	return resourceIPGeoProtectionRead(ctx, d, m)
}

func resourceIPGeoProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionRead")
	logger.Debugf("in resourceIPGeoProtectionRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
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

	response, err := client.GetIPGeoProtection(ctx, appsec.GetIPGeoProtectionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	})
	if err != nil {
		logger.Errorf("calling GetIPGeoProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	enabled := response.ApplyNetworkLayerControls

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", policyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("enabled", enabled); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "protections", response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceIPGeoProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionUpdate")
	logger.Debugf("in resourceIPGeoProtectionUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "networkProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	enabled, err := tf.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	_, err = client.UpdateIPGeoProtection(ctx, appsec.UpdateIPGeoProtectionRequest{
		ConfigID:                  configID,
		Version:                   version,
		PolicyID:                  policyID,
		ApplyNetworkLayerControls: enabled,
	})
	if err != nil {
		logger.Errorf("calling UpdateIPGeoProtection: %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceIPGeoProtectionRead(ctx, d, m)
}

func resourceIPGeoProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionDelete")
	logger.Debugf("in resourceIPGeoProtectionDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ipgeoProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	_, err = client.UpdateIPGeoProtection(ctx, appsec.UpdateIPGeoProtectionRequest{
		ConfigID:                  configID,
		Version:                   version,
		PolicyID:                  policyID,
		ApplyNetworkLayerControls: false,
	})
	if err != nil {
		logger.Errorf("calling UpdateIPGeoProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
