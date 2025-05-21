package cloudwrapper

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &locationDataSource{}
	_ datasource.DataSourceWithConfigure = &locationDataSource{}
)

type (
	locationDataSource struct {
		client cloudwrapper.CloudWrapper
	}

	locationDataSourceModel struct {
		LocationName  types.String `tfsdk:"location_name"`
		TrafficType   types.String `tfsdk:"traffic_type"`
		TrafficTypeID types.Int64  `tfsdk:"traffic_type_id"`
		LocationID    types.String `tfsdk:"location_id"`
	}
)

// NewLocationDataSource returns a new location's data source
func NewLocationDataSource() datasource.DataSource {
	return &locationDataSource{}
}

func (d *locationDataSource) setClient(client cloudwrapper.CloudWrapper) {
	d.client = client
}

// Metadata configures data source's meta information
func (d *locationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cloudwrapper_location"
}

// Configure configures data source at the beginning of the lifecycle
func (d *locationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *locationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CloudWrapper location",
		Attributes: map[string]schema.Attribute{
			"location_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the location.",
			},
			"traffic_type": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.OneOf("LIVE", "LIVE_VOD", "WEB_STANDARD_TLS", "WEB_ENHANCED_TLS")},
				Description: "Represents the traffic type. LIVE applies to low-latency media traffic, such as live streaming. " +
					"LIVE_VOD applies to redundant media traffic, like video on demand content. " +
					"WEB_STANDARD_TLS or WEB_ENHANCED_TLS applies to web content using Standard TLS security or Enhanced TLS security, respectively.",
			},
			"traffic_type_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier for the location and traffic type combination.",
			},
			"location_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the location.",
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state
func (d *locationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CloudWrapper Location DataSource Read")

	var data locationDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	locations, err := d.client.ListLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("reading %s", ErrCloudWrapperLocation), err.Error())
		return
	}

	for _, loc := range locations.Locations {
		if loc.LocationName == data.LocationName.ValueString() {
			if trafficType, ok := getMatchingTrafficType(loc.TrafficTypes, data.TrafficType.ValueString()); ok {
				data.LocationID = types.StringValue(trafficType.MapName)
				data.TrafficTypeID = types.Int64Value(int64(trafficType.TrafficTypeID))
				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
				return
			}
		}
	}
	resp.Diagnostics.AddError("No matching location", "no location with given location name and traffic type")
}

func getMatchingTrafficType(items []cloudwrapper.TrafficTypeItem, trafficType string) (*cloudwrapper.TrafficTypeItem, bool) {
	for _, item := range items {
		if item.TrafficType == trafficType {
			return &item, true
		}
	}
	return nil, false
}
