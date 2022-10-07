package cps

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceCPSUploadCertificate(t *testing.T) {
	tests := map[string]struct {
		init                func(*testing.T, *mockcps, *cps.Enrollment, int, int)
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
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)
				mockUpdate(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:          enrollmentDV2,
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			error:               nil,
		},
		"create with ch-mgmt true, update to false": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)
				mockUpdate(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 1)
			},
			enrollment:          enrollmentDV2,
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
			client := &mockcps{}
			test.init(t, client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
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

func TestCreateCPSUploadCertificate(t *testing.T) {
	tests := map[string]struct {
		init         func(*testing.T, *mockcps, *cps.Enrollment, int, int)
		enrollment   *cps.Enrollment
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"successful create - RSA cert": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			error:        nil,
		},
		"successful create - ECDSA cert, without trust chain": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, ECDSA, certECDSAForTests, "", enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certECDSAForTests, "", ECDSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_ecdsa.tf",
			checkFunc:    checkAttrs(createMockData(certECDSAForTests, "", "", "", false, true, false, false, nil)),
			error:        nil,
		},
		"successful create - both cert and trust chains": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadBothThirdPartyCertificateAndTrustChain(m,
					ECDSA,
					certECDSAForTests,
					trustChainECDSAForTests,
					RSA,
					certRSAForTests,
					trustChainRSAForTests,
					enrollmentID,
					changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certECDSAForTests, "", ECDSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_both_certificates.tf",
			checkFunc:    checkAttrs(createMockData(certECDSAForTests, trustChainECDSAForTests, certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			error:        nil,
		},
		"create: auto_approve_warnings match": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, true, false, []string{"CERTIFICATE_ADDED_TO_TRUST_CHAIN", "CERTIFICATE_ALREADY_LOADED", "CERTIFICATE_DATA_BLANK_OR_MISSING"})),
			error:        nil,
		},
		"create: auto_approve_warnings missing warnings error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockGetPostVerificationWarnings(m, blankAndNullWarnings, enrollmentID, changeID)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings.tf",
			error:        regexp.MustCompile(`Error: could not process post verification warnings: not every warning has been acknowledged: warnings cannot be approved: CERTIFICATE_HAS_NULL_ISSUER`),
		},
		"create: auto_approve_warnings not provided and empty warning list": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithEmptyWarningList(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        nil,
		},
		"required attribute not provided": {
			init:       func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, _, _ int) {},
			enrollment: nil,
			configPath: "testdata/TestResCPSUploadCertificate/certificates/no_certificates.tf",
			checkFunc:  nil,
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"create: auto_approve_warnings not provided and not empty warning list": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockGetPostVerificationWarnings(m, noKMIDataWarning, enrollmentID, changeID)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        regexp.MustCompile(`Error: could not process post verification warnings: not every warning has been acknowledged: warnings cannot be approved: CERTIFICATE_KMI_DATA_MISSING`),
		},
		"create: auto_approve_warnings empty list and warnings": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithEmptyWarningList(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_empty.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, true, false, nil)),
			error:        nil,
		},
		"create: change management wrong type": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.ChangeManagement = true
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockEmptyGetPostVerificationWarnings(m, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/auto_approve_warnings/auto_approve_warnings_not_provided.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, false, false, false, nil)),
			error:        nil,
		},
		"create: change management set to false or not specified": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.ChangeManagement = true
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockEmptyGetPostVerificationWarnings(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_not_specified.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, false, false, false, []string{"Warning 1", "Warning 2", "Warning 3"})),
			error:        nil,
		},
		"create: trust chain without certificate": {
			init:         func(_ *testing.T, _ *mockcps, _ *cps.Enrollment, _, _ int) {},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/trust_chain_without_cert.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("provided ECDSA trust chain without ECDSA certificate. Please remove it or add a certificate"),
		},
		"create: get enrollment error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, _ int) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollmentID,
				}).Return(nil, fmt.Errorf("could not get an erollments")).Once()
			},
			enrollment:   nil,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get an enrollment"),
		},
		"create: upload third party cert error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}

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
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not upload a certificate"),
		},
		"create: get change post verification warnings error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)

				m.On("GetChangePostVerificationWarnings", mock.Anything, cps.GetChangeRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				}).Return(nil, fmt.Errorf("could not get change post verification warnings")).Once()
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get change post verification warnings"),
		},
		"create: acknowledge post verification warnings error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockGetPostVerificationWarnings(m, "Some warning", enrollmentID, changeID)

				m.On("AcknowledgePostVerificationWarnings", mock.Anything, cps.AcknowledgementRequest{
					Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
					EnrollmentID:    enrollmentID,
					ChangeID:        changeID,
				}).Return(fmt.Errorf("could not acknowledge post verification warnings")).Once()
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/certificates/create_certificate_rsa.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not acknowledge post verification warnings"),
		},
		"create: get change status error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)

				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				}).Return(nil, fmt.Errorf("could not get change status")).Once()
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get change status"),
		},
		"create: acknowledge change management error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)

				m.On("AcknowledgeChangeManagement", mock.Anything, cps.AcknowledgementRequest{
					Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
					EnrollmentID:    enrollmentID,
					ChangeID:        changeID,
				}).Return(fmt.Errorf("could not acknowledge change management")).Once()
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not acknowledge change management"),
		},
		"create: no pending changes error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.ChangeManagement = true
				enrollment.PendingChanges = []string{}
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("Error: could not get change ID: no pending changes were found on enrollment"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockcps{}
			test.init(t, client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
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
		init         func(*testing.T, *mockcps, *cps.Enrollment, int, int)
		enrollment   *cps.Enrollment
		enrollmentID int
		changeID     int
		configPath   string
		checkFunc    resource.TestCheckFunc
		error        *regexp.Regexp
	}{
		"read: get certificate history": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 2, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 2, completeStatus)
				mockGetCertificateHistory(m, enrollmentID, 2, certRSAForTests, trustChainRSAForTests, RSA)
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_true.tf",
			checkFunc:    checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, true, nil)),
			error:        nil,
		},
		"read: get certificate history error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, completeStatus)

				m.On("GetCertificateHistory", mock.Anything, cps.GetCertificateHistoryRequest{
					EnrollmentID: enrollmentID,
				}).Return(nil, fmt.Errorf("could not get certificate history")).Once()
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_true.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get certificate history"),
		},
		"read: get change history error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				m.On("GetChangeHistory", mock.Anything, cps.GetChangeHistoryRequest{
					EnrollmentID: enrollmentID,
				}).Return(nil, fmt.Errorf("could not get certificate history")).Once()
			},
			enrollment:   enrollmentDV2,
			enrollmentID: 2,
			changeID:     22,
			configPath:   "testdata/TestResCPSUploadCertificate/wait_for_deployment/wait_for_deployment_false.tf",
			checkFunc:    nil,
			error:        regexp.MustCompile("could not get certificate history"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockcps{}
			test.init(t, client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
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
		init                func(*testing.T, *mockcps, *cps.Enrollment, int, int)
		enrollment          *cps.Enrollment
		enrollmentID        int
		changeID            int
		configPathForCreate string
		configPathForUpdate string
		checkFuncForCreate  resource.TestCheckFunc
		checkFuncForUpdate  resource.TestCheckFunc
		errorForCreate      *regexp.Regexp
		errorForUpdate      *regexp.Regexp
	}{
		"update: ignore change to acknowledge post verification warnings flag": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:          enrollmentDV2,
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
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
				mockUploadThirdPartyCertificateAndTrustChain(m, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
				mockGetPostVerificationWarnings(m, threeWarnings, enrollmentID, changeID)
				mockAcknowledgePostVerificationWarnings(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 2)
			},
			enrollment:          enrollmentDV2,
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
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 4)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
			},
			enrollment:          enrollmentDV2,
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_false.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			checkFuncForUpdate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      nil,
		},
		"update: change in already deployed certificate - error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
			},
			enrollment:          enrollmentDV2,
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/certificates/changed_certificate.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("Error: cannot make changes to certificate that is already on staging and/or production network, need to create new enrollment"),
		},
		"update: get enrollment error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)
				mockAcknowledgeChangeManagement(m, enrollmentID, changeID)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)

				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollmentID,
				}).Return(nil, fmt.Errorf("could not get an enrollment")).Times(1)
			},
			enrollment:          enrollmentDV2,
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_false.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, true, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("could not get an enrollment"),
		},
		"update: get change status error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)

				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				}).Return(nil, fmt.Errorf("could not get change status")).Once()
			},
			enrollment:          enrollmentDV2,
			enrollmentID:        2,
			changeID:            22,
			configPathForCreate: "testdata/TestResCPSUploadCertificate/change_management/change_management_false.tf",
			configPathForUpdate: "testdata/TestResCPSUploadCertificate/change_management/change_management_true.tf",
			checkFuncForCreate:  checkAttrs(createMockData("", "", certRSAForTests, trustChainRSAForTests, false, true, false, false, nil)),
			errorForCreate:      nil,
			errorForUpdate:      regexp.MustCompile("could not get change status"),
		},
		"update: acknowledge change management error": {
			init: func(t *testing.T, m *mockcps, enrollment *cps.Enrollment, enrollmentID, changeID int) {
				enrollment.ChangeManagement = true
				mockCreateWithACKPostWarnings(m, enrollmentID, changeID, enrollment)
				mockReadGetChangeHistory(m, enrollment, certRSAForTests, trustChainRSAForTests, RSA, enrollmentID, 3)
				mockGetEnrollment(m, enrollmentID, 1, enrollment)
				mockGetChangeStatus(m, enrollmentID, changeID, 1, waitAckChMgmtStatus)

				m.On("AcknowledgeChangeManagement", mock.Anything, cps.AcknowledgementRequest{
					Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
					EnrollmentID:    enrollmentID,
					ChangeID:        changeID,
				}).Return(fmt.Errorf("could not acknowledge change management")).Once()
			},
			enrollment:          enrollmentDV2,
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
			client := &mockcps{}
			test.init(t, client, test.enrollment, test.enrollmentID, test.changeID)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
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

var (
	certECDSAForTests          = "-----BEGIN CERTIFICATE ECDSA REQUEST-----\n...\n-----END CERTIFICATE ECDSA REQUEST-----"
	certRSAForTests            = "-----BEGIN CERTIFICATE RSA REQUEST-----\n...\n-----END CERTIFICATE RSA REQUEST-----"
	trustChainRSAForTests      = "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST-----\n...\n-----END CERTIFICATE TRUST-CHAIN RSA REQUEST-----"
	trustChainECDSAForTests    = "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----\n...\n-----END CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----"
	waitAckChMgmtStatus        = "wait-ack-change-management"
	waitReviewThirdPartyStatus = "wait-review-third-party-cert"
	completeStatus             = "complete"
	RSA                        = "RSA"
	ECDSA                      = "ECDSA"
	blankAndNullWarnings       = "Certificate data is blank or missing.\nCertificate has a null issuer"
	noKMIDataWarning           = "No KMI data is available for the new certificate."
	threeWarnings              = "Certificate Added to the new Trust Chain: TEST\nThere is a problem deploying the 'RSA' certificate.  Please contact your Akamai support team to resolve the issue.\nCertificate data is blank or missing."

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
	mockCreateWithEmptyWarningList = func(client *mockcps, enrollmentID, changeID int, enrollment *cps.Enrollment) {
		mockGetEnrollment(client, enrollmentID, 1, enrollment)
		enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/2/changes/%d", changeID)}
		mockUploadThirdPartyCertificateAndTrustChain(client, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
		mockEmptyGetPostVerificationWarnings(client, enrollmentID, changeID)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitAckChMgmtStatus)
		mockAcknowledgeChangeManagement(client, enrollmentID, changeID)
	}

	// mockUpdate mocks default approach when updating the resource
	mockUpdate = func(client *mockcps, enrollmentID, changeID int, enrollment *cps.Enrollment) {
		mockGetEnrollment(client, enrollmentID, 1, enrollment)
		enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/2/changes/%d", changeID)}
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitAckChMgmtStatus)
		mockAcknowledgeChangeManagement(client, enrollmentID, changeID)
	}

	// mockCreateWithACKPostWarnings mocks acknowledging post verification warnings along with creation of the resource
	mockCreateWithACKPostWarnings = func(client *mockcps, enrollmentID, changeID int, enrollment *cps.Enrollment) {
		mockGetEnrollment(client, enrollmentID, 1, enrollment)
		enrollment.PendingChanges = []string{fmt.Sprintf("/cps/v2/enrollments/%d/changes/%d", enrollmentID, changeID)}
		mockUploadThirdPartyCertificateAndTrustChain(client, RSA, certRSAForTests, trustChainRSAForTests, enrollmentID, changeID)
		mockGetChangeStatus(client, enrollmentID, changeID, 1, waitReviewThirdPartyStatus)
		mockGetPostVerificationWarnings(client, "Certificate Added to the new Trust Chain: TEST\nThere is a problem deploying the 'RSA' certificate.  Please contact your Akamai support team to resolve the issue.\nCertificate data is blank or missing.", enrollmentID, changeID)
		mockAcknowledgePostVerificationWarnings(client, enrollmentID, changeID)
	}

	// mockGetCertificateHistory mocks GetCertificateHistory call with provided values
	mockGetCertificateHistory = func(client *mockcps, enrollmentID, timesToRun int, certificate, trustChain, keyAlgorithm string) {
		client.On("GetCertificateHistory", mock.Anything, cps.GetCertificateHistoryRequest{
			EnrollmentID: enrollmentID,
		}).Return(&cps.GetCertificateHistoryResponse{
			Certificates: []cps.HistoryCertificate{
				{
					DeploymentStatus:         "",
					Geography:                "",
					MultiStackedCertificates: nil,
					PrimaryCertificate: cps.CertificateObject{
						Certificate:  certificate,
						Expiry:       "",
						KeyAlgorithm: keyAlgorithm,
						TrustChain:   trustChain,
					},
					RA:            "",
					Slots:         nil,
					StagingStatus: "inactive",
					Type:          "",
				},
				{
					DeploymentStatus:         "",
					Geography:                "",
					MultiStackedCertificates: nil,
					PrimaryCertificate: cps.CertificateObject{
						Certificate:  certRSAForTests,
						Expiry:       "",
						KeyAlgorithm: RSA,
						TrustChain:   trustChainRSAForTests,
					},
					RA:            "",
					Slots:         nil,
					StagingStatus: "active",
					Type:          "",
				},
			},
		}, nil).Times(timesToRun)
	}

	// mockReadGetChangeHistory mocks GetChangeHistory call with provided values
	mockReadGetChangeHistory = func(client *mockcps, enrollment *cps.Enrollment, certificate, trustChain, keyAlgorithm string, enrollmentID, timesToRun int) {
		mockGetEnrollment(client, enrollmentID, timesToRun, enrollment)
		mockGetChangeHistory(client, enrollmentID, timesToRun, enrollment, keyAlgorithm, certificate, trustChain)
	}

	// mockGetChangeHistory mocks GetChangeHistory call with provided data
	mockGetChangeHistory = func(client *mockcps, enrollmentID, timesToRun int, enrollment *cps.Enrollment, keyAlgorithm, certificate, trustChain string) {
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

	// mockGetEnrollment mocks GetEnrollment call with provided values
	mockGetEnrollment = func(client *mockcps, enrollmentID, timesToRun int, enrollment *cps.Enrollment) {
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
			EnrollmentID: enrollmentID,
		}).Return(enrollment, nil).Times(timesToRun)
	}

	// mockUploadThirdPartyCertificateAndTrustChain mocks UploadThirdPartyCertificateAndTrustChain call with provided values
	mockUploadThirdPartyCertificateAndTrustChain = func(client *mockcps, keyAlgorithm, certificate, trustChain string, enrollmentID, changeID int) {
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
	mockUploadBothThirdPartyCertificateAndTrustChain = func(client *mockcps, keyAlgorithm, certificate, trustChain, keyAlgorithm2, certificate2, trustChain2 string, enrollmentID, changeID int) {
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
	mockGetPostVerificationWarnings = func(client *mockcps, warnings string, enrollmentID, changeID int) {
		client.On("GetChangePostVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		}).Return(&cps.PostVerificationWarnings{
			Warnings: warnings,
		}, nil).Once()
	}

	// mockEmptyGetPostVerificationWarnings mocks GetPostVerificationWarnings call with empty warnings list
	mockEmptyGetPostVerificationWarnings = func(client *mockcps, enrollmentID, changeID int) {
		client.On("GetChangePostVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		}).Return(&cps.PostVerificationWarnings{}, nil).Once()
	}

	// mockAcknowledgePostVerificationWarnings mocks AcknowledgePostVerificationWarnings call with provided values
	mockAcknowledgePostVerificationWarnings = func(client *mockcps, enrollmentID, changeID int) {
		client.On("AcknowledgePostVerificationWarnings", mock.Anything, cps.AcknowledgementRequest{
			Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
			EnrollmentID:    enrollmentID,
			ChangeID:        changeID,
		}).Return(nil).Once()
	}

	// mockGetChangeStatus mocks GetChangeStatus call with provided values
	mockGetChangeStatus = func(client *mockcps, enrollmentID, changeID, timesToRun int, status string) {
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
	mockAcknowledgeChangeManagement = func(client *mockcps, enrollmentID, changeID int) {
		client.On("AcknowledgeChangeManagement", mock.Anything, cps.AcknowledgementRequest{
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
		if data.certificateECDSA != "" {
			certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "certificate_ecdsa_pem", data.certificateECDSA))
		} else {
			certificateCheck = append(certificateCheck, resource.TestCheckNoResourceAttr("akamai_cps_upload_certificate.test", "certificate_ecdsa_pem"))
		}
		if data.certificateRSA != "" {
			certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "certificate_rsa_pem", data.certificateRSA))
		} else {
			certificateCheck = append(certificateCheck, resource.TestCheckNoResourceAttr("akamai_cps_upload_certificate.test", "certificate_rsa_pem"))
		}
		if data.trustChainECDSA != "" {
			certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "trust_chain_ecdsa_pem", data.trustChainECDSA))
		} else {
			certificateCheck = append(certificateCheck, resource.TestCheckNoResourceAttr("akamai_cps_upload_certificate.test", "trust_chain_ecdsa_pem"))
		}
		if data.trustChainRSA != "" {
			certificateCheck = append(certificateCheck, resource.TestCheckResourceAttr("akamai_cps_upload_certificate.test", "trust_chain_rsa_pem", data.trustChainRSA))
		} else {
			certificateCheck = append(certificateCheck, resource.TestCheckNoResourceAttr("akamai_cps_upload_certificate.test", "trust_chain_rsa_pem"))
		}
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
