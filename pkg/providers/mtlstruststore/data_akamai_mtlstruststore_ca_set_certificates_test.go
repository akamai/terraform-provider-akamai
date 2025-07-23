package mtlstruststore

import (
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	tst "github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetCertificatesDataSource(t *testing.T) {
	t.Parallel()
	commonStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
		CheckEqual("id", "12345").
		CheckEqual("version", "1").
		CheckEqual("certificates.#", "2").
		CheckEqual("certificates.0.subject", "CN=example.com, O=Example Org, C=US").
		CheckEqual("certificates.0.issuer", "CN=Example CA, O=Example Org, C=US").
		CheckEqual("certificates.0.fingerprint", "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD").
		CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIID...FAKE...CERT==\n-----END CERTIFICATE-----").
		CheckEqual("certificates.0.serial_number", "1234567890ABCDEF").
		CheckEqual("certificates.0.signature_algorithm", "SHA256WITHRSA").
		CheckEqual("certificates.0.created_by", "admin@example.com").
		CheckEqual("certificates.0.created_date", "2023-05-20T12:34:56Z").
		CheckEqual("certificates.0.start_date", "2023-06-01T00:00:00Z").
		CheckEqual("certificates.0.end_date", "2025-06-01T00:00:00Z").
		CheckEqual("certificates.0.description", "Test certificate for example.com").
		CheckEqual("certificates.1.subject", "CN=api.example.org, O=Example Org, C=US").
		CheckEqual("certificates.1.issuer", "CN=Example CA, O=Example Org, C=US").
		CheckEqual("certificates.1.serial_number", "ABCDEF1234567890")

	tests := map[string]struct {
		init     func(*mtlstruststore.Mock, caSetTestData)
		testData caSetTestData
		steps    []resource.TestStep
		error    *regexp.Regexp
	}{
		"read successful": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				mockGetCASetVersionCertificates(m)
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
				mockGetCASetVersionCertificates(m)
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
				mockGetCASetVersionCertificates(m)
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
						}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"read successful - only expired certificates": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiredCert := mtlstruststore.ExpiredCert
				getResponse := &mtlstruststore.GetCASetVersionCertificatesResponse{
					CASetID:   "12345",
					CASetName: "example-ca-set",
					Version:   1,
					Certificates: []mtlstruststore.CertificateResponse{
						{
							CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
							Description:        ptr.To("Example Certificate"),
							CreatedBy:          "example user",
							CreatedDate:        tst.NewTimeFromStringMust("2023-04-16T12:08:34Z"),
							StartDate:          tst.NewTimeFromStringMust("2023-04-16T12:08:34Z"),
							EndDate:            tst.NewTimeFromStringMust("2024-04-16T12:08:34Z"),
							Fingerprint:        "AB:CD:EF:12:34:56:78:90",
							Issuer:             "CN=Example CA, O=Example Org, C=US",
							SerialNumber:       "1234567890ABCDEF",
							SignatureAlgorithm: "SHA256WITHRSA",
							Subject:            "CN=example.com, O=Example Org, C=US",
						},
					},
				}
				m.On("GetCASetVersionCertificates", testutils.MockContext, mtlstruststore.GetCASetVersionCertificatesRequest{
					CASetID:           "12345",
					Version:           1,
					CertificateStatus: &expiredCert,
				}).Return(getResponse, nil).Times(3)
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
		                  expired = true
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "1").
							CheckEqual("certificates.0.created_date", "2023-04-16T12:08:34Z").
							CheckEqual("certificates.0.start_date", "2023-04-16T12:08:34Z").
							CheckEqual("certificates.0.end_date", "2024-04-16T12:08:34Z").Build(),
					),
				},
			},
		},
		"read successful - with expired and expiry_threshold_in_days": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiredCert := mtlstruststore.ExpiredOrExpiringCert
				getResponse := &mtlstruststore.GetCASetVersionCertificatesResponse{
					CASetID:   "12345",
					CASetName: "example-ca-set",
					Version:   1,
					Certificates: []mtlstruststore.CertificateResponse{
						{
							CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
							Description:        ptr.To("Example Certificate Expired"),
							CreatedBy:          "example user",
							CreatedDate:        tst.NewTimeFromStringMust("2024-11-05T12:08:34.099457Z"),
							StartDate:          tst.NewTimeFromStringMust("2024-11-05T12:08:34.099457Z"),
							EndDate:            tst.NewTimeFromStringMust("2025-07-20T12:08:34.099457Z"),
							Fingerprint:        "AB:CD:EF:12:34:56:78:90",
							Issuer:             "CN=Example CA, O=Example Org, C=US",
							SerialNumber:       "ABCDEF1234567890",
							SignatureAlgorithm: "SHA256WITHRSA",
							Subject:            "CN=example.com, O=Example Org, C=US",
						},
						{
							CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
							Description:        ptr.To("Example Certificate Expiring"),
							CreatedBy:          "example user",
							CreatedDate:        tst.NewTimeFromStringMust("2025-01-01T12:08:34.099457Z"),
							StartDate:          tst.NewTimeFromStringMust("2025-01-01T12:08:34.099457Z"),
							EndDate:            tst.NewTimeFromStringMust("2025-07-10T12:08:34.099457Z"),
							Fingerprint:        "AB:CD:EF:12:34:56:78:90",
							Issuer:             "CN=Example CA, O=Example Org, C=US",
							SerialNumber:       "1234567890ABCDEF",
							SignatureAlgorithm: "SHA256WITHRSA",
							Subject:            "CN=example.com, O=Example Org, C=US",
						},
					},
				}
				m.On("GetCASetVersionCertificates", testutils.MockContext, mtlstruststore.GetCASetVersionCertificatesRequest{
					CASetID:               "12345",
					Version:               1,
					CertificateStatus:     &expiredCert,
					ExpiryThresholdInDays: ptr.To(20),
				}).Return(getResponse, nil).Times(3)
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
		                  expired = true
						  expiry_threshold_in_days= "20"
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "2").
							CheckEqual("certificates.0.created_date", "2024-11-05T12:08:34Z").
							CheckEqual("certificates.0.start_date", "2024-11-05T12:08:34Z").
							CheckEqual("certificates.0.end_date", "2025-07-20T12:08:34Z").
							CheckEqual("certificates.0.created_by", "example user").
							CheckEqual("certificates.0.description", "Example Certificate Expired").
							CheckEqual("certificates.0.serial_number", "ABCDEF1234567890").
							CheckEqual("certificates.1.created_date", "2025-01-01T12:08:34Z").
							CheckEqual("certificates.1.start_date", "2025-01-01T12:08:34Z").
							CheckEqual("certificates.1.end_date", "2025-07-10T12:08:34Z").
							CheckEqual("certificates.1.created_by", "example user").
							CheckEqual("certificates.1.description", "Example Certificate Expiring").
							CheckEqual("certificates.1.serial_number", "1234567890ABCDEF").Build(),
					),
				},
			},
		},
		"read successful - with only expiry_threshold_in_days": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				expiringCert := mtlstruststore.ExpiringCert
				getResponse := &mtlstruststore.GetCASetVersionCertificatesResponse{
					CASetID:   "12345",
					CASetName: "example-ca-set",
					Version:   1,
					Certificates: []mtlstruststore.CertificateResponse{
						{
							CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
							Description:        ptr.To("Example Certificate Expired"),
							CreatedBy:          "example user",
							CreatedDate:        tst.NewTimeFromStringMust("2024-11-05T12:08:34.099457Z"),
							StartDate:          tst.NewTimeFromStringMust("2024-11-05T12:08:34.099457Z"),
							EndDate:            tst.NewTimeFromStringMust("2025-08-01T12:08:34.099457Z"),
							Fingerprint:        "AB:CD:EF:12:34:56:78:90",
							Issuer:             "CN=Example CA, O=Example Org, C=US",
							SerialNumber:       "ABCDEF1234567890",
							SignatureAlgorithm: "SHA256WITHRSA",
							Subject:            "CN=example.com, O=Example Org, C=US",
						},
						{
							CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
							Description:        ptr.To("Example Certificate Expiring"),
							CreatedBy:          "example user",
							CreatedDate:        tst.NewTimeFromStringMust("2025-01-01T12:08:34.099457Z"),
							StartDate:          tst.NewTimeFromStringMust("2025-01-01T12:08:34.099457Z"),
							EndDate:            tst.NewTimeFromStringMust("2025-07-10T12:08:34.099457Z"),
							Fingerprint:        "AB:CD:EF:12:34:56:78:90",
							Issuer:             "CN=Example CA, O=Example Org, C=US",
							SerialNumber:       "1234567890ABCDEF",
							SignatureAlgorithm: "SHA256WITHRSA",
							Subject:            "CN=example.com, O=Example Org, C=US",
						},
					},
				}
				m.On("GetCASetVersionCertificates", testutils.MockContext, mtlstruststore.GetCASetVersionCertificatesRequest{
					CASetID:               "12345",
					Version:               1,
					CertificateStatus:     &expiringCert,
					ExpiryThresholdInDays: ptr.To(30),
				}).Return(getResponse, nil).Times(3)
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
						  expiry_threshold_in_days= "30"
						}`,
					Check: resource.ComposeTestCheckFunc(
						test.NewStateChecker("data.akamai_mtlstruststore_ca_set_certificates.test").
							CheckEqual("id", "12345").
							CheckEqual("version", "1").
							CheckEqual("certificates.#", "2").
							CheckEqual("certificates.0.created_date", "2024-11-05T12:08:34Z").
							CheckEqual("certificates.0.start_date", "2024-11-05T12:08:34Z").
							CheckEqual("certificates.0.end_date", "2025-08-01T12:08:34Z").
							CheckEqual("certificates.0.created_by", "example user").
							CheckEqual("certificates.0.description", "Example Certificate Expired").
							CheckEqual("certificates.0.serial_number", "ABCDEF1234567890").
							CheckEqual("certificates.1.created_date", "2025-01-01T12:08:34Z").
							CheckEqual("certificates.1.start_date", "2025-01-01T12:08:34Z").
							CheckEqual("certificates.1.end_date", "2025-07-10T12:08:34Z").
							CheckEqual("certificates.1.created_by", "example user").
							CheckEqual("certificates.1.description", "Example Certificate Expiring").
							CheckEqual("certificates.1.serial_number", "1234567890ABCDEF").Build(),
					),
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
		                  expired = true
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
					ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[(id,name|name,id)\] is required`),
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
		"validation error - invalid value for `expiry_threshold_in_days`": {
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
		                  expired = true
						  expiry_threshold_in_days = "-20"
						}`,
					ExpectError: regexp.MustCompile(`Attribute expiry_threshold_in_days value must be at least 0, got: -20`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
			if test.init != nil {
				test.init(client, test.testData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockGetCASetVersionCertificates(m *mtlstruststore.Mock) {
	getResponse := &mtlstruststore.GetCASetVersionCertificatesResponse{
		CASetID:   "12345",
		CASetName: "example-ca-set",
		Version:   1,
		Certificates: []mtlstruststore.CertificateResponse{
			{
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
			},
			{
				Subject:            "CN=api.example.org, O=Example Org, C=US",
				Issuer:             "CN=Example CA, O=Example Org, C=US",
				StartDate:          time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
				EndDate:            time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
				Fingerprint:        "11:22:33:44:55:66:77:88:99:00:AA:BB:CC:DD:EE:FF:11:22:33:44",
				CertificatePEM:     "-----BEGIN CERTIFICATE-----\nMIID...FAKE...API==\n-----END CERTIFICATE-----",
				SerialNumber:       "ABCDEF1234567890",
				SignatureAlgorithm: "SHA384WITHRSA",
				CreatedDate:        time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				CreatedBy:          "system@example.org",
				Description:        ptr.To("API certificate for internal usage"),
			},
		},
	}
	m.On("GetCASetVersionCertificates", testutils.MockContext, mtlstruststore.GetCASetVersionCertificatesRequest{
		CASetID: "12345",
		Version: 1,
	}).Return(getResponse, nil).Times(3)
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
