package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSecurityPolicyRenameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRenameCreate")
	logger.Debugf("!!! in resourceSecurityPolicyRenameCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "securityPolicyRename", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	policyname, err := tools.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createSecurityPolicy := appsec.UpdateSecurityPolicyRequest{}
	createSecurityPolicy.ConfigID = configid
	createSecurityPolicy.Version = version
	createSecurityPolicy.PolicyID = policyid
	createSecurityPolicy.PolicyName = policyname

	_, erru := client.UpdateSecurityPolicy(ctx, createSecurityPolicy)
	if erru != nil {
		logger.Errorf("calling 'createSecurityPolicy': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createSecurityPolicy.ConfigID, createSecurityPolicy.PolicyID))

	return resourceSecurityPolicyRenameRead(ctx, d, m)
}

func resourceSecurityPolicyRenameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRead")
	logger.Debugf("!!! in resourceSecurityPolicyRenameRead")

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
	if err := d.Set("security_policy_id", getSecurityPolicy.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_name", securitypolicy.PolicyName); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceSecurityPolicyRenameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyUpdate")
	logger.Debugf("!!! in resourceSecurityPolicyRenameRead")

	idParts, err := splitID(d.Id(), 2, "configid:policyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "securityPolicyRename", m)
	policyid := idParts[1]
	policyname, err := tools.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateSecurityPolicy := appsec.UpdateSecurityPolicyRequest{}
	updateSecurityPolicy.ConfigID = configid
	updateSecurityPolicy.Version = version
	updateSecurityPolicy.PolicyID = policyid
	updateSecurityPolicy.PolicyName = policyname

	_, erru := client.UpdateSecurityPolicy(ctx, updateSecurityPolicy)
	if erru != nil {
		logger.Errorf("calling 'updateSecurityPolicy': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSecurityPolicyRenameRead(ctx, d, m)
}

func resourceSecurityPolicyRenameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(context.TODO(), d, m)
}
