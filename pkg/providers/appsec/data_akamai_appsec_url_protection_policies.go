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
	urlProtectionPoliciesDataSource struct {
		meta.DataSource
	}

	// urlProtectionPoliciesDataSourceModel maps the url protection policies data source schema data
	urlProtectionPoliciesDataSourceModel struct {
		ConfigID              types.Int64                          `tfsdk:"config_id"`
		URLProtectionPolicies []urlProtectionPolicyDataSourceModel `tfsdk:"url_protection_policies"`
	}
)

var (
	_ datasource.DataSource              = &urlProtectionPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &urlProtectionPoliciesDataSource{}
)

// NewURLProtectionPoliciesDataSource returns a new url protection policies data source
func NewURLProtectionPoliciesDataSource() datasource.DataSource {
	return &urlProtectionPoliciesDataSource{}
}

// Metadata configures data source's meta information
func (d *urlProtectionPoliciesDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = "akamai_appsec_url_protection_policies"
}

// Schema is used to define data source's terraform schema
func (d *urlProtectionPoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL Protection Policies data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"url_protection_policies": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of url protection policies",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"config_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the security configuration",
						},
						"url_protection_policy_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Unique identifier of the URL protection policy",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the URL protection policy",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the URL protection policy",
						},
						"bypass_conditions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of bypass conditions for the URL protection policy",
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
							Description: "Maximum rate threshold for the URL protection policy",
						},
						"api_definitions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of API definitions associated with the URL protection policy",
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
										Description: "Hostname for the URL protection policy",
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
							Description: "Date when the URL protection policy was created",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who created the URL protection policy",
						},
						"update_date": schema.StringAttribute{
							Computed:    true,
							Description: "Date when the URL protection policy was last updated",
						},
						"updated_by": schema.StringAttribute{
							Computed:    true,
							Description: "User who last updated the URL protection policy",
						},
					},
				},
			},
		},
	}
}

func (d *urlProtectionPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "URLProtectionPoliciesDataSource Read")

	var data urlProtectionPoliciesDataSourceModel

	var urlProtectionPoliciesModel urlProtectionPoliciesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client.GetAPPSEC()
	configID := data.ConfigID.ValueInt64()

	version, err := getLatestConfigVersion(ctx, int(configID), d.Client.GetAPPSEC())
	if err != nil {
		resp.Diagnostics.AddError("invalid config version", err.Error())
		return
	}

	listURLProtectionPolicyRequest := appsec.ListURLProtectionPoliciesRequest{
		ConfigID:      configID,
		ConfigVersion: int64(version),
	}

	urlProtectionPolicies, err := client.ListURLProtectionPolicies(ctx, listURLProtectionPolicyRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'ListURLProtectionPolicies'", err.Error())
		return
	}

	if urlProtectionPolicies == nil {
		// Reset state when URL Protection Policies are not found
		urlProtectionPoliciesModel.ConfigID = data.ConfigID
		urlProtectionPoliciesModel.URLProtectionPolicies = []urlProtectionPolicyDataSourceModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &urlProtectionPoliciesModel)...)
		return
	}

	if len(urlProtectionPolicies.URLProtectionPolicies) > 0 {
		urlProtectionPoliciesModel.URLProtectionPolicies = make([]urlProtectionPolicyDataSourceModel, 0, len(urlProtectionPolicies.URLProtectionPolicies))
		for _, urlProtectionPolicy := range urlProtectionPolicies.URLProtectionPolicies {
			urlProtectionPolicyModel, diags := createURLProtectionPolicyModel(ctx, &urlProtectionPolicy)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			urlProtectionPoliciesModel.URLProtectionPolicies = append(urlProtectionPoliciesModel.URLProtectionPolicies, urlProtectionPolicyModel)

		}
	}

	urlProtectionPoliciesModel.ConfigID = data.ConfigID

	resp.Diagnostics.Append(resp.State.Set(ctx, &urlProtectionPoliciesModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
