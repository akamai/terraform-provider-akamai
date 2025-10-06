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

	commonStateChecker := test.NewStateChecker("data.akamai_domainownership_domains.test").
		CheckEqual("domains.#", "2").
		// --- Domain 1 ---
		CheckEqual("domains.0.domain_name", "example1.com").
		CheckEqual("domains.0.validation_scope", "ZONE").
		CheckEqual("domains.0.account_id", "act_1-ABCDEF").
		CheckEqual("domains.0.domain_status", "PENDING_VALIDATION").
		CheckEqual("domains.0.validation_method", "DNS").
		CheckEqual("domains.0.validation_requested_by", "user1@example.com").
		CheckEqual("domains.0.validation_requested_date", "2024-12-10T10:15:30Z").
		CheckEqual("domains.0.validation_challenge.challenge_token", "abcd1234").
		CheckEqual("domains.0.validation_challenge.dns_cname", "_acme-challenge.example.com").
		CheckEqual("domains.0.validation_challenge.http_redirect_from", "http://example.com/.well-known/acme-challenge/abcd1234").
		CheckEqual("domains.0.validation_challenge.http_redirect_to", "http://example.com/.well-known/acme-challenge/abcd1234").
		// --- Domain 2 ---
		CheckEqual("domains.1.domain_name", "test.org").
		CheckEqual("domains.1.validation_scope", "ZONE").
		CheckEqual("domains.1.account_id", "act_1-XYZ123").
		CheckEqual("domains.1.domain_status", "VALIDATED").
		CheckEqual("domains.1.validation_method", "DNS").
		CheckEqual("domains.1.validation_requested_by", "admin@test.org").
		CheckEqual("domains.1.validation_requested_date", "2024-12-11T09:45:10Z").
		CheckEqual("domains.1.validation_completed_date", "2024-12-12T10:00:00Z").
		CheckEqual("domains.1.validation_challenge.challenge_token", "efgh5678").
		CheckEqual("domains.1.validation_challenge.dns_cname", "_acme-challenge.test.org")

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
							ValidationScope:         "ZONE",
							ValidationMethod:        ptr.To("DNS"),
							ValidationRequestedBy:   "user1@example.com",
							ValidationRequestedDate: tst.NewTimeFromStringMust("2024-12-10T10:15:30Z"),
							ValidationChallenge: &domainownership.ValidationChallenge{
								ChallengeToken:            "abcd1234",
								ChallengeTokenExpiresDate: tst.NewTimeFromStringMust("2024-12-15T10:15:30Z"),
								DNSCname:                  "_acme-challenge.example.com",
								HTTPRedirectFrom:          ptr.To("http://example.com/.well-known/acme-challenge/abcd1234"),
								HTTPRedirectTo:            ptr.To("http://example.com/.well-known/acme-challenge/abcd1234"),
							},
						},
					},
					Links: []domainownership.Link{
						{Href: "/domain-validation-service/v1/domains?page=1&pageSize=1", Rel: "self"},
						{Href: "/domain-validation-service/v1/domains?page=2&pageSize=1", Rel: "next"},
					},
					Metadata: domainownership.Metadata{
						HasNext:     true,
						HasPrevious: false,
						Page:        1,
						PageSize:    1,
						TotalItems:  2,
					},
				}

				// --- Page 2 ---
				page2 := &domainownership.ListDomainsResponse{
					Domains: []domainownership.DomainItem{
						{
							AccountID:               "act_1-XYZ123",
							DomainName:              "test.org",
							DomainStatus:            "VALIDATED",
							ValidationScope:         "ZONE",
							ValidationMethod:        ptr.To("DNS"),
							ValidationRequestedBy:   "admin@test.org",
							ValidationRequestedDate: tst.NewTimeFromStringMust("2024-12-11T09:45:10Z"),
							ValidationCompletedDate: ptr.To(tst.NewTimeFromStringMust("2024-12-12T10:00:00Z")),
							ValidationChallenge: &domainownership.ValidationChallenge{
								ChallengeToken:            "efgh5678",
								ChallengeTokenExpiresDate: tst.NewTimeFromStringMust("2024-12-17T09:45:10Z"),
								DNSCname:                  "_acme-challenge.test.org",
							},
						},
					},
					Links: []domainownership.Link{
						{Href: "/domain-validation-service/v1/domains?page=1&pageSize=1", Rel: "prev"},
						{Href: "/domain-validation-service/v1/domains?page=2&pageSize=1", Rel: "self"},
					},
					Metadata: domainownership.Metadata{
						HasNext:     false,
						HasPrevious: true,
						Page:        2,
						PageSize:    1,
						TotalItems:  2,
					},
				}

				m.On("ListDomains", mock.Anything, mock.MatchedBy(func(req domainownership.ListDomainsRequest) bool {
					return req.Page != nil && *req.Page == 1 && *req.PageSize == 1000
				})).Return(page1, nil).Once()

				m.On("ListDomains", mock.Anything, mock.MatchedBy(func(req domainownership.ListDomainsRequest) bool {
					return req.Page != nil && *req.Page == 2 && *req.PageSize == 1000
				})).Return(page2, nil).Once()

				m.On("ListDomains", mock.Anything, mock.MatchedBy(func(req domainownership.ListDomainsRequest) bool {
					return req.Page != nil && *req.PageSize == 1000
				})).Return(&domainownership.ListDomainsResponse{
					Metadata: domainownership.Metadata{
						HasNext: false,
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"domains.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *domainownership.Mock) {
				m.On("ListDomains", mock.Anything, mock.MatchedBy(func(req domainownership.ListDomainsRequest) bool {
					return req.Page != nil && *req.Page == 1 && *req.PageSize == 1000
				})).Return(nil, &domainownership.Error{
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
							ValidationScope:         "ZONE",
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

				m.On("ListDomains", mock.Anything, mock.MatchedBy(func(req domainownership.ListDomainsRequest) bool {
					return req.Page != nil && *req.Page == 1 && *req.PageSize == 1000
				})).Return(page1, nil).Once()

				m.On("ListDomains", mock.Anything, mock.MatchedBy(func(req domainownership.ListDomainsRequest) bool {
					return req.Page != nil && *req.Page == 2 && *req.PageSize == 1000
				})).Return(nil, &domainownership.Error{
					Status: 500,
					Detail: "error on second page",
				}).Once()
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
