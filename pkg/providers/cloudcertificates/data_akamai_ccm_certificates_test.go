package cloudcertificates

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClientCertificateDataSource(t *testing.T) {
	t.Parallel()
	baseChecker := test.NewStateChecker("data.akamai_cloudcertificates_certificates.test").
		CheckEqual("certificates.#", "1").
		CheckEqual("certificates.0.certificate_id", "cert1_1234").
		CheckEqual("certificates.0.certificate_name", "test_certificate1").
		CheckEqual("certificates.0.sans.#", "1").
		CheckEqual("certificates.0.sans.0", "test1.example.com").
		CheckEqual("certificates.0.subject.common_name", "test1.example.com").
		CheckEqual("certificates.0.subject.country", "US").
		CheckEqual("certificates.0.subject.organization", "Test Org1").
		CheckEqual("certificates.0.subject.state", "CA").
		CheckEqual("certificates.0.subject.locality", "Test City").
		CheckEqual("certificates.0.certificate_type", "THIRD_PARTY").
		CheckEqual("certificates.0.contract_id", "A-123").
		CheckEqual("certificates.0.key_type", "RSA").
		CheckEqual("certificates.0.key_size", "2048").
		CheckEqual("certificates.0.secure_network", "ENHANCED_TLS").
		CheckEqual("certificates.0.account_id", "act_789").
		CheckEqual("certificates.0.created_date", "2024-01-01T12:00:00Z").
		CheckEqual("certificates.0.created_by", "test_user").
		CheckEqual("certificates.0.modified_date", "2024-05-01T12:00:00Z").
		CheckEqual("certificates.0.modified_by", "test_user2").
		CheckEqual("certificates.0.certificate_status", "ACTIVE").
		CheckEqual("certificates.0.csr_expiration_date", "2027-01-01T00:00:00Z").
		CheckEqual("certificates.0.signed_certificate_issuer", "O=Test Org1,L=Test City,ST=CA,C=US").
		CheckEqual("certificates.0.signed_certificate_not_valid_after_date", "2027-12-23T08:19:47Z").
		CheckEqual("certificates.0.signed_certificate_not_valid_before_date", "2025-09-23T07:19:47Z")

	baseCheckerMissingCertMaterials := baseChecker.
		CheckMissing("certificates.0.csr_pem").
		CheckMissing("certificates.0.signed_certificate_pem").
		CheckMissing("certificates.0.signed_certificate_serial_number").
		CheckMissing("certificates.0.signed_certificate_sha256_fingerprint").
		CheckMissing("certificates.0.trust_chain_pem")

	baseResponse := &ccm.ListCertificatesResponse{
		Certificates: []ccm.Certificate{
			{
				CertificateID:   "cert1_1234",
				CertificateName: "test_certificate1",
				SANs:            []string{"test1.example.com"},
				Subject: &ccm.Subject{
					CommonName:   "test1.example.com",
					Country:      "US",
					Organization: "Test Org1",
					State:        "CA",
					Locality:     "Test City",
				},
				CertificateType:                     "THIRD_PARTY",
				KeyType:                             ccm.CryptographicAlgorithmRSA,
				KeySize:                             ccm.KeySize2048,
				SecureNetwork:                       string(ccm.SecureNetworkEnhancedTLS),
				ContractID:                          "A-123",
				AccountID:                           "act_789",
				CreatedDate:                         tst.NewTimeFromStringMust("2024-01-01T12:00:00Z"),
				CreatedBy:                           "test_user",
				ModifiedDate:                        tst.NewTimeFromStringMust("2024-05-01T12:00:00Z"),
				ModifiedBy:                          "test_user2",
				CertificateStatus:                   string(ccm.StatusActive),
				CSRExpirationDate:                   tst.NewTimeFromStringMust("2027-01-01T00:00:00Z"),
				SignedCertificateIssuer:             ptr.To("O=Test Org1,L=Test City,ST=CA,C=US"),
				SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-12-23T08:19:47Z")),
				SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2025-09-23T07:19:47Z")),
			},
			{
				CertificateID:   "cert2_1234",
				CertificateName: "test_certificate2",
				SANs:            []string{"test2.example.com"},
				Subject: &ccm.Subject{
					CommonName:   "test2.example.com",
					Country:      "US",
					Organization: "Test Org2",
					State:        "CA",
					Locality:     "Test City",
				},
				CertificateType:                     "THIRD_PARTY",
				KeyType:                             ccm.CryptographicAlgorithmECDSA,
				KeySize:                             ccm.KeySizeP256,
				SecureNetwork:                       string(ccm.SecureNetworkEnhancedTLS),
				ContractID:                          "A-123",
				AccountID:                           "act_789",
				CreatedDate:                         tst.NewTimeFromStringMust("2024-05-01T12:00:00Z"),
				CreatedBy:                           "test_user",
				ModifiedDate:                        tst.NewTimeFromStringMust("2024-07-01T12:00:00Z"),
				ModifiedBy:                          "test_user2",
				CertificateStatus:                   string(ccm.StatusReadyForUse),
				CSRExpirationDate:                   tst.NewTimeFromStringMust("2027-05-01T00:00:00Z"),
				SignedCertificateIssuer:             ptr.To("O=Test Org2,L=Test City,ST=CA,C=US"),
				SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-11-15T10:30:00Z")),
				SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2025-08-15T10:30:00Z")),
			},
			{
				CertificateID:   "cert3_1234",
				CertificateName: "test_certificate3",
				SANs:            []string{"test3.example.com"},
				Subject: &ccm.Subject{
					CommonName:   "test3.example.com",
					Country:      "US",
					Organization: "Test Org3",
					State:        "CA",
					Locality:     "Test City",
				},
				CertificateType:   "THIRD_PARTY",
				KeyType:           ccm.CryptographicAlgorithmRSA,
				KeySize:           ccm.KeySize2048,
				SecureNetwork:     string(ccm.SecureNetworkEnhancedTLS),
				ContractID:        "A-123",
				AccountID:         "act_789",
				CreatedDate:       tst.NewTimeFromStringMust("2024-12-01T12:00:00Z"),
				CreatedBy:         "test_user",
				ModifiedDate:      tst.NewTimeFromStringMust("2025-01-01T12:00:00Z"),
				ModifiedBy:        "test_user2",
				CertificateStatus: string(ccm.StatusCSRReady),
				CSRExpirationDate: tst.NewTimeFromStringMust("2027-12-01T00:00:00Z"),
			},
		},
		Links: ccm.Links{
			Self:     "/ccm/v1/certificates?page=1&pageSize=100",
			Next:     nil,
			Previous: nil,
		},
	}

	certificates101 := make([]ccm.Certificate, 101)
	for i := 0; i < 101; i++ {
		certificates101[i] = ccm.Certificate{
			CertificateID:   fmt.Sprintf("cert_%d_1234", i+1),
			CertificateName: fmt.Sprintf("test_certificate_%d", i+1),
			SANs:            []string{fmt.Sprintf("test%d.example.com", i+1)},
			Subject: &ccm.Subject{
				CommonName:   fmt.Sprintf("test%d.example.com", i+1),
				Country:      "US",
				Organization: fmt.Sprintf("Test Org%d", i+1),
				State:        "CA",
				Locality:     "Test City",
			},
			CertificateType:                     "THIRD_PARTY",
			KeyType:                             ccm.CryptographicAlgorithmRSA,
			KeySize:                             ccm.KeySize2048,
			SecureNetwork:                       string(ccm.SecureNetworkEnhancedTLS),
			ContractID:                          "A-123",
			AccountID:                           "act_789",
			CreatedDate:                         tst.NewTimeFromStringMust("2024-01-01T12:00:00Z"),
			CreatedBy:                           "test_user",
			ModifiedDate:                        tst.NewTimeFromStringMust("2024-05-01T12:00:00Z"),
			ModifiedBy:                          "test_user2",
			CertificateStatus:                   string(ccm.StatusActive),
			CSRExpirationDate:                   tst.NewTimeFromStringMust("2027-01-01T00:00:00Z"),
			SignedCertificateIssuer:             ptr.To(fmt.Sprintf("O=Test Org%d,L=Test City,ST=CA,C=US", i+1)),
			SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-12-23T08:19:47Z")),
			SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2025-09-23T07:19:47Z")),
		}
	}

	tests := map[string]struct {
		init  func(*ccm.Mock)
		steps []resource.TestStep
	}{
		"happy path without optional params": {
			init: func(m *ccm.Mock) {
				mockListCertificates(m, ccm.ListCertificatesRequest{
					PageSize: 100,
					Page:     1,
				}, baseResponse, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/without_optional_params.tf"),
					Check: baseChecker.
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("certificate_status").
						CheckMissing("expiring_in_days").
						CheckMissing("domain").
						CheckMissing("certificate_name").
						CheckMissing("key_type").
						CheckMissing("issuer").
						CheckMissing("include_certificate_materials").
						CheckMissing("sort").
						CheckEqual("certificates.#", "3").
						CheckEqual("certificates.1.certificate_id", "cert2_1234").
						CheckEqual("certificates.1.certificate_name", "test_certificate2").
						CheckEqual("certificates.1.sans.#", "1").
						CheckEqual("certificates.1.sans.0", "test2.example.com").
						CheckEqual("certificates.1.subject.common_name", "test2.example.com").
						CheckEqual("certificates.1.subject.country", "US").
						CheckEqual("certificates.1.subject.organization", "Test Org2").
						CheckEqual("certificates.1.subject.state", "CA").
						CheckEqual("certificates.1.subject.locality", "Test City").
						CheckEqual("certificates.1.certificate_type", "THIRD_PARTY").
						CheckEqual("certificates.1.contract_id", "A-123").
						CheckEqual("certificates.1.key_type", "ECDSA").
						CheckEqual("certificates.1.key_size", "P-256").
						CheckEqual("certificates.1.secure_network", "ENHANCED_TLS").
						CheckEqual("certificates.1.account_id", "act_789").
						CheckEqual("certificates.1.created_date", "2024-05-01T12:00:00Z").
						CheckEqual("certificates.1.created_by", "test_user").
						CheckEqual("certificates.1.modified_date", "2024-07-01T12:00:00Z").
						CheckEqual("certificates.1.modified_by", "test_user2").
						CheckEqual("certificates.1.certificate_status", "READY_FOR_USE").
						CheckEqual("certificates.1.csr_expiration_date", "2027-05-01T00:00:00Z").
						CheckEqual("certificates.1.signed_certificate_issuer", "O=Test Org2,L=Test City,ST=CA,C=US").
						CheckEqual("certificates.1.signed_certificate_not_valid_after_date", "2027-11-15T10:30:00Z").
						CheckEqual("certificates.1.signed_certificate_not_valid_before_date", "2025-08-15T10:30:00Z").
						CheckEqual("certificates.2.certificate_id", "cert3_1234").
						CheckEqual("certificates.2.certificate_name", "test_certificate3").
						CheckEqual("certificates.2.sans.#", "1").
						CheckEqual("certificates.2.sans.0", "test3.example.com").
						CheckEqual("certificates.2.subject.common_name", "test3.example.com").
						CheckEqual("certificates.2.subject.country", "US").
						CheckEqual("certificates.2.subject.organization", "Test Org3").
						CheckEqual("certificates.2.subject.state", "CA").
						CheckEqual("certificates.2.subject.locality", "Test City").
						CheckEqual("certificates.2.certificate_type", "THIRD_PARTY").
						CheckEqual("certificates.2.contract_id", "A-123").
						CheckEqual("certificates.2.key_type", "RSA").
						CheckEqual("certificates.2.key_size", "2048").
						CheckEqual("certificates.2.secure_network", "ENHANCED_TLS").
						CheckEqual("certificates.2.account_id", "act_789").
						CheckEqual("certificates.2.created_date", "2024-12-01T12:00:00Z").
						CheckEqual("certificates.2.created_by", "test_user").
						CheckEqual("certificates.2.modified_date", "2025-01-01T12:00:00Z").
						CheckEqual("certificates.2.modified_by", "test_user2").
						CheckEqual("certificates.2.certificate_status", "CSR_READY").
						CheckEqual("certificates.2.csr_expiration_date", "2027-12-01T00:00:00Z").
						CheckMissing("certificates.2.signed_certificate_not_valid_after_date").
						CheckMissing("certificates.2.signed_certificate_not_valid_before_date").
						CheckMissing("certificates.2.signed_certificate_issuer").
						Build(),
				},
			},
		},
		"happy path with null subject": {
			init: func(m *ccm.Mock) {
				certWithNullSubject := baseResponse.Certificates[0]
				certWithNullSubject.Subject = nil
				mockListCertificates(m, ccm.ListCertificatesRequest{
					PageSize: 100,
					Page:     1,
				}, &ccm.ListCertificatesResponse{
					Certificates: []ccm.Certificate{certWithNullSubject},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/without_optional_params.tf"),
					Check: baseCheckerMissingCertMaterials.
						CheckMissing("certificates.0.subject").
						// Override missing subject fields checks because baseChecker expects them to be present
						CheckMissing("certificates.0.subject.common_name").
						CheckMissing("certificates.0.subject.country").
						CheckMissing("certificates.0.subject.organization").
						CheckMissing("certificates.0.subject.state").
						CheckMissing("certificates.0.subject.locality").
						Build(),
				},
			},
		},
		"happy path with certificate_name and certificate_status": {
			init: func(m *ccm.Mock) {
				mockListCertificates(m, ccm.ListCertificatesRequest{
					CertificateName:   "test_certificate1",
					CertificateStatus: []ccm.CertificateStatus{ccm.StatusActive, ccm.StatusReadyForUse},
					PageSize:          100,
					Page:              1,
				}, &ccm.ListCertificatesResponse{
					Certificates: []ccm.Certificate{baseResponse.Certificates[0]},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/with_certificate_name_and_status.tf"),
					Check: baseCheckerMissingCertMaterials.
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("expiring_in_days").
						CheckMissing("domain").
						CheckMissing("key_type").
						CheckMissing("issuer").
						CheckMissing("include_certificate_materials").
						CheckMissing("sort").
						CheckEqual("certificate_name", "test_certificate1").
						CheckEqual("certificate_status.#", "2").
						CheckEqual("certificate_status.0", "ACTIVE").
						CheckEqual("certificate_status.1", "READY_FOR_USE").
						Build(),
				},
			},
		},
		"happy path with contract_id, group_id and domain": {
			init: func(m *ccm.Mock) {
				mockListCertificates(m, ccm.ListCertificatesRequest{
					ContractID: "A-123",
					GroupID:    "1234",
					Domain:     "test2.example.com",
					PageSize:   100,
					Page:       1,
				}, &ccm.ListCertificatesResponse{
					Certificates: []ccm.Certificate{baseResponse.Certificates[1]},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/with_domain.tf"),
					Check: baseCheckerMissingCertMaterials.
						CheckMissing("certificate_status").
						CheckMissing("expiring_in_days").
						CheckMissing("certificate_name").
						CheckMissing("key_type").
						CheckMissing("issuer").
						CheckMissing("include_certificate_materials").
						CheckMissing("sort").
						CheckEqual("contract_id", "A-123").
						CheckEqual("group_id", "1234").
						CheckEqual("domain", "test2.example.com").
						CheckEqual("certificates.#", "1").
						CheckEqual("certificates.0.certificate_id", "cert2_1234").
						CheckEqual("certificates.0.certificate_name", "test_certificate2").
						CheckEqual("certificates.0.sans.#", "1").
						CheckEqual("certificates.0.sans.0", "test2.example.com").
						CheckEqual("certificates.0.subject.common_name", "test2.example.com").
						CheckEqual("certificates.0.subject.country", "US").
						CheckEqual("certificates.0.subject.organization", "Test Org2").
						CheckEqual("certificates.0.subject.state", "CA").
						CheckEqual("certificates.0.subject.locality", "Test City").
						CheckEqual("certificates.0.certificate_type", "THIRD_PARTY").
						CheckEqual("certificates.0.contract_id", "A-123").
						CheckEqual("certificates.0.key_type", "ECDSA").
						CheckEqual("certificates.0.key_size", "P-256").
						CheckEqual("certificates.0.secure_network", "ENHANCED_TLS").
						CheckEqual("certificates.0.account_id", "act_789").
						CheckEqual("certificates.0.created_date", "2024-05-01T12:00:00Z").
						CheckEqual("certificates.0.created_by", "test_user").
						CheckEqual("certificates.0.modified_date", "2024-07-01T12:00:00Z").
						CheckEqual("certificates.0.modified_by", "test_user2").
						CheckEqual("certificates.0.certificate_status", "READY_FOR_USE").
						CheckEqual("certificates.0.csr_expiration_date", "2027-05-01T00:00:00Z").
						CheckEqual("certificates.0.signed_certificate_issuer", "O=Test Org2,L=Test City,ST=CA,C=US").
						CheckEqual("certificates.0.signed_certificate_not_valid_after_date", "2027-11-15T10:30:00Z").
						CheckEqual("certificates.0.signed_certificate_not_valid_before_date", "2025-08-15T10:30:00Z").
						Build(),
				},
			},
		},
		"happy path with issuer and include_certificate_materials set to true": {
			init: func(m *ccm.Mock) {
				filteredCertificate := baseResponse.Certificates[0]
				filteredCertificate.CSRPEM = ptr.To("-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n")
				filteredCertificate.SignedCertificatePEM = ptr.To(testSignedCertificatePEM)
				filteredCertificate.SignedCertificateSerialNumber = ptr.To("12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF")
				filteredCertificate.SignedCertificateSHA256Fingerprint = ptr.To("FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10")
				filteredCertificate.TrustChainPEM = ptr.To(testTrustChainPEM)
				mockListCertificates(m, ccm.ListCertificatesRequest{
					IncludeCertificateMaterials: true,
					Issuer:                      "Test Org1",
					PageSize:                    100,
					Page:                        1,
				}, &ccm.ListCertificatesResponse{
					Certificates: []ccm.Certificate{filteredCertificate},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/with_include_certificate_materials_and_issuer.tf"),
					Check: baseChecker.
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("certificate_status").
						CheckMissing("expiring_in_days").
						CheckMissing("domain").
						CheckMissing("certificate_name").
						CheckMissing("key_type").
						CheckMissing("sort").
						CheckEqual("issuer", "Test Org1").
						CheckEqual("include_certificate_materials", "true").
						CheckEqual("certificates.0.csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nTEST-CSR-PEM\n-----END CERTIFICATE REQUEST-----\n").
						CheckEqual("certificates.0.signed_certificate_pem", testSignedCertificatePEM).
						CheckEqual("certificates.0.signed_certificate_serial_number", "12:34:56:78:90:AB:CD:EF:12:34:56:78:90:AB:CD:EF").
						CheckEqual("certificates.0.signed_certificate_sha256_fingerprint", "FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10:FE:DC:BA:98:76:54:32:10").
						CheckEqual("certificates.0.trust_chain_pem", testTrustChainPEM).
						Build(),
				},
			},
		},
		"happy path with expiring_in_days set to positive value": {
			init: func(m *ccm.Mock) {
				// Set dates to ensure the certificate is expiring within 980 days from modified date
				mockListCertificates(m, ccm.ListCertificatesRequest{
					ExpiringInDays: ptr.To[int64](980),
					PageSize:       100,
					Page:           1,
				}, &ccm.ListCertificatesResponse{
					Certificates: []ccm.Certificate{baseResponse.Certificates[0]},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/with_expiring_in_days_positive_value.tf"),
					Check: baseCheckerMissingCertMaterials.
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("certificate_status").
						CheckMissing("domain").
						CheckMissing("certificate_name").
						CheckMissing("key_type").
						CheckMissing("issuer").
						CheckMissing("include_certificate_materials").
						CheckMissing("sort").
						CheckEqual("expiring_in_days", "980").
						Build(),
				},
			},
		},
		"happy path with expiring_in_days set to 0 value": {
			init: func(m *ccm.Mock) {
				expiredCertificate := baseResponse.Certificates[0]
				expiredCertificate.CreatedDate = tst.NewTimeFromStringMust("2020-01-01T12:00:00Z")
				expiredCertificate.ModifiedDate = tst.NewTimeFromStringMust("2020-05-01T12:00:00Z")
				expiredCertificate.CSRExpirationDate = tst.NewTimeFromStringMust("2023-01-01T12:00:00Z")
				mockListCertificates(m, ccm.ListCertificatesRequest{
					ExpiringInDays: ptr.To[int64](0),
					PageSize:       100,
					Page:           1,
				}, &ccm.ListCertificatesResponse{
					Certificates: []ccm.Certificate{expiredCertificate},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/with_expiring_in_days_set_0_value.tf"),
					Check: baseCheckerMissingCertMaterials.
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("certificate_status").
						CheckMissing("domain").
						CheckMissing("certificate_name").
						CheckMissing("key_type").
						CheckMissing("issuer").
						CheckMissing("include_certificate_materials").
						CheckMissing("sort").
						CheckEqual("expiring_in_days", "0").
						CheckEqual("certificates.0.created_date", "2020-01-01T12:00:00Z").
						CheckEqual("certificates.0.modified_date", "2020-05-01T12:00:00Z").
						CheckEqual("certificates.0.csr_expiration_date", "2023-01-01T12:00:00Z").
						Build(),
				},
			},
		},
		"happy path with key_type and and sorting by createdDate": {
			init: func(m *ccm.Mock) {
				mockListCertificates(m, ccm.ListCertificatesRequest{
					KeyType:  ccm.CryptographicAlgorithmRSA,
					Sort:     "-createdDate",
					PageSize: 100,
					Page:     1,
				}, &ccm.ListCertificatesResponse{
					Certificates: []ccm.Certificate{baseResponse.Certificates[2], baseResponse.Certificates[0]},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/with_key_type_and_sort.tf"),
					Check: baseCheckerMissingCertMaterials.
						CheckMissing("contract_id").
						CheckMissing("group_id").
						CheckMissing("certificate_status").
						CheckMissing("expiring_in_days").
						CheckMissing("domain").
						CheckMissing("certificate_name").
						CheckMissing("issuer").
						CheckMissing("include_certificate_materials").
						CheckEqual("key_type", "RSA").
						CheckEqual("sort", "-createdDate").
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_id", "cert3_1234").
						CheckEqual("certificates.0.certificate_name", "test_certificate3").
						CheckEqual("certificates.0.sans.#", "1").
						CheckEqual("certificates.0.sans.0", "test3.example.com").
						CheckEqual("certificates.0.subject.common_name", "test3.example.com").
						CheckEqual("certificates.0.subject.country", "US").
						CheckEqual("certificates.0.subject.organization", "Test Org3").
						CheckEqual("certificates.0.subject.state", "CA").
						CheckEqual("certificates.0.subject.locality", "Test City").
						CheckEqual("certificates.0.certificate_type", "THIRD_PARTY").
						CheckEqual("certificates.0.contract_id", "A-123").
						CheckEqual("certificates.0.key_type", "RSA").
						CheckEqual("certificates.0.key_size", "2048").
						CheckEqual("certificates.0.secure_network", "ENHANCED_TLS").
						CheckEqual("certificates.0.contract_id", "A-123").
						CheckEqual("certificates.0.account_id", "act_789").
						CheckEqual("certificates.0.created_date", "2024-12-01T12:00:00Z").
						CheckEqual("certificates.0.created_by", "test_user").
						CheckEqual("certificates.0.modified_date", "2025-01-01T12:00:00Z").
						CheckEqual("certificates.0.modified_by", "test_user2").
						CheckEqual("certificates.0.certificate_status", "CSR_READY").
						CheckEqual("certificates.0.csr_expiration_date", "2027-12-01T00:00:00Z").
						CheckMissing("certificates.0.signed_certificate_not_valid_after_date").
						CheckMissing("certificates.0.signed_certificate_not_valid_before_date").
						CheckMissing("certificates.0.signed_certificate_issuer").
						CheckEqual("certificates.1.certificate_id", "cert1_1234").
						CheckEqual("certificates.1.certificate_name", "test_certificate1").
						CheckEqual("certificates.1.sans.#", "1").
						CheckEqual("certificates.1.sans.0", "test1.example.com").
						CheckEqual("certificates.1.subject.common_name", "test1.example.com").
						CheckEqual("certificates.1.subject.country", "US").
						CheckEqual("certificates.1.subject.organization", "Test Org1").
						CheckEqual("certificates.1.subject.state", "CA").
						CheckEqual("certificates.1.subject.locality", "Test City").
						CheckEqual("certificates.1.certificate_type", "THIRD_PARTY").
						CheckEqual("certificates.1.contract_id", "A-123").
						CheckEqual("certificates.1.key_type", "RSA").
						CheckEqual("certificates.1.key_size", "2048").
						CheckEqual("certificates.1.secure_network", "ENHANCED_TLS").
						CheckEqual("certificates.1.account_id", "act_789").
						CheckEqual("certificates.1.created_date", "2024-01-01T12:00:00Z").
						CheckEqual("certificates.1.created_by", "test_user").
						CheckEqual("certificates.1.modified_date", "2024-05-01T12:00:00Z").
						CheckEqual("certificates.1.modified_by", "test_user2").
						CheckEqual("certificates.1.certificate_status", "ACTIVE").
						CheckEqual("certificates.1.csr_expiration_date", "2027-01-01T00:00:00Z").
						CheckEqual("certificates.1.signed_certificate_issuer", "O=Test Org1,L=Test City,ST=CA,C=US").
						CheckEqual("certificates.1.signed_certificate_not_valid_after_date", "2027-12-23T08:19:47Z").
						CheckEqual("certificates.1.signed_certificate_not_valid_before_date", "2025-09-23T07:19:47Z").
						CheckMissing("certificates.1.csr_pem").
						CheckMissing("certificates.1.signed_certificate_pem").
						CheckMissing("certificates.1.signed_certificate_serial_number").
						CheckMissing("certificates.1.signed_certificate_sha256_fingerprint").
						CheckMissing("certificates.1.trust_chain_pem").
						Build(),
				},
			},
		},
		"happy path with more than 100 certificates": {
			init: func(m *ccm.Mock) {
				mockListCertificates(m, ccm.ListCertificatesRequest{
					PageSize: 100,
					Page:     1,
				}, &ccm.ListCertificatesResponse{
					Certificates: certificates101[:100],
					Links: ccm.Links{
						Self:     "/ccm/v1/certificates?page=1&pageSize=100",
						Next:     ptr.To("/ccm/v1/certificates?page=2&pageSize=100"),
						Previous: nil,
					},
				}, nil)
				mockListCertificates(m, ccm.ListCertificatesRequest{
					PageSize: 100,
					Page:     2,
				}, &ccm.ListCertificatesResponse{
					Certificates: certificates101[100:],
					Links: ccm.Links{
						Self:     "/ccm/v1/certificates?page=2&pageSize=100",
						Next:     nil,
						Previous: ptr.To("/ccm/v1/certificates?page=1&pageSize=100"),
					},
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/without_optional_params.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_certificates.test").
						CheckEqual("certificates.#", "101").
						CheckEqual("certificates.100.certificate_id", "cert_101_1234").
						CheckEqual("certificates.100.certificate_name", "test_certificate_101").
						CheckEqual("certificates.100.sans.#", "1").
						CheckEqual("certificates.100.sans.0", "test101.example.com").
						CheckEqual("certificates.100.subject.common_name", "test101.example.com").
						CheckEqual("certificates.100.subject.country", "US").
						CheckEqual("certificates.100.subject.organization", "Test Org101").
						CheckEqual("certificates.100.subject.state", "CA").
						CheckEqual("certificates.100.subject.locality", "Test City").
						CheckEqual("certificates.100.certificate_type", "THIRD_PARTY").
						CheckEqual("certificates.100.contract_id", "A-123").
						CheckEqual("certificates.100.key_type", "RSA").
						CheckEqual("certificates.100.key_size", "2048").
						CheckEqual("certificates.100.secure_network", "ENHANCED_TLS").
						CheckEqual("certificates.100.account_id", "act_789").
						CheckEqual("certificates.100.created_date", "2024-01-01T12:00:00Z").
						CheckEqual("certificates.100.created_by", "test_user").
						CheckEqual("certificates.100.modified_date", "2024-05-01T12:00:00Z").
						CheckEqual("certificates.100.modified_by", "test_user2").
						CheckEqual("certificates.100.certificate_status", "ACTIVE").
						CheckEqual("certificates.100.csr_expiration_date", "2027-01-01T00:00:00Z").
						CheckEqual("certificates.100.signed_certificate_issuer", "O=Test Org101,L=Test City,ST=CA,C=US").
						CheckEqual("certificates.100.signed_certificate_not_valid_after_date", "2027-12-23T08:19:47Z").
						CheckEqual("certificates.100.signed_certificate_not_valid_before_date", "2025-09-23T07:19:47Z").
						Build(),
				},
			},
		},
		"error response from ListCertificates": {
			init: func(m *ccm.Mock) {
				mockListCertificates(m, ccm.ListCertificatesRequest{
					PageSize: 100,
					Page:     1,
				}, &ccm.ListCertificatesResponse{}, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCertificates/without_optional_params.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"validation error - invalid key_type": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/invalid_key_type.tf"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value Match(.|\n)*` +
						`Attribute key_type value must be one of: \["RSA" "ECDSA"], got: "INVALID-TYPE"`),
				},
			},
		},
		"validation error - invalid certificate_status": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/invalid_certificate_status.tf"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value Match(.|\n)*` +
						`Attribute certificate_status\[Value\("INVALID_STATUS"\)\] value must be one of:\n` +
						`\["ACTIVE" "READY_FOR_USE" "CSR_READY"\], got: "INVALID_STATUS"`),
				},
			},
		},
		"validation error - invalid sort": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/invalid_sort.tf"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value Match(.|\n)*` +
						`Attribute sort must be a comma-separated list of fields with optional '\+' or\n` +
						`'-' prefix. Valid fields are "certificateName", "createdDate",\n` +
						`"expirationDate", and "modifiedDate", got: !invalidSort`),
				},
			},
		},
		"validation error - empty certificate_name": {
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCertificates/empty_certificate_name.tf"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value Length(.|\n)*` +
						`Attribute certificate_name string length must be at least 1, got: 0`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &ccm.Mock{}
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

func mockListCertificates(m *ccm.Mock, req ccm.ListCertificatesRequest, certificates *ccm.ListCertificatesResponse, err error) {
	if err != nil {
		m.On("ListCertificates", testutils.MockContext, req).Return(nil, err).Once()
		return
	}
	m.On("ListCertificates", testutils.MockContext, req).Return(certificates, nil).Times(3)
}
