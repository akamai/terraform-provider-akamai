package cloudcertificates

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/assert"
)

// rollbackActivationResponse builds a RollbackLineageResponse carrying a single PRODUCTION activation item.
func rollbackActivationResponse(lineageID, generationID, activationID int64, status string) *cloudcertificates.RollbackLineageResponse {
	return &cloudcertificates.RollbackLineageResponse{
		Items: []cloudcertificates.GetActivationStatusResponse{
			*activationStatus(lineageID, generationID, activationID, "PRODUCTION", status, time.Now(), time.Now()),
		},
	}
}

// rollbackStateNotUpdatedWarning mirrors the "Terraform State Not Updated" warning emitted after a successful
// rollback, for use as an expected value in test cases.
func rollbackStateNotUpdatedWarning(lineageID, generationID int64) testutils.Diagnostic {
	return testutils.Diagnostic{
		Summary: "Terraform State Not Updated",
		Detail: fmt.Sprintf("Lineage %d's production generation was rolled back to generation %d outside of Terraform. On "+
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
			lineageID, generationID, generationID, generationID, generationID),
	}
}

// rollbackProgress builds the expected progress messages for invoking the rollback action. The initial "rolling
// back..." message is always included; activationID > 0 adds the "submitted..." message, and a non-empty
// finishedStatus additionally adds the "finished with status..." message.
func rollbackProgress(lineageID, generationID, activationID int64, finishedStatus string) []string {
	progress := []string{fmt.Sprintf("Rolling back lineage %d to generation %d...", lineageID, generationID)}
	if activationID == 0 {
		return progress
	}
	progress = append(progress, fmt.Sprintf("Rollback submitted as activation %d for lineage %d; waiting for it to complete...", activationID, lineageID))
	if finishedStatus == "" {
		return progress
	}
	return append(progress, fmt.Sprintf("Rollback of lineage %d to generation %d finished with status %s.", lineageID, generationID, finishedStatus))
}

func TestRollbackActionInvoke(t *testing.T) {
	t.Parallel()

	const (
		testLineageID    = 700001
		testGenerationID = 2912
	)

	tests := map[string]struct {
		invokeTimeout    *string
		init             func(client *cloudcertificates.Mock)
		expectedErrors   []testutils.Diagnostic
		expectedWarnings []testutils.Diagnostic
		expectedProgress []string
		maxElapsed       time.Duration // optional: proves a shorter timeouts.invoke override actually took effect, rather than silently falling back to the 1-minute test default
	}{
		"successful rollback, activation completes": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(rollbackActivationResponse(testLineageID, testGenerationID, 9001, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9001}).
					Return(activationStatus(testLineageID, testGenerationID, 9001, "PRODUCTION", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
			expectedWarnings: []testutils.Diagnostic{rollbackStateNotUpdatedWarning(testLineageID, testGenerationID)},
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 9001, "COMPLETE"),
		},
		"rollback not needed (204/nil response)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, nil).Once()
			},
			expectedWarnings: testutils.Diagnostics("Rollback Not Needed",
				fmt.Sprintf("Generation %d is already the current production generation for lineage %d; nothing to roll back.", testGenerationID, testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation ends PARTIAL_SUCCESS, surfaced as a warning": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(rollbackActivationResponse(testLineageID, testGenerationID, 9002, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9002}).
					Return(activationStatus(testLineageID, testGenerationID, 9002, "PRODUCTION", "PARTIAL_SUCCESS", time.Now(), time.Now()), nil).Once()
			},
			expectedWarnings: []testutils.Diagnostic{
				{Summary: "Activation Partially Succeeded", Detail: fmt.Sprintf("Activation 9002 for lineage %d ended with status PARTIAL_SUCCESS.", testLineageID)},
				rollbackStateNotUpdatedWarning(testLineageID, testGenerationID),
			},
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 9002, "PARTIAL_SUCCESS"),
		},
		"activation starts PENDING and transitions to COMPLETE, exercising real polling": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(rollbackActivationResponse(testLineageID, testGenerationID, 9008, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9008}).
					Return(activationStatus(testLineageID, testGenerationID, 9008, "PRODUCTION", "PENDING", time.Now(), time.Now()), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9008}).
					Return(activationStatus(testLineageID, testGenerationID, 9008, "PRODUCTION", "IN_PROGRESS", time.Now(), time.Now()), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9008}).
					Return(activationStatus(testLineageID, testGenerationID, 9008, "PRODUCTION", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
			expectedWarnings: []testutils.Diagnostic{rollbackStateNotUpdatedWarning(testLineageID, testGenerationID)},
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 9008, "COMPLETE"),
		},
		"unexpected response: no production activation returned": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(&cloudcertificates.RollbackLineageResponse{
						Items: []cloudcertificates.GetActivationStatusResponse{
							*activationStatus(testLineageID, testGenerationID, 9003, "STAGING", "COMPLETE", time.Now(), time.Now()),
						},
					}, nil).Once()
			},
			expectedErrors: testutils.Diagnostics("Unexpected Response",
				fmt.Sprintf("Did not receive a production activation for lineage %d.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"lineage not found": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrLineageNotFound)).Once()
			},
			expectedErrors: testutils.Diagnostics("Lineage Not Found",
				fmt.Sprintf("No certificate lineage found with ID %d.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"rollback unavailable": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrRollbackUnavailable)).Once()
			},
			expectedErrors: testutils.Diagnostics("Rollback Unavailable",
				fmt.Sprintf("Lineage %d has no previous production generation to roll back to.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation not supported (single generation lineage)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrSingleGenerationActivationNotSupported)).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Not Supported",
				fmt.Sprintf("Lineage %d is a SINGLE_GENERATION lineage; activation operations are not supported.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation cooldown in effect": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrActivationCooldownInEffect)).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Cooldown In Effect",
				fmt.Sprintf("Lineage %d had a recent activation that is still propagating; wait before starting another.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"pending activation in progress": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrPendingActivationInProgress)).Once()
			},
			expectedErrors: testutils.Diagnostics("Pending Activation In Progress",
				fmt.Sprintf("Lineage %d already has a pending activation in progress; wait for it to complete before starting another.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"upstream activation error": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrUpstreamActivationError)).Once()
			},
			expectedErrors: testutils.Diagnostics("Upstream Activation Error",
				"An upstream service error occurred while processing the rollback; this is typically transient and retrying may succeed."),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"invalid rollback request (bad request)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrLineageBadRequest)).Once()
			},
			expectedErrors: testutils.Diagnostics("Invalid Rollback Request",
				fmt.Sprintf("Generation %d is not a valid rollback target for lineage %d - it must be the lineage's previous production generation. %s",
					testGenerationID, testLineageID, fmt.Errorf("rolling back lineage: %w", cloudcertificates.ErrLineageBadRequest).Error())),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"generic error falls back to default diagnostic": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("rolling back lineage: some unexpected error")).Once()
			},
			expectedErrors: testutils.Diagnostics("Failed To Roll Back Lineage",
				fmt.Sprintf("Failed to roll back lineage %d to generation %d: rolling back lineage: some unexpected error", testLineageID, testGenerationID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation times out using a shorter override timeout": {
			invokeTimeout: ptr.To("20ms"),
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(rollbackActivationResponse(testLineageID, testGenerationID, 9004, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9004}).
					Return(activationStatus(testLineageID, testGenerationID, 9004, "PRODUCTION", "PENDING", time.Now(), time.Now()), nil).Maybe()
			},
			expectedErrors: testutils.Diagnostics("Reached Activation Timeout",
				fmt.Sprintf("Activation 9004 for lineage %d did not reach a terminal status within the configured timeout.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 9004, ""),
			maxElapsed:       2 * time.Second,
		},
		"activation ends FAILED, no reconciliation warning": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(rollbackActivationResponse(testLineageID, testGenerationID, 9005, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9005}).
					Return(activationStatus(testLineageID, testGenerationID, 9005, "PRODUCTION", "FAILED", time.Now(), time.Now()), nil).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Failed",
				fmt.Sprintf("Activation 9005 for lineage %d ended with status FAILED.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 9005, ""),
		},
		"activation ends ABORTED, no reconciliation warning": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(rollbackActivationResponse(testLineageID, testGenerationID, 9006, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9006}).
					Return(activationStatus(testLineageID, testGenerationID, 9006, "PRODUCTION", "ABORTED", time.Now(), time.Now()), nil).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Aborted",
				fmt.Sprintf("Activation 9006 for lineage %d ended with status ABORTED.", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 9006, ""),
		},
		"GetActivationStatus read error, no reconciliation warning": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(rollbackActivationResponse(testLineageID, testGenerationID, 9007, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9007}).
					Return(nil, errors.New("transient read failure")).Once()
			},
			expectedErrors: testutils.Diagnostics("Failed To Read Activation Status",
				fmt.Sprintf("Failed to read activation 9007 for lineage %d: not retrying due to non-retriable error: transient read failure", testLineageID)),
			expectedProgress: rollbackProgress(testLineageID, testGenerationID, 9007, ""),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// a microsecond polling interval, fast enough that tests don't wait out the real production default.
			a := &RollbackAction{activationPollingConfig: activationPollingConfig{pollInterval: time.Microsecond, defaultTimeout: time.Minute}}
			client := testutils.ConfigureTestAction(t, a)
			tc.init(client.CloudCertificates)

			attrs := map[string]tftypes.Value{
				"lineage_id":    tftypes.NewValue(tftypes.Number, testLineageID),
				"generation_id": tftypes.NewValue(tftypes.Number, testGenerationID),
			}
			if tc.invokeTimeout != nil {
				attrs["timeouts"] = timeoutsValue(timeoutsObjectType(t, a), tc.invokeTimeout)
			}

			start := time.Now()
			resp, progress := testutils.InvokeTestAction(context.Background(), t, a, attrs)
			if tc.maxElapsed > 0 {
				assert.Less(t, time.Since(start), tc.maxElapsed,
					"invocation took too long; timeouts.invoke override may not have been applied")
			}

			testutils.AssertActionDiagnostics(t, resp.Diagnostics, tc.expectedErrors, tc.expectedWarnings)
			testutils.AssertActionProgress(t, tc.expectedProgress, progress)
			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func TestRollbackActionConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fixture     string
		init        func(client *cloudcertificates.Mock)
		expectError *regexp.Regexp
	}{
		"minimal: lineage_id and generation_id are set": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestRollbackAction/minimal.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: 800001, GenerationID: 3000}).
					Return(rollbackActivationResponse(800001, 3000, 9101, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 800001, ActivationID: 9101}).
					Return(activationStatus(800001, 3000, 9101, "PRODUCTION", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
		},
		"with_timeout: timeouts.invoke overrides the default": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestRollbackAction/with_timeout.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("RollbackLineage", testutils.MockContext, cloudcertificates.RollbackLineageRequest{LineageID: 800002, GenerationID: 3001}).
					Return(rollbackActivationResponse(800002, 3001, 9102, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 800002, ActivationID: 9102}).
					Return(activationStatus(800002, 3001, 9102, "PRODUCTION", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
		},
		"missing lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestRollbackAction/missing_lineage_id.tf"),
			expectError: tst.ErrPattern(`The argument "lineage_id" is required, but no definition was found.`),
		},
		"negative lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestRollbackAction/negative_lineage_id.tf"),
			expectError: tst.ErrPattern(`Attribute lineage_id value must be at least 1, got: -1`),
		},
		"missing generation_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestRollbackAction/missing_generation_id.tf"),
			expectError: tst.ErrPattern(`The argument "generation_id" is required, but no definition was found.`),
		},
		"negative generation_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestRollbackAction/negative_generation_id.tf"),
			expectError: tst.ErrPattern(`Attribute generation_id value must be at least 1, got: -1`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.CloudCertificates)
			}

			config := defaultSubproviderConfig()
			config.rollback.pollInterval = time.Microsecond

			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_14_0),
				},
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(config)),
				Steps: []resource.TestStep{
					{
						Config:      tc.fixture,
						ExpectError: tc.expectError,
					},
				},
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}
