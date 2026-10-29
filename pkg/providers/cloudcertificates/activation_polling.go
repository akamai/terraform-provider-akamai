package cloudcertificates

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v12/internal/retry"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// activationPollingConfig holds the default invocation timeout and poll interval shared by every CCM action
// that submits an activation and polls it to a terminal status (Rollback, ReplaceStaging). Kept distinct from
// activationResourceConfig - despite the identical shape - since these are action, not resource, settings.
type activationPollingConfig struct {
	defaultTimeout time.Duration
	pollInterval   time.Duration
}

// defaultActivationPollingConfig returns the production configuration for CCM polling actions.
func defaultActivationPollingConfig() activationPollingConfig {
	return activationPollingConfig{
		defaultTimeout: 30 * time.Minute,
		pollInterval:   15 * time.Second,
	}
}

// awaitNetworkActivationOpts aggregates the arguments to awaitNetworkActivation.
type awaitNetworkActivationOpts struct {
	client       cloudcertificates.CloudCertificates
	lineageID    int64
	network      cloudcertificates.TargetNetwork
	items        []cloudcertificates.GetActivationStatusResponse
	pollInterval time.Duration
	onFound      func(activationID int64) // called once the activation is located; may be nil
}

// awaitNetworkActivation locates the activation for opts.network in opts.items and polls it to a terminal status
// (COMPLETE, PARTIAL_SUCCESS, FAILED, or ABORTED), or until ctx's deadline elapses - set by the caller so
// multi-network callers (e.g. activationResource.promoteNetworksAndWait) can share one budget across calls.
// FAILED/ABORTED is reported as an error and PARTIAL_SUCCESS as a warning; the final response is still returned
// in all cases so callers can persist it.
func awaitNetworkActivation(ctx context.Context, opts awaitNetworkActivationOpts) (*cloudcertificates.GetActivationStatusResponse, diag.Diagnostics) {
	activation := findActivationForNetwork(opts.items, opts.network)
	if activation == nil {
		var diags diag.Diagnostics
		diags.AddError("Unexpected Response",
			fmt.Sprintf("Did not receive a %s activation for lineage %d.", strings.ToLower(string(opts.network)), opts.lineageID))
		return nil, diags
	}
	if opts.onFound != nil {
		opts.onFound(activation.ActivationID)
	}

	var diags diag.Diagnostics
	final, err := retry.Poll(ctx, retry.PollingOpts[cloudcertificates.GetActivationStatusResponse]{
		Fn: func(ctx context.Context) (*cloudcertificates.GetActivationStatusResponse, error) {
			return opts.client.GetActivationStatus(ctx, cloudcertificates.GetActivationStatusRequest{
				LineageID:    opts.lineageID,
				ActivationID: activation.ActivationID,
			})
		},
		ShouldRetryData: func(status cloudcertificates.GetActivationStatusResponse) bool {
			return !isTerminalActivationStatus(status.ActivationStatus)
		},
		Interval: opts.pollInterval,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			diags.AddError("Reached Activation Timeout",
				fmt.Sprintf("Activation %d for lineage %d did not reach a terminal status within the configured timeout.", activation.ActivationID, opts.lineageID))
			return nil, diags
		}
		diags.AddError("Failed To Read Activation Status",
			fmt.Sprintf("Failed to read activation %d for lineage %d: %s", activation.ActivationID, opts.lineageID, err.Error()))
		return nil, diags
	}

	diags.Append(activationTerminalStatusDiagnostics(*final)...)
	return final, diags
}

// sharedActivationErrorDiagnostics maps sentinel errors common to every CCM operation that triggers an
// activation (PromoteLineage, RollbackLineage, ReplaceStagingLineage) to clear diagnostics. operation names the
// specific operation for the upstream-error message only (e.g. "activation", "rollback", "replace-staging
// operation"); it doesn't affect the other mappings, which are identical across every caller. Returns empty
// diagnostics when err matches none of these sentinels, so callers can fall through to their own
// operation-specific mappings.
func sharedActivationErrorDiagnostics(err error, lineageID int64, operation string) diag.Diagnostics {
	var diags diag.Diagnostics
	switch {
	case errors.Is(err, cloudcertificates.ErrLineageNotFound):
		diags.AddError("Lineage Not Found", fmt.Sprintf("No certificate lineage found with ID %d.", lineageID))
	case errors.Is(err, cloudcertificates.ErrSingleGenerationActivationNotSupported):
		diags.AddError("Activation Not Supported", fmt.Sprintf("Lineage %d is a SINGLE_GENERATION lineage; activation operations are not supported.", lineageID))
	case errors.Is(err, cloudcertificates.ErrActivationCooldownInEffect):
		diags.AddError("Activation Cooldown In Effect", fmt.Sprintf("Lineage %d had a recent activation that is still propagating; wait before starting another.", lineageID))
	case errors.Is(err, cloudcertificates.ErrPendingActivationInProgress):
		diags.AddError("Pending Activation In Progress", fmt.Sprintf("Lineage %d already has a pending activation in progress; wait for it to complete before starting another.", lineageID))
	case errors.Is(err, cloudcertificates.ErrUpstreamActivationError):
		diags.AddError("Upstream Activation Error", fmt.Sprintf("An upstream service error occurred while processing the %s; this is typically transient and retrying may succeed.", operation))
	}
	return diags
}

// isTerminalActivationStatus reports whether status is one of the terminal activation statuses: COMPLETE,
// PARTIAL_SUCCESS, FAILED, or ABORTED.
func isTerminalActivationStatus(status string) bool {
	switch status {
	case "COMPLETE", "PARTIAL_SUCCESS", "FAILED", "ABORTED":
		return true
	default:
		return false
	}
}

// activationTerminalStatusDiagnostics reports FAILED/ABORTED activations as errors and PARTIAL_SUCCESS as a
// warning. In all three cases the activation is still a real, immutable record, so the caller should persist it
// to state regardless of the diagnostic added here.
func activationTerminalStatusDiagnostics(activation cloudcertificates.GetActivationStatusResponse) diag.Diagnostics {
	var diags diag.Diagnostics
	detail := fmt.Sprintf("Activation %d for lineage %d ended with status %s.", activation.ActivationID, activation.LineageID, activation.ActivationStatus)
	if len(activation.Errors) > 0 {
		messages := make([]string, len(activation.Errors))
		for i, e := range activation.Errors {
			messages[i] = fmt.Sprintf("%s: %s", e.Type, e.Message)
		}
		detail += fmt.Sprintf(" Errors: %s.", strings.Join(messages, "; "))
	}
	switch activation.ActivationStatus {
	case "FAILED":
		diags.AddError("Activation Failed", detail)
	case "ABORTED":
		diags.AddError("Activation Aborted", detail)
	case "PARTIAL_SUCCESS":
		diags.AddWarning("Activation Partially Succeeded", detail)
	}
	return diags
}

// findActivationForNetwork returns the item in items whose TargetEnvironment matches network, or nil if none does.
func findActivationForNetwork(items []cloudcertificates.GetActivationStatusResponse, network cloudcertificates.TargetNetwork) *cloudcertificates.GetActivationStatusResponse {
	for i := range items {
		if items[i].TargetEnvironment == string(network) {
			return &items[i]
		}
	}
	return nil
}
