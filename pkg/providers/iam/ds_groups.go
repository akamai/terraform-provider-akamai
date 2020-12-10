package iam

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource schema for akamai_iam_groups data source
func (p *provider) dsGroups() *schema.Resource {
	return &schema.Resource{
		Description: `List all groups in which you have a scope of "admin" for the current account and contract`,
		ReadContext: p.dsGroupsRead,
		Schema: map[string]*schema.Schema{
			"get_actions": {
				Type:        schema.TypeBool,
				Description: `When enabled, the response includes information about actions such as "edit" or "delete"`,
				Optional:    true,
			},
			"groups": NestedGroupsSchema(25), // Can handle group nesting up to 25 levels deep
		},
	}
}

// NestedGroupsSchema builds a nested groups schema to the given depth
func NestedGroupsSchema(depth int) *schema.Schema {
	schem := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The group's name",
			Computed:    true,
		},
		"group_id": {
			Type:        schema.TypeInt,
			Description: "A unique identifier for each group",
			Computed:    true,
		},
		"parent_group_id": {
			Type:        schema.TypeInt,
			Description: "Identifies the parent group to which a group belongs",
			Computed:    true,
		},
		"time_created": {
			Type:        schema.TypeString,
			Description: "ISO 8601 timestamp indicating when the group was originally created",
			Computed:    true,
		},
		"time_modified": {
			Type:        schema.TypeString,
			Description: "ISO 8601 timestamp indicating when the group was last updated",
			Computed:    true,
		},
		"modified_by": {
			Type:        schema.TypeString,
			Description: "The username or email of the last person to edit the group",
			Computed:    true,
		},
		"created_by": {
			Type:        schema.TypeString,
			Description: "The user name or email of the person who created the group",
			Computed:    true,
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
	}

	if depth > 1 {
		schem["sub_groups"] = NestedGroupsSchema(depth - 1)
	}

	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem:     &schema.Resource{Schema: schem},
	}
}

func (p *provider) dsGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
