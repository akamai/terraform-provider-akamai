package apidefinitions

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &apiDataSource{}
	_ datasource.DataSourceWithConfigure = &apiDataSource{}
)

type (
	apiDataSource struct {
		meta.DataSource
	}

	apiDataSourceModel struct {
		apiResourceModel
		Name types.String `tfsdk:"name"`
	}
)

// NewAPIDataSource returns a new API data source
func NewAPIDataSource() datasource.DataSource {
	return &apiDataSource{}
}

// Metadata configures data source's meta information
func (d apiDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_apidefinitions_api"
}

// Schema defines the schema for the API configuration data source
func (d apiDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "API Definition configuration",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Optional:    true,
				Description: "Unique identifier of the API",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the API",
			},
			"api": schema.StringAttribute{
				CustomType:  apiStateType{},
				Computed:    true,
				Description: "JSON-formatted information about the API latest version configuration",
			},
			"contract_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the contract (without the 'ctr_' prefix) assigned to the API.",
			},
			"group_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The unique identifier for the group (without the 'grp_' prefix) assigned to the API.",
			},
			"latest_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Latest version of the API",
			},
			"staging_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version of the API currently deployed in staging",
			},
			"production_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version of the API currently deployed in production",
			},
		},
	}
}

func (d apiDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d apiDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "API Definitions Configuration DataSource Read")

	var data apiDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &data)...); response.Diagnostics.HasError() {
		return
	}
	var ID int64
	if data.ID.IsNull() {
		byName, err := findAPIByName(ctx, d.Client.GetAPIDefinitions(), data.Name.ValueString())
		if err != nil {
			response.Diagnostics.AddError("Error retrieving API", err.Error())
			return
		}
		ID = *byName
	} else {
		ID = data.ID.ValueInt64()
	}

	apiResponse, err := d.Client.GetAPIDefinitions().ListEndpointVersions(ctx, apidefinitions.ListEndpointVersionsRequest{
		APIEndpointID: ID,
	})
	if err != nil {
		var apiErr *apidefinitions.Error
		if errors.As(err, &apiErr) && apiErr != nil && (apiErr.Status == 403 || apiErr.Status == 404) {
			response.Diagnostics.AddError("Error retrieving API", fmt.Sprintf("unable to find API with ID %d", ID))
			return
		}
		response.Diagnostics.AddError("Error retrieving API",
			fmt.Sprintf("Could not retrieve Endpoint Version List for API ID %d: %s", ID, err),
		)
		return
	}

	var latestVersion = getLatestAPIConfigVersion(apiResponse)
	endpoint, err := getEndpoint(ctx, d.Client.GetAPIDefinitions(), ID)
	if err != nil {
		response.Diagnostics.AddError("Error retrieving API", err.Error())
		return
	}

	endpointVersion, err := d.Client.GetAPIDefinitionsV0().GetAPIVersion(ctx, v0.GetAPIVersionRequest{
		Version: latestVersion.VersionNumber,
		ID:      ID,
	})
	if err != nil {
		response.Diagnostics.AddError("Error retrieving API", err.Error())
		return
	}

	data.populateFromEndpoint(endpoint)
	data.populateFromVersion((*v0.API)(endpointVersion), latestVersion.VersionNumber, true)
	data.ContractID = types.StringValue(endpointVersion.ContractID)
	data.GroupID = types.Int64Value(endpointVersion.GroupID)

	if response.Diagnostics.Append(response.State.Set(ctx, &data)...); response.Diagnostics.HasError() {
		return
	}
}

func findAPIByName(ctx context.Context, client apidefinitions.APIDefinitions, expectedName string) (*int64, error) {
	endpoints, err := client.ListEndpoints(ctx, apidefinitions.ListEndpointsRequest{
		Contains: expectedName,
		PageSize: 1000,
	})
	if err != nil {
		return nil, err
	}

	var API *apidefinitions.Endpoint

	for i := range endpoints.APIEndpoints {
		if endpoints.APIEndpoints[i].APIEndpointName == expectedName {
			API = &endpoints.APIEndpoints[i]
			break
		}
	}

	if API != nil {
		return &API.APIEndpointID, nil
	}

	return nil, fmt.Errorf("unable to find API with Name %s", expectedName)

}

func getLatestAPIConfigVersion(result *apidefinitions.ListEndpointVersionsResponse) apidefinitions.APIVersion {
	var response apidefinitions.APIVersion
	for _, v := range result.APIVersions {
		if v.VersionNumber > response.VersionNumber {
			response = v
		}
	}
	return response
}
