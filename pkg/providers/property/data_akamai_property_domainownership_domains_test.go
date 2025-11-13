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

func TestDomainOwnershipDomainsDataSource(t *testing.T) {
	t.Parallel()

	testDir := "testdata/TestDataDomainOwnershipDomains/"

	commonStateChecker := test.NewStateChecker("data.akamai_property_domainownership_domains.test").
		CheckEqual("domains.#", "4").
		// --- Domain 1 ---
		CheckEqual("domains.0.domain_name", "example1.com").
		CheckEqual("domains.0.validation_scope", "HOST").
		CheckEqual("domains.0.account_id", "act_1-ABCDEF").
		CheckEqual("domains.0.domain_status", "PENDING_VALIDATION").
		CheckEqual("domains.0.validation_method", "DNS").
		CheckEqual("domains.0.validation_requested_by", "user1@example.com").
		CheckEqual("domains.0.validation_requested_date", "2024-12-10T10:15:30Z").
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
		// --- Domain 2 --- missing validation_challenge
		CheckEqual("domains.1.domain_name", "example2.com").
		CheckEqual("domains.1.validation_scope", "DOMAIN").
		CheckEqual("domains.1.account_id", "act_1-PQRST").
		CheckEqual("domains.1.domain_status", "PENDING_VALIDATION").
		CheckEqual("domains.1.validation_method", "DNS").
		CheckEqual("domains.1.validation_requested_by", "user1@example.com").
		CheckEqual("domains.1.validation_requested_date", "2024-12-10T10:15:30Z").
		CheckMissing("domains.1.validation_challenge").
		// --- Domain 3 --- missing http_file and http_redirect in validation_challenge
		CheckEqual("domains.2.domain_name", "example3.com").
		CheckEqual("domains.2.validation_scope", "WILDCARD").
		CheckEqual("domains.2.account_id", "act_1-ABCDEF").
		CheckEqual("domains.2.domain_status", "PENDING_VALIDATION").
		CheckEqual("domains.2.validation_method", "DNS").
		CheckEqual("domains.2.validation_requested_by", "user1@example.com").
		CheckEqual("domains.2.validation_requested_date", "2024-12-10T10:15:30Z").
		CheckEqual("domains.2.validation_challenge.cname_record.name", "cname-name-3").
		CheckEqual("domains.2.validation_challenge.cname_record.target", "cname-target-3").
		CheckEqual("domains.2.validation_challenge.txt_record.name", "txt-name-3").
		CheckEqual("domains.2.validation_challenge.txt_record.value", "txt-value-3").
		CheckEqual("domains.2.validation_challenge.expiration_date", "2025-08-05T13:27:19Z").
		// --- Domain 4 ---
		CheckEqual("domains.3.domain_name", "test.org").
		CheckEqual("domains.3.validation_scope", "HOST").
		CheckEqual("domains.3.account_id", "act_1-XYZ123").
		CheckEqual("domains.3.domain_status", "VALIDATED").
		CheckEqual("domains.3.validation_method", "DNS").
		CheckEqual("domains.3.validation_requested_by", "admin@test.org").
		CheckEqual("domains.3.validation_requested_date", "2024-12-11T09:45:10Z").
		CheckEqual("domains.3.validation_completed_date", "2024-12-12T10:00:00Z").
		CheckEqual("domains.3.validation_challenge.cname_record.name", "cname-name-4").
		CheckEqual("domains.3.validation_challenge.cname_record.target", "cname-target-4").
		CheckEqual("domains.3.validation_challenge.txt_record.name", "txt-name-4").
		CheckEqual("domains.3.validation_challenge.txt_record.value", "txt-value-4").
		CheckEqual("domains.3.validation_challenge.http_file.path", "http-file-path-4").
		CheckEqual("domains.3.validation_challenge.http_file.content", "http-file-content-4").
		CheckEqual("domains.3.validation_challenge.http_file.content_type", "text/plain").
		CheckEqual("domains.3.validation_challenge.http_redirect.from", "http-redirect-from-4").
		CheckEqual("domains.3.validation_challenge.http_redirect.to", "http-redirect-to-4").
		CheckEqual("domains.3.validation_challenge.expiration_date", "2025-08-05T13:27:19Z")

	tests := map[string]struct {
		init  func(*domainownership.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{

		"happy path - list domains with pagination": {
			init: func(m *domainownership.Mock) {
				// --- Page 1 ---
				page1 := &domainownership.ListDomainsResponse{
					Domains: []domainownership.DomainItem{
						{
							AccountID:               "act_1-ABCDEF",
							DomainName:              "example1.com",
							DomainStatus:            "PENDING_VALIDATION",
							ValidationScope:         "HOST",
							ValidationMethod:        ptr.To("DNS"),
							ValidationRequestedBy:   "user1@example.com",
							ValidationRequestedDate: tst.NewTimeFromStringMust("2024-12-10T10:15:30Z"),
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
							AccountID:               "act_1-PQRST",
							DomainName:              "example2.com",
							DomainStatus:            "PENDING_VALIDATION",
							ValidationScope:         "DOMAIN",
							ValidationMethod:        ptr.To("DNS"),
							ValidationRequestedBy:   "user1@example.com",
							ValidationRequestedDate: tst.NewTimeFromStringMust("2024-12-10T10:15:30Z"),
						},
						{
							AccountID:               "act_1-ABCDEF",
							DomainName:              "example3.com",
							DomainStatus:            "PENDING_VALIDATION",
							ValidationScope:         "WILDCARD",
							ValidationMethod:        ptr.To("DNS"),
							ValidationRequestedBy:   "user1@example.com",
							ValidationRequestedDate: tst.NewTimeFromStringMust("2024-12-10T10:15:30Z"),
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
					Links: []domainownership.Link{
						{Href: "/domain-validation-service/v1/domains?page=1&pageSize=3", Rel: "self"},
						{Href: "/domain-validation-service/v1/domains?page=2&pageSize=1", Rel: "next"},
					},
					Metadata: domainownership.Metadata{
						HasNext:     true,
						HasPrevious: false,
						Page:        1,
						PageSize:    3, // "Even though we request 1000 items, the mock returns only 3 on the first page for testing pagination.
						TotalItems:  4,
					},
				}

				// --- Page 2 ---
				page2 := &domainownership.ListDomainsResponse{
					Domains: []domainownership.DomainItem{
						{
							AccountID:               "act_1-XYZ123",
							DomainName:              "test.org",
							DomainStatus:            "VALIDATED",
							ValidationScope:         "HOST",
							ValidationMethod:        ptr.To("DNS"),
							ValidationRequestedBy:   "admin@test.org",
							ValidationRequestedDate: tst.NewTimeFromStringMust("2024-12-11T09:45:10Z"),
							ValidationCompletedDate: ptr.To(tst.NewTimeFromStringMust("2024-12-12T10:00:00Z")),
							ValidationChallenge: &domainownership.ValidationChallenge{
								CnameRecord: domainownership.CnameRecord{
									Name:   "cname-name-4",
									Target: "cname-target-4",
								},
								TXTRecord: domainownership.TXTRecord{
									Name:  "txt-name-4",
									Value: "txt-value-4",
								},
								HTTPFile: &domainownership.HTTPFile{
									Path:        "http-file-path-4",
									Content:     "http-file-content-4",
									ContentType: "text/plain",
								},
								HTTPRedirect: &domainownership.HTTPRedirect{
									From: "http-redirect-from-4",
									To:   "http-redirect-to-4",
								},
								ExpirationDate: tst.NewTimeFromString(t, "2025-08-05T13:27:19Z"),
							},
						},
					},
					Links: []domainownership.Link{
						{Href: "/domain-validation-service/v1/domains?page=1&pageSize=3", Rel: "prev"},
						{Href: "/domain-validation-service/v1/domains?page=2&pageSize=1", Rel: "self"},
					},
					Metadata: domainownership.Metadata{
						HasNext:     false,
						HasPrevious: true,
						Page:        2,
						PageSize:    1, // The mock returns only 1 on the second page to test that pagination is handled correctly.
						TotalItems:  4,
					},
				}

				m.On("ListDomains", mock.Anything, domainownership.ListDomainsRequest{
					Paginate: ptr.To(true),
					Page:     1,
					PageSize: 1000,
				}).Return(page1, nil).Times(3)

				m.On("ListDomains", mock.Anything, domainownership.ListDomainsRequest{
					Paginate: ptr.To(true),
					Page:     2,
					PageSize: 1000,
				}).Return(page2, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"domains.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - API returns empty list": {
			init: func(m *domainownership.Mock) {
				emptyResponse := &domainownership.ListDomainsResponse{
					Domains: []domainownership.DomainItem{},
					Metadata: domainownership.Metadata{
						HasNext: false,
						Page:    1,
					},
				}

				m.On("ListDomains", mock.Anything, domainownership.ListDomainsRequest{
					Paginate: ptr.To(true),
					Page:     1,
					PageSize: 1000,
				}).Return(emptyResponse, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"domains.tf"),
					Check: test.NewStateChecker("data.akamai_property_domainownership_domains.test").
						CheckEqual("domains.#", "0").
						Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *domainownership.Mock) {
				m.On("ListDomains", mock.Anything, domainownership.ListDomainsRequest{
					Paginate: ptr.To(true),
					Page:     1,
					PageSize: 1000,
				}).Return(nil, &domainownership.Error{
					Status: 500,
					Detail: "oops",
				}).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"domains.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error - API error on second page": {
			init: func(m *domainownership.Mock) {
				page1 := &domainownership.ListDomainsResponse{
					Domains: []domainownership.DomainItem{
						{
							AccountID:               "act_1-ABCDEF",
							DomainName:              "example1.com",
							DomainStatus:            "PENDING_VALIDATION",
							ValidationScope:         "DOMAIN",
							ValidationMethod:        ptr.To("DNS"),
							ValidationRequestedBy:   "user1@example.com",
							ValidationRequestedDate: tst.NewTimeFromStringMust("2024-12-10T10:15:30Z"),
						},
					},
					Metadata: domainownership.Metadata{
						HasNext: true,
						Page:    1,
					},
				}

				m.On("ListDomains", mock.Anything, domainownership.ListDomainsRequest{
					Paginate: ptr.To(true),
					Page:     1,
					PageSize: 1000,
				}).Return(page1, nil).Times(1)

				m.On("ListDomains", mock.Anything, domainownership.ListDomainsRequest{
					Paginate: ptr.To(true),
					Page:     2,
					PageSize: 1000,
				}).Return(nil, &domainownership.Error{
					Status: 500,
					Detail: "error on second page"}).Times(1)
			},

			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"domains.tf"),
					ExpectError: regexp.MustCompile("error on second page"),
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
