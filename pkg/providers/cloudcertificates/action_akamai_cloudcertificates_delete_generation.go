package cloudcertificates

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action              = &DeleteGenerationAction{}
	_ action.ActionWithConfigure = &DeleteGenerationAction{}
)

// DeleteGenerationAction represents the akamai_cloudcertificates_delete_generation action.
type DeleteGenerationAction struct {
	meta.Action
}

// NewDeleteGenerationAction returns a new instance of the akamai_cloudcertificates_delete_generation action.
func NewDeleteGenerationAction() action.Action {
	return &DeleteGenerationAction{}
}

// DeleteGenerationActionModel represents the model for akamai_cloudcertificates_delete_generation action configuration.
type DeleteGenerationActionModel struct {
	LineageID                           types.Int64 `tfsdk:"lineage_id"`
	GenerationID                        types.Int64 `tfsdk:"generation_id"`
	AcknowledgeRollbackCandidateRemoval types.Bool  `tfsdk:"acknowledge_rollback_candidate_removal"`
}

// Metadata sets the full name of the akamai_cloudcertificates_delete_generation action.
func (a *DeleteGenerationAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_delete_generation"
}

// Schema sets the schema of the akamai_cloudcertificates_delete_generation action.
func (a *DeleteGenerationAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Permanently deletes a single certificate generation within a lineage. The generation must " +
			"not be currently deployed on staging or production, and must not be the only generation in the lineage.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the lineage the generation belongs to.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"generation_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the generation to delete.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"acknowledge_rollback_candidate_removal": schema.BoolAttribute{
				Optional: true,
				Description: "Whether to confirm deleting the generation even though it is the lineage's previous " +
					"production generation (the rollback candidate). Defaults to false.",
			},
		},
	}
}

// Invoke permanently deletes generation_id from lineage_id by calling DeleteGeneration.
func (a *DeleteGenerationAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data DeleteGenerationActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cloud Certificates Delete Generation Action Invoke")

	lineageID := data.LineageID.ValueInt64()
	generationID := data.GenerationID.ValueInt64()
	client := a.Client.GetCloudCertificates()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Deleting generation %d of lineage %d...", generationID, lineageID),
	})

	if err := client.DeleteGeneration(ctx, cloudcertificates.DeleteGenerationRequest{
		LineageID:                           lineageID,
		GenerationID:                        generationID,
		AcknowledgeRollbackCandidateRemoval: data.AcknowledgeRollbackCandidateRemoval.ValueBool(),
	}); err != nil {
		resp.Diagnostics.Append(deleteGenerationErrorDiagnostics(err, lineageID, generationID)...)
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Generation %d of lineage %d deleted.", generationID, lineageID),
	})

	if data.AcknowledgeRollbackCandidateRemoval.ValueBool() {
		resp.Diagnostics.AddWarning("Rollback Candidate Possibly Removed",
			fmt.Sprintf("Generation %d of lineage %d was deleted with acknowledge_rollback_candidate_removal set. "+
				"If it was the lineage's previous production generation (the rollback candidate), the "+
				"akamai_cloudcertificates_rollback action now has no eligible target for lineage %d until a new "+
				"promotion establishes a fresh one.", generationID, lineageID, lineageID))
		return
	}

	// A success here without the acknowledge flag means generationID could only have been head or an
	// already-ABANDONED/ARCHIVED generation (currentProduction/currentStaging/previousProduction are all
	// rejected without it). Hedge for the head case, since that's the one where an upload resource may exist.
	resp.Diagnostics.AddWarning("Check for a Bound Upload Resource",
		fmt.Sprintf("Generation %d of lineage %d was deleted. If it was the lineage's head generation and an "+
			"akamai_cloudcertificates_upload resource is bound to it, that resource's generation_id no longer "+
			"matches any valid target (head or current production) for lineage %d; remove or reconcile it before "+
			"the next terraform plan, or its ModifyPlan will hard-error.", generationID, lineageID, lineageID))
}

// deleteGenerationErrorDiagnostics maps well-documented DeleteGeneration sentinel errors to clear diagnostics.
func deleteGenerationErrorDiagnostics(err error, lineageID, generationID int64) diag.Diagnostics {
	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrLineageNotFound):
		diags.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d.", lineageID))
	case errors.Is(err, cloudcertificates.ErrGenerationNotFound):
		diags.AddError("Generation Not Found", fmt.Sprintf("No generation %d found for lineage %d.", generationID, lineageID))
	case errors.Is(err, cloudcertificates.ErrGenerationDeployed):
		diags.AddError("Generation Deployed",
			fmt.Sprintf("Generation %d of lineage %d is currently deployed on staging or production and cannot be deleted.",
				generationID, lineageID))
	case errors.Is(err, cloudcertificates.ErrGenerationDeleteConflict):
		diags.AddError("Confirmation Required",
			fmt.Sprintf("Generation %d of lineage %d is the lineage's previous production generation (the rollback "+
				"candidate). Set acknowledge_rollback_candidate_removal to true and invoke the action again to delete it.",
				generationID, lineageID))
	case errors.Is(err, cloudcertificates.ErrLastGeneration):
		diags.AddError("Last Generation",
			fmt.Sprintf("Generation %d is the only remaining generation in lineage %d; delete the lineage itself instead.",
				generationID, lineageID))
	default:
		diags.AddError("Failed To Delete Generation", fmt.Sprintf("Failed to delete generation %d of lineage %d: %s", generationID, lineageID, err))
	}
	return diags
}
