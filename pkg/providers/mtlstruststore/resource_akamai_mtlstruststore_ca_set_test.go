package mtlstruststore

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ptr"
	tst "github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/internal/text"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testprovider"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type commonDataForResource struct {
	description             *string
	name                    string
	certificates            []mtlstruststore.CertificateResponse
	certificatesForResponse []mtlstruststore.CertificateResponse // populate only if different from certificates.
	caSetID                 string
	versionDescription      *string
	version                 int64
	stagingVersion          *int64
	stagingStatus           string
	productionVersion       *int64
	productionStatus        string
	newVersion              int64
	allowInsecureSHA1       bool
	properties              []mtlstruststore.AssociationProperty
	enrollments             []mtlstruststore.AssociationEnrollment
	caSetStatus             string
	validation              mtlstruststore.Validation
}

func TestCASetResource(t *testing.T) {
	t.Parallel()
	createData := commonDataForResource{
		caSetID:     "123456789",
		version:     1,
		name:        "set-1",
		description: ptr.To("Test CA Set for validation"),
		certificates: []mtlstruststore.CertificateResponse{
			{
				CertificatePEM:     "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
				Description:        ptr.To("Test certificate"),
				CreatedBy:          "johndoe",
				CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T16:01:02.555444Z"),
				EndDate:            tst.NewTimeFromStringMust("2026-04-16T16:01:02.555444Z"),
				Fingerprint:        "1234567890abcdef1234567890abcdef",
				Issuer:             "CN=Dummy CA",
				SerialNumber:       "987654321fedcba987654321fedcba",
				SignatureAlgorithm: "SHA256WITHRSA",
				StartDate:          tst.NewTimeFromStringMust("2025-04-17T16:01:02.555444Z"),
				Subject:            "CN=Dummy CA test",
			},
		},
		versionDescription: ptr.To("Initial version for testing"),
		stagingVersion:     nil,
		stagingStatus:      "INACTIVE",
		caSetStatus:        "NOT_DELETED",
	}
	baseCheck := test.NewStateChecker("akamai_mtlstruststore_ca_set.test").
		CheckEqual("name", "set-1").
		CheckEqual("description", "Test CA Set for validation").
		CheckEqual("account_id", "ACC-123456").
		CheckEqual("id", "123456789").
		CheckEqual("created_by", "someone").
		CheckEqual("created_date", "2025-04-16T12:08:34.099457Z").
		CheckEqual("version_created_by", "someone").
		CheckEqual("version_created_date", "2025-04-16T12:08:34.099457Z").
		CheckMissing("version_modified_by").
		CheckMissing("version_modified_date").
		CheckEqual("allow_insecure_sha1", "false").
		CheckEqual("version_description", "Initial version for testing").
		CheckEqual("latest_version", "1").
		CheckMissing("staging_version").
		CheckMissing("production_version")
	check := baseCheck.
		CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
		CheckEqual("certificates.0.description", "Test certificate").
		CheckEqual("certificates.0.created_by", "johndoe").
		CheckEqual("certificates.0.created_date", "2025-04-16T16:01:02.555444Z").
		CheckEqual("certificates.0.end_date", "2026-04-16T16:01:02.555444Z").
		CheckEqual("certificates.0.fingerprint", "1234567890abcdef1234567890abcdef").
		CheckEqual("certificates.0.issuer", "CN=Dummy CA").
		CheckEqual("certificates.0.serial_number", "987654321fedcba987654321fedcba").
		CheckEqual("certificates.0.signature_algorithm", "SHA256WITHRSA").
		CheckEqual("certificates.0.start_date", "2025-04-17T16:01:02.555444Z").
		CheckEqual("certificates.0.subject", "CN=Dummy CA test").
		CheckMissing("timeouts.delete")

	tests := map[string]struct {
		configPath string
		init       func(*mtlstruststore.Mock, commonDataForResource)
		mockData   commonDataForResource
		steps      []resource.TestStep
		error      *regexp.Regexp
	}{
		"create a ca set": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)
				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
			},
		},
		"create a ca set with duplicated certificates": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				resourceData1 := resourceData
				resourceData1.certificates = []mtlstruststore.CertificateResponse{
					{
						CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
						Description:    ptr.To("Test certificate"),
					},
					{
						CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
					},
				}
				resourceData1.certificatesForResponse = []mtlstruststore.CertificateResponse{
					{
						CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
						Description:    ptr.To("Test certificate"),
					},
				}
				resourceData1.validation = mtlstruststore.Validation{
					Warnings: []mtlstruststore.Warning{
						{
							ContextInfo: map[string]any{
								"description": ptr.To("m3d"),
								"fingerprint": "aa7b651c620ac2afba6ee4afa9b9ca09adbcab62f3fbdd0597fba9d6c5047cc5",
							},
							Detail:  "The certificate with the fingerprint aa7b651c620ac2afba6ee4afa9b9ca09adbcab62f3fbdd0597fba9d6c5047cc5 has been submitted more than once. Duplicate certificates are not allowed.",
							Pointer: "/certificates/1",
							Title:   "Duplicate certificate has been submitted in the certificates.",
							Type:    "/mtls-edge-truststore/error-types/duplicate-certificate",
						},
					},
				}
				mockValidateCertificates(m, resourceData1, nil).Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_duplicated_certs.tf"),
					ExpectError: regexp.MustCompile("Error: Certificates validation failed - Duplicate certificate has been submitted in the certificates.(\n|.)+" +
						"-----BEGIN CERTIFICATE-----(\\n|.)+" +
						"MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV(\\n|.)+" +
						"-----END CERTIFICATE-----"),
				},
			},
		},
		"create a ca set with cert provided by another resource": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(4) // one less because cert is unknown during planning phase.
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)
				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_external_cert.tf"),
					Check:  check.Build(),
				},
			},
		},
		"create a ca set with no description and allow insecure SHA1": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				resourceData.description = nil
				resourceData.versionDescription = nil
				resourceData.allowInsecureSHA1 = true
				resourceData.certificates = []mtlstruststore.CertificateResponse{
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
						CreatedBy:          "johndoe",
						CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T16:01:02.555444Z"),
						EndDate:            tst.NewTimeFromStringMust("2026-04-16T16:01:02.555444Z"),
						Fingerprint:        "1234567890abcdef1234567890abcdef",
						Issuer:             "CN=Dummy CA",
						SerialNumber:       "987654321fedcba987654321fedcba",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          tst.NewTimeFromStringMust("2025-04-17T16:01:02.555444Z"),
						Subject:            "CN=Dummy CA test",
					},
				}
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)
				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_no_description_insecure_sha1.tf"),
					Check: check.
						CheckEqual("allow_insecure_sha1", "true").
						CheckMissing("description").
						CheckMissing("version_description").
						CheckMissing("certificates.0.description").
						Build(),
				},
			},
		},
		"create a ca set, but ca set deletion is already in progress in Delete": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)
				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				// mock that the deletion of the ca set is already in progress.
				// No delete ca set call is made.
				resourceData.caSetStatus = "DELETING"
				mockGetCASet(m, resourceData).Once()
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
			},
		},
		"create a ca set, but ca set deletion is already completed in Delete": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)
				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				// mock that the deletion of the ca set is already completed.
				// No delete ca set call is made.
				resourceData.caSetStatus = "DELETED"
				mockGetCASet(m, resourceData).Once()
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
			},
		},
		"create a ca set - when version creation fails - taint the resource for next apply": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(4)
				mockCreateCASet(m, resourceData).Times(1)
				m.On("CreateCASetVersion", testutils.MockContext, mtlstruststore.CreateCASetVersionRequest{
					CASetID: resourceData.caSetID,
					Body: mtlstruststore.CreateCASetVersionRequestBody{
						AllowInsecureSHA1: resourceData.allowInsecureSHA1,
						Description:       resourceData.versionDescription,
						Certificates: []mtlstruststore.CertificateRequest{
							{
								CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
								Description:    ptr.To("Test certificate"),
							},
						},
					},
				}).Return(nil, fmt.Errorf("error creating CA set version")).Times(1)

				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)

				// delete.
				mockListCASetAssociations(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)

				// create.
				mockValidateCertificates(m, resourceData, nil).Times(4)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)

				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					ExpectError: regexp.MustCompile("error creating CA set version"),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
			},
		},
		"expect error - create a ca set with expired certificate": {
			init: func(m *mtlstruststore.Mock, _ commonDataForResource) {
				mockValidateCertificates(m, commonDataForResource{
					certificates: []mtlstruststore.CertificateResponse{
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIC3jCCAcagAwIBAgIBATANBgkqhkiG9w0BAQsFADAAMB4XDTI0MDcwMTAwMDAwMFoXDTI1MDcwMjEwNDcxNlowADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL1xDWbGQoeGrUimkp7KUnlj1w0+aHNs9QH9sXygMxxis9cNZeBewJv9fL7n2MmFSkgmsAxJ8/90G19cyfWzNnPtO9PF9kBVarPU79CUVPx9D3hJBPIlDKozrbZYy2H4HRbQ41xM9DF4DjXIqX3Lk8YslTf8SOSxgIpQLKVrdvIxSTY3uH+u0E67dtcTcz6Ytop1Z0u4Q7GesC6iUoqWYNNPRGTETN++kTZ1XqVXWoVWML4ffeHpqUqHm/ITY0OKeXcIMTD/lg0zFdMqMYY01Y76Vddgts5utqmgt7qJ6mWlETHpVXNiIn/ooukxCsAgxgfqS/iXEyOvWrmCT/O4rqECAwEAAaNjMGEwHQYDVR0OBBYEFNG9BugMXYKr404m1c4nIIuCbB6VMB8GA1UdIwQYMBaAFNG9BugMXYKr404m1c4nIIuCbB6VMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgGGMA0GCSqGSIb3DQEBCwUAA4IBAQB+0lz3a78X+Eg0QqLBUJBtNzuVwNMb2Yq3in0s+91QOMymfc5uVl55S/8RxInbTdwD51jkW9akAl3fpQu2PBRwLPJPpewHnWZJgVK/xws0TDJYHe0iWPpNTfQRU5QuciTAp5lwhyyzpamZK2uE76lYTwZ8y6ZHblPDp1JCu6k2soH0YkvTzrKSUJUki70jhVajEFEUZ8S19PQ8+UeEycxn6c629ZPgw87aej8SEbPiY60J1vq+o4px/9HpW9pGeZNMilIvvY9ezDtqERmC4mKbYPNzSFwkYJ5mqG4yUlafITpl6/nvPi9rihMcf06rEhhPye+nztQrscMYcvr7qaDh\n-----END CERTIFICATE-----\n",
							Description:    ptr.To("Test certificate"),
						},
					},
				},
					&mtlstruststore.Error{
						Type:     "/mtls-edge-truststore/error-types/certificate-validation-failure",
						Title:    "Certificates have failed validation.",
						Status:   400,
						Instance: "/mtls-edge-truststore/error-types/certificate-validation-failure/85817ae44a4fc762",
						Errors: []mtlstruststore.ErrorItem{
							{
								Detail:  "The certificate with subject  and fingerprint 69ecacc778efa8564ffbaf81200667e30fdf8ed59d9887c80591d4889dc86275 has expired. Expiry date is 2025-07-02T10:47:16.000000Z. The check was performed on 2025-08-13T13:07:06.000000Z.",
								Pointer: "/certificates/0",
								Title:   "The certificate has expired.",
								Type:    "/mtls-edge-truststore/error-types/expired-certificate",
								ContextInfo: map[string]any{
									"certificatePem": "-----BEGIN CERTIFICATE-----\nMIIC3jCCAcagAwIBAgIBATANBgkqhkiG9w0BAQsFADAAMB4XDTI0MDcwMTAwMDAwMFoXDTI1MDcwMjEwNDcxNlowADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL1xDWbGQoeGrUimkp7KUnlj1w0+aHNs9QH9sXygMxxis9cNZeBewJv9fL7n2MmFSkgmsAxJ8/90G19cyfWzNnPtO9PF9kBVarPU79CUVPx9D3hJBPIlDKozrbZYy2H4HRbQ41xM9DF4DjXIqX3Lk8YslTf8SOSxgIpQLKVrdvIxSTY3uH+u0E67dtcTcz6Ytop1Z0u4Q7GesC6iUoqWYNNPRGTETN++kTZ1XqVXWoVWML4ffeHpqUqHm/ITY0OKeXcIMTD/lg0zFdMqMYY01Y76Vddgts5utqmgt7qJ6mWlETHpVXNiIn/ooukxCsAgxgfqS/iXEyOvWrmCT/O4rqECAwEAAaNjMGEwHQYDVR0OBBYEFNG9BugMXYKr404m1c4nIIuCbB6VMB8GA1UdIwQYMBaAFNG9BugMXYKr404m1c4nIIuCbB6VMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgGGMA0GCSqGSIb3DQEBCwUAA4IBAQB+0lz3a78X+Eg0QqLBUJBtNzuVwNMb2Yq3in0s+91QOMymfc5uVl55S/8RxInbTdwD51jkW9akAl3fpQu2PBRwLPJPpewHnWZJgVK/xws0TDJYHe0iWPpNTfQRU5QuciTAp5lwhyyzpamZK2uE76lYTwZ8y6ZHblPDp1JCu6k2soH0YkvTzrKSUJUki70jhVajEFEUZ8S19PQ8+UeEycxn6c629ZPgw87aej8SEbPiY60J1vq+o4px/9HpW9pGeZNMilIvvY9ezDtqERmC4mKbYPNzSFwkYJ5mqG4yUlafITpl6/nvPi9rihMcf06rEhhPye+nztQrscMYcvr7qaDh\n-----END CERTIFICATE-----\n",
									"checkDate":      "2025-08-13T13:07:06.000000Z",
									"description":    nil,
									"expiryDate":     "2025-07-02T10:47:16.000000Z",
									"fingerprint":    "69ecacc778efa8564ffbaf81200667e30fdf8ed59d9887c80591d4889dc86275",
									"subject":        "",
								},
							},
						},
					})
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_expired_cert.tf"),
					ExpectError: regexp.MustCompile(`Certificates validation failed - The certificate has expired.` +
						`\n[\s\S]*?-----BEGIN CERTIFICATE-----\nMIIC3jCCAcagAwIBAgIBATANBgkqhkiG9w0BAQsFADAAMB4XDTI0MDcwMTAwMDAwMFoXDTI1MDcwMjEwNDcxNlowADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL1xDWbGQoeGrUimkp7KUnlj1w0\+aHNs9QH9sXygMxxis9cNZeBewJv9fL7n2MmFSkgmsAxJ8/90G19cyfWzNnPtO9PF9kBVarPU79CUVPx9D3hJBPIlDKozrbZYy2H4HRbQ41xM9DF4DjXIqX3Lk8YslTf8SOSxgIpQLKVrdvIxSTY3uH\+u0E67dtcTcz6Ytop1Z0u4Q7GesC6iUoqWYNNPRGTETN\+\+kTZ1XqVXWoVWML4ffeHpqUqHm/ITY0OKeXcIMTD/lg0zFdMqMYY01Y76Vddgts5utqmgt7qJ6mWlETHpVXNiIn/ooukxCsAgxgfqS/iXEyOvWrmCT/O4rqECAwEAAaNjMGEwHQYDVR0OBBYEFNG9BugMXYKr404m1c4nIIuCbB6VMB8GA1UdIwQYMBaAFNG9BugMXYKr404m1c4nIIuCbB6VMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgGGMA0GCSqGSIb3DQEBCwUAA4IBAQB\+0lz3a78X\+Eg0QqLBUJBtNzuVwNMb2Yq3in0s\+91QOMymfc5uVl55S/8RxInbTdwD51jkW9akAl3fpQu2PBRwLPJPpewHnWZJgVK/xws0TDJYHe0iWPpNTfQRU5QuciTAp5lwhyyzpamZK2uE76lYTwZ8y6ZHblPDp1JCu6k2soH0YkvTzrKSUJUki70jhVajEFEUZ8S19PQ8\+UeEycxn6c629ZPgw87aej8SEbPiY60J1vq\+o4px/9HpW9pGeZNMilIvvY9ezDtqERmC4mKbYPNzSFwkYJ5mqG4yUlafITpl6/nvPi9rihMcf06rEhhPye\+nztQrscMYcvr7qaDh\n-----END CERTIFICATE-----\n`),
				},
			},
		},
		"expect error - create a ca set with empty certificate pem": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create_empty_pem.tf"),
					ExpectError: regexp.MustCompile("Certificate must be in PEM format, got:"),
				},
			},
		},
		"expect error - create a ca set without certificates": {
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create_no_certs.tf"),
					ExpectError: regexp.MustCompile(`The argument "certificates" is required, but no definition was found.`),
				},
			},
		},
		"expect error - create a ca set with empty certificates description": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create_empty_cert_description.tf"),
					ExpectError: regexp.MustCompile("description\nstring length must be between 1 and 255, got: 0"),
				},
			},
		},
		"expect error - create a ca set with empty description": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create_empty_description.tf"),
					ExpectError: regexp.MustCompile("description string length must be between 1 and 255, got: 0"),
				},
			},
		},
		"expect error - create a ca set with empty version description": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create_empty_version_description.tf"),
					ExpectError: regexp.MustCompile("version_description string length must be between 1 and 255, got: 0"),
				},
			},
		},
		"expect error - create a ca set, first attempt to delete failed, there is no retry logic": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)
				// First attempt to delete failed.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "FAILED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "FAILED", "COMPLETE", "FAILED").Times(1)

				// Fake second attempt to delete to fulfill tests cleanup requirements.
				mockValidateCertificates(m, resourceData, nil).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Destroy:     true,
					ExpectError: regexp.MustCompile("contact support team to resolve the issue"),
				},
			},
		},
		"expect error - create a ca set, delete with associated properties should fail": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)
				// delete.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				resourceData.properties = []mtlstruststore.AssociationProperty{
					{
						PropertyID: "2",
						Hostnames: []mtlstruststore.AssociationHostname{
							{
								Hostname: "example-3.com",
							},
						},
					},
				}
				mockListCASetAssociations(m, resourceData).Times(2)

				// Fake second attempt to delete to fulfill tests cleanup requirements.
				mockValidateCertificates(m, resourceData, nil).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				resourceData.properties = []mtlstruststore.AssociationProperty{}
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Destroy:     true,
					ExpectError: regexp.MustCompile(`CA set is in use by 1 properties:\s+'' \(2\)`),
				},
			},
		},
		"expect error - create a ca set, delete with associated enrollments should fail": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockListCASetActivations(m, resourceData, false).Times(2)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(1)
				// delete.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				resourceData.enrollments = []mtlstruststore.AssociationEnrollment{
					{
						CN:              "some.example.com",
						EnrollmentID:    10430,
						EnrollmentLink:  "/cps/v2/enrollments/10430",
						ProductionSlots: []int64{},
						StagingSlots:    []int64{39352},
					},
				}
				mockListCASetAssociations(m, resourceData).Times(2)

				// Fake second attempt to delete to fulfill tests cleanup requirements.
				mockValidateCertificates(m, resourceData, nil).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				resourceData.enrollments = []mtlstruststore.AssociationEnrollment{}
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Destroy:     true,
					ExpectError: regexp.MustCompile(`CA set is in use by 1 enrollments:(\n|\s)+'some.example.com' \(10430\)`),
				},
			},
		},
		"expect error - name update": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)

				// read
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// update
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockValidateCertificates(m, resourceData, nil).Times(1)

				// delete
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)

			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/update_name.tf"),
					ExpectError: regexp.MustCompile("updating field `name` is not possible"),
				},
			},
		},
		"expect error - description update": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)

				// read
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// update
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockValidateCertificates(m, resourceData, nil).Times(1)

				// delete
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)

			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/update_description.tf"),
					ExpectError: regexp.MustCompile("updating field `description` is not possible"),
				},
			},
		},
		"expect error - null description update": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create
				resourceData.description = nil
				resourceData.versionDescription = nil
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)

				// read
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// update
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockValidateCertificates(m, resourceData, nil).Times(1)

				// delete
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)

			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_no_descriptions.tf"),
					Check: check.
						CheckMissing("description").
						CheckMissing("version_description").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("description"), knownvalue.Null()),
						},
					},
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/update_description.tf"),
					ExpectError: regexp.MustCompile("updating field `description` is not possible"),
				},
			},
		},
		"expect a few errors with the same title - create a ca set": {
			init: func(m *mtlstruststore.Mock, _ commonDataForResource) {
				mockValidateCertificates(m, commonDataForResource{
					certificates: []mtlstruststore.CertificateResponse{
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
							Description:    ptr.To("Incorrect PEM format"),
						},
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIE8jCCAtqgAwIBAgICB+kwDQYJKoZIhvcNAQELBQAwPDEPMA0GA1UECgwGQWth\nbWFpMQswCQYDVQQGEwJVUzELMAkGA1UECAwCTlkxDzANBgNVBAcMBkJvc3RvbjAe\nFw0yNTA2MDkxMjMzMzNaFw0yNjA2MDkxMjMzMzNaMDwxDzANBgNVBAoMBkFrYW1h\naTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAk5ZMQ8wDQYDVQQHDAZCb3N0b24wggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQChcMd+fNC6Kw9ZWLvZT+ftsOM1\n6IyEf7PbfX/4aiHT/iZug30AXZb1gOXXyQ04KlTtML9FKs2ZGb+0wNXy/hl4cKDQ\nVvrBC+/YFq38RTGSMOgtxjWRgFUltiXMBL/8V2HAbZauyEk+u76Lb190DTdrZZ86\nn0sYi3sBWrDzNRvAq8P7SLj8uu4CtETRk3j/mvH88PWjXuMZBd9JihRejqZWIEcj\ndJ/0CmH5f4/9/bsKGHiZKIKx5uHlzA5rLaB+cm12tqC5/8hrKJst9lfF0XELVjbJ\nqjlqxyvwzMIuqTwdzXY/6sO2DGtDJCHlXeIpR5kmfBYXflPQUQXgj04+GEhtc5WS\nYR0tCSMEphcxz00FntIA7/Eb88tPGsu7zqbj8YjEQ3vEK8fBwqhQpRkODL3LG26V\nVCU/4vuBKgddeH5BJpP3oL6kEiBDN4ZQvxtVt3st1FGsXjBh/NyV9QsK5ya1C7If\nW9GUew43YBpDW/5ZJ8bj74wIRP2XCcyP/5fETKwCtsySGK3fW5MwXpE6vrAOzfXt\n4+vCIVI+25L4AhyU5Q4QEA6WKgz/1YcXJW3/kEofAY9HMRuutLe2HjyfxPSxC6DE\n2MeeRiCoJjDfVg1zlQMHkxgZhSkahS4ln8IPvyTFUthgNqIm43K7s/eXl0qrmVF6\nH/V55bKO7WxcrHqkPwIDAQABMA0GCSqGSIb3DQEBCwUAA4ICAQCN0BD6S7t2jzD9\n/c8BLziTrKFVel4Q6aYEFH8PvWhffDR814+zBF1+FHMCttNr+hjX/WGGv7T0wFDj\nL+TKNDQETE7lErGbdFYByKyuqjxL8SRVIkRdTxgyoZghPt8qMmyHmQAis1B4xA4d\n1TAQBoQ4gp74ehJGaLC2a+FeEq6ta9ts3rjjGAMVZX7S2SX6+wz9BF/1c5d88vSC\nt3dINAJXj14eZnAlEk2cpZjVF1P+G7ydnC6l32idjFoMKqZ0ChRX2O+4qFEjzszA\n+W0eIzkklvRvAzFgvPPoUHQZk8tVD9oei/0G0fCSpcFw7z20GxwRwmQ1/cJeKWm3\ngR1TEjRS12iZhJJTTZLPONJu/hBSwPKPpqf9SlxHy3wtG2qIJjYXcF1t5k1KJf3K\nvKGA7AtN57lqrgCxklc+1GxoDFrFouhs238a6qvkp3xT4ytETZXMIf0J/xwO4koh\nMHZS0hF0M2/r6G2v5XQ5lStZ8MRuwec01c6F21yGm+EH/MhSQRQFzQg6ZbCWbLhG\nmyLkOJRWmO7YYrlGn1Oaqewm8a8x5M2Sa0bLEQZAbJSmQ7KBcJPbR7ZqAaYvg1PD\nHu3ea2zMgNP90Hxja068RHWNIliXLfU9NtqoiCKQZ5OfT2lK088PIw1lnptzSwcr\nJupDTFbp/p2GavKA7VC/3sLifWs7Vg==\n-----END CERTIFICATE-----",
							Description:    ptr.To("Test Certificate"),
						},
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIE8jCCAtqgAwIBAgICB+kwDQYJKoZIhvcNAQELBQAwPDEPMA0GA1UECgwGQWth\nbWFpMQswCQYDVQQGEwJVUzELMAkGA1UECAwCTlkxDzANBgNVBAcMBkJvc3RvbjAe\nFw0yNTA2MDkxMjMzMzNaFw0yNjA2MDkxMjMzMzNaMDwxDzANBgNVBAoMBkFrYW1h\naTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAk5ZMQ8wDQYDVQQHDAZCb3N0b24wggIi\nMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCx8wl0DoQEHjXBbOjyxINUDhf0\nWm4o3oLe/yYtx3jWx/Alg3eRk4Koi02TJTvp56KSihDV53vV8iZsr88ZaN8fKjPt\nur/eMqyIbrqyjqQacD4Z+cNYX8Csko5Bxst0Gh7pqxJeSYlOIuwtfx3BZxfOvqss\n9hzCR6YrKABlqL1pOq75sc5Lsc4w8Jyi8dNhGfiHXVagfLo3kz4utmyDQntJSmBs\nIw0QeRdXqPvA1f07cF2pv2Lq6Smr+4pKuMB+mp5cOxrawFrpKqmWUYXsVFPHwiVY\n+8wIaH2ywfnOaFWLh0ts5om3szFRlzrDistdc+mk06ofO3oN7t0iX23zVP71/jMW\nnnXHZdZ2UpgwCfX0XI7Gn579r5qAn6SfCc6D7XbShLvViAXzM67uDFKZbZkdZziy\nY8vFbwQ7VjytYotdwJxH6dRlg0/jM/nkQJxV3gtnotN7pIZduV9xQwalermVqsWL\nm/9M7FldY0GjzdEMKc3fCK+VJ9lBtuvw0NtSqUbnyV7srNmmYCWkEISjipb0/Dwp\nbng4NU5edgfS463UgwIvmJt7Nkx/Mz4N+uKUhkYyvMtnJgS03ZxwfFjA3ttn/fk+\nDVl0LhaAyGE0hfa6/u2sdhK01K9ZIYlz/11OqNoYWmfYqHPGnudXwj6ZRqzBqis0\nazfIMGwpkI/OvsuUgQIDAQABMA0GCSqGSIb3DQEBCwUAA4ICAQB5dIwvNQuXsF68\nSB9QleUqkIAdSllJg+RApAMeD3w0/6+yBIYOBeaTt/VOIa7mu1aOb0f0elvEIHaa\nfRVGoOv7ZpCmi+s1fkZqlGuV6aSy2HaGLToApru2Bml5Uy5CoRg1XVD23Eo5rKve\n5DOkeL50b0rBXobPxCkvsuQGf5wfVNsipR05QRy4m3SjVanQGVyhRGjo35CWfjzQ\nGIXY+eN6fP/1v+apNEHwj4suuftRlg9h+ZJnmAGDtznfSFzr4pBveyEv9ztbGO97\nHsBlobYOH9K/InKuu6QD17bCZuGFhgqfgXqQ0Q7rQLe410wotVT+1Q9xOyXIKnwX\nZ2sqmpB8Y74hZucAiGxv7yQonTwG6PJjMlOHiBPu66iLHnWFEI6d5frhnDwILxTc\negVlxl0FhcW9ZescC5+N6MoHARpQmYzKmeSAErOAxNihTkz87tYozffDw2thkOVu\nHGZFUxCKnYve3JKX0GKR9NVo4yFYH/rTfZy8pLFEmHeTx+98KnhZvNszVtsIa/kb\nnsM5bu0BTYTWEHvnWWnzPSOa1l324+4LGUBtilJot1573F1H1RJ71cYVEsqqNiOA\nZCDqgZG6hbrTSxhkC6YnHPYWZueX33IIsAqFPjmnzwf1w76gdD3NLEV7hBtkdI/Y\n9u9Gex1D2e718dr4lgnceYug3A6mJQ==\n-----END CERTIFICATE-----",
							Description:    ptr.To("Test Certificate"),
						},
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMjUxODA1MTBTBkFrYW1haTEUMBIGaFw0yNjA1MjUxODA1MTBaMCcxDzANBgNVBAo\nAQDc/kIqoeeZozLtDfR3uCt5UvkxXrz+5Y/w2KlKWJwyMv5GcKbYhmeYZxF93+ys\n7AKuWq3Wga2D3XIlNQeurrnY/2/hWjjzfXI9fKd3eIMw29ALcheDRpnJyrrwXdYj\nNZ24gVLcy0Xp0tkj9SABzjJyBxyulbayPXnRVhOBsifAN65d+HlV9+gQXiX7M70U\nSXv27A3roeg6M9BiOxrc+x7GtzSRWB1/vGYLn0zHolEDrzBq0kRlH53YRnsYzGpG\n8UGs5cFKGmAw2zkpEKvoQKCrYXNQaoDo+qQ7kzdf72xVUm9EcVR4i1bFVvIvPQWk\ndCURqs/A0+jN0JDhvUnpDopXAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBpjAPBgNV\nHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRsC7RVcqVK/tLQmGMYl58af6DL2DANBgkq\nhkiG9w0BAQsFAAOCAQEAJtLx39c2rVNA0oGrpeg50QN1oAT6uXMvnLud+lVgu9Wi\nG7wrqMFZNnrlO7Ke91jJM7JEjAJCb7sd9yrfpq2D3MSinfFYSkrJppoUqYQnCMIP\nDNyaPuj7DgMbduiOX3f67FNSnndndPPR2MwNyvdT3UMR05itzsxUQK29cgSWwVhz\n-----END CERTIFICATE-----",
							Description:    ptr.To("Incorrect PEM format"),
						},
					},
				},
					&mtlstruststore.Error{
						Errors: []mtlstruststore.ErrorItem{
							{
								ContextInfo: map[string]any{
									"certificatePem": "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----",
									"description":    "Incorrect PEM format",
								},
								Detail:  "The certificate is not in a properly encoded PEM format. It cannot be parsed.",
								Pointer: "/certificates/0",
								Title:   "The certificate is malformed. It cannot be parsed.",
								Type:    "/mtls-edge-truststore/error-types/malformed-certificate",
							},
							{
								ContextInfo: map[string]any{
									"certificatePem": "-----BEGIN CERTIFICATE-----\nMjUxODA1MTBTBkFrYW1haTEUMBIGaFw0yNjA1MjUxODA1MTBaMCcxDzANBgNVBAo\nAQDc/kIqoeeZozLtDfR3uCt5UvkxXrz+5Y/w2KlKWJwyMv5GcKbYhmeYZxF93+ys\n7AKuWq3Wga2D3XIlNQeurrnY/2/hWjjzfXI9fKd3eIMw29ALcheDRpnJyrrwXdYj\nNZ24gVLcy0Xp0tkj9SABzjJyBxyulbayPXnRVhOBsifAN65d+HlV9+gQXiX7M70U\nSXv27A3roeg6M9BiOxrc+x7GtzSRWB1/vGYLn0zHolEDrzBq0kRlH53YRnsYzGpG\n8UGs5cFKGmAw2zkpEKvoQKCrYXNQaoDo+qQ7kzdf72xVUm9EcVR4i1bFVvIvPQWk\ndCURqs/A0+jN0JDhvUnpDopXAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBpjAPBgNV\nHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRsC7RVcqVK/tLQmGMYl58af6DL2DANBgkq\nhkiG9w0BAQsFAAOCAQEAJtLx39c2rVNA0oGrpeg50QN1oAT6uXMvnLud+lVgu9Wi\nG7wrqMFZNnrlO7Ke91jJM7JEjAJCb7sd9yrfpq2D3MSinfFYSkrJppoUqYQnCMIP\nDNyaPuj7DgMbduiOX3f67FNSnndndPPR2MwNyvdT3UMR05itzsxUQK29cgSWwVhz\n-----END CERTIFICATE-----",
									"description":    "Incorrect PEM format",
								},
								Detail:  "The certificate is not in a properly encoded PEM format. It cannot be parsed.",
								Pointer: "/certificates/3",
								Title:   "The certificate is malformed. It cannot be parsed.",
								Type:    "/mtls-edge-truststore/error-types/malformed-certificate",
							},
						},
						Instance: "/mtls-edge-truststore/error-types/certificate-validation-failure/56e1f179ad2d5b20",
						Status:   http.StatusBadRequest,
						Title:    "Certificates have failed validation.",
						Type:     "/mtls-edge-truststore/error-types/certificate-validation-failure",
					})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_invalid_certificates.tf"),
					ExpectError: regexp.MustCompile(`Certificates validation failed - The certificate is malformed. It cannot be parsed.` +
						`\n[\s\S]*?Incorrect PEM format\n-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n\n\n` +
						`Incorrect PEM format\n-----BEGIN CERTIFICATE-----\nMjUxODA1MTBTBkFrYW1haTEUMBIGaFw0yNjA1MjUxODA1MTBaMCcxDzANBgNVBAo\nAQDc/kIqoeeZozLtDfR3uCt5UvkxXrz\+5Y/w2KlKWJwyMv5GcKbYhmeYZxF93\+ys\n7AKuWq3Wga2D3XIlNQeurrnY/2/hWjjzfXI9fKd3eIMw29ALcheDRpnJyrrwXdYj\nNZ24gVLcy0Xp0tkj9SABzjJyBxyulbayPXnRVhOBsifAN65d\+HlV9\+gQXiX7M70U\nSXv27A3roeg6M9BiOxrc\+x7GtzSRWB1/vGYLn0zHolEDrzBq0kRlH53YRnsYzGpG\n8UGs5cFKGmAw2zkpEKvoQKCrYXNQaoDo\+qQ7kzdf72xVUm9EcVR4i1bFVvIvPQWk\ndCURqs/A0\+jN0JDhvUnpDopXAgMBAAGjQjBAMA4GA1UdDwEB/wQEAwIBpjAPBgNV\nHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRsC7RVcqVK/tLQmGMYl58af6DL2DANBgkq\nhkiG9w0BAQsFAAOCAQEAJtLx39c2rVNA0oGrpeg50QN1oAT6uXMvnLud\+lVgu9Wi\nG7wrqMFZNnrlO7Ke91jJM7JEjAJCb7sd9yrfpq2D3MSinfFYSkrJppoUqYQnCMIP\nDNyaPuj7DgMbduiOX3f67FNSnndndPPR2MwNyvdT3UMR05itzsxUQK29cgSWwVhz\n-----END CERTIFICATE-----`),
				},
			},
		},
		"expect a few errors with different title - create a ca set": {
			init: func(m *mtlstruststore.Mock, _ commonDataForResource) {
				mockValidateCertificates(m, commonDataForResource{
					certificates: []mtlstruststore.CertificateResponse{
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV123\n-----END CERTIFICATE-----",
							Description:    ptr.To("Incorrect PEM format first group"),
						},
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV345\n-----END CERTIFICATE-----",
							Description:    ptr.To("Incorrect PEM format first group"),
						},
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV567\n-----END CERTIFICATE-----",
							Description:    ptr.To("Incorrect PEM format second group"),
						},
						{
							CertificatePEM: "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV789\n-----END CERTIFICATE-----",
							Description:    ptr.To("Incorrect PEM format second group"),
						},
					},
				},
					&mtlstruststore.Error{
						Errors: []mtlstruststore.ErrorItem{
							{
								ContextInfo: map[string]any{
									"certificatePem": "-----BEGIN CERTIFICATE-----MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV123-----END CERTIFICATE-----",
									"description":    "Incorrect PEM format first group",
								},
								Detail:  "The certificate is not in a properly encoded PEM format. It cannot be parsed.",
								Pointer: "/certificates/0",
								Title:   "The certificate is malformed. It cannot be parsed.",
								Type:    "/mtls-edge-truststore/error-types/malformed-certificate",
							},
							{
								ContextInfo: map[string]any{
									"certificatePem": "-----BEGIN CERTIFICATE-----MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV345-----END CERTIFICATE-----",
									"description":    "Incorrect PEM format first group",
								},
								Detail:  "The certificate is not in a properly encoded PEM format. It cannot be parsed.",
								Pointer: "/certificates/1",
								Title:   "The certificate is malformed. It cannot be parsed.",
								Type:    "/mtls-edge-truststore/error-types/malformed-certificate",
							},
							{
								ContextInfo: map[string]any{
									"certificatePem": "-----BEGIN CERTIFICATE-----MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV567-----END CERTIFICATE-----",
									"description":    "Incorrect PEM format second group",
								},
								Detail:  "The certificate with subject  and fingerprint 69ecacc778efa8564ffbaf81200667e30fdf8ed59d9887c80591d4889dc86275 has expired. Expiry date is 2025-07-02T10:47:16.000000Z. The check was performed on 2025-08-13T13:07:06.000000Z.",
								Pointer: "/certificates/2",
								Title:   "The certificate has expired.",
								Type:    "/mtls-edge-truststore/error-types/expired-certificate",
							},
							{
								ContextInfo: map[string]any{
									"certificatePem": "-----BEGIN CERTIFICATE-----MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV789-----END CERTIFICATE-----",
									"description":    "Incorrect PEM format second group",
								},
								Detail:  "The certificate with subject  and fingerprint 69ecacc778efa8564ffbaf81200667e30fdf8ed59d9887c80591d4889dc86275 has expired. Expiry date is 2025-07-02T10:47:16.000000Z. The check was performed on 2025-08-13T13:07:06.000000Z.",
								Pointer: "/certificates/3",
								Title:   "The certificate has expired.",
								Type:    "/mtls-edge-truststore/error-types/expired-certificate",
							},
						},
						Instance: "/mtls-edge-truststore/error-types/certificate-validation-failure/0e0bc044ebda5a16",
						Status:   http.StatusBadRequest,
						Title:    "Certificates have failed validation.",
						Type:     "/mtls-edge-truststore/error-types/certificate-validation-failure",
					})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_invalid_certificates_with_different_titles.tf"),
					ExpectError: regexp.MustCompile(`Certificates validation failed - The certificate is malformed. It cannot be parsed.` +
						`\n[\s\S]*?Incorrect PEM format first group\n[\s\S]*?MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV123` +
						`\n[\s\S]*?Incorrect PEM format first group\n[\s\S]*?MIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV345`),
				},
			},
		},
		"expect error - update an activated ca set on staging, but still updated config contains certificates which have expired": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(2)

				// update.
				var one int64 = 1
				updateData := resourceData
				updateData.version = one
				updateData.stagingVersion = ptr.To(one)
				updateData.stagingStatus = "ACTIVE"
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.certificates = slices.Insert(updateData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				updateData.versionDescription = ptr.To("Second version for testing")
				mockValidateCertificates(m, updateData,
					&mtlstruststore.Error{
						Type:        "/mtls-edge-truststore/error-types/certificate-validation-failure",
						Detail:      "Certificates have failed validation.",
						Status:      400,
						ContextInfo: nil,
						Instance:    "/mtls-edge-truststore/error-types/certificate-validation-failure/f0490e43aa97d24b",
						Errors: []mtlstruststore.ErrorItem{
							{
								ContextInfo: map[string]any{
									"checkDate":   "2025-08-04T08:52:16.000000Z",
									"description": nil,
									"expiryDate":  "2024-10-30T17:07:04.000000Z",
									"fingerprint": "81bd3e1660199559a8513b5efa9c95f93fb7e6a61bd283b87f6e5c128178d9b9",
									"subject":     "CN=cert100",
								},
								Detail:  "The certificate with subject CN=cert100 and fingerprint 81bd3e1660199559a8513b5efa9c95f93fb7e6a61bd283b87f6e5c128178d9b9 has expired. Expiry date is 2024-10-30T17:07:04.000000Z. The check was performed on 2025-08-04T08:52:16.000000Z.",
								Pointer: "/certificates/0",
								Title:   "The certificate has expired.",
								Type:    "/mtls-edge-truststore/error-types/expired-certificate",
							},
						},
					},
				).Times(1)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update.tf"),
					ExpectError: regexp.MustCompile("Error: Certificates validation failed - The certificate has expired.(\n|.)+" +
						"second cert(\n|.)+" +
						"-----BEGIN CERTIFICATE-----(\n|.)+" +
						"FOO(\n|.)+" +
						"-----END CERTIFICATE-----"),
				},
			},
		},
		"expect error - update an activated ca set on staging when there are already 100 versions": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(3)

				// update right to version 100.
				var currentVersion int64 = 1
				var newVersion int64 = 100
				updateData := resourceData
				updateData.version = currentVersion
				updateData.stagingVersion = ptr.To(currentVersion)
				updateData.stagingStatus = "ACTIVE"
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.certificates = slices.Insert(updateData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				updateData.versionDescription = ptr.To("Second version for testing")
				mockValidateCertificates(m, updateData, nil).Times(5)
				updateData.newVersion = newVersion
				mockCloneCASetVersion(m, updateData).Times(1)
				updateData.version = newVersion
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(4)

				// try to update again - expect error now.
				currentVersion = 100
				updateData2 := resourceData
				updateData2.version = currentVersion
				updateData2.stagingVersion = ptr.To(currentVersion)
				updateData2.stagingStatus = "ACTIVE"
				mockListCASetActivations(m, updateData2, true).Times(1)
				mockGetCASet(m, updateData2).Times(1)
				mockGetCASetVersion(m, updateData2).Times(1)
				updateData2.allowInsecureSHA1 = true
				mockValidateCertificates(m, updateData2, nil).Times(1)

				// delete.
				mockListCASetAssociations(m, updateData2).Times(2)
				mockGetCASet(m, updateData2).Once()
				mockDeleteCASet(m, updateData2).Times(1)
				mockGetCASetDeletionStatus(m, updateData2, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData2, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update.tf"),
					Check: baseCheck.
						CheckEqual("version_description", "Second version for testing").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("certificates.1.created_by", "johndoe").
						CheckEqual("certificates.1.created_date", "2025-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.end_date", "2026-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.fingerprint", "1234567890abcdef1234567890abcdef").
						CheckEqual("certificates.1.issuer", "CN=Dummy CA").
						CheckEqual("certificates.1.serial_number", "987654321fedcba987654321fedcba").
						CheckEqual("certificates.1.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("certificates.1.start_date", "2025-04-17T16:01:02.555444Z").
						CheckEqual("certificates.1.subject", "CN=Dummy CA test").
						CheckEqual("latest_version", "100").
						CheckEqual("staging_version", "1").
						CheckMissing("production_version").
						CheckEqual("timeouts.delete", "5m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("5m")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("description"), knownvalue.StringExact("Test CA Set for validation")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_description"), knownvalue.StringExact("Second version for testing")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),
						},
					},
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCASet/update_allow_insecure_sha1.tf"),
					ExpectError: regexp.MustCompile("Cannot create more than 100 versions for a CA Set."),
				},
			},
		},
		"update a non activated ca set": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)
				// update.
				updateData := resourceData
				updateData.version = 1
				mockListCASetActivations(m, updateData, false).Times(1)
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.certificates = slices.Insert(updateData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				updateData.versionDescription = ptr.To("Second version for testing")
				mockValidateCertificates(m, updateData, nil).Times(5)
				mockUpdateCASetVersion(m, updateData).Times(1)

				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update.tf"),
					Check: baseCheck.
						CheckEqual("version_description", "Second version for testing").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("certificates.1.created_by", "johndoe").
						CheckEqual("certificates.1.created_date", "2025-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.end_date", "2026-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.fingerprint", "1234567890abcdef1234567890abcdef").
						CheckEqual("certificates.1.issuer", "CN=Dummy CA").
						CheckEqual("certificates.1.serial_number", "987654321fedcba987654321fedcba").
						CheckEqual("certificates.1.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("certificates.1.start_date", "2025-04-17T16:01:02.555444Z").
						CheckEqual("certificates.1.subject", "CN=Dummy CA test").
						CheckEqual("timeouts.delete", "5m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),

							// The added certificate will be the first in the state.
							// We expect all computed attributes to be unknown except the description.
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("certificate_pem"),
								knownvalue.StringExact("-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("description"),
								knownvalue.StringExact("second cert")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("created_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("created_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("end_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("fingerprint")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("issuer")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("serial_number")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("signature_algorithm")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("start_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("subject")),

							// The existing certificate must be completely known at the plan phase.
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("certificate_pem"),
								knownvalue.StringExact("-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("description"),
								knownvalue.StringExact("Test certificate")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("created_by"), knownvalue.StringExact("johndoe")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("created_date"), knownvalue.StringExact("2025-04-16T16:01:02.555444Z")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("end_date"), knownvalue.StringExact("2026-04-16T16:01:02.555444Z")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("fingerprint"), knownvalue.StringExact("1234567890abcdef1234567890abcdef")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("issuer"), knownvalue.StringExact("CN=Dummy CA")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("serial_number"), knownvalue.StringExact("987654321fedcba987654321fedcba")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("signature_algorithm"), knownvalue.StringExact("SHA256WITHRSA")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("start_date"), knownvalue.StringExact("2025-04-17T16:01:02.555444Z")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("subject"), knownvalue.StringExact("CN=Dummy CA test")),
						},
					},
				},
			},
		},
		"update a non activated ca set, changing only allow_insecure_sha1": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)
				// update.
				updateData := resourceData
				updateData.version = 1
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.allowInsecureSHA1 = true
				mockValidateCertificates(m, updateData, nil).Times(5)
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(4)
				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_allow_insecure_sha1.tf"),
					Check: check.
						CheckEqual("allow_insecure_sha1", "true").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),
						},
					},
				},
			},
		},
		"update a non activated ca set, removing certificates description": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)
				// update
				updateData := resourceData
				updateData.version = 1
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.certificates = []mtlstruststore.CertificateResponse{
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n",
						Description:        nil,
						CreatedBy:          "johndoe",
						CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T16:01:02.555444Z"),
						EndDate:            tst.NewTimeFromStringMust("2026-04-16T16:01:02.555444Z"),
						Fingerprint:        "1234567890abcdef1234567890abcdef",
						Issuer:             "CN=Dummy CA",
						SerialNumber:       "987654321fedcba987654321fedcba",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          tst.NewTimeFromStringMust("2025-04-17T16:01:02.555444Z"),
						Subject:            "CN=Dummy CA test",
					},
				}
				mockValidateCertificates(m, updateData, nil).Times(5)
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(4)
				// delete
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_no_cert_description.tf"),
					Check: check.
						CheckMissing("certificates.0.description").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),

							// We expect all computed attributes to be unknown
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("certificate_pem"),
								knownvalue.StringExact("-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("description"),
								knownvalue.Null()),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("created_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("created_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("end_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("fingerprint")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("issuer")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("serial_number")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("signature_algorithm")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("start_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("subject")),
						},
					},
				},
			},
		},
		// This case tests whether unknown config value for `description` and `version description`.
		// is properly set to null in both create and update phases.
		"update a non activated ca set with no description and version description, changing only allow_insecure_sha1": {
			init: func(m *mtlstruststore.Mock, resourceDataIn commonDataForResource) {
				resourceData := resourceDataIn
				resourceData.description = nil
				resourceData.versionDescription = nil
				certs := append([]mtlstruststore.CertificateResponse(nil), resourceData.certificates...)
				certs[0].Description = nil
				resourceData.certificates = certs
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)
				// update.
				updateData := resourceData
				updateData.version = 1
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.allowInsecureSHA1 = true
				mockValidateCertificates(m, updateData, nil).Times(5)
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(4)
				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_no_all_descriptions.tf"),
					Check: check.
						CheckMissing("description").
						CheckMissing("version_description").
						CheckMissing("certificates.0.description").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("description"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_description"), knownvalue.Null()),
						},
					},
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_allow_insecure_sha1_no_descriptions.tf"),
					Check: check.
						CheckMissing("description").
						CheckMissing("version_description").
						CheckMissing("certificates.0.description").
						CheckEqual("allow_insecure_sha1", "true").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("description"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_description"), knownvalue.Null()),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),
						},
					},
				},
			},
		},
		"update only timeout should do nothing (from default to some value)": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// update.
				updateData := resourceData
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockValidateCertificates(m, updateData, nil).Times(4)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_timeout.tf"),
					Check: check.
						CheckEqual("timeouts.delete", "5m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("5m")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by"), knownvalue.StringExact("someone")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date"), knownvalue.StringExact("2025-04-16T12:08:34.099457Z")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date"), knownvalue.Null()),
						},
					},
				},
			},
		},
		"update only timeout should do nothing (from some value to default)": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// update.
				updateData := resourceData
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockValidateCertificates(m, updateData, nil).Times(4)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_timeout.tf"),
					Check: check.
						CheckEqual("timeouts.delete", "5m").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check: check.
						CheckMissing("timeouts.delete").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("timeouts"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by"), knownvalue.StringExact("someone")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date"), knownvalue.StringExact("2025-04-16T12:08:34.099457Z")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date"), knownvalue.Null()),
						},
					},
				},
			},
		},
		"update only timeout should do nothing (from some value to other value)": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// update.
				updateData := resourceData
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockValidateCertificates(m, updateData, nil).Times(4)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_timeout.tf"),
					Check: check.
						CheckEqual("timeouts.delete", "5m").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_timeout2.tf"),
					Check: check.
						CheckEqual("timeouts.delete", "6m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("6m")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by"), knownvalue.StringExact("someone")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date"), knownvalue.StringExact("2025-04-16T12:08:34.099457Z")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date"), knownvalue.Null()),
						},
					},
				},
			},
		},
		"update a non activated ca set with only order change in the config": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				resourceData.certificates = slices.Insert(resourceData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// update.
				updateData := resourceData
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockValidateCertificates(m, updateData, nil).Times(3)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create_two_certs.tf"),
					Check: baseCheck.
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("certificates.1.created_by", "johndoe").
						CheckEqual("certificates.1.created_date", "2025-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.end_date", "2026-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.fingerprint", "1234567890abcdef1234567890abcdef").
						CheckEqual("certificates.1.issuer", "CN=Dummy CA").
						CheckEqual("certificates.1.serial_number", "987654321fedcba987654321fedcba").
						CheckEqual("certificates.1.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("certificates.1.start_date", "2025-04-17T16:01:02.555444Z").
						CheckEqual("certificates.1.subject", "CN=Dummy CA test").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_two_certs.tf"),
					Check: baseCheck.
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("certificates.1.created_by", "johndoe").
						CheckEqual("certificates.1.created_date", "2025-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.end_date", "2026-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.fingerprint", "1234567890abcdef1234567890abcdef").
						CheckEqual("certificates.1.issuer", "CN=Dummy CA").
						CheckEqual("certificates.1.serial_number", "987654321fedcba987654321fedcba").
						CheckEqual("certificates.1.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("certificates.1.start_date", "2025-04-17T16:01:02.555444Z").
						CheckEqual("certificates.1.subject", "CN=Dummy CA test").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(0).AtMapKey("description"),
								knownvalue.StringExact("second cert")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test",
								tfjsonpath.New("certificates").AtSliceIndex(1).AtMapKey("description"),
								knownvalue.StringExact("Test certificate")),
						},
					},
				},
			},
		},
		"update an activated ca set on staging": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(3)

				// update.
				var one int64 = 1
				updateData := resourceData
				updateData.version = one
				updateData.stagingVersion = ptr.To(one)
				updateData.stagingStatus = "ACTIVE"
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.certificates = slices.Insert(updateData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				updateData.versionDescription = ptr.To("Second version for testing")
				mockValidateCertificates(m, updateData, nil).Times(5)
				updateData.newVersion = 2
				mockCloneCASetVersion(m, updateData).Times(1)
				updateData.version = 2
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update.tf"),
					Check: baseCheck.
						CheckEqual("version_description", "Second version for testing").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("certificates.1.created_by", "johndoe").
						CheckEqual("certificates.1.created_date", "2025-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.end_date", "2026-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.fingerprint", "1234567890abcdef1234567890abcdef").
						CheckEqual("certificates.1.issuer", "CN=Dummy CA").
						CheckEqual("certificates.1.serial_number", "987654321fedcba987654321fedcba").
						CheckEqual("certificates.1.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("certificates.1.start_date", "2025-04-17T16:01:02.555444Z").
						CheckEqual("certificates.1.subject", "CN=Dummy CA test").
						CheckEqual("latest_version", "2").
						CheckEqual("staging_version", "1").
						CheckMissing("production_version").
						CheckEqual("timeouts.delete", "5m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("5m")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("description"), knownvalue.StringExact("Test CA Set for validation")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_description"), knownvalue.StringExact("Second version for testing")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),
						},
					},
				},
			},
		},
		"update an activated ca set on staging, original certificates have expired (from create), but new one are correct": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(3)

				// update.
				var one int64 = 1
				updateData := resourceData
				updateData.version = one
				updateData.stagingVersion = ptr.To(one)
				updateData.stagingStatus = "ACTIVE"
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.certificates = []mtlstruststore.CertificateResponse{
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
						Description:        ptr.To("second cert"),
						CreatedBy:          "johndoe",
						CreatedDate:        tst.NewTimeFromStringMust("2025-04-16T16:01:02.555444Z"),
						EndDate:            tst.NewTimeFromStringMust("2026-04-16T16:01:02.555444Z"),
						Fingerprint:        "999567890abcdef1234567890abcdef",
						Issuer:             "CN=Dummy CA",
						SerialNumber:       "999654321fedcba987654321fedcba",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          tst.NewTimeFromStringMust("2025-04-17T16:01:02.555444Z"),
						Subject:            "CN=Dummy CA test",
					},
					{
						CertificatePEM:     "-----BEGIN CERTIFICATE-----\nUPDATED\n-----END CERTIFICATE-----\n",
						Description:        ptr.To("Test certificate"),
						CreatedBy:          "johndoe",
						CreatedDate:        tst.NewTimeFromStringMust("2025-03-16T16:01:02.555444Z"),
						EndDate:            tst.NewTimeFromStringMust("2026-03-16T16:01:02.555444Z"),
						Fingerprint:        "777567890abcdef1234567890abcdef",
						Issuer:             "CN=Dummy CA",
						SerialNumber:       "777654321fedcba987654321fedcba",
						SignatureAlgorithm: "SHA256WITHRSA",
						StartDate:          tst.NewTimeFromStringMust("2025-03-17T16:01:02.555444Z"),
						Subject:            "CN=Dummy CA test",
					},
				}
				updateData.versionDescription = ptr.To("Second version for testing")
				mockValidateCertificates(m, updateData, nil).Times(5)
				updateData.newVersion = 2
				updateData.validation.Warnings = []mtlstruststore.Warning{
					{
						ContextInfo: map[string]any{
							"checkDate":   "2025-08-01T09:34:53.000000Z",
							"description": nil,
							"expiryDate":  "2024-10-30T17:07:04.000000Z",
							"fingerprint": "1234567890abcdef1234567890abcdef",
							"subject":     "CN=cert100",
						},
						Detail: "The certificate with subject CN=cert100 and fingerprint 1234567890abcdef1234567890abcdef has expired. Expiry date is 2024-10-30T17:07:04.000000Z. The check was performed on 2025-08-01T09:34:53.000000Z.",
						Title:  "The certificate has expired.",
						Type:   "/mtls-edge-truststore/error-types/expired-certificate",
					},
				}

				mockCloneCASetVersion(m, updateData).Times(1)
				updateData.version = 2
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, updateData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_expired.tf"),
					Check: baseCheck.
						CheckEqual("version_description", "Second version for testing").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nUPDATED\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("certificates.1.created_by", "johndoe").
						CheckEqual("certificates.1.created_date", "2025-03-16T16:01:02.555444Z").
						CheckEqual("certificates.1.end_date", "2026-03-16T16:01:02.555444Z").
						CheckEqual("certificates.1.fingerprint", "777567890abcdef1234567890abcdef").
						CheckEqual("certificates.1.issuer", "CN=Dummy CA").
						CheckEqual("certificates.1.serial_number", "777654321fedcba987654321fedcba").
						CheckEqual("certificates.1.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("certificates.1.start_date", "2025-03-17T16:01:02.555444Z").
						CheckEqual("certificates.1.subject", "CN=Dummy CA test").
						CheckEqual("latest_version", "2").
						CheckEqual("staging_version", "1").
						CheckMissing("production_version").
						CheckEqual("timeouts.delete", "5m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("5m")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("description"), knownvalue.StringExact("Test CA Set for validation")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_description"), knownvalue.StringExact("Second version for testing")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),
						},
					},
				},
			},
		},
		"update an activated ca set on staging, changing only allow_insecure_sha1 should also create a new version": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(2)
				// update.
				var one int64 = 1
				updateData := resourceData
				updateData.version = one
				updateData.stagingVersion = ptr.To(one)
				updateData.stagingStatus = "ACTIVE"
				mockListCASetActivations(m, updateData, true).Times(1)
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.allowInsecureSHA1 = true
				mockValidateCertificates(m, updateData, nil).Times(5)
				updateData.newVersion = 2
				mockCloneCASetVersion(m, updateData).Times(1)
				updateData.version = 2
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_allow_insecure_sha1.tf"),
					Check: check.
						CheckEqual("allow_insecure_sha1", "true").
						CheckEqual("latest_version", "2").
						CheckEqual("staging_version", "1").
						CheckMissing("production_version").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("description"), knownvalue.StringExact("Test CA Set for validation")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_description"), knownvalue.StringExact("Initial version for testing")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by")),
							plancheck.ExpectUnknownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date")),
						},
					},
				},
			},
		},
		"update only timeout should do nothing, also when activated on staging (from default to some value)": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(2)

				// update.
				updateData := resourceData
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockValidateCertificates(m, updateData, nil).Times(4)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update_timeout.tf"),
					Check: check.
						CheckEqual("timeouts.delete", "5m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("5m")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("staging_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("production_version"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("latest_version"), knownvalue.Int64Exact(1)),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_by"), knownvalue.StringExact("someone")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_created_date"), knownvalue.StringExact("2025-04-16T12:08:34.099457Z")),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_by"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_mtlstruststore_ca_set.test", tfjsonpath.New("version_modified_date"), knownvalue.Null()),
						},
					},
				},
			},
		},
		"update an activated ca set on production": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(2)

				// update.
				var one int64 = 1
				updateData := resourceData
				updateData.version = one
				updateData.productionVersion = ptr.To(one)
				updateData.productionStatus = "ACTIVE"
				mockListCASetActivations(m, updateData, true).Times(1)
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				updateData.certificates = slices.Insert(updateData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				updateData.versionDescription = ptr.To("Second version for testing")
				mockValidateCertificates(m, updateData, nil).Times(5)
				updateData.newVersion = 2
				mockCloneCASetVersion(m, updateData).Times(1)
				updateData.version = 2
				mockUpdateCASetVersion(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, resourceData, true).Times(4)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update.tf"),
					Check: baseCheck.
						CheckEqual("version_description", "Second version for testing").
						CheckEqual("version_modified_by", "someone").
						CheckEqual("version_modified_date", "2025-04-18T12:18:34Z").
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("latest_version", "2").
						CheckMissing("staging_version").
						CheckEqual("production_version", "1").
						CheckEqual("timeouts.delete", "5m").
						Build(),
				},
			},
		},
		"update ca set which was removed outside Terraform should remove resource and create a new one": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// create.
				mockValidateCertificates(m, resourceData, nil).Times(5)
				mockCreateCASet(m, resourceData).Times(1)
				mockCreateCASetVersion(m, resourceData).Times(1)
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(2)

				// attempt to update, but the CA Set was removed outside Terraform.
				mockGetCASet(m, resourceData).Return(nil, mtlstruststore.ErrGetCASetNotFound).Times(1)

				// create a new resource, as the old one was removed outside Terraform.
				updateData := resourceData
				updateData.caSetID = "777"
				updateData.versionDescription = ptr.To("Second version for testing")
				updateData.certificates = slices.Insert(updateData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				mockValidateCertificates(m, updateData, nil).Times(5)
				mockCreateCASet(m, updateData).Times(1)
				mockCreateCASetVersion(m, updateData).Times(1)
				mockGetCASet(m, updateData).Times(1)
				// read.
				mockGetCASet(m, updateData).Times(1)
				mockGetCASetVersion(m, updateData).Times(1)
				mockListCASetActivations(m, updateData, false).Times(2)

				// delete.
				mockListCASetAssociations(m, updateData).Times(2)
				mockGetCASet(m, updateData).Once()
				mockDeleteCASet(m, updateData).Times(1)
				mockGetCASetDeletionStatus(m, updateData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, updateData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/create.tf"),
					Check:  check.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/update.tf"),
					Check: baseCheck.
						CheckEqual("id", "777").
						CheckEqual("version_description", "Second version for testing").
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("timeouts.delete", "5m").
						Build(),
				},
			},
		},
		"import": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// import.
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(2)
				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResCASet/import.tf"),
					ImportState:   true,
					ImportStateId: "123456789",
					ResourceName:  "akamai_mtlstruststore_ca_set.test",
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						assert.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "set-1", rs.Attributes["name"])
						assert.Equal(t, "123456789", rs.Attributes["id"])
						assert.Equal(t, "Test CA Set for validation", rs.Attributes["description"])
						assert.Equal(t, "false", rs.Attributes["allow_insecure_sha1"])
						assert.Equal(t, "Initial version for testing", rs.Attributes["version_description"])
						assert.Equal(t, "1", rs.Attributes["certificates.#"])
						assert.Equal(t, "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n", rs.Attributes["certificates.0.certificate_pem"])
						assert.Equal(t, "Test certificate", rs.Attributes["certificates.0.description"])
						assert.Equal(t, "1", rs.Attributes["latest_version"])
						return nil
					},
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
				},
			},
		},
		"import with certs order change": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// import.
				resourceData.certificates = slices.Insert(resourceData.certificates, 0, mtlstruststore.CertificateResponse{
					CertificatePEM: "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n",
					Description:    ptr.To("second cert"),
				})
				mockGetCASet(m, resourceData).Times(1)
				// read.
				mockGetCASet(m, resourceData).Times(1)
				mockGetCASetVersion(m, resourceData).Times(1)
				mockListCASetActivations(m, resourceData, false).Times(3)
				// update.
				mockGetCASet(m, resourceData).Times(2)
				mockGetCASetVersion(m, resourceData).Times(3)
				mockValidateCertificates(m, resourceData, nil).Times(3)

				// delete.
				mockListCASetAssociations(m, resourceData).Times(2)
				mockGetCASet(m, resourceData).Once()
				mockDeleteCASet(m, resourceData).Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "IN_PROGRESS", "IN_PROGRESS", "COMPLETED").Times(1)
				mockGetCASetDeletionStatus(m, resourceData, "COMPLETE", "COMPLETE", "COMPLETE").Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResCASet/import_order_change.tf"),
					ImportState:   true,
					ImportStateId: "123456789",
					ResourceName:  "akamai_mtlstruststore_ca_set.test",
					ImportStateCheck: func(s []*terraform.InstanceState) error {
						assert.Len(t, s, 1)
						rs := s[0]
						assert.Equal(t, "set-1", rs.Attributes["name"])
						assert.Equal(t, "123456789", rs.Attributes["id"])
						assert.Equal(t, "Test CA Set for validation", rs.Attributes["description"])
						assert.Equal(t, "false", rs.Attributes["allow_insecure_sha1"])
						assert.Equal(t, "Initial version for testing", rs.Attributes["version_description"])
						assert.Equal(t, "2", rs.Attributes["certificates.#"])
						assert.Equal(t, "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n", rs.Attributes["certificates.0.certificate_pem"])
						assert.Equal(t, "second cert", rs.Attributes["certificates.0.description"])
						assert.Equal(t, "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n", rs.Attributes["certificates.1.certificate_pem"])
						assert.Equal(t, "Test certificate", rs.Attributes["certificates.1.description"])
						assert.Equal(t, "1", rs.Attributes["latest_version"])
						return nil
					},
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCASet/import_order_change.tf"),
					Check: baseCheck.
						CheckEqual("certificates.#", "2").
						CheckEqual("certificates.0.certificate_pem", "-----BEGIN CERTIFICATE-----\nFOO\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.0.description", "second cert").
						CheckEqual("certificates.1.certificate_pem", "-----BEGIN CERTIFICATE-----\nMIIDXTCCAkWgAwIBAgIJALa6Rz1u5z2OMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV\n-----END CERTIFICATE-----\n").
						CheckEqual("certificates.1.description", "Test certificate").
						CheckEqual("certificates.1.created_by", "johndoe").
						CheckEqual("certificates.1.created_date", "2025-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.end_date", "2026-04-16T16:01:02.555444Z").
						CheckEqual("certificates.1.fingerprint", "1234567890abcdef1234567890abcdef").
						CheckEqual("certificates.1.issuer", "CN=Dummy CA").
						CheckEqual("certificates.1.serial_number", "987654321fedcba987654321fedcba").
						CheckEqual("certificates.1.signature_algorithm", "SHA256WITHRSA").
						CheckEqual("certificates.1.start_date", "2025-04-17T16:01:02.555444Z").
						CheckEqual("certificates.1.subject", "CN=Dummy CA test").
						Build(),
				},
			},
		},
		"expect error - import ca set without version": {
			init: func(m *mtlstruststore.Mock, resourceData commonDataForResource) {
				// import.
				resourceData.version = 0
				mockGetCASet(m, resourceData).Times(1)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASet/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "123456789",
					ResourceName:                         "akamai_mtlstruststore_ca_set.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("The CA set does not have any version"),
				},
			},
		},
		"expect error - unknown id in import": {
			init: func(m *mtlstruststore.Mock, _ commonDataForResource) {
				m.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
					CASetID: "9999",
				}).
					Return(nil, mtlstruststore.ErrGetCASetNotFound)
			},
			mockData: createData,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResCASet/import.tf"),
					ImportState:                          true,
					ImportStateId:                        "9999",
					ResourceName:                         "akamai_mtlstruststore_ca_set.test",
					ImportStateVerifyIdentifierAttribute: "id",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("ca set not found"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mtlstruststore.Mock{}
			if tc.init != nil {
				tc.init(client, tc.mockData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ExternalProviders: map[string]resource.ExternalProvider{
						"random": {
							Source:            "registry.terraform.io/hashicorp/random",
							VersionConstraint: "3.1.0",
						},
					},
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider(), testprovider.NewMockSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockValidateCertificates(client *mtlstruststore.Mock, testData commonDataForResource, err error) *mock.Call {
	var certificates []mtlstruststore.ValidateCertificate
	for _, c := range testData.certificates {
		certificates = append(certificates, mtlstruststore.ValidateCertificate{
			CertificatePEM: c.CertificatePEM,
			Description:    c.Description,
		})
	}

	var certificatesResponse []mtlstruststore.ValidateCertificateResponse
	for _, c := range testData.certificates {
		certificatesResponse = append(certificatesResponse, mtlstruststore.ValidateCertificateResponse{
			// Certificates are trimmed in the API, so we do it here too.
			CertificatePEM: text.TrimRightWhitespace(c.CertificatePEM),
			Description:    c.Description,
		})
	}
	if err != nil {
		return client.On("ValidateCertificates", testutils.MockContext, mtlstruststore.ValidateCertificatesRequest{
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			Certificates:      certificates}).
			Return(nil, err).Once()
	}
	return client.On("ValidateCertificates", testutils.MockContext, mtlstruststore.ValidateCertificatesRequest{
		AllowInsecureSHA1: testData.allowInsecureSHA1,
		Certificates:      certificates}).
		Return(&mtlstruststore.ValidateCertificatesResponse{
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			Certificates:      certificatesResponse,
			Validation:        testData.validation,
		}, nil).Once()
}

func mockCreateCASet(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	return client.On("CreateCASet", testutils.MockContext, mtlstruststore.CreateCASetRequest{
		CASetName:   testData.name,
		Description: testData.description,
	}).
		Return(&mtlstruststore.CreateCASetResponse{
			CASetID:               testData.caSetID,
			CASetName:             testData.name,
			AccountID:             "ACC-123456",
			CASetLink:             "",
			CASetStatus:           "NOT_DELETED",
			Description:           testData.description,
			LatestVersionLink:     nil,
			LatestVersion:         nil,
			StagingVersionLink:    nil,
			StagingVersion:        nil,
			ProductionVersionLink: nil,
			ProductionVersion:     nil,
			VersionsLink:          "",
			CreatedDate:           tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
			CreatedBy:             "someone",
			DeletedDate:           nil,
			DeletedBy:             nil,
		}, nil)
}

func mockUpdateCASetVersion(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	modifiedDate := tst.NewTimeFromStringMust("2025-04-18T12:18:34Z")
	var certificateRequests []mtlstruststore.CertificateRequest
	for _, c := range testData.certificates {
		certificateRequests = append(certificateRequests, mtlstruststore.CertificateRequest{
			CertificatePEM: c.CertificatePEM,
			Description:    c.Description,
		})
	}
	var certificateResponse []mtlstruststore.CertificateResponse
	for _, c := range testData.certificates {
		certificateResponse = append(certificateResponse, mtlstruststore.CertificateResponse{
			// Certificates are trimmed in the API, so we do it here too.
			CertificatePEM:     text.TrimRightWhitespace(c.CertificatePEM),
			Description:        c.Description,
			CreatedBy:          c.CreatedBy,
			CreatedDate:        c.CreatedDate,
			EndDate:            c.EndDate,
			Fingerprint:        c.Fingerprint,
			Issuer:             c.Issuer,
			SerialNumber:       c.SerialNumber,
			SignatureAlgorithm: c.SignatureAlgorithm,
			StartDate:          c.StartDate,
			Subject:            c.Subject,
		})
	}
	return client.On("UpdateCASetVersion", testutils.MockContext, mtlstruststore.UpdateCASetVersionRequest{
		CASetID: testData.caSetID,
		Version: testData.version,
		Body: mtlstruststore.UpdateCASetVersionRequestBody{
			Description:       testData.versionDescription,
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			Certificates:      certificateRequests,
		},
	}).
		Return(&mtlstruststore.UpdateCASetVersionResponse{
			CASetID:           testData.caSetID,
			Version:           testData.version,
			CASetName:         testData.name,
			VersionLink:       "",
			Description:       testData.versionDescription,
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			StagingStatus:     "",
			ProductionStatus:  "",
			CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
			CreatedBy:         "someone",
			ModifiedDate:      &modifiedDate,
			ModifiedBy:        ptr.To("someone"),
			Certificates:      certificateResponse,
		}, nil)
}

func mockCloneCASetVersion(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	var certificateResponse []mtlstruststore.CertificateResponse
	for _, c := range testData.certificates {
		certificateResponse = append(certificateResponse, mtlstruststore.CertificateResponse{
			// Certificates are trimmed in the API, so we do it here too.
			CertificatePEM: text.TrimRightWhitespace(c.CertificatePEM),
			Description:    c.Description,
		})
	}
	modifiedDate := tst.NewTimeFromStringMust("2025-04-18T12:18:34Z")
	clonedDescription := fmt.Sprintf("This CA set version is cloned from version %d", testData.version)
	return client.On("CloneCASetVersion", testutils.MockContext, mtlstruststore.CloneCASetVersionRequest{
		CASetID: testData.caSetID,
		Version: testData.version,
	}).
		Return(&mtlstruststore.CloneCASetVersionResponse{
			CASetID:           testData.caSetID,
			Version:           testData.newVersion,
			CASetName:         testData.name,
			VersionLink:       "",
			Description:       ptr.To(clonedDescription),
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			StagingStatus:     testData.stagingStatus,
			ProductionStatus:  testData.productionStatus,
			CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
			CreatedBy:         "someone",
			ModifiedDate:      &modifiedDate,
			ModifiedBy:        ptr.To("someone"),
			Certificates:      certificateResponse,
		}, nil)
}

func mockGetCASet(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	var version *int64
	if testData.version > 0 {
		version = ptr.To(testData.version)
	}
	return client.On("GetCASet", testutils.MockContext, mtlstruststore.GetCASetRequest{
		CASetID: testData.caSetID,
	}).
		Return(&mtlstruststore.GetCASetResponse{
			CASetID:               testData.caSetID,
			CASetName:             testData.name,
			AccountID:             "ACC-123456",
			CASetLink:             "",
			CASetStatus:           testData.caSetStatus,
			Description:           testData.description,
			LatestVersionLink:     nil,
			LatestVersion:         version,
			StagingVersionLink:    nil,
			StagingVersion:        testData.stagingVersion,
			ProductionVersionLink: nil,
			ProductionVersion:     testData.productionVersion,
			VersionsLink:          "",
			CreatedDate:           tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
			CreatedBy:             "someone",
			DeletedDate:           nil,
			DeletedBy:             nil,
		}, nil)
}

func mockCreateCASetVersion(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	var certificateRequests []mtlstruststore.CertificateRequest
	for _, c := range testData.certificates {
		certificateRequests = append(certificateRequests, mtlstruststore.CertificateRequest{
			CertificatePEM: c.CertificatePEM,
			Description:    c.Description,
		})
	}
	var certificates []mtlstruststore.CertificateResponse
	if testData.certificatesForResponse != nil {
		certificates = testData.certificatesForResponse
	} else {
		certificates = testData.certificates
	}
	var certificateResponse []mtlstruststore.CertificateResponse
	for _, c := range certificates {
		certificateResponse = append(certificateResponse, mtlstruststore.CertificateResponse{
			// Certificates are trimmed in the API, so we do it here too.
			CertificatePEM:     text.TrimRightWhitespace(c.CertificatePEM),
			Description:        c.Description,
			CreatedBy:          c.CreatedBy,
			CreatedDate:        c.CreatedDate,
			EndDate:            c.EndDate,
			Fingerprint:        c.Fingerprint,
			Issuer:             c.Issuer,
			SerialNumber:       c.SerialNumber,
			SignatureAlgorithm: c.SignatureAlgorithm,
			StartDate:          c.StartDate,
			Subject:            c.Subject,
		})
	}
	return client.On("CreateCASetVersion", testutils.MockContext, mtlstruststore.CreateCASetVersionRequest{
		CASetID: testData.caSetID,
		Body: mtlstruststore.CreateCASetVersionRequestBody{
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			Description:       testData.versionDescription,
			Certificates:      certificateRequests,
		},
	}).
		Return(&mtlstruststore.CreateCASetVersionResponse{
			CASetID:           testData.caSetID,
			Version:           testData.version,
			CASetName:         testData.name,
			VersionLink:       "",
			Description:       testData.versionDescription,
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			StagingStatus:     testData.stagingStatus,
			ProductionStatus:  "INACTIVE",
			CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
			CreatedBy:         "someone",
			ModifiedDate:      nil,
			ModifiedBy:        nil,
			Certificates:      certificateResponse,
			Validation:        &testData.validation,
		}, nil)
}

func mockGetCASetVersion(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	var certificateResponse []mtlstruststore.CertificateResponse
	for _, c := range testData.certificates {
		certificateResponse = append(certificateResponse, mtlstruststore.CertificateResponse{
			// Certificates are trimmed in the API, so we do it here too.
			CertificatePEM:     text.TrimRightWhitespace(c.CertificatePEM),
			Description:        c.Description,
			CreatedBy:          c.CreatedBy,
			CreatedDate:        c.CreatedDate,
			EndDate:            c.EndDate,
			Fingerprint:        c.Fingerprint,
			Issuer:             c.Issuer,
			SerialNumber:       c.SerialNumber,
			SignatureAlgorithm: c.SignatureAlgorithm,
			StartDate:          c.StartDate,
			Subject:            c.Subject,
		})
	}
	return client.On("GetCASetVersion", testutils.MockContext, mtlstruststore.GetCASetVersionRequest{
		CASetID: testData.caSetID,
		Version: testData.version,
	}).
		Return(&mtlstruststore.GetCASetVersionResponse{
			CASetID:           testData.caSetID,
			Version:           testData.version,
			CASetName:         testData.name,
			VersionLink:       "",
			Description:       testData.versionDescription,
			AllowInsecureSHA1: testData.allowInsecureSHA1,
			StagingStatus:     "INACTIVE",
			ProductionStatus:  "INACTIVE",
			CreatedDate:       tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
			CreatedBy:         "someone",
			ModifiedDate:      nil,
			ModifiedBy:        nil,
			Certificates:      certificateResponse,
		}, nil)
}

func mockListCASetActivations(client *mtlstruststore.Mock, testData commonDataForResource, activated bool) *mock.Call {
	var activations []mtlstruststore.ActivateCASetVersionResponse
	var network string
	if testData.stagingVersion != nil {
		network = "STAGING"
	} else {
		network = "PRODUCTION"
	}
	if activated {
		activations = []mtlstruststore.ActivateCASetVersionResponse{
			{
				ActivationID:     10,
				ActivationLink:   "",
				CASetID:          testData.caSetID,
				CASetName:        testData.name,
				CASetLink:        "",
				CreatedBy:        "someone",
				CreatedDate:      tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z"),
				ModifiedBy:       nil,
				ModifiedDate:     nil,
				Network:          network,
				ActivationStatus: "",
				ActivationType:   "",
				Version:          testData.version,
				VersionLink:      "",
			},
		}
	} else {
		activations = []mtlstruststore.ActivateCASetVersionResponse{}
	}
	return client.On("ListCASetActivations", testutils.MockContext, mtlstruststore.ListCASetActivationsRequest{
		CASetID: testData.caSetID,
	}).
		Return(&mtlstruststore.ListCASetActivationsResponse{
			Activations: activations,
		}, nil)
}

func mockListCASetAssociations(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	return client.On("ListCASetAssociations", testutils.MockContext, mtlstruststore.ListCASetAssociationsRequest{
		CASetID: testData.caSetID,
	}).
		Return(&mtlstruststore.ListCASetAssociationsResponse{
			Associations: mtlstruststore.Associations{
				Properties:  testData.properties,
				Enrollments: testData.enrollments,
			},
		}, nil)
}

func mockDeleteCASet(client *mtlstruststore.Mock, testData commonDataForResource) *mock.Call {
	return client.On("DeleteCASet", testutils.MockContext, mtlstruststore.DeleteCASetRequest{
		CASetID: testData.caSetID,
	}).Return(nil)
}

func mockGetCASetDeletionStatus(client *mtlstruststore.Mock, testData commonDataForResource,
	status, stagingStatus, productionStatus string) *mock.Call {
	startTime := tst.NewTimeFromStringMust("2025-04-16T12:06:00Z")
	endTime := tst.NewTimeFromStringMust("2025-04-16T12:18:34Z")
	estimatedEndTime := tst.NewTimeFromStringMust("2025-04-16T12:08:34.099457Z")
	response := mtlstruststore.GetCASetDeletionStatusResponse{
		Status:           status,
		StatusLink:       "",
		CASetLink:        "",
		ResourceMethod:   nil,
		CASetID:          testData.caSetID,
		CASetName:        testData.name,
		EstimatedEndTime: &estimatedEndTime,
		EndTime:          &endTime,
		StartTime:        startTime,
		Deletions: []mtlstruststore.CASetNetworkDeleteStatus{
			{
				Network: "STAGING",
				Status:  stagingStatus,
			},
			{
				Network: "PRODUCTION",
				Status:  productionStatus,
			},
		},
	}
	if status == "FAILED" {
		response.FailureReason = ptr.To("Indication about deletion failure in a network")
	}
	if status == "IN_PROGRESS" {
		response.RetryAfter = time.Now().Add(1 * time.Millisecond)
	} else {
		response.RetryAfter = time.Time{}
	}
	setStatus(response, stagingStatus, 0)
	setStatus(response, productionStatus, 1)
	return client.On("GetCASetDeletionStatus", testutils.MockContext, mtlstruststore.GetCASetDeletionStatusRequest{
		CASetID: testData.caSetID,
	}).Return(&response, nil)
}

func setStatus(response mtlstruststore.GetCASetDeletionStatusResponse, status string, index int) {
	switch status {
	case "IN_PROGRESS":
		response.Deletions[index].FailureReason = nil
		response.Deletions[index].PercentComplete = 50
	case "FAILED":
		response.Deletions[index].FailureReason = ptr.To("Indication about deletion failure in a network")
		response.Deletions[index].PercentComplete = 0
	case "COMPLETED":
		response.Deletions[index].FailureReason = nil
		response.Deletions[index].PercentComplete = 100
	}
}
