package cloudcertificates

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestArchivedGenerationsDataSource(t *testing.T) {
	t.Parallel()

	sparseGeneration := cloudcertificates.ArchivedGeneration{
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{
				{
					AlgorithmInstanceID: 6123,
					CertificateStatus:   "ABANDONED",
					KeyType:             "ECDSA",
				},
				{
					AlgorithmInstanceID: 6122,
					CertificateStatus:   "ABANDONED",
					KeyType:             "RSA",
				},
			},
			GenerationCreatedBy:   ptr.To("user"),
			GenerationCreatedTime: ptr.To(time.Date(2026, 7, 30, 10, 10, 4, 0, time.UTC)),
			GenerationModifiedBy:  ptr.To("user"),
		},
		GenerationID:     4077,
		GenerationStatus: "ABANDONED",
	}
	newerSparseGeneration := sparseGeneration
	newerSparseGeneration.GenerationID = 4078
	newerSparseGeneration.GenerationCreatedTime = ptr.To(time.Date(2026, 7, 31, 10, 10, 4, 0, time.UTC))
	fullGeneration := sparseGeneration
	fullGeneration.Algorithms = append([]cloudcertificates.Algorithm(nil), sparseGeneration.Algorithms...)
	fullGeneration.Algorithms[0].AlgorithmInstanceCreatedBy = "user"
	fullGeneration.Algorithms[0].AlgorithmInstanceCreatedTime = ptr.To(time.Date(2026, 7, 30, 10, 10, 4, 0, time.UTC))
	fullGeneration.Algorithms[0].AlgorithmInstanceModifiedBy = ptr.To("ecdsa-user")
	fullGeneration.Algorithms[0].AlgorithmInstanceModifiedTime = ptr.To(time.Date(2026, 7, 30, 11, 10, 4, 0, time.UTC))
	fullGeneration.Algorithms[0].CSRExpirationDate = ptr.To(time.Date(2027, 10, 1, 10, 10, 1, 0, time.UTC))
	fullGeneration.Algorithms[0].CSRPEM = "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n"
	fullGeneration.Algorithms[0].SignedCertificateIssuer = ptr.To("CN=ECDSA Test CA")
	fullGeneration.Algorithms[0].SignedCertificateNotValidBeforeDate = ptr.To(time.Date(2026, 7, 30, 11, 0, 0, 0, time.UTC))
	fullGeneration.Algorithms[0].SignedCertificateNotValidAfterDate = ptr.To(time.Date(2027, 7, 30, 11, 0, 0, 0, time.UTC))
	fullGeneration.Algorithms[0].SignedCertificatePEM = ptr.To("-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n")
	fullGeneration.Algorithms[0].SignedCertificateSerialNumber = ptr.To("01:23:45:67")
	fullGeneration.Algorithms[0].SignedCertificateSHA256Fingerprint = ptr.To("ECDSA:SHA256:FINGERPRINT")
	fullGeneration.Algorithms[1].AlgorithmInstanceCreatedBy = "user"
	fullGeneration.Algorithms[1].AlgorithmInstanceCreatedTime = ptr.To(time.Date(2026, 7, 30, 10, 10, 4, 0, time.UTC))
	fullGeneration.Algorithms[1].AlgorithmInstanceModifiedBy = ptr.To("rsa-user")
	fullGeneration.Algorithms[1].AlgorithmInstanceModifiedTime = ptr.To(time.Date(2026, 7, 30, 12, 10, 4, 0, time.UTC))
	fullGeneration.Algorithms[1].CSRExpirationDate = ptr.To(time.Date(2027, 10, 1, 10, 10, 3, 0, time.UTC))
	fullGeneration.Algorithms[1].CSRPEM = "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n"
	fullGeneration.Algorithms[1].SignedCertificateIssuer = ptr.To("CN=RSA Test CA")
	fullGeneration.Algorithms[1].SignedCertificateNotValidBeforeDate = ptr.To(time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC))
	fullGeneration.Algorithms[1].SignedCertificateNotValidAfterDate = ptr.To(time.Date(2027, 7, 30, 12, 0, 0, 0, time.UTC))
	fullGeneration.Algorithms[1].SignedCertificatePEM = ptr.To("-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n")
	fullGeneration.Algorithms[1].SignedCertificateSerialNumber = ptr.To("89:ab:cd:ef")
	fullGeneration.Algorithms[1].SignedCertificateSHA256Fingerprint = ptr.To("RSA:SHA256:FINGERPRINT")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - sparse algorithms by default": {
			init: func(m *cloudcertificates.Mock) {
				mockListArchivedGenerations(m, false, &cloudcertificates.ListArchivedGenerationsResponse{
					Items: []cloudcertificates.ArchivedGeneration{sparseGeneration},
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/default.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_archived_generations.test").
						CheckEqual("lineage_id", "500005").
						CheckMissing("include_algorithms").
						CheckEqual("generations.#", "1").
						CheckEqual("generations.0.generation_id", "4077").
						CheckEqual("generations.0.generation_status", "ABANDONED").
						CheckEqual("generations.0.generation_created_by", "user").
						CheckEqual("generations.0.generation_created_time", "2026-07-30T10:10:04Z").
						CheckEqual("generations.0.generation_modified_by", "user").
						CheckMissing("generations.0.generation_modified_time").
						CheckMissing("generations.0.first_promoted_to_production_time").
						CheckEqual("generations.0.algorithms.%", "2").
						CheckEqual("generations.0.algorithms.ECDSA.algorithm_instance_id", "6123").
						CheckEqual("generations.0.algorithms.ECDSA.certificate_status", "ABANDONED").
						CheckMissing("generations.0.algorithms.ECDSA.algorithm_instance_created_by").
						CheckMissing("generations.0.algorithms.ECDSA.algorithm_instance_created_time").
						CheckMissing("generations.0.algorithms.ECDSA.algorithm_instance_modified_by").
						CheckMissing("generations.0.algorithms.ECDSA.algorithm_instance_modified_time").
						CheckMissing("generations.0.algorithms.ECDSA.csr_expiration_date").
						CheckMissing("generations.0.algorithms.ECDSA.csr_pem").
						CheckMissing("generations.0.algorithms.ECDSA.signed_certificate_issuer").
						CheckMissing("generations.0.algorithms.ECDSA.signed_certificate_not_valid_after_date").
						CheckMissing("generations.0.algorithms.ECDSA.signed_certificate_not_valid_before_date").
						CheckMissing("generations.0.algorithms.ECDSA.signed_certificate_pem").
						CheckMissing("generations.0.algorithms.ECDSA.signed_certificate_serial_number").
						CheckMissing("generations.0.algorithms.ECDSA.signed_certificate_sha256_fingerprint").
						CheckMissing("generations.0.algorithms.ECDSA.trust_chain_pem").
						CheckEqual("generations.0.algorithms.RSA.algorithm_instance_id", "6122").
						CheckEqual("generations.0.algorithms.RSA.certificate_status", "ABANDONED").
						CheckMissing("generations.0.algorithms.RSA.algorithm_instance_created_by").
						CheckMissing("generations.0.algorithms.RSA.algorithm_instance_created_time").
						CheckMissing("generations.0.algorithms.RSA.algorithm_instance_modified_by").
						CheckMissing("generations.0.algorithms.RSA.algorithm_instance_modified_time").
						CheckMissing("generations.0.algorithms.RSA.csr_expiration_date").
						CheckMissing("generations.0.algorithms.RSA.csr_pem").
						CheckMissing("generations.0.algorithms.RSA.signed_certificate_issuer").
						CheckMissing("generations.0.algorithms.RSA.signed_certificate_not_valid_after_date").
						CheckMissing("generations.0.algorithms.RSA.signed_certificate_not_valid_before_date").
						CheckMissing("generations.0.algorithms.RSA.signed_certificate_pem").
						CheckMissing("generations.0.algorithms.RSA.signed_certificate_serial_number").
						CheckMissing("generations.0.algorithms.RSA.signed_certificate_sha256_fingerprint").
						CheckMissing("generations.0.algorithms.RSA.trust_chain_pem").
						Build(),
				},
			},
		},
		"happy path - full algorithms included": {
			init: func(m *cloudcertificates.Mock) {
				mockListArchivedGenerations(m, true, &cloudcertificates.ListArchivedGenerationsResponse{
					Items: []cloudcertificates.ArchivedGeneration{fullGeneration},
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/include_algorithms.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_archived_generations.test").
						CheckEqual("include_algorithms", "true").
						CheckEqual("generations.0.algorithms.%", "2").
						CheckEqual("generations.0.algorithms.ECDSA.algorithm_instance_id", "6123").
						CheckEqual("generations.0.algorithms.ECDSA.algorithm_instance_created_by", "user").
						CheckEqual("generations.0.algorithms.ECDSA.algorithm_instance_created_time", "2026-07-30T10:10:04Z").
						CheckEqual("generations.0.algorithms.ECDSA.algorithm_instance_modified_by", "ecdsa-user").
						CheckEqual("generations.0.algorithms.ECDSA.algorithm_instance_modified_time", "2026-07-30T11:10:04Z").
						CheckEqual("generations.0.algorithms.ECDSA.certificate_status", "ABANDONED").
						CheckEqual("generations.0.algorithms.ECDSA.csr_expiration_date", "2027-10-01T10:10:01Z").
						CheckEqual("generations.0.algorithms.ECDSA.csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n").
						CheckEqual("generations.0.algorithms.ECDSA.signed_certificate_issuer", "CN=ECDSA Test CA").
						CheckEqual("generations.0.algorithms.ECDSA.signed_certificate_not_valid_before_date", "2026-07-30T11:00:00Z").
						CheckEqual("generations.0.algorithms.ECDSA.signed_certificate_not_valid_after_date", "2027-07-30T11:00:00Z").
						CheckEqual("generations.0.algorithms.ECDSA.signed_certificate_pem", "-----BEGIN CERTIFICATE-----\nECDSA-CERT\n-----END CERTIFICATE-----\n").
						CheckEqual("generations.0.algorithms.ECDSA.signed_certificate_serial_number", "01:23:45:67").
						CheckEqual("generations.0.algorithms.ECDSA.signed_certificate_sha256_fingerprint", "ECDSA:SHA256:FINGERPRINT").
						CheckMissing("generations.0.algorithms.ECDSA.trust_chain_pem").
						CheckEqual("generations.0.algorithms.RSA.algorithm_instance_id", "6122").
						CheckEqual("generations.0.algorithms.RSA.algorithm_instance_created_by", "user").
						CheckEqual("generations.0.algorithms.RSA.algorithm_instance_created_time", "2026-07-30T10:10:04Z").
						CheckEqual("generations.0.algorithms.RSA.algorithm_instance_modified_by", "rsa-user").
						CheckEqual("generations.0.algorithms.RSA.algorithm_instance_modified_time", "2026-07-30T12:10:04Z").
						CheckEqual("generations.0.algorithms.RSA.certificate_status", "ABANDONED").
						CheckEqual("generations.0.algorithms.RSA.csr_expiration_date", "2027-10-01T10:10:03Z").
						CheckEqual("generations.0.algorithms.RSA.csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n").
						CheckEqual("generations.0.algorithms.RSA.signed_certificate_issuer", "CN=RSA Test CA").
						CheckEqual("generations.0.algorithms.RSA.signed_certificate_not_valid_before_date", "2026-07-30T12:00:00Z").
						CheckEqual("generations.0.algorithms.RSA.signed_certificate_not_valid_after_date", "2027-07-30T12:00:00Z").
						CheckEqual("generations.0.algorithms.RSA.signed_certificate_pem", "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n").
						CheckEqual("generations.0.algorithms.RSA.signed_certificate_serial_number", "89:ab:cd:ef").
						CheckEqual("generations.0.algorithms.RSA.signed_certificate_sha256_fingerprint", "RSA:SHA256:FINGERPRINT").
						CheckMissing("generations.0.algorithms.RSA.trust_chain_pem").
						Build(),
				},
			},
		},
		"happy path - multiple archived generations ordered oldest first": {
			init: func(m *cloudcertificates.Mock) {
				mockListArchivedGenerations(m, false, &cloudcertificates.ListArchivedGenerationsResponse{
					Items: []cloudcertificates.ArchivedGeneration{sparseGeneration, newerSparseGeneration},
				}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/default.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_archived_generations.test").
						CheckEqual("generations.#", "2").
						CheckEqual("generations.0.generation_id", "4077").
						CheckEqual("generations.1.generation_id", "4078").
						Build(),
				},
			},
		},
		"happy path - no archived generations": {
			init: func(m *cloudcertificates.Mock) {
				mockListArchivedGenerations(m, false, &cloudcertificates.ListArchivedGenerationsResponse{}).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/default.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_archived_generations.test").
						CheckEqual("generations.#", "0").
						Build(),
				},
			},
		},
		"lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListArchivedGenerations", testutils.MockContext, cloudcertificates.ListArchivedGenerationsRequest{
					LineageID: 500005,
				}).Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/default.tf"),
					ExpectError: regexp.MustCompile("No certificate lineage found with ID 500005"),
				},
			},
		},
		"internal server error": {
			init: func(m *cloudcertificates.Mock) {
				m.On("ListArchivedGenerations", testutils.MockContext, cloudcertificates.ListArchivedGenerationsRequest{
					LineageID: 500005,
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListArchivedGenerations, &cloudcertificates.Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891022",
				})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/default.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve archived generations") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
						regexp.QuoteMeta(`"status": 500`)),
				},
			},
		},
		"validation error - missing lineage_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/missing_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "lineage_id" is required`),
				},
			},
		},
		"validation error - lineage_id lower than 1": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataArchivedGenerations/invalid_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute lineage_id value must be at least 1, got: 0`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.CloudCertificates)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps:                    tc.steps,
			})
			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func mockListArchivedGenerations(m *cloudcertificates.Mock, includeAlgorithms bool, response *cloudcertificates.ListArchivedGenerationsResponse) *mock.Call {
	return m.On("ListArchivedGenerations", testutils.MockContext, cloudcertificates.ListArchivedGenerationsRequest{
		LineageID:         500005,
		IncludeAlgorithms: includeAlgorithms,
	}).Return(response, nil)
}
