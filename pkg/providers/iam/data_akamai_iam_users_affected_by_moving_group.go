package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &usersAffectedByMovingGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &usersAffectedByMovingGroupDataSource{}
)

type (
	usersAffectedByMovingGroupDataSource struct {
		meta meta.Meta
	}

	usersAffectedByMovingGroupModel struct {
		SourceGroupID      types.Int64                      `tfsdk:"source_group_id"`
		DestinationGroupID types.Int64                      `tfsdk:"destination_group_id"`
		UserType           types.String                     `tfsdk:"user_type"`
		Users              []userAffectedByMovingGroupModel `tfsdk:"users"`
	}

	userAffectedByMovingGroupModel struct {
		AccountID     types.String `tfsdk:"account_id"`
		Email         types.String `tfsdk:"email"`
		FirstName     types.String `tfsdk:"first_name"`
		LastLoginDate types.String `tfsdk:"last_login_date"`
		LastName      types.String `tfsdk:"last_name"`
		UIIdentityID  types.String `tfsdk:"ui_identity_id"`
		UIUsername    types.String `tfsdk:"ui_username"`
	}
)

// NewUsersAffectedByMovingGroupDataSource returns new users affected by moving group data source.
func NewUsersAffectedByMovingGroupDataSource() datasource.DataSource {
	return &usersAffectedByMovingGroupDataSource{}
}

func (a *usersAffectedByMovingGroupDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_users_affected_by_moving_group"
}

func (a *usersAffectedByMovingGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (a *usersAffectedByMovingGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"source_group_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier for a group you want to move.",
			},
			"destination_group_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier for a group you're putting the other group into.",
			},
			"user_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("lostAccess", "gainAccess", ""),
				},
				Description: "Filters the list by users who have lostAccess or the reverse gainAccess.",
			},
			"users": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of affected users.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for each account.",
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "The user's email address.",
						},
						"first_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's first name.",
						},
						"last_login_date": schema.StringAttribute{
							Computed:    true,
							Description: "ISO 8601 timestamp indicating when the user last logged in.",
						},
						"last_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's surname.",
						},
						"ui_identity_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for each user, which corresponds to their Control Center profile or client ID. Also known as a contactId in other APIs.",
						},
						"ui_username": schema.StringAttribute{
							Computed:    true,
							Description: "The user's username in Control Center.",
						},
					},
				},
			},
		},
	}
}

func (a *usersAffectedByMovingGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Users affected by moving group Datasource Read")

	var data usersAffectedByMovingGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(a.meta)
	users, err := client.ListAffectedUsers(ctx, iam.ListAffectedUsersRequest{
		SourceGroupID:      data.SourceGroupID.ValueInt64(),
		DestinationGroupID: data.DestinationGroupID.ValueInt64(),
		UserType:           data.UserType.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Reading IAM Users affected by moving group failed", err.Error())
		return
	}

	data.Users = convertUsersAffectedByMove(users)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func convertUsersAffectedByMove(users []iam.GroupUser) []userAffectedByMovingGroupModel {
	var convertedUsers []userAffectedByMovingGroupModel
	for _, u := range users {
		convertedUser := userAffectedByMovingGroupModel{
			AccountID:     types.StringValue(u.AccountID),
			Email:         types.StringValue(u.Email),
			FirstName:     types.StringValue(u.FirstName),
			LastName:      types.StringValue(u.LastName),
			UIIdentityID:  types.StringValue(u.IdentityID),
			UIUsername:    types.StringValue(u.UserName),
			LastLoginDate: types.StringValue(date.FormatRFC3339Nano(u.LastLoginDate)),
		}
		convertedUsers = append(convertedUsers, convertedUser)
	}
	return convertedUsers
}
