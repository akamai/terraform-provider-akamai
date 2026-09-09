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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &KeyResource{}
	_ resource.ResourceWithConfigure      = &KeyResource{}
	_ resource.ResourceWithModifyPlan     = &KeyResource{}
	_ resource.ResourceWithImportState    = &KeyResource{}
	_ resource.ResourceWithValidateConfig = &KeyResource{}
)

const (
	readError                = "could not read access key from API"
	creationKeyFailed        = "access key creation failed"
	creationKeyVersionFailed = "access key version creation failed"
	assignedToPropertyError  = "cannot delete version: %d of access key %d assigned to property"
	diagErrAccessKeyNotFound = "cannot find access key: %d"
)

// KeyResource represents akamai_cloudaccess_key resource
type KeyResource struct {
	meta.Resource
	activationTimeout time.Duration
	updateTimeout     time.Duration
	deleteTimeout     time.Duration
	pollingInterval   time.Duration
}

// KeyResourceModel represents model of akamai_cloudaccess_key resource
type KeyResourceModel struct {
	AccessKeyName        types.String   `tfsdk:"access_key_name"`
	AuthenticationMethod types.String   `tfsdk:"authentication_method"`
	ContractID           types.String   `tfsdk:"contract_id"`
	GroupID              types.Int64    `tfsdk:"group_id"`
	PrimaryGUID          types.String   `tfsdk:"primary_guid"`
	CredentialsA         types.Object   `tfsdk:"credentials_a"`
	CredentialsB         types.Object   `tfsdk:"credentials_b"`
	NetworkConfig        types.Object   `tfsdk:"network_configuration"`
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

func credentialType() map[string]attr.Type {
	return credentialSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

// credentialsA extracts the CredentialsA object into a typed model. Returns nil when null or unknown.
func (m *KeyResourceModel) credentialsA(ctx context.Context) (*Credentials, diag.Diagnostics) {
	if m.CredentialsA.IsNull() || m.CredentialsA.IsUnknown() {
		return nil, nil
	}
	var cred Credentials
	return &cred, m.CredentialsA.As(ctx, &cred, basetypes.ObjectAsOptions{})
}

// setCredentialsA stores cred into the CredentialsA object field.
func (m *KeyResourceModel) setCredentialsA(ctx context.Context, cred *Credentials) diag.Diagnostics {
	if cred == nil {
		m.CredentialsA = types.ObjectNull(credentialType())
		return nil
	}
	var dd diag.Diagnostics
	m.CredentialsA, dd = types.ObjectValueFrom(ctx, credentialType(), cred)
	return dd
}

// credentialsB extracts the CredentialsB object into a typed model. Returns nil when null or unknown.
func (m *KeyResourceModel) credentialsB(ctx context.Context) (*Credentials, diag.Diagnostics) {
	if m.CredentialsB.IsNull() || m.CredentialsB.IsUnknown() {
		return nil, nil
	}
	var cred Credentials
	return &cred, m.CredentialsB.As(ctx, &cred, basetypes.ObjectAsOptions{})
}

// setCredentialsB stores cred into the CredentialsB object field.
func (m *KeyResourceModel) setCredentialsB(ctx context.Context, cred *Credentials) diag.Diagnostics {
	if cred == nil {
		m.CredentialsB = types.ObjectNull(credentialType())
		return nil
	}
	var dd diag.Diagnostics
	m.CredentialsB, dd = types.ObjectValueFrom(ctx, credentialType(), cred)
	return dd
}

// credentials extracts both CredentialsA and CredentialsB into typed models.
func (m *KeyResourceModel) credentials(ctx context.Context) (*Credentials, *Credentials, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	a, diagnostics := m.credentialsA(ctx)
	diags.Append(diagnostics...)
	b, diagnostics := m.credentialsB(ctx)
	diags.Append(diagnostics...)
	return a, b, diags
}

func networkConfigType() map[string]attr.Type {
	return networkConfigurationSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

// networkConfigModel extracts the NetworkConfig object into a typed model. Returns nil when null or unknown.
func (m *KeyResourceModel) networkConfigModel(ctx context.Context) (*NetworkConfig, diag.Diagnostics) {
	if m.NetworkConfig.IsNull() || m.NetworkConfig.IsUnknown() {
		return nil, nil
	}
	var networkConfig NetworkConfig
	return &networkConfig, m.NetworkConfig.As(ctx, &networkConfig, basetypes.ObjectAsOptions{})
}

// setNetworkConfig stores networkConfig into the NetworkConfig object field.
func (m *KeyResourceModel) setNetworkConfig(ctx context.Context, networkConfig *NetworkConfig) diag.Diagnostics {
	if networkConfig == nil {
		m.NetworkConfig = types.ObjectNull(networkConfigType())
		return nil
	}
	var dd diag.Diagnostics
	m.NetworkConfig, dd = types.ObjectValueFrom(ctx, networkConfigType(), networkConfig)
	return dd
}

// NewKeyResource returns new cloudaccess key resource
func NewKeyResource() resource.Resource {
	return &KeyResource{
		activationTimeout: 60 * time.Minute,
		updateTimeout:     60 * time.Minute,
		deleteTimeout:     60 * time.Minute,
		pollingInterval:   1 * time.Minute,
	}
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r *KeyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config KeyResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	credA, credB, dd := config.credentials(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	networkConfig, dd := config.networkConfigModel(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}

	verifyCloudAccessKeyIDPresence(config.AuthenticationMethod, credA, credB, &resp.Diagnostics)
	verifyAdditionalCDNPresence(config.AuthenticationMethod, networkConfig, &resp.Diagnostics)
	verifyCloudAccessKeyIDAndSecretLength(ctx, config, &resp.Diagnostics)
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

	var planCredA, planCredB *Credentials
	if plan != nil {
		var dd diag.Diagnostics
		planCredA, planCredB, dd = plan.credentials(ctx)
		if response.Diagnostics.Append(dd...); response.Diagnostics.HasError() {
			return
		}
	}

	if state == nil && planCredA == nil && planCredB == nil {
		response.Diagnostics.AddError("at least one credentials are required for creation", "`credentials_a` or `credentials_b` must be specified")
		return
	}

	if plan != nil && planCredA != nil && planCredB != nil &&
		planCredA.PrimaryKey.ValueBool() && planCredB.PrimaryKey.ValueBool() {
		response.Diagnostics.AddError("primary version of access key error", "only one pair of access key version can have 'primary_key' set as 'true'")
		return
	}

	if plan != nil && !cloudAccessKeyIDUnique(planCredA, planCredB) {
		response.Diagnostics.AddError("cloud access key id of access key error", "'cloud_access_key_id' should be unique for each pair of credentials")
		return
	}

	if state != nil && plan != nil {
		stateCredA, stateCredB, dd := state.credentials(ctx)
		if response.Diagnostics.Append(dd...); response.Diagnostics.HasError() {
			return
		}
		if changedOrderOfCredentials(stateCredA, stateCredB, planCredA, planCredB) {
			response.Diagnostics.AddError("access key credentials error", "cannot change order of `credentials_a` and `credentials_b`")
			return
		}
		if checkIfSecretChangedAndWasNotEmpty(stateCredA, stateCredB, planCredA, planCredB) {
			response.Diagnostics.AddError("access key credentials error", "cannot update cloud access secret without update of cloud access key id, expect update of secret after import with no API calls")
			return
		}
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
				Required: true,
				MarkdownDescription: "The type of signing process used to authenticate API requests:\n" +
					"  - `AOS4_HMAC_SHA256` — Akamai Object Storage\n" +
					"  - `AVM_CLOUDINARY` — Akamai Video Manager Cloudinary\n" +
					"  - `AWS4_HMAC_SHA256` — Amazon Web Services\n" +
					"  - `GOOG4_HMAC_SHA256` — Google Cloud Services\n" +
					"  - `G2O` — Akamai Signature Header Authentication\n" +
					"  - `VP_QUEUE_IT` — Akamai Visitor Prioritization powered by Queue-it",
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(string(cloudaccess.AuthAOS), string(cloudaccess.AuthAVMCloudinary), string(cloudaccess.AuthAWS), string(cloudaccess.AuthGOOG), string(cloudaccess.AuthG2O), string(cloudaccess.AuthVPQueueIt)),
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
			"credentials_a":         credentialSchema(),
			"credentials_b":         credentialSchema(),
			"network_configuration": networkConfigurationSchema(),
			"access_key_uid": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique identifier Akamai assigns to an access key.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func networkConfigurationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
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
	}
}

func credentialSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:    true,
		Description: "The combination of a `cloud_access_key_id` and a `cloud_secret_access_key` used to sign API requests. This pair can be identified as access key version. Access key can contain only two access key versions at specific time (defined as credentialsA and credentialsB).",
		Attributes: map[string]schema.Attribute{
			"cloud_access_key_id": schema.StringAttribute{
				Description: "Access key id from cloud provider which is used to sign API requests",
				Optional:    true,
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
	}
}

// isTimeoutChanged defines if timeout changed between plan and state
func isTimeoutChanged(state, plan *KeyResourceModel) bool {
	return state != nil && plan != nil &&
		!state.Timeouts.Equal(plan.Timeouts)
}

// hasEqualContent reports whether all user-managed (non-timeout) fields are equal between state and plan.
func (m *KeyResourceModel) hasEqualContent(ctx context.Context, plan *KeyResourceModel) (bool, diag.Diagnostics) {
	if m.AccessKeyName != plan.AccessKeyName {
		return false, nil
	}
	stateCredA, stateCredB, dd := m.credentials(ctx)
	if dd.HasError() {
		return false, dd
	}
	planCredA, planCredB, dd := plan.credentials(ctx)
	if dd.HasError() {
		return false, dd
	}
	if !credentialsEqualContent(stateCredA, planCredA) {
		return false, nil
	}
	if !credentialsEqualContent(stateCredB, planCredB) {
		return false, nil
	}
	return true, nil
}

// credentialsEqualContent compares user-managed credential fields (excludes computed version/version_guid).
func credentialsEqualContent(a, b *Credentials) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.CloudAccessKeyID == b.CloudAccessKeyID &&
		a.CloudSecretAccessKey == b.CloudSecretAccessKey &&
		a.PrimaryKey == b.PrimaryKey
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
	createTimeout, diags := plan.Timeouts.Create(ctx, r.activationTimeout)
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

	planCredA, planCredB, dd := plan.credentials(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	if planCredA != nil && planCredB != nil {
		plan, diags = r.createVersion(ctx, plan, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(r.setupPrimaryGUID(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// setupPrimaryGUID calculates `primary_guid` based on `primary_key` and `version_guid` parameters
func (r *KeyResource) setupPrimaryGUID(ctx context.Context, state *KeyResourceModel) diag.Diagnostics {
	credA, credB, dd := state.credentials(ctx)
	if dd.HasError() {
		return dd
	}
	if credA != nil && credA.PrimaryKey.ValueBool() {
		state.PrimaryGUID = credA.VersionGUID
		return nil
	}
	if credB != nil && credB.PrimaryKey.ValueBool() {
		state.PrimaryGUID = credB.VersionGUID
		return nil
	}
	state.PrimaryGUID = types.StringValue("")
	return nil
}

func (r *KeyResource) create(ctx context.Context, plan *KeyResourceModel) (*KeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.Client.GetCloudAccess()
	credA, dd := plan.credentialsA(ctx)
	if diags.Append(dd...); diags.HasError() {
		return nil, diags
	}
	creationKeyWithCredA := credA != nil
	req, dd := plan.buildCreateKeyRequest(ctx, creationKeyWithCredA)
	if diags.Append(dd...); diags.HasError() {
		return nil, diags
	}
	resp, err := client.CreateAccessKey(ctx, req)
	if err != nil {
		diags.AddError("create access key failed", err.Error())
		return nil, diags
	}

	return r.waitUntilActivationCompleted(ctx, resp.RequestID, resp.RetryAfter, plan, creationKeyWithCredA)
}

func (r *KeyResource) createVersion(ctx context.Context, plan *KeyResourceModel, useCredentialA bool) (*KeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.Client.GetCloudAccess()

	req, dd := plan.buildCreateKeyVersionRequest(ctx, useCredentialA)
	if diags.Append(dd...); diags.HasError() {
		return nil, diags
	}
	resp, err := client.CreateAccessKeyVersion(ctx, req)
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

	resp.Diagnostics.Append(r.setupPrimaryGUID(ctx, oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
}

func (r *KeyResource) read(ctx context.Context, data *KeyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	client := r.Client.GetCloudAccess()

	result, err := client.GetAccessKey(ctx, cloudaccess.AccessKeyRequest{
		AccessKeyUID: data.AccessKeyUID.ValueInt64(),
	})
	if errors.Is(err, cloudaccess.ErrAccessKeyNotFound) {
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
	return data.populateModelFromVersionsList(ctx, versions)
}

// Update implements resource.Resource.
func (r *KeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Access Key Resource")
	var diags diag.Diagnostics
	var plan *KeyResourceModel
	client := r.Client.GetCloudAccess()
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, r.updateTimeout)
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

	// If only timeouts changed, update state locally — no API calls needed.
	equalContent, dd := oldState.hasEqualContent(ctx, plan)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	if equalContent && isTimeoutChanged(oldState, plan) {
		oldState.Timeouts = plan.Timeouts
		resp.Diagnostics.AddWarning("Local update only",
			"Only the timeout settings have changed, no API call will be made.")
		resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
		return
	}

	resp.Diagnostics.Append(initializeCredentialVersions(ctx, oldState, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.AccessKeyUID = oldState.AccessKeyUID
	if oldState.AccessKeyName != plan.AccessKeyName {
		resp.Diagnostics.Append(r.updateAccessKey(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	stateCredA, stateCredB, dd := oldState.credentials(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	planCredA, planCredB, dd := plan.credentials(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}

	updateInStateCredentialA, updateInStateCredentialB := keyVersionRequiresUpdateInState(stateCredA, stateCredB, planCredA, planCredB)
	if updateInStateCredentialA {
		stateCredA.CloudSecretAccessKey = planCredA.CloudSecretAccessKey
		resp.Diagnostics.Append(oldState.setCredentialsA(ctx, stateCredA)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if updateInStateCredentialB {
		stateCredB.CloudSecretAccessKey = planCredB.CloudSecretAccessKey
		resp.Diagnostics.Append(oldState.setCredentialsB(ctx, stateCredB)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if updateInStateCredentialA || updateInStateCredentialB {
		resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
	}

	deleteCredentialsA, deleteCredentialsB := keyVersionRequiresDeletion(stateCredA, stateCredB, planCredA, planCredB)
	if deleteCredentialsA || deleteCredentialsB {
		diags = r.deleteVersion(ctx, oldState, stateCredA, stateCredB, client, resp, diags, deleteCredentialsA, deleteCredentialsB)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if deleteCredentialsA {
			stateCredA = nil
			resp.Diagnostics.Append(oldState.setCredentialsA(ctx, stateCredA)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		if deleteCredentialsB {
			stateCredB = nil
			resp.Diagnostics.Append(oldState.setCredentialsB(ctx, stateCredB)...)
			if resp.Diagnostics.HasError() {
				return
			}

		}
	}

	createCredentialsA, createCredentialsB := keyVersionRequiresCreation(stateCredA, stateCredB, planCredA, planCredB)
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

	resp.Diagnostics.Append(r.setupPrimaryGUID(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyResource) deleteVersion(ctx context.Context, oldState *KeyResourceModel, oldStateCredA, oldStateCredB *Credentials, client cloudaccess.CloudAccess, resp *resource.UpdateResponse, diags diag.Diagnostics, deleteCredentialsA, deleteCredentialsB bool) diag.Diagnostics {
	var versionsToDelete []int64
	if deleteCredentialsA {
		versionToDelete := oldStateCredA.Version.ValueInt64()
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
		versionToDelete := oldStateCredB.Version.ValueInt64()
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

// verifyAdditionalCDNPresence checks that AdditionalCDN is not specified for VP_QUEUE_IT and
// AVM_CLOUDINARY authentication methods.
func verifyAdditionalCDNPresence(authenticationMethod types.String, networkConfig *NetworkConfig, diags *diag.Diagnostics) {
	if authenticationMethod.IsUnknown() {
		return
	}
	authMethod := authenticationMethod.ValueString()

	if networkConfig == nil {
		return
	}
	if authMethod != string(cloudaccess.AuthVPQueueIt) &&
		authMethod != string(cloudaccess.AuthAVMCloudinary) {
		return
	}
	if !networkConfig.AdditionalCDN.IsNull() {
		diags.AddAttributeError(
			path.Root("network_configuration").AtName("additional_cdn"),
			"additional cdn not allowed error",
			"for the selected authentication method `additional_cdn` must not be specified",
		)
	}
}

// verifyCloudAccessKeyIDAndSecretLength checks length of CloudAccessKeyID and CloudSecretAccessKey for G2O authentication method.
func verifyCloudAccessKeyIDAndSecretLength(ctx context.Context, data KeyResourceModel, diags *diag.Diagnostics) {
	if data.AuthenticationMethod.IsUnknown() {
		return
	}
	if data.AuthenticationMethod.ValueString() != string(cloudaccess.AuthG2O) {
		return
	}
	credA, credB, dd := data.credentials(ctx)
	if dd.HasError() {
		diags.Append(dd...)
		return
	}
	validateG2OCredentials(credA, "credentials_a", diags)
	validateG2OCredentials(credB, "credentials_b", diags)
}

// validateG2OCredentials checks that CloudAccessKeyID is between 1 and 8 characters long and CloudSecretAccessKey is between 32 and 64 characters long for G2O authentication method.
func validateG2OCredentials(creds *Credentials, credName string, diags *diag.Diagnostics) {
	if creds == nil || creds.CloudAccessKeyID.IsUnknown() || creds.CloudSecretAccessKey.IsUnknown() {
		return
	}
	if len(creds.CloudAccessKeyID.ValueString()) < cloudaccess.G2OAccessKeyIDMinLength || len(creds.CloudAccessKeyID.ValueString()) > cloudaccess.G2OAccessKeyIDMaxLength {
		diags.AddAttributeError(
			path.Root(credName).AtName("cloud_access_key_id"),
			"cloud access key id value error",
			fmt.Sprintf("for the selected authentication method `cloud_access_key_id` in `%s` should be between %d and %d characters long", credName, cloudaccess.G2OAccessKeyIDMinLength, cloudaccess.G2OAccessKeyIDMaxLength),
		)
	}
	if len(creds.CloudSecretAccessKey.ValueString()) < cloudaccess.G2OSecretAccessKeyMinLength || len(creds.CloudSecretAccessKey.ValueString()) > cloudaccess.G2OSecretAccessKeyMaxLength {
		diags.AddAttributeError(
			path.Root(credName).AtName("cloud_secret_access_key"),
			"cloud secret access key value error",
			fmt.Sprintf("for the selected authentication method `cloud_secret_access_key` in `%s` should be between %d and %d characters long", credName, cloudaccess.G2OSecretAccessKeyMinLength, cloudaccess.G2OSecretAccessKeyMaxLength),
		)
	}
}

// verifyCloudAccessKeyIDPresence checks that CloudAccessKeyID is present for authentication
// methods other than VP_QUEUE_IT and AVM_CLOUDINARY.
func verifyCloudAccessKeyIDPresence(authenticationMethod types.String, credA, credB *Credentials, diags *diag.Diagnostics) {
	if authenticationMethod.IsUnknown() {
		return
	}
	authMethod := authenticationMethod.ValueString()

	if authMethod == string(cloudaccess.AuthVPQueueIt) ||
		authMethod == string(cloudaccess.AuthAVMCloudinary) {
		return
	}
	if credA != nil && credA.CloudAccessKeyID.IsNull() {
		diags.AddAttributeError(
			path.Root("credentials_a").AtName("cloud_access_key_id"),
			"cloud access key id missing error",
			"for the selected authentication method `cloud_access_key_id` in `credentials_a` cannot be empty",
		)
	}
	if credB != nil && credB.CloudAccessKeyID.IsNull() {
		diags.AddAttributeError(
			path.Root("credentials_b").AtName("cloud_access_key_id"),
			"cloud access key id missing error",
			"for the selected authentication method `cloud_access_key_id` in `credentials_b` cannot be empty",
		)
	}
}

func cloudAccessKeyIDUnique(credA, credB *Credentials) bool {
	if credA == nil || credB == nil {
		return true
	}
	if credA.CloudAccessKeyID.IsNull() || credB.CloudAccessKeyID.IsNull() {
		return true
	}
	return !credA.CloudAccessKeyID.Equal(credB.CloudAccessKeyID)
}

func changedOrderOfCredentials(stateCredA, stateCredB, planCredA, planCredB *Credentials) bool {
	anyCredentialNil := stateCredA == nil || planCredA == nil ||
		stateCredB == nil || planCredB == nil

	if anyCredentialNil {
		return false
	}

	keyIDsSwapped := stateCredA.CloudAccessKeyID.Equal(planCredB.CloudAccessKeyID) &&
		stateCredB.CloudAccessKeyID.Equal(planCredA.CloudAccessKeyID)

	// When both CloudAccessKeyIDs are present, we have the legacy behaviour - we check if they are swapped
	if !planCredA.CloudAccessKeyID.IsNull() && !planCredB.CloudAccessKeyID.IsNull() {
		return keyIDsSwapped
	}

	secretsSwapped := stateCredA.CloudSecretAccessKey.Equal(planCredB.CloudSecretAccessKey) &&
		stateCredB.CloudSecretAccessKey.Equal(planCredA.CloudSecretAccessKey)

	// When exactly one CloudAccessKeyID is present, use both key ID and secrets to confirm the swap
	if !planCredA.CloudAccessKeyID.IsNull() || !planCredB.CloudAccessKeyID.IsNull() {
		return keyIDsSwapped && secretsSwapped
	}

	// When both plan key IDs are null, additionally require secrets to be distinct to avoid false positives
	return keyIDsSwapped && secretsSwapped &&
		!planCredA.CloudSecretAccessKey.Equal(planCredB.CloudSecretAccessKey)
}

func checkIfSecretChangedAndWasNotEmpty(stateCredA, stateCredB, planCredA, planCredB *Credentials) bool {
	if stateCredA != nil && planCredA != nil &&
		!planCredA.CloudAccessKeyID.IsNull() &&
		stateCredA.CloudAccessKeyID.ValueString() == planCredA.CloudAccessKeyID.ValueString() &&
		stateCredA.CloudSecretAccessKey.ValueString() != "" && stateCredA.CloudSecretAccessKey.ValueString() != planCredA.CloudSecretAccessKey.ValueString() {
		return true
	}
	if stateCredB != nil && planCredB != nil &&
		!planCredB.CloudAccessKeyID.IsNull() &&
		stateCredB.CloudAccessKeyID.ValueString() == planCredB.CloudAccessKeyID.ValueString() &&
		stateCredB.CloudSecretAccessKey.ValueString() != "" && stateCredB.CloudSecretAccessKey.ValueString() != planCredB.CloudSecretAccessKey.ValueString() {
		return true
	}
	return false
}

// keyVersionRequiresUpdateInState reports whether a credential's secret needs to be written to state
// without making an API call. This covers the post-import scenario: after terraform import, secrets
// are absent from state because the API never returns them.
func keyVersionRequiresUpdateInState(stateCredA, stateCredB, planCredA, planCredB *Credentials) (bool, bool) {
	var updateCredA, updateCredB bool
	if stateCredA != nil && planCredA != nil && stateCredA.CloudAccessKeyID == planCredA.CloudAccessKeyID && stateCredA.CloudSecretAccessKey.ValueString() == "" {
		updateCredA = true
	}
	if stateCredB != nil && planCredB != nil && stateCredB.CloudAccessKeyID == planCredB.CloudAccessKeyID && stateCredB.CloudSecretAccessKey.ValueString() == "" {
		updateCredB = true
	}
	return updateCredA, updateCredB
}

func initializeCredentialVersions(ctx context.Context, oldState, data *KeyResourceModel) diag.Diagnostics {
	stateCredA, dd := oldState.credentialsA(ctx)
	if dd.HasError() {
		return dd
	}
	planCredA, dd := data.credentialsA(ctx)
	if dd.HasError() {
		return dd
	}
	if stateCredA != nil && stateCredA.Version.ValueInt64() != 0 && planCredA != nil {
		planCredA.Version = stateCredA.Version
		planCredA.VersionGUID = stateCredA.VersionGUID
		if dd = data.setCredentialsA(ctx, planCredA); dd.HasError() {
			return dd
		}
	}

	stateCredB, dd := oldState.credentialsB(ctx)
	if dd.HasError() {
		return dd
	}
	planCredB, dd := data.credentialsB(ctx)
	if dd.HasError() {
		return dd
	}
	if stateCredB != nil && stateCredB.Version.ValueInt64() != 0 && planCredB != nil {
		planCredB.Version = stateCredB.Version
		planCredB.VersionGUID = stateCredB.VersionGUID
		if dd = data.setCredentialsB(ctx, planCredB); dd.HasError() {
			return dd
		}
	}
	return nil
}

func keyVersionRequiresCreation(stateCredA, stateCredB, planCredA, planCredB *Credentials) (bool, bool) {
	return stateCredA == nil && planCredA != nil, stateCredB == nil && planCredB != nil
}

func keyVersionRequiresDeletion(stateCredA, stateCredB, planCredA, planCredB *Credentials) (bool, bool) {
	var deleteCredA, deleteCredB bool
	if stateCredA != nil &&
		(planCredA == nil ||
			(!stateCredA.CloudAccessKeyID.Equal(planCredA.CloudAccessKeyID)) ||
			(planCredA.CloudAccessKeyID.IsNull() &&
				!stateCredA.CloudSecretAccessKey.Equal(planCredA.CloudSecretAccessKey))) {
		deleteCredA = true
	}
	if stateCredB != nil &&
		(planCredB == nil ||
			(!stateCredB.CloudAccessKeyID.Equal(planCredB.CloudAccessKeyID)) ||
			(planCredB.CloudAccessKeyID.IsNull() &&
				!stateCredB.CloudSecretAccessKey.Equal(planCredB.CloudSecretAccessKey))) {
		deleteCredB = true
	}
	return deleteCredA, deleteCredB
}
func (r *KeyResource) updateAccessKey(ctx context.Context, plan *KeyResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	client := r.Client.GetCloudAccess()
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
	client := r.Client.GetCloudAccess()

	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteTimeout, diags := oldState.Timeouts.Delete(ctx, r.deleteTimeout)
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
	client := r.Client.GetCloudAccess()
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
	client := r.Client.GetCloudAccess()

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
	client := r.Client.GetCloudAccess()
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
		case <-time.After(r.pollingInterval):
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
	client := r.Client.GetCloudAccess()

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
		case <-time.After(r.pollingInterval):
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

func (m *KeyResourceModel) importModelFromAccessKey(ctx context.Context, response *cloudaccess.GetAccessKeyResponse, inputGroupID int64, inputContractID string) diag.Diagnostics {
	m.AccessKeyName = types.StringValue(response.AccessKeyName)
	m.AccessKeyUID = types.Int64Value(response.AccessKeyUID)
	m.AuthenticationMethod = types.StringValue(response.AuthenticationMethod)

	networkConfig := &NetworkConfig{
		SecurityNetwork: types.StringValue(string(response.NetworkConfiguration.SecurityNetwork)),
	}
	if response.NetworkConfiguration.AdditionalCDN != nil {
		networkConfig.AdditionalCDN = types.StringValue(string(*response.NetworkConfiguration.AdditionalCDN))
	}
	if dd := m.setNetworkConfig(ctx, networkConfig); dd.HasError() {
		return dd
	}
	m.PrimaryGUID = types.StringValue("")
	m.CredentialsA = types.ObjectNull(credentialType())
	m.CredentialsB = types.ObjectNull(credentialType())

	if inputGroupID != 0 && inputContractID != "" {
		if err := m.findAndSetGroupAndContract(response.Groups, inputGroupID, inputContractID); err != nil {
			var dd diag.Diagnostics
			dd.AddError("Cannot Find Access key for a given groupID and contractID", err.Error())
			return dd
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
	client := r.Client.GetCloudAccess()
	result, err := client.GetAccessKey(ctx, cloudaccess.AccessKeyRequest{
		AccessKeyUID: accessKeyID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Cannot Find Access key", err.Error())
		return
	}

	if resp.Diagnostics.Append(data.importModelFromAccessKey(ctx, result, groupID, contractID)...); resp.Diagnostics.HasError() {
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
		if versions.AccessKeyVersions[0].CloudAccessKeyID != nil && versions.AccessKeyVersions[1].CloudAccessKeyID != nil &&
			*versions.AccessKeyVersions[0].CloudAccessKeyID == *versions.AccessKeyVersions[1].CloudAccessKeyID {
			resp.Diagnostics.AddError("cloud access key id of access key error", "'cloud_access_key_id' should be unique for each pair of credentials")
			return
		}
	}

	resp.Diagnostics.Append(data.populateModelFromVersionsList(ctx, versions)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *KeyResourceModel) populateModelFromVersionsList(ctx context.Context, versions *cloudaccess.ListAccessKeyVersionsResponse) diag.Diagnostics {
	var diags diag.Diagnostics
	credAFromState := false
	credBFromState := false
	// sort in ascending order
	// it will allow firstly to check older versions from state and than process later versions which are related to drift
	slices.SortFunc(versions.AccessKeyVersions, func(a, b cloudaccess.AccessKeyVersion) int {
		return cmp.Compare(a.Version, b.Version)
	})

	credA, credB, dd := m.credentials(ctx)
	if diags.Append(dd...); diags.HasError() {
		return diags
	}

	for _, version := range versions.AccessKeyVersions {
		cloudAccessKeyID := types.StringNull()
		if version.CloudAccessKeyID != nil {
			cloudAccessKeyID = types.StringValue(*version.CloudAccessKeyID)
		}
		if credA != nil && version.Version == credA.Version.ValueInt64() {
			credA.CloudAccessKeyID = cloudAccessKeyID
			credAFromState = true
			continue
		}
		if credB != nil && version.Version == credB.Version.ValueInt64() {
			credB.CloudAccessKeyID = cloudAccessKeyID
			credBFromState = true
			continue
		}
		//This part of loop is reached when on server exist version which is not present in state, so we encounter drift
		//It should be assigned to first empty Credential pair in incremental order
		if !credAFromState {
			credA = &Credentials{
				CloudAccessKeyID: cloudAccessKeyID,
				// Cannot retrieve secret from server
				CloudSecretAccessKey: types.StringValue(""),
				Version:              types.Int64Value(version.Version),
				VersionGUID:          types.StringValue(version.VersionGUID),
				PrimaryKey:           types.BoolValue(false),
			}
			credAFromState = true
			continue
		}
		if !credBFromState {
			credB = &Credentials{
				CloudAccessKeyID: cloudAccessKeyID,
				// Cannot retrieve secret from server
				CloudSecretAccessKey: types.StringValue(""),
				Version:              types.Int64Value(version.Version),
				VersionGUID:          types.StringValue(version.VersionGUID),
				PrimaryKey:           types.BoolValue(false),
			}
			credBFromState = true
			continue
		}
	}
	if !credAFromState {
		credA = nil
	}
	if !credBFromState {
		credB = nil
	}
	diags.Append(m.setCredentialsA(ctx, credA)...)
	diags.Append(m.setCredentialsB(ctx, credB)...)
	return diags
}

func (m *KeyResourceModel) buildCreateKeyRequest(ctx context.Context, useCredA bool) (cloudaccess.CreateAccessKeyRequest, diag.Diagnostics) {
	var dd diag.Diagnostics
	networkConfig, dd := m.networkConfigModel(ctx)
	if dd.HasError() {
		return cloudaccess.CreateAccessKeyRequest{}, dd
	}
	creds, dd := m.credentialsForAccessKeyCreation(ctx, useCredA)
	if dd.HasError() {
		return cloudaccess.CreateAccessKeyRequest{}, dd
	}
	request := cloudaccess.CreateAccessKeyRequest{
		AccessKeyName:        m.AccessKeyName.ValueString(),
		AuthenticationMethod: m.AuthenticationMethod.ValueString(),
		ContractID:           m.ContractID.ValueString(),
		GroupID:              m.GroupID.ValueInt64(),
		Credentials:          creds,
		NetworkConfiguration: cloudaccess.SecureNetwork{
			SecurityNetwork: cloudaccess.NetworkType(networkConfig.SecurityNetwork.ValueString()),
		},
	}
	if networkConfig.AdditionalCDN.ValueString() != "" {
		request.NetworkConfiguration.AdditionalCDN = ptr.To(cloudaccess.CDNType(networkConfig.AdditionalCDN.ValueString()))
	}
	return request, dd
}

func (m *KeyResourceModel) credentialsForAccessKeyCreation(ctx context.Context, useCredA bool) (cloudaccess.Credentials, diag.Diagnostics) {
	if useCredA {
		cred, dd := m.credentialsA(ctx)
		if dd.HasError() {
			return cloudaccess.Credentials{}, dd
		}
		return cloudaccess.Credentials{
			CloudSecretAccessKey: cred.CloudSecretAccessKey.ValueString(),
			CloudAccessKeyID:     cred.CloudAccessKeyID.ValueString(),
		}, dd
	}
	cred, dd := m.credentialsB(ctx)
	if dd.HasError() {
		return cloudaccess.Credentials{}, dd
	}
	return cloudaccess.Credentials{
		CloudSecretAccessKey: cred.CloudSecretAccessKey.ValueString(),
		CloudAccessKeyID:     cred.CloudAccessKeyID.ValueString(),
	}, dd
}

func (m *KeyResourceModel) buildCreateKeyVersionRequest(ctx context.Context, useCredA bool) (cloudaccess.CreateAccessKeyVersionRequest, diag.Diagnostics) {
	var bodyParams cloudaccess.CreateAccessKeyVersionRequestBody
	if useCredA {
		cred, dd := m.credentialsA(ctx)
		if dd.HasError() {
			return cloudaccess.CreateAccessKeyVersionRequest{}, dd
		}
		bodyParams = cloudaccess.CreateAccessKeyVersionRequestBody{
			CloudAccessKeyID:     cred.CloudAccessKeyID.ValueString(),
			CloudSecretAccessKey: cred.CloudSecretAccessKey.ValueString(),
		}
	} else {
		cred, dd := m.credentialsB(ctx)
		if dd.HasError() {
			return cloudaccess.CreateAccessKeyVersionRequest{}, dd
		}
		bodyParams = cloudaccess.CreateAccessKeyVersionRequestBody{
			CloudAccessKeyID:     cred.CloudAccessKeyID.ValueString(),
			CloudSecretAccessKey: cred.CloudSecretAccessKey.ValueString(),
		}
	}
	return cloudaccess.CreateAccessKeyVersionRequest{
		AccessKeyUID: m.AccessKeyUID.ValueInt64(),
		Body:         bodyParams,
	}, nil
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
	client := r.Client.GetCloudAccess()
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
					cred, dd := plan.credentialsA(ctx)
					if diags.Append(dd...); diags.HasError() {
						return nil, diags
					}
					cred.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					cred.VersionGUID = types.StringValue(versionResp.VersionGUID)
					diags.Append(plan.setCredentialsA(ctx, cred)...)
				} else {
					cred, dd := plan.credentialsB(ctx)
					if diags.Append(dd...); diags.HasError() {
						return nil, diags
					}
					cred.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					cred.VersionGUID = types.StringValue(versionResp.VersionGUID)
					diags.Append(plan.setCredentialsB(ctx, cred)...)
				}
				if diags.HasError() {
					return nil, diags
				}
				return plan, diags
			}
		}
		if statusResp.ProcessingStatus == cloudaccess.ProcessingFailed {
			diags.AddError(creationKeyFailed, "Processing failed, retry the terraform apply and verify you've properly formatted the resource arguments.")
			return nil, diags
		}
		select {
		case <-time.After(r.pollingInterval):
			continue
		case <-ctx.Done():
			diags.AddError("reached activation timeout", ctx.Err().Error())
			return nil, diags
		}
	}
}

func (r *KeyResource) waitUntilVersionCreatedCompleted(ctx context.Context, requestID int64, statusTimeout int64, plan *KeyResourceModel, credentialA bool) (*KeyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := r.Client.GetCloudAccess()
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
					cred, dd := plan.credentialsA(ctx)
					if diags.Append(dd...); diags.HasError() {
						return nil, diags
					}
					cred.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					cred.VersionGUID = types.StringValue(versionResp.VersionGUID)
					diags.Append(plan.setCredentialsA(ctx, cred)...)
				} else {
					cred, dd := plan.credentialsB(ctx)
					if diags.Append(dd...); diags.HasError() {
						return nil, diags
					}
					cred.Version = types.Int64Value(statusResp.AccessKeyVersion.Version)
					cred.VersionGUID = types.StringValue(versionResp.VersionGUID)
					diags.Append(plan.setCredentialsB(ctx, cred)...)
				}
				if diags.HasError() {
					return nil, diags
				}
				return plan, diags
			}
		}
		if statusResp.ProcessingStatus == cloudaccess.ProcessingFailed {
			diags.AddError(creationKeyVersionFailed, "Processing failed, retry the terraform apply and verify you've properly formatted the resource arguments.")
			return nil, diags
		}

		select {
		case <-time.After(r.pollingInterval):
			continue
		case <-ctx.Done():
			diags.AddError("reached activation timeout", ctx.Err().Error())
			return nil, diags
		}
	}
}
