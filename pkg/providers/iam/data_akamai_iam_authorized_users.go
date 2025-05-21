package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &authorizedUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &authorizedUsersDataSource{}
)

type (
	authorizedUsersDataSource struct {
		meta meta.Meta
	}

	authorizedUsersModel struct {
		AuthorizedUsers []authorizedUserModel `tfsdk:"authorized_users"`
	}

	authorizedUserModel struct {
		Email        types.String `tfsdk:"email"`
		Username     types.String `tfsdk:"username"`
		FirstName    types.String `tfsdk:"first_name"`
		LastName     types.String `tfsdk:"last_name"`
		UIIdentityID types.String `tfsdk:"ui_identity_id"`
	}
)

// NewAuthorizedUsersDataSource returns the list of authorized API client users.
func NewAuthorizedUsersDataSource() datasource.DataSource {
	return &authorizedUsersDataSource{}
}

func (a *authorizedUsersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_authorized_users"
}

func (a *authorizedUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (a *authorizedUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of authorized API client users",
		Attributes: map[string]schema.Attribute{
			"authorized_users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Computed:    true,
							Description: "The user's email address.",
						},
						"first_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's first name.",
						},
						"last_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's surname.",
						},
						"ui_identity_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for each user, which corresponds to their Control Center profile or client ID. Also known as a contactId in other APIs.",
						},
						"username": schema.StringAttribute{
							Computed:    true,
							Description: "The user's Control Center sign-in name.",
						},
					},
				},
			},
		},
	}
}

func (a *authorizedUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Authorized Users DataSource Read")

	var data authorizedUsersModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := inst.Client(a.meta)

	users, err := client.ListAuthorizedUsers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch authorized users:", err.Error())
		return
	}

	data.read(users)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *authorizedUsersModel) read(authorisedUsers iam.ListAuthorizedUsersResponse) {
	var users []authorizedUserModel
	for _, u := range authorisedUsers {
		authorizedUser := authorizedUserModel{
			Email:        types.StringValue(u.Email),
			Username:     types.StringValue(u.Username),
			FirstName:    types.StringValue(u.FirstName),
			LastName:     types.StringValue(u.LastName),
			UIIdentityID: types.StringValue(u.UIIdentityID),
		}

		users = append(users, authorizedUser)
	}

	d.AuthorizedUsers = users
}
