package datastream

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v10/internal/customtypes"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &answerXServiceIDsDataSource{}
	_ datasource.DataSourceWithConfigure = &answerXServiceIDsDataSource{}
)

// answerXServiceIDItemAttrTypes defines the attribute types for a single service ID element.
var answerXServiceIDItemAttrTypes = map[string]attr.Type{
	"id":      types.Int64Type,
	"name":    types.StringType,
	"product": types.StringType,
}

type (
	answerXServiceIDsDataSource struct {
		meta   meta.Meta
		client datastream.DS
	}

	answerXServiceIDsDataSourceModel struct {
		ContractID customtypes.IgnorePrefixValue `tfsdk:"contract_id"`
		ServiceIDs types.Set                     `tfsdk:"service_ids"`
	}

	answerXServiceIDItemModel struct {
		ID      types.Int64  `tfsdk:"id"`
		Name    types.String `tfsdk:"name"`
		Product types.String `tfsdk:"product"`
	}
)

// NewAnswerXServiceIDsDataSource returns a new AnswerX service IDs data source.
func NewAnswerXServiceIDsDataSource() datasource.DataSource {
	return &answerXServiceIDsDataSource{}
}

func (d *answerXServiceIDsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_datastream_answerx_service_ids"
}

func (d *answerXServiceIDsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	m, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.meta = m
	d.client = inst.Client(m)
}

// Schema defines the Terraform schema for this data source.
func (d *answerXServiceIDsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of AnswerX service IDs available in the given contract.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifies the contract for which to list available AnswerX service IDs.",
				CustomType:  customtypes.IgnorePrefixType{Prefix: "ctr_"},
				Validators:  []validator.String{validators.NotEmptyString()},
			},
			"service_ids": schema.SetNestedAttribute{
				Computed:    true,
				Description: "Set of AnswerX service IDs available in the given contract.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "Service ID monitored in the stream.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the service ID.",
						},
						"product": schema.StringAttribute{
							Computed:    true,
							Description: "The product associated with the service ID.",
						},
					},
				},
			},
		},
	}
}

const answerXServiceIDsPageSize int64 = 5000

// Read is called when the provider must read data source values in order to update state.
func (d *answerXServiceIDsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "DataStream AnswerXServiceIDs DataSource Read")

	var data answerXServiceIDsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	ctx = session.ContextWithOptions(ctx, session.WithContextLog(d.meta.Log("datastream", "answerXServiceIDsDataSourceRead")))

	contractID := strings.TrimPrefix(data.ContractID.ValueString(), "ctr_")
	if contractID == "" {
		resp.Diagnostics.AddError("Invalid contract ID", "contract_id must not be blank")
		return
	}

	var allServiceIDs []datastream.AnswerXServiceDetail
	var page int64 = 1
	for {
		response, err := d.client.ListAnswerXServiceIDs(ctx, datastream.ListAnswerXServiceIDsRequest{
			ContractID: contractID,
			Page:       page,
			PageSize:   answerXServiceIDsPageSize,
		})
		if err != nil {
			resp.Diagnostics.AddError("Listing AnswerX service IDs failed", err.Error())
			return
		}

		if response == nil || response.Metadata == nil || response.Metadata.Page != page || response.Metadata.LastPage < page {
			resp.Diagnostics.AddError(
				"Listing AnswerX service IDs failed", "Failed to list AnswerX service IDs: invalid API response",
			)
			return
		}

		allServiceIDs = append(allServiceIDs, response.AnswerXServiceIDs...)

		if response.Metadata.Page >= response.Metadata.LastPage {
			break
		}
		page++
	}

	serviceIDs, diags := parseAnswerXServiceIDModels(ctx, allServiceIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ServiceIDs = serviceIDs

	if resp.Diagnostics.Append(resp.State.Set(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}
}

func parseAnswerXServiceIDModels(ctx context.Context, serviceIDs []datastream.AnswerXServiceDetail) (types.Set, diag.Diagnostics) {
	items := make([]answerXServiceIDItemModel, 0, len(serviceIDs))
	for _, s := range serviceIDs {
		items = append(items, answerXServiceIDItemModel{
			ID:      types.Int64Value(s.SSID),
			Name:    types.StringValue(s.Name),
			Product: types.StringValue(s.Product),
		})
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: answerXServiceIDItemAttrTypes}, items)
}
