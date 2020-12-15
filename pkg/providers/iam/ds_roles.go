package iam

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Get roles for the current account and contract",
		ReadContext: p.tfCRUD("ds:Roles:Read", p.dsRolesRead),
		Schema: map[string]*schema.Schema{
			// outputs
			"roles": {
				Type: schema.TypeSet,
				// Description: "TODO", // These descriptions were taken from the API docs
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:        schema.TypeString,
							Description: "A unique identifier for each role",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "The role's name",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The role's description",
							Computed:    true,
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Whether the role is a standard role or a custom role",
							Computed:    true,
						},
						"time_created": {
							Type:        schema.TypeString,
							Description: "ISO 8601 timestamp indicating when the role was originally created",
							Computed:    true,
						},
						"time_modified": {
							Type:        schema.TypeString,
							Description: "ISO 8601 timestamp indicating when the role was last updated",
							Computed:    true,
						},
						"modified_by": {
							Type:        schema.TypeString,
							Description: "The username or email of the last person to edit the role",
							Computed:    true,
						},
						"created_by": {
							Type:        schema.TypeString,
							Description: "The user name or email of the person who created the role",
							Computed:    true,
						},
						"granted_roles": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_id": {
										Type: schema.TypeString,
										// Description: "TODO",
										Computed: true,
									},
									"name": {
										Type: schema.TypeString,
										// Description: "TODO",
										Computed: true,
									},
									"description": {
										Type: schema.TypeString,
										// Description: "TODO",
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (p *provider) dsRolesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := p.log(ctx)

	logger.Debug("Fetching roles")
	res, err := p.client.ListRoles(ctx, iam.ListRolesRequest{})
	if err != nil {
		logger.WithError(err).Error("Could not get roles")
		return diag.FromErr(err)
	}

	if err := d.Set("roles", rolesToState(res)); err != nil {
		logger.WithError(err).Error("Could not set roles in state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_roles")
	return nil
}

func rolesToState(roles []iam.Role) []interface{} {
	var out []interface{}

	for _, r := range roles {
		out = append(out, roleToState(r))
	}

	return out
}

func roleToState(r iam.Role) map[string]interface{} {
	m := map[string]interface{}{}

	m["role_id"] = strconv.FormatInt(r.RoleID, 10)
	m["name"] = r.RoleName
	m["description"] = r.RoleDescription
	m["type"] = string(r.RoleType)
	m["time_created"] = r.CreatedDate
	m["time_modified"] = r.ModifiedDate
	m["modified_by"] = r.ModifiedBy
	m["created_by"] = r.CreatedBy
	m["granted_roles"] = grantedRolesToState(r.GrantedRoles)

	return m
}

func grantedRolesToState(roles []iam.RoleGrantedRole) []interface{} {
	var out []interface{}

	for _, r := range roles {
		out = append(out, grantedRoleToState(r))
	}

	return out
}

func grantedRoleToState(r iam.RoleGrantedRole) map[string]interface{} {
	m := map[string]interface{}{}

	m["name"] = r.RoleName
	m["role_id"] = strconv.FormatInt(r.RoleID, 10)
	m["description"] = r.Description

	return m
}
