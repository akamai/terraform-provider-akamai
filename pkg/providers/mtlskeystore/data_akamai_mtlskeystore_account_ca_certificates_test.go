package mtlskeystore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/mtlskeystore"
	tst "github.com/akamai/terraform-provider-akamai/v7/internal/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccountCACertificatesDataSource(t *testing.T) {
	t.Parallel()
	baseChecker := test.NewStateChecker("data.akamai_mtlskeystore_account_ca_certificates.test").
		CheckEqual("certificates.0.account_id", "test-id1").
		CheckEqual("certificates.0.certificate", "-----BEGIN CERTIFICATE-----\ntest0\n-----END CERTIFICATE-----").
		CheckEqual("certificates.0.common_name", "test-name").
		CheckEqual("certificates.0.created_by", "jkowalski").
		CheckEqual("certificates.0.created_date", "2024-05-18T23:08:07Z").
		CheckEqual("certificates.0.expiry_date", "2027-05-18T23:08:07Z").
		CheckEqual("certificates.0.id", "1234").
		CheckEqual("certificates.0.issued_date", "2024-05-18T23:08:07Z").
		CheckEqual("certificates.0.key_algorithm", "RSA").
		CheckEqual("certificates.0.key_size_in_bytes", "4096").
		CheckEqual("certificates.0.qualification_date", "2024-06-18T21:08:41Z").
		CheckEqual("certificates.0.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("certificates.0.status", "CURRENT").
		CheckEqual("certificates.0.subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/").
		CheckEqual("certificates.0.version", "123").
		CheckEqual("certificates.1.account_id", "test-id2").
		CheckEqual("certificates.1.certificate", "-----BEGIN CERTIFICATE-----\ntest1\n-----END CERTIFICATE-----").
		CheckEqual("certificates.1.common_name", "test-name").
		CheckEqual("certificates.1.created_by", "jkowalski").
		CheckEqual("certificates.1.created_date", "2024-05-18T22:17:41Z").
		CheckEqual("certificates.1.expiry_date", "2027-05-18T22:17:41Z").
		CheckEqual("certificates.1.id", "12345").
		CheckEqual("certificates.1.issued_date", "2024-05-18T22:17:41Z").
		CheckEqual("certificates.1.key_algorithm", "RSA").
		CheckEqual("certificates.1.key_size_in_bytes", "4096").
		CheckEqual("certificates.1.signature_algorithm", "SHA256_WITH_RSA").
		CheckEqual("certificates.1.status", "PREVIOUS").
		CheckEqual("certificates.1.subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/").
		CheckEqual("certificates.1.version", "124")

	baseResponse := &mtlskeystore.ListAccountCACertificatesResponse{
		Certificates: []mtlskeystore.AccountCACertificate{
			{
				AccountID:          "test-id1",
				Certificate:        "-----BEGIN CERTIFICATE-----\ntest0\n-----END CERTIFICATE-----",
				CommonName:         "test-name",
				CreatedBy:          "jkowalski",
				CreatedDate:        tst.NewTimeFromString(t, "2024-05-18T23:08:07Z"),
				ExpiryDate:         tst.NewTimeFromString(t, "2027-05-18T23:08:07Z"),
				ID:                 1234,
				IssuedDate:         tst.NewTimeFromString(t, "2024-05-18T23:08:07Z"),
				KeyAlgorithm:       "RSA",
				KeySizeInBytes:     4096,
				QualificationDate:  ptr.To(tst.NewTimeFromString(t, "2024-06-18T21:08:41Z")),
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "CURRENT",
				Subject:            "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/",
				Version:            123,
			},
			{
				AccountID:          "test-id2",
				Certificate:        "-----BEGIN CERTIFICATE-----\ntest1\n-----END CERTIFICATE-----",
				CommonName:         "test-name",
				CreatedBy:          "jkowalski",
				CreatedDate:        tst.NewTimeFromString(t, "2024-05-18T22:17:41Z"),
				ExpiryDate:         tst.NewTimeFromString(t, "2027-05-18T22:17:41Z"),
				ID:                 12345,
				IssuedDate:         tst.NewTimeFromString(t, "2024-05-18T22:17:41Z"),
				KeyAlgorithm:       "RSA",
				KeySizeInBytes:     4096,
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "PREVIOUS",
				Subject:            "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/",
				Version:            124,
			},
			{
				AccountID:          "test-id3",
				Certificate:        "-----BEGIN CERTIFICATE-----\ntest2\n-----END CERTIFICATE-----",
				CommonName:         "test-name",
				CreatedBy:          "jsmith",
				CreatedDate:        tst.NewTimeFromString(t, "2024-05-18T21:08:41Z"),
				ExpiryDate:         tst.NewTimeFromString(t, "2027-05-18T21:08:41Z"),
				ID:                 123455,
				IssuedDate:         tst.NewTimeFromString(t, "2024-05-18T21:08:41Z"),
				KeyAlgorithm:       "RSA",
				KeySizeInBytes:     4096,
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "EXPIRED",
				Subject:            "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/",
				Version:            125,
			},
		},
	}

	filteredResponse := &mtlskeystore.ListAccountCACertificatesResponse{
		Certificates: []mtlskeystore.AccountCACertificate{
			{
				AccountID:          "test-id1",
				Certificate:        "-----BEGIN CERTIFICATE-----\ntest0\n-----END CERTIFICATE-----",
				CommonName:         "test-name",
				CreatedBy:          "jkowalski",
				CreatedDate:        tst.NewTimeFromString(t, "2024-05-18T23:08:07Z"),
				ExpiryDate:         tst.NewTimeFromString(t, "2027-05-18T23:08:07Z"),
				ID:                 1234,
				IssuedDate:         tst.NewTimeFromString(t, "2024-05-18T23:08:07Z"),
				KeyAlgorithm:       "RSA",
				KeySizeInBytes:     4096,
				QualificationDate:  ptr.To(tst.NewTimeFromString(t, "2024-06-18T21:08:41Z")),
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "CURRENT",
				Subject:            "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/",
				Version:            123,
			},
			{
				AccountID:          "test-id2",
				Certificate:        "-----BEGIN CERTIFICATE-----\ntest1\n-----END CERTIFICATE-----",
				CommonName:         "test-name",
				CreatedBy:          "jkowalski",
				CreatedDate:        tst.NewTimeFromString(t, "2024-05-18T22:17:41Z"),
				ExpiryDate:         tst.NewTimeFromString(t, "2027-05-18T22:17:41Z"),
				ID:                 12345,
				IssuedDate:         tst.NewTimeFromString(t, "2024-05-18T22:17:41Z"),
				KeyAlgorithm:       "RSA",
				KeySizeInBytes:     4096,
				SignatureAlgorithm: "SHA256_WITH_RSA",
				Status:             "PREVIOUS",
				Subject:            "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/",
				Version:            124,
			},
		},
	}
	tests := map[string]struct {
		init  func(*mtlskeystore.Mock)
		steps []resource.TestStep
	}{
		"happy path without statuses": {
			init: func(m *mtlskeystore.Mock) {
				mockListAccountCACertificates(m, []string{}, baseResponse)
			},
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlskeystore_account_ca_certificates" "test" {
}`,
					Check: baseChecker.
						CheckMissing("status").
						CheckEqual("certificates.#", "3").
						CheckEqual("certificates.2.account_id", "test-id3").
						CheckEqual("certificates.2.certificate", "-----BEGIN CERTIFICATE-----\ntest2\n-----END CERTIFICATE-----").
						CheckEqual("certificates.2.common_name", "test-name").
						CheckEqual("certificates.2.created_by", "jsmith").
						CheckEqual("certificates.2.created_date", "2024-05-18T21:08:41Z").
						CheckEqual("certificates.2.expiry_date", "2027-05-18T21:08:41Z").
						CheckEqual("certificates.2.id", "123455").
						CheckEqual("certificates.2.issued_date", "2024-05-18T21:08:41Z").
						CheckEqual("certificates.2.key_algorithm", "RSA").
						CheckEqual("certificates.2.key_size_in_bytes", "4096").
						CheckEqual("certificates.2.signature_algorithm", "SHA256_WITH_RSA").
						CheckEqual("certificates.2.status", "EXPIRED").
						CheckEqual("certificates.2.subject", "/C=US/O=Akamai Technologies, Inc./OU=Akamai mTLS/CN=test/").
						CheckEqual("certificates.2.version", "125").
						Build(),
				},
			},
		},
		"happy path with filtering statuses": {
			init: func(m *mtlskeystore.Mock) {
				mockListAccountCACertificates(m, []string{"CURRENT", "PREVIOUS"}, filteredResponse)
			},
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlskeystore_account_ca_certificates" "test" {
	status = ["CURRENT","PREVIOUS"]
}`,
					Check: baseChecker.
						CheckEqual("status.#", "2").
						CheckEqual("status.0", "CURRENT").
						CheckEqual("status.1", "PREVIOUS").
						CheckEqual("certificates.#", "2").
						Build(),
				},
			},
		},
		"happy path with filtering statuses - empty response": {
			init: func(m *mtlskeystore.Mock) {
				mockListAccountCACertificates(m, []string{"EXPIRED"}, &mtlskeystore.ListAccountCACertificatesResponse{Certificates: []mtlskeystore.AccountCACertificate{}})
			},
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlskeystore_account_ca_certificates" "test" {
	status = ["EXPIRED"]
}`,
					Check: test.NewStateChecker("data.akamai_mtlskeystore_account_ca_certificates.test").
						CheckEqual("status.#", "1").
						CheckEqual("status.0", "EXPIRED").
						CheckEqual("certificates.#", "0").
						Build(),
				},
			},
		},
		"error response from API": {
			init: func(m *mtlskeystore.Mock) {
				m.On("ListAccountCACertificates", testutils.MockContext, mtlskeystore.ListAccountCACertificatesRequest{}).Return(nil, fmt.Errorf("oops")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlskeystore_account_ca_certificates" "test" {
}`,
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"validation error - invalid status value": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
	  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlskeystore_account_ca_certificates" "test" {
	status = ["CURRENT","invalid-status"]
}`,
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Value Match`),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockListAccountCACertificates(m *mtlskeystore.Mock, statuses []string, response *mtlskeystore.ListAccountCACertificatesResponse) {
	var statusList []mtlskeystore.CertificateStatus
	for _, s := range statuses {
		statusList = append(statusList, mtlskeystore.CertificateStatus(s))
	}

	m.On("ListAccountCACertificates", testutils.MockContext, mtlskeystore.ListAccountCACertificatesRequest{
		Status: statusList,
	}).Return(response, nil).Times(3)
}
