package mtlskeystore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	intTest "github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClientCertificatesDataSource(t *testing.T) {
	t.Parallel()
	commonStateChecker := test.NewStateChecker("data.akamai_mtlskeystore_client_certificates.client_certificates").
		CheckEqual("certificates.0.certificate_id", "1").
		CheckEqual("certificates.0.certificate_name", "testName").
		CheckEqual("certificates.0.created_by", "testUser").
		CheckEqual("certificates.0.created_date", "2001-01-01T01:11:11Z").
		CheckEqual("certificates.0.geography", "CORE").
		CheckEqual("certificates.0.key_algorithm", "RSA").
		CheckEqual("certificates.0.signer", "AKAMAI").
		CheckEqual("certificates.0.subject", "/CN=clientauth_example_com_SSL_2024/O=Akamai Technologies Inc./OU=Media/C=US/").
		CheckEqual("certificates.0.secure_network", "STANDARD_TLS").
		CheckEqual("certificates.0.notification_emails.0", "test@example.com").
		CheckEqual("certificates.0.notification_emails.1", "test1@example.com").
		CheckEqual("certificates.1.certificate_id", "2").
		CheckEqual("certificates.1.certificate_name", "testName2").
		CheckEqual("certificates.1.created_by", "testUser2").
		CheckEqual("certificates.1.created_date", "2001-01-01T01:11:11Z").
		CheckEqual("certificates.1.geography", "CHINA_AND_CORE").
		CheckEqual("certificates.1.key_algorithm", "ECDSA").
		CheckEqual("certificates.1.signer", "THIRD_PARTY").
		CheckEqual("certificates.1.subject", "/CN=clientauth_example_com_SSL_2025/O=Akamai Technologies Inc./OU=Media/C=CN/").
		CheckEqual("certificates.1.secure_network", "ENHANCED_TLS").
		CheckEqual("certificates.1.notification_emails.0", "test2@example.com").
		CheckEqual("certificates.1.notification_emails.1", "test21@example.com")

	commonListClientCertificatesResponse := &mtlskeystore.ListClientCertificatesResponse{
		Certificates: []mtlskeystore.Certificate{
			{
				CertificateID:      1,
				CertificateName:    "testName",
				CreatedBy:          "testUser",
				CreatedDate:        intTest.NewTimeFromString(t, "2001-01-01T01:11:11Z"),
				Geography:          "CORE",
				KeyAlgorithm:       "RSA",
				Signer:             "AKAMAI",
				Subject:            "/CN=clientauth_example_com_SSL_2024/O=Akamai Technologies Inc./OU=Media/C=US/",
				SecureNetwork:      "STANDARD_TLS",
				NotificationEmails: []string{"test@example.com", "test1@example.com"},
			},
			{
				CertificateID:      2,
				CertificateName:    "testName2",
				CreatedBy:          "testUser2",
				CreatedDate:        intTest.NewTimeFromString(t, "2001-01-01T01:11:11Z"),
				Geography:          "CHINA_AND_CORE",
				KeyAlgorithm:       "ECDSA",
				Signer:             "THIRD_PARTY",
				Subject:            "/CN=clientauth_example_com_SSL_2025/O=Akamai Technologies Inc./OU=Media/C=CN/",
				SecureNetwork:      "ENHANCED_TLS",
				NotificationEmails: []string{"test2@example.com", "test21@example.com"},
			},
		},
	}
	tests := map[string]struct {
		init  func(*mtlskeystore.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificates(m, commonListClientCertificatesResponse, false)
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
  								edgerc = "../../common/testutils/edgerc"
							}

							data "akamai_mtlskeystore_client_certificates" "client_certificates" {
							}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"error API response": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificates(m, &mtlskeystore.ListClientCertificatesResponse{}, true)
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
  								edgerc = "../../common/testutils/edgerc"
							}

							data "akamai_mtlskeystore_client_certificates" "client_certificates" {
							}`,
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error - empty list of client certificates": {
			init: func(m *mtlskeystore.Mock) {
				m.On("ListClientCertificates", testutils.MockContext).Return(nil, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
  								edgerc = "../../common/testutils/edgerc"
							}

							data "akamai_mtlskeystore_client_certificates" "client_certificates" {
							}`,
					ExpectError: regexp.MustCompile("No Client Certificates"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
			if test.init != nil {
				test.init(client)
			}

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

func mockListClientCertificates(m *mtlskeystore.Mock, clientCertificates *mtlskeystore.ListClientCertificatesResponse, throwsError bool) {
	if throwsError {
		m.On("ListClientCertificates", testutils.MockContext).Return(nil, fmt.Errorf("oops")).Once()
		return
	}
	m.On("ListClientCertificates", testutils.MockContext).Return(clientCertificates, nil).Times(3)
}
