package reportinggroups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/reportinggroups"
	"github.com/akamai/terraform-provider-akamai/v11/internal/customtypes"
	"github.com/akamai/terraform-provider-akamai/v11/internal/text"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &reportingGroupsResource{}
	_ resource.ResourceWithImportState = &reportingGroupsResource{}
	_ resource.ResourceWithConfigure   = &reportingGroupsResource{}
	_ resource.ResourceWithModifyPlan  = &reportingGroupsResource{}
)

type reportingGroupsResourceModel struct {
	ReportingGroupID   types.Int64  `tfsdk:"reporting_group_id"`
	ReportingGroupName types.String `tfsdk:"reporting_group_name"`
	AccessGroup        types.Object `tfsdk:"access_group"`
	Contract           types.Object `tfsdk:"contract"`
}

// accessGroup extracts the AccessGroup object into a typed model.
func (m *reportingGroupsResourceModel) accessGroup(ctx context.Context) (*accessGroupResourceModel, diag.Diagnostics) {
	var ag accessGroupResourceModel
	return &ag, m.AccessGroup.As(ctx, &ag, basetypes.ObjectAsOptions{})
}

// contract extracts the Contract object into a typed model.
func (m *reportingGroupsResourceModel) contract(ctx context.Context) (*contractResourceModel, diag.Diagnostics) {
	var c contractResourceModel
	return &c, m.Contract.As(ctx, &c, basetypes.ObjectAsOptions{})
}

// setAccessGroup stores ag into the AccessGroup object field.
func (m *reportingGroupsResourceModel) setAccessGroup(ctx context.Context, ag *accessGroupResourceModel) diag.Diagnostics {
	var dd diag.Diagnostics
	m.AccessGroup, dd = types.ObjectValueFrom(ctx, accessGroupType(), ag)
	return dd
}

// setContract stores c into the Contract object field.
func (m *reportingGroupsResourceModel) setContract(ctx context.Context, c *contractResourceModel) diag.Diagnostics {
	var dd diag.Diagnostics
	m.Contract, dd = types.ObjectValueFrom(ctx, contractType(), c)
	return dd
}

type accessGroupResourceModel struct {
	ContractID customtypes.IgnorePrefixValue `tfsdk:"contract_id"`
	GroupID    types.String                  `tfsdk:"group_id"`
}

type contractResourceModel struct {
	ContractID customtypes.IgnorePrefixValue `tfsdk:"contract_id"`
	CPCodes    types.Set                     `tfsdk:"cp_codes"`
}

// parsedCPCodeIDs returns the CP code IDs as int64, stripping the "cpc_" prefix.
func (m *contractResourceModel) parsedCPCodeIDs(ctx context.Context) ([]int64, diag.Diagnostics) {
	var diags diag.Diagnostics
	var cps []cpCodeResourceModel
	if diags = m.CPCodes.ElementsAs(ctx, &cps, false); diags.HasError() {
		return nil, diags
	}
	ids := make([]int64, 0, len(cps))
	for _, cp := range cps {
		id, err := str.GetInt64ID(cp.CPCodeID.ValueString(), "cpc_")
		if err != nil {
			diags.AddError(
				"Invalid CP Code ID",
				fmt.Sprintf(
					"Error with conversion CP Code ID string to int64. Got %q: %s",
					cp.CPCodeID.ValueString(), err,
				),
			)
			return nil, diags
		}
		ids = append(ids, id)
	}
	return ids, diags
}

// setCPCodes stores cps into the CPCodes set field.
func (m *contractResourceModel) setCPCodes(ctx context.Context, cps []cpCodeResourceModel) diag.Diagnostics {
	var dd diag.Diagnostics
	m.CPCodes, dd = types.SetValueFrom(ctx, cpcodeType(), cps)
	return dd
}

type cpCodeResourceModel struct {
	CPCodeID   customtypes.IgnorePrefixValue `tfsdk:"cp_code_id"`
	CPCodeName types.String                  `tfsdk:"cp_code_name"`
}

// NewReportingGroupsResource returns a new Reporting Groups resource.
func NewReportingGroupsResource() resource.Resource {
	return &reportingGroupsResource{}
}

type reportingGroupsResource struct {
	meta.Resource
}

// ModifyPlan enforces that access_group.group_id is immutable after creation.
func (r *reportingGroupsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !modifiers.IsUpdate(req) {
		return
	}

	var state, plan reportingGroupsResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	stateAG, dd := state.accessGroup(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	planAG, dd := plan.accessGroup(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}

	groupIDPath := path.Root("access_group").AtName("group_id")

	if stateAG.GroupID.IsNull() && !planAG.GroupID.IsNull() {
		resp.Diagnostics.AddAttributeError(
			groupIDPath,
			"Resource was imported without a group_id",
			"The resource was imported without a group_id, but the configuration now specifies one. "+
				"To fix this, remove the resource from state with: terraform state rm <address>, "+
				"then re-import it with the group_id included in the import ID: "+
				"terraform import <address> <reportingGroupID>,<groupID>",
		)
		return
	}

	if !stateAG.GroupID.Equal(planAG.GroupID) {
		resp.Diagnostics.AddAttributeError(
			groupIDPath,
			"Cannot update access_group.group_id after creation",
			"The group_id field cannot be changed once set. Remove the resource and recreate it instead.",
		)
	}
}

func (r *reportingGroupsResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_reportinggroups_group"
}

func (r *reportingGroupsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Reporting Group, which organizes CP codes under a named group for reporting purposes.",
		Attributes: map[string]schema.Attribute{
			"reporting_group_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier of the reporting group. Populated after creation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"reporting_group_name": schema.StringAttribute{
				Required:    true,
				Description: "The descriptive label for the reporting group.",
				Validators:  []validator.String{validators.NotEmptyString()},
			},
			"access_group": accessGroupSchema(),
			"contract":     contractSchema(),
		},
	}
}

func accessGroupSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required:    true,
		Description: "The access control group that controls access to specific CP codes.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifies the contract assigned to the access control group.",
				CustomType:  customtypes.IgnorePrefixType{Prefix: "ctr_"},
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{validators.NotEmptyString()},
			},
			"group_id": schema.StringAttribute{
				Optional: true,
				Description: "Identifies the access control group. It is required for " +
					"reporting group creation and cannot be updated.",
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("grp_")),
					modifiers.StringRequiredForCreate(),
				},
			},
		},
	}
}

func contractSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required: true,
		Description: "A collection of contracts and CP codes assigned to the reporting group. " +
			"Exactly one contract is allowed.",
		Attributes: map[string]schema.Attribute{
			"contract_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifies the contract assigned to the reporting group.",
				CustomType:  customtypes.IgnorePrefixType{Prefix: "ctr_"},
				PlanModifiers: []planmodifier.String{
					modifiers.StringUseStateIf(modifiers.EqualUpToPrefixFunc("ctr_")),
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{validators.NotEmptyString()},
			},
			"cp_codes": cpCodeSchema(),
		},
	}

}

func cpCodeSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Required:    true,
		Description: "A collection of CP codes assigned to the reporting group.",
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"cp_code_id": schema.StringAttribute{
					Required:    true,
					Description: "Identifies a CP code.",
					CustomType:  customtypes.IgnorePrefixType{Prefix: "cpc_"},
					Validators:  []validator.String{validators.NotEmptyString()},
				},
				"cp_code_name": schema.StringAttribute{
					Computed:    true,
					Description: "The descriptive label for the CP code.",
				},
			},
		},
	}
}

func (r *reportingGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Reporting Group")

	var plan reportingGroupsResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	contract, dd := plan.contract(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}

	cpCodeIDs, dd := contract.parsedCPCodeIDs(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}

	cpCodesCreateReq := make([]reportinggroups.CPCodeCreate, 0, len(cpCodeIDs))
	for _, id := range cpCodeIDs {
		cpCodesCreateReq = append(cpCodesCreateReq, reportinggroups.CPCodeCreate{CPCodeID: id})
	}

	accessGroup, dd := plan.accessGroup(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	groupID, err := str.GetInt64ID(accessGroup.GroupID.ValueString(), "grp_")
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Access Group ID",
			fmt.Sprintf("Error with conversion GroupID string to int64. Got %q: %s", accessGroup.GroupID.ValueString(), err),
		)
		return
	}

	createReq := reportinggroups.CreateReportingGroupRequest{
		ReportingGroupName: plan.ReportingGroupName.ValueString(),
		AccessGroup: reportinggroups.AccessGroup{
			ContractID: strings.TrimPrefix(accessGroup.ContractID.ValueString(), "ctr_"),
			GroupID:    ptr.To(groupID),
		},
		Contracts: []reportinggroups.ContractCreate{{
			ContractID: strings.TrimPrefix(contract.ContractID.ValueString(), "ctr_"),
			CPCodes:    cpCodesCreateReq,
		}},
	}
	result, err := r.Client.GetReportingGroups().CreateReportingGroup(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create Reporting Group", err.Error())
		return
	}
	tflog.Debug(ctx, "Reporting Group created", map[string]any{"reporting_group_id": result.ReportingGroupID})

	plan.ReportingGroupID = types.Int64Value(result.ReportingGroupID)
	// reporting_group_name and access_group are not repopulated from the Create response:
	// the API echoes them back unchanged, so the plan values are authoritative (drifts can be detected in Read).
	// access_group.group_id is skipped for a different reason: the API never returns it.
	if resp.Diagnostics.Append(plan.populateContract(ctx, result.Contracts)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *reportingGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Reporting Group")

	var state reportingGroupsResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}
	reportingGroupID := state.ReportingGroupID
	ctx = tflog.SetField(ctx, "reporting_group_id", reportingGroupID.ValueInt64())

	result, err := r.Client.GetReportingGroups().GetReportingGroup(ctx, reportinggroups.GetReportingGroupsRequest{
		ReportingGroupID: reportingGroupID.ValueInt64(),
	})
	if err != nil {
		if errors.Is(err, &reportinggroups.Error{HTTPStatus: http.StatusNotFound}) {
			resp.Diagnostics.AddWarning(
				"Reporting Group not found",
				fmt.Sprintf(
					"Error 404 returned while reading reporting group %d. "+
						"The resource may have been deleted outside Terraform "+
						"or the ID may be invalid. Terraform will remove it from state.",
					reportingGroupID.ValueInt64(),
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Reporting Group", err.Error())
		return
	}

	state.ReportingGroupName = types.StringValue(result.ReportingGroupName)

	if resp.Diagnostics.Append(state.populateAccessGroup(ctx, result.AccessGroup)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(state.populateContract(ctx, result.Contracts)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *reportingGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating Reporting Group")

	var plan reportingGroupsResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state reportingGroupsResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "reporting_group_id", state.ReportingGroupID.ValueInt64())

	contract, dd := plan.contract(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	cpCodeIDs, dd := contract.parsedCPCodeIDs(ctx)
	if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
		return
	}
	cpCodesReq := make([]reportinggroups.CPCode, 0, len(cpCodeIDs))
	for _, id := range cpCodeIDs {
		cpCodesReq = append(cpCodesReq, reportinggroups.CPCode{CPCodeID: id})
	}

	updateReq := reportinggroups.UpdateReportingGroupRequest{
		ReportingGroupID:   state.ReportingGroupID.ValueInt64(),
		ReportingGroupName: plan.ReportingGroupName.ValueString(),
		Contracts: []reportinggroups.Contract{{
			ContractID: strings.TrimPrefix(contract.ContractID.ValueString(), "ctr_"),
			CPCodes:    cpCodesReq,
		}},
	}
	result, err := r.Client.GetReportingGroups().UpdateReportingGroup(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update Reporting Group", err.Error())
		return
	}
	tflog.Debug(ctx, "Reporting Group updated")

	// access_group is not repopulated: contract_id is immutable so the plan value is already correct,
	// and group_id is never returned by the API so the plan value must be preserved as-is.
	if resp.Diagnostics.Append(plan.populateContract(ctx, result.Contracts)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *reportingGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Reporting Group")

	var state reportingGroupsResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "reporting_group_id", state.ReportingGroupID.ValueInt64())

	if err := r.Client.GetReportingGroups().DeleteReportingGroup(ctx, reportinggroups.DeleteReportingGroupRequest{
		ReportingGroupID: state.ReportingGroupID.ValueInt64(),
	}); err != nil {
		resp.Diagnostics.AddError("Unable to delete Reporting Group", err.Error())
		return
	}

	tflog.Debug(ctx, "Reporting Group deleted")
}

// ImportState implements resource's ImportState method. The import ID is the reporting_group_id.
func (r *reportingGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Reporting Group")

	parts, err := text.ImportIDSplitter("reportingGroupID[,groupID]").
		AcceptLen(1).
		AcceptLen(2).
		Split(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Incorrect import ID", err.Error())
		return
	}

	reportingGroupID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Reporting Group ID in import ID",
			fmt.Sprintf("Error parsing reporting group ID from import ID: %s. Expected format: reportingGroupID[,groupID]", err),
		)
		return
	}

	// Parse group_id early (before the API call) to fail fast on malformed input.
	var groupID int64
	if len(parts) == 2 {
		groupID, err = str.GetInt64ID(parts[1], "grp_")
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Group ID in import ID",
				fmt.Sprintf("Error parsing group ID from import ID: %s. Expected format: reportingGroupID[,groupID]", err),
			)
			return
		}
	}

	ctx = tflog.SetField(ctx, "reporting_group_id", reportingGroupID)

	result, err := r.Client.GetReportingGroups().GetReportingGroup(ctx, reportinggroups.GetReportingGroupsRequest{
		ReportingGroupID: reportingGroupID,
	})
	if err != nil {
		if errors.Is(err, &reportinggroups.Error{HTTPStatus: http.StatusNotFound}) {
			resp.Diagnostics.AddError(
				"Cannot import non-existent remote object",
				fmt.Sprintf("Reporting Group with ID %d was not found on the server. "+
					"Please verify the ID is correct.", reportingGroupID),
			)
			return
		}
		resp.Diagnostics.AddError("Unable to import Reporting Group", err.Error())
		return
	}

	tflog.Debug(ctx, "Reporting Group found, importing")

	var state reportingGroupsResourceModel
	state.ReportingGroupID = types.Int64Value(result.ReportingGroupID)
	state.ReportingGroupName = types.StringValue(result.ReportingGroupName)
	if resp.Diagnostics.Append(state.populateAccessGroup(ctx, result.AccessGroup)...); resp.Diagnostics.HasError() {
		return
	}
	if resp.Diagnostics.Append(state.populateContract(ctx, result.Contracts)...); resp.Diagnostics.HasError() {
		return
	}

	// The API never returns group_id; use the value parsed from the import ID, if provided.
	if len(parts) == 2 {
		ag, dd := state.accessGroup(ctx)
		if resp.Diagnostics.Append(dd...); resp.Diagnostics.HasError() {
			return
		}
		ag.GroupID = types.StringValue(strconv.FormatInt(groupID, 10))
		if resp.Diagnostics.Append(state.setAccessGroup(ctx, ag)...); resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (m *reportingGroupsResourceModel) populateAccessGroup(
	ctx context.Context, group reportinggroups.AccessGroup,
) diag.Diagnostics {
	var diags diag.Diagnostics
	currentAG := &accessGroupResourceModel{}

	// AccessGroup will be Null during import

	// GroupID is write-only from the API's perspective — it is never returned in
	// any API response, so always preserve the existing state value (or null on import).
	if tf.IsKnown(m.AccessGroup) {
		ag, dd := m.accessGroup(ctx)
		if diags.Append(dd...); diags.HasError() {
			return diags
		}
		currentAG = ag
	}

	currentAG.ContractID = customtypes.NewIgnorePrefixValue("ctr_", group.ContractID)
	diags.Append(m.setAccessGroup(ctx, currentAG)...)
	return diags
}

func (m *reportingGroupsResourceModel) populateContract(
	ctx context.Context, contracts []reportinggroups.Contract,
) diag.Diagnostics {
	var diags diag.Diagnostics
	if len(contracts) == 0 {
		diags.AddError(
			"No Contract Data",
			"Expected one contract in API response, but got none.",
		)
		return diags
	}
	c := contracts[0]

	contract := contractResourceModel{
		ContractID: customtypes.NewIgnorePrefixValue("ctr_", c.ContractID),
	}

	cpCodes := make([]cpCodeResourceModel, 0, len(c.CPCodes))
	for _, cp := range c.CPCodes {
		cpCodes = append(cpCodes, cpCodeResourceModel{
			CPCodeID:   customtypes.NewIgnorePrefixValue("cpc_", strconv.FormatInt(cp.CPCodeID, 10)),
			CPCodeName: types.StringValue(cp.CPCodeName),
		})
	}

	if diags.Append(contract.setCPCodes(ctx, cpCodes)...); diags.HasError() {
		return diags
	}
	diags.Append(m.setContract(ctx, &contract)...)
	return diags
}

func accessGroupType() map[string]attr.Type {
	return accessGroupSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func contractType() map[string]attr.Type {
	return contractSchema().GetType().(attr.TypeWithAttributeTypes).AttributeTypes()
}

func cpcodeType() types.ObjectType {
	return cpCodeSchema().NestedObject.Type().(types.ObjectType)
}
