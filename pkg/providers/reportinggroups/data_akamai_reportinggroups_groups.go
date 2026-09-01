package reportinggroups

import (
	"context"
	"sort"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &reportingGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &reportingGroupsDataSource{}
)

type (
	reportingGroupsDataSource struct {
		meta.DataSource
	}

	reportingGroupsDataSourceModel struct {
		ContractID         types.String          `tfsdk:"contract_id"`
		GroupID            types.String          `tfsdk:"group_id"`
		ReportingGroupName types.String          `tfsdk:"reporting_group_name"`
		CPCodeID           types.String          `tfsdk:"cp_code_id"`
		Groups             []reportingGroupModel `tfsdk:"groups"`
	}

	reportingGroupModel struct {
		ReportingGroupID   types.Int64                  `tfsdk:"reporting_group_id"`
		ReportingGroupName types.String                 `tfsdk:"reporting_group_name"`
		AccessGroup        *accessGroupModel            `tfsdk:"access_group"`
		Contract           *reportingGroupContractModel `tfsdk:"contract"`
	}
)

// NewReportingGroupsDataSource returns a new reportingGroupsDataSource.
func NewReportingGroupsDataSource() datasource.DataSource {
	return &reportingGroupsDataSource{}
}

// Metadata configures data source's meta information.
func (d *reportingGroupsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_reportinggroups_groups"
}

// Schema is used to define data source's terraform schema.
func (d *reportingGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List detailed information about reporting groups available for your account and contract.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies the contract to filter reporting groups by.",
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies the access group to filter reporting groups by. Accepts numeric ID or grp_ prefixed value.",
			},
			"reporting_group_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the reporting group to filter by.",
			},
			"cp_code_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifies the CP code to filter reporting groups by.",
			},
			"groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of reporting groups available for your account and contract.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"reporting_group_id": schema.Int64Attribute{
							Computed:    true,
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
				},
			},
		},
	}
}

// Read is called when the provider must read data source values in order to update state.
func (d *reportingGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reporting Groups Groups DataSource Read")

	var data reportingGroupsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	listResp, err := d.Client.GetReportingGroups().ListReportingGroups(ctx, reportinggroups.ListReportingGroupsRequest{
		ContractID:         data.ContractID.ValueString(),
		GroupID:            data.GroupID.ValueString(),
		ReportingGroupName: data.ReportingGroupName.ValueString(),
		CPCodeID:           data.CPCodeID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("List Reporting Groups failed", err.Error())
		return
	}

	data.convertReportingGroupsToModel(listResp.Groups)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (m *reportingGroupsDataSourceModel) convertReportingGroupsToModel(groups []reportinggroups.ReportingGroup) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ReportingGroupID < groups[j].ReportingGroupID
	})

	m.Groups = make([]reportingGroupModel, 0, len(groups))
	for _, group := range groups {
		var contractModel *reportingGroupContractModel
		if len(group.Contracts) == 1 {
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
			contractModel = &reportingGroupContractModel{
				ContractID: types.StringValue(contract.ContractID),
				CPCodes:    cpCodes,
			}
		}

		var groupIDVal types.String
		if group.AccessGroup.GroupID != nil {
			groupIDVal = types.StringValue(strconv.FormatInt(*group.AccessGroup.GroupID, 10))
		} else {
			groupIDVal = types.StringNull()
		}

		m.Groups = append(m.Groups, reportingGroupModel{
			ReportingGroupID:   types.Int64Value(group.ReportingGroupID),
			ReportingGroupName: types.StringValue(group.ReportingGroupName),
			AccessGroup: &accessGroupModel{
				ContractID: types.StringValue(group.AccessGroup.ContractID),
				GroupID:    groupIDVal,
			},
			Contract: contractModel,
		})
	}
}
