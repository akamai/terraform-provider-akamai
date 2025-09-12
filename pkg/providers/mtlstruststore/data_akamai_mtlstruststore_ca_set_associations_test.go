package mtlstruststore

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestCASetAssociationsDataSource(t *testing.T) {
	t.Parallel()
	testDir := "testdata/TestDataCASetAssociations/"

	mockGetCASet := func(client *mtlstruststore.Mock) *mock.Call {
		return client.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
			CASetID: "123",
		}).Return(&mtlstruststore.GetCASetResponse{
			CASetID:   "123",
			CASetName: "abc",
		}, nil).Times(3)
	}

	mockListCASets := func(client *mtlstruststore.Mock, resp *mtlstruststore.ListCASetsResponse, err error) *mock.Call {
		return client.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{CASetNamePrefix: "abc"}).
			Return(resp, err).Times(3)
	}

	mockListCASetAssociations := func(client *mtlstruststore.Mock, resp *mtlstruststore.ListCASetAssociationsResponse, err error) *mock.Call {
		return client.On("ListCASetAssociations", testutils.MockContext, mtlstruststore.ListCASetAssociationsRequest{
			CASetID: "123",
		}).Return(resp, err).Times(3)
	}

	defaultStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_associations.test").
		CheckEqual("id", "123").
		CheckEqual("name", "abc").
		CheckEqual("properties.#", "0").
		CheckEqual("enrollments.#", "0")

	tests := map[string]struct {
		init  func(client *mtlstruststore.Mock)
		steps []resource.TestStep
	}{
		"id provided, empty response": {
			init: func(client *mtlstruststore.Mock) {
				mockGetCASet(client)
				mockListCASetAssociations(client, &mtlstruststore.ListCASetAssociationsResponse{
					Associations: mtlstruststore.Associations{
						Properties:  nil,
						Enrollments: []mtlstruststore.AssociationEnrollment{},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"id.tf"),
					Check:  defaultStateChecker.Build(),
				},
			},
		},
		"name provided, empty response": {
			init: func(client *mtlstruststore.Mock) {
				mockListCASets(client, &mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "1234",
							CASetName:   "abcd",
							CASetStatus: "NOT_DELETED",
						},
						{
							CASetID:     "123",
							CASetName:   "abc",
							CASetStatus: "NOT_DELETED",
						},
					},
				}, nil)

				mockListCASetAssociations(client, &mtlstruststore.ListCASetAssociationsResponse{
					Associations: mtlstruststore.Associations{
						Properties:  nil,
						Enrollments: []mtlstruststore.AssociationEnrollment{},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"name.tf"),
					Check:  defaultStateChecker.Build(),
				},
			},
		},
		"associated to non-navigable property": {
			init: func(client *mtlstruststore.Mock) {
				mockGetCASet(client)
				mockListCASetAssociations(client, &mtlstruststore.ListCASetAssociationsResponse{
					Associations: mtlstruststore.Associations{
						Properties: []mtlstruststore.AssociationProperty{
							{
								PropertyID: "123",
								Hostnames: []mtlstruststore.AssociationHostname{
									{
										Hostname: "example.com",
										Network:  "STAGING",
										Status:   "ATTACHED",
									},
								},
							},
						}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"id.tf"),
					Check: defaultStateChecker.
						CheckEqual("properties.#", "1").
						CheckEqual("properties.0.property_id", "123").
						CheckMissing("properties.0.property_name").
						CheckMissing("properties.0.asset_id").
						CheckMissing("properties.0.group_id").
						CheckEqual("properties.0.hostnames.#", "1").
						CheckEqual("properties.0.hostnames.0.hostname", "example.com").
						CheckEqual("properties.0.hostnames.0.network", "STAGING").
						CheckEqual("properties.0.hostnames.0.status", "ATTACHED").
						Build(),
				},
			},
		},
		"associated to navigable property": {
			init: func(client *mtlstruststore.Mock) {
				mockGetCASet(client)
				mockListCASetAssociations(client, &mtlstruststore.ListCASetAssociationsResponse{
					Associations: mtlstruststore.Associations{
						Properties: []mtlstruststore.AssociationProperty{
							{
								PropertyID:   "123",
								PropertyName: ptr.To("test-prp-name"),
								AssetID:      ptr.To(int64(123456)),
								GroupID:      ptr.To(int64(345)),
								Hostnames: []mtlstruststore.AssociationHostname{
									{
										Hostname: "example.com",
										Network:  "STAGING",
										Status:   "ATTACHED",
									},
								},
							},
						},
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"id.tf"),
					Check: defaultStateChecker.
						CheckEqual("properties.#", "1").
						CheckEqual("properties.0.property_id", "123").
						CheckEqual("properties.0.property_name", "test-prp-name").
						CheckEqual("properties.0.asset_id", "123456").
						CheckEqual("properties.0.group_id", "345").
						CheckEqual("properties.0.hostnames.#", "1").
						CheckEqual("properties.0.hostnames.0.hostname", "example.com").
						CheckEqual("properties.0.hostnames.0.network", "STAGING").
						CheckEqual("properties.0.hostnames.0.status", "ATTACHED").
						Build(),
				},
			},
		},
		"associated to enrollments": {
			init: func(client *mtlstruststore.Mock) {
				mockGetCASet(client)
				mockListCASetAssociations(client, &mtlstruststore.ListCASetAssociationsResponse{
					Associations: mtlstruststore.Associations{
						Enrollments: []mtlstruststore.AssociationEnrollment{
							{
								EnrollmentID:    123456,
								EnrollmentLink:  "/cps/v2/enrollments/123456",
								StagingSlots:    []int64{78956},
								ProductionSlots: []int64{78956, 56478},
								CN:              "example1.com",
							},
							{
								EnrollmentID:    123457,
								EnrollmentLink:  "/cps/v2/enrollments/123457",
								StagingSlots:    []int64{23456},
								ProductionSlots: []int64{23456},
								CN:              "example2.com",
							},
						}},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"id.tf"),
					Check: defaultStateChecker.
						CheckEqual("enrollments.#", "2").
						CheckEqual("enrollments.0.enrollment_id", "123456").
						CheckEqual("enrollments.0.staging_slots.#", "1").
						CheckEqual("enrollments.0.staging_slots.0", "78956").
						CheckEqual("enrollments.0.production_slots.#", "2").
						CheckEqual("enrollments.0.production_slots.0", "78956").
						CheckEqual("enrollments.0.production_slots.1", "56478").
						CheckEqual("enrollments.0.cn", "example1.com").
						CheckEqual("enrollments.1.enrollment_id", "123457").
						CheckEqual("enrollments.1.staging_slots.#", "1").
						CheckEqual("enrollments.1.staging_slots.0", "23456").
						CheckEqual("enrollments.1.production_slots.#", "1").
						CheckEqual("enrollments.1.production_slots.0", "23456").
						CheckEqual("enrollments.1.cn", "example2.com").
						Build(),
				},
			},
		},
		"validation error - neither id or name provided": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_associations" "test" {
}
`,
					ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[id,name] is required`),
				},
			},
		},
		"validation error - both id and name provided": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_associations" "test" {
  id   = "123"
  name = "abc"
}
`,
					ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of \[id,name] is required`),
				},
			},
		},
		"validation error - empty id provided": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_associations" "test" {
  id   = ""
}
`,
					ExpectError: regexp.MustCompile(`Attribute id string length must be at least 1, got: 0`),
				},
			},
		},
		"validation error - empty name provided": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_associations" "test" {
  name   = ""
}
`,
					ExpectError: regexp.MustCompile(`Attribute name string length must be between 3 and 64, got: 0`),
				},
			},
		},
		"api error - non-existing id": {
			init: func(client *mtlstruststore.Mock) {
				client.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: "123",
				}).Return(nil, mtlstruststore.ErrGetCASetNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"id.tf"),
					ExpectError: regexp.MustCompile(`Error: Could not fetch CA set`),
				},
			},
		},
		"error - name not found": {
			init: func(client *mtlstruststore.Mock) {
				mockListCASets(client, &mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "1234",
							CASetName:   "abcd",
							CASetStatus: "NOT_DELETED",
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"name.tf"),
					ExpectError: regexp.MustCompile(`Error: Could not fetch CA set ID for provided name`),
				},
			},
		},
		"api error - fail on finding CA Set by name": {
			init: func(client *mtlstruststore.Mock) {
				mockListCASets(client, nil, &mtlstruststore.Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"name.tf"),
					ExpectError: regexp.MustCompile(`Error: Could not fetch CA set ID for provided name`),
				},
			},
		},
		"api error - fetching associations fails": {
			init: func(client *mtlstruststore.Mock) {
				mockGetCASet(client).Once()
				mockListCASetAssociations(client, nil, &mtlstruststore.Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"id.tf"),
					ExpectError: regexp.MustCompile(`Error: Error fetching CA set associations`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mtlstruststore.Mock{}
			if tc.init != nil {
				tc.init(client)
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
