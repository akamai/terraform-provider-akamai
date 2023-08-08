package cloudwrapper

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &propertiesDataSource{}
	_ datasource.DataSourceWithConfigure = &propertiesDataSource{}
)

type (
	propertiesDataSource struct {
		client cloudwrapper.CloudWrapper
	}

	propertiesDataSourceModel struct {
		ID          types.String    `tfsdk:"id"`
		ContractIDs types.List      `tfsdk:"contract_ids"`
		Unused      types.Bool      `tfsdk:"unused"`
		Properties  []propertyModel `tfsdk:"properties"`
	}

	propertyModel struct {
		PropertyID   types.Int64  `tfsdk:"property_id"`
		Type         types.String `tfsdk:"type"`
		PropertyName types.String `tfsdk:"property_name"`
		ContractID   types.String `tfsdk:"contract_id"`
		GroupID      types.Int64  `tfsdk:"group_id"`
	}
)

// NewPropertiesDataSource returns a new properties' data source
func NewPropertiesDataSource() datasource.DataSource {
	return &propertiesDataSource{}
}

// SetClient assigns given client to properties data source
func (d *propertiesDataSource) SetClient(client cloudwrapper.CloudWrapper) {
	d.client = client
}

// Metadata configures data source's meta information
func (d *propertiesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudwrapper_properties"
}

// Configure configures data source at the beginning of the lifecycle
func (d *propertiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	if d.client != nil {
		return
	}

	m, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	d.client = cloudwrapper.Client(m.Session())
}

// Schema is used to define data source's terraform schema
func (d *propertiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CloudWrapper properties",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source.",
			},
			"contract_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of contract IDs with Cloud Wrapper entitlement.",
			},
			"unused": schema.BoolAttribute{
				Optional:    true,
				Description: "Specify whether the response should contain only unused properties.",
			},
		},
		Blocks: map[string]schema.Block{
			"properties": schema.ListNestedBlock{
				Description: "List of all unused properties.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"property_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Property ID of the property.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of property. MEDIA applies to live or video on demand content. WEB applies to website or app content.",
						},
						"property_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the property belonging to the origin.",
						},
						"contract_id": schema.StringAttribute{
							Computed:    true,
							Description: "Contract ID having Cloud Wrapper entitlement.",
						},
						"group_id": schema.Int64Attribute{
							Computed:    true,
							Description: "ID of the group which the property belongs to.",
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *propertiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudWrapper Properties DataSource Read")

	var data propertiesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var contracts []string
	if resp.Diagnostics.Append(data.ContractIDs.ElementsAs(ctx, &contracts, false)...); resp.Diagnostics.HasError() {
		return
	}

	properties, err := d.client.ListProperties(ctx, cloudwrapper.ListPropertiesRequest{
		Unused:      data.Unused.ValueBool(),
		ContractIDs: contracts,
	})
	if err != nil {
		resp.Diagnostics.AddError("Reading CloudWrapper Properties", err.Error())
		return
	}

	for _, prp := range properties.Properties {
		data.Properties = append(data.Properties, propertyModel{
			PropertyID:   types.Int64Value(prp.PropertyID),
			Type:         types.StringValue(string(prp.Type)),
			PropertyName: types.StringValue(prp.PropertyName),
			ContractID:   types.StringValue(prp.ContractID),
			GroupID:      types.Int64Value(prp.GroupID),
		})
	}

	data.ID = types.StringValue(uuid.NewString())

	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}
