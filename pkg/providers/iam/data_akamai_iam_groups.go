package iam

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMGroups() *schema.Resource {
	return &schema.Resource{
		Description: `List all groups in which you have a scope of "admin" for the current account and contract.`,
		ReadContext: dataIAMGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": nestedGroupsSchema(maxSupportedGroupNesting), // Can handle groups with nesting up to 50 levels deep
		},
	}
}

// nestedGroupsSchema builds a nested groups schema to the given depth
func nestedGroupsSchema(depth int) *schema.Schema {
	nestedSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The group's name.",
			Computed:    true,
		},
		"group_id": {
			Type:        schema.TypeString,
			Description: "A unique identifier for each group.",
			Computed:    true,
		},
		"parent_group_id": {
			Type:        schema.TypeString,
			Description: "Identifies the parent group to which a group belongs.",
			Computed:    true,
		},
		"time_created": {
			Type:        schema.TypeString,
			Description: "ISO 8601 timestamp indicating when the group was originally created.",
			Computed:    true,
		},
		"time_modified": {
			Type:        schema.TypeString,
			Description: "ISO 8601 timestamp indicating when the group was last updated.",
			Computed:    true,
		},
		"modified_by": {
			Type:        schema.TypeString,
			Description: "The username or email of the last person to edit the group.",
			Computed:    true,
		},
		"created_by": {
			Type:        schema.TypeString,
			Description: "The user name or email of the person who created the group.",
			Computed:    true,
		},
	}

	if depth > 1 {
		nestedSchema["sub_groups"] = nestedGroupsSchema(depth - 1)
	}

	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem:     &schema.Resource{Schema: nestedSchema},
	}
}

func dataIAMGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "dataIAMGroupsRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching groups")
	res, err := client.ListGroups(ctx, iam.ListGroupsRequest{})
	if err != nil {
		logger.Error("Could not get groups", "error", err)
		return diag.FromErr(err)
	}

	groups := groupsToState(res)

	if err := d.Set("groups", groups); err != nil {
		logger.Error("Could not set groups in state", "error", err)
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
	m["time_created"] = date.FormatRFC3339Nano(g.CreatedDate)
	m["time_modified"] = date.FormatRFC3339Nano(g.ModifiedDate)
	m["modified_by"] = g.ModifiedBy
	m["created_by"] = g.CreatedBy

	if len(g.SubGroups) > 0 {
		m["sub_groups"] = groupsToState(g.SubGroups)
	}

	return m
}
