package appsec

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	wafAIRulesResourceName  = "wafAIRules"
	aiRuleStatusEnabled     = "ENABLED"
	aiRuleStatusDisabled    = "DISABLED"
	aiRuleStatusNotEnrolled = "NOT_ENROLLED"
)

var (
	_ resource.Resource                   = &wafAIRulesResource{}
	_ resource.ResourceWithConfigure      = &wafAIRulesResource{}
	_ resource.ResourceWithValidateConfig = &wafAIRulesResource{}
	_ resource.ResourceWithImportState    = &wafAIRulesResource{}
)

type (
	wafAIRulesResource struct {
		meta.Resource
	}

	wafAIRulesResourceModel struct {
		ConfigID         types.Int64  `tfsdk:"config_id"`
		SecurityPolicyID types.String `tfsdk:"security_policy_id"`
		AIRuleStatus     types.String `tfsdk:"ai_rule_status"`
		RuleID           types.Int64  `tfsdk:"rule_id"`
		RuleVersionID    types.Int64  `tfsdk:"rule_version_id"`
		Action           types.String `tfsdk:"action"`
	}
)

// NewWAFAIRulesResource returns a new WAF AI rules resource.
func NewWAFAIRulesResource() resource.Resource { return &wafAIRulesResource{} }

// Metadata sets the resource type name.
func (r *wafAIRulesResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_waf_ai_rules"
}

// Schema defines the resource Terraform schema.
func (r *wafAIRulesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "WAF AI rules resource. This resource manages either the AI rules enable/disable status for a security policy or the action for a specific AI rule.",
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
				Validators:  []validator.String{validators.NotEmptyString()},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"ai_rule_status": schema.StringAttribute{
				Optional:    true,
				Description: "AI rules status for the policy. Valid values: ENABLED, DISABLED. NOT_ENROLLED is returned by the API when the policy is not enrolled but cannot be set. Mutually exclusive with rule_id and action.",
				Validators:  []validator.String{stringvalidator.OneOf(aiRuleStatusEnabled, aiRuleStatusDisabled)},
			},
			"rule_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Unique identifier of the AI rule whose action to manage. Mutually exclusive with ai_rule_status.",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"rule_version_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Version of the AI rule. Looked up automatically from the rule list if not provided.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					modifiers.PreventInt64Update(),
				},
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Description: "Action to set for the AI rule specified by rule_id. Valid values: alert, deny, deny_<custom_deny_name>, none. Mutually exclusive with ai_rule_status.",
				Validators:  []validator.String{stringvalidator.RegexMatches(regexp.MustCompile(`^(alert|deny|deny_custom_.+|none)$`), "must be alert, deny, deny_custom_<id>, or none")},
			},
		},
	}
}

// ValidateConfig enforces mutual exclusivity between status mode and action mode.
func (r *wafAIRulesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data wafAIRulesResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	statusKnown := !data.AIRuleStatus.IsUnknown()
	ruleKnown := !data.RuleID.IsUnknown()
	actionKnown := !data.Action.IsUnknown()

	statusSet := statusKnown && !data.AIRuleStatus.IsNull()
	ruleSet := ruleKnown && !data.RuleID.IsNull()
	actionSet := actionKnown && !data.Action.IsNull()

	// Check conflicts provable from known values regardless of other unknowns.
	if statusSet && ruleSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("ai_rule_status"),
			"Conflicting attributes",
			"ai_rule_status cannot be set together with rule_id or action.",
		)
		return
	}
	if statusSet && actionSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("ai_rule_status"),
			"Conflicting attributes",
			"ai_rule_status cannot be set together with rule_id or action.",
		)
		return
	}

	// Skip checks whose outcome depends on values not yet resolved.
	if !statusKnown || !ruleKnown || !actionKnown {
		return
	}

	// All values are known — do full validation.
	statusSetFull := !data.AIRuleStatus.IsNull()
	ruleSetFull := !data.RuleID.IsNull()
	actionSetFull := !data.Action.IsNull()

	if !statusSetFull && !ruleSetFull && !actionSetFull {
		resp.Diagnostics.AddError(
			"Missing required attributes",
			"Either ai_rule_status or both rule_id and action must be provided.",
		)
		return
	}

	if ruleSetFull != actionSetFull {
		resp.Diagnostics.AddError(
			"Incomplete rule action configuration",
			"rule_id and action must be provided together.",
		)
	}
}

// Create implements resource's Create method.
func (r *wafAIRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating WAF AI Rules Resource")

	var data wafAIRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()
	client := r.Client.GetAPPSEC()

	version, err := getModifiableConfigVersion(ctx, int(configID), wafAIRulesResourceName, client)
	if err != nil {
		resp.Diagnostics.AddError("fetching modifiable config version", err.Error())
		return
	}

	if !data.AIRuleStatus.IsNull() {
		_, err = client.UpdateAIRulesStatus(ctx, appsec.UpdateAIRulesStatusRequest{
			ConfigID: configID, Version: version, PolicyID: policyID,
			Body: appsec.UpdateAIRulesStatusRequestBody{AIRuleStatus: data.AIRuleStatus.ValueString()},
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating AI Rules status", err.Error())
			return
		}
		data.RuleID = types.Int64Null()
		data.RuleVersionID = types.Int64Null()
		data.Action = types.StringNull()
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	ruleID := data.RuleID.ValueInt64()
	ruleVersionID, err := r.resolveRuleVersionID(ctx, client, configID, version, policyID, ruleID, data.RuleVersionID)
	if err != nil {
		resp.Diagnostics.AddError("resolving rule version", err.Error())
		return
	}

	_, err = client.UpdateAIRuleAction(ctx, appsec.UpdateAIRuleActionRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
		RuleID: ruleID, RuleVersionID: ruleVersionID,
		Body: appsec.UpdateAIRuleActionRequestBody{Action: data.Action.ValueString()},
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("updating action for AI rule %d", ruleID), err.Error())
		return
	}

	data.AIRuleStatus = types.StringNull()
	data.RuleVersionID = types.Int64Value(ruleVersionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read implements resource's Read method.
func (r *wafAIRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading WAF AI Rules Resource")

	var data wafAIRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()
	client := r.Client.GetAPPSEC()

	version, err := getLatestConfigVersion(ctx, int(configID), client)
	if err != nil {
		resp.Diagnostics.AddError("fetching latest config version", err.Error())
		return
	}

	if !data.AIRuleStatus.IsNull() {
		statusResp, err := client.GetAIRulesStatus(ctx, appsec.GetAIRulesStatusRequest{
			ConfigID: configID, Version: version, PolicyID: policyID,
		})
		if err != nil {
			var apiErr *appsec.Error
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddWarning(
					"AI rules are no longer enrolled",
					fmt.Sprintf("Security policy %q is no longer enrolled in AI rules; removing the resource from state.", policyID),
				)
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError("reading AI rules status", err.Error())
			return
		}
		if statusResp.AIRuleStatus == aiRuleStatusNotEnrolled {
			resp.Diagnostics.AddWarning(
				"AI rules are no longer enrolled",
				fmt.Sprintf("Security policy %q is no longer enrolled in AI rules; removing the resource from state.", policyID),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		data.AIRuleStatus = types.StringValue(statusResp.AIRuleStatus)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	actionResp, err := client.GetAIRuleAction(ctx, appsec.GetAIRuleActionRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
		RuleID: data.RuleID.ValueInt64(), RuleVersionID: data.RuleVersionID.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("reading AI rule action", err.Error())
		return
	}
	data.Action = types.StringValue(actionResp.Action)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource's Update method.
func (r *wafAIRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating WAF AI Rules Resource")

	var plan wafAIRulesResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := plan.ConfigID.ValueInt64()
	policyID := plan.SecurityPolicyID.ValueString()
	client := r.Client.GetAPPSEC()

	version, err := getModifiableConfigVersion(ctx, int(configID), wafAIRulesResourceName, client)
	if err != nil {
		resp.Diagnostics.AddError("fetching modifiable config version", err.Error())
		return
	}

	if !plan.AIRuleStatus.IsNull() {
		_, err = client.UpdateAIRulesStatus(ctx, appsec.UpdateAIRulesStatusRequest{
			ConfigID: configID, Version: version, PolicyID: policyID,
			Body: appsec.UpdateAIRulesStatusRequestBody{AIRuleStatus: plan.AIRuleStatus.ValueString()},
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating AI Rules status", err.Error())
			return
		}
		plan.RuleID = types.Int64Null()
		plan.RuleVersionID = types.Int64Null()
		plan.Action = types.StringNull()
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	ruleID := plan.RuleID.ValueInt64()
	_, err = client.UpdateAIRuleAction(ctx, appsec.UpdateAIRuleActionRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
		RuleID: ruleID, RuleVersionID: plan.RuleVersionID.ValueInt64(),
		Body: appsec.UpdateAIRuleActionRequestBody{Action: plan.Action.ValueString()},
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("updating action for AI rule %d", ruleID), err.Error())
		return
	}
	plan.AIRuleStatus = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource's Delete method.
func (r *wafAIRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting WAF AI Rules Resource")

	var data wafAIRulesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.AIRuleStatus.IsNull() {
		// Action mode: reset the rule action to "none" before removing from state.
		configID := data.ConfigID.ValueInt64()
		policyID := data.SecurityPolicyID.ValueString()
		client := r.Client.GetAPPSEC()
		version, err := getModifiableConfigVersion(ctx, int(configID), wafAIRulesResourceName, client)
		if err != nil {
			resp.Diagnostics.AddError("fetching modifiable config version", err.Error())
			return
		}
		_, err = client.UpdateAIRuleAction(ctx, appsec.UpdateAIRuleActionRequest{
			ConfigID: configID, Version: version, PolicyID: policyID,
			RuleID: data.RuleID.ValueInt64(), RuleVersionID: data.RuleVersionID.ValueInt64(),
			Body: appsec.UpdateAIRuleActionRequestBody{Action: "none"},
		})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("resetting action for AI rule %d", data.RuleID.ValueInt64()), err.Error())
		}
		return
	}

	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()
	client := r.Client.GetAPPSEC()

	version, err := getModifiableConfigVersion(ctx, int(configID), wafAIRulesResourceName, client)
	if err != nil {
		resp.Diagnostics.AddError("fetching modifiable config version", err.Error())
		return
	}

	_, err = client.UpdateAIRulesStatus(ctx, appsec.UpdateAIRulesStatusRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
		Body: appsec.UpdateAIRulesStatusRequestBody{AIRuleStatus: aiRuleStatusDisabled},
	})
	if err != nil {
		resp.Diagnostics.AddError("disabling AI rules", err.Error())
	}
}

// ImportState implements resource's ImportState method.
// Status mode import ID: CONFIG_ID:POLICY_ID
// Action mode import ID: CONFIG_ID:POLICY_ID:RULE_ID
func (r *wafAIRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing WAF AI Rules resource")

	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 && len(parts) != 3 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ID %q incorrectly formatted", req.ID),
			"Expected CONFIG_ID:POLICY_ID (status mode) or CONFIG_ID:POLICY_ID:RULE_ID (action mode).",
		)
		return
	}

	configID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid config id %q", parts[0]), "")
		return
	}
	policyID := parts[1]
	if policyID == "" {
		resp.Diagnostics.AddError("invalid security policy id", "policy ID cannot be empty")
		return
	}

	client := r.Client.GetAPPSEC()

	version, err := getLatestConfigVersion(ctx, int(configID), client)
	if err != nil {
		resp.Diagnostics.AddError("fetching latest config version", err.Error())
		return
	}

	if len(parts) == 2 {
		statusResp, err := client.GetAIRulesStatus(ctx, appsec.GetAIRulesStatusRequest{
			ConfigID: configID, Version: version, PolicyID: policyID,
		})
		if err != nil {
			var apiErr *appsec.Error
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
				resp.Diagnostics.AddError(
					"policy not enrolled in AI rules",
					fmt.Sprintf("security policy %q is not enrolled in AI rules and cannot be imported.", policyID),
				)
				return
			}
			resp.Diagnostics.AddError("reading AI rules status", err.Error())
			return
		}
		if statusResp.AIRuleStatus == aiRuleStatusNotEnrolled {
			resp.Diagnostics.AddError(
				"policy not enrolled in AI rules",
				fmt.Sprintf("security policy %q is not enrolled in AI rules and cannot be imported.", policyID),
			)
			return
		}
		resp.Diagnostics.Append(resp.State.Set(ctx, wafAIRulesResourceModel{
			ConfigID:         types.Int64Value(configID),
			SecurityPolicyID: types.StringValue(policyID),
			AIRuleStatus:     types.StringValue(statusResp.AIRuleStatus),
			RuleID:           types.Int64Null(),
			RuleVersionID:    types.Int64Null(),
			Action:           types.StringNull(),
		})...)
		return
	}

	ruleID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid rule id %q", parts[2]), "")
		return
	}

	rulesResp, err := client.ListAIRules(ctx, appsec.ListAIRulesRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
	})
	if err != nil {
		resp.Diagnostics.AddError("reading AI rules", err.Error())
		return
	}

	var ruleVersionID int64
	for _, rule := range rulesResp.AIRules {
		if rule.RuleID == ruleID {
			ruleVersionID = rule.RuleVersion
			break
		}
	}
	if ruleVersionID == 0 {
		resp.Diagnostics.AddError(fmt.Sprintf("AI rule %d not found in policy %q", ruleID, policyID), "")
		return
	}

	actionResp, err := client.GetAIRuleAction(ctx, appsec.GetAIRuleActionRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
		RuleID: ruleID, RuleVersionID: ruleVersionID,
	})
	if err != nil {
		resp.Diagnostics.AddError("reading AI rule action", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, wafAIRulesResourceModel{
		ConfigID:         types.Int64Value(configID),
		SecurityPolicyID: types.StringValue(policyID),
		AIRuleStatus:     types.StringNull(),
		RuleID:           types.Int64Value(ruleID),
		RuleVersionID:    types.Int64Value(ruleVersionID),
		Action:           types.StringValue(actionResp.Action),
	})...)
}

func (r *wafAIRulesResource) resolveRuleVersionID(ctx context.Context, client appsec.APPSEC, configID int64, version int, policyID string, ruleID int64, stateVersionID types.Int64) (int64, error) {
	if !stateVersionID.IsNull() && !stateVersionID.IsUnknown() {
		return stateVersionID.ValueInt64(), nil
	}
	rulesResp, err := client.ListAIRules(ctx, appsec.ListAIRulesRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
	})
	if err != nil {
		return 0, fmt.Errorf("looking up rule version: %w", err)
	}
	for _, rule := range rulesResp.AIRules {
		if rule.RuleID == ruleID {
			return rule.RuleVersion, nil
		}
	}
	return 0, fmt.Errorf("AI rule %d not found in policy %q", ruleID, policyID)
}
