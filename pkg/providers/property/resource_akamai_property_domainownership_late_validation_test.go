package property

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDomainOwnershipLateValidationResource(t *testing.T) {
	t.Parallel()

	commonStateChecker := test.NewStateChecker("akamai_property_domainownership_late_validation.test").
		CheckEqual("property_id", "prp_123").
		CheckEqual("contract_id", "ctr_1").
		CheckEqual("version", "1").
		CheckEqual("group_id", "grp_1").
		CheckEqual("validation_method", "DNS_CNAME")

	importChecker := test.NewImportChecker().
		CheckEqual("property_id", "prp_123").
		CheckEqual("contract_id", "ctr_1").
		CheckEqual("version", "1").
		CheckEqual("group_id", "grp_1").
		CheckEqual("validation_method", "DNS_CNAME")

	importCheckerWithoutPrefix := test.NewImportChecker().
		CheckEqual("property_id", "123").
		CheckEqual("contract_id", "1").
		CheckEqual("version", "1").
		CheckEqual("group_id", "1").
		CheckEqual("validation_method", "DNS_CNAME")

	key := domainOwnershipLateValidationResourceModel{
		PropertyID: types.StringValue("prp_123"),
		Version:    types.Int64Value(1),
		ContractID: types.StringValue("ctr_1"),
		GroupID:    types.StringValue("grp_1"),
	}

	tests := map[string]struct {
		init  func(*mockProperty, *domainownership.Mock)
		steps []resource.TestStep
	}{
		"create with all domains already validated": {
			init: func(m *mockProperty, _ *domainownership.Mock) {
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from1.test.domain",
							CnameTo:              "to1.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-12-31T23:59:59Z"),
								Status:                   "VALIDATED",
								ValidationCname: &papi.ValidationCname{
									Hostname: "example.com",
									Target:   "example.akamai.net",
								},
								ValidationHTTP: &papi.ValidationHTTP{
									FileContentMethod: papi.FileContentMethod{
										Body: "test-body",
										URL:  "test-url",
									},
									RedirectMethod: papi.RedirectMethod{
										HTTPRedirectFrom: "redirect-from",
										HTTPRedirectTo:   "redirect-to",
									},
								},
								ValidationTXT: &papi.ValidationTXT{
									ChallengeToken: "test-token",
									Hostname:       "test-example.com",
								},
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from2.test.domain",
							CnameTo:              "to2.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-12-31T23:59:59Z"),
								Status:                   "VALIDATED",
							},
						},
					},
				}
				// read, create
				m.mockGetPropertyVersionHostnames().Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"create with some domains requires validation, succeeds immediately": {
			init: func(m *mockProperty, domMock *domainownership.Mock) {
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from1.test.domain",
							CnameTo:              "to1.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-12-31T23:59:59Z"),
								Status:                   "VALIDATED",
								ValidationCname: &papi.ValidationCname{
									Hostname: "example.com",
									Target:   "example.akamai.net",
								},
								ValidationHTTP: &papi.ValidationHTTP{
									FileContentMethod: papi.FileContentMethod{
										Body: "test-body",
										URL:  "test-url",
									},
									RedirectMethod: papi.RedirectMethod{
										HTTPRedirectFrom: "redirect-from",
										HTTPRedirectTo:   "redirect-to",
									},
								},
								ValidationTXT: &papi.ValidationTXT{
									ChallengeToken: "test-token",
									Hostname:       "test-example.com",
								},
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from2.test.domain",
							CnameTo:              "to2.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-10-20T23:59:59Z"),
								Status:                   "PENDING",
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from3.test.domain",
							CnameTo:              "to3.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-11-15T23:59:59Z"),
								Status:                   "NOT_VALIDATED",
							},
						},
					},
				}
				// read
				m.mockGetPropertyVersionHostnames()

				req := domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from2.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
						{
							DomainName:       "from3.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				resp := domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from2.test.domain",
							DomainStatus:    "VALIDATED",
							ValidationScope: "HOST",
						},
						{
							DomainName:      "from3.test.domain",
							DomainStatus:    "VALIDATED",
							ValidationScope: "HOST",
						},
					},
				}
				// create
				domMock.On("ValidateDomains", testutils.MockContext, req).Return(&resp, nil).Once()
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from1.test.domain",
							CnameTo:              "to1.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-12-31T23:59:59Z"),
								Status:                   "VALIDATED",
								ValidationCname: &papi.ValidationCname{
									Hostname: "example.com",
									Target:   "example.akamai.net",
								},
								ValidationHTTP: &papi.ValidationHTTP{
									FileContentMethod: papi.FileContentMethod{
										Body: "test-body",
										URL:  "test-url",
									},
									RedirectMethod: papi.RedirectMethod{
										HTTPRedirectFrom: "redirect-from",
										HTTPRedirectTo:   "redirect-to",
									},
								},
								ValidationTXT: &papi.ValidationTXT{
									ChallengeToken: "test-token",
									Hostname:       "test-example.com",
								},
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from2.test.domain",
							CnameTo:              "to2.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-10-20T23:59:59Z"),
								Status:                   "VALIDATED",
							},
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from3.test.domain",
							CnameTo:              "to3.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{
								ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-11-15T23:59:59Z"),
								Status:                   "VALIDATED",
							},
						},
					},
				}

				//read
				m.mockGetPropertyVersionHostnames()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"create with one domain requiring polling": {
			init: func(m *mockProperty, domMock *domainownership.Mock) {
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				// read
				m.mockGetPropertyVersionHostnames()

				validateReq := domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from1.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				validateResp := domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from1.test.domain",
							DomainStatus:    "PENDING",
							ValidationScope: "HOST",
						},
					},
				}
				// create
				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(&validateResp, nil).Once()

				// create
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				// create
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				// read
				m.mockGetPropertyVersionHostnames()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					Check: resource.ComposeTestCheckFunc(
						commonStateChecker.Build(),
					),
				},
			},
		},
		"create with polling timeout": {
			init: func(m *mockProperty, domMock *domainownership.Mock) {
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				// read
				m.mockGetPropertyVersionHostnames()

				validateReq := domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from1.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				validateResp := domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from1.test.domain",
							DomainStatus:    "PENDING",
							ValidationScope: "HOST",
						},
					},
				}
				// create
				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(&validateResp, nil).Once()

				// create
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames().Times(0)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config_with_timeout.tf"),
					ExpectError: regexp.MustCompile("Timeout while waiting for domain validation"),
				},
			},
		},
		"create fails because of the API error": {
			init: func(m *mockProperty, domMock *domainownership.Mock) {
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				// read
				m.mockGetPropertyVersionHostnames()

				validateReq := domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from1.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				// create
				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(nil, fmt.Errorf("API error occurred")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					ExpectError: regexp.MustCompile("API error occurred"),
				},
			},
		},
		"update with domain requiring polling": {
			init: func(m *mockProperty, domMock *domainownership.Mock) {

				m.latestVersion = 1
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				// read
				m.mockGetPropertyVersionHostnames().Times(3)

				// update
				m.latestVersion = 2
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				validateReq := domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from1.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				validateResp := domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from1.test.domain",
							DomainStatus:    "PENDING",
							ValidationScope: "HOST",
						},
					},
				}
				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(&validateResp, nil).Once()

				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				// read
				m.mockGetPropertyVersionHostnames()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					Check: resource.ComposeTestCheckFunc(
						commonStateChecker.Build(),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config_update.tf"),
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("akamai_property_domainownership_late_validation.test").
							CheckEqual("property_id", "prp_123").
							CheckEqual("version", "2").
							CheckEqual("group_id", "grp_1").
							CheckEqual("contract_id", "ctr_1").
							Build(),
					),
				},
			},
		},
		"update property version having same hostname": {
			init: func(m *mockProperty, _ *domainownership.Mock) {
				m.latestVersion = 1
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames().Times(3)

				// update
				m.latestVersion = 2
				m.mockGetPropertyVersionHostnames()

				// read
				m.mockGetPropertyVersionHostnames()

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("akamai_property_domainownership_late_validation.test").
							CheckEqual("version", "1").
							Build(),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config_update.tf"),
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("akamai_property_domainownership_late_validation.test").
							CheckEqual("version", "2").
							Build(),
					),
				},
			},
		},
		"update property version having different hostname ": {
			init: func(m *mockProperty, domMock *domainownership.Mock) {
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				// read at the beginning of Create
				m.mockGetPropertyVersionHostnames()

				validateReq := domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from1.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				validateResp := domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from1.test.domain",
							DomainStatus:    "PENDING",
							ValidationScope: "HOST",
						},
					},
				}
				// create
				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(&validateResp, nil).Once()

				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				// Read
				m.mockGetPropertyVersionHostnames().Twice()

				// update
				m.latestVersion = 2
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from2.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				validateReq = domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from2.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				validateResp = domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from2.test.domain",
							DomainStatus:    "PENDING",
							ValidationScope: "HOST",
						},
					},
				}

				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(&validateResp, nil).Once()

				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				// read
				m.mockGetPropertyVersionHostnames()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					Check: resource.ComposeTestCheckFunc(
						commonStateChecker.Build(),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config_update.tf"),
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("akamai_property_domainownership_late_validation.test").
							CheckEqual("property_id", "prp_123").
							CheckEqual("contract_id", "ctr_1").
							CheckEqual("version", "2").
							CheckEqual("group_id", "grp_1").Build(),
					),
				},
			},
		},
		"update group id": {
			init: func(m *mockProperty, domMock *domainownership.Mock) {
				// 1. INITIAL STATE: One domain is PENDING
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				validateReq := domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from1.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				validateResp := domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from1.test.domain",
							DomainStatus:    "PENDING",
							ValidationScope: "HOST",
						},
					},
				}
				// create
				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(&validateResp, nil).Once()

				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				// Read
				m.mockGetPropertyVersionHostnames().Twice()

				// update
				m.groupID = "grp_2"
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from2.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "PENDING"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				validateReq = domainownership.ValidateDomainsRequest{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:       "from2.test.domain",
							ValidationMethod: "DNS_CNAME",
							ValidationScope:  "HOST",
						},
					},
				}
				validateResp = domainownership.ValidateDomainsResponse{
					Domains: []domainownership.ValidateDomainResponse{
						{
							DomainName:      "from2.test.domain",
							DomainStatus:    "PENDING",
							ValidationScope: "HOST",
						},
					},
				}
				domMock.On("ValidateDomains", testutils.MockContext, validateReq).Return(&validateResp, nil).Once()

				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from2.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()

				m.mockGetPropertyVersionHostnames()
			},
			steps: []resource.TestStep{
				{
					// Step 1: CREATE the resource with group_id "grp_1"
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					Check: resource.ComposeTestCheckFunc(
						commonStateChecker.Build(),
					),
				},
				{
					// Step 2: UPDATE the resource with property_id grp_2"
					Config: testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config_update_group_id.tf"),
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("akamai_property_domainownership_late_validation.test").
							CheckEqual("property_id", "prp_123").
							CheckEqual("contract_id", "ctr_1").
							CheckEqual("version", "1").
							CheckEqual("group_id", "grp_2").Build(),
					),
				},
			},
		},
		"error expected - missing contract_id": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_late_validation" "test" {
  group_id          = "grp_2"
  property_id       = "prp_123"
  version           = 1
  validation_method = "DNS_CNAME"
}
`,
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found.`),
				},
			},
		},
		"error expected - missing group_id": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_late_validation" "test" {
  contract_id       = "ctr_1"
  property_id       = "prp_123"
  version           = 1
  validation_method = "DNS_CNAME"
}
`,
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found.`),
				},
			},
		},
		"error expected - missing property_id": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_late_validation" "test" {
  contract_id       = "ctr_1"
  group_id          = "grp_2"
  version           = 1
  validation_method = "DNS_CNAME"
}
`,
					ExpectError: regexp.MustCompile(`The argument "property_id" is required, but no definition was found.`),
				},
			},
		},
		"error expected - missing version": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_late_validation" "test" {
  contract_id       = "ctr_1"
  group_id          = "grp_2"
  property_id       = "prp_123"
  validation_method = "DNS_CNAME"
}
`,
					ExpectError: regexp.MustCompile(`The argument "version" is required, but no definition was found.`),
				},
			},
		},
		"error expected - missing validation_method": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_late_validation" "test" {
  contract_id       = "ctr_1"
  group_id          = "grp_2"
  property_id       = "prp_123"
  version           = 1
}
`,
					ExpectError: regexp.MustCompile(`The argument "validation_method" is required, but no definition was found.`),
				},
			},
		},
		"error expected - incorrect validation_method": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_late_validation" "test" {
  contract_id       = "ctr_1"
  group_id          = "grp_2"
  property_id       = "prp_123"
  version           = 1
  validation_method = "incorrect"
}
`,
					ExpectError: regexp.MustCompile(`(?s)Attribute validation_method value must be one of: \["DNS_CNAME" "DNS_TXT".+"HTTP"], got: "incorrect"`),
				},
			},
		},
		"import successful": {
			init: func(m *mockProperty, _ *domainownership.Mock) {
				m.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:                   "from1.test.domain",
							DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
						},
					},
				}
				m.mockGetPropertyVersionHostnames()
			},
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					ImportStateId:    "prp_123,1,ctr_1,grp_1,DNS_CNAME",
					ImportStateCheck: importChecker.Build(),
					ImportState:      true,
					ResourceName:     "akamai_property_domainownership_late_validation.test",
				},
			},
		},
		"import without prefixes in the ID": {
			init: func(m *mockProperty, _ *domainownership.Mock) {
				req := papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "123",
					GroupID:           "1",
					ContractID:        "1",
					PropertyVersion:   1,
					IncludeCertStatus: true,
				}

				resp := papi.GetPropertyVersionHostnamesResponse{
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
					PropertyID:      "prp_123",
					PropertyVersion: 1,
					Hostnames: papi.HostnameResponseItems{
						Items: []papi.Hostname{
							{
								CnameFrom:                   "from1.test.domain",
								DomainOwnershipVerification: &papi.DomainOwnershipVerification{Status: "VALIDATED"},
							},
						},
					},
				}

				m.papiMock.On("GetPropertyVersionHostnames", testutils.MockContext, req).Return(&resp, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/import.tf"),
					ImportStateId:    "123,1,1,1,DNS_CNAME",
					ImportStateCheck: importCheckerWithoutPrefix.Build(),
					ImportState:      true,
					ResourceName:     "akamai_property_domainownership_late_validation.test",
				},
			},
		},
		"import fails with wrong number of parts": {
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/import.tf"),
					ImportStateId:    "prp_123,ctr_1,grp_1",
					ImportStateCheck: importChecker.Build(),
					ImportState:      true,
					ResourceName:     "akamai_property_domainownership_late_validation.test",
					ExpectError:      regexp.MustCompile("Error: Unexpected Import Identifier"),
				},
			},
		},
		"import fails with invalid property id": {
			init: func(m *mockProperty, _ *domainownership.Mock) {
				m.papiMock.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        m.propertyID,
					PropertyVersion:   m.latestVersion,
					ContractID:        m.contractID,
					GroupID:           m.groupID,
					IncludeCertStatus: true,
				}).Return(nil, fmt.Errorf("resource not found")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/import.tf"),
					ImportStateId:    "prp_123,1,ctr_1,grp_1,DNS_CNAME",
					ImportStateCheck: importChecker.Build(),
					ImportState:      true,
					ResourceName:     "akamai_property_domainownership_late_validation.test",
					ExpectError:      regexp.MustCompile("resource not found"),
				},
			},
		},
		"import fails because of the API error": {
			init: func(m *mockProperty, _ *domainownership.Mock) {
				m.papiMock.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        m.propertyID,
					PropertyVersion:   m.latestVersion,
					ContractID:        m.contractID,
					GroupID:           m.groupID,
					IncludeCertStatus: true,
				}).Return(nil, fmt.Errorf("API error occurred")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDOMLateValidation/config.tf"),
					ExpectError: regexp.MustCompile("API error occurred"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			mp := mockProperty{
				mockPropertyData: mockPropertyData{
					propertyID:    key.PropertyID.ValueString(),
					groupID:       key.GroupID.ValueString(),
					contractID:    key.ContractID.ValueString(),
					latestVersion: int(key.Version.ValueInt64()),
				},
				papiMock: client.PAPI,
			}
			if tc.init != nil {
				tc.init(&mp, client.DomainOwnership)
			}

			config := defaultSubproviderConfig()
			config.lateValidation.searchInterval = 1 * time.Millisecond

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, newSubproviderWithConfig(config)),
				Steps:                    tc.steps,
			})

			client.DomainOwnership.AssertExpectations(t)
			client.PAPI.AssertExpectations(t)
		})
	}
}
