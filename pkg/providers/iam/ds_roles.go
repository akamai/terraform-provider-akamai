package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) dsRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Get roles for the current account and contract",
		ReadContext: p.dsRolesRead,
		Schema: map[string]*schema.Schema{
			// inputs
			"group_id": {
				Type:        schema.TypeInt,
				Description: "A unique identifier for a group",
				Required:    true,
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
			"ignore_client_context": {
				Type:        schema.TypeBool,
				Description: "When enabled, returns all roles for the current account without regard the contract type associated with your API client",
				Optional:    true,
			},

			// outputs
			"role_id": {
				Type:        schema.TypeInt,
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
			"user_count": {
				Type:        schema.TypeInt,
				Description: "The number of users who use the role",
				Computed:    true,
			},
			"users": {
				Type:        schema.TypeSet,
				Description: "Permissions available to the user for this group",
				Computed:    true,
				Elem: schema.Resource{
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
			"actions": {
				Type:        schema.TypeSet,
				Description: "Permissions available to the user for this group",
				Computed:    true,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete": {
							Type:        schema.TypeBool,
							Description: "Indicates whether the user can remove items from the group",
							Computed:    true,
						},
						"edit": {
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

func (p *provider) dsRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p.handleMeta(m)
	return nil
}
