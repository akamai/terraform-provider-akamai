package cloudcertificates

import (
	"fmt"
	"net/http"
	"regexp"
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

func TestDataBindings(t *testing.T) {
	t.Parallel()

	maxPageSize := cloudcertificates.MaxListLineageBindingsPageSize

	stagingBinding := cloudcertificates.LineageBinding{
		Active:   true,
		Hostname: "staging.example.com",
		Networks: []string{"STAGING"},
	}
	prodBinding := cloudcertificates.LineageBinding{
		Active:   true,
		Hostname: "www.example.com",
		Networks: []string{"PRODUCTION"},
	}
	bothNetworksBinding := cloudcertificates.LineageBinding{
		Active:   false,
		Hostname: "api.example.com",
		Networks: []string{"STAGING", "PRODUCTION"},
	}

	allBindingsChecker := test.NewStateChecker("data.akamai_cloudcertificates_bindings.test").
		CheckEqual("lineage_id", "500005").
		CheckEqual("bindings.#", "3").
		CheckEqual("bindings.0.active", "true").
		CheckEqual("bindings.0.hostname", "www.example.com").
		CheckEqual("bindings.0.networks.#", "1").
		CheckEqual("bindings.0.networks.0", "PRODUCTION").
		CheckEqual("bindings.1.active", "false").
		CheckEqual("bindings.1.hostname", "api.example.com").
		CheckEqual("bindings.1.networks.#", "2").
		CheckEqual("bindings.1.networks.0", "STAGING").
		CheckEqual("bindings.1.networks.1", "PRODUCTION").
		CheckEqual("bindings.2.active", "true").
		CheckEqual("bindings.2.hostname", "staging.example.com").
		CheckEqual("bindings.2.networks.#", "1").
		CheckEqual("bindings.2.networks.0", "STAGING")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - no filters, single page": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageBindings(m, "", "", "", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{prodBinding, bothNetworksBinding, stagingBinding},
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					Check: allBindingsChecker.
						CheckMissing("network").
						CheckMissing("sort_order").
						Build(),
				},
			},
		},
		"happy path - all pages collected via cursor": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageBindings(m, "", "", "", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{prodBinding, bothNetworksBinding},
					NextCursor: ptr.To("MzE5NA=="),
					TotalCount: 3,
				}).Times(3)
				mockListLineageBindings(m, "", "", "MzE5NA==", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{stagingBinding},
					TotalCount: 3,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					Check: allBindingsChecker.
						CheckMissing("network").
						CheckMissing("sort_order").
						Build(),
				},
			},
		},
		"happy path - no bindings": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageBindings(m, "", "", "", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{},
					TotalCount: 0,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_bindings.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("bindings.#", "0").
						Build(),
				},
			},
		},
		"happy path - empty page with cursor stops pagination": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageBindings(m, "", "", "", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{},
					NextCursor: ptr.To("stale-cursor"),
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_bindings.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("bindings.#", "0").
						Build(),
				},
			},
		},
		"happy path - filtered by network and sort_order": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageBindings(m, cloudcertificates.TargetNetworkStaging, cloudcertificates.SortOrderDescending, "", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{stagingBinding},
					TotalCount: 1,
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBindings/filtered.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_bindings.test").
						CheckEqual("lineage_id", "500005").
						CheckEqual("network", "STAGING").
						CheckEqual("sort_order", "DESC").
						CheckEqual("bindings.#", "1").
						CheckEqual("bindings.0.active", "true").
						CheckEqual("bindings.0.hostname", "staging.example.com").
						CheckEqual("bindings.0.networks.#", "1").
						CheckEqual("bindings.0.networks.0", "STAGING").
						Build(),
				},
			},
		},
		"lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListLineageBindings", testutils.MockContext, cloudcertificates.ListLineageBindingsRequest{
					LineageID: 500005,
					PageSize:  maxPageSize,
				}).Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					ExpectError: regexp.MustCompile("No certificate lineage found with ID 500005"),
				},
			},
		},
		"internal server error": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListLineageBindings", testutils.MockContext, cloudcertificates.ListLineageBindingsRequest{
					LineageID: 500005,
					PageSize:  maxPageSize,
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListLineageBindings, &cloudcertificates.Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891090",
				})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve lineage bindings") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
						regexp.QuoteMeta(`"status": 500`)),
				},
			},
		},
		"invalid or expired cursor": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageBindings(m, "", "", "", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{prodBinding},
					NextCursor: ptr.To("invalid-or-expired-cursor"),
					TotalCount: 2,
				}).Once()
				m.On("ListLineageBindings", testutils.MockContext, cloudcertificates.ListLineageBindingsRequest{
					LineageID: 500005,
					PageSize:  maxPageSize,
					After:     "invalid-or-expired-cursor",
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListLineageBindings, &cloudcertificates.Error{
					Type:     "/error-types/invalid-cursor",
					Title:    "Invalid or expired pagination cursor.",
					Status:   http.StatusBadRequest,
					Detail:   "The 'after' cursor value is invalid or has expired. Please restart pagination without a cursor.",
					Instance: "/error-types/invalid-cursor?traceId=1234567891091",
				})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve lineage bindings") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/invalid-cursor"`) + `.*` +
						regexp.QuoteMeta(`"title": "Invalid`) + `.*` + regexp.QuoteMeta(`or expired`) + `.*` + regexp.QuoteMeta(`pagination cursor."`) + `.*` +
						regexp.QuoteMeta(`"status": 400`)),
				},
			},
		},
		"repeated nonempty cursor": {
			init: func(m *cloudcertificates.Mock) {
				mockListLineageBindings(m, "", "", "", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{prodBinding},
					NextCursor: ptr.To("repeated-cursor"),
					TotalCount: 2,
				}).Once()
				mockListLineageBindings(m, "", "", "repeated-cursor", &cloudcertificates.ListLineageBindingsResponse{
					Bindings:   []cloudcertificates.LineageBinding{bothNetworksBinding},
					NextCursor: ptr.To("repeated-cursor"),
					TotalCount: 2,
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBindings/default.tf"),
					ExpectError: regexp.MustCompile(regexp.QuoteMeta("lineage bindings pagination returned a non-advancing cursor")),
				},
			},
		},
		"validation error - missing lineage_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBindings/missing_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "lineage_id" is required`),
				},
			},
		},
		"validation error - lineage_id lower than 1": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBindings/invalid_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute lineage_id value must be at least 1, got: 0`),
				},
			},
		},
		"validation error - invalid network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBindings/invalid_network.tf"),
					ExpectError: tst.ErrPattern(regexp.QuoteMeta(`Attribute network value must be one of: ["STAGING" "PRODUCTION"], got: "NOT_A_NETWORK"`)),
				},
			},
		},
		"validation error - invalid sort_order": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBindings/invalid_sort_order.tf"),
					ExpectError: tst.ErrPattern(regexp.QuoteMeta(`Attribute sort_order value must be one of: ["ASC" "DESC"], got: "NOT_A_SORT_ORDER"`)),
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

func mockListLineageBindings(m *cloudcertificates.Mock, network cloudcertificates.TargetNetwork, sort cloudcertificates.SortOrder, cursor string, response *cloudcertificates.ListLineageBindingsResponse) *mock.Call {
	return m.On("ListLineageBindings", testutils.MockContext, cloudcertificates.ListLineageBindingsRequest{
		LineageID: 500005,
		Network:   network,
		PageSize:  cloudcertificates.MaxListLineageBindingsPageSize,
		After:     cursor,
		Sort:      sort,
	}).Return(response, nil)
}
