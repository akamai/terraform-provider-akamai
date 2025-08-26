package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlstruststore"
	tst "github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetVersionsDataSource(t *testing.T) {
	mockListCASetVersions := func(m *mtlstruststore.Mock, caSetID string, includeCertificates, activeVersionsOnly bool, versionsResponse mtlstruststore.ListCASetVersionsResponse) {
		m.On("ListCASetVersions", testutils.MockContext, mtlstruststore.ListCASetVersionsRequest{
			CASetID:             caSetID,
			IncludeCertificates: includeCertificates,
			ActiveVersionsOnly:  activeVersionsOnly,
		}).Return(&versionsResponse, nil).Times(3)
	}

	mockListCASets := func(m *mtlstruststore.Mock, testData caSetTestData) {
		m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
			CASetNamePrefix: testData.caSetName,
		}).Return(&mtlstruststore.ListCASetsResponse{
			CASets: testData.caSets,
		}, nil).Times(3)
	}
	t.Parallel()
	baseChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set_versions.test").
		CheckEqual("id", "12345").
		CheckEqual("name", "test-ca-set-name").
		CheckEqual("include_certificates", "true").
		CheckEqual("active_versions_only", "false").
		CheckEqual("versions.0.version", "1").
		CheckEqual("versions.0.version_description", "test-description-two-active").
		CheckEqual("versions.0.allow_insecure_sha1", "true").
		CheckEqual("versions.0.staging_status", "ACTIVE").
		CheckEqual("versions.0.production_status", "ACTIVE").
		CheckEqual("versions.0.created_by", "jkowalski").
		CheckEqual("versions.0.created_date", "2024-04-16T12:08:34.099457Z").
		CheckEqual("versions.0.modified_by", "jkowalski").
		CheckEqual("versions.0.modified_date", "2024-04-16T12:08:34.099457Z").
		CheckEqual("versions.1.version", "2").
		CheckEqual("versions.1.version_description", "test-description-one-active").
		CheckEqual("versions.1.allow_insecure_sha1", "true").
		CheckEqual("versions.1.staging_status", "INACTIVE").
		CheckEqual("versions.1.production_status", "ACTIVE").
		CheckEqual("versions.1.created_by", "jkowalski").
		CheckEqual("versions.1.created_date", "2024-04-16T12:08:34.099457Z").
		CheckEqual("versions.1.modified_by", "jkowalski").
		CheckEqual("versions.1.modified_date", "2024-04-16T12:08:34.099457Z")

	baseResponse := mtlstruststore.ListCASetVersionsResponse{
		Versions: []mtlstruststore.CASetVersion{
			{
				CASetID:           "12345",
				CASetName:         "test-ca-set-name",
				Version:           1,
				Description:       ptr.To("test-description-two-active"),
				AllowInsecureSHA1: true,
				StagingStatus:     "ACTIVE",
				ProductionStatus:  "ACTIVE",
				CreatedBy:         "jkowalski",
				CreatedDate:       tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
				ModifiedBy:        ptr.To("jkowalski"),
				ModifiedDate:      ptr.To(tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z")),
				Certificates:      nil,
			},
			{
				CASetID:           "12345",
				CASetName:         "test-ca-set-name",
				Version:           2,
				Description:       ptr.To("test-description-one-active"),
				AllowInsecureSHA1: true,
				StagingStatus:     "INACTIVE",
				ProductionStatus:  "ACTIVE",
				CreatedBy:         "jkowalski",
				CreatedDate:       tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
				ModifiedBy:        ptr.To("jkowalski"),
				ModifiedDate:      ptr.To(tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z")),
				Certificates:      nil,
			},
			{
				CASetID:           "12345",
				CASetName:         "test-ca-set-name",
				Version:           3,
				Description:       ptr.To("test-description-two-inactive"),
				AllowInsecureSHA1: true,
				StagingStatus:     "INACTIVE",
				ProductionStatus:  "INACTIVE",
				CreatedBy:         "jkowalski",
				CreatedDate:       tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
				ModifiedBy:        ptr.To("jkowalski"),
				ModifiedDate:      ptr.To(tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z")),
				Certificates:      nil,
			},
		},
	}

	baseResponseWithCertificates := mtlstruststore.ListCASetVersionsResponse{
		Versions: []mtlstruststore.CASetVersion{
			{
				CASetID:           "12345",
				CASetName:         "test-ca-set-name",
				Version:           1,
				Description:       ptr.To("test-description-two-active"),
				AllowInsecureSHA1: true,
				StagingStatus:     "ACTIVE",
				ProductionStatus:  "ACTIVE",
				CreatedBy:         "jkowalski",
				CreatedDate:       tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
				ModifiedBy:        ptr.To("jkowalski"),
				ModifiedDate:      ptr.To(tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z")),
				Certificates: []mtlstruststore.CertificateResponse{
					{
						Subject:            "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-subject-example.com",
						Issuer:             "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-issuer-example.com",
						EndDate:            tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
						StartDate:          tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
						Fingerprint:        "test-fingerprint",
						CertificatePEM:     "-----BEGIN CERTIFICATE-----test-----END CERTIFICATE-----",
						SerialNumber:       "1234",
						SignatureAlgorithm: "SHA256WITHRSA",
						CreatedDate:        tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
						CreatedBy:          "jkowalski",
						Description:        ptr.To("test-description1"),
					},
				},
			},
			{
				CASetID:           "12345",
				CASetName:         "test-ca-set-name",
				Version:           2,
				Description:       ptr.To("test-description-one-active"),
				AllowInsecureSHA1: true,
				StagingStatus:     "INACTIVE",
				ProductionStatus:  "ACTIVE",
				CreatedBy:         "jkowalski",
				CreatedDate:       tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
				ModifiedBy:        ptr.To("jkowalski"),
				ModifiedDate:      ptr.To(tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z")),
				Certificates: []mtlstruststore.CertificateResponse{
					{
						Subject:            "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-subject-example.com",
						Issuer:             "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-issuer-example.com",
						EndDate:            tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
						StartDate:          tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
						Fingerprint:        "test-fingerprint",
						CertificatePEM:     "-----BEGIN CERTIFICATE-----test-----END CERTIFICATE-----",
						SerialNumber:       "12345",
						SignatureAlgorithm: "SHA256WITHRSA",
						CreatedDate:        tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
						CreatedBy:          "jkowalski",
						Description:        ptr.To("test-description2"),
					},
				},
			},
			{
				CASetID:           "12345",
				CASetName:         "test-ca-set-name",
				Version:           3,
				Description:       ptr.To("test-description-two-inactive"),
				AllowInsecureSHA1: true,
				StagingStatus:     "INACTIVE",
				ProductionStatus:  "INACTIVE",
				CreatedBy:         "jkowalski",
				CreatedDate:       tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
				ModifiedBy:        ptr.To("jkowalski"),
				ModifiedDate:      ptr.To(tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z")),
				Certificates: []mtlstruststore.CertificateResponse{
					{
						Subject:            "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-subject-example.com",
						Issuer:             "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-issuer-example.com",
						EndDate:            tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
						StartDate:          tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
						Fingerprint:        "test-fingerprint",
						CertificatePEM:     "-----BEGIN CERTIFICATE-----test-----END CERTIFICATE-----",
						SerialNumber:       "123456",
						SignatureAlgorithm: "SHA256WITHRSA",
						CreatedDate:        tst.NewTimeFromString(t, "2024-04-16T12:08:34.099457Z"),
						CreatedBy:          "jkowalski",
						Description:        ptr.To("test-description3"),
					},
				},
			},
		},
	}

	tests := map[string]struct {
		init  func(*mtlstruststore.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - with include_certificates set to true": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetVersions(m, "12345", true, false, baseResponseWithCertificates)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id.tf"),
					Check: baseChecker.
						CheckEqual("versions.#", "3").
						CheckEqual("versions.0.certificates.#", "1").
						CheckEqual("versions.0.certificates.0.subject", "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-subject-example.com").
						CheckEqual("versions.0.certificates.0.issuer", "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-issuer-example.com").
						CheckEqual("versions.0.certificates.0.end_date", "2025-04-16T12:08:34.099457Z").
						CheckEqual("versions.0.certificates.0.start_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.0.certificates.0.fingerprint", "test-fingerprint").
						CheckEqual("versions.0.certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----test-----END CERTIFICATE-----").
						CheckEqual("versions.0.certificates.0.serial_number", "1234").
						CheckEqual("versions.0.certificates.0.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("versions.0.certificates.0.created_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.0.certificates.0.created_by", "jkowalski").
						CheckEqual("versions.0.certificates.0.description", "test-description1").
						CheckEqual("versions.1.certificates.#", "1").
						CheckEqual("versions.1.certificates.0.subject", "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-subject-example.com").
						CheckEqual("versions.1.certificates.0.issuer", "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-issuer-example.com").
						CheckEqual("versions.1.certificates.0.end_date", "2025-04-16T12:08:34.099457Z").
						CheckEqual("versions.1.certificates.0.start_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.1.certificates.0.fingerprint", "test-fingerprint").
						CheckEqual("versions.1.certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----test-----END CERTIFICATE-----").
						CheckEqual("versions.1.certificates.0.serial_number", "12345").
						CheckEqual("versions.1.certificates.0.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("versions.1.certificates.0.created_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.1.certificates.0.created_by", "jkowalski").
						CheckEqual("versions.1.certificates.0.description", "test-description2").
						CheckEqual("versions.2.certificates.#", "1").
						CheckEqual("versions.2.certificates.0.subject", "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-subject-example.com").
						CheckEqual("versions.2.certificates.0.issuer", "C=US,ST=MA,L=Cambridge,O=Akamai,CN=test-issuer-example.com").
						CheckEqual("versions.2.certificates.0.end_date", "2025-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.certificates.0.start_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.certificates.0.fingerprint", "test-fingerprint").
						CheckEqual("versions.2.certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----test-----END CERTIFICATE-----").
						CheckEqual("versions.2.certificates.0.serial_number", "123456").
						CheckEqual("versions.2.certificates.0.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("versions.2.certificates.0.created_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.certificates.0.created_by", "jkowalski").
						CheckEqual("versions.2.certificates.0.description", "test-description3").
						Build(),
				},
			},
		},
		"happy path - fetch by name": {
			init: func(m *mtlstruststore.Mock) {
				testData := caSetTestData{
					caSetName: "test-ca-set-name",
					caSets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "12345",
							CASetName:   "test-ca-set-name",
							CASetStatus: "NOT_DELETED",
						},
					},
				}
				mockListCASets(m, testData)
				mockListCASetVersions(m, "12345", true, false, baseResponseWithCertificates)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/name.tf"),
					Check: baseChecker.
						CheckEqual("versions.#", "3").
						CheckEqual("versions.0.certificates.#", "1").
						CheckEqual("versions.1.certificates.#", "1").
						CheckEqual("versions.2.version", "3").
						CheckEqual("versions.2.version_description", "test-description-two-inactive").
						CheckEqual("versions.2.allow_insecure_sha1", "true").
						CheckEqual("versions.2.staging_status", "INACTIVE").
						CheckEqual("versions.2.production_status", "INACTIVE").
						CheckEqual("versions.2.created_by", "jkowalski").
						CheckEqual("versions.2.created_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.modified_by", "jkowalski").
						CheckEqual("versions.2.modified_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.certificates.#", "1").
						Build(),
				},
			},
		},
		"happy path - fetch by id with include_certificates set to false": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetVersions(m, "12345", false, false, baseResponse)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id_include_certificates.tf"),
					Check: baseChecker.
						CheckEqual("include_certificates", "false").
						CheckEqual("versions.#", "3").
						CheckEqual("versions.0.certificates.#", "0").
						CheckEqual("versions.1.certificates.#", "0").
						CheckEqual("versions.2.version", "3").
						CheckEqual("versions.2.version_description", "test-description-two-inactive").
						CheckEqual("versions.2.allow_insecure_sha1", "true").
						CheckEqual("versions.2.staging_status", "INACTIVE").
						CheckEqual("versions.2.production_status", "INACTIVE").
						CheckEqual("versions.2.created_by", "jkowalski").
						CheckEqual("versions.2.created_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.modified_by", "jkowalski").
						CheckEqual("versions.2.modified_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.certificates.#", "0").
						Build(),
				},
			},
		},
		"happy path - fetch by id with active_versions_only set to true and include_certificates set to false": {
			init: func(m *mtlstruststore.Mock) {
				baseResponseWithActiveVersions := mtlstruststore.ListCASetVersionsResponse{
					Versions: baseResponse.Versions[:2],
				}
				mockListCASetVersions(m, "12345", false, true, baseResponseWithActiveVersions)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id_include_certificates_and_active_versions_only.tf"),
					Check: baseChecker.
						CheckEqual("active_versions_only", "true").
						CheckEqual("include_certificates", "false").
						CheckEqual("versions.#", "2").
						CheckEqual("versions.0.certificates.#", "0").
						CheckEqual("versions.1.certificates.#", "0").
						Build(),
				},
			},
		},
		"happy path - fetch by id with active versions only and certificates": {
			init: func(m *mtlstruststore.Mock) {
				baseResponseWithActiveVersions := mtlstruststore.ListCASetVersionsResponse{
					Versions: baseResponseWithCertificates.Versions[:2],
				}
				mockListCASetVersions(m, "12345", true, true, baseResponseWithActiveVersions)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id_active_versions_only.tf"),
					Check: baseChecker.
						CheckEqual("active_versions_only", "true").
						CheckEqual("versions.#", "2").
						CheckEqual("versions.0.certificates.#", "1").
						CheckEqual("versions.1.certificates.#", "1").
						Build(),
				},
			},
		},
		"happy path - fetch by id with active versions only and certificates, but no active versions": {
			init: func(m *mtlstruststore.Mock) {
				baseResponseWithActiveVersions := mtlstruststore.ListCASetVersionsResponse{
					Versions: []mtlstruststore.CASetVersion{},
				}
				mockListCASetVersions(m, "12345", true, true, baseResponseWithActiveVersions)
				m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: "12345",
				}).Return(&mtlstruststore.GetCASetResponse{
					CASetID:   "12345",
					CASetName: "test-ca-set-name",
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id_active_versions_only.tf"),
					Check: test.NewStateChecker("data.akamai_mtlstruststore_ca_set_versions.test").
						CheckEqual("id", "12345").
						CheckEqual("name", "test-ca-set-name").
						CheckEqual("active_versions_only", "true").
						CheckEqual("versions.#", "0").
						Build(),
				},
			},
		},
		"happy path - fetch by id but no certificates found": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetVersions(m, "12345", true, false, baseResponse)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id.tf"),
					Check: baseChecker.
						CheckEqual("versions.#", "3").
						CheckEqual("versions.0.certificates.#", "0").
						CheckEqual("versions.1.certificates.#", "0").
						CheckEqual("versions.2.version", "3").
						CheckEqual("versions.2.version_description", "test-description-two-inactive").
						CheckEqual("versions.2.allow_insecure_sha1", "true").
						CheckEqual("versions.2.staging_status", "INACTIVE").
						CheckEqual("versions.2.production_status", "INACTIVE").
						CheckEqual("versions.2.created_by", "jkowalski").
						CheckEqual("versions.2.created_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.modified_by", "jkowalski").
						CheckEqual("versions.2.modified_date", "2024-04-16T12:08:34.099457Z").
						CheckEqual("versions.2.certificates.#", "0").
						Build(),
				},
			},
		},
		"error: could not find by ca set name": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "test-ca-set-name",
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{
							CASetID:     "01234",
							CASetName:   "test-ca-set-name foo",
							CASetStatus: "NOT_DELETED",
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/name.tf"),
					ExpectError: regexp.MustCompile("no CA set found with name 'test-ca-set-name'"),
				},
			},
		},
		"error: empty CA set name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/empty_name.tf"),
					ExpectError: regexp.MustCompile(`Attribute name must not be empty or only whitespace`),
				},
			},
		},
		"error: empty CA set id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/empty_id.tf"),
					ExpectError: regexp.MustCompile("Attribute id string length must be at least 1, got: 0"),
				},
			},
		},
		"error: List CA Set Versions failed": {
			init: func(m *mtlstruststore.Mock) {
				m.On("ListCASetVersions", testutils.MockContext, mtlstruststore.ListCASetVersionsRequest{
					CASetID:             "12345",
					IncludeCertificates: true,
				}).Return(nil, fmt.Errorf("failed to retrieve CA set versions")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id.tf"),
					ExpectError: regexp.MustCompile("failed to retrieve CA set versions"),
				},
			},
		},
		"validation error - missing required argument id or name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/empty.tf"),
					ExpectError: regexp.MustCompile(`No attribute specified when one \(and only one\) of \[id,name] is\s+required`),
				},
			},
		},
		"validation error - both id and name are provided": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASetVersions/id_name.tf"),
					ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of \[name,id] is\s+required`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
