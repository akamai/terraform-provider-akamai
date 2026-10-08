package appsec

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &scriptBehaviorRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptBehaviorRuleDataSource{}
)

type scriptBehaviorRuleDataSource struct {
	meta.DataSource
}

type scriptBehaviorRuleDataSourceModel struct {
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

// NewScriptBehaviorRuleDataSource returns a new script behavior rule data source.
func NewScriptBehaviorRuleDataSource() datasource.DataSource { return &scriptBehaviorRuleDataSource{} }

func (d *scriptBehaviorRuleDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_script_behavior_rule"
}

func (d *scriptBehaviorRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a single script behavior rule for a security policy.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
			},
			"rule_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the script behavior rule.",
			},
			"rule_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the script behavior rule.",
			},
			"rule_description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the script behavior rule.",
			},
			"activity_type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of activity the rule controls.",
			},
			"action": schema.StringAttribute{
				Computed:    true,
				Description: "Action taken when the rule conditions are matched.",
			},
			"page_condition": schema.StringAttribute{
				Computed:    true,
				Description: "Condition type for matching page URLs.",
			},
			"page_values": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of page URL patterns to match.",
			},
			"script_type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of script to match.",
			},
			"script_condition": schema.StringAttribute{
				Computed:    true,
				Description: "Condition type for matching scripts.",
			},
			"script_values": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of script patterns to match.",
			},
			"activity_values_condition": schema.StringAttribute{
				Computed:    true,
				Description: "Condition type for matching activity type values.",
			},
			"activity_values": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of activity type values to match.",
			},
		},
	}
}

func (d *scriptBehaviorRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "ScriptBehaviorRuleDataSource Read")

	var data scriptBehaviorRuleDataSourceModel
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

	ruleID := data.RuleID.ValueInt64()
	var found *appsec.ClientSideProtectionRuleBody
	for i := range *rules {
		if (*rules)[i].RuleID != nil && *(*rules)[i].RuleID == ruleID {
			found = &(*rules)[i]
			break
		}
	}
	if found == nil {
		resp.Diagnostics.AddError("reading script behavior rule", fmt.Sprintf("rule %d not found in policy %q", ruleID, data.SecurityPolicyID.ValueString()))
		return
	}

	data.RuleName = types.StringValue(found.RuleName)
	data.RuleDescription = stringOrNull(found.Description)
	data.ActivityType = types.StringValue(found.ActivityType)
	data.Action = types.StringValue(found.Action)
	data.PageCondition = types.StringValue(found.Conditions.Page.Condition)
	data.ScriptType = stringOrNull(found.Conditions.Script.ScriptType)
	data.ScriptCondition = types.StringValue(found.Conditions.Script.Condition)
	data.ActivityValuesCondition = types.StringValue(found.Conditions.ActivityTypeValues.Condition)

	pageValues, diags := types.ListValueFrom(ctx, types.StringType, found.Conditions.Page.Values)
	resp.Diagnostics.Append(diags...)
	data.PageValues = pageValues

	scriptValues, diags := types.ListValueFrom(ctx, types.StringType, found.Conditions.Script.Values)
	resp.Diagnostics.Append(diags...)
	data.ScriptValues = scriptValues

	activityValues, diags := types.ListValueFrom(ctx, types.StringType, found.Conditions.ActivityTypeValues.Values)
	resp.Diagnostics.Append(diags...)
	data.ActivityValues = activityValues

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
