package cloudwrapper

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &locationsDataSource{}
	_ datasource.DataSourceWithConfigure = &locationsDataSource{}
)

type (
	locationsDataSource struct {
		client cloudwrapper.CloudWrapper
	}

	locationsDataSourceModel struct {
		Locations []locationModel `tfsdk:"locations"`
	}

	locationModel struct {
		LocationID         types.Int64        `tfsdk:"location_id"`
		LocationName       types.String       `tfsdk:"location_name"`
		MultiCDNLocationID types.String       `tfsdk:"multi_cdn_location_id"`
		TrafficTypes       []trafficTypeModel `tfsdk:"traffic_types"`
	}

	trafficTypeModel struct {
		TrafficType   types.String `tfsdk:"traffic_type"`
		TrafficTypeID types.Int64  `tfsdk:"traffic_type_id"`
		LocationID    types.String `tfsdk:"location_id"`
	}
)

// NewLocationsDataSource returns a new location's data source
func NewLocationsDataSource() datasource.DataSource {
	return &locationsDataSource{}
}

func (d *locationsDataSource) setClient(client cloudwrapper.CloudWrapper) {
	d.client = client
}

// Metadata configures data source's meta information
func (d *locationsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudwrapper_locations"
}

// Configure configures data source at the beginning of the lifecycle
func (d *locationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *locationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CloudWrapper locations",
		Blocks: map[string]schema.Block{
			"locations": schema.ListNestedBlock{
				Description: "List of the locations.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"location_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the location.",
						},
						"location_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the location.",
						},
						"multi_cdn_location_id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier of the multi CDN location.",
						},
					},
					Blocks: map[string]schema.Block{
						"traffic_types": schema.ListNestedBlock{
							Description: "List of traffic types for the location.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"traffic_type_id": schema.Int64Attribute{
										Computed:    true,
										Description: "Unique identifier for the location and traffic type combination.",
									},
									"traffic_type": schema.StringAttribute{
										Computed: true,
										Description: "Represents the traffic type. LIVE applies to low-latency media traffic, such as live streaming. " +
											"LIVE_VOD applies to redundant media traffic, like video on demand content. " +
											"WEB_STANDARD_TLS or WEB_ENHANCED_TLS applies to web content using Standard TLS security or Enhanced TLS security, respectively.",
									},
									"location_id": schema.StringAttribute{
										Computed:    true,
										Description: "Represents the failover map.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *locationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudWrapper Locations DataSource Read")

	var data locationsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	locations, err := d.client.ListLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("reading CloudWrapper Locations", err.Error())
		return
	}

	for _, loc := range locations.Locations {
		trafficTypes := make([]trafficTypeModel, 0)
		for _, trafficType := range loc.TrafficTypes {
			tt := trafficTypeModel{
				TrafficType:   types.StringValue(trafficType.TrafficType),
				TrafficTypeID: types.Int64Value(int64(trafficType.TrafficTypeID)),
				LocationID:    types.StringValue(trafficType.MapName),
			}
			trafficTypes = append(trafficTypes, tt)
		}
		location := locationModel{
			LocationID:         types.Int64Value(int64(loc.LocationID)),
			LocationName:       types.StringValue(loc.LocationName),
			MultiCDNLocationID: types.StringValue(loc.MultiCDNLocationID),
			TrafficTypes:       trafficTypes,
		}
		data.Locations = append(data.Locations, location)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
