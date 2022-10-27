package cps

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/cps/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

type testDataForCPSCSR struct {
	Enrollment               cps.Enrollment
	EnrollmentID             int
	ThirdPartyCSRResponse    cps.ThirdPartyCSRResponse
	GetChangeHistoryResponse *cps.GetChangeHistoryResponse
}

var (
	expectReadCPSCSR = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, timesToRun int) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		getEnrollmentRes := &data.Enrollment

		changeID, _ := tools.GetChangeIDFromPendingChanges(data.Enrollment.PendingChanges)
		getChangeThirdPartyCSRReq := cps.GetChangeRequest{
			EnrollmentID: data.EnrollmentID,
			ChangeID:     changeID,
		}
		getChangeThirdPartyCSRRes := &data.ThirdPartyCSRResponse

		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(getEnrollmentRes, nil).Times(timesToRun)
		client.On("GetChangeThirdPartyCSR", mock.Anything, getChangeThirdPartyCSRReq).Return(getChangeThirdPartyCSRRes, nil).Times(timesToRun)
	}

	expectReadCPSCSRGetEnrollmentError = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, errorMessage string) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(nil, fmt.Errorf(errorMessage)).Once()
	}

	expectReadCPSCSRGetThirdPartyError = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, errorMessage string) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		getEnrollmentRes := &data.Enrollment

		changeID, _ := tools.GetChangeIDFromPendingChanges(data.Enrollment.PendingChanges)
		getChangeThirdPartyCSRReq := cps.GetChangeRequest{
			EnrollmentID: data.EnrollmentID,
			ChangeID:     changeID,
		}

		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(getEnrollmentRes, nil).Once()
		client.On("GetChangeThirdPartyCSR", mock.Anything, getChangeThirdPartyCSRReq).Return(nil, fmt.Errorf(errorMessage)).Once()
	}

	expectReadCPSCSRNoPendingChanges = func(t *testing.T, client *cps.Mock, data testDataForCPSCSR, timesToRun int) {
		getEnrollmentReq := cps.GetEnrollmentRequest{
			EnrollmentID: data.EnrollmentID,
		}
		getEnrollmentRes := &data.Enrollment

		getChangeHistoryReq := cps.GetChangeHistoryRequest{EnrollmentID: data.EnrollmentID}
		getChangeHistoryRes := cps.GetChangeHistoryResponse{
			Changes: data.GetChangeHistoryResponse.Changes,
		}
		client.On("GetEnrollment", mock.Anything, getEnrollmentReq).Return(getEnrollmentRes, nil).Times(timesToRun)
		client.On("GetChangeHistory", mock.Anything, getChangeHistoryReq).Return(&getChangeHistoryRes, nil).Times(timesToRun)
	}

	bothAlgorithmsData = testDataForCPSCSR{
		Enrollment:   *enrollmentDV2,
		EnrollmentID: 2,
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{
				{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
				{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "ECDSA",
				},
			},
		},
	}

	RSAData = testDataForCPSCSR{
		Enrollment:   *enrollmentDV2,
		EnrollmentID: 2,
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{
				{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "RSA",
				},
			},
		},
	}

	ECDSAData = testDataForCPSCSR{
		Enrollment:   *enrollmentDV2,
		EnrollmentID: 2,
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{
				{
					CSR:          "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
					KeyAlgorithm: "ECDSA",
				},
			},
		},
	}

	noAlgorithmsData = testDataForCPSCSR{
		Enrollment:   *enrollmentDV2,
		EnrollmentID: 1,
		ThirdPartyCSRResponse: cps.ThirdPartyCSRResponse{
			CSRs: []cps.CertSigningRequest{},
		},
	}

	noPendingChanges = testDataForCPSCSR{
		Enrollment:   *enrollmentDV1,
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
)

func TestDataCPSCSR(t *testing.T) {
	tests := map[string]struct {
		init       func(*testing.T, *cps.Mock, testDataForCPSCSR)
		mockData   testDataForCPSCSR
		configPath string
		error      *regexp.Regexp
	}{
		"happy path with both algorithms": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 5)
			},
			mockData:   bothAlgorithmsData,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"happy path with RSA algorithm": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 5)
			},
			mockData:   RSAData,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"happy path with ECDSA algorithm": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 5)
			},
			mockData:   ECDSAData,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      nil,
		},
		"no algorithms": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSR(t, m, testData, 5)
			},
			mockData:   noAlgorithmsData,
			configPath: "testdata/TestDataCPSCSR/no_algorithms.tf",
			error:      nil,
		},
		"no pending changes": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRNoPendingChanges(t, m, testData, 5)
			},
			mockData:   noPendingChanges,
			configPath: "testdata/TestDataCPSCSR/no_algorithms.tf",
			error:      nil,
		},
		"enrollment_id not provided": {
			init:       func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {},
			mockData:   testDataForCPSCSR{},
			configPath: "testdata/TestDataCPSCSR/no_enrollment_id.tf",
			error:      regexp.MustCompile("Missing required argument"),
		},
		"could not fetch enrollment": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRGetEnrollmentError(t, m, testData, "could not get enrollment")
			},
			mockData:   bothAlgorithmsData,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      regexp.MustCompile("could not get enrollment"),
		},
		"could not fetch third party csr": {
			init: func(t *testing.T, m *cps.Mock, testData testDataForCPSCSR) {
				expectReadCPSCSRGetThirdPartyError(t, m, testData, "could not get third party csr")
			},
			mockData:   bothAlgorithmsData,
			configPath: "testdata/TestDataCPSCSR/default.tf",
			error:      regexp.MustCompile("could not get third party csr"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(t, client, test.mockData)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString(test.configPath),
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
	switch data.GetChangeHistoryResponse.Changes[0].PrimaryCertificate.KeyAlgorithm {
	case "RSA":
		{
			checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cps_csr.test", "csr_ecdsa"))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(
				"data.akamai_cps_csr.test", "csr_rsa", data.GetChangeHistoryResponse.Changes[0].PrimaryCertificate.CSR))
		}
	case "ECDSA":
		{
			checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cps_csr.test", "csr_rsa"))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(
				"data.akamai_cps_csr.test", "csr_ecdsa", data.GetChangeHistoryResponse.Changes[0].PrimaryCertificate.CSR))
		}
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
