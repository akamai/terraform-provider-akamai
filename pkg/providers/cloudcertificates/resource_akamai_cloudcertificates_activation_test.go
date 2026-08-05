package cloudcertificates

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// activationStatus builds a GetActivationStatusResponse for tests, mirroring the real captured shape confirmed in
// edgegrid-golang's TestPromoteLineage fixtures.
func activationStatus(lineageID, generationID, activationID int64, network, status string, created, modified time.Time) *cloudcertificates.GetActivationStatusResponse {
	return &cloudcertificates.GetActivationStatusResponse{
		ActivationID:           activationID,
		LineageID:              lineageID,
		GenerationID:           generationID,
		ActivationType:         "PROMOTE",
		ActivationStatus:       status,
		TargetEnvironment:      network,
		ActivationCreatedTime:  created,
		ActivationModifiedTime: modified,
		CreatedBy:              "terraform-dev",
		ModifiedBy:             "terraform-dev",
	}
}

func TestActivationResource(t *testing.T) {
	t.Parallel()

	// fast enough that this test doesn't wait out the real 15s production default between polls.
	config := defaultSubproviderConfig()
	config.activation.pollInterval = time.Microsecond

	dsName := "akamai_cloudcertificates_activation.test"
	t1 := time.Date(2026, 8, 13, 9, 56, 1, 0, time.UTC)
	t2 := time.Date(2026, 8, 13, 9, 57, 30, 0, time.UTC)
	t3 := time.Date(2026, 8, 20, 11, 10, 5, 0, time.UTC)
	t4 := time.Date(2026, 8, 20, 11, 11, 40, 0, time.UTC)

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - promote to production only, immediately COMPLETE": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500035, GenerationID: 7274, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500035, 7274, 6102, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				// poll (once, immediately terminal), post-apply refresh, pre-destroy refresh
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6102}).
					Return(activationStatus(500035, 7274, 6102, "PRODUCTION", "COMPLETE", t1, t2), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("lineage_id", "500035").
						CheckEqual("production_generation_id", "7274").
						CheckMissing("staging_generation_id").
						CheckMissing("staging.generation_id").
						CheckEqual("production.generation_id", "7274").
						CheckEqual("production.activation_id", "6102").
						CheckEqual("production.activation_type", "PROMOTE").
						CheckEqual("production.activation_status", "COMPLETE").
						CheckEqual("production.activation_created_time", "2026-08-13T09:56:01Z").
						CheckEqual("production.activation_modified_time", "2026-08-13T09:57:30Z").
						CheckEqual("production.created_by", "terraform-dev").
						CheckEqual("production.modified_by", "terraform-dev").
						CheckMissing("production.total_hostname_count").
						CheckMissing("production.in_progress_hostname_count").
						CheckMissing("production.pre_empted_by").
						CheckMissing("production.error_types").
						Build(),
				},
			},
		},
		"happy path - promote to staging only, needs multiple polls before COMPLETE": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500036, GenerationID: 7273, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500036, 7273, 6100, "STAGING", "IN_PROGRESS", t1, t1)},
				}, nil).Once()
				req := cloudcertificates.GetActivationStatusRequest{LineageID: 500036, ActivationID: 6100}
				m.On("GetActivationStatus", testutils.MockContext, req).
					Return(activationStatus(500036, 7273, 6100, "STAGING", "IN_PROGRESS", t1, t1), nil).Once()
				// remaining poll, post-apply refresh, pre-destroy refresh
				m.On("GetActivationStatus", testutils.MockContext, req).
					Return(activationStatus(500036, 7273, 6100, "STAGING", "COMPLETE", t1, t2), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_staging.tf"),
					Check: test.NewStateChecker(dsName).
						CheckMissing("production_generation_id").
						CheckMissing("production.generation_id").
						CheckEqual("staging.generation_id", "7273").
						CheckEqual("staging.activation_id", "6100").
						CheckEqual("staging.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		"happy path - both networks target the same generation, single combined PromoteLineage call": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500050, GenerationID: 6000,
					Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging, cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{
						*activationStatus(500050, 6000, 8500, "STAGING", "PENDING", t1, t1),
						*activationStatus(500050, 6000, 8501, "PRODUCTION", "PENDING", t1, t1),
					},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500050, ActivationID: 8500}).
					Return(activationStatus(500050, 6000, 8500, "STAGING", "COMPLETE", t1, t2), nil)
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500050, ActivationID: 8501}).
					Return(activationStatus(500050, 6000, 8501, "PRODUCTION", "COMPLETE", t1, t2), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_same_generation.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.generation_id", "6000").
						CheckEqual("staging.activation_id", "8500").
						CheckEqual("staging.activation_status", "COMPLETE").
						CheckEqual("production.generation_id", "6000").
						CheckEqual("production.activation_id", "8501").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		"happy path - activation ends PARTIAL_SUCCESS, warning only, resource still created": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500040, GenerationID: 8000, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500040, 8000, 7000, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				resp := activationStatus(500040, 8000, 7000, "PRODUCTION", "PARTIAL_SUCCESS", t1, t2)
				resp.ErrorTypes = "/error-types/upstream-activation-error"
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500040, ActivationID: 7000}).
					Return(resp, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production_partial.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("production.activation_status", "PARTIAL_SUCCESS").
						CheckEqual("production.error_types", "/error-types/upstream-activation-error").
						Build(),
				},
			},
		},
		"expect error - activation ends FAILED, not persisted, next apply retries from scratch": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500041, GenerationID: 8001, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500041, 8001, 7001, "STAGING", "PENDING", t1, t1)},
				}, nil).Once()
				resp := activationStatus(500041, 8001, 7001, "STAGING", "FAILED", t1, t2)
				resp.ErrorTypes = "/error-types/cert-parse-error"
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500041, ActivationID: 7001}).
					Return(resp, nil).Once()
				// a failed activation is never persisted (matching property_activation's own pollActivation, which
				// never calls d.SetId on ABORTED/FAILED), so the next apply submits a brand-new PromoteLineage call
				// from scratch rather than being stuck bound to the failed one.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500041, GenerationID: 8001, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500041, 8001, 7002, "STAGING", "PENDING", t2, t2)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500041, ActivationID: 7002}).
					Return(activationStatus(500041, 8001, 7002, "STAGING", "COMPLETE", t2, t2), nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_failed.tf"),
					ExpectError: regexp.MustCompile(`Activation 7001 for lineage 500041 ended with status FAILED`),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_failed.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.activation_id", "7002").
						CheckEqual("staging.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		// COD-08/TST-01: an already-tracked network whose promotion to a NEW generation ends FAILED must keep the
		// other, untouched network's activation exactly as it was (never overwritten, never nulled out), and let
		// a later apply - even with the exact same generation_id that just failed - retry only that network via
		// ModifyPlan, without requiring the user to pick a new value.
		"expect error - production's promotion to a new generation ends FAILED on update, retried next apply, staging untouched": {
			init: func(m *cloudcertificates.Mock) {
				// initial apply: both networks target the same generation, promoted together in one call.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500091, GenerationID: 9010,
					Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging, cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{
						*activationStatus(500091, 9010, 9600, "STAGING", "PENDING", t1, t1),
						*activationStatus(500091, 9010, 9601, "PRODUCTION", "PENDING", t1, t1),
					},
				}, nil).Once()
				// staging is never touched again for the rest of this test: any later PromoteLineage call for it
				// would be an unexpected mock call and panic.
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500091, ActivationID: 9600}).
					Return(activationStatus(500091, 9010, 9600, "STAGING", "COMPLETE", t1, t1), nil)
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500091, ActivationID: 9601}).
					Return(activationStatus(500091, 9010, 9601, "PRODUCTION", "COMPLETE", t1, t1), nil)

				// update: only production_generation_id changes; its fresh promotion ends FAILED.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500091, GenerationID: 9021, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500091, 9021, 9602, "PRODUCTION", "PENDING", t2, t2)},
				}, nil).Once()
				resp := activationStatus(500091, 9021, 9602, "PRODUCTION", "FAILED", t2, t2)
				resp.ErrorTypes = "/error-types/cert-parse-error"
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500091, ActivationID: 9602}).
					Return(resp, nil).Once()

				// retry, same config as the failed update: production_generation_id already matches state - only
				// ModifyPlan noticing production's tracked generation (9010) still differs from it forces a plan at
				// all, letting Update retry with the very same generation_id that just failed, and this time it succeeds.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500091, GenerationID: 9021, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500091, 9021, 9603, "PRODUCTION", "PENDING", t3, t3)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500091, ActivationID: 9603}).
					Return(activationStatus(500091, 9021, 9603, "PRODUCTION", "COMPLETE", t3, t3), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/partial_failure_initial.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.activation_id", "9600").
						CheckEqual("production.activation_id", "9601").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/partial_failure_production_change.tf"),
					ExpectError: regexp.MustCompile(`Activation 9602 for lineage 500091 ended with status FAILED`),
				},
				{
					// production, still asking for the same 9021 that just failed, is promoted again and this time
					// succeeds - proving ModifyPlan forces the retry on its own, without the user having to pick a
					// new generation_id. Staging is never promoted again (its mock only allows the initial call).
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/partial_failure_production_change.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.generation_id", "9010").
						CheckEqual("staging.activation_id", "9600").
						CheckEqual("production.generation_id", "9021").
						CheckEqual("production.activation_id", "9603").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		// a network added for the first time via Update (not Create) whose only promotion attempt ends FAILED
		// leaves nothing at all tracked for it (unlike a network that already had a prior successful activation).
		// ModifyPlan must still force a retry on the very next apply, even with the exact same generation_id.
		"expect error - production added for the first time via update ends FAILED, retried next apply": {
			init: func(m *cloudcertificates.Mock) {
				// initial apply: staging only.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500092, GenerationID: 9030, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500092, 9030, 9700, "STAGING", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500092, ActivationID: 9700}).
					Return(activationStatus(500092, 9030, 9700, "STAGING", "COMPLETE", t1, t1), nil)

				// update: production_generation_id is set for the first time; its only promotion attempt ends FAILED.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500092, GenerationID: 9040, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500092, 9040, 9701, "PRODUCTION", "PENDING", t2, t2)},
				}, nil).Once()
				resp := activationStatus(500092, 9040, 9701, "PRODUCTION", "FAILED", t2, t2)
				resp.ErrorTypes = "/error-types/cert-parse-error"
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500092, ActivationID: 9701}).
					Return(resp, nil).Once()

				// retry, same config: nothing was ever tracked for production, so ModifyPlan must force a plan for
				// it purely from staging_generation_id/production_generation_id vs. their tracked counterparts.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500092, GenerationID: 9040, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500092, 9040, 9702, "PRODUCTION", "PENDING", t3, t3)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500092, ActivationID: 9702}).
					Return(activationStatus(500092, 9040, 9702, "PRODUCTION", "COMPLETE", t3, t3), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_only_initial.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("staging.activation_id", "9700").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_and_new_production.tf"),
					ExpectError: regexp.MustCompile(`Activation 9701 for lineage 500092 ended with status FAILED`),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_and_new_production.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("production.generation_id", "9040").
						CheckEqual("production.activation_id", "9702").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		"expect error - generation already active on network (204 no-op)": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500025, GenerationID: 4090, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(nil, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/already_active.tf"),
					ExpectError: regexp.MustCompile(`(?s)already\s+active\s+on\s+STAGING;\s+there\s+is\s+nothing\s+new\s+to\s+activate`),
				},
			},
		},
		"expect error - lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 999999, GenerationID: 4202, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/lineage_not_found.tf"),
					ExpectError: regexp.MustCompile(`No certificate lineage found with ID 999999`),
				},
			},
		},
		"expect error - lineage has no head generation": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500043, GenerationID: 4200, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(nil, cloudcertificates.ErrLineageNoHeadGeneration).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/no_head_generation.tf"),
					ExpectError: regexp.MustCompile(`Lineage 500043 has no head generation to activate`),
				},
			},
		},
		"expect error - incomplete cert material": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500026, GenerationID: 4068, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(nil, cloudcertificates.ErrIncompleteCertMaterial).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/incomplete_cert_material.tf"),
					ExpectError: regexp.MustCompile(`(?s)no\s+algorithm\s+instance\s+ready\s+for\s+use;\s+upload\s+a\s+signed\s+certificate\s+before\s+activating`),
				},
			},
		},
		// RES-03/TST-01: a read failure right after a successful promotion is a hard error and is not persisted -
		// same as a FAILED/ABORTED activation. Recovery is a fresh PromoteLineage call for the same generation on
		// the next apply; the real API, not a locally-remembered guess, is what actually arbitrates whether that's
		// safe (e.g. it would reject a genuinely still-in-flight original attempt with ErrPendingActivationInProgress).
		"expect error - first status read after a successful promotion fails, not persisted, next apply retries with a fresh promotion": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500082, GenerationID: 9700, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500082, 9700, 9800, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500082, ActivationID: 9800}).
					Return(nil, errors.New("simulated transient failure")).Once()
				// not persisted, so the next apply submits a brand-new PromoteLineage call from scratch rather than
				// being stuck bound to the unconfirmed one.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500082, GenerationID: 9700, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500082, 9700, 9801, "PRODUCTION", "PENDING", t2, t2)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500082, ActivationID: 9801}).
					Return(activationStatus(500082, 9700, 9801, "PRODUCTION", "COMPLETE", t2, t2), nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_read_fails.tf"),
					ExpectError: regexp.MustCompile(`(?s)Failed To Read Activation Status.*simulated transient failure`),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_read_fails.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("production.generation_id", "9700").
						CheckEqual("production.activation_id", "9801").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		// a polling timeout is a hard error and is not persisted - same recovery mechanics as any other unconfirmed
		// promotion outcome (see the read-failure test above).
		"expect error - activation never reaches a terminal status within the configured timeout, not persisted, next apply retries": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500042, GenerationID: 8100, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500042, 8100, 7100, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500042, ActivationID: 7100}).
					Return(activationStatus(500042, 8100, 7100, "PRODUCTION", "IN_PROGRESS", t1, t1), nil)
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500042, GenerationID: 8100, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500042, 8100, 7101, "PRODUCTION", "PENDING", t2, t2)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500042, ActivationID: 7101}).
					Return(activationStatus(500042, 8100, 7101, "PRODUCTION", "COMPLETE", t2, t2), nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_short_timeout.tf"),
					ExpectError: regexp.MustCompile(`Reached Activation Timeout`),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_short_timeout.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("production.generation_id", "8100").
						CheckEqual("production.activation_id", "7101").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		// A tracked activation reaching FAILED/ABORTED/PARTIAL_SUCCESS between applies (e.g. one
		// persisted mid-flight after Create/Update's own poll gave up, or simply drifting remotely) must not be
		// silently reported as fine on the next refresh - Read must apply the same terminal-status diagnostics
		// as the polling path.
		"expect error - tracked activation reaches FAILED on a later refresh, Read surfaces it instead of a silent success": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500035, GenerationID: 7274, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500035, 7274, 6102, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				req := cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6102}
				// poll (immediately terminal), post-apply refresh
				m.On("GetActivationStatus", testutils.MockContext, req).
					Return(activationStatus(500035, 7274, 6102, "PRODUCTION", "COMPLETE", t1, t1), nil).Twice()
				// refresh during the plan-only step: the same activation now reports FAILED.
				resp := activationStatus(500035, 7274, 6102, "PRODUCTION", "FAILED", t1, t2)
				resp.ErrorTypes = "/error-types/cert-parse-error"
				m.On("GetActivationStatus", testutils.MockContext, req).
					Return(resp, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("production.activation_id", "6102").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(`(?s)Activation Failed.*Activation 6102 for lineage 500035 ended with status FAILED`),
				},
			},
		},
		// TFP-CRIT-01: PROMOTE always targets the lineage's single current head, so staging and production can
		// never both be freshly promoted to different generations in one apply - this must be rejected outright,
		// before any API call is attempted (no mock expectations are set up below; any call would panic).
		"expect error - both networks freshly promoted to different generations is rejected before any API call": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_diff_generation.tf"),
					ExpectError: tst.ErrPattern(
						`(?s)Conflicting Promotion Targets.*staging_generation_id \(6100\) and production_generation_id \(6200\) both need a fresh promotion`,
					),
				},
			},
		},
		// caught in ModifyPlan, not just reconcileActivations: PlanOnly proves this fails at plan time, before
		// Create ever runs (no mock expectations are set up below; any call would panic).
		"expect error - both networks freshly promoted to different generations, caught at plan time on create": {
			steps: []resource.TestStep{
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_diff_generation.tf"),
					PlanOnly: true,
					ExpectError: tst.ErrPattern(
						`(?s)Conflicting Promotion Targets.*staging_generation_id \(6100\) and production_generation_id \(6200\) both need a fresh promotion`,
					),
				},
			},
		},
		// same, but on an update: staging is already tracked (frozen, no fresh promotion needed for it at its
		// current generation), then both staging and production are changed to new, different generations in the
		// same apply - PlanOnly proves this is also caught at plan time before Update ever runs.
		"expect error - both networks freshly promoted to different generations, caught at plan time on update": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500098, GenerationID: 7100, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500098, 7100, 7199, "STAGING", "COMPLETE", t1, t1)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500098, ActivationID: 7199}).
					Return(activationStatus(500098, 7100, 7199, "STAGING", "COMPLETE", t1, t1), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/conflict_update_initial.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("staging.activation_id", "7199").Build(),
				},
				{
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/conflict_update_both_new.tf"),
					PlanOnly: true,
					ExpectError: tst.ErrPattern(
						`(?s)Conflicting Promotion Targets.*staging_generation_id \(7200\) and production_generation_id \(7300\) both need a fresh promotion`,
					),
				},
			},
		},
		"expect error - neither staging_generation_id nor production_generation_id is set": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/missing_both.tf"),
					ExpectError: regexp.MustCompile(`(?s)At\s+least\s+one\s+of\s+staging_generation_id\s+or\s+production_generation_id\s+must\s+be\s+set`),
				},
			},
		},
		// COD-08/TFP-07/RES-11/TST-01: ValidateConfig defers while production_generation_id is unknown; Terraform
		// Core guarantees it's called again once this resolves (here, to null, same as staging_generation_id) -
		// no mock is set up below, so any PromoteLineage/GetActivationStatus call would panic.
		"expect error - production_generation_id resolves to null at apply time, invariant re-enforced by ValidateConfig": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
							edgerc = "../../common/testutils/edgerc"
						}

						resource "terraform_data" "production_generation_id" {
							input = null
						}

						resource "akamai_cloudcertificates_activation" "test" {
							lineage_id                = 500096
							production_generation_id  = terraform_data.production_generation_id.output
						}
					`,
					ExpectError: regexp.MustCompile(`(?s)At\s+least\s+one\s+of\s+staging_generation_id\s+or\s+production_generation_id\s+must\s+be\s+set`),
				},
			},
		},
		// RES-11/TST-01: ValidateConfig must defer its check while a generation_id is unknown (e.g. it comes from
		// another resource not yet applied) rather than treating "not yet known" the same as "genuinely null."
		"happy path - validation defers when production_generation_id is unknown, staging_generation_id is set": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
							edgerc = "../../common/testutils/edgerc"
						}

						resource "terraform_data" "production_generation_id" {
							input = 9500
						}

						resource "akamai_cloudcertificates_activation" "test" {
							lineage_id                = 500094
							staging_generation_id     = 9400
							production_generation_id  = terraform_data.production_generation_id.output
						}
					`,
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - validation defers when both generation_ids are unknown": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
							edgerc = "../../common/testutils/edgerc"
						}

						resource "terraform_data" "staging_generation_id" {
							input = 9401
						}

						resource "terraform_data" "production_generation_id" {
							input = 9501
						}

						resource "akamai_cloudcertificates_activation" "test" {
							lineage_id                = 500095
							staging_generation_id     = terraform_data.staging_generation_id.output
							production_generation_id  = terraform_data.production_generation_id.output
						}
					`,
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"happy path - a tracked activation disappears remotely, Read warns and clears it for automatic re-promotion": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500035, GenerationID: 7274, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500035, 7274, 6102, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				req := cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6102}
				// poll (immediately terminal), post-apply refresh
				m.On("GetActivationStatus", testutils.MockContext, req).
					Return(activationStatus(500035, 7274, 6102, "PRODUCTION", "COMPLETE", t1, t2), nil).Twice()
				// refresh during the plan-only step: activation was deleted outside terraform
				m.On("GetActivationStatus", testutils.MockContext, req).
					Return(nil, cloudcertificates.ErrActivationNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("production.activation_id", "6102").Build(),
				},
				{
					// production_generation_id is unchanged in config, but its tracked activation is now gone: this
					// is exactly what networkNeedsPromotion detects, so the plan shows production being re-promoted
					// rather than erroring or clearing the whole resource (staging_generation_id would show the
					// same, but this fixture only manages production).
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		// COD-08/TST-01-adjacent (RES-05/TFP-10): a network disappearing must never affect the other network's own
		// tracked activation - proving the fix scopes clearing to just the missing network, not the whole resource.
		"happy path - only production activation disappears remotely, staging is left completely untouched": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500093, GenerationID: 9800, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500093, 9800, 9900, "STAGING", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500093, GenerationID: 9801, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500093, 9801, 9901, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()

				// staging: poll (immediately terminal), post-apply refresh, pre-plan refresh before step 2, and the
				// step 2 post-apply refresh - always found, never re-promoted.
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500093, ActivationID: 9900}).
					Return(activationStatus(500093, 9800, 9900, "STAGING", "COMPLETE", t1, t1), nil)

				// production: poll (immediately terminal), post-apply refresh, then gone during step 2's pre-plan refresh.
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500093, ActivationID: 9901}).
					Return(activationStatus(500093, 9801, 9901, "PRODUCTION", "COMPLETE", t1, t1), nil).Twice()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500093, ActivationID: 9901}).
					Return(nil, cloudcertificates.ErrActivationNotFound).Once()

				// step 2's Update automatically re-promotes production to the same generation_id it was already tracking.
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500093, GenerationID: 9801, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500093, 9801, 9902, "PRODUCTION", "PENDING", t2, t2)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500093, ActivationID: 9902}).
					Return(activationStatus(500093, 9801, 9902, "PRODUCTION", "COMPLETE", t2, t2), nil)
			},
			steps: []resource.TestStep{
				{
					// staging alone first: PROMOTE always targets the lineage's single current head, so staging
					// and production can never both be freshly promoted to different generations in one apply.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_only_9800.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("staging.activation_id", "9900").Build(),
				},
				{
					// production added in a later apply, once staging_generation_id is no longer changing.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_diff_generation_read_independence.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.activation_id", "9900").
						CheckEqual("production.activation_id", "9901").
						Build(),
				},
				{
					// production_generation_id is unchanged in config: only ModifyPlan noticing production's tracked
					// activation is gone forces a plan at all, and it's scoped to production alone.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_diff_generation_read_independence.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							// staging's own activation is completely unaffected by production's disappearance.
							plancheck.ExpectKnownValue(dsName, tfjsonpath.New("staging").AtMapKey("activation_id"), knownvalue.Int64Exact(9900)),
							// production is cleared and marked for automatic re-promotion.
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("production")),
						},
					},
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.activation_id", "9900").
						CheckEqual("production.activation_id", "9902").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				},
			},
		},
		// Read must refresh staging and production independently: one network's read failing must not skip the
		// other, matching reconcileActivations' own treatment of the two networks as unrelated to each other.
		"expect error - staging read fails during refresh, production is still refreshed independently": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500093, GenerationID: 9800, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500093, 9800, 9900, "STAGING", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500093, GenerationID: 9801, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500093, 9801, 9901, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()

				// staging: poll, post-apply refresh, pre-plan/post-apply refresh around adding production, then
				// errors during the final plan-only step.
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500093, ActivationID: 9900}).
					Return(activationStatus(500093, 9800, 9900, "STAGING", "COMPLETE", t1, t1), nil).Times(4)
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500093, ActivationID: 9900}).
					Return(nil, errors.New("simulated transient failure")).Once()

				// production: refreshed the same number of times regardless - proving staging's error never skips it.
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500093, ActivationID: 9901}).
					Return(activationStatus(500093, 9801, 9901, "PRODUCTION", "COMPLETE", t1, t1), nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					// staging alone first: PROMOTE always targets the lineage's single current head, so staging
					// and production can never both be freshly promoted to different generations in one apply.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_only_9800.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("staging.activation_id", "9900").Build(),
				},
				{
					// production added in a later apply, once staging_generation_id is no longer changing.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_diff_generation_read_independence.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.activation_id", "9900").
						CheckEqual("production.activation_id", "9901").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_diff_generation_read_independence.tf"),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(`(?s)Failed To Read STAGING Activation Status.*simulated transient failure`),
				},
			},
		},
		"happy path - only timeouts change, no additional API calls": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500035, GenerationID: 7274, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500035, 7274, 6102, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6102}).
					Return(activationStatus(500035, 7274, 6102, "PRODUCTION", "COMPLETE", t1, t2), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("production.activation_id", "6102").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production_timeout_45m.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("production.activation_id", "6102").Build(),
				},
			},
		},
		"happy path - changing production_generation_id on an already-active resource promotes a fresh activation in place": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500035, GenerationID: 7274, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500035, 7274, 6102, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				req1 := cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6102}
				// poll (immediately terminal), post-apply refresh, pre-plan refresh before the update step
				m.On("GetActivationStatus", testutils.MockContext, req1).
					Return(activationStatus(500035, 7274, 6102, "PRODUCTION", "COMPLETE", t1, t2), nil)
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500035, GenerationID: 9000, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500035, 9000, 6200, "PRODUCTION", "PENDING", t2, t2)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6200}).
					Return(activationStatus(500035, 9000, 6200, "PRODUCTION", "COMPLETE", t2, t2), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("production.activation_id", "6102").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production_regen.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("production_generation_id", "9000").
						CheckEqual("production.generation_id", "9000").
						CheckEqual("production.activation_id", "6200").
						Build(),
				},
			},
		},
		"happy path - staging-only renewal freezes production, re-adding the same production value is a no-op": {
			init: func(m *cloudcertificates.Mock) {
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500051, GenerationID: 6100, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500051, 6100, 8600, "STAGING", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500051, ActivationID: 8600}).
					Return(activationStatus(500051, 6100, 8600, "STAGING", "COMPLETE", t1, t2), nil)
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500051, GenerationID: 6200, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500051, 6200, 8601, "PRODUCTION", "PENDING", t1, t1)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500051, ActivationID: 8601}).
					Return(activationStatus(500051, 6200, 8601, "PRODUCTION", "COMPLETE", t1, t2), nil)
				m.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
					LineageID: 500051, GenerationID: 6300, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkStaging},
				}).Return(&cloudcertificates.PromoteLineageResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500051, 6300, 8602, "STAGING", "PENDING", t3, t3)},
				}, nil).Once()
				m.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500051, ActivationID: 8602}).
					Return(activationStatus(500051, 6300, 8602, "STAGING", "COMPLETE", t3, t4), nil)
			},
			steps: []resource.TestStep{
				{
					// staging alone first: PROMOTE always targets the lineage's single current head, so staging
					// and production can never both be freshly promoted to different generations in one apply.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_only_6100.tf"),
					Check:  test.NewStateChecker(dsName).CheckEqual("staging.generation_id", "6100").Build(),
				},
				{
					// production added in a later apply, once staging_generation_id is no longer changing.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/both_diff_generation.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("staging.generation_id", "6100").
						CheckEqual("production.generation_id", "6200").
						Build(),
				},
				{
					// staging-only renewal: production_generation_id removed from config entirely, freezing it.
					// RES-07: production isn't changing, so its plan must stay known rather than being reset to
					// unknown; staging is being freshly promoted (6100 -> 6300), so it's expected to go unknown.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_renewal_production_omitted.tf"),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectUnknownValue(dsName, tfjsonpath.New("staging")),
							plancheck.ExpectKnownValue(dsName, tfjsonpath.New("production").AtMapKey("activation_id"), knownvalue.Int64Exact(8601)),
						},
					},
					Check: test.NewStateChecker(dsName).
						CheckMissing("production_generation_id").
						CheckEqual("staging.generation_id", "6300").
						CheckEqual("staging.activation_id", "8602").
						CheckEqual("production.generation_id", "6200").
						CheckEqual("production.activation_id", "8601").
						Build(),
				},
				{
					// re-adding the same production value that was already tracked is a no-op: no new PromoteLineage call.
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/staging_renewal_production_readded_same.tf"),
					Check: test.NewStateChecker(dsName).
						CheckEqual("production_generation_id", "6200").
						CheckEqual("staging.generation_id", "6300").
						CheckEqual("production.generation_id", "6200").
						CheckEqual("production.activation_id", "8601").
						Build(),
				},
			},
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
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(config)),
				IsUnitTest:               true,
				Steps:                    tc.steps,
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

// TestActivationResourceDefaultTimeout proves activationResourceConfig.defaultTimeout is actually wired into
// plan.Timeouts.Create's fallback: the fixture below sets no timeouts block at all, so this only times out if
// that wiring is correct - unlike every timeout case in TestActivationResource, which instead exercises the
// timeouts block's own configured value, never falling back to defaultTimeout.
func TestActivationResourceDefaultTimeout(t *testing.T) {
	t.Parallel()

	t1 := time.Date(2026, 8, 13, 9, 56, 1, 0, time.UTC)

	client := edgegrid.NewTestClient()
	client.CloudCertificates.On("PromoteLineage", testutils.MockContext, cloudcertificates.PromoteLineageRequest{
		LineageID: 500035, GenerationID: 7274, Networks: []cloudcertificates.TargetNetwork{cloudcertificates.TargetNetworkProduction},
	}).Return(&cloudcertificates.PromoteLineageResponse{
		Items: []cloudcertificates.GetActivationStatusResponse{*activationStatus(500035, 7274, 6102, "PRODUCTION", "PENDING", t1, t1)},
	}, nil).Once()
	client.CloudCertificates.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6102}).
		Return(activationStatus(500035, 7274, 6102, "PRODUCTION", "IN_PROGRESS", t1, t1), nil)

	config := defaultSubproviderConfig()
	config.activation.defaultTimeout = 20 * time.Millisecond
	config.activation.pollInterval = time.Microsecond

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(config)),
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
				ExpectError: regexp.MustCompile(`Reached Activation Timeout`),
			},
		},
	})
	client.CloudCertificates.AssertExpectations(t)
}

func TestActivationResourceImport(t *testing.T) {
	t.Parallel()

	t1 := time.Date(2026, 8, 13, 9, 56, 1, 0, time.UTC)
	t2 := time.Date(2026, 8, 13, 9, 57, 30, 0, time.UTC)

	t.Run("self-detects a single managed network from the activation's target_environment", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CloudCertificates.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500035, ActivationID: 6102}).
			Return(activationStatus(500035, 7274, 6102, "PRODUCTION", "COMPLETE", t1, t2), nil).Times(3)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_activation.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					ImportState:        true,
					ImportStateId:      "500035,6102",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("lineage_id", "500035").
						CheckEqual("production_generation_id", "7274").
						CheckMissing("staging_generation_id").
						CheckEqual("production.generation_id", "7274").
						CheckEqual("production.activation_id", "6102").
						CheckEqual("production.activation_status", "COMPLETE").
						Build(),
				},
				{
					// RES-10: the first plan after import must be empty
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					PlanOnly: true,
				},
			},
		})
		client.CloudCertificates.AssertExpectations(t)
	})

	t.Run("self-detects both managed networks, order-independent", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CloudCertificates.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500060, ActivationID: 9100}).
			Return(activationStatus(500060, 7000, 9100, "STAGING", "COMPLETE", t1, t2), nil).Times(3)
		client.CloudCertificates.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500060, ActivationID: 9101}).
			Return(activationStatus(500060, 7001, 9101, "PRODUCTION", "COMPLETE", t1, t2), nil).Times(3)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_activation.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/import_both.tf"),
					ImportState:        true,
					ImportStateId:      "500060,9101,9100",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("lineage_id", "500060").
						CheckEqual("staging_generation_id", "7000").
						CheckEqual("production_generation_id", "7001").
						CheckEqual("staging.activation_id", "9100").
						CheckEqual("production.activation_id", "9101").
						Build(),
				},
				{
					// RES-10: the first plan after import must be empty
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/import_both.tf"),
					PlanOnly: true,
				},
			},
		})
		client.CloudCertificates.AssertExpectations(t)
	})

	t.Run("errors on a duplicate network across the given activation IDs", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		client.CloudCertificates.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500060, ActivationID: 9100}).
			Return(activationStatus(500060, 7000, 9100, "STAGING", "COMPLETE", t1, t2), nil).Once()
		client.CloudCertificates.On("GetActivationStatus", testutils.MockContext, cloudcertificates.GetActivationStatusRequest{LineageID: 500060, ActivationID: 9102}).
			Return(activationStatus(500060, 7002, 9102, "STAGING", "COMPLETE", t1, t2), nil).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_activation.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/import_both.tf"),
					ImportState:   true,
					ImportStateId: "500060,9100,9102",
					ExpectError:   regexp.MustCompile(`(?s)already\s+imported\s+from\s+an\s+earlier\s+activation\s+ID`),
				},
			},
		})
		client.CloudCertificates.AssertExpectations(t)
	})

	t.Run("rejects malformed import IDs", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_activation.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					ImportState:   true,
					ImportStateId: "notanumber,6102",
					ExpectError:   regexp.MustCompile(`expected a numeric lineage_id, got: "notanumber"`),
				},
				{
					ResourceName:  "akamai_cloudcertificates_activation.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					ImportState:   true,
					ImportStateId: "500035",
					ExpectError:   regexp.MustCompile(`invalid number of importID parts: 1`),
				},
			},
		})
		client.CloudCertificates.AssertExpectations(t)
	})

	t.Run("rejects a zero or negative lineage_id", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_activation.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					ImportState:   true,
					ImportStateId: "0,6102",
					ExpectError:   regexp.MustCompile(`lineage_id must be greater than 0, got: 0`),
				},
				{
					ResourceName:  "akamai_cloudcertificates_activation.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					ImportState:   true,
					ImportStateId: "-1,6102",
					ExpectError:   regexp.MustCompile(`lineage_id must be greater than 0, got: -1`),
				},
			},
		})
		client.CloudCertificates.AssertExpectations(t)
	})

	t.Run("rejects a zero or negative activation_id", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_activation.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					ImportState:   true,
					ImportStateId: "500035,0",
					ExpectError:   regexp.MustCompile(`activation_id must be greater than 0, got: 0`),
				},
				{
					ResourceName:  "akamai_cloudcertificates_activation.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesActivation/promote_production.tf"),
					ImportState:   true,
					ImportStateId: "500035,-1",
					ExpectError:   regexp.MustCompile(`activation_id must be greater than 0, got: -1`),
				},
			},
		})
		client.CloudCertificates.AssertExpectations(t)
	})
}
