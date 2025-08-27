package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

const maxSupportedGroupNesting = 50

// NewGroupDataSource returns all the details for a group.
func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type (
	groupDataSource struct {
		meta meta.Meta
	}

	groupModel struct {
		GroupID       types.Int64  `tfsdk:"group_id"`
		GroupName     types.String `tfsdk:"group_name"`
		Actions       *actions     `tfsdk:"actions"`
		CreatedBy     types.String `tfsdk:"created_by"`
		CreatedDate   types.String `tfsdk:"created_date"`
		ModifiedBy    types.String `tfsdk:"modified_by"`
		ModifiedDate  types.String `tfsdk:"modified_date"`
		ParentGroupID types.Int64  `tfsdk:"parent_group_id"`
		SubGroups     []groupModel `tfsdk:"sub_groups"`
	}
)

// groupSchemaAttributes returns the schema attributes for a group, with optional nesting for subgroups
func (d *groupDataSource) groupSchemaAttributes(remainingNesting int) map[string]schema.Attribute {
	groupAttributes := map[string]schema.Attribute{
		"group_id": schema.Int64Attribute{
			Required:    true,
			Description: "Unique identifier for each group.",
		},
		"group_name": schema.StringAttribute{
			Computed:    true,
			Description: "Descriptive label for the group.",
		},
		"created_by": schema.StringAttribute{
			Computed:    true,
			Description: "The user who created the group.",
		},
		"created_date": schema.StringAttribute{
			Computed:    true,
			Description: "ISO 8601 timestamp indicating when the group was created.",
		},
		"modified_by": schema.StringAttribute{
			Computed:    true,
			Description: "The user who last edited the group.",
		},
		"modified_date": schema.StringAttribute{
			Computed:    true,
			Description: "ISO 8601 timestamp indicating when the group was last updated.",
		},
		"parent_group_id": schema.Int64Attribute{
			Computed:    true,
			Description: "Unique identifier for the parent group within the subgroup tree.",
		},
		"actions": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Specifies activities available for the group.",
			Attributes: map[string]schema.Attribute{
				"delete": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether you can remove the group from the account. You can't remove a group that contains resources or subgroups, or if users have roles on that group.",
				},
				"edit": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether you can modify the group.",
				},
			},
		},
	}

	if remainingNesting > 0 {
		groupAttributes["sub_groups"] = schema.ListNestedAttribute{
			Computed:    true,
			Description: fmt.Sprintf("Children of the parent group. Maximal depth of subgroups is %d.", maxSupportedGroupNesting),
			NestedObject: schema.NestedAttributeObject{
				Attributes: d.groupSchemaAttributes(remainingNesting - 1),
			},
		}
	}

	return groupAttributes
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "IAM Group data source.",
		Attributes:          d.groupSchemaAttributes(maxSupportedGroupNesting + 1),
	}
}

func (d *groupDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_group"
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	d.meta = meta.Must(req.ProviderData)
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Group DataSource Read")

	var data groupModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)

	getGroupResp, err := client.GetGroup(ctx, iam.GetGroupRequest{
		GroupID: data.GroupID.ValueInt64(),
		Actions: true,
	})
	if err != nil {
		resp.Diagnostics.AddError("Fetching IAM group failed", err.Error())
		return
	}

	groupData, diags := d.convertGroupData(getGroupResp, data, maxSupportedGroupNesting)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &groupData)...)
}

// convertGroupData converts the response from the IAM client into the groupDataSourceModel structure
func (d *groupDataSource) convertGroupData(group *iam.Group, data groupModel, remainingNesting int) (groupModel, diag.Diagnostics) {
	data.GroupName = types.StringValue(group.GroupName)
	data.CreatedBy = types.StringValue(group.CreatedBy)
	data.CreatedDate = types.StringValue(date.FormatRFC3339Nano(group.CreatedDate))
	data.ModifiedBy = types.StringValue(group.ModifiedBy)
	data.ModifiedDate = types.StringValue(date.FormatRFC3339Nano(group.ModifiedDate))
	data.ParentGroupID = types.Int64Value(group.ParentGroupID)
	data.GroupID = types.Int64Value(group.GroupID)

	if group.Actions != nil {
		data.Actions = &actions{
			Delete: types.BoolValue(group.Actions.Delete),
			Edit:   types.BoolValue(group.Actions.Edit),
		}
	}

	var subGroups []groupModel
	if remainingNesting > 1 {
		for _, subGroup := range group.SubGroups {
			subGroupData, diags := d.convertGroupData(&subGroup, groupModel{}, remainingNesting-1)
			if diags.HasError() {
				return groupModel{}, diags
			}
			subGroups = append(subGroups, subGroupData)
		}
		data.SubGroups = subGroups
	} else if remainingNesting <= 1 && len(group.SubGroups) > 0 {
		return groupModel{}, diag.Diagnostics{diag.NewErrorDiagnostic(
			"unsupported subgroup depth",
			fmt.Sprintf("Subgroup %d contains more subgroups and exceeds the total supported limit of nesting %d.", group.GroupID, maxSupportedGroupNesting),
		)}
	}

	return data, nil
}
