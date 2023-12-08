package cps

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const (
	enrollment1ID = 1
	enrollment2ID = 2
	enrollment3ID = 3
	enrollment4ID = 4
	changeID      = 2848126
)

var (
	enrollmentDV1 = &cps.GetEnrollmentResponse{
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
			C:                   "US",
			CN:                  "test.akamai.com",
			L:                   "Cambridge",
			O:                   "Akamai",
			OU:                  "WebEx",
			PreferredTrustChain: "intermediate-a",
			SANS:                []string{"san.test.akamai.com"},
			ST:                  "MA",
		},
		Location:                       "/cps/v2/enrollments/1",
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
		OrgID:              tools.IntPtr(123),
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
		ValidationType:  "dv",
		AssignedSlots:   []int{1},
		StagingSlots:    []int{2},
		ProductionSlots: []int{3},
	}
	enrollmentDV2 = &cps.GetEnrollmentResponse{
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
		Location:             "/cps/v2/enrollments/2",
		CertificateChainType: "default",
		CertificateType:      "san",
		ChangeManagement:     false,
		CSR: &cps.CSR{
			C:    "US",
			CN:   "akatest.com",
			L:    "Cambridge",
			O:    "Akamai",
			OU:   "WebEx",
			SANS: []string{"san.test.akamai1.com", "san.test.akamai2.com", "san.test.akamai3.com"},
			ST:   "MA",
		},
		EnableMultiStackedCertificates: false,
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
		OrgID:              tools.IntPtr(123),
		RA:                 "lets-encrypt",
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
		PendingChanges: []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/27552/changes/2848126",
				ChangeType: "new-certificate",
			},
		},
		ValidationType:  "dv",
		AssignedSlots:   []int{1},
		StagingSlots:    []int{2},
		ProductionSlots: []int{3},
	}
	enrollmentThirdParty = &cps.GetEnrollmentResponse{
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
		Location:             "/cps/v2/enrollments/3",
		CertificateChainType: "default",
		CertificateType:      "third-party",
		ChangeManagement:     false,
		CSR: &cps.CSR{
			C:    "US",
			CN:   "akatest.com",
			L:    "Cambridge",
			O:    "Akamai",
			OU:   "WebEx",
			SANS: []string{"san.test.akamai1.com", "san.test.akamai2.com", "san.test.akamai3.com"},
			ST:   "MA",
		},
		EnableMultiStackedCertificates: false,
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
		OrgID:              tools.IntPtr(123),
		RA:                 "lets-encrypt",
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
		PendingChanges: []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/27552/changes/2848126",
				ChangeType: "new-certificate",
			},
		},
		ValidationType:  "third-party",
		AssignedSlots:   []int{1},
		StagingSlots:    []int{2},
		ProductionSlots: []int{3},
	}
	enrollmentEV = &cps.GetEnrollmentResponse{
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
		Location:             "/cps/v2/enrollments/4",
		CertificateChainType: "default",
		CertificateType:      "third-party",
		ChangeManagement:     false,
		CSR: &cps.CSR{
			C:    "US",
			CN:   "akatest.com",
			L:    "Cambridge",
			O:    "Akamai",
			OU:   "WebEx",
			SANS: []string{"san.test.akamai1.com", "san.test.akamai2.com", "san.test.akamai3.com"},
			ST:   "MA",
		},
		EnableMultiStackedCertificates: false,
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
		OrgID:              tools.IntPtr(123),
		RA:                 "lets-encrypt",
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
		PendingChanges: []cps.PendingChange{
			{
				Location:   "/cps/v2/enrollments/27552/changes/2848126",
				ChangeType: "new-certificate",
			},
		},
		ValidationType:  "ev",
		AssignedSlots:   []int{1},
		StagingSlots:    []int{2},
		ProductionSlots: []int{3},
	}
)

func TestDataEnrollment(t *testing.T) {
	tests := map[string]struct {
		enrollment   *cps.GetEnrollmentResponse
		enrollmentID int
		init         func(*testing.T, *cps.Mock)
		steps        []resource.TestStep
	}{
		"happy path without challenges": {
			enrollment:   enrollmentDV1,
			enrollmentID: enrollment1ID,
			init: func(t *testing.T, m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollment1ID,
				}).Return(enrollmentDV1, nil).Times(5)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataEnrollment/enrollment_without_challenges.tf"),
					Check:  checkAttributesForEnrollment(enrollmentDV1, enrollment1ID, mockEmptyChanges(), mockEmptyDVArray()),
				},
			},
		},
		"happy path with challenges": {
			enrollment:   enrollmentDV2,
			enrollmentID: enrollment2ID,
			init: func(t *testing.T, m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollment2ID,
				}).Return(enrollmentDV2, nil).Times(5)

				dvArray := mockDVArray()
				change := mockLetsEncryptChallenges()

				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
					ChangeID:     changeID,
					EnrollmentID: enrollment2ID,
				}).Return(change, nil).Times(5)

				m.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
					ChangeID:     changeID,
					EnrollmentID: enrollment2ID,
				}).Return(dvArray, nil).Times(5)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataEnrollment/enrollment_with_challenges.tf"),
					Check:  checkAttributesForEnrollment(enrollmentDV2, enrollment2ID, mockLetsEncryptChallenges(), mockDVArray()),
				},
			},
		},
		"could not fetch an enrollment": {
			enrollment:   enrollmentDV1,
			enrollmentID: enrollment1ID,
			init: func(t *testing.T, m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollment1ID,
				}).Return(nil, fmt.Errorf("could not get an enrollment")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataEnrollment/enrollment_without_challenges.tf"),
					ExpectError: regexp.MustCompile("could not get an enrollment"),
				},
			},
		},
		"could not fetch a change status": {
			enrollment:   enrollmentDV2,
			enrollmentID: enrollment2ID,
			init: func(t *testing.T, m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollment2ID,
				}).Return(enrollmentDV2, nil).Once()

				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
					ChangeID:     changeID,
					EnrollmentID: enrollment2ID,
				}).Return(nil, fmt.Errorf("could not get a change status")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataEnrollment/enrollment_with_challenges.tf"),
					ExpectError: regexp.MustCompile("could not get a change status"),
				},
			},
		},
		"no changes on lets encrypt challenges": {
			enrollment:   enrollmentDV2,
			enrollmentID: enrollment2ID,
			init: func(t *testing.T, m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollment2ID,
				}).Return(enrollmentDV2, nil).Once()

				change := mockLetsEncryptChallenges()

				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
					ChangeID:     changeID,
					EnrollmentID: enrollment2ID,
				}).Return(change, nil).Once()

				m.On("GetChangeLetsEncryptChallenges", mock.Anything, cps.GetChangeRequest{
					ChangeID:     changeID,
					EnrollmentID: enrollment2ID,
				}).Return(nil, fmt.Errorf("could not get LetsEncrypt challenges")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataEnrollment/enrollment_with_challenges.tf"),
					ExpectError: regexp.MustCompile("could not get LetsEncrypt challenges"),
				},
			},
		},
		"third party change type": {
			enrollment:   enrollmentThirdParty,
			enrollmentID: enrollment3ID,
			init: func(t *testing.T, m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollment3ID,
				}).Return(enrollmentThirdParty, nil).Times(5)

				change := mockThirdPartyCSRChallenges()

				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
					ChangeID:     changeID,
					EnrollmentID: enrollment3ID,
				}).Return(change, nil).Times(5)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataEnrollment/enrollment_with_third_party_challenges.tf"),
					Check:  checkAttributesForEnrollment(enrollmentThirdParty, enrollment3ID, mockThirdPartyCSRChallenges(), mockThirdPartyCSRDVArray()),
				},
			},
		},
		"ev change type": {
			enrollment:   enrollmentEV,
			enrollmentID: enrollment4ID,
			init: func(t *testing.T, m *cps.Mock) {
				m.On("GetEnrollment", mock.Anything, cps.GetEnrollmentRequest{
					EnrollmentID: enrollment4ID,
				}).Return(enrollmentEV, nil).Times(5)

				change := mockEVChallenges()

				m.On("GetChangeStatus", mock.Anything, cps.GetChangeStatusRequest{
					ChangeID:     changeID,
					EnrollmentID: enrollment4ID,
				}).Return(change, nil).Times(5)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataEnrollment/enrollment_with_ev_challenges.tf"),
					Check:  checkAttributesForEnrollment(enrollmentEV, enrollment4ID, mockEVChallenges(), mockEVDVArray()),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cps.Mock{}
			test.init(t, client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkAttributesForEnrollment(en *cps.GetEnrollmentResponse, enID int, changes *cps.Change, dvArray *cps.DVArray) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkCommonAttrs(en, enID),
		checkSetTypeAttrs(en),
		checkChallenges(changes, dvArray),
		checkPendingChangesEnrollment(en),
	)
}

func checkCommonAttrs(en *cps.GetEnrollmentResponse, enID int) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		// Admin Contact
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.#", "1"),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.address_line_one", en.AdminContact.AddressLineOne),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.address_line_two", en.AdminContact.AddressLineTwo),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.city", en.AdminContact.City),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.country_code", en.AdminContact.Country),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.email", en.AdminContact.Email),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.first_name", en.AdminContact.FirstName),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.last_name", en.AdminContact.LastName),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.organization", en.AdminContact.OrganizationName),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.phone", en.AdminContact.Phone),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.postal_code", en.AdminContact.PostalCode),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "admin_contact.0.region", en.AdminContact.Region),
		// CSR
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "csr.#", "1"),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "csr.0.country_code", en.CSR.C),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "csr.0.organization", en.CSR.O),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "csr.0.organizational_unit", en.CSR.OU),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "csr.0.state", en.CSR.ST),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "csr.0.city", en.CSR.L),
		// Network Configuration
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.#", "1"),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.clone_dns_names", strconv.FormatBool(en.NetworkConfiguration.DNSNameSettings.CloneDNSNames)),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.geography", en.NetworkConfiguration.Geography),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.must_have_ciphers", en.NetworkConfiguration.MustHaveCiphers),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.ocsp_stapling", string(en.NetworkConfiguration.OCSPStapling)),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.preferred_ciphers", en.NetworkConfiguration.PreferredCiphers),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.quic_enabled", strconv.FormatBool(en.NetworkConfiguration.QuicEnabled)),
		// Organization
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.#", "1"),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.address_line_one", en.Org.AddressLineOne),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.address_line_two", en.Org.AddressLineTwo),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.city", en.Org.City),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.country_code", en.Org.Country),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.name", en.Org.Name),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.phone", en.Org.Phone),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.postal_code", en.Org.PostalCode),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "organization.0.region", en.Org.Region),
		// TechContact
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.#", "1"),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.address_line_one", en.TechContact.AddressLineOne),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.address_line_two", en.TechContact.AddressLineTwo),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.city", en.TechContact.City),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.country_code", en.TechContact.Country),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.email", en.TechContact.Email),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.first_name", en.TechContact.FirstName),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.last_name", en.TechContact.LastName),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.organization", en.TechContact.OrganizationName),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.phone", en.TechContact.Phone),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.postal_code", en.TechContact.PostalCode),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.region", en.TechContact.Region),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "tech_contact.0.title", en.TechContact.Title),
		// Other
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "validation_type", en.ValidationType),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "id", strconv.Itoa(enID)),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "certificate_chain_type", en.CertificateChainType),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "certificate_type", en.CertificateType),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "common_name", en.CSR.CN),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "enable_multi_stacked_certificates", strconv.FormatBool(en.EnableMultiStackedCertificates)),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "registration_authority", en.RA),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "signature_algorithm", en.SignatureAlgorithm),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "sni_only", strconv.FormatBool(en.NetworkConfiguration.SNIOnly)),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "secure_network", en.NetworkConfiguration.SecureNetwork),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "org_id", strconv.Itoa(*en.OrgID)),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "assigned_slots.#", strconv.Itoa(len(en.AssignedSlots))),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "staging_slots.#", strconv.Itoa(len(en.StagingSlots))),
		resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "production_slots.#", strconv.Itoa(len(en.ProductionSlots))),
	)
}

func checkSetTypeAttrs(en *cps.GetEnrollmentResponse) resource.TestCheckFunc {
	var checkFunctions []resource.TestCheckFunc
	sansCount := len(en.CSR.SANS)
	checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "sans.#", strconv.Itoa(sansCount)))

	for i := 0; i < sansCount; i++ {
		checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("sans.%v", i), en.CSR.SANS[i]))
	}

	if en.NetworkConfiguration.ClientMutualAuthentication != nil {
		checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.client_mutual_authentication.#", "1"))
		checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.client_mutual_authentication.0.set_id", en.NetworkConfiguration.ClientMutualAuthentication.SetID))
		checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.client_mutual_authentication.0.send_ca_list_to_client", strconv.FormatBool(*en.NetworkConfiguration.ClientMutualAuthentication.AuthenticationOptions.SendCAListToClient)))
		checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.client_mutual_authentication.0.set_id", strconv.FormatBool(*en.NetworkConfiguration.ClientMutualAuthentication.AuthenticationOptions.OCSP.Enabled)))
	}

	disallowedTLSVersionsNumber := len(en.NetworkConfiguration.DisallowedTLSVersions)
	checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "network_configuration.0.disallowed_tls_versions.#", strconv.Itoa(disallowedTLSVersionsNumber)))
	for i := 0; i < disallowedTLSVersionsNumber; i++ {
		checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("network_configuration.0.disallowed_tls_versions.%v", i), en.NetworkConfiguration.DisallowedTLSVersions[i]))
	}
	return resource.ComposeAggregateTestCheckFunc(checkFunctions...)
}

func checkChallenges(changes *cps.Change, dvArray *cps.DVArray) resource.TestCheckFunc {
	if len(changes.AllowedInput) < 1 || changes.AllowedInput[0].Type != "lets-encrypt-challenges" || len(dvArray.DV) == 0 {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "http_challenges.#", "0"),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "dns_challenges.#", "0"),
		)
	}

	numOfHTTPChanges := calculateNumberOfChanges(dvArray, "http-01")
	numOfDNSChanges := calculateNumberOfChanges(dvArray, "dns-01")
	var checkFunctions []resource.TestCheckFunc

	checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "http_challenges.#", strconv.Itoa(numOfHTTPChanges)))
	checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "dns_challenges.#", strconv.Itoa(numOfDNSChanges)))

	for _, dv := range dvArray.DV {
		if dv.ValidationStatus == "VALIDATED" {
			continue
		}
		httpCounter := 0
		dnsCounter := 0
		for _, challenge := range dv.Challenges {
			if challenge.Status != "pending" {
				continue
			}
			if challenge.Type == "http-01" {
				checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("http_challenges.%v.full_path", httpCounter), challenge.FullPath))
				checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("http_challenges.%v.response_body", httpCounter), challenge.ResponseBody))
				checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("http_challenges.%v.domain", httpCounter), dv.Domain))
				httpCounter++
			}
			if challenge.Type == "dns-01" {
				checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("dns_challenges.%v.full_path", dnsCounter), challenge.FullPath))
				checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("dns_challenges.%v.response_body", dnsCounter), challenge.ResponseBody))
				checkFunctions = append(checkFunctions, resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", fmt.Sprintf("dns_challenges.%v.domain", dnsCounter), dv.Domain))
				dnsCounter++
			}
		}
	}
	return resource.ComposeAggregateTestCheckFunc(checkFunctions...)
}

func checkPendingChangesEnrollment(en *cps.GetEnrollmentResponse) resource.TestCheckFunc {
	if len(en.PendingChanges) > 0 {
		return resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "pending_changes", "true")
	}
	return resource.TestCheckResourceAttr("data.akamai_cps_enrollment.test", "pending_changes", "false")
}

func calculateNumberOfChanges(dvArray *cps.DVArray, changeType string) int {
	counter := 0
	for _, dv := range dvArray.DV {
		if dv.ValidationStatus == "VALIDATED" {
			continue
		}
		for _, challenge := range dv.Challenges {
			if challenge.Status != "pending" {
				continue
			}
			if challenge.Type == changeType {
				counter++
			}
		}
	}
	return counter
}
