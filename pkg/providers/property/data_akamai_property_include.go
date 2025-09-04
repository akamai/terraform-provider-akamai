package property

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
)

var _ datasource.DataSource = &includeDataSource{}
var _ datasource.DataSourceWithConfigure = &includeDataSource{}

// NewIncludeDataSource returns a new property include data source
func NewIncludeDataSource() datasource.DataSource {
	return &includeDataSource{}
}

// includeDataSource defines the data source implementation for fetching property include information.
type includeDataSource struct {
	meta meta.Meta
}

// includeDataSourceModel describes the data source data model for PropertyIncludeDataSource.
type includeDataSourceModel struct {
	AssetID           types.String `tfsdk:"asset_id"`
	ContractID        types.String `tfsdk:"contract_id"`
	GroupID           types.String `tfsdk:"group_id"`
	IncludeID         types.String `tfsdk:"include_id"`
	Name              types.String `tfsdk:"name"`
	Type              types.String `tfsdk:"type"`
	LatestVersion     types.Int64  `tfsdk:"latest_version"`
	StagingVersion    types.Int64  `tfsdk:"staging_version"`
	ProductionVersion types.Int64  `tfsdk:"production_version"`
	ID                types.String `tfsdk:"id"`
}

// Metadata configures data source's meta information
func (d *includeDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_property_include"
}

// Schema is used to define data source's terraform schema
func (d *includeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Property Include data source",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				MarkdownDescription: "Identifies the contract under which the include was created",
				Required:            true,
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "Identifies the group under which the include was created",
				Required:            true,
			},
			"include_id": schema.StringAttribute{
				MarkdownDescription: "Identifies the group under which the include was created",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A descriptive name for the include",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Specifies the type of the include, either 'MICROSERVICES' or 'COMMON_SETTINGS'",
				Computed:            true,
			},
			"latest_version": schema.Int64Attribute{
				MarkdownDescription: "Specifies the most recent version of the include",
				Computed:            true,
			},
			"staging_version": schema.Int64Attribute{
				MarkdownDescription: "The most recent version which was activated to the test network",
				Computed:            true,
			},
			"production_version": schema.Int64Attribute{
				MarkdownDescription: "The most recent version which was activated to the production network",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source",
				Computed:            true,
			},
			"asset_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the include in the Identity and Access Management API.",
			},
		},
	}
}

// Configure  configures data source at the beginning of the lifecycle
func (d *includeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	d.meta = meta.Must(req.ProviderData)
}

// Read is called when the provider must read data source values in order to update state
func (d *includeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "PropertyIncludeDataSource Read")

	var data includeDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	getIncludeResp, err := client.GetInclude(ctx, papi.GetIncludeRequest{
		ContractID: data.ContractID.ValueString(),
		GroupID:    data.GroupID.ValueString(),
		IncludeID:  data.IncludeID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching property include failed", err.Error())
		return
	}

	include := getIncludeResp.Include

	data.AssetID = types.StringValue(include.AssetID)
	data.Name = types.StringValue(include.IncludeName)
	data.Type = types.StringValue(string(include.IncludeType))
	data.LatestVersion = types.Int64Value(int64(include.LatestVersion))

	if include.StagingVersion != nil {
		data.StagingVersion = types.Int64Value(int64(*include.StagingVersion))
	}

	if include.ProductionVersion != nil {
		data.ProductionVersion = types.Int64Value(int64(*include.ProductionVersion))
	}

	data.ID = data.IncludeID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
