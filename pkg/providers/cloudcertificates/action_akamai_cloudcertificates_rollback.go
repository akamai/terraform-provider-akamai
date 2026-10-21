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
	_ action.Action              = &RollbackAction{}
	_ action.ActionWithConfigure = &RollbackAction{}
)

// RollbackAction represents the akamai_cloudcertificates_rollback action.
type RollbackAction struct {
	meta.Action
	activationPollingConfig
}

// NewRollbackAction returns a new instance of the akamai_cloudcertificates_rollback action.
func NewRollbackAction(config activationPollingConfig) func() action.Action {
	return func() action.Action {
		return &RollbackAction{activationPollingConfig: config}
	}
}

// RollbackActionModel represents the model for akamai_cloudcertificates_rollback action configuration.
type RollbackActionModel struct {
	LineageID    types.Int64    `tfsdk:"lineage_id"`
	GenerationID types.Int64    `tfsdk:"generation_id"`
	Timeouts     timeouts.Value `tfsdk:"timeouts"`
}

// Metadata sets the full name of the akamai_cloudcertificates_rollback action.
func (a *RollbackAction) Metadata(_ context.Context, _ action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = "akamai_cloudcertificates_rollback"
}

// Schema sets the schema of the akamai_cloudcertificates_rollback action.
func (a *RollbackAction) Schema(ctx context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Rolls a lineage's production generation back to its previous production generation.",
		Attributes: map[string]schema.Attribute{
			"lineage_id": schema.Int64Attribute{
				Required:    true,
				Description: "The unique identifier of the lineage to roll back.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"generation_id": schema.Int64Attribute{
				Required: true,
				Description: "The unique identifier of the generation to roll back to. It must be the lineage's " +
					"previous production generation.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"timeouts": timeouts.AttributesWithOpts(ctx, timeouts.Opts{
				InvokeDescription: "Timeout for waiting on the rollback's activation to reach a terminal status. " +
					"Defaults to 30m. Accepts a duration string, e.g. \"45m\" or \"1h\".",
			}),
		},
	}
}

// Invoke rolls the lineage's production generation back to generation_id, then waits for the resulting
// activation to reach a terminal status.
func (a *RollbackAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RollbackActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	invokeTimeout, diags := data.Timeouts.Invoke(ctx, a.defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Cloud Certificates Rollback Action Invoke")

	lineageID := data.LineageID.ValueInt64()
	generationID := data.GenerationID.ValueInt64()
	client := a.Client.GetCloudCertificates()

	ctx, cancel := context.WithTimeout(ctx, invokeTimeout)
	defer cancel()

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Rolling back lineage %d to generation %d...", lineageID, generationID),
	})

	result, err := client.RollbackLineage(ctx, cloudcertificates.RollbackLineageRequest{
		LineageID:    lineageID,
		GenerationID: generationID,
	})
	if err != nil {
		resp.Diagnostics.Append(rollbackLineageErrorDiagnostics(err, lineageID, generationID)...)
		return
	}
	if result == nil {
		resp.Diagnostics.AddWarning("Rollback Not Needed",
			fmt.Sprintf("Generation %d is already the current production generation for lineage %d; nothing to roll back.",
				generationID, lineageID))
		return
	}

	final, diags := awaitNetworkActivation(ctx, awaitNetworkActivationOpts{
		client:       client,
		lineageID:    lineageID,
		network:      cloudcertificates.TargetNetworkProduction,
		items:        result.Items,
		pollInterval: a.pollInterval,
		onFound: func(activationID int64) {
			resp.SendProgress(action.InvokeProgressEvent{
				Message: fmt.Sprintf("Rollback submitted as activation %d for lineage %d; waiting for it to complete...",
					activationID, lineageID),
			})
		},
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Rollback of lineage %d to generation %d finished with status %s.", lineageID, generationID, final.ActivationStatus),
	})

	resp.Diagnostics.AddWarning("Terraform State Not Updated",
		fmt.Sprintf("Lineage %d's production generation was rolled back to generation %d outside of Terraform. On "+
			"the next plan, an akamai_cloudcertificates_upload resource bound to the previously current production "+
			"generation will hard-error in ModifyPlan, since that generation no longer matches the lineage's head or "+
			"current production generation - and even if it could be silently reconciled, its configured "+
			"signed_certificate_pem/trust_chain_pem are still the old generation's real certificate text, which can't "+
			"express generation %d's content. An akamai_cloudcertificates_activation resource computing "+
			"production_generation_id from the lineage's current production generation won't warn you at plan time "+
			"either: it quietly recomputes to generation %d and the plan looks like a normal reconciliation - only "+
			"apply fails, rejected since promotion requires the target to already be head and generation %d reached "+
			"production via rollback, not a fresh promotion, with a generic error giving no hint a rollback happened. "+
			"Run a full export for this lineage to re-derive its Terraform-managed resources from live state before "+
			"running another apply, rather than hand-patching the existing config.",
			lineageID, generationID, generationID, generationID, generationID))
}

// rollbackLineageErrorDiagnostics maps well-documented RollbackLineage sentinel errors to clear diagnostics.
func rollbackLineageErrorDiagnostics(err error, lineageID, generationID int64) diag.Diagnostics {
	if diags := sharedActivationErrorDiagnostics(err, lineageID, "rollback"); diags.HasError() {
		return diags
	}

	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrRollbackUnavailable):
		diags.AddError("Rollback Unavailable", fmt.Sprintf("Lineage %d has no previous production generation to roll back to.", lineageID))
	case errors.Is(err, cloudcertificates.ErrLineageBadRequest):
		diags.AddError("Invalid Rollback Request",
			fmt.Sprintf("Generation %d is not a valid rollback target for lineage %d - it must be the lineage's previous production generation. %s",
				generationID, lineageID, err.Error()))
	default:
		diags.AddError("Failed To Roll Back Lineage", fmt.Sprintf("Failed to roll back lineage %d to generation %d: %s", lineageID, generationID, err))
	}
	return diags
}
