package cps

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var someStatusPostReviewWarning = "live-check-action"

func TestResourceCPSUploadCertificate(t *testing.T) {
	tests := map[string]struct {
		init                func(*testing.T, *cps.Mock, *cps.Enrollment, int, int)
		enrollment          *cps.Enrollment
		enrollmentID        int
		changeID            int
		configPathForCreate string
		configPathForUpdate string
		checkFuncForCreate  resource.TestCheckFunc
		checkFuncForUpdate  resource.TestCheckFunc
		error               *regexp.Regexp
	}{
		"create with ch-mgmt false, update to true": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockUpdate(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockUpdate(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, complete)
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
			test.init(t, client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPathForCreate),
							Check:       test.checkFuncForCreate,
							ExpectError: test.error,
						},
						{
							Config:      loadFixtureString(test.configPathForUpdate),
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
		init         func(*testing.T, *cps.Mock, *cps.Enrollment, int, int)
		enrollment   cps.Enrollment
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"create third_party_enrollment, create cps_upload_certificate with enrollment_id which is third_party_enrollment's resource dependency": {
			init: func(t *testing.T, client *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				// Create third party enrollment
				client.On("CreateEnrollment",
					mock.Anything,
					cps.CreateEnrollmentRequest{
						Enrollment: *enrollment,
						ContractID: "1",
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
				client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}).
					Return(enrollment, nil).Once()
				client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
				client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}).
					Return(enrollment, nil).Times(1)
				// CPS upload certificate
				mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 1, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(client, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(client, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 1, enrollment, waitAckChangeManagement)
				// Read third party enrollment
				mockGetEnrollment(client, enrollmentID, 1, enrollment)
				mockReadGetChangeHistory(client, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 2, enrollment, waitAckChangeManagement)
				// Delete third party enrollment
				client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
					EnrollmentID:              enrollmentID,
					AllowCancelPendingChanges: tools.BoolPtr(true),
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
			test.init(t, client, &test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPath),
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
		init                      func(*testing.T, *cps.Mock, *cps.Enrollment, *cps.Enrollment, int, int, int)
		enrollment                *cps.Enrollment
		enrollmentUpdated         *cps.Enrollment
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				// checkUnacknowledgedWarnings
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				// create
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				// read after create
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
				// checkUnacknowledgedWarnings
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, complete)
				// read before update
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				// checkUnacknowledgedWarnings
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, complete)
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
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, enrollmentUpdated, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				// checkUnacknowledgedWarnings
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitReviewThirdPartyCert)
				// update with wait-review status (after uploading cert)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				// upsert
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadBothThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, fourWarnings, enrollmentID, changeIDUpdated)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, someStatusPostReviewWarning)
				// read
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, enrollmentUpdated, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				// checkUnacknowledgedWarnings
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 1, enrollment, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, enrollmentUpdated, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				// checkUnacknowledgedWarnings
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 2, enrollment, waitAckChangeManagement)
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
			test.init(t, client, test.enrollment, test.enrollmentUpdated, test.enrollmentID, test.changeID, test.changeIDUpdated)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPathForCreate),
							Check:       test.checkFuncForCreate,
							ExpectError: test.errorForCreate,
						},
						{
							Config:      loadFixtureString(test.configPathForUpdate),
							Check:       test.checkFuncForUpdate,
							ExpectError: test.errorForUpdate,
						},
						{
							Config:      loadFixtureString(test.configPathForSecondUpdate),
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
		init         func(*testing.T, *cps.Mock, *cps.Enrollment, int, int)
		enrollment   *cps.Enrollment
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"successful create - RSA cert": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			error:        nil,
		},
		"successful create - ECDSA cert, without trust chain": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAForTests, "", enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, someStatusPostReviewWarning)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certECDSAForTests, "", ECDSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_ecdsa.tf",
			checkFunc:    checkAttrs(createMockData(certECDSAForTests, "", "", "", false, true, false, false, nil)),
			error:        nil,
		},
		"successful create - both cert and trust chains": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
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
				mockGetChangeStatus(m, enrollmentID, changeID, 1, someStatusPostReviewWarning)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, enrollment, ECDSA, certECDSAForTests, trustChainECDSAForTests, RSA, certRSAForTests, trustChainRSAForTests)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 1, enrollment, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, enrollment, ECDSA, certECDSAForTests, trustChainECDSAForTests, RSA, certRSAForTests, trustChainRSAForTests)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 2, enrollment, complete)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_both_certificates.tf",
			checkFunc:    checkAttrs(createMockData(certECDSAForTests, trustChainECDSAForTests, certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			error:        nil,
		},
		"create: auto_approve_warnings match": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, true, false, []string{"CERTIFICATE_ADDED_TO_TRUST_CHAIN", "CERTIFICATE_ALREADY_LOADED", "CERTIFICATE_DATA_BLANK_OR_MISSING"})),
			error:        nil,
		},
		"create: auto_approve_warnings missing warnings error": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithEmptyWarningList(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        nil,
		},
		"required attribute not provided": {
			init:       func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, _, _ int) {},
			enrollment: nil,
			configPath: "testdata/TestResCPSUploadCertificate/certificates/no_certificates.tf",
			checkFunc:  nil,
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"create: auto_approve_warnings not provided and not empty warning list": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithEmptyWarningList(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_empty.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, true, false, nil)),
			error:        nil,
		},
		"create: change management wrong type": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockEmptyGetPostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        nil,
		},
		"create: change management set to false or not specified": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockEmptyGetPostVerificationWarnings(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_not_specified.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, false, false, []string{"Warning 1", "Warning 2", "Warning 3"})),
			error:        nil,
		},
		"create: it takes some time to acknowledge warnings": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, "", enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, someStatusPostReviewWarning)
				//read's call from upsert
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, someStatusPostReviewWarning)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, complete)
				mockGetChangeHistory(m, enrollmentID, 1, enrollment, RSA, certRSAForTests, "")
				//rest of the flow
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 1, enrollment, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 2, complete)
				mockGetChangeHistory(m, enrollmentID, 1, enrollment, RSA, certRSAForTests, "")
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 2, enrollment, complete)

			},
			enrollment:   createEnrollment(2, 22, false, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/status_changes_slowly.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, "", false, true, false, true, nil)),
			error:        nil,
		},
		"create: trust chain without certificate": {
			init: func(_ *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, complete)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/trust_chain_without_cert.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("provided ECDSA trust chain without ECDSA certificate. Please remove it or add a certificate"),
		},
		"create: get enrollment error": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				m.On("UploadThirdPartyCertAndTrustChain", mock.Anything, cps.UploadThirdPartyCertAndTrustChainRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				m.On("GetChangePostVerificationWarnings", mock.Anything, cps.GetChangeRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, "Some warning", enrollmentID, changeID)
				m.On("AcknowledgePostVerificationWarnings", mock.Anything, cps.AcknowledgementRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				m.On("AcknowledgeChangeManagement", mock.Anything, cps.AcknowledgementRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 4, enrollment)
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
			test.init(t, client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPath),
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
		init         func(*testing.T, *cps.Mock, *cps.Enrollment, int, int)
		enrollment   *cps.Enrollment
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"read: get change history": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 1, enrollment, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 2, enrollment, complete)
				mockGetChangeStatus(m, enrollmentID, changeID, 2, complete)
				mockGetChangeStatus(m, enrollmentID, changeID, 2, complete)
				mockGetChangeHistory(m, enrollmentID, 2, enrollment, RSA, certRSAForTests, trustChainRSAForTests)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_true.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, true, nil)),
			error:        nil,
		},
		"read: get change history with wait-upload-third-party-status": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitUploadThirdParty)
				mockGetChangeHistory(m, enrollmentID, 1, enrollment, RSA, certRSAForTests, trustChainRSAForTests)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 1, enrollment, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitUploadThirdParty)
				mockGetChangeHistory(m, enrollmentID, 1, enrollment, RSA, certRSAForTests, trustChainRSAForTests)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 2, enrollment, waitAckChangeManagement)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_true.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, true, nil)),
			error:        nil,
		},
		"read: get change history with correct change": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeHistoryWithDifferentChanges(m, enrollmentID, 1, certRSAForTests, trustChainRSAForTests, RSA)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 1, enrollment, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeHistoryWithDifferentChanges(m, enrollmentID, 1, certRSAForTests, trustChainRSAForTests, RSA)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 2, enrollment, complete)
			},
			enrollment:   createEnrollment(2, 22, true, true),
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_false.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
		},
		"read: get change history error": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				m.On("GetChangeHistory", mock.Anything, cps.GetChangeHistoryRequest{
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
			test.init(t, client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPath),
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
		init                func(*testing.T, *cps.Mock, *cps.Enrollment, *cps.Enrollment, int, int, int)
		enrollment          *cps.Enrollment
		enrollmentUpdated   *cps.Enrollment
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, someStatusPostReviewWarning)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, someStatusPostReviewWarning)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, complete)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, complete)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAUpdatedForTests, "", enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, 222)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, someStatusPostReviewWarning)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeIDUpdated, enrollmentUpdated, certRSAUpdatedForTests, "", RSA, waitAckChangeManagement)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockUploadBothThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitReviewThirdPartyCert)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, 222)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeIDUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, someStatusPostReviewWarning)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, enrollmentUpdated, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeIDUpdated, 1, enrollmentUpdated, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeHistoryBothCerts(m, enrollmentID, 1, enrollmentUpdated, ECDSA, certECDSAUpdatedForTests, trustChainECDSAUpdatedForTests, RSA, certRSAUpdatedForTests, trustChainRSAUpdatedForTests)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeIDUpdated, 2, enrollmentUpdated, complete)
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
		"update: renewal with old certificate - api error": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarnings(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				// STEP 2
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				// UPDATE
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				mockGetChangeStatus(m, enrollmentID, changeIDUpdated, 1, waitUploadThirdParty)
				mockGetEnrollment(m, enrollmentID, 1, enrollmentUpdated)
				m.On("UploadThirdPartyCertAndTrustChain", mock.Anything, cps.UploadThirdPartyCertAndTrustChainRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeIDUpdated,
					Certificates: cps.ThirdPartyCertificates{
						CertificatesAndTrustChains: []cps.CertificateAndTrustChain{
							{
								Certificate:  certRSAForTests,
								TrustChain:   trustChainRSAForTests,
								KeyAlgorithm: RSA,
							},
						},
					},
				}).Return(fmt.Errorf("provided certificate is wrong")).Once()
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
			errorForUpdate:      regexp.MustCompile("provided certificate is wrong"),
		},
		"update: change in already deployed certificate; no pending changes - error": {
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, complete)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
			init: func(t *testing.T, m *cps.Mock, enrollment *cps.Enrollment, enrollmentUpdated *cps.Enrollment, enrollmentID, changeID, changeIDUpdated int) {
				mockCheckUnacknowledgedWarnings(m, enrollmentID, changeID, 3, enrollment, waitUploadThirdParty)
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadAndCheckUnacknowledgedWarningsForUpdate(m, enrollmentID, changeID, enrollment, certRSAForTests, trustChainRSAForTests, RSA, waitAckChangeManagement)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChangeManagement)
				m.On("AcknowledgeChangeManagement", mock.Anything, cps.AcknowledgementRequest{
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
			test.init(t, client, test.enrollment, test.enrollmentUpdated, test.enrollmentID, test.changeID, test.changeIDUpdated)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPathForCreate),
							Check:       test.checkFuncForCreate,
							ExpectError: test.errorForCreate,
						},
						{
							Config:      loadFixtureString(test.configPathForUpdate),
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

	tests := map[string]struct {
		init          func(*cps.Mock)
		expectedError *regexp.Regexp
		stateCheck    func(s []*terraform.InstanceState) error
	}{
		"import": {
			init: func(client *cps.Mock) {
				enrollment := cps.Enrollment{
					ValidationType: "third-party",
				}

				mockGetEnrollment(client, id, 2, &enrollment)
				mockGetChangeHistory(client, id, 2, &enrollment, ECDSA, certECDSAForTests, trustChainECDSAForTests)
			},
			stateCheck: func(s []*terraform.InstanceState) error {
				state := s[0]
				assertAttributeFor(state, t, "certificate_ecdsa_pem", certECDSAForTests)
				assertAttributeFor(state, t, "trust_chain_ecdsa_pem", trustChainECDSAForTests)
				assertAttributeFor(state, t, "acknowledge_post_verification_warnings", "false")
				assertAttributeFor(state, t, "wait_for_deployment", "false")
				assertAttributeFor(state, t, "acknowledge_change_management", "false")
				assertAttributeFor(state, t, "auto_approve_warnings.#", "0")
				return nil
			},
		},
		"import error when validation type is not third_party": {
			init: func(client *cps.Mock) {
				enrollment := cps.Enrollment{
					ValidationType: "dv",
				}

				mockGetEnrollment(client, id, 1, &enrollment)
			},
			expectedError: regexp.MustCompile("unable to import: wrong validation type: expected 'third-party', got 'dv'"),
		},
		"import error when no certificate yet uploaded": {
			init: func(client *cps.Mock) {
				enrollment := cps.Enrollment{
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
					ProviderFactories: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:           loadFixtureString("testdata/TestResCPSUploadCertificate/import/import_upload.tf"),
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

func assertAttributeFor(state *terraform.InstanceState, t *testing.T, key, value string) {
	valueInState, exist := state.Attributes[key]
	assert.True(t, exist, fmt.Sprintf("attribute '%s' was not present", key))
	assert.Equal(t, value, valueInState, fmt.Sprintf("attribute '%s' has incorrect value %s", key, valueInState))
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
	createEnrollment = func(enrollmentID, changeID int, changeManagement, pendingChangesPresent bool) *cps.Enrollment {
		pendingChanges := []cps.PendingChange{
			{
				Location:   fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID),
				ChangeType: "new-certificate",
			},
		}
		if !pendingChangesPresent {
			pendingChanges = []cps.PendingChange{}
		}
		return &cps.Enrollment{
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
	mockCreateWithEmptyWarningList = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.Enrollment) {
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
	mockUpdate = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.Enrollment) {
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
	mockCreateWithACKPostWarnings = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.Enrollment) {
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
		mockGetChangeStatus(client, enrollmentID, changeID, 1, someStatusPostReviewWarning)
	}

	// mockReadAndCheckUnacknowledgedWarnings mocks Read and CustomDiff functions
	mockReadAndCheckUnacknowledgedWarnings = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.Enrollment, certificate, trustChain, keyAlgorithm, status string) {
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 1, enrollment, status)
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 2, enrollment, status)
	}

	// mockReadAndCheckUnacknowledgedWarningsForUpdate mocks Read and CustomDiff functions during Update
	mockReadAndCheckUnacknowledgedWarningsForUpdate = func(client *cps.Mock, enrollmentID, changeID int, enrollment *cps.Enrollment, certificate, trustChain, keyAlgorithm, status string) {
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 1, enrollment, status)
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 2, enrollment, status)
		mockReadGetChangeHistory(client, enrollment, certificate, trustChain, keyAlgorithm, enrollmentID, 1)
		mockCheckUnacknowledgedWarnings(client, enrollmentID, changeID, 3, enrollment, status)
	}

	// mockReadGetChangeHistory mocks GetChangeHistory call with provided values
	mockReadGetChangeHistory = func(client *cps.Mock, enrollment *cps.Enrollment, certificate, trustChain, keyAlgorithm string, enrollmentID, timesToRun int) {
		mockGetEnrollment(client, enrollmentID, timesToRun, enrollment)
		mockGetChangeHistory(client, enrollmentID, timesToRun, enrollment, keyAlgorithm, certificate, trustChain)
	}

	// mockGetChangeHistory mocks GetChangeHistory call with provided data
	mockGetChangeHistory = func(client *cps.Mock, enrollmentID, timesToRun int, enrollment *cps.Enrollment, keyAlgorithm, certificate, trustChain string) {
		client.On("GetChangeHistory", mock.Anything, cps.GetChangeHistoryRequest{
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
	mockGetChangeHistoryBothCerts = func(client *cps.Mock, enrollmentID, timesToRun int, enrollment *cps.Enrollment, keyAlgorithm, certificate, trustChain, keyAlgorithm2, certificate2, trustChain2 string) {
		client.On("GetChangeHistory", mock.Anything, cps.GetChangeHistoryRequest{
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
		client.On("GetChangeHistory", mock.Anything, cps.GetChangeHistoryRequest{
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
		client.On("GetChangeHistory", mock.Anything, cps.GetChangeHistoryRequest{
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
	mockGetEnrollment = func(client *cps.Mock, enrollmentID, timesToRun int, enrollment *cps.Enrollment) {
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
			EnrollmentID: enrollmentID,
		}).Return(enrollment, nil).Times(timesToRun)
	}

	// mockUploadThirdPartyCertificateAndTrustChain mocks UploadThirdPartyCertificateAndTrustChain call with provided values
	mockUploadThirdPartyCertificateAndTrustChain = func(client *cps.Mock, keyAlgorithm, certificate, trustChain string, enrollmentID, changeID int) {
		client.On("UploadThirdPartyCertAndTrustChain", mock.Anything, cps.UploadThirdPartyCertAndTrustChainRequest{
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
		client.On("UploadThirdPartyCertAndTrustChain", mock.Anything, cps.UploadThirdPartyCertAndTrustChainRequest{
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
		client.On("GetChangePostVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		}).Return(&cps.PostVerificationWarnings{
			Warnings: warnings,
		}, nil).Once()
	}

	// mockEmptyGetPostVerificationWarnings mocks GetPostVerificationWarnings call with empty warnings list
	mockEmptyGetPostVerificationWarnings = func(client *cps.Mock, enrollmentID, changeID int) {
		client.On("GetChangePostVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		}).Return(&cps.PostVerificationWarnings{}, nil).Once()
	}

	// mockAcknowledgePostVerificationWarnings mocks AcknowledgePostVerificationWarnings call with provided values
	mockAcknowledgePostVerificationWarnings = func(client *cps.Mock, enrollmentID, changeID int) {
		client.On("AcknowledgePostVerificationWarnings", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
			EnrollmentID:    enrollmentID,
			ChangeID:        changeID,
		}).Return(nil).Once()
	}

	// mockGetChangeStatus mocks GetChangeStatus call with provided values
	mockGetChangeStatus = func(client *cps.Mock, enrollmentID, changeID, timesToRun int, status string) {
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("AcknowledgeChangeManagement", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
			EnrollmentID:    enrollmentID,
			ChangeID:        changeID,
		}).Return(nil).Once()
	}

	// mockCheckUnacknowledgedWarnings mocks CheckUnacknowledgedWarnings function with provided values
	mockCheckUnacknowledgedWarnings = func(client *cps.Mock, enrollmentID, changeID, timesToRun int, enrollment *cps.Enrollment, status string) {
		mockGetEnrollment(client, enrollmentID, timesToRun, enrollment)
		mockGetChangeStatus(client, enrollmentID, changeID, timesToRun, status)
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

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "enrollment_id", "2"),
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "acknowledge_post_verification_warnings", strconv.FormatBool(data.acknowledgePostVerificationWarnings)),
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "acknowledge_change_management", strconv.FormatBool(data.acknowledgeChangeManagement)),
			resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "wait_for_deployment", strconv.FormatBool(data.waitForDeployment)),
			resource.ComposeAggregateTestCheckFunc(warningsCheck...),
			resource.ComposeAggregateTestCheckFunc(certificateCheck...),
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
}
