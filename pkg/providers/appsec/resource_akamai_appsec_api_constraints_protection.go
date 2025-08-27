package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceAPIConstraintsProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIConstraintsProtectionCreate,
		ReadContext:   resourceAPIConstraintsProtectionRead,
		UpdateContext: resourceAPIConstraintsProtectionUpdate,
		DeleteContext: resourceAPIConstraintsProtectionDelete,
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
				Description: "Whether to enable API constraints protection",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func resourceAPIConstraintsProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionCreate")
	logger.Debugf("in resourceAPIConstraintsProtectionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "apiConstraintsProtection", m)
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

	request := appsec.UpdateAPIConstraintsProtectionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ApplyAPIConstraints: enabled,
	}
	_, err = client.UpdateAPIConstraintsProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateAPIConstraints: %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))
	return resourceAPIConstraintsProtectionRead(ctx, d, m)
}

func resourceAPIConstraintsProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionRead")
	logger.Debugf("in resourceAPIConstraintsProtectionRead")

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

	request := appsec.GetAPIConstraintsProtectionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	response, err := client.GetAPIConstraintsProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling GetAPIConstraintsProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	enabled := response.ApplyAPIConstraints

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

func resourceAPIConstraintsProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionUpdate")
	logger.Debugf("in resourceAPIConstraintsProtectionUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "apiConstraintsProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	enabled, err := tf.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := appsec.UpdateAPIConstraintsProtectionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ApplyAPIConstraints: enabled,
	}
	_, err = client.UpdateAPIConstraintsProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateAPIConstraintsProtection: %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAPIConstraintsProtectionRead(ctx, d, m)
}

func resourceAPIConstraintsProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIConstraintsProtectionDelete")
	logger.Debugf("in resourceAPIConstraintsProtectionDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "apiConstraintsProtection", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	request := appsec.UpdateAPIConstraintsProtectionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ApplyAPIConstraints: false,
	}
	_, err = client.UpdateAPIConstraintsProtection(ctx, request)
	if err != nil {
		logger.Errorf("calling UpdateAPIConstraintsProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
