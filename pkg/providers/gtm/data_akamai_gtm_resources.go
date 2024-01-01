package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &resourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &resourcesDataSource{}
)

// resourcesDataSourceModel describes the data source data model for GTM resources data source.
type resourcesDataSourceModel struct {
	ID        types.String     `tfsdk:"id"`
	Domain    types.String     `tfsdk:"domain"`
	Resources []domainResource `tfsdk:"resources"`
}

var (
	resourcesBlock = map[string]schema.Block{
		"resources": schema.SetNestedBlock{
			Description: "GTM Resource data source.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:    true,
						Description: "A descriptive label for the resource.",
					},
					"aggregation_type": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies how GTM handles different load numbers when multiple load servers are used for a data center or property.",
					},
					"constrained_property": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the name of the property that this resource constraints.",
					},
					"decay_rate": schema.Float64Attribute{
						Computed:    true,
						Description: "For internal use only.",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive note to help you track what the resource constraints.",
					},
					"host_header": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the host header used when fetching the load object.",
					},
					"leader_string": schema.StringAttribute{
						Computed:    true,
						Description: "Specifies the text that comes before the loadObject.",
					},
					"least_squares_decay": schema.Float64Attribute{
						Computed:    true,
						Description: "For internal use only.",
					},
					"load_imbalance_percentage": schema.Float64Attribute{
						Computed:    true,
						Description: "Indicates the percent of load imbalance factor for the domain.",
					},
					"max_u_multiplicative_increment": schema.Float64Attribute{
						Computed:    true,
						Description: "For internal use only.",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Indicates the kind of loadObject format used to determine the load on the resource.",
					},
					"upper_bound": schema.Int64Attribute{
						Computed:    true,
						Description: "An optional sanity check that specifies the maximum allowed value for any component of the load object.",
					},
				},
				Blocks: resourceBlocks,
			},
		},
	}
)

// NewGTMResourcesDataSource returns a new GTM resources data source.
func NewGTMResourcesDataSource() datasource.DataSource {
	return &resourcesDataSource{}
}

// resourcesDataSource defines the data source implementation for fetching GTM resources information.
type resourcesDataSource struct {
	meta meta.Meta
}

// Metadata configures data source's meta information
func (d *resourcesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_gtm_resources"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *resourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Schema is used to define data source's terraform schema
func (d *resourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GTM Resources data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source.",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "GTM domain name.",
			},
		},
		Blocks: resourcesBlock,
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *resourcesDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Resources DataSource Read")

	var data resourcesDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &data)...); response.Diagnostics.HasError() {
		return
	}

	client := frameworkInst.Client(d.meta)
	resources, err := client.ListResources(ctx, data.Domain.ValueString())
	if err != nil {
		response.Diagnostics.AddError("fetching GTM resources failed:", err.Error())
		return
	}

	data.ID = types.StringValue("akamai_gtm_resources")
	data.Resources = getResources(resources)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}
