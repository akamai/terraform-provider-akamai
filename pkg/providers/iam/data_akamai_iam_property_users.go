package iam

import (
	"context"
	"fmt"
	"regexp"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &propertyUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &propertyUsersDataSource{}
)

type (
	propertyUsersDataSource struct {
		meta meta.Meta
	}

	propertyUsersModel struct {
		AssetID  types.String   `tfsdk:"asset_id"`
		UserType types.String   `tfsdk:"user_type"`
		Users    []propertyUser `tfsdk:"users"`
	}

	propertyUser struct {
		FirstName    types.String `tfsdk:"first_name"`
		LastName     types.String `tfsdk:"last_name"`
		IsBlocked    types.Bool   `tfsdk:"is_blocked"`
		UIIdentityID types.String `tfsdk:"ui_identity_id"`
		UIUserName   types.String `tfsdk:"ui_user_name"`
	}
)

// NewPropertyUsersDataSource returns a new iam property users data source.
func NewPropertyUsersDataSource() datasource.DataSource {
	return &propertyUsersDataSource{}
}

func (d *propertyUsersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_property_users"
}

func (d *propertyUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *propertyUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management property users data source. Lists users for " +
			"property or include.",
		Attributes: map[string]schema.Attribute{
			"asset_id": schema.StringAttribute{
				Required: true,
				Description: "IAM identifier of the property or include. The optional aid_ prefix " +
					"is allowed.",
				Validators: []validator.String{stringvalidator.RegexMatches(
					regexp.MustCompile(`^(aid_)?\d+$`),
					`must be a number with the optional "aid_" prefix`)},
			},
			"user_type": schema.StringAttribute{
				Optional:    true,
				Description: "Filters the list based on users' access to the property.",
				Validators: []validator.String{stringvalidator.OneOf(
					string(iam.PropertyUserTypeAll),
					string(iam.PropertyUserTypeBlocked),
					string(iam.PropertyUserTypeAssigned))},
			},
			"users": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of users accessing the property.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"first_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's first name.",
						},
						"last_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's surname.",
						},
						"is_blocked": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether a user's access to a property is blocked.",
						},
						"ui_identity_id": schema.StringAttribute{
							Computed: true,
							Description: "Unique identifier for each user, which corresponds " +
								"to their Control Center profile or client ID. Also known " +
								"as a contactId in other APIs.",
						},
						"ui_user_name": schema.StringAttribute{
							Computed:    true,
							Description: "The user's username in Control Center.",
						},
					},
				},
			},
		},
	}
}

func (d *propertyUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Property Users DataSource Read")

	var data propertyUsersModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := str.GetIntID(data.AssetID.ValueString(), "aid_")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read IAM Property Users DataSource",
			fmt.Sprintf("Error occurred while parsing property id: %s. Error: %s",
				data.AssetID.ValueString(), err.Error()))
		return
	}

	client := inst.Client(d.meta)

	users, err := client.ListUsersForProperty(ctx, iam.ListUsersForPropertyRequest{
		PropertyID: int64(id),
		UserType:   iam.PropertyUserType(data.UserType.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read IAM Property Users DataSource",
			fmt.Sprintf(
				"Error occurred while listing users for property: %s (user type: %s). Error: %s",
				data.AssetID.ValueString(), data.UserType.ValueString(), err.Error()))
		return
	}

	for _, user := range users {
		data.Users = append(data.Users, newPropertyUser(user))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func newPropertyUser(user iam.UsersForProperty) propertyUser {
	return propertyUser{
		FirstName:    types.StringValue(user.FirstName),
		LastName:     types.StringValue(user.LastName),
		IsBlocked:    types.BoolValue(user.IsBlocked),
		UIIdentityID: types.StringValue(user.UIIdentityID),
		UIUserName:   types.StringValue(user.UIUserName),
	}
}
