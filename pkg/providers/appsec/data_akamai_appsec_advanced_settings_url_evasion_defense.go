package appsec

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type urlEvasionDefenseDataSource struct {
	meta.DataSource
}

var (
	_ datasource.DataSource              = &urlEvasionDefenseDataSource{}
	_ datasource.DataSourceWithConfigure = &urlEvasionDefenseDataSource{}
)

// NewURLEvasionDefenseDataSource returns a new URL Evasion Defense data source.
func NewURLEvasionDefenseDataSource() datasource.DataSource {
	return &urlEvasionDefenseDataSource{}
}

// Metadata configures data source's meta information.
func (d *urlEvasionDefenseDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_advanced_settings_url_evasion_defense"
}

// Schema is used to define data source's Terraform schema.
func (d *urlEvasionDefenseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "URL Evasion Defense advanced settings data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Security configuration ID.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Sets whether the feature is `enabled` or `disabled`.",
			},
			"bypass_lists": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Client list identifiers for trusted clients who are exempt from URL evasion mitigation rules.",
			},
			"rules": urlEvasionDefenseDataSourceRulesSchema(),
		},
	}
}

func urlEvasionDefenseDataSourceRulesSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "URL Evasion Defense rules.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"rule_id": schema.Int64Attribute{
					Computed:    true,
					Description: "Uniquely identifies the URL evasion mitigation rule.",
				},
				"action": schema.StringAttribute{
					Computed:    true,
					Description: "The URL evasion mitigation rule action.",
				},
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "The URL evasion mitigation rule name.",
				},
				"description": schema.StringAttribute{
					Computed:    true,
					Description: "The URL evasion mitigation rule description.",
				},
				"condition_operator": schema.StringAttribute{
					Computed:    true,
					Description: "Sets how the rule evaluates conditions. Use `OR` to match any condition, or `AND` to match on all conditions. When the specified conditions are met, the rule does not trigger.",
				},
				"conditions": urlEvasionDefenseDataSourceConditionsSchema(),
			},
		},
	}
}

func urlEvasionDefenseDataSourceConditionsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:    true,
		Description: "The list of match conditions.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Computed:    true,
					Description: "The condition type to match on.",
				},
				"extensions": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "The file extensions that trigger the condition. This only applies to the `extensionMatch` condition `type`.",
				},
				"filenames": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "The filenames that trigger the condition. This only applies to the `filenameMatch` condition `type`.",
				},
				"hosts": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "The hostnames that trigger the condition. This only applies to the `hostMatch` condition `type`.",
				},
				"ips": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "The IPs that trigger the condition. This only applies to the `ipMatch` condition `type`.",
				},
				"methods": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "The HTTP request methods that trigger the condition. The possible values are `GET`, `POST`, `HEAD`, `PUT`, `DELETE`, `OPTIONS`, `TRACE`, `CONNECT` and `PATCH`. This only applies to the `requestMethodMatch` condition `type`.",
				},
				"paths": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "The paths that trigger the condition. This only applies to the  `pathMatch` condition `type`.",
				},
				"client_lists": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
					Description: "The clientLists that trigger the condition. This only applies to the `clientListMatch` condition `type`.",
				},
				"header": schema.StringAttribute{
					Computed:    true,
					Description: "The HTTP header that triggers the condition. This only applies to the `requestHeaderMatch` condition `type`.",
				},
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "The query parameter name that triggers the condition. This only applies to the `uriQueryMatch` condition `type`.",
				},
				"value": schema.StringAttribute{
					Computed:    true,
					Description: "The query parameter value if the condition `type` is `uriQueryMatch` and header value if the condition `type` is `requestHeaderMatch`. This only applies when the condition `type` is `uriQueryMatch` or `requestHeaderMatch`.",
				},
				"positive_match": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether the condition should trigger on a match (`true`) or a lack of match (`false`).",
				},
				"use_headers": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether the condition should include `X-Forwarded-For` (XFF) header. This applies to the `ipMatch` and `clientListMatch` condition `type`.",
				},
				"name_case_sensitive": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether to consider the case-sensitivity of the provided query parameter `name`. This only applies to the `uriQueryMatch` condition `type`.",
				},
				"value_case_sensitive": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether to consider the case-sensitivity of the provided `value`. This only applies to the `requestHeaderMatch` and `uriQueryMatch` condition `type`.",
				},
				"value_wildcard": schema.BoolAttribute{
					Computed:    true,
					Description: "Whether the provided parameter `value` is a wildcard. This only applies to the `requestHeaderMatch and `uriQueryMatch` condition `type`.",
				},
			},
		},
	}
}

func (d *urlEvasionDefenseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "URLEvasionDefenseDataSource Read")

	var data urlEvasionDefenseResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := data.ConfigID.ValueInt64()
	client := d.Client.GetAPPSEC()
	version, err := getLatestConfigVersion(ctx, int(configID), client)
	if err != nil {
		resp.Diagnostics.AddError("invalid config version", err.Error())
		return
	}

	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{
		ConfigID: configID,
		Version:  int64(version),
	}

	result, err := client.GetAdvancedSettingsURLEvasionDefense(ctx, getRequest)
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetAdvancedSettingsURLEvasionDefense'", err.Error())
		return
	}

	resp.Diagnostics.Append(data.populateModelFromResponse(ctx, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
