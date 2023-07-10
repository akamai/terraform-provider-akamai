package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

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
		dr.Resource = &domain
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
		).Return(&domain)

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

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "20"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create, update domain name - error", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		dr := gtm.DomainResponse{}
		dr.Resource = &domain
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
		).Return(&domain)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("DeleteDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
		).Return(&completeResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/update_domain_name.tf"),
						ExpectError: regexp.MustCompile("Error: once the domain is created, updating its name is not allowed"),
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
		).Return(&domain)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
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
		dr.Resource = &domain
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
		).Return(&domain)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
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
		dr.Resource = &domain
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
		).Return(&domain)

		client.On("GetDomainStatus",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Times(2)

		client.On("DeleteDomain",
			mock.Anything,
			mock.AnythingOfType("*gtm.Domain"),
		).Return(&completeResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
						),
					},
					{
						Config:                  testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						ImportState:             true,
						ImportStateId:           gtmTestDomain,
						ResourceName:            resourceName,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"contract", "group"},
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestGTMDomainOrder(t *testing.T) {
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reordered emails - no diff": {
			client:        getGTMDomainMocks(),
			pathForCreate: "testdata/TestResGtmDomain/order/email_notification_list/create.tf",
			pathForUpdate: "testdata/TestResGtmDomain/order/email_notification_list/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered emails and update comment - diff only for comment": {
			client:        getGTMDomainMocks(),
			pathForCreate: "testdata/TestResGtmDomain/order/email_notification_list/create.tf",
			pathForUpdate: "testdata/TestResGtmDomain/order/email_notification_list/reorder_and_update_comment.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useClient(test.client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
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

// getGTMDomainMocks mocks creation and deletion calls for gtm_domain resource
func getGTMDomainMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetDomain := client.On("GetDomain",
		mock.Anything, // ctx is irrelevant for this test
		gtmTestDomain,
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	dr := gtm.DomainResponse{}
	dr.Resource = &domainWithOrderedEmails
	dr.Status = &pendingResponseStatus
	client.On("CreateDomain",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.Domain"),
		mock.AnythingOfType("map[string]string"),
	).Return(&dr, nil).Run(func(args mock.Arguments) {
		mockGetDomain.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Domain), nil}
	})

	client.On("NewDomain",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(&domain)

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	client.On("DeleteDomain",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.Domain"),
	).Return(&completeResponseStatus, nil)

	return client
}

var (
	// datacenters is gtm.Datacenter structure used in tests
	datacenters = []*gtm.Datacenter{
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
	}

	// links is gtm.Link structure used in tests
	links = []*gtm.Link{
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
	}

	// properties is gtm.Property structure used in tests
	properties = []*gtm.Property{
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
	}

	// status is gtm.ResponseStatus structure used in tests
	status = &gtm.ResponseStatus{
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

	// domainWithOrderedEmails is a gtm.Domain structure used in testing of email_notification_order list
	domainWithOrderedEmails = gtm.Domain{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSslClientCertificate: "",
		DefaultSslClientPrivateKey:  "",
		DefaultTimeoutPenalty:       25,
		EmailNotificationList:       []string{"email1@nomail.com", "email2@nomail.com", "email3@nomail.com"},
		LastModified:                "2019-04-25T14:53:12.000+00:00",
		LastModifiedBy:              "operator",
		Links:                       links,
		LoadFeedback:                false,
		LoadImbalancePercentage:     10.0,
		ModificationComments:        "Edit Property test_property",
		Name:                        gtmTestDomain,
		Properties:                  properties,
		Status:                      status,
		Type:                        "weighted",
	}

	domain = gtm.Domain{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSslClientCertificate: "",
		DefaultSslClientPrivateKey:  "",
		DefaultTimeoutPenalty:       25,
		EmailNotificationList:       make([]string, 0),
		LastModified:                "2019-04-25T14:53:12.000+00:00",
		LastModifiedBy:              "operator",
		Links:                       links,
		LoadFeedback:                false,
		LoadImbalancePercentage:     10.0,
		ModificationComments:        "Edit Property test_property",
		Name:                        gtmTestDomain,
		Properties:                  properties,
		Status:                      status,
		Type:                        "weighted",
	}

	deniedResponseStatus = gtm.ResponseStatus{
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

	pendingResponseStatus = gtm.ResponseStatus{
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

	completeResponseStatus = gtm.ResponseStatus{
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

	gtmTestDomain = "gtm_terra_testdomain.akadns.net"
	resourceName  = "akamai_gtm_domain.testdomain"
)
