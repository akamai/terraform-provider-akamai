package cps

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceDVEnrollment(t *testing.T) {
	t.Parallel()
	t.Run("lifecycle test", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:                   "US",
				CN:                  "test.akamai.com",
				L:                   "Cambridge",
				O:                   "Akamai",
				OU:                  "WebEx",
				SANS:                []string{"san.test.akamai.com"},
				ST:                  "MA",
				PreferredTrustChain: "intermediate-a",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(3)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(3)

		var enrollmentUpdate cps.GetEnrollmentResponse
		err := copier.CopyWithOption(&enrollmentUpdate, enrollment, copier.Option{DeepCopy: true})
		require.NoError(t, err)
		enrollmentUpdate.AdminContact.FirstName = "R5"
		enrollmentUpdate.AdminContact.LastName = "D5"
		enrollmentUpdate.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
		enrollmentUpdate.CSR.PreferredTrustChain = "dst-root-ca-x3"
		enrollmentUpdate.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com"}
		enrollmentUpdate.Location = ""
		enrollmentUpdate.PendingChanges = nil

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentUpdate.Location = "/cps/v2/enrollments/1"
		enrollmentUpdate.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdate, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Twice()
		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san2.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san2.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san2.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/lifecycle/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "2"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "2"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.san.test.akamai.com,_acme-challenge.test.akamai.com"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "timeouts.0.default", "2h"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/lifecycle/update_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "3"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "3"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "timeouts.#", "0"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.san.test.akamai.com,_acme-challenge.san2.test.akamai.com,_acme-challenge.test.akamai.com"),
					),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test, remove san, returns 'wait-review-cert-warning' status", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:                   "US",
				CN:                  "test.akamai.com",
				L:                   "Cambridge",
				O:                   "Akamai",
				OU:                  "WebEx",
				SANS:                []string{"san.test.akamai.com"},
				ST:                  "MA",
				PreferredTrustChain: "intermediate-a",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(3)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(3)

		var enrollmentUpdate cps.GetEnrollmentResponse
		err := copier.CopyWithOption(&enrollmentUpdate, enrollment, copier.Option{DeepCopy: true})
		require.NoError(t, err)
		enrollmentUpdate.AdminContact.FirstName = "R1"
		enrollmentUpdate.AdminContact.LastName = "D1"
		enrollmentUpdate.CSR.SANS = nil
		enrollmentUpdate.CSR.PreferredTrustChain = ""
		enrollmentUpdate.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com"}
		enrollmentUpdate.Location = ""
		enrollmentUpdate.PendingChanges = nil

		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)
		allowCancel := true
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollmentUpdate.Location = "/cps/v2/enrollments/1"
		enrollmentUpdate.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdate, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewCertWarning,
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewCertWarning,
			},
		}, nil).Twice()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/lifecycle/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "2"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "2"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.san.test.akamai.com,_acme-challenge.test.akamai.com"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "timeouts.0.default", "2h"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/empty_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "timeouts.#", "0"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.test.akamai.com"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test transitions DNS names settings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		allowCancel := true

		enrollment := getTestDVEnrollment()
		enrollment.CSR.SANS = []string{"san.test.akamai.com"}
		enrollment.CSR.PreferredTrustChain = "intermediate-a"
		enrollment.NetworkConfiguration.DisallowedTLSVersions = []string{"TLSv1", "TLSv1_1"}
		enrollment.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: createEnrollmentReqBodyFromEnrollment(enrollment),
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{{
			Location:   "/cps/v2/enrollments/1/changes/2",
			ChangeType: "new-certificate",
		}}
		currentEnrollment := enrollment
		currentEnrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
			DNSNames:      []string{"test.akamai.com", "san.test.akamai.com"},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&currentEnrollment, nil).Times(15)
		mockDVTransitionChangeStatus(client)

		enrollmentWithExplicitDNSNames := getTestDVEnrollment()
		require.NoError(t, copier.CopyWithOption(&enrollmentWithExplicitDNSNames, enrollment, copier.Option{DeepCopy: true}))
		enrollmentWithExplicitDNSNames.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: false,
			DNSNames:      []string{"test.akamai.com"},
		}
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     createEnrollmentReqBodyFromEnrollment(enrollmentWithExplicitDNSNames),
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Run(func(mock.Arguments) {
			currentEnrollment = enrollmentWithExplicitDNSNames
		}).Once()

		enrollmentWithAdditionalExplicitDNSName := getTestDVEnrollment()
		require.NoError(t, copier.CopyWithOption(&enrollmentWithAdditionalExplicitDNSName, enrollmentWithExplicitDNSNames, copier.Option{DeepCopy: true}))
		enrollmentWithAdditionalExplicitDNSName.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com", "san.test.akamai.com"}
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     createEnrollmentReqBodyFromEnrollment(enrollmentWithAdditionalExplicitDNSName),
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Run(func(mock.Arguments) {
			currentEnrollment = enrollmentWithAdditionalExplicitDNSName
		}).Once()

		enrollmentWithAllSANs := getTestDVEnrollment()
		require.NoError(t, copier.CopyWithOption(&enrollmentWithAllSANs, enrollment, copier.Option{DeepCopy: true}))
		enrollmentWithAllSANs.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
			DNSNames:      []string{"test.akamai.com", "san.test.akamai.com"},
		}
		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentRequestBody:     createEnrollmentReqBodyFromEnrollment(enrollmentWithAllSANs),
				EnrollmentID:              1,
				AllowCancelPendingChanges: &allowCancel,
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Run(func(mock.Arguments) {
			currentEnrollment = enrollmentWithAllSANs
		}).Once()

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_names_transition/enable_for_all_sans.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("network_configuration.0.enable_for_all_sans", "true").
						CheckEqual("network_configuration.0.dns_names.#", "2").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "san.test.akamai.com").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.dns_names.#", "1").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/update_explicit_dns_names.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.dns_names.#", "2").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "san.test.akamai.com").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_names_transition/enable_for_all_sans_after_explicit_dns_names.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("network_configuration.0.enable_for_all_sans", "true").
						CheckEqual("network_configuration.0.dns_names.#", "2").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "san.test.akamai.com").
						Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("reject enable_for_all_sans mode transition without changing dns_names", func(t *testing.T) {
		t.Parallel()
		testDVEnrollmentDNSModeTransitionRejected(t,
			testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/update_explicit_dns_names.tf"),
		)
	})

	t.Run("reject clone_dns_names mode transition without changing dns_names", func(t *testing.T) {
		t.Parallel()
		testDVEnrollmentDNSModeTransitionRejected(t,
			testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/clone_dns_names_unchanged.tf"),
		)
	})

	t.Run("reject DNS name settings for non-SNI enrollment", func(t *testing.T) {
		t.Parallel()
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(edgegrid.NewTestClient(), NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings_conflict/non_sni_with_enable_for_all_sans.tf"),
				ExpectError: regexp.MustCompile("'enable_for_all_sans', 'clone_dns_names', and 'dns_names' cannot be provided when 'sni_only' is false"),
			}},
		})
	})

	t.Run("create enrollment, empty sans", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		var enrollmentGet cps.GetEnrollmentResponse
		require.NoError(t, copier.CopyWithOption(&enrollmentGet, enrollment, copier.Option{DeepCopy: true}))
		enrollmentGet.CSR.SANS = []string{enrollment.CSR.CN}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(2)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/empty_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "1"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.test.akamai.com"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with empty sans and waiting for deletion", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		var enrollmentGet cps.GetEnrollmentResponse
		require.NoError(t, copier.CopyWithOption(&enrollmentGet, enrollment, copier.Option{DeepCopy: true}))
		enrollmentGet.CSR.SANS = []string{enrollment.CSR.CN}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(2)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the first get enrollment call still returns the enrollment.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// Mock that the second get enrollment is not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/empty_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "1"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.test.akamai.com"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment, MTLS", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		var enrollmentUpdate cps.GetEnrollmentResponse
		require.NoError(t, copier.CopyWithOption(&enrollmentUpdate, enrollment, copier.Option{DeepCopy: true, IgnoreEmpty: true}))
		enrollmentUpdate.NetworkConfiguration.ClientMutualAuthentication = &cps.ClientMutualAuthentication{
			AuthenticationOptions: &cps.AuthenticationOptions{
				OCSP: &cps.OCSP{
					Enabled: ptr.To(true),
				},
				SendCAListToClient: ptr.To(false),
			},
			SetID: "12345",
		}
		enrollmentUpdateReqBody := createEnrollmentReqBodyFromEnrollment(enrollmentUpdate)

		client.CPS.On("UpdateEnrollment",
			testutils.MockContext,
			cps.UpdateEnrollmentRequest{
				EnrollmentID:              1,
				EnrollmentRequestBody:     enrollmentUpdateReqBody,
				AllowCancelPendingChanges: ptr.To(true),
			},
		).Return(&cps.UpdateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/3"},
		}, nil).Once()

		enrollmentUpdate.Location = "/cps/v2/enrollments/1"
		enrollmentUpdate.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/3",
				ChangeType: "new-certificate",
			},
		}
		var enrollmentGet cps.GetEnrollmentResponse
		require.NoError(t, copier.CopyWithOption(&enrollmentGet, enrollmentUpdate, copier.Option{DeepCopy: true}))
		enrollmentGet.CSR.SANS = []string{enrollmentUpdate.CSR.CN}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(2)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/client_mutual_auth/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "network_configuration.0.client_mutual_authentication.0.set_id", "12345"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "network_configuration.0.client_mutual_authentication.0.ocsp_enabled", "true"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "network_configuration.0.client_mutual_authentication.0.send_ca_list_to_client", "false"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("lifecycle test with common name not empty, present in sans", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		commonName := "test.akamai.com"
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:    "US",
				CN:   commonName,
				L:    "Cambridge",
				O:    "Akamai",
				OU:   "WebEx",
				SANS: []string{commonName, "san.test.akamai.com"},
				ST:   "MA",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(4)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "2"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "2"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.san.test.akamai.com,_acme-challenge.test.akamai.com"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "2"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "2"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.san.test.akamai.com,_acme-challenge.test.akamai.com"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("set challenges arrays to empty if no allowedInput found", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:                   "US",
				CN:                  "test.akamai.com",
				L:                   "Cambridge",
				O:                   "Akamai",
				OU:                  "WebEx",
				PreferredTrustChain: "intermediate-a",
				SANS:                []string{"san.test.akamai.com"},
				ST:                  "MA",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(2)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/lifecycle/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "0"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "0"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("update with acknowledge warnings change, no enrollment update", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(3)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(3)

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Twice()
		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san2.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san2.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san2.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()
		client.CPS.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Twice()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "1"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns 404", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		// GetChangePreVerificationWarnings returns 404: the state suggests warnings exist but the API disagrees.
		// The provider should continue polling rather than fail.
		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(nil, cps.ErrNotFound).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Twice()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("certificate_type", "san").
						CheckEqual("validation_type", "dv").
						CheckEqual("registration_authority", "lets-encrypt").
						CheckEqual("dns_challenges.#", "1").
						CheckEqualBatch("dns_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).
						CheckEqual("http_challenges.#", "1").
						CheckEqualBatch("http_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns 404, later actual warns to acknowledge", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		// GetChangePreVerificationWarnings returns 404: the state suggests warnings exist but the API disagrees.
		// The provider should continue polling rather than fail.
		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(nil, cps.ErrNotFound).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Once()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		client.CPS.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Twice()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("certificate_type", "san").
						CheckEqual("validation_type", "dv").
						CheckEqual("registration_authority", "lets-encrypt").
						CheckEqual("dns_challenges.#", "1").
						CheckEqualBatch("dns_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).
						CheckEqual("http_challenges.#", "1").
						CheckEqualBatch("http_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns empty warnings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: ""}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Twice()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("certificate_type", "san").
						CheckEqual("validation_type", "dv").
						CheckEqual("registration_authority", "lets-encrypt").
						CheckEqual("dns_challenges.#", "1").
						CheckEqualBatch("dns_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).
						CheckEqual("http_challenges.#", "1").
						CheckEqualBatch("http_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("acknowledge warnings, pre-verification returns empty warnings, later actual warns to acknowledge", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: ""}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Once()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		client.CPS.On("AcknowledgePreVerificationWarnings", testutils.MockContext, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Twice()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Twice()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("allow_duplicate_common_name", "false").
						CheckEqual("acknowledge_pre_verification_warnings", "true").
						CheckEqual("common_name", "test.akamai.com").
						CheckEqual("sans.#", "0").
						CheckEqual("secure_network", "enhanced-tls").
						CheckEqual("sni_only", "true").
						CheckEqual("admin_contact.#", "1").
						CheckEqualBatch("admin_contact.0.", test.AttributeBatch{
							"first_name":       "R1",
							"last_name":        "D1",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r1d1@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("tech_contact.#", "1").
						CheckEqualBatch("tech_contact.0.", test.AttributeBatch{
							"first_name":       "R2",
							"last_name":        "D2",
							"title":            "",
							"organization":     "Akamai",
							"email":            "r2d2@akamai.com",
							"phone":            "123123123",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("csr.#", "1").
						CheckEqualBatch("csr.0.", test.AttributeBatch{
							"country_code":          "US",
							"city":                  "Cambridge",
							"organization":          "Akamai",
							"organizational_unit":   "WebEx",
							"preferred_trust_chain": "",
							"state":                 "MA",
						}).
						CheckEqual("network_configuration.#", "1").
						CheckEqualBatch("network_configuration.0.", test.AttributeBatch{
							"disallowed_tls_versions.#": "0",
							"clone_dns_names":           "true",
							"enable_for_all_sans":       "true",
							"geography":                 "core",
							"must_have_ciphers":         "ak-akamai-2020q1",
							"ocsp_stapling":             "on",
							"preferred_ciphers":         "ak-akamai-2020q1",
							"quic_enabled":              "false",
						}).
						CheckEqual("organization.#", "1").
						CheckEqualBatch("organization.0.", test.AttributeBatch{
							"name":             "Akamai",
							"phone":            "321321321",
							"address_line_one": "150 Broadway",
							"address_line_two": "",
							"city":             "Cambridge",
							"region":           "MA",
							"postal_code":      "12345",
							"country_code":     "US",
						}).
						CheckEqual("certificate_chain_type", "default").
						CheckEqual("signature_algorithm", "SHA-256").
						CheckEqual("certificate_type", "san").
						CheckEqual("validation_type", "dv").
						CheckEqual("registration_authority", "lets-encrypt").
						CheckEqual("dns_challenges.#", "1").
						CheckEqualBatch("dns_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).
						CheckEqual("http_challenges.#", "1").
						CheckEqualBatch("http_challenges.0.", test.AttributeBatch{
							"domain":        "test.akamai.com",
							"full_path":     "_acme-challenge.test.akamai.com",
							"response_body": "abc123",
						}).Build(),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment, allow duplicate common name", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
				AllowDuplicateCN:      true,
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		var enrollmentGet cps.GetEnrollmentResponse
		require.NoError(t, copier.CopyWithOption(&enrollmentGet, enrollment, copier.Option{DeepCopy: true}))
		enrollmentGet.CSR.SANS = []string{enrollment.CSR.CN}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		// second verification loop, valid status, empty allowed input array
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(2)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "http", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "dns", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(2)

		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/allow_duplicate_cn/create_enrollment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.0.domain", "test.akamai.com"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.0.full_path", "_acme-challenge.test.akamai.com"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "dns_challenges.0.response_body", "dns"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.#", "1"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.0.domain", "test.akamai.com"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.0.full_path", "_acme-challenge.test.akamai.com"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "http_challenges.0.response_body", "http"),
						resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "allow_duplicate_common_name", "true"),
						resource.TestCheckOutput("domains_to_validate", "_acme-challenge.test.akamai.com"),
					),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("verification failed with warnings, no acknowledgement", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitReviewPreVerificationSafetyChecks,
			},
		}, nil).Twice()

		client.CPS.On("GetChangePreVerificationWarnings", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
					ExpectError: regexp.MustCompile(`enrollment pre-verification returned warnings and the enrollment cannot be validated. Please fix the issues or set acknowledge_pre_verification_warnings flag to true then run 'terraform apply' again: some warning`),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment returns an error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(nil, fmt.Errorf("error creating enrollment")).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
					ExpectError: regexp.MustCompile(`error creating enrollment`),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with explicit dns_names", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:                   "US",
				CN:                  "test.akamai.com",
				L:                   "Cambridge",
				O:                   "Akamai",
				OU:                  "WebEx",
				SANS:                []string{"san.test.akamai.com"},
				ST:                  "MA",
				PreferredTrustChain: "intermediate-a",
			},
			EnableMultiStackedCertificates: false,
			NetworkConfiguration: &cps.NetworkConfiguration{
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
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
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBodyDNS := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBodyDNS,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		allowCancelDNS := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancelDNS,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/create_enrollment.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("contract_id", "ctr_1").
						CheckEqual("certificate_type", "san").
						CheckEqual("validation_type", "dv").
						CheckEqual("registration_authority", "lets-encrypt").
						CheckEqual("network_configuration.#", "1").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.dns_names.#", "1").
						CheckEqual("dns_challenges.#", "1").
						CheckEqual("http_challenges.#", "1").
						CheckTypeSetElemAttr("network_configuration.0.dns_names.*", "test.akamai.com").
						Build(),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with an explicitly empty dns_names attribute", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.CSR.SANS = []string{"san.test.akamai.com"}
		enrollment.CSR.PreferredTrustChain = "intermediate-a"
		enrollment.NetworkConfiguration.DisallowedTLSVersions = []string{"TLSv1", "TLSv1_1"}
		enrollment.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{CloneDNSNames: false}

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: createEnrollmentReqBodyFromEnrollment(enrollment),
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{{
			Location:   "/cps/v2/enrollments/1/changes/2",
			ChangeType: "new-certificate",
		}}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(10)
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()
		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Maybe()
		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{}, nil).Maybe()
		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/empty_dns_names.tf"),
				Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
					CheckEqual("network_configuration.0.enable_for_all_sans", "false").
					CheckEqual("network_configuration.0.dns_names.#", "0").
					Build(),
			}},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans and clone_dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings_conflict/create_enrollment.tf"),
					ExpectError: regexp.MustCompile("'enable_for_all_sans' and 'clone_dns_names' cannot both be provided at the same time"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans and dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings_conflict/create_enrollment_with_dns_names.tf"),
					ExpectError: regexp.MustCompile("'dns_names' cannot be provided when 'enable_for_all_sans' is true"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans unknown at plan and dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					Source:            "registry.terraform.io/hashicorp/random",
					VersionConstraint: "3.1.0",
				},
			},
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings_conflict/create_enrollment_with_unknown_enable_for_all_sans.tf"),
					ExpectError: regexp.MustCompile("'dns_names' cannot be provided when 'enable_for_all_sans' is true"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans unknown at plan and false without dns_names", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			ExternalProviders: map[string]resource.ExternalProvider{
				"random": {
					Source:            "registry.terraform.io/hashicorp/random",
					VersionConstraint: "3.1.0",
				},
			},
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings_conflict/create_enrollment_with_unknown_false_enable_for_all_sans.tf"),
					ExpectError: regexp.MustCompile("'dns_names' is required when 'enable_for_all_sans' or 'clone_dns_names' is false"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with clone_dns_names and dns_names conflict", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings_conflict/create_enrollment_with_clone_dns_names.tf"),
					ExpectError: regexp.MustCompile("'dns_names' cannot be provided when 'clone_dns_names' is true"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("create enrollment with enable_for_all_sans false without dns_names", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings_conflict/create_enrollment_with_false_without_dns_names.tf"),
					ExpectError: regexp.MustCompile("'dns_names' is required when 'enable_for_all_sans' or 'clone_dns_names' is false"),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("migration: create with clone_dns_names false - no conflict with enable_for_all_sans", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		enrollment := getTestDVEnrollment()
		enrollment.CSR.PreferredTrustChain = "intermediate-a"
		enrollment.CSR.SANS = []string{"san.test.akamai.com"}
		enrollment.NetworkConfiguration.DisallowedTLSVersions = []string{"TLSv1", "TLSv1_1"}
		enrollment.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: false,
			DNSNames:      []string{"test.akamai.com"},
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/clone_dns_names_migration/with_clone_dns_names_false.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("network_configuration.0.enable_for_all_sans", "false").
						CheckEqual("network_configuration.0.clone_dns_names", "false").
						Build(),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})

	t.Run("default: create without dns name settings defaults enable_for_all_sans to true", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		enrollment := getTestDVEnrollment()
		enrollment.CSR.PreferredTrustChain = "intermediate-a"
		enrollment.CSR.SANS = []string{"san.test.akamai.com"}
		enrollment.NetworkConfiguration.DisallowedTLSVersions = []string{"TLSv1", "TLSv1_1"}
		enrollment.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: true,
			DNSNames:      []string{"test.akamai.com", "san.test.akamai.com"},
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "pre-verification-safety-checks",
			},
		}, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(2)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/clone_dns_names_migration/without_dns_fields.tf"),
					Check: test.NewStateChecker("akamai_cps_dv_enrollment.dv").
						CheckEqual("network_configuration.0.enable_for_all_sans", "true").
						Build(),
				},
			},
		})

		client.CPS.AssertExpectations(t)
	})
	t.Run("non-SNI enrollment without DNS name settings has no drift", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := getTestDVEnrollment()
		enrollment.NetworkConfiguration.DNSNameSettings = nil
		enrollment.NetworkConfiguration.SNIOnly = false
		enrollment.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.DisallowedTLSVersions = []string{"TLSv1", "TLSv1_1"}
		enrollment.CSR.PreferredTrustChain = "intermediate-a"
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)
		assert.Nil(t, enrollmentReqBody.NetworkConfiguration.DNSNameSettings)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{ID: 1, Enrollment: "/cps/v2/enrollments/1"}, nil).Once()
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		allowCancel := true
		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		config := strings.Replace(
			testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/clone_dns_names_migration/without_dns_fields.tf"),
			"  sans = [\n    \"san.test.akamai.com\",\n  ]\n",
			"",
			1,
		)
		config = strings.Replace(config, "sni_only       = true", "sni_only       = false", 1)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{Config: config},
				{Config: config, PlanOnly: true},
			},
		})
		client.CPS.AssertExpectations(t)
	})
}

func getTestDVEnrollment() cps.GetEnrollmentResponse {
	return cps.GetEnrollmentResponse{
		AdminContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r1d1@akamai.com",
			FirstName:        "R1",
			LastName:         "D1",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		CertificateChainType: "default",
		CertificateType:      "san",
		CSR: &cps.CSR{
			C:  "US",
			CN: "test.akamai.com",
			L:  "Cambridge",
			O:  "Akamai",
			OU: "WebEx",
			ST: "MA",
		},
		NetworkConfiguration: &cps.NetworkConfiguration{
			DNSNameSettings: &cps.DNSNameSettings{
				CloneDNSNames: false,
				DNSNames:      []string{"test.akamai.com"},
			},
			Geography:        "core",
			MustHaveCiphers:  "ak-akamai-2020q1",
			OCSPStapling:     "on",
			PreferredCiphers: "ak-akamai-2020q1",
			QuicEnabled:      false,
			SecureNetwork:    "enhanced-tls",
			SNIOnly:          true,
		},
		Org: &cps.Org{
			AddressLineOne: "150 Broadway",
			City:           "Cambridge",
			Country:        "US",
			Name:           "Akamai",
			Phone:          "321321321",
			PostalCode:     "12345",
			Region:         "MA",
		},
		RA:                 "lets-encrypt",
		SignatureAlgorithm: "SHA-256",
		TechContact: &cps.Contact{
			AddressLineOne:   "150 Broadway",
			City:             "Cambridge",
			Country:          "US",
			Email:            "r2d2@akamai.com",
			FirstName:        "R2",
			LastName:         "D2",
			OrganizationName: "Akamai",
			Phone:            "123123123",
			PostalCode:       "12345",
			Region:           "MA",
		},
		ValidationType: "dv",
	}
}

func testDVEnrollmentDNSModeTransitionRejected(t *testing.T, config string) {
	t.Helper()

	client := edgegrid.NewTestClient()
	allowCancel := true
	enrollment := getTestDVEnrollment()
	enrollment.CSR.SANS = []string{"san.test.akamai.com"}
	enrollment.CSR.PreferredTrustChain = "intermediate-a"
	enrollment.NetworkConfiguration.DisallowedTLSVersions = []string{"TLSv1", "TLSv1_1"}
	enrollment.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
	enrollment.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
	enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{CloneDNSNames: true}

	client.CPS.On("CreateEnrollment",
		testutils.MockContext,
		cps.CreateEnrollmentRequest{
			EnrollmentRequestBody: createEnrollmentReqBodyFromEnrollment(enrollment),
			ContractID:            "1",
		},
	).Return(&cps.CreateEnrollmentResponse{
		ID:         1,
		Enrollment: "/cps/v2/enrollments/1",
		Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
	}, nil).Once()

	enrollment.Location = "/cps/v2/enrollments/1"
	enrollment.PendingChanges = []cps.PendingChange{{
		Location:   "/cps/v2/enrollments/1/changes/2",
		ChangeType: "new-certificate",
	}}
	enrollment.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"test.akamai.com", "san.test.akamai.com"}
	client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
		Return(&enrollment, nil).Times(10)
	client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
		Return(nil, cps.ErrNotFound).Once()
	client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
		EnrollmentID: 1,
		ChangeID:     2,
	}).Return(&cps.Change{
		AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
		StatusInfo: &cps.StatusInfo{
			State:  "awaiting-input",
			Status: coodinateDomainValidation,
		},
	}, nil).Maybe()
	client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
		EnrollmentID: 1,
		ChangeID:     2,
	}).Return(&cps.DVArray{}, nil).Maybe()
	client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
		EnrollmentID:              1,
		AllowCancelPendingChanges: &allowCancel,
	}).Return(&cps.RemoveEnrollmentResponse{Enrollment: "1"}, nil).Maybe()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
		Steps: []resource.TestStep{
			{
				Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_names_transition/enable_for_all_sans.tf"),
			},
			{
				Config:      config,
				ExpectError: regexp.MustCompile("'dns_names' must change when switching 'enable_for_all_sans' or 'clone_dns_names' from true to false"),
			},
		},
	})
	client.CPS.AssertExpectations(t)
}

func mockDVTransitionChangeStatus(client *edgegrid.TestClient) {
	client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
		EnrollmentID: 1,
		ChangeID:     2,
	}).Return(&cps.Change{
		AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
		StatusInfo: &cps.StatusInfo{
			State:  "awaiting-input",
			Status: coodinateDomainValidation,
		},
	}, nil).Times(15)
	client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
		EnrollmentID: 1,
		ChangeID:     2,
	}).Return(&cps.DVArray{}, nil).Times(11)
}

func TestDNSNameSettingsDefaultingEnabled(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		sniOnly  cty.Value
		expected bool
	}{
		"SNI enabled":     {sniOnly: cty.True, expected: true},
		"SNI disabled":    {sniOnly: cty.False, expected: false},
		"SNI unspecified": {sniOnly: cty.NullVal(cty.Bool), expected: false},
		"SNI unknown":     {sniOnly: cty.UnknownVal(cty.Bool), expected: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, dnsNameSettingsDefaultingEnabled(test.sniOnly))
		})
	}
}

func TestValidateDNSNameSettings(t *testing.T) {
	t.Parallel()

	nullBool := cty.NullVal(cty.Bool)
	nullDNSNames := cty.NullVal(cty.Set(cty.String))
	dnsNameSettingsConfig := func(sniOnly, cloneDNS, enableForAllSANs, dnsNames cty.Value) cty.Value {
		return cty.ObjectVal(map[string]cty.Value{
			"sni_only": sniOnly,
			"network_configuration": cty.ListVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"clone_dns_names":     cloneDNS,
					"enable_for_all_sans": enableForAllSANs,
					"dns_names":           dnsNames,
				}),
			}),
		})
	}

	tests := map[string]struct {
		config      cty.Value
		expectError string
	}{
		"non-SNI permits omitted DNS name settings": {
			config: dnsNameSettingsConfig(cty.False, nullBool, nullBool, nullDNSNames),
		},
		"non-SNI rejects clone_dns_names": {
			config:      dnsNameSettingsConfig(cty.False, cty.False, nullBool, nullDNSNames),
			expectError: "'enable_for_all_sans', 'clone_dns_names', and 'dns_names' cannot be provided when 'sni_only' is false",
		},
		"non-SNI rejects enable_for_all_sans": {
			config:      dnsNameSettingsConfig(cty.False, nullBool, cty.True, nullDNSNames),
			expectError: "'enable_for_all_sans', 'clone_dns_names', and 'dns_names' cannot be provided when 'sni_only' is false",
		},
		"non-SNI rejects dns_names": {
			config:      dnsNameSettingsConfig(cty.False, nullBool, nullBool, cty.SetVal([]cty.Value{cty.StringVal("example.com")})),
			expectError: "'enable_for_all_sans', 'clone_dns_names', and 'dns_names' cannot be provided when 'sni_only' is false",
		},
		"SNI with clone_dns_names false requires dns_names": {
			config:      dnsNameSettingsConfig(cty.True, cty.False, nullBool, nullDNSNames),
			expectError: "'dns_names' is required when 'enable_for_all_sans' or 'clone_dns_names' is false",
		},
		"SNI with enable_for_all_sans false requires dns_names": {
			config:      dnsNameSettingsConfig(cty.True, nullBool, cty.False, nullDNSNames),
			expectError: "'dns_names' is required when 'enable_for_all_sans' or 'clone_dns_names' is false",
		},
		"SNI permits clone_dns_names to be false with an empty dns_names attribute": {
			config: dnsNameSettingsConfig(cty.True, cty.False, nullBool, cty.SetValEmpty(cty.String)),
		},
		"SNI permits enable_for_all_sans to be false with an empty dns_names attribute": {
			config: dnsNameSettingsConfig(cty.True, nullBool, cty.False, cty.SetValEmpty(cty.String)),
		},
		"SNI rejects clone_dns_names and enable_for_all_sans together": {
			config:      dnsNameSettingsConfig(cty.True, cty.False, cty.True, nullDNSNames),
			expectError: "'enable_for_all_sans' and 'clone_dns_names' cannot both be provided at the same time",
		},
		"SNI permits omitted DNS name settings": {
			config: dnsNameSettingsConfig(cty.True, nullBool, nullBool, nullDNSNames),
		},
		"SNI permits clone_dns_names true": {
			config: dnsNameSettingsConfig(cty.True, cty.True, nullBool, nullDNSNames),
		},
		"SNI permits enable_for_all_sans true": {
			config: dnsNameSettingsConfig(cty.True, nullBool, cty.True, nullDNSNames),
		},
		"SNI permits clone_dns_names false with dns_names": {
			config: dnsNameSettingsConfig(cty.True, cty.False, nullBool, cty.SetVal([]cty.Value{cty.StringVal("example.com")})),
		},
		"SNI permits enable_for_all_sans false with dns_names": {
			config: dnsNameSettingsConfig(cty.True, nullBool, cty.False, cty.SetVal([]cty.Value{cty.StringVal("example.com")})),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := validateDNSNameSettings(test.config)
			if test.expectError == "" {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Equal(t, test.expectError, err.Error())
		})
	}
}

func TestResourceDVEnrollmentImport(t *testing.T) {
	t.Parallel()
	t.Run("import", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:                   "US",
				CN:                  "test.akamai.com",
				L:                   "Cambridge",
				O:                   "Akamai",
				OU:                  "WebEx",
				PreferredTrustChain: "intermediate-a",
				ST:                  "MA",
			},
			NetworkConfiguration: &cps.NetworkConfiguration{
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: true,
					DNSNames:      []string{"test.akamai.com"}},
				Geography:        "core",
				MustHaveCiphers:  "ak-akamai-2020q1",
				OCSPStapling:     "on",
				PreferredCiphers: "ak-akamai-2020q1",
				QuicEnabled:      false,
				SecureNetwork:    "enhanced-tls",
				SNIOnly:          true,
			},
			Org: &cps.Org{
				AddressLineOne: "150 Broadway",
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)
		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(3)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(3)
		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/import/import_enrollment.tf"),
				},
				{
					Config:            testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/import/import_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     "1,ctr_1",
					ResourceName:      "akamai_cps_dv_enrollment.dv",
					ImportStateVerify: true,
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						require.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						assert.Equal(t, "true", rs.Attributes["network_configuration.0.enable_for_all_sans"])
						assert.Equal(t, "true", rs.Attributes["network_configuration.0.clone_dns_names"])
						assert.Equal(t, "1", rs.Attributes["network_configuration.0.dns_names.#"])
						assert.ElementsMatch(t, []string{"test.akamai.com"}, dnsNamesFromState(rs.Attributes))
						return nil
					},
					ImportStateVerifyIgnore: []string{"network_configuration.0.quic_enabled"},
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("import with dns name settings", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		id := "1,ctr_1"

		enrollment := cps.GetEnrollmentResponse{
			AdminContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r1d1@akamai.com",
				FirstName:        "R1",
				LastName:         "D1",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			CertificateChainType: "default",
			CertificateType:      "san",
			CSR: &cps.CSR{
				C:                   "US",
				CN:                  "test.akamai.com",
				L:                   "Cambridge",
				O:                   "Akamai",
				OU:                  "WebEx",
				SANS:                []string{"san.test.akamai.com"},
				PreferredTrustChain: "intermediate-a",
				ST:                  "MA",
			},
			NetworkConfiguration: &cps.NetworkConfiguration{
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
					DNSNames:      []string{"test.akamai.com"},
				},
				DisallowedTLSVersions: []string{"TLSv1", "TLSv1_1"},
				Geography:             "core",
				MustHaveCiphers:       "ak-akamai-default",
				OCSPStapling:          "on",
				PreferredCiphers:      "ak-akamai-default",
				QuicEnabled:           false,
				SecureNetwork:         "enhanced-tls",
				SNIOnly:               true,
			},
			Org: &cps.Org{
				AddressLineOne: "150 Broadway",
				City:           "Cambridge",
				Country:        "US",
				Name:           "Akamai",
				Phone:          "321321321",
				PostalCode:     "12345",
				Region:         "MA",
			},
			RA:                 "lets-encrypt",
			SignatureAlgorithm: "SHA-256",
			TechContact: &cps.Contact{
				AddressLineOne:   "150 Broadway",
				City:             "Cambridge",
				Country:          "US",
				Email:            "r2d2@akamai.com",
				FirstName:        "R2",
				LastName:         "D2",
				OrganizationName: "Akamai",
				Phone:            "123123123",
				PostalCode:       "12345",
				Region:           "MA",
			},
			ValidationType: "dv",
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(3)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(3)
		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		// Mock that the enrollment in not found after removal.
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/create_enrollment.tf"),
				},
				{
					Config:            testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/create_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_cps_dv_enrollment.dv",
					ImportStateVerify: true,
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						require.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						assert.Equal(t, "false", rs.Attributes["network_configuration.0.enable_for_all_sans"])
						assert.Equal(t, "1", rs.Attributes["network_configuration.0.dns_names.#"])
						assert.ElementsMatch(t, []string{"test.akamai.com"}, dnsNamesFromState(rs.Attributes))
						return nil
					},
					// It looks that there bug in SDK that values for bool optional fields are not persisted on create
					ImportStateVerifyIgnore: []string{"network_configuration.0.quic_enabled"},
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("import with clone_dns_names false", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		id := "1,ctr_1"

		enrollment := getTestDVEnrollment()
		enrollment.CSR.PreferredTrustChain = "intermediate-a"
		enrollment.CSR.SANS = []string{"san.test.akamai.com"}
		enrollment.NetworkConfiguration.DisallowedTLSVersions = []string{"TLSv1", "TLSv1_1"}
		enrollment.NetworkConfiguration.MustHaveCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.PreferredCiphers = "ak-akamai-default"
		enrollment.NetworkConfiguration.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: false,
			DNSNames:      []string{"test.akamai.com"},
		}
		enrollmentReqBody := createEnrollmentReqBodyFromEnrollment(enrollment)

		client.CPS.On("CreateEnrollment",
			testutils.MockContext,
			cps.CreateEnrollmentRequest{
				EnrollmentRequestBody: enrollmentReqBody,
				ContractID:            "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		enrollment.Location = "/cps/v2/enrollments/1"
		enrollment.PendingChanges = []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/1/changes/2",
				ChangeType: "new-certificate",
			},
		}
		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		client.CPS.On("GetChangeStatus", testutils.MockContext, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: coodinateDomainValidation,
			},
		}, nil).Times(3)

		client.CPS.On("GetChangeLetsEncryptChallenges", testutils.MockContext, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenge{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
		}}, nil).Times(3)
		allowCancel := true

		client.CPS.On("RemoveEnrollment", testutils.MockContext, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(nil, cps.ErrNotFound).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/create_enrollment.tf"),
				},
				{
					Config:            testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/dns_name_settings/create_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_cps_dv_enrollment.dv",
					ImportStateVerify: true,
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						require.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						assert.Equal(t, "false", rs.Attributes["network_configuration.0.enable_for_all_sans"])
						assert.Equal(t, "false", rs.Attributes["network_configuration.0.clone_dns_names"])
						assert.ElementsMatch(t, []string{"test.akamai.com"}, dnsNamesFromState(rs.Attributes))
						return nil
					},
					ImportStateVerifyIgnore: []string{"network_configuration.0.quic_enabled"},
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})

	t.Run("import error when validation type is not dv", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		id := "1,ctr_1"

		enrollment := cps.GetEnrollmentResponse{
			ValidationType: "third-party",
		}

		client.CPS.On("GetEnrollment", testutils.MockContext, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(1)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewCustomPollingSubprovider(testPollChangeStatusInterval, testPollGetEnrollmentInterval)),
			Steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/import/import_enrollment.tf"),
					ImportState:   true,
					ImportStateId: id,
					ResourceName:  "akamai_cps_dv_enrollment.dv",
					ExpectError:   regexp.MustCompile("unable to import: wrong validation type: expected 'dv', got 'third-party'"),
				},
			},
		})
		client.CPS.AssertExpectations(t)
	})
}

func createEnrollmentReqBodyFromEnrollment(en cps.GetEnrollmentResponse) cps.EnrollmentRequestBody {
	networkConfiguration := en.NetworkConfiguration
	if networkConfiguration != nil && networkConfiguration.DNSNameSettings != nil && networkConfiguration.DNSNameSettings.CloneDNSNames {
		networkConfigurationCopy := *networkConfiguration
		dnsNameSettingsCopy := *networkConfiguration.DNSNameSettings
		dnsNameSettingsCopy.DNSNames = nil
		networkConfigurationCopy.DNSNameSettings = &dnsNameSettingsCopy
		networkConfiguration = &networkConfigurationCopy
	}
	return cps.EnrollmentRequestBody{
		AdminContact:                   en.AdminContact,
		AutoRenewalStartTime:           en.AutoRenewalStartTime,
		CertificateChainType:           en.CertificateChainType,
		CertificateType:                en.CertificateType,
		ChangeManagement:               en.ChangeManagement,
		CSR:                            en.CSR,
		EnableMultiStackedCertificates: en.EnableMultiStackedCertificates,
		NetworkConfiguration:           networkConfiguration,
		Org:                            en.Org,
		OrgID:                          en.OrgID,
		RA:                             en.RA,
		SignatureAlgorithm:             en.SignatureAlgorithm,
		TechContact:                    en.TechContact,
		ThirdParty:                     en.ThirdParty,
		ValidationType:                 en.ValidationType,
	}
}
