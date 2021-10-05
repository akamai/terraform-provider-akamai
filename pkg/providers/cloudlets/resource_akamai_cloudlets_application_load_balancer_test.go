package cloudlets

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceApplicationLoadBalancer(t *testing.T) {

	type loadBalancerAttributes struct {
		originID, version, description, balancingType string
	}

	var (
		expectCreateLoadBalancer = func(_ *testing.T, client *mockcloudlets, originID, description, balancingType string, version int64) (*cloudlets.Origin, *cloudlets.LoadBalancerVersion) {
			loadBalancerConfig := cloudlets.LoadBalancerOriginCreateRequest{
				OriginID:    originID,
				Description: cloudlets.Description{description},
			}
			loadBalancerVersionReq := cloudlets.LoadBalancerVersion{
				BalancingType: cloudlets.BalancingType(balancingType),
				DataCenters: []cloudlets.DataCenter{
					{
						City:            "Boston",
						CloudService:    true,
						Continent:       "NA",
						Country:         "US",
						Hostname:        "test-hostname",
						Latitude:        102.78108,
						LivenessHosts:   []string{"tf.test"},
						Longitude:       -116.07064,
						OriginID:        "test_origin",
						Percent:         10,
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
				Description: description,
				Type:        "APPLICATION_LOAD_BALANCER",
			}
			client.On("CreateOrigin", mock.Anything, loadBalancerConfig).Return(&origin, nil).Once()
			client.On("CreateLoadBalancerVersion", mock.Anything, cloudlets.CreateLoadBalancerVersionRequest{
				OriginID:            originID,
				LoadBalancerVersion: loadBalancerVersionReq,
			}).Return(&loadBalancerVersionResp, nil).Once()

			return &origin, &loadBalancerVersionResp
		}

		expectReadLoadBalancer = func(_ *testing.T, client *mockcloudlets, origin *cloudlets.Origin, loadBalancerVersion *cloudlets.LoadBalancerVersion, times int) {
			client.On("GetOrigin", mock.Anything, origin.OriginID).Return(origin, nil).Times(times)
			client.On("GetLoadBalancerVersion", mock.Anything, cloudlets.GetLoadBalancerVersionRequest{
				OriginID:       origin.OriginID,
				Version:        loadBalancerVersion.Version,
				ShouldValidate: true,
			}).Return(loadBalancerVersion, nil).Times(times)
		}

		expectUpdateOrigin = func(t *testing.T, client *mockcloudlets, origin *cloudlets.Origin, updatedDescription string) *cloudlets.Origin {
			var updatedOrigin cloudlets.Origin
			err := copier.CopyWithOption(&updatedOrigin, origin, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			updatedOrigin.Description = updatedDescription
			client.On("UpdateOrigin", mock.Anything, cloudlets.LoadBalancerOriginUpdateRequest{
				OriginID:    origin.OriginID,
				Description: cloudlets.Description{updatedDescription},
			}).Return(&updatedOrigin, nil).Once()
			return &updatedOrigin
		}

		expectCreateLoadBalancerVersion = func(t *testing.T, client *mockcloudlets, originID string, loadBalancerVersion *cloudlets.LoadBalancerVersion, newBalancingType string) *cloudlets.LoadBalancerVersion {
			client.On("GetLoadBalancerActivations", mock.Anything, originID).Return(cloudlets.ActivationsList{
				cloudlets.ActivationResponse{
					OriginID: originID,
					Network:  cloudlets.ActivationNetworkProd,
					Version:  loadBalancerVersion.Version,
					Status:   cloudlets.ActivationStatusActive,
				},
			}, nil)
			var newVersionReq, newVersionResp cloudlets.LoadBalancerVersion
			err := copier.CopyWithOption(&newVersionReq, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			err = copier.CopyWithOption(&newVersionResp, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			newVersionReq.Version = 0
			newVersionReq.Warnings = nil
			newVersionReq.BalancingType = cloudlets.BalancingType(newBalancingType)
			newVersionResp.BalancingType = cloudlets.BalancingType(newBalancingType)
			newVersionResp.Version++
			client.On("CreateLoadBalancerVersion", mock.Anything, cloudlets.CreateLoadBalancerVersionRequest{
				OriginID:            originID,
				LoadBalancerVersion: newVersionReq,
			}).Return(&newVersionResp, nil).Once()
			return &newVersionResp
		}

		expectUpdateLoadBalancerVersion = func(t *testing.T, client *mockcloudlets, originID string, loadBalancerVersion *cloudlets.LoadBalancerVersion, newBalancingType string) *cloudlets.LoadBalancerVersion {
			client.On("GetLoadBalancerActivations", mock.Anything, originID).Return(nil, nil)
			var updateVersionReq, updateVersionResp cloudlets.LoadBalancerVersion
			err := copier.CopyWithOption(&updateVersionReq, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			err = copier.CopyWithOption(&updateVersionResp, loadBalancerVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)

			updateVersionReq.Version = 0
			updateVersionReq.Warnings = nil
			updateVersionReq.BalancingType = cloudlets.BalancingTypePerformance
			updateVersionResp.BalancingType = cloudlets.BalancingType(newBalancingType)
			client.On("UpdateLoadBalancerVersion", mock.Anything, cloudlets.UpdateLoadBalancerVersionRequest{
				OriginID:            originID,
				ShouldValidate:      true,
				Version:             loadBalancerVersion.Version,
				LoadBalancerVersion: updateVersionReq,
			}).Return(&updateVersionResp, nil).Once()
			return &updateVersionResp
		}

		expectImportLoadBalancer = func(_ *testing.T, client *mockcloudlets, origin *cloudlets.Origin, numVersions int) {
			client.On("GetOrigin", mock.Anything, origin.OriginID).Return(origin, nil).Once()

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
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "description", attrs.description),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "balancing_type", attrs.balancingType),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.#", "1"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.cloud_server_host_header_override", "false"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.cloud_service", "true"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.country", "US"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.continent", "NA"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.latitude", "102.78108"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.longitude", "-116.07064"),
				resource.TestCheckResourceAttr("akamai_cloudlets_application_load_balancer.alb", "data_centers.0.percent", "10"),
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
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		origin = expectUpdateOrigin(t, client, origin, "test description updated")
		lbVersion = expectCreateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_update.tf", testDir)),
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
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		origin = expectUpdateOrigin(t, client, origin, "test description updated")
		lbVersion = expectUpdateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_update.tf", testDir)),
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

	t.Run("update only origin", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle_origin_update"
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		origin = expectUpdateOrigin(t, client, origin, "test description updated")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_update.tf", testDir)),
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
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		lbVersion = expectUpdateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_update.tf", testDir)),
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

	t.Run("attempt creating existing origin", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		origin := &cloudlets.Origin{
			OriginID:    "test_origin",
			Description: "some other description",
			Type:        "APPLICATION_LOAD_BALANCER",
		}
		lbVersion := &cloudlets.LoadBalancerVersion{
			BalancingType: "WEIGHTED",
			DataCenters: []cloudlets.DataCenter{
				{
					City:            "Boston",
					CloudService:    true,
					Continent:       "NA",
					Country:         "US",
					Hostname:        "test-hostname",
					Latitude:        102.78108,
					LivenessHosts:   []string{"tf.test"},
					Longitude:       -116.07064,
					OriginID:        "test_origin",
					Percent:         10,
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
			Version: 1,
			Warnings: []cloudlets.Warning{
				{
					Detail:      "test warning details",
					JSONPointer: "/path",
					Title:       "test warning",
					Type:        "test type",
				},
			},
		}
		client.On("GetOrigin", mock.Anything, "test_origin").Return(origin, nil).Once()
		client.On("ListLoadBalancerVersions", mock.Anything, cloudlets.ListLoadBalancerVersionsRequest{
			OriginID: origin.OriginID,
		}).Return([]cloudlets.LoadBalancerVersion{*lbVersion}, nil).Once()

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		origin = expectUpdateOrigin(t, client, origin, "test description")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "some other description",
							balancingType: "WEIGHTED",
						}),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("attempt creating existing origin without existing version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		origin := &cloudlets.Origin{
			OriginID:    "test_origin",
			Description: "test description",
			Type:        "APPLICATION_LOAD_BALANCER",
		}
		lbVersion := &cloudlets.LoadBalancerVersion{
			BalancingType: "WEIGHTED",
			DataCenters: []cloudlets.DataCenter{
				{
					City:            "Boston",
					CloudService:    true,
					Continent:       "NA",
					Country:         "US",
					Hostname:        "test-hostname",
					Latitude:        102.78108,
					LivenessHosts:   []string{"tf.test"},
					Longitude:       -116.07064,
					OriginID:        "test_origin",
					Percent:         10,
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
			Version: 1,
			Warnings: []cloudlets.Warning{
				{
					Detail:      "test warning details",
					JSONPointer: "/path",
					Title:       "test warning",
					Type:        "test type",
				},
			},
		}
		client.On("GetOrigin", mock.Anything, origin.OriginID).Return(origin, nil).Once()
		client.On("ListLoadBalancerVersions", mock.Anything, cloudlets.ListLoadBalancerVersionsRequest{
			OriginID: origin.OriginID,
		}).Return(nil, nil).Once()

		client.On("GetOrigin", mock.Anything, origin.OriginID).Return(origin, nil).Times(3)

		var lbVersionReq cloudlets.LoadBalancerVersion
		err := copier.CopyWithOption(&lbVersionReq, lbVersion, copier.Option{DeepCopy: true})
		require.NoError(t, err)

		lbVersionReq.Version = 0
		lbVersionReq.Warnings = nil
		client.On("CreateLoadBalancerVersion", mock.Anything, cloudlets.CreateLoadBalancerVersionRequest{
			OriginID:            origin.OriginID,
			LoadBalancerVersion: lbVersionReq,
		}).Return(lbVersion, nil).Once()

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:    "test_origin",
							description: "test description",
						}),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating origin", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		client.On("CreateOrigin", mock.Anything, cloudlets.LoadBalancerOriginCreateRequest{
			OriginID:    "test_origin",
			Description: cloudlets.Description{"test description"},
		}).Return(nil, fmt.Errorf("creating origin")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						ExpectError: regexp.MustCompile("creating origin"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		client.On("CreateOrigin", mock.Anything, cloudlets.LoadBalancerOriginCreateRequest{
			OriginID:    "test_origin",
			Description: cloudlets.Description{"test description"},
		}).Return(&cloudlets.Origin{OriginID: "test_origin"}, nil).Once()

		loadBalancerVersionReq := cloudlets.LoadBalancerVersion{
			BalancingType: "WEIGHTED",
			DataCenters: []cloudlets.DataCenter{
				{
					City:            "Boston",
					CloudService:    true,
					Continent:       "NA",
					Country:         "US",
					Hostname:        "test-hostname",
					Latitude:        102.78108,
					LivenessHosts:   []string{"tf.test"},
					Longitude:       -116.07064,
					OriginID:        "test_origin",
					Percent:         10,
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						ExpectError: regexp.MustCompile("creating version"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("load balancer lifecycle with create new version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		origin = expectUpdateOrigin(t, client, origin, "test description updated")
		lbVersion = expectCreateLoadBalancerVersion(t, client, origin.OriginID, lbVersion, "PERFORMANCE")

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						Check: checkAttributes(loadBalancerAttributes{
							originID:      "test_origin",
							version:       "1",
							description:   "test description",
							balancingType: "WEIGHTED",
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_update.tf", testDir)),
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

	t.Run("error fetching origin", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		_, _ = expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, fmt.Errorf("fetching origin")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						ExpectError: regexp.MustCompile("fetching origin"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error fetching version", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		_, _ = expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(&cloudlets.Origin{OriginID: "test_origin"}, nil).Once()
		client.On("GetLoadBalancerVersion", mock.Anything, cloudlets.GetLoadBalancerVersionRequest{
			OriginID:       "test_origin",
			Version:        1,
			ShouldValidate: true,
		}).Return(nil, fmt.Errorf("fetching version")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
						ExpectError: regexp.MustCompile("fetching version"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("import load balancer", func(t *testing.T) {
		testDir := "testdata/TestResLoadBalancerConfig/lifecycle"
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 3)

		expectImportLoadBalancer(t, client, origin, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
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
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		client.On("GetOrigin", mock.Anything, "not_existing_test_origin").Return(nil, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
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
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
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
		client := new(mockcloudlets)

		client.On("GetOrigin", mock.Anything, "test_origin").Return(nil, &cloudlets.Error{StatusCode: http.StatusNotFound}).Once()
		origin, lbVersion := expectCreateLoadBalancer(t, client, "test_origin", "test description", "WEIGHTED", 1)

		expectReadLoadBalancer(t, client, origin, lbVersion, 2)

		expectImportLoadBalancer(t, client, origin, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/alb_create.tf", testDir)),
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
}
