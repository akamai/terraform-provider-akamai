package property

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	validationTestData struct {
		create               map[domainKey]domainDetails
		postCreateValidation map[domainKey]domainDetails
		update               map[domainKey]domainDetails
		postUpdateValidation map[domainKey]domainDetails
	}
)

func getMinCreate() validationTestData {
	return validationTestData{
		create: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("REQUEST_ACCEPTED", "FQDN", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("REQUEST_ACCEPTED", "FQDN", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("REQUEST_ACCEPTED", "FQDN", nil),
		},
		postCreateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "FQDN", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "FQDN", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "FQDN", nil),
		},
	}
}

func get101Domains() validationTestData {
	createMap := make(map[domainKey]domainDetails)
	postCreateValidationMap := make(map[domainKey]domainDetails)
	for i := range 101 {
		createMap[newDomainKey(fmt.Sprintf("test%d.example.com", i), "HOST")] =
			newDomainDetails("REQUEST_ACCEPTED", "FQDN", nil)
		postCreateValidationMap[newDomainKey(fmt.Sprintf("test%d.example.com", i), "HOST")] =
			newDomainDetails("VALIDATED", "FQDN", nil)
	}
	return validationTestData{
		create:               createMap,
		postCreateValidation: postCreateValidationMap,
	}
}

func get501Domains() validationTestData {
	createMap := make(map[domainKey]domainDetails)
	postCreateValidationMap := make(map[domainKey]domainDetails)
	updateMap := make(map[domainKey]domainDetails)
	postUpdateValidationMap := make(map[domainKey]domainDetails)
	for i := range 501 {
		createMap[newDomainKey(fmt.Sprintf("test%d.example.com", i), "HOST")] =
			newDomainDetails("REQUEST_ACCEPTED", "FQDN", nil)
		postCreateValidationMap[newDomainKey(fmt.Sprintf("test%d.example.com", i), "HOST")] =
			newDomainDetails("VALIDATED", "FQDN", nil)
		updateMap[newDomainKey(fmt.Sprintf("update-test%d.example.com", i), "HOST")] =
			newDomainDetails("REQUEST_ACCEPTED", "FQDN", nil)
		postUpdateValidationMap[newDomainKey(fmt.Sprintf("update-test%d.example.com", i), "HOST")] =
			newDomainDetails("VALIDATED", "FQDN", nil)
	}
	return validationTestData{
		create:               createMap,
		postCreateValidation: postCreateValidationMap,
		update:               updateMap,
		postUpdateValidation: postUpdateValidationMap,
	}
}

func getFullCreate() validationTestData {
	return validationTestData{
		create: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("REQUEST_ACCEPTED", "FQDN", ptr.To("DNS_CNAME")),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("REQUEST_ACCEPTED", "FQDN", ptr.To("DNS_TXT")),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("REQUEST_ACCEPTED", "FQDN", ptr.To("HTTP")),
		},
		postCreateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "FQDN", ptr.To("DNS_CNAME")),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "FQDN", ptr.To("DNS_TXT")),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "FQDN", ptr.To("HTTP")),
		},
	}
}

func getMinCreateAddOne() validationTestData {
	return validationTestData{
		create: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("REQUEST_ACCEPTED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("REQUEST_ACCEPTED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("REQUEST_ACCEPTED", "", nil),
		},
		postCreateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "", nil),
		},
		update: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test4.example.com", "HOST"):     newDomainDetails("REQUEST_ACCEPTED", "", nil),
		},
		postUpdateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test4.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
		},
	}
}

func getMinCreateRemoveOne() validationTestData {
	return validationTestData{
		create: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("REQUEST_ACCEPTED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("REQUEST_ACCEPTED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("REQUEST_ACCEPTED", "", nil),
		},
		postCreateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "", nil),
		},
		postUpdateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"): newDomainDetails("VALIDATED", "", nil),
		},
	}
}

func getMinCreateAddOneRemoveOne() validationTestData {
	return validationTestData{
		create: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("REQUEST_ACCEPTED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("REQUEST_ACCEPTED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("REQUEST_ACCEPTED", "", nil),
		},
		postCreateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "", nil),
		},
		update: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test4.example.com", "WILDCARD"): newDomainDetails("REQUEST_ACCEPTED", "", nil),
		},
		postUpdateValidation: map[domainKey]domainDetails{
			newDomainKey("test1.example.com", "HOST"):     newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("VALIDATED", "", nil),
			newDomainKey("test4.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "", nil),
		},
	}
}

func TestDomainOwnershipValidationResource(t *testing.T) {
	searchInterval = 1 * time.Millisecond

	minCreateChecker := test.NewStateChecker("akamai_property_domainownership_validation.test").
		CheckEqual("domains.#", "3").
		CheckEqual("domains.0.domain_name", "test1.example.com").
		CheckEqual("domains.0.validation_scope", "HOST").
		CheckMissing("domains.0.validation_method").
		CheckEqual("domains.1.domain_name", "test2.example.com").
		CheckEqual("domains.1.validation_scope", "DOMAIN").
		CheckMissing("domains.1.validation_method").
		CheckEqual("domains.2.domain_name", "test3.example.com").
		CheckEqual("domains.2.validation_scope", "WILDCARD").
		CheckMissing("domains.2.validation_method")

	importChecker := test.NewImportChecker().
		CheckEqual("domains.#", "3").
		CheckEqual("domains.0.domain_name", "test1.example.com").
		CheckEqual("domains.0.validation_scope", "HOST").
		CheckMissing("domains.0.validation_method").
		CheckEqual("domains.1.domain_name", "test2.example.com").
		CheckEqual("domains.1.validation_scope", "DOMAIN").
		CheckMissing("domains.1.validation_method").
		CheckEqual("domains.2.domain_name", "test3.example.com").
		CheckEqual("domains.2.validation_scope", "WILDCARD").
		CheckMissing("domains.2.validation_method")

	domains101Checker := test.NewStateChecker("akamai_property_domainownership_validation.test").
		CheckEqual("domains.#", "101")

	tests := map[string]struct {
		init     func(*domainownership.Mock, validationTestData)
		mockData validationTestData
		steps    []resource.TestStep
	}{
		"create with 3 domains - no polling": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read before destroy
				mockSearchDomains(m, mockData.postCreateValidation)
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
			},
		},
		"create with 3 domains - with polling": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				defaultPollTimeout = 30 * time.Minute
				// Create
				mockSearchDomains(m, mockData.create)
				pending := map[domainKey]domainDetails{
					newDomainKey("test1.example.com", "HOST"):     newDomainDetails("PENDING", "FQDN", nil),
					newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("PENDING", "FQDN", nil),
					newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("PENDING", "FQDN", nil),
				}
				mockValidateDomains(m, mockData.create, pending)
				mockSearchDomains(m, mockData.postCreateValidation)
				// Read before destroy
				mockSearchDomains(m, mockData.postCreateValidation)
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
			},
		},
		"create with 3 domains - one domain already validated": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Simulate one domain already validated before creation.
				mockData.create[newDomainKey("test2.example.com", "DOMAIN")] = newDomainDetails("VALIDATED", "FQDN", nil)
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read before destroy
				mockSearchDomains(m, mockData.postCreateValidation)
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
			},
		},
		"create with 3 domains - all domains already validated": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create - use already validated statuses.
				mockSearchDomains(m, mockData.postCreateValidation)
				// Read before destroy
				mockSearchDomains(m, mockData.postCreateValidation)
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
			},
		},
		"create with 3 domains - with validation methods": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read before destroy
				mockSearchDomains(m, mockData.postCreateValidation)
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getFullCreate(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create_with_validation_methods.tf"),
					Check: minCreateChecker.
						CheckEqual("domains.0.validation_method", "DNS_CNAME").
						CheckEqual("domains.1.validation_method", "DNS_TXT").
						CheckEqual("domains.2.validation_method", "HTTP").
						Build(),
				},
			},
		},
		"create with 101 domains": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read before destroy
				mockSearchDomains(m, mockData.postCreateValidation)
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: get101Domains(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create_101.tf"),
					Check:  domains101Checker.Build(),
				},
			},
		},
		"create with 101 domains - poll for 100 domains": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				defaultPollTimeout = 30 * time.Minute
				// Create
				mockSearchDomains(m, mockData.create)
				// Simulate 100 domains are in PENDING status.
				pending := make(map[domainKey]domainDetails, len(mockData.postCreateValidation))
				maps.Copy(pending, mockData.postCreateValidation)
				for i := range 100 {
					pending[newDomainKey(fmt.Sprintf("test%d.example.com", i), "HOST")] =
						newDomainDetails("PENDING", "FQDN", nil)
				}
				// Second call to SearchDomains returns 100 domains as VALIDATED.
				searchedCompleted := make(map[domainKey]domainDetails)
				for i := range 100 {
					searchedCompleted[domainKey{
						domainName:      fmt.Sprintf("test%d.example.com", i),
						validationScope: string(domainownership.ValidationScopeHost),
					}] = domainDetails{
						validationStatus: "VALIDATED",
					}
				}
				mockValidateDomains(m, mockData.create, pending)
				mockSearchDomains(m, searchedCompleted)
				// Read before destroy
				mockSearchDomains(m, mockData.postCreateValidation)
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: get101Domains(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create_101.tf"),
					Check:  domains101Checker.Build(),
				},
			},
		},
		"create with 3 domains - resource drift by invalidation of a domain in read": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read
				mockSearchDomains(m, mockData.postCreateValidation)
				// Refresh with drift - one domain got invalidated outside of TF
				var searchDomainsBody []domainownership.Domain
				for k := range mockData.create {
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScope(k.validationScope),
					})
				}
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:       "test1.example.com",
						ValidationScope:  "HOST",
						ValidationMethod: nil,
						DomainStatus:     "VALIDATED",
						ValidationLevel:  "FQDN",
					},
					{
						DomainName:       "test2.example.com",
						ValidationScope:  "DOMAIN",
						ValidationMethod: nil,
						DomainStatus:     "INVALIDATED",
						ValidationLevel:  "FQDN",
					},
					{
						DomainName:       "test3.example.com",
						ValidationScope:  "WILDCARD",
						ValidationMethod: nil,
						DomainStatus:     "VALIDATED",
						ValidationLevel:  "FQDN",
					},
				}
				sortDomains(searchDomainsBody)
				sort.Slice(domainsResponse, func(i, j int) bool {
					if domainsResponse[i].DomainName == domainsResponse[j].DomainName {
						return domainsResponse[i].ValidationScope < domainsResponse[j].ValidationScope
					}
					return domainsResponse[i].DomainName < domainsResponse[j].DomainName
				})

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()

				drift := make(map[domainKey]domainDetails)
				maps.Copy(drift, mockData.postCreateValidation)
				delete(drift, domainKey{
					domainName:      "test2.example.com",
					validationScope: string(domainownership.ValidationScopeDomain),
				})
				// Read before destroy
				mockSearchDomains(m, drift)
				// Delete
				mockSearchDomains(m, drift)
				mockInvalidateDomains(m, drift)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Check: minCreateChecker.
						CheckEqual("domains.#", "2").
						CheckEqual("domains.1.validation_scope", "WILDCARD").
						CheckEqual("domains.1.domain_name", "test3.example.com").
						CheckMissing("domains.2.domain_name").
						CheckMissing("domains.2.validation_scope").
						CheckMissing("domains.2.validation_method").
						Build(),
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"create with 3 domains - resource drift - domains have different than VALIDATED status in read": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read
				mockSearchDomains(m, mockData.postCreateValidation)
				// Refresh with drift - domains have different than VALIDATED statuses
				var searchDomainsBody []domainownership.Domain
				for k := range mockData.create {
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScope(k.validationScope),
					})
				}
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:       "test1.example.com",
						ValidationScope:  "HOST",
						ValidationMethod: nil,
						DomainStatus:     "REQUEST_ACCEPTED",
						ValidationLevel:  "FQDN",
					},
					{
						DomainName:       "test2.example.com",
						ValidationScope:  "DOMAIN",
						ValidationMethod: nil,
						DomainStatus:     "INVALIDATED",
						ValidationLevel:  "FQDN",
					},
					{
						DomainName:       "test3.example.com",
						ValidationScope:  "WILDCARD",
						ValidationMethod: nil,
						DomainStatus:     "TOKEN_EXPIRED",
						ValidationLevel:  "FQDN",
					},
				}
				sortDomains(searchDomainsBody)
				sort.Slice(domainsResponse, func(i, j int) bool {
					if domainsResponse[i].DomainName == domainsResponse[j].DomainName {
						return domainsResponse[i].ValidationScope < domainsResponse[j].ValidationScope
					}
					return domainsResponse[i].DomainName < domainsResponse[j].DomainName
				})

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"create with 3 domains - update by adding one domain": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read x2
				mockSearchDomains(m, mockData.postCreateValidation).Twice()
				// Update
				mockSearchDomains(m, mockData.update)
				mockValidateDomains(m, mockData.update, mockData.postUpdateValidation)
				// Read after update
				mockSearchDomains(m, mockData.postUpdateValidation)
				// Delete
				mockSearchDomains(m, mockData.postUpdateValidation)
				mockInvalidateDomains(m, mockData.postUpdateValidation)
			},
			mockData: getMinCreateAddOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/add_1.tf"),
					Check: minCreateChecker.
						CheckEqual("domains.#", "4").
						CheckEqual("domains.3.domain_name", "test4.example.com").
						CheckEqual("domains.3.validation_scope", "HOST").
						CheckMissing("domains.3.validation_method").
						Build(),
				},
			},
		},
		"create with 501 domains - update by adding 501 new domains in the place of old, 2 requests are needed": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read x2
				mockSearchDomains(m, mockData.postCreateValidation).Twice()
				// Update
				domainsFromCreate := mockData.postCreateValidation
				domainsFromUpdate := mockData.update
				allDomains := make(map[domainKey]domainDetails, len(domainsFromCreate)+len(domainsFromUpdate))
				maps.Copy(allDomains, domainsFromCreate)
				maps.Copy(allDomains, domainsFromUpdate)

				var searchDomainsBody []domainownership.Domain
				var domainsResponse []domainownership.SearchDomainItem
				for k, v := range allDomains {
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScope(k.validationScope),
					})
					domainsResponse = append(domainsResponse, domainownership.SearchDomainItem{
						DomainName:       k.domainName,
						ValidationScope:  k.validationScope,
						ValidationMethod: v.validationMethod,
						ValidationLevel:  v.validationLevel,
						DomainStatus:     v.validationStatus,
					})
				}

				sortDomains(searchDomainsBody)
				sort.Slice(domainsResponse, func(i, j int) bool {
					if domainsResponse[i].DomainName == domainsResponse[j].DomainName {
						return domainsResponse[i].ValidationScope < domainsResponse[j].ValidationScope
					}
					return domainsResponse[i].DomainName < domainsResponse[j].DomainName
				})
				// fetchDomainsFromAPI requires 2 API calls, first with 1000 domains.
				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody[:1000],
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse[:1000],
				}, nil).Once()
				// Second API call with remaining 2 domains.
				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody[1000:],
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse[1000:],
				}, nil).Once()
				mockInvalidateDomains(m, mockData.postCreateValidation)
				mockValidateDomains(m, mockData.update, mockData.postUpdateValidation)
				// Read after update
				mockSearchDomains(m, mockData.postUpdateValidation)
				// Delete
				mockSearchDomains(m, mockData.postUpdateValidation)
				mockInvalidateDomains(m, mockData.postUpdateValidation)
			},
			mockData: get501Domains(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create_501.tf"),
					Check:  domains101Checker.CheckEqual("domains.#", "501").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/update_501.tf"),
					Check:  domains101Checker.CheckEqual("domains.#", "501").Build(),
				},
			},
		},
		"create with 3 domains - update by adding one domain which is already VALIDATED": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read x2
				mockSearchDomains(m, mockData.postCreateValidation).Twice()
				// Update - simulate the added domain is already validated.
				mockSearchDomains(m, mockData.postUpdateValidation)
				// Read after update
				mockSearchDomains(m, mockData.postUpdateValidation)
				// Delete
				mockSearchDomains(m, mockData.postUpdateValidation)
				mockInvalidateDomains(m, mockData.postUpdateValidation)
			},
			mockData: getMinCreateAddOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/add_1.tf"),
					Check: minCreateChecker.
						CheckEqual("domains.#", "4").
						CheckEqual("domains.3.domain_name", "test4.example.com").
						CheckEqual("domains.3.validation_scope", "HOST").
						CheckMissing("domains.3.validation_method").
						Build(),
				},
			},
		},
		"create with 3 domains - update by removing one domain": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read x2
				mockSearchDomains(m, mockData.postCreateValidation).Twice()
				// Update
				mockSearchDomains(m, mockData.postCreateValidation)
				toInvalidate := map[domainKey]domainDetails{
					newDomainKey("test3.example.com", string(domainownership.ValidationScopeWildcard)): newDomainDetails("VALIDATED", "FQDN", nil),
				}
				mockInvalidateDomains(m, toInvalidate)
				// Read after update
				mockSearchDomains(m, mockData.postUpdateValidation)
				// Delete
				mockSearchDomains(m, mockData.postUpdateValidation)
				mockInvalidateDomains(m, mockData.postUpdateValidation)
			},
			mockData: getMinCreateRemoveOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/remove_1.tf"),
					Check: minCreateChecker.
						CheckEqual("domains.#", "2").
						CheckMissing("domains.2.domain_name").
						CheckMissing("domains.2.validation_scope").
						CheckMissing("domains.2.validation_method").
						CheckMissing("domains.3.domain_name").
						CheckMissing("domains.3.validation_scope").
						CheckMissing("domains.3.validation_method").
						Build(),
				},
			},
		},
		"create with 3 domains - update by removing one domain and adding one domain": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read x2
				mockSearchDomains(m, mockData.postCreateValidation).Twice()
				// Update
				domainsToSearch := make(map[domainKey]domainDetails)
				maps.Copy(domainsToSearch, mockData.postCreateValidation)
				domainsToSearch[newDomainKey("test4.example.com", string(domainownership.ValidationScopeWildcard))] = domainDetails{}
				mockSearchDomains(m, domainsToSearch)

				toInvalidate := map[domainKey]domainDetails{
					newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("VALIDATED", "FQDN", nil),
				}
				mockInvalidateDomains(m, toInvalidate)
				mockValidateDomains(m, mockData.update, mockData.postUpdateValidation)
				// Read after update
				mockSearchDomains(m, mockData.postUpdateValidation)
				// Delete
				mockSearchDomains(m, mockData.postUpdateValidation)
				mockInvalidateDomains(m, mockData.postUpdateValidation)
			},
			mockData: getMinCreateAddOneRemoveOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/add_1_remove_1.tf"),
					Check: minCreateChecker.
						CheckEqual("domains.#", "3").
						CheckEqual("domains.2.domain_name", "test4.example.com").
						CheckEqual("domains.2.validation_scope", "WILDCARD").
						CheckMissing("domains.2.validation_method").
						Build(),
				},
			},
		},
		"import 3 unique domains - provide all import keys and return matching domains": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Import
				mockSearchDomains(m, mockData.postCreateValidation)
				// Read
				mockSearchDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId:    "test1.example.com:HOST,test2.example.com:DOMAIN,test3.example.com:WILDCARD",
					ImportStateCheck: importChecker.Build(),
					ImportState:      true,
					ResourceName:     "akamai_property_domainownership_validation.test",
				},
			},
		},
		"import 3 unique domains - provide only domain names - every domain is unique, match by FQDN": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Import
				var searchDomainsBody []domainownership.Domain
				var domainsResponse []domainownership.SearchDomainItem
				for k, v := range mockData.postCreateValidation {
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScopeHost,
					})
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScopeWildcard,
					})
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScopeDomain,
					})
					domainsResponse = append(domainsResponse, domainownership.SearchDomainItem{
						DomainName:       k.domainName,
						ValidationScope:  k.validationScope,
						ValidationMethod: v.validationMethod,
						DomainStatus:     v.validationStatus,
						ValidationLevel:  v.validationLevel,
					})
				}
				sortDomains(searchDomainsBody)

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
				// Read
				mockSearchDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId:    "test1.example.com,test2.example.com,test3.example.com",
					ImportStateCheck: importChecker.Build(),
					ImportState:      true,
					ResourceName:     "akamai_property_domainownership_validation.test",
				},
			},
		},
		"import 3 domains - provide only domain names - match by FQDN when other validation levels are present": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Import
				var searchDomainsBody []domainownership.Domain
				var domainsResponse []domainownership.SearchDomainItem
				for k, v := range mockData.postCreateValidation {
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScopeHost,
					})
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScopeWildcard,
					})
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScopeDomain,
					})
					domainsResponse = append(domainsResponse, domainownership.SearchDomainItem{
						DomainName:       k.domainName,
						ValidationScope:  k.validationScope,
						ValidationMethod: v.validationMethod,
						DomainStatus:     v.validationStatus,
						ValidationLevel:  "FQDN",
					})
				}

				sortDomains(searchDomainsBody)

				domainsResponse = append(domainsResponse, domainownership.SearchDomainItem{
					DomainName:      "test2.example.com",
					ValidationScope: string(domainownership.ValidationScopeHost),
					ValidationLevel: "ROOT/WILDCARD",
					DomainStatus:    "REQUEST_ACCEPTED",
				})
				domainsResponse = append(domainsResponse, domainownership.SearchDomainItem{
					DomainName:      "test2.example.com",
					ValidationScope: string(domainownership.ValidationScopeWildcard),
					ValidationLevel: "ROOT/WILDCARD",
					DomainStatus:    "REQUEST_ACCEPTED",
				})

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
				// Read
				mockSearchDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId:    "test1.example.com,test2.example.com,test3.example.com",
					ImportStateCheck: importChecker.Build(),
					ImportState:      true,
					ResourceName:     "akamai_property_domainownership_validation.test",
				},
			},
		},
		"expect error - import - invalid number of import parts": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com:HOST:FQDN",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`invalid import ID format. Expected format is a list of: domain\[:scope\], got(\n|.)test1.example.com:HOST:FQDN`),
				},
			},
		},
		"expect error - import - duplicated domains": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com,test1.example.com",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`domain 'test1.example.com' was already provided in the importID. Please(\n|.)remove duplicate domain entries`),
				},
			},
		},
		"expect error - import - duplicated domains with validation scope": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com:HOST,test1.example.com:HOST",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`domain 'test1.example.com' with validation scope 'HOST' was already provided(\n|.)in the importID. Please remove duplicate domain entries`),
				},
			},
		},
		"expect error - import - invalid validation scope": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com:INVALID",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`invalid validation scope 'INVALID' for domain 'test1.example.com'. Expected(\n|.)one of 'HOST', 'WILDCARD', 'DOMAIN'`),
				},
			},
		},
		"expect error - import - inconclusive domain without validation scope, multiple validation scope exists": {
			init: func(m *domainownership.Mock, _ validationTestData) {
				// Import
				searchDomainsBody := []domainownership.Domain{
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeHost,
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeWildcard,
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeDomain,
					},
				}
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeHost),
						ValidationLevel: "FQDN",
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeDomain),
						ValidationLevel: "FQDN",
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeWildcard),
						ValidationLevel: "FQDN",
					},
				}

				sortDomains(searchDomainsBody)

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`the domain 'test1.example.com' exists with multiple validation scopes. Please(\n|.)re-import specifying the validation scope for the domain`),
				},
			},
		},
		"expect error - import - domain without validation scope not found": {
			init: func(m *domainownership.Mock, _ validationTestData) {
				// Import
				searchDomainsBody := []domainownership.Domain{
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeHost,
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeWildcard,
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeDomain,
					},
				}
				domainsResponse := []domainownership.SearchDomainItem{}

				sortDomains(searchDomainsBody)

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`the domain 'test1.example.com' was not found`),
				},
			},
		},
		"expect error - import - domain with validation scope provided, no FQDN validation level found": {
			init: func(m *domainownership.Mock, _ validationTestData) {
				// Import
				searchDomainsBody := []domainownership.Domain{
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeHost,
					},
				}
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeHost),
						ValidationLevel: "ROOT/WILDCARD",
						DomainStatus:    "VALIDATED",
					},
				}
				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com:HOST",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`only domains with validation level FQDN can be imported, the requested domain(\n|.)'test1.example.com' with validation scope 'HOST' has validation level(\n|.)'ROOT/WILDCARD'`),
				},
			},
		},
		"expect error - import - domain with validation scope not found": {
			init: func(m *domainownership.Mock, _ validationTestData) {
				// Import
				searchDomainsBody := []domainownership.Domain{
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeHost,
					},
				}
				domainsResponse := []domainownership.SearchDomainItem{}

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com:HOST",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`the domain 'test1.example.com' with validation scope 'HOST' was not found`),
				},
			},
		},
		"expect error - import - domain with validation scope is INVALIDATED": {
			init: func(m *domainownership.Mock, _ validationTestData) {
				// Import
				searchDomainsBody := []domainownership.Domain{
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeHost,
					},
				}
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeHost),
						DomainStatus:    "INVALIDATED",
					},
				}

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com:HOST",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`the domain 'test1.example.com' with validation scope 'HOST' is in(\n|.)'INVALIDATED' status and cannot be imported`),
				},
			},
		},
		"expect error - import - domain without validation scope is in REQUEST_ACCEPTED status": {
			init: func(m *domainownership.Mock, _ validationTestData) {
				// Import
				searchDomainsBody := []domainownership.Domain{
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeDomain,
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeHost,
					},
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeWildcard,
					},
				}
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeHost),
						ValidationLevel: "FQDN",
						DomainStatus:    "REQUEST_ACCEPTED",
					},
				}

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`the domain 'test1.example.com' is in 'REQUEST_ACCEPTED' status and cannot be(\n|.)imported`),
				},
			},
		},
		"expect error - import - API error": {
			init: func(m *domainownership.Mock, _ validationTestData) {
				// Import
				searchDomainsBody := []domainownership.Domain{
					{
						DomainName:      "test1.example.com",
						ValidationScope: domainownership.ValidationScopeHost,
					},
				}
				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(nil, fmt.Errorf("API error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ImportStateId: "test1.example.com:HOST",
					ImportState:   true,
					ResourceName:  "akamai_property_domainownership_validation.test",
					ExpectError:   regexp.MustCompile(`API error`),
				},
			},
		},
		"expect error - create - timeout exceeded": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				defaultPollTimeout = 1 * time.Microsecond
				searchInterval = 5 * time.Second
				// Create
				mockSearchDomains(m, mockData.create)
				pending := map[domainKey]domainDetails{
					newDomainKey("test1.example.com", "HOST"):     newDomainDetails("PENDING", "FQDN", nil),
					newDomainKey("test2.example.com", "DOMAIN"):   newDomainDetails("PENDING", "FQDN", nil),
					newDomainKey("test3.example.com", "WILDCARD"): newDomainDetails("PENDING", "FQDN", nil),
				}
				mockValidateDomains(m, mockData.create, pending)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:       minCreateChecker.Build(),
					ExpectError: regexp.MustCompile(`Error: Timeout while waiting for domain validation`),
				},
			},
		},
		"expect error - create - one domain not found in API": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				var searchDomainsBody []domainownership.Domain
				for k := range mockData.create {
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScope(k.validationScope),
					})
				}

				// Simulate one domain not found in API.
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeHost),
						DomainStatus:    "REQUEST_ACCEPTED",
					},
					{
						DomainName:      "test2.example.com",
						ValidationScope: string(domainownership.ValidationScopeDomain),
						DomainStatus:    "REQUEST_ACCEPTED",
					},
				}

				sortDomains(searchDomainsBody)

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ExpectError: regexp.MustCompile("domain test3.example.com with scope WILDCARD is not found in API"),
				},
			},
		},
		"expect error - create - one domain is INVALIDATED": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				oneDomainInvalidated := mockData
				oneDomainInvalidated.create[newDomainKey("test3.example.com", "WILDCARD")] =
					newDomainDetails("INVALIDATED", "FQDN", nil)
				// Create
				mockSearchDomains(m, oneDomainInvalidated.create)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ExpectError: regexp.MustCompile(`domain test3.example.com with scope WILDCARD is in INVALIDATED status, cannot(\n|.)validate`),
				},
			},
		},
		"expect error - create - one domain is TOKEN_EXPIRED": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				oneDomainTokenExpired := mockData
				oneDomainTokenExpired.create[newDomainKey("test3.example.com", "WILDCARD")] =
					newDomainDetails("TOKEN_EXPIRED", "FQDN", nil)
				// Create
				mockSearchDomains(m, oneDomainTokenExpired.create)
			},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					ExpectError: regexp.MustCompile(`domain test3.example.com with scope WILDCARD is in TOKEN_EXPIRED status,(\n|.)cannot validate`),
				},
			},
		},
		"expect error - create with 3 domains - update by adding one domain - timeout exceeded": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				defaultPollTimeout = 1 * time.Microsecond
				searchInterval = 5 * time.Second
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read x2
				mockSearchDomains(m, mockData.postCreateValidation).Twice()
				// Update
				mockSearchDomains(m, mockData.update)
				pending := map[domainKey]domainDetails{
					newDomainKey("test4.example.com", "HOST"): newDomainDetails("PENDING", "FQDN", nil),
				}
				mockValidateDomains(m, mockData.update, pending)
				// Delete - use postCreateValidation as update failed.
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreateAddOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/add_1.tf"),
					ExpectError: regexp.MustCompile(`Error: Timeout while waiting for domain validation`),
				},
			},
		},
		"expect error - update - one domain not found in API": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read
				mockSearchDomains(m, mockData.postCreateValidation).Twice()

				// Update
				var searchDomainsBody []domainownership.Domain
				for k := range mockData.update {
					searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
						DomainName:      k.domainName,
						ValidationScope: domainownership.ValidationScope(k.validationScope),
					})
				}

				// Simulate one domain not found in API.
				domainsResponse := []domainownership.SearchDomainItem{
					{
						DomainName:      "test1.example.com",
						ValidationScope: string(domainownership.ValidationScopeHost),
						DomainStatus:    "VALIDATED",
					},
					{
						DomainName:      "test2.example.com",
						ValidationScope: string(domainownership.ValidationScopeDomain),
						DomainStatus:    "VALIDATED",
					},
					{
						DomainName:      "test3.example.com",
						ValidationScope: string(domainownership.ValidationScopeWildcard),
						DomainStatus:    "VALIDATED",
					},
				}

				sortDomains(searchDomainsBody)

				m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
					IncludeAll: true,
					Body: domainownership.SearchDomainsBody{
						Domains: searchDomainsBody,
					},
				}).Return(&domainownership.SearchDomainsResponse{
					Domains: domainsResponse,
				}, nil).Once()
				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreateAddOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/add_1.tf"),
					ExpectError: regexp.MustCompile("domain test4.example.com with scope HOST is not found in API"),
				},
			},
		},
		"expect error - update - one domain is INVALIDATED": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read
				mockSearchDomains(m, mockData.postCreateValidation).Twice()

				// Update
				oneDomainInvalidated := mockData.update
				oneDomainInvalidated[newDomainKey("test4.example.com", "HOST")] =
					newDomainDetails("INVALIDATED", "FQDN", nil)
				mockSearchDomains(m, oneDomainInvalidated)

				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreateAddOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/add_1.tf"),
					ExpectError: regexp.MustCompile(`domain test4.example.com with scope HOST is in INVALIDATED status, cannot(\n|.)validate`),
				},
			},
		},
		"expect error - update - one domain is TOKEN_EXPIRED": {
			init: func(m *domainownership.Mock, mockData validationTestData) {
				// Create
				mockSearchDomains(m, mockData.create)
				mockValidateDomains(m, mockData.create, mockData.postCreateValidation)
				// Read
				mockSearchDomains(m, mockData.postCreateValidation).Twice()

				// Update
				oneDomainTokenExpired := mockData.update
				oneDomainTokenExpired[newDomainKey("test4.example.com", "HOST")] =
					newDomainDetails("TOKEN_EXPIRED", "FQDN", nil)
				mockSearchDomains(m, oneDomainTokenExpired)

				// Delete
				mockSearchDomains(m, mockData.postCreateValidation)
				mockInvalidateDomains(m, mockData.postCreateValidation)
			},
			mockData: getMinCreateAddOne(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/create.tf"),
					Check:  minCreateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/add_1.tf"),
					ExpectError: regexp.MustCompile(`domain test4.example.com with scope HOST is in TOKEN_EXPIRED status, cannot(\n|.)validate`),
				},
			},
		},
		"validation error - no domains": {
			init:     func(_ *domainownership.Mock, _ validationTestData) {},
			mockData: getMinCreate(),
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/no_domains.tf"),
					ExpectError: regexp.MustCompile(`The argument "domains" is required, but no definition was found.`),
				},
			},
		},
		"validation error - empty domains": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/empty_domains.tf"),
					ExpectError: regexp.MustCompile(`Attribute domains set must contain at least 1 elements and at most 1000(\n|.)elements, got: 0`),
				},
			},
		},
		"validation error - more than 1000 domains": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/too_many_domains.tf"),
					ExpectError: regexp.MustCompile(`Attribute domains set must contain at least 1 elements and at most 1000(\n|.)elements, got: 1001`),
				},
			},
		},
		"validation error - wrong validation scope": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/wrong_validation_scope.tf"),
					ExpectError: regexp.MustCompile(`validation_scope(\n|.)value must be one of: \["HOST" "WILDCARD" "DOMAIN"\], got: "WRONG"`),
				},
			},
		},
		"validation error - wrong validation method": {
			init: func(_ *domainownership.Mock, _ validationTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMValidation/wrong_validation_method.tf"),
					ExpectError: regexp.MustCompile(`.validation_method(\n|.)value must be one of: \["DNS_CNAME" "DNS_TXT" "HTTP"\], got: "WRONG"`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &domainownership.Mock{}

			if tc.init != nil {
				tc.init(client, tc.mockData)
			}

			useDomainOwnership(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func mockSearchDomains(m *domainownership.Mock, domains map[domainKey]domainDetails) *mock.Call {
	var searchDomainsBody []domainownership.Domain
	var domainsResponse []domainownership.SearchDomainItem
	for k, v := range domains {
		searchDomainsBody = append(searchDomainsBody, domainownership.Domain{
			DomainName:      k.domainName,
			ValidationScope: domainownership.ValidationScope(k.validationScope),
		})
		domainsResponse = append(domainsResponse, domainownership.SearchDomainItem{
			DomainName:       k.domainName,
			ValidationScope:  k.validationScope,
			ValidationMethod: v.validationMethod,
			ValidationLevel:  v.validationLevel,
			DomainStatus:     v.validationStatus,
		})
	}

	sortDomains(searchDomainsBody)
	sort.Slice(domainsResponse, func(i, j int) bool {
		if domainsResponse[i].DomainName == domainsResponse[j].DomainName {
			return domainsResponse[i].ValidationScope < domainsResponse[j].ValidationScope
		}
		return domainsResponse[i].DomainName < domainsResponse[j].DomainName
	})

	return m.On("SearchDomains", testutils.MockContext, domainownership.SearchDomainsRequest{
		IncludeAll: true,
		Body: domainownership.SearchDomainsBody{
			Domains: searchDomainsBody,
		},
	}).Return(&domainownership.SearchDomainsResponse{
		Domains: domainsResponse,
	}, nil).Once()
}

func mockValidateDomains(m *domainownership.Mock, domainsToValidate map[domainKey]domainDetails, responseDomains map[domainKey]domainDetails) {
	var validateDomainsBody []domainownership.ValidateDomainRequest
	var validateDomainResponse []domainownership.ValidateDomainResponse
	for k, v := range domainsToValidate {
		if v.validationStatus == "REQUEST_ACCEPTED" {
			var validationMethod *domainownership.ValidationMethod
			if v.validationMethod != nil {
				validationMethod = ptr.To(domainownership.ValidationMethod(*v.validationMethod))
			}
			validateDomainsBody = append(validateDomainsBody, domainownership.ValidateDomainRequest{
				DomainName:       k.domainName,
				ValidationScope:  domainownership.ValidationScope(k.validationScope),
				ValidationMethod: validationMethod,
			})
		}
	}

	sort.Slice(validateDomainsBody, func(i, j int) bool {
		if validateDomainsBody[i].DomainName == validateDomainsBody[j].DomainName {
			return validateDomainsBody[i].ValidationScope < validateDomainsBody[j].ValidationScope
		}
		return validateDomainsBody[i].DomainName < validateDomainsBody[j].DomainName
	})

	for k, v := range responseDomains {
		validateDomainResponse = append(validateDomainResponse, domainownership.ValidateDomainResponse{
			DomainName:      k.domainName,
			ValidationScope: k.validationScope,
			DomainStatus:    v.validationStatus,
		})
	}
	sort.Slice(validateDomainResponse, func(i, j int) bool {
		if validateDomainResponse[i].DomainName == validateDomainResponse[j].DomainName {
			return validateDomainResponse[i].ValidationScope < validateDomainResponse[j].ValidationScope
		}
		return validateDomainResponse[i].DomainName < validateDomainResponse[j].DomainName
	})

	var responses []domainownership.ValidateDomainsResponse
	for chunk := range slices.Chunk(validateDomainResponse, 100) {
		responses = append(responses, domainownership.ValidateDomainsResponse{Domains: chunk})
	}

	var requests []domainownership.ValidateDomainsRequest
	for chunk := range slices.Chunk(validateDomainsBody, 100) {
		requests = append(requests, domainownership.ValidateDomainsRequest{Domains: chunk})
	}

	for i, req := range requests {
		m.On("ValidateDomains", testutils.MockContext, req).Return(&responses[i], nil).Once()
	}
}

func mockInvalidateDomains(m *domainownership.Mock, domains map[domainKey]domainDetails) {
	var invalidateDomainsBody []domainownership.Domain
	var invalidateDomainResponse []domainownership.InvalidateDomainResponse
	for k := range domains {
		invalidateDomainsBody = append(invalidateDomainsBody, domainownership.Domain{
			DomainName:      k.domainName,
			ValidationScope: domainownership.ValidationScope(k.validationScope),
		})
		invalidateDomainResponse = append(invalidateDomainResponse, domainownership.InvalidateDomainResponse{
			DomainName:      k.domainName,
			ValidationScope: k.validationScope,
			DomainStatus:    "INVALIDATED",
		})
	}

	sortDomains(invalidateDomainsBody)

	m.On("InvalidateDomains", testutils.MockContext, domainownership.InvalidateDomainsRequest{
		Domains: invalidateDomainsBody,
	}).Return(&domainownership.InvalidateDomainsResponse{
		Domains: invalidateDomainResponse,
	}, nil).Once()
}

func newDomainKey(name string, scope string) domainKey {
	return domainKey{
		domainName:      name,
		validationScope: scope,
	}
}

func newDomainDetails(status, level string, method *string) domainDetails {
	return domainDetails{
		validationStatus: status,
		validationLevel:  level,
		validationMethod: method,
	}
}
