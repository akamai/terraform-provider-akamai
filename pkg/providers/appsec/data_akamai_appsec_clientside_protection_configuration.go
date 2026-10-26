package appsec

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type clientsideProtectionConfigurationDataSource struct {
	meta.DataSource
}

var (
	_ datasource.DataSource              = &clientsideProtectionConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &clientsideProtectionConfigurationDataSource{}
)

// NewClientsideProtectionConfigurationDataSource returns a new clientside protection configuration data source.
func NewClientsideProtectionConfigurationDataSource() datasource.DataSource {
	return &clientsideProtectionConfigurationDataSource{}
}

func (d *clientsideProtectionConfigurationDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_clientside_protection_configuration"
}

func (d *clientsideProtectionConfigurationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Client-side Protection & Compliance configuration data source.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
			},
			"edge_injection": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Edge injection settings.",
				Attributes: map[string]schema.Attribute{
					"load_script_async": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether to load the protection script asynchronously.",
					},
				},
			},
			"edge_test_parameters": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Edge test parameters.",
				Attributes: map[string]schema.Attribute{
					"disable_injection_key": schema.StringAttribute{
						Computed:    true,
						Description: "Query string key used to disable script injection.",
					},
					"disable_injection_value": schema.StringAttribute{
						Computed:    true,
						Description: "Query string value used to disable script injection.",
					},
					"force_injection_key": schema.StringAttribute{
						Computed:    true,
						Description: "Query string key used to force script injection.",
					},
					"force_injection_value": schema.StringAttribute{
						Computed:    true,
						Description: "Query string value used to force script injection.",
					},
				},
			},
			"injection_criteria": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Injection criteria.",
				Attributes: map[string]schema.Attribute{
					"injection_path_includes": schema.StringAttribute{
						Computed:    true,
						Description: "Paths where script injection is enabled.",
					},
					"injection_path_excludes": schema.StringAttribute{
						Computed:    true,
						Description: "Paths where script injection is disabled.",
					},
					"injection_percent": schema.Int64Attribute{
						Computed:    true,
						Description: "Percentage of requests that receive script injection (0-100).",
					},
					"injection_policy": schema.StringAttribute{
						Computed:    true,
						Description: "Controls where in the page the script is injected.",
					},
					"script_tag_attributes": schema.ListAttribute{
						Computed:    true,
						ElementType: types.StringType,
						Description: "Custom HTML attributes to add to the injected script tag.",
					},
				},
			},
		},
	}
}

type clientsideProtectionConfigurationDataSourceModel struct {
	ConfigID           types.Int64                                  `tfsdk:"config_id"`
	SecurityPolicyID   types.String                                 `tfsdk:"security_policy_id"`
	EdgeInjection      *clientsideProtectionEdgeInjectionModel      `tfsdk:"edge_injection"`
	EdgeTestParameters *clientsideProtectionEdgeTestParametersModel `tfsdk:"edge_test_parameters"`
	InjectionCriteria  *clientsideProtectionInjectionCriteriaModel  `tfsdk:"injection_criteria"`
}

func (d *clientsideProtectionConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading clientside protection configuration data source")

	var data clientsideProtectionConfigurationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := int(data.ConfigID.ValueInt64())
	client := d.Client.GetAPPSEC()

	version, err := getLatestConfigVersion(ctx, configID, client)
	if err != nil {
		resp.Diagnostics.AddError("could not get latest config version", err.Error())
		return
	}

	cpcConfig, err := client.GetClientsideProtectionConfiguration(ctx, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         data.ConfigID.ValueInt64(),
		Version:          int64(version),
		SecurityPolicyID: data.SecurityPolicyID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("calling 'GetClientsideProtectionConfiguration'", err.Error())
		return
	}

	data.EdgeInjection = &clientsideProtectionEdgeInjectionModel{
		LoadScriptAsync: types.BoolValue(cpcConfig.EdgeInjection.LoadScriptAsync),
	}
	data.EdgeTestParameters = &clientsideProtectionEdgeTestParametersModel{
		DisableInjectionKey:   types.StringValue(cpcConfig.EdgeTestParameters.DisableInjectionKey),
		DisableInjectionValue: types.StringValue(cpcConfig.EdgeTestParameters.DisableInjectionValue),
		ForceInjectionKey:     types.StringValue(cpcConfig.EdgeTestParameters.ForceInjectionKey),
		ForceInjectionValue:   types.StringValue(cpcConfig.EdgeTestParameters.ForceInjectionValue),
	}
	var scriptTagAttrsList types.List
	if len(cpcConfig.InjectionCriteria.ScriptTagAttributes) == 0 {
		scriptTagAttrsList = types.ListNull(types.StringType)
	} else {
		l, d := types.ListValueFrom(ctx, types.StringType, cpcConfig.InjectionCriteria.ScriptTagAttributes)
		resp.Diagnostics.Append(d...)
		scriptTagAttrsList = l
	}

	data.InjectionCriteria = &clientsideProtectionInjectionCriteriaModel{
		InjectionPathIncludes: stringPtrOrNull(cpcConfig.InjectionCriteria.InjectionPathIncludes),
		InjectionPathExcludes: stringPtrOrNull(cpcConfig.InjectionCriteria.InjectionPathExcludes),
		InjectionPercent:      types.Int64Value(int64(cpcConfig.InjectionCriteria.InjectionPercent)),
		InjectionPolicy:       stringPtrOrNull(cpcConfig.InjectionCriteria.InjectionPolicy),
		ScriptTagAttributes:   scriptTagAttrsList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
