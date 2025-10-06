package appsec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	customRulesUsageDataSource struct {
		meta meta.Meta
	}

	// customRulesUsageDataSourceModel describes the data source data model for CustomRulesUsageDataSource.
	customRulesUsageDataSourceModel struct {
		ConfigID   types.Int64             `tfsdk:"config_id"`
		RuleIDs    types.Set               `tfsdk:"rule_ids"`
		Rules      []customRulesUsageModel `tfsdk:"rules"`
		JSON       types.String            `tfsdk:"json"`
		OutputText types.String            `tfsdk:"output_text"`
	}

	customRulesUsageModel struct {
		RuleID   types.Int64   `tfsdk:"rule_id"`
		Policies []policyModel `tfsdk:"policies"`
	}

	policyModel struct {
		PolicyID   types.String `tfsdk:"policy_id"`
		PolicyName types.String `tfsdk:"policy_name"`
	}

	customRuleUsageItem struct {
		RuleID     int64  `json:"ruleId"`
		PolicyID   string `json:"policyId"`
		PolicyName string `json:"policyName"`
	}
)

var (
	_ datasource.DataSource              = &customRulesUsageDataSource{}
	_ datasource.DataSourceWithConfigure = &customRulesUsageDataSource{}
)

// NewCustomRulesUsageDataSource returns a new custom rules usage data source
func NewCustomRulesUsageDataSource() datasource.DataSource { return &customRulesUsageDataSource{} }

// Metadata configures data source's meta information
func (d *customRulesUsageDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_custom_rules_usage"
}

// Schema is used to define data source's terraform schema
func (d *customRulesUsageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Custom rules usage data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "A security configuration ID.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rule_ids": schema.SetAttribute{
				Required:    true,
				ElementType: types.Int64Type,
				Description: "A custom rule IDs.",
				Validators: []validator.Set{
					setvalidator.IsRequired(),
				},
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A custom rules usage.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policies": schema.ListNestedAttribute{
							Computed:    true,
							Description: "A set of security policies in which a custom rule is used.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"policy_id": schema.StringAttribute{
										Computed:    true,
										Description: "The security policy ID.",
									},
									"policy_name": schema.StringAttribute{
										Computed:    true,
										Description: "The security policy name.",
									},
								},
							},
						},
						"rule_id": schema.Int64Attribute{
							Computed:    true,
							Description: "The ID of the custom rule.",
						},
					},
				},
			},
			"json": schema.StringAttribute{
				Computed:    true,
				Description: "JSON-formatted information about the custom rules usage.",
			},
			"output_text": schema.StringAttribute{
				Computed:    true,
				Description: "Tabular representation of the custom rules usage.",
			},
		},
	}
}

func (d *customRulesUsageDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring Custom Rules data source")

	if request.ProviderData == nil {
		return
	}

	m, ok := request.ProviderData.(meta.Meta)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	d.meta = m
}

func (d *customRulesUsageDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Custom Rules data source")

	var data customRulesUsageDataSourceModel
	response.Diagnostics.Append(request.Config.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	configID := data.ConfigID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), d.meta)
	if err != nil {
		response.Diagnostics.AddError("get latest config version error", err.Error())
		return
	}

	getCustomRuleUsage := appsec.GetCustomRulesUsageRequest{
		ConfigID: configID,
		Version:  version,
		RequestBody: appsec.RuleIDs{
			IDs: convertRuleIDsSetToSlice(data.RuleIDs),
		},
	}

	customRulesUsage, err := client.GetCustomRulesUsage(ctx, getCustomRuleUsage)
	if err != nil {
		response.Diagnostics.AddError("calling 'GetCustomRuleUsage'", err.Error())
		return
	}

	rules := customRulesUsage.Rules
	usage := createCustomRulesUsageModel(rules)

	jsonBody, err := json.MarshalIndent(customRulesUsage, "", "  ")
	if err != nil {
		response.Diagnostics.AddError("Error marshaling JSON", err.Error())
		return
	}

	outputText, diags := createOutputText(rules)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	data.Rules = usage
	data.JSON = types.StringValue(string(jsonBody))
	data.OutputText = types.StringValue(outputText)
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func createOutputText(rules []appsec.CustomRuleUsage) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	ots := OutputTemplates{}
	InitTemplates(ots)

	usageItems := make([]customRuleUsageItem, 0, len(rules))
	for _, rule := range rules {
		usageItem := customRuleUsageItem{
			RuleID: rule.RuleID,
		}
		for _, policy := range rule.Policies {
			usageItem.PolicyID = policy.PolicyID
			usageItem.PolicyName = policy.PolicyName
			usageItems = append(usageItems, usageItem)
		}
	}

	outputText, err := RenderTemplates(ots, "customRulesUsage", usageItems)
	if err != nil {
		diags.AddError("Error rendering output text", err.Error())
	}
	return outputText, diags
}

func createCustomRulesUsageModel(rules []appsec.CustomRuleUsage) []customRulesUsageModel {
	usage := make([]customRulesUsageModel, 0, len(rules))
	for _, rule := range rules {
		usageModel := customRulesUsageModel{
			RuleID: types.Int64Value(rule.RuleID),
		}
		for _, policy := range rule.Policies {
			policyModel := policyModel{
				PolicyID:   types.StringValue(policy.PolicyID),
				PolicyName: types.StringValue(policy.PolicyName),
			}
			usageModel.Policies = append(usageModel.Policies, policyModel)
		}
		usage = append(usage, usageModel)
	}
	return usage
}

func convertRuleIDsSetToSlice(set types.Set) []int64 {
	elems := set.Elements()
	result := make([]int64, len(elems))

	for i, e := range elems {
		intVal := e.(types.Int64)
		result[i] = intVal.ValueInt64()
	}

	return result
}
