package cloudcertificates

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestDeleteGenerationActionInvoke(t *testing.T) {
	t.Parallel()

	const (
		testLineageID    = 900001
		testGenerationID = 5000
	)

	tests := map[string]struct {
		acknowledgeRollbackCandidateRemoval *bool
		init                                func(client *cloudcertificates.Mock)
		expectedErrors                      []testutils.Diagnostic
		expectedWarnings                    []testutils.Diagnostic
		expectedProgress                    []string
	}{
		"successful deletion": {
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(nil).Once()
			},
			expectedWarnings: testutils.Diagnostics("Check for a Bound Upload Resource",
				fmt.Sprintf("Generation %d of lineage %d was deleted. If it was the lineage's head generation and an "+
					"akamai_cloudcertificates_upload resource is bound to it, that resource's generation_id no longer "+
					"matches any valid target (head or current production) for lineage %d; remove or reconcile it before "+
					"the next terraform plan, or its ModifyPlan will hard-error.", testGenerationID, testLineageID, testLineageID)),
			expectedProgress: []string{
				fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID),
				fmt.Sprintf("Generation %d of lineage %d deleted.", testGenerationID, testLineageID),
			},
		},
		"successful deletion with acknowledge_rollback_candidate_removal": {
			acknowledgeRollbackCandidateRemoval: ptr.To(true),
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{
					LineageID: testLineageID, GenerationID: testGenerationID, AcknowledgeRollbackCandidateRemoval: true,
				}).Return(nil).Once()
			},
			expectedWarnings: testutils.Diagnostics("Rollback Candidate Possibly Removed",
				fmt.Sprintf("Generation %d of lineage %d was deleted with acknowledge_rollback_candidate_removal set. "+
					"If it was the lineage's previous production generation (the rollback candidate), the "+
					"akamai_cloudcertificates_rollback action now has no eligible target for lineage %d until a new "+
					"promotion establishes a fresh one.", testGenerationID, testLineageID, testLineageID)),
			expectedProgress: []string{
				fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID),
				fmt.Sprintf("Generation %d of lineage %d deleted.", testGenerationID, testLineageID),
			},
		},
		"lineage not found": {
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(fmt.Errorf("deleting generation: %w", cloudcertificates.ErrLineageNotFound)).Once()
			},
			expectedErrors: testutils.Diagnostics("Lineage Not Found",
				fmt.Sprintf("No certificate lineage found with ID %d.", testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID)},
		},
		"generation not found": {
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(fmt.Errorf("deleting generation: %w", cloudcertificates.ErrGenerationNotFound)).Once()
			},
			expectedErrors: testutils.Diagnostics("Generation Not Found",
				fmt.Sprintf("No generation %d found for lineage %d.", testGenerationID, testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID)},
		},
		"generation deployed": {
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(fmt.Errorf("deleting generation: %w", cloudcertificates.ErrGenerationDeployed)).Once()
			},
			expectedErrors: testutils.Diagnostics("Generation Deployed",
				fmt.Sprintf("Generation %d of lineage %d is currently deployed on staging or production and cannot be deleted.",
					testGenerationID, testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID)},
		},
		"confirmation required (rollback candidate)": {
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(fmt.Errorf("deleting generation: %w", cloudcertificates.ErrGenerationDeleteConflict)).Once()
			},
			expectedErrors: testutils.Diagnostics("Confirmation Required",
				fmt.Sprintf("Generation %d of lineage %d is the lineage's previous production generation (the rollback "+
					"candidate). Set acknowledge_rollback_candidate_removal to true and invoke the action again to delete it.",
					testGenerationID, testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID)},
		},
		"last generation": {
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(fmt.Errorf("deleting generation: %w", cloudcertificates.ErrLastGeneration)).Once()
			},
			expectedErrors: testutils.Diagnostics("Last Generation",
				fmt.Sprintf("Generation %d is the only remaining generation in lineage %d; delete the lineage itself instead.",
					testGenerationID, testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID)},
		},
		"generic error falls back to default diagnostic": {
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: testLineageID, GenerationID: testGenerationID}).
					Return(fmt.Errorf("deleting generation: some unexpected error")).Once()
			},
			expectedErrors: testutils.Diagnostics("Failed To Delete Generation",
				fmt.Sprintf("Failed to delete generation %d of lineage %d: deleting generation: some unexpected error", testGenerationID, testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Deleting generation %d of lineage %d...", testGenerationID, testLineageID)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := &DeleteGenerationAction{}
			client := testutils.ConfigureTestAction(t, a)
			tc.init(client.CloudCertificates)

			resp, progress := testutils.InvokeTestAction(context.Background(), t, a, map[string]tftypes.Value{
				"lineage_id":                             tftypes.NewValue(tftypes.Number, testLineageID),
				"generation_id":                          tftypes.NewValue(tftypes.Number, testGenerationID),
				"acknowledge_rollback_candidate_removal": tftypes.NewValue(tftypes.Bool, tc.acknowledgeRollbackCandidateRemoval),
			})

			testutils.AssertActionDiagnostics(t, resp.Diagnostics, tc.expectedErrors, tc.expectedWarnings)
			testutils.AssertActionProgress(t, tc.expectedProgress, progress)
			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func TestDeleteGenerationActionConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fixture     string
		init        func(client *cloudcertificates.Mock)
		expectError *regexp.Regexp
	}{
		"minimal: lineage_id and generation_id are set": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestDeleteGenerationAction/minimal.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{LineageID: 920001, GenerationID: 6000}).
					Return(nil).Once()
			},
		},
		"full: acknowledge_rollback_candidate_removal is set": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestDeleteGenerationAction/full.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("DeleteGeneration", testutils.MockContext, cloudcertificates.DeleteGenerationRequest{
					LineageID: 920002, GenerationID: 6001, AcknowledgeRollbackCandidateRemoval: true,
				}).Return(nil).Once()
			},
		},
		"missing lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestDeleteGenerationAction/missing_lineage_id.tf"),
			expectError: tst.ErrPattern(`The argument "lineage_id" is required, but no definition was found.`),
		},
		"negative lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestDeleteGenerationAction/negative_lineage_id.tf"),
			expectError: tst.ErrPattern(`Attribute lineage_id value must be at least 1, got: -1`),
		},
		"missing generation_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestDeleteGenerationAction/missing_generation_id.tf"),
			expectError: tst.ErrPattern(`The argument "generation_id" is required, but no definition was found.`),
		},
		"negative generation_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestDeleteGenerationAction/negative_generation_id.tf"),
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

			resource.UnitTest(t, resource.TestCase{
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_14_0),
				},
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
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
