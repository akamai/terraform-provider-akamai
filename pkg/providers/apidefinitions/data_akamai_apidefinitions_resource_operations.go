package apidefinitions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SubproviderName defines name of the botman subprovider
const SubproviderName = "apidefinitions"

var (
	_                                 datasource.DataSource              = &resourceOperationsDataSource{}
	_                                 datasource.DataSourceWithConfigure = &resourceOperationsDataSource{}
	resourceOperationsDataSourceMutex sync.Mutex
)

type (
	resourceOperationsDataSource struct{}
	resourceOperationsModel      struct {
		apiResourceOperationModel
		ResourcePath types.String `tfsdk:"resource_path"`
		ResourceName types.String `tfsdk:"resource_name"`
	}

	operation struct {
		APIEndpointID      int64  `json:"api_id"`
		APIResourceID      int64  `json:"api_resource_id"`
		APIResourceLogicID int64  `json:"api_resource_logic_id"`
		OperationID        string `json:"id"`
		OperationName      string `json:"name"`
		OperationPurpose   string `json:"purpose"`
		ResourcePath       string `json:"resource_path"`
		ResourceName       string `json:"resource_name"`
	}
)

// NewResourceOperationsDataSource returns a new resource operations data source
func NewResourceOperationsDataSource() datasource.DataSource {
	return &resourceOperationsDataSource{}
}

// Metadata configures data source's meta information
func (r resourceOperationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_apidefinitions_resource_operations"
}

// Configure configures data source at the beginning of the lifecycle
func (r resourceOperationsDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	metaConfig, ok := request.ProviderData.(meta.Meta)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
	}
	if client == nil {
		client = apidefinitions.Client(metaConfig.Session())
	}
}

func (r resourceOperationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Retrieve resource operations for a specific API and version",
		Attributes: map[string]schema.Attribute{
			"api_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier for the endpoint",
			},
			"resource_path": schema.StringAttribute{
				Optional:    true,
				Description: "Resource path to search",
			},
			"resource_name": schema.StringAttribute{
				Optional:    true,
				Description: "Resource name to search",
			},
			"resource_operations": schema.StringAttribute{
				CustomType: operationsStateType{},
				Computed:   true,
				Validators: []validator.String{
					validators.NotEmptyString(),
					operationsStateValidator{},
				},
				Description: "JSON-formatted information about the API configuration",
			},
			"version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version of the endpoint",
			},
		},
	}
}

func (r resourceOperationsDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "API Definitions Resource Operations DataSource Read")

	var data resourceOperationsModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &data)...); response.Diagnostics.HasError() {
		return
	}

	var apiResponse *apidefinitions.SearchResourceOperationsResponse
	apiResponse, err := searchResourceAndOperations(ctx, data)
	if err != nil {
		response.Diagnostics.AddError(
			"Error retrieving resource operations",
			fmt.Sprintf("Could not retrieve resource operations for API ID %d: %s", data.APIID.ValueInt64(), err),
		)
		return
	}

	var version int64
	for _, ep := range apiResponse.APIEndpoints {
		if ep.APIEndpointID == data.APIID.ValueInt64() {
			version = ep.ProductionVersion.VersionNumber
			data.Version = types.Int64Value(version)
			break
		}
	}
	var operations []operation
	for _, op := range apiResponse.Operations {

		if (data.APIID.ValueInt64() != 0 && data.APIID.ValueInt64() != op.APIEndpointID) || !op.Metadata.IsActive {
			continue
		}

		resourcePath, resourceName := getResourcePathByLogicID(apiResponse.Resources, op.APIResourceLogicID)

		if data.ResourcePath.ValueString() != "" && data.ResourcePath.ValueString() != resourcePath {
			continue
		}

		if data.ResourceName.ValueString() != "" && data.ResourceName.ValueString() != resourceName {
			continue
		}

		operations = append(operations, operation{
			APIEndpointID:      op.APIEndpointID,
			APIResourceID:      op.APIResourceID,
			APIResourceLogicID: op.APIResourceLogicID,
			OperationID:        op.OperationID,
			OperationName:      op.OperationName,
			OperationPurpose:   op.OperationPurpose,
			ResourceName:       resourceName,
			ResourcePath:       resourcePath,
		})
	}
	operationsJSON, err := json.Marshal(operations)
	if err != nil {
		response.Diagnostics.AddError(
			"Error serializing operations to JSON",
			fmt.Sprintf("Could not serialize operations: %s", err),
		)
		return
	}
	data.ResourceOperations = operationsStateValue{types.StringValue(string(operationsJSON))}

	// Set the state
	if response.Diagnostics.Append(response.State.Set(ctx, &data)...); response.Diagnostics.HasError() {
		return
	}
}

func searchResourceAndOperations(ctx context.Context, data resourceOperationsModel) (*apidefinitions.SearchResourceOperationsResponse, error) {
	tflog.Debug(ctx, "API Definitions search Resource Operations")
	cacheKey := fmt.Sprintf("%s:%d", "searchResourceAndOperations", data.APIID.ValueInt64())

	apiResponse := &apidefinitions.SearchResourceOperationsResponse{}

	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, apiResponse)

	if err == nil {
		// Successfully retrieved from cache
		return apiResponse, nil
	}

	resourceOperationsDataSourceMutex.Lock()
	defer func() {
		tflog.Debug(ctx, "Unlocking mutex")
		resourceOperationsDataSourceMutex.Unlock()
	}()

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, apiResponse)
	if err == nil {
		return apiResponse, nil
	}

	apiResponse, err = client.SearchResourceOperations(ctx)
	if err != nil {
		return nil, err
	}

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, apiResponse)

	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		return nil, err
	}

	return apiResponse, nil
}

func getResourcePathByLogicID(resources []apidefinitions.Resource, logicID int64) (string, string) {
	for _, resource := range resources {
		if resource.APIResourceLogicID == logicID {
			return resource.ResourcePath, resource.APIResourceName
		}
	}
	return "", ""
}
