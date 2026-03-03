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
	urlProtectionRulesDataSource struct {
		meta meta.Meta
	}

	// urlProtectionRulesDataSourceModel maps the url protection rules data source schema data
	urlProtectionRulesDataSourceModel struct {
		ConfigID           types.Int64                        `tfsdk:"config_id"`
		URLProtectionRules []urlProtectionRuleDataSourceModel `tfsdk:"url_protection_rules"`
	}
)

var (
	_ datasource.DataSource              = &urlProtectionRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &urlProtectionRulesDataSource{}
)

// NewURLProtectionRulesDataSource returns a new url protection rule data source
func NewURLProtectionRulesDataSource() datasource.DataSource { return &urlProtectionRulesDataSource{} }

// Metadata configures data source's meta information
func (d *urlProtectionRulesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_url_protection_rules"
}

// Schema is used to define data source's terraform schema
func (d *urlProtectionRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL Protection Rules data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"url_protection_rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of url protection rules",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"config_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the security configuration",
						},
						"url_protection_rule_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the URL protection rule",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the URL protection rule",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the URL protection rule",
						},
						"bypass_conditions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of bypass conditions for the URL protection rule",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "Type of condition (e.g., RequestHeaderCondition, NetworkListCondition)",
									},
									"names": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
										Description: "List of header names for RequestHeaderCondition",
									},
									"name_wildcard": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether to use wildcard matching for header names",
									},
									"values": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
										Description: "List of values for the condition",
									},
									"value_case_sensitive": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether the value matching is case sensitive",
									},
									"value_wildcard": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether to use wildcard matching for values",
									},
								},
							},
						},
						"max_rate_threshold": schema.Int64Attribute{
							Computed:    true,
							Description: "Maximum rate threshold for the URL protection rule",
						},
						"api_definitions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of API definitions associated with the URL protection rule",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"api_definition_id": schema.Int64Attribute{
										Computed:    true,
										Description: "Unique identifier of the API definition",
									},
									"defined_resources": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether defined resources are included",
									},
									"resource_ids": schema.ListAttribute{
										ElementType: types.Int64Type,
										Computed:    true,
										Description: "List of resource IDs",
									},
									"undefined_resources": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether undefined resources are included",
									},
								},
							},
						},
						"hostname_paths": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of hostname and path configurations",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"hostname": schema.StringAttribute{
										Computed:    true,
										Description: "Hostname for the URL protection rule",
									},
									"paths": schema.ListAttribute{
										ElementType: types.StringType,
										Computed:    true,
										Description: "List of paths associated with the hostname",
									},
								},
							},
						},
						"intelligent_load_shedding": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Intelligent load shedding configuration",
							Attributes: map[string]schema.Attribute{
								"hits_per_sec": schema.Int64Attribute{
									Computed:    true,
									Description: "Number of hits per second threshold",
								},
								"categories": schema.ListAttribute{
									ElementType: types.StringType,
									Computed:    true,
									Description: "List of categories for intelligent load shedding",
								},
								"custom_criteria": schema.ListNestedAttribute{
									Computed:    true,
									Description: "Custom criteria for intelligent load shedding",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Computed:    true,
												Description: "Type of custom criteria (e.g., CLIENT_LIST)",
											},
											"list_ids": schema.ListAttribute{
												ElementType: types.StringType,
												Computed:    true,
												Description: "List of client list IDs",
											},
											"positive_match": schema.BoolAttribute{
												Computed:    true,
												Description: "Whether this is a positive match condition",
											},
										},
									},
								},
							},
						},
						"used": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether you're currently using the URL protection policy",
						},
						"create_date": schema.StringAttribute{
							Computed:    true,
							Description: "Date when the URL protection rule was created",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who created the URL protection rule",
						},
						"update_date": schema.StringAttribute{
							Computed:    true,
							Description: "Date when the URL protection rule was last updated",
						},
						"updated_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who last updated the URL protection rule",
						},
					},
				},
			},
		},
	}
}

// Configure configures data source at the beginning of the lifecycle
func (d *urlProtectionRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *urlProtectionRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "URLProtectionRulesDataSource Read")

	var data urlProtectionRulesDataSourceModel

	var urlProtectionRulesModel urlProtectionRulesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := inst.Client(d.meta)
	configID := data.ConfigID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), d.meta)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version", err.Error())
		return
	}

	listURLProtectionRuleRequest := appsec.ListURLProtectionRulesRequest{
		ConfigID:      configID,
		ConfigVersion: int64(version),
	}

	urlProtectionRules, err := client.ListURLProtectionRules(ctx, listURLProtectionRuleRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'ListURLProtectionRules'", err.Error())
		return
	}

	if urlProtectionRules == nil {
		// Reset state when URL Protection Rules are not found
		urlProtectionRulesModel.ConfigID = data.ConfigID
		urlProtectionRulesModel.URLProtectionRules = []urlProtectionRuleDataSourceModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &urlProtectionRulesModel)...)
		return
	}

	if len(urlProtectionRules.URLProtectionRules) > 0 {
		urlProtectionRulesModel.URLProtectionRules = make([]urlProtectionRuleDataSourceModel, 0, len(urlProtectionRules.URLProtectionRules))
		for _, urlProtectionRule := range urlProtectionRules.URLProtectionRules {
			urlProtectionRuleModel, diags := createURLProtectionRuleModel(ctx, &urlProtectionRule)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			urlProtectionRulesModel.URLProtectionRules = append(urlProtectionRulesModel.URLProtectionRules, urlProtectionRuleModel)

		}
	}

	urlProtectionRulesModel.ConfigID = data.ConfigID

	resp.Diagnostics.Append(resp.State.Set(ctx, &urlProtectionRulesModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
