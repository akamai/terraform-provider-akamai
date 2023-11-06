package cloudlets

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
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
		expectCreateLoadBalancer = func(_ *testing.T, client *cloudlets.Mock, originID, originDescription, description, balancingType string, version int64, livenessHosts []string) (*cloudlets.Origin, *cloudlets.LoadBalancerVersion) {
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
						Latitude:        tools.Float64Ptr(102.78108),
						LivenessHosts:   livenessHosts,
						Longitude:       tools.Float64Ptr(-116.07064),
						OriginID:        "test_origin",
						Percent:         tools.Float64Ptr(100),
						StateOrProvince: tools.StringPtr("MA"),
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
						OriginID:  "test_origin",
						Akamaized: false,
					},
				},
			}
			client.On("ListOrigins", mock.Anything, mock.Anything).Return(origins, nil)

			client.On("CreateOrigin", mock.Anything, loadBalancerConfig).Return(&origin, nil).Once()
			client.On("CreateLoadBalancerVersion", mock.Anything, cloudlets.CreateLoadBalancerVersionRequest{
				OriginID:            originID,
				LoadBalancerVersion: loadBalancerVersionReq,
			}).Return(&loadBalancerVersionResp, nil).Once()

			return &origin, &loadBalancerVersionResp
		}

		expectReadLoadBalancer = func(_ *testing.T, client *cloudlets.Mock, origin *cloudlets.Origin, loadBalancerVersion *cloudlets.LoadBalancerVersion, times int) {
			client.On("GetOrigin", mock.Anything, cloudlets.GetOriginRequest{
				OriginID: origin.OriginID,
			}).Return(origin, nil).Times(times)
			client.On("GetLoadBalancerVersion", mock.Anything, cloudlets.GetLoadBalancerVersionRequest{
				OriginID:       origin.OriginID,
				Version:        loadBalancerVersion.Version,
				ShouldValidate: true,
			}).Return(loadBalancerVersion, nil).Times(times)
		}

		expectCreateLoadBalancerVersion = func(t *testing.T, client *cloudlets.Mock, originID string, loadBalancerVersion *cloudlets.LoadBalancerVersion, newBalancingType, description string) *cloudlets.LoadBalancerVersion {
			client.On("ListLoadBalancerActivations", mock.Anything, cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID}).Return([]cloudlets.LoadBalancerActivation{
				{
					OriginID: originID,
					Network:  cloudlets.LoadBalancerActivationNetworkProduction,
					Version:  loadBalancerVersion.Version,
					Status:   cloudlets.LoadBalancerActivationStatusActive,
				},
			}, nil)
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
			client.On("CreateLoadBalancerVersion", mock.Anything, cloudlets.CreateLoadBalancerVersionRequest{
				OriginID:            originID,
				LoadBalancerVersion: newVersionReq,
			}).Return(&newVersionResp, nil).Once()
			return &newVersionResp
		}

		expectOriginDescriptionUpdate = func(t *testing.T, client *cloudlets.Mock, origin *cloudlets.Origin, description string) *cloudlets.Origin {
			var newOrigin cloudlets.Origin
			err := copier.CopyWithOption(&newOrigin, origin, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			newOrigin.Description = description

			client.On("UpdateOrigin", mock.Anything, cloudlets.UpdateOriginRequest{
				OriginID: origin.OriginID,
				Description: cloudlets.Description{
					Description: description,
				},
			}).Return(&newOrigin, nil)

			return &newOrigin
		}

		expectUpdateLoadBalancerVersion = func(t *testing.T, client *cloudlets.Mock, originID string, loadBalancerVersion *cloudlets.LoadBalancerVersion, newBalancingType, description string) *cloudlets.LoadBalancerVersion {
			client.On("ListLoadBalancerActivations", mock.Anything, cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID}).Return(nil, nil)
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
			client.On("UpdateLoadBalancerVersion", mock.Anything, cloudlets.UpdateLoadBalancerVersionRequest{
				OriginID:            originID,
				ShouldValidate:      true,
				Version:             loadBalancerVersion.Version,
				LoadBalancerVersion: updateVersionReq,
			}).Return(&updateVersionResp, nil).Once()
			return &updateVersionResp
		}

		expectImportLoadBalancer = func(_ *testing.T, client *cloudlets.Mock, origin *cloudlets.Origin, numVersions int) {
			client.On("GetOrigin", mock.Anything, cloudlets.GetOriginRequest{OriginID: origin.OriginID}).Return(origin, nil).Once()

			var versionList []cloudlets.LoadBalancerVersion
			for i := 1; i <= numVersions; i++ {
				versionList = append(versionList, cloudlets.LoadBalancerVersion{OriginID: origin.OriginID, Version: int64(i)})
			}
			client.On("ListLoadBalancerVersions", mock.Anything, cloudlets.ListLoadBalancerVersionsRequest{
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		lbVersion = expectCreateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		lbVersion = expectUpdateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		lbVersion = expectUpdateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "", "test description updated")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		lbVersion = expectUpdateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE", "")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		// no liveness_settings for update
		var lbVersionUpdate cloudlets.LoadBalancerVersion
		err := copier.CopyWithOption(&lbVersionUpdate, lbVersion, copier.Option{DeepCopy: true})
		require.NoError(t, err)

		lbVersionUpdate.LivenessSettings = nil
		lbVersionUpdate.Description = "test description updated"

		lbVersionUpdate = *expectUpdateLoadBalancerVersion(t, client, origin.OriginID, &lbVersionUpdate, "PERFORMANCE", "")

		expectReadLoadBalancer(t, client, origin, &lbVersionUpdate, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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
					OriginID:  "test_origin",
					Akamaized: false,
				},
			},
		}
		client.On("ListOrigins", mock.Anything, mock.Anything).Return(origins, nil)

		client.On("CreateOrigin", mock.Anything, cloudlets.CreateOriginRequest{OriginID: "test_origin"}).Return(nil, fmt.Errorf("creating origin")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						ExpectError: regexp.MustCompile("creating origin"),
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
					OriginID:  "test_origin",
					Akamaized: false,
				},
			},
		}
		client.On("ListOrigins", mock.Anything, mock.Anything).Return(origins, nil)

		client.On("CreateOrigin", mock.Anything, cloudlets.CreateOriginRequest{OriginID: "test_origin"}).Return(&cloudlets.Origin{OriginID: "test_origin"}, nil).Once()

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
					Latitude:        tools.Float64Ptr(102.78108),
					LivenessHosts:   []string{"tf.test"},
					Longitude:       tools.Float64Ptr(-116.07064),
					OriginID:        "test_origin",
					Percent:         tools.Float64Ptr(100),
					StateOrProvince: tools.StringPtr("MA"),
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
		client.On("CreateLoadBalancerVersion", mock.Anything, cloudlets.CreateLoadBalancerVersionRequest{
			OriginID:            "test_origin",
			LoadBalancerVersion: loadBalancerVersionReq,
		}).Return(nil, fmt.Errorf("creating version")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		lbVersion = expectCreateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "origin description", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		origin = expectOriginDescriptionUpdate(t, client, origin, "update origin description")
		lbVersion = expectCreateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE", "test description updated")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:          "test_origin",
							version:           "1",
							originDescription: "origin description",
							description:       "test description",
							balancingType:     "WEIGHTED",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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

		origin, _ := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		client.On("GetOrigin", mock.Anything, cloudlets.GetOriginRequest{
			OriginID: "test_origin",
		}).Return(origin, nil).Once()

		client.On("GetLoadBalancerVersion", mock.Anything, cloudlets.GetLoadBalancerVersionRequest{
			OriginID:       "test_origin",
			Version:        1,
			ShouldValidate: true,
		}).Return(nil, fmt.Errorf("fetching version")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		expectImportLoadBalancer(t, client, origin, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		client.On("GetOrigin", mock.Anything, cloudlets.GetOriginRequest{OriginID: "not_existing_test_origin"}).Return(nil, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
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

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
					},
					{
						ImportState: true,
						ImportStateIdFunc: func(state *terraform.State) (string, error) {
							return "", nil
						},
						ResourceName: "akamai_cloudlets_application_load_balancer.alb",
						ExpectError:  regexp.MustCompile("origin id cannot be empty"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error importing load balancer no version found", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(cloudlets.Mock)

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{"tf.test"})

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		expectImportLoadBalancer(t, client, origin, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
					},
					{
						ImportState:   true,
						ImportStateId: "test_origin",
						ResourceName:  "akamai_cloudlets_application_load_balancer.alb",
						ExpectError:   regexp.MustCompile("no load balancer version found for origin_id: test_origin"),
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
		client.On("ListOrigins", mock.Anything, mock.Anything).Return(origins, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
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

		origins := []cloudlets.OriginResponse{
			{
				Hostname: "test-hostname",
				Origin: cloudlets.Origin{
					OriginID:  "test_origin",
					Akamaized: true,
				},
			},
		}
		client.On("ListOrigins", mock.Anything, mock.Anything).Return(origins, nil).Times(2)

		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "", "test description", "WEIGHTED", 1, []string{})

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_update.tf", testDir)),
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
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/alb_create.tf", testDir)),
						ExpectError: regexp.MustCompile("the total data center percentage must be 100%: total=10.012%"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}
