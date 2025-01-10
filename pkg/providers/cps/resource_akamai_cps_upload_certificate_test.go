package cps

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceCPSUploadCertificate(t *testing.T) {
	tests := map[string]struct {
		init                func(*cps.Mock, *cps.GetEnrollmentResponse, int, int)
		enrollment          *cps.GetEnrollmentResponse
		enrollmentID        int
		changeID            int
		configPathForCreate string
		configPathForUpdate string
		checkFuncForCreate  resource.TestCheckFunc
		checkFuncForUpdate  resource.TestCheckFunc
		error               *regexp.Regexp
	}{
		"create with ch-mgmt false, update to true": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockUpdate(m, enrollmentID, changeID, enrollment)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadForComplete(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			error:               nil,
		},
		"create with ch-mgmt true, update to false": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				//Read
				mockReadForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)

				//Update
				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterCreate)

				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			error:               nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForCreate),
							Check:       test.checkFuncForCreate,
							ExpectError: test.error,
						},
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForUpdate),
							Check:       test.checkFuncForUpdate,
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResourceCPSUploadCertificateWithThirdPartyEnrollmentDependency(t *testing.T) {
	tests := map[string]struct {
		init         func(*cps.Mock, *cps.GetEnrollmentResponse, int, int)
		enrollment   cps.GetEnrollmentResponse
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"create third_party_enrollment, create cps_upload_certificate with enrollment_id which is third_party_enrollment's resource dependency": {
			init: func(client *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				// Create third party enrollment
				enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(*enrollment)
				client.On("CreateEnrollment",
					testutils.MockContext,
					cps.CreateEnrollmentRequest{
						EnrollmentRequestBody: enrollmentReqBody,
						ContractID:            "1",
					},
				).Return(&cps.CreateEnrollmentResponse{
					ID:         enrollmentID,
					Enrollment: fmt.Sprintf("/cps/v2/enrollments/%d", enrollmentID),
					Changes:    []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)},
				}, nil).Once()
				enrollment.Location = fmt.Sprintf("/cps/v2/enrollments/%d", enrollmentID)
				enrollment.PendingChanges = []cps.PendingChange{
					{
						Location:   fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID),
						ChangeType: "new-certificate",
					},
				}
				client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}).
					Return(enrollment, nil).Once()
				client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				}).Return(&cps.Change{
					AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
					StatusInfo: &cps.StatusInfo{
						State:  "awaiting-input",
						Status: waitUploadThirdParty,
					},
				}, nil).Once()
				// Read third party enrollment
				client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}).
					Return(enrollment, nil).Times(1)
				// CPS upload certificate
				mockCreateWithACKPostWarnings(client, enrollmentID, changeID, enrollment)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadGetChangeHistoryForComplete(client, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				// Read third party enrollment
				mockGetEnrollment(client, enrollmentID, 1, enrollmentAfterCreate)
				mockReadGetChangeHistoryForComplete(client, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				// Delete third party enrollment
				client.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
					EnrollmentID:              enrollmentID,
					AllowCancelPendingChanges: ptr.To(true),
				}).Return(&cps.RemoveEnrollmentResponse{
					Enrollment: fmt.Sprintf("%d", enrollmentID),
				}, nil).Once()
			},
			enrollment:   getSimpleEnrollment(),
			enrollmentID: 2,
			changeID:     2,
			configPath:   "testdata/TestResCPSUploadCertificate/third_party_enrollment_with_upload_cert.tf",
			checkFunc: resource.ComposeAggregateTestCheckFunc(
				checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
				resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.test_enrollment", "contract_id", "ctr_1"),
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(client, &test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							ExpectError: test.error,
							Check:       test.checkFunc,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResourceCPSUploadCertificateLifecycle(t *testing.T) {
	tests := map[string]struct {
		init                      func(*cps.Mock, *cps.GetEnrollmentResponse, *cps.GetEnrollmentResponse, int, int, int)
		enrollment                *cps.GetEnrollmentResponse
		enrollmentUpdated         *cps.GetEnrollmentResponse
		enrollmentID              int
		changeID                  int
		changeIDUpdated           int
		configPathForCreate       string
		configPathForUpdate       string
		configPathForSecondUpdate string
		checkFuncForCreate        resource.TestCheckFunc
		checkFuncForUpdate        resource.TestCheckFunc
		checkFuncForSecondUpdate  resource.TestCheckFunc
		errorForCreate            *regexp.Regexp
		errorForUpdate            *regexp.Regexp
		errorForSecondUpdate      *regexp.Regexp
	}{
		"create -> failed update after cert renewal due to missing post-ver-warnings -> update with accept all": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentUpdated *cps.GetEnrollmentResponse, enrollmentID, changeID, changeIDUpdated int) {
				// create
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				// read after create and before update
				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadCompleteForUpdate(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)
				// update
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitUploadThirdParty)
				// upsert
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadBothThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, fourWarnings, enrollmentID, changeIDUpdated)
				// expected error
				// read before update
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				// update with wait-review status (after uploading cert)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				// upsert
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadBothThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, fourWarnings, enrollmentID, changeIDUpdated)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, liveCheckAction)
				// read
				enrollmentAfterUpdate := copyEnrollmentWithEmptyPendingChanges(*enrollmentUpdated)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterUpdate)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterUpdate)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
			},
			enrollment:                createEnrollment(2, 22, true, true),
			enrollmentUpdated:         createEnrollment(2, 222, false, true),
			enrollmentID:              2,
			changeID:                  22,
			changeIDUpdated:           222,
			configPathForCreate:       "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_with_chM_true.tf",
			configPathForUpdate:       "testdata/TestResCPSUploadCertificate/certificates/changed_certificates_with_auto_approve_warnings.tf",
			configPathForSecondUpdate: "testdata/TestResCPSUploadCertificate/certificates/changed_certificates_with_auto_approve_warnings_accept.tf",
			checkFuncForCreate:        checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, true, false, []string{"CERTIFICATE_ADDED_TO_TRUST_CHAIN", "CERTIFICATE_ALREADY_LOADED", "CERTIFICATE_DATA_BLANK_OR_MISSING"})),
			checkFuncForUpdate:        nil,
			checkFuncForSecondUpdate:  checkAttrs(createMockData(certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, certRSAUpdatedForTests, trustChainRSAUpdatedForTests, true, false, true, false, []string{"CERTIFICATE_ADDED_TO_TRUST_CHAIN", "CERTIFICATE_ALREADY_LOADED", "CERTIFICATE_DATA_BLANK_OR_MISSING", "CERTIFICATE_HAS_NULL_ISSUER"})),
			errorForCreate:            nil,
			errorForUpdate:            regexp.MustCompile(`Error: could not process post verification warnings: not every warning has been acknowledged: warnings cannot be approved: "CERTIFICATE_HAS_NULL_ISSUER"`),
			errorForSecondUpdate:      nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(client, test.enrollment, test.enrollmentUpdated, test.enrollmentID, test.changeID, test.changeIDUpdated)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForCreate),
							Check:       test.checkFuncForCreate,
							ExpectError: test.errorForCreate,
						},
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForUpdate),
							Check:       test.checkFuncForUpdate,
							ExpectError: test.errorForUpdate,
						},
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForSecondUpdate),
							Check:       test.checkFuncForUpdate,
							ExpectError: test.errorForSecondUpdate,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestCreateCPSUploadCertificate(t *testing.T) {
	tests := map[string]struct {
		init         func(*cps.Mock, *cps.GetEnrollmentResponse, int, int)
		enrollment   *cps.GetEnrollmentResponse
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"successful create - RSA cert": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			error:        nil,
		},
		"successful create - ECDSA cert, without trust chain": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAForTests, "", enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, liveCheckAction)
				mockRead(m, enrollmentID, changeID, enrollment, certECDSAForTests, "", ECDSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_ecdsa.tf",
			checkFunc:    checkAttrs(createMockData(certECDSAForTests, "", "", "", false, true, false, false, nil)),
			error:        nil,
		},
		"successful create - both cert and trust chains": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadBothThirdPartyCertificateAndTrustChain(m,
					ECDSA,
					certECDSAForTests,
					trustChainECDSAForTests,
					RSA,
					certRSAForTests,
					trustChainRSAForTests,
					enrollmentID,
					changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, liveCheckAction)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				//Read after create
				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterCreate)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, ECDSA, certECDSAForTests, trustChainECDSAForTests, RSA, certRSAForTests, trustChainRSAForTests)

				//Refresh
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterCreate)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, ECDSA, certECDSAForTests, trustChainECDSAForTests, RSA, certRSAForTests, trustChainRSAForTests)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_both_certificates.tf",
			checkFunc:    checkAttrs(createMockData(certECDSAForTests, trustChainECDSAForTests, certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			error:        nil,
		},
		"successful create - with timeout": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_timeouts.tf",
			checkFunc:    checkAttrs(testDataForAttrs{certificateRSA: certRSAForTests, trustChainRSA: trustChainRSAForTests, acknowledgePostVerificationWarnings: true, timeouts: "3h"}),
			error:        nil,
		},
		"create: auto_approve_warnings match": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, true, false, []string{"CERTIFICATE_ADDED_TO_TRUST_CHAIN", "CERTIFICATE_ALREADY_LOADED", "CERTIFICATE_DATA_BLANK_OR_MISSING"})),
			error:        nil,
		},
		"create: auto_approve_warnings missing warnings error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, blankAndNullWarnings, enrollmentID, changeID)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings.tf",
			error:        regexp.MustCompile(`Error: could not process post verification warnings: not every warning has been acknowledged: warnings cannot be approved: "CERTIFICATE_HAS_NULL_ISSUER"`),
		},
		"create: auto_approve_warnings not provided and empty warning list": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithEmptyWarningList(m, enrollmentID, changeID, enrollment)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadForComplete(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        nil,
		},
		"required attribute not provided": {
			enrollment: nil,
			configPath: "testdata/TestResCPSUploadCertificate/certificates/no_certificates.tf",
			checkFunc:  nil,
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"create: auto_approve_warnings not provided and not empty warning list": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, noKMIDataWarning, enrollmentID, changeID)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        regexp.MustCompile(`Error: could not process post verification warnings: not every warning has been acknowledged: warnings cannot be approved: "CERTIFICATE_KMI_DATA_MISSING"`),
		},
		"create: auto_approve_warnings empty list and warnings": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithEmptyWarningList(m, enrollmentID, changeID, enrollment)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadForComplete(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_empty.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, true, false, nil)),
			error:        nil,
		},
		"create: change management wrong type": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockEmptyGetPostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadForComplete(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        nil,
		},
		"create: change management set to false or not specified": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockEmptyGetPostVerificationWarnings(m, enrollmentID, changeID)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_not_specified.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, false, false, []string{"Warning 1", "Warning 2", "Warning 3"})),
			error:        nil,
		},
		"create: it takes some time to acknowledge warnings": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, "", enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, liveCheckAction)
				//read's call from upsert
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, liveCheckAction)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, complete)
				mockGetChangeHistory(m, enrollmentID, 1, RSA, certRSAForTests, "")
				//rest of the flow
				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, "", RSA, enrollmentID, 1)

			},
			enrollment:   createEnrollment(2, 22, false, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/status_changes_slowly.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, "", false, true, false, true, nil)),
			error:        nil,
		},
		"create: trust chain without certificate": {
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/trust_chain_without_cert.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("provided ECDSA trust chain without ECDSA certificate. Please remove it or add a certificate"),
		},
		"create: get enrollment error": {
			init: func(m *cps.Mock, _ *cps.GetEnrollmentResponse, enrollmentID, _ int) {
				m.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{
					EnrollmentID: enrollmentID,
				}).Return(nil, fmt.Errorf("could not get an erollments")).Once()
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get an enrollment"),
		},
		"create: upload third party cert error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				m.On("UploadThirdPartyCertAndTrustChain", testutils.MockContext, cps.UploadThirdPartyCertAndTrustChainRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
					Certificates: cps.ThirdPartyCertificates{
						CertificatesAndTrustChains: []cps.CertificateAndTrustChain{
							{
								Certificate:  certRSAForTests,
								TrustChain:   trustChainRSAForTests,
								KeyAlgorithm: RSA,
							},
						},
					},
				}).Return(fmt.Errorf("could not upload a certificate")).Once()
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not upload a certificate"),
		},
		"create: get change post verification warnings error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				m.On("GetChangePostVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				}).Return(nil, fmt.Errorf("could not get change post verification warnings")).Once()
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get change post verification warnings"),
		},
		"create: acknowledge post verification warnings error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, "Some warning", enrollmentID, changeID)
				m.On("AcknowledgePostVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
					Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
					EnrollmentID:    enrollmentID,
					ChangeID:        changeID,
				}).Return(fmt.Errorf("could not acknowledge post verification warnings")).Once()
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not acknowledge post verification warnings"),
		},
		"create: get change status error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				m.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				}).Return(nil, fmt.Errorf("could not get change status")).Once()
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get change status"),
		},
		"create: acknowledge change management error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				m.On("AcknowledgeChangeManagement", testutils.MockContext, cps.AcknowledgementRequest{
					Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
					EnrollmentID:    enrollmentID,
					ChangeID:        changeID,
				}).Return(fmt.Errorf("could not acknowledge change management")).Once()
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not acknowledge change management"),
		},
		"create: no pending changes error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, _ int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
			},
			enrollment:   createEnrollment(2, 22, true, false),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("Error: could not get change ID: no pending changes were found on enrollment"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			if test.init != nil {
				test.init(client, test.enrollment, test.enrollmentID, test.changeID)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       test.checkFunc,
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestReadCPSUploadCertificate(t *testing.T) {
	tests := map[string]struct {
		init         func(*cps.Mock, *cps.GetEnrollmentResponse, int, int)
		enrollment   *cps.GetEnrollmentResponse
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"read: it takes time to get change history reach complete status but should wait": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, liveCheckAction)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, complete)
				mockGetChangeHistory(m, enrollmentID, 1, RSA, certRSAForTests, trustChainRSAForTests)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_true.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, true, nil)),
			error:        nil,
		},
		"read: get change history with correct change": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterCreate)
				mockGetChangeHistoryWithDifferentChanges(m, enrollmentID, 1, certRSAForTests, trustChainRSAForTests, RSA)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterCreate)
				mockGetChangeHistoryWithDifferentChanges(m, enrollmentID, 1, certRSAForTests, trustChainRSAForTests, RSA)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_false.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
		},
		"read: get change history error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterCreate)
				m.On("GetChangeHistory", testutils.MockContext, cps.GetChangeHistoryRequest{
					EnrollmentID: enrollmentID,
				}).Return(nil, fmt.Errorf("could not get certificate history")).Once()
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_false.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get certificate history"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       test.checkFunc,
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestUpdateCPSUploadCertificate(t *testing.T) {
	tests := map[string]struct {
		init                func(*cps.Mock, *cps.GetEnrollmentResponse, *cps.GetEnrollmentResponse, int, int, int)
		enrollment          *cps.GetEnrollmentResponse
		enrollmentUpdated   *cps.GetEnrollmentResponse
		enrollmentID        int
		changeID            int
		changeIDUpdated     int
		configPathForCreate string
		configPathForUpdate string
		checkFuncForCreate  resource.TestCheckFunc
		checkFuncForUpdate  resource.TestCheckFunc
		errorForCreate      *regexp.Regexp
		errorForUpdate      *regexp.Regexp
	}{
		"update: ignore change to acknowledge post verification warnings flag": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, liveCheckAction)
				mockReadForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/ack_post_verification_warnings/ack_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/ack_post_verification_warnings/ack_false.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: ignore change to auto_approve_warnings list": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, liveCheckAction)
				mockReadForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_updated.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, true, false, []string{"CERTIFICATE_ADDED_TO_TRUST_CHAIN", "CERTIFICATE_ALREADY_LOADED", "CERTIFICATE_DATA_BLANK_OR_MISSING"})),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, true, false, []string{"CERTIFICATE_ADDED_TO_TRUST_CHAIN", "CERTIFICATE_DATA_BLANK_OR_MISSING"})),
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: ignore ack change management set to false - warn": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadCompleteForUpdate(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterCreate)
				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_false.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: change in cert - upsert without trust chain": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentUpdated *cps.GetEnrollmentResponse, enrollmentID, changeID, changeIDUpdated int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadCompleteForUpdate(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)

				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAUpdatedForTests, "", enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeIDUpdated)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, liveCheckAction)

				enrollmentAfterUpdate := copyEnrollmentWithEmptyPendingChanges(*enrollmentUpdated)
				mockReadForComplete(m, enrollmentID, enrollmentAfterUpdate, certRSAUpdatedForTests, "", RSA)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentUpdated:   createEnrollment(2, 222, false, true),
			enrollmentID:        2,
			changeID:            22,
			changeIDUpdated:     222,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/certificates/changed_certificate.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAUpdatedForTests, "", true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: change in cert - upsert with both certs and trust chains": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentUpdated *cps.GetEnrollmentResponse, enrollmentID, changeID, changeIDUpdated int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadCompleteForUpdate(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)

				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadBothThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeIDUpdated)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, liveCheckAction)

				enrollmentAfterUpdate := copyEnrollmentWithEmptyPendingChanges(*enrollmentUpdated)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterUpdate)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentAfterUpdate)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentUpdated:   createEnrollment(2, 222, false, true),
			enrollmentID:        2,
			changeID:            22,
			changeIDUpdated:     222,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/certificates/changed_both_certificates.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData(certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, certRSAUpdatedForTests, trustChainRSAUpdatedForTests, true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: change in cert - enrollment was already updated": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentUpdated *cps.GetEnrollmentResponse, enrollmentID, changeID, changeIDUpdated int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				//Read after create
				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				//Refresh + read before update (enrollment was just updated)
				mockReadGetChangeHistory(m, enrollmentUpdated, certRSAForTests, trustChainRSAForTests, RSA, waitUploadThirdParty, enrollmentID, changeIDUpdated, 1)
				mockReadGetChangeHistory(m, enrollmentUpdated, certRSAForTests, trustChainRSAForTests, RSA, waitUploadThirdParty, enrollmentID, changeIDUpdated, 1)

				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAUpdatedForTests, "", enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeIDUpdated)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, liveCheckAction)

				enrollmentAfterUpdate := copyEnrollmentWithEmptyPendingChanges(*enrollmentUpdated)
				mockReadForComplete(m, enrollmentID, enrollmentAfterUpdate, certRSAUpdatedForTests, "", RSA)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentUpdated:   createEnrollment(2, 222, false, true),
			enrollmentID:        2,
			changeID:            22,
			changeIDUpdated:     222,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/certificates/changed_certificate.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAUpdatedForTests, "", true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: renewal with old certificate - no diff": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				//read
				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentUpdated:   createEnrollment(2, 222, true, true),
			enrollmentID:        2,
			changeID:            22,
			changeIDUpdated:     222,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  nil,
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: renewal with old certificate - enrollment was already updated, no diff, but warn": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentUpdated *cps.GetEnrollmentResponse, enrollmentID, changeID, changeIDUpdated int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockRead(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				//read
				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadGetChangeHistoryForComplete(m, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				//refresh
				mockReadGetChangeHistory(m, enrollmentUpdated, certRSAForTests, trustChainRSAForTests, RSA, waitUploadThirdParty, enrollmentID, changeIDUpdated, 1)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentUpdated:   createEnrollment(2, 222, true, true),
			enrollmentID:        2,
			changeID:            22,
			changeIDUpdated:     222,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  nil,
			errorForCreate:      nil,
			errorForUpdate:      nil, //ToDo: add warn checking once it's possible (e.g. https://github.com/hashicorp/terraform-plugin-testing/pull/17)
		},
		"update: change in already deployed certificate; no pending changes - error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, enrollmentUpdated *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadCompleteForUpdate(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentUpdated:   createEnrollment(2, 0, true, false),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/certificates/changed_certificate.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("Error: cannot make changes to certificate that is already on staging and/or production network, need to create new enrollment"),
		},
		"update: change in certificate with pending changes and wrong status - error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadCompleteForUpdate(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)

				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/certificates/changed_certificate.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("Error: cannot make changes to the certificate with current status: wait-ack-change-management"),
		},
		"update: get enrollment error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)

				enrollmentAfterCreate := copyEnrollmentWithEmptyPendingChanges(*enrollment)
				mockReadCompleteForUpdate(m, enrollmentID, enrollmentAfterCreate, certRSAForTests, trustChainRSAForTests, RSA)

				m.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{
					EnrollmentID: enrollmentID,
				}).Return(nil, fmt.Errorf("could not get an enrollment")).Times(1)
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_false.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("could not get an enrollment"),
		},
		"update: get change status error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)

				mockReadForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				m.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				}).Return(nil, fmt.Errorf("could not get change status")).Once()
			},
			enrollment:          createEnrollment(2, 22, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_false.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("could not get change status"),
		},
		"update: acknowledge change management error": {
			init: func(m *cps.Mock, enrollment *cps.GetEnrollmentResponse, _ *cps.GetEnrollmentResponse, enrollmentID, changeID, _ int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)

				mockReadForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				m.On("AcknowledgeChangeManagement", testutils.MockContext, cps.AcknowledgementRequest{
					Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
					EnrollmentID:    enrollmentID,
					ChangeID:        changeID,
				}).Return(fmt.Errorf("could not acknowledge change management")).Once()
			},
			enrollment:          createEnrollment(2, changeID, true, true),
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_false.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("could not acknowledge change management"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(client, test.enrollment, test.enrollmentUpdated, test.enrollmentID, test.changeID, test.changeIDUpdated)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForCreate),
							Check:       test.checkFuncForCreate,
							ExpectError: test.errorForCreate,
						},
						{
							Config:      testutils.LoadFixtureString(t, test.configPathForUpdate),
							Check:       test.checkFuncForUpdate,
							ExpectError: test.errorForUpdate,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResourceUploadCertificateImport(t *testing.T) {
	id := 1

	checker := test.NewImportChecker().
		CheckEqual("certificate_ecdsa_pem", certECDSAForTests).
		CheckEqual("trust_chain_ecdsa_pem", trustChainECDSAForTests).
		CheckEqual("acknowledge_post_verification_warnings", "false").
		CheckEqual("wait_for_deployment", "false").
		CheckEqual("acknowledge_change_management", "false").
		CheckEqual("auto_approve_warnings.#", "0")

	tests := map[string]struct {
		init          func(*cps.Mock)
		expectedError *regexp.Regexp
		stateCheck    func(s []*terraform.InstanceState) error
	}{
		"import": {
			init: func(client *cps.Mock) {
				enrollment := cps.GetEnrollmentResponse{
					ValidationType: "third-party",
				}
				mockGetEnrollment(client, id, 2, &enrollment)
				mockGetChangeHistory(client, id, 2, ECDSA, certECDSAForTests, trustChainECDSAForTests)
			},
			stateCheck: checker.Build(),
		},
		"import error when validation type is not third_party": {
			init: func(client *cps.Mock) {
				enrollment := cps.GetEnrollmentResponse{
					ValidationType: "dv",
				}

				mockGetEnrollment(client, id, 1, &enrollment)
			},
			expectedError: regexp.MustCompile("unable to import: wrong validation type: expected 'third-party', got 'dv'"),
		},
		"import error when no certificate yet uploaded": {
			init: func(client *cps.Mock) {
				enrollment := cps.GetEnrollmentResponse{
					ValidationType: "third-party",
				}

				mockGetEnrollment(client, id, 1, &enrollment)
				mockGetChangeHistoryWithoutCerts(client, id, 1)
			},
			expectedError: regexp.MustCompile("no certificate was yet uploaded"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(client)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:           testutils.LoadFixtureString(t, "testdata/TestResCPSUploadCertificate/import/import_upload.tf"),
							ImportState:      true,
							ImportStateId:    fmt.Sprintf("%d", id),
							ResourceName:     "akamai_cps_upload_certificate.import",
							ImportStateCheck: test.stateCheck,
							ExpectError:      test.expectedError,
						},
					},
				})
			})
		})
	}
}

// copyEnrollmentWithEmptyPendingChanges returns enrollment after enrollment reaches "complete" state - it's pending changes disappear
func copyEnrollmentWithEmptyPendingChanges(enrollment cps.GetEnrollmentResponse) *cps.GetEnrollmentResponse {
	enrollment.PendingChanges = []cps.PendingChange{}
	return &enrollment
}

var (
	certECDSAForTests              = "-----BEGIN CERTIFICATE ECDSA REQUEST-----\n...\n-----END CERTIFICATE ECDSA REQUEST-----"
	certRSAForTests                = "-----BEGIN CERTIFICATE RSA REQUEST-----\n...\n-----END CERTIFICATE RSA REQUEST-----"
	certRSAUpdatedForTests         = "-----BEGIN CERTIFICATE RSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE RSA REQUEST UPDATED-----"
	certECDSAUpdatedForTests       = "-----BEGIN CERTIFICATE ECDSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE ECDSA REQUEST UPDATED-----"
	trustChainRSAForTests          = "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST-----\n...\n-----END CERTIFICATE TRUST-CHAIN RSA REQUEST-----"
	trustChainECDSAForTests        = "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----\n...\n-----END CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----"
	trustChainRSAUpdatedForTests   = "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE TRUST-CHAIN RSA REQUEST UPDATED-----"
	trustChainECDSAUpdatedForTests = "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST UPDATED-----\n...\n-----END CERTIFICATE TRUST-CHAIN ECDSA REQUEST UPDATED-----"
	RSA                            = "RSA"
	ECDSA                          = "ECDSA"
	blankAndNullWarnings           = "Certificate data is blank or missing.\nCertificate has a null issuer"
	noKMIDataWarning               = "No KMI data is available for the new certificate."
	threeWarnings                  = "Certificate Added to the new Trust Chain: TEST\nThere is a problem deploying the 'RSA' certificate.  Please contact your Akamai support team to resolve the issue.\nCertificate data is blank or missing."
	fourWarnings                   = "Certificate Added to the new Trust Chain: TEST\nThere is a problem deploying the 'RSA' certificate.  Please contact your Akamai support team to resolve the issue.\nCertificate data is blank or missing.\nCertificate has a null issuer"

	// createEnrollment returns third-party enrollment with provided values
	createEnrollment = func(enrollmentID, changeID int, changeManagement, pendingChangesPresent bool) *cps.GetEnrollmentResponse {
		pendingChanges := []cps.PendingChange{
			{
				Location:   fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID),
				ChangeType: "new-certificate",
			},
		}
		if !pendingChangesPresent {
			pendingChanges = []cps.PendingChange{}
		}
		return &cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@terraform-test.net",
				FirstName:        "R5",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "000111222",
				PostalCode:       "02142",
				Region:           "MA",
				Title:            "Administrator",
			},
			Location:             fmt.Sprintf("/cps/v2/enrollments/%d", enrollmentID),
			CertificateChainType: "default",
			CertificateType:      "third-party",
			ChangeManagement:     changeManagement,
			CSR: &cps.CSR{
				C:    "US",
				CN:   "akatest.com",
				L:    "Cambridge",
				O:    "Akamai",
				OU:   "WebEx",
				SANS: []string{"san.test.akamai1.com", "san.test.akamai2.com", "san.test.akamai3.com"},
				ST:   "MA",
			},
			EnableMultiStackedCertificates: true,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1", "TLSv2_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"akatest.com"},
				},
				Geography:        "core",
				MustHaveCiphers:  "ak-akamai-default",
				OCSPStapling:     "on",
				PreferredCiphers: "ak-akamai-default",
				QuicEnabled:      false,
				SecureNetwork:    "enhanced-tls",
				SNIOnly:          true,
			},
			Org: &cps.Org{
				AddressLineOne: "150 Broadway",
				AddressLineTwo: "building 1",
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "55555",
				Region:         "MA",
			},
			RA:                 "third-party",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r5d2@testakamai.com",
				FirstName:        "R5",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
				Title:            "Technician",
			},
			MaxAllowedWildcardSanNames: 25,
			MaxAllowedSanNames:         100,
			PendingChanges:             pendingChanges,
			ValidationType:             "third-party",
		}
	}

	// createMockData return test data structure used in test with provided configuration
	createMockData = func(certECDSA, trustECDSA, certRSA, trustRSA string, ackChM, ackPost, autoAppWarnSet, waitForDeploy bool, warnings []string) testDataForAttrs {
		return testDataForAttrs{
			certificateECDSA:                    certECDSA,
			certificateRSA:                      certRSA,
			trustChainECDSA:                     trustECDSA,
			trustChainRSA:                       trustRSA,
			acknowledgeChangeManagement:         ackChM,
			acknowledgePostVerificationWarnings: ackPost,
			isAutoApproveWarningsSet:            autoAppWarnSet,
			waitForDeployment:                   waitForDeploy,
			autoApproveWarnings:                 warnings,
		}
	}

	// mockCreateWithEmptyWarningList mocks getting empty warnings list from API
	mockCreateWithEmptyWarningList = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.GetEnrollmentResponse) {
		mockGetEnrollment(client, enrollmentID, 1, enrollment)
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   fmt.Sprintf("/cps/v2/enrollments/2/changes/%d", changeID),
				ChangeType: "new-certificate",
			},
		}
		mockUploadThirdPartyCertificateAndTrustChain(client, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
		mockEmptyGetPostVerificationWarnings(client, enrollmentID, changeID)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitAckChangeManagement)
		mockAcknowledgeChangeManagement(client, enrollmentID, changeID)
	}

	// mockUpdate mocks default approach when updating the resource
	mockUpdate = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.GetEnrollmentResponse) {
		mockGetEnrollment(client, enrollmentID, 1, enrollment)
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   fmt.Sprintf("/cps/v2/enrollments/2/changes/%d", changeID),
				ChangeType: "new-certificate",
			},
		}
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitAckChangeManagement)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitAckChangeManagement)
		mockAcknowledgeChangeManagement(client, enrollmentID, changeID)
	}

	// mockCreateWithACKPostWarnings mocks acknowledging post verification warnings along with creation of the resource
	mockCreateWithACKPostWarnings = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.GetEnrollmentResponse) {
		mockGetEnrollment(client, enrollmentID, 1, enrollment)
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID),
				ChangeType: "new-certificate",
			},
		}
		mockUploadThirdPartyCertificateAndTrustChain(client, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
		mockGetPostVerificationWarnings(client, "Certificate Added to the new Trust Chain: TEST\nThere is a problem deploying the 'RSA' certificate.  Please contact your Akamai support team to resolve the issue.\nCertificate data is blank or missing.", enrollmentID, changeID)
		mockAcknowledgePostVerificationWarnings(client, enrollmentID, changeID)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, liveCheckAction)
	}

	// mockReadForComplete mocks Read functions when cert has been already deployed to production (status = complete)
	mockReadForComplete = func(client *cps.Mock, enrollmentID int, enrollment *cps.GetEnrollmentResponse, certificate, trustChain, keyAlgorithm string) {
		mockReadGetChangeHistoryForComplete(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockReadGetChangeHistoryForComplete(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
	}

	// mockRead mocks Read functions
	mockRead = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.GetEnrollmentResponse, certificate, trustChain, keyAlgorithm, status string) {
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, status, enrollmentID, changeID, 1)
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, status, enrollmentID, changeID, 1)
	}

	// mockReadForUpdate mocks Read functions during Update
	mockReadForUpdate = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.GetEnrollmentResponse, certificate, trustChain, keyAlgorithm, status string) {
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, status, enrollmentID, changeID, 1)
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, status, enrollmentID, changeID, 1)
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, status, enrollmentID, changeID, 1)
	}

	// mockReadCompleteForUpdate mocks Read functions during Update when cert has been already deployed to production (status = complete) and it does not change during all reads
	mockReadCompleteForUpdate = func(client *cps.Mock, enrollmentID int, enrollment *cps.GetEnrollmentResponse, certificate, trustChain, keyAlgorithm string) {
		mockReadGetChangeHistoryForComplete(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockReadGetChangeHistoryForComplete(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockReadGetChangeHistoryForComplete(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
	}

	// mockReadGetChangeHistory mocks Read functions across execution
	mockReadGetChangeHistory = func(client *cps.Mock, enrollment *cps.GetEnrollmentResponse, certificate, trustChain, keyAlgorithm, status string, enrollmentID, changeID, timesToRun int) {
		mockGetEnrollment(client, enrollmentID, timesToRun, enrollment)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, status)
		mockGetChangeHistory(client, enrollmentID, timesToRun, keyAlgorithm, certificate, trustChain)
	}

	// mockReadGetChangeHistoryForComplete mocks Read functions across execution when cert has been already deployed to production (status = complete)
	mockReadGetChangeHistoryForComplete = func(client *cps.Mock, enrollment *cps.GetEnrollmentResponse, certificate, trustChain, keyAlgorithm string, enrollmentID, timesToRun int) {
		mockGetEnrollment(client, enrollmentID, timesToRun, enrollment)
		mockGetChangeHistory(client, enrollmentID, timesToRun, keyAlgorithm, certificate, trustChain)
	}

	// mockGetChangeHistory mocks GetChangeHistory call with provided data
	mockGetChangeHistory = func(client *cps.Mock, enrollmentID, timesToRun int, keyAlgorithm, certificate, trustChain string) {
		client.On("GetChangeHistory", testutils.MockContext, cps.GetChangeHistoryRequest{
			EnrollmentID: enrollmentID,
		}).Return(&cps.GetChangeHistoryResponse{
			Changes: []cps.ChangeHistory{
				{
					Action:                   "",
					ActionDescription:        "",
					BusinessCaseID:           "",
					CreatedBy:                "",
					CreatedOn:                "",
					LastUpdated:              "",
					MultiStackedCertificates: nil,
					PrimaryCertificate: cps.CertificateChangeHistory{
						Certificate:  certificate,
						TrustChain:   trustChain,
						CSR:          "",
						KeyAlgorithm: keyAlgorithm,
					},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "active",
				},
				{
					Action:                   "",
					ActionDescription:        "",
					BusinessCaseID:           "",
					CreatedBy:                "",
					CreatedOn:                "",
					LastUpdated:              "",
					MultiStackedCertificates: nil,
					PrimaryCertificate: cps.CertificateChangeHistory{
						Certificate:  certECDSAForTests,
						TrustChain:   trustChainECDSAForTests,
						CSR:          "",
						KeyAlgorithm: ECDSA,
					},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "inactive",
				},
			},
		}, nil).Times(timesToRun)
	}

	// mockGetChangeHistoryBothCerts mocks GetChangeHistory call with provided data for both certs
	mockGetChangeHistoryBothCerts = func(client *cps.Mock, enrollmentID, timesToRun int, keyAlgorithm, certificate, trustChain, keyAlgorithm2, certificate2, trustChain2 string) {
		client.On("GetChangeHistory", testutils.MockContext, cps.GetChangeHistoryRequest{
			EnrollmentID: enrollmentID,
		}).Return(&cps.GetChangeHistoryResponse{
			Changes: []cps.ChangeHistory{
				{
					Action:            "",
					ActionDescription: "",
					BusinessCaseID:    "",
					CreatedBy:         "",
					CreatedOn:         "",
					LastUpdated:       "",
					MultiStackedCertificates: []cps.CertificateChangeHistory{
						{
							Certificate:  certificate2,
							TrustChain:   trustChain2,
							CSR:          "",
							KeyAlgorithm: keyAlgorithm2,
						},
					},
					PrimaryCertificate: cps.CertificateChangeHistory{
						Certificate:  certificate,
						TrustChain:   trustChain,
						CSR:          "",
						KeyAlgorithm: keyAlgorithm,
					},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "active",
				},
			},
		}, nil).Times(timesToRun)
	}

	// mockGetChangeHistoryWithDifferentChanges mocks GetChangeHistory call with different changes
	mockGetChangeHistoryWithDifferentChanges = func(client *cps.Mock, enrollmentID, timesToRun int, certificate, trustChain, keyAlgorithm string) {
		client.On("GetChangeHistory", testutils.MockContext, cps.GetChangeHistoryRequest{
			EnrollmentID: enrollmentID,
		}).Return(
			&cps.GetChangeHistoryResponse{
				Changes: []cps.ChangeHistory{
					{
						Action:                         "renew",
						ActionDescription:              "Renew Certificate",
						BusinessCaseID:                 "",
						CreatedBy:                      "user",
						CreatedOn:                      "2022-10-11T10:40:44Z",
						LastUpdated:                    "2022-10-11T10:40:44Z",
						MultiStackedCertificates:       []cps.CertificateChangeHistory{},
						PrimaryCertificate:             cps.CertificateChangeHistory{},
						PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
						RA:                             "third-party",
						Status:                         "incomplete",
					},
					{
						Action:                         "renew",
						ActionDescription:              "Renew Certificate",
						BusinessCaseID:                 "",
						CreatedBy:                      "<auto-renewal>",
						CreatedOn:                      "2022-10-11T10:40:44Z",
						LastUpdated:                    "2022-10-11T10:40:44Z",
						MultiStackedCertificates:       []cps.CertificateChangeHistory{},
						PrimaryCertificate:             cps.CertificateChangeHistory{},
						PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
						RA:                             "third-party",
						Status:                         "cancelled",
					},
					{
						Action:                   "",
						ActionDescription:        "",
						BusinessCaseID:           "",
						CreatedBy:                "",
						CreatedOn:                "",
						LastUpdated:              "",
						MultiStackedCertificates: nil,
						PrimaryCertificate: cps.CertificateChangeHistory{
							Certificate:  certificate,
							TrustChain:   trustChain,
							CSR:          "",
							KeyAlgorithm: keyAlgorithm,
						},
						PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
						RA:                             "",
						Status:                         "active",
					},
				},
			}, nil).Times(timesToRun)
	}

	// mockGetChangeHistoryWithoutCerts mocks GetChangeHistory call when no cert was not yet uploaded
	mockGetChangeHistoryWithoutCerts = func(client *cps.Mock, enrollmentID, timesToRun int) {
		client.On("GetChangeHistory", testutils.MockContext, cps.GetChangeHistoryRequest{
			EnrollmentID: enrollmentID,
		}).Return(
			&cps.GetChangeHistoryResponse{
				Changes: []cps.ChangeHistory{
					{
						Action:            "new-certificate",
						ActionDescription: "Create New Certificate",
						Status:            "wait-upload-third-party",
						LastUpdated:       "2022-07-21T21:40:00Z",
						CreatedBy:         "user",
						CreatedOn:         "2022-07-21T21:40:00Z",
						RA:                "third-party",
					},
				},
			}, nil).Times(timesToRun)
	}

	// mockGetEnrollment mocks GetEnrollment call with provided values
	mockGetEnrollment = func(client *cps.Mock, enrollmentID, timesToRun int, enrollment *cps.GetEnrollmentResponse) {
		client.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{
			EnrollmentID: enrollmentID,
		}).Return(enrollment, nil).Times(timesToRun)
	}

	// mockUploadThirdPartyCertificateAndTrustChain mocks UploadThirdPartyCertificateAndTrustChain call with provided values
	mockUploadThirdPartyCertificateAndTrustChain = func(client *cps.Mock, keyAlgorithm, certificate, trustChain string, enrollmentID, changeID int) {
		client.On("UploadThirdPartyCertAndTrustChain", testutils.MockContext, cps.UploadThirdPartyCertAndTrustChainRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
			Certificates: cps.ThirdPartyCertificates{
				CertificatesAndTrustChains: []cps.CertificateAndTrustChain{
					{
						Certificate:  certificate,
						TrustChain:   trustChain,
						KeyAlgorithm: keyAlgorithm,
					},
				},
			},
		}).Return(nil).Once()
	}

	// mockUploadBothThirdPartyCertificateAndTrustChain mocks UploadThirdPartyCertificateAndTrustChain call for both certificates
	mockUploadBothThirdPartyCertificateAndTrustChain = func(client *cps.Mock, keyAlgorithm, certificate, trustChain, keyAlgorithm2, certificate2, trustChain2 string, enrollmentID, changeID int) {
		client.On("UploadThirdPartyCertAndTrustChain", testutils.MockContext, cps.UploadThirdPartyCertAndTrustChainRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
			Certificates: cps.ThirdPartyCertificates{
				CertificatesAndTrustChains: []cps.CertificateAndTrustChain{
					{
						Certificate:  certificate,
						TrustChain:   trustChain,
						KeyAlgorithm: keyAlgorithm,
					},
					{
						Certificate:  certificate2,
						TrustChain:   trustChain2,
						KeyAlgorithm: keyAlgorithm2,
					},
				},
			},
		}).Return(nil).Once()
	}

	// mockGetPostVerificationWarnings mocks GetPostVerificationWarnings call with provided values
	mockGetPostVerificationWarnings = func(client *cps.Mock, warnings string, enrollmentID, changeID int) {
		client.On("GetChangePostVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		}).Return(&cps.PostVerificationWarnings{
			Warnings: warnings,
		}, nil).Once()
	}

	// mockEmptyGetPostVerificationWarnings mocks GetPostVerificationWarnings call with empty warnings list
	mockEmptyGetPostVerificationWarnings = func(client *cps.Mock, enrollmentID, changeID int) {
		client.On("GetChangePostVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		}).Return(&cps.PostVerificationWarnings{}, nil).Once()
	}

	// mockAcknowledgePostVerificationWarnings mocks AcknowledgePostVerificationWarnings call with provided values
	mockAcknowledgePostVerificationWarnings = func(client *cps.Mock, enrollmentID, changeID int) {
		client.On("AcknowledgePostVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
			EnrollmentID:    enrollmentID,
			ChangeID:        changeID,
		}).Return(nil).Once()
	}

	// mockGetChangeStatus mocks GetChangeStatus call with provided values
	mockGetChangeStatus = func(client *cps.Mock, enrollmentID, changeID, timesToRun int, status string) {
		client.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{
				{
					Info:              "Test",
					RequiredToProceed: false,
					Type:              "Type",
					Update:            "Update",
				},
			},
			StatusInfo: &cps.StatusInfo{
				DeploymentSchedule: nil,
				Description:        "Desc",
				Error:              nil,
				State:              "",
				Status:             status,
			},
		}, nil).Times(timesToRun)
	}

	// mockAcknowledgeChangeManagement mocks AcknowledgeChangeManagement call with provided values
	mockAcknowledgeChangeManagement = func(client *cps.Mock, enrollmentID, changeID int) {
		client.On("AcknowledgeChangeManagement", testutils.MockContext, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
			EnrollmentID:    enrollmentID,
			ChangeID:        changeID,
		}).Return(nil).Once()
	}

	// checkAttrs creates check functions for a resource based on received data
	checkAttrs = func(data testDataForAttrs) resource.TestCheckFunc {
		warningsLength := len(data.autoApproveWarnings)
		var warningsCheck []resource.TestCheckFunc
		if warningsLength == 0 {
			if data.isAutoApproveWarningsSet {
				warningsCheck = append(warningsCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "auto_approve_warnings.#", "0"))
			} else {
				warningsCheck = append(warningsCheck, resource.TestCheckNoResourceAttr("akamai_cps_upload_certificate.test", "auto_approve_warnings"))
			}
		} else {
			warningsCheck = append(warningsCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "auto_approve_warnings.#", strconv.Itoa(warningsLength)))
			for i := 0; i < warningsLength; i++ {
				warningsCheck = append(warningsCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", fmt.Sprintf("auto_approve_warnings.%d", i), data.autoApproveWarnings[i]))
			}
		}

		var certificateCheck []resource.TestCheckFunc
		certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "certificate_ecdsa_pem", data.certificateECDSA))
		certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "certificate_rsa_pem", data.certificateRSA))
		certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "trust_chain_ecdsa_pem", data.trustChainECDSA))
		certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "trust_chain_rsa_pem", data.trustChainRSA))

		checkFuncs := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "enrollment_id", "2"),
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "acknowledge_post_verification_warnings", strconv.FormatBool(data.acknowledgePostVerificationWarnings)),
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "acknowledge_change_management", strconv.FormatBool(data.acknowledgeChangeManagement)),
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "wait_for_deployment", strconv.FormatBool(data.waitForDeployment)),
			resource.ComposeAggregateTestCheckFunc(warningsCheck...),
			resource.ComposeAggregateTestCheckFunc(certificateCheck...),
		}

		if data.timeouts != "" {
			checkFuncs = append(checkFuncs,
				resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "timeouts.#", "1"),
				resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "timeouts.0.default", data.timeouts),
			)
		} else {
			checkFuncs = append(checkFuncs,
				resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "timeouts.#", "0"),
			)
		}

		return resource.ComposeAggregateTestCheckFunc(
			checkFuncs...,
		)
	}
)

// testDataForAttrs holds data used to create check functions
type testDataForAttrs struct {
	certificateECDSA                    string
	certificateRSA                      string
	trustChainECDSA                     string
	trustChainRSA                       string
	acknowledgeChangeManagement         bool
	acknowledgePostVerificationWarnings bool
	isAutoApproveWarningsSet            bool
	waitForDeployment                   bool
	autoApproveWarnings                 []string
	timeouts                            string
}
