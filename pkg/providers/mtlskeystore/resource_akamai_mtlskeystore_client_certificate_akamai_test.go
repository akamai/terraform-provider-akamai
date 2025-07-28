package mtlskeystore

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	clientCertificateData struct {
		certificateName    string
		certificateID      int64
		contractID         string
		geography          string
		groupID            int64
		keyAlgorithm       string
		notificationEmails []string
		secureNetwork      string
		subject            string
		createdBy          string
		createdDate        string
		versions           []clientCertificateVersionData
	}

	clientCertificateVersionData struct {
		version                  int64
		status                   string
		expiryDate               string
		issuer                   string
		keyAlgorithm             string
		certificateSubmittedBy   string
		certificateSubmittedDate string
		createdBy                string
		createdDate              string
		deleteRequestedDate      string
		deployedDate             string
		issuedDate               string
		ellipticCurve            string
		keySizeInBytes           string
		scheduledDeleteDate      string
		signatureAlgorithm       string
		subject                  string
		versionGUID              string
		certificateBlock         certificateBlock
	}
)

var (
	testClientCertificateWithoutSubject = clientCertificateData{
		certificateName: "test-certificate",
		certificateID:   123456789,
		contractID:      "123456789",
		geography:       "CORE",
		groupID:         987654321,
		notificationEmails: []string{
			"testemail1@example.com",
			"testemail2@example.com",
		},
		secureNetwork: "STANDARD_TLS",
		createdBy:     "joeDoe",
		createdDate:   "2023-01-01T12:00:00Z",
		versions: []clientCertificateVersionData{
			{
				version:            1,
				status:             "DEPLOYED",
				expiryDate:         "2024-12-31T23:59:59Z",
				issuer:             "Example Issuer",
				keyAlgorithm:       "RSA",
				createdBy:          "joeDoe",
				createdDate:        "2023-01-01T12:00:00Z",
				deployedDate:       "2023-01-02T12:00:00Z",
				issuedDate:         "2023-01-01T12:00:00Z",
				keySizeInBytes:     "2048",
				subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
				versionGUID:        "test_identifier_1-1",
				signatureAlgorithm: "SHA256_WITH_RSA",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
			},
		},
	}

	testClientCertificateWithOptionals = func() clientCertificateData {
		return clientCertificateData{
			certificateName: "test-certificate",
			certificateID:   123456789,
			contractID:      "123456789",
			geography:       "CORE",
			groupID:         987654321,
			keyAlgorithm:    "RSA",
			notificationEmails: []string{
				"testemail1@example.com",
				"testemail2@example.com",
			},
			secureNetwork: "STANDARD_TLS",
			createdBy:     "joeDoe",
			createdDate:   "2023-01-01T12:00:00Z",
			subject:       "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
			versions: []clientCertificateVersionData{
				{
					version:            1,
					status:             "DEPLOYED",
					expiryDate:         "2024-12-31T23:59:59Z",
					issuer:             "Example Issuer",
					keyAlgorithm:       "RSA",
					createdBy:          "joeDoe",
					createdDate:        "2023-01-01T12:00:00Z",
					deployedDate:       "2023-01-02T12:00:00Z",
					issuedDate:         "2023-01-01T12:00:00Z",
					keySizeInBytes:     "2048",
					subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
					versionGUID:        "test_identifier_1-1",
					signatureAlgorithm: "SHA256_WITH_RSA",
					certificateBlock: certificateBlock{
						certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					},
				},
			},
		}
	}

	testClientCertificateWithContractPrefix = clientCertificateData{
		certificateName: "test-certificate",
		certificateID:   123456789,
		contractID:      "ctr_123456789",
		geography:       "CORE",
		groupID:         987654321,
		notificationEmails: []string{
			"testemail1@example.com",
			"testemail2@example.com",
		},
		secureNetwork: "STANDARD_TLS",
		createdBy:     "joeDoe",
		createdDate:   "2023-01-01T12:00:00Z",
		versions: []clientCertificateVersionData{
			{
				version:            1,
				status:             "DEPLOYED",
				expiryDate:         "2024-12-31T23:59:59Z",
				issuer:             "Example Issuer",
				keyAlgorithm:       "RSA",
				createdBy:          "joeDoe",
				createdDate:        "2023-01-01T12:00:00Z",
				deployedDate:       "2023-01-02T12:00:00Z",
				issuedDate:         "2023-01-01T12:00:00Z",
				keySizeInBytes:     "2048",
				subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
				versionGUID:        "test_identifier_1-1",
				signatureAlgorithm: "SHA256_WITH_RSA",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
			},
		},
	}

	testClientCertificateWith2Versions = func() clientCertificateData {
		return clientCertificateData{
			certificateName: "test-certificate",
			certificateID:   123456789,
			contractID:      "123456789",
			geography:       "CORE",
			groupID:         987654321,
			notificationEmails: []string{
				"testemail1@example.com",
				"testemail2@example.com",
			},
			secureNetwork: "STANDARD_TLS",
			createdBy:     "joeDoe",
			createdDate:   "2023-01-01T12:00:00Z",
			subject:       "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
			versions: []clientCertificateVersionData{
				{
					version:            2,
					status:             "DEPLOYED",
					expiryDate:         "2025-12-31T23:59:59Z",
					issuer:             "Example Issuer",
					keyAlgorithm:       "RSA",
					createdBy:          "joeDoe",
					createdDate:        "2025-01-01T12:00:00Z",
					deployedDate:       "2025-01-02T12:00:00Z",
					issuedDate:         "2025-01-01T12:00:00Z",
					keySizeInBytes:     "2048",
					subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
					versionGUID:        "test_identifier_1-2",
					signatureAlgorithm: "SHA256_WITH_RSA",
					certificateBlock: certificateBlock{
						certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					},
				},
				{
					version:            1,
					status:             "DEPLOYED",
					expiryDate:         "2024-12-31T23:59:59Z",
					issuer:             "Example Issuer",
					keyAlgorithm:       "RSA",
					createdBy:          "joeDoe",
					createdDate:        "2023-01-01T12:00:00Z",
					deployedDate:       "2023-01-02T12:00:00Z",
					issuedDate:         "2023-01-01T12:00:00Z",
					keySizeInBytes:     "2048",
					subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
					versionGUID:        "test_identifier_1-1",
					signatureAlgorithm: "SHA256_WITH_RSA",
					certificateBlock: certificateBlock{
						certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					},
				},
			},
		}
	}

	testClientCertificateWith2RotatedVersions = clientCertificateData{
		certificateName: "test-certificate",
		certificateID:   123456789,
		contractID:      "123456789",
		geography:       "CORE",
		groupID:         987654321,
		notificationEmails: []string{
			"testemail1@example.com",
			"testemail2@example.com",
		},
		secureNetwork: "STANDARD_TLS",
		createdBy:     "joeDoe",
		createdDate:   "2023-01-01T12:00:00Z",
		subject:       "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
		versions: []clientCertificateVersionData{
			{
				version:            3,
				status:             "DEPLOYED",
				expiryDate:         "2025-12-31T23:59:59Z",
				issuer:             "Example Issuer",
				keyAlgorithm:       "RSA",
				createdBy:          "joeDoe",
				createdDate:        "2025-01-01T12:00:00Z",
				deployedDate:       "2025-01-02T12:00:00Z",
				issuedDate:         "2025-01-01T12:00:00Z",
				keySizeInBytes:     "2048",
				subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
				versionGUID:        "test_identifier_1-3",
				signatureAlgorithm: "SHA256_WITH_RSA",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
			},
			{
				version:            2,
				status:             "DEPLOYED",
				expiryDate:         "2025-12-31T23:59:59Z",
				issuer:             "Example Issuer",
				keyAlgorithm:       "RSA",
				createdBy:          "joeDoe",
				createdDate:        "2025-01-01T12:00:00Z",
				deployedDate:       "2025-01-02T12:00:00Z",
				issuedDate:         "2025-01-01T12:00:00Z",
				keySizeInBytes:     "2048",
				subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
				versionGUID:        "test_identifier_1-2",
				signatureAlgorithm: "SHA256_WITH_RSA",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
			},
		},
	}

	testClientCertificateDriftNameAndEmails = clientCertificateData{
		certificateName: "test-certificate-drift",
		certificateID:   123456789,
		contractID:      "123456789",
		geography:       "CORE",
		groupID:         987654321,
		notificationEmails: []string{
			"testemail3@example.com",
			"testemail4@example.com",
		},
		secureNetwork: "STANDARD_TLS",
		createdBy:     "joeDoe",
		createdDate:   "2023-01-01T12:00:00Z",
		versions: []clientCertificateVersionData{
			{
				version:            1,
				status:             "DEPLOYED",
				expiryDate:         "2024-12-31T23:59:59Z",
				issuer:             "Example Issuer",
				keyAlgorithm:       "RSA",
				createdBy:          "joeDoe",
				createdDate:        "2023-01-01T12:00:00Z",
				deployedDate:       "2023-01-02T12:00:00Z",
				issuedDate:         "2023-01-01T12:00:00Z",
				keySizeInBytes:     "2048",
				subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/",
				versionGUID:        "test_identifier_1-1",
				signatureAlgorithm: "SHA256_WITH_RSA",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
			},
		},
	}

	testClientCertificateMissedContractAndGroupInSubject = clientCertificateData{
		certificateName: "test-certificate",
		certificateID:   123456789,
		contractID:      "G-12RS3N4",
		geography:       "CORE",
		groupID:         123456,
		notificationEmails: []string{
			"testemail1@example.com",
			"testemail2@example.com",
		},
		secureNetwork: "STANDARD_TLS",
		createdBy:     "joeDoe",
		createdDate:   "2023-01-01T12:00:00Z",
		subject:       "/C=US/O=Akamai Technologies, Inc./OU=Example /CN=test-certificate/",
		keyAlgorithm:  "RSA",
		versions: []clientCertificateVersionData{
			{
				version:            1,
				status:             "DEPLOYED",
				expiryDate:         "2024-12-31T23:59:59Z",
				issuer:             "Example Issuer",
				keyAlgorithm:       "RSA",
				createdBy:          "joeDoe",
				createdDate:        "2023-01-01T12:00:00Z",
				deployedDate:       "2023-01-02T12:00:00Z",
				issuedDate:         "2023-01-01T12:00:00Z",
				keySizeInBytes:     "2048",
				subject:            "/C=US/O=Akamai Technologies, Inc./OU=Example ctr_123456789 987654321/CN=test-certificate/",
				versionGUID:        "test_identifier_1-1",
				signatureAlgorithm: "SHA256_WITH_RSA",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
			},
		},
	}
)

func TestClientCertificateAkamaiResource(t *testing.T) {
	t.Parallel()
	baseChecker := tst.NewStateChecker("akamai_mtlskeystore_client_certificate_akamai.test").
		CheckEqual("certificate_id", "123456789").
		CheckEqual("certificate_name", "test-certificate").
		CheckEqual("contract_id", "123456789").
		CheckEqual("geography", "CORE").
		CheckEqual("group_id", "987654321").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "testemail1@example.com").
		CheckEqual("notification_emails.1", "testemail2@example.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("created_by", "joeDoe").
		CheckEqual("created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.#", "1").
		CheckEqual("versions.0.version", "1").
		CheckEqual("versions.0.status", "DEPLOYED").
		CheckEqual("versions.0.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.0.issuer", "Example Issuer").
		CheckEqual("versions.0.key_algorithm", "RSA").
		CheckEqual("versions.0.created_by", "joeDoe").
		CheckEqual("versions.0.created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.0.issued_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.0.key_size_in_bytes", "2048").
		CheckEqual("versions.0.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.0.version_guid", "test_identifier_1-1").
		CheckEqual("versions.0.signature_algorithm", "SHA256_WITH_RSA")

	v2v1VersionsChecker := tst.NewStateChecker("akamai_mtlskeystore_client_certificate_akamai.test").
		CheckEqual("certificate_id", "123456789").
		CheckEqual("certificate_name", "test-certificate").
		CheckEqual("contract_id", "123456789").
		CheckEqual("geography", "CORE").
		CheckEqual("group_id", "987654321").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "testemail1@example.com").
		CheckEqual("notification_emails.1", "testemail2@example.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("created_by", "joeDoe").
		CheckEqual("created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.#", "2").
		CheckEqual("versions.0.version", "2").
		CheckEqual("versions.0.status", "DEPLOYED").
		CheckEqual("versions.0.expiry_date", "2025-12-31T23:59:59Z").
		CheckEqual("versions.0.issuer", "Example Issuer").
		CheckEqual("versions.0.key_algorithm", "RSA").
		CheckEqual("versions.0.created_by", "joeDoe").
		CheckEqual("versions.0.created_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.0.issued_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.0.key_size_in_bytes", "2048").
		CheckEqual("versions.0.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.0.version_guid", "test_identifier_1-2").
		CheckEqual("versions.0.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.1.version", "1").
		CheckEqual("versions.1.status", "DEPLOYED").
		CheckEqual("versions.1.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.1.issuer", "Example Issuer").
		CheckEqual("versions.1.key_algorithm", "RSA").
		CheckEqual("versions.1.created_by", "joeDoe").
		CheckEqual("versions.1.created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.1.issued_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.1.key_size_in_bytes", "2048").
		CheckEqual("versions.1.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.1.version_guid", "test_identifier_1-1").
		CheckEqual("versions.1.signature_algorithm", "SHA256_WITH_RSA")

	v3v2VersionsChecker := tst.NewStateChecker("akamai_mtlskeystore_client_certificate_akamai.test").
		CheckEqual("certificate_id", "123456789").
		CheckEqual("certificate_name", "test-certificate").
		CheckEqual("contract_id", "123456789").
		CheckEqual("geography", "CORE").
		CheckEqual("group_id", "987654321").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "testemail1@example.com").
		CheckEqual("notification_emails.1", "testemail2@example.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("created_by", "joeDoe").
		CheckEqual("created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.#", "2").
		CheckEqual("versions.0.version", "3").
		CheckEqual("versions.0.status", "DEPLOYED").
		CheckEqual("versions.0.expiry_date", "2025-12-31T23:59:59Z").
		CheckEqual("versions.0.issuer", "Example Issuer").
		CheckEqual("versions.0.key_algorithm", "RSA").
		CheckEqual("versions.0.created_by", "joeDoe").
		CheckEqual("versions.0.created_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.0.issued_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.0.key_size_in_bytes", "2048").
		CheckEqual("versions.0.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.0.version_guid", "test_identifier_1-3").
		CheckEqual("versions.0.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.1.version", "2").
		CheckEqual("versions.1.status", "DEPLOYED").
		CheckEqual("versions.1.expiry_date", "2025-12-31T23:59:59Z").
		CheckEqual("versions.1.issuer", "Example Issuer").
		CheckEqual("versions.1.key_algorithm", "RSA").
		CheckEqual("versions.1.created_by", "joeDoe").
		CheckEqual("versions.1.created_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.1.issued_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.1.key_size_in_bytes", "2048").
		CheckEqual("versions.1.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.1.version_guid", "test_identifier_1-2").
		CheckEqual("versions.1.signature_algorithm", "SHA256_WITH_RSA")

	tests := map[string]struct {
		init           func(*mtlskeystore.Mock, clientCertificateData, clientCertificateData)
		mockData       clientCertificateData
		mockUpdateData clientCertificateData
		steps          []resource.TestStep
		error          *regexp.Regexp
	}{
		"happy path - create client certificate without optionals": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"happy path - create client certificate with optionals": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithOptionals(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create_with_optionals.tf"),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"happy path - create client certificate with prefix on contract_id": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithContractPrefix,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create_with_contract_prefix.tf"),
					Check: baseChecker.
						CheckEqual("contract_id", "ctr_123456789").
						Build(),
				},
			},
		},
		"happy path - update client certificate notification e-mails only": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation, key algorithm is set as RSA if not provided in create request
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				// Read x2
				mockGetClientCertificateAkamai(m, testData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Twice()
				//Update - only emails changed
				mockPatchClientCertificateAkamai(m, testData, nil, []string{"testemail1@example.com", "testemail2@example.com", "testemail3@example.com"}).Once()
				testData.notificationEmails = []string{"testemail1@example.com", "testemail2@example.com", "testemail3@example.com"}
				mockGetClientCertificateAkamai(m, testData).Once()
				// Read after Emails update
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/update_notification_emails.tf"),
					Check: baseChecker.
						CheckEqual("notification_emails.#", "3").
						CheckEqual("notification_emails.0", "testemail1@example.com").
						CheckEqual("notification_emails.1", "testemail2@example.com").
						CheckEqual("notification_emails.2", "testemail3@example.com").
						Build(),
				},
			},
		},
		"happy path - update client certificate name update": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation, key algorithm is set as RSA if not provided in create request
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				// Read x2
				mockGetClientCertificateAkamai(m, testData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Twice()
				//Update - only name changed
				mockPatchClientCertificateAkamai(m, testData, ptr.To("test-certificate-changed-name"), testData.notificationEmails).Once()
				testData.certificateName = "test-certificate-changed-name"
				mockGetClientCertificateAkamai(m, testData).Once()
				// Read after name update
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/update_name.tf"),
					Check: baseChecker.
						CheckEqual("certificate_name", "test-certificate-changed-name").
						Build(),
				},
			},
		},
		"happy path - update client certificate name update and e-mails": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation, key algorithm is set as RSA if not provided in create request
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				// Read x2
				mockGetClientCertificateAkamai(m, testData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Twice()
				//Update - name and e-mails changed
				mockPatchClientCertificateAkamai(m, testData, ptr.To("test-certificate-changed-name"), []string{"testemail1@example.com", "testemail2@example.com", "testemail3@example.com"}).Once()
				testData.certificateName = "test-certificate-changed-name"
				testData.notificationEmails = []string{"testemail1@example.com", "testemail2@example.com", "testemail3@example.com"}
				mockGetClientCertificateAkamai(m, testData).Once()
				// Read after name update
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/update_name_and_emails.tf"),
					Check: baseChecker.
						CheckEqual("certificate_name", "test-certificate-changed-name").
						CheckEqual("notification_emails.#", "3").
						CheckEqual("notification_emails.0", "testemail1@example.com").
						CheckEqual("notification_emails.1", "testemail2@example.com").
						CheckEqual("notification_emails.2", "testemail3@example.com").
						Build(),
				},
			},
		},
		"happy path - 1 version created, second added during automatic rotation ": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, testUpdateData clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read after rotation
				// subject is set after automatic rotation
				testUpdateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testUpdateData.keyAlgorithm = "RSA"
				mockGetClientCertificateAkamai(m, testUpdateData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testUpdateData.versions, testUpdateData.certificateID).Twice()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testUpdateData.versions, testUpdateData.certificateID).Once()
				mockDeleteClientCertificateAkamaiVersion(m, testUpdateData.versions, testUpdateData.certificateID, 2).Once()
				mockDeleteClientCertificateAkamaiVersion(m, testUpdateData.versions, testUpdateData.certificateID, 1).Once()
			},
			mockData:       testClientCertificateWithoutSubject,
			mockUpdateData: testClientCertificateWith2Versions(),
			steps: []resource.TestStep{
				{
					Config:  testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Destroy: false,
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: v2v1VersionsChecker.
						Build(),
				},
			},
		},
		"happy path - 1 version created, both versions automatically rotated": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, testUpdateData clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read after rotation
				// subject is set after automatic rotation
				testUpdateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testUpdateData.keyAlgorithm = "RSA"
				mockGetClientCertificateAkamai(m, testUpdateData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testUpdateData.versions, testUpdateData.certificateID).Twice()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testUpdateData.versions, testUpdateData.certificateID).Once()
				mockDeleteClientCertificateAkamaiVersion(m, testUpdateData.versions, testUpdateData.certificateID, 2).Once()
				mockDeleteClientCertificateAkamaiVersion(m, testUpdateData.versions, testUpdateData.certificateID, 1).Once()
			},
			mockData:       testClientCertificateWithoutSubject,
			mockUpdateData: testClientCertificateWith2RotatedVersions,
			steps: []resource.TestStep{
				{
					Config:  testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Destroy: false,
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: v3v2VersionsChecker.
						Build(),
				},
			},
		},
		"happy path - drift on name and e-mails": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, testUpdateData clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read after rotation
				// subject is set after automatic rotation
				testUpdateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testUpdateData.keyAlgorithm = "RSA"
				mockGetClientCertificateAkamai(m, testUpdateData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testUpdateData.versions, testUpdateData.certificateID).Twice()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testUpdateData.versions, testUpdateData.certificateID).Once()
				mockDeleteClientCertificateAkamaiVersion(m, testUpdateData.versions, testUpdateData.certificateID, 1).Once()
			},
			mockData:       testClientCertificateWithoutSubject,
			mockUpdateData: testClientCertificateDriftNameAndEmails,
			steps: []resource.TestStep{
				{
					Config:  testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Destroy: false,
					Check: baseChecker.
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/drift_name_emails.tf"),
					Check: baseChecker.
						CheckEqual("certificate_name", "test-certificate-drift").
						CheckEqual("notification_emails.#", "2").
						CheckEqual("notification_emails.0", "testemail3@example.com").
						CheckEqual("notification_emails.1", "testemail4@example.com").
						Build(),
				},
			},
		},
		"happy path - check poling mechanism": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				pollingDuration = 1 * time.Second
				mockCreateClientCertificateAkamai(m, testData)
				// first tick of polling mechanism
				testData.versions[0].status = "DEPLOYMENT_PENDING"
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// second tick of polling mechanism
				testData.versions[0].status = "DEPLOYED"
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
			},
		},
		"error - fail on update contract": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				// Read
				mockGetClientCertificateAkamai(m, testData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/error_update_contract.tf"),
					ExpectError: regexp.MustCompile("updating field `contract_id` is not possible"),
				},
			},
		},
		"error - fail on update group": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				// Read
				mockGetClientCertificateAkamai(m, testData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/error_update_group.tf"),
					ExpectError: regexp.MustCompile("updating field `group_id` is not possible"),
				},
			},
		},
		"error - fail on update geography": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				// Read
				mockGetClientCertificateAkamai(m, testData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/error_update_geography.tf"),
					ExpectError: regexp.MustCompile("updating field `geography` is not possible"),
				},
			},
		},
		"error - fail on update subject": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Times(2)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Times(2)
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithOptionals(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create_with_optionals.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/error_update_subject.tf"),
					ExpectError: regexp.MustCompile("updating field `subject` is not possible"),
				},
			},
		},
		"error - fail on update key_algorithm": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Times(2)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Times(2)
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithOptionals(),
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create_with_optionals.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/error_update_key_algorithm.tf"),
					ExpectError: regexp.MustCompile("updating field `key_algorithm` is not possible"),
				},
			},
		},
		"error - fail on update secure_network": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData, _ clientCertificateData) {
				mockCreateClientCertificateAkamai(m, testData)
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				// Read
				mockGetClientCertificateAkamai(m, testData).Twice()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateAkamaiVersion(m, testData.versions, testData.certificateID, 1)
			},
			mockData: testClientCertificateWithoutSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/create.tf"),
					Check: baseChecker.
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/error_update_secure_network.tf"),
					ExpectError: regexp.MustCompile("updating field `secure_network` is not possible"),
				},
			},
		},
		"error - invalid key algorithm": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/invalid_key_algorithm.tf"),
					ExpectError: regexp.MustCompile(`Attribute key_algorithm value must be one of: \["RSA" "ECDSA"], got: "AAA"`),
				},
			},
		},
		"error - invalid secure network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/invalid_secure_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute secure_network value must be one of: \["STANDARD_TLS"\n"ENHANCED_TLS"], got: "WRONG_NETWORK"`),
				},
			},
		},
		"error - invalid geography": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/invalid_geography.tf"),
					ExpectError: regexp.MustCompile(`Attribute geography value must be one of: \["CHINA_AND_CORE" "RUSSIA_AND_CORE"\n"CORE"], got: "EUROPE"`),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.mockData, tc.mockUpdateData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}

}

func TestClientCertificateAkamaiResource_ImportState(t *testing.T) {
	t.Parallel()
	baseChecker := tst.NewImportChecker().
		CheckEqual("certificate_id", "123456789").
		CheckEqual("certificate_name", "test-certificate").
		CheckEqual("contract_id", "123456789").
		CheckEqual("geography", "CORE").
		CheckEqual("group_id", "987654321").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "testemail1@example.com").
		CheckEqual("notification_emails.1", "testemail2@example.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("created_by", "joeDoe").
		CheckEqual("created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.#", "1").
		CheckEqual("versions.0.version", "1").
		CheckEqual("versions.0.status", "DEPLOYED").
		CheckEqual("versions.0.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.0.issuer", "Example Issuer").
		CheckEqual("versions.0.key_algorithm", "RSA").
		CheckEqual("versions.0.created_by", "joeDoe").
		CheckEqual("versions.0.created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.0.issued_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.0.key_size_in_bytes", "2048").
		CheckEqual("versions.0.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.0.version_guid", "test_identifier_1-1").
		CheckEqual("versions.0.signature_algorithm", "SHA256_WITH_RSA")

	secnondVersionChecker := baseChecker.
		CheckEqual("versions.#", "2").
		CheckEqual("versions.0.version", "2").
		CheckEqual("versions.0.status", "DEPLOYED").
		CheckEqual("versions.0.expiry_date", "2025-12-31T23:59:59Z").
		CheckEqual("versions.0.issuer", "Example Issuer").
		CheckEqual("versions.0.key_algorithm", "RSA").
		CheckEqual("versions.0.created_by", "joeDoe").
		CheckEqual("versions.0.created_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.0.issued_date", "2025-01-01T12:00:00Z").
		CheckEqual("versions.0.key_size_in_bytes", "2048").
		CheckEqual("versions.0.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.0.version_guid", "test_identifier_1-2").
		CheckEqual("versions.0.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.1.version", "1").
		CheckEqual("versions.1.status", "DEPLOYED").
		CheckEqual("versions.1.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.1.issuer", "Example Issuer").
		CheckEqual("versions.1.key_algorithm", "RSA").
		CheckEqual("versions.1.created_by", "joeDoe").
		CheckEqual("versions.1.created_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.1.issued_date", "2023-01-01T12:00:00Z").
		CheckEqual("versions.1.key_size_in_bytes", "2048").
		CheckEqual("versions.1.subject", "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/").
		CheckEqual("versions.1.version_guid", "test_identifier_1-1").
		CheckEqual("versions.1.signature_algorithm", "SHA256_WITH_RSA")

	tests := map[string]struct {
		init     func(*mtlskeystore.Mock, clientCertificateData)
		mockData clientCertificateData
		steps    []resource.TestStep
		error    *regexp.Regexp
	}{
		"import - client certificate with one version": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				// Import
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
			},
			mockData: testClientCertificateWithOptionals(),
			steps: []resource.TestStep{
				{
					ImportStateCheck: baseChecker.Build(),
					ImportStateId:    "123456789",
					ImportState:      true,
					ResourceName:     "akamai_mtlskeystore_client_certificate_akamai.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"import - client certificate with two versions": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				// Import
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
			},
			mockData: testClientCertificateWith2Versions(),
			steps: []resource.TestStep{
				{
					ImportStateCheck: secnondVersionChecker.Build(),
					ImportStateId:    "123456789",
					ImportState:      true,
					ResourceName:     "akamai_mtlskeystore_client_certificate_akamai.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"import - no group and contract in certificate subject, but provided with importID": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				// Import
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
			},
			mockData: testClientCertificateMissedContractAndGroupInSubject,
			steps: []resource.TestStep{
				{
					ImportStateCheck: tst.NewImportChecker().
						CheckEqual("certificate_id", "123456789").
						CheckEqual("group_id", "123456").
						CheckEqual("contract_id", "G-12RS3N4").Build(),
					ImportStateId: "123456789,123456,G-12RS3N4",
					ImportState:   true,
					ResourceName:  "akamai_mtlskeystore_client_certificate_akamai.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"import - client certificate with two versions with one DELETE_PENDING": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				// Import
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				testData.versions[0].status = "DELETE_PENDING"
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
				// Read
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
			},
			mockData: testClientCertificateWith2Versions(),
			steps: []resource.TestStep{
				{
					ImportStateCheck: secnondVersionChecker.CheckEqual("versions.0.status", "DELETE_PENDING").Build(),
					ImportStateId:    "123456789",
					ImportState:      true,
					ResourceName:     "akamai_mtlskeystore_client_certificate_akamai.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"error - problem with parsing custom subject": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				// Import
				// subject is set after creation
				testData.subject = "/CN=NON-PARSABLE SUBJECT"
				testData.keyAlgorithm = "RSA"
				mockGetClientCertificateAkamai(m, testData).Once()
			},
			mockData: testClientCertificateWithOptionals(),
			error:    regexp.MustCompile(`parsing subject "/CN=NON-PARSABLE SUBJECT": invalid subject format`),
			steps: []resource.TestStep{
				{
					ImportStateCheck: secnondVersionChecker.Build(),
					ImportStateId:    "123456789",
					ImportState:      true,
					ExpectError:      regexp.MustCompile(`since it is not possible to extract contract and group from certificate\nsubject, you need to provide an importID in the format\n'certificateID,groupID,contractID'. Where certificate, groupID and contractID\nare required`),
					ResourceName:     "akamai_mtlskeystore_client_certificate_akamai.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"error - Get Client Certificate failed": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				// Import
				// subject is set after creation
				testData.subject = "/CN=NON-PARSABLE SUBJECT"
				testData.keyAlgorithm = "RSA"
				m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
					CertificateID: testData.certificateID,
				}).Return(nil, fmt.Errorf("unable to Get Client Certificate")).Once()
			},
			mockData: testClientCertificateWithOptionals(),
			error:    regexp.MustCompile(`Unable to Get Client Certificate`),
			steps: []resource.TestStep{
				{
					ImportStateCheck: secnondVersionChecker.Build(),
					ImportStateId:    "123456789",
					ImportState:      true,
					ExpectError:      regexp.MustCompile("Unable to Get Client Certificate"),
					ResourceName:     "akamai_mtlskeystore_client_certificate_akamai.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"error - no group and contract in certificate subject": {
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				mockGetClientCertificateAkamai(m, testData).Once()
			},
			mockData: testClientCertificateMissedContractAndGroupInSubject,
			steps: []resource.TestStep{
				{
					ImportStateId: "123456789",
					ImportState:   true,
					ResourceName:  "akamai_mtlskeystore_client_certificate_akamai.test",
					ExpectError:   regexp.MustCompile(`since it is not possible to extract contract and group from certificate\nsubject, you need to provide an importID in the format\n'certificateID,groupID,contractID'. Where certificate, groupID and contractID\nare required`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"error - incorrect number of parts in importID": {
			mockData: testClientCertificateMissedContractAndGroupInSubject,
			steps: []resource.TestStep{
				{
					ImportStateId: "123456789,123456,G-12RS3N4,123",
					ImportState:   true,
					ResourceName:  "akamai_mtlskeystore_client_certificate_akamai.test",
					ExpectError:   regexp.MustCompile(`you need to provide an importID in the format\n'certificateID,\[groupID,contractID]'. Where certificateID is required and\ngroupID and contractID are optional`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
		"error - importing one version DELETE_PENDING": {
			// Default subject is returned.
			init: func(m *mtlskeystore.Mock, testData clientCertificateData) {
				// Import
				// subject is set after creation
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
				testData.keyAlgorithm = "RSA"
				testData.versions[0].status = "DELETE_PENDING"
				mockGetClientCertificateAkamai(m, testData).Once()
				mockListClientCertificateAkamaiVersions(t, m, testData.versions, testData.certificateID).Once()
			},
			mockData: testClientCertificateWithOptionals(),
			steps: []resource.TestStep{
				{
					ImportStateId: "123456789",
					ImportState:   true,
					ExpectError:   regexp.MustCompile("Certificate in Delete Pending State"),
					ResourceName:  "akamai_mtlskeystore_client_certificate_akamai.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateAkamai/import_one_version.tf"),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.mockData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}

}

func mockPatchClientCertificateAkamai(m *mtlskeystore.Mock, data clientCertificateData, newName *string, emails []string) *mock.Call {
	request := mtlskeystore.PatchClientCertificateRequest{
		CertificateID: data.certificateID,
	}
	if newName != nil && data.certificateName != *newName {
		request.Body.CertificateName = newName
	}
	if len(emails) > 0 && !slices.Equal(data.notificationEmails, emails) {
		request.Body.NotificationEmails = emails
	}
	return m.On("PatchClientCertificate", testutils.MockContext, request).Return(nil).Once()
}

func mockCreateClientCertificateAkamai(m *mtlskeystore.Mock, testData clientCertificateData) *mock.Call {
	request := mtlskeystore.CreateClientCertificateRequest{
		CertificateName:    testData.certificateName,
		ContractID:         strings.TrimPrefix(testData.contractID, "ctr_"),
		Geography:          mtlskeystore.Geography(testData.geography),
		GroupID:            testData.groupID,
		NotificationEmails: testData.notificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(testData.secureNetwork),
		Signer:             mtlskeystore.SignerAkamai,
		Subject:            ptr.To(testData.subject),
	}

	response := mtlskeystore.CreateClientCertificateResponse{
		CertificateID:      testData.certificateID,
		CertificateName:    testData.certificateName,
		Geography:          testData.geography,
		NotificationEmails: testData.notificationEmails,
		SecureNetwork:      testData.secureNetwork,
		Signer:             string(mtlskeystore.SignerAkamai),
		Subject:            testData.subject,
		CreatedBy:          testData.createdBy,
		CreatedDate:        test.NewTimeFromStringMust(testData.createdDate),
	}

	if testData.subject == "" {
		response.Subject = "/C=US/O=Akamai Technologies, Inc./OU=Example 123456789 987654321/CN=test-certificate/"
	}

	if testData.keyAlgorithm == "" {
		testData.keyAlgorithm = "RSA"
		response.KeyAlgorithm = testData.keyAlgorithm
	} else {
		request.KeyAlgorithm = ptr.To(mtlskeystore.CryptographicAlgorithm(testData.keyAlgorithm))
		response.KeyAlgorithm = testData.keyAlgorithm
	}

	return m.On("CreateClientCertificate", testutils.MockContext, request).Return(&response, nil).Once()
}

func mockGetClientCertificateAkamai(m *mtlskeystore.Mock, testData clientCertificateData) *mock.Call {
	response := mtlskeystore.GetClientCertificateResponse{
		CertificateID:      testData.certificateID,
		CertificateName:    testData.certificateName,
		CreatedBy:          testData.createdBy,
		CreatedDate:        test.NewTimeFromStringMust(testData.createdDate),
		Geography:          testData.geography,
		KeyAlgorithm:       testData.keyAlgorithm,
		NotificationEmails: testData.notificationEmails,
		SecureNetwork:      testData.secureNetwork,
		Subject:            testData.subject,
	}

	return m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
		CertificateID: testData.certificateID,
	}).Return(&response, nil).Once()
}

func mockListClientCertificateAkamaiVersions(t *testing.T, m *mtlskeystore.Mock, versions []clientCertificateVersionData, certificateID int64) *mock.Call {
	responseVersions := make([]mtlskeystore.ClientCertificateVersion, 0, len(versions))

	for _, version := range versions {
		certificateVersions := mtlskeystore.ClientCertificateVersion{
			Version:            version.version,
			VersionGUID:        version.versionGUID,
			CreatedBy:          version.createdBy,
			CreatedDate:        test.NewTimeFromString(t, version.createdDate),
			ExpiryDate:         ptr.To(test.NewTimeFromString(t, version.expiryDate)),
			IssuedDate:         ptr.To(test.NewTimeFromString(t, version.issuedDate)),
			Issuer:             ptr.To(version.issuer),
			KeyAlgorithm:       version.keyAlgorithm,
			EllipticCurve:      ptr.To(version.ellipticCurve),
			KeySizeInBytes:     ptr.To(version.keySizeInBytes),
			SignatureAlgorithm: ptr.To(version.signatureAlgorithm),
			Status:             version.status,
			Subject:            ptr.To(version.subject),
		}

		if version.certificateBlock != (certificateBlock{}) {
			certificateVersions.CertificateBlock = &mtlskeystore.CertificateBlock{
				Certificate: version.certificateBlock.certificate,
				TrustChain:  version.certificateBlock.trustChain,
			}
		}

		if version.certificateSubmittedBy != "" {
			certificateVersions.CertificateSubmittedBy = ptr.To(version.certificateSubmittedBy)
		}
		if version.certificateSubmittedDate != "" {
			certificateVersions.CertificateSubmittedDate = ptr.To(test.NewTimeFromString(t, version.certificateSubmittedDate))
		}
		if version.deleteRequestedDate != "" {
			certificateVersions.DeleteRequestedDate = ptr.To(test.NewTimeFromString(t, version.deleteRequestedDate))
		}
		if version.scheduledDeleteDate != "" {
			certificateVersions.ScheduledDeleteDate = ptr.To(test.NewTimeFromString(t, version.scheduledDeleteDate))
		}

		responseVersions = append(responseVersions, certificateVersions)
	}

	return m.On("ListClientCertificateVersions", testutils.MockContext, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: certificateID,
	}).Return(&mtlskeystore.ListClientCertificateVersionsResponse{
		Versions: responseVersions,
	}, nil).Once()
}

func mockDeleteClientCertificateAkamaiVersion(m *mtlskeystore.Mock, versions []clientCertificateVersionData, certificateID, version int64) *mock.Call {
	return m.On("DeleteClientCertificateVersion", testutils.MockContext, mtlskeystore.DeleteClientCertificateVersionRequest{
		CertificateID: certificateID,
		Version:       versions[version-1].version,
	}).Return(nil).Once()
}
