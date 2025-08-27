package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
func resourceSecurityPolicyRename() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyRenameCreate,
		ReadContext:   resourceSecurityPolicyRenameRead,
		UpdateContext: resourceSecurityPolicyRenameUpdate,
		DeleteContext: resourceSecurityPolicyRenameDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
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
			"security_policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "New name to be given to the security policy",
			},
		},
	}
}

func resourceSecurityPolicyRenameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRenameCreate")
	logger.Debugf("in resourceSecurityPolicyRenameCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "securityPolicyRename", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyname, err := tf.GetStringValue("security_policy_name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createSecurityPolicy := appsec.UpdateSecurityPolicyRequest{
		ConfigID:   configID,
		Version:    version,
		PolicyID:   policyID,
		PolicyName: policyname,
	}

	_, err = client.UpdateSecurityPolicy(ctx, createSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'createSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createSecurityPolicy.ConfigID, createSecurityPolicy.PolicyID))

	return resourceSecurityPolicyRenameRead(ctx, d, m)
}

func resourceSecurityPolicyRenameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRenameRead")
	logger.Debugf("in resourceSecurityPolicyRenameRead")

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
	if err := d.Set("security_policy_id", getSecurityPolicy.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_name", securitypolicy.PolicyName); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceSecurityPolicyRenameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRenameUpdate")
	logger.Debugf("in resourceSecurityPolicyRenameUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:policyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "securityPolicyRename", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	policyname, err := tf.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSecurityPolicy := appsec.UpdateSecurityPolicyRequest{
		ConfigID:   configID,
		Version:    version,
		PolicyID:   policyID,
		PolicyName: policyname,
	}

	_, err = client.UpdateSecurityPolicy(ctx, updateSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'updateSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceSecurityPolicyRenameRead(ctx, d, m)
}

func resourceSecurityPolicyRenameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(ctx, d, m)
}
