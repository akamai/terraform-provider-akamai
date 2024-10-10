package cps

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/providers/cps/tools"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type testDataForCPSCSR struct {
	Enrollment               cps.GetEnrollmentResponse
	EnrollmentID             int
	GetChangeStatusResponse  cps.Change
	GetChangeHistoryResponse *cps.GetChangeHistoryResponse
	ThirdPartyCSRResponse    cps.ThirdPartyCSRResponse
}

var (
	expectReadCPSCSR = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, timesToRun int) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		getEnrollmentRes := data.Enrollment

		changeID, _ := tools.GetChangeIDFromPendingChanges(data.Enrollment.PendingChanges)

		getChangeStatusReq := cps.GetChangeStatusRequest{
			EnrollmentID: data.EnrollmentID,
			ChangeID:     changeID,
		}
		getChangeStatusRes := &data.GetChangeStatusResponse

		getChangeThirdPartyCSRReq := cps.GetChangeRequest{
			EnrollmentID: data.EnrollmentID,
			ChangeID:     changeID,
		}
		getChangeThirdPartyCSRRes := &data.ThirdPartyCSRResponse

		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(&getEnrollmentRes, nil).Times(timesToRun)
		client.On("GetChangeStatus", mock.Anything, getChangeStatusReq).Return(getChangeStatusRes, nil).Times(timesToRun)
		client.On("GetChangeThirdPartyCSR", mock.Anything, getChangeThirdPartyCSRReq).Return(getChangeThirdPartyCSRRes, nil).Times(timesToRun)
	}

	expectReadCPSCSRWithHistory = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, timesToRun int) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		getEnrollmentRes := data.Enrollment

		changeID, _ := tools.GetChangeIDFromPendingChanges(data.Enrollment.PendingChanges)

		getChangeStatusReq := cps.GetChangeStatusRequest{
			EnrollmentID: data.EnrollmentID,
			ChangeID:     changeID,
		}
		getChangeStatusRes := &data.GetChangeStatusResponse

		getChangeHistoryReq := cps.GetChangeHistoryRequest{EnrollmentID: data.EnrollmentID}
		getChangeHistoryRes := cps.GetChangeHistoryResponse{
			Changes: data.GetChangeHistoryResponse.Changes,
		}

		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(&getEnrollmentRes, nil).Times(timesToRun)
		client.On("GetChangeStatus", mock.Anything, getChangeStatusReq).Return(getChangeStatusRes, nil).Times(timesToRun)
		client.On("GetChangeHistory", mock.Anything, getChangeHistoryReq).Return(&getChangeHistoryRes, nil).Times(timesToRun)
	}

	expectReadCPSCSRGetEnrollmentError = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, errorMessage string) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(nil, fmt.Errorf(errorMessage)).Once()
	}

	expectReadDVEnrollment = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(enrollmentDV2, nil).Once()
	}

	expectReadCPSCSRGetThirdPartyError = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, errorMessage string) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		getEnrollmentRes := data.Enrollment

		changeID, _ := tools.GetChangeIDFromPendingChanges(data.Enrollment.PendingChanges)
		getChangeStatusReq := cps.GetChangeStatusRequest{
			EnrollmentID: data.EnrollmentID,
			ChangeID:     changeID,
		}
		getChangeStatusRes := &data.GetChangeStatusResponse

		getChangeThirdPartyCSRReq := cps.GetChangeRequest{
			EnrollmentID: data.EnrollmentID,
			ChangeID:     changeID,
		}

		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(&getEnrollmentRes, nil).Once()
		client.On("GetChangeStatus", mock.Anything, getChangeStatusReq).Return(getChangeStatusRes, nil).Once()
		client.On("GetChangeThirdPartyCSR", mock.Anything, getChangeThirdPartyCSRReq).Return(nil, fmt.Errorf(errorMessage)).Once()
	}

	expectReadCPSCSRNoPendingChanges = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, timesToRun int) {
		data.Enrollment.PendingChanges = []cps.PendingChange{}
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		getEnrollmentRes := data.Enrollment

		getChangeHistoryReq := cps.GetChangeHistoryRequest{EnrollmentID: data.EnrollmentID}
		getChangeHistoryRes := cps.GetChangeHistoryResponse{
			Changes: data.GetChangeHistoryResponse.Changes,
		}
		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(&getEnrollmentRes, nil).Times(timesToRun)
		client.On("GetChangeHistory", mock.Anything, getChangeHistoryReq).Return(&getChangeHistoryRes, nil).Times(timesToRun)
	}

	bothAlgorithmsDataFromCSR = testDataForCPSCSR{
		Enrollment:   *enrollmentThirdParty,
		EnrollmentID: 2,
		GetChangeStatusResponse: cps.Change{
			StatusInfo: &cps.StatusInfo{
				Status: "wait-upload-third-party",
			},
		},
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{
				{
					CSR:          "-----BEGIN CERTIFICATE CSR RSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
				{
					CSR:          "-----BEGIN CERTIFICATE CSR ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "ECDSA",
				},
			},
		},
	}

	bothAlgorithmsDataWithGetChangeHistory = testDataForCPSCSR{
		Enrollment:   *enrollmentThirdParty,
		EnrollmentID: 2,
		GetChangeStatusResponse: cps.Change{
			StatusInfo: &cps.StatusInfo{
				Status: "not-in-list",
			},
		},
		GetChangeHistoryResponse: &cps.GetChangeHistoryResponse{
			Changes: []cps.ChangeHistory{
				{
					Action:            "action1",
					ActionDescription: "description1",
					BusinessCaseID:    "id",
					CreatedBy:         "user",
					CreatedOn:         "",
					LastUpdated:       "",
					MultiStackedCertificates: []cps.CertificateChangeHistory{
						{
							Certificate:  "-----BEGIN CERTIFICATE RSA REQUEST-----\n...RSA...\n-----END CERTIFICATE REQUEST-----",
							TrustChain:   "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							CSR:          "-----BEGIN CERTIFICATE CSR RSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					PrimaryCertificate: cps.CertificateChangeHistory{
						Certificate:  "-----BEGIN CERTIFICATE ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						TrustChain:   "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						CSR:          "-----BEGIN CERTIFICATE CSR ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						KeyAlgorithm: "ECDSA",
					},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "",
				},
			},
		},
	}

	bothAlgorithmsDataWithGetLongerChangeHistory = testDataForCPSCSR{
		Enrollment:   *enrollmentThirdParty,
		EnrollmentID: 2,
		GetChangeStatusResponse: cps.Change{
			StatusInfo: &cps.StatusInfo{
				Status: "not-in-list",
			},
		},
		GetChangeHistoryResponse: &cps.GetChangeHistoryResponse{
			Changes: []cps.ChangeHistory{
				{
					Action:                         "action1",
					ActionDescription:              "description1",
					BusinessCaseID:                 "id",
					CreatedBy:                      "user",
					CreatedOn:                      "",
					LastUpdated:                    "",
					MultiStackedCertificates:       []cps.CertificateChangeHistory{},
					PrimaryCertificate:             cps.CertificateChangeHistory{},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "",
				},
				{
					Action:            "action1",
					ActionDescription: "description1",
					BusinessCaseID:    "id",
					CreatedBy:         "user",
					CreatedOn:         "",
					LastUpdated:       "",
					MultiStackedCertificates: []cps.CertificateChangeHistory{
						{
							Certificate:  "-----BEGIN CERTIFICATE RSA REQUEST-----\n...RSA...\n-----END CERTIFICATE REQUEST-----",
							TrustChain:   "-----BEGIN CERTIFICATE TRUST-CHAIN RSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							CSR:          "-----BEGIN CERTIFICATE CSR RSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
							KeyAlgorithm: "RSA",
						},
					},
					PrimaryCertificate: cps.CertificateChangeHistory{
						Certificate:  "-----BEGIN CERTIFICATE ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						TrustChain:   "-----BEGIN CERTIFICATE TRUST-CHAIN ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						CSR:          "-----BEGIN CERTIFICATE CSR ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						KeyAlgorithm: "ECDSA",
					},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "",
				},
			},
		},
	}

	RSAData = testDataForCPSCSR{
		Enrollment:   *enrollmentThirdParty,
		EnrollmentID: 2,
		GetChangeStatusResponse: cps.Change{
			StatusInfo: &cps.StatusInfo{
				Status: "wait-upload-third-party",
			},
		},
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{
				{
					CSR:          "-----BEGIN CERTIFICATE CSR RSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
			},
		},
	}

	ECDSAData = testDataForCPSCSR{
		Enrollment:   *enrollmentThirdParty,
		EnrollmentID: 2,
		GetChangeStatusResponse: cps.Change{
			StatusInfo: &cps.StatusInfo{
				Status: "wait-upload-third-party",
			},
		},
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{
				{
					CSR:          "-----BEGIN CERTIFICATE CSR ECDSA REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "ECDSA",
				},
			},
		},
	}

	noAlgorithmsData = testDataForCPSCSR{
		Enrollment:   *enrollmentThirdParty,
		EnrollmentID: 1,
		GetChangeStatusResponse: cps.Change{
			StatusInfo: &cps.StatusInfo{
				Status: "wait-upload-third-party",
			},
		},
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{},
		},
	}

	noPendingChanges = testDataForCPSCSR{
		Enrollment:   *enrollmentThirdParty,
		EnrollmentID: 1,
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{},
		},
		GetChangeHistoryResponse: &cps.GetChangeHistoryResponse{
			Changes: []cps.ChangeHistory{
				{
					Action:                   "action1",
					ActionDescription:        "description1",
					BusinessCaseID:           "id",
					CreatedBy:                "user",
					CreatedOn:                "",
					LastUpdated:              "",
					MultiStackedCertificates: nil,
					PrimaryCertificate: cps.CertificateChangeHistory{
						Certificate:  "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						TrustChain:   "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						KeyAlgorithm: "ECDSA",
					},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "",
				},
				{
					Action:                   "action2",
					ActionDescription:        "description2",
					BusinessCaseID:           "id",
					CreatedBy:                "user",
					CreatedOn:                "",
					LastUpdated:              "",
					MultiStackedCertificates: nil,
					PrimaryCertificate: cps.CertificateChangeHistory{
						Certificate:  "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						TrustChain:   "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
						KeyAlgorithm: "RSA",
					},
					PrimaryCertificateOrderDetails: cps.CertificateOrderDetails{},
					RA:                             "",
					Status:                         "",
				},
			},
		},
	}

	dvEnrollment = testDataForCPSCSR{
		Enrollment:   *enrollmentDV2,
		EnrollmentID: 2,
	}
)

func TestDataCPSCSR(t *testing.T) {
	tests := map[string]struct {
		init       func(*testing.T, *cps.Mock, testDataForCPSCSR)
		mockData   testDataForCPSCSR
		configPath string
		error      *regexp.Regexp
	}{
		"happy path with both algorithms with get change": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 3)
			},
			mockData:   bothAlgorithmsDataFromCSR,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"happy path with both algorithms with get change history": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRWithHistory(t, m, testData, 3)
			},
			mockData:   bothAlgorithmsDataWithGetChangeHistory,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"happy path with both algorithms with get longer change history": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRWithHistory(t, m, testData, 3)
			},
			mockData:   bothAlgorithmsDataWithGetLongerChangeHistory,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"happy path with RSA algorithm": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 3)
			},
			mockData:   RSAData,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"happy path with ECDSA algorithm": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 3)
			},
			mockData:   ECDSAData,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"no algorithms": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 3)
			},
			mockData:   noAlgorithmsData,
			configPath: "testdata/TestDataCPSCSR/no_algorithms.tf",
			error:      nil,
		},
		"no pending changes": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRNoPendingChanges(t, m, testData, 3)
			},
			mockData:   noPendingChanges,
			configPath: "testdata/TestDataCPSCSR/no_algorithms.tf",
			error:      nil,
		},
		"enrollment_id not provided": {
			init:       func(_ *testing.T, _ *cps.Mock, _ testDataForCPSCSR) {},
			mockData:   testDataForCPSCSR{},
			configPath: "testdata/TestDataCPSCSR/no_enrollment_id.tf",
			error:      regexp.MustCompile("Missing required argument"),
		},
		"could not fetch enrollment": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRGetEnrollmentError(t, m, testData, "could not get enrollment")
			},
			mockData:   bothAlgorithmsDataFromCSR,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      regexp.MustCompile("could not get enrollment"),
		},
		"could not fetch third party csr": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRGetThirdPartyError(t, m, testData, "could not get third party csr")
			},
			mockData:   bothAlgorithmsDataFromCSR,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      regexp.MustCompile("could not get third party csr"),
		},
		"enrollment is dv": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadDVEnrollment(t, m, testData)
			},
			mockData:   dvEnrollment,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      regexp.MustCompile("given enrollment has non third-party certificate type which is not supported by this data source"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(t, client, test.mockData)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkAttrsForCPSCSR(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkAttrsForCPSCSRFromHistory(data testDataForCPSCSR) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc
	var rsa, ecdsa bool
	certificateFound := false
	for _, change := range data.GetChangeHistoryResponse.Changes {
		for _, certificate := range append(change.MultiStackedCertificates, change.PrimaryCertificate) {
			switch certificate.KeyAlgorithm {
			case "RSA":
				{
					rsa = true
					checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cps_csr.test", "csr_rsa", certificate.CSR))
					certificateFound = true
				}
			case "ECDSA":
				{
					ecdsa = true
					checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cps_csr.test", "csr_ecdsa", certificate.CSR))
					certificateFound = true
				}
			}
		}
		if certificateFound {
			break
		}
	}
	if !rsa {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cps_csr.test", "csr_rsa"))
	}
	if !ecdsa {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cps_csr.test", "csr_ecdsa"))
	}
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func checkAttrsForCPSCSR(data testDataForCPSCSR) resource.TestCheckFunc {
	if data.GetChangeHistoryResponse != nil {
		return checkAttrsForCPSCSRFromHistory(data)
	}
	changeID, _ := tools.GetChangeIDFromPendingChanges(data.Enrollment.PendingChanges)
	var csrECDSA, csrRSA string
	for _, csr := range data.ThirdPartyCSRResponse.CSRs {
		if csr.KeyAlgorithm == "RSA" {
			csrRSA = csr.CSR
		} else if csr.KeyAlgorithm == "ECDSA" {
			csrECDSA = csr.CSR
		}
	}

	var checkFuncs []resource.TestCheckFunc
	if csrECDSA == "" {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cps_csr.test", "csr_ecdsa"))
	} else {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cps_csr.test", "csr_ecdsa", csrECDSA))
	}
	if csrRSA == "" {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cps_csr.test", "csr_rsa"))
	} else {
		resource.TestCheckResourceAttr("data.akamai_cps_csr.test", "csr_rsa", csrRSA)
	}
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cps_csr.test", "id", fmt.Sprintf("%d:%d", data.EnrollmentID, changeID)))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
