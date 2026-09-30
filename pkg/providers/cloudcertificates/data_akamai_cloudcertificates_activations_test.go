package cloudcertificates

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataActivations(t *testing.T) {
	t.Parallel()

	maxPageSize := 100

	promote := cloudcertificates.GetActivationStatusResponse{
		ActivationID:           1258,
		LineageID:              500022,
		ActivationType:         "PROMOTE",
		ActivationStatus:       "COMPLETE",
		GenerationID:           3032,
		ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 49, 12, 0, time.UTC),
		ActivationModifiedTime: time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
		TargetEnvironment:      "PRODUCTION",
		CreatedBy:              "terraform-dev",
		ModifiedBy:             "terraform-dev",
	}
	rollback := cloudcertificates.GetActivationStatusResponse{
		ActivationID:           1256,
		LineageID:              500022,
		ActivationType:         "ROLLBACK",
		ActivationStatus:       "INIT",
		GenerationID:           3029,
		ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 42, 29, 0, time.UTC),
		ActivationModifiedTime: time.Date(2026, 7, 23, 13, 42, 30, 0, time.UTC),
		TargetEnvironment:      "PRODUCTION",
		CreatedBy:              "terraform-dev",
		ModifiedBy:             "terraform-dev",
	}
	replaceStaging := cloudcertificates.GetActivationStatusResponse{
		ActivationID:           1254,
		LineageID:              500022,
		ActivationType:         "REPLACE_STAGING",
		ActivationStatus:       "INIT",
		GenerationID:           3028,
		ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 3, 24, 0, time.UTC),
		ActivationModifiedTime: time.Date(2026, 7, 23, 13, 3, 24, 0, time.UTC),
		TargetEnvironment:      "STAGING",
		CreatedBy:              "terraform-dev",
		ModifiedBy:             "terraform-dev",
	}

	allActivationsChecker := test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
		CheckEqual("lineage_id", "500022").
		CheckEqual("activations.#", "3").
		CheckEqual("activations.0.activation_id", "1258").
		CheckEqual("activations.0.lineage_id", "500022").
		CheckEqual("activations.0.activation_type", "PROMOTE").
		CheckEqual("activations.0.activation_status", "COMPLETE").
		CheckEqual("activations.0.generation_id", "3032").
		CheckEqual("activations.0.activation_created_time", "2026-07-23T13:49:12Z").
		CheckEqual("activations.0.activation_modified_time", "2026-07-23T13:49:13Z").
		CheckEqual("activations.0.target_environment", "PRODUCTION").
		CheckEqual("activations.0.created_by", "terraform-dev").
		CheckEqual("activations.0.modified_by", "terraform-dev").
		CheckMissing("activations.0.total_hostname_count").
		CheckMissing("activations.0.in_progress_hostname_count").
		CheckMissing("activations.0.pre_empted_by").
		CheckMissing("activations.0.error_types").
		CheckEqual("activations.1.activation_id", "1256").
		CheckEqual("activations.1.lineage_id", "500022").
		CheckEqual("activations.1.activation_type", "ROLLBACK").
		CheckEqual("activations.1.activation_status", "INIT").
		CheckEqual("activations.1.generation_id", "3029").
		CheckEqual("activations.1.activation_created_time", "2026-07-23T13:42:29Z").
		CheckEqual("activations.1.activation_modified_time", "2026-07-23T13:42:30Z").
		CheckEqual("activations.1.target_environment", "PRODUCTION").
		CheckEqual("activations.1.created_by", "terraform-dev").
		CheckEqual("activations.1.modified_by", "terraform-dev").
		CheckMissing("activations.1.total_hostname_count").
		CheckMissing("activations.1.in_progress_hostname_count").
		CheckMissing("activations.1.pre_empted_by").
		CheckMissing("activations.1.error_types").
		CheckEqual("activations.2.activation_id", "1254").
		CheckEqual("activations.2.lineage_id", "500022").
		CheckEqual("activations.2.activation_type", "REPLACE_STAGING").
		CheckEqual("activations.2.activation_status", "INIT").
		CheckEqual("activations.2.generation_id", "3028").
		CheckEqual("activations.2.activation_created_time", "2026-07-23T13:03:24Z").
		CheckEqual("activations.2.activation_modified_time", "2026-07-23T13:03:24Z").
		CheckEqual("activations.2.target_environment", "STAGING").
		CheckEqual("activations.2.created_by", "terraform-dev").
		CheckEqual("activations.2.modified_by", "terraform-dev").
		CheckMissing("activations.2.total_hostname_count").
		CheckMissing("activations.2.in_progress_hostname_count").
		CheckMissing("activations.2.pre_empted_by").
		CheckMissing("activations.2.error_types")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - no limit, single page": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote, rollback, replaceStaging},
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					Check: allActivationsChecker.
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - no limit, all pages collected via cursor": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote, rollback},
					NextCursor: ptr.To("MzE5NA=="),
					TotalCount: 3,
				}).Times(3)
				mockListActivations(m, maxPageSize, "MzE5NA==", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{replaceStaging},
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					Check: allActivationsChecker.
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - no limit, no activations": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{},
					TotalCount: 0,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
						CheckEqual("lineage_id", "500022").
						CheckEqual("activations.#", "0").
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - empty page with cursor continues pagination": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{},
					NextCursor: ptr.To("MzE5NA=="),
				}).Times(3)
				mockListActivations(m, maxPageSize, "MzE5NA==", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote, rollback, replaceStaging},
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					Check: allActivationsChecker.
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - empty next cursor stops pagination": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote, rollback, replaceStaging},
					NextCursor: ptr.To(""),
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					Check: allActivationsChecker.
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - non-advancing cursor stops pagination": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote},
					NextCursor: ptr.To("cursor-a"),
				}).Times(3)
				// The API returns the same cursor it was given, so it never advances; the second page must be
				// the last one requested, or pagination would loop forever.
				mockListActivations(m, maxPageSize, "cursor-a", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{},
					NextCursor: ptr.To("cursor-a"),
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
						CheckEqual("lineage_id", "500022").
						CheckMissing("limit").
						CheckEqual("activations.#", "1").
						CheckEqual("activations.0.activation_id", "1258").
						Build(),
				},
			},
		},
		"happy path - all optional fields set": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListActivations", testutils.MockContext, cloudcertificates.ListActivationsRequest{
					LineageID: 500022,
					PageSize:  maxPageSize,
				}).Return(&cloudcertificates.ListActivationsResponse{
					Items: []cloudcertificates.GetActivationStatusResponse{
						{
							ActivationID:            1258,
							LineageID:               500022,
							ActivationType:          "PROMOTE",
							ActivationStatus:        "PARTIAL_SUCCESS",
							GenerationID:            3032,
							ActivationCreatedTime:   time.Date(2026, 7, 23, 13, 49, 12, 0, time.UTC),
							ActivationModifiedTime:  time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
							TargetEnvironment:       "PRODUCTION",
							TotalHostnameCount:      ptr.To(int64(42)),
							InProgressHostnameCount: ptr.To(int64(5)),
							PreEmptedBy:             ptr.To(int64(1076)),
							CreatedBy:               "terraform-dev",
							ModifiedBy:              "terraform-dev",
							ErrorTypes:              "/error-types/upstream-activation-error",
						},
					},
					TotalCount: 1,
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
						CheckEqual("lineage_id", "500022").
						CheckEqual("activations.#", "1").
						CheckEqual("activations.0.activation_id", "1258").
						CheckEqual("activations.0.lineage_id", "500022").
						CheckEqual("activations.0.activation_type", "PROMOTE").
						CheckEqual("activations.0.activation_status", "PARTIAL_SUCCESS").
						CheckEqual("activations.0.generation_id", "3032").
						CheckEqual("activations.0.activation_created_time", "2026-07-23T13:49:12Z").
						CheckEqual("activations.0.activation_modified_time", "2026-07-23T13:49:13Z").
						CheckEqual("activations.0.target_environment", "PRODUCTION").
						CheckEqual("activations.0.created_by", "terraform-dev").
						CheckEqual("activations.0.modified_by", "terraform-dev").
						CheckEqual("activations.0.total_hostname_count", "42").
						CheckEqual("activations.0.in_progress_hostname_count", "5").
						CheckEqual("activations.0.pre_empted_by", "1076").
						CheckEqual("activations.0.error_types", "/error-types/upstream-activation-error").
						Build(),
				},
			},
		},
		"happy path - limit smaller than a single page": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, 2, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote, rollback, replaceStaging},
					NextCursor: ptr.To("MzE5NA=="),
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/with_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
						CheckEqual("lineage_id", "500022").
						CheckEqual("limit", "2").
						CheckEqual("activations.#", "2").
						CheckEqual("activations.0.activation_id", "1258").
						CheckEqual("activations.0.activation_type", "PROMOTE").
						CheckEqual("activations.0.generation_id", "3032").
						CheckEqual("activations.1.activation_id", "1256").
						CheckEqual("activations.1.activation_type", "ROLLBACK").
						CheckEqual("activations.1.generation_id", "3029").
						Build(),
				},
			},
		},
		"happy path - minimum limit": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, 1, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote, rollback, replaceStaging},
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/with_limit_one.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
						CheckEqual("lineage_id", "500022").
						CheckEqual("limit", "1").
						CheckEqual("activations.#", "1").
						CheckEqual("activations.0.activation_id", "1258").
						CheckEqual("activations.0.activation_type", "PROMOTE").
						Build(),
				},
			},
		},
		"happy path - limit exactly matches page size": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      generatedActivations(1, maxPageSize),
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 200,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/with_max_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
						CheckEqual("lineage_id", "500022").
						CheckEqual("limit", strconv.Itoa(maxPageSize)).
						CheckEqual("activations.#", strconv.Itoa(maxPageSize)).
						CheckEqual("activations.0.activation_id", "1").
						CheckEqual("activations.99.activation_id", strconv.Itoa(maxPageSize)).
						Build(),
				},
			},
		},
		"happy path - limit spanning multiple pages sizes final page to remaining count": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      generatedActivations(1, maxPageSize),
					NextCursor: ptr.To("MzE5NA=="),
					TotalCount: 500,
				}).Times(3)
				mockListActivations(m, 50, "MzE5NA==", &cloudcertificates.ListActivationsResponse{
					Items:      generatedActivations(maxPageSize+1, 50),
					NextCursor: ptr.To("NDE5NA=="),
					TotalCount: 500,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/with_limit_across_pages.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activations.test").
						CheckEqual("lineage_id", "500022").
						CheckEqual("limit", "150").
						CheckEqual("activations.#", "150").
						CheckEqual("activations.0.activation_id", "1").
						CheckEqual("activations.99.activation_id", strconv.Itoa(maxPageSize)).
						CheckEqual("activations.100.activation_id", strconv.Itoa(maxPageSize+1)).
						CheckEqual("activations.149.activation_id", "150").
						CheckEqual("activations.149.lineage_id", "500022").
						CheckEqual("activations.149.activation_type", "PROMOTE").
						CheckEqual("activations.149.activation_status", "COMPLETE").
						CheckEqual("activations.149.generation_id", "3032").
						CheckEqual("activations.149.activation_created_time", "2026-07-23T13:49:12Z").
						CheckEqual("activations.149.activation_modified_time", "2026-07-23T13:49:13Z").
						CheckEqual("activations.149.target_environment", "PRODUCTION").
						CheckEqual("activations.149.created_by", "terraform-dev").
						CheckEqual("activations.149.modified_by", "terraform-dev").
						CheckMissing("activations.149.total_hostname_count").
						CheckMissing("activations.149.in_progress_hostname_count").
						CheckMissing("activations.149.pre_empted_by").
						CheckMissing("activations.149.error_types").
						Build(),
				},
			},
		},
		"lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListActivations", testutils.MockContext, cloudcertificates.ListActivationsRequest{
					LineageID: 500022,
					PageSize:  maxPageSize,
				}).Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Lineage Not Found") + `.*` +
						regexp.QuoteMeta("No certificate lineage found with ID 500022")),
				},
			},
		},
		"internal server error": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListActivations", testutils.MockContext, cloudcertificates.ListActivationsRequest{
					LineageID: 500022,
					PageSize:  maxPageSize,
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListActivations, &cloudcertificates.Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891022",
				})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve lineage activations") + `.*` +
						regexp.QuoteMeta("listing activations: API error:") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
						regexp.QuoteMeta(`"title": "An unexpected error occurred."`) + `.*` +
						regexp.QuoteMeta(`"status": 500`) + `.*` +
						regexp.QuoteMeta(`"detail": ""`) + `.*` +
						regexp.QuoteMeta(`"instance": "/error-types/internal-error?traceId=1234567891022"`)),
				},
			},
		},
		"invalid or expired cursor": {
			init: func(m *cloudcertificates.Mock) {
				mockListActivations(m, maxPageSize, "", &cloudcertificates.ListActivationsResponse{
					Items:      []cloudcertificates.GetActivationStatusResponse{promote, rollback},
					NextCursor: ptr.To("invalid-or-expired-cursor"),
					TotalCount: 3,
				}).Once()
				m.On("ListActivations", testutils.MockContext, cloudcertificates.ListActivationsRequest{
					LineageID: 500022,
					PageSize:  maxPageSize,
					Cursor:    "invalid-or-expired-cursor",
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListActivations, &cloudcertificates.Error{
					Type:     "/error-types/invalid-cursor",
					Title:    "Invalid or expired pagination cursor.",
					Status:   http.StatusBadRequest,
					Detail:   "The 'cursor' value is invalid or has expired. Please restart pagination without a cursor.",
					Instance: "/error-types/invalid-cursor?traceId=1234567891027",
				})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivations/without_limit.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve lineage activations") + `.*` +
						regexp.QuoteMeta("listing activations: API error:") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/invalid-cursor"`) + `.*` +
						regexp.QuoteMeta(`"title": "Invalid or expired pagination cursor."`) + `.*` +
						regexp.QuoteMeta(`"status": 400`) + `.*` +
						regexp.QuoteMeta(`"detail": "The 'cursor' value is invalid or has expired.`) + `\s*` +
						regexp.QuoteMeta("Please restart") + `\s*` +
						regexp.QuoteMeta(`pagination without a cursor."`) + `.*` +
						regexp.QuoteMeta(`"instance": "/error-types/invalid-cursor?traceId=1234567891027"`)),
				},
			},
		},
		"validation error - missing lineage_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivations/missing_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "lineage_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - lineage_id lower than 1": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivations/invalid_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute lineage_id value must be at least 1, got: 0`),
				},
			},
		},
		"validation error - limit lower than 1": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivations/invalid_limit.tf"),
					ExpectError: regexp.MustCompile(`Attribute limit value must be at least 1, got: 0`),
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

func mockListActivations(m *cloudcertificates.Mock, pageSize int, cursor string, response *cloudcertificates.ListActivationsResponse) *mock.Call {
	return m.On("ListActivations", testutils.MockContext, cloudcertificates.ListActivationsRequest{
		LineageID: 500022,
		PageSize:  pageSize,
		Cursor:    cursor,
	}).Return(response, nil)
}

// generatedActivations returns the requested number of activation records with consecutive IDs starting at firstID.
func generatedActivations(firstID, count int) []cloudcertificates.GetActivationStatusResponse {
	activations := make([]cloudcertificates.GetActivationStatusResponse, 0, count)
	for i := range count {
		activations = append(activations, cloudcertificates.GetActivationStatusResponse{
			ActivationID:           int64(firstID + i),
			LineageID:              500022,
			ActivationType:         "PROMOTE",
			ActivationStatus:       "COMPLETE",
			GenerationID:           3032,
			ActivationCreatedTime:  time.Date(2026, 7, 23, 13, 49, 12, 0, time.UTC),
			ActivationModifiedTime: time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC),
			TargetEnvironment:      "PRODUCTION",
			CreatedBy:              "terraform-dev",
			ModifiedBy:             "terraform-dev",
		})
	}
	return activations
}
