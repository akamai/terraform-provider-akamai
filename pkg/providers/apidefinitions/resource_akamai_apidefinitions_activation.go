package apidefinitions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf/validators"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &activationResource{}
	_ resource.ResourceWithConfigure   = &activationResource{}
	_ resource.ResourceWithImportState = &activationResource{}
)

type activationResource struct{}

type activationResourceModel struct {
	EndpointID             types.Int64  `tfsdk:"api_id"`
	Version                types.Int64  `tfsdk:"version"`
	Network                types.String `tfsdk:"network"`
	Notes                  types.String `tfsdk:"notes"`
	NotificationRecipients types.Set    `tfsdk:"notification_recipients"`
	Status                 types.String `tfsdk:"status"`
	AutoAckWarnings        types.Bool   `tfsdk:"auto_acknowledge_warnings"`
}

// NewActivationResource returns new api definition activation resource
func NewActivationResource() resource.Resource {
	return &activationResource{}
}

// Metadata implements datasource.DataSource.
func (r *activationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_apidefinitions_activation"
}

// Configure implements datasource.DataSource.
func (r *activationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	metaConfig, ok := req.ProviderData.(meta.Meta)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if client == nil {
		client = apidefinitions.Client(metaConfig.Session())
	}
	if clientV0 == nil {
		clientV0 = v0.Client(metaConfig.Session())
	}
}

// Schema implements datasource.DataSource.
func (r *activationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "API Definition configuration activation",
		Attributes: map[string]schema.Attribute{
			"api_id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Description: "Unique identifier of the API",
			},
			"version": schema.Int64Attribute{
				Required:    true,
				Description: "Version of the API to be activated",
			},
			"network": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					modifiers.PreventStringUpdate(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						string(apidefinitions.ActivationNetworkStaging),
						string(apidefinitions.ActivationNetworkProduction),
					),
				},
				Description: "Network on which to activate the API version (STAGING or PRODUCTION)",
			},
			"notes": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
				Description: "Notes describing the activation",
			},
			"notification_recipients": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.SizeBetween(0, 30),
					setvalidator.ValueStringsAre(
						validators.EmailValidator{},
					),
				},
				Description: "List of email addresses to be notified with the results of the activation",
			},
			"auto_acknowledge_warnings": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Automatically acknowledge all warnings for activation to continue. Default is false",
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The results of the activation",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *activationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating API Definitions Activation Resource")

	var plan *activationResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.create(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *activationResource) create(ctx context.Context, data *activationResourceModel) diag.Diagnostics {
	return r.handleActivation(ctx, data)
}

// Read implements resource.Resource.
func (r *activationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading API Definitions Activation")

	var state *activationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.read(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *activationResource) read(ctx context.Context, data *activationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	endpoint, err := getEndpoint(ctx, data.EndpointID.ValueInt64())
	if err != nil {
		diags.AddError("Unable to read Endpoint", err.Error())
		return diags
	}

	data.populateFrom(*endpoint)

	return diags
}

// Update implements resource.Resource.
func (r *activationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating API Definitions Activation")

	var plan *activationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.update(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *activationResource) update(ctx context.Context, plan *activationResourceModel) diag.Diagnostics {
	return r.handleActivation(ctx, plan)
}

// Delete implements resource.Resource.
func (r *activationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting API Definitions Activation")

	var state *activationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	network := apidefinitions.NetworkType(state.Network.ValueString())

	diags := deactivateEndpointOnNetwork(ctx, state.EndpointID.ValueInt64(), network)
	resp.Diagnostics.Append(diags...)
}

// ImportState implements resource's ImportState method
func (r *activationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing API Definitions Activation resource")

	parts := strings.Split(req.ID, ":")

	if len(parts) != 2 {
		resp.Diagnostics.AddError(fmt.Sprintf("ID '%s' incorrectly formatted: should be 'API_ID:NETWORK'", req.ID), "")
		return
	}

	endpointID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid API id '%v'", parts[0]), "")
		return
	}

	network := apidefinitions.NetworkType(parts[1])
	if !(network == apidefinitions.ActivationNetworkStaging || network == apidefinitions.ActivationNetworkProduction) {
		resp.Diagnostics.AddError(fmt.Sprintf("invalid network value %s; must be either %s or %s", parts[1], apidefinitions.ActivationNetworkStaging, apidefinitions.ActivationNetworkProduction), "")
		return
	}

	endpoint, err := getEndpoint(ctx, endpointID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read Endpoint", err.Error())
		return
	}

	stateOnNetwork := getStateOnNetwork(network, *endpoint)

	if !stateOnNetwork.IsActive() {
		resp.Diagnostics.AddError(fmt.Sprintf("API is not active on the network %s", network), "")
		return
	}

	data := activationResourceModel{
		EndpointID:             types.Int64Value(endpointID),
		Version:                types.Int64Value(*stateOnNetwork.VersionNumber),
		Network:                types.StringValue(string(network)),
		Notes:                  types.StringNull(),
		NotificationRecipients: types.SetNull(types.StringType),
		Status:                 types.StringValue(string(apidefinitions.ActivationStatusActive)),
		AutoAckWarnings:        types.BoolValue(false),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *activationResource) handleActivation(ctx context.Context, state *activationResourceModel) diag.Diagnostics {
	tflog.Debug(ctx, "handleActivation")

	var diags diag.Diagnostics
	endpoint, err := getEndpoint(ctx, state.EndpointID.ValueInt64())
	if err != nil {
		diags.AddError("Unable to read Endpoint", err.Error())
		return diags
	}

	network := apidefinitions.NetworkType(state.Network.ValueString())
	var notificationRecipients []string
	if err := state.NotificationRecipients.ElementsAs(ctx, &notificationRecipients, false); err != nil {
		diags.Append(err...)
		return diags
	}

	endpointID := state.EndpointID.ValueInt64()
	versionToActivate := state.Version.ValueInt64()
	networkToActivate := apidefinitions.NetworkType(state.Network.ValueString())

	networkState := getStateOnNetwork(network, *endpoint)

	if networkState.Status != nil && *networkState.Status == apidefinitions.ActivationStatusFailed {
		resp, err := client.CloneEndpointVersion(ctx, apidefinitions.CloneEndpointVersionRequest{
			VersionNumber: versionToActivate,
			APIEndpointID: endpointID,
		})

		if err != nil {
			diags.AddError("Unable to clone an API Version", err.Error())
			return diags
		}

		versionToActivate = resp.VersionNumber
	}

	if shouldActivate(networkState, versionToActivate) {
		tflog.Debug(ctx, fmt.Sprintf("API Activation: Activating API %d Version %d on network %s", endpointID, versionToActivate, networkToActivate))

		var verifyRequest = apidefinitions.VerifyVersionRequest{
			VersionNumber: versionToActivate,
			APIEndpointID: endpointID,
			Body: apidefinitions.VerifyVersionRequestBody{
				Networks: []apidefinitions.NetworkType{networkToActivate},
			},
		}

		var activationRequest = apidefinitions.ActivateVersionRequest{
			VersionNumber: versionToActivate,
			APIEndpointID: endpointID,
			Body: apidefinitions.ActivationRequestBody{
				Networks:               []apidefinitions.NetworkType{networkToActivate},
				Notes:                  state.Notes.ValueString(),
				NotificationRecipients: notificationRecipients,
			},
		}

		alerts, err := client.VerifyVersion(ctx, verifyRequest)

		if err != nil {
			diags.AddError("Activation Verification Failed", err.Error())
			return diags
		}

		alertsBySeverity := groupAlerts(alerts)

		if len(alertsBySeverity[apidefinitions.SeverityError]) > 0 {
			errors := formatAlerts(alertsBySeverity[apidefinitions.SeverityError])
			diags.AddError("Unable to proceed due to Activation Errors, fix errors to proceed", errors)
			return diags
		}

		if len(alertsBySeverity[apidefinitions.SeverityWarning]) > 0 && !state.AutoAckWarnings.ValueBool() {
			warnings := formatAlerts(alertsBySeverity[apidefinitions.SeverityWarning])
			diags.AddError("Unable to proceed due to Activation Warnings, set auto_acknowledge_warnings to acknowledge warnings and proceed", warnings)
			return diags
		}

		err = startActivation(ctx, activationRequest)
		if err != nil {
			diags.AddError("Activation Failed", err.Error())
			return diags
		}

		endpoint, diags = pollActivation(ctx, endpointID, versionToActivate, network)
		if diags != nil {
			diags.Append(diags...)
			return diags
		}

	}

	state.populateFrom(*endpoint)
	return diags
}

func formatAlerts(alerts []apidefinitions.VerifyVersionAlert) string {
	output := ""
	for _, alert := range alerts {
		output += alert.Detail + "\n\n"
	}
	return output
}

func (m *activationResourceModel) isStaging() bool {
	network := apidefinitions.NetworkType(m.Network.ValueString())
	return network == apidefinitions.ActivationNetworkStaging
}

func (m *activationResourceModel) populateFrom(resp apidefinitions.EndpointDetail) diag.Diagnostics {
	var diags diag.Diagnostics

	if m.isStaging() {
		if resp.StagingVersion.VersionNumber == nil {
			m.Status = types.StringValue("")
		} else {
			m.Status = types.StringValue(string(*resp.StagingVersion.Status))
		}
	} else {
		if resp.ProductionVersion.VersionNumber == nil {
			m.Status = types.StringValue("")
		} else {
			m.Status = types.StringValue(string(*resp.ProductionVersion.Status))
		}
	}

	return diags
}

func groupAlerts(items []apidefinitions.VerifyVersionAlert) map[apidefinitions.Severity][]apidefinitions.VerifyVersionAlert {
	grouped := make(map[apidefinitions.Severity][]apidefinitions.VerifyVersionAlert)

	for _, item := range items {
		key := item.Severity
		grouped[key] = append(grouped[key], item)
	}

	return grouped
}
