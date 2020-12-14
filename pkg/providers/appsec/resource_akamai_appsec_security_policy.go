package appsec

import (
	"context"
	"errors"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_policy_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_policy_id": &schema.Schema{
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

	createSecurityPolicy := v2.CreateSecurityPolicyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicy.Version = version

	policyname, err := tools.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicy.PolicyName = policyname

	policyprefix, err := tools.GetStringValue("security_policy_prefix", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicy.PolicyPrefix = policyprefix

	spcr, errc := client.CreateSecurityPolicy(ctx, createSecurityPolicy)
	if errc != nil {
		logger.Errorf("calling 'createSecurityPolicy': %s", errc.Error())
		return diag.FromErr(errc)
	}

	d.Set("security_policy_id", spcr.PolicyID)
	d.Set("security_policy_name", spcr.PolicyName)
	d.SetId(spcr.PolicyID)

	return resourceSecurityPolicyRead(ctx, d, m)
}

func resourceSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyUpdate")

	updateSecurityPolicy := v2.UpdateSecurityPolicyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSecurityPolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSecurityPolicy.Version = version

	policyname, err := tools.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSecurityPolicy.PolicyName = policyname

	policyprefix, err := tools.GetStringValue("security_policy_prefix", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateSecurityPolicy.PolicyPrefix = policyprefix

	updateSecurityPolicy.PolicyID = d.Id()

	_, erru := client.UpdateSecurityPolicy(ctx, updateSecurityPolicy)
	if erru != nil {
		logger.Errorf("calling 'updateSecurityPolicy': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSecurityPolicyRead(ctx, d, m)
}

func resourceSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRemove")

	removeSecurityPolicy := v2.RemoveSecurityPolicyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSecurityPolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSecurityPolicy.Version = version

	removeSecurityPolicy.PolicyID = d.Id()

	_, errd := client.RemoveSecurityPolicy(ctx, removeSecurityPolicy)
	if errd != nil {
		logger.Errorf("calling 'removeSecurityPolicy': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRead")

	getSecurityPolicy := v2.GetSecurityPolicyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSecurityPolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSecurityPolicy.Version = version

	getSecurityPolicy.PolicyID = d.Id()

	securitypolicy, err := client.GetSecurityPolicy(ctx, getSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}
	d.Set("policy_name", securitypolicy.PolicyName)
	d.Set("security_policy_id", securitypolicy.PolicyID)
	d.SetId(securitypolicy.PolicyID)

	return nil
}
