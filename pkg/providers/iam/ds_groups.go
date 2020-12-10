package iam

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource schema for akamai_iam_groups data source
func (p *provider) dsGroups() *schema.Resource {
	return &schema.Resource{
		Description: `List all groups in which you have a scope of "admin" for the current account and contract`,
		ReadContext: p.tfCRUD("ds:Groups:Read", p.dsGroupsRead),
		Schema: map[string]*schema.Schema{
			"get_actions": {
				Type:        schema.TypeBool,
				Description: `When enabled, the response includes information about actions such as "edit" or "delete"`,
				Optional:    true,
			},
			"groups": NestedGroupsSchema(50), // Can handle groups with nesting up to 50 levels deep
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
			Type:        schema.TypeString,
			Description: "A unique identifier for each group",
			Computed:    true,
		},
		"parent_group_id": {
			Type:        schema.TypeString,
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

func (p *provider) dsGroupsRead(ctx context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	logger := log.FromContext(ctx)

	Actions := d.Get("get_actions").(bool)

	logger.Debug("Fetching groups")
	req := iam.ListGroupsRequest{Actions: Actions}
	res, err := p.client.ListGroups(ctx, req)
	if err != nil {
		logger.WithError(err).Error("Could not get groups")
		return diag.FromErr(err)
	}

	groups := groupsToState(res)

	if err := d.Set("groups", groups); err != nil {
		logger.WithError(err).Error("Could not set groups in state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_groups")
	return nil
}

// Convert many groups to a value that can be stored in state
func groupsToState(groups []iam.Group) []interface{} {
	var out []interface{}

	for _, g := range groups {
		out = append(out, groupToState(g))
	}

	return out
}

// Convert a group to a value that can be stored in state
func groupToState(g iam.Group) map[string]interface{} {
	m := map[string]interface{}{}

	m["name"] = g.GroupName
	m["group_id"] = strconv.FormatInt(g.GroupID, 10)
	m["parent_group_id"] = strconv.FormatInt(g.ParentGroupID, 10)
	m["time_created"] = g.CreatedDate
	m["time_modified"] = g.ModifiedDate
	m["modified_by"] = g.ModifiedBy
	m["created_by"] = g.CreatedBy

	if g.Actions != nil {
		m["edit_allowed"] = g.Actions.Edit
		m["delete_allowed"] = g.Actions.Delete
	}

	if len(g.SubGroups) > 0 {
		m["sub_groups"] = groupsToState(g.SubGroups)
	}

	return m
}
