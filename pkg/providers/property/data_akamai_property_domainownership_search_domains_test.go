package property

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDomainOwnershipSearchDomains(t *testing.T) {
	t.Parallel()
	testDir := "testdata/TestDSDomainOwnershipSearchDomains/"

	mockSearchDomains := func(client *domainownership.Mock, req domainownership.SearchDomainsRequest, resp *domainownership.SearchDomainsResponse) *mock.Call {
		return client.On("SearchDomains", testutils.MockContext, req).Return(resp, nil)
	}

	defaultStateChecker := test.NewStateChecker("data.akamai_property_domainownership_search_domains.test").
		CheckEqual("domains.#", "3").
		CheckEqual("domains.0.domain_name", "dom1.test").
		CheckEqual("domains.0.validation_scope", "HOST").
		CheckEqual("domains.0.account_id", "1-ACCOUN").
		CheckEqual("domains.0.domain_status", "REQUEST_ACCEPTED").
		CheckEqual("domains.0.validation_level", "FQDN").
		CheckEqual("domains.0.validation_method", "DNS_TXT").
		CheckEqual("domains.0.validation_requested_by", "someuser").
		CheckEqual("domains.0.validation_requested_date", "2025-08-04T13:27:19Z").
		CheckMissing("domains.0.validation_completed_date").
		CheckEqual("domains.0.validation_challenge.cname_record.name", "cname-name-1").
		CheckEqual("domains.0.validation_challenge.cname_record.target", "cname-target-1").
		CheckEqual("domains.0.validation_challenge.txt_record.name", "txt-name-1").
		CheckEqual("domains.0.validation_challenge.txt_record.value", "txt-value-1").
		CheckEqual("domains.0.validation_challenge.http_file.path", "http-file-path-1").
		CheckEqual("domains.0.validation_challenge.http_file.content", "http-file-content-1").
		CheckEqual("domains.0.validation_challenge.http_file.content_type", "text/plain").
		CheckEqual("domains.0.validation_challenge.http_redirect.from", "http-redirect-from-1").
		CheckEqual("domains.0.validation_challenge.http_redirect.to", "http-redirect-to-1").
		CheckEqual("domains.0.validation_challenge.expiration_date", "2025-08-05T13:27:19Z").
		CheckEqual("domains.1.domain_name", "dom2.test").
		CheckEqual("domains.1.validation_scope", "HOST").
		CheckEqual("domains.1.account_id", "1-ACCOUN").
		CheckEqual("domains.1.domain_status", "VALIDATED").
		CheckEqual("domains.1.validation_level", "FQDN").
		CheckEqual("domains.1.validation_method", "SYSTEM").
		CheckEqual("domains.1.validation_requested_by", "someuser").
		CheckEqual("domains.1.validation_requested_date", "2025-08-04T13:27:19Z").
		CheckEqual("domains.1.validation_completed_date", "2025-08-05T11:56:07Z").
		CheckMissing("domains.1.validation_challenge.cname_record.name").
		CheckMissing("domains.1.validation_challenge.cname_record.target").
		CheckMissing("domains.1.validation_challenge.txt_record.name").
		CheckMissing("domains.1.validation_challenge.txt_record.value").
		CheckMissing("domains.1.validation_challenge.http_file.path").
		CheckMissing("domains.1.validation_challenge.http_file.content").
		CheckMissing("domains.1.validation_challenge.http_file.content_type").
		CheckMissing("domains.1.validation_challenge.http_redirect.from").
		CheckMissing("domains.1.validation_challenge.http_redirect.to").
		CheckMissing("domains.1.validation_challenge.expiration_date").
		CheckEqual("domains.2.domain_name", "dom3.test").
		CheckEqual("domains.2.validation_scope", "HOST").
		CheckEqual("domains.2.account_id", "1-ACCOUN").
		CheckEqual("domains.2.domain_status", "VALIDATED").
		CheckEqual("domains.2.validation_level", "FQDN").
		CheckEqual("domains.2.validation_method", "DNS_TXT").
		CheckEqual("domains.2.validation_requested_by", "someuser").
		CheckEqual("domains.2.validation_requested_date", "2025-08-04T13:27:19Z").
		CheckMissing("domains.2.validation_completed_date").
		CheckEqual("domains.2.validation_challenge.cname_record.name", "cname-name-3").
		CheckEqual("domains.2.validation_challenge.cname_record.target", "cname-target-3").
		CheckEqual("domains.2.validation_challenge.txt_record.name", "txt-name-3").
		CheckEqual("domains.2.validation_challenge.txt_record.value", "txt-value-3").
		CheckMissing("domains.2.validation_challenge.http_file.path").
		CheckMissing("domains.2.validation_challenge.http_file.content").
		CheckMissing("domains.2.validation_challenge.http_file.content_type").
		CheckMissing("domains.2.validation_challenge.http_redirect.from").
		CheckMissing("domains.2.validation_challenge.http_redirect.to").
		CheckEqual("domains.2.validation_challenge.expiration_date", "2025-08-05T13:27:19Z")

	request := domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: []domainownership.Domain{
				{DomainName: "dom1.test",
					ValidationScope: "HOST"},
				{DomainName: "dom2.test",
					ValidationScope: "HOST"},
				{DomainName: "dom3.test",
					ValidationScope: "HOST"},
				{DomainName: "dom4.test", // this domain is not returned by the mock, to test partial responses
					ValidationScope: "HOST"},
			},
		},
	}

	tests := map[string]struct {
		init  func(client *domainownership.Mock)
		steps []resource.TestStep
	}{
		"full response": {
			init: func(client *domainownership.Mock) {
				mockSearchDomains(client, request, &domainownership.SearchDomainsResponse{
					Domains: []domainownership.SearchDomainItem{
						{
							AccountID:               ptr.To("1-ACCOUN"),
							DomainName:              "dom1.test",
							ValidationScope:         "HOST",
							DomainStatus:            "REQUEST_ACCEPTED",
							ValidationLevel:         "FQDN",
							ValidationMethod:        ptr.To("DNS_TXT"),
							ValidationRequestedBy:   ptr.To("someuser"),
							ValidationRequestedDate: ptr.To(tst.NewTimeFromString(t, "2025-08-04T13:27:19Z")),
							ValidationChallenge: &domainownership.ValidationChallenge{
								CnameRecord: domainownership.CnameRecord{
									Name:   "cname-name-1",
									Target: "cname-target-1",
								},
								TXTRecord: domainownership.TXTRecord{
									Name:  "txt-name-1",
									Value: "txt-value-1",
								},
								HTTPFile: &domainownership.HTTPFile{
									Path:        "http-file-path-1",
									Content:     "http-file-content-1",
									ContentType: "text/plain",
								},
								HTTPRedirect: &domainownership.HTTPRedirect{
									From: "http-redirect-from-1",
									To:   "http-redirect-to-1",
								},
								ExpirationDate: tst.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
							},
						},
						{
							AccountID:               ptr.To("1-ACCOUN"),
							DomainName:              "dom2.test",
							ValidationScope:         "HOST",
							DomainStatus:            "VALIDATED",
							ValidationLevel:         "FQDN",
							ValidationMethod:        ptr.To("SYSTEM"),
							ValidationRequestedBy:   ptr.To("someuser"),
							ValidationRequestedDate: ptr.To(tst.NewTimeFromString(t, "2025-08-04T13:27:19Z")),
							ValidationCompletedDate: ptr.To(tst.NewTimeFromString(t, "2025-08-05T11:56:07Z")),
						},
						{
							AccountID:               ptr.To("1-ACCOUN"),
							DomainName:              "dom3.test",
							ValidationScope:         "HOST",
							DomainStatus:            "VALIDATED",
							ValidationLevel:         "FQDN",
							ValidationMethod:        ptr.To("DNS_TXT"),
							ValidationRequestedBy:   ptr.To("someuser"),
							ValidationRequestedDate: ptr.To(tst.NewTimeFromString(t, "2025-08-04T13:27:19Z")),
							ValidationChallenge: &domainownership.ValidationChallenge{
								CnameRecord: domainownership.CnameRecord{
									Name:   "cname-name-3",
									Target: "cname-target-3",
								},
								TXTRecord: domainownership.TXTRecord{
									Name:  "txt-name-3",
									Value: "txt-value-3",
								},
								ExpirationDate: tst.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
							},
						},
					},
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"search.tf"),
					Check:  defaultStateChecker.Build(),
				},
			},
		},
		"validation error - no domains": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						}`,
					ExpectError: regexp.MustCompile(`The argument "domains" is required, but no definition was found.`),
				},
			},
		},
		"validation error - empty domains": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						  domains = []
						}`,
					ExpectError: regexp.MustCompile(`Attribute domains set must contain at least 1 elements, got: 0`),
				},
			},
		},
		"validation error - empty domain": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						  domains = [
							{
							}
						  ]
						}`,
					ExpectError: regexp.MustCompile(`Inappropriate value for attribute "domains": element 0: attributes\s+"domain_name" and "validation_scope" are required.`),
				},
			},
		},
		"validation error - missing domain name": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						  domains = [
							{
							  validation_scope = "HOST"
							}
						  ]
						}`,
					ExpectError: regexp.MustCompile(`Inappropriate value for attribute "domains": element 0: attribute\s+"domain_name" is required.`),
				},
			},
		},
		"validation error - empty domain name": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						  domains = [
							{
							  domain_name      = ""
							  validation_scope = "HOST"
							}
						  ]
						}`,
					ExpectError: regexp.MustCompile(`string length must be between 1 and 253, got: 0`),
				},
			},
		},
		"validation error - too long domain name": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						  domains = [
							{
							  domain_name      = "a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890a1234567890"
							  validation_scope = "HOST"
							}
						  ]
						}`,
					ExpectError: regexp.MustCompile(`string length must be between 1 and 253, got: 264`),
				},
			},
		},
		"validation error - missing validation scope": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						  domains = [
							{
							  domain_name      = "test.com"
							}
						  ]
						}`,
					ExpectError: regexp.MustCompile(`Inappropriate value for attribute "domains": element 0: attribute\s+"validation_scope" is required.`),
				},
			},
		},
		"validation error - invalid validation scope": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						 edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_property_domainownership_search_domains" "test" {
						  domains = [
							{
							  domain_name      = "test.com"
							  validation_scope = "FOO"
							}
						  ]
						}`,
					ExpectError: regexp.MustCompile(`value must be one of: \["HOST" "WILDCARD" "DOMAIN"], got: "FOO"`),
				},
			},
		},
		"api error": {
			init: func(client *domainownership.Mock) {
				mockSearchDomains(client, request, nil).Return(nil, &domainownership.Error{
					Type:     "internal-server-error",
					Title:    "Internal Server Error",
					Detail:   "Error processing request",
					Instance: "TestInstances",
					Status:   500,
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"search.tf"),
					ExpectError: regexp.MustCompile(`Error: searching domains failed`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &domainownership.Mock{}
			if tc.init != nil {
				tc.init(client)
			}

			useDomainOwnership(client, func() {
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
