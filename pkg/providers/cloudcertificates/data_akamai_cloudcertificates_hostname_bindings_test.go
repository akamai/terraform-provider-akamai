package cloudcertificates

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCloudCertificatesHostnameBindingsDataSource(t *testing.T) {
	t.Parallel()
	pageSize = 3

	mockListBindings := func(m *cloudcertificates.Mock, testData listBindingsTestData) {
		m.On("ListBindings", testutils.MockContext, testData.request).Return(&testData.response, nil).Times(3)
	}

	noFilteringStateChecker := test.NewStateChecker("data.akamai_cloudcertificates_hostname_bindings.test").
		CheckMissing("contract_id").
		CheckMissing("group_id").
		CheckMissing("domain").
		CheckMissing("network").
		CheckMissing("expiring_in_days").
		CheckEqual("bindings.#", "3").
		CheckEqual("bindings.0.certificate_id", "cert-123").
		CheckEqual("bindings.0.hostname", "example.com").
		CheckEqual("bindings.0.network", "PRODUCTION").
		CheckEqual("bindings.0.resource_type", "CDN_HOSTNAME").
		CheckEqual("bindings.1.certificate_id", "cert-456").
		CheckEqual("bindings.1.hostname", "test.example.com").
		CheckEqual("bindings.1.network", "STAGING").
		CheckEqual("bindings.1.resource_type", "CDN_HOSTNAME").
		CheckEqual("bindings.2.certificate_id", "cert-789").
		CheckEqual("bindings.2.hostname", "dev.example.com").
		CheckEqual("bindings.2.network", "PRODUCTION").
		CheckEqual("bindings.2.resource_type", "CDN_HOSTNAME")

	allFilteringStateChecker := test.NewStateChecker("data.akamai_cloudcertificates_hostname_bindings.test").
		CheckEqual("contract_id", "K-0N7RAK71").
		CheckEqual("group_id", "123456").
		CheckEqual("domain", "example.com").
		CheckEqual("network", "PRODUCTION").
		CheckEqual("expiring_in_days", "30").
		CheckEqual("bindings.#", "1").
		CheckEqual("bindings.0.certificate_id", "cert-123").
		CheckEqual("bindings.0.hostname", "example.com").
		CheckEqual("bindings.0.network", "PRODUCTION").
		CheckEqual("bindings.0.resource_type", "CDN_HOSTNAME")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"no filtering options": {
			init: func(m *cloudcertificates.Mock) {
				testData := listBindingsTestData{
					request: cloudcertificates.ListBindingsRequest{
						ExpiringInDays: nil,
						PageSize:       3,
						Page:           1,
					},
					response: cloudcertificates.ListBindingsResponse{
						Bindings: []cloudcertificates.CertificateBinding{
							{
								CertificateID: "cert-123",
								Hostname:      "example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
							{
								CertificateID: "cert-456",
								Hostname:      "test.example.com",
								Network:       "STAGING",
								ResourceType:  "CDN_HOSTNAME",
							},
							{
								CertificateID: "cert-789",
								Hostname:      "dev.example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
						},
						Links: cloudcertificates.Links{
							Next: nil,
						},
					},
				}
				mockListBindings(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesHostnameBindings/single.tf"),
					Check:  noFilteringStateChecker.Build(),
				},
			},
		},
		"no filtering options, but with API paging": {
			init: func(m *cloudcertificates.Mock) {
				// First page
				testData := listBindingsTestData{
					request: cloudcertificates.ListBindingsRequest{
						ExpiringInDays: nil,
						PageSize:       3,
						Page:           1,
					},
					response: cloudcertificates.ListBindingsResponse{
						Bindings: []cloudcertificates.CertificateBinding{
							{
								CertificateID: "cert-123",
								Hostname:      "example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
							{
								CertificateID: "cert-456",
								Hostname:      "test.example.com",
								Network:       "STAGING",
								ResourceType:  "CDN_HOSTNAME",
							},
							{
								CertificateID: "cert-789",
								Hostname:      "dev.example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
						},
						Links: cloudcertificates.Links{
							Next: ptr.To("/ccm/v2/bindings?page=2&pageSize=2"),
						},
					},
				}
				mockListBindings(m, testData)
				// Second page
				testData = listBindingsTestData{
					request: cloudcertificates.ListBindingsRequest{
						ExpiringInDays: nil,
						PageSize:       3,
						Page:           2,
					},
					response: cloudcertificates.ListBindingsResponse{
						Bindings: []cloudcertificates.CertificateBinding{
							{
								CertificateID: "cert-777",
								Hostname:      "foo.example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
						},
						Links: cloudcertificates.Links{
							Next: nil,
						},
					},
				}
				mockListBindings(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesHostnameBindings/single.tf"),
					Check: noFilteringStateChecker.
						CheckEqual("bindings.#", "4").Build(),
				},
			},
		},
		"single filtering option - domain": {
			init: func(m *cloudcertificates.Mock) {
				testData := listBindingsTestData{
					request: cloudcertificates.ListBindingsRequest{
						ExpiringInDays: nil,
						PageSize:       3,
						Page:           1,
					},
					response: cloudcertificates.ListBindingsResponse{
						Bindings: []cloudcertificates.CertificateBinding{
							{
								CertificateID: "cert-123",
								Hostname:      "example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
							{
								CertificateID: "cert-456",
								Hostname:      "test.example.com",
								Network:       "STAGING",
								ResourceType:  "CDN_HOSTNAME",
							},
							{
								CertificateID: "cert-789",
								Hostname:      "dev.example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
						},
						Links: cloudcertificates.Links{
							Next: nil,
						},
					},
				}
				mockListBindings(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesHostnameBindings/single.tf"),
					Check:  noFilteringStateChecker.Build(),
				},
			},
		},
		"all filtering options": {
			init: func(m *cloudcertificates.Mock) {
				testData := listBindingsTestData{
					request: cloudcertificates.ListBindingsRequest{
						ContractID:     "K-0N7RAK71",
						GroupID:        "123456",
						Domain:         "example.com",
						Network:        "PRODUCTION",
						ExpiringInDays: ptr.To(int64(30)),
						PageSize:       3,
						Page:           1,
					},
					response: cloudcertificates.ListBindingsResponse{
						Bindings: []cloudcertificates.CertificateBinding{
							{
								CertificateID: "cert-123",
								Hostname:      "example.com",
								Network:       "PRODUCTION",
								ResourceType:  "CDN_HOSTNAME",
							},
						},
						Links: cloudcertificates.Links{
							Next: nil,
						},
					},
				}
				mockListBindings(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesHostnameBindings/all.tf"),
					Check:  allFilteringStateChecker.Build(),
				},
			},
		},
		"all filtering options - no data returned": {
			init: func(m *cloudcertificates.Mock) {
				testData := listBindingsTestData{
					request: cloudcertificates.ListBindingsRequest{
						ContractID:     "K-0N7RAK71",
						GroupID:        "123456",
						Domain:         "example.com",
						Network:        "PRODUCTION",
						ExpiringInDays: ptr.To(int64(30)),
						PageSize:       3,
						Page:           1,
					},
					response: cloudcertificates.ListBindingsResponse{
						Bindings: []cloudcertificates.CertificateBinding{},
						Links: cloudcertificates.Links{
							Next: nil,
						},
					},
				}
				mockListBindings(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesHostnameBindings/all.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_hostname_bindings.test").
						CheckEqual("contract_id", "K-0N7RAK71").
						CheckEqual("group_id", "123456").
						CheckEqual("domain", "example.com").
						CheckEqual("network", "PRODUCTION").
						CheckEqual("expiring_in_days", "30").
						CheckEqual("bindings.#", "0").Build(),
				},
			},
		},
		"expect error: API returns error": {
			init: func(m *cloudcertificates.Mock) {
				testData := listBindingsTestData{
					request: cloudcertificates.ListBindingsRequest{
						ContractID:     "K-0N7RAK71",
						GroupID:        "123456",
						Domain:         "example.com",
						Network:        "PRODUCTION",
						ExpiringInDays: ptr.To(int64(30)),
						PageSize:       3,
						Page:           1,
					},
				}
				m.On("ListBindings", testutils.MockContext, testData.request).Return(nil, fmt.Errorf("not found")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesHostnameBindings/all.tf"),
					ExpectError: regexp.MustCompile("List bindings failed"),
				},
			},
		},
		"invalid network": {
			init: func(_ *cloudcertificates.Mock) {
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesHostnameBindings/invalid_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute network value must be one of: \["STAGING" "PRODUCTION"], got: "foo"`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &cloudcertificates.Mock{}
			if tc.init != nil {
				tc.init(client)
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

type listBindingsTestData struct {
	request  cloudcertificates.ListBindingsRequest
	response cloudcertificates.ListBindingsResponse
}
