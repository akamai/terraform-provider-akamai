package apidefinitions

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &apiResource{}
	_ resource.ResourceWithConfigure   = &apiResource{}
	_ resource.ResourceWithImportState = &apiResource{}
)

type apiResource struct{}

type apiResourceModel struct {
	ID                types.Int64   `tfsdk:"id"`
	API               apiStateValue `tfsdk:"api"`
	LatestVersion     types.Int64   `tfsdk:"latest_version"`
	StagingVersion    types.Int64   `tfsdk:"staging_version"`
	ProductionVersion types.Int64   `tfsdk:"production_version"`
}

// NewAPIResource returns new api definition API resource
func NewAPIResource() resource.Resource {
	return &apiResource{}
}

// Metadata implements resource.Resource.
func (r *apiResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_apidefinitions_api"
}

// Configure implements resource.Resource.
func (r *apiResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	metaConfig, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if client == nil {
		client = apidefinitions.Client(metaConfig.Session())
	}
	if clientV0 == nil {
		clientV0 = v0.Client(metaConfig.Session())
	}
}

// Schema implements resource.Resource.
func (r *apiResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "API Definition configuration",
		Attributes: map[string]schema.Attribute{
			"api": schema.StringAttribute{
				CustomType: apiStateType{},
				Required:   true,
				Validators: []validator.String{
					validators.NotEmptyString(),
				},
				Description: "JSON-formatted information about the API configuration",
			},
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Unique identifier of the API",
			},
			"latest_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Latest version of the API",
			},
			"staging_version": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Version of the API currently deployed in staging",
			},
			"production_version": schema.Int64Attribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Version of the API currently deployed in production",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *apiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating API Definitions API Resource")

	var data *apiResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.create(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiResource) create(ctx context.Context, data *apiResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	APIState := []byte(data.API.ValueString())

	var registerEndpointRequest = v0.RegisterAPIRequest{}

	err := json.Unmarshal(APIState, &registerEndpointRequest)
	if err != nil {
		diags.AddError("Create API Failed, Unable to deserialize state", err.Error())
		return diags
	}

	resp, err := clientV0.RegisterAPI(ctx, registerEndpointRequest)
	if err != nil {
		diags.AddError("Create API Failed", err.Error())
		return diags
	}

	data.StagingVersion = types.Int64Null()
	data.ProductionVersion = types.Int64Null()
	newEndpointVersion := int64(1)
	return data.populateFromVersion((*v0.API)(resp), newEndpointVersion, false)
}

// Read implements resource.Resource.
func (r *apiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading API Resource")

	var data *apiResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.read(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiResource) read(ctx context.Context, data *apiResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	result, err := client.ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: data.ID.ValueInt64(),
	})

	if err != nil {
		diags.AddError("Reading API Failed", err.Error())
		return diags
	}
	var latestVersion int64
	for _, v := range result.APIVersions {
		if v.VersionNumber > latestVersion {
			latestVersion = v.VersionNumber
		}
	}

	endpoint, err := getEndpoint(ctx, data.ID.ValueInt64())

	if err != nil {
		diags.AddError("Unable to read Endpoint", err.Error())
		return diags
	}

	endpointVersion, err := clientV0.GetAPIVersion(ctx, v0.GetAPIVersionRequest{
		Version: latestVersion,
		ID:      data.ID.ValueInt64(),
	})

	if err != nil {
		diags.AddError("Unable to read Endpoint", err.Error())
		return diags
	}

	data.populateFromEndpoint(endpoint)
	return data.populateFromVersion((*v0.API)(endpointVersion), latestVersion, false)
}

// Update implements resource.Resource.
func (r *apiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating API Resource")

	var data *apiResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var oldState *apiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.StagingVersion = oldState.StagingVersion
	data.ProductionVersion = oldState.ProductionVersion
	resp.Diagnostics.Append(r.update(ctx, oldState, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiResource) update(ctx context.Context, state *apiResourceModel, data *apiResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	jsonPayloadRaw := []byte(data.API.ValueString())

	var body = v0.API{}

	err := json.Unmarshal(jsonPayloadRaw, &body)
	if err != nil {
		diags.AddError("Update API Failed", err.Error())
		return diags
	}

	id := state.ID.ValueInt64()
	endpointVersion, err := client.GetEndpointVersion(ctx, apidefinitions.GetEndpointVersionRequest{
		VersionNumber: state.LatestVersion.ValueInt64(),
		APIEndpointID: id,
	})

	if err != nil {
		diags.AddError("Unable to Get API Version", err.Error())
		return diags
	}

	var versionNumber = state.LatestVersion.ValueInt64()

	if endpointVersion.Locked {
		resp, err := client.CloneEndpointVersion(ctx, apidefinitions.CloneEndpointVersionRequest{
			VersionNumber: versionNumber,
			APIEndpointID: id,
		})

		if err != nil {
			diags.AddError("Unable to clone an API Version", err.Error())
			return diags
		}

		versionNumber = resp.VersionNumber
	}

	body.RecordVersion = ptr.To(endpointVersion.LockVersion)

	updateEndpointVersionReq := v0.UpdateAPIVersionRequest{
		ID:      id,
		Version: versionNumber,
		Body:    v0.UpdateAPIVersionRequestBody(body),
	}

	resp, err := clientV0.UpdateAPIVersion(ctx, updateEndpointVersionReq)
	if err != nil {
		diags.AddError("Update API Failed", err.Error())
		return diags
	}

	return data.populateFromVersion((*v0.API)(resp), versionNumber, false)
}

// Delete implements resource.Resource.
func (r *apiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting API Resource")

	var data *apiResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := getEndpoint(ctx, data.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Endpoint", err.Error())
		return
	}

	diags := deactivateEndpoint(ctx, *endpoint)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	if endpoint.StagingVersion.Status == nil && endpoint.ProductionVersion.Status == nil {
		err := client.DeleteEndpoint(ctx, apidefinitions.DeleteEndpointRequest{APIEndpointID: data.ID.ValueInt64()})
		if err != nil {
			resp.Diagnostics.AddError("Deletion of API Failed", err.Error())
			return
		}
	} else {
		_, err := client.HideEndpoint(ctx, apidefinitions.HideEndpointRequest{APIEndpointID: data.ID.ValueInt64()})
		if err != nil {
			resp.Diagnostics.AddError("Deletion of API Failed", err.Error())
			return
		}
		tflog.Info(ctx, "API has been hidden.")
	}
}

// ImportState implements resource's ImportState method
func (r *apiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing API Definitions API resource")

	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 {
		resp.Diagnostics.AddError(fmt.Sprintf("ID '%s' incorrectly formatted: should be 'API_ID:VERSION'", req.ID), "")
		return
	}

	endpointID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid API id '%v'", parts[0]), "")
		return
	}
	versionNumber, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid API version '%v'", parts[1]), "")
		return
	}

	endpoint, err := getEndpoint(ctx, endpointID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Endpoint", err.Error())
		return
	}

	version, err := clientV0.GetAPIVersion(ctx, v0.GetAPIVersionRequest{ID: endpointID, Version: versionNumber})
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Version", err.Error())
		return
	}

	data := apiResourceModel{}

	data.populateFromEndpoint(endpoint)
	data.populateFromVersion((*v0.API)(version), versionNumber, true)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *apiResourceModel) populateFromEndpoint(endpoint *apidefinitions.EndpointDetail) diag.Diagnostics {
	var diags diag.Diagnostics

	m.StagingVersion = types.Int64Null()
	m.ProductionVersion = types.Int64Null()

	if endpoint.StagingVersion.IsActive() {
		m.StagingVersion = types.Int64Value(*endpoint.StagingVersion.VersionNumber)
	}

	if endpoint.ProductionVersion.IsActive() {
		m.ProductionVersion = types.Int64Value(*endpoint.ProductionVersion.VersionNumber)
	}

	return diags
}

func (m *apiResourceModel) populateFromVersion(resp *v0.API, versionNumber int64, importState bool) diag.Diagnostics {
	var diags diag.Diagnostics
	m.ID = types.Int64Value(*resp.ID)
	m.LatestVersion = types.Int64Value(versionNumber)
	var apiState *string
	var err error
	if importState {
		apiState, err = serializeIndent(resp.RegisterAPIRequest)
	} else {
		apiState, err = serialize(resp.RegisterAPIRequest)
	}
	if err != nil {
		diags.AddError("error parsing API", err.Error())
		return diags
	}
	m.API = apiStateValue{types.StringValue(*apiState)}
	return diags
}

func serializeIndent(version v0.RegisterAPIRequest) (*string, error) {
	jsonBody, err := json.MarshalIndent(version, "", "  ")
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil

}

func serialize(version v0.RegisterAPIRequest) (*string, error) {
	jsonBody, err := json.Marshal(version)
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil
}

func deserialize(body string) (*v0.RegisterAPIRequest, error) {
	endpoint := v0.RegisterAPIRequest{}

	err := json.Unmarshal([]byte(body), &endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}
