package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

type (
	usersDataSource struct {
		meta meta.Meta
	}

	usersModel struct {
		GroupID types.Int64 `tfsdk:"group_id"`
		Users   []userItem  `tfsdk:"users"`
	}

	userItem struct {
		AccountID                          types.String      `tfsdk:"account_id"`
		Actions                            *userActions      `tfsdk:"actions"`
		AdditionalAuthentication           types.String      `tfsdk:"additional_authentication"`
		AdditionalAuthenticationConfigured types.Bool        `tfsdk:"additional_authentication_configured"`
		AuthGrants                         []authGrantsModel `tfsdk:"auth_grants"`
		Email                              types.String      `tfsdk:"email"`
		FirstName                          types.String      `tfsdk:"first_name"`
		LastName                           types.String      `tfsdk:"last_name"`
		IsLocked                           types.Bool        `tfsdk:"is_locked"`
		LastLoginDate                      types.String      `tfsdk:"last_login_date"`
		TFAConfigured                      types.Bool        `tfsdk:"tfa_configured"`
		TFAEnabled                         types.Bool        `tfsdk:"tfa_enabled"`
		UIIdentityID                       types.String      `tfsdk:"ui_identity_id"`
		UIUserName                         types.String      `tfsdk:"ui_user_name"`
	}

	userActions struct {
		Delete           types.Bool `tfsdk:"delete"`
		APIClient        types.Bool `tfsdk:"api_client"`
		CanEditMFA       types.Bool `tfsdk:"can_edit_mfa"`
		CanEditNone      types.Bool `tfsdk:"can_edit_none"`
		CanEditTFA       types.Bool `tfsdk:"can_edit_tfa"`
		Edit             types.Bool `tfsdk:"edit"`
		EditProfile      types.Bool `tfsdk:"edit_profile"`
		IsCloneable      types.Bool `tfsdk:"is_cloneable"`
		ResetPassword    types.Bool `tfsdk:"reset_password"`
		ThirdPartyAccess types.Bool `tfsdk:"third_party_access"`
	}
)

// NewUsersDataSource returns a new iam users data source
func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// Metadata configures data source's meta information
func (d *usersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_users"
}

// Configure configures data source at the beginning of the lifecycle
func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.meta = meta.Must(req.ProviderData)
}

// Schema is used to define data source's terraform schema
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management users",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Filters users for a specific group.",
			},
			"users": schema.ListNestedAttribute{
				Description: "List of users in the account",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier of the account.",
						},
						"actions": schema.SingleNestedAttribute{
							Description: "Specifies permissions available to the user for this group.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"delete": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether the user is deletable.",
								},
								"api_client": schema.BoolAttribute{
									Computed:    true,
									Description: "Enables the administrator to create an API client.",
								},
								"can_edit_mfa": schema.BoolAttribute{
									Computed:    true,
									Description: "With a true value, the user can turn their MFA setting on or off.",
								},
								"can_edit_none": schema.BoolAttribute{
									Computed:    true,
									Description: "With a true value, the user can turn their None setting on or off.",
								},
								"can_edit_tfa": schema.BoolAttribute{
									Computed:    true,
									Description: "With a true value, the user can turn their TFA setting on or off.",
								},
								"edit": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether the user is editable.",
								},
								"edit_profile": schema.BoolAttribute{
									Computed:    true,
									Description: "With a true value, the user can edit their user profile.",
								},
								"is_cloneable": schema.BoolAttribute{
									Computed:    true,
									Description: "Enables an administrator to create a new user with permissions cloned from this user.",
								},
								"reset_password": schema.BoolAttribute{
									Computed:    true,
									Description: "Enables an administrator to send a user a password by email or see a one-time token.",
								},
								"third_party_access": schema.BoolAttribute{
									Computed:    true,
									Description: "Enables the administrator to manage extended access.",
								},
							},
						},
						"additional_authentication": schema.StringAttribute{
							Computed:    true,
							Description: "Specifies the user's multi-factor authentication method, confirming their identity.",
						},
						"additional_authentication_configured": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the user has multi-factor authentication configured.",
						},
						"auth_grants": schema.ListNestedAttribute{
							Description: "A user's role assignments, per group.",
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
									"sub_groups": nestedAuthGrant(50),
								}},
						},
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "The user's email address.",
						},
						"first_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's first name.",
						},
						"is_locked": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the user's account is locked.",
						},
						"last_login_date": schema.StringAttribute{
							Computed:    true,
							Description: "ISO 8601 timestamp indicating when the user last logged in.",
						},
						"last_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's surname.",
						},
						"tfa_configured": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether TFA is configured.",
						},
						"tfa_enabled": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether TFA is enabled.",
						},
						"ui_identity_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for each user, which corresponds to their Control Center profile or client ID. Also known as a contactId in other APIs.",
						},
						"ui_user_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's username in Control Center.",
						},
					},
				},
			}}}
}

// Read is called when the provider must read data source values in order to update state
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Users DataSource Read")

	var data usersModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := inst.Client(d.meta)

	groupID := data.GroupID.ValueInt64Pointer()

	users, err := client.ListUsers(ctx, iam.ListUsersRequest{
		GroupID:    groupID,
		AuthGrants: true,
		Actions:    true,
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s:", ErrIAMListUsers), err.Error())
		return
	}

	if resp.Diagnostics.Append(data.read(users)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (d *usersModel) read(users []iam.UserListItem) diag.Diagnostics {
	for _, user := range users {
		authGrants, diags := readAuthGrantSubGroups(user.AuthGrants, maxSupportedGroupNesting)
		if diags.HasError() {
			return diags
		}
		userItem := userItem{
			AccountID:                          types.StringValue(user.AccountID),
			AdditionalAuthentication:           types.StringValue(string(user.AdditionalAuthentication)),
			AdditionalAuthenticationConfigured: types.BoolValue(user.AdditionalAuthenticationConfigured),
			AuthGrants:                         authGrants,
			Email:                              types.StringValue(user.Email),
			FirstName:                          types.StringValue(user.FirstName),
			IsLocked:                           types.BoolValue(user.IsLocked),
			LastName:                           types.StringValue(user.LastName),
			LastLoginDate:                      types.StringValue(date.FormatRFC3339Nano(user.LastLoginDate)),
			TFAConfigured:                      types.BoolValue(user.TFAConfigured),
			TFAEnabled:                         types.BoolValue(user.TFAEnabled),
			UIIdentityID:                       types.StringValue(user.IdentityID),
			UIUserName:                         types.StringValue(user.UserName),
		}

		if user.Actions != nil {
			userItem.Actions = &userActions{
				Delete:           types.BoolValue(user.Actions.Delete),
				APIClient:        types.BoolValue(user.Actions.APIClient),
				CanEditMFA:       types.BoolValue(user.Actions.CanEditMFA),
				CanEditNone:      types.BoolValue(user.Actions.CanEditNone),
				CanEditTFA:       types.BoolValue(user.Actions.CanEditTFA),
				Edit:             types.BoolValue(user.Actions.Edit),
				EditProfile:      types.BoolValue(user.Actions.EditProfile),
				IsCloneable:      types.BoolValue(user.Actions.IsCloneable),
				ResetPassword:    types.BoolValue(user.Actions.ResetPassword),
				ThirdPartyAccess: types.BoolValue(user.Actions.ThirdPartyAccess),
			}
		}
		d.Users = append(d.Users, userItem)
	}
	return nil
}
