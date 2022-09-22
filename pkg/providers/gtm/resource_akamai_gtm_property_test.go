package gtm

import (
	"net/http"
	"regexp"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

var prop = gtm.Property{
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
			Name:         "",
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

func TestResGtmProperty(t *testing.T) {

	t.Run("create property", func(t *testing.T) {
		client := &mockgtm{}

		getCall := client.On("GetProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.PropertyResponse{}
		resp.Resource = &prop
		resp.Status = &pendingResponseStatus
		client.On("CreateProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		})

		client.On("NewProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("NewTrafficTarget",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything, // ctx is irrelevant for this test
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
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Property), nil}
		})

		client.On("DeleteProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_property.tfexample_prop_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmProperty/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted-round-robin"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmProperty/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_prop_1"),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted-round-robin"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property failed", func(t *testing.T) {
		client := &mockgtm{}

		client.On("CreateProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("NewProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		})

		client.On("NewTrafficTarget",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything, // ctx is irrelevant for this test
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
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmProperty/create_basic.tf"),
						ExpectError: regexp.MustCompile("property Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create property denied", func(t *testing.T) {
		client := &mockgtm{}

		dr := gtm.PropertyResponse{}
		dr.Resource = &prop
		dr.Status = &deniedResponseStatus
		client.On("CreateProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Property"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("NewProperty",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gtm.Property{
			Name: "tfexample_prop_1",
		})

		client.On("NewTrafficTarget",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&gtm.TrafficTarget{})

		client.On("NewStaticRRSet",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&gtm.StaticRRSet{})

		liveCall := client.On("NewLivenessTest",
			mock.Anything, // ctx is irrelevant for this test
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
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmProperty/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("no diff with re-ordered traffic targets", func(t *testing.T) {
		client := getMocks()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmProperty/create_multiple_traffic_targets.tf"),
					},
					{
						Config:   loadFixtureString("testdata/TestResGtmProperty/create_multiple_traffic_targets.tf"),
						PlanOnly: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("no diff with re-ordered traffic targets servers", func(t *testing.T) {
		client := getMocks()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmProperty/multiple_servers.tf"),
					},
					{
						Config:   loadFixtureString("testdata/TestResGtmProperty/multiple_servers_in_diff_order.tf"),
						PlanOnly: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func getMocks() *mockgtm {
	client := &mockgtm{}

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
