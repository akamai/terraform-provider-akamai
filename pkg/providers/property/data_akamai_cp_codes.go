package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &cpCodesDataSource{}
	_ datasource.DataSourceWithConfigure = &cpCodesDataSource{}
)

type (
	cpCodesDataSource struct {
		meta meta.Meta
	}
	// cpCodesDataSourceModel provide whole model of CP Codes datasource
	cpCodesDataSourceModel struct {
		ContractID        types.String  `tfsdk:"contract_id"`
		GroupID           types.String  `tfsdk:"group_id"`
		AccountID         types.String  `tfsdk:"account_id"`
		FilterByName      types.String  `tfsdk:"filter_by_name"`
		FilterByProductID types.String  `tfsdk:"filter_by_product_id"`
		CpCodes           []cpCodeModel `tfsdk:"cp_codes"`
	}

	// cpCodeModel provides the model of one specific CP code
	cpCodeModel struct {
		Name        types.String `tfsdk:"name"`
		CPCodeID    types.String `tfsdk:"cp_code_id"`
		CreatedDate types.String `tfsdk:"created_date"`
		ProductIDs  types.List   `tfsdk:"product_ids"`
	}
)

// NewCPCodesDataSource returns a new CP codes data source
func NewCPCodesDataSource() datasource.DataSource {
	return &cpCodesDataSource{}
}

func (c *cpCodesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	c.meta = meta.Must(req.ProviderData)
}

func (c *cpCodesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_cp_codes"
}

func (c *cpCodesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "CP codes schema",
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the group.",
			},
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the contract.",
			},
			"account_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifies the prevailing account under which you requested the data.",
			},
			"filter_by_name": schema.StringAttribute{
				Optional:    true,
				Description: "Allows you to filter CP codes by a specific CP code's name.",
			},
			"filter_by_product_id": schema.StringAttribute{
				Optional:    true,
				Description: "Allows you to filter CP codes by a specific product ID.",
			},

			"cp_codes": schema.ListNestedAttribute{
				Description: "CP codes list.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cp_code_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of a specific CP code.",
						},
						"name": schema.StringAttribute{
							Description: "Name of the CP code.",
							Computed:    true,
						},
						"created_date": schema.StringAttribute{
							Description: "The date and time when the CP code was created.",
							Computed:    true,
						},
						"product_ids": schema.ListAttribute{
							Computed:    true,
							Description: "A list of of product IDs for a given CP code.",
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (c *cpCodesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "CPCodesDataSource Read")
	var diags diag.Diagnostics
	var data cpCodesDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
	client := Client(c.meta)
	cpCodes, err := client.GetCPCodes(ctx, papi.GetCPCodesRequest{
		ContractID: data.ContractID.ValueString(),
		GroupID:    data.GroupID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("fetching cp codes failed", err.Error())
		return
	}

	data.GroupID = types.StringValue(strings.TrimPrefix(cpCodes.GroupID, "grp_"))
	data.ContractID = types.StringValue(strings.TrimPrefix(cpCodes.ContractID, "ctr_"))
	data.AccountID = types.StringValue(strings.TrimPrefix(cpCodes.AccountID, "act_"))

	filteredCodes := filterCPCodes(data.FilterByProductID, data.FilterByName, cpCodes.CPCodes.Items)
	data.CpCodes, diags = mapCodesToModel(ctx, filteredCodes)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
func mapCodesToModel(ctx context.Context, cpCodes []papi.CPCode) ([]cpCodeModel, diag.Diagnostics) {
	var codes []cpCodeModel
	for _, cpCodeItem := range cpCodes {
		var cpCodeModel cpCodeModel

		cpCodeModel.CPCodeID = types.StringValue(strings.TrimPrefix(cpCodeItem.ID, "cpc_"))
		cpCodeModel.Name = types.StringValue(cpCodeItem.Name)
		cpCodeModel.CreatedDate = types.StringValue(cpCodeItem.CreatedDate)
		products, diags := types.ListValueFrom(ctx, types.StringType, cpCodeItem.ProductIDs)
		if diags.HasError() {
			return nil, diags
		}
		cpCodeModel.ProductIDs = products
		codes = append(codes, cpCodeModel)
	}
	return codes, nil
}

func filterCPCodes(productID, name types.String, cpCodes []papi.CPCode) []papi.CPCode {
	var filteredCodes []papi.CPCode
	for _, cpCodeItem := range cpCodes {
		if !name.IsNull() && !name.IsUnknown() && name.ValueString() != cpCodeItem.Name {
			continue
		}
		if !productID.IsNull() && !productID.IsUnknown() && !containsWithoutPrefix(cpCodeItem.ProductIDs, productID.ValueString(), "cpc_") {
			continue
		}
		filteredCodes = append(filteredCodes, cpCodeItem)
	}
	return filteredCodes
}

func containsWithoutPrefix(list []string, element, prefix string) bool {
	for _, listElement := range list {
		if strings.EqualFold(strings.TrimPrefix(element, prefix), strings.TrimPrefix(listElement, prefix)) {
			return true
		}
	}
	return false
}
