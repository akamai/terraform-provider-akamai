package reportinggroups

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &cpCodesDataSource{}
	_ datasource.DataSourceWithConfigure = &cpCodesDataSource{}
)

type (
	cpCodesDataSource struct {
		meta.DataSource
	}

	cpCodesDataSourceModel struct {
		ContractID types.String            `tfsdk:"contract_id"`
		GroupID    types.String            `tfsdk:"group_id"`
		ProductID  types.String            `tfsdk:"product_id"`
		CPCodeName types.String            `tfsdk:"cp_code_name"`
		CPCodes    []cpCodeDataSourceModel `tfsdk:"cp_codes"`
	}
)

// NewCPCodesDataSource returns a new reporting groups CP codes data source.
func NewCPCodesDataSource() datasource.DataSource {
	return &cpCodesDataSource{}
}

func (d *cpCodesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_reportinggroups_cp_codes"
}

// Schema is used to define data source's terraform schema.
func (d *cpCodesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists detailed information about CP codes available within account and contract.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies the contract to filter data by.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies the access group to filter data by.",
			},
			"product_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies the product or service to filter data by.",
			},
			"cp_code_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the CP code to filter data by.",
			},
			"cp_codes": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A collection of CP codes available for your contract.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cp_code_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique numeric identifier of the CP code.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the CP code.",
						},
						"purgeable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the CP code can be used to purge content.",
						},
						"account_id": schema.StringAttribute{
							Computed:    true,
							Description: "The account identifier associated with the CP code.",
						},
						"default_time_zone": schema.StringAttribute{
							Computed:    true,
							Description: "The default time zone of the CP code.",
						},
						"override_time_zone": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The override time zone of the CP code.",
							Attributes: map[string]schema.Attribute{
								"time_zone_id": schema.StringAttribute{
									Computed:    true,
									Description: "The time zone identifier.",
								},
								"time_zone_value": schema.StringAttribute{
									Computed:    true,
									Description: "The time zone value.",
								},
							},
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the CP code.",
						},
						"contracts": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of contracts associated with the CP code.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"contract_id": schema.StringAttribute{
										Computed:    true,
										Description: "The contract identifier.",
									},
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "The contract status.",
									},
								},
							},
						},
						"products": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of products associated with the CP code.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"product_id": schema.StringAttribute{
										Computed:    true,
										Description: "The product identifier.",
									},
									"product_name": schema.StringAttribute{
										Computed:    true,
										Description: "The product name.",
									},
								},
							},
						},
						"access_group": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "The access control group associated with the CP code.",
							Attributes: map[string]schema.Attribute{
								"contract_id": schema.StringAttribute{
									Computed:    true,
									Description: "The contract identifier assigned to the access control group.",
								},
								"group_id": schema.Int64Attribute{
									Computed:    true,
									Description: "The access control group identifier.",
								},
							},
						},
					},
				},
			},
		}}
}

func (d *cpCodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reporting Groups CPCodes DataSource Read")

	var data cpCodesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Calling ListCPCodes", map[string]any{
		"contract_id":  data.ContractID.ValueString(),
		"group_id":     data.GroupID.ValueString(),
		"product_id":   data.ProductID.ValueString(),
		"cp_code_name": data.CPCodeName.ValueString(),
	})

	cpCodesDetail, err := d.Client.GetReportingGroups().ListCPCodes(ctx, reportinggroups.ListCPCodesRequest{
		ContractID: data.ContractID.ValueString(),
		GroupID:    data.GroupID.ValueString(),
		ProductID:  data.ProductID.ValueString(),
		CPCodeName: data.CPCodeName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to retrieve CP codes detail", err.Error())
		return
	}

	data.populateFromListCPCodes(cpCodesDetail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *cpCodesDataSourceModel) populateFromListCPCodes(cpcodes *reportinggroups.ListCPCodesResponse) {
	m.CPCodes = make([]cpCodeDataSourceModel, 0, len(cpcodes.CPCodes))
	for _, cpCode := range cpcodes.CPCodes {
		contracts := make([]cpCodeContractModel, len(cpCode.Contracts))
		for i, contract := range cpCode.Contracts {
			contracts[i] = cpCodeContractModel{
				ContractID: types.StringValue(contract.ContractID),
				Status:     types.StringValue(contract.Status),
			}
		}

		products := make([]cpCodeProductModel, len(cpCode.Products))
		for i, product := range cpCode.Products {
			products[i] = cpCodeProductModel{
				ProductID:   types.StringValue(product.ProductID),
				ProductName: types.StringValue(product.ProductName),
			}
		}

		m.CPCodes = append(m.CPCodes, cpCodeDataSourceModel{
			CPCodeID:        types.Int64Value(cpCode.CPCodeID),
			Name:            types.StringValue(cpCode.CPCodeName),
			Purgeable:       types.BoolValue(cpCode.Purgeable),
			AccountID:       types.StringValue(cpCode.AccountID),
			DefaultTimeZone: types.StringValue(cpCode.DefaultTimeZone),
			OverrideTimeZone: &cpCodeTimeZoneModel{
				TimeZoneID:    types.StringValue(cpCode.OverrideTimeZone.TimeZoneID),
				TimeZoneValue: types.StringValue(cpCode.OverrideTimeZone.TimeZoneValue),
			},
			Type:      types.StringValue(cpCode.Type),
			Contracts: contracts,
			Products:  products,
			AccessGroup: &cpCodeAccessGroupModel{
				ContractID: types.StringValue(cpCode.AccessGroup.ContractID),
				GroupID:    types.Int64PointerValue(cpCode.AccessGroup.GroupID),
			},
		})
	}
}
