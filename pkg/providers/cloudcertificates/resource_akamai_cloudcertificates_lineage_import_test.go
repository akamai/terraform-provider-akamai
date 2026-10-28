package cloudcertificates

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLineageResourceImport(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		init  func(*cloudcertificates.Mock)
		steps []resource.TestStep
	}{
		"happy path - import by lineage_id": {
			init: func(m *cloudcertificates.Mock) {
				// once for ImportState's Read, once for the follow-up PlanOnly step's refresh
				mockGetLineage(m, 500001, createLineageResponseFull(), 2)
				// cleanup after ImportStatePersist leaves the resource in state
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_lineage.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:        true,
					ImportStateId:      "500001",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("lineage_id", "500001").
						CheckEqual("account_id", "A-CCT1234").
						CheckEqual("contract_id", "C-0N7RAC7").
						CheckEqual("group_id", "12345").
						CheckEqual("geo_class", "STANDARD_WORLDWIDE").
						CheckEqual("lineage_name", "test non-validated domain").
						CheckEqual("lineage_type", "MULTIPLE_GENERATION").
						CheckEqual("secure_network", "ENHANCED_TLS").
						CheckEqual("stack_mode", "MULTIPLE_STACK").
						CheckEqual("key_specs.%", "2").
						CheckEqual("key_specs.RSA", "2048").
						CheckEqual("key_specs.ECDSA", "P-256").
						CheckEqual("sans.#", "2").
						CheckSetContains("sans", "www.example.com").
						CheckSetContains("sans", "example.com").
						CheckEqual("subject.common_name", "example.com").
						CheckEqual("subject.organization", "Example Corp.").
						CheckEqual("subject.organizational_unit", "IT").
						CheckEqual("subject.country", "US").
						CheckEqual("subject.state", "Massachusetts").
						CheckEqual("subject.locality", "Cambridge").
						CheckEqual("lineage_created_by", "terraform-dev").
						CheckEqual("lineage_created_time", "2026-07-07T10:18:46Z").
						CheckEqual("lineage_modified_by", "terraform-dev").
						CheckEqual("lineage_modified_time", "2026-07-07T10:18:46Z").
						CheckEqual("head.generation_id", "2912").
						CheckEqual("head.generation_status", "CSR_READY").
						CheckEqual("head.algorithms.%", "2").
						CheckEqual("head.algorithms.RSA.algorithm_instance_id", "6122").
						CheckMissing("head.algorithms.RSA.algorithm_instance_modified_by").
						CheckMissing("head.algorithms.RSA.signed_certificate_pem").
						CheckMissing("current_production").
						CheckMissing("previous_production").
						CheckMissing("current_staging").
						Build(),
				},
				{
					// verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					PlanOnly: true,
				},
			},
		},
		"happy path - import by lineage_id, minimal fields": {
			init: func(m *cloudcertificates.Mock) {
				// once for ImportState's Read, once for the follow-up PlanOnly step's refresh
				mockGetLineage(m, 500006, createLineageResponseMinimal(), 2)
				// cleanup after ImportStatePersist leaves the resource in state
				mockDeleteLineage(m, 500006)
			},
			steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_lineage.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_minimal.tf"),
					ImportState:        true,
					ImportStateId:      "500006",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("lineage_id", "500006").
						CheckEqual("account_id", "A-CCT1234").
						CheckEqual("stack_mode", "SINGLE_STACK").
						CheckEqual("lineage_name", "www.example.com20260720103154092771").
						CheckEqual("key_specs.%", "1").
						CheckEqual("key_specs.RSA", "2048").
						CheckEqual("sans.#", "2").
						CheckEqual("head.generation_id", "1323").
						CheckEqual("head.algorithms.%", "1").
						CheckMissing("subject.common_name").
						CheckMissing("subject.organization").
						CheckMissing("subject.organizational_unit").
						CheckMissing("subject.country").
						CheckMissing("subject.state").
						CheckMissing("subject.locality").
						// RES-10: the mapper must collapse an all-empty API subject to a null object, not one with
						// every field individually null.
						CheckMissing("subject.%").
						CheckMissing("current_production").
						CheckMissing("previous_production").
						CheckMissing("current_staging").
						Build(),
				},
				{
					// verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_minimal.tf"),
					PlanOnly: true,
				},
			},
		},
		// mirrors importing a real lineage that has been rolled back on production, so
		// previous_production (the rollback candidate) is populated alongside current_production.
		"happy path - import by lineage_id, previous_production populated": {
			init: func(m *cloudcertificates.Mock) {
				resp := createLineageResponseFull()
				currentProductionAlgorithms := resp.Head.Algorithms
				resp.Head = nil
				resp.CurrentProduction = &cloudcertificates.ProductionGeneration{
					ProductionGenerationID:     2913,
					ProductionGenerationStatus: "ACTIVE",
					Generation: cloudcertificates.Generation{
						Algorithms:            currentProductionAlgorithms,
						GenerationCreatedBy:   ptr.To("terraform-dev"),
						GenerationCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-08T09:00:00Z")),
						GenerationModifiedBy:  ptr.To("terraform-dev"),
					},
				}
				previousProductionAlgorithms := []cloudcertificates.Algorithm{
					{
						AlgorithmInstanceCreatedBy:          "terraform-dev",
						AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						AlgorithmInstanceID:                 6122,
						CertificateStatus:                   "READY_FOR_USE",
						CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                             "RSA",
						SignedCertificateIssuer:             ptr.To("CN=Previous Production CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nPREVIOUS-RSA-CERT-6122\n-----END CERTIFICATE-----\n"),
						SignedCertificateSerialNumber:       ptr.To("aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99"),
						SignedCertificateSHA256Fingerprint:  ptr.To("AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99"),
						TrustChainPEM:                       ptr.To("-----BEGIN CERTIFICATE-----\nPREVIOUS-RSA-CHAIN-6122\n-----END CERTIFICATE-----\n"),
					},
					{
						AlgorithmInstanceCreatedBy:          "terraform-dev",
						AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						AlgorithmInstanceID:                 6123,
						CertificateStatus:                   "READY_FOR_USE",
						CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                             "ECDSA",
						SignedCertificateIssuer:             ptr.To("CN=Previous Production CA"),
						SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-07T10:18:46Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nPREVIOUS-ECDSA-CERT-6123\n-----END CERTIFICATE-----\n"),
						SignedCertificateSerialNumber:       ptr.To("99:88:77:66:55:44:33:22:11:00:ff:ee:dd:cc:bb:aa"),
						SignedCertificateSHA256Fingerprint:  ptr.To("99:88:77:66:55:44:33:22:11:00:FF:EE:DD:CC:BB:AA:99:88:77:66:55:44:33:22:11:00:FF:EE:DD:CC:BB"),
						TrustChainPEM:                       ptr.To("-----BEGIN CERTIFICATE-----\nPREVIOUS-ECDSA-CHAIN-6123\n-----END CERTIFICATE-----\n"),
					},
				}
				resp.PreviousProduction = &cloudcertificates.PreviousProductionGeneration{
					PreviousProductionGenerationID:     2912,
					PreviousProductionGenerationStatus: "READY_FOR_USE",
					Generation: cloudcertificates.Generation{
						Algorithms:                    previousProductionAlgorithms,
						FirstPromotedToProductionTime: ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						GenerationCreatedBy:           ptr.To("terraform-dev"),
						GenerationCreatedTime:         ptr.To(tst.NewTimeFromStringMust("2026-07-07T10:18:46Z")),
						GenerationModifiedBy:          ptr.To("terraform-dev"),
						GenerationModifiedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-08T09:00:00Z")),
					},
				}
				// once for ImportState's Read, once for the follow-up PlanOnly step's refresh
				mockGetLineage(m, 500001, resp, 2)
				// cleanup after ImportStatePersist leaves the resource in state
				mockDeleteLineage(m, 500001)
			},
			steps: []resource.TestStep{
				{
					ResourceName:       "akamai_cloudcertificates_lineage.test",
					Config:             testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:        true,
					ImportStateId:      "500001",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("lineage_id", "500001").
						CheckMissing("head").
						CheckEqual("current_production.generation_id", "2913").
						CheckEqual("current_production.generation_status", "ACTIVE").
						CheckMissing("current_production.algorithms.RSA.signed_certificate_issuer").
						CheckMissing("current_production.algorithms.RSA.signed_certificate_pem").
						CheckEqual("previous_production.generation_id", "2912").
						CheckEqual("previous_production.generation_status", "READY_FOR_USE").
						CheckEqual("previous_production.algorithms.%", "2").
						CheckEqual("previous_production.algorithms.RSA.certificate_status", "READY_FOR_USE").
						CheckEqual("previous_production.algorithms.RSA.signed_certificate_issuer", "CN=Previous Production CA").
						CheckEqual("previous_production.algorithms.RSA.signed_certificate_not_valid_before_date", "2026-07-07T10:18:46Z").
						CheckEqual("previous_production.algorithms.RSA.signed_certificate_not_valid_after_date", "2027-07-07T10:18:46Z").
						CheckEqual("previous_production.algorithms.RSA.signed_certificate_pem", "-----BEGIN CERTIFICATE-----\nPREVIOUS-RSA-CERT-6122\n-----END CERTIFICATE-----\n").
						CheckEqual("previous_production.algorithms.RSA.signed_certificate_serial_number", "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99").
						CheckEqual("previous_production.algorithms.RSA.signed_certificate_sha256_fingerprint", "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99").
						CheckEqual("previous_production.algorithms.ECDSA.signed_certificate_issuer", "CN=Previous Production CA").
						CheckEqual("previous_production.algorithms.ECDSA.signed_certificate_serial_number", "99:88:77:66:55:44:33:22:11:00:ff:ee:dd:cc:bb:aa").
						CheckEqual("previous_production.first_promoted_to_production_time", "2026-07-07T10:18:46Z").
						CheckEqual("previous_production.generation_created_by", "terraform-dev").
						CheckEqual("previous_production.generation_created_time", "2026-07-07T10:18:46Z").
						CheckEqual("previous_production.generation_modified_by", "terraform-dev").
						CheckEqual("previous_production.generation_modified_time", "2026-07-08T09:00:00Z").
						CheckMissing("current_staging").
						Build(),
				},
				{
					// verifies the persisted imported state matches config with no drift
					Config:   testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					PlanOnly: true,
				},
			},
		},
		"error - lineage_id is zero": {
			// ImportState rejects the ID before ever calling the API, so nothing needs to be mocked.
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "0",
					ExpectError:   tst.ErrPattern(`lineage_id must be greater than 0, got: 0`),
				},
			},
		},
		"error - lineage_id is negative": {
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "-5",
					ExpectError:   tst.ErrPattern(`lineage_id must be greater than 0, got: -5`),
				},
			},
		},
		"error - malformed multipart import ID": {
			// ImportIDSplitter rejects the extra part before ever calling the API, so nothing needs to be mocked.
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "500001,extra",
					ExpectError:   tst.ErrPattern(`invalid number of importID parts: 2; you need to provide an importID in the format 'lineageID'`),
				},
			},
		},
		"error - post-import refresh finds lineage not found": {
			init: func(m *cloudcertificates.Mock) {
				mockGetLineageFails(m, 500001, cloudcertificates.ErrLineageNotFound, 1)
			},
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "500001",
					ExpectError:   tst.ErrPattern(`Cannot import non-existent remote object`),
				},
			},
		},
		"error - post-import refresh fails, generic error": {
			init: func(m *cloudcertificates.Mock) {
				mockGetLineageFails(m, 500001, cloudcertificates.ErrInternalError, 1)
			},
			steps: []resource.TestStep{
				{
					ResourceName:  "akamai_cloudcertificates_lineage.test",
					Config:        testutils.LoadFixtureString(t, "testdata/TestResourceCloudCertificatesLineage/create_full.tf"),
					ImportState:   true,
					ImportStateId: "500001",
					ExpectError:   tst.ErrPattern(`Failed to read lineage`),
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
				Steps:                    tc.steps,
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}
