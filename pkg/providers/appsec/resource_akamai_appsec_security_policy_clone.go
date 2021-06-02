package appsec

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSecurityPolicyClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyCloneCreate,
		ReadContext:   resourceSecurityPolicyCloneRead,
		UpdateContext: resourceSecurityPolicyCloneUpdate,
		DeleteContext: resourceSecurityPolicyCloneDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"create_from_security_policy_id": {
				Type:     schema.TypeString,
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
				Description: "Policy ID for clone",
			},
		},
	}
}

func resourceSecurityPolicyCloneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyCloneCreate")

	createSecurityPolicyClone := appsec.CreateSecurityPolicyCloneRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicyClone.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicyClone.Version = version

	createfromsecuritypolicy, err := tools.GetStringValue("create_from_security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicyClone.CreateFromSecurityPolicy = createfromsecuritypolicy

	policyname, err := tools.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicyClone.PolicyName = policyname

	policyprefix, err := tools.GetStringValue("security_policy_prefix", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createSecurityPolicyClone.PolicyPrefix = policyprefix

	spcr, err := client.CreateSecurityPolicyClone(ctx, createSecurityPolicyClone)
	if err != nil {
		logger.Errorf("calling 'createSecurityPolicyClone': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("security_policy_id", spcr.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_name", spcr.PolicyName); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_prefix", createSecurityPolicyClone.PolicyPrefix); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(spcr.PolicyID)

	return resourceSecurityPolicyCloneRead(ctx, d, m)
}

func resourceSecurityPolicyCloneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyCloneRead")

	getSecurityPolicyClone := appsec.GetSecurityPolicyCloneRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSecurityPolicyClone.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSecurityPolicyClone.Version = version

	if d.HasChange("version") {
		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSecurityPolicyClone.Version = version
	}

	getSecurityPolicyClone.PolicyID = d.Id()

	securitypolicyclone, err := client.GetSecurityPolicyClone(ctx, getSecurityPolicyClone)
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicyClone': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("security_policy_name", securitypolicyclone.PolicyName); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", securitypolicyclone.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(securitypolicyclone.PolicyID)

	return nil
}

func resourceSecurityPolicyCloneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyCloneRead")

	removeSecurityPolicyClone := appsec.RemoveSecurityPolicyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSecurityPolicyClone.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeSecurityPolicyClone.Version = version

	removeSecurityPolicyClone.PolicyID = d.Id()

	_, errd := client.RemoveSecurityPolicy(ctx, removeSecurityPolicyClone)
	if errd != nil {
		logger.Errorf("calling 'removeSecurityPolicyClone': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceSecurityPolicyCloneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}
