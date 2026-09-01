package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &edgeHostnamesDataSource{}
	_ datasource.DataSourceWithConfigure = &edgeHostnamesDataSource{}
)

type (
	edgeHostnamesDataSource struct {
		meta.DataSource
	}

	edgeHostnamesDataSourceModel struct {
		ContractID    types.String        `tfsdk:"contract_id"`
		GroupID       types.String        `tfsdk:"group_id"`
		Options       types.List          `tfsdk:"options"`
		AccountID     types.String        `tfsdk:"account_id"`
		EdgeHostnames []edgeHostnameModel `tfsdk:"edge_hostnames"`
	}

	edgeHostnameModel struct {
		DomainPrefix        types.String    `tfsdk:"domain_prefix"`
		DomainSuffix        types.String    `tfsdk:"domain_suffix"`
		EdgeHostnameDomain  types.String    `tfsdk:"edge_hostname_domain"`
		EdgeHostnameID      types.String    `tfsdk:"edge_hostname_id"`
		IPVersionBehaviour  types.String    `tfsdk:"ip_version_behaviour"`
		ProductID           types.String    `tfsdk:"product_id"`
		Secure              types.Bool      `tfsdk:"secure"`
		Status              types.String    `tfsdk:"status"`
		UseCases            []useCasesModel `tfsdk:"use_cases"`
		HTTPSServiceBinding types.String    `tfsdk:"https_service_binding"`
	}

	useCasesModel struct {
		Option  types.String `tfsdk:"option"`
		Type    types.String `tfsdk:"type"`
		UseCase types.String `tfsdk:"use_case"`
	}
)

// NewEdgeHostnamesDataSource returns a new Edge Hostnames data source
func NewEdgeHostnamesDataSource() datasource.DataSource {
	return &edgeHostnamesDataSource{}
}

// Metadata configures data source's meta information.
func (d *edgeHostnamesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_edge_hostnames"
}

// Schema is used to define data source's terraform schema.
func (d *edgeHostnamesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve a list of edge hostnames for a given contract and group.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the contract.",
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the group.",
			},
			"options": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Optional query parameters that enables extra mapping-related information.",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifies the prevailing account under which you requested the data.",
			},
			"edge_hostnames": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The set of requested edge hostnames",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"domain_prefix": schema.StringAttribute{
							Computed:    true,
							Description: "The origin domain portion of the edge hostname.",
						},
						"domain_suffix": schema.StringAttribute{
							Computed:    true,
							Description: "The Akamai-specific portion of the edge hostname.",
						},
						"edge_hostname_domain": schema.StringAttribute{
							Computed:    true,
							Description: "The full edge domain name formed from the domainPrefix and domainSuffix.",
						},
						"edge_hostname_id": schema.StringAttribute{
							Computed:    true,
							Description: "The edge hostname's unique identifier.",
						},
						"ip_version_behaviour": schema.StringAttribute{
							Computed:    true,
							Description: "IP version behavior of the edge hostname.",
						},
						"product_id": schema.StringAttribute{
							Computed:    true,
							Description: "The product you created the edge hostname for.",
						},
						"secure": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to use the edge hostname with SSL.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Status of the edge hostname.",
						},
						"use_cases": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Available use cases for edge hostnames assigned to the product.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"option": schema.StringAttribute{
										Computed:    true,
										Description: "Specifies one of the available options available in the response object.",
									},
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "Identifies the type of network over which traffic deploys.",
									},
									"use_case": schema.StringAttribute{
										Computed:    true,
										Description: "Identifies each use case.",
									},
								},
							},
						},
						"https_service_binding": schema.StringAttribute{
							Computed:    true,
							Description: "HTTPS service binding of the edge hostname.",
						},
					},
				}},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *edgeHostnamesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "PAPI Edge Hostnames DataSource Read")

	var data edgeHostnamesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var options []string
	if tf.IsKnown(data.Options) {
		resp.Diagnostics.Append(data.Options.ElementsAs(ctx, &options, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, "Calling GetEdgeHostnames", map[string]any{
		"contract_id": data.ContractID.ValueString(),
		"group_id":    data.GroupID.ValueString(),
		"options":     options,
	})

	edgehostnames, err := d.Client.GetPAPI().GetEdgeHostnames(ctx, papi.GetEdgeHostnamesRequest{
		ContractID: data.ContractID.ValueString(),
		GroupID:    data.GroupID.ValueString(),
		Options:    options,
	})
	if err != nil {
		resp.Diagnostics.AddError("List Edge Hostnames failed", err.Error())
		return
	}

	data.convertEdgeHostnamesToModel(edgehostnames)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *edgeHostnamesDataSourceModel) convertEdgeHostnamesToModel(edgehostnames *papi.GetEdgeHostnamesResponse) {
	m.AccountID = types.StringValue(edgehostnames.AccountID)

	for _, eh := range edgehostnames.EdgeHostnames.Items {
		var useCases []useCasesModel
		for _, uc := range eh.UseCases {
			useCases = append(useCases, useCasesModel{
				Option:  types.StringValue(uc.Option),
				Type:    types.StringValue(uc.Type),
				UseCase: types.StringValue(uc.UseCase),
			})
		}

		productID := types.StringValue(eh.ProductID)
		if eh.ProductID == "" {
			productID = types.StringNull()
		}

		m.EdgeHostnames = append(m.EdgeHostnames, edgeHostnameModel{
			DomainPrefix:        types.StringValue(eh.DomainPrefix),
			DomainSuffix:        types.StringValue(eh.DomainSuffix),
			EdgeHostnameDomain:  types.StringValue(eh.Domain),
			EdgeHostnameID:      types.StringValue(eh.ID),
			IPVersionBehaviour:  types.StringValue(eh.IPVersionBehavior),
			ProductID:           productID,
			Secure:              types.BoolValue(eh.Secure),
			Status:              types.StringValue(eh.Status),
			UseCases:            useCases,
			HTTPSServiceBinding: types.StringPointerValue(eh.HTTPSServiceBinding),
		})
	}
}
