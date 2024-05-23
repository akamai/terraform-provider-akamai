package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type geoMapsDataSource struct {
	meta meta.Meta
}

type geoMapsDataSourceModel struct {
	ID      types.String    `tfsdk:"id"`
	Domain  types.String    `tfsdk:"domain"`
	GeoMaps []geographicMap `tfsdk:"geo_maps"`
}

var (
	_ datasource.DataSource              = &geoMapsDataSource{}
	_ datasource.DataSourceWithConfigure = &geoMapsDataSource{}
)

// NewGTMGeoMapsDataSource returns a new GTM Geographic maps data source
func NewGTMGeoMapsDataSource() datasource.DataSource { return &geoMapsDataSource{} }

func (d *geoMapsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_gtm_geomaps"
}

func (d *geoMapsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	d.meta = meta.Must(req.ProviderData)
}

func (d *geoMapsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GTM Geographic maps data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source.",
				DeprecationMessage:  "Required by the terraform plugin testing framework, always set to `gtm_geomaps`.",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "GTM domain name.",
			},
			"geo_maps": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of geographic maps within the domain.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "A descriptive label for the Geographic map.",
						},
						"assignments": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Contains information about the geographic zone groupings of countries.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"datacenter_id": schema.Int64Attribute{
										Computed:    true,
										Description: "A unique identifier for an existing data center in the domain.",
									},
									"countries": schema.SetAttribute{
										Computed: true,
										Description: "Specifies an array of two-letter ISO 3166 country codes, or " +
											"for finer subdivisions, the two-letter country code and the two-letter " +
											"state or province code separated by a forward slash.",
										ElementType: types.StringType,
									},
									"nickname": schema.StringAttribute{
										Computed:    true,
										Description: "A descriptive label for the group.",
									},
								},
							},
						},
						"default_datacenter": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "A placeholder for all other geographic zones, countries not found in these geographic zones.",
							Attributes: map[string]schema.Attribute{
								"datacenter_id": schema.Int64Attribute{
									Computed:    true,
									Description: "For each property, an identifier for all other geographic zones.",
								},
								"nickname": schema.StringAttribute{
									Computed:    true,
									Description: "A descriptive label for all other geographic zones.",
								},
							},
						},
						"links": schema.SetNestedAttribute{
							Computed:    true,
							Description: "Specifies the URL path that allows direct navigation to the Geographic maps.",
							NestedObject: schema.NestedAttributeObject{
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
					},
				},
			},
		},
	}
}

func (d *geoMapsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Geomaps DataSource Read")

	var data geoMapsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	geoMaps, err := client.ListGeoMaps(ctx, data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("fetching GTM Geographic maps failed: ", err.Error())
		return
	}

	data.ID = types.StringValue("gtm_geomaps")
	geographicMaps, diags := getGeographicMaps(ctx, geoMaps)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.GeoMaps = geographicMaps
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
