package appsec

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &clientsideProtectionConfigurationResource{}
	_ resource.ResourceWithConfigure   = &clientsideProtectionConfigurationResource{}
	_ resource.ResourceWithImportState = &clientsideProtectionConfigurationResource{}
)

type clientsideProtectionConfigurationResource struct {
	meta.Resource
}

type clientsideProtectionConfigurationModel struct {
	ConfigID           types.Int64                                 `tfsdk:"config_id"`
	SecurityPolicyID   types.String                                `tfsdk:"security_policy_id"`
	EdgeInjection      clientsideProtectionEdgeInjectionModel      `tfsdk:"edge_injection"`
	EdgeTestParameters clientsideProtectionEdgeTestParametersModel `tfsdk:"edge_test_parameters"`
	InjectionCriteria  clientsideProtectionInjectionCriteriaModel  `tfsdk:"injection_criteria"`
}

type clientsideProtectionEdgeInjectionModel struct {
	LoadScriptAsync types.Bool `tfsdk:"load_script_async"`
}

type clientsideProtectionEdgeTestParametersModel struct {
	DisableInjectionKey   types.String `tfsdk:"disable_injection_key"`
	DisableInjectionValue types.String `tfsdk:"disable_injection_value"`
	ForceInjectionKey     types.String `tfsdk:"force_injection_key"`
	ForceInjectionValue   types.String `tfsdk:"force_injection_value"`
}

type clientsideProtectionInjectionCriteriaModel struct {
	InjectionPathIncludes types.String `tfsdk:"injection_path_includes"`
	InjectionPathExcludes types.String `tfsdk:"injection_path_excludes"`
	InjectionPercent      types.Int64  `tfsdk:"injection_percent"`
	InjectionPolicy       types.String `tfsdk:"injection_policy"`
}

// NewClientsideProtectionConfigurationResource returns a new clientside protection configuration resource.
func NewClientsideProtectionConfigurationResource() resource.Resource {
	return &clientsideProtectionConfigurationResource{}
}

func (r *clientsideProtectionConfigurationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_appsec_clientside_protection_configuration"
}

func (r *clientsideProtectionConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Client-side Protection & Compliance (CPC) configuration resource.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"edge_injection": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Edge injection settings.",
				Attributes: map[string]schema.Attribute{
					"load_script_async": schema.BoolAttribute{
						Required:    true,
						Description: "Whether to load the protection script asynchronously.",
					},
				},
			},
			"edge_test_parameters": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Edge test parameters.",
				Attributes: map[string]schema.Attribute{
					"disable_injection_key": schema.StringAttribute{
						Required:    true,
						Description: "Query string key used to disable script injection.",
					},
					"disable_injection_value": schema.StringAttribute{
						Required:    true,
						Description: "Query string value used to disable script injection.",
					},
					"force_injection_key": schema.StringAttribute{
						Required:    true,
						Description: "Query string key used to force script injection.",
					},
					"force_injection_value": schema.StringAttribute{
						Required:    true,
						Description: "Query string value used to force script injection.",
					},
				},
			},
			"injection_criteria": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Injection criteria.",
				Attributes: map[string]schema.Attribute{
					"injection_path_includes": schema.StringAttribute{
						Optional:    true,
						Description: "Paths where script injection is enabled.",
					},
					"injection_path_excludes": schema.StringAttribute{
						Optional:    true,
						Description: "Paths where script injection is disabled.",
					},
					"injection_percent": schema.Int64Attribute{
						Required:    true,
						Description: "Percentage of requests that receive script injection (0-100).",
						Validators: []validator.Int64{
							int64validator.Between(0, 100),
						},
					},
					"injection_policy": schema.StringAttribute{
						Optional:    true,
						Description: "Controls where in the page the script is injected.",
					},
				},
			},
		},
	}
}

func (r *clientsideProtectionConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating clientside protection configuration resource")

	var plan clientsideProtectionConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.upsert(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clientsideProtectionConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading clientside protection configuration resource")

	var state clientsideProtectionConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.read(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *clientsideProtectionConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating clientside protection configuration resource")

	var plan clientsideProtectionConfigurationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.upsert(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete resets the configuration to defaults, since there is no delete API endpoint for this resource.
func (r *clientsideProtectionConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting clientside protection configuration resource")

	var state clientsideProtectionConfigurationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	configID := int(state.ConfigID.ValueInt64())
	client := r.Client.GetAPPSEC()

	version, err := getModifiableConfigVersion(ctx, configID, "clientsideProtectionConfiguration", client)
	if err != nil {
		resp.Diagnostics.AddError("could not get modifiable config version", err.Error())
		return
	}

	_, err = client.UpdateClientsideProtectionConfiguration(ctx, appsec.UpdateClientsideProtectionConfigurationRequest{
		ConfigID:         state.ConfigID.ValueInt64(),
		Version:          int64(version),
		SecurityPolicyID: state.SecurityPolicyID.ValueString(),
		Body: appsec.UpdateClientsideProtectionConfigurationRequestBody{
			EdgeInjection:      appsec.ClientsideProtectionEdgeInjectionRequest{LoadScriptAsync: false},
			EdgeTestParameters: appsec.ClientsideProtectionEdgeTestParameters{},
			InjectionCriteria: appsec.ClientsideProtectionInjectionCriteria{
				InjectionPathIncludes: ptr.To(""),
				InjectionPathExcludes: ptr.To(""),
				InjectionPercent:      0,
				InjectionPolicy:       ptr.To(""),
			},
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("reset clientside protection configuration to defaults failed", err.Error())
	}
}

// ImportState parses the import ID in the format "CONFIG_ID:POLICY_ID".
func (r *clientsideProtectionConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing clientside protection configuration resource")

	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("invalid import ID %q: expected format is 'CONFIG_ID:POLICY_ID'", req.ID),
			"",
		)
		return
	}

	configID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("invalid import ID %q: CONFIG_ID must be an integer", req.ID),
			"",
		)
		return
	}

	state := clientsideProtectionConfigurationModel{
		ConfigID:         types.Int64Value(configID),
		SecurityPolicyID: types.StringValue(parts[1]),
	}

	resp.Diagnostics.Append(r.read(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *clientsideProtectionConfigurationResource) upsert(ctx context.Context, data *clientsideProtectionConfigurationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	configID := int(data.ConfigID.ValueInt64())
	client := r.Client.GetAPPSEC()

	version, err := getModifiableConfigVersion(ctx, configID, "clientsideProtectionConfiguration", client)
	if err != nil {
		diags.AddError("could not get modifiable config version", err.Error())
		return diags
	}

	updateReq := appsec.UpdateClientsideProtectionConfigurationRequest{
		ConfigID:         data.ConfigID.ValueInt64(),
		Version:          int64(version),
		SecurityPolicyID: data.SecurityPolicyID.ValueString(),
		Body: appsec.UpdateClientsideProtectionConfigurationRequestBody{
			EdgeInjection: appsec.ClientsideProtectionEdgeInjectionRequest{
				LoadScriptAsync: data.EdgeInjection.LoadScriptAsync.ValueBool(),
			},
			EdgeTestParameters: appsec.ClientsideProtectionEdgeTestParameters{
				DisableInjectionKey:   data.EdgeTestParameters.DisableInjectionKey.ValueString(),
				DisableInjectionValue: data.EdgeTestParameters.DisableInjectionValue.ValueString(),
				ForceInjectionKey:     data.EdgeTestParameters.ForceInjectionKey.ValueString(),
				ForceInjectionValue:   data.EdgeTestParameters.ForceInjectionValue.ValueString(),
			},
			InjectionCriteria: appsec.ClientsideProtectionInjectionCriteria{
				InjectionPathIncludes: data.InjectionCriteria.InjectionPathIncludes.ValueStringPointer(),
				InjectionPathExcludes: data.InjectionCriteria.InjectionPathExcludes.ValueStringPointer(),
				InjectionPercent:      int(data.InjectionCriteria.InjectionPercent.ValueInt64()),
				InjectionPolicy:       data.InjectionCriteria.InjectionPolicy.ValueStringPointer(),
			},
		},
	}

	_, err = client.UpdateClientsideProtectionConfiguration(ctx, updateReq)
	if err != nil {
		diags.AddError("update clientside protection configuration request failed", err.Error())
		return diags
	}

	return r.read(ctx, data)
}

func (r *clientsideProtectionConfigurationResource) read(ctx context.Context, data *clientsideProtectionConfigurationModel) diag.Diagnostics {
	var diags diag.Diagnostics

	configID := int(data.ConfigID.ValueInt64())
	client := r.Client.GetAPPSEC()

	version, err := getLatestConfigVersion(ctx, configID, client)
	if err != nil {
		diags.AddError("could not get latest config version", err.Error())
		return diags
	}

	cpcConfig, err := client.GetClientsideProtectionConfiguration(ctx, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         data.ConfigID.ValueInt64(),
		Version:          int64(version),
		SecurityPolicyID: data.SecurityPolicyID.ValueString(),
	})
	if err != nil {
		diags.AddError("calling 'GetClientsideProtectionConfiguration'", err.Error())
		return diags
	}

	data.EdgeInjection = clientsideProtectionEdgeInjectionModel{
		LoadScriptAsync: types.BoolValue(cpcConfig.EdgeInjection.LoadScriptAsync),
	}
	data.EdgeTestParameters = clientsideProtectionEdgeTestParametersModel{
		DisableInjectionKey:   types.StringValue(cpcConfig.EdgeTestParameters.DisableInjectionKey),
		DisableInjectionValue: types.StringValue(cpcConfig.EdgeTestParameters.DisableInjectionValue),
		ForceInjectionKey:     types.StringValue(cpcConfig.EdgeTestParameters.ForceInjectionKey),
		ForceInjectionValue:   types.StringValue(cpcConfig.EdgeTestParameters.ForceInjectionValue),
	}
	data.InjectionCriteria = clientsideProtectionInjectionCriteriaModel{
		InjectionPathIncludes: stringPtrOrNull(cpcConfig.InjectionCriteria.InjectionPathIncludes),
		InjectionPathExcludes: stringPtrOrNull(cpcConfig.InjectionCriteria.InjectionPathExcludes),
		InjectionPercent:      types.Int64Value(int64(cpcConfig.InjectionCriteria.InjectionPercent)),
		InjectionPolicy:       stringPtrOrNull(cpcConfig.InjectionCriteria.InjectionPolicy),
	}

	return diags
}

func stringPtrOrNull(value *string) types.String {
	if value == nil {
		return types.StringNull()
	}
	return tf.StringValueOrNullIfEmpty(*value)
}
