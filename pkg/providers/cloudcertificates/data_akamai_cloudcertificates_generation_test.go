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
)

func TestGenerationDataSource(t *testing.T) {
	t.Parallel()

	generationReq := cloudcertificates.GetGenerationRequest{
		LineageID:    500005,
		GenerationID: 929,
	}

	commonStateChecker := test.NewStateChecker("data.akamai_cloudcertificates_generation.test").
		CheckEqual("lineage_id", "500005").
		CheckEqual("generation_id", "929").
		CheckEqual("generation_status", "READY_FOR_USE").
		CheckEqual("generation_created_by", "terraform-dev").
		CheckEqual("generation_created_time", "2026-07-17T12:28:59Z").
		CheckEqual("generation_modified_by", "terraform-dev").
		CheckEqual("generation_modified_time", "2026-07-17T12:50:44Z")

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - generation with signed certificate": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetGeneration", testutils.MockContext, generationReq).Return(&cloudcertificates.GetGenerationResponse{
					Generation: cloudcertificates.Generation{
						Algorithms: []cloudcertificates.Algorithm{
							{
								AlgorithmInstanceCreatedBy:          "terraform-dev",
								AlgorithmInstanceCreatedTime:        ptr.To(time.Date(2026, 7, 17, 12, 28, 59, 0, time.UTC)),
								AlgorithmInstanceID:                 1714,
								AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
								AlgorithmInstanceModifiedTime:       ptr.To(time.Date(2026, 7, 17, 12, 50, 43, 0, time.UTC)),
								CertificateStatus:                   "READY_FOR_USE",
								CSRExpirationDate:                   ptr.To(time.Date(2027, 7, 17, 12, 28, 59, 0, time.UTC)),
								CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
								KeyType:                             "RSA",
								SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
								SignedCertificateNotValidAfterDate:  ptr.To(time.Date(2027, 7, 17, 12, 50, 4, 0, time.UTC)),
								SignedCertificateNotValidBeforeDate: ptr.To(time.Date(2026, 7, 17, 12, 50, 4, 0, time.UTC)),
								SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"),
								SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:03"),
								SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:03"),
								TrustChainPEM:                       ptr.To("-----BEGIN CERTIFICATE-----\nRSA-TRUST-CHAIN\n-----END CERTIFICATE-----\n"),
							},
						},
						GenerationCreatedBy:    ptr.To("terraform-dev"),
						GenerationCreatedTime:  ptr.To(time.Date(2026, 7, 17, 12, 28, 59, 0, time.UTC)),
						GenerationModifiedBy:   ptr.To("terraform-dev"),
						GenerationModifiedTime: ptr.To(time.Date(2026, 7, 17, 12, 50, 44, 0, time.UTC)),
					},
					GenerationID:     929,
					GenerationStatus: "READY_FOR_USE",
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataGeneration/basic.tf"),
					Check: commonStateChecker.
						CheckMissing("first_promoted_to_production_time").
						CheckEqual("algorithms.%", "1").
						CheckEqual("algorithms.RSA.algorithm_instance_id", "1714").
						CheckEqual("algorithms.RSA.algorithm_instance_created_by", "terraform-dev").
						CheckEqual("algorithms.RSA.algorithm_instance_created_time", "2026-07-17T12:28:59Z").
						CheckEqual("algorithms.RSA.algorithm_instance_modified_by", "terraform-dev").
						CheckEqual("algorithms.RSA.algorithm_instance_modified_time", "2026-07-17T12:50:43Z").
						CheckEqual("algorithms.RSA.certificate_status", "READY_FOR_USE").
						CheckEqual("algorithms.RSA.csr_expiration_date", "2027-07-17T12:28:59Z").
						CheckEqual("algorithms.RSA.csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n").
						CheckEqual("algorithms.RSA.signed_certificate_issuer", "CN=Test Certificate Authority").
						CheckEqual("algorithms.RSA.signed_certificate_not_valid_after_date", "2027-07-17T12:50:04Z").
						CheckEqual("algorithms.RSA.signed_certificate_not_valid_before_date", "2026-07-17T12:50:04Z").
						CheckEqual("algorithms.RSA.signed_certificate_pem", "-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n").
						CheckEqual("algorithms.RSA.signed_certificate_serial_number", "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:03").
						CheckEqual("algorithms.RSA.signed_certificate_sha256_fingerprint", "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:03").
						CheckEqual("algorithms.RSA.trust_chain_pem", "-----BEGIN CERTIFICATE-----\nRSA-TRUST-CHAIN\n-----END CERTIFICATE-----\n").
						Build(),
				},
			},
		},
		"happy path - generation with pending CSR algorithms and optional fields absent": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetGeneration", testutils.MockContext, cloudcertificates.GetGenerationRequest{
					LineageID:    500014,
					GenerationID: 3218,
				}).Return(&cloudcertificates.GetGenerationResponse{
					Generation: cloudcertificates.Generation{
						Algorithms: []cloudcertificates.Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(time.Date(2026, 7, 27, 11, 26, 38, 0, time.UTC)),
								AlgorithmInstanceID:          4632,
								CertificateStatus:            "CSR_READY",
								CSRExpirationDate:            ptr.To(time.Date(2027, 9, 28, 11, 26, 38, 0, time.UTC)),
								CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
								KeyType:                      "ECDSA",
							},
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(time.Date(2026, 7, 27, 11, 26, 38, 0, time.UTC)),
								AlgorithmInstanceID:          4631,
								CertificateStatus:            "CSR_READY",
								CSRExpirationDate:            ptr.To(time.Date(2027, 9, 28, 11, 26, 38, 0, time.UTC)),
								CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
								KeyType:                      "RSA",
							},
						},
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(time.Date(2026, 7, 27, 11, 26, 38, 0, time.UTC)),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
					GenerationID:     3218,
					GenerationStatus: "CSR_READY",
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataGeneration/csr_ready.tf"),
					Check: test.NewStateChecker("data.akamai_cloudcertificates_generation.test").
						CheckEqual("lineage_id", "500014").
						CheckEqual("generation_id", "3218").
						CheckEqual("generation_status", "CSR_READY").
						CheckEqual("generation_created_by", "terraform-dev").
						CheckEqual("generation_created_time", "2026-07-27T11:26:38Z").
						CheckEqual("generation_modified_by", "terraform-dev").
						CheckMissing("generation_modified_time").
						CheckMissing("first_promoted_to_production_time").
						CheckEqual("algorithms.%", "2").
						CheckEqual("algorithms.ECDSA.algorithm_instance_id", "4632").
						CheckEqual("algorithms.ECDSA.algorithm_instance_created_by", "terraform-dev").
						CheckEqual("algorithms.ECDSA.algorithm_instance_created_time", "2026-07-27T11:26:38Z").
						CheckEqual("algorithms.ECDSA.csr_expiration_date", "2027-09-28T11:26:38Z").
						CheckEqual("algorithms.ECDSA.csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n").
						CheckEqual("algorithms.ECDSA.certificate_status", "CSR_READY").
						CheckMissing("algorithms.ECDSA.algorithm_instance_modified_by").
						CheckMissing("algorithms.ECDSA.algorithm_instance_modified_time").
						CheckMissing("algorithms.ECDSA.signed_certificate_issuer").
						CheckMissing("algorithms.ECDSA.signed_certificate_not_valid_after_date").
						CheckMissing("algorithms.ECDSA.signed_certificate_not_valid_before_date").
						CheckMissing("algorithms.ECDSA.signed_certificate_pem").
						CheckMissing("algorithms.ECDSA.signed_certificate_serial_number").
						CheckMissing("algorithms.ECDSA.signed_certificate_sha256_fingerprint").
						CheckMissing("algorithms.ECDSA.trust_chain_pem").
						CheckEqual("algorithms.RSA.algorithm_instance_id", "4631").
						CheckEqual("algorithms.RSA.algorithm_instance_created_by", "terraform-dev").
						CheckEqual("algorithms.RSA.algorithm_instance_created_time", "2026-07-27T11:26:38Z").
						CheckEqual("algorithms.RSA.csr_expiration_date", "2027-09-28T11:26:38Z").
						CheckEqual("algorithms.RSA.csr_pem", "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n").
						CheckEqual("algorithms.RSA.certificate_status", "CSR_READY").
						CheckMissing("algorithms.RSA.signed_certificate_issuer").
						CheckMissing("algorithms.RSA.signed_certificate_not_valid_after_date").
						CheckMissing("algorithms.RSA.signed_certificate_not_valid_before_date").
						CheckMissing("algorithms.RSA.signed_certificate_pem").
						CheckMissing("algorithms.RSA.signed_certificate_serial_number").
						CheckMissing("algorithms.RSA.signed_certificate_sha256_fingerprint").
						CheckMissing("algorithms.RSA.trust_chain_pem").
						Build(),
				},
			},
		},
		"happy path - generation with first_promoted_to_production_time set": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetGeneration", testutils.MockContext, generationReq).Return(&cloudcertificates.GetGenerationResponse{
					Generation: cloudcertificates.Generation{
						Algorithms: []cloudcertificates.Algorithm{
							{
								AlgorithmInstanceCreatedBy:   "terraform-dev",
								AlgorithmInstanceCreatedTime: ptr.To(time.Date(2026, 7, 17, 12, 28, 59, 0, time.UTC)),
								AlgorithmInstanceID:          1714,
								CertificateStatus:            "READY_FOR_USE",
								CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
								KeyType:                      "RSA",
							},
						},
						FirstPromotedToProductionTime: ptr.To(time.Date(2026, 8, 19, 11, 47, 44, 0, time.UTC)),
						GenerationCreatedBy:           ptr.To("terraform-dev"),
						GenerationCreatedTime:         ptr.To(time.Date(2026, 7, 17, 12, 28, 59, 0, time.UTC)),
						GenerationModifiedBy:          ptr.To("terraform-dev"),
						GenerationModifiedTime:        ptr.To(time.Date(2026, 7, 17, 12, 50, 44, 0, time.UTC)),
					},
					GenerationID:     929,
					GenerationStatus: "READY_FOR_USE",
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataGeneration/basic.tf"),
					Check: commonStateChecker.
						CheckEqual("first_promoted_to_production_time", "2026-08-19T11:47:44Z").
						CheckEqual("algorithms.%", "1").
						Build(),
				},
			},
		},
		"generation not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetGeneration", testutils.MockContext, generationReq).
					Return(nil, cloudcertificates.ErrGenerationNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataGeneration/basic.tf"),
					ExpectError: regexp.MustCompile("No generation found with ID 929 for lineage 500005"),
				},
			},
		},
		"lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetGeneration", testutils.MockContext, generationReq).
					Return(nil, cloudcertificates.ErrLineageNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataGeneration/basic.tf"),
					ExpectError: regexp.MustCompile("No certificate lineage found with ID 500005"),
				},
			},
		},
		"internal server error": {
			init: func(m *cloudcertificates.Mock) {
				m.On("GetGeneration", testutils.MockContext, generationReq).
					Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrGetGeneration, &cloudcertificates.Error{
						Type:     "/error-types/internal-error",
						Title:    "An unexpected error occurred.",
						Status:   http.StatusInternalServerError,
						Instance: "/error-types/internal-error?traceId=1234567891081",
					})).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataGeneration/basic.tf"),
					ExpectError: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Failed to retrieve generation") + `.*` +
						regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
						regexp.QuoteMeta(`"status": 500`)),
				},
			},
		},
		"validation error - lineage_id missing": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataGeneration/no_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "lineage_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - generation_id missing": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataGeneration/no_generation_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "generation_id" is required, but no definition was found.`),
				},
			},
		},
		"validation error - lineage_id is 0": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataGeneration/zero_lineage_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute lineage_id value must be at least 1, got: 0`),
				},
			},
		},
		"validation error - generation_id is 0": {
			init: func(_ *cloudcertificates.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataGeneration/zero_generation_id.tf"),
					ExpectError: regexp.MustCompile(`Attribute generation_id value must be at least 1, got: 0`),
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
