package iam

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIAMGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage a group in your account",
		CreateContext: resourceIAMGroupCreate,
		ReadContext:   resourceIAMGroupRead,
		UpdateContext: resourceIAMGroupUpdate,
		DeleteContext: resourceIAMGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"parent_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier for the parent group",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human readable name for a group",
			},
			"sub_groups": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: "Subgroups IDs",
			},
		},
	}
}

func resourceIAMGroupCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMGroupCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Creating group")

	parentGroupID, err := tools.GetIntValue("parent_group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupName, err := tools.GetStringValue("name", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := client.CreateGroup(ctx, iam.GroupRequest{GroupID: int64(parentGroupID), GroupName: groupName})
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(strconv.FormatInt(group.GroupID, 10))

	return resourceIAMGroupRead(ctx, rd, m)
}

func resourceIAMGroupRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMGroupRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	groupID, err := strconv.ParseInt(rd.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := client.GetGroup(ctx, iam.GetGroupRequest{GroupID: groupID})
	if err != nil {
		return diag.FromErr(err)
	}

	subGroups := make([]int64, len(group.SubGroups))
	for i, g := range group.SubGroups {
		subGroups[i] = g.GroupID
	}

	data := map[string]interface{}{
		"parent_group_id": group.ParentGroupID,
		"name":            group.GroupName,
		"sub_groups":      subGroups,
	}
	if err = tools.SetAttrs(rd, data); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIAMGroupUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMGroupUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	groupID, err := strconv.ParseInt(rd.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	parentGroupID, err := tools.GetIntValue("parent_group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	if rd.HasChange("name") {
		groupName, err := tools.GetStringValue("name", rd)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.UpdateGroupName(ctx, iam.GroupRequest{GroupName: groupName, GroupID: groupID})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if rd.HasChange("parent_group_id") {
		err := client.MoveGroup(ctx, iam.MoveGroupRequest{DestinationGroupID: int64(parentGroupID), SourceGroupID: groupID})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIAMGroupRead(ctx, rd, m)
}

func resourceIAMGroupDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("IAM", "resourceIAMGroupDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	groupID, err := strconv.ParseInt(rd.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.RemoveGroup(ctx, iam.RemoveGroupRequest{GroupID: groupID}); err != nil {
		return diag.FromErr(err)
	}

	rd.SetId("")

	return nil
}
