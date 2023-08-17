package cps

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/tools"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceDVEnrollment(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := cps.Enrollment{
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(3)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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

		var enrollmentUpdate cps.Enrollment
		err := copier.CopyWithOption(&enrollmentUpdate, enrollment, copier.Option{DeepCopy: true})
		require.NoError(t, err)
		enrollmentUpdate.AdminContact.FirstName = "R5"
		enrollmentUpdate.AdminContact.LastName = "D5"
		enrollmentUpdate.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
		enrollmentUpdate.CSR.PreferredTrustChain = "dst-root-ca-x3"
		enrollmentUpdate.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
		enrollmentUpdate.Location = ""
		enrollmentUpdate.PendingChanges = nil
		allowCancel := true
		client.On("UpdateEnrollment",
			mock.Anything,
			cps.UpdateEnrollmentRequest{
				Enrollment:                enrollmentUpdate,
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdate, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Twice()
		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
							resource.TestCheckOutput("domains_to_validate", "_acme-challenge.san.test.akamai.com,_acme-challenge.san2.test.akamai.com,_acme-challenge.test.akamai.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create enrollment, empty sans", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := cps.Enrollment{
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		var enrollmentGet cps.Enrollment
		require.NoError(t, copier.CopyWithOption(&enrollmentGet, enrollment, copier.Option{DeepCopy: true}))
		enrollmentGet.CSR.SANS = []string{enrollment.CSR.CN}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(2)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})
	t.Run("create enrollment, MTLS", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := cps.Enrollment{
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(&cps.CreateEnrollmentResponse{
			ID:         1,
			Enrollment: "/cps/v2/enrollments/1",
			Changes:    []string{"/cps/v2/enrollments/1/changes/2"},
		}, nil).Once()

		var enrollmentUpdate cps.Enrollment
		require.NoError(t, copier.CopyWithOption(&enrollmentUpdate, enrollment, copier.Option{DeepCopy: true, IgnoreEmpty: true}))
		enrollmentUpdate.NetworkConfiguration.ClientMutualAuthentication = &cps.ClientMutualAuthentication{
			AuthenticationOptions: &cps.AuthenticationOptions{
				OCSP: &cps.OCSP{
					Enabled: tools.BoolPtr(true),
				},
				SendCAListToClient: tools.BoolPtr(false),
			},
			SetID: "12345",
		}
		client.On("UpdateEnrollment",
			mock.Anything,
			cps.UpdateEnrollmentRequest{
				EnrollmentID:              1,
				Enrollment:                enrollmentUpdate,
				AllowCancelPendingChanges: tools.BoolPtr(true),
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
		var enrollmentGet cps.Enrollment
		require.NoError(t, copier.CopyWithOption(&enrollmentGet, enrollmentUpdate, copier.Option{DeepCopy: true}))
		enrollmentGet.CSR.SANS = []string{enrollmentUpdate.CSR.CN}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     3,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(2)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})

	t.Run("lifecycle test with common name not empty, present in sans", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		commonName := "test.akamai.com"
		enrollment := cps.Enrollment{
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(4)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})

	t.Run("set challenges arrays to empty if no allowedInput found", func(t *testing.T) {
		client := &cps.Mock{}
		enrollment := cps.Enrollment{
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(2)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(2)

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})

	t.Run("update with acknowledge warnings change, no enrollment update", func(t *testing.T) {
		client := &cps.Mock{}
		enrollment := cps.Enrollment{
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
				Geography:     "core",
				SecureNetwork: "enhanced-tls",
				SNIOnly:       true,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(3)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Twice()
		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})

	t.Run("acknowledge warnings", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := cps.Enrollment{
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
				Geography:     "core",
				SecureNetwork: "enhanced-tls",
				SNIOnly:       true,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-acknowledgement"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusVerificationWarnings,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()
		client.On("AcknowledgePreVerificationWarnings", mock.Anything, cps.AcknowledgementRequest{
			EnrollmentID:    1,
			ChangeID:        2,
			Acknowledgement: cps.Acknowledgement{Acknowledgement: "acknowledge"},
		}).Return(nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: inputTypePreVerificationWarningsAck,
			},
		}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Twice()

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})

	t.Run("create enrollment, allow duplicate common name", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := cps.Enrollment{
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment:       enrollment,
				ContractID:       "1",
				AllowDuplicateCN: true,
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
		var enrollmentGet cps.Enrollment
		require.NoError(t, copier.CopyWithOption(&enrollmentGet, enrollment, copier.Option{DeepCopy: true}))
		enrollmentGet.CSR.SANS = []string{enrollment.CSR.CN}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		// first verification loop, invalid status
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
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
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		// final verification loop, everything in place
		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(2)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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

		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})

	t.Run("verification failed with warnings, no acknowledgement", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := cps.Enrollment{
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
				Geography:     "core",
				SecureNetwork: "enhanced-tls",
				SNIOnly:       true,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: inputTypePreVerificationWarningsAck}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusVerificationWarnings,
			},
		}, nil).Twice()

		client.On("GetChangePreVerificationWarnings", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.PreVerificationWarnings{Warnings: "some warning"}, nil).Once()

		allowCancel := true
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`enrollment pre-verification returned warnings and the enrollment cannot be validated. Please fix the issues or set acknowledge_pre_validation_warnings flag to true then run 'terraform apply' again: some warning`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create enrollment returns an error", func(t *testing.T) {
		client := &cps.Mock{}
		PollForChangeStatusInterval = 1 * time.Millisecond
		enrollment := cps.Enrollment{
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
				Geography:     "core",
				SecureNetwork: "enhanced-tls",
				SNIOnly:       true,
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

		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
			},
		).Return(nil, fmt.Errorf("error creating enrollment")).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`error creating enrollment`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestResourceDVEnrollmentImport(t *testing.T) {
	t.Run("import", func(t *testing.T) {
		client := &cps.Mock{}
		id := "1,ctr_1"

		enrollment := cps.Enrollment{
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
				Geography:     "core",
				SecureNetwork: "enhanced-tls",
				SNIOnly:       true,
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
		client.On("CreateEnrollment",
			mock.Anything,
			cps.CreateEnrollmentRequest{
				Enrollment: enrollment,
				ContractID: "1",
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
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: statusCoordinateDomainValidation,
			},
		}, nil).Times(3)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
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
		client.On("RemoveEnrollment", mock.Anything, cps.RemoveEnrollmentRequest{
			EnrollmentID:              1,
			AllowCancelPendingChanges: &allowCancel,
		}).Return(&cps.RemoveEnrollmentResponse{
			Enrollment: "1",
		}, nil).Once()
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/import/import_enrollment.tf"),
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "1", rs.Attributes["id"])
							return nil
						},
					},
					{
						Config:            testutils.LoadFixtureString(t, "testdata/TestResDVEnrollment/import/import_enrollment.tf"),
						ImportState:       true,
						ImportStateId:     id,
						ResourceName:      "akamai_cps_dv_enrollment.dv",
						ImportStateVerify: true,
					},
				},
			})
		})
	})

	t.Run("import error when validation type is not dv", func(t *testing.T) {
		client := &cps.Mock{}
		id := "1,ctr_1"

		enrollment := cps.Enrollment{
			ValidationType: "third-party",
		}

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})
	})
}
