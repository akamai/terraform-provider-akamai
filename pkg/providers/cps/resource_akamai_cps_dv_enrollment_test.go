package cps

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceDVEnrollment(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &mockcps{}
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
			ChangeManagement:     false,
			CSR: &cps.CSR{
				C:    "US",
				CN:   "test.akamai.com",
				L:    "Cambridge",
				O:    "Akamai",
				OU:   "WebEx",
				SANS: []string{"san.test.akamai.com"},
				ST:   "MA",
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
		enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
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
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
			},
		}, nil).Times(3)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenges{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenges{
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
		enrollmentUpdate.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdate, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			},
		}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			},
		}, nil).Twice()
		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenges{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenges{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenges{
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVEnrollment/lifecycle/create_enrollment.tf"),
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
						Config: loadFixtureString("testdata/TestResDVEnrollment/lifecycle/update_enrollment.tf"),
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
		client := &mockcps{}
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
			ChangeManagement:     false,
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
		enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
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
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
			},
		}, nil).Times(2)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenges{
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVEnrollment/empty_sans/create_enrollment.tf"),
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

	t.Run("lifecycle test with common name not empty, present in sans", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &mockcps{}
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
			ChangeManagement:     false,
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
		enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
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
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
			},
		}, nil).Times(4)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenges{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenges{
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
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
						Config: loadFixtureString("testdata/TestResDVEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
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
		client := &mockcps{}
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
			ChangeManagement:     false,
			CSR: &cps.CSR{
				C:    "US",
				CN:   "test.akamai.com",
				L:    "Cambridge",
				O:    "Akamai",
				OU:   "WebEx",
				SANS: []string{"san.test.akamai.com"},
				ST:   "MA",
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
		enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVEnrollment/lifecycle/create_enrollment.tf"),
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
		client := &mockcps{}
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
		enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
			},
		}, nil).Times(3)

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenges{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenges{
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
				Status: "coodinate-domain-validation",
			},
		}, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
			},
		}, nil).Twice()
		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenges{
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
					{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenges{
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
					{FullPath: "_acme-challenge.san.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
				},
				Domain:           "san.test.akamai.com",
				ValidationStatus: "IN_PROGRESS",
			},
			{
				Challenges: []cps.Challenges{
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "contract_id", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "certificate_type", "san"),
							resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "validation_type", "dv"),
							resource.TestCheckResourceAttr("akamai_cps_dv_enrollment.dv", "registration_authority", "lets-encrypt"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
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
		client := &mockcps{}
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
		enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-ack"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "wait-review-pre-verification-safety-checks",
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
			AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "coodinate-domain-validation",
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
				Status: "coodinate-domain-validation",
			},
		}, nil).Twice()

		client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.DVArray{DV: []cps.DV{
			{
				Challenges: []cps.Challenges{
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDVEnrollment/acknowledge_warnings/create_enrollment.tf"),
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

	t.Run("verification failed with warnings, no acknowledgement", func(t *testing.T) {
		client := &mockcps{}
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
		enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "pre-verification-warnings-ack"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: "wait-review-pre-verification-safety-checks",
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`enrollment pre-verification returned warnings and the enrollment cannot be validated. Please fix the issues or set acknowledge_pre_validation_warnings flag to true then run 'terraform apply' again: some warning`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create enrollment returns an error", func(t *testing.T) {
		client := &mockcps{}
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResDVEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`error creating enrollment`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestResourceDVEnrollmentImport(t *testing.T) {
	client := &mockcps{}
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
	enrollment.PendingChanges = []string{"/cps/v2/enrollments/1/changes/2"}
	client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
		Return(&enrollment, nil).Once()

	client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
		EnrollmentID: 1,
		ChangeID:     2,
	}).Return(&cps.Change{
		AllowedInput: []cps.AllowedInput{{Type: "lets-encrypt-challenges"}},
		StatusInfo: &cps.StatusInfo{
			State:  "awaiting-input",
			Status: "coodinate-domain-validation",
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
			Status: "coodinate-domain-validation",
		},
	}, nil).Times(3)

	client.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
		EnrollmentID: 1,
		ChangeID:     2,
	}).Return(&cps.DVArray{DV: []cps.DV{
		{
			Challenges: []cps.Challenges{
				{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "http-01", Status: "pending"},
				{FullPath: "_acme-challenge.test.akamai.com", ResponseBody: "abc123", Type: "dns-01", Status: "pending"},
			},
			Domain:           "test.akamai.com",
			ValidationStatus: "IN_PROGRESS",
		},
		{
			Challenges: []cps.Challenges{
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
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestResDVEnrollment/import/import_enrollment.tf"),
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						assert.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
						assert.Equal(t, "1", rs.Attributes["id"])
						return nil
					},
				},
				{
					Config:            loadFixtureString("testdata/TestResDVEnrollment/import/import_enrollment.tf"),
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_cps_dv_enrollment.dv",
					ImportStateVerify: true,
				},
			},
		})
	})
}
