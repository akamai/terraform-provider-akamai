package iam

import (
	"context"

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

	groups := []interface{}{}
	for _, ct := range res {
		// groups = append(groups, ct)
	}

	if err := d.Set("groups", groups); err != nil {
		logger.WithError(err).Error("Could not set groups in state")
	}

	d.SetId("akamai_iam_groups")
	return nil
}
