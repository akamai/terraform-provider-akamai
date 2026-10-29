package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &edgeHostnameDataSource{}
	_ datasource.DataSourceWithConfigure = &edgeHostnameDataSource{}
)

type (
	edgeHostnameDataSource struct {
		meta.DataSource
	}

	edgeHostnameDataSourceModel struct {
		EdgeHostnameID types.String `tfsdk:"id"`
		ContractID     types.String `tfsdk:"contract_id"`
		GroupID        types.String `tfsdk:"group_id"`
		Options        types.String `tfsdk:"options"`
		AccountID      types.String `tfsdk:"account_id"`
		EdgeHostname   types.Object `tfsdk:"edge_hostname"`
	}
)

// NewEdgeHostnameDataSource returns a new Edge Hostname data source.
func NewEdgeHostnameDataSource() datasource.DataSource {
	return &edgeHostnameDataSource{}
}

// Metadata configures data source's meta information.
func (d *edgeHostnameDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_edge_hostname"
}

// Schema is used to define data source's terraform schema.
func (d *edgeHostnameDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve an edge hostname for a given contract, group, and edge hostname ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the edge hostname.",
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the contract.",
			},
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the group.",
			},
			"options": schema.StringAttribute{
				Optional:    true,
				Description: "Optional query parameter that enables extra mapping-related information.",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifies the prevailing account under which you requested the data.",
			},
			"edge_hostname": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Requested edge hostname",
				Attributes: map[string]schema.Attribute{
					"delivery_policy": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Delivery policy of the edge hostname.",
						Attributes: map[string]schema.Attribute{
							"policy_type": schema.StringAttribute{
								Computed:    true,
								Description: "Type of the delivery policy.",
							},
							"policy_attributes": schema.ListNestedAttribute{
								Computed:    true,
								Description: "Attributes of the delivery policy.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"attribute_type": schema.StringAttribute{
											Computed:    true,
											Description: "Type of the policy attribute.",
										},
										"code": schema.StringAttribute{
											Computed:    true,
											Description: "The code of the geographical region",
										},
										"name": schema.StringAttribute{
											Computed:    true,
											Description: "Name of the policy attribute.",
										},
									},
								},
							},
						},
					},
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
					"map_details_serial_number": schema.Int64Attribute{
						Computed:    true,
						Description: "Serial number of the map details.",
					},
					"map_details_slot_number": schema.Int64Attribute{
						Computed:    true,
						Description: "Slot number of the map details.",
					},
					"map_details_map_domain": schema.StringAttribute{
						Computed:    true,
						Description: "Map domain of the map details.",
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
			},
		},
	}
}

// Read is called to retrieve data for the resource and populate the Terraform state with it.
func (d *edgeHostnameDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "PAPI Edge Hostname DataSource Read")

	var data edgeHostnameDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var options string
	if tf.IsKnown(data.Options) {
		options = data.Options.ValueString()
	}

	tflog.Debug(ctx, "Calling GetEdgeHostname", map[string]any{
		"id":          data.EdgeHostnameID.ValueString(),
		"contract_id": data.ContractID.ValueString(),
		"group_id":    data.GroupID.ValueString(),
		"options":     options,
	})

	edgeHostname, err := d.Client.GetPAPI().GetEdgeHostname(ctx, papi.GetEdgeHostnameRequest{
		EdgeHostnameID: data.EdgeHostnameID.ValueString(),
		ContractID:     data.ContractID.ValueString(),
		GroupID:        data.GroupID.ValueString(),
		Options:        options,
	})
	if err != nil {
		resp.Diagnostics.AddError("Get Edge Hostname failed", err.Error())
		return
	}

	if resp.Diagnostics.Append(data.populateEdgeHostnameModel(ctx, edgeHostname)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *edgeHostnameDataSourceModel) populateEdgeHostnameModel(ctx context.Context, edgeHostname *papi.GetEdgeHostnamesResponse) diag.Diagnostics {
	m.AccountID = types.StringValue(edgeHostname.AccountID)

	eh := edgeHostname.EdgeHostnames.Items[0]

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

	var deliveryPolicy *deliverPolicyModel

	if eh.DeliveryPolicy != nil {
		var attrs []policyAttributesModel
		for _, pa := range eh.DeliveryPolicy.PolicyAttributes {
			attrs = append(attrs, policyAttributesModel{
				AttributeType: types.StringValue(pa.AttributeType),
				Code:          types.StringValue(pa.Code),
				Name:          types.StringValue(pa.Name),
			})
		}

		deliveryPolicy = &deliverPolicyModel{
			PolicyType:       types.StringValue(eh.DeliveryPolicy.PolicyType),
			PolicyAttributes: attrs,
		}
	}

	edgeHostnameModelValue := edgeHostnameModel{
		DeliveryPolicy:         deliveryPolicy,
		DomainPrefix:           types.StringValue(eh.DomainPrefix),
		DomainSuffix:           types.StringValue(eh.DomainSuffix),
		EdgeHostnameDomain:     types.StringValue(eh.Domain),
		EdgeHostnameID:         types.StringValue(eh.ID),
		IPVersionBehaviour:     types.StringValue(eh.IPVersionBehavior),
		ProductID:              productID,
		Secure:                 types.BoolValue(eh.Secure),
		Status:                 types.StringValue(eh.Status),
		MapDetailsSerialNumber: types.Int64PointerValue(eh.MapDetailsSerialNumber),
		MapDetailsSlotNumber:   types.Int64PointerValue(eh.MapDetailsSlotNumber),
		MapDetailsMapDomain:    types.StringPointerValue(eh.MapDetailsMapDomain),
		UseCases:               useCases,
		HTTPSServiceBinding:    types.StringPointerValue(eh.HTTPSServiceBinding),
	}

	obj, diags := types.ObjectValueFrom(ctx, edgeHostnameModel{}.attrTypes(), edgeHostnameModelValue)
	if diags.HasError() {
		return diags
	}

	m.EdgeHostname = obj
	return nil
}

func (e edgeHostnameModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"delivery_policy":           types.ObjectType{AttrTypes: deliverPolicyModel{}.attrTypes()},
		"domain_prefix":             types.StringType,
		"domain_suffix":             types.StringType,
		"edge_hostname_domain":      types.StringType,
		"edge_hostname_id":          types.StringType,
		"ip_version_behaviour":      types.StringType,
		"product_id":                types.StringType,
		"secure":                    types.BoolType,
		"status":                    types.StringType,
		"map_details_serial_number": types.Int64Type,
		"map_details_slot_number":   types.Int64Type,
		"map_details_map_domain":    types.StringType,
		"use_cases":                 types.ListType{ElemType: types.ObjectType{AttrTypes: useCasesModel{}.attrTypes()}},
		"https_service_binding":     types.StringType,
	}
}

func (u useCasesModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"option":   types.StringType,
		"type":     types.StringType,
		"use_case": types.StringType,
	}
}

func (d deliverPolicyModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"policy_type":       types.StringType,
		"policy_attributes": types.ListType{ElemType: types.ObjectType{AttrTypes: policyAttributesModel{}.attrTypes()}},
	}
}

func (p policyAttributesModel) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"attribute_type": types.StringType,
		"code":           types.StringType,
		"name":           types.StringType,
	}
}
