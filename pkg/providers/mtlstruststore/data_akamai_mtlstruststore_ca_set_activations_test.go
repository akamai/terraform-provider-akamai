package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetActivationsDataSource(t *testing.T) {
	t.Parallel()

	twoActivations := []mtlstruststore.ActivateCASetVersionResponse{
		{
			CASetName:        "example-ca-set",
			ActivationID:     1,
			Version:          1,
			Network:          "STAGING",
			CreatedBy:        "example user",
			CreatedDate:      tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
			ModifiedBy:       ptr.To("example user"),
			ModifiedDate:     ptr.To(tst.NewTimeFromString(t, "2025-05-16T12:08:34.099457Z")),
			ActivationStatus: "COMPLETE",
			ActivationType:   "ACTIVATE",
		},
		{
			CASetName:        "example-ca-set",
			ActivationID:     2,
			Version:          2,
			Network:          "PRODUCTION",
			CreatedBy:        "example user",
			CreatedDate:      tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
			ActivationStatus: "IN_PROGRESS",
			ActivationType:   "DEACTIVATE",
		},
	}

	oneActivationChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_activations.test").
		CheckEqual("ca_set_id", "12345").
		CheckEqual("ca_set_name", "example-ca-set").
		CheckMissing("version").
		CheckMissing("type").
		CheckMissing("status").
		CheckMissing("network").
		CheckEqual("activations.#", "1").
		CheckEqual("activations.0.id", "1").
		CheckEqual("activations.0.version", "1").
		CheckEqual("activations.0.network", "STAGING").
		CheckEqual("activations.0.status", "COMPLETE").
		CheckEqual("activations.0.type", "ACTIVATE").
		CheckEqual("activations.0.created_by", "example user").
		CheckEqual("activations.0.created_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("activations.0.modified_by", "example user").
		CheckEqual("activations.0.modified_date", "2025-05-16T12:08:34.099457Z")

	mockListCASetActivations := func(m *mtlstruststore.Mock, testData caSetTestData) {
		m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
			CASetID: testData.caSetID,
		}).Return(&mtlstruststore.ListCASetActivationsResponse{
			Activations: testData.caSetActivations,
		}, nil).Times(3)
	}

	mockListCASets := func(m *mtlstruststore.Mock, testData caSetTestData) {
		m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
			CASetNamePrefix: testData.caSetName,
		}).Return(&mtlstruststore.ListCASetsResponse{
			CASets: testData.caSets,
		}, nil).Times(3)
	}

	tests := map[string]struct {
		init     func(*mtlstruststore.Mock, caSetTestData)
		testData caSetTestData
		steps    []resource.TestStep
		error    *regexp.Regexp
	}{
		"happy path - fetch by id, staging and production, max and min attributes": {
			testData: caSetTestData{
				caSetActivations: twoActivations,
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/id.tf"),
					Check: oneActivationChecker.
						CheckEqual("activations.#", "2").
						CheckEqual("activations.1.id", "2").
						CheckEqual("activations.1.version", "2").
						CheckEqual("activations.1.network", "PRODUCTION").
						CheckEqual("activations.1.status", "IN_PROGRESS").
						CheckEqual("activations.1.type", "DEACTIVATE").
						CheckEqual("activations.1.created_by", "example user").
						CheckEqual("activations.1.created_date", "2025-04-16T12:08:34.099457Z").
						CheckMissing("activations.1.modified_by").
						CheckMissing("activations.1.modified_date").
						Build(),
				},
			},
		},
		"happy path - fetch by id, no activations": {
			testData: caSetTestData{
				caSetActivations: []mtlstruststore.ActivateCASetVersionResponse{},
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/id.tf"),
					Check: oneActivationChecker.
						CheckEqual("activations.#", "0").
						CheckMissing("ca_set_name").
						CheckMissing("activations.0.id").
						CheckMissing("activations.0.version").
						CheckMissing("activations.0.network").
						CheckMissing("activations.0.status").
						CheckMissing("activations.0.type").
						CheckMissing("activations.0.created_by").
						CheckMissing("activations.0.created_date").
						CheckMissing("activations.0.modified_by").
						CheckMissing("activations.0.modified_date").
						Build(),
				},
			},
		},
		"happy path - fetch by name, single activation": {
			testData: caSetTestData{
				caSets: []mtlstruststore.CASetResponse{
					{
						CASetID:   "12345",
						CASetName: "example-ca-set",
					},
				},
				caSetActivations: twoActivations[:1],
				caSetName:        "example-ca-set",
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASets(m, testData)
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/name.tf"),
					Check:  oneActivationChecker.Build(),
				},
			},
		},
		"happy path - fetch by id, filter by network": {
			testData: caSetTestData{
				caSetActivations: twoActivations,
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/network.tf"),
					Check: oneActivationChecker.
						CheckEqual("network", "STAGING").
						Build(),
				},
			},
		},
		"happy path - fetch by id, filter by status": {
			testData: caSetTestData{
				caSetActivations: twoActivations,
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/status.tf"),
					Check: oneActivationChecker.
						CheckEqual("status", "COMPLETE").
						Build(),
				},
			},
		},
		"happy path - fetch by id, filter by version": {
			testData: caSetTestData{
				caSetActivations: twoActivations,
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/version.tf"),
					Check: oneActivationChecker.
						CheckEqual("version", "1").
						Build(),
				},
			},
		},
		"happy path - fetch by id, filter by type": {
			testData: caSetTestData{
				caSetActivations: twoActivations,
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/type.tf"),
					Check: oneActivationChecker.
						CheckEqual("type", "ACTIVATE").
						Build(),
				},
			},
		},
		"happy path - fetch by id, all filters": {
			testData: caSetTestData{
				caSetActivations: twoActivations,
				caSetID:          "12345",
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/all_filters.tf"),
					Check: oneActivationChecker.
						CheckEqual("type", "ACTIVATE").
						CheckEqual("status", "COMPLETE").
						CheckEqual("network", "STAGING").
						CheckEqual("version", "1").
						Build(),
				},
			},
		},
		"expect error: List CA sets failed": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "example-ca-set",
				}).Return(nil, fmt.Errorf("List CA sets failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/name.tf"),
					ExpectError: regexp.MustCompile("List CA sets failed"),
				},
			},
		},
		"expect error: List CA set activations failed": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: "12345",
				}).Return(nil, fmt.Errorf("List CA set activations failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/id.tf"),
					ExpectError: regexp.MustCompile("List CA set activations failed"),
				},
			},
		},
		"expect error: both name and id provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/id_name.tf"),
					ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of \[ca_set_name,ca_set_id] is\nrequired`),
				},
			},
		},
		"expect error: wrong network provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/wrong_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute network value must be one of: \["STAGING" "PRODUCTION"], got: "foo"`),
				},
			},
		},
		"expect error: wrong status provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/wrong_status.tf"),
					ExpectError: regexp.MustCompile(`Attribute status value must be one of: \["IN_PROGRESS" "COMPLETE" "FAILED"],\ngot: "foo"`),
				},
			},
		},
		"expect error: wrong type provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/wrong_type.tf"),
					ExpectError: regexp.MustCompile(`Attribute type value must be one of: \["ACTIVATE" "DEACTIVATE" "DELETE"], got:\n"foo"`),
				},
			},
		},
		"expect error: too short name provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/short_name.tf"),
					ExpectError: regexp.MustCompile(`Attribute ca_set_name string length must be between 3 and 64, got: 2`),
				},
			},
		},
		"expect error: invalid name provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/invalid_name.tf"),
					ExpectError: regexp.MustCompile(`Attribute ca_set_name allowed characters are alphanumerics \(a-z, A-Z, 0-9\),\nunderscore \(_\), hyphen \(-\), percent \(%\) and period \(\.\), got: ###`),
				},
			},
		},
		"expect error: could not find ca set by name": {
			testData: caSetTestData{
				caSetName: "example-ca-set",
				caSets:    []mtlstruststore.CASetResponse{},
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: testData.caSetName,
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: testData.caSets,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivations/name.tf"),
					ExpectError: regexp.MustCompile(`no CA set found with name 'example-ca-set'`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.testData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func Test_activationFilter_matches(t *testing.T) {
	tests := map[string]struct {
		filter        activationFilter
		activation    mtlstruststore.ActivateCASetVersionResponse
		expectedMatch bool
	}{
		"no filters matches any activation": {
			filter: activationFilter{},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Network:          "STAGING",
				Version:          1,
				ActivationStatus: "COMPLETE",
				ActivationType:   "ACTIVATE",
			},
			expectedMatch: true,
		},
		"network filter matches": {
			filter: activationFilter{network: "STAGING"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Network: "STAGING",
			},
			expectedMatch: true,
		},
		"network filter does not match": {
			filter: activationFilter{network: "PRODUCTION"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Network: "STAGING",
			},
			expectedMatch: false,
		},
		"network filter case-insensitive match": {
			filter: activationFilter{network: "staging"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Network: "STAGING",
			},
			expectedMatch: true,
		},
		"status filter matches": {
			filter: activationFilter{status: "COMPLETE"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				ActivationStatus: "COMPLETE",
			},
			expectedMatch: true,
		},
		"status filter does not match": {
			filter: activationFilter{status: "FAILED"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				ActivationStatus: "COMPLETE",
			},
			expectedMatch: false,
		},
		"status filter case-insensitive match": {
			filter: activationFilter{status: "complete"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				ActivationStatus: "COMPLETE",
			},
			expectedMatch: true,
		},
		"type filter matches": {
			filter: activationFilter{activationType: "ACTIVATE"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				ActivationType: "ACTIVATE",
			},
			expectedMatch: true,
		},
		"type filter does not match": {
			filter: activationFilter{activationType: "DEACTIVATE"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				ActivationType: "ACTIVATE",
			},
			expectedMatch: false,
		},
		"type filter case-insensitive match": {
			filter: activationFilter{activationType: "activate"},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				ActivationType: "ACTIVATE",
			},
			expectedMatch: true,
		},
		"version filter matches": {
			filter: activationFilter{version: 2},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Version: 2,
			},
			expectedMatch: true,
		},
		"version filter does not match": {
			filter: activationFilter{version: 3},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Version: 2,
			},
			expectedMatch: false,
		},
		"all filters match": {
			filter: activationFilter{
				network:        "PRODUCTION",
				status:         "COMPLETE",
				activationType: "ACTIVATE",
				version:        5,
			},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Network:          "PRODUCTION",
				ActivationStatus: "COMPLETE",
				ActivationType:   "ACTIVATE",
				Version:          5,
			},
			expectedMatch: true,
		},
		"all filters, one does not match": {
			filter: activationFilter{
				network:        "PRODUCTION",
				status:         "COMPLETE",
				activationType: "ACTIVATE",
				version:        5,
			},
			activation: mtlstruststore.ActivateCASetVersionResponse{
				Network:          "PRODUCTION",
				ActivationStatus: "FAILED",
				ActivationType:   "ACTIVATE",
				Version:          5,
			},
			expectedMatch: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.filter.matches(tt.activation)
			if got != tt.expectedMatch {
				t.Errorf("activationFilter.matches() = %v, want %v", got, tt.expectedMatch)
			}
		})
	}
}
