package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &asmapDataSource{}
	_ datasource.DataSourceWithConfigure = &asmapDataSource{}
)

// NewGTMAsmapDataSource returns a new GTM asmap data source
func NewGTMAsmapDataSource() datasource.DataSource {
	return &asmapDataSource{}
}

type asmapDataSourceModel struct {
	ID                types.String       `tfsdk:"id"`
	Domain            types.String       `tfsdk:"domain"`
	Name              types.String       `tfsdk:"map_name"`
	DefaultDatacenter *defaultDatacenter `tfsdk:"default_datacenter"`
	Assignments       []asMapAssignment  `tfsdk:"assignments"`
	Links             []link             `tfsdk:"links"`
}

type asmapDataSource struct {
	meta meta.Meta
}

// Configure configures data source at the beginning of the lifecycle.
func (d *asmapDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
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

// Metadata configures data source's meta information.
func (*asmapDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_gtm_asmap"
}

var (
	asmapBlock = map[string]schema.Block{
		"assignments": schema.ListNestedBlock{
			Description: "Contains information about the AS zone groupings of AS IDs.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"datacenter_id": schema.Int64Attribute{
						Computed:    true,
						Description: "A unique identifier for an existing data center in the domain.",
					},
					"as_numbers": schema.SetAttribute{
						Computed:    true,
						Description: "Specifies an array of AS numbers.",
						ElementType: types.Int64Type,
					},
					"nickname": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for the group.",
					},
				},
			},
		},
		"default_datacenter": schema.SingleNestedBlock{
			Description: "A placeholder for all other AS zones, AS IDs not found in these AS zones.",
			Attributes: map[string]schema.Attribute{
				"datacenter_id": schema.Int64Attribute{
					Computed:    true,
					Description: "For each property, an identifier for all other AS zones",
				},
				"nickname": schema.StringAttribute{
					Computed:    true,
					Description: "A descriptive label for all other AS zones",
				},
			},
		},
		"links": schema.ListNestedBlock{
			Description: "Specifies the URL path that allows direct navigation to the AS map.",
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
	}
)

// Schema is used to define data source's terraform schema.
func (d *asmapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GTM AS map data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				DeprecationMessage:  "Required by the terraform plugin testing framework, always set to `gtm_asmap`.",
				MarkdownDescription: "Identifier of the data source.",
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "A descriptive label for the AS map.",
			},
			"map_name": schema.StringAttribute{
				Required:    true,
				Description: "A descriptive label for the AS map",
			},
		},
		Blocks: asmapBlock,
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *asmapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM AS map DataSource Read")

	var data asmapDataSourceModel

	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	asMap, err := client.GetAsMap(ctx, data.Name.ValueString(), data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("fetching GTM ASmap failed: ", err.Error())
		return
	}

	data.setAttributes(asMap)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (m *asmapDataSourceModel) setAttributes(asmap *gtm.AsMap) {
	m.Name = types.StringValue(asmap.Name)
	m.setDefaultDatacenter(asmap.DefaultDatacenter)
	m.setAssignments(asmap.Assignments)
	m.setLinks(asmap.Links)
	m.ID = types.StringValue("gtm_asmap")

}

func (m *asmapDataSourceModel) setDefaultDatacenter(d *gtm.DatacenterBase) {
	m.DefaultDatacenter = &defaultDatacenter{
		DatacenterID: types.Int64Value(int64(d.DatacenterId)),
		Nickname:     types.StringValue(d.Nickname),
	}
}

func (m *asmapDataSourceModel) setAssignments(assignments []*gtm.AsAssignment) {
	toBasetypesInt64Slice := func(n []int64) []basetypes.Int64Value {
		out := make([]basetypes.Int64Value, 0, len(n))
		for _, number := range n {
			out = append(out, types.Int64Value(number))
		}
		return out
	}

	for _, a := range assignments {
		assignmentObject := asMapAssignment{
			DatacenterID: types.Int64Value(int64(a.DatacenterId)),
			Nickname:     types.StringValue(a.Nickname),
			ASNumbers:    toBasetypesInt64Slice(a.AsNumbers),
		}
		m.Assignments = append(m.Assignments, assignmentObject)
	}
}

func (m *asmapDataSourceModel) setLinks(links []*gtm.Link) {
	for _, l := range links {
		linkObject := link{
			Rel:  types.StringValue(l.Rel),
			Href: types.StringValue(l.Href),
		}

		m.Links = append(m.Links, linkObject)
	}
}
