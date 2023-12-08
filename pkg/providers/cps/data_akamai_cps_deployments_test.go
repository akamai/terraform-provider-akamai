package cps

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataCPSDeployments(t *testing.T) {
	tests := map[string]struct {
		configPath     string
		checkFunctions []resource.TestCheckFunc
		withError      *regexp.Regexp
		init           func(*cps.Mock)
	}{
		"validate schema with ECDSA primary certificate": {
			configPath: "testdata/TestDataDeployments/deployments.tf",
			init: func(m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: 123,
				}).Return(&cps.GetEnrollmentResponse{
					Location:             "/cps/v2/enrollments/123",
					AutoRenewalStartTime: "2024-10-17T12:08:13Z",
				}, nil)
				m.On("ListDeployments", mock.Anything, cps.ListDeploymentsRequest{
					EnrollmentID: 123,
				}).Return(&cps.ListDeploymentsResponse{
					Production: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "ECDSA certificate as string on production",
							KeyAlgorithm: "ECDSA",
							Expiry:       "2024-12-16T12:08:13Z",
						},
						MultiStackedCertificates: []cps.DeploymentCertificate{
							{
								Certificate:  "RSA certificate as string on production",
								KeyAlgorithm: "RSA",
								Expiry:       "2024-12-16T12:08:13Z",
							},
						},
					},
					Staging: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "ECDSA certificate as string on staging",
							KeyAlgorithm: "ECDSA",
						},
						MultiStackedCertificates: []cps.DeploymentCertificate{
							{
								Certificate:  "RSA certificate as string on staging",
								KeyAlgorithm: "RSA",
								Expiry:       "2024-12-16T12:08:13Z",
							},
						},
					},
				}, nil)
			},
			checkFunctions: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "enrollment_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_rsa", "RSA certificate as string on production"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_ecdsa", "ECDSA certificate as string on production"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_rsa", "RSA certificate as string on staging"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_ecdsa", "ECDSA certificate as string on staging"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "expiry_date", "2024-12-16T12:08:13Z"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "auto_renewal_start_time", "2024-10-17T12:08:13Z"),
			},
		},
		"validate schema with RSA primary certificate": {
			configPath: "testdata/TestDataDeployments/deployments.tf",
			init: func(m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: 123,
				}).Return(&cps.GetEnrollmentResponse{
					Location:             "/cps/v2/enrollments/123",
					AutoRenewalStartTime: "2024-10-17T12:08:13Z",
				}, nil)
				m.On("ListDeployments", mock.Anything, cps.ListDeploymentsRequest{
					EnrollmentID: 123,
				}).Return(&cps.ListDeploymentsResponse{
					Production: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "RSA certificate as string on production",
							KeyAlgorithm: "RSA",
							Expiry:       "2024-12-16T12:08:13Z",
						},
						MultiStackedCertificates: []cps.DeploymentCertificate{
							{
								Certificate:  "ECDSA certificate as string on production",
								KeyAlgorithm: "ECDSA",
								Expiry:       "2024-12-16T12:08:13Z",
							},
						},
					},
					Staging: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "RSA certificate as string on staging",
							KeyAlgorithm: "RSA",
						},
						MultiStackedCertificates: []cps.DeploymentCertificate{
							{
								Certificate:  "ECDSA certificate as string on staging",
								KeyAlgorithm: "ECDSA",
								Expiry:       "2024-12-16T12:08:13Z",
							},
						},
					},
				}, nil)
			},
			checkFunctions: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "enrollment_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_rsa", "RSA certificate as string on production"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_ecdsa", "ECDSA certificate as string on production"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_rsa", "RSA certificate as string on staging"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_ecdsa", "ECDSA certificate as string on staging"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "expiry_date", "2024-12-16T12:08:13Z"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "auto_renewal_start_time", "2024-10-17T12:08:13Z"),
			},
		},
		"no RSA MultiStackedCertificate": {
			configPath: "testdata/TestDataDeployments/deployments.tf",
			init: func(m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: 123,
				}).Return(&cps.GetEnrollmentResponse{
					Location:             "/cps/v2/enrollments/123",
					AutoRenewalStartTime: "2024-10-17T12:08:13Z",
				}, nil)
				m.On("ListDeployments", mock.Anything, cps.ListDeploymentsRequest{
					EnrollmentID: 123,
				}).Return(&cps.ListDeploymentsResponse{
					Production: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "ECDSA certificate as string on production",
							KeyAlgorithm: "ECDSA",
							Expiry:       "2024-12-16T12:08:13Z",
						},
					},
					Staging: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "ECDSA certificate as string on staging",
							KeyAlgorithm: "ECDSA",
						},
					},
				}, nil)
			},
			checkFunctions: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "enrollment_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_rsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_ecdsa", "ECDSA certificate as string on production"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_rsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_ecdsa", "ECDSA certificate as string on staging"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "expiry_date", "2024-12-16T12:08:13Z"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "auto_renewal_start_time", "2024-10-17T12:08:13Z"),
			},
		},
		"no ECDSA MultiStackedCertificate": {
			configPath: "testdata/TestDataDeployments/deployments.tf",
			init: func(m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: 123,
				}).Return(&cps.GetEnrollmentResponse{
					Location:             "/cps/v2/enrollments/123",
					AutoRenewalStartTime: "2024-10-17T12:08:13Z",
				}, nil)
				m.On("ListDeployments", mock.Anything, cps.ListDeploymentsRequest{
					EnrollmentID: 123,
				}).Return(&cps.ListDeploymentsResponse{
					Production: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "RSA certificate as string on production",
							KeyAlgorithm: "RSA",
							Expiry:       "2024-12-16T12:08:13Z",
						},
					},
					Staging: &cps.Deployment{
						PrimaryCertificate: cps.DeploymentCertificate{
							Certificate:  "RSA certificate as string on staging",
							KeyAlgorithm: "RSA",
						},
					},
				}, nil)
			},
			checkFunctions: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "enrollment_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_rsa", "RSA certificate as string on production"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_ecdsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_rsa", "RSA certificate as string on staging"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_ecdsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "expiry_date", "2024-12-16T12:08:13Z"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "auto_renewal_start_time", "2024-10-17T12:08:13Z"),
			},
		},
		"enrollment not found": {
			configPath: "testdata/TestDataDeployments/deployments.tf",
			init: func(m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: 123,
				}).Return(nil, &cps.Error{
					Type:       "not-found",
					Title:      "Not Found",
					StatusCode: http.StatusNotFound,
				})
			},
			withError: regexp.MustCompile("could not fetch enrollment by id 123"),
		},
		"no deployed certificates": {
			configPath: "testdata/TestDataDeployments/deployments.tf",
			init: func(m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: 123,
				}).Return(&cps.GetEnrollmentResponse{
					Location:             "/cps/v2/enrollments/123",
					AutoRenewalStartTime: "2024-10-17T12:08:13Z",
				}, nil)
				m.On("ListDeployments", mock.Anything, cps.ListDeploymentsRequest{
					EnrollmentID: 123,
				}).Return(&cps.ListDeploymentsResponse{
					Production: nil,
					Staging:    nil,
				}, nil)
			},
			checkFunctions: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "enrollment_id", "123"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "auto_renewal_start_time", "2024-10-17T12:08:13Z"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_rsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "production_certificate_ecdsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_rsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "staging_certificate_ecdsa.#", "0"),
				resource.TestCheckResourceAttr("data.akamai_cps_deployments.test", "expiry_date.#", "0"),
			},
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       resource.ComposeAggregateTestCheckFunc(test.checkFunctions...),
							ExpectError: test.withError,
						},
					},
				})
				client.AssertExpectations(t)
			})
		})
	}
}
