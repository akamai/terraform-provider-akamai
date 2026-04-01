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
	urlProtectionPoliciesActionsDataSource struct {
		meta meta.Meta
	}

	urlProtectionPoliciesActionsDataSourceModel struct {
		ConfigID                     types.Int64                      `tfsdk:"config_id"`
		SecurityPolicyID             types.String                     `tfsdk:"security_policy_id"`
		URLProtectionPoliciesActions []urlProtectionPolicyActionModel `tfsdk:"url_protection_policies_actions"`
	}

	urlProtectionPolicyActionModel struct {
		URLProtectionPolicyID  types.Int64  `tfsdk:"url_protection_policy_id"`
		MaxRateThresholdAction types.String `tfsdk:"max_rate_threshold_action"`
		LoadSheddingAction     types.String `tfsdk:"load_shedding_action"`
	}
)

var (
	_ datasource.DataSource              = &urlProtectionPoliciesActionsDataSource{}
	_ datasource.DataSourceWithConfigure = &urlProtectionPoliciesActionsDataSource{}
)

// NewURLProtectionPoliciesActionsDataSource returns a new URL protection policies actions data source.
func NewURLProtectionPoliciesActionsDataSource() datasource.DataSource {
	return &urlProtectionPoliciesActionsDataSource{}
}

// Metadata configures data source's meta information.
func (d *urlProtectionPoliciesActionsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_url_protection_policies_actions"
}

// Schema is used to define data source's terraform schema.
func (d *urlProtectionPoliciesActionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL protection policies actions data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
			},
			"url_protection_policies_actions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of URL protection policies actions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"url_protection_policy_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the URL protection policy.",
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
				},
			},
		},
	}
}

// Configure configures data source at the beginning of the lifecycle.
func (d *urlProtectionPoliciesActionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *urlProtectionPoliciesActionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "URLProtectionPoliciesActionsDataSource Read")

	var data urlProtectionPoliciesActionsDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	configID := data.ConfigID.ValueInt64()
	policyID := data.SecurityPolicyID.ValueString()

	version, err := getLatestConfigVersion(ctx, int(configID), d.meta)
	if err != nil {
		resp.Diagnostics.AddError("Read URL Protection Policies Actions failed", err.Error())
		return
	}

	getURLProtectionPoliciesActions := appsec.ListURLProtectionPoliciesActionsRequest{
		ConfigID:      configID,
		ConfigVersion: int64(version),
		PolicyID:      policyID,
	}

	urlProtectionPoliciesActions, err := client.ListURLProtectionPoliciesActions(ctx, getURLProtectionPoliciesActions)
	if err != nil {
		resp.Diagnostics.AddError("Read URL Protection Policies Actions failed", err.Error())
		return
	}

	actions := make([]urlProtectionPolicyActionModel, 0)
	if urlProtectionPoliciesActions != nil && len(urlProtectionPoliciesActions.URLProtectionPoliciesActions) > 0 {
		for _, action := range urlProtectionPoliciesActions.URLProtectionPoliciesActions {
			actions = append(actions, urlProtectionPolicyActionModel{
				URLProtectionPolicyID:  types.Int64Value(action.URLProtectionPolicyID),
				MaxRateThresholdAction: types.StringValue(action.MaxRateThresholdAction),
				LoadSheddingAction:     types.StringValue(action.LoadSheddingAction),
			})
		}
	}
	data.URLProtectionPoliciesActions = actions

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
