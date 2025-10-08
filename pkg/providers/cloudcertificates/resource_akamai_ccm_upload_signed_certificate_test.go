package cloudcertificates

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testprovider"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const (
	testSignedCertificatePEM = "-----BEGIN CERTIFICATE-----\ntestsignedcertificate\n-----END CERTIFICATE-----\n"

	testTrustChainPEM = "-----BEGIN CERTIFICATE-----\ntesttrustchaincertificate1\n-----END CERTIFICATE-----\n" +
		"-----BEGIN CERTIFICATE-----\ntesttrustchaincertificate2\n-----END CERTIFICATE-----\n"

	testRenewedCertificatePEM = "-----BEGIN CERTIFICATE-----\ntestrenewedsignedcertificate\n-----END CERTIFICATE-----\n"
)

var mockCerts = mockCertificates{
	minimumCertificate: mockCertificate{
		cert: ccm.Certificate{
			CertificateID:     "12345",
			CertificateStatus: "CSR_READY",
		},
	},
	signedCertificate: mockCertificate{
		cert: ccm.Certificate{
			CertificateID:                       "12345",
			ModifiedDate:                        tst.NewTimeFromStringMust("2025-09-23T07:26:30.616267Z"),
			ModifiedBy:                          "jsmith",
			CertificateStatus:                   "READY_FOR_USE",
			SignedCertificatePEM:                ptr.To(testSignedCertificatePEM),
			SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-12-23T08:19:47Z")),
			SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2025-09-23T07:19:47Z")),
			SignedCertificateSerialNumber:       ptr.To("12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF"),
			SignedCertificateSHA256Fingerprint:  ptr.To("FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10"),
			SignedCertificateIssuer:             ptr.To("CN=Test Issuer, O=Test Org, C=US"),
		},
	},
	minimumRenewedCertificate: mockCertificate{
		cert: ccm.Certificate{
			CertificateID:     "23456",
			CertificateStatus: "CSR_READY",
		},
	},
	renewedCertificate: mockCertificate{
		cert: ccm.Certificate{
			CertificateID:                       "23456",
			ModifiedDate:                        tst.NewTimeFromStringMust("2025-10-23T07:26:30.616267Z"),
			ModifiedBy:                          "janesmith",
			CertificateStatus:                   "READY_FOR_USE",
			SignedCertificatePEM:                ptr.To(testRenewedCertificatePEM),
			SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2028-01-23T08:19:47Z")),
			SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2025-10-23T07:19:47Z")),
			SignedCertificateSerialNumber:       ptr.To("34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF:12"),
			SignedCertificateSHA256Fingerprint:  ptr.To("DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE"),
			SignedCertificateIssuer:             ptr.To("CN=Test Issuer 2, O=Test Inc, C=US"),
		},
	},
}

var signedCertChecker = test.NewStateChecker("akamai_cloudcertificates_upload_signed_certificate.upload").
	CheckEqual("certificate_id", "12345").
	CheckEqual("signed_certificate_pem", testSignedCertificatePEM).
	CheckEqual("acknowledge_warnings", "false").
	CheckEqual("modified_date", "2025-09-23T07:26:30.616267Z").
	CheckEqual("modified_by", "jsmith").
	CheckEqual("certificate_status", "READY_FOR_USE").
	CheckEqual("signed_certificate_not_valid_after_date", "2027-12-23T08:19:47Z").
	CheckEqual("signed_certificate_not_valid_before_date", "2025-09-23T07:19:47Z").
	CheckEqual("signed_certificate_serial_number", "12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF").
	CheckEqual("signed_certificate_sha256_fingerprint", "FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10").
	CheckEqual("signed_certificate_issuer", "CN=Test Issuer, O=Test Org, C=US")

var renewedCertChecker = test.NewStateChecker("akamai_cloudcertificates_upload_signed_certificate.upload").
	CheckEqual("certificate_id", "23456").
	CheckEqual("signed_certificate_pem", testRenewedCertificatePEM).
	CheckEqual("acknowledge_warnings", "false").
	CheckEqual("modified_date", "2025-10-23T07:26:30.616267Z").
	CheckEqual("modified_by", "janesmith").
	CheckEqual("certificate_status", "READY_FOR_USE").
	CheckEqual("signed_certificate_not_valid_after_date", "2028-01-23T08:19:47Z").
	CheckEqual("signed_certificate_not_valid_before_date", "2025-10-23T07:19:47Z").
	CheckEqual("signed_certificate_serial_number", "34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF:12").
	CheckEqual("signed_certificate_sha256_fingerprint", "DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE").
	CheckEqual("signed_certificate_issuer", "CN=Test Issuer 2, O=Test Inc, C=US")

var signedCertImportChecker = test.NewImportChecker().
	CheckEqual("certificate_id", "12345").
	CheckEqual("signed_certificate_pem", testSignedCertificatePEM).
	CheckEqual("acknowledge_warnings", "false").
	CheckEqual("modified_date", "2025-09-23T07:26:30.616267Z").
	CheckEqual("modified_by", "jsmith").
	CheckEqual("certificate_status", "READY_FOR_USE").
	CheckEqual("signed_certificate_not_valid_after_date", "2027-12-23T08:19:47Z").
	CheckEqual("signed_certificate_not_valid_before_date", "2025-09-23T07:19:47Z").
	CheckEqual("signed_certificate_serial_number", "12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF").
	CheckEqual("signed_certificate_sha256_fingerprint", "FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10").
	CheckEqual("signed_certificate_issuer", "CN=Test Issuer, O=Test Org, C=US")

func TestUploadSignedCertificateResource(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		init  func(*ccm.Mock, mockCertificates)
		steps []resource.TestStep
	}{
		"happy path - upload signed certificate PEM": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					Check: signedCertChecker.
						CheckMissing("trust_chain_pem").
						Build(),
				},
			},
		},
		"happy path - upload signed certificate PEM and trustchain": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.cert.TrustChainPEM = ptr.To(testTrustChainPEM)
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_trustchain.tf"),
					Check: signedCertChecker.
						CheckEqual("trust_chain_pem", testTrustChainPEM).
						Build(),
				},
			},
		},
		"happy path - upload signed certificate with acknowledge_warnings enabled": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.AcknowledgeWarnings = true
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_acknowledge_warnings.tf"),
					Check: signedCertChecker.
						CheckMissing("trust_chain_pem").
						CheckEqual("acknowledge_warnings", "true").
						Build(),
				},
			},
		},
		"happy path - certificate renewal": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)

				// UPDATE
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan x 2
				mc.minimumRenewedCertificate.mockGet(m).Twice()
				// Update
				mc.renewedCertificate.mockPatch(m)
				// Plan
				mc.renewedCertificate.mockGet(m)
				// Read
				mc.renewedCertificate.mockGet(m)
				// Plan
				mc.renewedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					Check: signedCertChecker.
						CheckMissing("trust_chain_pem").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/renewal.tf"),
					Check: renewedCertChecker.
						CheckMissing("trust_chain_pem").
						Build(),
				},
			},
		},
		"happy path - upload signed certificate PEM with params provided by another resource": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan
				mc.minimumCertificate.mockGet(m)
				// Create
				mc.signedCertificate.cert.TrustChainPEM = ptr.To(testTrustChainPEM)
				mc.signedCertificate.AcknowledgeWarnings = true
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/another_resource.tf"),
					Check: signedCertChecker.
						CheckEqual("trust_chain_pem", testTrustChainPEM).
						CheckEqual("acknowledge_warnings", "true").
						Build(),
				},
			},
		},
		"import a signed certificate": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Import
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)

				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: signedCertImportChecker.
						CheckMissing("trust_chain_pem").
						Build(),
					ImportStateId:      "12345",
					ImportState:        true,
					ResourceName:       "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ImportStatePersist: true,
				},
				{
					// Confirm idempotency after import
					Config:   testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					PlanOnly: true,
				},
			},
		},
		"import a signed certificate with trustchain": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				mc.signedCertificate.cert.TrustChainPEM = ptr.To(testTrustChainPEM)
				// Import
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)

				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: signedCertImportChecker.
						CheckEqual("trust_chain_pem", testTrustChainPEM).
						Build(),
					ImportStateId:      "12345",
					ImportState:        true,
					ResourceName:       "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_trustchain.tf"),
					ImportStatePersist: true,
				},
				{
					// Confirm idempotency after import
					Config:   testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_trustchain.tf"),
					PlanOnly: true,
				},
			},
		},
		"import a signed certificate with acknowledge_warnings": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Import
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)

				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					ImportStateCheck: signedCertImportChecker.
						CheckEqual("acknowledge_warnings", "true").
						CheckMissing("trust_chain_pem").
						Build(),
					ImportStateId:      "12345,true",
					ImportState:        true,
					ResourceName:       "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_acknowledge_warnings.tf"),
					ImportStatePersist: true,
				},
				{
					// Confirm idempotency after import
					Config:   testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_acknowledge_warnings.tf"),
					PlanOnly: true,
				},
			},
		},
		"error - certificate already in state READY_FOR_USE": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 1
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(`The certificate '12345' has status 'READY_FOR_USE'`),
				},
			},
		},
		"error - trying to update signed certificate PEM": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)

				// UPDATE
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					Check: signedCertChecker.
						CheckMissing("trust_chain_pem").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_different_cert_pem.tf"),
					ExpectError: regexp.MustCompile(`The certificate '12345' has status 'READY_FOR_USE'`),
				},
			},
		},
		"error - trying to change only acknowledge_warnings": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)

				// UPDATE
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					Check: signedCertChecker.
						CheckMissing("trust_chain_pem").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/with_acknowledge_warnings.tf"),
					ExpectError: regexp.MustCompile(`The certificate '12345' has status 'READY_FOR_USE'`),
				},
			},
		},
		"error - 404 certificate not found in plan": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 1
				mc.minimumCertificate.err = &ccm.Error{
					Type:  "/error-types/certificate-not-found",
					Title: "Certificate subscription is not found.",
				}
				mc.minimumCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(
						`(?s)Cannot upload signed certificate to a non-existent CCM Certificate object.+` +
							`The certificate '12345' was not found on the server.+` +
							`API error:.+` +
							`Certificate subscription is not found.`),
				},
			},
		},
		"error - GetCertificate fails generally in plan": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 1
				mc.minimumCertificate.err = fmt.Errorf("API failed")
				mc.minimumCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(`(?s)Unable to get CCM Certificate for signed certificate upload.+` +
						`Error retrieving certificate '12345': API failed`),
				},
			},
		},
		"error - PatchCertificate fails in create": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.err = fmt.Errorf("API failed")
				mc.signedCertificate.mockPatch(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(`(?s)Error uploading signed certificate during resource creation.+API failed`),
				},
			},
		},
		"error - GetCertificate fails in read": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.err = fmt.Errorf("API failed")
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(`(?s)Error reading CCM Certificate.+` +
						`Error retrieving certificate '12345': API failed`),
				},
			},
		},
		"error - PatchCertificate fails in update": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Plan x 2
				mc.minimumCertificate.mockGet(m).Twice()
				// Create
				mc.signedCertificate.mockPatch(m)
				// Plan
				mc.signedCertificate.mockGet(m)
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan
				mc.signedCertificate.mockGet(m)

				// UPDATE
				// Read
				mc.signedCertificate.mockGet(m)
				// Plan x 2
				mc.minimumRenewedCertificate.mockGet(m).Twice()
				// Update
				mc.renewedCertificate.err = fmt.Errorf("API failed")
				mc.renewedCertificate.mockPatch(m)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					Check: signedCertChecker.
						CheckMissing("trust_chain_pem").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/renewal.tf"),
					ExpectError: regexp.MustCompile(`(?s)Error uploading signed certificate during resource update.+API failed`),
				},
			},
		},
		"error - missing certificate ID": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/no_certificate_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "certificate_id" is required, but no definition was found`),
				},
			},
		},
		"error - empty certificate ID": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/empty_certificate_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute certificate_id string length must be at least 1, got: 0`),
				},
			},
		},
		"error - missing signed certificate PEM": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/no_signed_certificate_pem.tf"),
					ExpectError: regexp.MustCompile(`The argument "signed_certificate_pem" is required, but no definition`),
				},
			},
		},
		"error - invalid format for signed certificate PEM": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/invalid_certificate_pem_format.tf"),
					ExpectError: regexp.MustCompile(`Attribute signed_certificate_pem must be in PEM format`),
				},
			},
		},
		"error - invalid format for trust chain PEM": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/invalid_trustchain_format.tf"),
					ExpectError: regexp.MustCompile(`Attribute trust_chain_pem must be in PEM format`),
				},
			},
		},
		"error - trying to import for cert where PEM was never uploaded": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Read
				mc.minimumCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					ImportStateId: "12345",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(
						`(?s)Cannot import CCM Certificate in 'CSR_READY' status.+` +
							`The certificate '12345' has status 'CSR_READY'.+` +
							`signed certificate PEM has not been uploaded yet.`),
				},
			},
		},
		"error - 404 cert not found in import": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Read
				mc.signedCertificate.err = &ccm.Error{
					Type:  "/error-types/certificate-not-found",
					Title: "Certificate subscription is not found.",
				}
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					ImportStateId: "12345",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(
						`(?s)CCM Certificate for import not found.+` +
							`The certificate '12345' was not found on the server.+` +
							`API error:.+` +
							`Certificate subscription is not found.`),
				},
			},
		},
		"error - GetCertificate fails generally in import": {
			init: func(m *ccm.Mock, mc mockCertificates) {
				// Read
				mc.signedCertificate.err = fmt.Errorf("API failed")
				mc.signedCertificate.mockGet(m)
			},
			steps: []resource.TestStep{
				{
					ImportStateId: "12345",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(`(?s)Error reading CCM Certificate for import.+` +
						`Error retrieving certificate '12345': API failed`),
				},
			},
		},
		"error - whitespace-only import ID": {
			steps: []resource.TestStep{
				{
					ImportStateId: "   ",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(
						`(?s)Incorrect import ID.+importID cannot be empty.+` +
							`'certificateID\[,acknowledge_warnings\]'`),
				},
			},
		},
		"error - three parts of import ID when maximum 2 allowed": {
			steps: []resource.TestStep{
				{
					ImportStateId: "12345,true,false",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(
						`(?s)Incorrect import ID.+invalid number of importID parts: 3.+` +
							`'certificateID\[,acknowledge_warnings\]'`),
				},
			},
		},
		"error - second part of import ID is not boolean": {
			steps: []resource.TestStep{
				{
					ImportStateId: "12345,foo",
					ImportState:   true,
					ResourceName:  "akamai_cloudcertificates_upload_signed_certificate.upload",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResUploadSignedCertificate/basic.tf"),
					ExpectError: regexp.MustCompile(
						`(?s)Incorrect import ID.+acknowledge_warnings must be 'true' or 'false'`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &ccm.Mock{}

			if tc.init != nil {
				tc.init(client, mockCerts)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(
						NewSubprovider(), testprovider.NewMockSubprovider()),
					Steps: tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

type mockCertificate struct {
	cert                ccm.Certificate
	AcknowledgeWarnings bool
	err                 error
}

func (c mockCertificate) mockGet(m *ccm.Mock) *mock.Call {
	return m.On("GetCertificate", testutils.MockContext, ccm.GetCertificateRequest{
		CertificateID: c.cert.CertificateID,
	}).Return(&c.cert, c.err).Once()
}

func (c mockCertificate) mockPatch(m *ccm.Mock) *mock.Call {
	var trustChainPEM string
	if c.cert.TrustChainPEM != nil {
		trustChainPEM = *c.cert.TrustChainPEM
	}
	return m.On("PatchCertificate", testutils.MockContext, ccm.PatchCertificateRequest{
		CertificateID:        c.cert.CertificateID,
		SignedCertificatePEM: *c.cert.SignedCertificatePEM,
		TrustChainPEM:        trustChainPEM,
		AcknowledgeWarnings:  c.AcknowledgeWarnings,
	}).Return(&c.cert, c.err).Once()
}

type mockCertificates struct {
	minimumCertificate        mockCertificate
	signedCertificate         mockCertificate
	minimumRenewedCertificate mockCertificate
	renewedCertificate        mockCertificate
}
