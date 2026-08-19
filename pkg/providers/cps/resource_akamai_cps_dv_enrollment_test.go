package cps

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
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
					DNSNames:      []string{"san.test.akamai.com"},
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
		enrollmentUpdate.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
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
					DNSNames:      []string{"san.test.akamai.com"},
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
		enrollmentUpdate.NetworkConfiguration.DNSNameSettings.DNSNames = nil
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
					DNSNames:      []string{commonName, "san.test.akamai.com"},
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
					DNSNames:      []string{"san.test.akamai.com"},
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
							"clone_dns_names":           "false",
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
							"clone_dns_names":           "false",
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
							"clone_dns_names":           "false",
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
							"clone_dns_names":           "false",
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

func TestResourceDVEnrollmentImport(t *testing.T) {
	t.Parallel()
	t.Run("import", func(t *testing.T) {
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
				PreferredTrustChain: "intermediate-a",
				ST:                  "MA",
			},
			NetworkConfiguration: &cps.NetworkConfiguration{
				DNSNameSettings: &cps.DNSNameSettings{
					CloneDNSNames: false,
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
					Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/import/import_enrollment.tf"),
				},
				{
					Config:            testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/import/import_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_cps_dv_enrollment.dv",
					ImportStateVerify: true,
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						assert.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						return nil
					},
					// It looks that there bug in SDK that values for bool optional fields are not persisted on create
					ImportStateVerifyIgnore: []string{"network_configuration.0.clone_dns_names", "network_configuration.0.quic_enabled"},
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
	return cps.EnrollmentRequestBody{
		AdminContact:                   en.AdminContact,
		AutoRenewalStartTime:           en.AutoRenewalStartTime,
		CertificateChainType:           en.CertificateChainType,
		CertificateType:                en.CertificateType,
		ChangeManagement:               en.ChangeManagement,
		CSR:                            en.CSR,
		EnableMultiStackedCertificates: en.EnableMultiStackedCertificates,
		NetworkConfiguration:           en.NetworkConfiguration,
		Org:                            en.Org,
		OrgID:                          en.OrgID,
		RA:                             en.RA,
		SignatureAlgorithm:             en.SignatureAlgorithm,
		TechContact:                    en.TechContact,
		ThirdParty:                     en.ThirdParty,
		ValidationType:                 en.ValidationType,
	}
}
