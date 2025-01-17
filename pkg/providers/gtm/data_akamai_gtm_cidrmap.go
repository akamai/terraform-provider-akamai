package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type cidrMapDataSource struct {
	meta meta.Meta
}

type cidrMapDataSourceModel struct {
	ID                types.String        `tfsdk:"id"`
	Domain            types.String        `tfsdk:"domain"`
	Name              types.String        `tfsdk:"map_name"`
	DefaultDatacenter *defaultDatacenter  `tfsdk:"default_datacenter"`
	Assignments       []cidrMapAssignment `tfsdk:"assignments"`
	Links             []link              `tfsdk:"links"`
}

var (
	_ datasource.DataSource              = &cidrMapDataSource{}
	_ datasource.DataSourceWithConfigure = &cidrMapDataSource{}
)

// NewGTMCIDRMapDataSource returns a new GTM CIDRMap data source
func NewGTMCIDRMapDataSource() datasource.DataSource { return &cidrMapDataSource{} }

func (d *cidrMapDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_gtm_cidrmap"
}

func (d *cidrMapDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	cidrMapBlocks = map[string]schema.Block{
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

func (d *cidrMapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
		Blocks: cidrMapBlocks,
	}
}

func (d *cidrMapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "GTM Cidrmap DataSource Read")

	var data cidrMapDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := Client(d.meta)
	cidrMap, err := client.GetCIDRMap(ctx, gtm.GetCIDRMapRequest{
		DomainName: data.Domain.ValueString(),
		MapName:    data.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching GTM CIDRMap failed: ", err.Error())
		return
	}

	diags := data.setAttributes(ctx, cidrMap)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *cidrMapDataSourceModel) setAttributes(ctx context.Context, cidrMap *gtm.GetCIDRMapResponse) diag.Diagnostics {
	m.Name = types.StringValue(cidrMap.Name)
	m.DefaultDatacenter = newDefaultDatacenter(*cidrMap.DefaultDatacenter)
	m.Links = getLinks(cidrMap.Links)
	diags := m.setAssignments(ctx, cidrMap.Assignments)
	if diags.HasError() {
		return diags
	}
	m.ID = types.StringValue("gtm_cidrmap")

	return nil
}

func (m *cidrMapDataSourceModel) setAssignments(ctx context.Context, assignments []gtm.CIDRAssignment) diag.Diagnostics {
	for _, a := range assignments {
		blocks, diags := types.SetValueFrom(ctx, types.StringType, a.Blocks)
		if diags.HasError() {
			return diags
		}
		assignmentObj := cidrMapAssignment{
			DatacenterID: types.Int64Value(int64(a.DatacenterID)),
			Nickname:     types.StringValue(a.Nickname),
			Blocks:       blocks,
		}

		m.Assignments = append(m.Assignments, assignmentObj)
	}
	return nil
}
