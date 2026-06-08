package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                   = &wafRulesetResource{}
	_ resource.ResourceWithConfigure      = &wafRulesetResource{}
	_ resource.ResourceWithValidateConfig = &wafRulesetResource{}
	_ resource.ResourceWithImportState    = &wafRulesetResource{}
)

// wafRulesetResource represents akamai_appsec_waf_ruleset resource
type wafRulesetResource struct {
	meta meta.Meta
}

// wafRulesetResourceModel is a model for akamai_appsec_waf_ruleset resource
type wafRulesetResourceModel struct {
	ConfigID     types.Int64  `tfsdk:"config_id"`
	PolicyID     types.String `tfsdk:"security_policy_id"`
	Rules        types.Set    `tfsdk:"rules"`
	AttackGroups types.Set    `tfsdk:"attack_groups"`
}

// wafRuleResourceModel represents a single WAF rule configuration
type wafRuleResourceModel struct {
	RuleID             types.Int64                      `tfsdk:"rule_id"`
	RuleAction         types.String                     `tfsdk:"rule_action"`
	ConditionException ruleConditionExceptionStateValue `tfsdk:"condition_exception"`
}

// attackGroupResourceModel represents a single attack group configuration
type attackGroupResourceModel struct {
	AttackGroup        types.String                            `tfsdk:"attack_group"`
	AttackGroupAction  types.String                            `tfsdk:"attack_group_action"`
	ConditionException attackGroupConditionExceptionStateValue `tfsdk:"condition_exception"`
}

const (
	readWAFRulesetError                     = "Unable to read WAF ruleset"
	updateWAFCompositeRulesetError          = "Unable to update WAF ruleset"
	readConfigVersionError                  = "Unable to read latest config version from API"
	invalidWAFRulesetConfigurationAttribute = "Invalid configuration attribute"
	wafRulesetValidationError               = "WAF ruleset validation error: %s"
	wafRulesetResourceName                  = "wafRuleset"
)

// NewWAFRulesetResource returns new appsec WAF ruleset resource
func NewWAFRulesetResource() resource.Resource {
	return &wafRulesetResource{}
}

// Metadata implements resource.Resource.
func (r *wafRulesetResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_waf_ruleset"
}

// Schema implements resource's Schema
func (r *wafRulesetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "WAF ruleset resource.",
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
					stringplanmodifier.UseStateForUnknown(),
					modifiers.PreventStringUpdateIfKnown("security_policy_id"),
				},
			},
			"rules": schema.SetNestedAttribute{
				Optional:    true,
				Description: "List of rule objects including action and condition exceptions",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.Int64Attribute{
							Required:    true,
							Description: "Unique identifier for a rule",
						},
						"rule_action": schema.StringAttribute{
							Required:    true,
							Description: "Action taken when the rule is triggered (alert, deny, deny_custom_{custom_deny_id}, none)",
							Validators:  []validator.String{validators.NotEmptyString()},
						},
						"condition_exception": schema.StringAttribute{
							CustomType:  ruleConditionExceptionStateType{},
							Optional:    true,
							Description: "Conditions and exceptions associated with the rule",
							PlanModifiers: []planmodifier.String{
								NormalizeRuleConditionException(),
							},
						},
					},
				},
			},
			"attack_groups": schema.SetNestedAttribute{
				Optional:    true,
				Description: "List of attack group objects including action and condition exceptions",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attack_group": schema.StringAttribute{
							Required:    true,
							Description: "Unique name of the attack group",
							Validators:  []validator.String{validators.NotEmptyString()},
						},
						"attack_group_action": schema.StringAttribute{
							Required:    true,
							Description: "Action taken when the attack group is triggered (alert, deny, deny_custom_{custom_deny_id}, none)",
							Validators:  []validator.String{validators.NotEmptyString()},
						},
						"condition_exception": schema.StringAttribute{
							CustomType:  attackGroupConditionExceptionStateType{},
							Description: "JSON-formatted conditions and exceptions associated with the attack group",
							Optional:    true,
							PlanModifiers: []planmodifier.String{
								NormalizeAttackGroupConditionException(),
							},
						},
					},
				},
			},
		},
	}
}

// rulesType returns the ObjectType for the rules nested attribute
func rulesType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"rule_id":             types.Int64Type,
			"rule_action":         types.StringType,
			"condition_exception": ruleConditionExceptionStateType{},
		},
	}
}

// attackGroupsType returns the ObjectType for the attack_groups nested attribute
func attackGroupsType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"attack_group":        types.StringType,
			"attack_group_action": types.StringType,
			"condition_exception": attackGroupConditionExceptionStateType{},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *wafRulesetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *wafRulesetResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data wafRulesetResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate rules
	if tf.IsKnown(data.Rules) {
		var rules []wafRuleResourceModel
		resp.Diagnostics.Append(data.Rules.ElementsAs(ctx, &rules, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(rules) > 0 {
			diags := validateWAFRules(rules)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(path.Root("rules"), invalidWAFRulesetConfigurationAttribute, extractErrors(diags.Errors()))
			}
		}
	}

	// Validate attack groups
	if tf.IsKnown(data.AttackGroups) {
		var attackGroups []attackGroupResourceModel
		resp.Diagnostics.Append(data.AttackGroups.ElementsAs(ctx, &attackGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(attackGroups) > 0 {
			diags := validateAttackGroups(attackGroups)
			if diags.HasError() {
				resp.Diagnostics.AddAttributeError(path.Root("attack_groups"), invalidWAFRulesetConfigurationAttribute, extractErrors(diags.Errors()))
			}
		}
	}
}

// Create implements resource's Create method
func (r *wafRulesetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating WAF Ruleset Resource")

	var data *wafRulesetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()

	// Get modifiable version of the configuration
	version, err := getModifiableConfigVersion(ctx, int(configID), wafRulesetResourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError(readConfigVersionError, err.Error())
		return
	}

	client := inst.Client(r.meta)

	// Build composite ruleset update request with all rules and attack groups
	updateRequest := appsec.UpdateWAFCompositeRulesetRequest{
		ConfigID: configID,
		Version:  int64(version),
		PolicyID: data.PolicyID.ValueString(),
	}

	// Fetch current WAF ruleset to validate existence of rules and attack groups, which can be updated
	getWAFRulesetRequest := buildGetWAFRulesetRequest(data, version)
	currentRuleset, err := client.GetWAFCompositeRuleset(ctx, getWAFRulesetRequest)
	if err != nil {
		resp.Diagnostics.AddError(readWAFRulesetError, err.Error())
		return
	}

	// Convert rules from plan to API format
	if tf.IsKnown(data.Rules) {
		var rules []wafRuleResourceModel
		resp.Diagnostics.Append(data.Rules.ElementsAs(ctx, &rules, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(rules) > 0 {
			// Build sets of valid rule IDs and attack group names
			validRuleIDs := make(map[int64]bool)
			for _, rule := range currentRuleset.Rules {
				validRuleIDs[rule.RuleID] = true
			}

			ruleUpdates, diags := convertResourceModelToWAFCompositeRules(rules, validRuleIDs)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			updateRequest.Rules = ruleUpdates
		}
	}

	// Convert attack groups from plan to API format
	if tf.IsKnown(data.AttackGroups) {
		var attackGroups []attackGroupResourceModel
		resp.Diagnostics.Append(data.AttackGroups.ElementsAs(ctx, &attackGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(attackGroups) > 0 {
			validAttackGroups := make(map[string]bool)
			for _, ag := range currentRuleset.AttackGroups {
				validAttackGroups[ag.Group] = true
			}

			attackGroupUpdates, diags := convertResourceModelToWAFCompositeAttackGroups(attackGroups, validAttackGroups)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			updateRequest.AttackGroups = attackGroupUpdates
		}
	}

	wafRuleset, err := client.UpdateWAFCompositeRuleset(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(updateWAFCompositeRulesetError, err.Error())
		return
	}

	data, diags := populateWAFRulesetState(ctx, data, wafRuleset)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Read implements resource's Read method
func (r *wafRulesetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading WAF Ruleset Resource")

	var data *wafRulesetResourceModel

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

	getWAFRulesetRequest := buildGetWAFRulesetRequest(data, version)

	client := inst.Client(r.meta)
	wafRuleset, err := client.GetWAFCompositeRuleset(ctx, getWAFRulesetRequest)
	if err != nil {
		// If the WAF Ruleset or security policy is not found, remove the resource from state.
		// This may happen if the security policy was deleted outside of Terraform.
		if strings.Contains(err.Error(), "not found") {
			tflog.Warn(ctx, "WAF Ruleset not found, removing from state", map[string]interface{}{
				"config_id":          configID,
				"security_policy_id": data.PolicyID.ValueString(),
			})
			resp.Diagnostics.AddWarning("Could not find WAF Ruleset in config version", err.Error())
			// Remove resource from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(readWAFRulesetError, err.Error())
		return
	}

	updatedData, diags := populateWAFRulesetState(ctx, data, wafRuleset)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedData)...)
}

// Update implements resource's Update method
func (r *wafRulesetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating WAF Ruleset Resource")

	var plan, state *wafRulesetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := int(plan.ConfigID.ValueInt64())

	// Get modifiable version of the configuration
	version, err := getModifiableConfigVersion(ctx, configID, wafRulesetResourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError(readConfigVersionError, err.Error())
		return
	}

	client := inst.Client(r.meta)

	// Build composite ruleset update request
	updateRequest := appsec.UpdateWAFCompositeRulesetRequest{
		ConfigID: plan.ConfigID.ValueInt64(),
		Version:  int64(version),
		PolicyID: plan.PolicyID.ValueString(),
	}

	// Fetch current WAF ruleset to validate existence of rules and attack groups
	getWAFRulesetRequest := buildGetWAFRulesetRequest(plan, version)
	currentRuleset, err := client.GetWAFCompositeRuleset(ctx, getWAFRulesetRequest)
	if err != nil {
		resp.Diagnostics.AddError(readWAFRulesetError, err.Error())
		return
	}

	// Only process rules if they are configured in plan or state
	planRulesConfigured := tf.IsKnown(plan.Rules)
	stateRulesConfigured := tf.IsKnown(state.Rules)

	if planRulesConfigured || stateRulesConfigured {
		var planRules, stateRules []wafRuleResourceModel

		if planRulesConfigured {
			resp.Diagnostics.Append(plan.Rules.ElementsAs(ctx, &planRules, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if stateRulesConfigured {
			resp.Diagnostics.Append(state.Rules.ElementsAs(ctx, &stateRules, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		validRuleIDs := make(map[int64]bool)
		for _, rule := range currentRuleset.Rules {
			validRuleIDs[rule.RuleID] = true
		}

		// Determine which rules need to be updated
		ruleUpdates, ruleDiags := buildRuleUpdates(stateRules, planRules, validRuleIDs)
		resp.Diagnostics.Append(ruleDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.Rules = ruleUpdates
	}

	// Only process attack groups if they are configured in plan or state
	planAttackGroupsConfigured := tf.IsKnown(plan.AttackGroups)
	stateAttackGroupsConfigured := tf.IsKnown(state.AttackGroups)

	if planAttackGroupsConfigured || stateAttackGroupsConfigured {
		var planAttackGroups, stateAttackGroups []attackGroupResourceModel

		if planAttackGroupsConfigured {
			resp.Diagnostics.Append(plan.AttackGroups.ElementsAs(ctx, &planAttackGroups, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if stateAttackGroupsConfigured {
			resp.Diagnostics.Append(state.AttackGroups.ElementsAs(ctx, &stateAttackGroups, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		validAttackGroups := make(map[string]bool)
		for _, ag := range currentRuleset.AttackGroups {
			validAttackGroups[ag.Group] = true
		}

		// Determine which attack groups need to be updated
		attackGroupUpdates, groupDiags := buildAttackGroupUpdates(stateAttackGroups, planAttackGroups, validAttackGroups)
		resp.Diagnostics.Append(groupDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateRequest.AttackGroups = attackGroupUpdates
	}

	wafRuleset, err := client.UpdateWAFCompositeRuleset(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(updateWAFCompositeRulesetError, err.Error())
		return
	}

	// Update plan state with API response
	stateDiags := updatePlanWithAPIResponse(ctx, plan, wafRuleset, planRulesConfigured, planAttackGroupsConfigured)
	resp.Diagnostics.Append(stateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// updatePlanWithAPIResponse updates the plan state with API response data
func updatePlanWithAPIResponse(ctx context.Context, plan *wafRulesetResourceModel, wafRuleset *appsec.CompositeRulesetResponse, planRulesConfigured, planAttackGroupsConfigured bool) diag.Diagnostics {
	var diags diag.Diagnostics

	// Handle rules - filter API response to only managed resources
	if planRulesConfigured {
		var planRules []wafRuleResourceModel
		diags.Append(plan.Rules.ElementsAs(ctx, &planRules, false)...)
		if diags.HasError() {
			return diags
		}

		managedRules := filterManagedRules(wafRuleset.Rules, planRules)
		rules, stateDiags := convertAPIRulesToResourceModel(managedRules, planRules)
		if stateDiags.HasError() {
			return stateDiags
		}
		plan.Rules, diags = types.SetValueFrom(ctx, rulesType(), rules)
		if diags.HasError() {
			return diags
		}
	}

	// Handle attack groups - filter API response to only managed resources
	if planAttackGroupsConfigured {
		var planAttackGroups []attackGroupResourceModel
		diags.Append(plan.AttackGroups.ElementsAs(ctx, &planAttackGroups, false)...)
		if diags.HasError() {
			return diags
		}

		managedAttackGroups := filterManagedAttackGroups(wafRuleset.AttackGroups, planAttackGroups)
		attackGroups, stateDiags := convertAPIAttackGroupsToResourceModel(managedAttackGroups, planAttackGroups)
		if stateDiags.HasError() {
			return stateDiags
		}
		plan.AttackGroups, diags = types.SetValueFrom(ctx, attackGroupsType(), attackGroups)
		if diags.HasError() {
			return diags
		}
	}

	return diags
}

// Delete implements resource's Delete method.
// This method removes the WAF ruleset by resetting only the managed rules and attack groups to action="none".
func (r *wafRulesetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting WAF Ruleset Resource")

	var data *wafRulesetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := int(data.ConfigID.ValueInt64())

	// Get modifiable version of the configuration
	version, err := getModifiableConfigVersion(ctx, configID, wafRulesetResourceName, r.meta)
	if err != nil {
		resp.Diagnostics.AddError(readConfigVersionError, err.Error())
		return
	}

	client := inst.Client(r.meta)

	// Build update request to reset only managed rules and attack groups to action="none"
	updateRequest := appsec.UpdateWAFCompositeRulesetRequest{
		ConfigID: data.ConfigID.ValueInt64(),
		Version:  int64(version),
		PolicyID: data.PolicyID.ValueString(),
	}

	// Reset only managed rules to action="none" with no condition exceptions
	if tf.IsKnown(data.Rules) {
		var rules []wafRuleResourceModel
		resp.Diagnostics.Append(data.Rules.ElementsAs(ctx, &rules, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(rules) > 0 {
			ruleUpdates := make([]appsec.WAFCompositeRuleUpdate, 0, len(rules))
			for _, rule := range rules {
				ruleUpdates = append(ruleUpdates, appsec.WAFCompositeRuleUpdate{
					RuleID:             rule.RuleID.ValueInt64(),
					Action:             "none",
					ConditionException: nil,
				})
			}
			updateRequest.Rules = ruleUpdates
		}
	}

	// Reset only managed attack groups to action="none" with no condition exceptions
	if tf.IsKnown(data.AttackGroups) {
		var attackGroups []attackGroupResourceModel
		resp.Diagnostics.Append(data.AttackGroups.ElementsAs(ctx, &attackGroups, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if len(attackGroups) > 0 {
			attackGroupUpdates := make([]appsec.WAFCompositeAttackGroupUpdate, 0, len(attackGroups))
			for _, ag := range attackGroups {
				attackGroupUpdates = append(attackGroupUpdates, appsec.WAFCompositeAttackGroupUpdate{
					Group:              ag.AttackGroup.ValueString(),
					Action:             "none",
					ConditionException: nil,
				})
			}
			updateRequest.AttackGroups = attackGroupUpdates
		}
	}

	// Execute the update to reset only managed configurations
	_, err = client.UpdateWAFCompositeRuleset(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(updateWAFCompositeRulesetError, err.Error())
		return
	}

	tflog.Info(ctx, "WAF Ruleset resource removed from state")
}

// ImportState implements resource's ImportState method
func (r *wafRulesetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing WAF Ruleset resource")

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

	data := wafRulesetResourceModel{
		ConfigID: types.Int64Value(configID),
		PolicyID: types.StringValue(policyID),
	}

	getWAFRulesetRequest := appsec.GetWAFCompositeRulesetRequest{
		ConfigID: configID,
		Version:  int64(version),
		PolicyID: policyID,
	}

	client := inst.Client(r.meta)
	wafRuleset, err := client.GetWAFCompositeRuleset(ctx, getWAFRulesetRequest)
	if err != nil {
		resp.Diagnostics.AddError(readWAFRulesetError, err.Error())
		return
	}

	// Convert all rules and attack groups from API response
	rules, diags := importAPIRulesToResourceModel(wafRuleset.Rules)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	attackGroups, diags := importAPIAttackGroupsToResourceModel(wafRuleset.AttackGroups)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Populate state with all rules and attack groups
	data.Rules, diags = types.SetValueFrom(ctx, rulesType(), rules)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	data.AttackGroups, diags = types.SetValueFrom(ctx, attackGroupsType(), attackGroups)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// Helper functions

// buildGetWAFRulesetRequest builds the request for getting WAF ruleset
func buildGetWAFRulesetRequest(m *wafRulesetResourceModel, version int) appsec.GetWAFCompositeRulesetRequest {
	return appsec.GetWAFCompositeRulesetRequest{
		ConfigID: m.ConfigID.ValueInt64(),
		Version:  int64(version),
		PolicyID: m.PolicyID.ValueString(),
	}
}

// populateWAFRulesetState filters and converts API response to resource model, then populates the state
// Only populates rules/attack groups that are configured in the model
func populateWAFRulesetState(
	ctx context.Context,
	model *wafRulesetResourceModel,
	apiRuleset *appsec.CompositeRulesetResponse,
) (*wafRulesetResourceModel, diag.Diagnostics) {
	// Create a copy of the model to avoid modifying the input
	result := &wafRulesetResourceModel{
		ConfigID:     model.ConfigID,
		PolicyID:     model.PolicyID,
		Rules:        model.Rules,
		AttackGroups: model.AttackGroups,
	}

	var diags diag.Diagnostics

	// Only populate rules if they are configured in model
	if tf.IsKnown(model.Rules) {
		var modelRules []wafRuleResourceModel
		diags.Append(model.Rules.ElementsAs(ctx, &modelRules, false)...)
		if diags.HasError() {
			return nil, diags
		}

		// Filter API response to only include managed rules
		managedRules := filterManagedRules(apiRuleset.Rules, modelRules)
		var rules []wafRuleResourceModel

		rules, diags = convertAPIRulesToResourceModel(managedRules, modelRules)
		if diags.HasError() {
			return nil, diags
		}
		result.Rules, diags = types.SetValueFrom(ctx, rulesType(), rules)
		if diags.HasError() {
			return nil, diags
		}
	}

	// Only populate attack groups if they are configured in model
	if tf.IsKnown(model.AttackGroups) {
		var modelAttackGroups []attackGroupResourceModel
		diags.Append(model.AttackGroups.ElementsAs(ctx, &modelAttackGroups, false)...)
		if diags.HasError() {
			return nil, diags
		}

		// Filter API response to only include managed attack groups
		managedAttackGroups := filterManagedAttackGroups(apiRuleset.AttackGroups, modelAttackGroups)
		var attackGroups []attackGroupResourceModel

		attackGroups, diags = convertAPIAttackGroupsToResourceModel(managedAttackGroups, modelAttackGroups)
		if diags.HasError() {
			return nil, diags
		}
		result.AttackGroups, diags = types.SetValueFrom(ctx, attackGroupsType(), attackGroups)
		if diags.HasError() {
			return nil, diags
		}
	}

	return result, diags
}

// filterManagedRules filters API rules to only include those managed in the plan/config
func filterManagedRules(apiRules []appsec.WAFCompositeRule, managedRules []wafRuleResourceModel) []appsec.WAFCompositeRule {
	// Build map of managed rule IDs from plan/config
	managedRuleIDs := make(map[int64]bool)
	for _, rule := range managedRules {
		managedRuleIDs[rule.RuleID.ValueInt64()] = true
	}

	// Filter API response to only include managed rules
	filtered := make([]appsec.WAFCompositeRule, 0)
	for _, rule := range apiRules {
		if managedRuleIDs[rule.RuleID] {
			filtered = append(filtered, rule)
		}
	}

	return filtered
}

// filterManagedAttackGroups filters API attack groups to only include those managed in the plan/config
func filterManagedAttackGroups(apiAttackGroups []appsec.WAFCompositeAttackGroup, managedAttackGroups []attackGroupResourceModel) []appsec.WAFCompositeAttackGroup {
	// Build map of managed attack group names from plan/config
	managedGroupNames := make(map[string]bool)
	for _, ag := range managedAttackGroups {
		managedGroupNames[ag.AttackGroup.ValueString()] = true
	}

	// Filter API response to only include managed attack groups
	filtered := make([]appsec.WAFCompositeAttackGroup, 0)
	for _, ag := range apiAttackGroups {
		if managedGroupNames[ag.Group] {
			filtered = append(filtered, ag)
		}
	}

	return filtered
}

// importAPIRulesToResourceModel converts all API rules to resource model for import.
// Rules with action "none" are excluded as they are not managed.
func importAPIRulesToResourceModel(apiRules []appsec.WAFCompositeRule) ([]wafRuleResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	output := make([]wafRuleResourceModel, 0, len(apiRules))
	for _, rule := range apiRules {

		if rule.Action == "none" {
			// If action is "none", do not include them in the state as they are not managed
			continue
		}
		var conditionException ruleConditionExceptionStateValue

		conditionExceptionStr, err := marshalJSON(rule.ConditionException)
		if err != nil {
			diags.AddError("marshaling rule condition exception", err.Error())
			return output, diags
		} else if conditionExceptionStr != "" {
			conditionException = ruleConditionExceptionStateValue{StringValue: types.StringValue(conditionExceptionStr)}
		}

		outputRule := wafRuleResourceModel{
			RuleID:             types.Int64Value(rule.RuleID),
			RuleAction:         types.StringValue(rule.Action),
			ConditionException: conditionException,
		}
		output = append(output, outputRule)
	}

	if len(output) == 0 {
		return nil, diags
	}
	return output, diags
}

// convertAPIRulesToResourceModel converts API rules to resource model
// Rules are returned in the same order as existingRules to preserve user's config order & prevent drift caused by api response order changes
func convertAPIRulesToResourceModel(apiRules []appsec.WAFCompositeRule, existingRules []wafRuleResourceModel) ([]wafRuleResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build map of API rules by ID for quick lookup
	apiRulesMap := make(map[int64]appsec.WAFCompositeRule)
	for _, rule := range apiRules {
		apiRulesMap[rule.RuleID] = rule
	}

	// Build output in the same order as existingRules
	output := make([]wafRuleResourceModel, 0, len(existingRules))
	for _, existingRule := range existingRules {
		ruleID := existingRule.RuleID.ValueInt64()
		apiRule, found := apiRulesMap[ruleID]
		if !found {
			continue
		}

		var conditionException ruleConditionExceptionStateValue
		conditionExceptionStr, err := marshalJSON(apiRule.ConditionException)

		if err != nil {
			diags.AddError("marshaling rule condition exception", err.Error())
			return output, diags
		} else if conditionExceptionStr != "" {
			conditionException = ruleConditionExceptionStateValue{StringValue: types.StringValue(conditionExceptionStr)}
		}

		outputRule := wafRuleResourceModel{
			RuleID:             types.Int64Value(apiRule.RuleID),
			RuleAction:         types.StringValue(apiRule.Action),
			ConditionException: conditionException,
		}
		output = append(output, outputRule)
	}

	return output, diags
}

// importAPIAttackGroupsToResourceModel converts all API attack groups to resource model for import.
// Attack groups with action "none" are excluded as they are not managed.
func importAPIAttackGroupsToResourceModel(apiAttackGroups []appsec.WAFCompositeAttackGroup) ([]attackGroupResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	output := make([]attackGroupResourceModel, 0, len(apiAttackGroups))
	for _, attackGroup := range apiAttackGroups {

		if attackGroup.Action == "none" {
			// If action is "none", do not include them in the state as they are not managed
			continue
		}

		var conditionException attackGroupConditionExceptionStateValue

		conditionExceptionStr, err := marshalJSON(attackGroup.ConditionException)
		if err != nil {
			diags.AddError("marshaling attack group condition exception", err.Error())
			return output, diags
		} else if conditionExceptionStr != "" {
			conditionException = attackGroupConditionExceptionStateValue{StringValue: types.StringValue(conditionExceptionStr)}
		}

		outputAttackGroup := attackGroupResourceModel{
			AttackGroup:        types.StringValue(attackGroup.Group),
			AttackGroupAction:  types.StringValue(attackGroup.Action),
			ConditionException: conditionException,
		}
		output = append(output, outputAttackGroup)
	}

	if len(output) == 0 {
		return nil, diags
	}
	return output, diags
}

// convertAPIAttackGroupsToResourceModel converts API attack groups to resource model
// Attack groups are returned in the same order as existingGroups to preserve user's config order & prevent drift caused by api response order changes
func convertAPIAttackGroupsToResourceModel(apiAttackGroups []appsec.WAFCompositeAttackGroup, existingGroups []attackGroupResourceModel) ([]attackGroupResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Build map of API attack groups by name for quick lookup
	apiGroupsMap := make(map[string]appsec.WAFCompositeAttackGroup)
	for _, ag := range apiAttackGroups {
		apiGroupsMap[ag.Group] = ag
	}

	// Build output in the same order as existingGroups
	output := make([]attackGroupResourceModel, 0, len(existingGroups))
	for _, existingGroup := range existingGroups {
		groupName := existingGroup.AttackGroup.ValueString()
		apiGroup, found := apiGroupsMap[groupName]
		if !found {
			// Attack group no longer exists in API response, skip it
			continue
		}

		var conditionException attackGroupConditionExceptionStateValue
		conditionExceptionStr, err := marshalJSON(apiGroup.ConditionException)
		if err != nil {
			diags.AddError("marshaling attack group condition exception", err.Error())
			return output, diags
		} else if conditionExceptionStr != "" {
			conditionException = attackGroupConditionExceptionStateValue{StringValue: types.StringValue(conditionExceptionStr)}
		}

		outputAttackGroup := attackGroupResourceModel{
			AttackGroup:        types.StringValue(apiGroup.Group),
			AttackGroupAction:  types.StringValue(apiGroup.Action),
			ConditionException: conditionException,
		}
		output = append(output, outputAttackGroup)
	}

	return output, diags
}

func marshalJSON(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}

	// Use compact JSON (no indentation)
	jsonBytes, err := json.Marshal(v)
	if err != nil || string(jsonBytes) == "{}" {
		return "", err
	}

	return string(jsonBytes), nil
}

// validateWAFRules validates WAF rules configuration
func validateWAFRules(rules []wafRuleResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	ruleIDMap := make(map[int64]int)

	for i, rule := range rules {

		// Check for duplicate rule_ids
		if !rule.RuleID.IsNull() && rule.RuleID.ValueInt64() != 0 {
			ruleID := rule.RuleID.ValueInt64()
			if firstIndex, exists := ruleIDMap[ruleID]; exists {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("duplicate rule_id %d found at indices %d and %d", ruleID, firstIndex, i)),
					"Each rule_id must be unique within the rules list.",
				)
			} else {
				ruleIDMap[ruleID] = i
			}
		}

		if rule.RuleID.IsNull() || rule.RuleID.ValueInt64() == 0 {
			diags.AddError(fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("rule[%d]: rule_id cannot be empty or zero", i)), "")
		}

		if tf.IsKnown(rule.RuleAction) {
			action := rule.RuleAction.ValueString()
			if !isValidWAFAction(action) {
				diags.AddError(fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("rule[%d]: invalid rule_action '%s'", i, action)), "Action must be one of: alert, deny, deny_custom_{custom_deny_id}, none")
			}
		}

		if tf.IsKnown(rule.ConditionException) {
			conditionException := rule.ConditionException.ValueString()
			if conditionException == "" || conditionException == "{}" {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("rule[%d]: condition_exception cannot be empty or an empty object", i)),
					"Condition exception must contain valid conditions or exceptions, or be omitted entirely.",
				)
			} else if err := validateRuleConditionExceptionSchema(conditionException); err != nil {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("rule[%d]: condition_exception must be valid JSON matching the rule condition exception schema", i)),
					err.Error(),
				)
			}
		}

		// Cross-field validation: action="none" cannot have condition_exception
		if tf.IsKnown(rule.RuleAction) && tf.IsKnown(rule.ConditionException) {
			action := rule.RuleAction.ValueString()
			conditionException := rule.ConditionException.ValueString()

			if action == "none" && conditionException != "" {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("rule[%d]: condition exceptions cannot be applied when action is 'none'", i)),
					"Remove the condition_exception or change the rule_action to a value other than 'none'.",
				)
			}
		}
	}

	return diags
}

// validateAttackGroups validates attack groups configuration
func validateAttackGroups(attackGroups []attackGroupResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	groupNameMap := make(map[string]int)

	for i, attackGroup := range attackGroups {
		// Check for duplicate attack_group names
		if !attackGroup.AttackGroup.IsNull() && attackGroup.AttackGroup.ValueString() != "" {
			groupName := attackGroup.AttackGroup.ValueString()
			if firstIndex, exists := groupNameMap[groupName]; exists {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("duplicate attack_group '%s' found at indices %d and %d", groupName, firstIndex, i)),
					"Each attack_group name must be unique within the attack_groups list.",
				)
			} else {
				groupNameMap[groupName] = i
			}
		}

		if attackGroup.AttackGroup.IsNull() || attackGroup.AttackGroup.ValueString() == "" {
			diags.AddError(fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("attack_group[%d]: attack_group cannot be empty", i)), "")
		}

		if tf.IsKnown(attackGroup.AttackGroupAction) {
			action := attackGroup.AttackGroupAction.ValueString()
			if !isValidWAFAction(action) {
				diags.AddError(fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("attack_group[%d]: invalid attack_group_action '%s'", i, action)), "Action must be one of: alert, deny, deny_custom_{custom_deny_id}, none")
			}
		}

		if tf.IsKnown(attackGroup.ConditionException) {
			conditionException := attackGroup.ConditionException.ValueString()
			if conditionException == "" || conditionException == "{}" {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("attack_group[%d]: condition_exception cannot be empty or an empty object", i)),
					"Condition exception must contain valid conditions or exceptions, or be omitted entirely.",
				)
			} else if err := validateAttackGroupConditionExceptionSchema(conditionException); err != nil {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("attack_group[%d]: condition_exception must be valid JSON matching the attack group condition exception schema", i)),
					err.Error(),
				)
			}
		}

		// Cross-field validation: action="none" cannot have condition_exception
		if tf.IsKnown(attackGroup.AttackGroupAction) && tf.IsKnown(attackGroup.ConditionException) {
			action := attackGroup.AttackGroupAction.ValueString()
			conditionException := attackGroup.ConditionException.ValueString()

			if action == "none" && conditionException != "" {
				diags.AddError(
					fmt.Sprintf(wafRulesetValidationError, fmt.Sprintf("attack_group[%d]: condition exceptions cannot be applied when action is 'none'", i)),
					"Remove the condition_exception or change the attack_group_action to a value other than 'none'.",
				)
			}
		}
	}

	return diags
}

// isValidWAFAction checks if the action is valid for WAF rules/attack groups
func isValidWAFAction(action string) bool {
	validActions := map[string]bool{
		"alert": true,
		"deny":  true,
		"none":  true,
	}

	// Check for custom deny actions (deny_custom_{custom_deny_id})
	if strings.HasPrefix(action, "deny_custom_") {
		return true
	}

	return validActions[action]
}

// validateRuleConditionExceptionSchema validates that JSON matches appsec.RuleConditionException schema
func validateRuleConditionExceptionSchema(value string) error {
	var ruleConditionException appsec.RuleConditionException
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.DisallowUnknownFields()
	return decoder.Decode(&ruleConditionException)
}

// validateAttackGroupConditionExceptionSchema validates that JSON matches appsec.AttackGroupConditionException schema
func validateAttackGroupConditionExceptionSchema(value string) error {
	var attackGroupConditionException appsec.AttackGroupConditionException
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.DisallowUnknownFields()
	return decoder.Decode(&attackGroupConditionException)
}

// convertResourceModelToWAFCompositeRules converts resource model rules to WAF composite rule update format
func convertResourceModelToWAFCompositeRules(rules []wafRuleResourceModel, validRuleIDs map[int64]bool) ([]appsec.WAFCompositeRuleUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics
	output := make([]appsec.WAFCompositeRuleUpdate, 0, len(rules))

	for _, rule := range rules {
		ruleID := rule.RuleID.ValueInt64()
		action := rule.RuleAction.ValueString()
		conditionExceptionStr := rule.ConditionException.ValueString()

		// Validate that the rule exists in the WAF ruleset
		if !validRuleIDs[ruleID] {
			diags.AddError(
				fmt.Sprintf(
					"Rule %d does not exist in the WAF ruleset for the specified config version and policy. "+
						"Terraform can only manage rules that are already present in the ruleset. "+
						"Remove this rule from the Terraform configuration or ensure it exists in the WAF ruleset.",
					ruleID,
				),
				"",
			)
			return output, diags
		}

		// Parse condition exception if provided, otherwise set to nil
		var conditionException *appsec.RuleConditionException
		if conditionExceptionStr != "" {
			conditionException = &appsec.RuleConditionException{}
			if err := json.Unmarshal([]byte(conditionExceptionStr), conditionException); err != nil {
				diags.AddError(fmt.Sprintf("parsing rule %d condition exception", ruleID), err.Error())
				return output, diags
			}
		} else if action != "none" {
			// if condition exception is not set for rule when action is not none, clear the condition exception set
			conditionException = &appsec.RuleConditionException{}
		}

		ruleUpdate := appsec.WAFCompositeRuleUpdate{
			RuleID:             ruleID,
			Action:             action,
			ConditionException: conditionException,
		}
		output = append(output, ruleUpdate)
	}

	return output, diags
}

// convertResourceModelToWAFCompositeAttackGroups converts resource model attack groups to WAF composite attack group update format
func convertResourceModelToWAFCompositeAttackGroups(attackGroups []attackGroupResourceModel, validAttackGroups map[string]bool) ([]appsec.WAFCompositeAttackGroupUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics
	output := make([]appsec.WAFCompositeAttackGroupUpdate, 0, len(attackGroups))

	for _, ag := range attackGroups {
		groupName := ag.AttackGroup.ValueString()
		action := ag.AttackGroupAction.ValueString()
		conditionExceptionStr := ag.ConditionException.ValueString()

		// Validate that the attack group exists in the WAF ruleset
		if !validAttackGroups[groupName] {
			diags.AddError(
				fmt.Sprintf(
					"Attack group %s does not exist in the WAF ruleset for the specified config version and policy. "+
						"Terraform can only manage attack groups that are already present in the ruleset. "+
						"Remove this attack group from the Terraform configuration or ensure it exists in the WAF ruleset.",
					groupName,
				), "",
			)
			return output, diags
		}

		// Parse condition exception if provided, otherwise set to nil
		var conditionException *appsec.AttackGroupConditionException
		if conditionExceptionStr != "" {
			conditionException = &appsec.AttackGroupConditionException{}
			if err := json.Unmarshal([]byte(conditionExceptionStr), conditionException); err != nil {
				diags.AddError(fmt.Sprintf("parsing attack group %s condition exception", groupName), err.Error())
				return output, diags
			}
		} else if action != "none" {
			// if condition exception is not set for attack group when action is not none, clear the condition exception set
			conditionException = &appsec.AttackGroupConditionException{}
		}

		attackGroupUpdate := appsec.WAFCompositeAttackGroupUpdate{
			Group:              groupName,
			Action:             action,
			ConditionException: conditionException,
		}
		output = append(output, attackGroupUpdate)
	}

	return output, diags
}

// buildRuleUpdates determines which rules need to be updated based on state vs plan
func buildRuleUpdates(stateRules, planRules []wafRuleResourceModel, validRuleIDs map[int64]bool) ([]appsec.WAFCompositeRuleUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics
	updates := make([]appsec.WAFCompositeRuleUpdate, 0)

	// Create maps rules for lookup
	stateRulesMap := make(map[int64]wafRuleResourceModel)
	for _, rule := range stateRules {
		stateRulesMap[rule.RuleID.ValueInt64()] = rule
	}

	planRulesMap := make(map[int64]wafRuleResourceModel)
	for _, rule := range planRules {
		planRulesMap[rule.RuleID.ValueInt64()] = rule
	}

	// Add/update rules from plan
	for _, planRule := range planRules {
		ruleID := planRule.RuleID.ValueInt64()
		stateRule, existsInState := stateRulesMap[ruleID]

		// Add rule if it changed or doesn't exist in state
		if !existsInState || ruleHasChanged(stateRule, planRule) {
			ruleUpdates, ruleDiags := convertResourceModelToWAFCompositeRules([]wafRuleResourceModel{planRule}, validRuleIDs)
			if ruleDiags.HasError() {
				diags.Append(ruleDiags...)
				return updates, diags
			}
			updates = append(updates, ruleUpdates...)
		}
	}

	// Handle rules removed from plan: set action="none"
	for _, stateRule := range stateRules {
		ruleID := stateRule.RuleID.ValueInt64()
		_, existsInPlan := planRulesMap[ruleID]

		// If rule exists in state but not in plan, reset it
		if !existsInPlan {
			updates = append(updates, appsec.WAFCompositeRuleUpdate{
				RuleID:             ruleID,
				Action:             "none",
				ConditionException: nil,
			})
		}
	}

	return updates, diags
}

// buildAttackGroupUpdates determines which attack groups need to be updated based on state vs plan
func buildAttackGroupUpdates(stateGroups, planGroups []attackGroupResourceModel, validAttackGroups map[string]bool) ([]appsec.WAFCompositeAttackGroupUpdate, diag.Diagnostics) {
	var diags diag.Diagnostics
	updates := make([]appsec.WAFCompositeAttackGroupUpdate, 0)

	// Create maps for easy lookup
	stateGroupsMap := make(map[string]attackGroupResourceModel)
	for _, group := range stateGroups {
		stateGroupsMap[group.AttackGroup.ValueString()] = group
	}

	planGroupsMap := make(map[string]attackGroupResourceModel)
	for _, group := range planGroups {
		planGroupsMap[group.AttackGroup.ValueString()] = group
	}

	// Add/update attack groups from plan
	for _, planGroup := range planGroups {
		groupName := planGroup.AttackGroup.ValueString()
		stateGroup, existsInState := stateGroupsMap[groupName]

		// Add attack group if it changed or doesn't exist in state
		if !existsInState || attackGroupHasChanged(stateGroup, planGroup) {
			groupUpdates, groupDiags := convertResourceModelToWAFCompositeAttackGroups([]attackGroupResourceModel{planGroup}, validAttackGroups)
			if groupDiags.HasError() {
				diags.Append(groupDiags...)
				return updates, diags
			}
			updates = append(updates, groupUpdates...)
		}
	}

	// Handle attack groups removed from plan: set action="none"
	for _, stateGroup := range stateGroups {
		groupName := stateGroup.AttackGroup.ValueString()
		_, existsInPlan := planGroupsMap[groupName]

		// If attack group exists in state but not in plan, reset it
		if !existsInPlan {
			updates = append(updates, appsec.WAFCompositeAttackGroupUpdate{
				Group:              groupName,
				Action:             "none",
				ConditionException: nil,
			})
		}
	}

	return updates, diags
}

// ruleHasChanged checks if a rule has changed between state and plan
func ruleHasChanged(stateRule, planRule wafRuleResourceModel) bool {
	if stateRule.RuleAction.ValueString() != planRule.RuleAction.ValueString() {
		return true
	}

	// Compare condition exceptions using semantic equality (deep equal on parsed JSON)
	equal, _ := stateRule.ConditionException.StringSemanticEquals(context.Background(), planRule.ConditionException)
	return !equal
}

// attackGroupHasChanged checks if an attack group has changed between state and plan
func attackGroupHasChanged(stateGroup, planGroup attackGroupResourceModel) bool {
	if stateGroup.AttackGroupAction.ValueString() != planGroup.AttackGroupAction.ValueString() {
		return true
	}

	// Compare condition exceptions using semantic equality (deep equal on parsed JSON)
	equal, _ := stateGroup.ConditionException.StringSemanticEquals(context.Background(), planGroup.ConditionException)
	return !equal
}
