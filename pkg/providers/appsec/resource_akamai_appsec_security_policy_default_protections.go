package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecurityPolicyDefaultProtections() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyDefaultProtectionsCreate,
		ReadContext:   resourceSecurityPolicyDefaultProtectionsRead,
		UpdateContext: resourceSecurityPolicyDefaultProtectionsUpdate,
		DeleteContext: resourceSecurityPolicyDefaultProtectionsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the new security policy",
			},
			"security_policy_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Four-character alphanumeric string prefix used in creating the security policy ID",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of the new security policy",
			},
		},
	}
}

func resourceSecurityPolicyDefaultProtectionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyDefaultProtectionsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "securityPolicy", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyName, err := tf.GetStringValue("security_policy_name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyPrefix, err := tf.GetStringValue("security_policy_prefix", d)
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := client.CreateSecurityPolicyWithDefaultProtections(ctx, appsec.CreateSecurityPolicyWithDefaultProtectionsRequest{
		ConfigVersion: appsec.ConfigVersion{
			ConfigID: int64(configID),
			Version:  version,
		},
		PolicyName:   policyName,
		PolicyPrefix: policyPrefix,
	})
	if err != nil {
		logger.Errorf("calling 'createSecurityPolicyWithDefaultProtections': %s", err.Error())
		return diag.FromErr(err)
	}
	if err := d.Set("security_policy_id", response.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	d.SetId(fmt.Sprintf("%d:%s", configID, response.PolicyID))

	return resourceSecurityPolicyDefaultProtectionsRead(ctx, d, m)
}

func resourceSecurityPolicyDefaultProtectionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyDefaultProtectionsRead")

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

	getSecurityPolicyReq := appsec.GetSecurityPolicyRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	securityPolicy, err := client.GetSecurityPolicy(ctx, getSecurityPolicyReq)
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getSecurityPolicyReq.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_name", securityPolicy.PolicyName); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	policyParts := strings.Split(securityPolicy.PolicyID, "_")
	if err := d.Set("security_policy_prefix", policyParts[0]); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", securityPolicy.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceSecurityPolicyDefaultProtectionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyDefaultProtectionsUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:policyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "securityPolicy", m)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID := iDParts[1]

	// Prevent an update call with the same policy name since API will reject it.
	if !d.HasChange("security_policy_name") {
		logger.Errorf("name in use - please specify a unique security_policy_name that is not used in any version of this security configuration")
		return resourceSecurityPolicyDefaultProtectionsRead(ctx, d, m)
	}

	policyName, err := tf.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSecurityPolicyReq := appsec.UpdateSecurityPolicyRequest{
		ConfigID:   configID,
		Version:    version,
		PolicyID:   securityPolicyID,
		PolicyName: policyName,
	}

	_, err = client.UpdateSecurityPolicy(ctx, updateSecurityPolicyReq)
	if err != nil {
		logger.Errorf("calling 'updateSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}
	if err := d.Set("security_policy_name", policyName); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return resourceSecurityPolicyDefaultProtectionsRead(ctx, d, m)
}

func resourceSecurityPolicyDefaultProtectionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyDefaultProtectionsDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:policyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "securityPolicy", m)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID := iDParts[1]

	latestVersion, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	stagingVersion, productionVersion, err := getActiveConfigVersions(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	if latestVersion == stagingVersion || latestVersion == productionVersion {
		return diag.Errorf("latest version %d is active, DeleteContext is a no-op", latestVersion)
	}

	removeSecurityPolicyReq := appsec.RemoveSecurityPolicyRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: securityPolicyID,
	}

	_, err = client.RemoveSecurityPolicy(ctx, removeSecurityPolicyReq)
	if err != nil {
		logger.Errorf("calling 'removeSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	return nil
}
