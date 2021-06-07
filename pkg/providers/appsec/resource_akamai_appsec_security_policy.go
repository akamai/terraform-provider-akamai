package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyCreate,
		ReadContext:   resourceSecurityPolicyRead,
		UpdateContext: resourceSecurityPolicyUpdate,
		DeleteContext: resourceSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_policy_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_settings": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_from_security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID for new policy",
			},
		},
	}
}

func resourceSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyCreate")
	logger.Debugf("!!! in resourceSecurityPolicyCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "securityPolicy", m)
	policyname, err := tools.GetStringValue("security_policy_name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyprefix, err := tools.GetStringValue("security_policy_prefix", d)
	if err != nil {
		return diag.FromErr(err)
	}
	defaultSettings, err := tools.GetBoolValue("default_settings", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createfromsecuritypolicy, err := tools.GetStringValue("create_from_security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if len(createfromsecuritypolicy) > 0 {
		createSecurityPolicyClone := appsec.CreateSecurityPolicyCloneRequest{}
		createSecurityPolicyClone.ConfigID = configid
		createSecurityPolicyClone.Version = version
		createSecurityPolicyClone.CreateFromSecurityPolicy = createfromsecuritypolicy
		createSecurityPolicyClone.PolicyName = policyname
		createSecurityPolicyClone.PolicyPrefix = policyprefix

		spcr, err := client.CreateSecurityPolicyClone(ctx, createSecurityPolicyClone)
		if err != nil {
			logger.Errorf("calling 'createSecurityPolicyClone': %s", err.Error())
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%d:%s", createSecurityPolicyClone.ConfigID, spcr.PolicyID))

	} else {
		createSecurityPolicy := appsec.CreateSecurityPolicyRequest{}
		createSecurityPolicy.ConfigID = configid
		createSecurityPolicy.Version = version
		createSecurityPolicy.PolicyName = policyname
		createSecurityPolicy.DefaultSettings = defaultSettings
		createSecurityPolicy.PolicyPrefix = policyprefix

		spcr, errc := client.CreateSecurityPolicy(ctx, createSecurityPolicy)

		if errc != nil {
			logger.Errorf("calling 'createSecurityPolicy': %s", errc.Error())
			return diag.FromErr(errc)
		}
		if err := d.Set("security_policy_id", spcr.PolicyID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		if err := d.Set("security_policy_name", spcr.PolicyName); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		d.SetId(fmt.Sprintf("%d:%s", createSecurityPolicy.ConfigID, spcr.PolicyID))
	}

	return resourceSecurityPolicyRead(ctx, d, m)
}

func resourceSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRead")
	logger.Debugf("!!! in resourceSecurityPolicyRead")

	idParts, err := splitID(d.Id(), 2, "configid:policyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]

	getSecurityPolicy := appsec.GetSecurityPolicyRequest{}
	getSecurityPolicy.ConfigID = configid
	getSecurityPolicy.Version = version
	getSecurityPolicy.PolicyID = policyid

	securitypolicy, err := client.GetSecurityPolicy(ctx, getSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getSecurityPolicy.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_name", securitypolicy.PolicyName); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	policyparts := strings.Split(securitypolicy.PolicyID, "_")
	if err := d.Set("security_policy_prefix", policyparts[0]); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", securitypolicy.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("default_settings", true); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyUpdate")
	logger.Debugf("!!! in resourceSecurityPolicyUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:policyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "securityPolicy", m)
	securitypolicyid := idParts[1]

	// Prevent an update call with the same policy name since API will reject it.
	if !d.HasChange("security_policy_name") {
		return resourceSecurityPolicyRead(ctx, d, m)
	}

	policyname, err := tools.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSecurityPolicy := appsec.UpdateSecurityPolicyRequest{}
	updateSecurityPolicy.ConfigID = configid
	updateSecurityPolicy.Version = version
	updateSecurityPolicy.PolicyID = securitypolicyid
	updateSecurityPolicy.PolicyName = policyname

	_, err = client.UpdateSecurityPolicy(ctx, updateSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'updateSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceSecurityPolicyRead(ctx, d, m)
}

func resourceSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyDelete")
	logger.Debugf("!!! in resourceSecurityPolicyDelete")

	idParts, err := splitID(d.Id(), 2, "configid:policyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "securityPolicy", m)
	securitypolicyid := idParts[1]

	latestVersion := getLatestConfigVersion(ctx, configid, m)
	stagingVersion, productionVersion := getActiveConfigVersions(ctx, configid, m)
	if latestVersion == stagingVersion || latestVersion == productionVersion {
		logger.Debugf("latest version %d is active, DeleteContext is a no-op", latestVersion)
	} else {
		removeSecurityPolicy := appsec.RemoveSecurityPolicyRequest{}
		removeSecurityPolicy.ConfigID = configid
		removeSecurityPolicy.Version = version
		removeSecurityPolicy.PolicyID = securitypolicyid
		_, errd := client.RemoveSecurityPolicy(ctx, removeSecurityPolicy)
		if errd != nil {
			logger.Errorf("calling 'removeSecurityPolicy': %s", errd.Error())
			return diag.FromErr(errd)
		}

		d.SetId("")
	}

	return nil
}
