package reportinggroups

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/convert"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &cpCodeDataSource{}
	_ datasource.DataSourceWithConfigure = &cpCodeDataSource{}
)

type (
	cpCodeDataSource struct {
		meta.DataSource
	}

	cpCodeDataSourceModel struct {
		CPCodeID         types.Int64             `tfsdk:"cp_code_id"`
		Name             types.String            `tfsdk:"name"`
		Purgeable        types.Bool              `tfsdk:"purgeable"`
		AccountID        types.String            `tfsdk:"account_id"`
		DefaultTimeZone  types.String            `tfsdk:"default_time_zone"`
		OverrideTimeZone *cpCodeTimeZoneModel    `tfsdk:"override_time_zone"`
		Type             types.String            `tfsdk:"type"`
		Contracts        []cpCodeContractModel   `tfsdk:"contracts"`
		Products         []cpCodeProductModel    `tfsdk:"products"`
		AccessGroup      *cpCodeAccessGroupModel `tfsdk:"access_group"`
	}

	cpCodeTimeZoneModel struct {
		TimeZoneID    types.String `tfsdk:"time_zone_id"`
		TimeZoneValue types.String `tfsdk:"time_zone_value"`
	}

	cpCodeContractModel struct {
		ContractID types.String `tfsdk:"contract_id"`
		Status     types.String `tfsdk:"status"`
	}

	cpCodeProductModel struct {
		ProductID   types.String `tfsdk:"product_id"`
		ProductName types.String `tfsdk:"product_name"`
	}

	cpCodeAccessGroupModel struct {
		ContractID types.String `tfsdk:"contract_id"`
		GroupID    types.String `tfsdk:"group_id"`
	}
)

// NewCPCodeDataSource returns a new reporting groups CP code data source.
func NewCPCodeDataSource() datasource.DataSource {
	return &cpCodeDataSource{}
}

func (d *cpCodeDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_reportinggroups_cp_code"
}

// Schema is used to define data source's terraform schema.
func (d *cpCodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve detailed information about a specific CP code.",
		Attributes: map[string]schema.Attribute{
			"cp_code_id": schema.Int64Attribute{
				Required:    true,
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
					"group_id": schema.StringAttribute{
						Computed:    true,
						Description: "The access control group identifier.",
					},
				},
			},
		},
	}
}

func (d *cpCodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reporting Groups CPCode DataSource Read")

	var data cpCodeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cpCodeDetail, err := d.Client.GetReportingGroups().GetCPCode(ctx, reportinggroups.GetCPCodeRequest{CPCodeID: data.CPCodeID.ValueInt64()})
	if err != nil {
		resp.Diagnostics.AddError("Failed to retrieve CP code detail", err.Error())
		return
	}

	data.populateFromCPCodeDetail(cpCodeDetail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (m *cpCodeDataSourceModel) populateFromCPCodeDetail(detail *reportinggroups.GetCPCodeResponse) {
	m.Name = types.StringValue(detail.CPCodeName)
	m.Purgeable = types.BoolValue(detail.Purgeable)
	m.AccountID = types.StringValue(detail.AccountID)
	m.DefaultTimeZone = types.StringValue(detail.DefaultTimeZone)
	m.Type = types.StringValue(detail.Type)

	m.OverrideTimeZone = &cpCodeTimeZoneModel{
		TimeZoneID:    types.StringValue(detail.OverrideTimeZone.TimeZoneID),
		TimeZoneValue: types.StringValue(detail.OverrideTimeZone.TimeZoneValue),
	}

	m.Contracts = make([]cpCodeContractModel, len(detail.Contracts))
	for i, contract := range detail.Contracts {
		m.Contracts[i] = cpCodeContractModel{
			ContractID: types.StringValue(contract.ContractID),
			Status:     types.StringValue(contract.Status),
		}
	}

	m.Products = make([]cpCodeProductModel, len(detail.Products))
	for i, product := range detail.Products {
		m.Products[i] = cpCodeProductModel{
			ProductID:   types.StringValue(product.ProductID),
			ProductName: types.StringValue(product.ProductName),
		}
	}

	m.AccessGroup = &cpCodeAccessGroupModel{
		ContractID: types.StringValue(detail.AccessGroup.ContractID),
		GroupID:    convert.Int64PtrToStringValue(detail.AccessGroup.GroupID),
	}
}
