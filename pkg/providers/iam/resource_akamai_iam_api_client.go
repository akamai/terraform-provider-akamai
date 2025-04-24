package iam

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &apiClientResource{}
	_ resource.ResourceWithImportState = &apiClientResource{}
)

type apiClientResource struct {
	meta meta.Meta
}

const (
	invalidConfigurationAttribute = "Invalid configuration attribute"
)

// NewAPIClientResource returns new akamai_iam_api_client resource.
func NewAPIClientResource() resource.Resource { return &apiClientResource{} }

type (
	apiClientGroupModel struct {
		GroupID         types.Int64  `tfsdk:"group_id"`
		GroupName       types.String `tfsdk:"group_name"`
		IsBlocked       types.Bool   `tfsdk:"is_blocked"`
		ParentGroupID   types.Int64  `tfsdk:"parent_group_id"`
		RoleDescription types.String `tfsdk:"role_description"`
		RoleID          types.Int64  `tfsdk:"role_id"`
		RoleName        types.String `tfsdk:"role_name"`
		Subgroups       types.List   `tfsdk:"sub_groups"`
	}

	groupAccessModel struct {
		CloneAuthorizedUserGroups types.Bool `tfsdk:"clone_authorized_user_groups"`
		Groups                    types.List `tfsdk:"groups"`
	}

	apiClientAPIModel struct {
		AccessLevel      types.String `tfsdk:"access_level"`
		APIID            types.Int64  `tfsdk:"api_id"`
		APIName          types.String `tfsdk:"api_name"`
		Description      types.String `tfsdk:"description"`
		DocumentationURL types.String `tfsdk:"documentation_url"`
		Endpoint         types.String `tfsdk:"endpoint"`
	}

	apiAccessModel struct {
		AllAccessibleAPIs types.Bool `tfsdk:"all_accessible_apis"`
		APIs              types.Set  `tfsdk:"apis"`
	}

	purgeOptionsModel struct {
		CanPurgeByCacheTag types.Bool        `tfsdk:"can_purge_by_cache_tag"`
		CanPurgeByCPCode   types.Bool        `tfsdk:"can_purge_by_cp_code"`
		CPCodeAccess       cpCodeAccessModel `tfsdk:"cp_code_access"`
	}

	cpCodeAccessModel struct {
		AllCurrentAndNewCpCodes types.Bool `tfsdk:"all_current_and_new_cp_codes"`
		CPCodes                 types.List `tfsdk:"cp_codes"`
	}

	apiClientActionsModel struct {
		Delete            types.Bool `tfsdk:"delete"`
		DeactivateAll     types.Bool `tfsdk:"deactivate_all"`
		Edit              types.Bool `tfsdk:"edit"`
		EditAPIs          types.Bool `tfsdk:"edit_apis"`
		EditAuth          types.Bool `tfsdk:"edit_auth"`
		EditGroups        types.Bool `tfsdk:"edit_groups"`
		EditIPACL         types.Bool `tfsdk:"edit_ip_acl"`
		EditSwitchAccount types.Bool `tfsdk:"edit_switch_account"`
		Lock              types.Bool `tfsdk:"lock"`
		Transfer          types.Bool `tfsdk:"transfer"`
		Unlock            types.Bool `tfsdk:"unlock"`
	}

	credentialActionsModel struct {
		Delete          types.Bool `tfsdk:"delete"`
		Activate        types.Bool `tfsdk:"activate"`
		Deactivate      types.Bool `tfsdk:"deactivate"`
		EditDescription types.Bool `tfsdk:"edit_description"`
		EditExpiration  types.Bool `tfsdk:"edit_expiration"`
	}

	credentialsModel struct {
		Actions      credentialActionsModel `tfsdk:"actions"`
		ClientToken  types.String           `tfsdk:"client_token"`
		ClientSecret types.String           `tfsdk:"client_secret"`
		CreatedOn    types.String           `tfsdk:"created_on"`
		CredentialID types.Int64            `tfsdk:"credential_id"`
		Description  types.String           `tfsdk:"description"`
		ExpiresOn    types.String           `tfsdk:"expires_on"`
		Status       types.String           `tfsdk:"status"`
	}

	ipACLModel struct {
		CIDR   types.List `tfsdk:"cidr"`
		Enable types.Bool `tfsdk:"enable"`
	}

	apiClientResourceModel struct {
		AllowAccountSwitch      types.Bool         `tfsdk:"allow_account_switch"`
		APIAccess               apiAccessModel     `tfsdk:"api_access"`
		AuthorizedUsers         types.List         `tfsdk:"authorized_users"`
		CanCreateAutoCredential types.Bool         `tfsdk:"can_create_auto_credential"`
		ClientDescription       types.String       `tfsdk:"client_description"`
		ClientName              types.String       `tfsdk:"client_name"`
		ClientType              types.String       `tfsdk:"client_type"`
		GroupAccess             *groupAccessModel  `tfsdk:"group_access"`
		IPACL                   *ipACLModel        `tfsdk:"ip_acl"`
		NotificationEmails      types.List         `tfsdk:"notification_emails"`
		PurgeOptions            *purgeOptionsModel `tfsdk:"purge_options"`
		Lock                    types.Bool         `tfsdk:"lock"`
		AccessToken             types.String       `tfsdk:"access_token"`
		Actions                 types.Object       `tfsdk:"actions"`
		ActiveCredentialCount   types.Int64        `tfsdk:"active_credential_count"`
		BaseURL                 types.String       `tfsdk:"base_url"`
		ClientID                types.String       `tfsdk:"client_id"`
		CreatedBy               types.String       `tfsdk:"created_by"`
		CreatedDate             types.String       `tfsdk:"created_date"`
		Credential              types.Object       `tfsdk:"credential"`
		ID                      types.String       `tfsdk:"id"`
	}
)

func (r *apiClientResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_iam_api_client"
}

func (r *apiClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"allow_account_switch": schema.BoolAttribute{
				Optional:    true,
				Computed:    true, // needed for default
				Default:     booldefault.StaticBool(false),
				Description: "Enables the API client to manage more than one account.",
			},
			"api_access": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The APIs the API client can access.",
				Attributes: map[string]schema.Attribute{
					"all_accessible_apis": schema.BoolAttribute{
						Required:    true,
						Description: "Enables the API client to access a full set of available APIs.",
					},
					"apis": apisSchema(),
				},
			},
			"authorized_users": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The API client's valid users. When the 'client_type' is either 'CLIENT' or 'USER_CLIENT', you need to specify a single username in an array.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
			},
			"can_create_auto_credential": schema.BoolAttribute{
				Optional:    true,
				Computed:    true, // needed for default
				Default:     booldefault.StaticBool(false),
				Description: "Whether the API client can create a credential for a new API client. The default is false.",
			},
			"client_description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "A human-readable description of the API client.",
				Validators:  []validator.String{validators.NotEmptyString()},
			},
			"client_name": schema.StringAttribute{
				Required:    true,
				Description: "A human-readable name for the API client.",
				Validators:  []validator.String{validators.NotEmptyString()},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"client_type": schema.StringAttribute{
				Required: true,
				Description: "Specifies the API client's ownership and credential management. " +
					"'CLIENT' indicates the creator owns and manages the credentials. " +
					"'USER_CLIENT' indicates another user owns the client and manages the credentials.",
			},
			"group_access": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Specifies the API client's group access.",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"clone_authorized_user_groups": schema.BoolAttribute{
						Required:    true,
						Description: "Sets the API client's group access the same as the authorized user.",
					},
					"groups": groupsSchema(maxSupportedGroupNesting),
				},
			},
			"ip_acl": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Specifies the API client's IP list restriction.",
				Attributes: map[string]schema.Attribute{
					"cidr": schema.ListAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "IP addresses or CIDR blocks the API client can access.",
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"enable": schema.BoolAttribute{
						Required:    true,
						Description: "Enables the API client to access the IP access control list (ACL).",
					},
				},
			},
			"notification_emails": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(types.ListValueMust(
					types.StringType, []attr.Value{})),
				Description: "Email addresses to notify users when credentials expire.",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"purge_options": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Configures the API client to access the Fast Purge API. Provide it only if the `apis` attribute includes an `api_name` of `CCU API`.",
				Attributes: map[string]schema.Attribute{
					"can_purge_by_cache_tag": schema.BoolAttribute{
						Required:    true,
						Description: "Whether the API client can purge content by cache tag.",
					},
					"can_purge_by_cp_code": schema.BoolAttribute{
						Required:    true,
						Description: "Whether the API client can purge content by CP code.",
					},
					"cp_code_access": schema.SingleNestedAttribute{
						Required:    true,
						Description: "CP codes the API client can purge.",
						Attributes: map[string]schema.Attribute{
							"all_current_and_new_cp_codes": schema.BoolAttribute{
								Required:    true,
								Description: "Whether the API can purge content by all current and new CP codes.",
							},
							"cp_codes": schema.ListAttribute{
								ElementType: types.Int64Type,
								Optional:    true,
								Computed:    true, // needed for default
								Default:     listdefault.StaticValue(types.ListValueMust(types.Int64Type, []attr.Value{})),
								Description: "CP codes the API client can purge.",
							},
						},
					},
				},
			},
			"lock": schema.BoolAttribute{
				Optional:    true,
				Computed:    true, // needed for default
				Default:     booldefault.StaticBool(false),
				Description: "Whether to lock or unlock the API client.",
			},
			"access_token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Part of the client secret that identifies your API client and lets you access applications and resources.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"actions": actionsSchema(),
			"active_credential_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of credentials active for the API client. When the count is zero, you can delete the API client without interruption.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"base_url": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The base URL for the service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Computed:    true,
				Description: "A unique identifier for the API client.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who created the API client.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 timestamp indicating when the API client was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"credential": credentialSchema(),
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the resource, same as 'client_id'",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func apisSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Optional:    true,
		Computed:    true, // When all_accessible_apis is true, full list is provided in the response from the API
		Description: "The set of APIs the API client can access when `all_accessible_apis` is `false`.",
		PlanModifiers: []planmodifier.Set{
			setplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"access_level": schema.StringAttribute{
					Required:    true,
					Description: "The API client's access level on an API basis, either 'READ', 'READ-ONLY', 'READ-WRITE', 'CREDENTIAL-READ-ONLY', or 'CREDENTIAL-READ-WRITE'.",
					Validators: []validator.String{
						stringvalidator.OneOf("READ", "READ-ONLY", "READ-WRITE", "CREDENTIAL-READ-ONLY", "CREDENTIAL-READ-WRITE"),
					},
				},
				"api_id": schema.Int64Attribute{
					Required:    true,
					Description: "A unique identifier of the API.",
				},
				"api_name": schema.StringAttribute{
					Required:    true,
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "A human-readable name for the API.",
				},
				"description": schema.StringAttribute{
					Required:    true,
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "A human-readable description for the API.",
				},
				"documentation_url": schema.StringAttribute{
					Required:    true,
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "A link to more information about the API.",
				},
				"endpoint": schema.StringAttribute{
					Required:    true,
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "Specifies where the API can access resources.",
				},
			},
		},
	}
}

// nestedGroupsSchema builds a nested groups schema to the given depth
func groupsSchema(depth int) schema.ListNestedAttribute {
	nestedSchema := schema.ListNestedAttribute{
		Optional:    true,
		Computed:    true,
		Description: "Groups the API client can access.",
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"group_id": schema.Int64Attribute{
					Required:    true,
					Description: "A unique identifier for the group.",
				},
				"group_name": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "A human-readable name for the group.",
				},
				"is_blocked": schema.BoolAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
					Description: "Blocks the API client access to the group's child groups.",
				},
				"parent_group_id": schema.Int64Attribute{
					Computed: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
					Description: "A unique identifier for the parent group within the group tree. ",
				},
				"role_description": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "A human-readable description for the role to convey its use.",
				},
				"role_id": schema.Int64Attribute{
					Required:    true,
					Description: "A unique identifier for the role.",
				},
				"role_name": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "A human-readable name for the role.",
				},
				"sub_groups": schema.ListNestedAttribute{
					Optional:    true,
					Computed:    true,
					Description: "Children of the parent group.",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{},
					},
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	}

	if depth > 1 {
		nestedSchema.NestedObject.Attributes["sub_groups"] = groupsSchema(depth - 1)
	}

	return nestedSchema
}

func credentialSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
		Attributes: map[string]schema.Attribute{
			"actions": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Actions available on the API client's credentials.",
				Attributes: map[string]schema.Attribute{
					"delete": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can remove the credential.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"activate": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can activate the credential.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"deactivate": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can deactivate the credential.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"edit_description": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can modify the credential's description.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"edit_expiration": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether you can modify the credential's expiration date.",
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"client_token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The part of the credential that identifies the API client.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The client secret.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 timestamp indicating when the credential was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"credential_id": schema.Int64Attribute{
				Computed:    true,
				Description: "A unique identifier of the credential.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A human-readable description for the credential.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expires_on": schema.StringAttribute{
				Computed:    true,
				Description: "The ISO 8601 timestamp indicating when the credential expires. The default expiration date is two years from the creation date.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Whether a credential is 'ACTIVE', 'INACTIVE', or 'DELETED'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func actionsSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
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
				Description: "Whether you can update the APIs the API client can access.",
			},
			"edit_auth": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether you can update the 'groups' or 'apis' the API client can access.",
			},
			"edit_groups": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether you can update the groups the API client can access.",
			},
			"edit_ip_acl": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether you can update the 'ip acl' the API client can access.",
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
	}
}

func (r *apiClientResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	r.meta = meta.Must(req.ProviderData)
}

// ModifyPlan performs plan modification on a resource level.
func (r *apiClientResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "Modifying plan for API Client Resource")
	var state, plan *apiClientResourceModel
	if !request.Plan.Raw.IsNull() {
		response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
	if !request.State.Raw.IsNull() {
		response.Diagnostics.Append(request.State.Get(ctx, &state)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	// If all_accessible_apis is true on create, we should disallow purge_options.cp_code_access.cp_codes
	if request.State.Raw.IsNull() &&
		isKnown(plan.APIAccess.AllAccessibleAPIs) &&
		plan.APIAccess.AllAccessibleAPIs.ValueBool() &&
		isKnown(plan.PurgeOptions.CPCodeAccess.CPCodes) {

		tflog.Debug(ctx, "If 'all_accessible_apis' is true on create, we should disallow 'purge_options.cp_code_access.cp_codes'")
		var cpCodes []types.Int64
		plan.PurgeOptions.CPCodeAccess.CPCodes.ElementsAs(ctx, &cpCodes, false)
		if len(cpCodes) > 0 {
			response.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(path.Root("purge_options"), "invalid fields combination", "purge_options.cp_code_access.cp_codes cannot be provided when all_accessible_apis is true"))
			return
		}

		plan.APIAccess.APIs = types.SetUnknown(apiType())
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("api_access"), plan.APIAccess)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	if request.Plan.Raw.IsNull() || request.State.Raw.IsNull() {
		tflog.Debug(ctx, "Plan and state are null, skipping plan modification")
		return
	}

	// After import `credential` is nil, so we need to suppress changes on the terraform plan
	if state.Credential.IsNull() && plan.Credential.IsUnknown() {
		tflog.Debug(ctx, "After import 'credential' is nil, so we need to suppress changes on the terraform plan")
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("credential"), state.Credential)...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	modifyActions(ctx, state, plan, response)
	if response.Diagnostics.HasError() {
		return
	}

	if !request.Plan.Raw.Equal(request.State.Raw) &&
		isKnown(state.Lock) && state.Lock.ValueBool() &&
		isKnown(plan.Lock) && plan.Lock.ValueBool() {
		tflog.Debug(ctx, "Cannot change API client without unlocking it first")
		response.Diagnostics.AddError("lock", "You cannot change API client without unlocking it first.")
		return
	}

	modifyGroupAccess(ctx, state, plan, response)
	if response.Diagnostics.HasError() {
		return
	}

	// If all_accessible_apis is true, we need to remove the apis from the plan
	if isKnown(state.APIAccess.AllAccessibleAPIs) &&
		isKnown(plan.APIAccess.AllAccessibleAPIs) &&
		state.APIAccess.AllAccessibleAPIs.ValueBool() != plan.APIAccess.AllAccessibleAPIs.ValueBool() && plan.APIAccess.AllAccessibleAPIs.ValueBool() {
		tflog.Debug(ctx, "If 'all_accessible_apis' is true, we need to remove the 'apis' from the plan")
		plan.APIAccess.APIs = types.SetUnknown(apiType())
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("api_access"), plan.APIAccess)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
}

func isKnown(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func modifyGroupAccess(ctx context.Context, state *apiClientResourceModel, plan *apiClientResourceModel, response *resource.ModifyPlanResponse) {
	// If clone_authorized_user_groups is changed to false, we need to invalidate first element of `groups` from the plan
	if isKnown(state.GroupAccess.CloneAuthorizedUserGroups) &&
		isKnown(plan.GroupAccess.CloneAuthorizedUserGroups) &&
		state.GroupAccess.CloneAuthorizedUserGroups.ValueBool() != plan.GroupAccess.CloneAuthorizedUserGroups.ValueBool() && !plan.GroupAccess.CloneAuthorizedUserGroups.ValueBool() {
		tflog.Debug(ctx, "If 'clone_authorized_user_groups' is changed to false, we need to invalidate first element of 'groups' from the plan")
		var groups []apiClientGroupModel
		response.Diagnostics.Append(plan.GroupAccess.Groups.ElementsAs(ctx, &groups, false)...)
		if response.Diagnostics.HasError() {
			return
		}
		groups[0].RoleName = types.StringUnknown()
		groups[0].RoleDescription = types.StringUnknown()
		groups[0].GroupName = types.StringUnknown()
		groups[0].ParentGroupID = types.Int64Unknown()

		g, diags := types.ListValueFrom(ctx, groupsType(maxSupportedGroupNesting), groups)
		if diags.HasError() {
			response.Diagnostics.Append(diags...)
			return
		}
		plan.GroupAccess.Groups = g
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("group_access"), plan.GroupAccess)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
	// If clone_authorized_user_groups is changed to true, we need to invalidate all `groups` from the plan
	if isKnown(state.GroupAccess.CloneAuthorizedUserGroups) &&
		isKnown(plan.GroupAccess.CloneAuthorizedUserGroups) &&
		state.GroupAccess.CloneAuthorizedUserGroups.ValueBool() != plan.GroupAccess.CloneAuthorizedUserGroups.ValueBool() && plan.GroupAccess.CloneAuthorizedUserGroups.ValueBool() {
		tflog.Debug(ctx, "If 'clone_authorized_user_groups' is changed to true, we need to invalidate all 'groups' from the plan")
		plan.GroupAccess.Groups = types.ListUnknown(groupsType(maxSupportedGroupNesting))
		response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("group_access"), plan.GroupAccess)...)
		if response.Diagnostics.HasError() {
			return
		}
	}
}

func modifyActions(ctx context.Context, state *apiClientResourceModel, plan *apiClientResourceModel, response *resource.ModifyPlanResponse) {
	if isKnown(state.Lock) && isKnown(plan.Lock) {
		if state.Lock.ValueBool() != plan.Lock.ValueBool() {
			tflog.Debug(ctx, "If 'lock' is changed, we need to invalidate 'actions' from the plan")
			var actions apiClientActionsModel
			response.Diagnostics.Append(state.Actions.As(ctx, &actions, basetypes.ObjectAsOptions{})...)
			if response.Diagnostics.HasError() {
				return
			}
			// during lock or unlock those fields in `actions` are toggled
			actions.Lock = types.BoolUnknown()
			actions.Unlock = types.BoolUnknown()

			actionsObject, diags := types.ObjectValueFrom(ctx, actionsType(), actions)
			response.Diagnostics.Append(diags...)
			if response.Diagnostics.HasError() {
				return
			}

			response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("actions"), actionsObject)...)
			if response.Diagnostics.HasError() {
				return
			}
		} else {
			response.Diagnostics.Append(response.Plan.SetAttribute(ctx, path.Root("actions"), state.Actions)...)
			if response.Diagnostics.HasError() {
				return
			}
		}
	}
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r *apiClientResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	tflog.Debug(ctx, "Validating API Client Resource")
	var data apiClientResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groups []apiClientGroupModel
	resp.Diagnostics.Append(data.GroupAccess.Groups.ElementsAs(ctx, &groups, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.GroupAccess != nil && !data.GroupAccess.CloneAuthorizedUserGroups.ValueBool() && len(groups) == 0 {
		resp.Diagnostics.AddAttributeError(path.Root("group_access"), invalidConfigurationAttribute, "You must specify at least one group when 'clone_authorized_user_groups' is false")
		return
	}

	if data.GroupAccess != nil && data.GroupAccess.CloneAuthorizedUserGroups.ValueBool() && len(groups) != 0 {
		resp.Diagnostics.AddAttributeError(path.Root("group_access"), invalidConfigurationAttribute, "You cannot specify any group when 'clone_authorized_user_groups' is true")
		return
	}

	apis, diags := data.apisFromModel(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.APIAccess.AllAccessibleAPIs.ValueBool() && len(apis) == 0 {
		resp.Diagnostics.AddAttributeError(path.Root("api_access"), invalidConfigurationAttribute, "You must specify at least one API when 'all_accessible_apis' is false")
		return
	}

	if data.APIAccess.AllAccessibleAPIs.ValueBool() && len(apis) != 0 {
		resp.Diagnostics.AddAttributeError(path.Root("api_access"), invalidConfigurationAttribute, "You cannot specify any API when 'all_accessible_apis' is true")
		return
	}

	if data.APIAccess.AllAccessibleAPIs.ValueBool() && data.PurgeOptions == nil {
		resp.Diagnostics.AddAttributeError(path.Root("api_access"), invalidConfigurationAttribute, "You must specify 'purge_options' when 'all_accessible_apis' is true")
		return
	}

	for _, api := range apis {
		if api.APIName == "CCU APIs" && data.PurgeOptions == nil {
			resp.Diagnostics.AddAttributeError(path.Root("purge_options"), invalidConfigurationAttribute, "You must specify 'purge_options' when 'api_name' is 'CCU APIs'")
			return
		}
	}

	var cpCodes []int64
	if data.PurgeOptions != nil && isKnown(data.PurgeOptions.CPCodeAccess.CPCodes) {
		resp.Diagnostics.Append(data.PurgeOptions.CPCodeAccess.CPCodes.ElementsAs(ctx, &cpCodes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if data.PurgeOptions != nil && data.PurgeOptions.CPCodeAccess.AllCurrentAndNewCpCodes.ValueBool() && len(cpCodes) != 0 {
		resp.Diagnostics.AddAttributeError(path.Root("purge_options"), invalidConfigurationAttribute, "You cannot specify any CP Code when 'all_current_and_new_cp_codes' is true")
	}
}

func (r *apiClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating API Client Resource")
	var plan apiClientResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.create(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Creating API Client Resource failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiClientResource) create(ctx context.Context, plan *apiClientResourceModel) error {
	var notificationEmails []string
	if isKnown(plan.NotificationEmails) {
		diags := plan.NotificationEmails.ElementsAs(ctx, &notificationEmails, false)
		if diags.HasError() {
			return fmt.Errorf("failed to get notification emails: %v", diags)
		}
	}

	groups, diags := getGroupsFromModel(ctx, plan.GroupAccess.Groups)
	if diags.HasError() {
		return fmt.Errorf("failed to get groups: %v", diags)
	}
	access, diags := plan.getAPIAccess(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to get api access: %v", diags)
	}
	purgeOptions, diags := plan.getPurgeOptions(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to get purge options: %v", diags)
	}

	var authorizedUsers []string
	if isKnown(plan.AuthorizedUsers) {
		diags := plan.AuthorizedUsers.ElementsAs(ctx, &authorizedUsers, false)
		if diags.HasError() {
			return fmt.Errorf("failed to get authorized users: %v", diags)
		}
	}
	ipACL, diags := plan.getIPACL(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to get ip acl: %v", diags)
	}
	client := inst.Client(r.meta)
	createAPIClientResponse, err := client.CreateAPIClient(ctx, iam.CreateAPIClientRequest{
		AllowAccountSwitch:      plan.AllowAccountSwitch.ValueBool(),
		APIAccess:               *access,
		AuthorizedUsers:         authorizedUsers,
		CanAutoCreateCredential: plan.CanCreateAutoCredential.ValueBool(),
		ClientDescription:       plan.ClientDescription.ValueString(),
		ClientName:              plan.ClientName.ValueString(),
		ClientType:              iam.ClientType(plan.ClientType.ValueString()),
		CreateCredential:        true,
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: plan.GroupAccess.CloneAuthorizedUserGroups.ValueBool(),
			Groups:                    groups,
		},
		IPACL:              ipACL,
		NotificationEmails: notificationEmails,
		PurgeOptions:       purgeOptions,
	})
	if err != nil {
		return err
	}

	// If the notification emails are empty, we need to update the API client as the
	// the API fills the emails by default with the email of the user who created the API client.
	if len(notificationEmails) == 0 {
		_, err := client.UpdateAPIClient(ctx, iam.UpdateAPIClientRequest{
			ClientID: createAPIClientResponse.ClientID,
			Body: iam.UpdateAPIClientRequestBody{
				NotificationEmails: []string{},
				ClientName:         plan.ClientName.ValueString(),
				APIAccess:          *access,
				AuthorizedUsers:    authorizedUsers,
				ClientType:         iam.ClientType(plan.ClientType.ValueString()),
				GroupAccess: iam.GroupAccess{
					CloneAuthorizedUserGroups: plan.GroupAccess.CloneAuthorizedUserGroups.ValueBool(),
					Groups:                    groups,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	if plan.Lock.ValueBool() {
		_, err := client.LockAPIClient(ctx, iam.LockAPIClientRequest{
			ClientID: createAPIClientResponse.ClientID,
		})
		if err != nil {
			return err
		}
	}

	getAPIClientResponse, err := client.GetAPIClient(ctx, iam.GetAPIClientRequest{
		ClientID:    createAPIClientResponse.ClientID,
		Actions:     true,
		GroupAccess: true,
		APIAccess:   true,
		Credentials: true,
		IPACL:       true,
	})
	if err != nil {
		return err
	}

	diags = plan.setData(ctx, getAPIClientResponse, createAPIClientResponse)
	if diags.HasError() {
		return fmt.Errorf("failed to set data: %v", diags)
	}

	return nil
}

func (r *apiClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading API Client Resource")
	var state apiClientResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.read(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Reading API Client Resource failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *apiClientResource) read(ctx context.Context, data *apiClientResourceModel) error {
	client := inst.Client(r.meta)

	apiClient, err := client.GetAPIClient(ctx, iam.GetAPIClientRequest{
		ClientID:    data.ClientID.ValueString(),
		Actions:     true,
		GroupAccess: true,
		APIAccess:   true,
		Credentials: true,
		IPACL:       true,
	})
	if err != nil {
		return err
	}

	diags := data.setData(ctx, apiClient, nil)
	if diags.HasError() {
		return fmt.Errorf("failed to set data: %v", diags)
	}

	return nil
}

func (r *apiClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating API Client Resource")
	var plan, state apiClientResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// during update the lock value is taken from `IsLocked` and set, despite the fact that it is still an old value before lock/unlock
	plannedLock := plan.Lock.ValueBool()

	if !plannedLock && state.Lock.ValueBool() {
		tflog.Debug(ctx, "Unlocking API Client Resource")
		client := inst.Client(r.meta)
		_, err := client.UnlockAPIClient(ctx, iam.UnlockAPIClientRequest{
			ClientID: plan.ClientID.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Updating API Client Resource failed", err.Error())
			return
		}
	}

	if err := r.update(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Updating API Client Resource failed", err.Error())
		return
	}

	if plannedLock && !state.Lock.ValueBool() {
		tflog.Debug(ctx, "Locking API Client Resource")
		client := inst.Client(r.meta)
		_, err := client.LockAPIClient(ctx, iam.LockAPIClientRequest{
			ClientID: plan.ClientID.ValueString(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Updating API Client Resource failed", err.Error())
			return
		}
	}
	if plannedLock {
		if err := r.read(ctx, &plan); err != nil {
			resp.Diagnostics.AddError("Updating API Client Resource failed", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *apiClientResource) update(ctx context.Context, plan *apiClientResourceModel) error {
	var notificationEmails []string
	if isKnown(plan.NotificationEmails) {
		diags := plan.NotificationEmails.ElementsAs(ctx, &notificationEmails, false)
		if diags.HasError() {
			return fmt.Errorf("failed to get notification emails: %v", diags)
		}
	}

	groups, diags := getGroupsFromModel(ctx, plan.GroupAccess.Groups)
	if diags.HasError() {
		return fmt.Errorf("failed to get groups: %v", diags)
	}
	access, diags := plan.getAPIAccess(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to get api access: %v", diags)
	}
	purgeOptions, diags := plan.getPurgeOptions(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to get purge options: %v", diags)
	}
	var authorizedUsers []string
	if isKnown(plan.AuthorizedUsers) {
		diags := plan.AuthorizedUsers.ElementsAs(ctx, &authorizedUsers, false)
		if diags.HasError() {
			return fmt.Errorf("failed to get authorized users: %v", diags)
		}
	}
	ipACL, diags := plan.getIPACL(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to get ip acl: %v", diags)
	}
	client := inst.Client(r.meta)
	resp, err := client.UpdateAPIClient(ctx, iam.UpdateAPIClientRequest{
		ClientID: plan.ClientID.ValueString(),
		Body: iam.UpdateAPIClientRequestBody{
			AllowAccountSwitch:      plan.AllowAccountSwitch.ValueBool(),
			APIAccess:               *access,
			AuthorizedUsers:         authorizedUsers,
			CanAutoCreateCredential: plan.CanCreateAutoCredential.ValueBool(),
			ClientDescription:       plan.ClientDescription.ValueString(),
			ClientName:              plan.ClientName.ValueString(),
			ClientType:              iam.ClientType(plan.ClientType.ValueString()),
			GroupAccess: iam.GroupAccess{
				CloneAuthorizedUserGroups: plan.GroupAccess.CloneAuthorizedUserGroups.ValueBool(),
				Groups:                    groups,
			},
			IPACL:              ipACL,
			NotificationEmails: notificationEmails,
			PurgeOptions:       purgeOptions,
		},
	})
	if err != nil {
		return err
	}

	diags = plan.setData(ctx, (*iam.GetAPIClientResponse)(resp), nil)
	if diags.HasError() {
		return fmt.Errorf("failed to set data: %v", diags)
	}

	return err
}

func (r *apiClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting API Client Resource")

	var state *apiClientResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(r.meta)

	if err := client.DeleteAPIClient(ctx, iam.DeleteAPIClientRequest{
		ClientID: state.ClientID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Deleting API Client Resource %s failed", state.ClientID), err.Error())
		return
	}
}

func (r *apiClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing API Client Resource")

	data := &apiClientResourceModel{}

	// in import, we only need to set api client ID to allow read function to fill other attributes
	clientID := strings.TrimSpace(req.ID)
	data.ClientID = types.StringValue(clientID)
	data.NotificationEmails = types.ListUnknown(types.StringType)
	data.Credential = types.ObjectUnknown(credentialType())
	data.Actions = types.ObjectUnknown(actionsType())
	data.APIAccess.APIs = types.SetUnknown(apiType())
	data.AuthorizedUsers = types.ListUnknown(types.StringType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *apiClientResourceModel) getAPIAccess(ctx context.Context) (*iam.APIAccess, diag.Diagnostics) {
	var apiAccess iam.APIAccess
	var apis []iam.API
	// we should modify the list of apis only when all_accessible_apis is false
	if !m.APIAccess.AllAccessibleAPIs.ValueBool() {
		var diags diag.Diagnostics
		apis, diags = m.apisFromModel(ctx)
		if diags.HasError() {
			return nil, diags
		}
	}
	apiAccess = iam.APIAccess{
		AllAccessibleAPIs: m.APIAccess.AllAccessibleAPIs.ValueBool(),
		APIs:              apis,
	}

	return &apiAccess, nil
}

func (m *apiClientResourceModel) apisFromModel(ctx context.Context) ([]iam.API, diag.Diagnostics) {
	var apis []iam.API

	var apiModel []apiClientAPIModel
	if isKnown(m.APIAccess.APIs) {
		diags := m.APIAccess.APIs.ElementsAs(ctx, &apiModel, false)
		if diags.HasError() {
			return nil, diags
		}
	}

	apis = make([]iam.API, 0, len(apiModel))
	for _, api := range apiModel {
		apis = append(apis, iam.API{
			AccessLevel:      iam.AccessLevel(api.AccessLevel.ValueString()),
			APIID:            api.APIID.ValueInt64(),
			APIName:          api.APIName.ValueString(),
			Description:      api.Description.ValueString(),
			DocumentationURL: api.DocumentationURL.ValueString(),
			Endpoint:         api.Endpoint.ValueString(),
		})
	}

	return apis, nil
}

func (m *apiClientResourceModel) getIPACL(ctx context.Context) (*iam.IPACL, diag.Diagnostics) {
	var ipACL *iam.IPACL
	if m.IPACL != nil {

		var cidr []string
		diags := m.IPACL.CIDR.ElementsAs(ctx, &cidr, false)
		if diags.HasError() {
			return nil, diags
		}

		ipACL = &iam.IPACL{
			CIDR:   cidr,
			Enable: m.IPACL.Enable.ValueBool(),
		}
	}
	return ipACL, nil
}

func getGroupsFromModel(ctx context.Context, groupModels types.List) ([]iam.ClientGroup, diag.Diagnostics) {
	if groupModels.IsNull() || groupModels.IsUnknown() {
		return nil, nil
	}

	var groups []apiClientGroupModel
	diags := groupModels.ElementsAs(ctx, &groups, false)
	if diags.HasError() {
		return nil, diags
	}

	if len(groups) == 0 {
		return nil, nil
	}
	result := make([]iam.ClientGroup, 0, len(groups))
	for _, group := range groups {
		subGroups, diags := getGroupsFromModel(ctx, group.Subgroups)
		if diags.HasError() {
			return nil, diags
		}
		result = append(result, iam.ClientGroup{
			GroupID:       group.GroupID.ValueInt64(),
			IsBlocked:     group.IsBlocked.ValueBool(),
			ParentGroupID: group.ParentGroupID.ValueInt64(),
			RoleID:        group.RoleID.ValueInt64(),
			Subgroups:     subGroups,
		})
	}
	return result, nil
}

func (m *apiClientResourceModel) getPurgeOptions(ctx context.Context) (*iam.PurgeOptions, diag.Diagnostics) {
	var purgeOptions *iam.PurgeOptions

	if m.PurgeOptions != nil {
		var cpCodes []int64
		if m.PurgeOptions != nil && isKnown(m.PurgeOptions.CPCodeAccess.CPCodes) {
			diags := m.PurgeOptions.CPCodeAccess.CPCodes.ElementsAs(ctx, &cpCodes, false)
			if diags.HasError() {
				return nil, diags
			}
		}

		purgeOptions = &iam.PurgeOptions{
			CanPurgeByCacheTag: m.PurgeOptions.CanPurgeByCacheTag.ValueBool(),
			CanPurgeByCPCode:   m.PurgeOptions.CanPurgeByCPCode.ValueBool(),
			CPCodeAccess: iam.CPCodeAccess{
				AllCurrentAndNewCPCodes: m.PurgeOptions.CPCodeAccess.AllCurrentAndNewCpCodes.ValueBool(),
				CPCodes:                 cpCodes,
			},
		}
	}
	return purgeOptions, nil
}

func (m *apiClientResourceModel) setData(ctx context.Context, getResponse *iam.GetAPIClientResponse, createResponse *iam.CreateAPIClientResponse) diag.Diagnostics {
	m.ClientID = types.StringValue(getResponse.ClientID)
	m.ClientType = types.StringValue(string(getResponse.ClientType))
	m.ClientName = types.StringValue(getResponse.ClientName)
	m.ClientDescription = types.StringValue(getResponse.ClientDescription)
	m.CreatedBy = types.StringValue(getResponse.CreatedBy)
	m.CreatedDate = types.StringValue(getResponse.CreatedDate.Format(time.RFC3339Nano))
	m.CanCreateAutoCredential = types.BoolValue(getResponse.CanAutoCreateCredential)
	m.AccessToken = types.StringValue(getResponse.AccessToken)
	m.Lock = types.BoolValue(getResponse.IsLocked)
	m.AllowAccountSwitch = types.BoolValue(getResponse.AllowAccountSwitch)
	m.BaseURL = types.StringValue(getResponse.BaseURL)
	m.ActiveCredentialCount = types.Int64Value(getResponse.ActiveCredentialCount)

	authorizedUsers, diags := types.ListValueFrom(ctx, types.StringType, getResponse.AuthorizedUsers)
	if diags.HasError() {
		return diags
	}
	m.AuthorizedUsers = authorizedUsers

	if getResponse.PurgeOptions != nil {
		cpCodes, diags := types.ListValueFrom(ctx, types.Int64Type, getResponse.PurgeOptions.CPCodeAccess.CPCodes)
		if diags.HasError() {
			return diags
		}
		m.PurgeOptions = &purgeOptionsModel{
			CanPurgeByCacheTag: types.BoolValue(getResponse.PurgeOptions.CanPurgeByCacheTag),
			CanPurgeByCPCode:   types.BoolValue(getResponse.PurgeOptions.CanPurgeByCPCode),
			CPCodeAccess: cpCodeAccessModel{
				AllCurrentAndNewCpCodes: types.BoolValue(getResponse.PurgeOptions.CPCodeAccess.AllCurrentAndNewCPCodes),
				CPCodes:                 cpCodes,
			},
		}
	}

	notificationEmailsObject, diags := types.ListValueFrom(ctx, types.StringType, getResponse.NotificationEmails)
	if diags.HasError() {
		return diags
	}
	m.NotificationEmails = notificationEmailsObject

	if getResponse.IPACL != nil {
		cidrObject, diags := types.ListValueFrom(ctx, types.StringType, getResponse.IPACL.CIDR)
		if diags.HasError() {
			return diags
		}

		m.IPACL = &ipACLModel{
			CIDR:   cidrObject,
			Enable: types.BoolValue(getResponse.IPACL.Enable),
		}
	}

	subGroups, diags := readGroups(ctx, getResponse.GroupAccess.Groups, maxSupportedGroupNesting)
	if diags.HasError() {
		return diags
	}
	m.GroupAccess = &groupAccessModel{
		CloneAuthorizedUserGroups: types.BoolValue(getResponse.GroupAccess.CloneAuthorizedUserGroups),
		Groups:                    subGroups,
	}

	var actions apiClientActionsModel
	if getResponse.Actions != nil {
		actions = apiClientActionsModel{
			Delete:            types.BoolValue(getResponse.Actions.Delete),
			DeactivateAll:     types.BoolValue(getResponse.Actions.DeactivateAll),
			Edit:              types.BoolValue(getResponse.Actions.Edit),
			EditAPIs:          types.BoolValue(getResponse.Actions.EditAPIs),
			EditAuth:          types.BoolValue(getResponse.Actions.EditAuth),
			EditGroups:        types.BoolValue(getResponse.Actions.EditGroups),
			EditIPACL:         types.BoolValue(getResponse.Actions.EditIPACL),
			EditSwitchAccount: types.BoolValue(getResponse.Actions.EditSwitchAccount),
			Lock:              types.BoolValue(getResponse.Actions.Lock),
			Transfer:          types.BoolValue(getResponse.Actions.Transfer),
			Unlock:            types.BoolValue(getResponse.Actions.Unlock),
		}

		actionsObject, diags := types.ObjectValueFrom(ctx, actionsType(), actions)
		if diags.HasError() {
			return diags
		}
		m.Actions = actionsObject
	}

	var credModel credentialsModel
	var found bool
	if createResponse != nil {
		if len(createResponse.Credentials) == 0 {
			diags.AddError("Error setting response", "no credentials found in the create response")
			return diags
		}
		credential := createResponse.Credentials[0]
		credModel = credentialsModel{
			Actions: credentialActionsModel{
				Delete:          types.BoolValue(credential.Actions.Delete),
				Activate:        types.BoolValue(credential.Actions.Activate),
				Deactivate:      types.BoolValue(credential.Actions.Deactivate),
				EditDescription: types.BoolValue(credential.Actions.EditDescription),
				EditExpiration:  types.BoolValue(credential.Actions.EditExpiration),
			},
			ClientSecret: types.StringValue(credential.ClientSecret),
			ClientToken:  types.StringValue(credential.ClientToken),
			CreatedOn:    types.StringValue(credential.CreatedOn.Format(time.RFC3339Nano)),
			CredentialID: types.Int64Value(credential.CredentialID),
			Description:  types.StringValue(credential.Description),
			ExpiresOn:    types.StringValue(credential.ExpiresOn.Format(time.RFC3339Nano)),
			Status:       types.StringValue(string(credential.Status)),
		}
		found = true
	} else {
		var oldCredential credentialsModel
		if isKnown(m.Credential) {
			diags = m.Credential.As(ctx, &oldCredential, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return diags
			}

			for _, credential := range getResponse.Credentials {
				if oldCredential.CredentialID.ValueInt64() == credential.CredentialID {

					credModel = credentialsModel{
						Actions: credentialActionsModel{
							Delete:          types.BoolValue(credential.Actions.Delete),
							Activate:        types.BoolValue(credential.Actions.Activate),
							Deactivate:      types.BoolValue(credential.Actions.Deactivate),
							EditDescription: types.BoolValue(credential.Actions.EditDescription),
							EditExpiration:  types.BoolValue(credential.Actions.EditExpiration),
						},
						ClientSecret: oldCredential.ClientSecret,
						ClientToken:  types.StringValue(credential.ClientToken),
						CreatedOn:    types.StringValue(credential.CreatedOn.Format(time.RFC3339Nano)),
						CredentialID: types.Int64Value(credential.CredentialID),
						Description:  types.StringValue(credential.Description),
						ExpiresOn:    types.StringValue(credential.ExpiresOn.Format(time.RFC3339Nano)),
						Status:       types.StringValue(string(credential.Status)),
					}
					found = true
					break
				}
			}
		}
	}
	if found {
		credentialObject, diags := types.ObjectValueFrom(ctx, credentialType(), credModel)
		if diags.HasError() {
			return diags
		}
		m.Credential = credentialObject
	} else {
		m.Credential = types.ObjectNull(credentialType())
	}

	apis := make([]apiClientAPIModel, 0, len(getResponse.APIAccess.APIs))
	for _, api := range getResponse.APIAccess.APIs {
		apis = append(apis, apiClientAPIModel{
			AccessLevel:      types.StringValue(string(api.AccessLevel)),
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
	m.APIAccess = apiAccessModel{
		AllAccessibleAPIs: types.BoolValue(getResponse.APIAccess.AllAccessibleAPIs),
		APIs:              apisObject,
	}

	m.ID = types.StringValue(getResponse.ClientID)

	return nil
}

func actionsType() map[string]attr.Type {
	return actionsSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func apiType() types.ObjectType {
	return apisSchema().NestedObject.Type().(types.ObjectType)
}

func credentialType() map[string]attr.Type {
	return credentialSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func groupsType(level int) types.ObjectType {
	return groupsSchema(level).NestedObject.Type().(types.ObjectType)
}

func readGroups(ctx context.Context, resp []iam.ClientGroup, level int) (types.List, diag.Diagnostics) {
	if len(resp) == 0 {
		result := make([]apiClientGroupModel, 0)
		return types.ListValueFrom(ctx, groupsSchema(level).NestedObject.Type().(types.ObjectType), result)
	}
	result := make([]apiClientGroupModel, 0, len(resp))
	for _, group := range resp {
		subGroups, diags := readGroups(ctx, group.Subgroups, level-1)
		if diags.HasError() {
			return types.ListNull(groupsType(level)), diags
		}
		result = append(result, apiClientGroupModel{
			GroupID:         types.Int64Value(group.GroupID),
			GroupName:       types.StringValue(group.GroupName),
			IsBlocked:       types.BoolValue(group.IsBlocked),
			ParentGroupID:   types.Int64Value(group.ParentGroupID),
			RoleDescription: types.StringValue(group.RoleDescription),
			RoleID:          types.Int64Value(group.RoleID),
			RoleName:        types.StringValue(group.RoleName),
			Subgroups:       subGroups,
		})
	}
	return types.ListValueFrom(ctx, groupsType(level), result)
}
