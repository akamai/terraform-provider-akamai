package cloudcertificates

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func defaultCertResp() *cloudcertificates.GetCertificateResponse {
	return &cloudcertificates.GetCertificateResponse{
		Certificate: cloudcertificates.Certificate{
			AccountID:                           "test_account",
			CertificateName:                     "example-cert",
			CertificateStatus:                   "ACTIVE",
			CertificateType:                     "THIRD_PARTY",
			ContractID:                          "test_contract",
			CreatedBy:                           "user1",
			CreatedDate:                         tst.NewTimeFromStringMust("2023-11-01T02:40:20Z"),
			CSRExpirationDate:                   tst.NewTimeFromStringMust("2025-01-03T00:00:00Z"),
			CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\ntest-csr\n-----END CERTIFICATE REQUEST-----",
			KeySize:                             "2048",
			KeyType:                             "RSA",
			ModifiedBy:                          "user2",
			ModifiedDate:                        tst.NewTimeFromStringMust("2024-06-02T05:06:08Z"),
			SANs:                                []string{"example.com"},
			SecureNetwork:                       "STANDARD_TLS",
			SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----"),
			SignedCertificateIssuer:             ptr.To("Test CA"),
			SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2023-01-02T00:00:00Z")),
			SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2025-01-02T00:00:00Z")),
			SignedCertificateSerialNumber:       ptr.To("123456789"),
			SignedCertificateSHA256Fingerprint:  ptr.To("aa:bb:cc:dd:ee:ff"),
			TrustChainPEM:                       ptr.To("-----BEGIN CERTIFICATE-----\ntrust-chain\n-----END CERTIFICATE-----"),
			Subject: &cloudcertificates.Subject{
				CommonName:   "example.com",
				Organization: "Example Org",
				Country:      "US",
				State:        "California",
				Locality:     "San Francisco",
			},
		},
	}
}
func TestCertificateDataSource(t *testing.T) {
	testDir := "testdata/TestDataCertificate/"
	t.Parallel()

	commonStateChecker := test.NewStateChecker("data.akamai_cloudcertificates_certificate.testcert").
		CheckEqual("certificate_id", "12345").
		CheckEqual("account_id", "test_account").
		CheckEqual("certificate_name", "example-cert").
		CheckEqual("certificate_status", "ACTIVE").
		CheckEqual("certificate_type", "THIRD_PARTY").
		CheckEqual("contract_id", "test_contract").
		CheckEqual("created_by", "user1").
		CheckEqual("created_date", "2023-11-01T02:40:20Z").
		CheckEqual("csr_expiration_date", "2025-01-03T00:00:00Z").
		CheckEqual("csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\ntest-csr\n-----END CERTIFICATE REQUEST-----").
		CheckEqual("key_size", "2048").
		CheckEqual("key_type", "RSA").
		CheckEqual("modified_by", "user2").
		CheckEqual("modified_date", "2024-06-02T05:06:08Z").
		CheckEqual("sans.#", "2").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("signed_certificate_pem", "-----BEGIN CERTIFICATE-----\ntest-cert\n-----END CERTIFICATE-----").
		CheckEqual("signed_certificate_issuer", "Test CA").
		CheckEqual("signed_certificate_not_valid_before_date", "2023-01-02T00:00:00Z").
		CheckEqual("signed_certificate_not_valid_after_date", "2025-01-02T00:00:00Z").
		CheckEqual("signed_certificate_serial_number", "123456789").
		CheckEqual("signed_certificate_sha256_fingerprint", "aa:bb:cc:dd:ee:ff").
		CheckEqual("trust_chain_pem", "-----BEGIN CERTIFICATE-----\ntrust-chain\n-----END CERTIFICATE-----").
		CheckEqual("subject.common_name", "example.com").
		CheckEqual("subject.organization", "Example Org").
		CheckEqual("subject.country", "US").
		CheckEqual("subject.state", "California").
		CheckEqual("subject.locality", "San Francisco").
		CheckEqual("bindings.#", "2").
		CheckEqual("bindings.0.certificate_id", "12345").
		CheckEqual("bindings.0.hostname", "www.example.com").
		CheckEqual("bindings.0.network", "PRODUCTION").
		CheckEqual("bindings.0.resource_type", "CDN_HOSTNAME").
		CheckEqual("bindings.1.certificate_id", "12345").
		CheckEqual("bindings.1.hostname", "api.example.com").
		CheckEqual("bindings.1.network", "STAGING").
		CheckEqual("bindings.1.resource_type", "CDN_HOSTNAME")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - get certificate without bindings": {
			init: func(m *cloudcertificates.Mock) {
				certReq := cloudcertificates.GetCertificateRequest{
					CertificateID: "12345",
				}

				certResp := defaultCertResp()
				certResp.Certificate.SANs = []string{"example.com"}

				m.On("GetCertificate", mock.Anything, certReq).Return(certResp, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"certificate.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_certificate.testcert").
						CheckEqual("certificate_id", "12345").
						CheckMissing("include_hostname_bindings").
						CheckEqual("account_id", "test_account").
						CheckEqual("bindings.#", "0").
						Build(),
				},
			},
		},
		"happy path - get certificate with bindings": {
			init: func(m *cloudcertificates.Mock) {
				certReq := cloudcertificates.GetCertificateRequest{
					CertificateID: "12345",
				}

				certResp := defaultCertResp()
				certResp.Certificate.SANs = []string{"example.com", "www.example.com"}

				bindingsReq := cloudcertificates.ListCertificateBindingsRequest{
					CertificateID: "12345",
					PageSize:      100,
					Page:          1,
				}

				bindingsResp := &cloudcertificates.ListCertificateBindingsResponse{
					Bindings: []cloudcertificates.CertificateBinding{
						{
							CertificateID: "12345",
							Hostname:      "www.example.com",
							Network:       "PRODUCTION",
							ResourceType:  "CDN_HOSTNAME",
						},
						{
							CertificateID: "12345",
							Hostname:      "api.example.com",
							Network:       "STAGING",
							ResourceType:  "CDN_HOSTNAME",
						},
					},
					Links: cloudcertificates.Links{
						Next: nil,
					},
				}

				m.On("GetCertificate", mock.Anything, certReq).Return(certResp, nil).Times(3)
				m.On("ListCertificateBindings", mock.Anything, bindingsReq).Return(bindingsResp, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"certificate_with_bindings.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - get certificate with pagination": {
			init: func(m *cloudcertificates.Mock) {
				certReq := cloudcertificates.GetCertificateRequest{
					CertificateID: "12345",
				}

				certResp := defaultCertResp()

				bindingsReqPage1 := cloudcertificates.ListCertificateBindingsRequest{
					CertificateID: "12345",
					PageSize:      100, // "Even though we request 100 items in each API call, the mock returns only 1 item here for testing pagination.
					Page:          1,
				}

				bindingsRespPage1 := &cloudcertificates.ListCertificateBindingsResponse{
					Bindings: []cloudcertificates.CertificateBinding{
						{
							CertificateID: "12345",
							Hostname:      "www1.example.com",
							Network:       "PRODUCTION",
							ResourceType:  "CDN_HOSTNAME",
						},
					},
					Links: cloudcertificates.Links{
						Next: ptr.To("next-page-url"),
					},
				}

				bindingsReqPage2 := cloudcertificates.ListCertificateBindingsRequest{
					CertificateID: "12345",
					PageSize:      100, // "Even though we request 100 items in each API call, the mock returns only 1 item here for testing pagination.
					Page:          2,
				}

				bindingsRespPage2 := &cloudcertificates.ListCertificateBindingsResponse{
					Bindings: []cloudcertificates.CertificateBinding{
						{
							CertificateID: "12345",
							Hostname:      "www2.example.com",
							Network:       "STAGING",
							ResourceType:  "CDN_HOSTNAME",
						},
					},
					Links: cloudcertificates.Links{
						Next: nil,
					},
				}

				m.On("GetCertificate", mock.Anything, certReq).Return(certResp, nil).Times(3)
				m.On("ListCertificateBindings", mock.Anything, bindingsReqPage1).Return(bindingsRespPage1, nil).Times(3)
				m.On("ListCertificateBindings", mock.Anything, bindingsReqPage2).Return(bindingsRespPage2, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"certificate_with_bindings.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_certificate.testcert").
						CheckEqual("certificate_id", "12345").
						CheckEqual("bindings.#", "2").
						CheckEqual("bindings.0.hostname", "www1.example.com").
						CheckEqual("bindings.1.hostname", "www2.example.com").
						Build(),
				},
			},
		},
		"happy path - get certificate with null subject": {
			init: func(m *cloudcertificates.Mock) {
				certReq := cloudcertificates.GetCertificateRequest{
					CertificateID: "12345",
				}

				certResp := defaultCertResp()
				certResp.Certificate.Subject = nil

				m.On("GetCertificate", mock.Anything, certReq).Return(certResp, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"certificate.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_certificate.testcert").
						CheckEqual("certificate_id", "12345").
						CheckEqual("account_id", "test_account").
						Build(),
				},
			},
		},
		"error - certificate not found": {
			init: func(m *cloudcertificates.Mock) {
				certReq := cloudcertificates.GetCertificateRequest{
					CertificateID: "12345",
				}
				m.On("GetCertificate", mock.Anything, certReq).Return((*cloudcertificates.GetCertificateResponse)(nil), cloudcertificates.ErrCertificateNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"certificate.tf"),
					ExpectError: regexp.MustCompile("Certificate Not Found"),
				},
			},
		},
		"error - API error on get certificate": {
			init: func(m *cloudcertificates.Mock) {
				certReq := cloudcertificates.GetCertificateRequest{
					CertificateID: "12345",
				}
				m.On("GetCertificate", mock.Anything, certReq).Return((*cloudcertificates.GetCertificateResponse)(nil), fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"certificate.tf"),
					ExpectError: regexp.MustCompile("Failed to retrieve certificate"),
				},
			},
		},
		"error - API error on list bindings": {
			init: func(m *cloudcertificates.Mock) {
				certReq := cloudcertificates.GetCertificateRequest{
					CertificateID: "12345",
				}

				certResp := &cloudcertificates.GetCertificateResponse{
					Certificate: cloudcertificates.Certificate{
						AccountID:         "test_account",
						CertificateName:   "example-cert",
						CertificateStatus: "ACTIVE",
						CertificateType:   "THIRD_PARTY",
						ContractID:        "test_contract",
						CreatedBy:         "user1",
						CreatedDate:       tst.NewTimeFromStringMust("2023-01-01T00:00:00Z"),
						CSRExpirationDate: tst.NewTimeFromStringMust("2025-01-01T00:00:00Z"),
						KeySize:           "2048",
						KeyType:           "RSA",
						ModifiedBy:        "user2",
						ModifiedDate:      tst.NewTimeFromStringMust("2023-06-01T00:00:00Z"),
						SANs:              []string{"example.com"},
						SecureNetwork:     "STANDARD_TLS",
						Subject: &cloudcertificates.Subject{
							CommonName: "example.com",
						},
					},
				}

				bindingsReq := cloudcertificates.ListCertificateBindingsRequest{
					CertificateID: "12345",
					PageSize:      100,
					Page:          1,
				}

				m.On("GetCertificate", mock.Anything, certReq).Return(certResp, nil).Once()
				m.On("ListCertificateBindings", mock.Anything, bindingsReq).Return((*cloudcertificates.ListCertificateBindingsResponse)(nil), fmt.Errorf("bindings error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"certificate_with_bindings.tf"),
					ExpectError: regexp.MustCompile("Failed to retrieve bindings"),
				},
			},
		},
		"validation error - certificate_id missing": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"certificate_id_is_missing.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument(\n|.)+` + `The argument "certificate_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - certificate_id empty": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"certificate_id_is_empty.tf"),
					ExpectError: regexp.MustCompile(`Attribute certificate_id string length must be at least 1`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
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
