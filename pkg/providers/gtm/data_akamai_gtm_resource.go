package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &resourceDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceDataSource{}
)

// NewGTMResourceDataSource returns a new GTM resource data source.
func NewGTMResourceDataSource() datasource.DataSource {
	return &resourceDataSource{}
}

// resourceDataSource defines the data source implementation for fetching GTM resource information.
type resourceDataSource struct {
	meta meta.Meta
}

// resourceDataSourceModel describes the data source data model for GTM resource data source.
type resourceDataSourceModel struct {
	ID                          types.String       `tfsdk:"id"`
	Domain                      types.String       `tfsdk:"domain"`
	ResourceName                types.String       `tfsdk:"resource_name"`
	AggregationType             types.String       `tfsdk:"aggregation_type"`
	ConstrainedProperty         types.String       `tfsdk:"constrained_property"`
	DecayRate                   types.Float64      `tfsdk:"decay_rate"`
	Description                 types.String       `tfsdk:"description"`
	HostHeader                  types.String       `tfsdk:"host_header"`
	LeaderString                types.String       `tfsdk:"leader_string"`
	LeastSquaresDecay           types.Float64      `tfsdk:"least_squares_decay"`
	LoadImbalancePercentage     types.Float64      `tfsdk:"load_imbalance_percentage"`
	MaxUMultiplicativeIncrement types.Float64      `tfsdk:"max_u_multiplicative_increment"`
	Type                        types.String       `tfsdk:"type"`
	UpperBound                  types.Int64        `tfsdk:"upper_bound"`
	Links                       []link             `tfsdk:"links"`
	ResourceInstances           []resourceInstance `tfsdk:"resource_instances"`
}

// Metadata configures data source's meta information.
func (d *resourceDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_gtm_resource"
}

// Configure configures data source at the beginning of the lifecycle.
func (d *resourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

var (
	resourceBlocks = map[string]schema.Block{
		"links": schema.SetNestedBlock{
			Description: "Specifies the URL path that allows direct navigation to the resource.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"rel": schema.StringAttribute{
						Computed:    true,
						Description: "Indicates the link relationship of the object.",
					},
					"href": schema.StringAttribute{
						Computed:    true,
						Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
					},
				},
			},
		},
		"resource_instances": schema.SetNestedBlock{
			Description: "Instances of the resource.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"load_object": schema.StringAttribute{
						Computed:    true,
						Description: "Identifies the load object file used to report real-time information about the current load, maximum allowable load and target load on each resource.",
					},
					"load_object_port": schema.Int64Attribute{
						Computed:    true,
						Description: "Specifies the TCP port of the loadObject.",
					},
					"load_servers": schema.SetAttribute{
						Computed:    true,
						Description: "Specifies the list of servers to requests the load object from.",
						ElementType: types.StringType,
					},
					"datacenter_id": schema.Int64Attribute{
						Computed:    true,
						Description: "A unique identifier for an existing data center in the domain.",
					},
					"use_default_load_object": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to use default loadObject.",
					},
				},
			},
		},
	}
)

// Schema is used to define data source's terraform schema.
func (d *resourceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GTM Resource data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source.",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "GTM domain name.",
			},
			"resource_name": schema.StringAttribute{
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
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *resourceDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Resource DataSource Read")

	var data resourceDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &data)...); response.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	resource, err := client.GetResource(ctx, data.ResourceName.ValueString(), data.Domain.ValueString())
	if err != nil {
		response.Diagnostics.AddError("fetching GTM resource failed:", err.Error())
		return
	}

	data.setAttributes(resource)

	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (m *resourceDataSourceModel) setAttributes(resource *gtm.Resource) {

	m.AggregationType = types.StringValue(resource.AggregationType)
	m.ConstrainedProperty = types.StringValue(resource.ConstrainedProperty)
	m.DecayRate = types.Float64Value(float64(resource.DecayRate))
	m.Description = types.StringValue(resource.Description)
	m.HostHeader = types.StringValue(resource.HostHeader)
	m.LeaderString = types.StringValue(resource.LeaderString)
	m.LeastSquaresDecay = types.Float64Value(float64(resource.LeastSquaresDecay))
	m.LoadImbalancePercentage = types.Float64Value(float64(resource.LoadImbalancePercentage))
	m.MaxUMultiplicativeIncrement = types.Float64Value(float64(resource.MaxUMultiplicativeIncrement))
	m.Type = types.StringValue(resource.Type)
	m.UpperBound = types.Int64Value(int64(resource.UpperBound))
	m.setLinks(resource.Links)
	m.setResourceInstances(resource.ResourceInstances)
	m.ID = types.StringValue(resource.Name)
}

func (m *resourceDataSourceModel) setLinks(links []*gtm.Link) {

	for _, l := range links {
		linkObject := link{
			Rel:  types.StringValue(l.Rel),
			Href: types.StringValue(l.Href),
		}

		m.Links = append(m.Links, linkObject)
	}
}

func (m *resourceDataSourceModel) setResourceInstances(resourceInstances []*gtm.ResourceInstance) {

	for _, res := range resourceInstances {
		resourceInstanceObject := resourceInstance{
			DataCenterID:         types.Int64Value(int64(res.DatacenterID)),
			UseDefaultLoadObject: types.BoolValue(res.UseDefaultLoadObject),
			LoadObject:           types.StringValue(res.LoadObject.LoadObject),
			LoadObjectPort:       types.Int64Value(int64(res.LoadObject.LoadObjectPort)),
		}

		for _, server := range res.LoadServers {
			resourceInstanceObject.LoadServers = append(resourceInstanceObject.LoadServers, types.StringValue(server))
		}

		m.ResourceInstances = append(m.ResourceInstances, resourceInstanceObject)
	}
}
