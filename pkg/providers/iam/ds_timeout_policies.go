package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsTimeoutPolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all session timeout policies Akamai supports",
		ReadContext: p.tfCRUD("ds:TimeoutPolicies:Read", p.dsTimeoutPoliciesRead),
		Schema: map[string]*schema.Schema{
			"policies": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Session timeout policies",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func (p *provider) dsTimeoutPoliciesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	logger.Debug("Fetching supported timeout policies")
	res, err := p.client.ListTimeoutPolicies(ctx)
	if err != nil {
		logger.WithError(err).Error("Could not get supported timeout policies")
		return diag.FromErr(err)
	}

	policies := map[string]interface{}{}
	for _, policy := range res {
		policies[policy.Name] = policy.Value
	}

	if err := d.Set("policies", policies); err != nil {
		logger.WithError(err).Error("Could not set timeout policies in state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_timeout_policies")
	return nil
}
