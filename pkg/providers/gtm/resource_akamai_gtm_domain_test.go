package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResGTMDomain(t *testing.T) {

	t.Run("create domain", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		dr := testCreateDomain
		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateDomainRequest"),
		).Return(&gtm.CreateDomainResponse{
			Resource: testDomain,
			Status:   testCreateDomain.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&dr, nil}
		})

		client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainRequest"),
		).Return(&testCreateDomain, nil).Times(3)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil)

		client.On("UpdateDomain",
			mock.Anything, // ctx is irrelevant for this test
			gtm.UpdateDomainRequest{
				Domain: testUpdateDomain,
				QueryArgs: &gtm.DomainQueryArgs{
					ContractID: "1-2ABCDEF",
					GroupID:    "123ABC",
				},
			},
		).Return(&updateDomainResponseStatus, nil)

		client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainRequest"),
		).Return(&testUpdateGetDomain, nil).Times(3)

		client.On("DeleteDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.DeleteDomainRequest"),
		).Return(&deleteDomainResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "false"),
							resource.TestCheckNoResourceAttr(resourceName, "sign_and_serve_algorithm"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "20"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "false"),
							resource.TestCheckNoResourceAttr(resourceName, "sign_and_serve_algorithm"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create domain with sign and serve", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		dr := testDomainWithSignAndServe
		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateDomainRequest"),
		).Return(&gtm.CreateDomainResponse{
			Resource: testDomain,
			Status:   testGetDomain.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&dr, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil)

		client.On("DeleteDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.DeleteDomainRequest"),
		).Return(&deleteDomainResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic_with_sign_and_serve.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "true"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve_algorithm", "RSA-SHA1"),
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
			mock.AnythingOfType("gtm.GetDomainRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		dr := testGetDomain
		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateDomainRequest"),
		).Return(&gtm.CreateDomainResponse{
			Resource: testDomain,
			Status:   testGetDomain.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&dr, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil)

		client.On("DeleteDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.DeleteDomainRequest"),
		).Return(&deleteDomainResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "false"),
							resource.TestCheckNoResourceAttr(resourceName, "sign_and_serve_algorithm"),
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
			mock.AnythingOfType("gtm.CreateDomainRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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

		dr := gtm.CreateDomainResponse{}
		dr.Resource = testDomain
		dr.Status = &deniedResponseStatus
		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateDomainRequest"),
		).Return(&dr, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
			mock.AnythingOfType("gtm.GetDomainRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		dr := testGetDomain
		client.On("CreateDomain",
			mock.Anything,
			mock.AnythingOfType("gtm.CreateDomainRequest"),
		).Return(&gtm.CreateDomainResponse{
			Resource: testDomain,
			Status:   testGetDomain.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&dr, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything,
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil).Times(2)

		client.On("DeleteDomain",
			mock.Anything,
			mock.AnythingOfType("gtm.DeleteDomainRequest"),
		).Return(&deleteDomainResponseStatus, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "false"),
							resource.TestCheckNoResourceAttr(resourceName, "sign_and_serve_algorithm"),
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

// getGTMDomainMocks mocks creation and deletion calls for gtm_domain resource
func getGTMDomainMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetDomain := client.On("GetDomain",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.GetDomainRequest"),
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	dr := domainWithOrderedEmails
	client.On("CreateDomain",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.CreateDomainRequest"),
	).Return(&gtm.CreateDomainResponse{
		Resource: domainWithOrderedEmailsDomain,
		Status:   domainWithOrderedEmails.Status,
	}, nil).Run(func(args mock.Arguments) {
		mockGetDomain.ReturnArguments = mock.Arguments{&dr, nil}
	})

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.GetDomainStatusRequest"),
	).Return(getDomainStatusResponseStatus, nil)

	client.On("DeleteDomain",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.DeleteDomainRequest"),
	).Return(&deleteDomainResponseStatus, nil)

	return client
}

var (
	// datacenters is gtm.Datacenter structure used in tests
	datacenters = []gtm.Datacenter{
		{
			City:                 "Snæfellsjökull",
			CloudServerTargeting: false,
			Continent:            "EU",
			Country:              "IS",
			DatacenterID:         3132,
			DefaultLoadObject: &gtm.LoadObject{
				LoadObject:     "",
				LoadObjectPort: 0,
				LoadServers:    make([]string, 0),
			},
			Latitude: 64.808,
			Links: []gtm.Link{
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

	// links is gtm.link structure used in tests
	links = []gtm.Link{
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
	properties = []gtm.Property{
		{
			BackupCName:            "",
			BackupIP:               "",
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
			IPv6:                   false,
			LastModified:           "2019-04-25T14:53:12.000+00:00",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/properties/test_property",
					Rel:  "self",
				},
			},
			LivenessTests: []gtm.LivenessTest{
				{
					DisableNonstandardPortWarning: false,
					HTTPError3xx:                  true,
					HTTPError4xx:                  true,
					HTTPError5xx:                  true,
					Name:                          "health check",
					RequestString:                 "",
					ResponseString:                "",
					SSLClientCertificate:          "",
					SSLClientPrivateKey:           "",
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
			TrafficTargets: []gtm.TrafficTarget{
				{
					DatacenterID: 3131,
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

	// testStatus is gtm.ResponseStatus structure used in tests
	testStatus = &gtm.ResponseStatus{
		ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
		Links: []gtm.Link{
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

	domainWithOrderedEmailsDomain = &gtm.Domain{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
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
		Status:                      testStatus,
		Type:                        "weighted",
	}

	// domainWithOrderedEmails is a gtm.Domain structure used in testing of email_notification_order list
	domainWithOrderedEmails = gtm.GetDomainResponse{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
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
		Status:                      testStatus,
		Type:                        "weighted",
	}

	testDomain = &gtm.Domain{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
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
		Status:                      testStatus,
		Type:                        "weighted",
	}

	testUpdateGetDomain = gtm.GetDomainResponse{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
		DefaultTimeoutPenalty:       25,
		EmailNotificationList:       make([]string, 0),
		LastModified:                "2019-04-25T14:53:12.000+00:00",
		LastModifiedBy:              "operator",
		Links:                       links,
		LoadFeedback:                false,
		LoadImbalancePercentage:     20.0,
		ModificationComments:        "Edit Property test_property",
		Name:                        gtmTestDomain,
		Properties:                  properties,
		Status:                      testStatus,
		Type:                        "weighted",
	}

	testUpdateDomain = &gtm.Domain{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
		DefaultTimeoutPenalty:       25,
		EmailNotificationList:       make([]string, 0),
		LastModified:                "2019-04-25T14:53:12.000+00:00",
		LastModifiedBy:              "operator",
		Links:                       links,
		LoadFeedback:                false,
		LoadImbalancePercentage:     20.0,
		ModificationComments:        "Edit Property test_property",
		Name:                        gtmTestDomain,
		Properties:                  properties,
		Status:                      testStatus,
		Type:                        "weighted",
	}

	testGetDomain = gtm.GetDomainResponse{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
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
		Status:                      testStatus,
		Type:                        "weighted",
	}

	testCreateDomain = gtm.GetDomainResponse{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
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
		Status:                      testStatus,
		Type:                        "weighted",
	}

	testDomainWithSignAndServe = gtm.GetDomainResponse{
		Datacenters:                 datacenters,
		DefaultErrorPenalty:         75,
		DefaultSSLClientCertificate: "",
		DefaultSSLClientPrivateKey:  "",
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
		Status:                      testStatus,
		Type:                        "weighted",
		SignAndServeAlgorithm:       ptr.To("RSA-SHA1"),
		SignAndServe:                true,
	}

	deniedResponseStatus = gtm.ResponseStatus{
		ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
		Links: []gtm.Link{
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
		ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
		Links: []gtm.Link{
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

	updateDomainResponseStatus = gtm.UpdateDomainResponse{
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
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
	}

	deleteDomainResponseStatus = gtm.DeleteDomainResponse{
		ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
		Links: []gtm.Link{
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
