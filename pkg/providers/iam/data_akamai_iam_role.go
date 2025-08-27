package iam

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

// NewRoleDataSource returns all the details for a role.
func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

type (
	roleDataSource struct {
		meta meta.Meta
	}

	roleModel struct {
		RoleID          types.Int64   `tfsdk:"role_id"`
		Actions         *actions      `tfsdk:"actions"`
		CreatedBy       types.String  `tfsdk:"created_by"`
		CreatedDate     types.String  `tfsdk:"created_date"`
		ModifiedBy      types.String  `tfsdk:"modified_by"`
		ModifiedDate    types.String  `tfsdk:"modified_date"`
		RoleDescription types.String  `tfsdk:"role_description"`
		RoleName        types.String  `tfsdk:"role_name"`
		Type            types.String  `tfsdk:"type"`
		GrantedRoles    []grantedRole `tfsdk:"granted_roles"`
		Users           []user        `tfsdk:"users"`
	}

	grantedRole struct {
		GrantedRoleID          types.Int64  `tfsdk:"granted_role_id"`
		GrantedRoleName        types.String `tfsdk:"granted_role_name"`
		GrantedRoleDescription types.String `tfsdk:"granted_role_description"`
	}

	user struct {
		AccountID     types.String `tfsdk:"account_id"`
		Email         types.String `tfsdk:"email"`
		FirstName     types.String `tfsdk:"first_name"`
		LastName      types.String `tfsdk:"last_name"`
		LastLoginDate types.String `tfsdk:"last_login_date"`
		UIIdentityID  types.String `tfsdk:"ui_identity_id"`
	}
)

func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Role data source.",
		Attributes: map[string]schema.Attribute{
			"role_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier for each role.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the granted role.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the granted role was originally created.",
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who last edited the granted role.",
			},
			"modified_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the granted role was last updated.",
			},
			"role_description": schema.StringAttribute{
				Computed:    true,
				Description: "Descriptive label for the role to convey its use.",
			},
			"role_name": schema.StringAttribute{
				Computed:    true,
				Description: "Descriptive label for the role.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Whether it's a standard role defined by Akamai or a custom role created by a user of your account. You can't modify or delete a standard role.",
				Validators: []validator.String{
					stringvalidator.OneOf("standard", "custom"),
				},
			},
			"actions": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Specifies activities available for the role.",
				Attributes: map[string]schema.Attribute{
					"delete": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can remove the role.",
					},
					"edit": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can modify the role.",
					},
				},
			},
			"granted_roles": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Granted roles for the account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"granted_role_description": schema.StringAttribute{
							Computed:    true,
							Description: "Descriptive label for the role to convey its use.",
						},
						"granted_role_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier for each granted role.",
						},
						"granted_role_name": schema.StringAttribute{
							Computed:    true,
							Description: "Descriptive label for the granted role.",
						},
					},
				},
			},
			"users": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Users on the account who share the same role.",
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
					},
				},
			},
		},
	}
}

func (d *roleDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_role"
}

func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Role DataSource Read")

	var data roleModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)

	getRoleResp, err := client.GetRole(ctx, iam.GetRoleRequest{
		ID:           data.RoleID.ValueInt64(),
		Actions:      true,
		GrantedRoles: true,
		Users:        true,
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching iam role failed", err.Error())
		return
	}

	data.setAttributes(getRoleResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *roleModel) setAttributes(role *iam.Role) {

	m.CreatedBy = types.StringValue(role.CreatedBy)
	m.CreatedDate = types.StringValue(date.FormatRFC3339Nano(role.CreatedDate))
	m.ModifiedBy = types.StringValue(role.ModifiedBy)
	m.ModifiedDate = types.StringValue(date.FormatRFC3339Nano(role.ModifiedDate))
	m.RoleName = types.StringValue(role.RoleName)
	m.RoleDescription = types.StringValue(role.RoleDescription)
	m.Type = types.StringValue(string(role.RoleType))

	if role.Actions != nil {
		m.Actions = &actions{
			Delete: types.BoolValue(role.Actions.Delete),
			Edit:   types.BoolValue(role.Actions.Edit),
		}
	}

	for _, r := range role.GrantedRoles {
		grantedRoleObject := grantedRole{
			GrantedRoleID:          types.Int64Value(r.RoleID),
			GrantedRoleName:        types.StringValue(r.RoleName),
			GrantedRoleDescription: types.StringValue(r.Description),
		}
		m.GrantedRoles = append(m.GrantedRoles, grantedRoleObject)
	}
	for _, u := range role.Users {
		userObject := user{
			AccountID:     types.StringValue(u.AccountID),
			Email:         types.StringValue(u.Email),
			FirstName:     types.StringValue(u.FirstName),
			LastName:      types.StringValue(u.LastName),
			LastLoginDate: types.StringValue(u.LastLoginDate.String()),
			UIIdentityID:  types.StringValue(u.UIIdentityID),
		}

		m.Users = append(m.Users, userObject)
	}
}
