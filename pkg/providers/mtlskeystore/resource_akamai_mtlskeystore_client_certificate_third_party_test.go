package mtlskeystore

import (
	"fmt"
	"regexp"
	"sort"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	commonDataForResource struct {
		CertificateName    string
		CertificateID      int64
		ContractID         string
		Geography          string
		GroupID            int64
		KeyAlgorithm       string
		NotificationEmails []string
		SecureNetwork      string
		Subject            string
		Versions           map[string]versionData
	}

	versionData struct {
		Version                  int64
		Status                   string
		ExpiryDate               string
		Issuer                   string
		KeyAlgorithm             string
		CertificateSubmittedBy   string
		CertificateSubmittedDate string
		CreatedBy                string
		CreatedDate              string
		DeleteRequestedDate      string
		DeployedDate             string
		IssuedDate               string
		KeyEllipticCurve         string
		KeySizeInBytes           string
		ScheduledDeleteDate      string
		SignatureAlgorithm       string
		Subject                  string
		VersionGUID              string
		CertificateBlock         certificateBlock
		CSRBlock                 csrBlock
		Validation               validation
		AssociatedProperties     []associatedProperty
	}

	certificateBlock struct {
		Certificate string
		TrustChain  string
	}

	csrBlock struct {
		CSR          string
		KeyAlgorithm string
	}

	validation struct {
		Errors   []validationError
		Warnings []validationError
	}

	validationError struct {
		Message string
		Reason  string
		Type    string
	}

	associatedProperty struct {
		AssetID         int64
		GroupID         int64
		PropertyVersion int64
		PropertyName    string
	}
)

var (
	testTwoVersions = commonDataForResource{
		CertificateName: "test-certificate",
		CertificateID:   12345,
		ContractID:      "ctr_12345",
		Geography:       "CORE",
		GroupID:         1234,
		KeyAlgorithm:    "RSA",
		NotificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		SecureNetwork: "STANDARD_TLS",
		Subject:       "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/",
		Versions: map[string]versionData{
			"v2": {
				Version:                  2,
				Status:                   "ACTIVE",
				ExpiryDate:               "2024-12-31T23:59:59Z",
				Issuer:                   "Example CA",
				KeyAlgorithm:             "ECDSA",
				CertificateSubmittedBy:   "jkowalski",
				CertificateSubmittedDate: "2023-01-01T00:00:00Z",
				CreatedBy:                "jkowalski",
				CreatedDate:              "2023-01-01T00:00:00Z",
				DeployedDate:             "2023-01-02T00:00:00Z",
				IssuedDate:               "2023-01-03T00:00:00Z",
				KeyEllipticCurve:         "test-ecdsa",
				KeySizeInBytes:           "2048",
				SignatureAlgorithm:       "SHA256_WITH_RSA",
				Subject:                  "CN=test.example.com",
				VersionGUID:              "v2-guid-12345",
				CertificateBlock: certificateBlock{
					Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				CSRBlock: csrBlock{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
			},
			"v1": {
				Version:                  1,
				Status:                   "ACTIVE",
				ExpiryDate:               "2024-12-31T23:59:59Z",
				Issuer:                   "Example CA",
				KeyAlgorithm:             "RSA",
				CertificateSubmittedBy:   "jkowalski",
				CertificateSubmittedDate: "2023-01-01T00:00:00Z",
				CreatedBy:                "jkowalski",
				CreatedDate:              "2023-01-01T00:00:00Z",
				DeployedDate:             "2023-01-02T00:00:00Z",
				IssuedDate:               "2023-01-03T00:00:00Z",
				KeySizeInBytes:           "2048",
				SignatureAlgorithm:       "SHA256_WITH_RSA",
				Subject:                  "CN=test.example.com",
				VersionGUID:              "v1-guid-12345",
				CertificateBlock: certificateBlock{
					Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				CSRBlock: csrBlock{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
			},
		},
	}

	testOneVersion = commonDataForResource{
		CertificateName: "test-certificate",
		CertificateID:   12345,
		ContractID:      "ctr_12345",
		Geography:       "CORE",
		GroupID:         1234,
		KeyAlgorithm:    "RSA",
		NotificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		SecureNetwork: "STANDARD_TLS",
		Versions: map[string]versionData{
			"v1": {
				Version:                  1,
				Status:                   "ACTIVE",
				ExpiryDate:               "2024-12-31T23:59:59Z",
				Issuer:                   "Example CA",
				KeyAlgorithm:             "RSA",
				CertificateSubmittedBy:   "jkowalski",
				CertificateSubmittedDate: "2023-01-01T00:00:00Z",
				CreatedBy:                "jkowalski",
				CreatedDate:              "2023-01-01T00:00:00Z",
				DeployedDate:             "2023-01-02T00:00:00Z",
				IssuedDate:               "2023-01-03T00:00:00Z",
				KeySizeInBytes:           "2048",
				SignatureAlgorithm:       "SHA256_WITH_RSA",
				Subject:                  "CN=test.example.com",
				VersionGUID:              "v1-guid-12345",
				CertificateBlock: certificateBlock{
					Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				CSRBlock: csrBlock{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
				Validation: validation{
					Errors:   nil,
					Warnings: nil,
				},
				AssociatedProperties: []associatedProperty{},
			},
		},
	}

	testUpdateNotificationEmailsAndCertificateName = commonDataForResource{
		CertificateName: "updated-certificate-name",
		CertificateID:   12345,
		ContractID:      "ctr_12345",
		Geography:       "CORE",
		GroupID:         1234,
		KeyAlgorithm:    "RSA",
		NotificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith-new@akamai.com",
			"test@akamai.com",
		},
		SecureNetwork: "STANDARD_TLS",
		Versions: map[string]versionData{
			"v1": {
				Version:                  1,
				Status:                   "ACTIVE",
				ExpiryDate:               "2024-12-31T23:59:59Z",
				Issuer:                   "Example CA",
				KeyAlgorithm:             "RSA",
				CertificateSubmittedBy:   "jkowalski",
				CertificateSubmittedDate: "2023-01-01T00:00:00Z",
				CreatedBy:                "jkowalski",
				CreatedDate:              "2023-01-01T00:00:00Z",
				DeployedDate:             "2023-01-02T00:00:00Z",
				IssuedDate:               "2023-01-03T00:00:00Z",
				KeySizeInBytes:           "2048",
				SignatureAlgorithm:       "SHA256_WITH_RSA",
				Subject:                  "CN=test.example.com",
				VersionGUID:              "v1-guid-12345",
				CertificateBlock: certificateBlock{
					Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				CSRBlock: csrBlock{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
			},
		},
	}

	testDataVersionWithAssociatedProperties = commonDataForResource{
		CertificateName: "test-certificate",
		CertificateID:   12345,
		ContractID:      "ctr_12345",
		Geography:       "CORE",
		GroupID:         1234,
		KeyAlgorithm:    "RSA",
		NotificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		SecureNetwork: "STANDARD_TLS",
		Versions: map[string]versionData{
			"v1": {
				Version:                  1,
				Status:                   "ACTIVE",
				ExpiryDate:               "2024-12-31T23:59:59Z",
				Issuer:                   "Example CA",
				KeyAlgorithm:             "RSA",
				CertificateSubmittedBy:   "jkowalski",
				CertificateSubmittedDate: "2023-01-01T00:00:00Z",
				CreatedBy:                "jkowalski",
				CreatedDate:              "2023-01-01T00:00:00Z",
				DeployedDate:             "2023-01-02T00:00:00Z",
				IssuedDate:               "2023-01-03T00:00:00Z",
				KeySizeInBytes:           "2048",
				SignatureAlgorithm:       "SHA256_WITH_RSA",
				Subject:                  "CN=test.example.com",
				VersionGUID:              "v1-guid-12345",
				CertificateBlock: certificateBlock{
					Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
					TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
				},
				CSRBlock: csrBlock{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
				Validation: validation{
					Errors:   nil,
					Warnings: nil,
				},
				AssociatedProperties: []associatedProperty{
					{
						AssetID:         123456,
						GroupID:         1234,
						PropertyVersion: 1,
						PropertyName:    "test-property-1",
					},
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
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
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
				// Get
				mockGetClientCertificate(m, testData).Times(2)
				mockListClientCertificateVersions(m, testData.Versions, testData.CertificateID).Times(2)
				// Delete
				mockListClientCertificateVersions(m, testData.Versions, testData.CertificateID)
				mockDeleteClientCertificateVersion(m, testData.Versions, nil, testData.CertificateID, "v1")
			},
			mockCreateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},
		"happy path - with optional params and multiple versions": {
			init: func(m *mtlskeystore.Mock, testData, _ commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testData)
				// Rotate version
				mockRotateClientCertificateVersion(m, testData.Versions, testData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testData).Times(2)
				mockListClientCertificateVersions(m, testData.Versions, testData.CertificateID).Times(2)
				// Delete
				mockListClientCertificateVersions(m, testData.Versions, testData.CertificateID)
				mockDeleteClientCertificateVersion(m, testData.Versions, nil, testData.CertificateID, "v1")
				mockDeleteClientCertificateVersion(m, testData.Versions, nil, testData.CertificateID, "v2")
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
			},
		},
		"happy path - update certificate name and notification emails": {
			init: func(m *mtlskeystore.Mock, testCreateData, testUpdateData commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Patch
				mockPatchClientCertificate(m, testUpdateData)
				// Get
				mockGetClientCertificate(m, testUpdateData).Times(2)
				mockListClientCertificateVersions(m, testUpdateData.Versions, testUpdateData.CertificateID).Times(2)
				// Delete
				mockListClientCertificateVersions(m, testUpdateData.Versions, testUpdateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testUpdateData.Versions, nil, testUpdateData.CertificateID, "v1")
			},
			mockCreateData: testOneVersion,
			mockUpdateData: testUpdateNotificationEmailsAndCertificateName,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update_certificate_name_and_notification_emails.tf"),
					Check: baseChecker.
						CheckEqual("certificate_name", "updated-certificate-name").
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						CheckEqual("notification_emails.#", "3").
						CheckEqual("notification_emails.0", "jkowalski@akamai.com").
						CheckEqual("notification_emails.1", "jsmith-new@akamai.com").
						CheckEqual("notification_emails.2", "test@akamai.com").
						Build(),
				},
			},
		},
		"happy path - add new version": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v4": {
						Version:                  4,
						Status:                   "ACTIVE",
						ExpiryDate:               "2025-12-31T23:59:59Z",
						Issuer:                   "Example CA",
						KeyAlgorithm:             "RSA",
						CertificateSubmittedBy:   "jkowalski",
						CertificateSubmittedDate: "2023-01-01T00:00:00Z",
						CreatedBy:                "jkowalski",
						CreatedDate:              "2023-01-01T00:00:00Z",
						DeployedDate:             "2023-01-02T00:00:00Z",
						IssuedDate:               "2023-01-03T00:00:00Z",
						KeySizeInBytes:           "2048",
						SignatureAlgorithm:       "SHA256_WITH_RSA",
						Subject:                  "CN=test.example.com",
						VersionGUID:              "v4-guid-12345",
						CertificateBlock: certificateBlock{
							Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						CSRBlock: csrBlock{
							CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					"v3": {
						Version:                  3,
						Status:                   "ACTIVE",
						ExpiryDate:               "2025-12-31T23:59:59Z",
						Issuer:                   "Example CA",
						KeyAlgorithm:             "RSA",
						CertificateSubmittedBy:   "jkowalski",
						CertificateSubmittedDate: "2023-01-01T00:00:00Z",
						CreatedBy:                "jkowalski",
						CreatedDate:              "2023-01-01T00:00:00Z",
						DeployedDate:             "2023-01-02T00:00:00Z",
						IssuedDate:               "2023-01-03T00:00:00Z",
						KeySizeInBytes:           "2048",
						SignatureAlgorithm:       "SHA256_WITH_RSA",
						Subject:                  "CN=test.example.com",
						VersionGUID:              "v3-guid-12345",
						CertificateBlock: certificateBlock{
							Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						CSRBlock: csrBlock{
							CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					"v2": testCreateData.Versions["v2"],
					"v1": testCreateData.Versions["v1"],
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(m, testCreateData.Versions, testCreateData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Add new versions (Rotate)
				mockRotateClientCertificateVersion(m, newVersions, testCreateData.CertificateID, "v3")
				mockRotateClientCertificateVersion(m, newVersions, testCreateData.CertificateID, "v4")
				// Get
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID).Times(2)
				// Delete
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v1")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v2")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v3")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v4")
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update_new_version.tf"),
					Check: secondVersionChecker.
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
				mockRotateClientCertificateVersion(m, testCreateData.Versions, testCreateData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Update (delete version)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v2")
				// Read
				mockGetClientCertificate(m, testUpdateData).Times(2)
				mockListClientCertificateVersions(m, testUpdateData.Versions, testUpdateData.CertificateID).Times(2)
				// Delete
				mockListClientCertificateVersions(m, testUpdateData.Versions, testUpdateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testUpdateData.Versions, nil, testUpdateData.CertificateID, "v1")
			},
			mockCreateData: testTwoVersions,
			mockUpdateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update_remove_one_version.tf"),
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
						Version:                  3,
						Status:                   "ACTIVE",
						ExpiryDate:               "2025-12-31T23:59:59Z",
						Issuer:                   "Example CA",
						KeyAlgorithm:             "RSA",
						CertificateSubmittedBy:   "jkowalski",
						CertificateSubmittedDate: "2023-01-01T00:00:00Z",
						CreatedBy:                "jkowalski",
						CreatedDate:              "2023-01-01T00:00:00Z",
						DeployedDate:             "2023-01-02T00:00:00Z",
						IssuedDate:               "2023-01-03T00:00:00Z",
						KeySizeInBytes:           "2048",
						SignatureAlgorithm:       "SHA256_WITH_RSA",
						Subject:                  "CN=test.example.com",
						VersionGUID:              "v3-guid-12345",
						CertificateBlock: certificateBlock{
							Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						CSRBlock: csrBlock{
							CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					"v1": testCreateData.Versions["v1"],
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(m, testCreateData.Versions, testCreateData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Update (delete version + add new version)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v2")
				mockRotateClientCertificateVersion(m, newVersions, testCreateData.CertificateID, "v3")
				// Get
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID).Times(2)
				// Delete
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v1")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v3")
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update_remove_one_and_add_one_version.tf"),
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
						Version:                  4,
						Status:                   "ACTIVE",
						ExpiryDate:               "2025-12-31T23:59:59Z",
						Issuer:                   "Example CA",
						KeyAlgorithm:             "RSA",
						CertificateSubmittedBy:   "jkowalski",
						CertificateSubmittedDate: "2023-01-01T00:00:00Z",
						CreatedBy:                "jkowalski",
						CreatedDate:              "2023-01-01T00:00:00Z",
						DeployedDate:             "2023-01-02T00:00:00Z",
						IssuedDate:               "2023-01-03T00:00:00Z",
						KeySizeInBytes:           "2048",
						SignatureAlgorithm:       "SHA256_WITH_RSA",
						Subject:                  "CN=test.example.com",
						VersionGUID:              "v4-guid-12345",
						CertificateBlock: certificateBlock{
							Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						CSRBlock: csrBlock{
							CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					"v3": {
						Version:                  3,
						Status:                   "ACTIVE",
						ExpiryDate:               "2025-12-31T23:59:59Z",
						Issuer:                   "Example CA",
						KeyAlgorithm:             "RSA",
						CertificateSubmittedBy:   "jkowalski",
						CertificateSubmittedDate: "2023-01-01T00:00:00Z",
						CreatedBy:                "jkowalski",
						CreatedDate:              "2023-01-01T00:00:00Z",
						DeployedDate:             "2023-01-02T00:00:00Z",
						IssuedDate:               "2023-01-03T00:00:00Z",
						KeySizeInBytes:           "2048",
						SignatureAlgorithm:       "SHA256_WITH_RSA",
						Subject:                  "CN=test.example.com",
						VersionGUID:              "v3-guid-12345",
						CertificateBlock: certificateBlock{
							Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						CSRBlock: csrBlock{
							CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(m, testCreateData.Versions, testCreateData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Update (delete all versions + add new versions)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v2")
				mockRotateClientCertificateVersion(m, newVersions, testCreateData.CertificateID, "v3")
				mockRotateClientCertificateVersion(m, newVersions, testCreateData.CertificateID, "v4")
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v1")
				// Get
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID).Times(2)
				// Delete
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v3")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v4")
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update_remove_all_and_add_new_versions.tf"),
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
		"happy path - refresh (new version)": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v3": {
						Version:                  3,
						Status:                   "ACTIVE",
						ExpiryDate:               "2025-12-31T23:59:59Z",
						Issuer:                   "Example CA",
						KeyAlgorithm:             "RSA",
						CertificateSubmittedBy:   "jkowalski",
						CertificateSubmittedDate: "2023-01-01T00:00:00Z",
						CreatedBy:                "jkowalski",
						CreatedDate:              "2024-01-01T00:00:00Z",
						DeployedDate:             "2023-01-02T00:00:00Z",
						IssuedDate:               "2023-01-03T00:00:00Z",
						KeySizeInBytes:           "2048",
						SignatureAlgorithm:       "SHA256_WITH_RSA",
						Subject:                  "CN=test.example.com",
						VersionGUID:              "v3-guid-12345",
						CertificateBlock: certificateBlock{
							Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						CSRBlock: csrBlock{
							CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					"v2": testCreateData.Versions["v2"],
					"v1": testCreateData.Versions["v1"],
				}
				// Step 1
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(m, testCreateData.Versions, testCreateData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Step 2
				// Read - mock that the new version was created outside terraform
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID)
				// Delete previous versions to allow delete
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v1")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v2")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v3")
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
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
		"error create - API error": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				mockCreateClientCertificate(m, testCreateData).Return(nil, fmt.Errorf("create failed"))
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					ExpectError: regexp.MustCompile("create failed"),
				},
			},
		},
		"error update - API error (Rotation failed)": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v3": {
						Version:                  3,
						Status:                   "ACTIVE",
						ExpiryDate:               "2025-12-31T23:59:59Z",
						Issuer:                   "Example CA",
						KeyAlgorithm:             "RSA",
						CertificateSubmittedBy:   "jkowalski",
						CertificateSubmittedDate: "2023-01-01T00:00:00Z",
						CreatedBy:                "jkowalski",
						CreatedDate:              "2023-01-01T00:00:00Z",
						DeployedDate:             "2023-01-02T00:00:00Z",
						IssuedDate:               "2023-01-03T00:00:00Z",
						KeySizeInBytes:           "2048",
						SignatureAlgorithm:       "SHA256_WITH_RSA",
						Subject:                  "CN=test.example.com",
						VersionGUID:              "v3-guid-12345",
						CertificateBlock: certificateBlock{
							Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
							TrustChain:  "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
						},
						CSRBlock: csrBlock{
							CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					"v1": testCreateData.Versions["v1"],
				}
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(m, testCreateData.Versions, testCreateData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Update (delete version + add new version)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v2")
				mockRotateClientCertificateVersion(m, newVersions, testCreateData.CertificateID, "v3").Return(nil, fmt.Errorf("update failed"))
				// Delete - with old versions to allow deletion
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v1")
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v2")
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update_remove_one_and_add_one_version.tf"),
					ExpectError: regexp.MustCompile("update failed"),
				},
			},
		},
		"error update - version status is DELETE_PENDING": {
			init: func(m *mtlskeystore.Mock, testCreateData, _ commonDataForResource) {
				newVersions := map[string]versionData{
					"v2": testCreateData.Versions["v2"],
					"v1": testCreateData.Versions["v1"],
				}
				v2 := newVersions["v2"]
				v2.Status = "DELETE_PENDING"
				v2.DeleteRequestedDate = "2024-01-01T00:00:00Z"
				v2.ScheduledDeleteDate = "2024-01-02T00:00:00Z"
				newVersions["v2"] = v2
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Rotate version
				mockRotateClientCertificateVersion(m, testCreateData.Versions, testCreateData.CertificateID, "v2")
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Update (delete version)
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID)
				// Delete - set to `ACTIVE` to allow deletion
				v2.Status = "ACTIVE"
				newVersions["v2"] = v2
				mockListClientCertificateVersions(m, newVersions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v1")
				mockDeleteClientCertificateVersion(m, newVersions, nil, testCreateData.CertificateID, "v2")
			},
			mockCreateData: testTwoVersions,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create_with_optional_params.tf"),
					Check: secondVersionChecker.
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/update_remove_one_version.tf"),
					ExpectError: regexp.MustCompile("cannot delete client certificate version with status DELETE_PENDING"),
				},
			},
		},
		"error delete - version has Associated Properties": {
			init: func(m *mtlskeystore.Mock, testCreateData, testAllowDelete commonDataForResource) {
				// Create
				mockCreateClientCertificate(m, testCreateData)
				// Get
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				// Read
				mockGetClientCertificate(m, testCreateData).Times(2)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Times(2)
				// Delete version with properties
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
				mockDeleteClientCertificateVersion(m, testCreateData.Versions, nil, testCreateData.CertificateID, "v1")
				// Delete - with old versions to allow deletion
				mockListClientCertificateVersions(m, testAllowDelete.Versions, testCreateData.CertificateID)
			},
			mockCreateData: testDataVersionWithAssociatedProperties,
			mockUpdateData: testOneVersion,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/create.tf"),
					Destroy:     true,
					ExpectError: regexp.MustCompile("cannot delete client certificate version 1 with associated properties"),
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
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "jkowalski@akamai.com").
		CheckEqual("notification_emails.1", "jsmith@akamai.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/").
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
			init: func(m *mtlskeystore.Mock, testCreateData commonDataForResource) {
				// Import
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
			},
			importData: testOneVersion,
			steps: []resource.TestStep{
				{
					ImportStateCheck: importChecker.Build(),
					ImportStateId:    "12345",
					ImportState:      true,
					ResourceName:     "akamai_mtlskeystore_client_certificate_third_party.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import_one_version.tf"),
				},
			},
		},
		"happy path - import with two versions": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testCreateData commonDataForResource) {
				// Import
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID)
			},
			importData: testTwoVersions,
			steps: []resource.TestStep{
				{
					ImportStateCheck: secondVersionChecker.Build(),
					ImportStateId:    "12345",
					ImportState:      true,
					ResourceName:     "akamai_mtlskeystore_client_certificate_third_party.test",
					Config:           testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import_two_versions.tf"),
				},
			},
		},
		"error - wrong import ID": {
			importID:    "wrong-id",
			expectError: regexp.MustCompile(`Error: could not convert import ID to int`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "wrong-id",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`Error: could not convert import ID to int`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import_one_version.tf"),
				},
			},
		},
		"error - Get Client Certificate failed": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testCreateData commonDataForResource) {
				// Import
				mockGetClientCertificate(m, testCreateData).Return(nil, fmt.Errorf("get failed"))
			},
			importData:  testOneVersion,
			expectError: regexp.MustCompile(`get failed`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`get failed`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import_one_version.tf"),
				},
			},
		},
		"error - List Client Certificate Versions failed": {
			importID: "12345",
			init: func(m *mtlskeystore.Mock, testCreateData commonDataForResource) {
				// Import
				mockGetClientCertificate(m, testCreateData)
				mockListClientCertificateVersions(m, testCreateData.Versions, testCreateData.CertificateID).Return(nil, fmt.Errorf("list failed"))
			},
			importData:  testOneVersion,
			expectError: regexp.MustCompile(`list failed`),
			steps: []resource.TestStep{
				{
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_mtlskeystore_client_certificate_third_party.test",
					ExpectError:   regexp.MustCompile(`list failed`),
					Config:        testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/import_one_version.tf"),
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
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/missing_certificate_name.tf"),
					ExpectError: regexp.MustCompile(`The argument "certificate_name" is required, but no definition was found`),
				},
			},
		},
		"error - missing contract_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/missing_contract_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
				},
			},
		},
		"error - missing geography": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/missing_geography.tf"),
					ExpectError: regexp.MustCompile(`The argument "geography" is required, but no definition was found`),
				},
			},
		},
		"error - missing group_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/missing_group_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
				},
			},
		},
		"error - missing secure_network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/missing_secure_network.tf"),
					ExpectError: regexp.MustCompile(`The argument "secure_network" is required, but no definition was found`),
				},
			},
		},
		"error - missing version": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/missing_version.tf"),
					ExpectError: regexp.MustCompile(`Attribute versions map must contain at least 1 elements and at most 5\nelements, got: 0`),
				},
			},
		},
		"error - invalid key algorithm": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/invalid_key_algorithm.tf"),
					ExpectError: regexp.MustCompile(`Attribute key_algorithm value must be one of: \["RSA" "ECDSA"], got: "INVALID"`),
				},
			},
		},
		"error - invalid secure network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/invalid_secure_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute secure_network value must be one of: \["STANDARD_TLS"\n"ENHANCED_TLS"], got: "INVALID"`),
				},
			},
		},
		"error - invalid geography": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/invalid_geography.tf"),
					ExpectError: regexp.MustCompile(`Attribute geography value must be one of: \["CORE" "RUSSIA_AND_CORE"\n"CHINA_AND_CORE"], got: "INVALID"`),
				},
			},
		},
		"error - invalid subject": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/invalid_subject.tf"),
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
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientCertificateThirdParty/more_than_5_versions.tf"),
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
	if testData.KeyAlgorithm == "" {
		testData.KeyAlgorithm = "RSA"
	}
	keyAlgorithm := mtlskeystore.CryptographicAlgorithm(testData.KeyAlgorithm)

	request := mtlskeystore.CreateClientCertificateRequest{
		CertificateName:    testData.CertificateName,
		ContractID:         testData.ContractID,
		Geography:          mtlskeystore.Geography(testData.Geography),
		GroupID:            testData.GroupID,
		KeyAlgorithm:       &keyAlgorithm,
		NotificationEmails: testData.NotificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(testData.SecureNetwork),
		Subject:            ptr.To(testData.Subject),
		Signer:             mtlskeystore.SignerThirdParty,
	}
	//if testData.Subject != "" {
	//	request.Subject = &testData.Subject
	//}

	return m.On("CreateClientCertificate", testutils.MockContext, request).Return(&mtlskeystore.CreateClientCertificateResponse{
		CertificateID:      testData.CertificateID,
		CertificateName:    testData.CertificateName,
		Geography:          mtlskeystore.Geography(testData.Geography),
		KeyAlgorithm:       keyAlgorithm,
		NotificationEmails: testData.NotificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(testData.SecureNetwork),
		Signer:             mtlskeystore.SignerThirdParty,
		Subject:            testData.Subject,
	}, nil).Once()
}

func mockRotateClientCertificateVersion(m *mtlskeystore.Mock, testData map[string]versionData, certificateID int64, versionKey string) *mock.Call {
	response := mtlskeystore.RotateClientCertificateVersionResponse{
		Version:            testData[versionKey].Version,
		VersionGUID:        testData[versionKey].VersionGUID,
		CreatedBy:          testData[versionKey].CreatedBy,
		CreatedDate:        testData[versionKey].CreatedDate,
		ExpiryDate:         testData[versionKey].ExpiryDate,
		IssuedDate:         testData[versionKey].IssuedDate,
		Issuer:             testData[versionKey].Issuer,
		KeyAlgorithm:       mtlskeystore.CryptographicAlgorithm(testData[versionKey].KeyAlgorithm),
		KeyEllipticCurve:   testData[versionKey].KeyEllipticCurve,
		KeySizeInBytes:     testData[versionKey].KeySizeInBytes,
		SignatureAlgorithm: testData[versionKey].SignatureAlgorithm,
		Status:             mtlskeystore.CertificateVersionStatus(testData[versionKey].Status),
		Subject:            testData[versionKey].Subject,
	}
	if testData[versionKey].CertificateBlock != (certificateBlock{}) {
		response.CertificateBlock = &mtlskeystore.CertificateBlock{
			Certificate: testData[versionKey].CertificateBlock.Certificate,
			TrustChain:  testData[versionKey].CertificateBlock.TrustChain,
		}
	}
	if testData[versionKey].CSRBlock != (csrBlock{}) {
		response.CSRBlock = &mtlskeystore.CSRBlock{
			CSR:          testData[versionKey].CSRBlock.CSR,
			KeyAlgorithm: mtlskeystore.CryptographicAlgorithm(testData[versionKey].CSRBlock.KeyAlgorithm),
		}
	}
	if testData[versionKey].CertificateSubmittedBy != "" {
		response.CertificateSubmittedBy = ptr.To(testData[versionKey].CertificateSubmittedBy)
	}
	if testData[versionKey].CertificateSubmittedDate != "" {
		response.CertificateSubmittedDate = ptr.To(testData[versionKey].CertificateSubmittedDate)
	}
	if testData[versionKey].DeleteRequestedDate != "" {
		response.DeleteRequestedDate = ptr.To(testData[versionKey].DeleteRequestedDate)
	}
	if testData[versionKey].DeployedDate != "" {
		response.DeployedDate = ptr.To(testData[versionKey].DeployedDate)
	}
	if testData[versionKey].ScheduledDeleteDate != "" {
		response.ScheduledDeleteDate = ptr.To(testData[versionKey].ScheduledDeleteDate)
	}

	return m.On("RotateClientCertificateVersion", testutils.MockContext, mtlskeystore.RotateClientCertificateVersionRequest{
		CertificateID: certificateID,
	}).Return(&response, nil).Once()
}

func mockPatchClientCertificate(m *mtlskeystore.Mock, testData commonDataForResource) *mock.Call {
	return m.On("PatchClientCertificate", testutils.MockContext, mtlskeystore.PatchClientCertificateRequest{
		CertificateID: testData.CertificateID,
		Body: mtlskeystore.PatchClientCertificateRequestBody{
			CertificateName:    &testData.CertificateName,
			NotificationEmails: testData.NotificationEmails,
		},
	}).Return(nil).Once()
}

func mockDeleteClientCertificateVersion(m *mtlskeystore.Mock, versions map[string]versionData, resp *mtlskeystore.DeleteClientCertificateVersionResponse, certificateID int64, versionKey string) *mock.Call {
	return m.On("DeleteClientCertificateVersion", testutils.MockContext, mtlskeystore.DeleteClientCertificateVersionRequest{
		CertificateID: certificateID,
		Version:       versions[versionKey].Version,
	}).Return(resp, nil).Once()
}

func mockGetClientCertificate(m *mtlskeystore.Mock, testData commonDataForResource) *mock.Call {
	return m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
		CertificateID: testData.CertificateID,
	}).Return(&mtlskeystore.GetClientCertificateResponse{
		CertificateID:      testData.CertificateID,
		CertificateName:    testData.CertificateName,
		Geography:          mtlskeystore.Geography(testData.Geography),
		KeyAlgorithm:       mtlskeystore.CryptographicAlgorithm(testData.KeyAlgorithm),
		NotificationEmails: testData.NotificationEmails,
		SecureNetwork:      mtlskeystore.SecureNetwork(testData.SecureNetwork),
		Subject:            "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test-certificate/",
	}, nil).Once()
}

func mockListClientCertificateVersions(m *mtlskeystore.Mock, versions map[string]versionData, certificateID int64) *mock.Call {
	responseVersions := make([]mtlskeystore.ClientCertificateVersion, 0, len(versions))

	for _, version := range versions {
		certificateVersions := mtlskeystore.ClientCertificateVersion{
			Version:            version.Version,
			VersionGUID:        version.VersionGUID,
			CreatedBy:          version.CreatedBy,
			CreatedDate:        version.CreatedDate,
			ExpiryDate:         version.ExpiryDate,
			IssuedDate:         version.IssuedDate,
			Issuer:             version.Issuer,
			KeyAlgorithm:       mtlskeystore.CryptographicAlgorithm(version.KeyAlgorithm),
			KeyEllipticCurve:   version.KeyEllipticCurve,
			KeySizeInBytes:     version.KeySizeInBytes,
			SignatureAlgorithm: version.SignatureAlgorithm,
			Status:             mtlskeystore.CertificateVersionStatus(version.Status),
			Subject:            version.Subject,
		}

		if version.CertificateBlock != (certificateBlock{}) {
			certificateVersions.CertificateBlock = &mtlskeystore.CertificateBlock{
				Certificate: version.CertificateBlock.Certificate,
				TrustChain:  version.CertificateBlock.TrustChain,
			}
		}
		if version.CSRBlock != (csrBlock{}) {
			certificateVersions.CSRBlock = &mtlskeystore.CSRBlock{
				CSR:          version.CSRBlock.CSR,
				KeyAlgorithm: mtlskeystore.CryptographicAlgorithm(version.CSRBlock.KeyAlgorithm),
			}
		}
		if version.CertificateSubmittedBy != "" {
			certificateVersions.CertificateSubmittedBy = ptr.To(version.CertificateSubmittedBy)
		}
		if version.CertificateSubmittedDate != "" {
			certificateVersions.CertificateSubmittedDate = ptr.To(version.CertificateSubmittedDate)
		}
		if version.DeleteRequestedDate != "" {
			certificateVersions.DeleteRequestedDate = ptr.To(version.DeleteRequestedDate)
		}
		if version.DeployedDate != "" {
			certificateVersions.DeployedDate = ptr.To(version.DeployedDate)
		}
		if version.ScheduledDeleteDate != "" {
			certificateVersions.ScheduledDeleteDate = ptr.To(version.ScheduledDeleteDate)
		}
		for _, property := range version.AssociatedProperties {
			certificateVersions.AssociatedProperties = append(certificateVersions.AssociatedProperties, mtlskeystore.AssociatedProperty{
				AssetID:         property.AssetID,
				GroupID:         property.GroupID,
				PropertyName:    property.PropertyName,
				PropertyVersion: property.PropertyVersion,
			})
		}
		for _, err := range version.Validation.Errors {
			certificateVersions.Validation.Errors = append(certificateVersions.Validation.Errors, mtlskeystore.ValidationDetail{
				Message: err.Message,
				Reason:  err.Reason,
				Type:    err.Type,
			})
		}
		for _, warning := range version.Validation.Warnings {
			certificateVersions.Validation.Warnings = append(certificateVersions.Validation.Warnings, mtlskeystore.ValidationDetail{
				Message: warning.Message,
				Reason:  warning.Reason,
				Type:    warning.Type,
			})
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
