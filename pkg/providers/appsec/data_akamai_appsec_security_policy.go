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

func dataSourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID",
			},
			"policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Policy ID List",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRead")

	getSecurityPolicy := appsec.GetSecurityPoliciesRequest{}

	configName := d.Get("name").(string)

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

	securitypolicy, err := client.GetSecurityPolicies(ctx, getSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	secpolicylist := make([]string, 0, len(securitypolicy.Policies))

	for _, configval := range securitypolicy.Policies {
		secpolicylist = append(secpolicylist, configval.PolicyID)
		if configval.PolicyName == configName {
			d.Set("policy_id", configval.PolicyID)
		}
	}

	d.Set("policy_list", secpolicylist)

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "securityPoliciesDS", securitypolicy)

	if err == nil {
		d.Set("output_text", outputtext)
	}
	if len(securitypolicy.Policies) > 0 {
		if err := d.Set("policy_id", securitypolicy.Policies[0].PolicyID); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	d.SetId(fmt.Sprintf("%d:%d", getSecurityPolicy.ConfigID, getSecurityPolicy.Version))

	return nil
}
