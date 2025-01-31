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
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type apiResourceOperation struct{}

type apiResourceOperationModel struct {
	APIEndpointID      types.Int64  `tfsdk:"endpoint_id"`
	ResourceOperations types.String `tfsdk:"resource_operations"`
}

// NewAPIResourceOperationResource returns new api resource operations
func NewAPIResourceOperationResource() resource.Resource {
	return &apiResourceOperation{}
}

// Metadata implements resource.Resource.
func (r *apiResourceOperation) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_apidefinitions_resource_operations"
}

// Configure implements resource.Resource.
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

// Schema implements resource.Resource.
func (r *apiResourceOperation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint_id": schema.Int64Attribute{
				Required:    true,
				Description: "(needed)",
			},
			"resource_operations": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "JSON-formatted information about the API Resource Operations",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *apiResourceOperation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating API Definitions API Resource Operation")

	var data *apiResourceOperationModel

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

func (r *apiResourceOperation) create(ctx context.Context, data *apiResourceOperationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	result, err := client.ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: data.APIEndpointID.ValueInt64(),
	})

	if err != nil {
		diags.AddError("Reading Resource Operations Failed", err.Error())
		return diags
	}

	var latestVersion int64
	for _, v := range result.APIVersions {
		if v.VersionNumber > latestVersion {
			latestVersion = v.VersionNumber
		}
	}

	// Convert stored JSON string into Go struct
	var requestBody = v0.UpdateResourceOperationResponse{}

	err = json.Unmarshal([]byte(data.ResourceOperations.ValueString()), &requestBody)

	if err != nil {
		diags.AddError("Create Resource Operations Failed, Unable to deserialize state", err.Error())
		return diags
	}
	// Prepare request
	var resourceOperationRequest = v0.UpdateResourceOperationRequest{
		APIEndpointID: data.APIEndpointID.ValueInt64(),
		VersionNumber: latestVersion,
		Body:          requestBody,
	}

	resp, err := clientV0.UpdateResourceOperation(ctx, resourceOperationRequest)

	if err != nil || resp == nil {
		diags.AddError("Create Resource Operations Failed", err.Error())
		return diags
	}

	jsonBytes, err := json.Marshal(resp)

	if err != nil {
		diags.AddError("Error reading the response", err.Error())
		return diags
	}

	// Normalize JSON before storing in Terraform state
	normalizeJSON, err := normalizeJSON(string(jsonBytes))
	if err != nil {
		diags.AddError("Error normalizing JSON", err.Error())
		return diags
	}

	data.APIEndpointID = types.Int64Value(data.APIEndpointID.ValueInt64())
	data.ResourceOperations = types.StringValue(normalizeJSON)
	return diags
}

// Read implements resource.Resource.
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
	result, err := client.ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: data.APIEndpointID.ValueInt64(),
	})

	if err != nil {
		diags.AddError("Reading Resource Operations Failed", err.Error())
		return diags
	}
	var latestVersion int64
	for _, v := range result.APIVersions {
		if v.VersionNumber > latestVersion {
			latestVersion = v.VersionNumber
		}
	}
	resourceOperation, err := clientV0.GetResourceOperation(ctx, v0.GetResourceOperationRequest{
		APIEndpointID: data.APIEndpointID.ValueInt64(),
		VersionNumber: latestVersion, //looked up above in ListEndpointVersions
	})

	if err != nil {
		diags.AddError("Unable to read Resource Operations", err.Error())
		return diags
	}

	jsonBody, err := json.Marshal(resourceOperation)

	if err != nil {
		diags.AddError("Error marshalling resource operation JSON", err.Error())
		return diags
	}

	//  Normalize JSON before storing state
	normalizeJSON, err := normalizeJSON(string(jsonBody))
	if err != nil {
		diags.AddError("Error normalizing JSON", err.Error())
		return diags
	}

	data.ResourceOperations = types.StringValue(normalizeJSON)
	data.APIEndpointID = types.Int64Value(data.APIEndpointID.ValueInt64())

	return diags
}

// Update implements resource.Resource.
func (r *apiResourceOperation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating API Resource")

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

	resp.Diagnostics.Append(r.create(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (r *apiResourceOperation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting API Resource")
	var data apiResourceOperationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := client.ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: data.APIEndpointID.ValueInt64(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Reading Resource Operations Failed", err.Error())
	}
	var latestVersion int64
	for _, v := range result.APIVersions {
		if v.VersionNumber > latestVersion {
			latestVersion = v.VersionNumber
		}
	}

	var deleteResourceOperationRequest = v0.DeleteResourceOperationRequest{
		APIEndpointID: data.APIEndpointID.ValueInt64(),
		VersionNumber: latestVersion,
	}

	deleteResponse, err := clientV0.DeleteResourceOperation(ctx, deleteResourceOperationRequest)

	if err != nil {
		resp.Diagnostics.AddError("Deletion of API Resource Operations Failed", err.Error())
		return
	}

	jsonBody, err := json.Marshal(deleteResponse)

	if err != nil {
		resp.Diagnostics.AddError("Unable to parse Deletion of API Resource Operations response", err.Error())
		return
	}

	tflog.Info(ctx, string(jsonBody))

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
		resp.Diagnostics.AddError(fmt.Sprintf("invalid API id '%v'", parts[0]), "")
		return
	}

	versionNumber, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid API version '%v'", parts[1]), "")
		return
	}

	resourceOperation, err := clientV0.GetResourceOperation(ctx, v0.GetResourceOperationRequest{
		APIEndpointID: endpointID,
		VersionNumber: versionNumber,
	})

	if err != nil {
		resp.Diagnostics.AddError("Unable to read Resource Operations", err.Error())
		return
	}

	jsonBody, err := json.Marshal(resourceOperation)

	if err != nil {
		resp.Diagnostics.AddError("Error marshalling resource operation JSON", err.Error())
		return
	}

	//  Normalize JSON before storing state
	normalizeJSON, err := normalizeJSON(string(jsonBody))
	if err != nil {
		resp.Diagnostics.AddError("Error normalizing JSON", err.Error())
		return
	}

	data := apiResourceOperationModel{}
	data.ResourceOperations = types.StringValue(normalizeJSON)
	data.APIEndpointID = types.Int64Value(endpointID)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func serializeWithIndent(version v0.RegisterAPIRequest) (*string, error) {
	jsonBody, err := json.MarshalIndent(version, "", "  ")
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil

}

// Function to normalize JSON for Terraform state
func normalizeJSON(input string) (string, error) {
	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(input), &jsonObj)
	if err != nil {
		return "", err
	}

	normalizedBytes, err := json.Marshal(jsonObj) // Ensures consistent key order
	if err != nil {
		return "", err
	}
	return string(normalizedBytes), nil
}
