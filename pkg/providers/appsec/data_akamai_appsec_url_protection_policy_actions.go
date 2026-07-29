package appsec

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type (
	urlProtectionPolicyActionsDataSource struct {
		meta.DataSource
	}

	urlProtectionPolicyActionsDataSourceModel struct {
		ConfigID               types.Int64  `tfsdk:"config_id"`
		SecurityPolicyID       types.String `tfsdk:"security_policy_id"`
		URLProtectionPolicyID  types.Int64  `tfsdk:"url_protection_policy_id"`
		MaxRateThresholdAction types.String `tfsdk:"max_rate_threshold_action"`
		LoadSheddingAction     types.String `tfsdk:"load_shedding_action"`
	}
)

var (
	_ datasource.DataSource              = &urlProtectionPolicyActionsDataSource{}
	_ datasource.DataSourceWithConfigure = &urlProtectionPolicyActionsDataSource{}
)

// NewURLProtectionPolicyActionsDataSource returns a new URL protection policy actions data source.
func NewURLProtectionPolicyActionsDataSource() datasource.DataSource {
	return &urlProtectionPolicyActionsDataSource{}
}

// Metadata configures data source's meta information.
func (d *urlProtectionPolicyActionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_url_protection_policy_actions"
}

// Schema is used to define data source's terraform schema.
func (d *urlProtectionPolicyActionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL protection policy actions data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
			},
			"url_protection_policy_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the URL protection policy for which to return action information.",
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

// Read is called when the provider must read data source values in order to update state.
func (d *urlProtectionPolicyActionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "URLProtectionPolicyActionsDataSource Read")

	var data urlProtectionPolicyActionsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := d.Client.GetAPPSEC()
	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()
	urlProtectionID := data.URLProtectionPolicyID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), d.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("Read URL Protection Policy Actions failed", err.Error())
		return
	}

	getURLProtectionPolicyActions := appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              configID,
		ConfigVersion:         int64(version),
		PolicyID:              policyID,
		URLProtectionPolicyID: urlProtectionID,
	}

	urlProtectionPolicyActions, err := client.GetURLProtectionPolicyActions(ctx, getURLProtectionPolicyActions)
	if err != nil {
		resp.Diagnostics.AddError("Read URL Protection Policy Actions failed", err.Error())
		return
	}

	data.MaxRateThresholdAction = types.StringValue(urlProtectionPolicyActions.MaxRateThresholdAction)
	data.LoadSheddingAction = types.StringValue(urlProtectionPolicyActions.LoadSheddingAction)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
