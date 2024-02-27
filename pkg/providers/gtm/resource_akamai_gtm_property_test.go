package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

var (
	prop = gtm.Property{
		BackupCName:            "",
		BackupIp:               "",
		BalanceByDownloadScore: false,
		CName:                  "www.boo.wow",
		Comments:               "",
		DynamicTTL:             300,
		FailbackDelay:          0,
		FailoverDelay:          0,
		HandoutMode:            "normal",
		HealthMax:              0,
		HealthMultiplier:       0,
		HealthThreshold:        0,
		Ipv6:                   false,
		LastModified:           "2019-04-25T14:53:12.000+00:00",
		Links: []*gtm.Link{
			{
				Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/test_property",
				Rel:  "self",
			},
		},
		LivenessTests: []*gtm.LivenessTest{
			{
				DisableNonstandardPortWarning: false,
				HttpError3xx:                  true,
				HttpError4xx:                  true,
				HttpError5xx:                  true,
				Name:                          "health check",
				RequestString:                 "",
				ResponseString:                "",
				SslClientCertificate:          "",
				SslClientPrivateKey:           "",
				TestInterval:                  60,
				TestObject:                    "/status",
				TestObjectPassword:            "",
				TestObjectPort:                80,
				TestObjectProtocol:            "HTTP",
				TestObjectUsername:            "",
				TestTimeout:                   25.0,
			},
		},
		LoadImbalancePercentage:   10.0,
		MapName:                   "",
		MaxUnreachablePenalty:     0,
		Name:                      "tfexample_prop_1",
		ScoreAggregationType:      "mean",
		StaticTTL:                 600,
		StickinessBonusConstant:   0,
		StickinessBonusPercentage: 50,
		TrafficTargets: []*gtm.TrafficTarget{
			{
				DatacenterId: 3131,
				Enabled:      true,
				HandoutCName: "",
				Servers: []string{
					"1.2.3.4",
					"1.2.3.5",
				},
				Weight: 50.0,
			},
		},
		Type:                 "weighted-round-robin",
		UnreachableThreshold: 0,
		UseComputedTargets:   false,
	}

	prop2 = gtm.Property{
		BackupCName:            "",
		BackupIp:               "",
		BalanceByDownloadScore: false,
		CName:                  "www.boo.wow",
		Comments:               "",
		DynamicTTL:             300,
		FailbackDelay:          0,
		FailoverDelay:          0,
		HandoutMode:            "normal",
		HealthMax:              0,
		HealthMultiplier:       0,
		HealthThreshold:        0,
		Ipv6:                   false,
		LastModified:           "2019-04-25T14:53:12.000+00:00",
		Links: []*gtm.Link{
			{
				Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/test_property",
				Rel:  "self",
			},
		},
		LivenessTests: []*gtm.LivenessTest{
			{
				DisableNonstandardPortWarning: false,
				HttpError3xx:                  true,
				HttpError4xx:                  true,
				HttpError5xx:                  true,
				Name:                          "health check",
				RequestString:                 "",
				ResponseString:                "",
				SslClientCertificate:          "",
				SslClientPrivateKey:           "",
				TestInterval:                  60,
				TestObject:                    "/status",
				TestObjectPassword:            "",
				TestObjectPort:                80,
				TestObjectProtocol:            "HTTP",
				TestObjectUsername:            "",
				TestTimeout:                   25.0,
			},
		},
		LoadImbalancePercentage:   10.0,
		MapName:                   "",
		MaxUnreachablePenalty:     0,
		Name:                      "tfexample_prop_1-updated",
		ScoreAggregationType:      "mean",
		StaticTTL:                 600,
		StickinessBonusConstant:   0,
		StickinessBonusPercentage: 50,
		TrafficTargets: []*gtm.TrafficTarget{
			{
				DatacenterId: 3131,
				Enabled:      true,
				HandoutCName: "",
				Servers: []string{
					"1.2.3.4",
					"1.2.3.5",
				},
				Weight: 50.0,
			},
		},
		Type:                 "weighted-round-robin",
		UnreachableThreshold: 0,
		UseComputedTargets:   false,
	}

	propertyResourceName = "akamai_gtm_property.tfexample_prop_1"
)

func TestResGtmProperty(t *testing.T) {

	t.Run("create property", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.PropertyResponse{}
		resp.Resource = &prop
		resp.Status = &pendingResponseStatus
		client.On("CreateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		})

		client.On("NewProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		})

		client.On("GetDomainStatus",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("NewTrafficTarget",
			mock.Anything,
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything,
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("float32"),
		)

		liveCall.RunFn = func(args mock.Arguments) {
			liveCall.ReturnArguments = mock.Arguments{
				&gtm.LivenessTest{
					Name:               args.String(1),
					TestObjectProtocol: args.String(2),
					TestInterval:       args.Int(3),
					TestTimeout:        args.Get(4).(float32),
				},
			}
		}

		client.On("UpdateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		})

		client.On("DeleteProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("NewProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		})

		client.On("NewTrafficTarget",
			mock.Anything,
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything,
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("float32"),
		)

		liveCall.RunFn = func(args mock.Arguments) {
			liveCall.ReturnArguments = mock.Arguments{
				&gtm.LivenessTest{
					Name:               args.String(1),
					TestObjectProtocol: args.String(2),
					TestInterval:       args.Int(3),
					TestTimeout:        args.Get(4).(float32),
				},
			}
		}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
						ExpectError: regexp.MustCompile("property Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.PropertyResponse{}
		dr.Resource = &prop
		dr.Status = &deniedResponseStatus
		client.On("CreateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("NewProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		})

		client.On("NewTrafficTarget",
			mock.Anything,
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything,
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("float32"),
		)

		liveCall.RunFn = func(args mock.Arguments) {
			liveCall.ReturnArguments = mock.Arguments{
				&gtm.LivenessTest{
					Name:               args.String(1),
					TestObjectProtocol: args.String(2),
					TestInterval:       args.Int(3),
					TestTimeout:        args.Get(4).(float32),
				},
			}
		}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property and update name - force new", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.PropertyResponse{}
		resp.Resource = &prop
		resp.Status = &pendingResponseStatus
		client.On("CreateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		}).Once()

		client.On("NewProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		}).Once()

		client.On("NewTrafficTarget",
			mock.Anything,
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything,
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("float32"),
		)

		liveCall.RunFn = func(args mock.Arguments) {
			liveCall.ReturnArguments = mock.Arguments{
				&gtm.LivenessTest{
					Name:               args.String(1),
					TestObjectProtocol: args.String(2),
					TestInterval:       args.Int(3),
					TestTimeout:        args.Get(4).(float32),
				},
			}
		}

		client.On("DeleteProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		// Create new property with updated name
		client.On("NewProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1-updated",
		}).Once()

		resp2 := gtm.PropertyResponse{
			Resource: &prop2,
			Status:   &pendingResponseStatus,
		}

		client.On("CreateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&resp2, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		}).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/update_name.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1-updated"),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("test_object_protocol different than HTTP, HTTPS or FTP", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.PropertyResponse{}
		resp.Resource = &prop
		resp.Status = &pendingResponseStatus
		client.On("CreateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		})

		client.On("NewProperty",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		})

		client.On("GetDomainStatus",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("NewTrafficTarget",
			mock.Anything,
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything,
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int"),
			mock.AnythingOfType("float32"),
		)

		liveCall.RunFn = func(args mock.Arguments) {
			liveCall.ReturnArguments = mock.Arguments{
				&gtm.LivenessTest{
					Name:               args.String(1),
					TestObjectProtocol: args.String(2),
					TestInterval:       args.Int(3),
					TestTimeout:        args.Get(4).(float32),
				},
			}
		}

		client.On("UpdateProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		})

		client.On("DeleteProperty",
			mock.Anything,
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_not_required.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(propertyResourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(propertyResourceName, "type", "weighted-round-robin"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property with test_object_protocol set to 'FTP' - test_object required error", func(t *testing.T) {
		client := &gtm.Mock{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_protocol_ftp.tf"),
						ExpectError: regexp.MustCompile(`Error: attribute 'test_object' is required when 'test_object_protocol' is set to 'HTTP', 'HTTPS' or 'FTP'`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property with test_object_protocol set to 'HTTP' - test_object required error", func(t *testing.T) {
		client := &gtm.Mock{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_protocol_http.tf"),
						ExpectError: regexp.MustCompile(`Error: attribute 'test_object' is required when 'test_object_protocol' is set to 'HTTP', 'HTTPS' or 'FTP'`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property with test_object_protocol set to 'HTTPS' - test_object required error", func(t *testing.T) {
		client := &gtm.Mock{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmProperty/test_object/test_object_protocol_https.tf"),
						ExpectError: regexp.MustCompile(`Error: attribute 'test_object' is required when 'test_object_protocol' is set to 'HTTP', 'HTTPS' or 'FTP'`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestResourceGTMTrafficTargetOrder(t *testing.T) {
	// To see actual plan when diff is expected, change 'nonEmptyPlan' to false in test case
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"second apply - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/create_multiple_traffic_targets.tf",
			pathForUpdate: "testdata/TestResGtmProperty/create_multiple_traffic_targets.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered traffic targets - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered traffic target with no datacenter_id - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/traffic_target/no_datacenter_id.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/no_datacenter_id_diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"added traffic target - diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/add_traffic_target.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"removed traffic target - diff (messy)": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/remove_traffic_target.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed 'enabled' field in traffic target - diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/change_enabled_field.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed 'enabled' field in re-ordered traffic target - diff (messy)": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/change_enabled_field_diff_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"re-ordered servers in traffic targets - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered servers and re-ordered traffic targets - no diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/traffic_target/diff_order.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"re-ordered and changed servers in traffic target - diff in one traffic target": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/changed_and_diff_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed servers - diff": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/change_server.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"changed servers and re-ordered traffic target - diff (messy)": {
			client:        getMocks(),
			pathForCreate: "testdata/TestResGtmProperty/multiple_servers.tf",
			pathForUpdate: "testdata/TestResGtmProperty/servers/change_server_and_diff_traffic_target_order.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useClient(test.client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.pathForCreate),
						},
						{
							Config:             testutils.LoadFixtureString(t, test.pathForUpdate),
							PlanOnly:           test.planOnly,
							ExpectNonEmptyPlan: test.nonEmptyPlan,
						},
					},
				})
			})
			test.client.AssertExpectations(t)
		})
	}
}

func getMocks() *gtm.Mock {
	client := &gtm.Mock{}

	// read
	getPropertyCall := client.On("GetProperty", mock.Anything, "tfexample_prop_1", "gtm_terra_testdomain.akadns.net").
		Return(nil, &gtm.Error{StatusCode: http.StatusNotFound})

	// create
	client.On("NewProperty", mock.Anything, "tfexample_prop_1").Return(&gtm.Property{
		Name: "tfexample_prop_1",
	}).Once()

	client.On("NewTrafficTarget", mock.Anything).Return(&gtm.TrafficTarget{}).Times(1)
	client.On("NewTrafficTarget", mock.Anything).Return(&gtm.TrafficTarget{}).Times(1)
	client.On("NewTrafficTarget", mock.Anything).Return(&gtm.TrafficTarget{}).Times(1)

	client.On("NewLivenessTest", mock.Anything, "lt5", "HTTP", 40, float32(30.0)).Return(&gtm.LivenessTest{
		Name:               "lt5",
		TestObjectProtocol: "HTTP",
		TestInterval:       40,
		TestTimeout:        30.0,
	}).Once()

	client.On("CreateProperty", mock.Anything, mock.AnythingOfType("*gtm.Property"), mock.AnythingOfType("string")).Return(&gtm.PropertyResponse{
		Resource: &prop,
		Status:   &pendingResponseStatus,
	}, nil).Run(func(args mock.Arguments) {
		getPropertyCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
	})

	// delete
	client.On("DeleteProperty",
		mock.Anything,
		mock.AnythingOfType("*gtm.Property"),
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	return client
}
