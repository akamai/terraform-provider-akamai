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

func TestRenewActionInvoke(t *testing.T) {
	t.Parallel()

	const testLineageID = 500001

	tests := map[string]struct {
		confirmAbandonHead *bool
		init               func(client *cloudcertificates.Mock)
		expectedErrors     []testutils.Diagnostic
		expectedWarnings   []testutils.Diagnostic
		expectedProgress   []string
	}{
		"successful renewal, no confirmation needed": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: testLineageID}).
					Return(&cloudcertificates.RenewLineageResponse{Head: &cloudcertificates.HeadGeneration{HeadGenerationID: 3000}}, nil).Once()
			},
			expectedWarnings: testutils.Diagnostics("Manual Follow-Up Required",
				fmt.Sprintf("Lineage %d's new head generation 3000 still needs to be signed before it can be deployed. "+
					"Retrieve its CSRs (for example via the signing_target attribute of the akamai_cloudcertificates_lineage "+
					"resource or data source), get them signed by your CA, then update akamai_cloudcertificates_upload with the "+
					"newly signed certificates and run terraform apply - the upload resource will detect that the head moved and "+
					"recreate itself for generation 3000. If akamai_cloudcertificates_activation's staging_generation_id/"+
					"production_generation_id reference the upload resource's generation_id (the recommended wiring), promotion "+
					"to staging and/or production then happens automatically on that same apply; otherwise, update those IDs "+
					"manually to 3000.", testLineageID)),
			expectedProgress: []string{
				fmt.Sprintf("Renewing lineage %d...", testLineageID),
				fmt.Sprintf("Lineage %d renewed: new head generation 3000 created with fresh CSRs for every algorithm.", testLineageID),
			},
		},
		"successful renewal with confirm_abandon_head": {
			confirmAbandonHead: ptr.To(true),
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: testLineageID, ConfirmAbandonHead: true}).
					Return(&cloudcertificates.RenewLineageResponse{Head: &cloudcertificates.HeadGeneration{HeadGenerationID: 3001}}, nil).Once()
			},
			expectedWarnings: testutils.Diagnostics("Manual Follow-Up Required",
				fmt.Sprintf("Lineage %d's new head generation 3001 still needs to be signed before it can be deployed. "+
					"Retrieve its CSRs (for example via the signing_target attribute of the akamai_cloudcertificates_lineage "+
					"resource or data source), get them signed by your CA, then update akamai_cloudcertificates_upload with the "+
					"newly signed certificates and run terraform apply - the upload resource will detect that the head moved and "+
					"recreate itself for generation 3001. If akamai_cloudcertificates_activation's staging_generation_id/"+
					"production_generation_id reference the upload resource's generation_id (the recommended wiring), promotion "+
					"to staging and/or production then happens automatically on that same apply; otherwise, update those IDs "+
					"manually to 3001. If a previous head generation existed, it was abandoned (along with any uploaded "+
					"certificates) and cannot be recovered.", testLineageID)),
			expectedProgress: []string{
				fmt.Sprintf("Renewing lineage %d...", testLineageID),
				fmt.Sprintf("Lineage %d renewed: new head generation 3001 created with fresh CSRs for every algorithm.", testLineageID),
			},
		},
		"unexpected response: no head generation returned": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: testLineageID}).
					Return(&cloudcertificates.RenewLineageResponse{}, nil).Once()
			},
			expectedErrors: testutils.Diagnostics("Unexpected Response",
				fmt.Sprintf("Lineage %d was renewed, but the response did not include a new head generation.", testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Renewing lineage %d...", testLineageID)},
		},
		"lineage not found": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: testLineageID}).
					Return(nil, fmt.Errorf("renewing lineage: %w", cloudcertificates.ErrLineageNotFound)).Once()
			},
			expectedErrors: testutils.Diagnostics("Lineage Not Found",
				fmt.Sprintf("No certificate lineage found with ID %d.", testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Renewing lineage %d...", testLineageID)},
		},
		"head deployed on staging": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: testLineageID}).
					Return(nil, fmt.Errorf("renewing lineage: %w", cloudcertificates.ErrRenewBlockedHeadOnStaging)).Once()
			},
			expectedErrors: testutils.Diagnostics("Head Deployed On Staging",
				fmt.Sprintf("Lineage %d's head generation is currently deployed on staging. Use the "+
					"akamai_cloudcertificates_replace_staging action to free it up, then renew again.", testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Renewing lineage %d...", testLineageID)},
		},
		"confirmation required": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: testLineageID}).
					Return(nil, fmt.Errorf("renewing lineage: %w", cloudcertificates.ErrConfirmationRequired)).Once()
			},
			expectedErrors: testutils.Diagnostics("Confirmation Required",
				fmt.Sprintf("Lineage %d's head generation has a still-valid, unexpired certificate. Set "+
					"confirm_abandon_head to true and invoke the action again to abandon it and renew.", testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Renewing lineage %d...", testLineageID)},
		},
		"generic error falls back to default diagnostic": {
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: testLineageID}).
					Return(nil, fmt.Errorf("renewing lineage: some unexpected error")).Once()
			},
			expectedErrors: testutils.Diagnostics("Failed To Renew Lineage",
				fmt.Sprintf("Failed to renew lineage %d: renewing lineage: some unexpected error", testLineageID)),
			expectedProgress: []string{fmt.Sprintf("Renewing lineage %d...", testLineageID)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			a := &RenewAction{}
			client := testutils.ConfigureTestAction(t, a)
			tc.init(client.CloudCertificates)

			resp, progress := testutils.InvokeTestAction(context.Background(), t, a, map[string]tftypes.Value{
				"lineage_id":           tftypes.NewValue(tftypes.Number, testLineageID),
				"confirm_abandon_head": tftypes.NewValue(tftypes.Bool, tc.confirmAbandonHead),
			})

			testutils.AssertActionDiagnostics(t, resp.Diagnostics, tc.expectedErrors, tc.expectedWarnings)
			testutils.AssertActionProgress(t, tc.expectedProgress, progress)
			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func TestRenewActionConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fixture     string
		init        func(client *cloudcertificates.Mock)
		expectError *regexp.Regexp
	}{
		"minimal: only lineage_id is set": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestRenewAction/minimal.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: 600001}).
					Return(&cloudcertificates.RenewLineageResponse{Head: &cloudcertificates.HeadGeneration{HeadGenerationID: 4000}}, nil).Once()
			},
		},
		"full: lineage_id and confirm_abandon_head are set": {
			fixture: testutils.LoadFixtureString(t, "testdata/TestRenewAction/full.tf"),
			init: func(client *cloudcertificates.Mock) {
				client.On("RenewLineage", testutils.MockContext, cloudcertificates.RenewLineageRequest{LineageID: 600002, ConfirmAbandonHead: true}).
					Return(&cloudcertificates.RenewLineageResponse{Head: &cloudcertificates.HeadGeneration{HeadGenerationID: 4001}}, nil).Once()
			},
		},
		"missing lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestRenewAction/missing_lineage_id.tf"),
			expectError: tst.ErrPattern(`The argument "lineage_id" is required, but no definition was found.`),
		},
		"negative lineage_id fails plan-time validation": {
			fixture:     testutils.LoadFixtureString(t, "testdata/TestRenewAction/negative_lineage_id.tf"),
			expectError: tst.ErrPattern(`Attribute lineage_id value must be at least 1, got: -1`),
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
