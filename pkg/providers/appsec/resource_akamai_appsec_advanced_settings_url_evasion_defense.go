package appsec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &urlEvasionDefenseResource{}
	_ resource.ResourceWithConfigure   = &urlEvasionDefenseResource{}
	_ resource.ResourceWithImportState = &urlEvasionDefenseResource{}
)

const (
	readError        = "could not read URL Evasion Defense from API"
	readVersionError = "could not read Configuration Version from API"
	updateFailed     = "URL Evasion Defense update failed"
	importIDError    = "ID '%s' incorrectly formatted: should be 'CONFIG_ID'"
)

type urlEvasionDefenseResource struct {
	meta.Resource
}

type urlEvasionDefenseResourceModel struct {
	ConfigID    types.Int64  `tfsdk:"config_id"`
	Status      types.String `tfsdk:"status"`
	BypassLists types.Set    `tfsdk:"bypass_lists"`
	Rules       types.List   `tfsdk:"rules"`
}

type urlEvasionDefenseRuleModel struct {
	RuleID            types.Int64  `tfsdk:"rule_id"`
	Action            types.String `tfsdk:"action"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	ConditionOperator types.String `tfsdk:"condition_operator"`
	Conditions        types.List   `tfsdk:"conditions"`
}

type urlEvasionDefenseConditionModel struct {
	Type               types.String `tfsdk:"type"`
	Extensions         types.Set    `tfsdk:"extensions"`
	Filenames          types.Set    `tfsdk:"filenames"`
	Hosts              types.Set    `tfsdk:"hosts"`
	IPs                types.Set    `tfsdk:"ips"`
	Methods            types.Set    `tfsdk:"methods"`
	Paths              types.Set    `tfsdk:"paths"`
	ClientLists        types.Set    `tfsdk:"client_lists"`
	Header             types.String `tfsdk:"header"`
	Name               types.String `tfsdk:"name"`
	Value              types.String `tfsdk:"value"`
	PositiveMatch      types.Bool   `tfsdk:"positive_match"`
	NameCaseSensitive  types.Bool   `tfsdk:"name_case_sensitive"`
	ValueCaseSensitive types.Bool   `tfsdk:"value_case_sensitive"`
	ValueWildcard      types.Bool   `tfsdk:"value_wildcard"`
	UseHeaders         types.Bool   `tfsdk:"use_headers"`
}

// NewURLEvasionDefenseResource returns new Url Evasion Defense resource
func NewURLEvasionDefenseResource() resource.Resource {
	return &urlEvasionDefenseResource{}
}

// Metadata implements resource.Resource
func (r *urlEvasionDefenseResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_advanced_settings_url_evasion_defense"
}

// Schema implements resource.Resource
func (r *urlEvasionDefenseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL Evasion Defense advanced settings resource.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Security configuration ID.",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"status": schema.StringAttribute{
				Required:    true,
				Description: "Sets whether the feature is `enabled` or `disabled`.",
				Validators: []validator.String{
					stringvalidator.OneOf(string(appsec.AdvancedSettingsURLEvasionDefenseStatusEnabled), string(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled)),
				},
			},
			"bypass_lists": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Client list identifiers for trusted clients who are exempt from URL evasion mitigation rules.",
				Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
			},
			"rules": urlEvasionDefenseRulesSchema(),
		},
	}
}

func urlEvasionDefenseRulesSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Description: "URL Evasion Defense rules.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"rule_id": schema.Int64Attribute{
					Required:    true,
					Description: "Uniquely identifies the URL evasion mitigation rule.",
				},
				"action": schema.StringAttribute{
					Required:    true,
					Description: "The URL evasion mitigation rule action.",
					Validators:  []validator.String{validators.NotEmptyString()},
				},
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "The URL evasion mitigation rule name.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"description": schema.StringAttribute{
					Computed:    true,
					Description: "The URL evasion mitigation rule description.",
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"condition_operator": schema.StringAttribute{
					Optional:    true,
					Description: "Sets how the rule evaluates conditions. Use `OR` to match any condition, or `AND` to match on all conditions. When the specified conditions are met, the rule does not trigger.",
					Validators: []validator.String{
						stringvalidator.OneOf(string(appsec.AdvancedSettingsURLEvasionDefenseConditionOperatorOr), string(appsec.AdvancedSettingsURLEvasionDefenseConditionOperatorAnd)),
					},
				},
				"conditions": urlEvasionDefenseConditionsSchema(),
			},
		},
	}
}

func urlEvasionDefenseConditionsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Description: "The list of match conditions.",
		Validators:  []validator.List{ValidateCondition()},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:    true,
					Description: "The condition type to match on.",
				},
				"extensions": schema.SetAttribute{
					Optional:    true,
					ElementType: types.StringType,
					Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
					Description: "The file extensions that trigger the condition. This only applies to the `extensionMatch` condition `type`.",
				},
				"filenames": schema.SetAttribute{
					Optional:    true,
					ElementType: types.StringType,
					Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
					Description: "The filenames that trigger the condition. This only applies to the `filenameMatch` condition `type`.",
				},
				"hosts": schema.SetAttribute{
					Optional:    true,
					ElementType: types.StringType,
					Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
					Description: "The hostnames that trigger the condition. This only applies to the `hostMatch` condition `type`.",
				},
				"ips": schema.SetAttribute{
					Optional:    true,
					ElementType: types.StringType,
					Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
					Description: "The IPs that trigger the condition. This only applies to the `ipMatch` condition `type`.",
				},
				"methods": schema.SetAttribute{
					Optional:    true,
					ElementType: types.StringType,
					Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
					Description: "The HTTP request methods that trigger the condition. The possible values are `GET`, `POST`, `HEAD`, `PUT`, `DELETE`, `OPTIONS`, `TRACE`, `CONNECT` and `PATCH`. This only applies to the `requestMethodMatch` condition `type`.",
				},
				"paths": schema.SetAttribute{
					Optional:    true,
					ElementType: types.StringType,
					Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
					Description: "The paths that trigger the condition. This only applies to the  `pathMatch` condition `type`.",
				},
				"client_lists": schema.SetAttribute{
					Optional:    true,
					ElementType: types.StringType,
					Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
					Description: "The clientLists that trigger the condition. This only applies to the `clientListMatch` condition `type`.",
				},
				"header": schema.StringAttribute{
					Optional:    true,
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "The HTTP header that triggers the condition. This only applies to the `requestHeaderMatch` condition `type`.",
				},
				"name": schema.StringAttribute{
					Optional:    true,
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "The query parameter name that triggers the condition. This only applies to the `uriQueryMatch` condition `type`.",
				},
				"value": schema.StringAttribute{
					Optional:    true,
					Validators:  []validator.String{validators.NotEmptyString()},
					Description: "The query parameter value if the condition `type` is `uriQueryMatch` and header value if the condition `type` is `requestHeaderMatch`. This only applies when the condition `type` is `uriQueryMatch` or `requestHeaderMatch`.",
				},
				"positive_match": schema.BoolAttribute{
					Optional:    true,
					Validators:  []validator.Bool{boolvalidator.Equals(true)},
					Description: "Whether the condition should trigger on a match (`true`) or a lack of match (`false`).",
				},
				"use_headers": schema.BoolAttribute{
					Optional:    true,
					Validators:  []validator.Bool{boolvalidator.Equals(true)},
					Description: "Whether the condition should include `X-Forwarded-For` (XFF) header. This applies to the `ipMatch` and `clientListMatch` condition `type`.",
				},
				"name_case_sensitive": schema.BoolAttribute{
					Optional:    true,
					Description: "Whether to consider the case-sensitivity of the provided query parameter `name`. This only applies to the `uriQueryMatch` condition `type`.",
					Validators:  []validator.Bool{boolvalidator.Equals(true)},
				},
				"value_case_sensitive": schema.BoolAttribute{
					Optional:    true,
					Validators:  []validator.Bool{boolvalidator.Equals(true)},
					Description: "Whether to consider the case-sensitivity of the provided `value`. This only applies to the `requestHeaderMatch` and `uriQueryMatch` condition `type`.",
				},
				"value_wildcard": schema.BoolAttribute{
					Optional:    true,
					Validators:  []validator.Bool{boolvalidator.Equals(true)},
					Description: "Whether the provided parameter `value` is a wildcard. This only applies to the `requestHeaderMatch` and `uriQueryMatch` condition `type`.",
				},
			},
		},
	}
}

// Create implements resource.Resource
func (r *urlEvasionDefenseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating URL Evasion Defense resource")

	var plan urlEvasionDefenseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.updateAndRefresh(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *urlEvasionDefenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading URL Evasion Defense resource")

	var state urlEvasionDefenseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := state.ConfigID.ValueInt64()
	client := r.Client.GetAPPSEC()
	version, err := getLatestConfigVersion(ctx, int(configID), client)
	if err != nil {
		resp.Diagnostics.AddError(readVersionError, err.Error())
		return
	}
	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{
		ConfigID: configID,
		Version:  int64(version),
	}
	result, err := client.GetAdvancedSettingsURLEvasionDefense(ctx, getRequest)
	if err != nil {
		resp.Diagnostics.AddError(readError, err.Error())
		return
	}

	diags := state.populateModelFromResponse(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *urlEvasionDefenseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating URL Evasion Defense resource")

	var plan urlEvasionDefenseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags := r.updateAndRefresh(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *urlEvasionDefenseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting URL Evasion Defense resource")

	var state urlEvasionDefenseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Status = types.StringValue(string(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled))

	resp.Diagnostics.Append(r.update(ctx, &state)...)
}

// ImportState implements resource.ResourceWithImportState. The expected import ID format is 'CONFIG_ID'
func (r *urlEvasionDefenseResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing URL Evasion Defense resource")

	configID, err := strconv.ParseInt(request.ID, 10, 64)
	if err != nil {
		response.Diagnostics.AddError(
			fmt.Sprintf(importIDError, request.ID),
			"",
		)
		return
	}

	client := r.Client.GetAPPSEC()
	version, err := getLatestConfigVersion(ctx, int(configID), client)
	if err != nil {
		response.Diagnostics.AddError(readVersionError, err.Error())
		return
	}

	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{
		ConfigID: configID,
		Version:  int64(version),
	}

	result, err := client.GetAdvancedSettingsURLEvasionDefense(ctx, getRequest)
	if err != nil {
		response.Diagnostics.AddError(readError, err.Error())
		return
	}

	state := urlEvasionDefenseResourceModel{
		ConfigID: types.Int64Value(configID),
	}

	response.Diagnostics.Append(state.populateModelFromResponse(ctx, result)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (r *urlEvasionDefenseResource) update(ctx context.Context, data *urlEvasionDefenseResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	configID := data.ConfigID.ValueInt64()
	client := r.Client.GetAPPSEC()
	version, err := getModifiableConfigVersion(ctx, int(configID), "urlEvasionDefense", client)
	if err != nil {
		diags.AddError(readVersionError, err.Error())
		return diags
	}

	getReq := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{
		ConfigID: configID,
		Version:  int64(version),
	}
	response, err := client.GetAdvancedSettingsURLEvasionDefense(ctx, getReq)
	if err != nil {
		diags.AddError(readError, err.Error())
		return diags
	}

	updateRequest, requestDiags := buildURLDefenseUpdateRequest(ctx, response, data, int64(version))
	diags.Append(requestDiags...)
	if diags.HasError() {
		return diags
	}

	_, err = client.UpdateAdvancedSettingsURLEvasionDefense(ctx, *updateRequest)
	if err != nil {
		diags.AddError(updateFailed, err.Error())
	}

	return diags
}

func (r *urlEvasionDefenseResource) updateAndRefresh(ctx context.Context, data *urlEvasionDefenseResourceModel) diag.Diagnostics {
	diags := r.update(ctx, data)
	if diags.HasError() {
		return diags
	}

	configID := data.ConfigID.ValueInt64()
	client := r.Client.GetAPPSEC()
	version, err := getLatestConfigVersion(ctx, int(configID), client)
	if err != nil {
		diags.AddError(readVersionError, err.Error())
		return diags
	}

	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{
		ConfigID: configID,
		Version:  int64(version),
	}
	response, err := client.GetAdvancedSettingsURLEvasionDefense(ctx, getRequest)
	if err != nil {
		diags.AddError(readError, err.Error())
		return diags
	}

	diags.Append(data.populateModelFromResponse(ctx, response)...)
	return diags
}

func buildURLDefenseUpdateRequest(ctx context.Context, response *appsec.GetAdvancedSettingsURLEvasionDefenseResponse, plan *urlEvasionDefenseResourceModel, version int64) (*appsec.UpdateAdvancedSettingsURLEvasionDefenseRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	updateRequest := &appsec.UpdateAdvancedSettingsURLEvasionDefenseRequest{
		ConfigID: plan.ConfigID.ValueInt64(),
		Version:  version,
		Body: appsec.UpdateAdvancedSettingsURLEvasionDefenseRequestBody{
			Status:      appsec.AdvancedSettingsURLEvasionDefenseStatus(plan.Status.ValueString()),
			LockVersion: response.LockVersion,
		},
	}

	if tf.IsKnown(plan.BypassLists) {
		var bypassLists []string
		diags.Append(plan.BypassLists.ElementsAs(ctx, &bypassLists, false)...)
		if diags.HasError() {
			return nil, diags
		}
		updateRequest.Body.BypassLists = bypassLists
	}

	if tf.IsKnown(plan.Rules) {
		var rules []urlEvasionDefenseRuleModel
		diags.Append(plan.Rules.ElementsAs(ctx, &rules, false)...)
		if diags.HasError() {
			return nil, diags
		}

		apiRules := make([]appsec.AdvancedSettingsURLEvasionDefenseRuleRequest, 0, len(rules))
		for _, rule := range rules {
			apiRule := appsec.AdvancedSettingsURLEvasionDefenseRuleRequest{
				RuleID: rule.RuleID.ValueInt64(),
			}
			if tf.IsKnown(rule.Action) {
				apiRule.Action = rule.Action.ValueString()
			}
			if tf.IsKnown(rule.ConditionOperator) {
				apiRule.ConditionOperator = ptr.To(appsec.AdvancedSettingsURLEvasionDefenseConditionOperator(rule.ConditionOperator.ValueString()))
			}

			if tf.IsKnown(rule.Conditions) {
				var conditionModels []urlEvasionDefenseConditionModel
				diags.Append(rule.Conditions.ElementsAs(ctx, &conditionModels, false)...)
				if diags.HasError() {
					return nil, diags
				}

				apiConditions := make([]appsec.AdvancedSettingsURLEvasionDefenseRuleCondition, 0, len(conditionModels))
				for _, c := range conditionModels {
					cond, d := buildConditionFromModel(ctx, c)
					diags.Append(d...)
					if diags.HasError() {
						return nil, diags
					}
					apiConditions = append(apiConditions, cond)
				}
				apiRule.Conditions = apiConditions
			}

			apiRules = append(apiRules, apiRule)
		}

		updateRequest.Body.Rules = apiRules
	}

	return updateRequest, diags
}

func buildConditionFromModel(ctx context.Context, c urlEvasionDefenseConditionModel) (appsec.AdvancedSettingsURLEvasionDefenseRuleCondition, diag.Diagnostics) {
	var diags diag.Diagnostics

	cond := appsec.AdvancedSettingsURLEvasionDefenseRuleCondition{
		Type:          c.Type.ValueString(),
		Header:        c.Header.ValueString(),
		Name:          c.Name.ValueString(),
		NameCase:      c.NameCaseSensitive.ValueBool(),
		PositiveMatch: c.PositiveMatch.ValueBool(),
		Value:         c.Value.ValueString(),
		UseHeaders:    c.UseHeaders.ValueBool(),
	}

	if c.Type.ValueString() == "requestHeaderMatch" {
		cond.ValueCase = c.ValueCaseSensitive.ValueBool()
		cond.ValueWildcard = c.ValueWildcard.ValueBool()
	} else if c.Type.ValueString() == "uriQueryMatch" {
		cond.CaseSensitive = c.ValueCaseSensitive.ValueBool()
		cond.Wildcard = c.ValueWildcard.ValueBool()
	}

	var d diag.Diagnostics
	cond.Extensions, d = getStringSetValues(ctx, c.Extensions)
	diags.Append(d...)
	cond.Filenames, d = getStringSetValues(ctx, c.Filenames)
	diags.Append(d...)
	cond.Hosts, d = getStringSetValues(ctx, c.Hosts)
	diags.Append(d...)
	cond.IPs, d = getStringSetValues(ctx, c.IPs)
	diags.Append(d...)
	cond.Methods, d = getStringSetValues(ctx, c.Methods)
	diags.Append(d...)
	cond.Paths, d = getStringSetValues(ctx, c.Paths)
	diags.Append(d...)
	cond.ClientLists, d = getStringSetValues(ctx, c.ClientLists)
	diags.Append(d...)

	return cond, diags
}

func (m *urlEvasionDefenseResourceModel) populateModelFromResponse(ctx context.Context, out *appsec.GetAdvancedSettingsURLEvasionDefenseResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Status = types.StringValue(out.Status)

	bypassLists, d := types.SetValueFrom(ctx, types.StringType, out.BypassLists)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.BypassLists = bypassLists

	rules := make([]urlEvasionDefenseRuleModel, 0, len(out.Rules))

	for _, rule := range out.Rules {
		if rule.Action == nil {
			continue
		}

		conditions, d := mapRuleConditionsToList(ctx, rule.Conditions)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		ruleState := urlEvasionDefenseRuleModel{
			RuleID:            types.Int64Value(rule.RuleID),
			Action:            stringOrNull(*rule.Action),
			Name:              stringOrNull(rule.Name),
			Description:       stringOrNull(rule.Description),
			ConditionOperator: types.StringPointerValue(rule.ConditionOperator),
			Conditions:        conditions,
		}

		rules = append(rules, ruleState)
	}

	ruleList, d := types.ListValueFrom(ctx, urlEvasionDefenseRuleType(), rules)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Rules = ruleList

	return diags
}

func mapRuleConditionsToList(ctx context.Context, conditions []appsec.AdvancedSettingsURLEvasionDefenseRuleCondition) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(conditions) == 0 {
		return types.ListNull(urlEvasionDefenseConditionType()), diags
	}

	conditionsModel := make([]urlEvasionDefenseConditionModel, 0, len(conditions))
	for _, c := range conditions {
		// Keep switch-based population, but pre-type Set fields to avoid DynamicPseudoType errors.
		conditionModel := urlEvasionDefenseConditionModel{
			Type:          types.StringValue(c.Type),
			PositiveMatch: boolOrNull(c.PositiveMatch),
			Extensions:    types.SetNull(types.StringType),
			Filenames:     types.SetNull(types.StringType),
			Hosts:         types.SetNull(types.StringType),
			IPs:           types.SetNull(types.StringType),
			Methods:       types.SetNull(types.StringType),
			Paths:         types.SetNull(types.StringType),
			ClientLists:   types.SetNull(types.StringType),
		}

		switch c.Type {
		case "extensionMatch":
			extensions, d := stringSetOrNull(ctx, c.Extensions)
			diags.Append(d...)
			conditionModel.Extensions = extensions
		case "filenameMatch":
			filenames, d := stringSetOrNull(ctx, c.Filenames)
			diags.Append(d...)
			conditionModel.Filenames = filenames
		case "hostMatch":
			hosts, d := stringSetOrNull(ctx, c.Hosts)
			diags.Append(d...)
			conditionModel.Hosts = hosts
		case "ipMatch":
			ips, d := stringSetOrNull(ctx, c.IPs)
			diags.Append(d...)
			conditionModel.IPs = ips
			conditionModel.UseHeaders = boolOrNull(c.UseHeaders)
		case "requestMethodMatch":
			methods, d := stringSetOrNull(ctx, c.Methods)
			diags.Append(d...)
			conditionModel.Methods = methods
		case "pathMatch":
			paths, d := stringSetOrNull(ctx, c.Paths)
			diags.Append(d...)
			conditionModel.Paths = paths
		case "clientListMatch":
			clientLists, d := stringSetOrNull(ctx, c.ClientLists)
			diags.Append(d...)
			conditionModel.ClientLists = clientLists
			conditionModel.UseHeaders = boolOrNull(c.UseHeaders)
		case "requestHeaderMatch":
			conditionModel.Header = stringOrNull(c.Header)
			conditionModel.Value = stringOrNull(c.Value)
			conditionModel.ValueCaseSensitive = boolOrNull(c.ValueCase)
			conditionModel.ValueWildcard = boolOrNull(c.ValueWildcard)
		case "uriQueryMatch":
			conditionModel.Name = stringOrNull(c.Name)
			conditionModel.Value = stringOrNull(c.Value)
			conditionModel.NameCaseSensitive = boolOrNull(c.NameCase)
			conditionModel.ValueCaseSensitive = boolOrNull(c.CaseSensitive)
			conditionModel.ValueWildcard = boolOrNull(c.Wildcard)
		}

		if diags.HasError() {
			return types.ListNull(urlEvasionDefenseConditionType()), diags
		}

		conditionsModel = append(conditionsModel, conditionModel)

	}

	listValue, d := types.ListValueFrom(ctx, urlEvasionDefenseConditionType(), conditionsModel)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(urlEvasionDefenseConditionType()), diags
	}

	return listValue, diags
}

func stringOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func boolOrNull(value bool) types.Bool {
	if !value {
		return types.BoolNull()
	}
	return types.BoolValue(true)
}

func stringSetOrNull(ctx context.Context, values []string) (types.Set, diag.Diagnostics) {
	if len(values) == 0 {
		return types.SetNull(types.StringType), nil
	}

	return types.SetValueFrom(ctx, types.StringType, values)
}

func urlEvasionDefenseRuleType() types.ObjectType {
	return urlEvasionDefenseRulesSchema().NestedObject.Type().(types.ObjectType)
}

func urlEvasionDefenseConditionType() types.ObjectType {
	return urlEvasionDefenseConditionsSchema().NestedObject.Type().(types.ObjectType)
}

func getStringSetValues(ctx context.Context, value types.Set) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !tf.IsKnown(value) {
		return nil, diags
	}

	var out []string
	diags.Append(value.ElementsAs(ctx, &out, false)...)
	return out, diags
}
