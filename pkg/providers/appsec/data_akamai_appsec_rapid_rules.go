package appsec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	rapidRulesDataSource struct {
		meta meta.Meta
	}

	// rapidRulesDataSourceModel describes the data source data model for RapidRulesDataSource.
	rapidRulesDataSourceModel struct {
		ID                   types.String     `tfsdk:"id"`
		ConfigID             types.Int64      `tfsdk:"config_id"`
		PolicyID             types.String     `tfsdk:"security_policy_id"`
		RuleID               types.Int64      `tfsdk:"rule_id"`
		Enabled              types.Bool       `tfsdk:"enabled"`
		DefaultAction        types.String     `tfsdk:"default_action"`
		RapidRules           []rapidRuleModel `tfsdk:"rapid_rules"`
		IncludeExpiryDetails types.Bool       `tfsdk:"include_expiry_details"`
		OutputText           types.String     `tfsdk:"output_text"`
	}

	rapidRuleModel struct {
		ID                   types.Int64  `tfsdk:"id"`
		Action               types.String `tfsdk:"action"`
		Lock                 types.Bool   `tfsdk:"lock"`
		Name                 types.String `tfsdk:"name"`
		AttackGroup          types.String `tfsdk:"attack_group"`
		AttackGroupException types.String `tfsdk:"attack_group_exception"`
		ConditionException   types.String `tfsdk:"condition_exception"`
		Expired              types.Bool   `tfsdk:"expired"`
		ExpireInDays         types.Int64  `tfsdk:"expire_in_days"`
	}
)

var (
	_ datasource.DataSource              = &rapidRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &rapidRulesDataSource{}
)

// NewRapidRulesDataSource returns a new rapid rules data source
func NewRapidRulesDataSource() datasource.DataSource { return &rapidRulesDataSource{} }

// Metadata configures data source's meta information
func (d *rapidRulesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_rapid_rules"
}

// Schema is used to define data source's terraform schema
func (d *rapidRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rapid rules data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"rule_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Unique identifier of a specific rapid rule for which to retrieve information",
			},
			"include_expiry_details": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to return expiry details, including `expired` and `expire_in_days` attributes, for each rapid rule. Defaults to `false` if not set.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "The rapid rules status.",
			},
			"default_action": schema.StringAttribute{
				Computed:    true,
				Description: "The default action for new rapid rules.",
			},
			"rapid_rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of rapid rules detailed information include action, lock, attack group, exceptions",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "The unique identifier of rapid rule.",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "The rapid rule action.",
						},
						"lock": schema.BoolAttribute{
							Computed:    true,
							Description: "The the rapid rule action lock.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The rapid rule name.",
						},
						"attack_group": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier of attack group, rapid rule belongs to.",
						},
						"attack_group_exception": schema.StringAttribute{
							Computed:    true,
							Description: "The attack group exception.",
						},
						"condition_exception": schema.StringAttribute{
							Computed:    true,
							Description: "The rapid rule exception.",
						},
						"expired": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the rule has already expired.",
						},
						"expire_in_days": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of days remaining before the rule expires. This field is present only if the rule has not yet expired.",
						},
					},
				},
			},
			"output_text": schema.StringAttribute{
				Computed:    true,
				Description: "Text representation",
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the data source",
				Computed:            true,
			},
		},
	}
}

// Configure configures data source at the beginning of the lifecycle
func (d *rapidRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Data Source Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.",
					req.ProviderData))
		}
	}()
	d.meta = meta.Must(req.ProviderData)
}

func (d *rapidRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "RapidRulesDataSource Read")

	var data rapidRulesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	configID := data.ConfigID.ValueInt64()
	ruleID := getRuleID(data.RuleID)
	policyID := data.PolicyID.ValueString()
	includeExpiry := !data.IncludeExpiryDetails.IsNull() && data.IncludeExpiryDetails.ValueBool()

	getRulesRequest := appsec.GetRapidRulesRequest{
		ConfigID:             configID,
		PolicyID:             policyID,
		RuleID:               ruleID,
		IncludeExpiryDetails: includeExpiry,
	}

	version, err := getLatestConfigVersion(ctx, int(configID), d.meta)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version", err.Error())
		return
	}
	getRulesRequest.Version = version

	getRapidRulesStatusRequest := appsec.GetRapidRulesStatusRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	status, err := client.GetRapidRulesStatus(ctx, getRapidRulesStatusRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetRapidRulesStatus'", err.Error())
		return
	}

	if !status.Enabled {
		data.populateState("No default action. Rapid rules is turned off.", status.Enabled, "Rapid rules is turned off.", nil)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	rules, err := client.GetRapidRules(ctx, getRulesRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetRapidRules'", err.Error())
		return
	}

	getRapidRulesDefaultActionRequest := appsec.GetRapidRulesDefaultActionRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	defaultActionResponse, err := client.GetRapidRulesDefaultAction(ctx, getRapidRulesDefaultActionRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'getRapidRulesDefaultAction'", err.Error())
		return
	}

	getAttackGroupsRequest := appsec.GetAttackGroupsRequest{
		ConfigID: int(configID),
		Version:  version,
		PolicyID: policyID,
	}

	attackGroups, err := client.GetAttackGroups(ctx, getAttackGroupsRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetAttackGroups'", err.Error())
		return
	}

	rapidRules := convertGetRapidRulesResponseToRapidRules(rules, attackGroups, includeExpiry)

	outputText, diags := generateOutputText(rapidRules)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	output, err := generateOutputRapidRules(*rapidRules)
	if err != nil {
		resp.Diagnostics.AddError("generating output_rapid_rules error", err.Error())
		return
	}

	data.populateState(defaultActionResponse.Action, status.Enabled, outputText, *output)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func generateOutputText(rapidRules *[]appsec.RapidRuleDetails) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	ots := OutputTemplates{}
	InitTemplates(ots)
	templateName := "RapidRulesWithConditionExceptionDS"
	outputText, err := RenderTemplates(ots, templateName, rapidRules)
	if err != nil {
		diags.AddError(err.Error(), "")
		return "", diags
	}
	return outputText, diags
}

func (m *rapidRulesDataSourceModel) populateState(defaultAction string, enabled bool, outputText string, output []rapidRuleModel) {
	m.ID = types.StringValue(fmt.Sprintf("%s:%s", m.ConfigID.String(), m.PolicyID.String()))
	m.Enabled = types.BoolValue(enabled)
	m.DefaultAction = types.StringValue(defaultAction)
	m.OutputText = types.StringValue(outputText)
	m.RapidRules = output
}

func toJSONString(value any) (string, error) {
	jsonString, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

func generateOutputRapidRules(input []appsec.RapidRuleDetails) (*[]rapidRuleModel, error) {
	output := make([]rapidRuleModel, 0, len(input))
	for _, rule := range input {
		attackGroupException, err := toJSONString(rule.AttackGroupException)
		if err != nil {
			return &output, err
		}

		conditionException, err := toJSONString(rule.ConditionException)
		if err != nil {
			return &output, err
		}

		outputRule := rapidRuleModel{
			ID:                   types.Int64Value(rule.ID),
			Action:               types.StringValue(rule.Action),
			Lock:                 types.BoolValue(rule.Lock),
			Name:                 types.StringValue(rule.Name),
			AttackGroup:          types.StringValue(rule.AttackGroup),
			AttackGroupException: types.StringValue(attackGroupException),
			ConditionException:   types.StringValue(conditionException),
			Expired:              types.BoolPointerValue(rule.Expired),
			ExpireInDays:         types.Int64PointerValue(rule.ExpireInDays),
		}
		output = append(output, outputRule)
	}
	return &output, nil
}

func convertGetRapidRulesResponseToRapidRules(input *appsec.GetRapidRulesResponse, attackGroups *appsec.GetAttackGroupsResponse, includeExpiry bool) *[]appsec.RapidRuleDetails {

	output := make([]appsec.RapidRuleDetails, 0, len(input.Rules))

	for _, rule := range input.Rules {
		attackGroup := getAttackGroup(rule.RiskScoreGroups)
		outputRule := appsec.RapidRuleDetails{
			ID:                   rule.ID,
			Action:               rule.Action,
			Lock:                 rule.Lock,
			Name:                 rule.Name,
			AttackGroup:          attackGroup,
			AttackGroupException: getAttackGroupException(attackGroup, attackGroups),
			ConditionException:   rule.ConditionException,
		}

		if includeExpiry {
			if rule.Expired != nil {
				expired := types.BoolPointerValue(rule.Expired)
				if !expired.IsNull() && expired.ValueBool() {
					outputRule.Expired = expired.ValueBoolPointer()
				}
			}

			if rule.ExpireInDays != nil {
				expireInDays := types.Int64PointerValue(rule.ExpireInDays)
				if !expireInDays.IsNull() {
					outputRule.ExpireInDays = expireInDays.ValueInt64Pointer()
				}
			}
		}

		output = append(output, outputRule)
	}
	return &output
}

func getAttackGroup(riskScoreGroups []string) string {
	if len(riskScoreGroups) > 0 {
		return riskScoreGroups[0]
	}
	return ""
}

func getAttackGroupException(attackGroup string, attackGroups *appsec.GetAttackGroupsResponse) *appsec.AttackGroupConditionException {
	for _, groupInfo := range attackGroups.AttackGroups {
		if groupInfo.Group == attackGroup {
			return groupInfo.ConditionException
		}
	}
	return nil
}

func getRuleID(value types.Int64) *int64 {
	if value.IsNull() {
		return nil
	}
	ruleID := value.ValueInt64()
	return &ruleID
}
