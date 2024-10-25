package cloudaccess

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &KeyResource{}
	_ resource.ResourceWithConfigure   = &KeyResource{}
	_ resource.ResourceWithModifyPlan  = &KeyResource{}
	_ resource.ResourceWithImportState = &KeyResource{}

	activationTimeout = 60 * time.Minute
	updateTimeout     = 60 * time.Minute
	deleteTimeout     = 60 * time.Minute
	pollingInterval   = 1 * time.Minute
)

const (
	readError = "could not read access key from API"

	assignedToPropertyError = "cannot delete version: %d of access key %d assigned to property"

	diagErrAccessKeyNotFound = "Cannot Find Access key: %d"
)

// KeyResource represents akamai_cloudaccess_key resource
type KeyResource struct {
	meta meta.Meta
}

// KeyResourceModel represents model of akamai_cloudaccess_key resource
type KeyResourceModel struct {
	AccessKeyName        types.String   `tfsdk:"access_key_name"`
	AuthenticationMethod types.String   `tfsdk:"authentication_method"`
	ContractID           types.String   `tfsdk:"contract_id"`
	GroupID              types.Int64    `tfsdk:"group_id"`
	PrimaryGUID          types.String   `tfsdk:"primary_guid"`
	CredentialsA         *Credentials   `tfsdk:"credentials_a"`
	CredentialsB         *Credentials   `tfsdk:"credentials_b"`
	NetworkConfig        NetworkConfig  `tfsdk:"network_configuration"`
	AccessKeyUID         types.Int64    `tfsdk:"access_key_uid"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

// Credentials represent set of attributes for specific access key versions
type Credentials struct {
	CloudAccessKeyID     types.String `tfsdk:"cloud_access_key_id"`
	CloudSecretAccessKey types.String `tfsdk:"cloud_secret_access_key"`
	PrimaryKey           types.Bool   `tfsdk:"primary_key"`
	Version              types.Int64  `tfsdk:"version"`
	VersionGUID          types.String `tfsdk:"version_guid"`
}

// NetworkConfig represents set of attributes for network configuration
type NetworkConfig struct {
	AdditionalCDN   types.String `tfsdk:"additional_cdn"`
	SecurityNetwork types.String `tfsdk:"security_network"`
}

// NewKeyResource returns new cloudaccess key resource
func NewKeyResource() resource.Resource {
	return &KeyResource{}
}

// ModifyPlan implements resource.ResourceWithModifyPlan
func (r *KeyResource) ModifyPlan(ctx context.Context, request resource.ModifyPlanRequest, response *resource.ModifyPlanResponse) {

	var state, plan *KeyResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	if state == nil && plan.CredentialsA == nil && plan.CredentialsB == nil {
		response.Diagnostics.AddError("at least one credentials are required for creation", "`credentials_a` or `credentials_b` must be specified")
		return
	}

	if plan != nil && plan.CredentialsA != nil && plan.CredentialsB != nil &&
		plan.CredentialsA.PrimaryKey.ValueBool() && plan.CredentialsB.PrimaryKey.ValueBool() {
		response.Diagnostics.AddError("primary version of access key error", "only one pair of access key version can have 'primary_key' set as 'true'")
		return
	}

	if plan != nil && plan.CredentialsA != nil && plan.CredentialsB != nil &&
		plan.CredentialsA.CloudAccessKeyID == plan.CredentialsB.CloudAccessKeyID {
		response.Diagnostics.AddError("cloud access key id of access key error", "'cloud_access_key_id' should be unique for each pair of credentials")
		return
	}

	if state != nil && plan != nil && changedOrderOfCredentials(state, plan) {
		response.Diagnostics.AddError("access key credentials error", "cannot change order of `credentials_a` and `credentials_b`")
		return
	}

	if state != nil && plan != nil && checkIfSecretChangedAndWasNotEmpty(state, plan) {
		response.Diagnostics.AddError("access key credentials error", "cannot update cloud access secret without update of cloud access key id, expect update of secret after import with no API calls")
		return
	}

	if state != nil && plan != nil && onlyTimeoutChanged(state, plan) {
		state.Timeouts = plan.Timeouts
		return
	}
}

// Metadata implements resource.Resource.
func (r *KeyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_cloudaccess_key"
}

// Schema implements resource.Resource.
func (r *KeyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"access_key_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the access key.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"authentication_method": schema.StringAttribute{
				Required:    true,
				Description: "The type of cloud provider signing process used to authenticate API requests. Two options are available: \"AWS4_HMAC_SHA256\" or \"GOOG4_HMAC_SHA256\".",
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(cloudaccess.AuthAWS), string(cloudaccess.AuthGOOG)),
				},
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier for the contract assigned to the access key",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					modifiers.PreventStringUpdate(),
				},
			},
			"group_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier assigned to the access control group assigned to the access key",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"primary_guid": schema.StringAttribute{
				Computed:    true,
				Description: "Value of `version_guid` field for credentials marked as primary",
			},
			"credentials_a": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "The combination of a `cloud_access_key_id` and a `cloud_secret_access_key` used to sign API requests. This pair can be identified as access key version. Access key can contain only two access key versions at specific time (defined as credentialsA and credentialsB).",
				Attributes: map[string]schema.Attribute{
					"cloud_access_key_id": schema.StringAttribute{
						Description: "Access key id from cloud provider which is used to sign API requests",
						Required:    true,
					},
					"cloud_secret_access_key": schema.StringAttribute{
						Description: "Cloud Access secret from cloud provider which is used to sign API requests",
						Required:    true,
						Sensitive:   true,
					},
					"primary_key": schema.BoolAttribute{
						Description: "Boolean value which helps to define if credentials should be assigned to property",
						Required:    true,
					},
					"version": schema.Int64Attribute{
						Description: "Numeric access key version associated with specific pair of cloud access credentials used to sign API requests",
						Computed:    true,
					},
					"version_guid": schema.StringAttribute{
						Description: "The unique identifier assigned to specific access key version",
						Computed:    true,
					},
				},
			},
			"credentials_b": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "The combination of a `cloud_access_key_id` and a `cloud_secret_access_key` used to sign API requests. This pair can be identified as access key version. Access key can contain only two access key versions at specific time (defined as credentialsA and credentialsB).",
				Attributes: map[string]schema.Attribute{
					"cloud_access_key_id": schema.StringAttribute{
						Description: "Access key id from cloud provider which is used to sign API requests",
						Required:    true,
					},
					"cloud_secret_access_key": schema.StringAttribute{
						Description: "Cloud Access secret from cloud provider which is used to sign API requests",
						Required:    true,
						Sensitive:   true,
					},
					"primary_key": schema.BoolAttribute{
						Description: "Boolean value which helps to define if credentials should be assigned to property",
						Required:    true,
					},
					"version": schema.Int64Attribute{
						Description: "Numeric access key version associated with specific pair of cloud access credentials used to sign API requests",
						Computed:    true,
					},
					"version_guid": schema.StringAttribute{
						Description: "The unique identifier assigned to specific access key version",
						Computed:    true,
					},
				},
			},
			"network_configuration": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The secure networks that you assigned the access key to during creation",
				Attributes: map[string]schema.Attribute{
					"additional_cdn": schema.StringAttribute{
						Optional:    true,
						Description: "Additional type of the deployment network that the access key will be deployed to.",
						Validators: []validator.String{
							stringvalidator.OneOf(string(cloudaccess.ChinaCDN), string(cloudaccess.RussiaCDN)),
						},
						PlanModifiers: []planmodifier.String{
							modifiers.PreventStringUpdate(),
						},
					},
					"security_network": schema.StringAttribute{
						Required:    true,
						Description: "The API deploys the access key to this secure network",
						Validators: []validator.String{
							stringvalidator.OneOf(string(cloudaccess.NetworkStandard), string(cloudaccess.NetworkEnhanced)),
						},
						PlanModifiers: []planmodifier.String{
							modifiers.PreventStringUpdate(),
						},
					},
				},
			},
			"access_key_uid": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique identifier Akamai assigns to an access key.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Delete:            true,
				Update:            true,
				Create:            true,
				CreateDescription: "Optional configurable resource create timeout. By default it's 60 minutes with 1 minute polling interval.",
				DeleteDescription: "Optional configurable resource delete timeout. By default it's 60 minutes with 1 minute polling interval.",
				UpdateDescription: "Optional configurable resource update timeout. By default it's 60 minutes with 1 minute polling interval.",
			}),
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *KeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"unexpected resource configure type",
				fmt.Sprintf("expected meta.Meta, got: %T. please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	r.meta = meta.Must(req.ProviderData)
}

// onlyTimeoutChanged defines if timeout is the only parameter which changed between plan and state
func onlyTimeoutChanged(state, plan *KeyResourceModel) bool {
	return state != nil && plan != nil &&
		!state.Timeouts.Equal(plan.Timeouts)
}

// Create implements resource.Resource.
func (r *KeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Access Key Resource")
	var diags diag.Diagnostics
	var plan *KeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	createTimeout, diags := plan.Timeouts.Create(ctx, activationTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	plan, diagnostics := r.create(ctx, plan)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	// save partial data to state - it will allow taint flow after further failure
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if plan.CredentialsA != nil && plan.CredentialsB != nil {
		plan, diags = r.createVersion(ctx, plan, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan = r.setupPrimaryGUID(plan)
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// setupPrimaryGUID calculates `primary_guid`based on `primary_key` and 'version_guid` parameters
func (r *KeyResource) setupPrimaryGUID(state *KeyResourceModel) *KeyResourceModel {
	// setting `primaryGuid` based on primary_key flag
	if state.CredentialsA != nil && state.CredentialsA.PrimaryKey.ValueBool() {
		state.PrimaryGUID = state.CredentialsA.VersionGUID
		return state
	}
	if state.CredentialsB != nil && state.CredentialsB.PrimaryKey.ValueBool() {
		state.PrimaryGUID = state.CredentialsB.VersionGUID
		return state
	}
	state.PrimaryGUID = types.StringValue("")
	return state
}

func (r *KeyResource) create(ctx context.Context, plan *KeyResourceModel) (*KeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := Client(r.meta)
	creationKeyWithCredA := plan.CredentialsA != nil
	resp, err := client.CreateAccessKey(ctx, plan.buildCreateKeyRequest(creationKeyWithCredA))
	if err != nil {
		diags.AddError("create access key failed", err.Error())
		return nil, diags
	}

	return r.waitUntilActivationCompleted(ctx, resp.RequestID, resp.RetryAfter, plan, creationKeyWithCredA)
}

func (r *KeyResource) createVersion(ctx context.Context, plan *KeyResourceModel, useCredentialA bool) (*KeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := Client(r.meta)

	resp, err := client.CreateAccessKeyVersion(ctx, plan.buildCreateKeyVersionRequest(useCredentialA))
	if err != nil {
		// If version creation fails whole resource should be tainted
		diags.AddError("create access key version failed", err.Error())
		return nil, diags
	}

	return r.waitUntilVersionCreatedCompleted(ctx, resp.RequestID, resp.RetryAfter, plan, useCredentialA)
}

// Read implements resource.Resource.
func (r *KeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Access Key Resource")
	var oldState *KeyResourceModel
	var diags diag.Diagnostics
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = r.read(ctx, oldState)
	keyNotFoundDiags := diag.NewErrorDiagnostic("get access key error", fmt.Sprintf(diagErrAccessKeyNotFound, oldState.AccessKeyUID))
	if diags.Contains(keyNotFoundDiags) {
		resp.State.RemoveResource(ctx)
		return
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	oldState = r.setupPrimaryGUID(oldState)

	resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
}

func (r *KeyResource) read(ctx context.Context, data *KeyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	client := Client(r.meta)

	result, err := client.GetAccessKey(ctx, cloudaccess.AccessKeyRequest{
		AccessKeyUID: data.AccessKeyUID.ValueInt64(),
	})
	if errors.Is(err, cloudaccess.ErrGetAccessKey) {
		diags.AddError("get access key error", fmt.Sprintf(diagErrAccessKeyNotFound, data.AccessKeyUID))
		return diags
	}
	if err != nil {
		diags.AddError("get access key failed", err.Error())
		return diags
	}
	data.populateModelFromAccessKey(result)

	versions, err := client.ListAccessKeyVersions(ctx, cloudaccess.ListAccessKeyVersionsRequest{AccessKeyUID: data.AccessKeyUID.ValueInt64()})
	if err != nil {
		diags.AddError("list access key versions failed", err.Error())
		return diags
	}
	return data.populateModelFromVersionsList(versions)
}

// Update implements resource.Resource.
func (r *KeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Access Key Resource")
	var diags diag.Diagnostics
	var plan *KeyResourceModel
	client := Client(r.meta)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, updateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var oldState *KeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	initializeCredentialVersions(oldState, plan)

	if onlyTimeoutChanged(oldState, plan) {
		oldState.Timeouts = plan.Timeouts
		resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
		return
	}
	plan.AccessKeyUID = oldState.AccessKeyUID
	if oldState.AccessKeyName != plan.AccessKeyName {
		resp.Diagnostics.Append(r.updateAccessKey(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	updateInStateCredentialA, updateInStateCredentialB := keyVersionRequireUpdateInState(oldState, plan)
	if updateInStateCredentialA {
		oldState.CredentialsA.CloudSecretAccessKey = plan.CredentialsA.CloudSecretAccessKey
	}
	if updateInStateCredentialB {
		oldState.CredentialsB.CloudSecretAccessKey = plan.CredentialsB.CloudSecretAccessKey
	}
	if updateInStateCredentialA || updateInStateCredentialB {
		resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
	}

	deleteCredentialsA, deleteCredentialsB := keyVersionRequiresDeletion(oldState, plan)
	if deleteCredentialsA || deleteCredentialsB {
		diags = r.deleteVersion(ctx, oldState, client, resp, diags, deleteCredentialsA, deleteCredentialsB)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if deleteCredentialsA {
			oldState.CredentialsA = nil
		}
		if deleteCredentialsB {
			oldState.CredentialsB = nil
		}
	}

	createCredentialsA, createCredentialsB := keyVersionRequiresCreation(oldState, plan)
	if createCredentialsA {
		plan, diags = r.createVersion(ctx, plan, true)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if createCredentialsB {
		plan, diags = r.createVersion(ctx, plan, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	plan = r.setupPrimaryGUID(plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyResource) deleteVersion(ctx context.Context, oldState *KeyResourceModel, client cloudaccess.CloudAccess, resp *resource.UpdateResponse, diags diag.Diagnostics, deleteCredentialsA, deleteCredentialsB bool) diag.Diagnostics {
	var versionsToDelete []int64
	if deleteCredentialsA {
		versionToDelete := oldState.CredentialsA.Version.ValueInt64()
		hasProperty, diags := isVersionAssignedToProperty(ctx, client, oldState.AccessKeyUID.ValueInt64(), versionToDelete)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return diags
		}
		if hasProperty {
			resp.Diagnostics.AddError("version assigned to property error", fmt.Sprintf(assignedToPropertyError, versionToDelete, oldState.AccessKeyUID.ValueInt64()))
			return diags
		}
		versionsToDelete = append(versionsToDelete, versionToDelete)
	}
	if deleteCredentialsB {
		versionToDelete := oldState.CredentialsB.Version.ValueInt64()
		hasProperty, diags := isVersionAssignedToProperty(ctx, client, oldState.AccessKeyUID.ValueInt64(), versionToDelete)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return diags
		}
		if hasProperty {
			resp.Diagnostics.AddError("version assigned to property error", fmt.Sprintf(assignedToPropertyError, versionToDelete, oldState.AccessKeyUID.ValueInt64()))
			return diags
		}
		versionsToDelete = append(versionsToDelete, versionToDelete)
	}
	for _, version := range versionsToDelete {
		diags := r.deleteKeyVersion(ctx, oldState, version, diags)
		if diags != nil {
			resp.Diagnostics.Append(diags...)
			return diags
		}
	}
	return diags
}

func changedOrderOfCredentials(oldState, plan *KeyResourceModel) bool {
	if oldState.CredentialsA != nil && plan.CredentialsA != nil && oldState.CredentialsB != nil && plan.CredentialsB != nil && (oldState.CredentialsA.CloudAccessKeyID == plan.CredentialsB.CloudAccessKeyID && oldState.CredentialsB.CloudAccessKeyID == plan.CredentialsA.CloudAccessKeyID) {
		return true
	}
	return false
}

func checkIfSecretChangedAndWasNotEmpty(oldState, plan *KeyResourceModel) bool {
	if oldState.CredentialsA != nil && plan.CredentialsA != nil &&
		oldState.CredentialsA.CloudAccessKeyID.ValueString() == plan.CredentialsA.CloudAccessKeyID.ValueString() &&
		oldState.CredentialsA.CloudSecretAccessKey.ValueString() != "" && oldState.CredentialsA.CloudSecretAccessKey.ValueString() != plan.CredentialsA.CloudSecretAccessKey.ValueString() {
		return true
	}
	if oldState.CredentialsB != nil && plan.CredentialsB != nil &&
		oldState.CredentialsB.CloudAccessKeyID.ValueString() == plan.CredentialsB.CloudAccessKeyID.ValueString() &&
		oldState.CredentialsB.CloudSecretAccessKey.ValueString() != "" && oldState.CredentialsB.CloudSecretAccessKey.ValueString() != plan.CredentialsB.CloudSecretAccessKey.ValueString() {
		return true
	}
	return false
}

func keyVersionRequireUpdateInState(oldState *KeyResourceModel, data *KeyResourceModel) (bool, bool) {
	var updateCredA, updateCredB bool
	if oldState.CredentialsA != nil && data.CredentialsA != nil && oldState.CredentialsA.CloudAccessKeyID == data.CredentialsA.CloudAccessKeyID && oldState.CredentialsA.CloudSecretAccessKey.ValueString() == "" {
		updateCredA = true
	}
	if oldState.CredentialsB != nil && data.CredentialsB != nil && oldState.CredentialsB.CloudAccessKeyID == data.CredentialsB.CloudAccessKeyID && oldState.CredentialsB.CloudSecretAccessKey.ValueString() == "" {
		updateCredB = true
	}
	return updateCredA, updateCredB
}

func initializeCredentialVersions(oldState *KeyResourceModel, data *KeyResourceModel) {
	if oldState.CredentialsA != nil && oldState.CredentialsA.Version.ValueInt64() != 0 && data.CredentialsA != nil {
		data.CredentialsA.Version = oldState.CredentialsA.Version
		data.CredentialsA.VersionGUID = oldState.CredentialsA.VersionGUID
	}
	if oldState.CredentialsB != nil && oldState.CredentialsB.Version.ValueInt64() != 0 && data.CredentialsB != nil {
		data.CredentialsB.Version = oldState.CredentialsB.Version
		data.CredentialsB.VersionGUID = oldState.CredentialsB.VersionGUID
	}
}

func keyVersionRequiresCreation(oldState *KeyResourceModel, data *KeyResourceModel) (bool, bool) {
	var createCredA, createCredB bool
	if oldState.CredentialsA == nil && data.CredentialsA != nil {
		createCredA = true
	}
	if oldState.CredentialsB == nil && data.CredentialsB != nil {
		createCredB = true
	}
	return createCredA, createCredB
}

func keyVersionRequiresDeletion(oldState *KeyResourceModel, data *KeyResourceModel) (bool, bool) {
	var deleteCredA, deleteCredB bool
	if oldState.CredentialsA != nil && (data.CredentialsA == nil || (oldState.CredentialsA.CloudAccessKeyID != data.CredentialsA.CloudAccessKeyID)) {
		deleteCredA = true
	}
	if oldState.CredentialsB != nil && (data.CredentialsB == nil || (oldState.CredentialsB.CloudAccessKeyID != data.CredentialsB.CloudAccessKeyID)) {
		deleteCredB = true
	}
	return deleteCredA, deleteCredB
}
func (r *KeyResource) updateAccessKey(ctx context.Context, plan *KeyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	client := Client(r.meta)
	resp, err := client.UpdateAccessKey(ctx, plan.buildUpdateRequest(), plan.buildFetchRequest())
	if err != nil {
		diags.AddError("update access key failed", err.Error())
		return diags
	}
	plan.AccessKeyName = types.StringValue(resp.AccessKeyName)
	plan.AccessKeyUID = types.Int64Value(resp.AccessKeyUID)

	return diags
}

func isVersionAssignedToProperty(ctx context.Context, client cloudaccess.CloudAccess, accessKeyUID int64, version int64) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	properties, err := client.LookupProperties(ctx, cloudaccess.LookupPropertiesRequest{
		AccessKeyUID: accessKeyUID,
		Version:      version,
	})
	if err != nil {
		diags.AddError("lookup properties failed ", err.Error())
		// As list of properties cannot be fetched this action should be blocked
		return false, diags
	}
	if len(properties.Properties) > 0 {
		return true, diags
	}
	return false, diags
}

// Delete implements resource.Resource.
func (r *KeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Access Key Resource")
	var oldState *KeyResourceModel
	client := Client(r.meta)

	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteTimeout, diags := oldState.Timeouts.Delete(ctx, deleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()
	versions, err := client.ListAccessKeyVersions(ctx, oldState.buildListKeyVersionsRequest())
	if err != nil {
		resp.Diagnostics.AddError("list access key versions failed", err.Error())
		return
	}
	for _, version := range versions.AccessKeyVersions {
		hasProperty, diags := isVersionAssignedToProperty(ctx, client, oldState.AccessKeyUID.ValueInt64(), version.Version)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		if hasProperty {
			resp.Diagnostics.AddError("version assigned to property error", fmt.Sprintf(assignedToPropertyError, version.Version, oldState.AccessKeyUID.ValueInt64()))
			return
		}
	}
	for _, version := range versions.AccessKeyVersions {
		versionToDelete := version.Version
		diags := r.deleteKeyVersion(ctx, oldState, versionToDelete, diags)
		if diags != nil {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if err = client.DeleteAccessKey(ctx, cloudaccess.AccessKeyRequest{
		AccessKeyUID: oldState.AccessKeyUID.ValueInt64(),
	}); err != nil {
		resp.Diagnostics.AddError("delete access key failed", err.Error())
		return
	}

	resp.Diagnostics.Append(r.waitForDelete(ctx, oldState.AccessKeyUID.ValueInt64())...)
}

func (r *KeyResource) deleteKeyVersion(ctx context.Context, oldState *KeyResourceModel, versionToDelete int64, diags diag.Diagnostics) diag.Diagnostics {
	client := Client(r.meta)
	_, err := client.DeleteAccessKeyVersion(ctx, oldState.buildDeleteKeyVersionRequest(versionToDelete))
	if err != nil {
		diags.AddError(fmt.Sprintf("delete access key version %d failed", versionToDelete), err.Error())
		return diags
	}
	isPending, diags := r.isPendingDelete(ctx, oldState.AccessKeyUID.ValueInt64(), versionToDelete)
	if diags.HasError() {
		return diags
	}
	if isPending {
		successfulDelete, diags := r.waitForVersionDelete(ctx, oldState.AccessKeyUID.ValueInt64(), versionToDelete)
		if !successfulDelete {
			return diags
		}

	}
	return diags
}

func (r *KeyResource) isPendingDelete(ctx context.Context, accessKeyUID int64, version int64) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := Client(r.meta)

	resp, err := client.GetAccessKeyVersion(ctx, cloudaccess.GetAccessKeyVersionRequest{
		AccessKeyUID: accessKeyUID,
		Version:      version,
	})
	if err != nil {
		diags.AddError(fmt.Sprintf("get access key version %d failed", version), err.Error())
		return false, diags
	}

	return resp.DeploymentStatus == cloudaccess.PendingDeletion, diags
}

func (r *KeyResource) waitForDelete(ctx context.Context, accessKeyUID int64) diag.Diagnostics {
	var diags diag.Diagnostics
	client := Client(r.meta)
	for {
		keys, err := client.ListAccessKeys(ctx, cloudaccess.ListAccessKeysRequest{})
		if err != nil {
			diags.AddError("list access keys failed", err.Error())
			return diags
		}
		keyDeleted := true
		for _, key := range keys.AccessKeys {
			if accessKeyUID == key.AccessKeyUID {
				keyDeleted = false
				break
			}
		}
		if keyDeleted {
			return diags
		}

		select {
		case <-time.Tick(pollingInterval):
			continue
		case <-ctx.Done():
			diags.AddError("deletion terminated",
				"context terminated the wait for deletion to finish")
			return diags
		}
	}
}

func (r *KeyResource) waitForVersionDelete(ctx context.Context, accessKeyUID int64, version int64) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := Client(r.meta)

	for {
		versions, err := client.ListAccessKeyVersions(ctx, cloudaccess.ListAccessKeyVersionsRequest{
			AccessKeyUID: accessKeyUID,
		})
		if err != nil {
			diags.AddError("list access key versions failed", err.Error())
			return false, diags
		}

		keyVersionDeleted := true
		for _, keyVersion := range versions.AccessKeyVersions {
			if version == keyVersion.Version {
				keyVersionDeleted = false
				break
			}
		}

		if keyVersionDeleted {
			return true, diags
		}

		select {
		case <-time.Tick(pollingInterval):
			continue
		case <-ctx.Done():
			diags.AddError("deletion terminated",
				ctx.Err().Error())
			return false, diags
		}
	}
}

func (m *KeyResourceModel) populateModelFromAccessKey(response *cloudaccess.GetAccessKeyResponse) {
	m.AccessKeyName = types.StringValue(response.AccessKeyName)
	m.AccessKeyUID = types.Int64Value(response.AccessKeyUID)
}

func (m *KeyResourceModel) importModelFromAccessKey(response *cloudaccess.GetAccessKeyResponse, inputGroupID int64, inputContractID string) error {
	m.AccessKeyName = types.StringValue(response.AccessKeyName)
	m.AccessKeyUID = types.Int64Value(response.AccessKeyUID)
	m.AuthenticationMethod = types.StringValue(response.AuthenticationMethod)

	m.NetworkConfig = NetworkConfig{
		SecurityNetwork: types.StringValue(string(response.NetworkConfiguration.SecurityNetwork)),
	}
	if response.NetworkConfiguration.AdditionalCDN != nil {
		m.NetworkConfig.AdditionalCDN = types.StringValue(string(*response.NetworkConfiguration.AdditionalCDN))
	}
	m.PrimaryGUID = types.StringValue("")

	if inputGroupID != 0 && inputContractID != "" {
		if err := m.findAndSetGroupAndContract(response.Groups, inputGroupID, inputContractID); err != nil {
			return err
		}
	} else {
		m.setDefaultGroupAndContract(response.Groups)
	}

	return nil
}

func (m *KeyResourceModel) findAndSetGroupAndContract(groups []cloudaccess.Group, inputGroupID int64, inputContractID string) error {
	for _, group := range groups {
		if group.GroupID == inputGroupID {
			for _, contractID := range group.ContractIDs {
				if contractID == inputContractID {
					m.GroupID = types.Int64Value(group.GroupID)
					m.ContractID = types.StringValue(contractID)
					return nil
				}
			}
			return fmt.Errorf("contractID %s not found in groupID %d", inputContractID, inputGroupID)
		}
	}
	return fmt.Errorf("groupID %d not found", inputGroupID)
}

func (m *KeyResourceModel) setDefaultGroupAndContract(groups []cloudaccess.Group) {
	if len(groups) == 0 {
		return
	}
	slices.SortFunc(groups, func(a, b cloudaccess.Group) int {
		return cmp.Compare(a.GroupID, b.GroupID)
	})
	lastGroup := groups[len(groups)-1]
	m.GroupID = types.Int64Value(lastGroup.GroupID)
	if len(lastGroup.ContractIDs) > 0 {
		m.ContractID = types.StringValue(lastGroup.ContractIDs[0])
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *KeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Access key Resource")

	// User-supplied import ID is a comma-separated list of accessKeyUID[,groupID[,contractID]]
	// contractID and groupID are optional as long as the accessKeyUID is sufficient to fetch the access key.
	var accessKeyUID, contractID string
	var groupID int64
	var err error
	parts := strings.Split(req.ID, ",")

	switch len(parts) {
	case 3:
		// All 3 parameters are present and valid
		accessKeyUID = parts[0]

		// Parse groupID safely
		groupID, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			resp.Diagnostics.AddError("Incorrect groupID", fmt.Sprintf("Couldn't parse provided groupID, \"%s\" is invalid", parts[1]))
			return
		}
		if groupID <= 0 {
			// Check if group ID is less than or equal to 0
			resp.Diagnostics.AddError("Invalid group ID", "group ID must be greater than 0")
		}
		// Validate contractID
		contractID = parts[2]
		if contractID == "" {
			resp.Diagnostics.AddError("Invalid contractID", "contractID cannot be empty")
			return
		}
	case 2:
		// contractID is absent but groupID is given
		resp.Diagnostics.AddError(
			"Incomplete Access Key Identifier",
			fmt.Sprintf("The identifier '%s' for Access Key '%s' is incomplete. Please provide both a valid Group ID and its corresponding Contract ID.", req.ID, accessKeyUID),
		)
		return
	case 1:
		// 1 parameter is valid
		accessKeyUID = parts[0]
	default:
		// Handle invalid cases with an error
		resp.Diagnostics.AddError(
			"Invalid Access Key Identifier",
			fmt.Sprintf("Unexpected number of parts in Access Key Identifier: %q", req.ID),
		)
		return
	}

	accessKeyID, err := strconv.ParseInt(accessKeyUID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Incorrect ID", fmt.Sprintf("Couldn't parse provided ID, \"%s\" is invalid", accessKeyUID))
		return
	}

	var data = &KeyResourceModel{}
	client := Client(r.meta)
	r.meta.OperationID()
	result, err := client.GetAccessKey(ctx, cloudaccess.AccessKeyRequest{
		AccessKeyUID: accessKeyID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Cannot Find Access key", err.Error())
		return
	}

	err = data.importModelFromAccessKey(result, groupID, contractID)
	if err != nil {
		resp.Diagnostics.AddError("Cannot Find Access key for a given groupID and contractID", err.Error())
		return
	}

	data.Timeouts = timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"delete": types.StringType,
			"create": types.StringType,
			"update": types.StringType,
		}),
	}

	versions, err := client.ListAccessKeyVersions(ctx, cloudaccess.ListAccessKeyVersionsRequest{AccessKeyUID: data.AccessKeyUID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError("Reading Access Key list Failed", err.Error())
		return
	}
	if len(versions.AccessKeyVersions) > 1 {
		if *versions.AccessKeyVersions[0].CloudAccessKeyID == *versions.AccessKeyVersions[1].CloudAccessKeyID {
			resp.Diagnostics.AddError("cloud access key id of access key error", "'cloud_access_key_id' should be unique for each pair of credentials")
			return
		}
	}

	data.populateModelFromVersionsList(versions)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *KeyResourceModel) populateModelFromVersionsList(versions *cloudaccess.ListAccessKeyVersionsResponse) diag.Diagnostics {
	var diags diag.Diagnostics
	credAFromState := false
	credBFromState := false
	// sort in ascending order
	// it will allow firstly to check older versions from state and than process later versions which are related to drift
	slices.SortFunc(versions.AccessKeyVersions, func(a, b cloudaccess.AccessKeyVersion) int {
		return cmp.Compare(a.Version, b.Version)
	})
	for _, version := range versions.AccessKeyVersions {
		if m.CredentialsA != nil && version.Version == m.CredentialsA.Version.ValueInt64() {
			m.CredentialsA.CloudAccessKeyID = types.StringValue(*version.CloudAccessKeyID)
			credAFromState = true
			continue
		}
		if m.CredentialsB != nil && version.Version == m.CredentialsB.Version.ValueInt64() {
			m.CredentialsB.CloudAccessKeyID = types.StringValue(*version.CloudAccessKeyID)
			credBFromState = true
			continue
		}
		//This part of loop is reached when on server exist version which is not present in state, so we encounter drift
		//It should be assigned to first empty Credential pair in incremental order
		if !credAFromState {
			m.CredentialsA = &Credentials{}
			m.CredentialsA.CloudAccessKeyID = types.StringValue(*version.CloudAccessKeyID)
			// Cannot retrieve secret form server
			m.CredentialsA.CloudSecretAccessKey = types.StringValue("")
			m.CredentialsA.Version = types.Int64Value(version.Version)
			m.CredentialsA.VersionGUID = types.StringValue(version.VersionGUID)
			m.CredentialsA.PrimaryKey = types.BoolValue(false)
			credAFromState = true
			continue
		}
		if !credBFromState {
			m.CredentialsB = &Credentials{}
			m.CredentialsB.CloudAccessKeyID = types.StringValue(*version.CloudAccessKeyID)
			// Cannot retrieve secret form server
			m.CredentialsB.CloudSecretAccessKey = types.StringValue("")
			m.CredentialsB.Version = types.Int64Value(version.Version)
			m.CredentialsB.VersionGUID = types.StringValue(version.VersionGUID)
			m.CredentialsB.PrimaryKey = types.BoolValue(false)
			credBFromState = true
			continue
		}
	}
	return diags
}

func (m *KeyResourceModel) buildCreateKeyRequest(useCredA bool) cloudaccess.CreateAccessKeyRequest {
	request := cloudaccess.CreateAccessKeyRequest{
		AccessKeyName:        m.AccessKeyName.ValueString(),
		AuthenticationMethod: m.AuthenticationMethod.ValueString(),
		ContractID:           m.ContractID.ValueString(),
		GroupID:              m.GroupID.ValueInt64(),
		Credentials:          m.setCredentialsForAccessKeyCreation(useCredA),
		NetworkConfiguration: cloudaccess.SecureNetwork{
			SecurityNetwork: cloudaccess.NetworkType(m.NetworkConfig.SecurityNetwork.ValueString()),
		},
	}
	if m.NetworkConfig.AdditionalCDN.ValueString() != "" {
		request.NetworkConfiguration.AdditionalCDN = ptr.To(cloudaccess.CDNType(m.NetworkConfig.AdditionalCDN.ValueString()))
	}
	return request
}

func (m *KeyResourceModel) setCredentialsForAccessKeyCreation(useCredA bool) cloudaccess.Credentials {
	if useCredA {
		return cloudaccess.Credentials{
			CloudSecretAccessKey: m.CredentialsA.CloudSecretAccessKey.ValueString(),
			CloudAccessKeyID:     m.CredentialsA.CloudAccessKeyID.ValueString(),
		}
	}
	return cloudaccess.Credentials{
		CloudSecretAccessKey: m.CredentialsB.CloudSecretAccessKey.ValueString(),
		CloudAccessKeyID:     m.CredentialsB.CloudAccessKeyID.ValueString(),
	}
}

func (m *KeyResourceModel) buildCreateKeyVersionRequest(useCredA bool) cloudaccess.CreateAccessKeyVersionRequest {
	var bodyParams cloudaccess.CreateAccessKeyVersionRequestBody
	if useCredA {
		bodyParams = cloudaccess.CreateAccessKeyVersionRequestBody{
			CloudAccessKeyID:     m.CredentialsA.CloudAccessKeyID.ValueString(),
			CloudSecretAccessKey: m.CredentialsA.CloudSecretAccessKey.ValueString(),
		}
	} else {
		bodyParams = cloudaccess.CreateAccessKeyVersionRequestBody{
			CloudAccessKeyID:     m.CredentialsB.CloudAccessKeyID.ValueString(),
			CloudSecretAccessKey: m.CredentialsB.CloudSecretAccessKey.ValueString(),
		}
	}
	return cloudaccess.CreateAccessKeyVersionRequest{
		AccessKeyUID: m.AccessKeyUID.ValueInt64(),
		Body:         bodyParams,
	}
}

func (m *KeyResourceModel) buildListKeyVersionsRequest() cloudaccess.ListAccessKeyVersionsRequest {
	return cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: m.AccessKeyUID.ValueInt64(),
	}
}

func (m *KeyResourceModel) buildDeleteKeyVersionRequest(version int64) cloudaccess.DeleteAccessKeyVersionRequest {
	return cloudaccess.DeleteAccessKeyVersionRequest{
		AccessKeyUID: m.AccessKeyUID.ValueInt64(),
		Version:      version,
	}
}

func (m *KeyResourceModel) buildUpdateRequest() cloudaccess.UpdateAccessKeyRequest {
	return cloudaccess.UpdateAccessKeyRequest{
		AccessKeyName: m.AccessKeyName.ValueString(),
	}
}

func (m *KeyResourceModel) buildFetchRequest() cloudaccess.AccessKeyRequest {
	return cloudaccess.AccessKeyRequest{
		AccessKeyUID: m.AccessKeyUID.ValueInt64(),
	}
}

func (r *KeyResource) waitUntilActivationCompleted(ctx context.Context, requestID int64, statusTimeout int64, plan *KeyResourceModel, credA bool) (*KeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := Client(r.meta)
	time.Sleep(time.Duration(statusTimeout) * time.Millisecond)
	for {
		statusResp, err := client.GetAccessKeyStatus(ctx, cloudaccess.GetAccessKeyStatusRequest{RequestID: requestID})
		if err != nil {
			diags.AddError(readError, err.Error())
			return nil, diags
		}
		if statusResp.ProcessingStatus == cloudaccess.ProcessingDone {
			plan.AccessKeyUID = types.Int64Value(statusResp.AccessKey.AccessKeyUID)
			versionResp, err := client.GetAccessKeyVersion(ctx, cloudaccess.GetAccessKeyVersionRequest{
				AccessKeyUID: statusResp.AccessKey.AccessKeyUID,
				Version:      statusResp.AccessKeyVersion.Version,
			})
			if err != nil {
				diags.AddError(readError, err.Error())
				return nil, diags
			}
			if versionResp.DeploymentStatus == cloudaccess.Active {
				if credA {
					plan.CredentialsA.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					plan.CredentialsA.VersionGUID = types.StringValue(versionResp.VersionGUID)
				} else {
					plan.CredentialsB.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					plan.CredentialsB.VersionGUID = types.StringValue(versionResp.VersionGUID)
				}
				return plan, diags
			}
		}
		select {
		case <-time.After(pollingInterval):
			continue
		case <-ctx.Done():
			diags.AddError("reached activation timeout", ctx.Err().Error())
			return nil, diags
		}
	}
}

func (r *KeyResource) waitUntilVersionCreatedCompleted(ctx context.Context, requestID int64, statusTimeout int64, plan *KeyResourceModel, credentialA bool) (*KeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := Client(r.meta)
	time.Sleep(time.Duration(statusTimeout) * time.Millisecond)

	for {
		statusResp, err := client.GetAccessKeyVersionStatus(ctx, cloudaccess.GetAccessKeyVersionStatusRequest{RequestID: requestID})
		if err != nil {
			diags.AddError(readError, err.Error())
			return nil, diags
		}
		if statusResp.ProcessingStatus == cloudaccess.ProcessingDone {
			versionResp, versionErr := client.GetAccessKeyVersion(ctx, cloudaccess.GetAccessKeyVersionRequest{
				AccessKeyUID: statusResp.AccessKeyVersion.AccessKeyUID,
				Version:      statusResp.AccessKeyVersion.Version,
			})
			if versionErr != nil {
				diags.AddError(readError, err.Error())
				return nil, diags
			}
			if versionResp.DeploymentStatus == cloudaccess.Active {
				if credentialA {
					plan.CredentialsA.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					plan.CredentialsA.VersionGUID = types.StringValue(versionResp.VersionGUID)
				} else {
					plan.CredentialsB.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					plan.CredentialsB.VersionGUID = types.StringValue(versionResp.VersionGUID)
				}
				return plan, diags
			}
		}

		select {
		case <-time.After(pollingInterval):
			continue
		case <-ctx.Done():
			diags.AddError("reached activation timeout", ctx.Err().Error())
			return nil, diags
		}
	}
}
