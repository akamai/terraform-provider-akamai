package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v10/internal/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetDataSource(t *testing.T) {
	mockGetCASet := func(m *mtlstruststore.Mock, testData caSetTestData) {
		m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
			CASetID: testData.caSetID,
		}).Return(&testData.caSetResponse, nil).Times(3)
	}

	mockGetCASetVersion := func(m *mtlstruststore.Mock, testData caSetTestData) {
		m.On("GetCASetVersion", testutils.MockContext, mtlstruststore.GetCASetVersionRequest{
			CASetID: testData.caSetID,
			Version: testData.caSetVersion,
		}).Return(&testData.caSetVersionResponse, nil).Times(3)
	}

	mockListCASets := func(m *mtlstruststore.Mock, testData caSetTestData) {
		m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
			CASetNamePrefix: testData.caSetName,
			CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
		}).Return(&mtlstruststore.ListCASetsResponse{
			CASets: testData.caSets,
		}, nil).Times(3)
	}

	t.Parallel()
	commonStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_set.test").
		CheckEqual("id", "12345").
		CheckEqual("name", "example-ca-set").
		CheckEqual("version", "1").
		CheckEqual("description", "Example CA Set").
		CheckEqual("account_id", "account-123").
		CheckEqual("created_by", "example user").
		CheckEqual("created_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("version_created_by", "example user").
		CheckEqual("version_created_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("version_modified_by", "example user").
		CheckEqual("version_modified_date", "2025-05-16T12:08:34.099457Z").
		CheckEqual("allow_insecure_sha1", "false").
		CheckEqual("version_description", "Version 1 description").
		CheckEqual("version_status", "NOT_DELETED").
		CheckEqual("status", "NOT_DELETED").
		CheckEqual("staging_version", "1").
		CheckEqual("production_version", "1").
		CheckEqual("certificates.#", "1").
		CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----...").
		CheckEqual("certificates.0.description", "Example Certificate").
		CheckEqual("certificates.0.created_by", "example user").
		CheckEqual("certificates.0.created_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("certificates.0.start_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("certificates.0.end_date", "2026-04-16T12:08:34.099457Z").
		CheckEqual("certificates.0.fingerprint", "AB:CD:EF:12:34:56:78:90").
		CheckEqual("certificates.0.issuer", "Example Issuer").
		CheckEqual("certificates.0.serial_number", "123456789").
		CheckEqual("certificates.0.signature_algorithm", "SHA256").
		CheckEqual("certificates.0.subject", "Example Subject").
		CheckMissing("removal_date").
		CheckMissing("version_removal_date")

	tests := map[string]struct {
		init     func(*mtlstruststore.Mock, caSetTestData)
		testData caSetTestData
		steps    []resource.TestStep
		error    *regexp.Regexp
	}{
		"happy path - fetch by id": {
			testData: commonTestData,
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockGetCASet(m, testData)
				mockGetCASetVersion(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASet/id.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - fetch by id and version": {
			testData: commonTestData,
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockGetCASet(m, testData)
				mockGetCASetVersion(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASet/id_version.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - fetch by name": {
			testData: commonTestData,
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASets(m, testData)
				mockGetCASet(m, testData)
				mockGetCASetVersion(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASet/name.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - fetch by name and version": {
			testData: commonTestData,
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockListCASets(m, testData)
				mockGetCASet(m, testData)
				mockGetCASetVersion(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASet/name_version.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path - fetch by id, no version for ca set": {
			testData: caSetTestData{
				caSetID:      "12345",
				caSetVersion: 0,
				caSetName:    "example-ca-set",
				caSets: []mtlstruststore.CASetResponse{
					{
						CASetID:   "12345",
						CASetName: "example-ca-set",
					},
				},
				caSetResponse: mtlstruststore.GetCASetResponse{
					CASetID:     "12345",
					CASetName:   "example-ca-set",
					Description: ptr.To("Example CA Set"),
					AccountID:   "account-123",
					CASetStatus: mtlstruststore.CASetStatusNotDeleted,
					CreatedBy:   "example user",
					CreatedDate: tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
				},
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockGetCASet(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASet/id.tf"),
					Check: commonStateChecker.
						CheckMissing("certificates.#").
						CheckMissing("certificates.0.certificate_pem").
						CheckMissing("certificates.0.description").
						CheckMissing("certificates.0.created_by").
						CheckMissing("certificates.0.created_date").
						CheckMissing("certificates.0.start_date").
						CheckMissing("certificates.0.end_date").
						CheckMissing("certificates.0.fingerprint").
						CheckMissing("certificates.0.issuer").
						CheckMissing("certificates.0.serial_number").
						CheckMissing("certificates.0.signature_algorithm").
						CheckMissing("certificates.0.subject").
						CheckMissing("version").
						CheckMissing("version_created_date").
						CheckMissing("version_created_by").
						CheckMissing("production_version").
						CheckMissing("staging_version").
						CheckMissing("version_description").
						CheckMissing("version_modified_by").
						CheckMissing("version_modified_date").
						CheckMissing("version_removal_date").
						CheckMissing("version_status").
						CheckMissing("allow_insecure_sha1").
						Build(),
				},
			},
		},
		"happy path - fetch by id with removal dates and DELETED status": {
			testData: caSetTestData{
				caSetID:      "12345",
				caSetVersion: 1,
				caSetResponse: mtlstruststore.GetCASetResponse{
					CASetID:           "12345",
					CASetName:         "example-ca-set",
					Description:       ptr.To("Example CA Set"),
					AccountID:         "account-123",
					CreatedBy:         "example user",
					CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
					CASetStatus:       mtlstruststore.CASetStatusDeleted,
					StagingVersion:    ptr.To(int64(1)),
					ProductionVersion: ptr.To(int64(1)),
					LatestVersion:     ptr.To(int64(1)),
					RemovalDate:       ptr.To(tst.NewTimeFromStringMust("2025-08-01T00:00:00Z")),
				},
				caSetVersionResponse: mtlstruststore.GetCASetVersionResponse{
					Description:        ptr.To("Version 1 description"),
					AllowInsecureSHA1:  false,
					CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
					CreatedBy:          "example user",
					ModifiedBy:         ptr.To("example user"),
					ModifiedDate:       ptr.To(tst.NewTimeFromStringMust("2025-05-16T12:08:34.099457Z")),
					CaSetVersionStatus: "DELETED",
					RemovalDate:        ptr.To(tst.NewTimeFromStringMust("2025-09-01T00:00:00Z")),
					Certificates:       commonTestData.caSetVersionResponse.Certificates,
				},
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				mockGetCASet(m, testData)
				mockGetCASetVersion(m, testData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASet/id.tf"),
					Check: commonStateChecker.
						CheckEqual("removal_date", "2025-08-01T00:00:00Z").
						CheckEqual("version_removal_date", "2025-09-01T00:00:00Z").
						CheckEqual("version_status", "DELETED").
						CheckEqual("status", "DELETED").
						Build(),
				},
			},
		},
		"expect error: CA set not found": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: "12345",
				}).Return(nil, fmt.Errorf("CA set not found")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASet/id.tf"),
					ExpectError: regexp.MustCompile("CA set not found"),
				},
			},
		},
		"expect error: Get CA set version failed": {
			testData: caSetTestData{
				caSetID:      "12345",
				caSetVersion: 1,
				caSetResponse: mtlstruststore.GetCASetResponse{
					CASetID:           "12345",
					CASetName:         "example-ca-set",
					Description:       ptr.To("Example CA Set"),
					AccountID:         "account-123",
					CreatedBy:         "example user",
					CreatedDate:       tst.NewTimeFromString(t, "2025-04-16T12:08:34.099457Z"),
					StagingVersion:    ptr.To(int64(1)),
					ProductionVersion: ptr.To(int64(1)),
					LatestVersion:     ptr.To(int64(1)),
				},
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: testData.caSetID,
				}).Return(&testData.caSetResponse, nil).Once()
				m.On("GetCASetVersion", testutils.MockContext, mtlstruststore.GetCASetVersionRequest{
					CASetID: "12345",
					Version: 1,
				}).Return(nil, fmt.Errorf("Get CA set version failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASet/id.tf"),
					ExpectError: regexp.MustCompile("Get CA set version failed"),
				},
			},
		},
		"expect error: List CA sets failed": {
			init: func(m *mtlstruststore.Mock, _ caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: "example-ca-set",
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(nil, fmt.Errorf("List CA sets failed")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASet/name.tf"),
					ExpectError: regexp.MustCompile(`(?s)List.+CA sets.+failed`),
				},
			},
		},
		"expect error: both name and id provided": {
			init: func(_ *mtlstruststore.Mock, _ caSetTestData) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASet/id_name.tf"),
					ExpectError: regexp.MustCompile(`2 attributes specified when one \(and only one\) of \[id,name] is required`),
				},
			},
		},
		"expect error: could not find ca set by name": {
			testData: caSetTestData{
				caSetName: "example-ca-set",
				caSets:    []mtlstruststore.CASetResponse{},
			},
			init: func(m *mtlstruststore.Mock, testData caSetTestData) {
				m.On("ListCASets", testutils.MockContext, mtlstruststore.ListCASetsRequest{
					CASetNamePrefix: testData.caSetName,
					CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
				}).Return(&mtlstruststore.ListCASetsResponse{
					CASets: testData.caSets,
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASet/name.tf"),
					ExpectError: regexp.MustCompile(`(?s)no CA set found with the name 'example-ca-set'`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.testData)
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

type caSetTestData struct {
	caSetID              string
	caSetVersion         int64
	caSetName            string
	caSets               []mtlstruststore.CASetResponse
	caSetActivations     []mtlstruststore.ActivateCASetVersionResponse
	caSetResponse        mtlstruststore.GetCASetResponse
	caSetVersionResponse mtlstruststore.GetCASetVersionResponse
}

var commonTestData = caSetTestData{
	caSetName:    "example-ca-set",
	caSetID:      "12345",
	caSetVersion: 1,
	caSetResponse: mtlstruststore.GetCASetResponse{
		CASetID:           "12345",
		CASetName:         "example-ca-set",
		Description:       ptr.To("Example CA Set"),
		AccountID:         "account-123",
		CreatedBy:         "example user",
		CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
		CASetStatus:       mtlstruststore.CASetStatusNotDeleted,
		StagingVersion:    ptr.To(int64(1)),
		ProductionVersion: ptr.To(int64(1)),
		LatestVersion:     ptr.To(int64(1)),
	},
	caSetVersionResponse: mtlstruststore.GetCASetVersionResponse{
		Description:        ptr.To("Version 1 description"),
		AllowInsecureSHA1:  false,
		CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
		CreatedBy:          "example user",
		ModifiedBy:         ptr.To("example user"),
		ModifiedDate:       ptr.To(tst.NewTimeFromStringMust("2025-05-16T12:08:34.099457Z")),
		CaSetVersionStatus: "NOT_DELETED",
		Certificates: []mtlstruststore.CertificateResponse{
			{
				CertificatePEM:     "-----BEGIN CERTIFICATE-----...",
				Description:        ptr.To("Example Certificate"),
				CreatedBy:          "example user",
				CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
				StartDate:          tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
				EndDate:            tst.NewTimeFromStringMust("2026-04-16T12:08:34.099457Z"),
				Fingerprint:        "AB:CD:EF:12:34:56:78:90",
				Issuer:             "Example Issuer",
				SerialNumber:       "123456789",
				SignatureAlgorithm: "SHA256",
				Subject:            "Example Subject",
			},
		},
	},
	caSets: []mtlstruststore.CASetResponse{
		{
			CASetID:   "12345",
			CASetName: "example-ca-set",
		},
	},
}
