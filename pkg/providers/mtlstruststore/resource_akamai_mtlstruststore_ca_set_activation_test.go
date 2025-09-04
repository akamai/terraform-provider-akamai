package mtlstruststore

import (
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestCASetActivationResource(t *testing.T) {
	pollingInterval = 1 * time.Millisecond
	mockListCASetActivations := func(client *mtlstruststore.Mock, testData commonDataForResource, activated bool) *mock.Call {
		var activations []mtlstruststore.ActivateCASetVersionResponse
		if activated {
			activations = []mtlstruststore.ActivateCASetVersionResponse{
				{
					ActivationID:     1,
					CASetID:          testData.caSetID,
					Version:          testData.version,
					Network:          "STAGING",
					ActivationStatus: "COMPLETE",
					ActivationType:   "ACTIVATE",
					CreatedBy:        "user1",
					CreatedDate:      time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					ModifiedBy:       ptr.To("user2"),
					ModifiedDate:     ptr.To(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
				{
					ActivationID:     2,
					CASetID:          testData.caSetID,
					Version:          testData.version,
					Network:          "STAGING",
					ActivationStatus: "FAILED",
					ActivationType:   "DEACTIVATE",
					CreatedBy:        "user1",
					CreatedDate:      time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
					ModifiedBy:       ptr.To("user2"),
					ModifiedDate:     ptr.To(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)),
				},
			}
		} else {
			activations = []mtlstruststore.ActivateCASetVersionResponse{}
		}
		return client.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
			CASetID: testData.caSetID,
		}).
			Return(&mtlstruststore.ListCASetActivationsResponse{
				Activations: activations,
			}, nil)
	}
	mockGetCASetVersion := func(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
		var certificateResponse []mtlstruststore.CertificateResponse
		for _, c := range testData.certificates {
			certificateResponse = append(certificateResponse, mtlstruststore.CertificateResponse{
				CertificatePEM: c.CertificatePEM,
				Description:    c.Description,
			})
		}
		return client.On("GetCASetVersion", testutils.MockContext, mtlstruststore.GetCASetVersionRequest{
			CASetID: testData.caSetID,
			Version: testData.version,
		}).
			Return(&mtlstruststore.GetCASetVersionResponse{
				CASetID:           testData.caSetID,
				Version:           testData.version,
				CASetName:         testData.name,
				VersionLink:       "",
				Description:       testData.versionDescription,
				AllowInsecureSHA1: testData.allowInsecureSHA1,
				StagingStatus:     testData.stagingStatus,
				ProductionStatus:  testData.productionStatus,
				CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
				CreatedBy:         "someone",
				ModifiedDate:      nil,
				ModifiedBy:        nil,
				Certificates:      certificateResponse,
			}, nil)
	}

	createActivationData := commonDataForResource{
		caSetID:            "12345",
		version:            1,
		versionDescription: ptr.To("Test Version"),
		stagingStatus:      "INACTIVE",
		stagingVersion:     ptr.To(int64(1)),
		productionStatus:   "ACTIVE",
	}

	tests := map[string]struct {
		init     func(*mtlstruststore.Mock, commonDataForResource)
		mockData commonDataForResource
		steps    []resource.TestStep
	}{
		// create.
		"create ca set activation - successful": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASet(m, resourceData).Once()
				mockListCASetVersionActivations(m, resourceData, true)
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockDeactivateCASetActivation(m, resourceData, 1)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("id", "1").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_by", "user2").Build(),
				},
			},
		},
		"create ca set activation - delay in activation": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				// Simulate delay in activation response by returning "IN_PROGRESS" status.
				mockGetCASetVersionActivation(m, resourceData, 1, "IN_PROGRESS", "ACTIVATE", 15)
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASet(m, resourceData).Once()
				mockListCASetVersionActivations(m, resourceData, true)
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockDeactivateCASetActivation(m, resourceData, 1)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("id", "1").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_by", "user2").Build(),
				},
			},
		},
		"create ca set activation - delay in response - context deadline exceeded": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				m.On("GetCASetVersionActivation", testutils.MockContext, mtlstruststore.GetCASetVersionActivationRequest{
					ActivationID: 1,
					CASetID:      resourceData.caSetID,
					Version:      resourceData.version,
				}).Return(&mtlstruststore.GetCASetVersionActivationResponse{
					ActivationID:     1,
					CASetID:          resourceData.caSetID,
					Version:          resourceData.version,
					ActivationStatus: "IN_PROGRESS",
					ActivationType:   "ACTIVATE",
					CreatedBy:        "user1",
					CreatedDate:      time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
					ModifiedBy:       ptr.To("user2"),
					ModifiedDate:     ptr.To(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)),
				}, nil).Run(func(_ mock.Arguments) {
					time.Sleep(30 * time.Millisecond) // Simulate delay in response - create time out is set to 20ms in the config.
				}).Maybe()
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_with_timeouts.tf"),
					ExpectError: regexp.MustCompile("CA Set 12345, Version 1: context deadline exceeded"),
				},
			},
		},
		"create ca set activation - CA set version already active": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASetVersion(m, resourceData).Once()
				mockGetCASet(m, resourceData).Once()
				m.On("ListCASetVersionActivations", testutils.MockContext, mtlstruststore.ListCASetVersionActivationsRequest{
					CASetID: resourceData.caSetID,
					Version: 1,
				}).Return(&mtlstruststore.ListCASetVersionActivationsResponse{
					Activations: []mtlstruststore.ActivateCASetVersionResponse{
						{
							ActivationID:     3,
							CASetID:          resourceData.caSetID,
							Version:          1,
							Network:          "STAGING",
							ActivationStatus: "COMPLETE",
							ActivationType:   "ACTIVATE",
							CreatedBy:        "user2",
							CreatedDate:      time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
							ModifiedBy:       ptr.To("user2"),
							ModifiedDate:     ptr.To(time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC)),
						},
						{
							ActivationID:     1,
							CASetID:          resourceData.caSetID,
							Version:          1,
							Network:          "STAGING",
							ActivationStatus: "COMPLETE",
							ActivationType:   "DEACTIVATE",
							CreatedBy:        "user1",
							CreatedDate:      time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
							ModifiedBy:       ptr.To("user1"),
							ModifiedDate:     ptr.To(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)),
						},
					},
				}, nil).Times(2)
				mockListCASetActivations(m, resourceData, true)
				mockGetCASetVersionActivation(m, resourceData, 3, "COMPLETE", "DEACTIVATE", 1)
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockDeactivateCASetActivation(m, resourceData, 3)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("id", "3").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user2").
						CheckEqual("modified_by", "user2").Build(),
				},
			},
		},
		"create ca set activation - error - CA set not found": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				m.On("GetCASetVersion", testutils.MockContext, mtlstruststore.GetCASetVersionRequest{
					CASetID: resourceData.caSetID,
					Version: resourceData.version,
				}).Return(nil, &mtlstruststore.Error{
					Detail:   "Cannot get CA set as the CA set with caSetId 8036389 is not found.",
					Instance: "/mtls-edge-truststore/error-types/ca-set-not-found?traceId=-7121205729277116366",
					Status:   404,
					Title:    "CA set not found.",
					Type:     "/mtls-edge-truststore/error-types/ca-set-not-found",
				})
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					ExpectError: regexp.MustCompile("activation failed"),
				},
			},
		},
		"create ca set activation - error - CA set version activation/deactivation In-Progress": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockGetCASetVersion(m, resourceData).Once()
				resourceData.stagingStatus = "IN_PROGRESS"
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: resourceData.caSetID,
				}).
					Return(&mtlstruststore.ListCASetActivationsResponse{
						Activations: []mtlstruststore.ActivateCASetVersionResponse{
							{
								ActivationID:     1,
								CASetID:          resourceData.caSetID,
								Version:          2,
								Network:          "STAGING",
								ActivationStatus: "IN_PROGRESS",
								ActivationType:   "ACTIVATE",
								CreatedBy:        "user1",
								CreatedDate:      tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
								ModifiedBy:       ptr.To("user2"),
								ModifiedDate:     ptr.To(tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z")),
							},
						},
					}, nil)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					ExpectError: regexp.MustCompile("activation already in progress for version 2"),
				},
			},
		},
		"create ca set activation - error - CA set activation failed": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				m.On("ActivateCASetVersion", testutils.MockContext, mtlstruststore.ActivateCASetVersionRequest{
					CASetID: resourceData.caSetID,
					Version: resourceData.version,
					Network: "STAGING",
				}).Return(nil, &mtlstruststore.Error{Detail: "Activation failed due to an unexpected error.",
					Status: 500,
					Title:  "Activation Error",
					Type:   "/mtls-edge-truststore/error-types/activation-failed"})
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					ExpectError: regexp.MustCompile("activation failed"),
				},
			},
		},
		"create ca set activation - error - GetCASetVersionActivation polling failed": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				m.On("GetCASetVersionActivation", testutils.MockContext, mtlstruststore.GetCASetVersionActivationRequest{
					ActivationID: 1,
					CASetID:      resourceData.caSetID,
					Version:      resourceData.version,
				}).Return(nil, &mtlstruststore.Error{
					Instance: "/mtls-edge-truststore/error-types/internal-error?traceId=-22899532639427393",
					Status:   500,
					Title:    "An unexpected error occurred.",
					Type:     "/mtls-edge-truststore/error-types/internal-error"}).Once()
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					ExpectError: regexp.MustCompile("activation polling failed: error checking activation status for CA Set 12345"),
				},
			},
		},

		// validation.
		"create ca set activation - without `ca_set_id`": {
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_no_ca_set_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "ca_set_id" is required, but no definition was found.`),
				},
			},
		},
		"create ca set activation - without `network`": {
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_no_network.tf"),
					ExpectError: regexp.MustCompile(`The argument "network" is required, but no definition was found.`),
				},
			},
		},
		"create ca set activation - without `version`": {
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_no_version.tf"),
					ExpectError: regexp.MustCompile(`The argument "version" is required, but no definition was found.`),
				},
			},
		},

		// Drift Handling.
		// current version in the state is now deactivated and later the same version was newly activated outside terraform.
		"update ca set activation - same version reactivated outside terraform": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true).Once()
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"

				// read.
				mockGetCASet(m, resourceData).Twice()
				mockListCASetVersionActivationsDrift(m, resourceData, true).Once()

				// read.
				resourceData.stagingVersion = ptr.To(int64(1))
				mockGetCASet(m, resourceData).Once()
				m.On("ListCASetVersionActivations", testutils.MockContext, mtlstruststore.ListCASetVersionActivationsRequest{
					CASetID: resourceData.caSetID,
					Version: resourceData.version,
				}).Return(&mtlstruststore.ListCASetVersionActivationsResponse{
					Activations: []mtlstruststore.ActivateCASetVersionResponse{
						{
							ActivationID:     2,
							CASetID:          resourceData.caSetID,
							Version:          1,
							Network:          "STAGING",
							ActivationStatus: "COMPLETE",
							ActivationType:   "ACTIVATE",
							CreatedBy:        "user1",
							CreatedDate:      time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
							ModifiedBy:       ptr.To("user1"),
							ModifiedDate:     ptr.To(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)),
						},
					},
				}, nil)

				// delete.
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, resourceData, 1)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("version", "1").
						CheckEqual("id", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_date", "2023-01-01T00:00:00Z").
						CheckEqual("created_date", "2023-01-01T00:00:00Z").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("version", "1").
						CheckEqual("id", "2").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_date", "2024-01-01T00:00:00Z").
						CheckEqual("created_date", "2024-01-01T00:00:00Z").Build(),
				},
			},
		},
		// current version in the state is now deactivated and new version was newly activated outside terraform.
		"update ca set activation - new version activated outside terraform": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true).Once()
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"

				// read.
				mockGetCASet(m, resourceData).Once()
				mockListCASetVersionActivationsDrift(m, resourceData, true).Once()

				// read.
				resourceData.stagingVersion = ptr.To(int64(2))
				mockGetCASet(m, resourceData).Once()
				mockListCASetVersionActivationsDrift(m, resourceData, false).Once()

				// update.
				resourceData.stagingStatus = "INACTIVE"
				resourceData.version = 1
				mockGetCASetVersion(m, resourceData).Once()
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: resourceData.caSetID,
				}).Return(&mtlstruststore.ListCASetActivationsResponse{
					Activations: []mtlstruststore.ActivateCASetVersionResponse{
						{
							ActivationID:     2,
							CASetID:          resourceData.caSetID,
							Version:          2,
							Network:          "STAGING",
							ActivationStatus: "COMPLETE",
							ActivationType:   "ACTIVATE",
							CreatedBy:        "user1",
							CreatedDate:      time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
							ModifiedBy:       ptr.To("user1"),
							ModifiedDate:     ptr.To(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)),
						},
					},
				}, nil).Once()
				mockActivateCASetVersion(m, resourceData, 3, "STAGING")
				m.On("GetCASetVersionActivation", testutils.MockContext, mtlstruststore.GetCASetVersionActivationRequest{
					ActivationID: 3,
					CASetID:      resourceData.caSetID,
					Version:      1,
				}).Return(&mtlstruststore.GetCASetVersionActivationResponse{
					ActivationID:     3,
					CASetID:          resourceData.caSetID,
					Version:          1,
					ActivationStatus: "COMPLETE",
					ActivationType:   "ACTIVATE",
					CreatedBy:        "user1",
					CreatedDate:      time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
					ModifiedBy:       ptr.To("user2"),
					ModifiedDate:     ptr.To(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)),
				}, nil).Once()

				// read.
				mockGetCASet(m, resourceData).Once()
				m.On("ListCASetVersionActivations", testutils.MockContext, mtlstruststore.ListCASetVersionActivationsRequest{
					CASetID: resourceData.caSetID,
					Version: 1,
				}).Return(&mtlstruststore.ListCASetVersionActivationsResponse{
					Activations: []mtlstruststore.ActivateCASetVersionResponse{
						{
							ActivationID:     3,
							CASetID:          resourceData.caSetID,
							Version:          1,
							Network:          "STAGING",
							ActivationStatus: "COMPLETE",
							ActivationType:   "ACTIVATE",
							CreatedBy:        "user1",
							CreatedDate:      time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
							ModifiedBy:       ptr.To("user1"),
							ModifiedDate:     ptr.To(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)),
						},
					},
				}, nil)

				// delete.
				resourceData.stagingStatus = "INACTIVE"
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: resourceData.caSetID,
				}).Return(&mtlstruststore.ListCASetActivationsResponse{
					Activations: []mtlstruststore.ActivateCASetVersionResponse{
						{
							ActivationID:     3,
							CASetID:          resourceData.caSetID,
							Version:          1,
							Network:          "STAGING",
							ActivationStatus: "COMPLETE",
							ActivationType:   "ACTIVATE",
							CreatedBy:        "user1",
							CreatedDate:      time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
							ModifiedBy:       ptr.To("user1"),
							ModifiedDate:     ptr.To(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)),
						},
					},
				}, nil).Once()
				mockGetCASetVersion(m, resourceData).Once()
				mockGetCASetVersionActivation(m, resourceData, 3, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, resourceData, 3).Once()
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("version", "1").
						CheckEqual("id", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_date", "2023-01-01T00:00:00Z").
						CheckEqual("created_date", "2023-01-01T00:00:00Z").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("version", "1").
						CheckEqual("id", "3").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_date", "2025-01-01T00:00:00Z").
						CheckEqual("created_date", "2025-01-01T00:00:00Z").Build(),
				},
			},
		},
		// ca-set is removed outside terraform during the update - since not found it clears the state.
		"update ca set activation - ca set dropped outside terraform": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true).Once()
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"

				// read.
				mockGetCASet(m, resourceData).Once()
				mockListCASetVersionActivationsDrift(m, resourceData, true).Once()

				// read.
				resourceData.stagingVersion = ptr.To(int64(2))
				m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: resourceData.caSetID,
				}).Return(nil, mtlstruststore.ErrGetCASetNotFound).Once()

				m.On("GetCASetVersion", testutils.MockContext, mtlstruststore.GetCASetVersionRequest{
					CASetID: resourceData.caSetID,
					Version: resourceData.version,
				}).Return(nil, mtlstruststore.ErrGetCASetNotFound).Once().Once()
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("version", "1").
						CheckEqual("id", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_date", "2023-01-01T00:00:00Z").
						CheckEqual("created_date", "2023-01-01T00:00:00Z").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update_with_timeouts.tf"),
					ExpectError: regexp.MustCompile("activation failed")},
			},
		},
		// update.
		"update ca set activation - successful": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true).Once()
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"

				// read.
				mockGetCASet(m, resourceData).Twice()
				mockListCASetVersionActivations(m, resourceData, true)

				// update.
				updateActivationData := resourceData
				updateActivationData.version = 2
				updateActivationData.stagingStatus = "INACTIVE"
				mockGetCASetVersion(m, updateActivationData).Once()
				mockListCASetActivations(m, updateActivationData, true)
				mockActivateCASetVersion(m, updateActivationData, 2, "STAGING")
				mockGetCASetVersionActivation(m, updateActivationData, 2, "COMPLETE", "ACTIVATE", 1)

				// read.
				mockGetCASet(m, updateActivationData).Once()
				mockListCASetVersionActivations(m, updateActivationData, true)

				// delete.
				mockListCASetActivations(m, updateActivationData, true)
				mockGetCASetVersion(m, updateActivationData).Once()
				mockGetCASetVersionActivation(m, updateActivationData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, updateActivationData, 1)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("id", "1").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_by", "user2").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("id", "2").
						CheckEqual("version", "2").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_by", "user2").Build(),
				},
			},
		},
		"update ca set activation - updating `timeout` only ": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASet(m, resourceData).Times(3)
				mockListCASetVersionActivations(m, resourceData, true)

				// delete.
				resourceData.stagingStatus = "ACTIVE"
				mockListCASetActivations(m, resourceData, true)
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, resourceData, 1)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("timeouts.create", "2m").
						CheckEqual("timeouts.update", "2m").
						CheckEqual("timeouts.delete", "1m").
						CheckEqual("modified_by", "user2").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update_with_timeouts.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("timeouts.create", "10m").
						CheckEqual("timeouts.update", "11m").
						CheckEqual("timeouts.delete", "12m").
						CheckEqual("modified_by", "user2").Build(),
				},
			},
		},
		"update ca set activation - failed": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASet(m, resourceData).Twice()
				mockListCASetVersionActivations(m, resourceData, true)
				// update.
				resourceData.stagingStatus = "INACTIVE"
				updateActivationData := resourceData
				updateActivationData.version = 2
				mockGetCASetVersion(m, updateActivationData).Once()
				mockListCASetActivations(m, updateActivationData, true)
				m.On("ActivateCASetVersion", testutils.MockContext, mtlstruststore.ActivateCASetVersionRequest{
					CASetID: updateActivationData.caSetID,
					Version: updateActivationData.version,
					Network: "STAGING",
				}).Return(nil, &mtlstruststore.Error{Detail: "Activation failed due to an unexpected error.",
					Status: 500,
					Title:  "Activation Error",
					Type:   "/mtls-edge-truststore/error-types/activation-failed"})

				// delete.
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, resourceData, 1)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_mtlstruststore_ca_set_activation.test", "ca_set_id", "12345"),
						resource.TestCheckResourceAttr("akamai_mtlstruststore_ca_set_activation.test", "id", "1"),
						resource.TestCheckResourceAttr("akamai_mtlstruststore_ca_set_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_mtlstruststore_ca_set_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_mtlstruststore_ca_set_activation.test", "created_by", "user1"),
						resource.TestCheckResourceAttr("akamai_mtlstruststore_ca_set_activation.test", "modified_by", "user2"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update.tf"),
					ExpectError: regexp.MustCompile("activation failed"),
				},
			},
		},
		"update ca set activation - failed - updating field `ca_set_id` is not possible ": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASet(m, resourceData).Times(2)
				mockListCASetVersionActivations(m, resourceData, true)

				// delete.
				resourceData.stagingStatus = "ACTIVE"
				mockListCASetActivations(m, resourceData, true)
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, resourceData, 1)

			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("id", "1").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_by", "user2").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update_ca_set_id.tf"),
					ExpectError: regexp.MustCompile("updating field `ca_set_id` is not possible"),
				},
			},
		},
		"update ca set activation - failed - updating field `network` is not possible ": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetActivations(m, resourceData, true)
				mockActivateCASetVersion(m, resourceData, 1, "STAGING")
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "ACTIVATE", 1)
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASet(m, resourceData).Times(2)
				mockListCASetVersionActivations(m, resourceData, true)

				// delete.
				resourceData.stagingStatus = "ACTIVE"
				mockListCASetActivations(m, resourceData, true)
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, resourceData, 1)

			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/create.tf"),
					Check: test.NewStateChecker("akamai_mtlstruststore_ca_set_activation.test").
						CheckEqual("ca_set_id", "12345").
						CheckEqual("id", "1").
						CheckEqual("version", "1").
						CheckEqual("network", "STAGING").
						CheckEqual("created_by", "user1").
						CheckEqual("modified_by", "user2").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/update_network.tf"),
					ExpectError: regexp.MustCompile("updating field `network` is not possible"),
				},
			},
		},

		// import.
		"import ca set activation - successful": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// import.
				mockListCASetActivations(m, resourceData, true)
				// read.
				resourceData.stagingStatus = "ACTIVE"
				mockGetCASet(m, resourceData).Once()
				mockListCASetVersionActivations(m, resourceData, true)

				// delete.
				mockListCASetActivations(m, resourceData, true)
				mockGetCASetVersion(m, resourceData).Once()
				mockListCASetAssociations(m, resourceData).Once()
				mockGetCASetVersionActivation(m, resourceData, 1, "COMPLETE", "DEACTIVATE", 1)
				mockDeactivateCASetActivation(m, resourceData, 1)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "12345:STAGING",
					ResourceName:                         "akamai_mtlstruststore_ca_set_activation.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
				},
			},
		},
		"import ca set activation - failed - version not active": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// import.
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: resourceData.caSetID,
				}).
					Return(&mtlstruststore.ListCASetActivationsResponse{
						Activations: []mtlstruststore.ActivateCASetVersionResponse{
							{
								ActivationID:     1,
								CASetID:          resourceData.caSetID,
								Version:          resourceData.version,
								Network:          "STAGING",
								ActivationStatus: "COMPLETE",
								ActivationType:   "DEACTIVATE",
								CreatedBy:        "user1",
								CreatedDate:      tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
								ModifiedBy:       ptr.To("user2"),
								ModifiedDate:     ptr.To(tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z")),
							},
						},
					}, nil)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "12345:STAGING",
					ResourceName:                         "akamai_mtlstruststore_ca_set_activation.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile(`CA set with ID 12345 is not active in the STAGING network. Only completed(\s|\n)+activations can be imported.`),
				},
			},
		},
		"import ca set activation - failed - version activation In-Progress": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// import.
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: resourceData.caSetID,
				}).
					Return(&mtlstruststore.ListCASetActivationsResponse{
						Activations: []mtlstruststore.ActivateCASetVersionResponse{
							{
								ActivationID:     1,
								CASetID:          resourceData.caSetID,
								Version:          resourceData.version,
								Network:          "STAGING",
								ActivationStatus: "IN_PROGRESS",
								ActivationType:   "ACTIVATE",
								CreatedBy:        "user1",
								CreatedDate:      tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
								ModifiedBy:       ptr.To("user2"),
								ModifiedDate:     ptr.To(tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z")),
							},
						},
					}, nil)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "12345:STAGING",
					ResourceName:                         "akamai_mtlstruststore_ca_set_activation.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("Error: Operation in progress"),
				},
			},
		},
		"import ca set activation - invalid import ID format": {
			init: func(_ *mtlstruststore.Mock, _ commonDataForResource) {
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "123456789",
					ResourceName:                         "akamai_mtlstruststore_ca_set_activation.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("Error: Invalid import ID format"),
				},
			},
		},
		"import ca set activation - validation error - invalid network": {
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "12345:SANDBOX",
					ResourceName:                         "akamai_mtlstruststore_ca_set_activation.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("Error: Invalid network"),
				},
			},
		},
		"import ca set activation - validation error - invalid ca_set_id": {
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/import.tf"),
					ImportState:                          true,
					ImportStateId:                        ":SANDBOX",
					ResourceName:                         "akamai_mtlstruststore_ca_set_activation.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("Error: Invalid CA set ID"),
				},
			},
		},
		"import ca set activation - failed - ca_set not found": {
			init: func(m *mtlstruststore.Mock, _ commonDataForResource) {
				// import.
				m.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
					CASetID: createActivationData.caSetID,
				}).Return(nil, mtlstruststore.ErrGetCASetNotFound)
			},
			mockData: createActivationData,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASetActivation/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "12345:STAGING",
					ResourceName:                         "akamai_mtlstruststore_ca_set_activation.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("CA set with ID 12345 not found: ca set not found"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mtlstruststore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.mockData)
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

func mockActivateCASetVersion(m *mtlstruststore.Mock, data commonDataForResource, ActivationID int64, Network string) {
	m.On("ActivateCASetVersion", testutils.MockContext, mtlstruststore.ActivateCASetVersionRequest{
		CASetID: data.caSetID,
		Version: data.version,
		Network: mtlstruststore.ActivationNetwork(Network),
	}).Return(&mtlstruststore.ActivateCASetVersionResponse{
		ActivationID: ActivationID,
		CASetID:      data.caSetID,
		Version:      data.version,
	}, nil).Once()
}

func mockGetCASetVersionActivation(m *mtlstruststore.Mock, data commonDataForResource, ActivationID int64, ActivationStatus, ActivationType string, times int) {
	m.On("GetCASetVersionActivation", testutils.MockContext, mtlstruststore.GetCASetVersionActivationRequest{
		ActivationID: ActivationID,
		CASetID:      data.caSetID,
		Version:      data.version,
	}).Return(&mtlstruststore.GetCASetVersionActivationResponse{
		ActivationID:     ActivationID,
		CASetID:          data.caSetID,
		Version:          data.version,
		ActivationStatus: ActivationStatus,
		ActivationType:   ActivationType,
		CreatedBy:        "user1",
		CreatedDate:      time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		ModifiedBy:       ptr.To("user2"),
		ModifiedDate:     ptr.To(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)),
	}, nil).Times(times)
}

func mockDeactivateCASetActivation(client *mtlstruststore.Mock, testData commonDataForResource, ActivationID int64) *mock.Call {
	return client.On("DeactivateCASetVersion", testutils.MockContext, mtlstruststore.DeactivateCASetVersionRequest{
		CASetID: testData.caSetID,
		Version: testData.version,
		Network: mtlstruststore.ActivationNetworkStaging,
	}).Return(&mtlstruststore.DeactivateCASetVersionResponse{
		ActivationID: ActivationID,
		CASetID:      testData.caSetID,
		Version:      testData.version,
	}, nil).Once()
}

func mockListCASetVersionActivations(client *mtlstruststore.Mock, testData commonDataForResource, activated bool) *mock.Call {
	var activations []mtlstruststore.ActivateCASetVersionResponse
	var network string
	if testData.stagingVersion != nil {
		network = "STAGING"
	} else {
		network = "PRODUCTION"
	}
	if activated {
		activations = []mtlstruststore.ActivateCASetVersionResponse{
			{
				ActivationID:     1,
				CASetID:          testData.caSetID,
				Version:          testData.version,
				Network:          network,
				ActivationStatus: "COMPLETE",
				ActivationType:   "DEACTIVATE",
				CreatedBy:        "user1",
				CreatedDate:      time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
				ModifiedBy:       ptr.To("user2"),
				ModifiedDate:     ptr.To(time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
			{
				ActivationID:     2,
				CASetID:          testData.caSetID,
				Version:          testData.version,
				Network:          network,
				ActivationStatus: "COMPLETE",
				ActivationType:   "ACTIVATE",
				CreatedBy:        "user1",
				CreatedDate:      time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				ModifiedBy:       ptr.To("user2"),
				ModifiedDate:     ptr.To(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
	return client.On("ListCASetVersionActivations", testutils.MockContext, mtlstruststore.ListCASetVersionActivationsRequest{
		CASetID: testData.caSetID,
		Version: testData.version,
	}).Return(&mtlstruststore.ListCASetVersionActivationsResponse{
		Activations: activations,
	}, nil)
}

func mockListCASetVersionActivationsDrift(client *mtlstruststore.Mock, testData commonDataForResource, beforeDrift bool) *mock.Call {
	var activations []mtlstruststore.ActivateCASetVersionResponse
	var network string
	if testData.stagingVersion != nil {
		network = "STAGING"
	} else {
		network = "PRODUCTION"
	}
	if beforeDrift {
		activations = []mtlstruststore.ActivateCASetVersionResponse{
			{
				ActivationID:     1,
				CASetID:          testData.caSetID,
				Version:          1,
				Network:          network,
				ActivationStatus: "COMPLETE",
				ActivationType:   "ACTIVATE",
				CreatedBy:        "user1",
				CreatedDate:      time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				ModifiedBy:       ptr.To("user1"),
				ModifiedDate:     ptr.To(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	} else { // after drift.
		activations = []mtlstruststore.ActivateCASetVersionResponse{
			{
				ActivationID:     1,
				CASetID:          testData.caSetID,
				Version:          1,
				Network:          network,
				ActivationStatus: "COMPLETE",
				ActivationType:   "DEACTIVATE",
				CreatedBy:        "user1",
				CreatedDate:      time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
				ModifiedBy:       ptr.To("user1"),
				ModifiedDate:     ptr.To(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
			{
				ActivationID:     2,
				CASetID:          testData.caSetID,
				Version:          2,
				Network:          network,
				ActivationStatus: "COMPLETE",
				ActivationType:   "ACTIVATE",
				CreatedBy:        "user1",
				CreatedDate:      time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
				ModifiedBy:       ptr.To("user1"),
				ModifiedDate:     ptr.To(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)),
			},
		}
	}
	return client.On("ListCASetVersionActivations", testutils.MockContext, mtlstruststore.ListCASetVersionActivationsRequest{
		CASetID: testData.caSetID,
		Version: testData.version,
	}).Return(&mtlstruststore.ListCASetVersionActivationsResponse{
		Activations: activations,
	}, nil)
}
