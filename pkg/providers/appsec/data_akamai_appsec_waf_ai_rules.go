package appsec

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	wafAIRulesDataSource struct {
		meta.DataSource
	}

	wafAIRulesDataSourceModel struct {
		ID               types.String `tfsdk:"id"`
		ConfigID         types.Int64  `tfsdk:"config_id"`
		SecurityPolicyID types.String `tfsdk:"security_policy_id"`
		AIRuleStatus     types.String `tfsdk:"ai_rule_status"`
		AIRules          types.List   `tfsdk:"ai_rules"`
		OutputText       types.String `tfsdk:"output_text"`
	}

	wafAIRuleModel struct {
		RuleID             types.Int64  `tfsdk:"rule_id"`
		RuleVersion        types.Int64  `tfsdk:"rule_version"`
		Title              types.String `tfsdk:"title"`
		RiskScoreGroup     types.String `tfsdk:"risk_score_group"`
		RuleDescription    types.String `tfsdk:"rule_description"`
		Action             types.String `tfsdk:"action"`
		ConditionException types.String `tfsdk:"condition_exception"`
	}
)

var (
	_ datasource.DataSource              = &wafAIRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &wafAIRulesDataSource{}

	wafAIRuleObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"rule_id":             types.Int64Type,
			"rule_version":        types.Int64Type,
			"title":               types.StringType,
			"risk_score_group":    types.StringType,
			"rule_description":    types.StringType,
			"action":              types.StringType,
			"condition_exception": types.StringType,
		},
	}
)

// NewWAFAIRulesDataSource returns a new WAF AI rules data source.
func NewWAFAIRulesDataSource() datasource.DataSource { return &wafAIRulesDataSource{} }

// Metadata sets the data source type name.
func (d *wafAIRulesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_waf_ai_rules"
}

// Schema defines the data source Terraform schema.
func (d *wafAIRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "WAF AI rules data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
			},
			"ai_rule_status": schema.StringAttribute{
				Computed:    true,
				Description: "Whether AI rules are enabled, disabled, or not enrolled for the policy.",
			},
			"ai_rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of AI rules for the security policy.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the AI rule.",
						},
						"rule_version": schema.Int64Attribute{
							Computed:    true,
							Description: "Version of the AI rule.",
						},
						"title": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the AI rule.",
						},
						"risk_score_group": schema.StringAttribute{
							Computed:    true,
							Description: "Risk score group the AI rule belongs to.",
						},
						"rule_description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of what the AI rule detects.",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "Action taken when the AI rule is triggered.",
						},
						"condition_exception": schema.StringAttribute{
							Computed:    true,
							Description: "JSON-encoded list of condition exceptions for the AI rule.",
						},
					},
				},
			},
			"output_text": schema.StringAttribute{
				Computed:    true,
				Description: "Text representation of the AI rules.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the data source.",
			},
		},
	}
}

// Read fetches AI rules from the API and populates state.
func (d *wafAIRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "WAFAIRulesDataSource Read")

	var data wafAIRulesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client.GetAPPSEC()
	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()

	version, err := getLatestConfigVersion(ctx, int(configID), client)
	if err != nil {
		resp.Diagnostics.AddError("fetching latest config version", err.Error())
		return
	}

	// ENROLLED (ENABLED/DISABLED) policies: GET /ai-rules returns 404; use GET /ai-rules/status.
	// NOT_ENROLLED policies: GET /ai-rules returns 200 with status + rules; GET /ai-rules/status returns 404.
	aiRuleStatus := aiRuleStatusNotEnrolled
	var rules []appsec.PolicyAIRule

	result, err := client.ListAIRules(ctx, appsec.ListAIRulesRequest{
		ConfigID: configID, Version: version, PolicyID: policyID,
	})
	var apiErr *appsec.Error
	switch {
	case err == nil:
		aiRuleStatus = result.AIRuleStatus
		rules = result.AIRules
	case errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound:
		// ENROLLED policy: fall back to /ai-rules/status for the enrollment status.
		statusResp, statusErr := client.GetAIRulesStatus(ctx, appsec.GetAIRulesStatusRequest{
			ConfigID: configID, Version: version, PolicyID: policyID,
		})
		var statusAPIErr *appsec.Error
		if statusErr == nil {
			aiRuleStatus = statusResp.AIRuleStatus
		} else if !errors.As(statusErr, &statusAPIErr) || statusAPIErr.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError("calling 'GetAIRulesStatus'", statusErr.Error())
			return
		}
		// Both endpoints 404 → aiRuleStatus stays NOT_ENROLLED, rules stays empty.
	default:
		resp.Diagnostics.AddError("calling 'ListAIRules'", err.Error())
		return
	}

	outputText, err := generateAIRulesOutputText(rules)
	if err != nil {
		resp.Diagnostics.AddError("generating output_text", err.Error())
		return
	}

	aiRuleModels, err := toWAFAIRuleModels(ctx, rules)
	if err != nil {
		resp.Diagnostics.AddError("mapping AI rules", err.Error())
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d:%s", configID, policyID))
	data.AIRuleStatus = types.StringValue(aiRuleStatus)
	data.AIRules = aiRuleModels
	data.OutputText = types.StringValue(outputText)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func toWAFAIRuleModels(ctx context.Context, rules []appsec.PolicyAIRule) (types.List, error) {
	out := make([]wafAIRuleModel, 0, len(rules))
	for _, r := range rules {
		exceptions := r.ConditionExceptions
		if exceptions == nil {
			exceptions = []appsec.AIRuleConditionException{}
		}
		conditionException, err := toJSONString(exceptions)
		if err != nil {
			return types.ListNull(wafAIRuleObjectType), fmt.Errorf("marshaling condition exceptions for rule %d: %w", r.RuleID, err)
		}
		out = append(out, wafAIRuleModel{
			RuleID:             types.Int64Value(r.RuleID),
			RuleVersion:        types.Int64Value(r.RuleVersion),
			Title:              types.StringValue(r.Title),
			RiskScoreGroup:     types.StringValue(r.RiskScoreGroup),
			RuleDescription:    types.StringValue(r.RuleDescription),
			Action:             types.StringValue(r.Action),
			ConditionException: types.StringValue(conditionException),
		})
	}
	result, diags := types.ListValueFrom(ctx, wafAIRuleObjectType, out)
	if diags.HasError() {
		return types.ListNull(wafAIRuleObjectType), fmt.Errorf("converting AI rules to list: %s", diags)
	}
	return result, nil
}

func generateAIRulesOutputText(rules []appsec.PolicyAIRule) (string, error) {
	ots := OutputTemplates{}
	InitTemplates(ots)
	return RenderTemplates(ots, "AIRulesDS", rules)
}
