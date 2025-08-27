package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyCreate,
		ReadContext:   resourceSecurityPolicyRead,
		UpdateContext: resourceSecurityPolicyUpdate,
		DeleteContext: resourceSecurityPolicyDelete,
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
			"default_settings": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				DiffSuppressFunc: func(_, oldStr, _ string, _ *schema.ResourceData) bool {
					return oldStr != "" // read & use this setting on policy creation, otherwise ignore it
				},
				Description: "Whether to assign default settings to the new security policy",
			},
			"create_from_security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of the existing security policy being cloned",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of the new security policy",
			},
		},
	}
}

func resourceSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyCreate")
	logger.Debugf("in resourceSecurityPolicyCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "securityPolicy", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyname, err := tf.GetStringValue("security_policy_name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyprefix, err := tf.GetStringValue("security_policy_prefix", d)
	if err != nil {
		return diag.FromErr(err)
	}
	defaultSettings, err := tf.GetBoolValue("default_settings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	createfromsecuritypolicy, err := tf.GetStringValue("create_from_security_policy_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	if len(createfromsecuritypolicy) > 0 {
		createSecurityPolicyClone := appsec.CreateSecurityPolicyCloneRequest{
			ConfigID:                 configID,
			Version:                  version,
			CreateFromSecurityPolicy: createfromsecuritypolicy,
			PolicyName:               policyname,
			PolicyPrefix:             policyprefix,
		}

		spcr, err := client.CreateSecurityPolicyClone(ctx, createSecurityPolicyClone)
		if err != nil {
			logger.Errorf("calling 'createSecurityPolicyClone': %s", err.Error())
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%d:%s", createSecurityPolicyClone.ConfigID, spcr.PolicyID))

	} else {
		createSecurityPolicy := appsec.CreateSecurityPolicyRequest{
			ConfigID:        configID,
			Version:         version,
			PolicyName:      policyname,
			DefaultSettings: defaultSettings,
			PolicyPrefix:    policyprefix,
		}

		spcr, err := client.CreateSecurityPolicy(ctx, createSecurityPolicy)
		if err != nil {
			logger.Errorf("calling 'createSecurityPolicy': %s", err.Error())
			return diag.FromErr(err)
		}
		if err := d.Set("security_policy_id", spcr.PolicyID); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
		if err := d.Set("security_policy_name", spcr.PolicyName); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		d.SetId(fmt.Sprintf("%d:%s", createSecurityPolicy.ConfigID, spcr.PolicyID))
	}

	return resourceSecurityPolicyRead(ctx, d, m)
}

func resourceSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRead")
	logger.Debugf("in resourceSecurityPolicyRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:policyID")
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

	getSecurityPolicy := appsec.GetSecurityPolicyRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	securitypolicy, err := client.GetSecurityPolicy(ctx, getSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getSecurityPolicy.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_name", securitypolicy.PolicyName); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	policyparts := strings.Split(securitypolicy.PolicyID, "_")
	if err := d.Set("security_policy_prefix", policyparts[0]); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", securitypolicy.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("default_settings", true); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyUpdate")
	logger.Debugf("in resourceSecurityPolicyUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:policyID")
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
		return resourceSecurityPolicyRead(ctx, d, m)
	}

	policyname, err := tf.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSecurityPolicy := appsec.UpdateSecurityPolicyRequest{
		ConfigID:   configID,
		Version:    version,
		PolicyID:   securityPolicyID,
		PolicyName: policyname,
	}

	_, err = client.UpdateSecurityPolicy(ctx, updateSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'updateSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceSecurityPolicyRead(ctx, d, m)
}

func resourceSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyDelete")
	logger.Debugf("in resourceSecurityPolicyDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:policyID")
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
		logger.Debugf("latest version %d is active, DeleteContext is a no-op", latestVersion)
	} else {
		removeSecurityPolicy := appsec.RemoveSecurityPolicyRequest{
			ConfigID: configID,
			Version:  version,
			PolicyID: securityPolicyID,
		}

		_, err = client.RemoveSecurityPolicy(ctx, removeSecurityPolicy)
		if err != nil {
			logger.Errorf("calling 'removeSecurityPolicy': %s", err.Error())
			return diag.FromErr(err)
		}
	}

	return nil
}
