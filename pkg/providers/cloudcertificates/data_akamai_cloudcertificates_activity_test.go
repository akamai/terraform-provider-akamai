package cloudcertificates

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataCertificatesActivity(t *testing.T) {
	t.Parallel()

	maxPageSize := 100

	lineageRenamed := cloudcertificates.LineageActivityEvent{
		ActivityID:  3242,
		EventType:   cloudcertificates.ActivityEventTypeLineageRenamed,
		LineageID:   500005,
		CreatedBy:   "terraform-dev",
		CreatedTime: tst.NewTimeFromStringMust("2026-07-21T12:18:03Z"),
	}
	certUploaded := cloudcertificates.LineageActivityEvent{
		ActivityID:   1805,
		EventType:    cloudcertificates.ActivityEventTypeCertUploaded,
		LineageID:    500005,
		GenerationID: ptr.To(int64(929)),
		CreatedBy:    "terraform-dev",
		Network:      ptr.To("STAGING"),
		Outcome:      ptr.To("ALL_SUCCESS"),
		CreatedTime:  tst.NewTimeFromStringMust("2026-07-17T12:50:44Z"),
	}
	lineageCreated := cloudcertificates.LineageActivityEvent{
		ActivityID:   1783,
		EventType:    cloudcertificates.ActivityEventTypeLineageCreated,
		LineageID:    500005,
		GenerationID: ptr.To(int64(929)),
		CreatedBy:    "terraform-dev",
		CreatedTime:  tst.NewTimeFromStringMust("2026-07-17T12:28:59Z"),
	}

	allEventsChecker := test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
		CheckEqual("lineage_id", "500005").
		CheckEqual("activities.#", "3").
		CheckEqual("activities.0.activity_id", "3242").
		CheckEqual("activities.0.event_type", cloudcertificates.ActivityEventTypeLineageRenamed).
		CheckEqual("activities.0.lineage_id", "500005").
		CheckEqual("activities.0.created_by", "terraform-dev").
		CheckEqual("activities.0.created_time", "2026-07-21T12:18:03Z").
		CheckMissing("activities.0.generation_id").
		CheckMissing("activities.0.network").
		CheckMissing("activities.0.outcome").
		CheckEqual("activities.1.activity_id", "1805").
		CheckEqual("activities.1.event_type", cloudcertificates.ActivityEventTypeCertUploaded).
		CheckEqual("activities.1.lineage_id", "500005").
		CheckEqual("activities.1.generation_id", "929").
		CheckEqual("activities.1.created_by", "terraform-dev").
		CheckEqual("activities.1.network", "STAGING").
		CheckEqual("activities.1.outcome", "ALL_SUCCESS").
		CheckEqual("activities.1.created_time", "2026-07-17T12:50:44Z").
		CheckEqual("activities.2.activity_id", "1783").
		CheckEqual("activities.2.event_type", cloudcertificates.ActivityEventTypeLineageCreated).
		CheckEqual("activities.2.lineage_id", "500005").
		CheckEqual("activities.2.generation_id", "929").
		CheckEqual("activities.2.created_by", "terraform-dev").
		CheckEqual("activities.2.created_time", "2026-07-17T12:28:59Z").
		CheckMissing("activities.2.network").
		CheckMissing("activities.2.outcome")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - no limit, single page": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, maxPageSize, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{lineageRenamed, certUploaded, lineageCreated},
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/without_limit.tf"),
					Check: allEventsChecker.
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - no limit, all pages collected via cursor": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, maxPageSize, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{lineageRenamed, certUploaded},
					NextCursor: ptr.To("MzE5NA=="),
					TotalCount: 3,
				}).Times(3)
				mockListLineageActivity(m, maxPageSize, "MzE5NA==", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{lineageCreated},
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/without_limit.tf"),
					Check: allEventsChecker.
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - no limit, no events": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, maxPageSize, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{},
					TotalCount: 0,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/without_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("activities.#", "0").
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - empty page with cursor stops pagination": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, maxPageSize, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{},
					NextCursor: ptr.To("stale-cursor"),
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/without_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("activities.#", "0").
						CheckMissing("limit").
						Build(),
				},
			},
		},
		"happy path - limit smaller than a single page": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, 2, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{lineageRenamed, certUploaded},
					NextCursor: ptr.To("MzE5NA=="),
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/with_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("limit", "2").
						CheckEqual("activities.#", "2").
						CheckEqual("activities.0.activity_id", "3242").
						CheckEqual("activities.0.event_type", cloudcertificates.ActivityEventTypeLineageRenamed).
						CheckEqual("activities.0.lineage_id", "500005").
						CheckEqual("activities.0.created_by", "terraform-dev").
						CheckEqual("activities.0.created_time", "2026-07-21T12:18:03Z").
						CheckMissing("activities.0.generation_id").
						CheckMissing("activities.0.network").
						CheckMissing("activities.0.outcome").
						CheckEqual("activities.1.activity_id", "1805").
						CheckEqual("activities.1.event_type", cloudcertificates.ActivityEventTypeCertUploaded).
						CheckEqual("activities.1.lineage_id", "500005").
						CheckEqual("activities.1.generation_id", "929").
						CheckEqual("activities.1.created_by", "terraform-dev").
						CheckEqual("activities.1.network", "STAGING").
						CheckEqual("activities.1.outcome", "ALL_SUCCESS").
						CheckEqual("activities.1.created_time", "2026-07-17T12:50:44Z").
						Build(),
				},
			},
		},
		"happy path - minimum limit": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, 1, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{lineageRenamed},
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/with_limit_one.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("limit", "1").
						CheckEqual("activities.#", "1").
						CheckEqual("activities.0.activity_id", "3242").
						CheckEqual("activities.0.event_type", cloudcertificates.ActivityEventTypeLineageRenamed).
						CheckEqual("activities.0.lineage_id", "500005").
						CheckEqual("activities.0.created_by", "terraform-dev").
						CheckEqual("activities.0.created_time", "2026-07-21T12:18:03Z").
						CheckMissing("activities.0.generation_id").
						CheckMissing("activities.0.network").
						CheckMissing("activities.0.outcome").
						Build(),
				},
			},
		},
		"happy path - limit exactly matches page size": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, maxPageSize, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     generatedEvents(1, maxPageSize),
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 200,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/with_max_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("limit", strconv.Itoa(maxPageSize)).
						CheckEqual("activities.#", strconv.Itoa(maxPageSize)).
						CheckEqual("activities.0.activity_id", "1").
						CheckEqual("activities.99.activity_id", strconv.Itoa(maxPageSize)).
						Build(),
				},
			},
		},
		"happy path - limit spanning multiple pages caps the last page size": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, maxPageSize, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     generatedEvents(1, maxPageSize),
					NextCursor: ptr.To("MzE5NA=="),
					TotalCount: 500,
				}).Times(3)
				mockListLineageActivity(m, 50, "MzE5NA==", &cloudcertificates.ListLineageActivityResponse{
					Events:     generatedEvents(maxPageSize+1, 50),
					NextCursor: ptr.To("NDE5NA=="),
					TotalCount: 500,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/with_limit_across_pages.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("limit", "150").
						CheckEqual("activities.#", "150").
						CheckEqual("activities.0.activity_id", "1").
						CheckEqual("activities.99.activity_id", strconv.Itoa(maxPageSize)).
						CheckEqual("activities.100.activity_id", strconv.Itoa(maxPageSize+1)).
						CheckEqual("activities.149.activity_id", "150").
						CheckEqual("activities.149.event_type", cloudcertificates.ActivityEventTypeLineageRenamed).
						CheckEqual("activities.149.lineage_id", "500005").
						CheckEqual("activities.149.created_by", "terraform-dev").
						CheckEqual("activities.149.created_time", "2026-07-21T12:18:03Z").
						CheckMissing("activities.149.generation_id").
						CheckMissing("activities.149.network").
						CheckMissing("activities.149.outcome").
						Build(),
				},
			},
		},
		"happy path - response larger than limit is truncated": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, 2, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{lineageRenamed, certUploaded, lineageCreated},
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/with_limit.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_activity.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("limit", "2").
						CheckEqual("activities.#", "2").
						CheckEqual("activities.0.activity_id", "3242").
						CheckEqual("activities.1.activity_id", "1805").
						Build(),
				},
			},
		},
		"lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListLineageActivity", testutils.MockContext, cloudcertificates.ListLineageActivityRequest{
					LineageID: 500005,
					PageSize:  maxPageSize,
				}).Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivity/without_limit.tf"),
					ExpectError: regexp.MustCompile("No certificate lineage found with ID 500005"),
				},
			},
		},
		"internal server error": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListLineageActivity", testutils.MockContext, cloudcertificates.ListLineageActivityRequest{
					LineageID: 500005,
					PageSize:  maxPageSize,
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListLineageActivity, &cloudcertificates.Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891022",
				})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/without_limit.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve lineage activity") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
						regexp.QuoteMeta(`"status": 500`)),
				},
			},
		},
		"invalid or expired cursor": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageActivity(m, maxPageSize, "", &cloudcertificates.ListLineageActivityResponse{
					Events:     []cloudcertificates.LineageActivityEvent{lineageRenamed, certUploaded},
					NextCursor: ptr.To("invalid-or-expired-cursor"),
					TotalCount: 3,
				}).Once()
				m.On("ListLineageActivity", testutils.MockContext, cloudcertificates.ListLineageActivityRequest{
					LineageID: 500005,
					PageSize:  maxPageSize,
					Cursor:    "invalid-or-expired-cursor",
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListLineageActivity, &cloudcertificates.Error{
					Type:     "/error-types/invalid-cursor",
					Title:    "Invalid or expired pagination cursor.",
					Status:   http.StatusBadRequest,
					Detail:   "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
					Instance: "/error-types/invalid-cursor?traceId=1234567891027",
				})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataActivity/without_limit.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve lineage activity") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/invalid-cursor"`) + `.*` +
						regexp.QuoteMeta(`"title": "Invalid`) + `.*` + regexp.QuoteMeta(`or expired`) + `.*` + regexp.QuoteMeta(`pagination cursor."`) + `.*` +
						regexp.QuoteMeta(`"status": 400`)),
				},
			},
		},
		"validation error - missing lineage_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivity/missing_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "lineage_id" is required`),
				},
			},
		},
		"validation error - lineage_id lower than 1": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivity/invalid_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute lineage_id value must be at least 1, got: 0`),
				},
			},
		},
		"validation error - limit lower than 1": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataActivity/invalid_limit.tf"),
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

func mockListLineageActivity(m *cloudcertificates.Mock, pageSize int, cursor string, response *cloudcertificates.ListLineageActivityResponse) *mock.Call {
	return m.On("ListLineageActivity", testutils.MockContext, cloudcertificates.ListLineageActivityRequest{
		LineageID: 500005,
		PageSize:  pageSize,
		Cursor:    cursor,
	}).Return(response, nil)
}

// generatedEvents returns count lineage-level events with consecutive activity IDs starting at firstID.
func generatedEvents(firstID, count int) []cloudcertificates.LineageActivityEvent {
	events := make([]cloudcertificates.LineageActivityEvent, 0, count)
	for i := range count {
		events = append(events, cloudcertificates.LineageActivityEvent{
			ActivityID:  int64(firstID + i),
			EventType:   cloudcertificates.ActivityEventTypeLineageRenamed,
			LineageID:   500005,
			CreatedBy:   "terraform-dev",
			CreatedTime: tst.NewTimeFromStringMust("2026-07-21T12:18:03Z"),
		})
	}
	return events
}
