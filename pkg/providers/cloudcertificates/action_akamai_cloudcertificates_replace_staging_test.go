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

// replaceStagingActivationResponse builds a ReplaceStagingLineageResponse carrying a single STAGING activation item.
func replaceStagingActivationResponse(lineageID, generationID, activationID int64, status string) *cloudcertificates.ReplaceStagingLineageResponse {
	return &cloudcertificates.ReplaceStagingLineageResponse{
		Items: []cloudcertificates.GetActivationStatusResponse{
			*activationStatus(lineageID, generationID, activationID, "STAGING", status, time.Now(), time.Now()),
		},
	}
}

// newTestReplaceStagingAction builds a ReplaceStagingAction wired to a mock CCM client via the real Configure
// flow with a microsecond polling interval, fast enough that tests don't wait out the real production default.
// replaceStagingStateNotUpdatedWarning mirrors the "Terraform State Not Updated" warning emitted after a
// successful replace-staging, for use as an expected value in test cases.
func replaceStagingStateNotUpdatedWarning(lineageID, generationID int64) testutils.Diagnostic {
	return testutils.Diagnostic{
		Summary: "Terraform State Not Updated",
		Detail: fmt.Sprintf("Lineage %d's staging generation was replaced with generation %d outside of Terraform. Unlike "+
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
			"hand-patching the existing config.", lineageID, generationID, generationID),
	}
}

// replaceStagingProgress builds the expected progress messages for invoking the replace-staging action. The
// initial "replacing staging..." message is always included; activationID > 0 adds the "submitted..." message,
// and a non-empty finishedStatus additionally adds the "finished with status..." message.
func replaceStagingProgress(lineageID, generationID, activationID int64, finishedStatus string) []string {
	progress := []string{fmt.Sprintf("Replacing staging generation with generation %d for lineage %d...", generationID, lineageID)}
	if activationID == 0 {
		return progress
	}
	progress = append(progress, fmt.Sprintf("Replace-staging submitted as activation %d for lineage %d; waiting for it to complete...", activationID, lineageID))
	if finishedStatus == "" {
		return progress
	}
	return append(progress, fmt.Sprintf("Replace-staging of lineage %d to generation %d finished with status %s.", lineageID, generationID, finishedStatus))
}

func TestReplaceStagingActionInvoke(t *testing.T) {
	t.Parallel()

	const (
		testLineageID    = 710001
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
		"successful replace staging, activation completes": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(replaceStagingActivationResponse(testLineageID, testGenerationID, 9201, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9201}).
					Return(activationStatus(testLineageID, testGenerationID, 9201, "STAGING", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
			expectedWarnings: []testutils.Diagnostic{replaceStagingStateNotUpdatedWarning(testLineageID, testGenerationID)},
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 9201, "COMPLETE"),
		},
		"replace staging not needed (204/nil response)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, nil).Once()
			},
			expectedWarnings: testutils.Diagnostics("Replace Staging Not Needed",
				fmt.Sprintf("Generation %d is already the current staging generation for lineage %d; nothing to replace.", testGenerationID, testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation starts PENDING and transitions to COMPLETE, exercising real polling": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(replaceStagingActivationResponse(testLineageID, testGenerationID, 9207, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9207}).
					Return(activationStatus(testLineageID, testGenerationID, 9207, "STAGING", "PENDING", time.Now(), time.Now()), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9207}).
					Return(activationStatus(testLineageID, testGenerationID, 9207, "STAGING", "IN_PROGRESS", time.Now(), time.Now()), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9207}).
					Return(activationStatus(testLineageID, testGenerationID, 9207, "STAGING", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
			expectedWarnings: []testutils.Diagnostic{replaceStagingStateNotUpdatedWarning(testLineageID, testGenerationID)},
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 9207, "COMPLETE"),
		},
		"unexpected response: no staging activation returned": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(&cloudcertificates.ReplaceStagingLineageResponse{
						Items: []cloudcertificates.GetActivationStatusResponse{
							*activationStatus(testLineageID, testGenerationID, 9202, "PRODUCTION", "COMPLETE", time.Now(), time.Now()),
						},
					}, nil).Once()
			},
			expectedErrors: testutils.Diagnostics("Unexpected Response",
				fmt.Sprintf("Did not receive a staging activation for lineage %d.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"lineage not found": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrLineageNotFound)).Once()
			},
			expectedErrors: testutils.Diagnostics("Lineage Not Found",
				fmt.Sprintf("No certificate lineage found with ID %d.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"no current staging": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrNoCurrentStaging)).Once()
			},
			expectedErrors: testutils.Diagnostics("No Current Staging",
				fmt.Sprintf("Lineage %d has no current staging generation to replace.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation not supported (single generation lineage)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrSingleGenerationActivationNotSupported)).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Not Supported",
				fmt.Sprintf("Lineage %d is a SINGLE_GENERATION lineage; activation operations are not supported.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation cooldown in effect": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrActivationCooldownInEffect)).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Cooldown In Effect",
				fmt.Sprintf("Lineage %d had a recent activation that is still propagating; wait before starting another.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"pending activation in progress": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrPendingActivationInProgress)).Once()
			},
			expectedErrors: testutils.Diagnostics("Pending Activation In Progress",
				fmt.Sprintf("Lineage %d already has a pending activation in progress; wait for it to complete before starting another.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"upstream activation error": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrUpstreamActivationError)).Once()
			},
			expectedErrors: testutils.Diagnostics("Upstream Activation Error",
				"An upstream service error occurred while processing the replace-staging operation; this is typically transient and retrying may succeed."),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"head not on staging (invalid replace staging request)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrReplaceStagingHeadNotOnStaging)).Once()
			},
			expectedErrors: testutils.Diagnostics("Invalid Replace Staging Request",
				fmt.Sprintf("Replacing staging for lineage %d is not valid in its current state. %s",
					testLineageID, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrReplaceStagingHeadNotOnStaging).Error())),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"invalid replace staging request (bad request)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrLineageBadRequest)).Once()
			},
			expectedErrors: testutils.Diagnostics("Invalid Replace Staging Request",
				fmt.Sprintf("Generation %d is not a valid replace-staging source for lineage %d - it must be the lineage's current or previous production generation. %s",
					testGenerationID, testLineageID, fmt.Errorf("replacing lineage staging generation: %w", cloudcertificates.ErrLineageBadRequest).Error())),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"generic error falls back to default diagnostic": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil, fmt.Errorf("replacing lineage staging generation: some unexpected error")).Once()
			},
			expectedErrors: testutils.Diagnostics("Failed To Replace Staging",
				fmt.Sprintf("Failed to replace staging for lineage %d with generation %d: replacing lineage staging generation: some unexpected error", testLineageID, testGenerationID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 0, ""),
		},
		"activation times out using a shorter override timeout": {
			invokeTimeout: ptr.To("20ms"),
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(replaceStagingActivationResponse(testLineageID, testGenerationID, 9203, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9203}).
					Return(activationStatus(testLineageID, testGenerationID, 9203, "STAGING", "PENDING", time.Now(), time.Now()), nil).Maybe()
			},
			expectedErrors: testutils.Diagnostics("Reached Activation Timeout",
				fmt.Sprintf("Activation 9203 for lineage %d did not reach a terminal status within the configured timeout.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 9203, ""),
			maxElapsed:       2 * time.Second,
		},
		"activation ends FAILED, no reconciliation warning": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(replaceStagingActivationResponse(testLineageID, testGenerationID, 9204, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9204}).
					Return(activationStatus(testLineageID, testGenerationID, 9204, "STAGING", "FAILED", time.Now(), time.Now()), nil).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Failed",
				fmt.Sprintf("Activation 9204 for lineage %d ended with status FAILED.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 9204, ""),
		},
		"activation ends ABORTED, no reconciliation warning": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(replaceStagingActivationResponse(testLineageID, testGenerationID, 9205, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9205}).
					Return(activationStatus(testLineageID, testGenerationID, 9205, "STAGING", "ABORTED", time.Now(), time.Now()), nil).Once()
			},
			expectedErrors: testutils.Diagnostics("Activation Aborted",
				fmt.Sprintf("Activation 9205 for lineage %d ended with status ABORTED.", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 9205, ""),
		},
		"GetActivationStatus read error, no reconciliation warning": {
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(replaceStagingActivationResponse(testLineageID, testGenerationID, 9206, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: testLineageID, ActivationID: 9206}).
					Return(nil, errors.New("transient read failure")).Once()
			},
			expectedErrors: testutils.Diagnostics("Failed To Read Activation Status",
				fmt.Sprintf("Failed to read activation 9206 for lineage %d: not retrying due to non-retriable error: transient read failure", testLineageID)),
			expectedProgress: replaceStagingProgress(testLineageID, testGenerationID, 9206, ""),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// a microsecond polling interval, fast enough that tests don't wait out the real production default.
			a := &ReplaceStagingAction{activationPollingConfig: activationPollingConfig{pollInterval: time.Microsecond, defaultTimeout: time.Minute}}
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

func TestReplaceStagingActionConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fixture     string
		init        func(client *cloudcertificates.Mock)
		expectError *regexp.Regexp
	}{
		"minimal: lineage_id and generation_id are set": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestReplaceStagingAction/minimal.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: 810001, GenerationID: 3100}).
					Return(replaceStagingActivationResponse(810001, 3100, 9301, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 810001, ActivationID: 9301}).
					Return(activationStatus(810001, 3100, 9301, "STAGING", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
		},
		"with_timeout: timeouts.invoke overrides the default": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestReplaceStagingAction/with_timeout.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("ReplaceStagingLineage", testutils.MockContext, cloudcertificates.ReplaceStagingLineageRequest{LineageID: 810002, GenerationID: 3101}).
					Return(replaceStagingActivationResponse(810002, 3101, 9302, "PENDING"), nil).Once()
				client.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 810002, ActivationID: 9302}).
					Return(activationStatus(810002, 3101, 9302, "STAGING", "COMPLETE", time.Now(), time.Now()), nil).Once()
			},
		},
		"missing lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestReplaceStagingAction/missing_lineage_id.tf"),
			expectError: tst.ErrPattern(`The argument "lineage_id" is required, but no definition was found.`),
		},
		"negative lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestReplaceStagingAction/negative_lineage_id.tf"),
			expectError: tst.ErrPattern(`Attribute lineage_id value must be at least 1, got: -1`),
		},
		"missing generation_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestReplaceStagingAction/missing_generation_id.tf"),
			expectError: tst.ErrPattern(`The argument "generation_id" is required, but no definition was found.`),
		},
		"negative generation_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestReplaceStagingAction/negative_generation_id.tf"),
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
			config.replaceStaging.pollInterval = time.Microsecond

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
