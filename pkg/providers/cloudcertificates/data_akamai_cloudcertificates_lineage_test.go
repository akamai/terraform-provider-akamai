package cloudcertificates

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestDataLineage(t *testing.T) {
	t.Parallel()

	baseLineage := baseLineageResponse()

	tests := map[string]struct {
		configPath  string
		init        func(*testing.T, *cloudcertificates.Mock)
		check       resource.TestCheckFunc
		stateChecks []statecheck.StateCheck
		error       *regexp.Regexp
	}{
		"happy path - without generation expansion": {
			configPath: "testdata/TestDataCloudCertificatesLineage/default.tf",
			init: func(_ *testing.T, m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID: 500002,
				}).Return(baseLineage, nil).Times(3)
			},
			check: checkLineageBaseAttrs(),
			// pins head.algorithms as null rather than an empty known list: TestCheckResourceAttr's "0" check
			// above tolerates both, so it alone wouldn't catch a regression here.
			stateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"data.akamai_cloudcertificates_lineage.test",
					tfjsonpath.New("head").AtMapKey("algorithms"),
					knownvalue.Null(),
				),
				statecheck.ExpectKnownValue(
					"data.akamai_cloudcertificates_lineage.test",
					tfjsonpath.New("signing_target").AtMapKey("algorithms"),
					knownvalue.Null(),
				),
			},
		},
		"happy path - with generation expansion": {
			configPath: "testdata/TestDataCloudCertificatesLineage/expanded.tf",
			init: func(_ *testing.T, m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID: 500022,
					ExpandGenerations: []cloudcertificates.ExpandGenerations{
						cloudcertificates.ExpandGenerationsHead,
						cloudcertificates.ExpandGenerationsCurrentProduction,
						cloudcertificates.ExpandGenerationsPreviousProduction,
						cloudcertificates.ExpandGenerationsCurrentStaging,
					},
				}).Return(expandedLineageResponse(), nil).Times(3)
			},
			check: checkLineageExpandedAttrs(),
		},
		"happy path - lineage created without a subject": {
			configPath: "testdata/TestDataCloudCertificatesLineage/default.tf",
			init: func(_ *testing.T, m *cloudcertificates.Mock) {
				resp := *baseLineageResponse()
				resp.Subject = cloudcertificates.Subject{}
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID: 500002,
				}).Return(&resp, nil).Times(3)
			},
			check: test.NewStateChecker("data.akamai_cloudcertificates_lineage.test").
				CheckMissing("subject.common_name").
				CheckMissing("subject.organization").
				CheckMissing("subject.organizational_unit").
				CheckMissing("subject.country").
				CheckMissing("subject.state").
				CheckMissing("subject.locality").
				Build(),
		},
		"lineage not found": {
			configPath: "testdata/TestDataCloudCertificatesLineage/default.tf",
			init: func(_ *testing.T, m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID: 500002,
				}).Return(nil, cloudcertificates.ErrLineageNotFound).Times(1)
			},
			error: regexp.MustCompile("No certificate lineage found with ID 500002"),
		},
		"internal server error": {
			configPath: "testdata/TestDataCloudCertificatesLineage/default.tf",
			init: func(_ *testing.T, m *cloudcertificates.Mock) {
				m.On("GetLineage", testutils.MockContext, cloudcertificates.GetLineageRequest{
					LineageID: 500002,
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrGetLineage, &cloudcertificates.Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891081",
				})).Times(1)
			},
			error: regexp.MustCompile(`(?s)` + regexp.QuoteMeta("Read Cloud Certificates Lineage failed") + `.*` +
				regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
				regexp.QuoteMeta(`"status": 500`)),
		},
		"validation error - lineage_id is 0": {
			configPath: "testdata/TestDataCloudCertificatesLineage/zero_lineage_id.tf",
			error:      regexp.MustCompile(`Attribute lineage_id value must be at least 1, got: 0`),
		},
		"validation error - lineage_id is not a number": {
			configPath: "testdata/TestDataCloudCertificatesLineage/invalid_id.tf",
			error:      regexp.MustCompile(`Inappropriate value for attribute "lineage_id": a number is required`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if test.init != nil {
				test.init(t, client.CloudCertificates)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:            testutils.LoadFixtureString(t, test.configPath),
						Check:             test.check,
						ConfigStateChecks: test.stateChecks,
						ExpectError:       test.error,
					},
				},
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func baseLineageResponse() *cloudcertificates.GetLineageResponse {
	resp := cloudcertificates.GetLineageResponse{
		AccountID:           "A-CCT1234",
		ContractID:          "C-0N7RAC7",
		GeoClass:            "STANDARD_WORLDWIDE",
		GroupID:             12345,
		LineageID:           500002,
		LineageName:         "www.example.com20260714095208534842",
		LineageType:         "MULTIPLE_GENERATION",
		SecureNetwork:       "ENHANCED_TLS",
		StackMode:           "MULTIPLE_STACK",
		SANs:                []string{"www.example.com", "example.com"},
		LineageCreatedBy:    "terraform-dev",
		LineageCreatedTime:  tst.NewTimeFromStringMust("2026-07-14T09:52:09Z"),
		LineageModifiedBy:   "terraform-dev",
		LineageModifiedTime: tst.NewTimeFromStringMust("2026-07-14T09:52:09Z"),
		Subject: cloudcertificates.Subject{
			CommonName:         "example.com",
			Organization:       "Example Corp.",
			OrganizationalUnit: "IT",
			Country:            "US",
			State:              "Massachusetts",
			Locality:           "Cambridge",
		},
		KeySpecs: []cloudcertificates.KeySpecResponse{
			{KeyType: "RSA", KeySize: "2048"},
			{KeyType: "ECDSA", KeySize: "P-256"},
		},
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     43,
			HeadGenerationStatus: "CSR_READY",
		},
	}
	return &resp
}

func expandedLineageResponse() *cloudcertificates.GetLineageResponse {
	resp := cloudcertificates.GetLineageResponse{
		AccountID:           "A-CCT1234",
		ContractID:          "C-0N7RAC7",
		GeoClass:            "STANDARD_WORLDWIDE",
		GroupID:             12345,
		LineageID:           500022,
		LineageName:         "alpha 2",
		LineageType:         "MULTIPLE_GENERATION",
		SecureNetwork:       "ENHANCED_TLS",
		StackMode:           "MULTIPLE_STACK",
		SANs:                []string{"www.example.com", "example.com"},
		LineageCreatedBy:    "terraform-dev",
		LineageCreatedTime:  tst.NewTimeFromStringMust("2026-07-22T09:21:09Z"),
		LineageModifiedBy:   "terraform-dev",
		LineageModifiedTime: tst.NewTimeFromStringMust("2026-07-23T11:26:01Z"),
		Subject: cloudcertificates.Subject{
			CommonName:         "example.com",
			Organization:       "Example Corp.",
			OrganizationalUnit: "IT",
			Country:            "US",
			State:              "Massachusetts",
			Locality:           "Cambridge",
		},
		KeySpecs: []cloudcertificates.KeySpecResponse{
			{KeyType: "RSA", KeySize: "2048"},
			{KeyType: "ECDSA", KeySize: "P-256"},
		},
		Head: &cloudcertificates.HeadGeneration{
			HeadGenerationID:     3029,
			HeadGenerationStatus: "CSR_READY",
			Generation: cloudcertificates.Generation{
				GenerationCreatedBy:   ptr.To("terraform-dev"),
				GenerationCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:26:01Z")),
				GenerationModifiedBy:  ptr.To("terraform-dev"),
				Algorithms: []cloudcertificates.Algorithm{
					{
						AlgorithmInstanceID:          4345,
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:26:01Z")),
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:26:01Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "RSA",
					},
					{
						AlgorithmInstanceID:          4346,
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:26:01Z")),
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:26:01Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "ECDSA",
					},
				},
			},
		},
		CurrentProduction: &cloudcertificates.ProductionGeneration{
			ProductionGenerationID:     3028,
			ProductionGenerationStatus: "ACTIVE",
			Generation: cloudcertificates.Generation{
				FirstPromotedToProductionTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:25:05Z")),
				GenerationCreatedBy:           ptr.To("terraform-dev"),
				GenerationCreatedTime:         ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:20:11Z")),
				GenerationModifiedBy:          ptr.To("terraform-dev"),
				GenerationModifiedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:21:42Z")),
				Algorithms: []cloudcertificates.Algorithm{
					{
						AlgorithmInstanceID:                 4343,
						AlgorithmInstanceCreatedBy:          "terraform-dev",
						AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:20:11Z")),
						AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
						AlgorithmInstanceModifiedTime:       ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:21:41Z")),
						CertificateStatus:                   "READY_FOR_USE",
						CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:20:11Z")),
						CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                             "RSA",
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-23T11:21:00Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:21:00Z")),
						SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:04"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:04"),
					},
					{
						AlgorithmInstanceID:          4344,
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:20:11Z")),
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:20:11Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "ECDSA",
					},
				},
			},
		},
		CurrentStaging: &cloudcertificates.StagingGeneration{
			StagingGenerationID:     3028,
			StagingGenerationStatus: "ACTIVE",
			Generation: cloudcertificates.Generation{
				FirstPromotedToProductionTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:25:05Z")),
				GenerationCreatedBy:           ptr.To("terraform-dev"),
				GenerationCreatedTime:         ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:20:11Z")),
				GenerationModifiedBy:          ptr.To("terraform-dev"),
				GenerationModifiedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:21:42Z")),
				Algorithms: []cloudcertificates.Algorithm{
					{
						AlgorithmInstanceID:                 4343,
						AlgorithmInstanceCreatedBy:          "terraform-dev",
						AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:20:11Z")),
						AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
						AlgorithmInstanceModifiedTime:       ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:21:41Z")),
						CertificateStatus:                   "READY_FOR_USE",
						CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:20:11Z")),
						CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                             "RSA",
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-23T11:21:00Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:21:00Z")),
						SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:04"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:04"),
					},
					{
						AlgorithmInstanceID:          4344,
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:20:11Z")),
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:20:11Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "ECDSA",
					},
				},
			},
		},
		PreviousProduction: &cloudcertificates.PreviousProductionGeneration{
			PreviousProductionGenerationID:     3027,
			PreviousProductionGenerationStatus: "READY_FOR_USE",
			Generation: cloudcertificates.Generation{
				FirstPromotedToProductionTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:18:53Z")),
				GenerationCreatedBy:           ptr.To("terraform-dev"),
				GenerationCreatedTime:         ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:13:41Z")),
				GenerationModifiedBy:          ptr.To("terraform-dev"),
				GenerationModifiedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:16:01Z")),
				Algorithms: []cloudcertificates.Algorithm{
					{
						AlgorithmInstanceID:                 4341,
						AlgorithmInstanceCreatedBy:          "terraform-dev",
						AlgorithmInstanceCreatedTime:        ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:13:41Z")),
						AlgorithmInstanceModifiedBy:         ptr.To("terraform-dev"),
						AlgorithmInstanceModifiedTime:       ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:16:00Z")),
						CertificateStatus:                   "READY_FOR_USE",
						CSRExpirationDate:                   ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:13:41Z")),
						CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                             "RSA",
						SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
						SignedCertificateNotValidAfterDate:  ptr.To(tst.NewTimeFromStringMust("2027-07-23T11:14:45Z")),
						SignedCertificateNotValidBeforeDate: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:14:45Z")),
						SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nRSA-CERT\n-----END CERTIFICATE-----\n"),
						SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:05"),
						SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:05"),
					},
					{
						AlgorithmInstanceID:          4342,
						AlgorithmInstanceCreatedBy:   "terraform-dev",
						AlgorithmInstanceCreatedTime: ptr.To(tst.NewTimeFromStringMust("2026-07-23T11:13:41Z")),
						CertificateStatus:            "CSR_READY",
						CSRExpirationDate:            ptr.To(tst.NewTimeFromStringMust("2027-09-24T11:13:41Z")),
						CSRPEM:                       "-----BEGIN CERTIFICATE REQUEST-----\nECDSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
						KeyType:                      "ECDSA",
					},
				},
			},
		},
	}
	return &resp
}

func checkLineageBaseAttrs() resource.TestCheckFunc {
	return test.NewStateChecker("data.akamai_cloudcertificates_lineage.test").
		CheckEqual("lineage_id", "500002").
		CheckEqual("account_id", "A-CCT1234").
		CheckEqual("contract_id", "C-0N7RAC7").
		CheckEqual("geo_class", "STANDARD_WORLDWIDE").
		CheckEqual("group_id", "12345").
		CheckEqual("lineage_name", "www.example.com20260714095208534842").
		CheckEqual("lineage_type", "MULTIPLE_GENERATION").
		CheckEqual("secure_network", "ENHANCED_TLS").
		CheckEqual("stack_mode", "MULTIPLE_STACK").
		CheckEqual("sans.#", "2").
		CheckEqual("sans.0", "www.example.com").
		CheckEqual("sans.1", "example.com").
		CheckEqual("subject.common_name", "example.com").
		CheckEqual("subject.organization", "Example Corp.").
		CheckEqual("subject.organizational_unit", "IT").
		CheckEqual("subject.country", "US").
		CheckEqual("subject.state", "Massachusetts").
		CheckEqual("subject.locality", "Cambridge").
		CheckEqual("key_specs.%", "2").
		CheckEqual("key_specs.RSA", "2048").
		CheckEqual("key_specs.ECDSA", "P-256").
		CheckEqual("lineage_created_by", "terraform-dev").
		CheckEqual("lineage_created_time", "2026-07-14T09:52:09Z").
		CheckEqual("lineage_modified_by", "terraform-dev").
		CheckEqual("lineage_modified_time", "2026-07-14T09:52:09Z").
		CheckEqual("head.generation_id", "43").
		CheckEqual("head.generation_status", "CSR_READY").
		CheckEqual("head.algorithms.%", "0").
		CheckEqual("signing_target.generation_id", "43").
		CheckEqual("signing_target.generation_status", "CSR_READY").
		CheckMissing("head.generation_created_by").
		CheckMissing("head.generation_created_time").
		CheckMissing("head.generation_modified_by").
		CheckMissing("head.generation_modified_time").
		CheckMissing("head.first_promoted_to_production_time").
		CheckMissing("current_production").
		CheckMissing("previous_production").
		CheckMissing("current_staging").
		Build()
}

func checkLineageExpandedAttrs() resource.TestCheckFunc {
	return test.NewStateChecker("data.akamai_cloudcertificates_lineage.test").
		CheckEqual("lineage_id", "500022").
		CheckEqual("lineage_name", "alpha 2").
		CheckEqual("head.generation_id", "3029").
		CheckEqual("head.generation_status", "CSR_READY").
		CheckEqual("head.generation_created_by", "terraform-dev").
		CheckEqual("head.generation_created_time", "2026-07-23T11:26:01Z").
		CheckEqual("head.algorithms.%", "2").
		CheckEqual("head.algorithms.RSA.algorithm_instance_id", "4345").
		CheckEqual("head.algorithms.RSA.certificate_status", "CSR_READY").
		CheckEqual("head.algorithms.ECDSA.algorithm_instance_id", "4346").
		CheckMissing("head.algorithms.RSA.algorithm_instance_modified_by").
		CheckEqual("current_production.generation_id", "3028").
		CheckEqual("current_production.generation_status", "ACTIVE").
		CheckEqual("current_production.algorithms.%", "2").
		CheckEqual("current_production.algorithms.RSA.certificate_status", "READY_FOR_USE").
		CheckEqual("current_production.algorithms.RSA.algorithm_instance_modified_by", "terraform-dev").
		CheckMissing("current_production.algorithms.ECDSA.algorithm_instance_modified_by").
		CheckEqual("current_production.algorithms.RSA.signed_certificate_issuer", "CN=Test Certificate Authority").
		CheckEqual("current_production.algorithms.RSA.signed_certificate_serial_number", "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:04").
		CheckEqual("current_production.first_promoted_to_production_time", "2026-07-23T11:25:05Z").
		CheckEqual("current_staging.generation_id", "3028").
		CheckEqual("current_staging.generation_status", "ACTIVE").
		CheckEqual("current_staging.algorithms.%", "2").
		CheckEqual("previous_production.generation_id", "3027").
		CheckEqual("previous_production.generation_status", "READY_FOR_USE").
		CheckEqual("previous_production.algorithms.%", "2").
		CheckEqual("previous_production.algorithms.RSA.signed_certificate_issuer", "CN=Test Certificate Authority").
		CheckEqual("previous_production.first_promoted_to_production_time", "2026-07-23T11:18:53Z").
		CheckEqual("signing_target.generation_id", "3029").
		CheckEqual("signing_target.generation_status", "CSR_READY").
		Build()
}
