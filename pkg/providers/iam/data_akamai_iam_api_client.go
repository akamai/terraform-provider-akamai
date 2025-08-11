package iam

import (
	"context"
	"fmt"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &apiClientDataSource{}
	_ datasource.DataSourceWithConfigure = &apiClientDataSource{}
)

type apiClientDataSource struct {
	meta meta.Meta
}

type (
	clientModel struct {
		ClientID                types.String           `tfsdk:"client_id"`
		AccessToken             types.String           `tfsdk:"access_token"`
		Actions                 *apiClientActionsModel `tfsdk:"actions"`
		ActiveCredentialCount   types.Int64            `tfsdk:"active_credential_count"`
		AllowAccountSwitch      types.Bool             `tfsdk:"allow_account_switch"`
		APIAccess               *apiAccessModel        `tfsdk:"api_access"`
		AuthorizedUsers         types.List             `tfsdk:"authorized_users"`
		BaseURL                 types.String           `tfsdk:"base_url"`
		CanAutoCreateCredential types.Bool             `tfsdk:"can_auto_create_credential"`
		ClientDescription       types.String           `tfsdk:"client_description"`
		ClientName              types.String           `tfsdk:"client_name"`
		ClientType              types.String           `tfsdk:"client_type"`
		CreatedBy               types.String           `tfsdk:"created_by"`
		CreatedDate             types.String           `tfsdk:"created_date"`
		Credentials             []getCredentialsModel  `tfsdk:"credentials"`
		GroupAccess             *groupAccessModel      `tfsdk:"group_access"`
		IPACL                   *ipACLModel            `tfsdk:"ip_acl"`
		NotificationEmails      types.List             `tfsdk:"notification_emails"`
		PurgeOptions            *purgeOptionsModel     `tfsdk:"purge_options"`
		IsLocked                types.Bool             `tfsdk:"is_locked"`
	}

	getCredentialsModel struct {
		Actions      types.Object `tfsdk:"actions"`
		ClientToken  types.String `tfsdk:"client_token"`
		CreatedOn    types.String `tfsdk:"created_on"`
		CredentialID types.Int64  `tfsdk:"credential_id"`
		Description  types.String `tfsdk:"description"`
		ExpiresOn    types.String `tfsdk:"expires_on"`
		Status       types.String `tfsdk:"status"`
	}
)

// NewAPIClientDataSource returns a new iam API client data source.
func NewAPIClientDataSource() datasource.DataSource {
	return &apiClientDataSource{}
}

func (d *apiClientDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_iam_api_client"
}

func (d *apiClientDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *apiClientDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Identity and Access Management API clients.",
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "A unique identifier for the API client.",
			},
			"access_token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The part of the client secret that identifies your API client and lets you access applications and resources.",
			},
			"actions": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Specifies activities available for the API client.",
				Attributes: map[string]schema.Attribute{
					"delete": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can remove the API client.",
					},
					"deactivate_all": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can deactivate the API client's credentials.",
					},
					"edit": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can update the API client.",
					},
					"edit_apis": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can update the `apis` the API client can access, same as `edit_auth`.",
					},
					"edit_auth": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can update the `apis` the API client can access, same as `edit_apis`.",
					},
					"edit_groups": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can update the `groups` the API client can access.",
					},
					"edit_ip_acl": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can update the `ip_acl` the API client can access.",
					},
					"edit_switch_account": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can update the API client's option to manage many accounts.",
					},
					"lock": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can lock the API client.",
					},
					"transfer": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can transfer the API client to a new owner.",
					},
					"unlock": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can unlock the API client.",
					},
				},
			},
			"active_credential_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of credentials active for the API client.",
			},
			"allow_account_switch": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the API client can manage more than one account.",
			},
			"api_access": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The APIs the API client can access.",
				Attributes: map[string]schema.Attribute{
					"all_accessible_apis": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the API client has access to a full set of available APIs.",
					},
					"apis": schema.SetNestedAttribute{
						Computed:    true,
						Description: "The set of APIs the API client can access when `all_accessible_apis` is disabled.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"access_level": schema.StringAttribute{
									Computed:    true,
									Description: "The API client's access level on an API basis, either `READ-ONLY`, `READ-WRITE`, `CREDENTIAL-READ-ONLY`, or `CREDENTIAL-READ-WRITE`.",
								},
								"api_id": schema.Int64Attribute{
									Computed:    true,
									Description: "A unique identifier for each API.",
								},
								"api_name": schema.StringAttribute{
									Computed:    true,
									Description: "A human-readable name for the API.",
								},
								"description": schema.StringAttribute{
									Computed:    true,
									Description: "A human-readable description for the API.",
								},
								"documentation_url": schema.StringAttribute{
									Computed:    true,
									Description: "A link to more information about the API.",
								},
								"endpoint": schema.StringAttribute{
									Computed:    true,
									Description: "Specifies where the API can access resources.",
								},
							},
						},
					},
				},
			},
			"authorized_users": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "The API client's valid users.",
			},
			"base_url": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The base URL for the service.",
			},
			"can_auto_create_credential": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the API client can create a credential for a new API client.",
			},
			"client_description": schema.StringAttribute{
				Computed:    true,
				Description: "A human-readable description of the API client.",
			},
			"client_name": schema.StringAttribute{
				Computed:    true,
				Description: "A human-readable name for the API client.",
			},
			"client_type": schema.StringAttribute{
				Computed: true,
				Description: "Specifies the API client's ownership and credential management. " +
					"'CLIENT' indicates the creator owns and manages the credentials. " +
					"'USER_CLIENT' indicates another user owns the client and manages the credentials.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the API client.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 timestamp indicating when the API client was created.",
			},
			"credentials": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The API client's credentials.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"actions": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Specifies activities available on the API client's credentials.",
							Attributes: map[string]schema.Attribute{
								"delete": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can remove the credential.",
								},
								"activate": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can activate the credential.",
								},
								"deactivate": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can deactivate the credential.",
								},
								"edit_description": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can modify the credential's description.",
								},
								"edit_expiration": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether you can modify the credential's expiration date.",
								},
							},
						},
						"client_token": schema.StringAttribute{
							Computed:    true,
							Sensitive:   true,
							Description: "The part of the credential that identifies the API client.",
						},
						"created_on": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 timestamp indicating when the credential was created.",
						},
						"credential_id": schema.Int64Attribute{
							Computed:    true,
							Description: "A unique identifier for each credential.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A human-readable description for the API client.",
						},
						"expires_on": schema.StringAttribute{
							Computed:    true,
							Description: "The ISO 8601 timestamp indicating when the credential expires.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Whether a credential is `ACTIVE`, `INACTIVE`, or `DELETED`.",
						},
					},
				},
			},
			"group_access": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Specifies the API client's group access.",
				Attributes: map[string]schema.Attribute{
					"clone_authorized_user_groups": schema.BoolAttribute{
						Computed:    true,
						Description: "Sets the API client's group access the same as the authorized user.",
					},
					"groups": clientGroupsSchema(maxSupportedGroupNesting),
				},
			},
			"ip_acl": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Specifies the API client's IP list restriction.",
				Attributes: map[string]schema.Attribute{
					"cidr": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "IP addresses or CIDR blocks the API client can access.",
					},
					"enable": schema.BoolAttribute{
						Computed:    true,
						Description: "Enables the API client to access the IP access control list (ACL).",
					},
				},
			},
			"notification_emails": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Email addresses to notify users when credentials expire.",
			},
			"purge_options": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Configures the API client's access to the Fast Purge API defined in `apis`.",
				Attributes: map[string]schema.Attribute{
					"can_purge_by_cache_tag": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the API client can purge content by cache tag.",
					},
					"can_purge_by_cp_code": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the API client can purge content by CP code.",
					},
					"cp_code_access": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "CP codes the API client can purge.",
						Attributes: map[string]schema.Attribute{
							"all_current_and_new_cp_codes": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether the API can purge content by all current and new CP codes.",
							},
							"cp_codes": schema.ListAttribute{
								ElementType: types.Int64Type,
								Computed:    true,
								Description: "CP codes the API client can purge.",
							},
						},
					},
				},
			},
			"is_locked": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the API client is locked.",
			},
		},
	}
}

func (d *apiClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "IAM API Client DataSource Read")

	client := inst.Client(d.meta)

	var data clientModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := client.GetAPIClient(ctx, iam.GetAPIClientRequest{
		ClientID:    data.ClientID.ValueString(),
		Actions:     true,
		GroupAccess: true,
		APIAccess:   true,
		Credentials: true,
		IPACL:       true,
	})

	if err != nil {
		resp.Diagnostics.AddError("IAM get API Client failed", err.Error())
		return
	}

	diags := data.read(ctx, apiClient)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *clientModel) read(ctx context.Context, res *iam.GetAPIClientResponse) diag.Diagnostics {
	m.ClientID = types.StringValue(res.ClientID)
	m.AccessToken = types.StringValue(res.AccessToken)
	m.ActiveCredentialCount = types.Int64Value(res.ActiveCredentialCount)
	m.AllowAccountSwitch = types.BoolValue(res.AllowAccountSwitch)
	m.BaseURL = types.StringValue(res.BaseURL)
	m.CanAutoCreateCredential = types.BoolValue(res.CanAutoCreateCredential)
	m.ClientDescription = types.StringValue(res.ClientDescription)
	m.ClientName = types.StringValue(res.ClientName)
	m.ClientType = types.StringValue(string(res.ClientType))
	m.CreatedBy = types.StringValue(res.CreatedBy)
	m.CreatedDate = types.StringValue(res.CreatedDate.Format(time.RFC3339Nano))
	m.IsLocked = types.BoolValue(res.IsLocked)

	notificationEmailsObject, diags := types.ListValueFrom(ctx, types.StringType, res.NotificationEmails)
	if diags.HasError() {
		return diags
	}
	m.NotificationEmails = notificationEmailsObject

	authorizedUsers, diags := types.ListValueFrom(ctx, types.StringType, res.AuthorizedUsers)
	if diags.HasError() {
		return diags
	}
	m.AuthorizedUsers = authorizedUsers

	if res.Actions != nil {
		m.Actions = &apiClientActionsModel{
			Delete:            types.BoolValue(res.Actions.Delete),
			DeactivateAll:     types.BoolValue(res.Actions.DeactivateAll),
			Edit:              types.BoolValue(res.Actions.Edit),
			EditAPIs:          types.BoolValue(res.Actions.EditAPIs),
			EditAuth:          types.BoolValue(res.Actions.EditAuth),
			EditGroups:        types.BoolValue(res.Actions.EditGroups),
			EditIPACL:         types.BoolValue(res.Actions.EditIPACL),
			EditSwitchAccount: types.BoolValue(res.Actions.EditSwitchAccount),
			Lock:              types.BoolValue(res.Actions.Lock),
			Transfer:          types.BoolValue(res.Actions.Transfer),
			Unlock:            types.BoolValue(res.Actions.Unlock),
		}
	}

	apis := make([]apiClientAPIModel, 0, len(res.APIAccess.APIs))
	for _, api := range res.APIAccess.APIs {
		apis = append(apis, apiClientAPIModel{
			AccessLevel:      types.StringValue(api.AccessLevel),
			APIID:            types.Int64Value(api.APIID),
			APIName:          types.StringValue(api.APIName),
			Description:      types.StringValue(api.Description),
			DocumentationURL: types.StringValue(api.DocumentationURL),
			Endpoint:         types.StringValue(api.Endpoint),
		})
	}
	apisObject, diags := types.SetValueFrom(ctx, apiType(), apis)
	if diags.HasError() {
		return diags
	}
	m.APIAccess = &apiAccessModel{
		AllAccessibleAPIs: types.BoolValue(res.APIAccess.AllAccessibleAPIs),
		APIs:              apisObject,
	}

	if res.IPACL != nil {
		cidrObject, diags := types.ListValueFrom(ctx, types.StringType, res.IPACL.CIDR)
		if diags.HasError() {
			return diags
		}

		m.IPACL = &ipACLModel{
			CIDR:   cidrObject,
			Enable: types.BoolValue(res.IPACL.Enable),
		}
	}

	if res.PurgeOptions != nil {
		cpCodes, diags := types.ListValueFrom(ctx, types.Int64Type, res.PurgeOptions.CPCodeAccess.CPCodes)
		if diags.HasError() {
			return diags
		}
		cpCodeAccess := cpCodeAccessModel{
			AllCurrentAndNewCPCodes: types.BoolValue(res.PurgeOptions.CPCodeAccess.AllCurrentAndNewCPCodes),
			CPCodes:                 cpCodes,
		}
		cpCodesAccessObject, diags := types.ObjectValueFrom(ctx, cpCodeAccessType(), cpCodeAccess)
		if diags.HasError() {
			return diags
		}

		m.PurgeOptions = &purgeOptionsModel{
			CanPurgeByCacheTag: types.BoolValue(res.PurgeOptions.CanPurgeByCacheTag),
			CanPurgeByCPCode:   types.BoolValue(res.PurgeOptions.CanPurgeByCPCode),
			CPCodeAccess:       cpCodesAccessObject,
		}
	}

	subGroups, diags := readGroups(ctx, res.GroupAccess.Groups, maxSupportedGroupNesting)
	if diags.HasError() {
		return diags
	}
	m.GroupAccess = &groupAccessModel{
		CloneAuthorizedUserGroups: types.BoolValue(res.GroupAccess.CloneAuthorizedUserGroups),
		Groups:                    subGroups,
	}

	var credentials []getCredentialsModel
	for _, cred := range res.Credentials {
		credentialActions, diags := types.ObjectValueFrom(ctx, credentialActionsType(), credentialActionsModel{
			Delete:          types.BoolValue(cred.Actions.Delete),
			Activate:        types.BoolValue(cred.Actions.Activate),
			Deactivate:      types.BoolValue(cred.Actions.Deactivate),
			EditDescription: types.BoolValue(cred.Actions.EditDescription),
			EditExpiration:  types.BoolValue(cred.Actions.EditExpiration),
		})
		if diags.HasError() {
			return diags
		}
		credentials = append(credentials, getCredentialsModel{
			Actions:      credentialActions,
			ClientToken:  types.StringValue(cred.ClientToken),
			CreatedOn:    types.StringValue(cred.CreatedOn.Format(time.RFC3339Nano)),
			CredentialID: types.Int64Value(cred.CredentialID),
			Description:  types.StringValue(cred.Description),
			ExpiresOn:    types.StringValue(cred.ExpiresOn.Format(time.RFC3339Nano)),
			Status:       types.StringValue(string(cred.Status)),
		})
	}
	m.Credentials = credentials

	return nil
}

// nestedGroupsSchema builds a nested groups schema to the given depth
func clientGroupsSchema(depth int) schema.ListNestedAttribute {
	nestedSchema := schema.ListNestedAttribute{
		Computed:    true,
		Description: "Groups the API client can access.",
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
					Description: "Blocks the API client access to the group's child groups.",
				},
				"parent_group_id": schema.Int64Attribute{
					Computed:    true,
					Description: "Unique identifier for the parent group within the group tree.",
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
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{},
					},
				},
			},
		},
	}

	if depth > 1 {
		nestedSchema.NestedObject.Attributes["sub_groups"] = clientGroupsSchema(depth - 1)
	}

	return nestedSchema
}
