package appsec

import (
	"context"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"

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
			"create_from_security_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_id": &schema.Schema{
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

	createSecurityPolicyClone := v2.CreateSecurityPolicyCloneRequest{}

	createSecurityPolicyClone.ConfigID = d.Get("config_id").(int)
	createSecurityPolicyClone.Version = d.Get("version").(int)
	createSecurityPolicyClone.CreateFromSecurityPolicy = d.Get("create_from_security_policy").(string)
	createSecurityPolicyClone.PolicyName = d.Get("policy_name").(string)
	createSecurityPolicyClone.PolicyPrefix = d.Get("policy_prefix").(string)

	spcr, err := client.CreateSecurityPolicyClone(ctx, createSecurityPolicyClone)
	if err != nil {
		logger.Warnf("calling 'createSecurityPolicyClone': %s", err.Error())
	}

	d.Set("policy_id", spcr.PolicyID)
	d.Set("policy_name", spcr.PolicyName)
	d.Set("policy_prefix", createSecurityPolicyClone.PolicyPrefix)
	d.SetId(spcr.PolicyID)

	return resourceSecurityPolicyCloneRead(ctx, d, m)
}

func resourceSecurityPolicyCloneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyCloneRead")

	getSecurityPolicyClone := v2.GetSecurityPolicyCloneRequest{}

	getSecurityPolicyClone.ConfigID = d.Get("config_id").(int)
	getSecurityPolicyClone.Version = d.Get("version").(int)
	getSecurityPolicyClone.PolicyID = d.Id()

	securitypolicyclone, err := client.GetSecurityPolicyClone(ctx, getSecurityPolicyClone)
	if err != nil {
		logger.Warnf("calling 'getSecurityPolicyClone': %s", err.Error())
	}

	d.Set("policy_name", securitypolicyclone.PolicyName)
	d.Set("policy_id", securitypolicyclone.PolicyID)
	d.SetId(securitypolicyclone.PolicyID)

	return nil
}

func resourceSecurityPolicyCloneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}

func resourceSecurityPolicyCloneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}
