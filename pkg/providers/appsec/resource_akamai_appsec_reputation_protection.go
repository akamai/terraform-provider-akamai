package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceReputationProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProtectionCreate,
		ReadContext:   resourceReputationProtectionRead,
		UpdateContext: resourceReputationProtectionUpdate,
		DeleteContext: resourceReputationProtectionDelete,
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
				Description: "Whether to enable reputation protection",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func resourceReputationProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionCreate")
	logger.Debugf("in resourceReputationProtectionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	enabled, err := tf.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := appsec.UpdateReputationProtectionRequest{
		ConfigID:                configID,
		Version:                 version,
		PolicyID:                policyID,
		ApplyReputationControls: enabled,
	}
	_, err = client.UpdateReputationProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateReputationProtection: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))

	return resourceReputationProtectionRead(ctx, d, m)
}

func resourceReputationProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionRead")
	logger.Debugf("in resourceReputationProtectionRead")

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

	request := appsec.GetReputationProtectionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	response, err := client.GetReputationProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling GetReputationProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	enabled := response.ApplyReputationControls

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

func resourceReputationProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionUpdate")
	logger.Debugf("in resourceReputationProtectionUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	enabled, err := tf.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := appsec.UpdateReputationProtectionRequest{
		ConfigID:                configID,
		Version:                 version,
		PolicyID:                policyID,
		ApplyReputationControls: enabled,
	}
	_, err = client.UpdateReputationProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateReputationProtection: %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceReputationProtectionRead(ctx, d, m)
}

func resourceReputationProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProtectionDelete")
	logger.Debugf("in resourceReputationProtectionDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	request := appsec.UpdateReputationProtectionRequest{
		ConfigID:                configID,
		Version:                 version,
		PolicyID:                policyID,
		ApplyReputationControls: false,
	}
	_, err = client.UpdateReputationProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateReputationProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
