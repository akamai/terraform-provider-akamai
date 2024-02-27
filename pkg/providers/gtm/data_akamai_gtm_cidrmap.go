package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type cidrmapDataSource struct {
	meta meta.Meta
}

type cidrmapDataSourceModel struct {
	ID                types.String        `tfsdk:"id"`
	Domain            types.String        `tfsdk:"domain"`
	Name              types.String        `tfsdk:"map_name"`
	DefaultDatacenter *defaultDatacenter  `tfsdk:"default_datacenter"`
	Assignments       []cidrMapAssignment `tfsdk:"assignments"`
	Links             []link              `tfsdk:"links"`
}

var (
	_ datasource.DataSource              = &cidrmapDataSource{}
	_ datasource.DataSourceWithConfigure = &cidrmapDataSource{}
)

// NewGTMCidrmapDataSource returns a new GTM cidrmap data source
func NewGTMCidrmapDataSource() datasource.DataSource { return &cidrmapDataSource{} }

func (d *cidrmapDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_gtm_cidrmap"
}

func (d *cidrmapDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

var (
	cidrmapBlocks = map[string]schema.Block{
		"assignments": schema.ListNestedBlock{
			Description: "Contains information about the CIDR zone groupings of CIDR blocks.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"datacenter_id": schema.Int64Attribute{
						Computed:    true,
						Description: "A unique identifier for an existing data center in the domain.",
					},
					"blocks": schema.SetAttribute{
						Computed:    true,
						Description: "Specifies an array of CIDR blocks.",
						ElementType: types.StringType,
					},
					"nickname": schema.StringAttribute{
						Computed:    true,
						Description: "A descriptive label for the CIDR zone group.",
					},
				},
			},
		},
		"default_datacenter": schema.SingleNestedBlock{
			Description: "A placeholder for all other CIDR zones, CIDR blocks not found in these CIDR zones.",
			Attributes: map[string]schema.Attribute{
				"datacenter_id": schema.Int64Attribute{
					Computed:    true,
					Description: "For each property, an identifier for all other CIDR zones' CNAME.",
				},
				"nickname": schema.StringAttribute{
					Computed:    true,
					Description: "A descriptive label for all other CIDR blocks.",
				},
			},
		},
		"links": schema.SetNestedBlock{
			Description: "Specifies the URL path that allows direct navigation to the CIDR map.",
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

func (d *cidrmapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "GTM CIDR map data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source.",
				DeprecationMessage:  "Required by the terraform plugin testing framework, always set to `gtm_cidrmap`.",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				Required:    true,
				Description: "GTM domain name.",
			},
			"map_name": schema.StringAttribute{
				Required:    true,
				Description: "A descriptive label for the CIDR map.",
			},
		},
		Blocks: cidrmapBlocks,
	}
}

func (d *cidrmapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Cidrmap DataSource Read")

	var data cidrmapDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	cidrMap, err := client.GetCidrMap(ctx, data.Name.ValueString(), data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("fetching GTM CIDRmap failed: ", err.Error())
		return
	}

	diags := data.setAttributes(ctx, cidrMap)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *cidrmapDataSourceModel) setAttributes(ctx context.Context, cidrmap *gtm.CidrMap) diag.Diagnostics {
	m.Name = types.StringValue(cidrmap.Name)
	m.setDefaultDatacenter(cidrmap.DefaultDatacenter)
	m.setLinks(cidrmap.Links)
	diags := m.setAssignments(ctx, cidrmap.Assignments)
	if diags.HasError() {
		return diags
	}
	m.ID = types.StringValue("gtm_cidrmap")

	return nil
}

func (m *cidrmapDataSourceModel) setDefaultDatacenter(b *gtm.DatacenterBase) {
	m.DefaultDatacenter = &defaultDatacenter{
		DatacenterID: types.Int64Value(int64(b.DatacenterId)),
		Nickname:     types.StringValue(b.Nickname),
	}
}

func (m *cidrmapDataSourceModel) setLinks(links []*gtm.Link) {
	for _, l := range links {
		linkObj := link{
			Rel:  types.StringValue(l.Rel),
			Href: types.StringValue(l.Href),
		}

		m.Links = append(m.Links, linkObj)
	}
}

func (m *cidrmapDataSourceModel) setAssignments(ctx context.Context, assignments []*gtm.CidrAssignment) diag.Diagnostics {
	for _, a := range assignments {
		blocks, diags := types.SetValueFrom(ctx, types.StringType, a.Blocks)
		if diags.HasError() {
			return diags
		}
		assignmentObj := cidrMapAssignment{
			DatacenterID: types.Int64Value(int64(a.DatacenterId)),
			Nickname:     types.StringValue(a.Nickname),
			Blocks:       blocks,
		}

		m.Assignments = append(m.Assignments, assignmentObj)
	}
	return nil
}
