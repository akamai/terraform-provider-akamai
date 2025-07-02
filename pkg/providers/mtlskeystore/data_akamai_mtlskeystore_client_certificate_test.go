package mtlskeystore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	tst "github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClientCertificateDataSource(t *testing.T) {
	t.Parallel()
	baseChecker := test.NewStateChecker("data.akamai_mtlskeystore_client_certificate.test").
		CheckEqual("certificate_id", "1234").
		CheckEqual("certificate_name", "test-name").
		CheckEqual("created_by", "jkowalski").
		CheckEqual("created_date", "2024-05-18T23:08:07Z").
		CheckEqual("geography", "CORE").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "jkowalski@akamai.com").
		CheckEqual("notification_emails.1", "jsmith@akamai.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("signer", "THIRD_PARTY").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/").
		CheckEqual("versions.#", "2").
		CheckEqual("versions.0.version", "1").
		CheckEqual("versions.0.version_guid", "test1234").
		CheckEqual("versions.0.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----").
		CheckEqual("versions.0.certificate_block.key_algorithm", "RSA").
		CheckEqual("versions.0.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----").
		CheckEqual("versions.0.certificate_submitted_by", "jkowalski").
		CheckEqual("versions.0.certificate_submitted_date", "2024-05-18T23:08:07Z").
		CheckEqual("versions.0.created_by", "jkowalski").
		CheckEqual("versions.0.created_date", "2024-05-17T23:08:07Z").
		CheckEqual("versions.0.csr_block.csr", "test-csr").
		CheckEqual("versions.0.csr_block.key_algorithm", "RSA").
		CheckEqual("versions.0.deployed_date", "2024-05-18T23:08:07Z").
		CheckEqual("versions.0.expiry_date", "2027-05-18T23:08:07Z").
		CheckEqual("versions.0.issued_date", "2027-05-18T23:08:07Z").
		CheckEqual("versions.0.issuer", "test-issuer").
		CheckEqual("versions.0.key_algorithm", "RSA").
		CheckEqual("versions.0.key_elliptic_curve", "test-rsa").
		CheckEqual("versions.0.key_size_in_bytes", "2048").
		CheckEqual("versions.0.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("versions.0.status", "DEPLOYED").
		CheckEqual("versions.0.subject", "/C=US/O=Akamai Technologies/OU=test/CN=test/").
		CheckEqual("versions.1.validation.errors.#", "0").
		CheckEqual("versions.1.validation.warnings.#", "0").
		CheckMissing("include_associated_properties").
		CheckMissing("versions.0.current").
		CheckMissing("versions.0.previous").
		CheckMissing("versions.1.current").
		CheckMissing("versions.1.previous")

	currentAndPreviousChecker := test.NewStateChecker("data.akamai_mtlskeystore_client_certificate.test").
		CheckEqual("certificate_id", "1234").
		CheckEqual("certificate_name", "test-name").
		CheckEqual("created_by", "jkowalski").
		CheckEqual("created_date", "2024-05-18T23:08:07Z").
		CheckEqual("geography", "CORE").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "jkowalski@akamai.com").
		CheckEqual("notification_emails.1", "jsmith@akamai.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("signer", "THIRD_PARTY").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/").
		CheckEqual("versions.#", "0").
		CheckEqual("current.version", "1").
		CheckEqual("current.csr_block.csr", "test-csr").
		CheckEqual("previous.version", "2").
		CheckEqual("previous.csr_block.csr", "test-csr")

	akamaiSignerNoCSRBlockChecker := test.NewStateChecker("data.akamai_mtlskeystore_client_certificate.test").
		CheckEqual("certificate_id", "1234").
		CheckEqual("certificate_name", "test-name").
		CheckEqual("created_by", "jkowalski").
		CheckEqual("created_date", "2024-05-18T23:08:07Z").
		CheckEqual("geography", "CORE").
		CheckEqual("key_algorithm", "RSA").
		CheckEqual("notification_emails.#", "2").
		CheckEqual("notification_emails.0", "jkowalski@akamai.com").
		CheckEqual("notification_emails.1", "jsmith@akamai.com").
		CheckEqual("secure_network", "STANDARD_TLS").
		CheckEqual("signer", "AKAMAI").
		CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/").
		CheckEqual("versions.#", "0").
		CheckEqual("current.version", "1").
		CheckMissing("current.csr_block").
		CheckMissing("previous")

	baseThirdPartyResponse := &mtlskeystore.GetClientCertificateResponse{
		CertificateID:   1234,
		CertificateName: "test-name",
		CreatedBy:       "jkowalski",
		CreatedDate:     tst.NewTimeFromString(t, "2024-05-18T23:08:07Z"),
		Geography:       "CORE",
		KeyAlgorithm:    "RSA",
		NotificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		SecureNetwork: "STANDARD_TLS",
		Signer:        "THIRD_PARTY",
		Subject:       "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/",
	}

	baseAkamaiResponse := &mtlskeystore.GetClientCertificateResponse{
		CertificateID:   1234,
		CertificateName: "test-name",
		CreatedBy:       "jkowalski",
		CreatedDate:     tst.NewTimeFromString(t, "2024-05-18T23:08:07Z"),
		Geography:       "CORE",
		KeyAlgorithm:    "RSA",
		NotificationEmails: []string{
			"jkowalski@akamai.com",
			"jsmith@akamai.com",
		},
		SecureNetwork: "STANDARD_TLS",
		Signer:        "AKAMAI",
		Subject:       "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/",
	}

	baseVersionsResponse := &mtlskeystore.ListClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				VersionGUID: "test1234",
				CertificateBlock: &mtlskeystore.CertificateBlock{
					Certificate:  "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
					KeyAlgorithm: "RSA",
					TrustChain:   "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----",
				},
				CertificateSubmittedBy:   ptr.To("jkowalski"),
				CertificateSubmittedDate: ptr.To("2024-05-18T23:08:07Z"),
				CreatedBy:                "jkowalski",
				CreatedDate:              "2024-05-17T23:08:07Z",
				CSRBlock: &mtlskeystore.CSRBlock{
					CSR:          "test-csr",
					KeyAlgorithm: "RSA",
				},
				DeployedDate:       ptr.To("2024-05-18T23:08:07Z"),
				ExpiryDate:         "2027-05-18T23:08:07Z",
				IssuedDate:         "2027-05-18T23:08:07Z",
				Issuer:             "test-issuer",
				KeyAlgorithm:       "RSA",
				KeyEllipticCurve:   "test-rsa",
				KeySizeInBytes:     "2048",
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "DEPLOYED",
				Subject:            "/C=US/O=Akamai Technologies/OU=test/CN=test/",
				Validation: mtlskeystore.ValidationResult{
					Errors:   []mtlskeystore.ValidationDetail{},
					Warnings: []mtlskeystore.ValidationDetail{},
				},
			},
			{
				Version:     2,
				VersionGUID: "test12345",
				CertificateBlock: &mtlskeystore.CertificateBlock{
					Certificate:  "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
					KeyAlgorithm: "RSA",
					TrustChain:   "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----",
				},
				CertificateSubmittedBy:   ptr.To("jkowalski"),
				CertificateSubmittedDate: ptr.To("2024-05-18T23:08:07Z"),
				CreatedBy:                "jkowalski",
				CreatedDate:              "2024-05-17T23:08:07Z",
				CSRBlock: &mtlskeystore.CSRBlock{
					CSR:          "test-csr",
					KeyAlgorithm: "RSA",
				},
				DeleteRequestedDate: ptr.To("2024-05-19T23:08:07Z"),
				DeployedDate:        ptr.To("2024-05-18T23:08:07Z"),
				ExpiryDate:          "2027-05-18T23:08:07Z",
				IssuedDate:          "2027-05-18T23:08:07Z",
				Issuer:              "test-issuer",
				KeyAlgorithm:        "ECDSA",
				KeyEllipticCurve:    "test-ecdsa",
				ScheduledDeleteDate: ptr.To("2024-05-20T23:08:07Z"),
				SignatureAlgorithm:  "SHA256_WITH_RSA",
				Status:              "DELETE_PENDING",
				Subject:             "/C=US/O=Akamai Technologies/OU=test/CN=test/",
				Validation: mtlskeystore.ValidationResult{
					Errors: []mtlskeystore.ValidationDetail{
						{
							Message: "test-error-message",
							Reason:  "test-error-reason",
							Type:    "test-error-type",
						},
					},
					Warnings: []mtlskeystore.ValidationDetail{
						{
							Message: "test-warning-message",
							Reason:  "test-warning-reason",
							Type:    "test-warning-type",
						},
					},
				},
			},
		},
	}

	versionsResponseWithAssociatedProperties := &mtlskeystore.ListClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				VersionGUID: "test1234",
				CertificateBlock: &mtlskeystore.CertificateBlock{
					Certificate:  "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
					KeyAlgorithm: "RSA",
					TrustChain:   "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----",
				},
				CertificateSubmittedBy:   ptr.To("jkowalski"),
				CertificateSubmittedDate: ptr.To("2024-05-18T23:08:07Z"),
				CreatedBy:                "jkowalski",
				CreatedDate:              "2024-05-17T23:08:07Z",
				CSRBlock: &mtlskeystore.CSRBlock{
					CSR:          "test-csr",
					KeyAlgorithm: "RSA",
				},
				DeployedDate:       ptr.To("2024-05-18T23:08:07Z"),
				ExpiryDate:         "2027-05-18T23:08:07Z",
				IssuedDate:         "2027-05-18T23:08:07Z",
				Issuer:             "test-issuer",
				KeyAlgorithm:       "RSA",
				KeyEllipticCurve:   "test-rsa",
				KeySizeInBytes:     "2048",
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "DEPLOYED",
				Subject:            "/C=US/O=Akamai Technologies/OU=test/CN=test/",
				Validation: mtlskeystore.ValidationResult{
					Errors:   []mtlskeystore.ValidationDetail{},
					Warnings: []mtlskeystore.ValidationDetail{},
				},
				AssociatedProperties: []mtlskeystore.AssociatedProperty{
					{
						AssetID:         1234,
						GroupID:         12345,
						PropertyName:    "test-property-name",
						PropertyVersion: 1,
					},
				},
			},
		},
	}

	versionsResponseWithCurrentAndPrevious := &mtlskeystore.ListClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				VersionGUID: "test1234",
				CertificateBlock: &mtlskeystore.CertificateBlock{
					Certificate:  "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
					KeyAlgorithm: "RSA",
					TrustChain:   "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----",
				},
				VersionAlias:             ptr.To("CURRENT"),
				CertificateSubmittedBy:   ptr.To("jkowalski"),
				CertificateSubmittedDate: ptr.To("2024-05-18T23:08:07Z"),
				CreatedBy:                "jkowalski",
				CreatedDate:              "2024-05-17T23:08:07Z",
				CSRBlock: &mtlskeystore.CSRBlock{
					CSR:          "test-csr",
					KeyAlgorithm: "RSA",
				},
				DeployedDate:       ptr.To("2024-05-18T23:08:07Z"),
				ExpiryDate:         "2027-05-18T23:08:07Z",
				IssuedDate:         "2027-05-18T23:08:07Z",
				Issuer:             "test-issuer",
				KeyAlgorithm:       "RSA",
				KeyEllipticCurve:   "test-rsa",
				KeySizeInBytes:     "2048",
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "DEPLOYED",
				Subject:            "/C=US/O=Akamai Technologies/OU=test/CN=test/",
				Validation: mtlskeystore.ValidationResult{
					Errors:   []mtlskeystore.ValidationDetail{},
					Warnings: []mtlskeystore.ValidationDetail{},
				},
			},
			{
				Version:     2,
				VersionGUID: "test12345",
				CertificateBlock: &mtlskeystore.CertificateBlock{
					Certificate:  "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
					KeyAlgorithm: "RSA",
					TrustChain:   "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----",
				},
				VersionAlias:             ptr.To("PREVIOUS"),
				CertificateSubmittedBy:   ptr.To("jkowalski"),
				CertificateSubmittedDate: ptr.To("2024-05-18T23:08:07Z"),
				CreatedBy:                "jkowalski",
				CreatedDate:              "2024-05-17T23:08:07Z",
				CSRBlock: &mtlskeystore.CSRBlock{
					CSR:          "test-csr",
					KeyAlgorithm: "RSA",
				},
				DeleteRequestedDate: ptr.To("2024-05-19T23:08:07Z"),
				DeployedDate:        ptr.To("2024-05-18T23:08:07Z"),
				ExpiryDate:          "2027-05-18T23:08:07Z",
				IssuedDate:          "2027-05-18T23:08:07Z",
				Issuer:              "test-issuer",
				KeyAlgorithm:        "ECDSA",
				KeyEllipticCurve:    "test-ecdsa",
				ScheduledDeleteDate: ptr.To("2024-05-20T23:08:07Z"),
				SignatureAlgorithm:  "SHA256_WITH_RSA",
				Status:              "DELETE_PENDING",
				Subject:             "/C=US/O=Akamai Technologies/OU=test/CN=test/",
				Validation: mtlskeystore.ValidationResult{
					Errors: []mtlskeystore.ValidationDetail{
						{
							Message: "test-error-message",
							Reason:  "test-error-reason",
							Type:    "test-error-type",
						},
					},
					Warnings: []mtlskeystore.ValidationDetail{
						{
							Message: "test-warning-message",
							Reason:  "test-warning-reason",
							Type:    "test-warning-type",
						},
					},
				},
			},
		},
	}

	versionsResponseWithoutCSRBlock := &mtlskeystore.ListClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				VersionGUID: "test1234",
				CertificateBlock: &mtlskeystore.CertificateBlock{
					Certificate:  "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
					KeyAlgorithm: "RSA",
					TrustChain:   "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----",
				},
				VersionAlias:             ptr.To("CURRENT"),
				CertificateSubmittedBy:   ptr.To("jkowalski"),
				CertificateSubmittedDate: ptr.To("2024-05-18T23:08:07Z"),
				CreatedBy:                "jkowalski",
				CreatedDate:              "2024-05-17T23:08:07Z",
				DeployedDate:             ptr.To("2024-05-18T23:08:07Z"),
				ExpiryDate:               "2027-05-18T23:08:07Z",
				IssuedDate:               "2027-05-18T23:08:07Z",
				Issuer:                   "test-issuer",
				KeyAlgorithm:             "RSA",
				KeyEllipticCurve:         "test-rsa",
				KeySizeInBytes:           "2048",
				SignatureAlgorithm:       "SHA256_WITH_RSA",
				Status:                   "DEPLOYED",
				Subject:                  "/C=US/O=Akamai Technologies/OU=test/CN=test/",
				Validation: mtlskeystore.ValidationResult{
					Errors:   []mtlskeystore.ValidationDetail{},
					Warnings: []mtlskeystore.ValidationDetail{},
				},
			},
		},
	}

	tests := map[string]struct {
		init  func(*mtlskeystore.Mock)
		steps []resource.TestStep
	}{
		"happy path without associated properties": {
			init: func(m *mtlskeystore.Mock) {
				mockClientCertificate(m, 1234, false, baseThirdPartyResponse, baseVersionsResponse)
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 1234
						}`,
					Check: baseChecker.
						CheckEqual("versions.1.version", "2").
						CheckEqual("versions.1.version_guid", "test12345").
						CheckEqual("versions.1.certificate_block.certificate", "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----").
						CheckEqual("versions.1.certificate_block.key_algorithm", "RSA").
						CheckEqual("versions.1.certificate_block.trust_chain", "-----BEGIN CERTIFICATE-----\ntest-trust\n-----END CERTIFICATE-----").
						CheckEqual("versions.1.certificate_submitted_by", "jkowalski").
						CheckEqual("versions.1.certificate_submitted_date", "2024-05-18T23:08:07Z").
						CheckEqual("versions.1.created_by", "jkowalski").
						CheckEqual("versions.1.created_date", "2024-05-17T23:08:07Z").
						CheckEqual("versions.1.delete_requested_date", "2024-05-19T23:08:07Z").
						CheckEqual("versions.1.deployed_date", "2024-05-18T23:08:07Z").
						CheckEqual("versions.1.expiry_date", "2027-05-18T23:08:07Z").
						CheckEqual("versions.1.issued_date", "2027-05-18T23:08:07Z").
						CheckEqual("versions.1.issuer", "test-issuer").
						CheckEqual("versions.1.key_algorithm", "ECDSA").
						CheckEqual("versions.1.key_elliptic_curve", "test-ecdsa").
						CheckEqual("versions.1.scheduled_delete_date", "2024-05-20T23:08:07Z").
						CheckEqual("versions.1.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("versions.1.status", "DELETE_PENDING").
						CheckEqual("versions.1.subject", "/C=US/O=Akamai Technologies/OU=test/CN=test/").
						CheckEqual("versions.1.validation.errors.#", "1").
						CheckEqual("versions.1.validation.warnings.#", "1").
						CheckEqual("versions.1.validation.errors.0.message", "test-error-message").
						CheckEqual("versions.1.validation.errors.0.reason", "test-error-reason").
						CheckEqual("versions.1.validation.errors.0.type", "test-error-type").
						CheckEqual("versions.1.validation.warnings.0.message", "test-warning-message").
						CheckEqual("versions.1.validation.warnings.0.reason", "test-warning-reason").
						CheckEqual("versions.1.validation.warnings.0.type", "test-warning-type").
						CheckMissing("versions.0.associated_properties").
						CheckMissing("versions.1.associated_properties").
						Build(),
				},
			},
		},
		"happy path with associated properties": {
			init: func(m *mtlskeystore.Mock) {
				mockClientCertificate(m, 1234, true, baseThirdPartyResponse, versionsResponseWithAssociatedProperties)
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 1234
							include_associated_properties = true
						}`,
					Check: baseChecker.
						CheckEqual("include_associated_properties", "true").
						CheckEqual("versions.#", "1").
						CheckEqual("versions.0.properties.#", "1").
						CheckEqual("versions.0.properties.0.asset_id", "1234").
						CheckEqual("versions.0.properties.0.group_id", "12345").
						CheckEqual("versions.0.properties.0.property_name", "test-property-name").
						CheckEqual("versions.0.properties.0.property_version", "1").
						Build(),
				},
			},
		},
		"happy path - empty versions": {
			init: func(m *mtlskeystore.Mock) {
				mockClientCertificate(m, 12344, false, baseThirdPartyResponse, &mtlskeystore.ListClientCertificateVersionsResponse{Versions: []mtlskeystore.ClientCertificateVersion{}})
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 12344
						}`,
					Check: test.NewStateChecker("data.akamai_mtlskeystore_client_certificate.test").
						CheckEqual("certificate_id", "12344").
						CheckEqual("certificate_name", "test-name").
						CheckEqual("created_by", "jkowalski").
						CheckEqual("created_date", "2024-05-18T23:08:07Z").
						CheckEqual("geography", "CORE").
						CheckEqual("key_algorithm", "RSA").
						CheckEqual("notification_emails.#", "2").
						CheckEqual("notification_emails.0", "jkowalski@akamai.com").
						CheckEqual("notification_emails.1", "jsmith@akamai.com").
						CheckEqual("secure_network", "STANDARD_TLS").
						CheckEqual("signer", "THIRD_PARTY").
						CheckEqual("subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/").
						CheckEqual("versions.#", "0").
						Build(),
				},
			},
		},
		"happy path with CURRENT and PREVIOUS versions": {
			init: func(m *mtlskeystore.Mock) {
				mockClientCertificate(m, 1234, false, baseThirdPartyResponse, versionsResponseWithCurrentAndPrevious)
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 1234
						}`,
					Check: currentAndPreviousChecker.Build(),
				},
			},
		},
		"happy path with AKAMAI signer certificate without CSR block": {
			init: func(m *mtlskeystore.Mock) {
				mockClientCertificate(m, 1234, false, baseAkamaiResponse, versionsResponseWithoutCSRBlock)
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 1234
						}`,
					Check: akamaiSignerNoCSRBlockChecker.Build(),
				},
			},
		},
		"error response from GetClientCertificate": {
			init: func(m *mtlskeystore.Mock) {
				m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
					CertificateID: 123,
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 123
						}`,
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error response from ListClientCertificateVersions": {
			init: func(m *mtlskeystore.Mock) {
				m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
					CertificateID: 1233,
				}).Return(baseThirdPartyResponse, nil).Once()

				m.On("ListClientCertificateVersions", testutils.MockContext, mtlskeystore.ListClientCertificateVersionsRequest{
					CertificateID: 1233,
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
						  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 1233
						}`,
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error - nil versions response": {
			init: func(m *mtlskeystore.Mock) {
				m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
					CertificateID: 12345,
				}).Return(baseThirdPartyResponse, nil).Once()

				m.On("ListClientCertificateVersions", testutils.MockContext, mtlskeystore.ListClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(nil, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
							  edgerc = "../../common/testutils/edgerc"	
						}
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 12345
						}`,
					ExpectError: regexp.MustCompile("Unexpected nil response for client certificate versions."),
				},
			},
		},
		"validation error - missing certificateID": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
							  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
						}`,
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
			},
		},
		"validation error - wrong certificateID": {
			steps: []resource.TestStep{
				{
					Config: `
						provider "akamai" {
							  edgerc = "../../common/testutils/edgerc"
						}
						
						data "akamai_mtlskeystore_client_certificate" "test" {
							certificate_id = 0
						}`,
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
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

func mockClientCertificate(m *mtlskeystore.Mock, certificateID int64, includeProp bool, response *mtlskeystore.GetClientCertificateResponse, versionResponse *mtlskeystore.ListClientCertificateVersionsResponse) {
	m.On("GetClientCertificate", testutils.MockContext, mtlskeystore.GetClientCertificateRequest{
		CertificateID: certificateID,
	}).Return(response, nil).Times(3)

	m.On("ListClientCertificateVersions", testutils.MockContext, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID:               certificateID,
		IncludeAssociatedProperties: includeProp,
	}).Return(versionResponse, nil).Times(3)
}
