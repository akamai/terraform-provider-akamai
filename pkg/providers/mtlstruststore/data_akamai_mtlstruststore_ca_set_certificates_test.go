package mtlstruststore

import (
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetCertificatesDataSource(t *testing.T) {
	getGetCASetVersionCertificatesRequest := func(certificateStatus *mtlstruststore.CertificateStatus) mtlstruststore.GetCASetVersionCertificatesRequest {
		return mtlstruststore.GetCASetVersionCertificatesRequest{
			CASetID:           "12345",
			Version:           1,
			CertificateStatus: certificateStatus,
		}
	}

	getGetCASetVersionCertificatesResponse := func(certificates []mtlstruststore.CertificateResponse) *mtlstruststore.GetCASetVersionCertificatesResponse {
		return &mtlstruststore.GetCASetVersionCertificatesResponse{
			CASetID:      "12345",
			CASetName:    "example-ca-set",
			Version:      1,
			Certificates: certificates,
		}
	}

	var (
		expiredCertificate = mtlstruststore.CertificateResponse{
			Subject:            "CN=example.com, O=Example Org, C=US",
			Issuer:             "CN=Example CA, O=Example Org, C=US",
			StartDate:          time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			EndDate:            time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Fingerprint:        "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD",
			CertificatePEM:     "-----BEGIN CERTIFICATE-----\nMIID...FAKE...CERT==\n-----END CERTIFICATE-----",
			SerialNumber:       "1234567890ABCDEF",
			SignatureAlgorithm: "SHA256WITHRSA",
			CreatedDate:        time.Date(2023, 5, 20, 12, 34, 56, 0, time.UTC),
			CreatedBy:          "admin@example.com",
			Description:        ptr.To("Test certificate for example.com"),
		}

		activeCertificate = mtlstruststore.CertificateResponse{
			Subject:            "CN=api.example.org, O=Example Org, C=US",
			Issuer:             "CN=Example CA, O=Example Org, C=US",
			StartDate:          time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
			EndDate:            time.Date(3026, 1, 10, 0, 0, 0, 0, time.UTC),
			Fingerprint:        "11:22:33:44:55:66:77:88:99:00:AA:BB:CC:DD:EE:FF:11:22:33:44",
			CertificatePEM:     "-----BEGIN CERTIFICATE-----\nMIID...FAKE...API==\n-----END CERTIFICATE-----",
			SerialNumber:       "ABCDEF1234567890",
			SignatureAlgorithm: "SHA384WITHRSA",
			CreatedDate:        time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			CreatedBy:          "system@example.org",
			Description:        ptr.To("API certificate for internal usage"),
		}

		expiringCertificate = mtlstruststore.CertificateResponse{
			CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
			Description:        ptr.To("Example Certificate Expired"),
			CreatedBy:          "example user",
			CreatedDate:        tst.NewTimeFromStringMust("2024-11-05T12:08:34.099457Z"),
			StartDate:          tst.NewTimeFromStringMust("2024-11-05T12:08:34.099457Z"),
			EndDate:            tst.NewTimeFromStringMust("2525-07-20T12:08:34.099457Z"),
			Fingerprint:        "AB:CD:EF:12:34:56:78:90",
			Issuer:             "CN=Example CA, O=Example Org, C=US",
			SerialNumber:       "ABCDEF1234567890",
			SignatureAlgorithm: "SHA256WITHRSA",
			Subject:            "CN=example.com, O=Example Org, C=US",
		}
	)

	t.Parallel()
	expiredCertChecker := test.AttributeBatch{
		"subject":             "CN=example.com, O=Example Org, C=US",
		"issuer":              "CN=Example CA, O=Example Org, C=US",
		"fingerprint":         "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD",
		"certificate_pem":     "-----BEGIN CERTIFICATE-----\nMIID...FAKE...CERT==\n-----END CERTIFICATE-----",
		"serial_number":       "1234567890ABCDEF",
		"signature_algorithm": "SHA256WITHRSA",
		"created_by":          "admin@example.com",
		"created_date":        "2023-05-20T12:34:56Z",
		"start_date":          "2023-06-01T00:00:00Z",
		"end_date":            "2025-06-01T00:00:00Z",
		"description":         "Test certificate for example.com",
	}

	activeCertChecker := test.AttributeBatch{
		"subject":             "CN=api.example.org, O=Example Org, C=US",
		"issuer":              "CN=Example CA, O=Example Org, C=US",
		"start_date":          "2024-01-10T00:00:00Z",
		"end_date":            "3026-01-10T00:00:00Z",
		"fingerprint":         "11:22:33:44:55:66:77:88:99:00:AA:BB:CC:DD:EE:FF:11:22:33:44",
		"certificate_pem":     "-----BEGIN CERTIFICATE-----\nMIID...FAKE...API==\n-----END CERTIFICATE-----",
		"serial_number":       "ABCDEF1234567890",
		"signature_algorithm": "SHA384WITHRSA",
		"created_date":        "2024-01-01T09:00:00Z",
		"created_by":          "system@example.org",
		"description":         "API certificate for internal usage",
	}

	expiringCertChecker := test.AttributeBatch{
		"certificate_pem":     "-----BEGIN CERTIFICATE-----...",
		"description":         "Example Certificate Expired",
		"created_by":          "example user",
		"created_date":        "2024-11-05T12:08:34Z",
		"start_date":          "2024-11-05T12:08:34Z",
		"end_date":            "2525-07-20T12:08:34Z",
		"fingerprint":         "AB:CD:EF:12:34:56:78:90",
		"issuer":              "CN=Example CA, O=Example Org, C=US",
		"serial_number":       "ABCDEF1234567890",
		"signature_algorithm": "SHA256WITHRSA",
		"subject":             "CN=example.com, O=Example Org, C=US",
	}

	commonStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
		CheckEqual("id", "12345").
		CheckEqual("version", "1").
		CheckEqual("certificates.#", "3").
		CheckEqualBatch("certificates.0.", expiredCertChecker).
		CheckEqualBatch("certificates.1.", activeCertChecker).
		CheckEqualBatch("certificates.2.", expiringCertChecker)

	tests := map[string]struct {
		init     func(*mtlstruststore.Mock, caSetTestData)
		testData caSetTestData
		steps    []resource.TestStep
		error    *regexp.Regexp
	}{
		"read successful": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ActiveOrExpiredCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiredCertificate, activeCertificate, expiringCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_active = true
						  include_expired = true
						}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"read successful - read by ca set name": {
			testData: mockCASetData,
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: testData.caSetName,
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "12345",
							CASetName:   "example-ca-set",
							CASetStatus: "NOT_DELETED",
						},
					},
				}, nil)
				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ActiveOrExpiredCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiredCertificate, activeCertificate, expiringCertificate}))
			},

			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  name = "example-ca-set"
						  version = 1
						  include_active = true
						  include_expired = true
						}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"read successful - version not provided, use the latest version": {
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: testData.caSetID,
				}).Return(&testData.caSetResponse, nil).Times(3)
				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ActiveOrExpiredCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiredCertificate, activeCertificate, expiringCertificate}))
			},
			testData: mockCASetData,
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  include_active = true
						  include_expired = true
						}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"read successful - only expired certificates": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiredCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiredCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_expired = true
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "1").
							CheckEqualBatch("certificates.0.", expiredCertChecker).Build(),
					),
				},
			},
		},
		"read successful - only active certificates": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ActiveCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{activeCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_active = true
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "1").
							CheckEqualBatch("certificates.0.", activeCertChecker).Build(),
					),
				},
			},
		},
		"read successful - with expired and include_expiring_in_days": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiringRequest := getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiringCert))
				expiringRequest.ExpiryThresholdInDays = ptr.To(20)
				mockGetCASetVersionCertificates(m, expiringRequest, getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiringCertificate}))

				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiredCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiredCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
		                  include_expired = true
						  include_expiring_in_days = 20
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "2").
							CheckEqualBatch("certificates.0.", expiringCertChecker).
							CheckEqualBatch("certificates.1.", expiredCertChecker).Build(),
					),
				},
			},
		},
		"read successful - with expired and include_expiring_by_date": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiringRequest := getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiringCert))
				expiringRequest.ExpiryThresholdTimestamp = tst.NewTimeFromString(t, "2625-07-01T12:00:00Z")
				mockGetCASetVersionCertificates(m, expiringRequest, getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiringCertificate}))

				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiredCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiredCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
		                  include_expired = true
						  include_expiring_by_date = "2625-07-01T12:00:00Z"
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "2").
							CheckEqualBatch("certificates.0.", expiringCertChecker).
							CheckEqualBatch("certificates.1.", expiredCertChecker).Build(),
					),
				},
			},
		},
		"read successful - with only include_expiring_in_days": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiringRequest := getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiringCert))
				expiringRequest.ExpiryThresholdInDays = ptr.To(30)
				mockGetCASetVersionCertificates(m, expiringRequest, getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiringCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_expiring_in_days = 30
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "1").
							CheckEqualBatch("certificates.0.", expiringCertChecker).Build(),
					),
				},
			},
		},
		"read successful - with only include_expiring_by_date": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiringRequest := getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiringCert))
				expiringRequest.ExpiryThresholdTimestamp = tst.NewTimeFromString(t, "2625-07-01T12:00:00Z")
				mockGetCASetVersionCertificates(m, expiringRequest, getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiringCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_expiring_by_date = "2625-07-01T12:00:00Z"
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "1").
							CheckEqualBatch("certificates.0.", expiringCertChecker).Build(),
					),
				},
			},
		},
		"read successful - with include_expiring_by_date as RFC3339Nano": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiringRequest := getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ExpiringCert))
				expiringRequest.ExpiryThresholdTimestamp = tst.NewTimeFromString(t, "2625-07-01T12:00:06.825235Z")
				mockGetCASetVersionCertificates(m, expiringRequest, getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiringCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_expiring_by_date = "2625-07-01T12:00:06.825235Z"
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "1").
							CheckEqualBatch("certificates.0.", expiringCertChecker).Build(),
					),
				},
			},
		},
		"read successful - with include_active and include_expired provided externally": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				mockGetCASetVersionCertificates(m, getGetCASetVersionCertificatesRequest(ptr.To(mtlstruststore.ActiveOrExpiredCert)), getGetCASetVersionCertificatesResponse([]mtlstruststore.CertificateResponse{expiredCertificate, activeCertificate, expiringCertificate}))
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						variable "isActive" {
							type    = bool
							default = true
						}

						variable "isExpired" {
							type    = bool
							default = true
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_active = var.isActive
						  include_expired = var.isExpired
						}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"read failed - ca set not found for given name": {
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				testData.caSets[0].CASetStatus = ""
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: testData.caSetName,
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: testData.caSets,
				}, nil).Once()
			},
			testData: mockCASetData,
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  name = "example-ca-set"
						  version = 1
		                  include_expired = true
						}`,
					ExpectError: regexp.MustCompile("failed to find CA set by name 'example-ca-set': no CA set found with name"),
				},
			},
		},
		"read failed - version not provided, the latest version is not available": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				mockCASetData.caSetResponse.LatestVersion = nil
				m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: mockCASetData.caSetID,
				}).Return(&mockCASetData.caSetResponse, nil).Times(1)
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  include_active = true
						  include_expired = true
						}`,
					ExpectError: regexp.MustCompile("no version provided and CA set has no latest version available"),
				},
			},
		},
		"validation error - one of `id` or `name` required": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  version = 1
						}`,
					ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[(id,name|name,id)] is required`),
				},
			},
		},
		"validation error - empty `name`": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  version = 1
                          name = ""
						}`,
					ExpectError: regexp.MustCompile(`Attribute name must not be empty or only whitespace`),
				},
			},
		},
		"validation error - empty `id`": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  version = 1
                          id = ""
						}`,
					ExpectError: regexp.MustCompile(`Attribute id string length must be at least 1, got: 0`),
				},
			},
		},
		"validation error - invalid value for `include_expiring_in_days`": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
		                  include_expired = true
						  include_expiring_in_days = -20
						}`,
					ExpectError: regexp.MustCompile(`Attribute include_expiring_in_days value must be at least 1, got: -20`),
				},
			},
		},
		"validation error - both `include_expiring_in_days` and `include_expiring_by_date` provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
		                  include_expired = true
						  include_expiring_in_days = 20
						  include_expiring_by_date = "2625-01-01T00:00:00Z"
						}`,
					ExpectError: regexp.MustCompile(`Attribute "include_expiring_by_date" cannot be specified when\n.*"include_expiring_in_days" is specified`),
				},
			},
		},
		"validation error - `include_expiring_by_date` provided with incorrect value": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
		                  include_expired = true
						  include_expiring_by_date = "yesterday"
						}`,
					ExpectError: regexp.MustCompile(`The provided expiring timestamp 'yesterday' is not a valid RFC3339 or\n.*RFC3339Nano formatted date`),
				},
			},
		},
		"validation error - scope of listing not provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						}`,
					ExpectError: regexp.MustCompile(`At least one attribute out of 'include_active', 'include_expired',\n.*'include_expiring_in_days', or 'include_expiring_by_date' must be specified\n.*with 'true' value for booleans, or some value for the rest`),
				},
			},
		},
		"validation error - scope of listing failed - only include_active = false": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_active = false
						}`,
					ExpectError: regexp.MustCompile(`At least one attribute out of 'include_active', 'include_expired',\n.*'include_expiring_in_days', or 'include_expiring_by_date' must be specified\n.*with 'true' value for booleans, or some value for the rest`),
				},
			},
		},
		"validation error - scope of listing failed - only include_expired = false": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_expired = false
						}`,
					ExpectError: regexp.MustCompile(`At least one attribute out of 'include_active', 'include_expired',\n.*'include_expiring_in_days', or 'include_expiring_by_date' must be specified\n.*with 'true' value for booleans, or some value for the rest`),
				},
			},
		},
		"validation error - `include_active` and `include_expiring_in_days` provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_active = true
						  include_expiring_in_days = 20
						}`,
					ExpectError: regexp.MustCompile(`Attribute "include_expiring_in_days" cannot be specified when\n.*"include_active" is specified`),
				},
			},
		},
		"validation error - `include_active` and `include_expiring_by_date` provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_active = true
						  include_expiring_by_date = "2625-07-01T12:00:00Z"
						}`,
					ExpectError: regexp.MustCompile(`Attribute "include_expiring_by_date" cannot be specified when\n.*"include_active" is specified`),
				},
			},
		},
		"validation error - `include_expiring_by_date` with old date": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}

						data "akamai_mtlstruststore_ca_set_certificates" "test" {
						  id = 12345
						  version = 1
						  include_expiring_by_date = "2024-07-01T12:00:00Z"
						}`,
					ExpectError: regexp.MustCompile(`The provided expiring threshold timestamp '2024-07-01T12:00:00Z' cannot be in\n.*the past`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.testData)
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

var mockCASetData = caSetTestData{
	caSetName:    "example-ca-set",
	caSetID:      "12345",
	caSetVersion: 1,
	caSetResponse: mtlstruststore.GetCASetResponse{
		CASetID:           "12345",
		CASetName:         "example-ca-set",
		Description:       ptr.To("Example CA Set"),
		AccountID:         "account-123",
		CreatedBy:         "example user",
		CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
		StagingVersion:    ptr.To(int64(1)),
		ProductionVersion: ptr.To(int64(1)),
		LatestVersion:     ptr.To(int64(1)),
	},
	caSetVersionResponse: mtlstruststore.GetCASetVersionResponse{
		Description:       ptr.To("Version 1 description"),
		AllowInsecureSHA1: false,
		CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
		CreatedBy:         "example user",
		ModifiedBy:        ptr.To("example user"),
		ModifiedDate:      ptr.To(tst.NewTimeFromStringMust("2025-05-16T12:08:34.099457Z")),
		Certificates: []mtlstruststore.CertificateResponse{
			{
				CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
				Description:        ptr.To("Example Certificate"),
				CreatedBy:          "example user",
				CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
				StartDate:          tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
				EndDate:            tst.NewTimeFromStringMust("2026-04-16T12:08:34.099457Z"),
				Fingerprint:        "AB:CD:EF:12:34:56:78:90",
				Issuer:             "Example Issuer",
				SerialNumber:       "123456789",
				SignatureAlgorithm: "SHA256",
				Subject:            "Example Subject",
			},
		},
	},
	caSets: []mtlstruststore.CASetResponse{
		{
			CASetID:     "12345",
			CASetName:   "example-ca-set",
			CASetStatus: "NOT_DELETED",
		},
	},
}

func mockGetCASetVersionCertificates(m *mtlstruststore.Mock, req mtlstruststore.GetCASetVersionCertificatesRequest, resp *mtlstruststore.GetCASetVersionCertificatesResponse) {
	m.On("GetCASetVersionCertificates", testutils.MockContext, req).Return(resp, nil).Times(3)
}
