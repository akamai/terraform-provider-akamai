package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMGrantableRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Get roles for the current account and contract",
		ReadContext: dataIAMGrantableRolesRead,
		Schema: map[string]*schema.Schema{
			"grantable_roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of grantable roles",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"granted_role_id": {
							Type:        schema.TypeInt,
							Description: "Granted role ID",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Granted role name",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Granted role description",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataIAMGrantableRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "dataIAMGrantableRolesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching grantable roles")

	roles, err := client.ListGrantableRoles(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get grantable roles")
		return diag.FromErr(err)
	}

	grantableRoles := []interface{}{}
	for _, role := range roles {
		grantableRoles = append(grantableRoles, toState(role))
	}

	if err := d.Set("grantable_roles", grantableRoles); err != nil {
		logger.WithError(err).Error("could not set grantable roles")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_grantable_roles")
	return nil
}

func toState(role iam.RoleGrantedRole) map[string]interface{} {
	result := map[string]interface{}{}

	result["granted_role_id"] = role.RoleID
	result["name"] = role.RoleName
	result["description"] = role.Description

	return result
}
