package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type geoMapDataSource struct {
	meta meta.Meta
}

type geoMapDataSourceModel struct {
	Domain            types.String              `tfsdk:"domain"`
	Name              types.String              `tfsdk:"map_name"`
	DefaultDatacenter *defaultDatacenter        `tfsdk:"default_datacenter"`
	Assignments       []geographicMapAssignment `tfsdk:"assignments"`
	Links             []link                    `tfsdk:"links"`
}

var (
	_ datasource.DataSource              = &geoMapDataSource{}
	_ datasource.DataSourceWithConfigure = &geoMapDataSource{}
)

// NewGTMGeoMapDataSource returns a new GTM Geographic map data source
func NewGTMGeoMapDataSource() datasource.DataSource { return &geoMapDataSource{} }

func (d *geoMapDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_gtm_geomap"
}

func (d *geoMapDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *geoMapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GTM Geographic map data source.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "GTM domain name.",
			},
			"map_name": schema.StringAttribute{
				Required:    true,
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
				Description: "Specifies the URL path that allows direct navigation to the Geographic map.",
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
	}
}

func (d *geoMapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Geomap DataSource Read")

	var data geoMapDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	geoMap, err := client.GetGeoMap(ctx, gtm.GetGeoMapRequest{
		MapName:    data.Name.ValueString(),
		DomainName: data.Domain.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("fetching GTM Geographic map failed: ", err.Error())
		return
	}

	diags := data.setAttributes(ctx, geoMap)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *geoMapDataSourceModel) setAttributes(ctx context.Context, geoMap *gtm.GetGeoMapResponse) diag.Diagnostics {
	m.Name = types.StringValue(geoMap.Name)
	m.DefaultDatacenter = newDefaultDatacenter(*geoMap.DefaultDatacenter)
	m.Links = getLinks(geoMap.Links)
	diags := m.setAssignments(ctx, geoMap.Assignments)
	if diags.HasError() {
		return diags
	}

	return nil
}

func (m *geoMapDataSourceModel) setAssignments(ctx context.Context, assignments []gtm.GeoAssignment) diag.Diagnostics {
	for _, a := range assignments {
		countries, diags := types.SetValueFrom(ctx, types.StringType, a.Countries)
		if diags.HasError() {
			return diags
		}
		assignmentObj := geographicMapAssignment{
			DatacenterID: types.Int64Value(int64(a.DatacenterID)),
			Nickname:     types.StringValue(a.Nickname),
			Countries:    countries,
		}

		m.Assignments = append(m.Assignments, assignmentObj)
	}
	return nil
}
