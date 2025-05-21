package cloudlets

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceApplicationLoadBalancer(t *testing.T) {

	type loadBalancerAttributes struct {
		originID, version, originDescription, description, balancingType string
	}
	var (
		expectCreateLoadBalancer = func(client *cloudlets.Mock, originID, originDescription, description, balancingType string, version int64, livenessHosts []string) (*cloudlets.Origin, *cloudlets.LoadBalancerVersion) {
			loadBalancerConfig := cloudlets.CreateOriginRequest{
				OriginID: originID,
				Description: cloudlets.Description{
					Description: originDescription,
				},
			}
			loadBalancerVersionReq := cloudlets.LoadBalancerVersion{
				Description:   description,
				BalancingType: cloudlets.BalancingType(balancingType),
				DataCenters: []cloudlets.DataCenter{
					{
						City:            "Boston",
						CloudService:    true,
						Continent:       "NA",
						Country:         "US",
						Hostname:        "test-hostname",
						Latitude:        ptr.To(102.78108),
						LivenessHosts:   livenessHosts,
						Longitude:       ptr.To(-116.07064),
						OriginID:        "test_origin",
						Percent:         ptr.To(100.0),
						StateOrProvince: ptr.To("MA"),
					},
				},
				LivenessSettings: &cloudlets.LivenessSettings{
					HostHeader:        "header",
					AdditionalHeaders: map[string]string{"abc": "123"},
					Interval:          10,
					Path:              "/status",
					Port:              1234,
					Protocol:          "HTTP",
					RequestString:     "test_request_string",
					ResponseString:    "test_response_string",
					Timeout:           60,
				},
			}
			loadBalancerVersionResp := loadBalancerVersionReq
			loadBalancerVersionResp.Version = version
			loadBalancerVersionResp.Warnings = []cloudlets.Warning{
				{
					Detail:      "test warning details",
					JSONPointer: "/path",
					Title:       "test warning",
					Type:        "test type",
				},
			}
			origin := cloudlets.Origin{
				OriginID:    originID,
				Description: originDescription,
				Type:        "APPLICATION_LOAD_BALANCER",
			}
			origins := []cloudlets.OriginResponse{
				{
					Hostname: "test-hostname",
					Origin: cloudlets.Origin{
						OriginID:  "test_origin_2",
						Akamaized: false,
					},
				},
			}
			client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()

			client.On("CreateOrigin", testutils.MockContext, loadBalancerConfig).Return(&origin, nil).Once()
			client.On("CreateLoadBalancerVersion", testutils.MockContext, cloudlets.CreateLoadBalancerVersionRequest{
				OriginID:            originID,
				LoadBalancerVersion: loadBalancerVersionReq,
			}).Return(&loadBalancerVersionResp, nil).Once()

			return &origin, &loadBalancerVersionResp
		}

		expectReadLoadBalancer = func(client *cloudlets.Mock, origin *cloudlets.Origin, loadBalancerVersion *cloudlets.LoadBalancerVersion, times int) {
			client.On("GetOrigin", testutils.MockContext, cloudlets.GetOriginRequest{
				OriginID: origin.OriginID,
			}).Return(origin, nil).Times(times)
			client.On("GetLoadBalancerVersion", testutils.MockContext, cloudlets.GetLoadBalancerVersionRequest{
				OriginID:       origin.OriginID,
				Version:        loadBalancerVersion.Version,
				ShouldValidate: true,
			}).Return(loadBalancerVersion, nil).Times(times)
		}

		expectCreateLoadBalancerVersion = func(client *cloudlets.Mock, originID string, loadBalancerVersion *cloudlets.LoadBalancerVersion, newBalancingType, description string) *cloudlets.LoadBalancerVersion {
			origins := []cloudlets.OriginResponse{
				{
					Hostname: "test-hostname",
					Origin: cloudlets.Origin{
						OriginID:  "test_origin",
						Akamaized: false,
					},
				},
				{
					Hostname: "test-hostname",
					Origin: cloudlets.Origin{
						OriginID:  "test_origin_2",
						Akamaized: false,
					},
				},
			}
			client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()
			client.On("ListLoadBalancerActivations", testutils.MockContext, cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID}).Return([]cloudlets.LoadBalancerActivation{
				{
					OriginID: originID,
					Network:  cloudlets.LoadBalancerActivationNetworkProduction,
					Version:  loadBalancerVersion.Version,
					Status:   cloudlets.LoadBalancerActivationStatusActive,
				},
			}, nil).Once()
			var newVersionReq, newVersionResp cloudlets.LoadBalancerVersion
			err := copier.CopyWithOption(&newVersionReq, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			err = copier.CopyWithOption(&newVersionResp, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			newVersionReq.Description = description
			newVersionReq.Version = 0
			newVersionReq.Warnings = nil
			newVersionReq.BalancingType = cloudlets.BalancingType(newBalancingType)
			newVersionResp.Description = description
			newVersionResp.BalancingType = cloudlets.BalancingType(newBalancingType)
			newVersionResp.Version++
			client.On("CreateLoadBalancerVersion", testutils.MockContext, cloudlets.CreateLoadBalancerVersionRequest{
				OriginID:            originID,
				LoadBalancerVersion: newVersionReq,
			}).Return(&newVersionResp, nil).Once()
			return &newVersionResp
		}

		expectOriginDescriptionUpdate = func(client *cloudlets.Mock, origin *cloudlets.Origin, description string) *cloudlets.Origin {
			var newOrigin cloudlets.Origin
			err := copier.CopyWithOption(&newOrigin, origin, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			newOrigin.Description = description

			client.On("UpdateOrigin", testutils.MockContext, cloudlets.UpdateOriginRequest{
				OriginID: origin.OriginID,
				Description: cloudlets.Description{
					Description: description,
				},
			}).Return(&newOrigin, nil).Once()

			return &newOrigin
		}

		expectUpdateLoadBalancerVersion = func(client *cloudlets.Mock, originID string, loadBalancerVersion *cloudlets.LoadBalancerVersion, newBalancingType, description string) *cloudlets.LoadBalancerVersion {
			origins := []cloudlets.OriginResponse{
				{
					Hostname: "test-hostname",
					Origin: cloudlets.Origin{
						OriginID:  "test_origin",
						Akamaized: false,
					},
				},
				{
					Hostname: "test-hostname",
					Origin: cloudlets.Origin{
						OriginID:  "test_origin_2",
						Akamaized: false,
					},
				},
			}
			client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()
			client.On("ListLoadBalancerActivations", testutils.MockContext, cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID}).Return(nil, nil).Once()
			var updateVersionReq, updateVersionResp cloudlets.LoadBalancerVersion
			err := copier.CopyWithOption(&updateVersionReq, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			err = copier.CopyWithOption(&updateVersionResp, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			if newBalancingType != "" {
				updateVersionReq.BalancingType = cloudlets.BalancingTypePerformance
				updateVersionResp.BalancingType = cloudlets.BalancingType(newBalancingType)
			}

			if description != "" {
				updateVersionReq.Description = description
				updateVersionResp.Description = description
			}

			updateVersionReq.Version = 0
			updateVersionReq.Warnings = nil
			client.On("UpdateLoadBalancerVersion", testutils.MockContext, cloudlets.UpdateLoadBalancerVersionRequest{
				OriginID:            originID,
				ShouldValidate:      true,
				Version:             loadBalancerVersion.Version,
				LoadBalancerVersion: updateVersionReq,
			}).Return(&updateVersionResp, nil).Once()
			return &updateVersionResp
		}

		expectImportLoadBalancer = func(client *cloudlets.Mock, origin *cloudlets.Origin, numVersions int) {
			client.On("GetOrigin", testutils.MockContext, cloudlets.GetOriginRequest{OriginID: origin.OriginID}).Return(origin, nil).Once()

			var versionList []cloudlets.LoadBalancerVersion
			for i := 1; i <= numVersions; i++ {
				versionList = append(versionList, cloudlets.LoadBalancerVersion{OriginID: origin.OriginID, Version: int64(i)})
			}
			client.On("ListLoadBalancerVersions", testutils.MockContext, cloudlets.ListLoadBalancerVersionsRequest{
				OriginID: origin.OriginID,
			}).Return(versionList, nil).Once()
		}

		checkAttributes = func(attrs loadBalancerAttributes) resource.TestCheckFunc {
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "id", attrs.originID),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "origin_description", attrs.originDescription),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "description", attrs.description),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "balancing_type", attrs.balancingType),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.#", "1"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.cloud_server_host_header_override", "false"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.cloud_service", "true"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.country", "US"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.continent", "NA"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.latitude", "102.78108"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.longitude", "-116.07064"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.percent", "100"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.hostname", "test-hostname"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.state_or_province", "MA"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.city", "Boston"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.origin_id", "test_origin"),
			}
			if attrs.version == "" {
				checks = append(checks, resource.TestCheckNoResourceAttr("akamai_cloudlets_application_load_balancer.alb", "version"))
			} else {
				checks = append(checks, resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "version", attrs.version))
			}
			return resource.ComposeAggregateTestCheckFunc(checks...)
		}
	)

	t.Run("load balancer lifecycle with create new version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		lbVersion = expectCreateLoadBalancerVersion(client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "2",
							description:   "test description updated",
							balancingType: "PERFORMANCE",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("load balancer lifecycle with update existing version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		lbVersion = expectUpdateLoadBalancerVersion(client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description updated",
							balancingType: "PERFORMANCE",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only description", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle_origin_update"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		lbVersion = expectUpdateLoadBalancerVersion(client, origin.OriginID, lbVersion, "", "test description updated")

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description updated",
							balancingType: "WEIGHTED",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle_version_update"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		lbVersion = expectUpdateLoadBalancerVersion(client, origin.OriginID, lbVersion, "PERFORMANCE", "")

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "PERFORMANCE",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update version + empty liveness_settings", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle_no_liveness_settings"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		// no liveness_settings for update
		var lbVersionUpdate cloudlets.LoadBalancerVersion
		err := copier.CopyWithOption(&lbVersionUpdate, lbVersion, copier.Option{DeepCopy: true})
		require.NoError(t, err)

		lbVersionUpdate.LivenessSettings = nil
		lbVersionUpdate.Description = "test description updated"

		lbVersionUpdate = *expectUpdateLoadBalancerVersion(client, origin.OriginID, &lbVersionUpdate, "PERFORMANCE", "")

		expectReadLoadBalancer(client, origin, &lbVersionUpdate, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description updated",
							balancingType: "PERFORMANCE",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating origin", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origins := []cloudlets.OriginResponse{
			{
				Hostname: "test-hostname",
				Origin: cloudlets.Origin{
					OriginID:  "test_origin_2",
					Akamaized: false,
				},
			},
		}
		client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()

		client.On("CreateOrigin", testutils.MockContext, cloudlets.CreateOriginRequest{OriginID: "test_origin"}).Return(nil, fmt.Errorf("creating origin")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						ExpectError: regexp.MustCompile("creating origin"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating origin which already exist", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origins := []cloudlets.OriginResponse{
			{
				Hostname: "test-hostname",
				Origin: cloudlets.Origin{
					OriginID:  "test_origin",
					Akamaized: false,
				},
			},
		}
		client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						ExpectError: regexp.MustCompile("already exists"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origins := []cloudlets.OriginResponse{
			{
				Hostname: "test-hostname",
				Origin: cloudlets.Origin{
					OriginID:  "test_origin_2",
					Akamaized: false,
				},
			},
		}
		client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()

		client.On("CreateOrigin", testutils.MockContext, cloudlets.CreateOriginRequest{OriginID: "test_origin"}).Return(&cloudlets.Origin{OriginID: "test_origin"}, nil).Once()

		loadBalancerVersionReq := cloudlets.LoadBalancerVersion{
			Description:   "test description",
			BalancingType: "WEIGHTED",
			DataCenters: []cloudlets.DataCenter{
				{
					City:            "Boston",
					CloudService:    true,
					Continent:       "NA",
					Country:         "US",
					Hostname:        "test-hostname",
					Latitude:        ptr.To(102.78108),
					LivenessHosts:   []string{"tf.test"},
					Longitude:       ptr.To(-116.07064),
					OriginID:        "test_origin",
					Percent:         ptr.To(100.0),
					StateOrProvince: ptr.To("MA"),
				},
			},
			LivenessSettings: &cloudlets.LivenessSettings{
				HostHeader:        "header",
				AdditionalHeaders: map[string]string{"abc": "123"},
				Interval:          10,
				Path:              "/status",
				Port:              1234,
				Protocol:          "HTTP",
				RequestString:     "test_request_string",
				ResponseString:    "test_response_string",
				Timeout:           60,
			},
		}
		client.On("CreateLoadBalancerVersion", testutils.MockContext, cloudlets.CreateLoadBalancerVersionRequest{
			OriginID:            "test_origin",
			LoadBalancerVersion: loadBalancerVersionReq,
		}).Return(nil, fmt.Errorf("creating version")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						ExpectError: regexp.MustCompile("creating version"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("load balancer lifecycle with create new version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		lbVersion = expectCreateLoadBalancerVersion(client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "2",
							description:   "test description updated",
							balancingType: "PERFORMANCE",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("load balancer lifecycle with origin description", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle_origin_desc"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "origin description", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		origin = expectOriginDescriptionUpdate(client, origin, "update origin description")
		lbVersion = expectCreateLoadBalancerVersion(client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:          "test_origin",
							version:           "1",
							originDescription: "origin description",
							description:       "test description",
							balancingType:     "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						Check: checkAttributes(loadBalancerAttributes{
							originID:          "test_origin",
							version:           "2",
							originDescription: "update origin description",
							description:       "test description updated",
							balancingType:     "PERFORMANCE",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error fetching version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, _ := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		client.On("GetOrigin", testutils.MockContext, cloudlets.GetOriginRequest{
			OriginID: "test_origin",
		}).Return(origin, nil).Once()

		client.On("GetLoadBalancerVersion", testutils.MockContext, cloudlets.GetLoadBalancerVersionRequest{
			OriginID:       "test_origin",
			Version:        1,
			ShouldValidate: true,
		}).Return(nil, fmt.Errorf("fetching version")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						ExpectError: regexp.MustCompile("fetching version"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("import load balancer", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		expectImportLoadBalancer(client, origin, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_origin",
						ResourceName:      "akamai_cloudlets_application_load_balancer.alb",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error importing load balancer not found", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		client.On("GetOrigin", testutils.MockContext, cloudlets.GetOriginRequest{OriginID: "not_existing_test_origin"}).Return(nil, fmt.Errorf("could not find origin with origin_id: not_existing_test_origin")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
					},
					{
						ImportState:   true,
						ImportStateId: "not_existing_test_origin",
						ResourceName:  "akamai_cloudlets_application_load_balancer.alb",
						ExpectError:   regexp.MustCompile("could not find origin with origin_id: not_existing_test_origin"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error importing load balancer origin_id cannot be empty", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
					},
					{
						ImportState: true,
						ImportStateIdFunc: func(_ *terraform.State) (string, error) {
							return "", nil
						},
						ResourceName: "akamai_cloudlets_application_load_balancer.alb",
						ExpectError:  regexp.MustCompile("the origin ID cannot be empty"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error importing load balancer no version found", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(client, origin, lbVersion, 2)

		expectImportLoadBalancer(client, origin, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
					},
					{
						ImportState:   true,
						ImportStateId: "test_origin",
						ResourceName:  "akamai_cloudlets_application_load_balancer.alb",
						ExpectError:   regexp.MustCompile("no load balancer version found for the origin_id: test_origin"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating origin with akamaized dc", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origins := []cloudlets.OriginResponse{
			{
				Hostname: "test-hostname",
				Origin: cloudlets.Origin{
					OriginID:  "test_origin",
					Akamaized: true,
				},
			},
		}
		client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						ExpectError: regexp.MustCompile("'liveness_hosts' field should be omitted for GTM hostname: \"test-hostname\". " +
							"Liveness tests for this host can be configured in DNS traffic management"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error updating origin with akamized dc", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle_dc_update"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(client, "test_origin", "", "test description", "WEIGHTED", 1, []string{})

		expectReadLoadBalancer(client, origin, lbVersion, 3)

		origins := []cloudlets.OriginResponse{
			{
				Hostname: "test-hostname",
				Origin: cloudlets.Origin{
					OriginID:  "test_origin",
					Akamaized: true,
				},
			},
			{
				Hostname: "test-hostname",
				Origin: cloudlets.Origin{
					OriginID:  "test_origin_2",
					Akamaized: false,
				},
			},
		}
		client.On("ListOrigins", testutils.MockContext, mock.Anything).Return(origins, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/alb_update.tf", testDir),
						ExpectError: regexp.MustCompile("'liveness_hosts' field should be omitted for GTM hostname: \"test-hostname\". " +
							"Liveness tests for this host can be configured in DNS traffic management"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating origin with sum of percentages other than 100", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/percentage_validation"
		client := new(cloudlets.Mock)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/alb_create.tf", testDir),
						ExpectError: regexp.MustCompile("the total data center percentage must be 100%: total=10.012%"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}
