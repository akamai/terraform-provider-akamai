package property

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/mock"
)

func TestDomainOwnershipDomainsRes(t *testing.T) {

	test1 := test.AttributeBatch{
		"domain_name":                                 "test1.com",
		"validation_scope":                            "HOST",
		"domain_status":                               "REQUEST_ACCEPTED",
		"account_id":                                  "ACC123",
		"validation_requested_by":                     "someone",
		"validation_requested_date":                   "2024-01-01T00:00:00Z",
		"validation_challenge.cname_record.name":      "challenge_test1.com",
		"validation_challenge.cname_record.target":    "target_test1.com",
		"validation_challenge.txt_record.name":        "challenge_test1.com",
		"validation_challenge.txt_record.value":       "value_test1.com",
		"validation_challenge.expiration_date":        "2024-01-10T00:00:00Z",
		"validation_challenge.http_file.path":         "/challenge_path",
		"validation_challenge.http_file.content":      "file_content_test1.com",
		"validation_challenge.http_file.content_type": "text/plain",
		"validation_challenge.http_redirect.from":     "/from_path",
		"validation_challenge.http_redirect.to":       "/to_path_test1.com",
	}

	test2 := test.AttributeBatch{
		"domain_name":                                 "test2.com",
		"validation_scope":                            "HOST",
		"domain_status":                               "REQUEST_ACCEPTED",
		"account_id":                                  "ACC123",
		"validation_requested_by":                     "someone",
		"validation_requested_date":                   "2024-01-01T00:00:00Z",
		"validation_challenge.cname_record.name":      "challenge_test2.com",
		"validation_challenge.cname_record.target":    "target_test2.com",
		"validation_challenge.txt_record.name":        "challenge_test2.com",
		"validation_challenge.txt_record.value":       "value_test2.com",
		"validation_challenge.expiration_date":        "2024-01-10T00:00:00Z",
		"validation_challenge.http_file.path":         "/challenge_path",
		"validation_challenge.http_file.content":      "file_content_test2.com",
		"validation_challenge.http_file.content_type": "text/plain",
		"validation_challenge.http_redirect.from":     "/from_path",
		"validation_challenge.http_redirect.to":       "/to_path_test2.com",
	}

	test3 := test.AttributeBatch{
		"domain_name":                                 "test3.com",
		"validation_scope":                            "HOST",
		"domain_status":                               "REQUEST_ACCEPTED",
		"account_id":                                  "ACC123",
		"validation_requested_by":                     "someone",
		"validation_requested_date":                   "2024-01-01T00:00:00Z",
		"validation_challenge.cname_record.name":      "challenge_test3.com",
		"validation_challenge.cname_record.target":    "target_test3.com",
		"validation_challenge.txt_record.name":        "challenge_test3.com",
		"validation_challenge.txt_record.value":       "value_test3.com",
		"validation_challenge.expiration_date":        "2024-01-10T00:00:00Z",
		"validation_challenge.http_file.path":         "/challenge_path",
		"validation_challenge.http_file.content":      "file_content_test3.com",
		"validation_challenge.http_file.content_type": "text/plain",
		"validation_challenge.http_redirect.from":     "/from_path",
		"validation_challenge.http_redirect.to":       "/to_path_test3.com",
	}

	test4 := test.AttributeBatch{
		"domain_name":                                 "test4.com",
		"validation_scope":                            "HOST",
		"domain_status":                               "REQUEST_ACCEPTED",
		"account_id":                                  "ACC123",
		"validation_requested_by":                     "someone",
		"validation_requested_date":                   "2024-01-01T00:00:00Z",
		"validation_challenge.cname_record.name":      "challenge_test4.com",
		"validation_challenge.cname_record.target":    "target_test4.com",
		"validation_challenge.txt_record.name":        "challenge_test4.com",
		"validation_challenge.txt_record.value":       "value_test4.com",
		"validation_challenge.expiration_date":        "2024-01-10T00:00:00Z",
		"validation_challenge.http_file.path":         "/challenge_path",
		"validation_challenge.http_file.content":      "file_content_test4.com",
		"validation_challenge.http_file.content_type": "text/plain",
		"validation_challenge.http_redirect.from":     "/from_path",
		"validation_challenge.http_redirect.to":       "/to_path_test4.com",
	}

	twoDomainsChecker := test.NewStateChecker("akamai_property_domainownership_domains.test").
		CheckEqual("domains.#", "2").
		CheckEqualBatch("domains.0.", test1).
		CheckEqualBatch("domains.1.", test2)

	fourDomainsChecker := test.NewStateChecker("akamai_property_domainownership_domains.test").
		CheckEqual("domains.#", "4").
		CheckEqualBatch("domains.0.", test1).
		CheckEqualBatch("domains.1.", test2).
		CheckEqualBatch("domains.2.", test3).
		CheckEqualBatch("domains.3.", test4)

	tests := map[string]struct {
		init  func(mock *domainownership.Mock)
		steps []resource.TestStep
	}{
		"create multiple domains": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Setting up state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"create multiple domains with different validation scopes": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeDomain,
					"test3.com", domainownership.ValidationScopeWildcard,
				)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeDomain,
					"test3.com", domainownership.ValidationScopeWildcard),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "DOMAIN"),
							getAddDomainSuccess(t, "test3.com", "WILDCARD"),
						},
					}, nil)
				// Setting up state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "DOMAIN", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "WILDCARD", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "DOMAIN", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "WILDCARD", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "DOMAIN", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "WILDCARD", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeDomain,
					"test3.com", domainownership.ValidationScopeWildcard), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/different_validation_scopes.tf"),
					Check: test.NewStateChecker("akamai_property_domainownership_domains.test").
						CheckEqual("domains.#", "3").
						CheckEqualBatch("domains.0.", test1).
						CheckEqual("domains.1.domain_name", "test2.com").
						CheckEqual("domains.1.validation_scope", "DOMAIN").
						CheckEqual("domains.1.domain_status", "REQUEST_ACCEPTED").
						CheckEqual("domains.1.account_id", "ACC123").
						CheckEqual("domains.1.validation_requested_by", "someone").
						CheckEqual("domains.1.validation_requested_date", "2024-01-01T00:00:00Z").
						CheckEqual("domains.1.validation_challenge.cname_record.name", "challenge_test2.com").
						CheckEqual("domains.1.validation_challenge.cname_record.target", "target_test2.com").
						CheckEqual("domains.1.validation_challenge.txt_record.name", "challenge_test2.com").
						CheckEqual("domains.1.validation_challenge.txt_record.value", "value_test2.com").
						CheckEqual("domains.1.validation_challenge.expiration_date", "2024-01-10T00:00:00Z").
						CheckEqual("domains.2.domain_name", "test3.com").
						CheckEqual("domains.2.validation_scope", "WILDCARD").
						CheckEqual("domains.2.domain_status", "REQUEST_ACCEPTED").
						CheckEqual("domains.2.account_id", "ACC123").
						CheckEqual("domains.2.validation_requested_by", "someone").
						CheckEqual("domains.2.validation_requested_date", "2024-01-01T00:00:00Z").
						CheckEqual("domains.2.validation_challenge.cname_record.name", "challenge_test3.com").
						CheckEqual("domains.2.validation_challenge.cname_record.target", "target_test3.com").
						CheckEqual("domains.2.validation_challenge.txt_record.name", "challenge_test3.com").
						CheckEqual("domains.2.validation_challenge.txt_record.value", "value_test3.com").
						CheckEqual("domains.2.validation_challenge.expiration_date", "2024-01-10T00:00:00Z").
						Build(),
				},
			},
		},
		"create multiple domains but one fails - successful rollback": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
						},
						Errors: []domainownership.AddDomainError{
							{
								DomainName:      "test2.com",
								ValidationScope: "HOST",
								Detail:          "Domain already exists.",
								Type:            "error-type",
								Title:           "Domain Add Error",
							},
						},
					}, nil)
				// Rollback
				mockDeleteDomains(mock, buildDeleteDomainsRequest("test1.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					ExpectError: regexp.MustCompile(`(?s)error adding domains.+\{\s+domainName: test2.com,\s+validationScope: HOST,\s+title: Domain Add Error,\s+detail: Domain already exists.\s+}\s+Rollback was successful`),
				},
			},
		},
		"create multiple domains but one fails - unsuccessful rollback": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
						},
						Errors: []domainownership.AddDomainError{
							{
								DomainName:      "test2.com",
								ValidationScope: "HOST",
								Detail:          "Domain already exists.",
								Type:            "error-type",
								Title:           "Domain Add Error",
							},
						},
					}, nil)
				// Rollback
				mockDeleteDomains(mock, buildDeleteDomainsRequest("test1.com", domainownership.ValidationScopeHost), &domainownership.Error{Type: "error-type", Title: "Domain Delete Error", Detail: "Domain cannot be deleted."})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					ExpectError: regexp.MustCompile(`(?s)error adding domains.+\{\s+domainName: test2.com,\s+validationScope: HOST,\s+title: Domain Add Error,\s+detail: Domain already exists.\s+}\s+Rollback was not successful: API error:`),
				},
			},
		},
		"create multiple domains - error on create": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					nil, &domainownership.Error{Type: "error-type", Title: "Domain Add Error", Detail: "Service unavailable."})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					ExpectError: regexp.MustCompile(`(?s)error adding domains.+API error`),
				},
			},
		},
		"create multiple domains but all were already existing": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"create multiple domains but one already existed": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"create multiple domains but one already existed as VALIDATED ROOT/WILDCARD": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				wildcardItem := getSearchDomainItem(t, "test1.com", "HOST", "VALIDATED")
				wildcardItem.ValidationLevel = "ROOT/WILDCARD"
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					wildcardItem,
				})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					ExpectError: regexp.MustCompile(`(?s)error adding domains.+domain test1.com with validation scope HOST is already part of other, already.+validated domain/wildcard and cannot be added again`),
				},
			},
		},
		"create multiple domains but one already existed as not VALIDATED ROOT/WILDCARD": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				wildcardItem := getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED")
				wildcardItem.ValidationLevel = "ROOT/WILDCARD"
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					wildcardItem,
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"create multiple domains but all were already existing but with a status that required refreshing": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "INVALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "TOKEN_EXPIRED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"create multiple domains but their status changed right before destroy": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Setting up state for Create (normally the domains wouldn't be VALIDATED right away - keeping this for test simplicity)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "VALIDATED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "VALIDATED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "INVALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "INVALIDATED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check: twoDomainsChecker.
						CheckEqual("domains.0.domain_status", "VALIDATED").
						CheckEqual("domains.1.domain_status", "VALIDATED").
						Build(),
				},
			},
		},
		"create multiple domains with same domain name but different validation scopes": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "DOMAIN"),
							getAddDomainSuccess(t, "test1.com", "HOST"),
						},
					}, nil)
				// Setting up state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "DOMAIN", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "DOMAIN", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "DOMAIN", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/same_domain_name.tf"),
					Check: test.NewStateChecker("akamai_property_domainownership_domains.test").
						CheckEqualBatch("domains.0.", test1).
						CheckEqualBatch("domains.1.", test1).
						CheckEqual("domains.1.validation_scope", "DOMAIN").
						CheckMissing("domains.1.validation_challenge.http_file.path").
						CheckMissing("domains.1.validation_challenge.http_file.content").
						CheckMissing("domains.1.validation_challenge.http_file.content_type").
						CheckMissing("domains.1.validation_challenge.http_redirect.from").
						CheckMissing("domains.1.validation_challenge.http_redirect.to").
						Build(),
				},
			},
		},
		"error - incorrect validationScope": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_domains" "test" {
  domains                = [
    {
      domain_name = "test1.com"
      validation_scope = "INVALID"
    },
  ]
}`,
					ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.+value must be one of: \["HOST" "WILDCARD" "DOMAIN"], got: "INVALID"`),
				},
			},
		},
		"delete has to invalidate domains": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				// Normally the domains wouldn't be VALIDATED right away - keeping this for test simplicity
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "VALIDATED"),
				})
				// Read after create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "VALIDATED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "VALIDATED"),
				})
				mockInvalidateDomains(mock, map[domainKey]domainDetails{
					{domainName: "test1.com", validationScope: "HOST"}: {validationStatus: "VALIDATED"},
					{domainName: "test2.com", validationScope: "HOST"}: {validationStatus: "VALIDATED"},
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Check:  resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr("akamai_property_domainownership_domains.test", "domains.#", "2")),
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
				},
			},
		},
		"update - multiple new domains": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create + Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test3.com", "HOST"),
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Read after update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
			},
		},
		"update - multiple new domains but one fails - successful rollback": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test3.com", "HOST"),
						},
						Errors: []domainownership.AddDomainError{
							{
								DomainName:      "test4.com",
								ValidationScope: "HOST",
								Detail:          "Domain already exists.",
								Type:            "error-type",
								Title:           "Domain Add Error",
							},
						},
					}, nil)
				mockDeleteDomains(mock, buildDeleteDomainsRequest("test3.com", domainownership.ValidationScopeHost), nil)
				// Delete
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					ExpectError: regexp.MustCompile(`(?s)error adding domains.+\{\s+domainName: test4.com,\s+validationScope: HOST,\s+title: Domain Add Error,\s+detail: Domain already exists.\s+}\s+Rollback was successful`),
				},
			},
		},
		"update - multiple new domains but one fails - unsuccessful rollback": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test3.com", "HOST"),
						},
						Errors: []domainownership.AddDomainError{
							{
								DomainName:      "test4.com",
								ValidationScope: "HOST",
								Detail:          "Domain already exists.",
								Type:            "error-type",
								Title:           "Domain Add Error",
							},
						},
					}, nil)
				mockDeleteDomains(mock, buildDeleteDomainsRequest("test3.com", domainownership.ValidationScopeHost), &domainownership.Error{Type: "error-type", Title: "Domain Delete Error", Detail: "Domain cannot be deleted."})
				// Delete
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Once()
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					ExpectError: regexp.MustCompile(`(?s)error adding domains.+\{\s+domainName: test4.com,\s+validationScope: HOST,\s+title: Domain Add Error,\s+detail: Domain already exists.\s+}\s+Rollback was not successful: API error:`),
				},
			},
		},
		"update - multiple new domains but all new were already existing": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Set state after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
			},
		},
		"update - multiple new domains but one already existed": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest("test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Set state after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
			},
		},
		"update - multiple new domains but one already existed as VALIDATED ROOT/WILDCARD": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				wildcard := getSearchDomainItem(t, "test3.com", "HOST", "VALIDATED")
				wildcard.ValidationLevel = "ROOT/WILDCARD"
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					wildcard,
				})
				// Delete
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Once()
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					ExpectError: regexp.MustCompile(`(?s)error adding domains.+domain test3.com with validation scope HOST is already part of other, already.+validated domain/wildcard and cannot be added again`),
				},
			},
		},
		"update - multiple new domains but one already existed as not VALIDATED ROOT/WILDCARD": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				wildcard := getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED")
				wildcard.ValidationLevel = "ROOT/WILDCARD"
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					wildcard,
				})
				mockAddDomains(mock,
					buildAddDomainsRequest("test3.com", domainownership.ValidationScopeHost,
						"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test3.com", "HOST"),
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Set state after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
			},
		},
		"update - multiple new domains but all new were already existing but with a status that required refreshing": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "INVALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "TOKEN_EXPIRED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test3.com", "HOST"),
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Set state after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
			},
		},
		"update - dropping multiple domains": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
							getAddDomainSuccess(t, "test3.com", "HOST"),
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
				// Set state after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"update - dropping multiple domains but the status changed right before entering update": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
							getAddDomainSuccess(t, "test3.com", "HOST"),
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					// Normally wouldn't be VALIDATED right after creation, but having it for tests simplicity
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test4.com", "HOST", "VALIDATED"),
				})
				// Read after Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					// Normally wouldn't be VALIDATED right after creation, but having it for tests simplicity
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test4.com", "HOST", "VALIDATED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "INVALIDATED"),
					getSearchDomainItem(t, "test4.com", "HOST", "INVALIDATED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
				// Set state after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check: fourDomainsChecker.
						CheckEqual("domains.2.domain_status", "VALIDATED").
						CheckEqual("domains.3.domain_status", "VALIDATED").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"update - adding multiple domains and dropping multiple domains": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
							getAddDomainSuccess(t, "test3.com", "HOST"),
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test5.com", domainownership.ValidationScopeHost,
					"test6.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test5.com", domainownership.ValidationScopeHost,
					"test6.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test5.com", "HOST"),
							getAddDomainSuccess(t, "test6.com", "HOST"),
						},
					}, nil)
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
				// Set state after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test5.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test6.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test5.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test6.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test5.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test6.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test5.com", domainownership.ValidationScopeHost,
					"test6.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_different_domains.tf"),
					Check: fourDomainsChecker.
						CheckEqual("domains.2.domain_name", "test5.com").
						CheckEqual("domains.2.validation_challenge.cname_record.name", "challenge_test5.com").
						CheckEqual("domains.2.validation_challenge.cname_record.target", "target_test5.com").
						CheckEqual("domains.2.validation_challenge.txt_record.name", "challenge_test5.com").
						CheckEqual("domains.2.validation_challenge.txt_record.value", "value_test5.com").
						CheckEqual("domains.2.validation_challenge.http_file.content", "file_content_test5.com").
						CheckEqual("domains.2.validation_challenge.http_redirect.to", "/to_path_test5.com").
						CheckEqual("domains.3.domain_name", "test6.com").
						CheckEqual("domains.3.validation_challenge.cname_record.name", "challenge_test6.com").
						CheckEqual("domains.3.validation_challenge.cname_record.target", "target_test6.com").
						CheckEqual("domains.3.validation_challenge.txt_record.name", "challenge_test6.com").
						CheckEqual("domains.3.validation_challenge.txt_record.value", "value_test6.com").
						CheckEqual("domains.3.validation_challenge.http_file.content", "file_content_test6.com").
						CheckEqual("domains.3.validation_challenge.http_redirect.to", "/to_path_test6.com").
						Build(),
				},
			},
		},
		"update - dropping multiple domains but one is VALIDATED": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
							getAddDomainSuccess(t, "test3.com", "HOST"),
							getAddDomainSuccess(t, "test4.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read before Update + fetch during Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "VALIDATED"),
				})
				// Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test3.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test4.com", "HOST", "VALIDATED"),
				})
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockInvalidateDomains(mock, map[domainKey]domainDetails{
					{domainName: "test4.com", validationScope: "HOST"}: {validationStatus: "VALIDATED"},
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test3.com", domainownership.ValidationScopeHost,
					"test4.com", domainownership.ValidationScopeHost), nil)
				// Set state after Update
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/four_domains.tf"),
					Check:  fourDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
		"read on update returns domains as VALIDATED": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read before first Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "VALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "VALIDATED"),
				}).Twice()
				// Read before second Update to set correct state before Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "INVALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "INVALIDATED"),
				}).Twice()
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "INVALIDATED"),
					getSearchDomainItem(t, "test2.com", "HOST", "INVALIDATED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					RefreshState: true,
					Check: twoDomainsChecker.
						CheckEqual("domains.0.domain_status", "VALIDATED").
						CheckEqual("domains.0.validation_method", "DNS_CNAME").
						CheckEqual("domains.0.validation_completed_date", "2024-01-02T00:00:00Z").
						CheckEqual("domains.1.domain_status", "VALIDATED").
						CheckEqual("domains.1.validation_method", "DNS_CNAME").
						CheckEqual("domains.1.validation_completed_date", "2024-01-02T00:00:00Z").
						Build(),
				},
				{
					RefreshState: true,
				},
			},
		},
		"read on update returns no domains": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read before Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read that returns no domains
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"read dropped one existing domain": {
			init: func(mock *domainownership.Mock) {
				// Create
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{})
				mockAddDomains(mock, buildAddDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test1.com", "HOST"),
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state for Create
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				}).Once()
				// Update
				// Check for domains from state
				searchRequest = buildSearchDomainsRequestTest("test1.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				})
				//Check for domains from plan
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockAddDomains(mock, buildAddDomainsRequest("test2.com", domainownership.ValidationScopeHost),
					&domainownership.AddDomainsResponse{
						Successes: []domainownership.AddDomainSuccess{
							getAddDomainSuccess(t, "test2.com", "HOST"),
						},
					}, nil)
				// Set state after update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Read after Update
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				// Delete
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				})
				mockDeleteDomains(mock, buildDeleteDomainsRequest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost,
				), nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
					Check:  twoDomainsChecker.Build(),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &domainownership.Mock{}
			if tc.init != nil {
				tc.init(m)
			}

			useDomainOwnership(m, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})

			m.AssertExpectations(t)
		})
	}
}

func TestDomainOwnershipDomainsImport(t *testing.T) {
	importChecker := test.NewImportChecker().
		CheckEqual("domains.#", "2")

	tests := map[string]struct {
		importID    string
		init        func(mock *domainownership.Mock)
		config      string
		stateCheck  func(s []*terraform.InstanceState) error
		expectError *regexp.Regexp
	}{
		"import with multiple domains with validationScopes": {
			config:     testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			stateCheck: importChecker.Build(),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
			},
			importID: "test1.com:HOST,test2.com:HOST",
		},
		"import with multiple domains of same name but with different validationScopes": {
			config:     testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/same_domain_name.tf"),
			stateCheck: importChecker.Build(),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "DOMAIN", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				}).Twice()
			},
			importID: "test1.com:HOST,test1.com:DOMAIN",
		},
		"import with multiple domains without validationScopes": {
			config:     testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			stateCheck: importChecker.Build(),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost,
					"test1.com", domainownership.ValidationScopeWildcard,
					"test2.com", domainownership.ValidationScopeDomain,
					"test2.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeWildcard)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Times(1)
				// Read
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Times(1)
			},
			importID: "test1.com,test2.com",
		},
		"import with domains with and without validationScopes": {
			config:     testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			stateCheck: importChecker.Build(),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeDomain,
					"test2.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeWildcard)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Times(1)
				// Read
				searchRequest = buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
				}).Times(1)
			},
			importID: "test1.com:HOST,test2.com",
		},
		"import with multiple domains with validationScopes but one is not found": {
			config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				}).Times(1)
			},
			importID:    "test1.com:HOST,test2.com:HOST",
			expectError: regexp.MustCompile(`(?s)Error verifying domains.+the domain 'test2.com' with validation scope 'HOST' was not found`),
		},
		"import with multiple domains without validationScopes but one is not found": {
			config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			expectError: regexp.MustCompile(`(?s)Error verifying domains.+the domain 'test2.com' was not found`),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost,
					"test1.com", domainownership.ValidationScopeWildcard,
					"test2.com", domainownership.ValidationScopeDomain,
					"test2.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeWildcard)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
				}).Times(1)
			},
			importID: "test1.com,test2.com",
		},
		"import with multiple domains without validationScopes but one is ambiguous": {
			config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			expectError: regexp.MustCompile(`(?s)Error verifying domains.+the domain 'test2.com' exists with multiple validation scopes. Please.+re-import specifying the validation scope for the domain`),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost,
					"test1.com", domainownership.ValidationScopeWildcard,
					"test2.com", domainownership.ValidationScopeDomain,
					"test2.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeWildcard)
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED"),
					getSearchDomainItem(t, "test2.com", "WILDCARD", "REQUEST_ACCEPTED"),
				}).Times(1)
			},
			importID: "test1.com,test2.com",
		},
		"import with multiple domains with too lengthy ID with validationScope": {
			config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_domains" "test" {
  domains                = [
    for i in range(1, 1001): {
      domain_name = "test${i}.example.com"
      validation_scope = "HOST"
    }
  ]
}`,
			expectError: regexp.MustCompile("the maximum number of domains that can be imported is 1000, got 1001"),
			importID:    generateLongImportID("HOST"),
		},
		"import with multiple domains with too lengthy ID without validationScope": {
			config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_domains" "test" {
  domains                = [
    for i in range(1, 1001): {
      domain_name = "test${i}.example.com"
    }
  ]
}`,
			expectError: regexp.MustCompile("the maximum number of domains that can be imported is 1000, got 1001"),
			importID:    generateLongImportID(""),
		},
		"import with multiple domains with too lengthy ID with mixed validationScope": {
			config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_domains" "test" {
  domains                = concat([
    for i in range(1, 1000): {
      domain_name = "test${i}.example.com"
    }
  ], [
    {
      domain_name = "test1001.example.com"
      validation_scope = "HOST"
    }
  ])
}
`,
			expectError: regexp.MustCompile("the maximum number of domains that can be imported is 1000, got 1001"),
			importID:    generateLongImportID("") + ":HOST",
		},
		"import with multiple domains with validationScopes but one is not FQDN": {
			config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeHost)
				wildcard := getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED")
				wildcard.ValidationLevel = "ROOT/WILDCARD"
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					wildcard,
				}).Once()
			},
			importID:    "test1.com:HOST,test2.com:HOST",
			expectError: regexp.MustCompile(`(?s)Error verifying domains.+only domains with validation level FQDN can be imported, the requested domain.+'test2.com' with validation scope 'HOST' has validation level 'ROOT/WILDCARD'`),
		},
		"import with multiple domains without validationScopes but one is not FQDN": {
			config: testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/two_domains.tf"),
			init: func(mock *domainownership.Mock) {
				searchRequest := buildSearchDomainsRequestTest(
					"test1.com", domainownership.ValidationScopeDomain,
					"test1.com", domainownership.ValidationScopeHost,
					"test1.com", domainownership.ValidationScopeWildcard,
					"test2.com", domainownership.ValidationScopeDomain,
					"test2.com", domainownership.ValidationScopeHost,
					"test2.com", domainownership.ValidationScopeWildcard)
				wildcard := getSearchDomainItem(t, "test2.com", "HOST", "REQUEST_ACCEPTED")
				wildcard.ValidationLevel = "ROOT/WILDCARD"
				mockBasicSearchDomains(mock, searchRequest, []domainownership.SearchDomainItem{
					getSearchDomainItem(t, "test1.com", "HOST", "REQUEST_ACCEPTED"),
					wildcard,
				}).Times(1)
			},
			importID:    "test1.com,test2.com",
			expectError: regexp.MustCompile(`(?s)Error verifying domains.+the domain 'test2.com' was not found or it was found only without FQDN.+validationLevel`),
		},
		"import with multiple domains with same domainName but with validation scope was already provided": {
			config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/same_domain_name.tf"),
			expectError: regexp.MustCompile(`(?s)Error parsing import ID.+domain 'test1.com' was already provided in the importID with validation.+scope: HOST - such combination is not allowed. Please remove duplicate domain.+entries`),
			importID:    "test1.com:HOST,test1.com",
		},
		"import with multiple domains with same domainName but without validation scope was already provided": {
			config:      testutils.LoadFixtureString(t, "testdata/TestResDOMDomains/same_domain_name.tf"),
			expectError: regexp.MustCompile(`(?s)Error parsing import ID.+domain 'test1.com' with validation scope 'HOST' was already provided in the.+importID without validation scope - such combination is not allowed. Please.+remove duplicate domain entries`),
			importID:    "test1.com,test1.com:HOST",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := &domainownership.Mock{}
			if tc.init != nil {
				tc.init(m)
			}

			useDomainOwnership(m, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: tc.stateCheck,
							ImportStateId:    tc.importID,
							ImportState:      true,
							ResourceName:     "akamai_property_domainownership_domains.test",
							Config:           tc.config,
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			m.AssertExpectations(t)
		})
	}
}

func validateAndSplit(domains ...any) []domainownership.Domain {
	if len(domains)%2 != 0 {
		panic("each domain must have a validation scope")
	}
	var domainList []domainownership.Domain
	for i := 0; i < len(domains); i += 2 {
		domainList = append(domainList, domainownership.Domain{
			DomainName:      domains[i].(string),
			ValidationScope: domains[i+1].(domainownership.ValidationScope),
		})
	}
	return domainList
}

func buildSearchDomainsRequestTest(domains ...any) domainownership.SearchDomainsRequest {
	result := domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: validateAndSplit(domains...),
		},
	}
	return result
}

func buildDeleteDomainsRequest(domains ...any) domainownership.DeleteDomainsRequest {
	result := domainownership.DeleteDomainsRequest{
		Domains: validateAndSplit(domains...),
	}
	return result
}

func buildAddDomainsRequest(domains ...any) domainownership.AddDomainsRequest {
	result := domainownership.AddDomainsRequest{
		Domains: validateAndSplit(domains...),
	}
	return result
}

func getSearchDomainItem(t *testing.T, domainName, validationScope, domainStatus string) domainownership.SearchDomainItem {
	result := domainownership.SearchDomainItem{
		DomainName:              domainName,
		DomainStatus:            domainStatus,
		ValidationScope:         validationScope,
		ValidationLevel:         "FQDN",
		AccountID:               ptr.To("ACC123"),
		ValidationRequestedBy:   ptr.To("someone"),
		ValidationRequestedDate: ptr.To(tst.NewTimeFromString(t, "2024-01-01T00:00:00Z")),
		ValidationChallenge: &domainownership.ValidationChallenge{
			CnameRecord: domainownership.CnameRecord{
				Name:   fmt.Sprintf("challenge_%s", domainName),
				Target: fmt.Sprintf("target_%s", domainName),
			},
			TXTRecord: domainownership.TXTRecord{
				Name:  fmt.Sprintf("challenge_%s", domainName),
				Value: fmt.Sprintf("value_%s", domainName),
			},
			ExpirationDate: tst.NewTimeFromString(t, "2024-01-10T00:00:00Z"),
		},
	}
	if domainStatus == "VALIDATED" {
		result.ValidationMethod = ptr.To("DNS_CNAME")
		result.ValidationCompletedDate = ptr.To(tst.NewTimeFromString(t, "2024-01-02T00:00:00Z"))
	}
	if validationScope == "HOST" {
		result.ValidationChallenge.HTTPFile = &domainownership.HTTPFile{
			Path:        "/challenge_path",
			Content:     fmt.Sprintf("file_content_%s", domainName),
			ContentType: "text/plain",
		}
		result.ValidationChallenge.HTTPRedirect = &domainownership.HTTPRedirect{
			From: "/from_path",
			To:   fmt.Sprintf("/to_path_%s", domainName),
		}
	}

	return result
}

func getAddDomainSuccess(t *testing.T, domainName, validationScope string) domainownership.AddDomainSuccess {
	searchResult := getSearchDomainItem(t, domainName, validationScope, "REQUEST_ACCEPTED")
	return domainownership.AddDomainSuccess{
		DomainName:              domainName,
		ValidationScope:         validationScope,
		AccountID:               *searchResult.AccountID,
		DomainStatus:            searchResult.DomainStatus,
		ValidationMethod:        searchResult.ValidationMethod,
		ValidationRequestedBy:   *searchResult.ValidationRequestedBy,
		ValidationRequestedDate: *searchResult.ValidationRequestedDate,
		ValidationCompletedDate: searchResult.ValidationCompletedDate,
		ValidationChallenge:     *searchResult.ValidationChallenge,
	}
}

func mockBasicSearchDomains(mock *domainownership.Mock, request domainownership.SearchDomainsRequest, respDomains []domainownership.SearchDomainItem) *mock.Call {
	return mock.On("SearchDomains", testutils.MockContext, request).Return(&domainownership.SearchDomainsResponse{Domains: respDomains}, nil).Once()
}

func mockAddDomains(mock *domainownership.Mock, request domainownership.AddDomainsRequest, resp *domainownership.AddDomainsResponse, err error) *mock.Call {
	return mock.On("AddDomains", testutils.MockContext, request).Return(resp, err).Once()
}

func mockDeleteDomains(mock *domainownership.Mock, request domainownership.DeleteDomainsRequest, err error) *mock.Call {
	return mock.On("DeleteDomains", testutils.MockContext, request).Return(err).Once()
}

func generateLongImportID(validationScope string) string {
	var importID strings.Builder
	for i := 0; i < 1001; i++ {
		if i > 0 {
			importID.WriteString(",")
		}
		importID.WriteString(fmt.Sprintf("test%d.example.com", i))
		if validationScope != "" {
			importID.WriteString(":" + validationScope)
		}
	}
	return importID.String()
}
