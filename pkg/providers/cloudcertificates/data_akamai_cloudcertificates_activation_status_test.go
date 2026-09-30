package cloudcertificates

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestActivationStatusDataSource(t *testing.T) {
	t.Parallel()

	activationStatusReq := cloudcertificates.GetActivationStatusRequest{
		LineageID:    500019,
		ActivationID: 100,
	}

	commonStateChecker := test.NewStateChecker("data.akamai_cloudcertificates_activation_status.test").
		CheckEqual("lineage_id", "500019").
		CheckEqual("activation_id", "100").
		CheckEqual("activation_type", "PROMOTE").
		CheckEqual("activation_status", "INIT").
		CheckEqual("generation_id", "200").
		CheckEqual("activation_created_time", "2026-07-22T12:47:55Z").
		CheckEqual("activation_modified_time", "2026-07-22T12:47:56Z").
		CheckEqual("target_environment", "STAGING").
		CheckEqual("created_by", "terraform-dev").
		CheckEqual("modified_by", "terraform-dev")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - activation status without optional fields set": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetActivationStatus", testutils.MockContext, activationStatusReq).Return(&cloudcertificates.GetActivationStatusResponse{
					ActivationID:           100,
					LineageID:              500019,
					ActivationType:         "PROMOTE",
					ActivationStatus:       "INIT",
					GenerationID:           200,
					ActivationCreatedTime:  time.Date(2026, 7, 22, 12, 47, 55, 0, time.UTC),
					ActivationModifiedTime: time.Date(2026, 7, 22, 12, 47, 56, 0, time.UTC),
					TargetEnvironment:      "STAGING",
					CreatedBy:              "terraform-dev",
					ModifiedBy:             "terraform-dev",
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/basic.tf"),
					Check: commonStateChecker.
						CheckMissing("total_hostname_count").
						CheckMissing("in_progress_hostname_count").
						CheckMissing("pre_empted_by").
						CheckMissing("error_types").
						Build(),
				},
			},
		},
		"happy path - activation status with all optional fields set": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetActivationStatus", testutils.MockContext, activationStatusReq).Return(&cloudcertificates.GetActivationStatusResponse{
					ActivationID:            100,
					LineageID:               500019,
					ActivationType:          "PROMOTE",
					ActivationStatus:        "PARTIAL_SUCCESS",
					GenerationID:            200,
					ActivationCreatedTime:   time.Date(2026, 7, 22, 12, 47, 55, 0, time.UTC),
					ActivationModifiedTime:  time.Date(2026, 7, 22, 12, 47, 56, 0, time.UTC),
					TargetEnvironment:       "STAGING",
					TotalHostnameCount:      ptr.To(int64(42)),
					InProgressHostnameCount: ptr.To(int64(5)),
					PreEmptedBy:             ptr.To(int64(1076)),
					CreatedBy:               "terraform-dev",
					ModifiedBy:              "terraform-dev",
					ErrorTypes:              "/error-types/upstream-activation-error",
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/basic.tf"),
					Check: commonStateChecker.
						CheckEqual("activation_status", "PARTIAL_SUCCESS").
						CheckEqual("total_hostname_count", "42").
						CheckEqual("in_progress_hostname_count", "5").
						CheckEqual("pre_empted_by", "1076").
						CheckEqual("error_types", "/error-types/upstream-activation-error").
						Build(),
				},
			},
		},
		"activation not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetActivationStatus", testutils.MockContext, activationStatusReq).
					Return(nil, cloudcertificates.ErrActivationNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/basic.tf"),
					ExpectError: regexp.MustCompile("No activation found with ID 100 for lineage 500019"),
				},
			},
		},
		"lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetActivationStatus", testutils.MockContext, activationStatusReq).
					Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/basic.tf"),
					ExpectError: regexp.MustCompile("No certificate lineage found with ID 500019"),
				},
			},
		},
		"internal server error": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetActivationStatus", testutils.MockContext, activationStatusReq).
					Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrGetActivationStatus, &cloudcertificates.Error{
						Type:     "/error-types/internal-error",
						Title:    "An unexpected error occurred.",
						Status:   http.StatusInternalServerError,
						Instance: "/error-types/internal-error?traceId=1234567891081",
					})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/basic.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve activation status") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
						regexp.QuoteMeta(`"status": 500`)),
				},
			},
		},
		"validation error - lineage_id missing": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/no_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "lineage_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - activation_id missing": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/no_activation_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "activation_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - lineage_id is 0": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/zero_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute lineage_id value must be at least 1, got: 0`),
				},
			},
		},
		"validation error - activation_id is 0": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivationStatus/zero_activation_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute activation_id value must be at least 1, got: 0`),
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
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps:                    tc.steps,
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}
