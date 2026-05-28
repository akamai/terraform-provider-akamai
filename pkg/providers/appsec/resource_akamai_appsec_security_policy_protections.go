package appsec

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &securityPolicyProtectionsResource{}
	_ resource.ResourceWithConfigure   = &securityPolicyProtectionsResource{}
	_ resource.ResourceWithImportState = &securityPolicyProtectionsResource{}
)

const securityPolicyProtectionsResourceName = "akamai_appsec_security_policy_protections"

// securityPolicyProtectionsResource manages all policy protection controls.
// Use this resource to manage all or any subset of protection flags for a security policy.
// Related resources: this is the consolidated counterpart of single-control resources
// such as akamai_appsec_waf_protection, akamai_appsec_rate_protection,
// akamai_appsec_reputation_protection, and akamai_appsec_url_protection.
type securityPolicyProtectionsResource struct {
	meta.Resource
}

// securityPolicyProtectionsModel maps Terraform state/plan for full policy protections.
type securityPolicyProtectionsModel struct {
	ConfigID                       types.Int64  `tfsdk:"config_id"`
	SecurityPolicyID               types.String `tfsdk:"security_policy_id"`
	ApplyAccountProtectionControls types.Bool   `tfsdk:"apply_account_protection_controls"`
	ApplyAPIConstraints            types.Bool   `tfsdk:"apply_api_constraints"`
	ApplyApplicationLayerControls  types.Bool   `tfsdk:"apply_application_layer_controls"`
	ApplyBotmanControls            types.Bool   `tfsdk:"apply_botman_controls"`
	ApplyMalwareControls           types.Bool   `tfsdk:"apply_malware_controls"`
	ApplyNetworkLayerControls      types.Bool   `tfsdk:"apply_network_layer_controls"`
	ApplyRateControls              types.Bool   `tfsdk:"apply_rate_controls"`
	ApplyReputationControls        types.Bool   `tfsdk:"apply_reputation_controls"`
	ApplySlowPostControls          types.Bool   `tfsdk:"apply_slow_post_controls"`
	ApplyURLProtectionControls     types.Bool   `tfsdk:"apply_url_protection_controls"`
}

// NewSecurityPolicyProtectionsResource returns the consolidated resource for
// managing policy protections in one place.
func NewSecurityPolicyProtectionsResource() resource.Resource {
	return &securityPolicyProtectionsResource{}
}

// Metadata sets the Terraform type name for the full protections resource.
func (r *securityPolicyProtectionsResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = securityPolicyProtectionsResourceName
}

// Schema defines the full protections resource schema.
// Related resources: this schema consolidates controls that were historically
// managed through separate resources.
func (r *securityPolicyProtectionsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage all security policy protection controls in one resource.",
		Attributes: map[string]schema.Attribute{
			"config_id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier of the security configuration.",
				PlanModifiers: []planmodifier.Int64{
					modifiers.PreventInt64Update(),
				},
			},
			"security_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier of the security policy.",
				Validators: []validator.String{
					validators.NotEmptyString(),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
			},
			"apply_account_protection_controls": schema.BoolAttribute{Required: true, Description: "Whether to enable account protection controls."},
			"apply_api_constraints":             schema.BoolAttribute{Required: true, Description: "Whether to enable API constraints."},
			"apply_application_layer_controls":  schema.BoolAttribute{Required: true, Description: "Whether to enable application layer controls."},
			"apply_botman_controls":             schema.BoolAttribute{Required: true, Description: "Whether to enable botman controls."},
			"apply_malware_controls":            schema.BoolAttribute{Required: true, Description: "Whether to enable malware controls."},
			"apply_network_layer_controls":      schema.BoolAttribute{Required: true, Description: "Whether to enable network layer controls."},
			"apply_rate_controls":               schema.BoolAttribute{Required: true, Description: "Whether to enable rate controls."},
			"apply_reputation_controls":         schema.BoolAttribute{Required: true, Description: "Whether to enable reputation controls."},
			"apply_slow_post_controls":          schema.BoolAttribute{Required: true, Description: "Whether to enable slow post controls."},
			"apply_url_protection_controls":     schema.BoolAttribute{Required: true, Description: "Whether to enable URL protection controls."},
		},
	}
}

// Create applies all requested protection control values.
func (r *securityPolicyProtectionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan securityPolicyProtectionsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.updatePolicyProtections(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readState(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes state for the full protections resource from the API.
func (r *securityPolicyProtectionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state securityPolicyProtectionsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readState(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update reapplies all requested protection controls.
func (r *securityPolicyProtectionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan securityPolicyProtectionsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.updatePolicyProtections(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.readState(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete disables managed protection controls remotely and removes the resource from state.
func (r *securityPolicyProtectionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state securityPolicyProtectionsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ApplyAccountProtectionControls = types.BoolValue(false)
	state.ApplyAPIConstraints = types.BoolValue(false)
	state.ApplyApplicationLayerControls = types.BoolValue(false)
	state.ApplyBotmanControls = types.BoolValue(false)
	state.ApplyMalwareControls = types.BoolValue(false)
	state.ApplyNetworkLayerControls = types.BoolValue(false)
	state.ApplyRateControls = types.BoolValue(false)
	state.ApplyReputationControls = types.BoolValue(false)
	state.ApplySlowPostControls = types.BoolValue(false)
	state.ApplyURLProtectionControls = types.BoolValue(false)

	r.updatePolicyProtections(ctx, state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports full protections resource state from CONFIG_ID:SECURITY_POLICY_ID.
func (r *securityPolicyProtectionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("ID '%s' incorrectly formatted: should be 'CONFIG_ID:SECURITY_POLICY_ID'", req.ID),
			"",
		)
		return
	}

	configID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid configuration id '%v'", parts[0]), "")
		return
	}

	securityPolicyID := parts[1]
	if securityPolicyID == "" {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid security policy id '%v'", parts[1]), "")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("config_id"), configID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("security_policy_id"), securityPolicyID)...)
}

// updatePolicyProtections updates all control fields managed by the
// akamai_appsec_security_policy_protections resource.
// Related resources: intended as a single update path instead of multiple
// per-control resources.
func (r *securityPolicyProtectionsResource) updatePolicyProtections(ctx context.Context, data securityPolicyProtectionsModel, diags *diag.Diagnostics) {
	version, err := getModifiableConfigVersion(ctx, int(data.ConfigID.ValueInt64()), "securityPolicyProtections", r.Client.GetAPPSEC())
	if err != nil {
		diags.AddError("Failed to retrieve modifiable config version", err.Error())
		return
	}

	client := r.Client.GetAPPSEC()

	_, err = client.UpdatePolicyProtections(ctx, appsec.UpdatePolicyProtectionsRequest{
		ConfigID:                       int(data.ConfigID.ValueInt64()),
		Version:                        version,
		PolicyID:                       data.SecurityPolicyID.ValueString(),
		ApplyAccountProtectionControls: data.ApplyAccountProtectionControls.ValueBool(),
		ApplyAPIConstraints:            data.ApplyAPIConstraints.ValueBool(),
		ApplyApplicationLayerControls:  data.ApplyApplicationLayerControls.ValueBool(),
		ApplyBotmanControls:            data.ApplyBotmanControls.ValueBool(),
		ApplyMalwareControls:           data.ApplyMalwareControls.ValueBool(),
		ApplyNetworkLayerControls:      data.ApplyNetworkLayerControls.ValueBool(),
		ApplyRateControls:              data.ApplyRateControls.ValueBool(),
		ApplyReputationControls:        data.ApplyReputationControls.ValueBool(),
		ApplySlowPostControls:          data.ApplySlowPostControls.ValueBool(),
		ApplyURLProtectionControls:     data.ApplyURLProtectionControls.ValueBool(),
	})

	if err != nil {
		diags.AddError("Failed to update security policy protections", err.Error())
	}
}

// readState reads all current controls from the API into the full protections model.
func (r *securityPolicyProtectionsResource) readState(ctx context.Context, data *securityPolicyProtectionsModel) diag.Diagnostics {
	var diags diag.Diagnostics

	version, err := getLatestConfigVersion(ctx, int(data.ConfigID.ValueInt64()), r.Client.GetAPPSEC())
	if err != nil {
		diags.AddError("Unable to read latest config version", err.Error())
		return diags
	}

	client := r.Client.GetAPPSEC()
	result, err := client.GetPolicyProtections(ctx, appsec.GetPolicyProtectionsRequest{
		ConfigID: int(data.ConfigID.ValueInt64()),
		Version:  version,
		PolicyID: data.SecurityPolicyID.ValueString(),
	})
	if err != nil {
		diags.AddError("Failed to read security policy protections", err.Error())
		return diags
	}

	data.ApplyAccountProtectionControls = types.BoolValue(result.ApplyAccountProtectionControls)
	data.ApplyAPIConstraints = types.BoolValue(result.ApplyAPIConstraints)
	data.ApplyApplicationLayerControls = types.BoolValue(result.ApplyApplicationLayerControls)
	data.ApplyBotmanControls = types.BoolValue(result.ApplyBotmanControls)
	data.ApplyMalwareControls = types.BoolValue(result.ApplyMalwareControls)
	data.ApplyNetworkLayerControls = types.BoolValue(result.ApplyNetworkLayerControls)
	data.ApplyRateControls = types.BoolValue(result.ApplyRateControls)
	data.ApplyReputationControls = types.BoolValue(result.ApplyReputationControls)
	data.ApplySlowPostControls = types.BoolValue(result.ApplySlowPostControls)
	data.ApplyURLProtectionControls = types.BoolValue(result.ApplyURLProtectionControls)

	return diags
}
