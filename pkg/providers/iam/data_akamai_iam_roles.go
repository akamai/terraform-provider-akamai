package iam

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Get roles for the current account and contract",
		ReadContext: dataIAMRolesRead,
		Schema: map[string]*schema.Schema{
			// outputs
			"roles": {
				Type:     schema.TypeSet,
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
					},
				},
			},
		},
	}
}

func dataIAMRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "dataIAMRolesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching roles")

	res, err := client.ListRoles(ctx, iam.ListRolesRequest{})
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

	return m
}
