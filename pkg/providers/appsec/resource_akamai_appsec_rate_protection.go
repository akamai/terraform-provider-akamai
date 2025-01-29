package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceRateProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRateProtectionCreate,
		ReadContext:   resourceRateProtectionRead,
		UpdateContext: resourceRateProtectionUpdate,
		DeleteContext: resourceRateProtectionDelete,
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
				Type:     schema.TypeBool,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func resourceRateProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionCreate")
	logger.Debugf("in resourceRateProtectionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "rateProtection", m)
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

	request := appsec.UpdateRateProtectionRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ApplyRateControls: enabled,
	}
	_, err = client.UpdateRateProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateRateProtection: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))

	return resourceRateProtectionRead(ctx, d, m)
}

func resourceRateProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionRead")
	logger.Debugf("in resourceRateProtectionRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
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

	request := appsec.GetRateProtectionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	response, err := client.GetRateProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling GetRateProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	enabled := response.ApplyRateControls

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

func resourceRateProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionUpdate")
	logger.Debugf("in resourceRateProtectionUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "rateProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	enabled, err := tf.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := appsec.UpdateRateProtectionRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ApplyRateControls: enabled,
	}
	_, err = client.UpdateRateProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateRateProtection: %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceRateProtectionRead(ctx, d, m)
}

func resourceRateProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRateProtectionDelete")
	logger.Debugf("in resourceRateProtectionDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "rateProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	request := appsec.UpdateRateProtectionRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ApplyRateControls: false,
	}
	_, err = client.UpdateRateProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateRateProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
