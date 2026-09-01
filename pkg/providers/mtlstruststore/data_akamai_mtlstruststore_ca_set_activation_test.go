package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetActivationDataSource(t *testing.T) {
	t.Parallel()

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
			CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
		}).Return(&mtlstruststore.ListCASetsResponse{
			CASets: testData.caSets,
		}, nil).Times(3)
	}

	twoActivations := []mtlstruststore.ActivateCASetVersionResponse{
		{
			ActivationID:     321,
			CASetID:          "123",
			CASetName:        "test name",
			Version:          1,
			Network:          "PRODUCTION",
			CreatedBy:        "example user",
			CreatedDate:      tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
			ModifiedBy:       ptr.To("example user"),
			ModifiedDate:     ptr.To(tst.NewTimeFromString(t, "2025-04-16T12:10:00.000000Z")),
			ActivationStatus: "COMPLETE",
			ActivationType:   "ACTIVATE",
		},
		{
			ActivationID:     4321,
			CASetID:          "123",
			CASetName:        "test name 2",
			Version:          2,
			Network:          "STAGING",
			CreatedBy:        "example user",
			CreatedDate:      tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
			ModifiedBy:       ptr.To("example user"),
			ModifiedDate:     ptr.To(tst.NewTimeFromString(t, "2025-04-16T12:10:00.000000Z")),
			ActivationStatus: "IN_PROGRESS",
			ActivationType:   "ACTIVATE",
		},
	}

	commonStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_activation.test").
		CheckEqual("id", "321").
		CheckEqual("ca_set_id", "123").
		CheckEqual("ca_set_name", "test name").
		CheckEqual("version", "1").
		CheckEqual("network", "PRODUCTION").
		CheckEqual("created_by", "example user").
		CheckEqual("created_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("modified_by", "example user").
		CheckEqual("modified_date", "2025-04-16T12:10:00Z").
		CheckEqual("status", "COMPLETE").
		CheckEqual("type", "ACTIVATE")

	tests := map[string]struct {
		testData caSetTestData
		init     func(*mtlstruststore.Mock, caSetTestData)
		steps    []resource.TestStep
		error    *regexp.Regexp
	}{
		"happy path - ca_set_id provided": {
			testData: caSetTestData{
				caSetID:          "123",
				caSetActivations: twoActivations,
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/ca_set_id.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - ca_set_id provided - minimal data": {
			testData: caSetTestData{
				caSetID: "123",
				caSetActivations: []mtlstruststore.ActivateCASetVersionResponse{
					{
						ActivationID:     321,
						CASetID:          "123",
						CASetName:        "test name",
						Version:          1,
						Network:          "PRODUCTION",
						CreatedBy:        "example user",
						CreatedDate:      tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
						ActivationStatus: "FAILED",
						ActivationType:   "DEACTIVATE",
					},
				},
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/ca_set_id.tf"),
					Check: commonStateChecker.
						CheckEqual("status", "FAILED").
						CheckEqual("type", "DEACTIVATE").
						CheckMissing("modified_by").
						CheckMissing("modified_date").
						Build(),
				},
			},
		},
		"happy path - ca_set_name provided": {
			testData: caSetTestData{
				caSetID:   "123",
				caSetName: "test_name",
				caSets: []mtlstruststore.CASetResponse{
					{
						CASetID:   "123",
						CASetName: "test_name",
					},
					{
						CASetID:   "1234",
						CASetName: "test_name_2",
					},
				},
				caSetActivations: twoActivations,
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASets(m, testData)
				mockListCASetActivations(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/ca_set_name.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"API error: ListCASets failed": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test_name",
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(nil, fmt.Errorf("ListCASets failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/ca_set_name.tf"),
					ExpectError: regexp.MustCompile(`(?s)API.+error: ListCASets failed`),
				},
			},
		},
		"API error: ListCASetActivations failed": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: "123",
				}).Return(nil, fmt.Errorf("API error: ListCASetActivations failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/ca_set_id.tf"),
					ExpectError: regexp.MustCompile(`(?s)API error: ListCASetActivations.+failed`),
				},
			},
		},
		"error: multiple ca sets with the same name": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test_name",
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:   "123",
							CASetName: "test_name",
						},
						{
							CASetID:   "1234",
							CASetName: "test_name",
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/ca_set_name.tf"),
					ExpectError: regexp.MustCompile(`multiple CA sets found with the name 'test_name' and status 'NOT_DELETED':\s+\[123 1234\]. Use the ID\s+to fetch a specific CA set`),
				},
			},
		},
		"error: cannot find activation by id": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: "123",
				}).Return(&mtlstruststore.ListCASetActivationsResponse{
					Activations: []mtlstruststore.ActivateCASetVersionResponse{
						{
							ActivationID: 789,
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/ca_set_id.tf"),
					ExpectError: regexp.MustCompile("activation with ID 321 not found"),
				},
			},
		},
		"validation error - missing required argument id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/no_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing one of required arguments: ca_set_id or ca_set_name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/no_ca_set_name_and_ca_set_id.tf"),
					ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[ca_set_id,ca_set_name] is\nrequired`),
				},
			},
		},
		"validation error - empty ca_set_name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/empty_name.tf"),
					ExpectError: regexp.MustCompile(`Attribute ca_set_name string length must be between 3 and 64, got: 0`),
				},
			},
		},
		"validation error - too short ca_set_name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetActivation/short_name.tf"),
					ExpectError: regexp.MustCompile(`Attribute ca_set_name string length must be between 3 and 64, got: 2`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.MTLSTruststore, tc.testData)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps:                    tc.steps,
			})
			client.MTLSTruststore.AssertExpectations(t)
		})
	}
}
