package appsec

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &scriptBehaviorRuleResource{}
	_ resource.ResourceWithConfigure   = &scriptBehaviorRuleResource{}
	_ resource.ResourceWithImportState = &scriptBehaviorRuleResource{}
)

type scriptBehaviorRuleResource struct {
	meta.Resource
}

type scriptBehaviorRuleModel struct {
	ConfigID                types.Int64  `tfsdk:"config_id"`
	SecurityPolicyID        types.String `tfsdk:"security_policy_id"`
	RuleID                  types.Int64  `tfsdk:"rule_id"`
	RuleName                types.String `tfsdk:"rule_name"`
	RuleDescription         types.String `tfsdk:"rule_description"`
	ActivityType            types.String `tfsdk:"activity_type"`
	Action                  types.String `tfsdk:"action"`
	PageCondition           types.String `tfsdk:"page_condition"`
	PageValues              types.List   `tfsdk:"page_values"`
	ScriptType              types.String `tfsdk:"script_type"`
	ScriptCondition         types.String `tfsdk:"script_condition"`
	ScriptValues            types.List   `tfsdk:"script_values"`
	ActivityValuesCondition types.String `tfsdk:"activity_values_condition"`
	ActivityValues          types.List   `tfsdk:"activity_values"`
}

// NewScriptBehaviorRuleResource returns a new script behavior rule resource.
func NewScriptBehaviorRuleResource() resource.Resource { return &scriptBehaviorRuleResource{} }

func (r *scriptBehaviorRuleResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_script_behavior_rule"
}

func (r *scriptBehaviorRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a script behavior rule for a security policy.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier of the script behavior rule, assigned by the API upon creation.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"rule_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the script behavior rule.",
			},
			"rule_description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the script behavior rule.",
			},
			"activity_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of activity the rule controls.",
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "Action taken when the rule conditions are matched.",
			},
			"page_condition": schema.StringAttribute{
				Required:    true,
				Description: "Condition type for matching page URLs.",
			},
			"page_values": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of page URL patterns to match.",
			},
			"script_type": schema.StringAttribute{
				Optional:    true,
				Description: "Type of script to match.",
			},
			"script_condition": schema.StringAttribute{
				Required:    true,
				Description: "Condition type for matching scripts.",
			},
			"script_values": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of script patterns to match.",
			},
			"activity_values_condition": schema.StringAttribute{
				Required:    true,
				Description: "Condition type for matching activity type values.",
			},
			"activity_values": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of activity type values to match.",
			},
		},
	}
}

func (r *scriptBehaviorRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "creating script behavior rule")

	var plan scriptBehaviorRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client.GetAPPSEC()
	configID := plan.ConfigID.ValueInt64()

	version, err := getModifiableConfigVersion(ctx, configID, "scriptBehaviorRule", client)
	if err != nil {
		resp.Diagnostics.AddError("getting modifiable config version", err.Error())
		return
	}

	body, diags := scriptBehaviorRuleBodyFromModel(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := client.CreateClientSideProtectionRule(ctx, appsec.CreateClientSideProtectionRuleRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: plan.SecurityPolicyID.ValueString(),
		Body:     body,
	})
	if err != nil {
		resp.Diagnostics.AddError("creating script behavior rule", err.Error())
		return
	}
	if created.RuleID == nil {
		resp.Diagnostics.AddError("creating script behavior rule", "API did not return a rule ID")
		return
	}
	plan.RuleID = types.Int64Value(*created.RuleID)

	found, diags := r.read(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("reading created rule", fmt.Sprintf("rule %d not found after creation", plan.RuleID.ValueInt64()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scriptBehaviorRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "reading script behavior rule")

	var state scriptBehaviorRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, diags := r.read(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		tflog.Warn(ctx, "script behavior rule not found, removing from state", map[string]any{
			"config_id":          state.ConfigID.ValueInt64(),
			"security_policy_id": state.SecurityPolicyID.ValueString(),
			"rule_id":            state.RuleID.ValueInt64(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scriptBehaviorRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "updating script behavior rule")

	var plan scriptBehaviorRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state scriptBehaviorRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.RuleID = state.RuleID

	client := r.Client.GetAPPSEC()
	configID := plan.ConfigID.ValueInt64()
	ruleID := plan.RuleID.ValueInt64()

	version, err := getModifiableConfigVersion(ctx, configID, "scriptBehaviorRule", client)
	if err != nil {
		resp.Diagnostics.AddError("getting modifiable config version", err.Error())
		return
	}

	body, diags := scriptBehaviorRuleBodyFromModel(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	body.RuleID = &ruleID

	_, err = client.UpdateClientSideProtectionRule(ctx, appsec.UpdateClientSideProtectionRuleRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: plan.SecurityPolicyID.ValueString(),
		RuleID:   ruleID,
		Body:     body,
	})
	if err != nil {
		resp.Diagnostics.AddError("updating script behavior rule", err.Error())
		return
	}

	found, diags := r.read(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("reading updated rule", fmt.Sprintf("rule %d not found after update", ruleID))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scriptBehaviorRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "deleting script behavior rule")

	var state scriptBehaviorRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client.GetAPPSEC()
	configID := state.ConfigID.ValueInt64()

	version, err := getModifiableConfigVersion(ctx, configID, "scriptBehaviorRule", client)
	if err != nil {
		resp.Diagnostics.AddError("getting modifiable config version", err.Error())
		return
	}

	err = client.DeleteClientSideProtectionRule(ctx, appsec.DeleteClientSideProtectionRuleRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: state.SecurityPolicyID.ValueString(),
		RuleID:   state.RuleID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("deleting script behavior rule", err.Error())
	}
}

func (r *scriptBehaviorRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 3 {
		resp.Diagnostics.AddError("invalid import ID", fmt.Sprintf("invalid import ID %q: expected format is 'CONFIG_ID:POLICY_ID:RULE_ID'", req.ID))
		return
	}

	var configID, ruleID int64
	if _, err := fmt.Sscan(parts[0], &configID); err != nil {
		resp.Diagnostics.AddError("invalid import ID", fmt.Sprintf("config_id must be an integer, got %q", parts[0]))
		return
	}
	if _, err := fmt.Sscan(parts[2], &ruleID); err != nil {
		resp.Diagnostics.AddError("invalid import ID", fmt.Sprintf("rule_id must be an integer, got %q", parts[2]))
		return
	}

	state := scriptBehaviorRuleModel{
		ConfigID:         types.Int64Value(configID),
		SecurityPolicyID: types.StringValue(parts[1]),
		RuleID:           types.Int64Value(ruleID),
	}

	found, diags := r.read(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.Diagnostics.AddError("importing script behavior rule", fmt.Sprintf("rule %d not found in policy %q", ruleID, parts[1]))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scriptBehaviorRuleResource) read(ctx context.Context, data *scriptBehaviorRuleModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	client := r.Client.GetAPPSEC()
	configID := data.ConfigID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, configID, client)
	if err != nil {
		diags.AddError("getting latest config version", err.Error())
		return false, diags
	}

	rules, err := client.ListClientSideProtectionRules(ctx, appsec.ListClientSideProtectionRulesRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: data.SecurityPolicyID.ValueString(),
	})
	if err != nil {
		diags.AddError("listing script behavior rules", err.Error())
		return false, diags
	}

	ruleID := data.RuleID.ValueInt64()
	var found *appsec.ClientSideProtectionRuleBody
	for i := range *rules {
		if (*rules)[i].RuleID != nil && *(*rules)[i].RuleID == ruleID {
			found = &(*rules)[i]
			break
		}
	}
	if found == nil {
		return false, diags
	}

	data.RuleName = types.StringValue(found.RuleName)
	data.RuleDescription = stringOrNull(found.Description)
	data.ActivityType = types.StringValue(found.ActivityType)
	data.Action = types.StringValue(found.Action)
	data.PageCondition = types.StringValue(found.Conditions.Page.Condition)
	data.ScriptType = stringOrNull(found.Conditions.Script.ScriptType)
	data.ScriptCondition = types.StringValue(found.Conditions.Script.Condition)
	data.ActivityValuesCondition = types.StringValue(found.Conditions.ActivityTypeValues.Condition)

	pageValues, d := types.ListValueFrom(ctx, types.StringType, found.Conditions.Page.Values)
	diags.Append(d...)
	data.PageValues = pageValues

	scriptValues, d := types.ListValueFrom(ctx, types.StringType, found.Conditions.Script.Values)
	diags.Append(d...)
	data.ScriptValues = scriptValues

	activityValues, d := types.ListValueFrom(ctx, types.StringType, found.Conditions.ActivityTypeValues.Values)
	diags.Append(d...)
	data.ActivityValues = activityValues

	return !diags.HasError(), diags
}

func scriptBehaviorRuleBodyFromModel(ctx context.Context, data scriptBehaviorRuleModel) (appsec.ClientSideProtectionRuleBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	var pageValues []string
	diags.Append(data.PageValues.ElementsAs(ctx, &pageValues, false)...)

	var scriptValues []string
	diags.Append(data.ScriptValues.ElementsAs(ctx, &scriptValues, false)...)

	var activityValues []string
	diags.Append(data.ActivityValues.ElementsAs(ctx, &activityValues, false)...)

	if diags.HasError() {
		return appsec.ClientSideProtectionRuleBody{}, diags
	}

	return appsec.ClientSideProtectionRuleBody{
		RuleName:     data.RuleName.ValueString(),
		Description:  data.RuleDescription.ValueString(),
		ActivityType: data.ActivityType.ValueString(),
		Action:       data.Action.ValueString(),
		Conditions: appsec.ClientSideProtectionConditions{
			Page: appsec.ClientSideProtectionConditionCriteria{
				Condition: data.PageCondition.ValueString(),
				Values:    pageValues,
			},
			Script: appsec.ClientSideProtectionScriptConditionCriteria{
				ScriptType: data.ScriptType.ValueString(),
				Condition:  data.ScriptCondition.ValueString(),
				Values:     scriptValues,
			},
			ActivityTypeValues: appsec.ClientSideProtectionConditionCriteria{
				Condition: data.ActivityValuesCondition.ValueString(),
				Values:    activityValues,
			},
		},
	}, diags
}
