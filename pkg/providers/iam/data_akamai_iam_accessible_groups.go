package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &accessibleGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &accessibleGroupsDataSource{}
)

type (
	accessibleGroupsDataSource struct {
		meta meta.Meta
	}

	accessibleGroupsModel struct {
		Username         types.String           `tfsdk:"username"`
		AccessibleGroups []accessibleGroupModel `tfsdk:"accessible_groups"`
	}

	accessibleGroupModel struct {
		GroupID         types.Int64               `tfsdk:"group_id"`
		GroupName       types.String              `tfsdk:"group_name"`
		IsBlocked       types.Bool                `tfsdk:"is_blocked"`
		RoleDescription types.String              `tfsdk:"role_description"`
		RoleID          types.Int64               `tfsdk:"role_id"`
		RoleName        types.String              `tfsdk:"role_name"`
		SubGroups       []accessibleSubgroupModel `tfsdk:"sub_groups"`
	}

	accessibleSubgroupModel struct {
		GroupID       types.Int64               `tfsdk:"group_id"`
		GroupName     types.String              `tfsdk:"group_name"`
		ParentGroupID types.Int64               `tfsdk:"parent_group_id"`
		SubGroups     []accessibleSubgroupModel `tfsdk:"sub_groups"`
	}
)

// NewAccessibleGroupsDataSource returns new accessible groups data source.
func NewAccessibleGroupsDataSource() datasource.DataSource {
	return &accessibleGroupsDataSource{}
}

func (a *accessibleGroupsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_accessible_groups"
}

func (a *accessibleGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	a.meta = meta.Must(req.ProviderData)
}

func (a *accessibleGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Required:    true,
				Description: "User's username for which the accessible groups will be listed.",
			},
			"accessible_groups": schema.ListNestedAttribute{
				Description: "List of accessible groups for the user",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier for each group.",
						},
						"group_name": schema.StringAttribute{
							Computed:    true,
							Description: "Descriptive label for the group.",
						},
						"is_blocked": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether a user's access to a group is blocked.",
						},
						"role_description": schema.StringAttribute{
							Computed:    true,
							Description: "Descriptive label for the role to convey its use.",
						},
						"role_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier for each role.",
						},
						"role_name": schema.StringAttribute{
							Computed:    true,
							Description: "Descriptive label for the role.",
						},
						"sub_groups": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Children of the parent group.",
							// Schema effectively must be +1 size nested, to support correctly marshalling of last level (null of type []accessibleSubgroupModel onto `sub_groups`)
							NestedObject: a.subgroupsSchema(maxSupportedGroupNesting + 1),
						},
					},
				},
			},
		},
	}
}

func (a *accessibleGroupsDataSource) subgroupsSchema(remainingNesting int) schema.NestedAttributeObject {
	subgroupSchema := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"group_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier for each group.",
			},
			"group_name": schema.StringAttribute{
				Computed:    true,
				Description: "Descriptive label for the group.",
			},
			"parent_group_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier for the parent group.",
			},
		},
	}

	var subgroupsAttribute schema.Attribute
	if remainingNesting > 0 {
		subgroupsAttribute = schema.ListNestedAttribute{
			Computed:     true,
			Description:  fmt.Sprintf("Children of the parent group. Current maximal depth of subgroups is %d.", maxSupportedGroupNesting),
			NestedObject: a.subgroupsSchema(remainingNesting - 1),
		}
		subgroupSchema.Attributes["sub_groups"] = subgroupsAttribute
	}

	return subgroupSchema
}

func (a *accessibleGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Accessible Groups Datasource Read")

	var data accessibleGroupsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(a.meta)
	groups, err := client.ListAccessibleGroups(ctx, iam.ListAccessibleGroupsRequest{UserName: data.Username.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Reading IAM Accessible Groups failed", err.Error())
		return
	}

	newData, diags := a.convertAccessibleGroups(groups, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newData)...)
}

func (a *accessibleGroupsDataSource) convertAccessibleGroups(groups iam.ListAccessibleGroupsResponse, data accessibleGroupsModel) (*accessibleGroupsModel, diag.Diagnostics) {
	data.AccessibleGroups = []accessibleGroupModel{}
	for _, group := range groups {
		newGroup := accessibleGroupModel{
			GroupID:         types.Int64Value(group.GroupID),
			GroupName:       types.StringValue(group.GroupName),
			IsBlocked:       types.BoolValue(group.IsBlocked),
			RoleDescription: types.StringValue(group.RoleDescription),
			RoleID:          types.Int64Value(group.RoleID),
			RoleName:        types.StringValue(group.RoleName),
		}

		subgroups, diags := a.convertAccessibleSubgroups(group.SubGroups, maxSupportedGroupNesting)
		if diags.HasError() {
			return nil, diags
		}
		newGroup.SubGroups = subgroups

		data.AccessibleGroups = append(data.AccessibleGroups, newGroup)

	}
	return &data, nil
}

func (a *accessibleGroupsDataSource) convertAccessibleSubgroups(subgroups []iam.AccessibleSubGroup, remainingNesting int) ([]accessibleSubgroupModel, diag.Diagnostics) {
	var groups []accessibleSubgroupModel

	for _, subgroup := range subgroups {
		newSubgroup := accessibleSubgroupModel{
			GroupID:       types.Int64Value(subgroup.GroupID),
			GroupName:     types.StringValue(subgroup.GroupName),
			ParentGroupID: types.Int64Value(subgroup.ParentGroupID),
		}
		if remainingNesting > 1 {
			newSubgroups, diags := a.convertAccessibleSubgroups(subgroup.SubGroups, remainingNesting-1)
			if diags.HasError() {
				return nil, diags
			}
			newSubgroup.SubGroups = newSubgroups
		} else if len(subgroup.SubGroups) > 0 { //	We already processed the last level (remainingNesting = 1). If it still has subgroups, we cannot handle it
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("unsupported subgroup depth",
				fmt.Sprintf("Subgroup %d contains more subgroups and exceed total supported limit of nesting %d.", subgroup.GroupID, maxSupportedGroupNesting))}
		}
		groups = append(groups, newSubgroup)
	}

	return groups, nil
}
