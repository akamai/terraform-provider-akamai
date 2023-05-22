package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMTimeoutPolicies() *schema.Resource {
	return &schema.Resource{
		Description: "Lists all session timeout policies Akamai supports",
		ReadContext: dataIAMTimeoutPoliciesRead,
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

func dataIAMTimeoutPoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "dataIAMTimeoutPoliciesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching supported timeout policies")

	res, err := client.ListTimeoutPolicies(ctx)
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
