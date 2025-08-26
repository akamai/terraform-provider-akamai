package mtlstruststore

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &caSetActivationResource{}
	_ resource.ResourceWithConfigure   = &caSetActivationResource{}
	_ resource.ResourceWithModifyPlan  = &caSetActivationResource{}
	_ resource.ResourceWithImportState = &caSetActivationResource{}

	pollingInterval = 5 * time.Second
)

var (
	onlyTimeoutChangeWarn = diag.NewWarningDiagnostic("Update with no API calls", "requested only timeout change; API won't be called")
)

type caSetActivationResource struct {
	meta              meta.Meta
	deleteTimeout     time.Duration
	activationTimeout time.Duration
}

// NewCASetActivationResource returns a new akamai_mtlstruststore_ca_set_activation resource.
func NewCASetActivationResource() resource.Resource {
	return &caSetActivationResource{
		deleteTimeout:     1 * time.Hour,
		activationTimeout: 1 * time.Hour,
	}
}

type caSetActivationResourceModel struct {
	CASetID      types.String   `tfsdk:"ca_set_id"`
	Version      types.Int64    `tfsdk:"version"`
	Network      types.String   `tfsdk:"network"`
	ID           types.Int64    `tfsdk:"id"`
	CreatedBy    types.String   `tfsdk:"created_by"`
	CreatedDate  types.String   `tfsdk:"created_date"`
	ModifiedBy   types.String   `tfsdk:"modified_by"`
	ModifiedDate types.String   `tfsdk:"modified_date"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

type caSetActivationInfo struct {
	ExpectedType string
	Label        string
	CASetID      string
	Version      int64
	ActivationID int64
	RetryAfter   time.Time
}

func (c *caSetActivationResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "akamai_mtlstruststore_ca_set_activation"
}

func (c *caSetActivationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ca_set_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{modifiers.PreventStringUpdate()},
				Description:   "Uniquely Identifies a CA set.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"version": schema.Int64Attribute{
				Required:    true,
				Description: "Identifies the version of the CA set.",
			},
			"network": schema.StringAttribute{
				Required:      true,
				Description:   "Indicates the network for any activation-related activities, either 'STAGING' or 'PRODUCTION'.",
				PlanModifiers: []planmodifier.String{modifiers.PreventStringUpdate()},
				Validators: []validator.String{
					stringvalidator.OneOf("STAGING", "PRODUCTION"),
				},
			},
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Uniquely Identifies a CA set Activation.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who submitted the activation request.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date the activation request was submitted in ISO-8601 format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "The user who completed the activation.",
			},
			"modified_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date the activation request was modified in ISO-8601 format.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Delete:            true,
				Update:            true,
				Create:            true,
				CreateDescription: "Optional configurable resource create timeout. By default it's 1h.",
				DeleteDescription: "Optional configurable resource delete timeout. By default it's 1h.",
				UpdateDescription: "Optional configurable resource update timeout. By default it's 1h.",
			}),
		},
	}
}

func (c *caSetActivationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unexpected Resource Configure Type",
				fmt.Sprintf("Expected meta.Meta, got: %T. Please report this issue to the provider developers.", req.ProviderData),
			)
		}
	}()

	c.meta = meta.Must(req.ProviderData)
}

func (c *caSetActivationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating CASetActivation resource")

	var plan caSetActivationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	activationTimeout, diags := plan.Timeouts.Create(ctx, c.activationTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := c.upsert(ctx, &plan, activationTimeout); err != nil {
		resp.Diagnostics.AddError("create CA set activation failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *caSetActivationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating CASetActivation Resource")

	var plan caSetActivationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var oldState caSetActivationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if onlyTimeoutChanged(&oldState, &plan) {
		oldState.Timeouts = plan.Timeouts
		resp.Diagnostics.Append(resp.State.Set(ctx, &oldState)...)
		return
	}

	activationTimeout, diags := plan.Timeouts.Update(ctx, c.activationTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := c.upsert(ctx, &plan, activationTimeout); err != nil {
		resp.Diagnostics.AddError("update a CA set activation failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c *caSetActivationResource) upsert(ctx context.Context, plan *caSetActivationResourceModel, activationTimeout time.Duration) error {
	client = Client(c.meta)

	caSetID := plan.CASetID.ValueString()
	version := plan.Version.ValueInt64()
	network := plan.Network.ValueString()

	resp, err := client.GetCASetVersion(ctx, mtlstruststore.GetCASetVersionRequest{
		CASetID: caSetID,
		Version: version,
	})
	if err != nil {
		switch {
		case errors.Is(err, mtlstruststore.ErrGetCASetNotFound):
			return fmt.Errorf("CA set with ID %s not found: %w", caSetID, err)
		case errors.Is(err, mtlstruststore.ErrGetCASetVersionNotFound):
			return fmt.Errorf("CA set version %d not found for CA set ID %s: %w", version, caSetID, err)
		default:
			return fmt.Errorf("failed to get CA set version for ID %s version %d: %w", caSetID, version, err)
		}
	}

	if (strings.EqualFold(network, string(mtlstruststore.NetworkStaging)) && resp.StagingStatus == "ACTIVE") ||
		strings.EqualFold(network, string(mtlstruststore.NetworkProduction)) && resp.ProductionStatus == "ACTIVE" {
		activation, err := c.findLatestActivation(ctx, client, caSetID, version, network)
		if err != nil {
			return fmt.Errorf("error fetching activations for version %d : %w", version, err)
		}
		if activation != nil {
			plan.setActivateCASetActivationData(activation)
			return nil
		}
	}

	ongoingActivation, err := checkOngoingCASetOperation(ctx, client, plan)
	if err != nil {
		return err
	}

	var activation *mtlstruststore.ActivateCASetVersionResponse

	if ongoingActivation != nil {
		switch ongoingActivation.ActivationType {
		case "ACTIVATE":
			if ongoingActivation.Version == version {
				// Fetch ongoing activation for the same version.
				activation = ongoingActivation
			} else {
				// Return error for different version activation in progress.
				return fmt.Errorf("activation already in progress for version %d", ongoingActivation.Version)
			}
		case "DEACTIVATE":
			// Return error for ongoing deactivation.
			return fmt.Errorf("deactivation in progress for version %d, cannot activate", ongoingActivation.Version)
		default:
			return fmt.Errorf("unsupported activation type %s", ongoingActivation.ActivationType)
		}
	} else {
		// Create a new activation request.
		activation, err = client.ActivateCASetVersion(ctx, mtlstruststore.ActivateCASetVersionRequest{
			CASetID: caSetID,
			Version: version,
			Network: mtlstruststore.ActivationNetwork(network),
		})
		if err != nil {
			return err
		}
	}

	activationInfo := caSetActivationInfo{
		ExpectedType: "ACTIVATE",
		Label:        "activation",
		CASetID:      activation.CASetID,
		Version:      activation.Version,
		ActivationID: activation.ActivationID,
		RetryAfter:   activation.RetryAfter,
	}

	status, err := waitForActivationOrDeactivation(ctx, activationTimeout, client, activationInfo)
	if err != nil {
		return fmt.Errorf("activation polling failed: %w", err)
	}

	plan.setCASetActivationData(status)

	return nil
}

func (c *caSetActivationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading CASetActivation Resource")
	var state caSetActivationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	shouldRemove, err := c.read(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError("read CA set activation failed", err.Error())
		return
	}

	if shouldRemove {
		tflog.Warn(ctx, "Removing resource from state: either missing or inactive")
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *caSetActivationResource) read(ctx context.Context, data *caSetActivationResourceModel) (remove bool, err error) {
	client := Client(c.meta)

	caSetID := data.CASetID.ValueString()
	version := data.Version.ValueInt64()
	network := data.Network.ValueString()

	caSetResp, err := client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
		CASetID: caSetID,
	})
	if err != nil {
		if errors.Is(err, mtlstruststore.ErrGetCASetNotFound) {
			return true, nil
		}
		return false, fmt.Errorf("failed to get CA set: %w", err)
	}

	var activeVersion *int64
	switch network {
	case "STAGING":
		activeVersion = caSetResp.StagingVersion
	case "PRODUCTION":
		activeVersion = caSetResp.ProductionVersion
	default:
		return true, fmt.Errorf("unsupported network: %s", network)
	}

	if activeVersion == nil {
		return true, nil
	}

	activation, err := c.findLatestActivation(ctx, client, caSetID, version, network)
	if err != nil {
		return false, fmt.Errorf("failed to find activation: %w", err)
	}
	if activation != nil {
		data.setActivateCASetActivationData(activation)
		return false, nil
	}

	return false, fmt.Errorf("no activation found for CASetID %s, version %d, network %s", caSetID, version, network)
}

func (c *caSetActivationResource) findLatestActivation(ctx context.Context, client mtlstruststore.MTLSTruststore, caSetID string, version int64, network string) (*mtlstruststore.ActivateCASetVersionResponse, error) {
	resp, err := client.ListCASetVersionActivations(ctx, mtlstruststore.ListCASetVersionActivationsRequest{
		CASetID: caSetID,
		Version: version,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list activations: %w", err)
	}

	for _, activation := range resp.Activations {
		if activation.Network == network && activation.ActivationStatus == "COMPLETE" && activation.ActivationType == "ACTIVATE" {
			return &activation, nil
		}
	}

	return nil, nil
}

func (c *caSetActivationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting CASetActivation Resource")

	var state *caSetActivationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client = Client(c.meta)

	caSetID := state.CASetID.ValueString()
	version := state.Version.ValueInt64()
	network := state.Network.ValueString()

	ongoingOperation, err := checkOngoingCASetOperation(ctx, client, state)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to deactivate CA set ID %s version %d", state.CASetID.ValueString(), state.Version.ValueInt64()), err.Error())
		return
	}

	var deactivation *mtlstruststore.DeactivateCASetVersionResponse

	if ongoingOperation != nil {
		switch ongoingOperation.ActivationType {
		case "DEACTIVATE":
			if ongoingOperation.Version == version {
				// Fetch ongoing deactivation for the same version.
				deactivation = (*mtlstruststore.DeactivateCASetVersionResponse)(ongoingOperation)
			} else {
				// Return error for different version deactivation in progress.
				resp.Diagnostics.AddError("Deactivation Error", fmt.Sprintf("deactivation already in progress for version %d", ongoingOperation.Version))
				return
			}
		case "ACTIVATE":
			// Return error for ongoing activation.
			resp.Diagnostics.AddError("Deactivation Error", fmt.Sprintf("activation in progress for version %d, cannot deactivate", ongoingOperation.Version))
			return
		default:
			resp.Diagnostics.AddError("Deactivation Error", fmt.Sprintf("unsupported activation type %s", ongoingOperation.ActivationType))
			return
		}
	} else {
		// Create a new deactivation request.
		deactivation, err = client.DeactivateCASetVersion(ctx, mtlstruststore.DeactivateCASetVersionRequest{
			CASetID: caSetID,
			Version: version,
			Network: mtlstruststore.ActivationNetwork(network),
		})
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to deactivate CA set ID %s version %d", state.CASetID.ValueString(), state.Version.ValueInt64()), err.Error())
			return
		}
	}

	deactivationInfo := caSetActivationInfo{
		ExpectedType: "DEACTIVATE",
		Label:        "deactivation",
		CASetID:      caSetID,
		Version:      version,
		ActivationID: deactivation.ActivationID,
		RetryAfter:   deactivation.RetryAfter,
	}

	timeout, diags := state.Timeouts.Delete(ctx, c.deleteTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err = waitForActivationOrDeactivation(ctx, timeout, client, deactivationInfo)
	if err != nil {
		resp.Diagnostics.AddError("Failed to deactivate CA set version %d", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func checkOngoingCASetOperation(ctx context.Context, client mtlstruststore.MTLSTruststore, res *caSetActivationResourceModel) (*mtlstruststore.ActivateCASetVersionResponse, error) {
	caSetID := res.CASetID.ValueString()
	network := res.Network.ValueString()

	activationsResp, err := client.ListCASetActivations(ctx, mtlstruststore.ListCASetActivationsRequest{
		CASetID: caSetID,
	})
	if err != nil {
		switch {
		case errors.Is(err, mtlstruststore.ErrGetCASetNotFound):
			return nil, fmt.Errorf("CA set with ID %s not found: %w", caSetID, err)
		default:
			return nil, fmt.Errorf("could not retrieve activation details for CA set ID %s: %w", caSetID, err)
		}
	}

	for _, activation := range activationsResp.Activations {
		if activation.Network != network {
			continue
		}

		if activation.ActivationStatus == "IN_PROGRESS" {
			return &activation, nil
		}
	}

	return nil, nil
}

func waitForActivationOrDeactivation(ctx context.Context, timeout time.Duration, client mtlstruststore.MTLSTruststore, activation caSetActivationInfo) (*mtlstruststore.GetCASetVersionActivationResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	activationPollInterval := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while waiting for %s of CA Set %s, Version %d: %w",
				activation.Label, activation.CASetID, activation.Version, ctx.Err())
		case <-time.After(activationPollInterval):
			activationResp, err := client.GetCASetVersionActivation(ctx, mtlstruststore.GetCASetVersionActivationRequest{
				ActivationID: activation.ActivationID,
				CASetID:      activation.CASetID,
				Version:      activation.Version,
			})
			if err != nil {
				return nil, fmt.Errorf("error checking %s status for CA Set %s, Version %d, ActivationID %d: %w",
					activation.Label, activation.CASetID, activation.Version, activation.ActivationID, err)
			}

			switch activationResp.ActivationStatus {
			case "COMPLETE":
				if activationResp.ActivationType == activation.ExpectedType {
					return activationResp, nil
				}
				return nil, fmt.Errorf("unexpected activation type: %s (expected %s)", activationResp.ActivationType, activation.ExpectedType)
			case "IN_PROGRESS":
				tflog.Debug(ctx, fmt.Sprintf("%s of CA set %s version %d in progress", activation.Label, activation.CASetID, activation.Version))
				if !activation.RetryAfter.IsZero() {
					activationPollInterval = time.Until(activation.RetryAfter)
				} else {
					activationPollInterval = pollingInterval
				}
			case "FAILED":
				return nil, fmt.Errorf("%s failed for CA Set %s, Version %d",
					activation.Label, activation.CASetID, activation.Version)
			default:
				return nil, fmt.Errorf("unknown activation status: %v for CA Set %s, Version %d",
					activationResp.ActivationStatus, activation.CASetID, activation.Version)
			}
		}
	}
}

func (c *caSetActivationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "Importing CA Set Activation Resource")

	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected format: 'caSetID:network', got: %q", req.ID),
		)
		return
	}

	caSetID, network := parts[0], parts[1]

	if caSetID == "" {
		resp.Diagnostics.AddError("Invalid CA set ID", "CA set ID cannot be empty.")
		return
	}

	if network != "STAGING" && network != "PRODUCTION" {
		resp.Diagnostics.AddError("Invalid network", fmt.Sprintf("Network must be 'STAGING' or 'PRODUCTION', got: %s", network))
		return
	}

	client := Client(c.meta)

	activationsResp, err := client.ListCASetActivations(ctx, mtlstruststore.ListCASetActivationsRequest{
		CASetID: caSetID,
	})
	if err != nil {
		switch {
		case errors.Is(err, mtlstruststore.ErrGetCASetNotFound):
			resp.Diagnostics.AddError("CA set not found", fmt.Sprintf("CA set with ID %s not found: %s", caSetID, err.Error()))
			return
		default:
			resp.Diagnostics.AddError("Failed to retrieve activation details", fmt.Sprintf("could not retrieve activation details for CA set ID %s: %s", caSetID, err.Error()))
			return
		}
	}

	var activatedVersion *mtlstruststore.ActivateCASetVersionResponse
	for _, activation := range activationsResp.Activations {
		if activation.Network != network {
			continue
		}
		if activation.ActivationStatus == "IN_PROGRESS" {
			resp.Diagnostics.AddError("Operation in progress", fmt.Sprintf("A CA set operation is already in progress: %s. Can only import completed activations.", caSetID))
			return
		}
		if activation.ActivationStatus == "COMPLETE" && activation.ActivationType == "ACTIVATE" {
			activatedVersion = &activation
			break
		}
	}

	if activatedVersion == nil {
		resp.Diagnostics.AddError("No active CA set", fmt.Sprintf("CA set with ID %s is not active in the %s network. Only completed activations can be imported.", caSetID, network))
		return
	}

	data := caSetActivationResourceModel{
		CASetID: types.StringValue(caSetID),
		Version: types.Int64Value(activatedVersion.Version),
		Network: types.StringValue(network),
		ID:      types.Int64Value(activatedVersion.ActivationID),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"delete": types.StringType,
				"create": types.StringType,
				"update": types.StringType,
			}),
		},
	}
	data.setActivateCASetActivationData(activatedVersion)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	tflog.Info(ctx, "Successfully imported CA Set Activation", map[string]any{
		"ca_set_id":     caSetID,
		"version":       activatedVersion.Version,
		"activation_id": activatedVersion.ActivationID,
		"network":       network,
	})
}

func (m *caSetActivationResourceModel) setActivateCASetActivationData(resp *mtlstruststore.ActivateCASetVersionResponse) {
	m.CASetID = types.StringValue(resp.CASetID)
	m.ID = types.Int64Value(resp.ActivationID)
	m.Version = types.Int64Value(resp.Version)
	m.CreatedBy = types.StringValue(resp.CreatedBy)
	m.CreatedDate = date.TimeRFC3339NanoValue(resp.CreatedDate)
	m.ModifiedBy = types.StringPointerValue(resp.ModifiedBy)
	m.ModifiedDate = date.TimeRFC3339NanoPointerValue(resp.ModifiedDate)
}

func (m *caSetActivationResourceModel) setCASetActivationData(resp *mtlstruststore.GetCASetVersionActivationResponse) {
	m.CASetID = types.StringValue(resp.CASetID)
	m.ID = types.Int64Value(resp.ActivationID)
	m.Version = types.Int64Value(resp.Version)
	m.CreatedBy = types.StringValue(resp.CreatedBy)
	m.CreatedDate = date.TimeRFC3339NanoValue(resp.CreatedDate)
	m.ModifiedBy = types.StringPointerValue(resp.ModifiedBy)
	m.ModifiedDate = date.TimeRFC3339NanoPointerValue(resp.ModifiedDate)
}

func (c *caSetActivationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if modifiers.IsDelete(req) {
		// Verify if CA set version is active and is in use before deleting.
		var state caSetActivationResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		client := Client(c.meta)

		versionResp, err := client.GetCASetVersion(ctx, mtlstruststore.GetCASetVersionRequest{
			CASetID: state.CASetID.ValueString(),
			Version: state.Version.ValueInt64(),
		})
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to get CA set version",
				fmt.Sprintf("Could not check associations for CA set ID %s version %d: %s",
					state.CASetID.ValueString(),
					state.Version.ValueInt64(),
					err.Error()),
			)
		} else {
			network := state.Network.ValueString()
			isActive := (network == "STAGING" && versionResp.StagingStatus == "ACTIVE") ||
				(network == "PRODUCTION" && versionResp.ProductionStatus == "ACTIVE")
			if isActive {
				caSetAssociations, err := client.ListCASetAssociations(ctx, mtlstruststore.ListCASetAssociationsRequest{
					CASetID: state.CASetID.ValueString(),
				})
				if err != nil {
					resp.Diagnostics.AddError("Listing CA set associations failed", err.Error())
					return
				}

				if len(caSetAssociations.Associations.Enrollments) > 0 || len(caSetAssociations.Associations.Properties) > 0 {
					tflog.Warn(ctx, "CA set with current version is in use and cannot be deleted", map[string]interface{}{
						"ca_set_id":   state.CASetID.ValueString(),
						"version":     state.Version.ValueInt64(),
						"enrollments": len(caSetAssociations.Associations.Enrollments),
						"properties":  len(caSetAssociations.Associations.Properties),
					})
					resp.Diagnostics.AddWarning(
						"CA Set in Use",
						fmt.Sprintf("The CA set with ID %s and version %d is still associated with one or more enrollments or properties and cannot be deleted. Details: %s",
							state.CASetID.ValueString(),
							state.Version.ValueInt64(),
							getAssociationDetails(caSetAssociations),
						),
					)
					return
				}
			}
		}
	}

	if modifiers.IsUpdate(req) {
		var state, plan caSetActivationResourceModel
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if onlyTimeoutChanged(&state, &plan) {
			resp.Diagnostics.Append(onlyTimeoutChangeWarn)
		}
		return
	}
}

func onlyTimeoutChanged(state, plan *caSetActivationResourceModel) bool {
	return state != nil && plan != nil &&
		state.Version == plan.Version &&
		!state.Timeouts.Equal(plan.Timeouts)
}
