package appsec

import (
	"context"
	"fmt"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	wafRulesetDataSource struct {
		meta meta.Meta
	}

	// wafRulesetDataSourceModel describes the data source data model for WAFRulesetDataSource.
	wafRulesetDataSourceModel struct {
		ConfigID         types.Int64        `tfsdk:"config_id"`
		SecurityPolicyID types.String       `tfsdk:"security_policy_id"`
		AttackGroups     []attackGroupModel `tfsdk:"attack_groups"`
		Rules            []wafRuleModel     `tfsdk:"rules"`
	}

	attackGroupModel struct {
		AttackGroup        types.String `tfsdk:"attack_group"`
		AttackGroupAction  types.String `tfsdk:"attack_group_action"`
		ConditionException types.String `tfsdk:"condition_exception"`
	}

	wafRuleModel struct {
		RuleID             types.Int64  `tfsdk:"rule_id"`
		RuleAction         types.String `tfsdk:"rule_action"`
		ConditionException types.String `tfsdk:"condition_exception"`
	}
)

var (
	_ datasource.DataSource              = &wafRulesetDataSource{}
	_ datasource.DataSourceWithConfigure = &wafRulesetDataSource{}
)

// NewWAFRulesetDataSource returns a new WAF ruleset data source
func NewWAFRulesetDataSource() datasource.DataSource {
	return &wafRulesetDataSource{}
}

// Metadata configures data source's meta information
func (d *wafRulesetDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_waf_ruleset"
}

// Schema is used to define data source's terraform schema
func (d *wafRulesetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "WAF ruleset data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"attack_groups": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of attack group configurations",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attack_group": schema.StringAttribute{
							Computed:    true,
							Description: "Unique name of the attack group",
						},
						"attack_group_action": schema.StringAttribute{
							Computed:    true,
							Description: "Action taken when the attack group is triggered (alert, deny, deny_custom_{custom_deny_id}, none)",
						},
						"condition_exception": schema.StringAttribute{
							Computed:    true,
							Description: "Conditions and exceptions associated with the attack group",
						},
					},
				},
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of rule objects including action and condition exceptions",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier for a rule",
						},
						"rule_action": schema.StringAttribute{
							Computed:    true,
							Description: "Action taken when the rule is triggered (alert, deny, deny_custom_{custom_deny_id}, none)",
						},
						"condition_exception": schema.StringAttribute{
							Computed:    true,
							Description: "Conditions and exceptions associated with the rule",
						},
					},
				},
			},
		},
	}
}

// Configure configures data source at the beginning of the lifecycle
func (d *wafRulesetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read is called when the provider must read data source values in order to update state
func (d *wafRulesetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "WAFRulesetDataSource Read")

	var data wafRulesetDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()

	// Get the latest version for this configuration
	version, err := getLatestConfigVersion(ctx, int(configID), d.meta)
	if err != nil {
		resp.Diagnostics.AddError("retrieving config version", err.Error())
		return
	}

	// Call GetWAFCompositeRuleset API to retrieve the complete WAF ruleset
	getWAFRulesetReq := appsec.GetWAFCompositeRulesetRequest{
		ConfigID: configID,
		Version:  int64(version),
		PolicyID: policyID,
	}

	wafRuleset, err := client.GetWAFCompositeRuleset(ctx, getWAFRulesetReq)
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetWAFCompositeRuleset'", err.Error())
		return
	}

	// Convert attack groups to model
	attackGroups, diags := convertAttackGroups(wafRuleset.AttackGroups)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Convert rules to model
	rules, diags := convertRules(wafRuleset.Rules)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Populate state
	data.AttackGroups = attackGroups
	data.Rules = rules

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// convertAttackGroups converts edge-grid SDK attack groups to framework model
// Attack groups are sorted alphabetically by group name to ensure consistent ordering
func convertAttackGroups(attackGroups []appsec.WAFCompositeAttackGroup) ([]attackGroupModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	output := make([]attackGroupModel, 0, len(attackGroups))

	for _, attackGroup := range attackGroups {

		conditionException, err := toJSONString(attackGroup.ConditionException)
		if err != nil {
			diags.AddError("marshaling attack group condition exception", err.Error())
			return output, diags
		}

		outputAttackGroup := attackGroupModel{
			AttackGroup:        types.StringValue(attackGroup.Group),
			AttackGroupAction:  types.StringValue(attackGroup.Action),
			ConditionException: types.StringValue(conditionException),
		}
		output = append(output, outputAttackGroup)
	}

	// Sort attack groups alphabetically by group name to prevent diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].AttackGroup.ValueString() < output[j].AttackGroup.ValueString()
	})

	return output, diags
}

// convertRules converts edge-grid SDK rules to framework model
// Rules are sorted by rule ID to ensure consistent ordering
func convertRules(rules []appsec.WAFCompositeRule) ([]wafRuleModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	output := make([]wafRuleModel, 0, len(rules))

	for _, rule := range rules {
		conditionException, err := toJSONString(rule.ConditionException)
		if err != nil {
			diags.AddError("marshaling rule condition exception", err.Error())
			return output, diags
		}

		outputRule := wafRuleModel{
			RuleID:             types.Int64Value(rule.RuleID),
			RuleAction:         types.StringValue(rule.Action),
			ConditionException: types.StringValue(conditionException),
		}
		output = append(output, outputRule)
	}

	// Sort rules by rule ID to prevent diffs
	sort.Slice(output, func(i, j int) bool {
		return output[i].RuleID.ValueInt64() < output[j].RuleID.ValueInt64()
	})

	return output, diags
}
