package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

type (
	//userDataSource defines iam user datasource
	userDataSource struct {
		meta meta.Meta
	}
	//userDataSourceModel defines structure of datasource attributes
	userDataSourceModel struct {
		UIIdentityID                       types.String        `tfsdk:"ui_identity_id"`
		AccountID                          types.String        `tfsdk:"account_id"`
		Actions                            *actionsModel       `tfsdk:"actions"`
		AdditionalAuthentication           types.String        `tfsdk:"additional_authentication"`
		AdditionalAuthenticationConfigured types.Bool          `tfsdk:"additional_authentication_configured"`
		Address                            types.String        `tfsdk:"address"`
		AuthGrants                         []*authGrantsModel  `tfsdk:"auth_grants"`
		City                               types.String        `tfsdk:"city"`
		ContactType                        types.String        `tfsdk:"contact_type"`
		Country                            types.String        `tfsdk:"country"`
		Email                              types.String        `tfsdk:"email"`
		EmailUpdatePending                 types.Bool          `tfsdk:"email_update_pending"`
		FirstName                          types.String        `tfsdk:"first_name"`
		IsLocked                           types.Bool          `tfsdk:"is_locked"`
		JobTitle                           types.String        `tfsdk:"job_title"`
		LastLoginDate                      types.String        `tfsdk:"last_login_date"`
		LastName                           types.String        `tfsdk:"last_name"`
		MobilePhone                        types.String        `tfsdk:"mobile_phone"`
		Notifications                      *notificationsModel `tfsdk:"notifications"`
		PasswordExpiryDate                 types.String        `tfsdk:"password_expiry_date"`
		Phone                              types.String        `tfsdk:"phone"`
		PreferredLanguage                  types.String        `tfsdk:"preferred_language"`
		SecondaryEmail                     types.String        `tfsdk:"secondary_email"`
		SessionTimeout                     types.Int64         `tfsdk:"session_timeout"`
		State                              types.String        `tfsdk:"state"`
		TFAConfigured                      types.Bool          `tfsdk:"tfa_configured"`
		TFAEnabled                         types.Bool          `tfsdk:"tfa_enabled"`
		TimeZone                           types.String        `tfsdk:"time_zone"`
		UIUserName                         types.String        `tfsdk:"ui_user_name"`
		ZIPCode                            types.String        `tfsdk:"zip_code"`
	}

	actionsModel struct {
		Delete           types.Bool `tfsdk:"delete"`
		APIClient        types.Bool `tfsdk:"api_client"`
		Edit             types.Bool `tfsdk:"edit"`
		IsCloneable      types.Bool `tfsdk:"is_cloneable"`
		ResetPassword    types.Bool `tfsdk:"reset_password"`
		ThirdPartyAccess types.Bool `tfsdk:"third_party_access"`
	}
	authGrantsModel struct {
		GroupID         types.Int64        `tfsdk:"group_id"`
		GroupName       types.String       `tfsdk:"group_name"`
		IsBlocked       types.Bool         `tfsdk:"is_blocked"`
		RoleDescription types.String       `tfsdk:"role_description"`
		RoleID          types.Int64        `tfsdk:"role_id"`
		RoleName        types.String       `tfsdk:"role_name"`
		SubGroups       []*authGrantsModel `tfsdk:"sub_groups"`
	}
	notificationsModel struct {
		Options                  optionsModel `tfsdk:"options"`
		EnableEmailNotifications types.Bool   `tfsdk:"enable_email_notifications"`
	}

	optionsModel struct {
		APIClientCredentialExpiryNotification types.Bool     `tfsdk:"api_client_credential_expiry_notification"`
		NewUserNotification                   types.Bool     `tfsdk:"new_user_notification"`
		PasswordExpiry                        types.Bool     `tfsdk:"password_expiry"`
		Proactive                             []types.String `tfsdk:"proactive"`
		Upgrade                               []types.String `tfsdk:"upgrade"`
	}
)

// NewUserDataSource returns a new iam allowed APIs data source
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// Metadata configures data source's meta information
func (d *userDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_user"
}

// Configure configures data source at the beginning of the lifecycle
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "IAM User Data Source",
		Attributes: map[string]schema.Attribute{
			"ui_identity_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for each user.",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for each account.",
			},
			"actions": schema.SingleNestedAttribute{
				Description: "A user's role assignments, per group.",
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
					"edit": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the user is editable.",
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
			"address": schema.StringAttribute{
				Computed:    true,
				Description: "The user's street address.",
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
						"sub_groups": nestedAuthGrant(maxSupportedGroupNesting),
					},
				},
			},
			"city": schema.StringAttribute{
				Computed:    true,
				Description: "The user's city.",
			},
			"contact_type": schema.StringAttribute{
				Computed:    true,
				Description: "User's contact type.",
			},
			"country": schema.StringAttribute{
				Computed:    true,
				Description: "User's country.",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "The user's email address.",
			},
			"email_update_pending": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether there are any pending changes to the email address.",
			},
			"first_name": schema.StringAttribute{
				Computed:    true,
				Description: "The user's first name.",
			},
			"is_locked": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the user's account is locked.",
			},
			"job_title": schema.StringAttribute{
				Computed:    true,
				Description: "The user's position at the company.",
			},
			"last_login_date": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 timestamp indicating when the user last logged in.",
			},
			"last_name": schema.StringAttribute{
				Computed:    true,
				Description: "The user's surname.",
			},
			"mobile_phone": schema.StringAttribute{
				Computed:    true,
				Description: "The user's mobile phone number, represented as a ten-digit integer within a string.",
			},
			"notifications": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Specifies email notifications the user receives for products.",
				Attributes: map[string]schema.Attribute{
					"options": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"api_client_credential_expiry_notification": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether the user gets notifications for expiring API client credentials.",
							},
							"new_user_notification": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether group administrators get notifications when the user creates other new users.",
							},
							"password_expiry": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether the user gets notifications for password expiration.",
							},
							"proactive": schema.ListAttribute{
								Computed:    true,
								Description: "Products for which the user gets notifications for service issues.",
								ElementType: types.StringType,
							},
							"upgrade": schema.ListAttribute{
								Computed:    true,
								Description: "Products for which the user receives notifications for upgrades.",
								ElementType: types.StringType,
							},
						},
					},
					"enable_email_notifications": schema.BoolAttribute{
						Computed:    true,
						Description: "Enables email notifications.",
					},
				},
			},
			"password_expiry_date": schema.StringAttribute{
				Computed:    true,
				Description: "The date a user's password expires.",
			},
			"phone": schema.StringAttribute{
				Computed:    true,
				Description: "The user's main phone number.",
			},
			"preferred_language": schema.StringAttribute{
				Computed:    true,
				Description: "The user's language.",
			},
			"secondary_email": schema.StringAttribute{
				Computed:    true,
				Description: "The user's alternate email address.",
			},
			"session_timeout": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of seconds it takes for the user's Control Center session to time out after no activity.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The user's state or province if user's country is USA or Canada.",
			},
			"tfa_configured": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether TFA is configured.",
			},
			"tfa_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether TFA is enabled.",
			},
			"time_zone": schema.StringAttribute{
				Computed:    true,
				Description: "The user's time zone.",
			},
			"ui_user_name": schema.StringAttribute{
				Computed:    true,
				Description: "The user's username in Control Center.",
			},
			"zip_code": schema.StringAttribute{
				Computed:    true,
				Description: "The user's postal code, represented as a string.",
			},
		},
	}
}

func nestedAuthGrant(depth int) *schema.ListNestedAttribute {
	authGrant := schema.ListNestedAttribute{
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
			},
		},
	}
	if depth > 1 {
		authGrant.NestedObject.Attributes["sub_groups"] = nestedAuthGrant(depth - 1)
	}
	return &authGrant
}

// Read is called when the provider must read data source values in order to update state
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM Allowed APIs DataSource Read")

	var data userDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := inst.Client(d.meta)

	user, err := client.GetUser(ctx, iam.GetUserRequest{
		IdentityID:    data.UIIdentityID.ValueString(),
		Actions:       true,
		AuthGrants:    true,
		Notifications: true,
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s:", ErrIAMGetUser), err.Error())
		return
	}

	if resp.Diagnostics.Append(data.setAttributes(user)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (d *userDataSourceModel) setAttributes(user *iam.User) diag.Diagnostics {
	d.UIIdentityID = types.StringValue(user.IdentityID)
	d.AccountID = types.StringValue(user.AccountID)
	d.AdditionalAuthentication = types.StringValue(string(user.AdditionalAuthentication))
	d.AdditionalAuthenticationConfigured = types.BoolValue(user.AdditionalAuthenticationConfigured)
	d.Address = types.StringValue(user.Address)
	d.City = types.StringValue(user.City)
	d.ContactType = types.StringValue(user.ContactType)
	d.Country = types.StringValue(user.Country)
	d.Email = types.StringValue(user.Email)
	d.State = types.StringValue(user.State)
	d.EmailUpdatePending = types.BoolValue(user.EmailUpdatePending)
	d.FirstName = types.StringValue(user.FirstName)
	d.IsLocked = types.BoolValue(user.IsLocked)
	d.JobTitle = types.StringValue(user.JobTitle)
	if !user.LastLoginDate.IsZero() {
		d.LastLoginDate = types.StringValue(user.LastLoginDate.Format(time.RFC3339Nano))
	}
	d.LastName = types.StringValue(user.LastName)
	d.MobilePhone = types.StringValue(user.MobilePhone)
	if !user.PasswordExpiryDate.IsZero() {
		d.PasswordExpiryDate = types.StringValue(user.PasswordExpiryDate.Format(time.RFC3339Nano))
	}
	d.Phone = types.StringValue(user.Phone)
	d.PreferredLanguage = types.StringValue(user.PreferredLanguage)
	d.SecondaryEmail = types.StringValue(user.SecondaryEmail)
	if user.SessionTimeOut != nil {
		d.SessionTimeout = types.Int64Value(int64(*user.SessionTimeOut))
	}
	d.TFAEnabled = types.BoolValue(user.TFAEnabled)
	d.TFAConfigured = types.BoolValue(user.TFAConfigured)
	d.TimeZone = types.StringValue(user.TimeZone)
	d.UIUserName = types.StringValue(user.UserName)
	d.ZIPCode = types.StringValue(user.ZipCode)
	d.Notifications = &notificationsModel{
		Options: optionsModel{
			APIClientCredentialExpiryNotification: types.BoolValue(user.Notifications.Options.APIClientCredentialExpiry),
			NewUserNotification:                   types.BoolValue(user.Notifications.Options.NewUser),
			PasswordExpiry:                        types.BoolValue(user.Notifications.Options.PasswordExpiry),
		},
		EnableEmailNotifications: types.BoolValue(user.Notifications.EnableEmail),
	}
	proactiveList := make([]types.String, 0, len(user.Notifications.Options.Proactive))
	for _, proactive := range user.Notifications.Options.Proactive {
		proactiveList = append(proactiveList, types.StringValue(proactive))
	}
	upgradeList := make([]types.String, 0, len(user.Notifications.Options.Upgrade))
	for _, upgrade := range user.Notifications.Options.Upgrade {
		upgradeList = append(upgradeList, types.StringValue(upgrade))
	}
	d.Notifications.Options.Proactive = proactiveList
	d.Notifications.Options.Upgrade = upgradeList
	d.Actions = &actionsModel{
		Delete:           types.BoolValue(user.Actions.Delete),
		APIClient:        types.BoolValue(user.Actions.APIClient),
		Edit:             types.BoolValue(user.Actions.Edit),
		IsCloneable:      types.BoolValue(user.Actions.IsCloneable),
		ResetPassword:    types.BoolValue(user.Actions.ResetPassword),
		ThirdPartyAccess: types.BoolValue(user.Actions.ThirdPartyAccess),
	}
	subGroups, diags := readAuthGrantSubGroups(user.AuthGrants, maxSupportedGroupNesting)
	if diags.HasError() {
		return diags
	}
	d.AuthGrants = subGroups
	return nil
}

func readAuthGrantSubGroups(authGrants []iam.AuthGrant, depth int) ([]*authGrantsModel, diag.Diagnostics) {
	authGrantModelList := make([]*authGrantsModel, 0, len(authGrants))
	for _, authGrant := range authGrants {
		grantModel := authGrantsModel{
			GroupID:         types.Int64Value(authGrant.GroupID),
			GroupName:       types.StringValue(authGrant.GroupName),
			IsBlocked:       types.BoolValue(authGrant.IsBlocked),
			RoleDescription: types.StringValue(authGrant.RoleDescription),
			RoleName:        types.StringValue(authGrant.RoleName),
		}
		if authGrant.RoleID != nil {
			grantModel.RoleID = types.Int64Value(int64(*authGrant.RoleID))
		}
		if depth > 1 {
			grants, diags := readAuthGrantSubGroups(authGrant.Subgroups, depth-1)
			if diags.HasError() {
				return nil, diags
			}
			grantModel.SubGroups = grants
		} else if depth == 1 && authGrant.Subgroups != nil && len(authGrant.Subgroups) > 0 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("unsupported subgroup depth",
				fmt.Sprintf("AuthGrant %d contains more subgroups and exceed total supported limit of nesting %d.", authGrant.GroupID, maxSupportedGroupNesting))}
		}
		authGrantModelList = append(authGrantModelList, &grantModel)
	}
	return authGrantModelList, nil
}
