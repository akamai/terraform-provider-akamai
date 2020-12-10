package iam

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Get roles for the current account and contract",
		ReadContext: p.tfCRUD("ds:Roles:Read", p.dsRolesRead),
		Schema: map[string]*schema.Schema{
			// inputs
			"group_id": {
				Type:        schema.TypeString,
				Description: "A unique identifier for a group",
				Optional:    true,
			},
			"get_actions": {
				Type:        schema.TypeBool,
				Description: `When enabled, the response includes information about actions such as "edit" or "delete"`,
				Optional:    true,
			},
			"get_users": {
				Type:        schema.TypeBool,
				Description: "When enabled, returns users assigned to the roles",
				Optional:    true,
			},
			"ignore_context": {
				Type:        schema.TypeBool,
				Description: "When enabled, returns all roles for the current account without regard the contract type associated with your API client",
				Optional:    true,
			},

			// outputs
			"roles": {
				Type:        schema.TypeSet,
				Description: "TODO", // These descriptions were taken from the API docs
				Computed:    true,
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
						"users": {
							Type:        schema.TypeSet,
							Description: "Permissions available to the user for this group", // These descriptions were taken from the API docs
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:        schema.TypeString,
										Description: "A unique identifier for a user's profile",
										Computed:    true,
									},
									"first_name": {
										Type:        schema.TypeString,
										Description: "The user's first name",
										Computed:    true,
									},
									"last_name": {
										Type:        schema.TypeString,
										Description: "The user's last name",
										Computed:    true,
									},
									"account_id": {
										Type:        schema.TypeString,
										Description: "A unique identifier for an account",
										Computed:    true,
									},
									"email": {
										Type:        schema.TypeString,
										Description: "The user's email address",
										Computed:    true,
									},
									"last_login": {
										Type:        schema.TypeString,
										Description: "ISO 8601 timestamp indicating when the user last logged in",
										Computed:    true,
									},
								},
							},
						},
						"granted_roles": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_id": {
										Type:        schema.TypeString,
										Description: "TODO",
										Computed:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "TODO",
										Computed:    true,
									},
									"description": {
										Type:        schema.TypeString,
										Description: "TODO",
										Computed:    true,
									},
								},
							},
						},
						"delete_allowed": {
							Type:        schema.TypeBool,
							Description: "Indicates whether the user can remove items from the group",
							Computed:    true,
						},
						"edit_allowed": {
							Type:        schema.TypeBool,
							Description: "Indicates whether the user can modify items in the group",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (p *provider) dsRolesRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := log.FromContext(ctx)

	getActions := d.Get("get_actions").(bool)
	getUsers := d.Get("get_users").(bool)
	IgnoreContext := d.Get("ignore_context").(bool)

	GroupID, err := strconv.ParseInt(d.Get("group_id").(string), 10, 64)
	if err != nil {
		logger.WithError(err).Error("group_id must be an integer")
		return diag.FromErr(err)
	}

	logger.Debug("Fetching roles")
	req := iam.ListRolesRequest{
		GroupID:       &GroupID,
		Actions:       getActions,
		IgnoreContext: IgnoreContext,
		Users:         getUsers,
	}
	res, err := p.client.ListRoles(ctx, req)
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
	m["users"] = usersToState(r.Users)
	m["granted_roles"] = grantedRolesToState(r.GrantedRoles)

	if r.Actions != nil {
		m["edit_allowed"] = r.Actions.Edit
		m["delete_allowed"] = r.Actions.Delete
	}

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

func usersToState(users []iam.RoleUser) []interface{} {
	var out []interface{}

	for _, r := range users {
		out = append(out, userToState(r))
	}

	return out
}

func userToState(u iam.RoleUser) map[string]interface{} {
	m := map[string]interface{}{}

	m["user_id"] = u.UIIdentityID
	m["first_name"] = u.FirstName
	m["last_name"] = u.LastName
	m["account_id"] = u.AccountID
	m["email"] = u.Email
	m["last_login"] = u.LastLoginDate

	return m
}
