package cloudcertificates

import (
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/mock"
)

// TestCloudCertificatesLineageE2E exercises lineage, upload and activation (both networks) together across two
// applies, mirroring a real-world renewal cycle (DXE-7096): create a lineage, sign and upload a certificate,
// promote it to both STAGING and PRODUCTION, then apply the exact same config again and confirm the upload and
// both activations replace cleanly onto a new generation. The second apply's replace is driven entirely by the
// (mocked) API reporting that the lineage's head has moved on to a fresh generation out-of-band - the same drift
// that upload's own ModifyPlan already handles for a single-resource renewal - which is why both steps can share one
// config file. terraform_data stands in for a real CA-signing chain (e.g. tls_private_key/tls_self_signed_cert/
// tls_locally_signed_cert): only its output content matters here, not how it's produced.
func TestCloudCertificatesLineageE2E(t *testing.T) {
	t.Parallel()

	const (
		lineageID           = 500001
		initialGenerationID = 2912
		renewedGenerationID = 3050
	)

	activationDS := "akamai_cloudcertificates_activation.activation"
	uploadDS := "akamai_cloudcertificates_upload.upload"
	lineageDS := "akamai_cloudcertificates_lineage.test"

	t1 := time.Date(2026, 8, 13, 9, 56, 1, 0, time.UTC)
	t2 := time.Date(2026, 8, 13, 9, 57, 30, 0, time.UTC)
	t3 := time.Date(2026, 8, 20, 11, 10, 5, 0, time.UTC)
	t4 := time.Date(2026, 8, 20, 11, 11, 40, 0, time.UTC)

	client := edgegrid.NewTestClient()
	m := client.CloudCertificates

	// lineage: created once, never changes for the lifetime of this test.
	mockCreateLineage(m, createLineageRequestFull(), createLineageResponseFull())
	m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
		LineageID:         lineageID,
		ExpandGenerations: allExpandGenerations,
	}).Return((*cloudcertificates.GetLineageResponse)(createLineageResponseFull()), nil).Times(3)
	mockDeleteLineage(m, lineageID)

	// upload's own head-drift check (ModifyPlan): sees the original head (2912) for the first few calls (Create's
	// own lookup, plus the post-apply drift checks after step 1), then - simulating a renewal that happened
	// out-of-band between the two applies - the lineage's head has moved on to a fresh generation (3050).
	headResp := getLineageHeadOnly()
	headCalls := 0
	m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
		LineageID:         lineageID,
		ExpandGenerations: []cloudcertificates.ExpandGenerations{cloudcertificates.ExpandGenerationsHead, cloudcertificates.ExpandGenerationsCurrentProduction},
	}).Return(headResp, nil).Run(func(mock.Arguments) {
		headCalls++
		if headCalls > 3 {
			headResp.Head = getLineageNewHead().Head
		}
	}).Times(7)

	// upload: initial certificate on generation 2912.
	m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: initialGenerationID}).
		Return(getGenerationNoneSigned(), nil).Once()
	m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequest()).
		Return(uploadRSAResponse(), nil).Once()
	m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: initialGenerationID}).
		Return(getGenerationRSAOnly(), nil).Times(3)

	// upload: same certificate content re-uploaded to generation 3050, once the resource replaces onto the new head.
	m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: renewedGenerationID}).
		Return(getGenerationNoneSignedForGeneration(renewedGenerationID), nil).Once()
	m.On("UploadSignedCertificate", testutils.MockContext, uploadRSARequestForGeneration(renewedGenerationID)).
		Return(uploadRSAResponse(), nil).Once()
	m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{LineageID: lineageID, GenerationID: renewedGenerationID}).
		Return(getGenerationRSAOnlyForGeneration(renewedGenerationID), nil).Twice()

	// Activations on generation 2912 target the same generation on both networks, so the activation resource issues
	// a single combined PromoteLineage call for STAGING and PRODUCTION, avoiding the ordering issue that
	// motivated merging staging and production into one resource (DXE-7096).
	m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
		LineageID: lineageID, GenerationID: initialGenerationID,
		Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging, cloudcertificates.TargetNetworkProduction},
	}).Return(&cloudcertificates.PromoteLineageResponse{
		Items: []cloudcertificates.GetActivationStatusResponse{
			*activationStatus(lineageID, initialGenerationID, 9001, "STAGING", "PENDING", t1, t1),
			*activationStatus(lineageID, initialGenerationID, 9002, "PRODUCTION", "PENDING", t1, t1),
		},
	}, nil).Once()
	m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: lineageID, ActivationID: 9001}).
		Return(activationStatus(lineageID, initialGenerationID, 9001, "STAGING", "COMPLETE", t1, t2), nil).Times(3)
	m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: lineageID, ActivationID: 9002}).
		Return(activationStatus(lineageID, initialGenerationID, 9002, "PRODUCTION", "COMPLETE", t1, t2), nil).Times(3)

	// activations on generation 3050, once the activation resource updates in place onto the renewed generation
	// (staging_generation_id/production_generation_id no longer force replacement).
	m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
		LineageID: lineageID, GenerationID: renewedGenerationID,
		Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging, cloudcertificates.TargetNetworkProduction},
	}).Return(&cloudcertificates.PromoteLineageResponse{
		Items: []cloudcertificates.GetActivationStatusResponse{
			*activationStatus(lineageID, renewedGenerationID, 9003, "STAGING", "PENDING", t3, t3),
			*activationStatus(lineageID, renewedGenerationID, 9004, "PRODUCTION", "PENDING", t3, t3),
		},
	}, nil).Once()
	m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: lineageID, ActivationID: 9003}).
		Return(activationStatus(lineageID, renewedGenerationID, 9003, "STAGING", "COMPLETE", t3, t4), nil).Twice()
	m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: lineageID, ActivationID: 9004}).
		Return(activationStatus(lineageID, renewedGenerationID, 9004, "PRODUCTION", "COMPLETE", t3, t4), nil).Twice()

	// fast enough that this test doesn't wait out the real 15s production default between polls.
	config := defaultSubproviderConfig()
	config.activation.pollInterval = time.Microsecond

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(config)),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: testutils.LoadFixtureString(t, "testdata/TestCloudCertificatesLineageE2E/initial.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					test.NewStateChecker(lineageDS).CheckEqual("lineage_id", "500001").Build(),
					test.NewStateChecker(uploadDS).CheckEqual("generation_id", "2912").Build(),
					test.NewStateChecker(activationDS).
						CheckEqual("staging_generation_id", "2912").
						CheckEqual("staging.activation_id", "9001").
						CheckEqual("staging.activation_status", "COMPLETE").
						CheckEqual("production_generation_id", "2912").
						CheckEqual("production.activation_id", "9002").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				),
			},
			{
				Config: testutils.LoadFixtureString(t, "testdata/TestCloudCertificatesLineageE2E/initial.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(uploadDS, plancheck.ResourceActionReplace),
						plancheck.ExpectResourceAction(activationDS, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					test.NewStateChecker(uploadDS).CheckEqual("generation_id", "3050").Build(),
					test.NewStateChecker(activationDS).
						CheckEqual("staging_generation_id", "3050").
						CheckEqual("staging.activation_id", "9003").
						CheckEqual("staging.activation_status", "COMPLETE").
						CheckEqual("production_generation_id", "3050").
						CheckEqual("production.activation_id", "9004").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				),
			},
		},
	})

	client.CloudCertificates.AssertExpectations(t)
}
