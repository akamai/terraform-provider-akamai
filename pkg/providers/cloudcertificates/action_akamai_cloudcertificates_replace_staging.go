package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action              = &ReplaceStagingAction{}
	_ action.ActionWithConfigure = &ReplaceStagingAction{}
)

// ReplaceStagingAction represents the akamai_cloudcertificates_replace_staging action.
type ReplaceStagingAction struct {
	meta.Action
	activationPollingConfig
}

// NewReplaceStagingAction returns a new instance of the akamai_cloudcertificates_replace_staging action.
func NewReplaceStagingAction(config activationPollingConfig) func() action.Action {
	return func() action.Action {
		return &ReplaceStagingAction{activationPollingConfig: config}
	}
}

// ReplaceStagingActionModel represents the model for akamai_cloudcertificates_replace_staging action configuration.
type ReplaceStagingActionModel struct {
	LineageID    types.Int64    `tfsdk:"lineage_id"`
	GenerationID types.Int64    `tfsdk:"generation_id"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

// Metadata sets the full name of the akamai_cloudcertificates_replace_staging action.
func (a *ReplaceStagingAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_replace_staging"
}

// Schema sets the schema of the akamai_cloudcertificates_replace_staging action.
func (a *ReplaceStagingAction) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Replaces the generation currently deployed on staging with the lineage's current production " +
			"or previous production generation.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the lineage.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"generation_id": schema.Int64Attribute{
				Required: true,
				Description: "The unique identifier of the generation to use as the source for the replace-staging " +
					"operation. It must be the lineage's current production or previous production generation.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"timeouts": timeouts.AttributesWithOpts(ctx, timeouts.Opts{
				InvokeDescription: "Timeout for waiting on the replace-staging activation to reach a terminal " +
					"status. Defaults to 30m. Accepts a duration string, e.g. \"45m\" or \"1h\".",
			}),
		},
	}
}

// Invoke replaces the lineage's staging generation with generation_id, then waits for the resulting activation
// to reach a terminal status.
func (a *ReplaceStagingAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data ReplaceStagingActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invokeTimeout, diags := data.Timeouts.Invoke(ctx, a.defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cloud Certificates Replace Staging Action Invoke")

	lineageID := data.LineageID.ValueInt64()
	generationID := data.GenerationID.ValueInt64()
	client := a.Client.GetCloudCertificates()

	ctx, cancel := context.WithTimeout(ctx, invokeTimeout)
	defer cancel()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Replacing staging generation with generation %d for lineage %d...", generationID, lineageID),
	})

	result, err := client.ReplaceStagingLineage(ctx, cloudcertificates.ReplaceStagingLineageRequest{
		LineageID:    lineageID,
		GenerationID: generationID,
	})
	if err != nil {
		resp.Diagnostics.Append(replaceStagingLineageErrorDiagnostics(err, lineageID, generationID)...)
		return
	}
	if result == nil {
		resp.Diagnostics.AddWarning("Replace Staging Not Needed",
			fmt.Sprintf("Generation %d is already the current staging generation for lineage %d; nothing to replace.",
				generationID, lineageID))
		return
	}

	final, diags := awaitNetworkActivation(ctx, awaitNetworkActivationOpts{
		client:       client,
		lineageID:    lineageID,
		network:      cloudcertificates.TargetNetworkStaging,
		items:        result.Items,
		pollInterval: a.pollInterval,
		onFound: func(activationID int64) {
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("Replace-staging submitted as activation %d for lineage %d; waiting for it to complete...",
					activationID, lineageID),
			})
		},
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Replace-staging of lineage %d to generation %d finished with status %s.", lineageID, generationID, final.ActivationStatus),
	})

	resp.Diagnostics.AddWarning("Terraform State Not Updated",
		fmt.Sprintf("Lineage %d's staging generation was replaced with generation %d outside of Terraform. Unlike "+
			"akamai_cloudcertificates_rollback, this doesn't disturb akamai_cloudcertificates_upload - replace-staging "+
			"never touches head, so upload stays correctly bound. If you renew and sign a new head next, staging "+
			"reconciles forward automatically, but only when staging_generation_id references the upload resource's "+
			"generation_id (the recommended wiring), not a fixed value - promoting a fresh head to staging always "+
			"succeeds regardless of what staging transiently served here. If staging_generation_id is pinned to a "+
			"fixed value instead, a plain terraform apply sees no difference and does nothing until you update it "+
			"manually to the newly signed head's generation ID. Absent a renewal, manually setting "+
			"staging_generation_id to %d on an akamai_cloudcertificates_activation resource will fail since promotion "+
			"requires the target to already be the lineage's head; run a full export for this lineage to re-derive the "+
			"activation resource from live state (finding the real current staging activation) instead of "+
			"hand-patching the existing config.", lineageID, generationID, generationID))
}

// replaceStagingLineageErrorDiagnostics maps well-documented ReplaceStagingLineage sentinel errors to clear diagnostics.
func replaceStagingLineageErrorDiagnostics(err error, lineageID, generationID int64) diag.Diagnostics {
	if diags := sharedActivationErrorDiagnostics(err, lineageID, "replace-staging operation"); diags.HasError() {
		return diags
	}

	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrNoCurrentStaging):
		diags.AddError("No Current Staging", fmt.Sprintf("Lineage %d has no current staging generation to replace.", lineageID))
	case errors.Is(err, cloudcertificates.ErrReplaceStagingHeadNotOnStaging):
		diags.AddError("Invalid Replace Staging Request",
			fmt.Sprintf("Replacing staging for lineage %d is not valid in its current state. %s", lineageID, err.Error()))
	case errors.Is(err, cloudcertificates.ErrLineageBadRequest):
		diags.AddError("Invalid Replace Staging Request",
			fmt.Sprintf("Generation %d is not a valid replace-staging source for lineage %d - it must be the lineage's current or previous production generation. %s",
				generationID, lineageID, err.Error()))
	default:
		diags.AddError("Failed To Replace Staging", fmt.Sprintf("Failed to replace staging for lineage %d with generation %d: %s", lineageID, generationID, err))
	}
	return diags
}
