package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

	authorizedUsersDataSourceModel struct {
		AuthorizedUsers []userModel `tfsdk:"authorized_users"`
	}

	userModel struct {
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

// Metadata configures data source's meta information
func (a *authorizedUsersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_iam_authorized_users"
}

// Configure configures data source at the beginning of the lifecycle
func (a *authorizedUsersDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			response.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					request.ProviderData))
		}
	}()
	a.meta = meta.Must(request.ProviderData)
}

func (a *authorizedUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
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

func (a *authorizedUsersDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Authorized Users DataSource Read")

	var data authorizedUsersDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &data)...); response.Diagnostics.HasError() {
		return
	}
	client := inst.Client(a.meta)

	users, err := client.ListAuthorizedUsers(ctx)
	if err != nil {
		response.Diagnostics.AddError("Failed to fetch authorized users:", err.Error())
		return
	}

	if response.Diagnostics.Append(data.read(users)...); response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (d *authorizedUsersDataSourceModel) read(authorisedUsers iam.ListAuthorizedUsersResponse) diag.Diagnostics {
	var users []userModel
	for _, user := range authorisedUsers {
		authorizedUser := userModel{
			Email:        types.StringValue(user.Email),
			Username:     types.StringValue(user.Username),
			FirstName:    types.StringValue(user.FirstName),
			LastName:     types.StringValue(user.LastName),
			UIIdentityID: types.StringValue(user.UIIdentityID),
		}

		users = append(users, authorizedUser)
	}

	d.AuthorizedUsers = users
	return nil
}
