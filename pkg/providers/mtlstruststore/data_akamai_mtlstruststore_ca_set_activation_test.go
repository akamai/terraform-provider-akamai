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

func TestCASetActivationDataSource(t *testing.T) {
	t.Parallel()
	testDir := "testdata/TestDataCASetActivation/"

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
					Config: testutils.LoadFixtureString(t, testDir+"ca_set_id.tf"),
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
					Config: testutils.LoadFixtureString(t, testDir+"ca_set_id.tf"),
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
				caSetName: "test name",
				caSets: []mtlstruststore.CASetResponse{
					{
						CASetID:   "123",
						CASetName: "test name",
					},
					{
						CASetID:   "1234",
						CASetName: "test name 2",
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
					Config: testutils.LoadFixtureString(t, testDir+"ca_set_name.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"API error: ListCASets failed": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test name",
				}).Return(nil, fmt.Errorf("API error: ListCASets failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"ca_set_name.tf"),
					ExpectError: regexp.MustCompile("API error: ListCASets failed"),
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
					Config:      testutils.LoadFixtureString(t, testDir+"ca_set_id.tf"),
					ExpectError: regexp.MustCompile("API error: ListCASetActivations failed"),
				},
			},
		},
		"error: multiple ca sets with the same name": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test name",
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:   "123",
							CASetName: "test name",
						},
						{
							CASetID:   "1234",
							CASetName: "test name",
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"ca_set_name.tf"),
					ExpectError: regexp.MustCompile(`multiple CA sets IDs found with name 'test name': map\[123: 1234:]. Use the ID\nto fetch a specific CA set`),
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
					Config:      testutils.LoadFixtureString(t, testDir+"ca_set_id.tf"),
					ExpectError: regexp.MustCompile("activation with ID 321 not found"),
				},
			},
		},
		"validation error - missing required argument id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"no_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - missing one of required arguments: ca_set_id or ca_set_name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"no_ca_set_name_and_ca_set_id.tf"),
					ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[ca_set_id,ca_set_name] is\nrequired`),
				},
			},
		},
		"validation error - empty ca_set_name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"empty_name.tf"),
					ExpectError: regexp.MustCompile(`Attribute ca_set_name must not be empty or only whitespace`),
				},
			},
		},
		"validation error - too short ca_set_name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"short_name.tf"),
					ExpectError: regexp.MustCompile(`Attribute ca_set_name must not be empty or only whitespace`),
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
