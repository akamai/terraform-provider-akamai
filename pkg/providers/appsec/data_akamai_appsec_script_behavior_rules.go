package appsec

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &scriptBehaviorRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptBehaviorRulesDataSource{}
)

type scriptBehaviorRulesDataSource struct {
	meta.DataSource
}

type scriptBehaviorRulesDataSourceModel struct {
	ConfigID         types.Int64                   `tfsdk:"config_id"`
	SecurityPolicyID types.String                  `tfsdk:"security_policy_id"`
	Rules            []scriptBehaviorRuleItemModel `tfsdk:"rules"`
}

type scriptBehaviorRuleItemModel struct {
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

// NewScriptBehaviorRulesDataSource returns a new script behavior rules data source.
func NewScriptBehaviorRulesDataSource() datasource.DataSource {
	return &scriptBehaviorRulesDataSource{}
}

func (d *scriptBehaviorRulesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_script_behavior_rules"
}

func (d *scriptBehaviorRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves all script behavior rules for a security policy.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of script behavior rules for the security policy.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"rule_id":          schema.Int64Attribute{Computed: true, Description: "Unique identifier of the script behavior rule."},
						"rule_name":        schema.StringAttribute{Computed: true, Description: "Name of the script behavior rule."},
						"rule_description": schema.StringAttribute{Computed: true, Description: "Description of the script behavior rule."},
						"activity_type":    schema.StringAttribute{Computed: true, Description: "Type of activity the rule controls."},
						"action":           schema.StringAttribute{Computed: true, Description: "Action taken when the rule conditions are matched."},
						"page_condition":   schema.StringAttribute{Computed: true, Description: "Condition type for matching page URLs."},
						"page_values": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of page URL patterns to match.",
						},
						"script_type":      schema.StringAttribute{Computed: true, Description: "Type of script to match."},
						"script_condition": schema.StringAttribute{Computed: true, Description: "Condition type for matching scripts."},
						"script_values": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of script patterns to match.",
						},
						"activity_values_condition": schema.StringAttribute{Computed: true, Description: "Condition type for matching activity type values."},
						"activity_values": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "List of activity type values to match.",
						},
					},
				},
			},
		},
	}
}

func (d *scriptBehaviorRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "ScriptBehaviorRulesDataSource Read")

	var data scriptBehaviorRulesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client.GetAPPSEC()
	configID := int(data.ConfigID.ValueInt64())

	version, err := getLatestConfigVersion(ctx, configID, client)
	if err != nil {
		resp.Diagnostics.AddError("getting latest config version", err.Error())
		return
	}

	rules, err := client.ListClientSideProtectionRules(ctx, appsec.ListClientSideProtectionRulesRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: data.SecurityPolicyID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("listing script behavior rules", err.Error())
		return
	}

	ruleModels := make([]scriptBehaviorRuleItemModel, 0, len(*rules))
	for _, rule := range *rules {
		item, diags := scriptBehaviorRuleItemFromBody(ctx, rule)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		ruleModels = append(ruleModels, item)
	}
	data.Rules = ruleModels

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func scriptBehaviorRuleItemFromBody(ctx context.Context, rule appsec.ClientSideProtectionRuleBody) (scriptBehaviorRuleItemModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ruleID int64
	if rule.RuleID != nil {
		ruleID = *rule.RuleID
	}

	pageValues, d := types.ListValueFrom(ctx, types.StringType, rule.Conditions.Page.Values)
	diags.Append(d...)

	scriptValues, d := types.ListValueFrom(ctx, types.StringType, rule.Conditions.Script.Values)
	diags.Append(d...)

	activityValues, d := types.ListValueFrom(ctx, types.StringType, rule.Conditions.ActivityTypeValues.Values)
	diags.Append(d...)

	return scriptBehaviorRuleItemModel{
		RuleID:                  types.Int64Value(ruleID),
		RuleName:                types.StringValue(rule.RuleName),
		RuleDescription:         stringOrNull(rule.Description),
		ActivityType:            types.StringValue(rule.ActivityType),
		Action:                  types.StringValue(rule.Action),
		PageCondition:           types.StringValue(rule.Conditions.Page.Condition),
		PageValues:              pageValues,
		ScriptType:              stringOrNull(rule.Conditions.Script.ScriptType),
		ScriptCondition:         types.StringValue(rule.Conditions.Script.Condition),
		ScriptValues:            scriptValues,
		ActivityValuesCondition: types.StringValue(rule.Conditions.ActivityTypeValues.Condition),
		ActivityValues:          activityValues,
	}, diags
}
