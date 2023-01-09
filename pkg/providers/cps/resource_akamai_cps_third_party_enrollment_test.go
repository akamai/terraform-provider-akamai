package cps

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResourceThirdPartyEnrollment(t *testing.T) {
	t.Run("lifecycle test", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := getSimpleEnrollment()

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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		var enrollmentUpdate cps.Enrollment
		err := copier.CopyWithOption(&enrollmentUpdate, enrollment, copier.Option{DeepCopy: true})
		require.NoError(t, err)
		enrollmentUpdate.AdminContact.FirstName = "R5"
		enrollmentUpdate.AdminContact.LastName = "D5"
		enrollmentUpdate.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
		enrollmentUpdate.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
		enrollmentUpdate.Location = ""
		enrollmentUpdate.PendingChanges = nil
		enrollmentUpdate.SignatureAlgorithm = "SHA-1"
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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/lifecycle/update_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
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
			CertificateType:      "third-party",
			ChangeManagement:     false,
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/empty_sans/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
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
			CertificateType:      "third-party",
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
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/lifecycle_cn_in_sans/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("set challenges arrays to empty if no allowedInput found", func(t *testing.T) {
		client := &cps.Mock{}
		enrollment := getSimpleEnrollment()

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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(2)

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/acknowledge_warnings/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
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
			CertificateType:      "third-party",
			ChangeManagement:     false,
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(2)

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/allow_duplicate_cn/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "allow_duplicate_common_name", "true"),
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
						Config:      loadFixtureString("testdata/TestResThirdPartyEnrollment/no_acknowledge_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`error creating enrollment`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("auto approve warnings - all warnings on the list to auto approve", func(t *testing.T) {
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
		}).Return(&cps.PreVerificationWarnings{Warnings: "Certificate data is blank or missing for 'RSA'.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Twice()

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("auto approve warnings - some warnings not on the list to auto approve", func(t *testing.T) {
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
		}).Return(&cps.PreVerificationWarnings{Warnings: "Certificate data is blank or missing for 'RSA'.\nError parsing expected trust chains.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

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
						Config:      loadFixtureString("testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`warnings cannot be approved: "FIXED_TRUST_CHAIN_PARSING_ERROR"`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("auto approve warnings - some warnings are unknown", func(t *testing.T) {
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
		}).Return(&cps.PreVerificationWarnings{Warnings: "Certificate data is blank or missing for 'RSA'.\nThis is unknown warning.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate."}, nil).Once()

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
						Config:      loadFixtureString("testdata/TestResThirdPartyEnrollment/auto_approve_warnings/create_enrollment.tf"),
						ExpectError: regexp.MustCompile(`received warning\(s\) does not match any known warning: 'This is unknown warning.'`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestResourceThirdPartyEnrollmentImport(t *testing.T) {
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
			CertificateType:      "third-party",
			CSR: &cps.CSR{
				C:  "US",
				CN: "test.akamai.com",
				L:  "Cambridge",
				O:  "Akamai",
				OU: "WebEx",
				ST: "MA",
			},
			EnableMultiStackedCertificates: true,
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
			RA:                 "third-party",
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
			ValidationType: "third-party",
			ThirdParty: &cps.ThirdParty{
				ExcludeSANS: false,
			},
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
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(4)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Times(3)

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
						ImportStateCheck: func(s []*terraform.InstanceState) error {
							assert.Len(t, s, 1)
							rs := s[0]
							assert.Equal(t, "ctr_1", rs.Attributes["contract_id"])
							assert.Equal(t, "1", rs.Attributes["id"])
							return nil
						},
					},
					{
						Config:            loadFixtureString("testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
						ImportState:       true,
						ImportStateId:     id,
						ResourceName:      "akamai_cps_third_party_enrollment.third_party",
						ImportStateVerify: true,
					},
				},
			})
		})
	})

	t.Run("import error when validation type is not third_party", func(t *testing.T) {
		client := &cps.Mock{}
		id := "1,ctr_1"

		enrollment := cps.Enrollment{
			ValidationType: "dv",
		}

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollment, nil).Times(1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:        loadFixtureString("testdata/TestResThirdPartyEnrollment/import/import_enrollment.tf"),
						ImportState:   true,
						ImportStateId: id,
						ResourceName:  "akamai_cps_third_party_enrollment.third_party",
						ExpectError:   regexp.MustCompile("unable to import: wrong validation type: expected 'third-party', got 'dv'"),
					},
				},
			})
		})
	})
}

func TestSuppressingSignatureAlgorithm(t *testing.T) {
	t.Run("suppress signature algorithm", func(t *testing.T) {
		PollForChangeStatusInterval = 1 * time.Millisecond
		client := &cps.Mock{}
		enrollment := getSimpleEnrollment()

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
		err := copier.CopyWithOption(&enrollmentGet, enrollment, copier.Option{DeepCopy: true})
		require.NoError(t, err)
		enrollmentGet.SignatureAlgorithm = ""

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Once()

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentGet, nil).Times(3)

		var enrollmentUpdate cps.Enrollment
		err = copier.CopyWithOption(&enrollmentUpdate, enrollment, copier.Option{DeepCopy: true})
		require.NoError(t, err)
		enrollmentUpdate.AdminContact.FirstName = "R5"
		enrollmentUpdate.AdminContact.LastName = "D5"
		enrollmentUpdate.CSR.SANS = []string{"san2.test.akamai.com", "san.test.akamai.com"}
		enrollmentUpdate.NetworkConfiguration.DNSNameSettings.DNSNames = []string{"san2.test.akamai.com", "san.test.akamai.com"}
		enrollmentUpdate.Location = ""
		enrollmentUpdate.PendingChanges = nil
		enrollmentUpdate.SignatureAlgorithm = ""
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

		var enrollmentUpdateGet cps.Enrollment
		err = copier.CopyWithOption(&enrollmentUpdateGet, enrollmentUpdate, copier.Option{DeepCopy: true})
		require.NoError(t, err)
		enrollmentUpdateGet.SignatureAlgorithm = ""

		client.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{EnrollmentID: 1}).
			Return(&enrollmentUpdateGet, nil).Times(3)

		client.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
			EnrollmentID: 1,
			ChangeID:     2,
		}).Return(&cps.Change{
			AllowedInput: []cps.AllowedInput{{Type: "third-party-certificate"}},
			StatusInfo: &cps.StatusInfo{
				State:  "awaiting-input",
				Status: waitUploadThirdParty,
			},
		}, nil).Once()

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
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/lifecycle/create_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResThirdPartyEnrollment/lifecycle/update_enrollment.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_cps_third_party_enrollment.third_party", "contract_id", "ctr_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func getSimpleEnrollment() cps.Enrollment {
	return cps.Enrollment{
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
		CertificateType:      "third-party",
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
		EnableMultiStackedCertificates: true,
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
		RA:                 "third-party",
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
		ValidationType: "third-party",
		ThirdParty: &cps.ThirdParty{
			ExcludeSANS: false,
		},
	}
}
