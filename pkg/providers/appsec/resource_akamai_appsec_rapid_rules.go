package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &rapidRulesResource{}
	_ resource.ResourceWithConfigure      = &rapidRulesResource{}
	_ resource.ResourceWithValidateConfig = &rapidRulesResource{}
	_ resource.ResourceWithImportState    = &rapidRulesResource{}
)

// rapidRulesResource represents akamai_appsec_rapid_rule resource
type rapidRulesResource struct {
	meta meta.Meta
}

// rapidRulesResourceModel is a model for akamai_appsec_rapid_rule resource
type rapidRulesResourceModel struct {
	ID              types.String                   `tfsdk:"id"`
	Enabled         types.Bool                     `tfsdk:"enabled"`
	ConfigID        types.Int64                    `tfsdk:"config_id"`
	PolicyID        types.String                   `tfsdk:"security_policy_id"`
	DefaultAction   types.String                   `tfsdk:"default_action"`
	RuleDefinitions rapidRuleDefinitionsStateValue `tfsdk:"rule_definitions"`
}

const (
	ruleDefinitionsValidationError      = "Rule definitions validation error: %s"
	invalidConfigurationAttribute       = "Invalid configuration attribute"
	readRapidRulesError                 = "Unable to read rapid rules"
	updateRapidRulesStatusError         = "Unable to update rapid rules status"
	updateRapidRulesDefaultActionFailed = "Unable to update rapid rules default action"
	resourceName                        = "rapidRules"
)

// NewRapidRulesResource returns new appsec rapid rules resource
func NewRapidRulesResource() resource.Resource {
	return &rapidRulesResource{}
}

// Metadata implements resource.Resource.
func (r *rapidRulesResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_rapid_rules"
}

// Schema implements resource's Schema
func (r *rapidRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rapid rule resource.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy",
				Validators:  []validator.String{validators.NotEmptyString()},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"default_action": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Default action that applies to violations of all rapid rules",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_definitions": schema.StringAttribute{
				CustomType:  rapidRuleDefinitionsStateType{},
				Optional:    true,
				Computed:    true,
				Description: "JSON-formatted list of rule definition (ID, action, action lock and exception)",
				Validators:  []validator.String{validators.NotEmptyString()},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString("null"),
			},
			"id": schema.StringAttribute{
				Description: "Identifier of the resource",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Hidden attribute containing information about rapid rules status enabled/disabled",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *rapidRulesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// ProviderData is nil when Configure is run first time as part of ValidateDataSourceConfig in framework provider
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	r.meta = meta.Must(req.ProviderData)
}

// ValidateConfig implements resource.ResourceWithValidateConfig.
func (r *rapidRulesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data rapidRulesResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultAction := data.DefaultAction.ValueString()
	if !data.DefaultAction.IsNull() && !data.DefaultAction.IsUnknown() {
		diags := validateDefaultAction(&defaultAction)
		if diags.HasError() {
			resp.Diagnostics.AddAttributeError(path.Root("default_action"), invalidConfigurationAttribute, extractErrors(diags.Errors()))
		}
	}

	ruleDefinitions, diags := getRuleDefinitionsFromModel(&data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if ruleDefinitions != nil && len(ruleDefinitions) == 0 {
		resp.Diagnostics.AddAttributeError(path.Root("rule_definitions"), "JSON cannot be empty", "Invalid rule definition file")
		return
	}

	if len(ruleDefinitions) > 0 {
		diags = validateRuleDefinitions(ruleDefinitions)
		if diags.HasError() {
			resp.Diagnostics.AddAttributeError(path.Root("rule_definitions"), invalidConfigurationAttribute, extractErrors(diags.Errors()))
		}
	}
}

// Create implements resource's Create method
func (r *rapidRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating Rapid Rules Resource")

	var data *rapidRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	defaultAction := data.DefaultAction.ValueString()

	version, err := getModifiableConfigVersion(ctx, int(configID), resourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read latest config version from API", err.Error())
		return
	}

	enableRapidRules := buildUpdateRapidRulesStatusRequest(data, version, true)
	client := inst.Client(r.meta)
	_, err = client.UpdateRapidRulesStatus(ctx, enableRapidRules)
	if err != nil {
		resp.Diagnostics.AddError(updateRapidRulesStatusError, err.Error())
		return
	}

	if defaultAction != "" {
		updateDefaultAction := buildUpdateRapidRulesDefaultActionRequest(data, version, defaultAction)
		_, err = client.UpdateRapidRulesDefaultAction(ctx, updateDefaultAction)
		if err != nil {
			resp.Diagnostics.AddError(updateRapidRulesDefaultActionFailed, err.Error())
			return
		}
	} else {
		getRapidRulesDefaultActionRequest := buildGetRapidRuleDefaultActionRequest(data, version)
		defaultActionResponse, err := client.GetRapidRulesDefaultAction(ctx, getRapidRulesDefaultActionRequest)
		if err != nil {
			resp.Diagnostics.AddError("calling 'getRapidRulesStatus'", err.Error())
			return
		}
		data.DefaultAction = types.StringValue(defaultActionResponse.Action)
	}

	ruleDefinitions, diags := getRuleDefinitionsFromModel(data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if len(ruleDefinitions) > 0 {
		getRulesRequest := buildGetRapidRulesRequest(data, version)
		rules, err := client.GetRapidRules(ctx, getRulesRequest)
		if err != nil {
			resp.Diagnostics.AddError(readRapidRulesError, err.Error())
			return
		}

		for _, definition := range ruleDefinitions {
			if definition.Action != nil {
				err = r.updateRapidRuleActionLock(ctx, data, *definition.ID, false, version)
				if err != nil {
					resp.Diagnostics.AddError("Update rapid rule action lock failure", err.Error())
					return
				}

				err = r.updateRapidRuleAction(ctx, data, rules, *definition.ID, *definition.Action, version)
				if err != nil {
					resp.Diagnostics.AddError("Update rapid rule action failure", err.Error())
					return
				}
			}

			if definition.Lock != nil {
				err = r.updateRapidRuleActionLock(ctx, data, *definition.ID, *definition.Lock, version)
				if err != nil {
					resp.Diagnostics.AddError("Update rapid rule action lock failure", err.Error())
					return
				}
			}

			if definition.ConditionException != nil {
				err = r.updateRapidRuleException(ctx, data, *definition.ID, version, *definition.ConditionException)
				if err != nil {
					resp.Diagnostics.AddError("Update rapid rule exception failure", err.Error())
					return
				}
			}
		}
	} else {
		var empty []appsec.RuleDefinition
		ruleDefinitionsJSON, err := serializeRuleDefinitions(empty)
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), "")
			return
		}
		data.RuleDefinitions = rapidRuleDefinitionsStateValue{types.StringValue(*ruleDefinitionsJSON)}
	}

	populateResourceID(data)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Read implements resource's Read method
func (r *rapidRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading Rapid Rules Resource")

	var data *rapidRulesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), r.meta)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version: ", err.Error())
		return
	}
	getRulesRequest := buildGetRapidRulesRequest(data, version)

	getRapidRulesStatusRequest := appsec.GetRapidRulesStatusRequest{
		ConfigID: getRulesRequest.ConfigID,
		Version:  getRulesRequest.Version,
		PolicyID: getRulesRequest.PolicyID,
	}

	client := inst.Client(r.meta)
	status, err := client.GetRapidRulesStatus(ctx, getRapidRulesStatusRequest)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read rapid rules status", err.Error())
		return
	}

	data.Enabled = types.BoolValue(status.Enabled)

	defaultAction, diags := r.readDefaultAction(ctx, status.Enabled, data, version)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	ruleDefinitionsState, diags := getRuleDefinitionsFromModel(data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	getRapidRulesRequest := buildGetRapidRulesRequest(data, version)
	rulesResponse, err := client.GetRapidRules(ctx, getRapidRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError(readRapidRulesError, err.Error())
		return
	}

	ruleDefinitionsRemoteState := toRuleDefinitions(rulesResponse)

	var definitions []appsec.RuleDefinition
	for _, definition := range ruleDefinitionsState {
		for _, val := range ruleDefinitionsRemoteState {
			if *definition.ID == *val.ID {
				definitions = append(definitions, val)
			}
		}
	}

	ruleDefinitionsJSON, err := serializeRuleDefinitions(definitions)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	populateState(data, defaultAction, *ruleDefinitionsJSON)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Update implements resource's Update method
func (r *rapidRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *rapidRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := plan.ConfigID.ValueInt64()
	defaultAction := plan.DefaultAction.ValueString()

	version, err := getModifiableConfigVersion(ctx, int(configID), resourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version: ", err.Error())
		return
	}

	client := inst.Client(r.meta)
	if plan.Enabled.ValueBool() != state.Enabled.ValueBool() {
		enableRapidRules := buildUpdateRapidRulesStatusRequest(plan, version, true)
		_, err = client.UpdateRapidRulesStatus(ctx, enableRapidRules)
		if err != nil {
			resp.Diagnostics.AddError(updateRapidRulesStatusError, err.Error())
			return
		}
	}

	if defaultAction != "" && areDefaultActionsDifferent(defaultAction, state.DefaultAction.ValueString()) {
		updateDefaultAction := buildUpdateRapidRulesDefaultActionRequest(plan, version, defaultAction)
		_, err = client.UpdateRapidRulesDefaultAction(ctx, updateDefaultAction)
		if err != nil {
			resp.Diagnostics.AddError(updateRapidRulesDefaultActionFailed, err.Error())
			return
		}
	}

	ruleDefinitionsOldState, diags := getRuleDefinitionsFromModel(state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	ruleDefinitionsPlan, diags := getRuleDefinitionsFromModel(plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	getRapidRulesRequest := buildGetRapidRulesRequest(plan, version)
	rulesResponse, err := client.GetRapidRules(ctx, getRapidRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError(readRapidRulesError, err.Error())
		return
	}

	if ruleDefinitionsOldState == nil && ruleDefinitionsPlan == nil {
		populateResourceID(plan)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	diags = r.manageRulesRemovedFromPlan(ctx, ruleDefinitionsOldState, ruleDefinitionsPlan, plan, version, rulesResponse)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if shouldUpdateRapidRules(&ruleDefinitionsPlan) {
		for _, definition := range ruleDefinitionsPlan {
			oldStateDefinition := getRuleDefinition(*definition.ID, &ruleDefinitionsOldState)

			if shouldUpdateRapidRuleAction(definition, oldStateDefinition) {
				err = r.updateRapidRuleActionLock(ctx, plan, *definition.ID, false, version)
				if err != nil {
					resp.Diagnostics.AddError("Update rapid rule action lock failure", err.Error())
					return
				}

				err = r.updateRapidRuleAction(ctx, plan, rulesResponse, *definition.ID, *definition.Action, version)
				if err != nil {
					resp.Diagnostics.AddError("Update rapid rule action failure", err.Error())
					return
				}
			}

			if shouldUpdateRapidRuleActionLock(definition, oldStateDefinition) {
				err = r.updateRapidRuleActionLock(ctx, plan, *definition.ID, *definition.Lock, version)
				if err != nil {
					resp.Diagnostics.AddError("Update rapid rule action lock failure", err.Error())
					return
				}
			}

			if shouldUpdateRapidRuleException(definition, oldStateDefinition) {
				err = r.updateRapidRuleException(ctx, plan, *definition.ID, version, *definition.ConditionException)
				if err != nil {
					resp.Diagnostics.AddError("update rapid rule exception failure:", err.Error())
					return
				}
			}
		}
	}

	populateResourceID(plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource's Delete method
func (r *rapidRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting Rapid Rules Resource")

	var data *rapidRulesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := int(data.ConfigID.ValueInt64())

	version, err := getModifiableConfigVersion(ctx, int(configID), resourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version: ", err.Error())
		return
	}

	disableRapidRules := buildUpdateRapidRulesStatusRequest(data, version, false)
	client := inst.Client(r.meta)
	_, err = client.UpdateRapidRulesStatus(ctx, disableRapidRules)
	if err != nil {
		resp.Diagnostics.AddError(updateRapidRulesStatusError, err.Error())
		return
	}
}

// ImportState implements resource's ImportState method
func (r *rapidRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing Rapid Rules resource")

	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 {
		resp.Diagnostics.AddError(fmt.Sprintf("ID '%s' incorrectly formatted: should be 'CONFIG_ID:SECURITY_POLICY_ID'", req.ID), "")
		return
	}

	configID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid configuration id '%v'", parts[0]), "")
		return
	}

	policyID := parts[1]
	if policyID == "" {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid security policy id '%v'", parts[1]), "")
		return
	}

	version, err := getLatestConfigVersion(ctx, int(configID), r.meta)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version: ", err.Error())
		return
	}

	data := rapidRulesResourceModel{
		ID:       types.StringValue(req.ID),
		ConfigID: types.Int64Value(configID),
		PolicyID: types.StringValue(policyID),
	}

	enableRapidRulesRequest := buildUpdateRapidRulesStatusRequest(&data, version, true)
	client := inst.Client(r.meta)
	_, err = client.UpdateRapidRulesStatus(ctx, enableRapidRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError(updateRapidRulesStatusError, err.Error())
		return
	}

	getRapidRulesDefaultActionRequest := buildGetRapidRuleDefaultActionRequest(&data, version)
	defaultActionResponse, err := client.GetRapidRulesDefaultAction(ctx, getRapidRulesDefaultActionRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'getRapidRulesStatus'", err.Error())
		return
	}

	getRapidRulesRequest := appsec.GetRapidRulesRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	rulesResponse, err := client.GetRapidRules(ctx, getRapidRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError(readRapidRulesError, err.Error())
		return
	}

	ruleDefinitions := toRuleDefinitions(rulesResponse)
	ruleDefinitionsJSON, err := serializeIndentRuleDefinitions(ruleDefinitions)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), "")
		return
	}

	populateState(&data, defaultActionResponse.Action, *ruleDefinitionsJSON)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *rapidRulesResource) manageRulesRemovedFromPlan(ctx context.Context, ruleDefinitionsOldState []appsec.RuleDefinition, ruleDefinitionsPlan []appsec.RuleDefinition, plan *rapidRulesResourceModel, version int, rulesResponse *appsec.GetRapidRulesResponse) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, definitionOldState := range ruleDefinitionsOldState {
		definitionPlan := getRuleDefinition(*definitionOldState.ID, &ruleDefinitionsPlan)

		if shouldRapidRuleRemovedFromPlan(definitionPlan) {
			err := r.updateRapidRuleActionLock(ctx, plan, *definitionOldState.ID, false, version)
			if err != nil {
				diags.AddError("Update rapid rule action lock failure", err.Error())
			}

			err = r.updateRapidRuleAction(ctx, plan, rulesResponse, *definitionOldState.ID, "none", version)
			if err != nil {
				diags.AddError("Update rapid rule action failure", err.Error())
			}

			err = r.updateRapidRuleActionLock(ctx, plan, *definitionOldState.ID, false, version)
			if err != nil {
				diags.AddError("Update rapid rule action lock failure", err.Error())
			}
		}

		if shouldExceptionRemovedFromPlan(definitionOldState, definitionPlan) {
			err := r.updateRapidRuleException(ctx, plan, *definitionPlan.ID, version, appsec.RuleConditionException{})
			if err != nil {
				diags.AddError("update rapid rule exception failure:", err.Error())
			}
		}
	}
	return diags
}

func buildGetRapidRuleDefaultActionRequest(m *rapidRulesResourceModel, version int) appsec.GetRapidRulesDefaultActionRequest {
	return appsec.GetRapidRulesDefaultActionRequest{
		ConfigID: m.ConfigID.ValueInt64(),
		Version:  version,
		PolicyID: m.PolicyID.ValueString(),
	}
}

func buildGetRapidRulesRequest(m *rapidRulesResourceModel, version int) appsec.GetRapidRulesRequest {
	return appsec.GetRapidRulesRequest{
		ConfigID: m.ConfigID.ValueInt64(),
		Version:  version,
		PolicyID: m.PolicyID.ValueString(),
	}
}

func buildUpdateRapidRulesStatusRequest(m *rapidRulesResourceModel, version int, enabled bool) appsec.UpdateRapidRulesStatusRequest {
	return appsec.UpdateRapidRulesStatusRequest{
		ConfigID: m.ConfigID.ValueInt64(),
		Version:  version,
		PolicyID: m.PolicyID.ValueString(),
		Body: appsec.UpdateRapidRulesStatusRequestBody{
			Enabled: &enabled,
		},
	}
}

func buildUpdateRapidRulesDefaultActionRequest(m *rapidRulesResourceModel, version int, defaultAction string) appsec.UpdateRapidRulesDefaultActionRequest {
	return appsec.UpdateRapidRulesDefaultActionRequest{
		ConfigID: m.ConfigID.ValueInt64(),
		Version:  version,
		PolicyID: m.PolicyID.ValueString(),
		Body: appsec.UpdateRapidRulesDefaultActionRequestBody{
			Action: defaultAction,
		},
	}
}

func buildUpdateRapidRuleActionRequest(m *rapidRulesResourceModel, version int, ruleVersion int, ruleID int64, action string) appsec.UpdateRapidRuleActionRequest {
	return appsec.UpdateRapidRuleActionRequest{
		ConfigID:    m.ConfigID.ValueInt64(),
		Version:     version,
		PolicyID:    m.PolicyID.ValueString(),
		RuleID:      ruleID,
		RuleVersion: ruleVersion,
		Body: appsec.UpdateRapidRuleActionRequestBody{
			Action: action,
		},
	}
}

func buildUpdateRapidRuleActionLockRequest(m *rapidRulesResourceModel, version int, lock bool, ruleID int64) appsec.UpdateRapidRuleActionLockRequest {
	return appsec.UpdateRapidRuleActionLockRequest{
		ConfigID: m.ConfigID.ValueInt64(),
		Version:  version,
		PolicyID: m.PolicyID.ValueString(),
		RuleID:   ruleID,
		Body: appsec.UpdateRapidRuleActionLockRequestBody{
			Enabled: &lock,
		},
	}
}

func buildUpdateRapidRuleExceptionRequest(m *rapidRulesResourceModel, version int, ruleID int64, exception appsec.RuleConditionException) appsec.UpdateRapidRuleExceptionRequest {
	return appsec.UpdateRapidRuleExceptionRequest{
		ConfigID: m.ConfigID.ValueInt64(),
		Version:  version,
		PolicyID: m.PolicyID.ValueString(),
		RuleID:   ruleID,
		Body:     exception,
	}
}

func (r *rapidRulesResource) readDefaultAction(ctx context.Context, status bool, data *rapidRulesResourceModel, version int) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	defaultAction := "unknown"
	if !status {
		return defaultAction, diags
	}
	client := inst.Client(r.meta)
	getRapidRulesDefaultActionRequest := buildGetRapidRuleDefaultActionRequest(data, version)
	defaultActionResponse, err := client.GetRapidRulesDefaultAction(ctx, getRapidRulesDefaultActionRequest)
	if err != nil {
		diags.AddError("calling 'GetRapidRulesDefaultAction'", err.Error())
		return defaultAction, diags
	}
	return defaultActionResponse.Action, diags
}

func (r *rapidRulesResource) updateRapidRuleAction(ctx context.Context, data *rapidRulesResourceModel, rules *appsec.GetRapidRulesResponse, ruleID int64, action string, version int) error {
	ruleVersion, err := getLatestRapidRuleVersion(ruleID, rules)
	if err != nil {
		return err
	}

	updateActionReq := buildUpdateRapidRuleActionRequest(data, version, *ruleVersion, ruleID, action)
	client := inst.Client(r.meta)
	_, err = client.UpdateRapidRuleAction(ctx, updateActionReq)
	if err != nil {
		return fmt.Errorf("calling 'UpdateRapidRuleAction': %s", err.Error())
	}
	return nil
}

func (r *rapidRulesResource) updateRapidRuleActionLock(ctx context.Context, data *rapidRulesResourceModel, ruleID int64, lock bool, version int) error {
	updateLockReq := buildUpdateRapidRuleActionLockRequest(data, version, lock, ruleID)
	client := inst.Client(r.meta)
	_, err := client.UpdateRapidRuleActionLock(ctx, updateLockReq)
	if err != nil {
		return fmt.Errorf("calling 'UpdateRapidRuleActionLock': %s", err.Error())
	}
	return nil
}

func (r *rapidRulesResource) updateRapidRuleException(ctx context.Context, data *rapidRulesResourceModel, ruleID int64, version int, exception appsec.RuleConditionException) error {
	updateExceptionReq := buildUpdateRapidRuleExceptionRequest(data, version, ruleID, exception)
	client := inst.Client(r.meta)
	_, err := client.UpdateRapidRuleException(ctx, updateExceptionReq)
	if err != nil {
		return fmt.Errorf("calling 'UpdateRapidRuleException': %s", err.Error())
	}
	return nil
}

func populateResourceID(m *rapidRulesResourceModel) {
	m.ID = types.StringValue(fmt.Sprintf("%d:%s", m.ConfigID.ValueInt64(), m.PolicyID.ValueString()))
}

func populateState(m *rapidRulesResourceModel, defaultAction string, ruleDefinitions string) {
	populateResourceID(m)
	m.DefaultAction = types.StringValue(defaultAction)
	m.RuleDefinitions = rapidRuleDefinitionsStateValue{types.StringValue(ruleDefinitions)}
}

func toRuleDefinitions(input *appsec.GetRapidRulesResponse) []appsec.RuleDefinition {
	definitions := make([]appsec.RuleDefinition, 0, len(input.Rules))

	for _, rule := range input.Rules {
		outputRule := appsec.RuleDefinition{
			ID:                 ptr.To(rule.ID),
			Action:             ptr.To(rule.Action),
			Lock:               ptr.To(rule.Lock),
			ConditionException: nil,
		}
		if rule.ConditionException != nil && (rule.ConditionException.Exception != nil || rule.ConditionException.AdvancedExceptionsList != nil) {
			outputRule.ConditionException = rule.ConditionException
		}
		definitions = append(definitions, outputRule)
	}

	return definitions
}

func areDefaultActionsDifferent(firstDefaultAction, secondDefaultAction string) bool {
	return firstDefaultAction != secondDefaultAction
}

func getRuleDefinitionsFromModel(data *rapidRulesResourceModel) ([]appsec.RuleDefinition, diag.Diagnostics) {
	var diags diag.Diagnostics
	rules := data.RuleDefinitions.ValueString()

	if data.RuleDefinitions.IsNull() || data.RuleDefinitions.IsUnknown() {
		return nil, diags
	}

	ruleDefinitions, err := deserializeRuleDefinitions(rules)
	if err != nil {
		diags.AddError(err.Error(), "Invalid rule definition JSON file: The configuration contains undefined or unrecognized fields.\nPlease review and ensure the JSON conforms to the expected format.")
		return nil, diags
	}
	return *ruleDefinitions, diags
}

func getLatestRapidRuleVersion(ruleID int64, rules *appsec.GetRapidRulesResponse) (*int, error) {
	for _, val := range rules.Rules {
		if val.ID == ruleID {
			return &val.Version, nil
		}
	}
	return nil, fmt.Errorf("cannot find latest rule version: rapid rule with ID: %d doesn't exist", ruleID)
}

func getRuleDefinition(ruleID int64, definitions *[]appsec.RuleDefinition) *appsec.RuleDefinition {
	if definitions == nil || len(*definitions) == 0 {
		return nil
	}
	for _, val := range *definitions {
		if *val.ID == ruleID {
			return &val
		}
	}
	return nil
}

func shouldUpdateRapidRules(ruleDefinitionsPlan *[]appsec.RuleDefinition) bool {
	return ruleDefinitionsPlan != nil && len(*ruleDefinitionsPlan) > 0
}

func shouldUpdateRapidRuleAction(definition appsec.RuleDefinition, oldStateDefinition *appsec.RuleDefinition) bool {
	return definition.Action != nil && (oldStateDefinition == nil || definition.Action != oldStateDefinition.Action)
}

func shouldUpdateRapidRuleActionLock(definition appsec.RuleDefinition, oldStateDefinition *appsec.RuleDefinition) bool {
	return definition.Lock != nil && (oldStateDefinition == nil || definition.Lock != oldStateDefinition.Lock)
}

func shouldUpdateRapidRuleException(definition appsec.RuleDefinition, oldStateDefinition *appsec.RuleDefinition) bool {
	return definition.ConditionException != nil && (oldStateDefinition == nil || definition.ConditionException != oldStateDefinition.ConditionException)
}

func shouldExceptionRemovedFromPlan(definitionOldState appsec.RuleDefinition, definitionPlan *appsec.RuleDefinition) bool {
	return definitionOldState.ConditionException != nil && definitionPlan != nil && definitionPlan.ConditionException == nil
}

func shouldRapidRuleRemovedFromPlan(definitionPlan *appsec.RuleDefinition) bool {
	return definitionPlan == nil
}

func extractErrors(diagnostics diag.Diagnostics) string {
	var errors []string
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity() == diag.SeverityError {
			errors = append(errors, diagnostic.Summary())
			if detail := diagnostic.Detail(); detail != "" {
				errors = append(errors, detail)
			}
		}
	}
	return strings.Join(errors, "\n")
}

func validateRuleDefinitions(definitions []appsec.RuleDefinition) diag.Diagnostics {
	var diags diag.Diagnostics
	if len(definitions) == 0 {
		return diags
	}

	for _, rule := range definitions {
		diags.Append(validateRuleDefinition(rule)...)
	}

	return diags
}

func validateRuleDefinition(v appsec.RuleDefinition) diag.Diagnostics {
	var diags diag.Diagnostics
	if v.ID == nil {
		diags.AddError(fmt.Sprintf(ruleDefinitionsValidationError, "rapid rule id: cannot be blank"), "")
	}
	diags.Append(validateRapidRuleAction(v.Action)...)
	if v.Lock == nil {
		diags.AddError(fmt.Sprintf(ruleDefinitionsValidationError, "rapid rule action lock: cannot be blank"), "")
	}
	return diags
}

func validateAction(action *string, allowedActions map[string]bool, errorFormat string) diag.Diagnostics {
	var diags diag.Diagnostics
	if action == nil {
		diags.AddError(fmt.Sprintf(errorFormat, "action: cannot be blank"), "")
		return diags
	}
	if allowedActions[*action] || strings.LastIndex(*action, "deny_custom_") == 0 {
		return diags
	}
	diags.AddError(fmt.Sprintf(errorFormat, toAllowedActionsMessage(allowedActions)), "")
	return diags
}

func toAllowedActionsMessage(allowedActions map[string]bool) string {
	concatenatedKeys := "action may only contain "
	first := true
	for key := range allowedActions {
		if !first {
			concatenatedKeys += ", "
		}
		concatenatedKeys += key
		first = false
	}
	return concatenatedKeys
}

func validateRapidRuleAction(action *string) diag.Diagnostics {
	allowedActions := map[string]bool{"alert": true, "deny": true, "deny_custom_{custom_deny_id}": false, "none": true}
	return validateAction(action, allowedActions, ruleDefinitionsValidationError)
}

func validateDefaultAction(action *string) diag.Diagnostics {
	allowedActions := map[string]bool{"akamai_managed": true, "alert": true, "deny": true, "deny_custom_{custom_deny_id}": false, "none": true}
	return validateAction(action, allowedActions, "Rapid rules default action validation error: %s")
}

func serializeIndentRuleDefinitions(definitions []appsec.RuleDefinition) (*string, error) {
	jsonBody, err := json.MarshalIndent(definitions, "", "  ")
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil

}

func serializeRuleDefinitions(definitions []appsec.RuleDefinition) (*string, error) {
	jsonBody, err := json.Marshal(definitions)
	if err != nil {
		return nil, err
	}
	return ptr.To(string(jsonBody)), nil
}

func deserializeRuleDefinitions(body string) (*[]appsec.RuleDefinition, error) {
	var definitions []appsec.RuleDefinition

	if body != "" {
		decoder := json.NewDecoder(strings.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&definitions); err != nil {
			return nil, err
		}
	}

	return &definitions, nil
}
