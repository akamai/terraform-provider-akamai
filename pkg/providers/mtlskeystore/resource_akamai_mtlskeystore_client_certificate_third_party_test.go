package mtlskeystore

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	commonDataForResource struct {
		certificateName    string
		certificateID      int64
		contractID         string
		geography          string
		groupID            int64
		keyAlgorithm       string
		notificationEmails []string
		secureNetwork      string
		subject            string
		preferredCA        *string
		versions           map[string]versionData
	}

	versionData struct {
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
		keyEllipticCurve         string
		keySizeInBytes           string
		scheduledDeleteDate      string
		signatureAlgorithm       string
		subject                  string
		versionGUID              string
		certificateBlock         certificateBlock
		csrBlock                 csrBlock
	}

	certificateBlock struct {
		certificate string
		trustChain  string
	}

	csrBlock struct {
		csr          string
		keyAlgorithm string
	}
)

var (
	testFiveVersions = commonDataForResource{
		certificateName: "test-certificate",
		certificateID:   12345,
		contractID:      "ctr_12345",
		geography:       "CORE",
		groupID:         1234,
		keyAlgorithm:    "RSA",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		versions: map[string]versionData{
			"v5": {
				version:                  5,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v5-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
			"v4": {
				version:                  4,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v4-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
			"v3": {
				version:                  3,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v3-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
			"v2": testTwoVersionsWithOptionalParams.versions["v2"],
			"v1": testTwoVersionsWithOptionalParams.versions["v1"],
		},
	}

	testTwoVersionsWithOptionalParams = commonDataForResource{
		certificateName: "test-certificate",
		certificateID:   12345,
		contractID:      "ctr_12345",
		geography:       "CORE",
		groupID:         1234,
		keyAlgorithm:    "RSA",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		subject:       "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/",
		versions: map[string]versionData{
			"v2": {
				version:                  2,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "ECDSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keyEllipticCurve:         "test-ecdsa",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v2-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
			"v1": {
				version:                  1,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v1-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
		},
	}

	testOneVersion = commonDataForResource{
		certificateName: "test-certificate",
		certificateID:   12345,
		contractID:      "ctr_12345",
		geography:       "CORE",
		groupID:         1234,
		keyAlgorithm:    "RSA",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		versions: map[string]versionData{
			"v1": {
				version:                  1,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v1-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
		},
	}

	testOneVersionWithSubject = commonDataForResource{
		certificateName: "test-certificate",
		certificateID:   12345,
		contractID:      "ctr_12345",
		geography:       "CORE",
		groupID:         1234,
		keyAlgorithm:    "RSA",
		subject:         "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		versions: map[string]versionData{
			"v1": {
				version:                  1,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v1-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
		},
	}

	testUpdateNotificationEmailsAndCertificateName = commonDataForResource{
		certificateName: "updated-certificate-name",
		certificateID:   12345,
		contractID:      "ctr_12345",
		geography:       "CORE",
		groupID:         1234,
		keyAlgorithm:    "RSA",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith-new@akamai.com",
			"test@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		versions: map[string]versionData{
			"v1": {
				version:                  1,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v1-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
		},
	}

	testUpdateNotificationEmails = commonDataForResource{
		certificateName: "test-certificate",
		certificateID:   12345,
		contractID:      "ctr_12345",
		geography:       "CORE",
		groupID:         1234,
		keyAlgorithm:    "RSA",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith-new@akamai.com",
			"test@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		versions: map[string]versionData{
			"v1": {
				version:                  1,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v1-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
		},
	}

	testUpdateCertificateName = commonDataForResource{
		certificateName: "updated-certificate-name",
		certificateID:   12345,
		contractID:      "ctr_12345",
		geography:       "CORE",
		groupID:         1234,
		keyAlgorithm:    "RSA",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		versions: map[string]versionData{
			"v1": {
				version:                  1,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v1-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
		},
	}

	testClientCertificateMissedContractAndGroupInSubjectTP = commonDataForResource{
		certificateName: "test-certificate",
		certificateID:   12345,
		contractID:      "G-12RS3N4",
		geography:       "CORE",
		groupID:         123456,
		keyAlgorithm:    "RSA",
		notificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		secureNetwork: "STANDARD_TLS",
		subject:       "/C=US/O=Akamai Technologies, Inc./OU=Example /CN=test-certificate/",
		versions: map[string]versionData{
			"v1": {
				version:                  1,
				status:                   "ACTIVE",
				expiryDate:               "2024-12-31T23:59:59Z",
				issuer:                   "Example CA",
				keyAlgorithm:             "RSA",
				certificateSubmittedBy:   "jkowalski",
				certificateSubmittedDate: "2023-01-01T00:00:00Z",
				createdBy:                "jkowalski",
				createdDate:              "2023-01-01T00:00:00Z",
				deployedDate:             "2023-01-02T00:00:00Z",
				issuedDate:               "2023-01-03T00:00:00Z",
				keySizeInBytes:           "2048",
				signatureAlgorithm:       "SHA256_WITH_RSA",
				subject:                  "CN=test.example.com",
				versionGUID:              "v1-guid-12345",
				certificateBlock: certificateBlock{
					certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				csrBlock: csrBlock{
					csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					keyAlgorithm: "RSA",
				},
			},
		},
	}
)

func TestClientCertificateThirdPartyResource(t *testing.T) {
	t.Parallel()
	baseChecker := test.NewStateChecker("akamai_mtlskeystore_client_certificate_third_party.test").
		CheckEqual("certificate_id", "12345").
		CheckEqual("certificate_name", "test-certificate").
		CheckEqual("contract_id", "ctr_12345").
		CheckEqual("geography", "CORE").
		CheckEqual("group_id", "1234").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "jkowalski@akamai.com").
		CheckEqual("notification_emails.1", "jsmith@akamai.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/").
		CheckEqual("versions.v1.version", "1").
		CheckEqual("versions.v1.status", "ACTIVE").
		CheckEqual("versions.v1.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.v1.issuer", "Example CA").
		CheckEqual("versions.v1.key_algorithm", "RSA").
		CheckEqual("versions.v1.certificate_submitted_by", "jkowalski").
		CheckEqual("versions.v1.certificate_submitted_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.v1.created_by", "jkowalski").
		CheckEqual("versions.v1.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.v1.deployed_date", "2023-01-02T00:00:00Z").
		CheckEqual("versions.v1.issued_date", "2023-01-03T00:00:00Z").
		CheckEqual("versions.v1.key_size_in_bytes", "2048").
		CheckEqual("versions.v1.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.v1.subject", "CN=test.example.com").
		CheckEqual("versions.v1.version_guid", "v1-guid-12345").
		CheckEqual("versions.v1.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.v1.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.v1.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
		CheckEqual("versions.v1.csr_block.key_algorithm", "RSA").
		CheckMissing("versions.v1.delete_requested_date").
		CheckMissing("versions.v1.scheduled_delete_date")

	secondVersionChecker := baseChecker.
		CheckEqual("versions.v2.version", "2").
		CheckEqual("versions.v2.status", "ACTIVE").
		CheckEqual("versions.v2.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.v2.issuer", "Example CA").
		CheckEqual("versions.v2.key_algorithm", "ECDSA").
		CheckEqual("versions.v2.certificate_submitted_by", "jkowalski").
		CheckEqual("versions.v2.certificate_submitted_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.v2.created_by", "jkowalski").
		CheckEqual("versions.v2.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.v2.deployed_date", "2023-01-02T00:00:00Z").
		CheckEqual("versions.v2.issued_date", "2023-01-03T00:00:00Z").
		CheckEqual("versions.v2.key_elliptic_curve", "test-ecdsa").
		CheckEqual("versions.v2.key_size_in_bytes", "2048").
		CheckEqual("versions.v2.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.v2.subject", "CN=test.example.com").
		CheckEqual("versions.v2.version_guid", "v2-guid-12345").
		CheckEqual("versions.v2.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.v2.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.v2.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
		CheckEqual("versions.v2.csr_block.key_algorithm", "RSA").
		CheckMissing("versions.v2.delete_requested_date").
		CheckMissing("versions.v2.scheduled_delete_date")

	tests := map[string]struct {
		init           func(*mtlskeystore.Mock, commonDataForResource, commonDataForResource)
		mockCreateData commonDataForResource
		mockUpdateData commonDataForResource
		steps          []resource.TestStep
		error          *regexp.Regexp
	}{
		"happy path - without optional params and one version": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Default subject is returned.
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},
		"happy path - with unprefixed contract_id": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Default subject is returned.
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_unprefixed_contract.tf"),
					Check: baseChecker.
						CheckEqual("contract_id", "12345").
						Build(),
				},
			},
		},
		"happy path - with optional params and multiple versions": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testData.versions, testData.certificateID, "v2")
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v2")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
			},
		},
		"happy path - update notification emails only": {
			init: func(m *mtlskeystore.Mock, testCreateData, testUpdateData commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Default subject is returned.
				testCreateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update
				mockPatchClientCertificate(m, 12345, nil, []string{"jkowalski@akamai.com", "jsmith-new@akamai.com", "test@akamai.com"})
				// Default subject is returned.
				testUpdateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testUpdateData)
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				// Read
				mockGetClientCertificate(m, testUpdateData)
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				mockDeleteClientCertificateVersion(m, testUpdateData.versions, nil, testUpdateData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			mockUpdateData: testUpdateNotificationEmails,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_notification_emails.tf"),
					Check: baseChecker.
						CheckEqual("notification_emails.#", "3").
						CheckEqual("notification_emails.0", "jkowalski@akamai.com").
						CheckEqual("notification_emails.1", "jsmith-new@akamai.com").
						CheckEqual("notification_emails.2", "test@akamai.com").
						Build(),
				},
			},
		},
		"happy path - update certificate name only": {
			init: func(m *mtlskeystore.Mock, testCreateData, testUpdateData commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Default subject is returned.
				testCreateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update
				mockPatchClientCertificate(m, 12345, ptr.To("updated-certificate-name"), nil)
				// Default subject is returned.
				testUpdateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testUpdateData).Twice()
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				mockDeleteClientCertificateVersion(m, testUpdateData.versions, nil, testUpdateData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			mockUpdateData: testUpdateCertificateName,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_certificate_name.tf"),
					Check: baseChecker.
						CheckEqual("certificate_name", "updated-certificate-name").
						Build(),
				},
			},
		},
		"happy path - update certificate name and notification emails": {
			init: func(m *mtlskeystore.Mock, testCreateData, testUpdateData commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Default subject is returned.
				testCreateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update
				mockPatchClientCertificate(m, 12345, ptr.To("updated-certificate-name"), []string{"jkowalski@akamai.com", "jsmith-new@akamai.com", "test@akamai.com"})
				// Default subject is returned.
				testUpdateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testUpdateData)
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				// Read
				mockGetClientCertificate(m, testUpdateData)
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				mockDeleteClientCertificateVersion(m, testUpdateData.versions, nil, testUpdateData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			mockUpdateData: testUpdateNotificationEmailsAndCertificateName,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_certificate_name_and_notification_emails.tf"),
					Check: baseChecker.
						CheckEqual("certificate_name", "updated-certificate-name").
						CheckEqual("notification_emails.#", "3").
						CheckEqual("notification_emails.0", "jkowalski@akamai.com").
						CheckEqual("notification_emails.1", "jsmith-new@akamai.com").
						CheckEqual("notification_emails.2", "test@akamai.com").
						Build(),
				},
			},
		},
		"happy path - add new versions": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v4": {
						version:                  4,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2023-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2023-01-01T00:00:00Z",
						deployedDate:             "2023-01-02T00:00:00Z",
						issuedDate:               "2023-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v4-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v3": {
						version:                  3,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2023-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2023-01-01T00:00:00Z",
						deployedDate:             "2023-01-02T00:00:00Z",
						issuedDate:               "2023-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v3-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v2": testCreateData.versions["v2"],
					"v1": testCreateData.versions["v1"],
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update: Add new versions (Rotate)
				mockRotateClientCertificateVersion(t, m, newVersions, testCreateData.certificateID, "v3")
				mockRotateClientCertificateVersion(t, m, newVersions, testCreateData.certificateID, "v4")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v1")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v2")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v3")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v4")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_new_version.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						CheckEqual("versions.v3.version", "3").
						CheckEqual("versions.v3.status", "ACTIVE").
						CheckEqual("versions.v3.expiry_date", "2025-12-31T23:59:59Z").
						CheckEqual("versions.v3.issuer", "Example CA").
						CheckEqual("versions.v3.key_algorithm", "RSA").
						CheckEqual("versions.v3.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.v3.certificate_submitted_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v3.created_by", "jkowalski").
						CheckEqual("versions.v3.created_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v3.deployed_date", "2023-01-02T00:00:00Z").
						CheckEqual("versions.v3.issued_date", "2023-01-03T00:00:00Z").
						CheckEqual("versions.v3.key_size_in_bytes", "2048").
						CheckEqual("versions.v3.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.v3.subject", "CN=test.example.com").
						CheckEqual("versions.v3.version_guid", "v3-guid-12345").
						CheckEqual("versions.v3.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v3.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v3.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
						CheckEqual("versions.v3.csr_block.key_algorithm", "RSA").
						CheckMissing("versions.v3.delete_requested_date").
						CheckMissing("versions.v3.scheduled_delete_date").
						CheckEqual("versions.v4.version", "4").
						CheckEqual("versions.v4.status", "ACTIVE").
						CheckEqual("versions.v4.expiry_date", "2025-12-31T23:59:59Z").
						CheckEqual("versions.v4.issuer", "Example CA").
						CheckEqual("versions.v4.key_algorithm", "RSA").
						CheckEqual("versions.v4.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.v4.certificate_submitted_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v4.created_by", "jkowalski").
						CheckEqual("versions.v4.created_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v4.deployed_date", "2023-01-02T00:00:00Z").
						CheckEqual("versions.v4.issued_date", "2023-01-03T00:00:00Z").
						CheckEqual("versions.v4.key_size_in_bytes", "2048").
						CheckEqual("versions.v4.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.v4.subject", "CN=test.example.com").
						CheckEqual("versions.v4.version_guid", "v4-guid-12345").
						CheckEqual("versions.v4.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v4.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v4.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
						CheckEqual("versions.v4.csr_block.key_algorithm", "RSA").
						CheckMissing("versions.v4.delete_requested_date").
						CheckMissing("versions.v4.scheduled_delete_date").
						Build(),
				},
			},
		},
		"happy path - remove version": {
			init: func(m *mtlskeystore.Mock, testCreateData, testUpdateData commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update (delete version)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v2")
				// Modify update test data so the custom subject is returned.
				testUpdateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/"
				mockGetClientCertificate(m, testUpdateData)
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				// Read
				mockGetClientCertificate(m, testUpdateData)
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, testUpdateData.versions, testUpdateData.certificateID)
				mockDeleteClientCertificateVersion(m, testUpdateData.versions, nil, testUpdateData.certificateID, "v1")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			mockUpdateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_remove_one_version.tf"),
					Check: baseChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
			},
		},
		"happy path - remove one and add new version": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v3": {
						version:                  3,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2023-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2023-01-01T00:00:00Z",
						deployedDate:             "2023-01-02T00:00:00Z",
						issuedDate:               "2023-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v3-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v1": testCreateData.versions["v1"],
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update (delete version + add new version)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v2")
				mockRotateClientCertificateVersion(t, m, newVersions, testCreateData.certificateID, "v3")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v1")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v3")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_remove_one_and_add_one_version.tf"),
					Check: baseChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						CheckEqual("versions.v3.version", "3").
						CheckEqual("versions.v3.status", "ACTIVE").
						CheckEqual("versions.v3.expiry_date", "2025-12-31T23:59:59Z").
						CheckEqual("versions.v3.issuer", "Example CA").
						CheckEqual("versions.v3.key_algorithm", "RSA").
						CheckEqual("versions.v3.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.v3.certificate_submitted_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v3.created_by", "jkowalski").
						CheckEqual("versions.v3.created_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v3.deployed_date", "2023-01-02T00:00:00Z").
						CheckEqual("versions.v3.issued_date", "2023-01-03T00:00:00Z").
						CheckEqual("versions.v3.key_size_in_bytes", "2048").
						CheckEqual("versions.v3.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.v3.subject", "CN=test.example.com").
						CheckEqual("versions.v3.version_guid", "v3-guid-12345").
						CheckEqual("versions.v3.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v3.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v3.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
						CheckEqual("versions.v3.csr_block.key_algorithm", "RSA").
						CheckMissing("versions.v3.delete_requested_date").
						CheckMissing("versions.v3.scheduled_delete_date").
						Build(),
				},
			},
		},
		"happy path - remove all versions and add new ones": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v4": {
						version:                  4,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2023-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2023-01-01T00:00:00Z",
						deployedDate:             "2023-01-02T00:00:00Z",
						issuedDate:               "2023-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v4-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v3": {
						version:                  3,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2023-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2023-01-01T00:00:00Z",
						deployedDate:             "2023-01-02T00:00:00Z",
						issuedDate:               "2023-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v3-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update (delete all versions + add new versions)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v2")
				mockRotateClientCertificateVersion(t, m, newVersions, testCreateData.certificateID, "v3")
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v1")
				mockRotateClientCertificateVersion(t, m, newVersions, testCreateData.certificateID, "v4")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v3")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v4")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_remove_all_and_add_new_versions.tf"),
					Check: test.NewStateChecker("akamai_mtlskeystore_client_certificate_third_party.test").
						CheckEqual("certificate_id", "12345").
						CheckEqual("certificate_name", "test-certificate").
						CheckEqual("contract_id", "ctr_12345").
						CheckEqual("geography", "CORE").
						CheckEqual("group_id", "1234").
						CheckEqual("key_algorithm", "RSA").
						CheckEqual("notification_emails.#", "2").
						CheckEqual("notification_emails.0", "jkowalski@akamai.com").
						CheckEqual("notification_emails.1", "jsmith@akamai.com").
						CheckEqual("secure_network", "STANDARD_TLS").
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						CheckEqual("versions.v3.version", "3").
						CheckEqual("versions.v3.status", "ACTIVE").
						CheckEqual("versions.v3.expiry_date", "2025-12-31T23:59:59Z").
						CheckEqual("versions.v3.issuer", "Example CA").
						CheckEqual("versions.v3.key_algorithm", "RSA").
						CheckEqual("versions.v3.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.v3.certificate_submitted_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v3.created_by", "jkowalski").
						CheckEqual("versions.v3.created_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v3.deployed_date", "2023-01-02T00:00:00Z").
						CheckEqual("versions.v3.issued_date", "2023-01-03T00:00:00Z").
						CheckEqual("versions.v3.key_size_in_bytes", "2048").
						CheckEqual("versions.v3.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.v3.subject", "CN=test.example.com").
						CheckEqual("versions.v3.version_guid", "v3-guid-12345").
						CheckEqual("versions.v3.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v3.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v3.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
						CheckEqual("versions.v3.csr_block.key_algorithm", "RSA").
						CheckMissing("versions.v3.delete_requested_date").
						CheckMissing("versions.v3.scheduled_delete_date").
						CheckEqual("versions.v4.version", "4").
						CheckEqual("versions.v4.status", "ACTIVE").
						CheckEqual("versions.v4.expiry_date", "2025-12-31T23:59:59Z").
						CheckEqual("versions.v4.issuer", "Example CA").
						CheckEqual("versions.v4.key_algorithm", "RSA").
						CheckEqual("versions.v4.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.v4.certificate_submitted_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v4.created_by", "jkowalski").
						CheckEqual("versions.v4.created_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.v4.deployed_date", "2023-01-02T00:00:00Z").
						CheckEqual("versions.v4.issued_date", "2023-01-03T00:00:00Z").
						CheckEqual("versions.v4.key_size_in_bytes", "2048").
						CheckEqual("versions.v4.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.v4.subject", "CN=test.example.com").
						CheckEqual("versions.v4.version_guid", "v4-guid-12345").
						CheckEqual("versions.v4.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v4.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v4.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
						CheckEqual("versions.v4.csr_block.key_algorithm", "RSA").
						CheckMissing("versions.v4.delete_requested_date").
						CheckMissing("versions.v4.scheduled_delete_date").
						Build(),
				},
			},
		},
		"happy path - remove all 5 versions and add new 5 versions": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newFiveVersions := map[string]versionData{
					"v5a": {
						version:                  10,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2024-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2024-01-01T00:00:00Z",
						deployedDate:             "2024-01-02T00:00:00Z",
						issuedDate:               "2024-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v10-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v4a": {
						version:                  9,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2024-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2024-01-01T00:00:00Z",
						deployedDate:             "2024-01-02T00:00:00Z",
						issuedDate:               "2024-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v9-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v3a": {
						version:                  8,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2024-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2024-01-01T00:00:00Z",
						deployedDate:             "2024-01-02T00:00:00Z",
						issuedDate:               "2024-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v8-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v2a": {
						version:                  7,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2024-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2024-01-01T00:00:00Z",
						deployedDate:             "2024-01-02T00:00:00Z",
						issuedDate:               "2024-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v7-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v1a": {
						version:                  6,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2024-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2024-01-01T00:00:00Z",
						deployedDate:             "2024-01-02T00:00:00Z",
						issuedDate:               "2024-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v6-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				testCreateData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v3")
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v4")
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v5")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update (delete all versions + add new versions)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v2")
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v3")
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v4")
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v5")
				mockRotateClientCertificateVersion(t, m, newFiveVersions, testCreateData.certificateID, "v1a")
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v1")
				mockRotateClientCertificateVersion(t, m, newFiveVersions, testCreateData.certificateID, "v2a")
				mockRotateClientCertificateVersion(t, m, newFiveVersions, testCreateData.certificateID, "v3a")
				mockRotateClientCertificateVersion(t, m, newFiveVersions, testCreateData.certificateID, "v4a")
				mockRotateClientCertificateVersion(t, m, newFiveVersions, testCreateData.certificateID, "v5a")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newFiveVersions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newFiveVersions, testCreateData.certificateID)
				// Delete
				mockListClientCertificateVersions(t, m, newFiveVersions, testCreateData.certificateID)
				mockDeleteClientCertificateVersion(m, newFiveVersions, nil, testCreateData.certificateID, "v1a")
				mockDeleteClientCertificateVersion(m, newFiveVersions, nil, testCreateData.certificateID, "v2a")
				mockDeleteClientCertificateVersion(m, newFiveVersions, nil, testCreateData.certificateID, "v3a")
				mockDeleteClientCertificateVersion(m, newFiveVersions, nil, testCreateData.certificateID, "v4a")
				mockDeleteClientCertificateVersion(m, newFiveVersions, nil, testCreateData.certificateID, "v5a")
			},
			mockCreateData: testFiveVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_five_versions.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_remove_all_and_add_new_five_versions.tf"),
					Check: test.NewStateChecker("akamai_mtlskeystore_client_certificate_third_party.test").
						CheckEqual("certificate_id", "12345").
						CheckEqual("certificate_name", "test-certificate").
						CheckEqual("contract_id", "ctr_12345").
						CheckEqual("geography", "CORE").
						CheckEqual("group_id", "1234").
						CheckEqual("key_algorithm", "RSA").
						CheckEqual("notification_emails.#", "2").
						CheckEqual("notification_emails.0", "jkowalski@akamai.com").
						CheckEqual("notification_emails.1", "jsmith@akamai.com").
						CheckEqual("secure_network", "STANDARD_TLS").
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/").
						CheckEqual("versions.v1a.version", "6").
						CheckEqual("versions.v1a.status", "ACTIVE").
						CheckEqual("versions.v1a.expiry_date", "2025-12-31T23:59:59Z").
						CheckEqual("versions.v1a.issuer", "Example CA").
						CheckEqual("versions.v1a.key_algorithm", "RSA").
						CheckEqual("versions.v1a.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.v1a.certificate_submitted_date", "2024-01-01T00:00:00Z").
						CheckEqual("versions.v1a.created_by", "jkowalski").
						CheckEqual("versions.v1a.created_date", "2024-01-01T00:00:00Z").
						CheckEqual("versions.v1a.deployed_date", "2024-01-02T00:00:00Z").
						CheckEqual("versions.v1a.issued_date", "2024-01-03T00:00:00Z").
						CheckEqual("versions.v1a.key_size_in_bytes", "2048").
						CheckEqual("versions.v1a.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.v1a.subject", "CN=test.example.com").
						CheckEqual("versions.v1a.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v1a.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.v1a.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
						CheckEqual("versions.v1a.csr_block.key_algorithm", "RSA").
						CheckMissing("versions.v1a.delete_requested_date").
						CheckMissing("versions.v1a.scheduled_delete_date").
						CheckEqual("versions.v1a.version_guid", "v6-guid-12345").
						CheckEqual("versions.v2a.version", "7").
						CheckEqual("versions.v2a.version_guid", "v7-guid-12345").
						CheckEqual("versions.v3a.version", "8").
						CheckEqual("versions.v3a.version_guid", "v8-guid-12345").
						CheckEqual("versions.v4a.version", "9").
						CheckEqual("versions.v4a.version_guid", "v9-guid-12345").
						CheckEqual("versions.v5a.version", "10").
						CheckEqual("versions.v5a.version_guid", "v10-guid-12345").
						Build(),
				},
			},
		},
		"happy path - refresh (new version)": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v3": {
						version:                  3,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2023-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2024-01-01T00:00:00Z",
						deployedDate:             "2023-01-02T00:00:00Z",
						issuedDate:               "2023-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v3-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v2": testCreateData.versions["v2"],
					"v1": testCreateData.versions["v1"],
				}
				// Step 1
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Step 2
				// Read - mock that the new version was created outside terraform
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, newVersions, testCreateData.certificateID)
				// Delete previous versions to allow delete
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v1")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v2")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.certificateID, "v3")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						CheckEqual("versions.2024-01-01T00:00:00_v3.version", "3").
						CheckEqual("versions.2024-01-01T00:00:00_v3.status", "ACTIVE").
						CheckEqual("versions.2024-01-01T00:00:00_v3.expiry_date", "2025-12-31T23:59:59Z").
						CheckEqual("versions.2024-01-01T00:00:00_v3.issuer", "Example CA").
						CheckEqual("versions.2024-01-01T00:00:00_v3.key_algorithm", "RSA").
						CheckEqual("versions.2024-01-01T00:00:00_v3.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.2024-01-01T00:00:00_v3.certificate_submitted_date", "2023-01-01T00:00:00Z").
						CheckEqual("versions.2024-01-01T00:00:00_v3.created_by", "jkowalski").
						CheckEqual("versions.2024-01-01T00:00:00_v3.created_date", "2024-01-01T00:00:00Z").
						CheckEqual("versions.2024-01-01T00:00:00_v3.deployed_date", "2023-01-02T00:00:00Z").
						CheckEqual("versions.2024-01-01T00:00:00_v3.issued_date", "2023-01-03T00:00:00Z").
						CheckEqual("versions.2024-01-01T00:00:00_v3.key_size_in_bytes", "2048").
						CheckEqual("versions.2024-01-01T00:00:00_v3.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.2024-01-01T00:00:00_v3.subject", "CN=test.example.com").
						CheckEqual("versions.2024-01-01T00:00:00_v3.version_guid", "v3-guid-12345").
						CheckEqual("versions.2024-01-01T00:00:00_v3.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.2024-01-01T00:00:00_v3.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
						CheckEqual("versions.2024-01-01T00:00:00_v3.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
						CheckEqual("versions.2024-01-01T00:00:00_v3.csr_block.key_algorithm", "RSA").
						CheckMissing("versions.2024-01-01T00:00:00_v3.delete_requested_date").
						CheckMissing("versions.2024-01-01T00:00:00_v3.scheduled_delete_date").
						Build(),
				},
			},
		},
		"happy path - refresh (removed version)": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				versionsWithoutV2 := map[string]versionData{
					"v1": testCreateData.versions["v1"],
				}
				// Step 1
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Step 2
				// Read - mock that the version was removed outside terraform
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, versionsWithoutV2, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, versionsWithoutV2, testCreateData.certificateID)
				// Delete previous versions to allow delete
				mockListClientCertificateVersions(t, m, versionsWithoutV2, testCreateData.certificateID)
				mockDeleteClientCertificateVersion(m, versionsWithoutV2, nil, testCreateData.certificateID, "v1")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						CheckMissing("versions.v2.version").
						CheckMissing("versions.v2.status").
						CheckMissing("versions.v2.expiry_date").
						CheckMissing("versions.v2.issuer").
						CheckMissing("versions.v2.key_algorithm").
						CheckMissing("versions.v2.certificate_submitted_by").
						CheckMissing("versions.v2.certificate_submitted_date").
						CheckMissing("versions.v2.created_by").
						CheckMissing("versions.v2.created_date").
						CheckMissing("versions.v2.deployed_date").
						CheckMissing("versions.v2.issued_date").
						CheckMissing("versions.v2.key_elliptic_curve").
						CheckMissing("versions.v2.key_size_in_bytes").
						CheckMissing("versions.v2.signature_algorithm").
						CheckMissing("versions.v2.subject").
						CheckMissing("versions.v2.version_guid").
						CheckMissing("versions.v2.certificate_block.certificate").
						CheckMissing("versions.v2.certificate_block.trust_chain").
						CheckMissing("versions.v2.csr_block.csr").
						CheckMissing("versions.v2.csr_block.key_algorithm").
						CheckMissing("versions.v2.delete_requested_date").
						CheckMissing("versions.v2.scheduled_delete_date").
						Build(),
				},
			},
		},
		"happy path - refresh (update certificate name and notification emails)": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newCertificateNameAndEmails := commonDataForResource{
					certificateID:   12345,
					certificateName: "updated-certificate",
					contractID:      "ctr_12345",
					geography:       "CORE",
					groupID:         1234,
					keyAlgorithm:    "RSA",
					notificationEmails: []string{
						"jkowalski-updated@akamai.com",
						"jsmith-updated@akamai.com",
					},
					secureNetwork: "STANDARD_TLS",
					subject:       "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/",
					versions:      testCreateData.versions,
				}
				// Step 1
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Step 2
				// Read - mock that the version was removed outside terraform
				mockGetClientCertificate(m, newCertificateNameAndEmails)
				mockListClientCertificateVersions(t, m, newCertificateNameAndEmails.versions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, newCertificateNameAndEmails)
				mockListClientCertificateVersions(t, m, newCertificateNameAndEmails.versions, testCreateData.certificateID)
				// Delete previous versions to allow delete
				mockListClientCertificateVersions(t, m, newCertificateNameAndEmails.versions, testCreateData.certificateID)
				mockDeleteClientCertificateVersion(m, newCertificateNameAndEmails.versions, nil, testCreateData.certificateID, "v1")
				mockDeleteClientCertificateVersion(m, newCertificateNameAndEmails.versions, nil, testCreateData.certificateID, "v2")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						CheckEqual("certificate_name", "updated-certificate").
						CheckEqual("notification_emails.#", "2").
						CheckEqual("notification_emails.0", "jkowalski-updated@akamai.com").
						CheckEqual("notification_emails.1", "jsmith-updated@akamai.com").
						Build(),
				},
			},
		},
		"error create - API error": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				mockCreateClientCertificate(m, testCreateData).Return(nil, fmt.Errorf("create failed"))
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					ExpectError: regexp.MustCompile("create failed"),
				},
			},
		},
		"error update - API error (Rotation failed)": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v3": {
						version:                  3,
						status:                   "ACTIVE",
						expiryDate:               "2025-12-31T23:59:59Z",
						issuer:                   "Example CA",
						keyAlgorithm:             "RSA",
						certificateSubmittedBy:   "jkowalski",
						certificateSubmittedDate: "2023-01-01T00:00:00Z",
						createdBy:                "jkowalski",
						createdDate:              "2023-01-01T00:00:00Z",
						deployedDate:             "2023-01-02T00:00:00Z",
						issuedDate:               "2023-01-03T00:00:00Z",
						keySizeInBytes:           "2048",
						signatureAlgorithm:       "SHA256_WITH_RSA",
						subject:                  "CN=test.example.com",
						versionGUID:              "v3-guid-12345",
						certificateBlock: certificateBlock{
							certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							trustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						csrBlock: csrBlock{
							csr:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							keyAlgorithm: "RSA",
						},
					},
					"v1": testCreateData.versions["v1"],
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testCreateData.versions, testCreateData.certificateID, "v2")
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Twice()
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				// Update (delete version + add new version)
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID).Twice()
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v2")
				mockRotateClientCertificateVersion(t, m, newVersions, testCreateData.certificateID, "v3").Return(nil, fmt.Errorf("update failed"))
				// Delete - with old versions to allow deletion
				mockListClientCertificateVersions(t, m, testCreateData.versions, testCreateData.certificateID)
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v1")
				mockDeleteClientCertificateVersion(m, testCreateData.versions, nil, testCreateData.certificateID, "v2")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update/update_remove_one_and_add_one_version.tf"),
					ExpectError: regexp.MustCompile("update failed"),
				},
			},
		},
		"error - update contract": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Default subject is returned.
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read
				mockGetClientCertificate(m, testData).Twice()
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/error_update_contract.tf"),
					ExpectError: regexp.MustCompile("updating field `contract_id` is not possible"),
				},
			},
		},
		"error - update group": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Default subject is returned.
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read
				mockGetClientCertificate(m, testData).Twice()
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/error_update_group.tf"),
					ExpectError: regexp.MustCompile("updating field `group_id` is not possible"),
				},
			},
		},
		"error - update geography": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Default subject is returned.
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read
				mockGetClientCertificate(m, testData).Twice()
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/error_update_geography.tf"),
					ExpectError: regexp.MustCompile("updating field `geography` is not possible"),
				},
			},
		},
		"error - update secure_network": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Default subject is returned.
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read
				mockGetClientCertificate(m, testData).Twice()
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/error_update_secure_network.tf"),
					ExpectError: regexp.MustCompile("updating field `secure_network` is not possible"),
				},
			},
		},
		"error - update key_algorithm": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Rotate version
				mockRotateClientCertificateVersion(t, m, testData.versions, testData.certificateID, "v2")
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testData).Twice()
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v2")
			},
			mockCreateData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/error_update_key_algorithm.tf"),
					ExpectError: regexp.MustCompile("updating field `key_algorithm` is not possible"),
				},
			},
		},
		"error - update subject when it was not provided during create": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Default subject is returned.
				testData.subject = "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS 12345 1234/CN=test-certificate/"
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testData).Twice()
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/error_update_subject.tf"),
					ExpectError: regexp.MustCompile("Error: Cannot Update 'subject'"),
				},
			},
		},
		"error - update subject when it was provided during create": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				mockGetClientCertificate(m, testData)
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				// Read x2
				mockGetClientCertificate(m, testData).Twice()
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID).Twice()
				// Delete
				mockListClientCertificateVersions(t, m, testData.versions, testData.certificateID)
				mockDeleteClientCertificateVersion(m, testData.versions, nil, testData.certificateID, "v1")
			},
			mockCreateData: testOneVersionWithSubject,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create/create_with_subject.tf"),
					Check: baseChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/error_update_subject.tf"),
					ExpectError: regexp.MustCompile("Error: Cannot Update 'subject'"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.mockCreateData, tc.mockUpdateData)
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

func TestClientCertificateThirdPartyResource_Import(t *testing.T) {
	t.Parallel()
	importChecker := test.NewImportChecker().
		CheckEqual("certificate_id", "12345").
		CheckEqual("certificate_name", "test-certificate").
		CheckEqual("geography", "CORE").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("contract_id", "12345").
		CheckEqual("group_id", "1234").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "jkowalski@akamai.com").
		CheckEqual("notification_emails.1", "jsmith@akamai.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=17471 12345 1234/CN=test.example.com").
		CheckEqual("versions.2023-01-01T00:00:00_v1.version", "1").
		CheckEqual("versions.2023-01-01T00:00:00_v1.status", "ACTIVE").
		CheckEqual("versions.2023-01-01T00:00:00_v1.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.2023-01-01T00:00:00_v1.issuer", "Example CA").
		CheckEqual("versions.2023-01-01T00:00:00_v1.key_algorithm", "RSA").
		CheckEqual("versions.2023-01-01T00:00:00_v1.certificate_submitted_by", "jkowalski").
		CheckEqual("versions.2023-01-01T00:00:00_v1.certificate_submitted_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v1.created_by", "jkowalski").
		CheckEqual("versions.2023-01-01T00:00:00_v1.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v1.deployed_date", "2023-01-02T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v1.issued_date", "2023-01-03T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v1.key_size_in_bytes", "2048").
		CheckEqual("versions.2023-01-01T00:00:00_v1.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.2023-01-01T00:00:00_v1.subject", "CN=test.example.com").
		CheckEqual("versions.2023-01-01T00:00:00_v1.version_guid", "v1-guid-12345").
		CheckEqual("versions.2023-01-01T00:00:00_v1.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.2023-01-01T00:00:00_v1.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.2023-01-01T00:00:00_v1.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
		CheckEqual("versions.2023-01-01T00:00:00_v1.csr_block.key_algorithm", "RSA").
		CheckMissing("versions.2023-01-01T00:00:00_v1.delete_requested_date").
		CheckMissing("versions.2023-01-01T00:00:00_v1.scheduled_delete_date")

	secondVersionChecker := importChecker.
		CheckEqual("versions.2023-01-01T00:00:00_v2.version", "2").
		CheckEqual("versions.2023-01-01T00:00:00_v2.status", "ACTIVE").
		CheckEqual("versions.2023-01-01T00:00:00_v2.expiry_date", "2024-12-31T23:59:59Z").
		CheckEqual("versions.2023-01-01T00:00:00_v2.issuer", "Example CA").
		CheckEqual("versions.2023-01-01T00:00:00_v2.key_algorithm", "ECDSA").
		CheckEqual("versions.2023-01-01T00:00:00_v2.certificate_submitted_by", "jkowalski").
		CheckEqual("versions.2023-01-01T00:00:00_v2.certificate_submitted_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v2.created_by", "jkowalski").
		CheckEqual("versions.2023-01-01T00:00:00_v2.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v2.deployed_date", "2023-01-02T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v2.issued_date", "2023-01-03T00:00:00Z").
		CheckEqual("versions.2023-01-01T00:00:00_v2.key_size_in_bytes", "2048").
		CheckEqual("versions.2023-01-01T00:00:00_v2.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.2023-01-01T00:00:00_v2.subject", "CN=test.example.com").
		CheckEqual("versions.2023-01-01T00:00:00_v2.version_guid", "v2-guid-12345").
		CheckEqual("versions.2023-01-01T00:00:00_v2.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.2023-01-01T00:00:00_v2.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----").
		CheckEqual("versions.2023-01-01T00:00:00_v2.csr_block.csr", "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----").
		CheckEqual("versions.2023-01-01T00:00:00_v2.csr_block.key_algorithm", "RSA").
		CheckMissing("versions.2023-01-01T00:00:00_v2.delete_requested_date").
		CheckMissing("versions.2023-01-01T00:00:00_v2.scheduled_delete_date")

	tests := map[string]struct {
		importID    string
		init        func(*mtlskeystore.Mock, commonDataForResource)
		importData  commonDataForResource
		expectError *regexp.Regexp
		steps       []resource.TestStep
	}{
		"happy path - import with one version": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testImportData commonDataForResource) {
				// Import
				// modify mock data to return default subject that allows to parse contract and group
				testImportData.subject = "/C=US/O=Akamai Technologies, Inc./OU=17471 12345 1234/CN=test.example.com"
				mockGetClientCertificate(m, testImportData)
				// Read
				mockGetClientCertificate(m, testImportData)
				mockListClientCertificateVersions(t, m, testImportData.versions, testImportData.certificateID)
			},
			importData: testOneVersion,
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.Build(),
					ImportStateId:    "12345",
					ImportState:      true,
					ResourceName:     "akamai_mtlskeystore_client_certificate_third_party.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"happy path - import with two versions": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testImportData commonDataForResource) {
				// Import
				// modify mock data to return default subject that allows to parse contract and group.
				testImportData.subject = "/C=US/O=Akamai Technologies, Inc./OU=17471 12345 1234/CN=test.example.com"
				mockGetClientCertificate(m, testImportData)
				// Read
				mockGetClientCertificate(m, testImportData)
				mockListClientCertificateVersions(t, m, testImportData.versions, testImportData.certificateID)
			},
			importData: testTwoVersionsWithOptionalParams,
			steps: []resource.TestStep{
				{
					ImportStateCheck: secondVersionChecker.Build(),
					ImportStateId:    "12345",
					ImportState:      true,
					ResourceName:     "akamai_mtlskeystore_client_certificate_third_party.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"import - no group and contract in certificate subject, but provided with importID": {
			importID: "12345,123456,G-12RS3N4",
			init: func(m *mtlskeystore.Mock, testImportData commonDataForResource) {
				// Import
				mockGetClientCertificate(m, testImportData)
				// Read
				mockGetClientCertificate(m, testImportData)
				mockListClientCertificateVersions(t, m, testImportData.versions, testImportData.certificateID)
			},
			importData: testClientCertificateMissedContractAndGroupInSubjectTP,
			steps: []resource.TestStep{
				{
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("certificate_id", "12345").
						CheckEqual("group_id", "123456").
						CheckEqual("contract_id", "G-12RS3N4").Build(),
					ImportStateId: "12345,123456,G-12RS3N4",
					ImportState:   true,
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"error - wrong import ID": {
			importID:    "wrong-id",
			expectError: regexp.MustCompile(`failed to parse certificate ID as an integer: wrong-id`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "wrong-id",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`failed to parse certificate ID as an integer: wrong-id`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"error - custom subject cannot be parsed": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testImportData commonDataForResource) {
				// Import
				// modify mock data to return custom subject that cannot be parsed.
				testImportData.subject = "some custom subject cannot parse contract and group/CN=test.example.com"
				mockGetClientCertificate(m, testImportData)
			},
			importData:  testOneVersion,
			expectError: regexp.MustCompile(`get failed`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`since it is not possible to extract contract and group from certificate\nsubject, you need to provide an importID in the format\n'certificateID,groupID,contractID'. Where certificate, groupID and contractID\nare required`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"error - Get Client Certificate failed": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testImportData commonDataForResource) {
				// Import
				mockGetClientCertificate(m, testImportData).Return(nil, fmt.Errorf("get failed"))
			},
			importData:  testOneVersion,
			expectError: regexp.MustCompile(`get failed`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`get failed`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"error - List Client Certificate Versions failed": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testImportData commonDataForResource) {
				// Import
				// modify mock data to return default subject that allows to parse contract and group
				testImportData.subject = "/C=US/O=Akamai Technologies, Inc./OU=17471 12345 1234/CN=test.example.com"
				mockGetClientCertificate(m, testImportData)
				// Read
				mockGetClientCertificate(m, testImportData)
				mockListClientCertificateVersions(t, m, testImportData.versions, testImportData.certificateID).Return(nil, fmt.Errorf("list failed"))
			},
			importData:  testOneVersion,
			expectError: regexp.MustCompile(`list failed`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`list failed`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"error - no group and contract in certificate subject": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testImportData commonDataForResource) {
				mockGetClientCertificate(m, testImportData)
			},
			importData:  testClientCertificateMissedContractAndGroupInSubjectTP,
			expectError: regexp.MustCompile(`list failed`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`since it is not possible to extract contract and group from certificate\nsubject, you need to provide an importID in the format\n'certificateID,groupID,contractID'. Where certificate, groupID and contractID\nare required`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
		"error - incorrect number of parts in importID": {
			importID:    "123456789,123456,G-12RS3N4,123",
			importData:  testClientCertificateMissedContractAndGroupInSubjectTP,
			expectError: regexp.MustCompile(`list failed`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "123456789,123456,G-12RS3N4,123",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`you need to provide an importID in the format\n'certificateID,\[groupID,contractID]'. Where certificateID is required and\ngroupID and contractID are optional`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import/import.tf"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.importData)
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

func TestClientCertificateThirdPartyResource_ValidationErrors(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"error - missing certificate_name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/missing_certificate_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "certificate_name" is required, but no definition was found`),
				},
			},
		},
		"error - missing contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/missing_contract_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
			},
		},
		"error - missing geography": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/missing_geography.tf"),
					ExpectError: regexp.MustCompile(`The argument "geography" is required, but no definition was found`),
				},
			},
		},
		"error - missing group_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/missing_group_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
				},
			},
		},
		"error - missing notification_emails": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/missing_notification_emails.tf"),
					ExpectError: regexp.MustCompile(`Attribute notification_emails list must contain at least 1 elements, got: 0`),
				},
			},
		},
		"error - missing secure_network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/missing_secure_network.tf"),
					ExpectError: regexp.MustCompile(`The argument "secure_network" is required, but no definition was found`),
				},
			},
		},
		"error - missing version": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/missing_version.tf"),
					ExpectError: regexp.MustCompile(`Attribute versions map must contain at least 1 elements and at most 5\nelements, got: 0`),
				},
			},
		},
		"error - invalid key algorithm": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/invalid_key_algorithm.tf"),
					ExpectError: regexp.MustCompile(`Attribute key_algorithm value must be one of: \["RSA" "ECDSA"], got: "INVALID"`),
				},
			},
		},
		"error - invalid secure network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/invalid_secure_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute secure_network value must be one of: \["STANDARD_TLS"\n"ENHANCED_TLS"], got: "INVALID"`),
				},
			},
		},
		"error - invalid geography": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/invalid_geography.tf"),
					ExpectError: regexp.MustCompile(`Attribute geography value must be one of: \["CORE" "RUSSIA_AND_CORE"\n"CHINA_AND_CORE"], got: "INVALID"`),
				},
			},
		},
		"error - invalid subject": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/invalid_subject.tf"),
					ExpectError: regexp.MustCompile(
						"Attribute subject The `subject` must contain a valid `CN` " +
							"attribute with a\nmaximum length of 64 characters., got: /C=US/O=Akamai Technologies,\nInc./OU=Akamai\nmTLS/" +
							"CN=test-certificateAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA/"),
				},
			},
		},
		"error - more than 5 versions provided": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/error/more_than_5_versions.tf"),
					ExpectError: regexp.MustCompile(`Attribute versions map must contain at least 1 elements and at most 5\nelements, got: 6`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}

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

func mockCreateClientCertificate(m *mtlskeystore.Mock, testData commonDataForResource) *mock.Call {
	if testData.keyAlgorithm == "" {
		testData.keyAlgorithm = "RSA"
	}
	keyAlgorithm := mtlskeystore.CryptographicAlgorithm(testData.keyAlgorithm)

	request := mtlskeystore.CreateClientCertificateRequest{
		CertificateName:    testData.certificateName,
		ContractID:         strings.TrimPrefix(testData.contractID, "ctr_"),
		Geography:          mtlskeystore.Geography(testData.geography),
		GroupID:            testData.groupID,
		KeyAlgorithm:       &keyAlgorithm,
		NotificationEmails: testData.notificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(testData.secureNetwork),
		Subject:            ptr.To(testData.subject),
		Signer:             mtlskeystore.SignerThirdParty,
		PreferredCA:        testData.preferredCA,
	}

	return m.On("CreateClientCertificate", testutils.MockContext, request).Return(&mtlskeystore.CreateClientCertificateResponse{
		CertificateID:      testData.certificateID,
		CertificateName:    testData.certificateName,
		Geography:          testData.geography,
		KeyAlgorithm:       string(keyAlgorithm),
		NotificationEmails: testData.notificationEmails,
		SecureNetwork:      testData.secureNetwork,
		Signer:             string(mtlskeystore.SignerThirdParty),
		Subject:            testData.subject,
	}, nil).Once()
}

func mockRotateClientCertificateVersion(t *testing.T, m *mtlskeystore.Mock, testData map[string]versionData, certificateID int64, versionKey string) *mock.Call {
	response := mtlskeystore.RotateClientCertificateVersionResponse{
		Version:            testData[versionKey].version,
		VersionGUID:        testData[versionKey].versionGUID,
		CreatedBy:          testData[versionKey].createdBy,
		CreatedDate:        tst.NewTimeFromString(t, testData[versionKey].createdDate),
		ExpiryDate:         ptr.To(tst.NewTimeFromString(t, testData[versionKey].expiryDate)),
		IssuedDate:         ptr.To(tst.NewTimeFromString(t, testData[versionKey].issuedDate)),
		Issuer:             ptr.To(testData[versionKey].issuer),
		KeyAlgorithm:       testData[versionKey].keyAlgorithm,
		EllipticCurve:      ptr.To(testData[versionKey].keyEllipticCurve),
		KeySizeInBytes:     ptr.To(testData[versionKey].keySizeInBytes),
		SignatureAlgorithm: ptr.To(testData[versionKey].signatureAlgorithm),
		Status:             testData[versionKey].status,
		Subject:            ptr.To(testData[versionKey].subject),
	}
	if testData[versionKey].certificateBlock != (certificateBlock{}) {
		response.CertificateBlock = &mtlskeystore.CertificateBlock{
			Certificate: testData[versionKey].certificateBlock.certificate,
			TrustChain:  testData[versionKey].certificateBlock.trustChain,
		}
	}
	if testData[versionKey].csrBlock != (csrBlock{}) {
		response.CSRBlock = &mtlskeystore.CSRBlock{
			CSR:          testData[versionKey].csrBlock.csr,
			KeyAlgorithm: testData[versionKey].csrBlock.keyAlgorithm,
		}
	}
	if testData[versionKey].certificateSubmittedBy != "" {
		response.CertificateSubmittedBy = ptr.To(testData[versionKey].certificateSubmittedBy)
	}
	if testData[versionKey].certificateSubmittedDate != "" {
		response.CertificateSubmittedDate = ptr.To(tst.NewTimeFromString(t, testData[versionKey].certificateSubmittedDate))
	}
	if testData[versionKey].deleteRequestedDate != "" {
		response.DeleteRequestedDate = ptr.To(tst.NewTimeFromString(t, testData[versionKey].deleteRequestedDate))
	}
	if testData[versionKey].deployedDate != "" {
		response.DeployedDate = ptr.To(tst.NewTimeFromString(t, testData[versionKey].deployedDate))
	}
	if testData[versionKey].scheduledDeleteDate != "" {
		response.ScheduledDeleteDate = ptr.To(tst.NewTimeFromString(t, testData[versionKey].scheduledDeleteDate))
	}

	return m.On("RotateClientCertificateVersion", testutils.MockContext, mtlskeystore.RotateClientCertificateVersionRequest{
		CertificateID: certificateID,
	}).Return(&response, nil).Once()
}

func mockPatchClientCertificate(m *mtlskeystore.Mock, certID int64, certName *string, emails []string) *mock.Call {
	return m.On("PatchClientCertificate", testutils.MockContext, mtlskeystore.PatchClientCertificateRequest{
		CertificateID: certID,
		Body: mtlskeystore.PatchClientCertificateRequestBody{
			CertificateName:    certName,
			NotificationEmails: emails,
		},
	}).Return(nil).Once()
}

func mockDeleteClientCertificateVersion(m *mtlskeystore.Mock, versions map[string]versionData, resp *mtlskeystore.DeleteClientCertificateVersionResponse, certificateID int64, versionKey string) *mock.Call {
	return m.On("DeleteClientCertificateVersion", testutils.MockContext, mtlskeystore.DeleteClientCertificateVersionRequest{
		CertificateID: certificateID,
		Version:       versions[versionKey].version,
	}).Return(resp, nil).Once()
}

func mockGetClientCertificate(m *mtlskeystore.Mock, testData commonDataForResource) *mock.Call {
	return m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
		CertificateID: testData.certificateID,
	}).Return(&mtlskeystore.GetClientCertificateResponse{
		CertificateID:      testData.certificateID,
		CertificateName:    testData.certificateName,
		Geography:          testData.geography,
		KeyAlgorithm:       testData.keyAlgorithm,
		NotificationEmails: testData.notificationEmails,
		SecureNetwork:      testData.secureNetwork,
		Subject:            testData.subject,
	}, nil).Once()
}

func mockListClientCertificateVersions(t *testing.T, m *mtlskeystore.Mock, versions map[string]versionData, certificateID int64) *mock.Call {
	responseVersions := make([]mtlskeystore.ClientCertificateVersion, 0, len(versions))

	for _, version := range versions {
		certificateVersions := mtlskeystore.ClientCertificateVersion{
			Version:            version.version,
			VersionGUID:        version.versionGUID,
			CreatedBy:          version.createdBy,
			CreatedDate:        tst.NewTimeFromString(t, version.createdDate),
			ExpiryDate:         ptr.To(tst.NewTimeFromString(t, version.expiryDate)),
			IssuedDate:         ptr.To(tst.NewTimeFromString(t, version.issuedDate)),
			Issuer:             ptr.To(version.issuer),
			KeyAlgorithm:       version.keyAlgorithm,
			EllipticCurve:      ptr.To(version.keyEllipticCurve),
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
		if version.csrBlock != (csrBlock{}) {
			certificateVersions.CSRBlock = &mtlskeystore.CSRBlock{
				CSR:          version.csrBlock.csr,
				KeyAlgorithm: version.csrBlock.keyAlgorithm,
			}
		}
		if version.certificateSubmittedBy != "" {
			certificateVersions.CertificateSubmittedBy = ptr.To(version.certificateSubmittedBy)
		}
		if version.certificateSubmittedDate != "" {
			certificateVersions.CertificateSubmittedDate = ptr.To(tst.NewTimeFromString(t, version.certificateSubmittedDate))
		}
		if version.deleteRequestedDate != "" {
			certificateVersions.DeleteRequestedDate = ptr.To(tst.NewTimeFromString(t, version.deleteRequestedDate))
		}
		if version.deployedDate != "" {
			certificateVersions.DeployedDate = ptr.To(tst.NewTimeFromString(t, version.deployedDate))
		}
		if version.scheduledDeleteDate != "" {
			certificateVersions.ScheduledDeleteDate = ptr.To(tst.NewTimeFromString(t, version.scheduledDeleteDate))
		}

		responseVersions = append(responseVersions, certificateVersions)
	}

	sort.Slice(responseVersions, func(i, j int) bool {
		return responseVersions[i].Version > responseVersions[j].Version
	})

	return m.On("ListClientCertificateVersions", testutils.MockContext, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               certificateID,
		IncludeAssociatedProperties: true,
	}).Return(&mtlskeystore.ListClientCertificateVersionsResponse{
		Versions: responseVersions,
	}, nil).Once()
}
