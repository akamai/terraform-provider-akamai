package property

import (
	"fmt"
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

func TestDomainOwnershipDataSource(t *testing.T) {
	testDir := "testdata/TestDataDomainOwnershipDomain/"
	t.Parallel()

	commonStateChecker := test.NewStateChecker("data.akamai_property_domainownership_domain.testdomain").
		CheckEqual("domain_name", "example.com").
		CheckEqual("validation_scope", "DOMAIN").
		CheckEqual("account_id", "test_account").
		CheckEqual("domain_status", "VALIDATED").
		CheckEqual("validation_method", "DNS_CNAME").
		CheckEqual("validation_requested_by", "user1").
		CheckEqual("validation_requested_date", "2023-01-03T00:00:00Z").
		CheckEqual("validation_completed_date", "2023-01-03T00:00:00Z").
		CheckEqual("validation_challenge.cname_record.name", "cname-name-1").
		CheckEqual("validation_challenge.cname_record.target", "cname-target-1").
		CheckEqual("validation_challenge.txt_record.name", "txt-name-1").
		CheckEqual("validation_challenge.txt_record.value", "txt-value-1").
		CheckEqual("validation_challenge.http_file.path", "http-file-path-1").
		CheckEqual("validation_challenge.http_file.content", "http-file-content-1").
		CheckEqual("validation_challenge.http_file.content_type", "text/plain").
		CheckEqual("validation_challenge.http_redirect.from", "http-redirect-from-1").
		CheckEqual("validation_challenge.http_redirect.to", "http-redirect-to-1").
		CheckEqual("validation_challenge.expiration_date", "2025-08-05T13:27:19Z").
		CheckEqual("domain_status_history.#", "1").
		CheckEqual("domain_status_history.0.domain_status", "VALIDATED").
		CheckEqual("domain_status_history.0.modified_date", "2023-01-01T00:00:00Z").
		CheckEqual("domain_status_history.0.modified_user", "user1").
		CheckEqual("domain_status_history.0.message", "Domain validated successfully")

	tests := map[string]struct {
		init  func(*domainownership.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - get domain": {
			init: func(m *domainownership.Mock) {
				req := domainownership.GetDomainRequest{
					DomainName:                 "example.com",
					ValidationScope:            domainownership.ValidationScope("DOMAIN"),
					IncludeDomainStatusHistory: true,
				}

				resp := &domainownership.GetDomainResponse{
					DomainName:              "example.com",
					ValidationScope:         "DOMAIN",
					AccountID:               "test_account",
					DomainStatus:            "VALIDATED",
					ValidationMethod:        ptr.To("DNS_CNAME"),
					ValidationRequestedBy:   "user1",
					ValidationRequestedDate: tst.NewTimeFromStringMust("2023-01-03T00:00:00Z"),
					ValidationCompletedDate: ptr.To(tst.NewTimeFromStringMust("2023-01-03T00:00:00Z")),
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
					DomainStatusHistory: []domainownership.DomainStatusHistory{
						{
							DomainStatus: "VALIDATED",
							ModifiedDate: tst.NewTimeFromStringMust("2023-01-01T00:00:00Z"),
							ModifiedUser: "user1",
							Message:      ptr.To("Domain validated successfully"),
						},
					},
				}

				m.On("GetDomain", mock.Anything, req).Return(resp, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"domainownership.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *domainownership.Mock) {
				req := domainownership.GetDomainRequest{
					DomainName:                 "example.com",
					ValidationScope:            domainownership.ValidationScope("DOMAIN"),
					IncludeDomainStatusHistory: true,
				}
				m.On("GetDomain", mock.Anything, req).Return((*domainownership.GetDomainResponse)(nil), fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"domainownership.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"validation error - domain_name missing": {
			init: func(_ *domainownership.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"domain_name_missing.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument(\n|.)+` + `The argument "domain_name" is required, but no definition was found.`),
				},
			},
		},
		"validation error - validation_scope missing": {
			init: func(_ *domainownership.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"validation_scope_missing.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument(\n|.)+` + `The argument "validation_scope" is required, but no definition was found.`),
				},
			},
		},
		"validation error - invalid `validation_scope`": {
			init: func(_ *domainownership.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"invalid_validation_scope.tf"),
					ExpectError: regexp.MustCompile(`Attribute validation_scope value must be one of: \["HOST" "WILDCARD"\n"DOMAIN"\], got: "DNS"`),
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
