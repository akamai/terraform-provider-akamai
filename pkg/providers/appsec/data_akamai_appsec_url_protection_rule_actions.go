package appsec

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	urlProtectionRuleActionsDataSource struct {
		meta meta.Meta
	}

	urlProtectionRuleActionsDataSourceModel struct {
		ConfigID               types.Int64  `tfsdk:"config_id"`
		SecurityPolicyID       types.String `tfsdk:"security_policy_id"`
		URLProtectionRuleID    types.Int64  `tfsdk:"url_protection_rule_id"`
		MaxRateThresholdAction types.String `tfsdk:"max_rate_threshold_action"`
		LoadSheddingAction     types.String `tfsdk:"load_shedding_action"`
	}
)

var (
	_ datasource.DataSource              = &urlProtectionRuleActionsDataSource{}
	_ datasource.DataSourceWithConfigure = &urlProtectionRuleActionsDataSource{}
)

// NewURLProtectionRuleActionsDataSource returns a new URL protection rule actions data source.
func NewURLProtectionRuleActionsDataSource() datasource.DataSource {
	return &urlProtectionRuleActionsDataSource{}
}

// Metadata configures data source's meta information.
func (d *urlProtectionRuleActionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_url_protection_rule_actions"
}

// Schema is used to define data source's terraform schema.
func (d *urlProtectionRuleActionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL protection rule actions data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
			},
			"url_protection_rule_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the URL protection rule for which to return action information.",
			},
			"max_rate_threshold_action": schema.StringAttribute{
				Computed:    true,
				Description: "Action to take when the max rate threshold is exceeded.",
			},
			"load_shedding_action": schema.StringAttribute{
				Computed:    true,
				Description: "Action to take for load shedding.",
			},
		},
	}
}

// Configure configures data source at the beginning of the lifecycle.
func (d *urlProtectionRuleActionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read is called when the provider must read data source values in order to update state.
func (d *urlProtectionRuleActionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "URLProtectionRuleActionsDataSource Read")

	var data urlProtectionRuleActionsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()
	urlProtectionID := data.URLProtectionRuleID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), d.meta)
	if err != nil {
		resp.Diagnostics.AddError("Read URL Protection Rule Actions failed", err.Error())
		return
	}

	getURLProtectionRuleActions := appsec.GetURLProtectionRuleActionsRequest{
		ConfigID:            configID,
		ConfigVersion:       int64(version),
		PolicyID:            policyID,
		URLProtectionRuleID: urlProtectionID,
	}

	urlProtectionRuleActions, err := client.GetURLProtectionRuleActions(ctx, getURLProtectionRuleActions)
	if err != nil {
		resp.Diagnostics.AddError("Read URL Protection Rule Actions failed", err.Error())
		return
	}

	data.MaxRateThresholdAction = types.StringValue(urlProtectionRuleActions.MaxRateThresholdAction)
	data.LoadSheddingAction = types.StringValue(urlProtectionRuleActions.LoadSheddingAction)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
