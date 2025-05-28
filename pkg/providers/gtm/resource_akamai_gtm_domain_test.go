package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const testDomainName = "gtm_terra_testdomain.akadns.net"

func TestResGTMDomain(t *testing.T) {
	const resourceName = "akamai_gtm_domain.testdomain"

	t.Run("create domain", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateDomain(client, getTestDomain(), &gtm.CreateDomainResponse{
			Resource: getReturnedTestDomain(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetDomain(client, getReturnedTestDomain(), nil, testutils.FourTimes)

		mockGetDomainStatus(client, testutils.Twice)

		mockUpdateDomain(client, &gtm.UpdateDomainResponse{Status: getDefaultResponseStatus()}, nil)

		mockGetDomain(client, getTestUpdateDomain(), nil, testutils.Twice)

		mockDeleteDomain(client, nil)

		mockDeleteDomainStatus(client, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "gtm_terra_testdomain.akadns.net"),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "false"),
							resource.TestCheckNoResourceAttr(resourceName, "sign_and_serve_algorithm"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "gtm_terra_testdomain.akadns.net"),
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

	t.Run("update domain failed", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateDomain(client, getTestDomain(), &gtm.CreateDomainResponse{
			Resource: getReturnedTestDomain(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetDomain(client, getReturnedTestDomain(), nil, testutils.FourTimes)

		mockGetDomainStatus(client, testutils.Once)

		mockUpdateDomain(client, nil, &gtm.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error updating domain",
			StatusCode: http.StatusInternalServerError,
		})

		mockDeleteDomain(client, nil)

		mockDeleteDomainStatus(client, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "gtm_terra_testdomain.akadns.net"),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "false"),
							resource.TestCheckNoResourceAttr(resourceName, "sign_and_serve_algorithm"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/update_basic.tf"),
						ExpectError: regexp.MustCompile("API error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create domain, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateDomain(client, getTestDomain(), &gtm.CreateDomainResponse{
			Resource: getReturnedTestDomain(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetDomain(client, getReturnedTestDomain(), nil, testutils.Twice)

		mockGetDomainStatus(client, testutils.Once)

		// Mock that the domain was deleted outside terraform
		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockDeleteDomain(client, nil)

		mockDeleteDomainStatus(client, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "gtm_terra_testdomain.akadns.net"),
							resource.TestCheckResourceAttr(resourceName, "type", "weighted"),
							resource.TestCheckResourceAttr(resourceName, "load_imbalance_percentage", "10"),
							resource.TestCheckResourceAttr(resourceName, "sign_and_serve", "false"),
							resource.TestCheckNoResourceAttr(resourceName, "sign_and_serve_algorithm"),
						),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						ExpectNonEmptyPlan: true,
						PlanOnly:           true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create domain with sign and serve", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateDomain(client, getTestDomainWithSignAndServe(), &gtm.CreateDomainResponse{
			Resource: getTestDomainWithSignAndServe(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetDomain(client, getTestDomainWithSignAndServe(), nil, testutils.Twice)

		mockGetDomainStatus(client, testutils.Once)

		mockDeleteDomain(client, nil)

		mockDeleteDomainStatus(client, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic_with_sign_and_serve.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "gtm_terra_testdomain.akadns.net"),
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

		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateDomain(client, getTestDomain(), &gtm.CreateDomainResponse{
			Resource: getReturnedTestDomain(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetDomain(client, getReturnedTestDomain(), nil, testutils.FourTimes)

		mockGetDomainStatus(client, testutils.Once)

		mockDeleteDomain(client, nil)

		mockDeleteDomainStatus(client, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "gtm_terra_testdomain.akadns.net"),
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

		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateDomain(client, getTestDomain(), nil, &gtm.Error{StatusCode: http.StatusBadRequest})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						ExpectError: regexp.MustCompile("Domain create error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create domain failed - domain already exists", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetDomain(client, getReturnedTestDomain(), nil, testutils.Once)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/create_basic.tf"),
						ExpectError: regexp.MustCompile("domain already exists error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create domain denied", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateDomain(client, getTestDomain(), &gtm.CreateDomainResponse{
			Resource: getReturnedTestDomain(),
			Status:   getDeniedResponseStatus(),
		}, nil)

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
}

func TestGTMDomainOrder(t *testing.T) {
	tests := map[string]struct {
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reordered emails - no diff": {
			pathForUpdate: "testdata/TestResGtmDomain/order/email_notification_list/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered emails and update comment - diff only for comment": {
			pathForUpdate: "testdata/TestResGtmDomain/order/email_notification_list/reorder_and_update_comment.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := getDomainOrderClient()
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/order/email_notification_list/create.tf"),
						},
						{
							Config:             testutils.LoadFixtureString(t, test.pathForUpdate),
							PlanOnly:           test.planOnly,
							ExpectNonEmptyPlan: test.nonEmptyPlan,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResGTMDomainImport(t *testing.T) {
	tests := map[string]struct {
		domainName  string
		init        func(*gtm.Mock)
		expectError *regexp.Regexp
		stateCheck  resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName: testDomainName,
			init: func(m *gtm.Mock) {
				// Read
				mockGetDomain(m, getImportedDomain(), nil, testutils.Once)
			},
			stateCheck: test.NewImportChecker().
				CheckEqual("name", "gtm_terra_testdomain.akadns.net").
				CheckEqual("type", "weighted").
				CheckEqual("default_unreachable_threshold", "5").
				CheckEqual("email_notification_list.0", "email1@nomail.com").
				CheckEqual("email_notification_list.1", "email2@nomail.com").
				CheckEqual("min_pingable_region_fraction", "1").
				CheckEqual("default_timeout_penalty", "25").
				CheckEqual("servermonitor_liveness_count", "1").
				CheckEqual("round_robin_prefix", "test prefix").
				CheckEqual("servermonitor_load_count", "1").
				CheckEqual("ping_interval", "10").
				CheckEqual("max_ttl", "10").
				CheckEqual("load_imbalance_percentage", "10").
				CheckEqual("default_health_max", "100").
				CheckEqual("map_update_interval", "123").
				CheckEqual("max_properties", "123").
				CheckEqual("max_resources", "123").
				CheckEqual("default_ssl_client_private_key", "test key").
				CheckEqual("default_error_penalty", "75").
				CheckEqual("max_test_timeout", "123").
				CheckEqual("cname_coalescing_enabled", "true").
				CheckEqual("default_health_multiplier", "10").
				CheckEqual("servermonitor_pool", "test pool").
				CheckEqual("load_feedback", "false").
				CheckEqual("min_ttl", "5").
				CheckEqual("default_max_unreachable_penalty", "123").
				CheckEqual("default_health_threshold", "5").
				CheckEqual("comment", "test comment").
				CheckEqual("min_test_interval", "10").
				CheckEqual("ping_packet_size", "5").
				CheckEqual("default_ssl_client_certificate", "test ssl").
				CheckEqual("sign_and_serve", "true").
				CheckEqual("sign_and_serve_algorithm", "RSA-SHA1").
				CheckEqual("end_user_mapping_enabled", "false").Build(),
		},
		"expect error - read": {
			domainName: testDomainName,
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetDomain(m, nil, fmt.Errorf("get failed"), testutils.Once)
			},
			expectError: regexp.MustCompile(`get failed`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			tc.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: tc.stateCheck,
							ImportStateId:    tc.domainName,
							ImportState:      true,
							ResourceName:     "akamai_gtm_domain.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResGtmDomain/import_basic.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// getDomainOrderClient mocks creation and deletion calls for gtm_domain resource
func getDomainOrderClient() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetDomain(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

	mockCreateDomain(client, getTestDomainWithNotifications(), &gtm.CreateDomainResponse{
		Resource: getTestDomainWithNotifications(),
		Status:   getDefaultResponseStatus(),
	}, nil)

	mockGetDomain(client, getTestDomainWithNotifications(), nil, testutils.ThreeTimes)

	mockGetDomainStatus(client, testutils.Once)

	mockDeleteDomain(client, nil)

	mockDeleteDomainStatus(client, nil)

	return client
}

func mockUpdateDomain(client *gtm.Mock, resp *gtm.UpdateDomainResponse, err error) *mock.Call {
	return client.On("UpdateDomain",
		testutils.MockContext,
		gtm.UpdateDomainRequest{
			Domain: getTestUpdateDomain(),
			QueryArgs: &gtm.DomainQueryArgs{
				ContractID: "1-2ABCDEF",
				GroupID:    "123ABC",
			},
		},
	).Return(resp, err).Once()
}

func mockCreateDomain(client *gtm.Mock, domain *gtm.Domain, resp *gtm.CreateDomainResponse, err error) {
	client.On("CreateDomain",
		testutils.MockContext,
		gtm.CreateDomainRequest{
			Domain: domain,
			QueryArgs: &gtm.DomainQueryArgs{
				ContractID: "1-2ABCDEF",
				GroupID:    "123ABC",
			},
		},
	).Return(resp, err).Once()
}

func mockDeleteDomain(client *gtm.Mock, err error) *mock.Call {
	resp := gtm.DeleteDomainsResponse{
		RequestID:      "e585a640-0849-4b87-8dd9-91afdaf8851c",
		ExpirationDate: "2050-01-03T12:00:00Z",
	}

	return client.On("DeleteDomains",
		testutils.MockContext,
		gtm.DeleteDomainsRequest{
			Body: gtm.DeleteDomainsRequestBody{
				DomainNames: []string{
					testDomainName,
				},
			},
		},
	).Return(&resp, err).Once()
}

func mockDeleteDomainStatus(client *gtm.Mock, err error) *mock.Call {
	resp := gtm.DeleteDomainsStatusResponse{
		DomainsSubmitted: 1,
		SuccessCount:     1,
		FailureCount:     0,
		IsComplete:       true,
		RequestID:        "e585a640-0849-4b87-8dd9-91afdaf8851c",
		ExpirationDate:   "2050-01-03T12:00:00Z",
	}
	if err != nil {
		resp.SuccessCount = 0
		resp.FailureCount = 1
	}

	return client.On("GetDeleteDomainsStatus",
		testutils.MockContext,
		gtm.DeleteDomainsStatusRequest{
			RequestID: "e585a640-0849-4b87-8dd9-91afdaf8851c",
		},
	).Return(&resp, err).Once()
}

func mockGetDomain(m *gtm.Mock, domain *gtm.Domain, err error, times int) *mock.Call {
	var resp *gtm.GetDomainResponse
	if domain != nil {
		r := gtm.GetDomainResponse(*domain)
		resp = &r
	}
	return m.On("GetDomain", testutils.MockContext, gtm.GetDomainRequest{
		DomainName: testDomainName,
	}).Return(resp, err).Times(times)
}

func getImportedDomain() *gtm.Domain {
	return &gtm.Domain{
		Name:                         testDomainName,
		Type:                         "weighted",
		DefaultUnreachableThreshold:  5.0,
		EmailNotificationList:        []string{"email1@nomail.com", "email2@nomail.com"},
		MinPingableRegionFraction:    1,
		DefaultTimeoutPenalty:        25,
		ServermonitorLivenessCount:   1,
		RoundRobinPrefix:             "test prefix",
		ServermonitorLoadCount:       1,
		PingInterval:                 10,
		MaxTTL:                       10,
		LoadImbalancePercentage:      10.0,
		DefaultHealthMax:             100,
		MapUpdateInterval:            123,
		MaxProperties:                123,
		MaxResources:                 123,
		DefaultSSLClientPrivateKey:   "test key",
		DefaultErrorPenalty:          75,
		MaxTestTimeout:               123,
		CNameCoalescingEnabled:       true,
		DefaultHealthMultiplier:      10,
		ServermonitorPool:            "test pool",
		LoadFeedback:                 false,
		MinTTL:                       5,
		DefaultMaxUnreachablePenalty: 123,
		DefaultHealthThreshold:       5,
		ModificationComments:         "test comment",
		MinTestInterval:              10,
		PingPacketSize:               5,
		DefaultSSLClientCertificate:  "test ssl",
		EndUserMappingEnabled:        false,
		SignAndServe:                 true,
		SignAndServeAlgorithm:        ptr.To("RSA-SHA1"),
	}
}

func getTestDomain() *gtm.Domain {
	return &gtm.Domain{
		Name:                    testDomainName,
		Type:                    "weighted",
		LoadImbalancePercentage: 10,
		DefaultErrorPenalty:     75,
		DefaultTimeoutPenalty:   25,
		ModificationComments:    "Edit Property test_property",
	}
}

func getReturnedTestDomain() *gtm.Domain {
	domain := getTestDomain()
	domain.Datacenters = getTestDatacenters()
	domain.Properties = getDomainTestProperties()
	domain.Links = getTestDomainLinks()
	domain.Status = getDefaultResponseStatus()
	domain.EmailNotificationList = make([]string, 0)
	domain.LastModifiedBy = "operator"
	domain.LastModified = "2019-04-25T14:53:12.000+00:00"
	return domain
}

func getTestDomainWithSignAndServe() *gtm.Domain {
	domain := getTestDomain()
	domain.SignAndServe = true
	domain.SignAndServeAlgorithm = ptr.To("RSA-SHA1")
	return domain
}

func getTestDomainWithNotifications() *gtm.Domain {
	domain := getTestDomain()
	domain.EmailNotificationList = []string{"email1@nomail.com", "email3@nomail.com", "email2@nomail.com"}
	return domain
}

func getTestUpdateDomain() *gtm.Domain {
	domain := getTestDomain()
	domain.LoadImbalancePercentage = 20
	domain.Datacenters = getTestDatacenters()
	domain.Properties = getDomainTestProperties()
	domain.Links = getTestDomainLinks()
	domain.Status = getDefaultResponseStatus()
	domain.EmailNotificationList = make([]string, 0)
	domain.LastModifiedBy = "operator"
	domain.LastModified = "2019-04-25T14:53:12.000+00:00"
	return domain
}

func getTestDatacenters() []gtm.Datacenter {
	return []gtm.Datacenter{
		{
			City:                 "Snæfellsjökull",
			CloudServerTargeting: false,
			Continent:            "EU",
			Country:              "IS",
			DatacenterID:         datacenterID3132,
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
}

func getTestDomainLinks() []gtm.Link {
	return []gtm.Link{
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
}

func getDomainTestProperties() []gtm.Property {
	return []gtm.Property{
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
					DatacenterID: datacenterID3131,
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
}

func getDeniedResponseStatus() *gtm.ResponseStatus {
	return &gtm.ResponseStatus{
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
}

func getPendingResponseStatus() *gtm.ResponseStatus {
	return &gtm.ResponseStatus{
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
}
