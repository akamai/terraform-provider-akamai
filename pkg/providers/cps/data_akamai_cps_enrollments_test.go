package cps

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/providers/cps/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

const contractID = "testing"

var (
	enrollmentsList = &cps.ListEnrollmentsResponse{
		Enrollments: []cps.Enrollment{*enrollmentDV1, *enrollmentDV2},
	}
	emptyEnrollmentList       = &cps.ListEnrollmentsResponse{}
	enrollmentsThirdPartyList = &cps.ListEnrollmentsResponse{
		Enrollments: []cps.Enrollment{*enrollmentDV1, *enrollmentDV2, *enrollmentThirdParty, *enrollmentEV},
	}
)

func TestDataEnrollments(t *testing.T) {
	tests := map[string]struct {
		contractID  string
		enrollments cps.ListEnrollmentsResponse
		init        func(*testing.T, *cps.Mock)
		steps       []resource.TestStep
	}{
		"happy path": {
			enrollments: cps.ListEnrollmentsResponse{Enrollments: []cps.Enrollment{*enrollmentDV1, *enrollmentDV2}},
			init: func(t *testing.T, m *cps.Mock) {
				m.On("ListEnrollments", mock.Anything, cps.ListEnrollmentsRequest{
					ContractID: contractID,
				}).Return(enrollmentsList, nil).Times(5)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataEnrollments/enrollments.tf"),
					Check:  checkAttrsForListEnrollments(enrollmentsList),
				},
			},
		},
		"could not fetch list of enrollments": {
			enrollments: cps.ListEnrollmentsResponse{Enrollments: []cps.Enrollment{*enrollmentDV1, *enrollmentDV2}},
			init: func(t *testing.T, m *cps.Mock) {
				m.On("ListEnrollments", mock.Anything, cps.ListEnrollmentsRequest{
					ContractID: contractID,
				}).Return(nil, fmt.Errorf("could not get list of enrollments")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("testdata/TestDataEnrollments/enrollments.tf"),
					ExpectError: regexp.MustCompile("could not get list of enrollments"),
				},
			},
		},
		"different change type enrollments": {
			enrollments: cps.ListEnrollmentsResponse{Enrollments: []cps.Enrollment{*enrollmentDV1, *enrollmentDV2, *enrollmentThirdParty, *enrollmentEV}},
			init: func(t *testing.T, m *cps.Mock) {
				m.On("ListEnrollments", mock.Anything, cps.ListEnrollmentsRequest{
					ContractID: contractID,
				}).Return(enrollmentsThirdPartyList, nil).Times(5)

			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataEnrollments/enrollments.tf"),
					Check:  checkAttrsForListEnrollments(enrollmentsThirdPartyList),
				},
			},
		},
		"no enrollments for given contract": {
			enrollments: cps.ListEnrollmentsResponse{},
			init: func(t *testing.T, m *cps.Mock) {
				m.On("ListEnrollments", mock.Anything, cps.ListEnrollmentsRequest{
					ContractID: contractID,
				}).Return(emptyEnrollmentList, nil).Times(10)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataEnrollments/enrollments.tf"),
					Check:  resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", "contract_id", contractID),
				},
				{
					Config: loadFixtureString("testdata/TestDataEnrollments/enrollments.tf"),
					Check:  resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", "enrollments.#", "0"),
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
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps:      test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkAttrsForListEnrollments(enrollments *cps.ListEnrollmentsResponse) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkCommonAttrsForListEnrollments(enrollments),
		checkSetTypeAttrsForListEnrollments(enrollments),
		checkPendingChangesAttr(enrollments),
	)
}

func checkCommonAttrsForListEnrollments(enrollments *cps.ListEnrollmentsResponse) resource.TestCheckFunc {
	var enrollmentsComposedCheckFuncs []resource.TestCheckFunc
	for i, en := range enrollments.Enrollments {
		enID, err := tools.GetEnrollmentID(en.Location)
		if err != nil {
			return nil
		}
		enrollmentCheckFuncs := resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.#", i), "1"),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.address_line_one", i), en.AdminContact.AddressLineOne),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.address_line_two", i), en.AdminContact.AddressLineTwo),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.city", i), en.AdminContact.City),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.country_code", i), en.AdminContact.Country),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.email", i), en.AdminContact.Email),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.first_name", i), en.AdminContact.FirstName),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.last_name", i), en.AdminContact.LastName),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.organization", i), en.AdminContact.OrganizationName),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.phone", i), en.AdminContact.Phone),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.postal_code", i), en.AdminContact.PostalCode),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.admin_contact.0.region", i), en.AdminContact.Region),
			// CSR
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.csr.#", i), "1"),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.csr.0.country_code", i), en.CSR.C),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.csr.0.organization", i), en.CSR.O),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.csr.0.organizational_unit", i), en.CSR.OU),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.csr.0.state", i), en.CSR.ST),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.csr.0.city", i), en.CSR.L),
			// Network Configuration
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.#", i), "1"),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.clone_dns_names", i), strconv.FormatBool(en.NetworkConfiguration.DNSNameSettings.CloneDNSNames)),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.geography", i), en.NetworkConfiguration.Geography),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.must_have_ciphers", i), en.NetworkConfiguration.MustHaveCiphers),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.ocsp_stapling", i), string(en.NetworkConfiguration.OCSPStapling)),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.preferred_ciphers", i), en.NetworkConfiguration.PreferredCiphers),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.quic_enabled", i), strconv.FormatBool(en.NetworkConfiguration.QuicEnabled)),
			// Organization
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.#", i), "1"),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.address_line_one", i), en.Org.AddressLineOne),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.address_line_two", i), en.Org.AddressLineTwo),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.city", i), en.Org.City),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.country_code", i), en.Org.Country),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.name", i), en.Org.Name),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.phone", i), en.Org.Phone),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.postal_code", i), en.Org.PostalCode),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.organization.0.region", i), en.Org.Region),
			// TechContact
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.#", i), "1"),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.address_line_one", i), en.TechContact.AddressLineOne),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.address_line_two", i), en.TechContact.AddressLineTwo),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.city", i), en.TechContact.City),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.country_code", i), en.TechContact.Country),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.email", i), en.TechContact.Email),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.first_name", i), en.TechContact.FirstName),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.last_name", i), en.TechContact.LastName),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.organization", i), en.TechContact.OrganizationName),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.phone", i), en.TechContact.Phone),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.postal_code", i), en.TechContact.PostalCode),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.region", i), en.TechContact.Region),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.tech_contact.0.title", i), en.TechContact.Title),
			// Other
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.validation_type", i), en.ValidationType),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.enrollment_id", i), strconv.Itoa(enID)),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.certificate_chain_type", i), en.CertificateChainType),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.certificate_type", i), en.CertificateType),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.common_name", i), en.CSR.CN),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.enable_multi_stacked_certificates", i), strconv.FormatBool(en.EnableMultiStackedCertificates)),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.registration_authority", i), en.RA),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.signature_algorithm", i), en.SignatureAlgorithm),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.sni_only", i), strconv.FormatBool(en.NetworkConfiguration.SNIOnly)),
			resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.secure_network", i), en.NetworkConfiguration.SecureNetwork),
		)
		enrollmentsComposedCheckFuncs = append(enrollmentsComposedCheckFuncs, enrollmentCheckFuncs)
	}
	return resource.ComposeAggregateTestCheckFunc(enrollmentsComposedCheckFuncs...)
}

func checkSetTypeAttrsForListEnrollments(enrollments *cps.ListEnrollmentsResponse) resource.TestCheckFunc {
	var enrollmentsComposedCheckFuncs []resource.TestCheckFunc
	for i, en := range enrollments.Enrollments {
		var enrollmentCheckFuncs []resource.TestCheckFunc
		sansCount := len(en.CSR.SANS)
		enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.sans.#", i), strconv.Itoa(sansCount)))

		for j := 0; j < sansCount; j++ {
			enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.sans.%v", i, j), en.CSR.SANS[j]))
		}
		if en.NetworkConfiguration.ClientMutualAuthentication != nil {
			enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.client_mutual_authentication.#", i), "1"))
			enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.client_mutual_authentication.0.set_id", i), en.NetworkConfiguration.ClientMutualAuthentication.SetID))
			enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.client_mutual_authentication.0.send_ca_list_to_client", i), strconv.FormatBool(*en.NetworkConfiguration.ClientMutualAuthentication.AuthenticationOptions.SendCAListToClient)))
			enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.client_mutual_authentication.0.set_id", i), strconv.FormatBool(*en.NetworkConfiguration.ClientMutualAuthentication.AuthenticationOptions.OCSP.Enabled)))
		}
		disallowedTLSVersionsNumber := len(en.NetworkConfiguration.DisallowedTLSVersions)
		enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.disallowed_tls_versions.#", i), strconv.Itoa(disallowedTLSVersionsNumber)))
		for j := 0; j < disallowedTLSVersionsNumber; j++ {
			enrollmentCheckFuncs = append(enrollmentCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.network_configuration.0.disallowed_tls_versions.%v", i, j), en.NetworkConfiguration.DisallowedTLSVersions[j]))
		}
		enrollmentsComposedCheckFuncs = append(enrollmentsComposedCheckFuncs, enrollmentCheckFuncs...)
	}
	return resource.ComposeAggregateTestCheckFunc(enrollmentsComposedCheckFuncs...)
}

func checkPendingChangesAttr(enrollments *cps.ListEnrollmentsResponse) resource.TestCheckFunc {
	var enrollmentComposedCheckFuncs []resource.TestCheckFunc
	for i, en := range enrollments.Enrollments {
		if len(en.PendingChanges) > 0 {
			enrollmentComposedCheckFuncs = append(enrollmentComposedCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.pending_changes", i), "true"))
		} else {
			enrollmentComposedCheckFuncs = append(enrollmentComposedCheckFuncs, resource.TestCheckResourceAttr("data.akamai_cps_enrollments.test", fmt.Sprintf("enrollments.%v.pending_changes", i), "false"))
		}
	}
	return resource.ComposeAggregateTestCheckFunc(enrollmentComposedCheckFuncs...)
}
