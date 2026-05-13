package reportinggroups

import (
	"context"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &reportingGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &reportingGroupDataSource{}
)

type (
	reportingGroupDataSource struct {
		meta.DataSource
	}

	reportingGroupDataSourceModel struct {
		ReportingGroupID   types.Int64                  `tfsdk:"reporting_group_id"`
		ReportingGroupName types.String                 `tfsdk:"reporting_group_name"`
		AccessGroup        *accessGroupModel            `tfsdk:"access_group"`
		Contract           *reportingGroupContractModel `tfsdk:"contract"`
	}

	accessGroupModel struct {
		ContractID types.String `tfsdk:"contract_id"`
		GroupID    types.String `tfsdk:"group_id"`
	}

	reportingGroupContractModel struct {
		ContractID types.String                `tfsdk:"contract_id"`
		CPCodes    []reportingGroupCPCodeModel `tfsdk:"cp_codes"`
	}

	reportingGroupCPCodeModel struct {
		CPCodeID   types.String `tfsdk:"cp_code_id"`
		CPCodeName types.String `tfsdk:"cp_code_name"`
	}
)

// NewReportingGroupDataSource returns a new reportingGroupDataSource.
func NewReportingGroupDataSource() datasource.DataSource {
	return &reportingGroupDataSource{}
}

// Metadata configures data source's meta information.
func (d *reportingGroupDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_reportinggroups_group"
}

// Schema is used to define data source's terraform schema.
func (d *reportingGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve details of a Reporting Group.",
		Attributes: map[string]schema.Attribute{
			"reporting_group_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier for the reporting group.",
			},
			"reporting_group_name": schema.StringAttribute{
				Computed:    true,
				Description: "The descriptive label for the reporting group.",
			},
			"access_group": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The access control group that controls access to specific CP codes.",
				Attributes: map[string]schema.Attribute{
					"contract_id": schema.StringAttribute{
						Computed:    true,
						Description: "Identifies the contract assigned to the access control group.",
					},
					"group_id": schema.StringAttribute{
						Computed:    true,
						Description: "Identifies the access control group. May be null if the reporting group belongs to many groups.",
					},
				},
			},
			"contract": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The contract and CP codes assigned to the reporting group.",
				Attributes: map[string]schema.Attribute{
					"contract_id": schema.StringAttribute{
						Computed:    true,
						Description: "Identifies the contract assigned to the reporting group.",
					},
					"cp_codes": schema.ListNestedAttribute{
						Computed:    true,
						Description: "A collection of CP codes assigned to the reporting group.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"cp_code_id": schema.StringAttribute{
									Computed:    true,
									Description: "Identifies a CP code.",
								},
								"cp_code_name": schema.StringAttribute{
									Computed:    true,
									Description: "The descriptive label for the CP code.",
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *reportingGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reporting Groups Group DataSource Read")

	var data reportingGroupDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	group, err := d.Client.GetReportingGroups().GetReportingGroup(ctx, reportinggroups.GetReportingGroupsRequest{
		ReportingGroupID: data.ReportingGroupID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Read Reporting Group failed", err.Error())
		return
	}

	data.convertReportingGroupToModel(*group)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *reportingGroupDataSourceModel) convertReportingGroupToModel(group reportinggroups.GetReportingGroupResponse) {
	m.ReportingGroupID = types.Int64Value(group.ReportingGroupID)
	m.ReportingGroupName = types.StringValue(group.ReportingGroupName)

	var groupIDVal types.String
	if group.AccessGroup.GroupID != nil {
		groupIDVal = types.StringValue(strconv.FormatInt(*group.AccessGroup.GroupID, 10))
	} else {
		groupIDVal = types.StringNull()
	}
	m.AccessGroup = &accessGroupModel{
		ContractID: types.StringValue(group.AccessGroup.ContractID),
		GroupID:    groupIDVal,
	}

	m.Contract = nil
	if len(group.Contracts) > 0 {
		contract := group.Contracts[0]
		cpCodes := make([]reportingGroupCPCodeModel, 0, len(contract.CPCodes))
		for _, cp := range contract.CPCodes {
			cpCodes = append(cpCodes, reportingGroupCPCodeModel{
				CPCodeID:   types.StringValue(strconv.FormatInt(cp.CPCodeID, 10)),
				CPCodeName: types.StringValue(cp.CPCodeName),
			})
		}
		sort.Slice(cpCodes, func(i, j int) bool {
			idI, _ := strconv.ParseInt(cpCodes[i].CPCodeID.ValueString(), 10, 64)
			idJ, _ := strconv.ParseInt(cpCodes[j].CPCodeID.ValueString(), 10, 64)
			return idI < idJ
		})
		m.Contract = &reportingGroupContractModel{
			ContractID: types.StringValue(contract.ContractID),
			CPCodes:    cpCodes,
		}
	}
}
