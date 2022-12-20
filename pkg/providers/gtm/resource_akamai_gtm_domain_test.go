package gtm

import (
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

var gtmTestDomain = "gtm_terra_testdomain.akadns.net"
var contract = "1-2ABCDEF"
var group = "123ABC"

var dom = gtm.Domain{
	Datacenters: []*gtm.Datacenter{
		{
			City:                 "Snæfellsjökull",
			CloudServerTargeting: false,
			Continent:            "EU",
			Country:              "IS",
			DatacenterId:         3132,
			DefaultLoadObject: &gtm.LoadObject{
				LoadObject:     "",
				LoadObjectPort: 0,
				LoadServers:    make([]string, 0),
			},
			Latitude: 64.808,
			Links: []*gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3132",
					Rel:  "self",
				},
			},
			Longitude:       -23.776,
			Nickname:        "property_test_dc2",
			StateOrProvince: "",
			Virtual:         true,
		},
	},
	DefaultErrorPenalty:         75,
	DefaultSslClientCertificate: "",
	DefaultSslClientPrivateKey:  "",
	DefaultTimeoutPenalty:       25,
	EmailNotificationList:       make([]string, 0),
	LastModified:                "2019-04-25T14:53:12.000+00:00",
	LastModifiedBy:              "operator",
	Links: []*gtm.Link{
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net",
			Rel:  "self",
		},
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters",
			Rel:  "datacenters",
		},
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties",
			Rel:  "properties",
		},
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/geographic-maps",
			Rel:  "geographic-maps",
		},
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/cidr-maps",
			Rel:  "cidr-maps",
		},
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/resources",
			Rel:  "resources",
		},
	},
	LoadFeedback:            false,
	LoadImbalancePercentage: 10.0,
	ModificationComments:    "Edit Property test_property",
	Name:                    gtmTestDomain,
	Properties: []*gtm.Property{
		{
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
			Name:                      "test_property",
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
		},
	},
	Status: &gtm.ResponseStatus{
		ChangeId: "40e36abd-bfb2-4635-9fca-62175cf17007",
		Links: &[]gtm.Link{
			{
				Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
				Rel:  "self",
			},
		},
		Message:               "Current configuration has been propagated to all GTM nameservers",
		PassingValidation:     true,
		PropagationStatus:     "COMPLETE",
		PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
	},
	Type: "weighted",
}

var deniedResponseStatus = gtm.ResponseStatus{
	ChangeId: "40e36abd-bfb2-4635-9fca-62175cf17007",
	Links: &[]gtm.Link{
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
			Rel:  "self",
		},
	},
	Message:               "Request could not be completed. Invalid credentials.",
	PassingValidation:     true,
	PropagationStatus:     "DENIED",
	PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
}
var pendingResponseStatus = gtm.ResponseStatus{
	ChangeId: "40e36abd-bfb2-4635-9fca-62175cf17007",
	Links: &[]gtm.Link{
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
			Rel:  "self",
		},
	},
	Message:               "Current configuration has been propagated to all GTM nameservers",
	PassingValidation:     true,
	PropagationStatus:     "PENDING",
	PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
}
var completeResponseStatus = gtm.ResponseStatus{
	ChangeId: "40e36abd-bfb2-4635-9fca-62175cf17007",
	Links: &[]gtm.Link{
		{
			Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
			Rel:  "self",
		},
	},
	Message:               "Current configuration has been propagated to all GTM nameservers",
	PassingValidation:     true,
	PropagationStatus:     "COMPLETE",
	PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
}

func TestResGtmDomain(t *testing.T) {

	t.Run("create domain", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		dr := gtm.DomainResponse{}
		dr.Resource = &dom
		dr.Status = &pendingResponseStatus
		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
			mock.AnythingOfType("map[string]string"),
		).Return(&dr, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Domain), nil}
		})

		client.On("NewDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&dom)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
			mock.AnythingOfType("map[string]string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Domain), nil}
		})

		client.On("DeleteDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_domain.testdomain"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(dataSourceName, "load_imbalance_percentage", "10"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmDomain/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(dataSourceName, "load_imbalance_percentage", "20"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create domain failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
			mock.AnythingOfType("map[string]string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("NewDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&dom)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmDomain/create_basic.tf"),
						ExpectError: regexp.MustCompile("Domain Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create domain denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.DomainResponse{}
		dr.Resource = &dom
		dr.Status = &deniedResponseStatus
		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
			mock.AnythingOfType("map[string]string"),
		).Return(&dr, nil)

		client.On("NewDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&dom)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmDomain/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("import domain", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetDomain",
			mock.Anything,
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		dr := gtm.DomainResponse{}
		dr.Resource = &dom
		dr.Status = &pendingResponseStatus
		client.On("CreateDomain",
			mock.Anything,
			mock.AnythingOfType("*gtm.Domain"),
			mock.AnythingOfType("map[string]string"),
		).Return(&dr, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Domain), nil}
		})

		client.On("NewDomain",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&dom)

		client.On("GetDomainStatus",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Times(2)

		client.On("DeleteDomain",
			mock.Anything,
			mock.AnythingOfType("*gtm.Domain"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_domain.testdomain"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(dataSourceName, "load_imbalance_percentage", "10"),
						),
					},
					{
						Config:                  loadFixtureString("testdata/TestResGtmDomain/create_basic.tf"),
						ImportState:             true,
						ImportStateId:           gtmTestDomain,
						ResourceName:            dataSourceName,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"contract", "group"},
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

// Sets a Hack flag so cn work with existing Domains (only Admin can Delete)
func testAccPreCheckTF(_ *testing.T) {

	// by definition, we are running acceptance tests. ;-)
	log.Printf("[DEBUG] [Akamai GTMV1] Setting HashiAcc true")
	HashiAcc = true

}
