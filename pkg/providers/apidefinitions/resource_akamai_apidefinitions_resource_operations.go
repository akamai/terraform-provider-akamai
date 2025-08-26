package apidefinitions

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type apiResourceOperation struct{}

type apiResourceOperationModel struct {
	APIID              types.Int64          `tfsdk:"api_id"`
	Version            types.Int64          `tfsdk:"version"`
	ResourceOperations operationsStateValue `tfsdk:"resource_operations"`
}

// NewAPIResourceOperationResource returns new api resource operations
func NewAPIResourceOperationResource() resource.Resource {
	return &apiResourceOperation{}
}

// Metadata implements resource.Resource Operations.
func (r *apiResourceOperation) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_apidefinitions_resource_operations"
}

// Configure implements resource.Resource Operations.
func (r *apiResourceOperation) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

// Schema implements resource.Resource Operations.
func (r *apiResourceOperation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier for the endpoint",
			},
			"version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version of the endpoint",
			},
			"resource_operations": schema.StringAttribute{
				CustomType: operationsStateType{},
				Required:   true,
				Validators: []validator.String{
					validators.NotEmptyString(),
					OperationsStateValidator(),
				},
				Description: "JSON-formatted information about the API configuration",
			},
		},
	}
}

// Create implements resource.Resource Operations.
func (r *apiResourceOperation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating API Definitions API Resource Operations")

	var data *apiResourceOperationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.upsert(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *apiResourceOperation) upsert(ctx context.Context, data *apiResourceOperationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	result, err := client.ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: data.APIID.ValueInt64(),
	})
	if err != nil {
		diags.AddError("Upsert Resource Operations Failed", err.Error())
		return diags
	}

	var latestEndpointVersion = getLatestVersion(result)
	var latestVersionNumber = latestEndpointVersion.VersionNumber

	if latestEndpointVersion.IsVersionLocked {
		resp, err := client.CloneEndpointVersion(ctx, apidefinitions.CloneEndpointVersionRequest{
			VersionNumber: latestVersionNumber,
			APIEndpointID: data.APIID.ValueInt64(),
		})

		if err != nil {
			diags.AddError("Unable to clone an API Version", err.Error())
			return diags
		}

		latestVersionNumber = resp.VersionNumber
	}

	// Convert stored JSON string into Go struct
	var requestBody = v0.ResourceOperationsRequestBody{}

	err = json.Unmarshal([]byte(data.ResourceOperations.ValueString()), &requestBody)
	if err != nil {
		diags.AddError("Upsert Resource Operations Failed, Unable to deserialize state", err.Error())
		return diags
	}
	// Prepare request
	var resourceOperationRequest = v0.UpdateResourceOperationRequest{
		APIID:         data.APIID.ValueInt64(),
		VersionNumber: latestVersionNumber,
		Body:          requestBody,
	}

	resp, err := clientV0.UpdateResourceOperation(ctx, resourceOperationRequest)
	if err != nil || resp == nil {
		diags.AddError("Upsert Resource Operations Failed", err.Error())
		return diags
	}

	operationsContent, err := serializeResourceOperationResponseIndent((*v0.ResourceOperationResponse)(resp))
	if err != nil {
		diags.AddError("Upsert Resource Operations Failed, Unable to serialize response", err.Error())
		return diags
	}

	data.ResourceOperations = operationsStateValue{types.StringValue(*operationsContent)}
	data.APIID = types.Int64Value(data.APIID.ValueInt64())
	data.Version = types.Int64Value(latestVersionNumber)
	return diags
}

// Read implements resource.Resource Operations.
func (r *apiResourceOperation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading API Resource Operations")

	var data *apiResourceOperationModel

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

func (r *apiResourceOperation) read(ctx context.Context, data *apiResourceOperationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	resourceOperation, err := clientV0.GetResourceOperation(ctx, v0.GetResourceOperationRequest{
		APIID:         data.APIID.ValueInt64(),
		VersionNumber: data.Version.ValueInt64(),
	})
	if err != nil {
		diags.AddError("Reading Resource Operations Failed", err.Error())
		return diags
	}

	operationsContent, err := serializeResourceOperationResponseIndent((*v0.ResourceOperationResponse)(resourceOperation))
	if err != nil {
		diags.AddError("Reading Resource Operations Failed : Unable to serialize response", err.Error())
		return diags
	}

	data.ResourceOperations = operationsStateValue{types.StringValue(*operationsContent)}
	data.APIID = types.Int64Value(data.APIID.ValueInt64())

	return diags
}

// Update implements resource.Resource Operations.
func (r *apiResourceOperation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating API Resource Operations")

	var data *apiResourceOperationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var oldState *apiResourceOperationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.upsert(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource Operations.
func (r *apiResourceOperation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting API Resource Operations")
	var data apiResourceOperationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := client.ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: data.APIID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Deleting Resource Operations Failed", err.Error())
		return
	}

	var latestEndpointVersion = getLatestVersion(result)
	var latestVersionNumber = latestEndpointVersion.VersionNumber

	deleteResponse, err := clientV0.DeleteResourceOperation(ctx, v0.DeleteResourceOperationRequest{
		APIID:         data.APIID.ValueInt64(),
		VersionNumber: latestVersionNumber,
	})
	if err != nil {
		resp.Diagnostics.AddError("Deleting Resource Operations Failed", err.Error())
		return
	}

	respJSON, err := serializeDeleteResourceOperationResponseIndent(deleteResponse)
	if err != nil {
		resp.Diagnostics.AddError("Deleting Resource Operations Failed : Unable to serialize response", err.Error())
		return
	}
	tflog.Info(ctx, string(*respJSON))
}

// ImportState implements resource's ImportState method
func (r *apiResourceOperation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing API Definitions API resource")

	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 {
		resp.Diagnostics.AddError(fmt.Sprintf("ID '%s' incorrectly formatted: should be 'API_ID:VERSION'", req.ID), "")
		return
	}

	endpointID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid API ID '%v'", parts[0]), "")
		return
	}

	versionNumber, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid API version '%v'", parts[1]), "")
		return
	}

	data := apiResourceOperationModel{}
	data.APIID = types.Int64Value(endpointID)
	data.Version = types.Int64Value(versionNumber)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func serializeResourceOperationResponseIndent(response *v0.ResourceOperationResponse) (*string, error) {
	jsonBody, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil
}

func serializeDeleteResourceOperationResponseIndent(response *v0.DeleteResourceOperationResponse) (*string, error) {
	jsonBody, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil
}

func getLatestVersion(result *apidefinitions.ListEndpointVersionsResponse) apidefinitions.APIVersion {
	var response apidefinitions.APIVersion
	for _, v := range result.APIVersions {
		if v.VersionNumber > response.VersionNumber {
			response = v
		}
	}
	return response
}
