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
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataLineages(t *testing.T) {
	t.Parallel()

	firstLineage := listLineageResponse(500002, 903, "www.example.com20260714095208534842", "CSR_READY")
	secondLineage := listLineageResponse(500003, 830, "www.example.com20260716073714851272", "CSR_READY")
	secondLineage.Subject = cloudcertificates.Subject{
		CommonName: "example.com", // Other fields left empty.
	}
	secondLineage.Head = nil
	expandedLineage := expandedListLineageResponse(listLineageResponse(500002, 903, "www.example.com20260714095208534842", "CSR_READY"))
	expandedSecondLineage := expandedListLineageResponse(listLineageResponse(500003, 830, "www.example.com20260716073714851272", "CSR_READY"))

	expiringInDays := 30
	filteredRequest := cloudcertificates.ListLineagesRequest{
		ContractID:    "C-0N7RAC7",
		LineageName:   "example",
		LineageIDs:    []int64{500001, 500002},
		SecureNetwork: cloudcertificates.SecureNetworkEnhancedTLS,
		StackMode:     cloudcertificates.StackModeMultipleStack,
		LineageType:   cloudcertificates.LineageTypeMultipleGeneration,
		Domain:        "example.com",
		GenerationStatus: []cloudcertificates.GenerationStatus{
			cloudcertificates.GenerationStatusActive,
			cloudcertificates.GenerationStatusReadyForUse,
		},
		ExpiringInDays: ptr.To(expiringInDays),
		KeyType:        cloudcertificates.CryptographicAlgorithmRSA,
		Issuer:         "Test Certificate Authority",
		ExpandGenerations: []cloudcertificates.ExpandGenerations{
			cloudcertificates.ExpandGenerationsHead,
			cloudcertificates.ExpandGenerationsCurrentProduction,
			cloudcertificates.ExpandGenerationsPreviousProduction,
			cloudcertificates.ExpandGenerationsCurrentStaging,
		},
		PageSize: cloudcertificates.MaxListLineagesPageSize,
		Sort:     "-modifiedDate",
	}

	tests := map[string]struct {
		config func(*testing.T) string
		init   func(*cloudcertificates.Mock)
		check  resource.TestCheckFunc
		error  *regexp.Regexp
	}{
		"happy path - without filters": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/default.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: cloudcertificates.MaxListLineagesPageSize,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{firstLineage, secondLineage},
					TotalCount: 2,
				}, nil).Times(3)
			},
			check: checkLineagesBaseAttrs(),
		},
		"happy path - filters, expansion, and pagination": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/filtered_expanded.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, filteredRequest).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{expandedLineage},
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 2,
				}, nil).Times(3)

				nextPageRequest := filteredRequest
				nextPageRequest.After = "next-cursor"
				// read
				m.On("ListLineages", testutils.MockContext, nextPageRequest).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{expandedSecondLineage},
					TotalCount: 2,
				}, nil).Times(3)
			},
			check: checkLineagesFilteredExpandedAttrs(),
		},
		"happy path - expiring in days is zero": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/zero_expiring_in_days.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				expiringInDaysZero := 0
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					ExpiringInDays: ptr.To(expiringInDaysZero),
					PageSize:       cloudcertificates.MaxListLineagesPageSize,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{firstLineage},
					TotalCount: 1,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "expiring_in_days", "0"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "1"),
			),
		},
		"happy path - no matching lineages": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/default.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: cloudcertificates.MaxListLineagesPageSize,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{},
					TotalCount: 0,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "0"),
			),
		},
		"happy path - empty string next cursor stops pagination": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/default.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: cloudcertificates.MaxListLineagesPageSize,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{firstLineage},
					NextCursor: ptr.To(""),
					TotalCount: 1,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.0.lineage_id", "500002"),
			),
		},
		"happy path - empty lineages with non-empty cursor stops pagination": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/default.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: cloudcertificates.MaxListLineagesPageSize,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{},
					NextCursor: ptr.To("stale-cursor"),
					TotalCount: 0,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "0"),
			),
		},
		"happy path - limit stops pagination once satisfied": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/with_limit.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: 1,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{firstLineage},
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 2,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "limit", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.0.lineage_id", "500002"),
			),
		},
		"happy path - maximum limit caps page size": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/max_limit.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: cloudcertificates.MaxListLineagesPageSize,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{firstLineage},
					TotalCount: 1,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "limit", "9223372036854775807"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.0.lineage_id", "500002"),
			),
		},
		"happy path - response larger than limit is truncated": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/with_limit.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: 1,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{firstLineage, secondLineage},
					TotalCount: 2,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "limit", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.0.lineage_id", "500002"),
			),
		},
		"happy path - limit spanning multiple pages caps the last page size": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/with_limit_two.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: 2,
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{firstLineage},
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 3,
				}, nil).Times(3)

				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: 1,
					After:    "next-cursor",
				}).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{secondLineage},
					TotalCount: 3,
				}, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "limit", "2"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.#", "2"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.0.lineage_id", "500002"),
				resource.TestCheckResourceAttr("data.akamai_cloudcertificates_lineages.test", "lineages.1.lineage_id", "500003"),
			),
		},
		"pagination error": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/filtered_expanded.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, filteredRequest).Return(&cloudcertificates.ListLineagesResponse{
					Lineages:   []cloudcertificates.Lineage{expandedLineage},
					NextCursor: ptr.To("next-cursor"),
					TotalCount: 2,
				}, nil).Once()

				nextPageRequest := filteredRequest
				nextPageRequest.After = "next-cursor"
				// read
				m.On("ListLineages", testutils.MockContext, nextPageRequest).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListLineages, &cloudcertificates.Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891082",
				})).Once()
			},
			error: regexp.MustCompile(`(?s)` +
				regexp.QuoteMeta("Read Cloud Certificates Lineages failed") + `.*` +
				regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
				regexp.QuoteMeta(`"title": "An unexpected error occurred."`) + `.*` +
				regexp.QuoteMeta(`"status": 500`) + `.*` +
				regexp.QuoteMeta(`"instance": "/error-types/internal-error?traceId=1234567891082"`)),
		},
		"internal server error": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/default.tf")
			},
			init: func(m *cloudcertificates.Mock) {
				// read
				m.On("ListLineages", testutils.MockContext, cloudcertificates.ListLineagesRequest{
					PageSize: cloudcertificates.MaxListLineagesPageSize,
				}).Return(nil, fmt.Errorf("%w: %w", cloudcertificates.ErrListLineages, &cloudcertificates.Error{
					Type:     "/error-types/internal-error",
					Title:    "An unexpected error occurred.",
					Status:   http.StatusInternalServerError,
					Instance: "/error-types/internal-error?traceId=1234567891081",
				})).Once()
			},
			error: regexp.MustCompile(`(?s)` +
				regexp.QuoteMeta("Read Cloud Certificates Lineages failed") + `.*` +
				regexp.QuoteMeta(`"type": "/error-types/internal-error"`) + `.*` +
				regexp.QuoteMeta(`"title": "An unexpected error occurred."`) + `.*` +
				regexp.QuoteMeta(`"status": 500`) + `.*` +
				regexp.QuoteMeta(`"instance": "/error-types/internal-error?traceId=1234567891081"`)),
		},
		"validation error - generation status": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_generation_status.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute generation_status[Value("NOT_A_STATUS")] value must be one of: ["CSR_READY" "READY_FOR_USE" "ACTIVE" "ARCHIVED" "ABANDONED"], got: "NOT_A_STATUS"`)),
		},
		"validation error - empty contract ID": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/empty_contract_id.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute contract_id string length must be at least 1, got: 0`)),
		},
		"validation error - empty lineage name": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/empty_lineage_name.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute lineage_name string length must be at least 1, got: 0`)),
		},
		"validation error - empty domain": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/empty_domain.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute domain string length must be at least 1, got: 0`)),
		},
		"validation error - empty generation status": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/empty_generation_status.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute generation_status set must contain at least 1 elements, got: 0`)),
		},
		"validation error - empty issuer": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/empty_issuer.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute issuer string length must be at least 1, got: 0`)),
		},
		"validation error - empty lineage IDs": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/empty_lineage_ids.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute lineage_ids set must contain at least 1 elements, got: 0`)),
		},
		"validation error - null lineage IDs": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/null_lineage_ids.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`This attribute contains a null value.`)),
		},
		"validation error - zero lineage ID": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/zero_lineage_id.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute lineage_ids[Value(0)] value must be at least 1, got: 0`)),
		},
		"validation error - null generation status": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/null_generation_status.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`This attribute contains a null value.`)),
		},
		"validation error - negative expiration window": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/negative_expiring_in_days.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute expiring_in_days value must be between 0 and 2147483647, got: -1`)),
		},
		"validation error - expiration window exceeds API integer range": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_expiring_in_days_upper_bound.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute expiring_in_days value must be between 0 and 2147483647, got: 2147483648`)),
		},
		"validation error - sort": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_sort.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute sort value must be one of: ["-modifiedDate" "+modifiedDate" "-createdDate" "+createdDate" "-lineageName" "+lineageName" "-expirationDate" "+expirationDate"], got: "modifiedDate"`)),
		},
		"validation error - secure network": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_secure_network.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute secure_network value must be one of: ["ENHANCED_TLS" "STANDARD_TLS"], got: "INVALID_NETWORK"`)),
		},
		"validation error - stack mode": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_stack_mode.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute stack_mode value must be one of: ["SINGLE_STACK" "MULTIPLE_STACK"], got: "INVALID_STACK_MODE"`)),
		},
		"validation error - lineage type": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_lineage_type.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute lineage_type value must be one of: ["MULTIPLE_GENERATION" "SINGLE_GENERATION"], got: "INVALID_LINEAGE_TYPE"`)),
		},
		"validation error - key type": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_key_type.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute key_type value must be one of: ["RSA" "ECDSA"], got: "INVALID_KEY_TYPE"`)),
		},
		"validation error - limit lower than 1": {
			config: func(t *testing.T) string {
				return testutils.LoadFixtureString(t, "testdata/TestDataCloudCertificatesLineages/invalid_limit.tf")
			},
			error: tst.ErrPattern(regexp.QuoteMeta(`Attribute limit value must be at least 1, got: 0`)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if test.init != nil {
				test.init(client.CloudCertificates)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      test.config(t),
						Check:       test.check,
						ExpectError: test.error,
					},
				},
			})

			client.CloudCertificates.AssertExpectations(t)
		})
	}
}

func listLineageResponse(lineageID, headGenerationID int64, lineageName, headStatus string) cloudcertificates.Lineage {
	createdTime := time.Date(2026, 7, 14, 9, 52, 9, 0, time.UTC)
	modifiedTime := time.Date(2026, 7, 17, 7, 58, 13, 0, time.UTC)

	return cloudcertificates.Lineage{
		AccountID:           "A-CCT1234",
		ContractID:          "C-0N7RAC7",
		GeoClass:            "STANDARD_WORLDWIDE",
		GroupID:             12345,
		LineageID:           lineageID,
		LineageName:         lineageName,
		LineageType:         "MULTIPLE_GENERATION",
		SecureNetwork:       "ENHANCED_TLS",
		StackMode:           "MULTIPLE_STACK",
		SANs:                []string{"www.example.com", "example.com"},
		LineageCreatedBy:    "terraform-dev",
		LineageCreatedTime:  createdTime,
		LineageModifiedBy:   "terraform-dev",
		LineageModifiedTime: modifiedTime,
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
			HeadGenerationID:     headGenerationID,
			HeadGenerationStatus: headStatus,
		},
	}
}

func expandedListLineageResponse(lineage cloudcertificates.Lineage) cloudcertificates.Lineage {
	headCreatedTime := time.Date(2026, 7, 27, 8, 38, 52, 0, time.UTC)
	headModifiedTime := time.Date(2026, 7, 27, 8, 39, 52, 0, time.UTC)
	headPromotedTime := time.Date(2026, 7, 27, 8, 40, 52, 0, time.UTC)
	lineage.Head.Generation = cloudcertificates.Generation{
		Algorithms: []cloudcertificates.Algorithm{{
			AlgorithmInstanceID:                 4544,
			AlgorithmInstanceCreatedBy:          "head-user",
			AlgorithmInstanceCreatedTime:        ptr.To(headCreatedTime),
			AlgorithmInstanceModifiedBy:         ptr.To("head-user"),
			AlgorithmInstanceModifiedTime:       ptr.To(headModifiedTime),
			CertificateStatus:                   "READY_FOR_USE",
			CSRExpirationDate:                   ptr.To(time.Date(2027, 10, 27, 8, 38, 52, 0, time.UTC)),
			CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
			KeyType:                             "RSA",
			SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
			SignedCertificateNotValidBeforeDate: ptr.To(time.Date(2026, 7, 27, 8, 40, 9, 0, time.UTC)),
			SignedCertificateNotValidAfterDate:  ptr.To(time.Date(2027, 7, 27, 8, 40, 9, 0, time.UTC)),
			SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nHEAD-RSA-CERT\n-----END CERTIFICATE-----\n"),
			SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11"),
			SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11"),
		}},
		FirstPromotedToProductionTime: ptr.To(headPromotedTime),
		GenerationCreatedBy:           ptr.To("head-user"),
		GenerationCreatedTime:         ptr.To(headCreatedTime),
		GenerationModifiedBy:          ptr.To("head-user"),
		GenerationModifiedTime:        ptr.To(headModifiedTime),
	}
	lineage.Head.HeadGenerationStatus = "CSR_READY"

	productionCreatedTime := time.Date(2026, 7, 23, 13, 44, 38, 0, time.UTC)
	productionModifiedTime := time.Date(2026, 7, 23, 13, 45, 38, 0, time.UTC)
	productionPromotedTime := time.Date(2026, 7, 23, 13, 49, 13, 0, time.UTC)
	lineage.CurrentProduction = &cloudcertificates.ProductionGeneration{
		ProductionGenerationID:     3032,
		ProductionGenerationStatus: "ACTIVE",
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{{
				AlgorithmInstanceID:                 4351,
				AlgorithmInstanceCreatedBy:          "production-user",
				AlgorithmInstanceCreatedTime:        ptr.To(productionCreatedTime),
				AlgorithmInstanceModifiedBy:         ptr.To("production-user"),
				AlgorithmInstanceModifiedTime:       ptr.To(productionModifiedTime),
				CertificateStatus:                   "READY_FOR_USE",
				CSRExpirationDate:                   ptr.To(time.Date(2027, 9, 23, 13, 44, 38, 0, time.UTC)),
				CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
				KeyType:                             "RSA",
				SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
				SignedCertificateNotValidBeforeDate: ptr.To(time.Date(2026, 7, 23, 13, 45, 9, 0, time.UTC)),
				SignedCertificateNotValidAfterDate:  ptr.To(time.Date(2027, 7, 23, 13, 45, 9, 0, time.UTC)),
				SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nPRODUCTION-RSA-CERT\n-----END CERTIFICATE-----\n"),
				SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:12"),
				SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:12"),
			}},
			FirstPromotedToProductionTime: ptr.To(productionPromotedTime),
			GenerationCreatedBy:           ptr.To("production-user"),
			GenerationCreatedTime:         ptr.To(productionCreatedTime),
			GenerationModifiedBy:          ptr.To("production-user"),
			GenerationModifiedTime:        ptr.To(productionModifiedTime),
		},
	}

	previousProductionCreatedTime := time.Date(2026, 7, 23, 11, 26, 1, 0, time.UTC)
	previousProductionModifiedTime := time.Date(2026, 7, 23, 12, 20, 52, 0, time.UTC)
	previousProductionPromotedTime := time.Date(2026, 7, 23, 12, 21, 36, 0, time.UTC)
	lineage.PreviousProduction = &cloudcertificates.PreviousProductionGeneration{
		PreviousProductionGenerationID:     3029,
		PreviousProductionGenerationStatus: "READY_FOR_USE",
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{{
				AlgorithmInstanceID:                 4345,
				AlgorithmInstanceCreatedBy:          "previous-production-user",
				AlgorithmInstanceCreatedTime:        ptr.To(previousProductionCreatedTime),
				AlgorithmInstanceModifiedBy:         ptr.To("previous-production-user"),
				AlgorithmInstanceModifiedTime:       ptr.To(previousProductionModifiedTime),
				CertificateStatus:                   "READY_FOR_USE",
				CSRExpirationDate:                   ptr.To(time.Date(2027, 9, 24, 11, 26, 1, 0, time.UTC)),
				CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
				KeyType:                             "RSA",
				SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
				SignedCertificateNotValidBeforeDate: ptr.To(time.Date(2026, 7, 23, 12, 20, 9, 0, time.UTC)),
				SignedCertificateNotValidAfterDate:  ptr.To(time.Date(2027, 7, 23, 12, 20, 9, 0, time.UTC)),
				SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nPREVIOUS-PRODUCTION-RSA-CERT\n-----END CERTIFICATE-----\n"),
				SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:13"),
				SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:13"),
			}},
			FirstPromotedToProductionTime: ptr.To(previousProductionPromotedTime),
			GenerationCreatedBy:           ptr.To("previous-production-user"),
			GenerationCreatedTime:         ptr.To(previousProductionCreatedTime),
			GenerationModifiedBy:          ptr.To("previous-production-user"),
			GenerationModifiedTime:        ptr.To(previousProductionModifiedTime),
		},
	}

	stagingCreatedTime := time.Date(2026, 7, 24, 14, 44, 38, 0, time.UTC)
	stagingModifiedTime := time.Date(2026, 7, 24, 14, 45, 38, 0, time.UTC)
	stagingPromotedTime := time.Date(2026, 7, 24, 14, 49, 13, 0, time.UTC)
	lineage.CurrentStaging = &cloudcertificates.StagingGeneration{
		StagingGenerationID:     3032,
		StagingGenerationStatus: "ACTIVE",
		Generation: cloudcertificates.Generation{
			Algorithms: []cloudcertificates.Algorithm{{
				AlgorithmInstanceID:                 4353,
				AlgorithmInstanceCreatedBy:          "staging-user",
				AlgorithmInstanceCreatedTime:        ptr.To(stagingCreatedTime),
				AlgorithmInstanceModifiedBy:         ptr.To("staging-user"),
				AlgorithmInstanceModifiedTime:       ptr.To(stagingModifiedTime),
				CertificateStatus:                   "READY_FOR_USE",
				CSRExpirationDate:                   ptr.To(time.Date(2027, 9, 25, 14, 44, 38, 0, time.UTC)),
				CSRPEM:                              "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
				KeyType:                             "RSA",
				SignedCertificateIssuer:             ptr.To("CN=Test Certificate Authority"),
				SignedCertificateNotValidBeforeDate: ptr.To(time.Date(2026, 7, 24, 14, 45, 9, 0, time.UTC)),
				SignedCertificateNotValidAfterDate:  ptr.To(time.Date(2027, 7, 24, 14, 45, 9, 0, time.UTC)),
				SignedCertificatePEM:                ptr.To("-----BEGIN CERTIFICATE-----\nSTAGING-RSA-CERT\n-----END CERTIFICATE-----\n"),
				SignedCertificateSerialNumber:       ptr.To("12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:14"),
				SignedCertificateSHA256Fingerprint:  ptr.To("12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:14"),
			}},
			FirstPromotedToProductionTime: ptr.To(stagingPromotedTime),
			GenerationCreatedBy:           ptr.To("staging-user"),
			GenerationCreatedTime:         ptr.To(stagingCreatedTime),
			GenerationModifiedBy:          ptr.To("staging-user"),
			GenerationModifiedTime:        ptr.To(stagingModifiedTime),
		},
	}

	return lineage
}

func checkLineagesBaseAttrs() resource.TestCheckFunc {
	dsName := "data.akamai_cloudcertificates_lineages.test"

	checker := test.NewStateChecker(dsName).
		CheckEqual("lineages.#", "2").
		CheckEqualBatch("lineages.0.", lineageAttributeBatch("500002", "www.example.com20260714095208534842")).
		CheckEqualBatch("lineages.0.subject.", fullSubjectAttributeBatch()).
		CheckEqualBatch("lineages.0.head.", generationPointerAttributeBatch("generation_id", "903", "generation_status", "CSR_READY")).
		CheckMissing("lineages.0.current_production").
		CheckMissing("lineages.0.previous_production").
		CheckMissing("lineages.0.current_staging").
		CheckEqualBatch("lineages.1.", lineageAttributeBatch("500003", "www.example.com20260716073714851272")).
		CheckEqual("lineages.1.subject.common_name", "example.com").
		CheckMissing("lineages.1.subject.organization").
		CheckMissing("lineages.1.subject.organizational_unit").
		CheckMissing("lineages.1.subject.country").
		CheckMissing("lineages.1.subject.state").
		CheckMissing("lineages.1.subject.locality").
		CheckMissing("lineages.1.head").
		CheckMissing("lineages.1.current_production").
		CheckMissing("lineages.1.previous_production").
		CheckMissing("lineages.1.current_staging")

	checker = checkUnexpandedGenerationFields(checker, "lineages.0.head.")
	checker = checkMissingLineagesFilters(checker)

	return checker.Build()
}

func checkLineagesFilteredExpandedAttrs() resource.TestCheckFunc {
	dsName := "data.akamai_cloudcertificates_lineages.test"

	checker := test.NewStateChecker(dsName).
		CheckEqual("contract_id", "C-0N7RAC7").
		CheckEqual("lineage_name", "example").
		CheckEqual("lineage_ids.#", "2").
		CheckEqual("secure_network", "ENHANCED_TLS").
		CheckEqual("stack_mode", "MULTIPLE_STACK").
		CheckEqual("lineage_type", "MULTIPLE_GENERATION").
		CheckEqual("domain", "example.com").
		CheckEqual("generation_status.#", "2").
		CheckEqual("expiring_in_days", "30").
		CheckEqual("key_type", "RSA").
		CheckEqual("issuer", "Test Certificate Authority").
		CheckEqual("expand_generations", "true").
		CheckEqual("sort", "-modifiedDate").
		CheckMissing("limit").
		CheckEqual("lineages.#", "2")

	checker = checkExpandedLineage(checker, "lineages.0.", "500002", "www.example.com20260714095208534842", "903", "CSR_READY")
	checker = checkExpandedLineage(checker, "lineages.1.", "500003", "www.example.com20260716073714851272", "830", "CSR_READY")

	return resource.ComposeAggregateTestCheckFunc(
		checker.Build(),
		resource.TestCheckTypeSetElemAttr(dsName, "lineage_ids.*", "500001"),
		resource.TestCheckTypeSetElemAttr(dsName, "lineage_ids.*", "500002"),
		resource.TestCheckTypeSetElemAttr(dsName, "generation_status.*", "READY_FOR_USE"),
		resource.TestCheckTypeSetElemAttr(dsName, "generation_status.*", "ACTIVE"),
	)
}

func lineageAttributeBatch(lineageID, lineageName string) test.AttributeBatch {
	return test.AttributeBatch{
		"lineage_id":            lineageID,
		"account_id":            "A-CCT1234",
		"contract_id":           "C-0N7RAC7",
		"geo_class":             "STANDARD_WORLDWIDE",
		"group_id":              "12345",
		"lineage_name":          lineageName,
		"lineage_type":          "MULTIPLE_GENERATION",
		"secure_network":        "ENHANCED_TLS",
		"stack_mode":            "MULTIPLE_STACK",
		"lineage_created_by":    "terraform-dev",
		"lineage_created_time":  "2026-07-14T09:52:09Z",
		"lineage_modified_by":   "terraform-dev",
		"lineage_modified_time": "2026-07-17T07:58:13Z",
		"sans.#":                "2",
		"sans.0":                "www.example.com",
		"sans.1":                "example.com",
		"key_specs.%":           "2",
		"key_specs.RSA":         "2048",
		"key_specs.ECDSA":       "P-256",
	}
}

func fullSubjectAttributeBatch() test.AttributeBatch {
	return test.AttributeBatch{
		"common_name":         "example.com",
		"organization":        "Example Corp.",
		"organizational_unit": "IT",
		"country":             "US",
		"state":               "Massachusetts",
		"locality":            "Cambridge",
	}
}

func generationPointerAttributeBatch(idAttribute, id, statusAttribute, status string) test.AttributeBatch {
	return test.AttributeBatch{
		idAttribute:     id,
		statusAttribute: status,
	}
}

func expandedGenerationAttributeBatch(user, createdTime, modifiedTime, promotedTime string) test.AttributeBatch {
	return test.AttributeBatch{
		"algorithms.%":                      "1",
		"first_promoted_to_production_time": promotedTime,
		"generation_created_by":             user,
		"generation_created_time":           createdTime,
		"generation_modified_by":            user,
		"generation_modified_time":          modifiedTime,
	}
}

func expandedAlgorithmAttributeBatch(
	algorithmID, user, createdTime, modifiedTime, csrExpirationDate, issuer, notValidBeforeDate, notValidAfterDate, certPEM, serialNumber, fingerprint string,
) test.AttributeBatch {
	return test.AttributeBatch{
		"algorithm_instance_id":                    algorithmID,
		"algorithm_instance_created_by":            user,
		"algorithm_instance_created_time":          createdTime,
		"algorithm_instance_modified_by":           user,
		"algorithm_instance_modified_time":         modifiedTime,
		"certificate_status":                       "READY_FOR_USE",
		"csr_expiration_date":                      csrExpirationDate,
		"csr_pem":                                  "-----BEGIN CERTIFICATE REQUEST-----\nRSA-CSR\n-----END CERTIFICATE REQUEST-----\n",
		"signed_certificate_issuer":                issuer,
		"signed_certificate_not_valid_before_date": notValidBeforeDate,
		"signed_certificate_not_valid_after_date":  notValidAfterDate,
		"signed_certificate_pem":                   certPEM,
		"signed_certificate_serial_number":         serialNumber,
		"signed_certificate_sha256_fingerprint":    fingerprint,
	}
}

func checkExpandedLineage(checker test.StateChecker, prefix, lineageID, lineageName, headGenerationID, headStatus string) test.StateChecker {
	checker = checker.
		CheckEqualBatch(prefix, lineageAttributeBatch(lineageID, lineageName)).
		CheckEqualBatch(prefix+"subject.", fullSubjectAttributeBatch())

	checker = checkExpandedGeneration(
		checker, prefix+"head.", "generation_id", headGenerationID, "generation_status", headStatus,
		"4544", "head-user", "2026-07-27T08:38:52Z", "2026-07-27T08:39:52Z", "2026-07-27T08:40:52Z",
		"2027-10-27T08:38:52Z", "CN=Test Certificate Authority", "2026-07-27T08:40:09Z", "2027-07-27T08:40:09Z", "-----BEGIN CERTIFICATE-----\nHEAD-RSA-CERT\n-----END CERTIFICATE-----\n", "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:11", "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:11",
	)
	checker = checkExpandedGeneration(
		checker, prefix+"current_production.", "generation_id", "3032", "generation_status", "ACTIVE",
		"4351", "production-user", "2026-07-23T13:44:38Z", "2026-07-23T13:45:38Z", "2026-07-23T13:49:13Z",
		"2027-09-23T13:44:38Z", "CN=Test Certificate Authority", "2026-07-23T13:45:09Z", "2027-07-23T13:45:09Z", "-----BEGIN CERTIFICATE-----\nPRODUCTION-RSA-CERT\n-----END CERTIFICATE-----\n", "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:12", "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:12",
	)
	checker = checkExpandedGeneration(
		checker, prefix+"previous_production.", "generation_id", "3029", "generation_status", "READY_FOR_USE",
		"4345", "previous-production-user", "2026-07-23T11:26:01Z", "2026-07-23T12:20:52Z", "2026-07-23T12:21:36Z",
		"2027-09-24T11:26:01Z", "CN=Test Certificate Authority", "2026-07-23T12:20:09Z", "2027-07-23T12:20:09Z", "-----BEGIN CERTIFICATE-----\nPREVIOUS-PRODUCTION-RSA-CERT\n-----END CERTIFICATE-----\n", "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:13", "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:13",
	)

	return checkExpandedGeneration(
		checker, prefix+"current_staging.", "generation_id", "3032", "generation_status", "ACTIVE",
		"4353", "staging-user", "2026-07-24T14:44:38Z", "2026-07-24T14:45:38Z", "2026-07-24T14:49:13Z",
		"2027-09-25T14:44:38Z", "CN=Test Certificate Authority", "2026-07-24T14:45:09Z", "2027-07-24T14:45:09Z", "-----BEGIN CERTIFICATE-----\nSTAGING-RSA-CERT\n-----END CERTIFICATE-----\n", "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0:12:34:56:14", "12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:F0:12:34:56:78:9A:BC:DE:14",
	)
}

func checkExpandedGeneration(
	checker test.StateChecker,
	prefix, idAttribute, id, statusAttribute, status, algorithmID, user, createdTime, modifiedTime, promotedTime string,
	csrExpirationDate, issuer, notValidBeforeDate, notValidAfterDate, certPEM, serialNumber, fingerprint string,
) test.StateChecker {
	return checker.
		CheckEqualBatch(prefix, generationPointerAttributeBatch(idAttribute, id, statusAttribute, status)).
		CheckEqualBatch(prefix, expandedGenerationAttributeBatch(user, createdTime, modifiedTime, promotedTime)).
		CheckEqualBatch(prefix+"algorithms.RSA.", expandedAlgorithmAttributeBatch(algorithmID, user, createdTime, modifiedTime, csrExpirationDate, issuer, notValidBeforeDate, notValidAfterDate, certPEM, serialNumber, fingerprint)).
		CheckMissing(prefix + "algorithms.RSA.trust_chain_pem")
}

func checkUnexpandedGenerationFields(checker test.StateChecker, prefix string) test.StateChecker {
	return checker.
		CheckMissing(prefix + "algorithms").
		CheckMissing(prefix + "first_promoted_to_production_time").
		CheckMissing(prefix + "generation_created_by").
		CheckMissing(prefix + "generation_created_time").
		CheckMissing(prefix + "generation_modified_by").
		CheckMissing(prefix + "generation_modified_time")
}

func checkMissingLineagesFilters(checker test.StateChecker) test.StateChecker {
	return checker.
		CheckMissing("contract_id").
		CheckMissing("lineage_name").
		CheckMissing("lineage_ids").
		CheckMissing("secure_network").
		CheckMissing("stack_mode").
		CheckMissing("lineage_type").
		CheckMissing("domain").
		CheckMissing("generation_status").
		CheckMissing("expiring_in_days").
		CheckMissing("key_type").
		CheckMissing("issuer").
		CheckMissing("expand_generations").
		CheckMissing("sort").
		CheckMissing("limit")
}
