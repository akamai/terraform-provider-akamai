package mtlskeystore

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClientCertificateUpload(t *testing.T) {
	t.Parallel()

	commonStateChecker := test.NewStateChecker("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload").
		CheckEqual("client_certificate_id", "12345").
		CheckEqual("version_number", "1").
		CheckEqual("signed_certificate", "certificate-data").
		CheckEqual("trust_chain", "trustchain-data").
		CheckEqual("version_guid", "unique-guid").
		CheckEqual("wait_for_deployment", "true")

	tests := map[string]struct {
		init  func(*mtlskeystore.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificate(m, nil)
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}

							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"happy path - without trust chain": {
			init: func(m *mtlskeystore.Mock) {
				func() {
					m.On("UploadSignedClientCertificate",
						testutils.MockContext,
						mtlskeystore.UploadSignedClientCertificateRequest{
							CertificateID: 12345,
							Version:       1,
							Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
								Certificate: "certificate-data",
							},
						}).Return(nil).Once()

					m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
						CertificateID: 12345,
					}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
						Versions: []mtlskeystore.ClientCertificateVersion{{Version: 1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed},
						},
					}, nil).Twice()
				}()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}

							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
							}`,
					Check: test.NewStateChecker("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload").
						CheckEqual("client_certificate_id", "12345").
						CheckEqual("version_number", "1").
						CheckEqual("signed_certificate", "certificate-data").
						CheckMissing("trust_chain").
						CheckEqual("version_guid", "unique-guid").Build(),
				},
			},
		},
		"happy path with polling": {
			init: func(m *mtlskeystore.Mock) {
				// Upload returns nil (success)
				m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
					CertificateID: 12345,
					Version:       1,
					Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
						Certificate: "certificate-data",
						TrustChain:  ptr.To("trustchain-data"),
					},
				}).Return(nil).Once()
				// First poll returns pending
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.DeploymentPending,
						},
					},
				}, nil).Once()
				// Second poll returns deployed
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
				// Read
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}

							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					Check: commonStateChecker.Build(),
				},
			},
		},
		"happy path without polling": {
			init: func(m *mtlskeystore.Mock) {
				m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
					CertificateID: 12345,
					Version:       1,
					Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
						Certificate: "certificate-data",
						TrustChain:  ptr.To("trustchain-data"),
					},
				}).Return(nil).Once()
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.DeploymentPending,
						},
					},
				}, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = false
							}`,
					Check: test.NewStateChecker("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload").
						CheckEqual("client_certificate_id", "12345").
						CheckEqual("version_number", "1").
						CheckEqual("signed_certificate", "certificate-data").
						CheckEqual("trust_chain", "trustchain-data").
						CheckEqual("version_guid", "unique-guid").
						CheckEqual("wait_for_deployment", "false").Build(),
				},
			},
		},
		"error API response": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificate(m, fmt.Errorf("API error"))
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}

							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("API error"),
				},
			},
		},
		"error - timeout during deployment": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificateWithTimeout(m)
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}

							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("timeout waiting for client certificate deployment:"),
				},
			},
		},
		"error - unexpected status": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificateUnexpectedStatus(m)
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("unexpected client certificate version status"),
				},
			},
		},
		"error - API error during polling": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificatePollingError(m)
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("error retrieving client certificate versions: API error"),
				},
			},
		},
		"error - no version returned during polling": {
			init: func(m *mtlskeystore.Mock) {
				m.On("UploadSignedClientCertificate",
					testutils.MockContext,
					mtlskeystore.UploadSignedClientCertificateRequest{
						CertificateID: 12345,
						Version:       1,
						Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
							Certificate: "certificate-data",
							TrustChain:  ptr.To("trustchain-data"),
						},
					},
				).Return(nil).Once()

				// Return a pending version on the first call and first poll
				m.On("GetClientCertificateVersions",
					testutils.MockContext,
					mtlskeystore.GetClientCertificateVersionsRequest{
						CertificateID: 12345,
					},
				).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.DeploymentPending,
						},
					},
				}, nil).Twice()
				// On the next call, return no versions
				m.On("GetClientCertificateVersions",
					testutils.MockContext,
					mtlskeystore.GetClientCertificateVersionsRequest{
						CertificateID: 12345,
					},
				).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("no client certificate versions found"),
				},
			},
		},
		"update - happy path with new version": {
			init: func(m *mtlskeystore.Mock) {
				// Initial create
				mockUploadSignedClientCertificate(m, nil)
				// Read before update
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
				// Update with new version
				m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
					CertificateID: 12345,
					Version:       2,
					Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
						Certificate: "certificate-data-updated",
						TrustChain:  ptr.To("trustchain-data"),
					},
				}).Return(nil).Once()

				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     2,
							VersionGUID: "unique-guid-updated",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload", "version_number", "1"),
					),
				},
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 2
								signed_certificate    = "certificate-data-updated"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload", "version_number", "2"),
						resource.TestCheckResourceAttr("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload", "version_guid", "unique-guid-updated"),
					),
				},
			},
		},
		"update - happy path with polling": {
			init: func(m *mtlskeystore.Mock) {
				// Initial create
				mockUploadSignedClientCertificate(m, nil)
				// Read before update
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
				// Update with new version
				m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
					CertificateID: 12345,
					Version:       2,
					Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
						Certificate: "certificate-data-updated",
						TrustChain:  ptr.To("trustchain-data"),
					},
				}).Return(nil).Once()
				// Polling: first call returns pending, second returns deployed
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     2,
							VersionGUID: "unique-guid-updated",
							Status:      mtlskeystore.DeploymentPending,
						},
					},
				}, nil).Once()
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     2,
							VersionGUID: "unique-guid-updated",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()

				// Final read
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     2,
							VersionGUID: "unique-guid-updated",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload", "version_number", "1"),
					),
				},
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 2
								signed_certificate    = "certificate-data-updated"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload", "version_number", "2"),
						resource.TestCheckResourceAttr("akamai_mtlskeystore_client_certificate_upload.client_certificate_upload", "version_guid", "unique-guid-updated"),
					),
				},
			},
		},
		"error update - same version number": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificate(m, nil)
				// Read before update
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
				},
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data-updated"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("Only updates with a different version_number are supported"),
				},
			},
		},
		"error update - different certificate id": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificate(m, nil)
				// Read before update
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
				},
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 99999
								version_number        = 2
								signed_certificate    = "certificate-data-updated"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("updating field `client_certificate_id` is not possible"),
				},
			},
		},
		"update - API error on upload": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificate(m, nil)
				// Read before update
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
				m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
					CertificateID: 12345,
					Version:       2,
					Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
						Certificate: "certificate-data-updated",
						TrustChain:  ptr.To("trustchain-data"),
					},
				}).Return(fmt.Errorf("update upload error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
				},
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 2
								signed_certificate    = "certificate-data-updated"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("update upload error"),
				},
			},
		},
		"update - API error on version retrieval": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificate(m, nil)
				// Read before update
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()
				m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
					CertificateID: 12345,
					Version:       2,
					Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
						Certificate: "certificate-data-updated",
						TrustChain:  ptr.To("trustchain-data"),
					},
				}).Return(nil).Once()
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(nil, fmt.Errorf("update get versions error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
				},
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 2
								signed_certificate    = "certificate-data-updated"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("update get versions error"),
				},
			},
		},
		"update - uploaded version not found": {
			init: func(m *mtlskeystore.Mock) {
				mockUploadSignedClientCertificate(m, nil)
				// Read before update
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{
							Version:     1,
							VersionGUID: "unique-guid",
							Status:      mtlskeystore.Deployed,
						},
					},
				}, nil).Once()

				m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
					CertificateID: 12345,
					Version:       2,
					Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
						Certificate: "certificate-data-updated",
						TrustChain:  ptr.To("trustchain-data"),
					},
				}).Return(nil).Once()
				m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{Version: 3, Status: mtlskeystore.Deployed, VersionGUID: "other-guid"},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 1
								signed_certificate    = "certificate-data"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
				},
				{
					Config: `provider "akamai" {
								edgerc = "../../common/testutils/edgerc"
							}
							resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
								client_certificate_id = 12345
								version_number        = 2
								signed_certificate    = "certificate-data-updated"
								trust_chain           = "trustchain-data"
								wait_for_deployment   = true
							}`,
					ExpectError: regexp.MustCompile("uploaded version not found"),
				},
			},
		},
	}
	pollingInterval = 1 * time.Millisecond
	numberOfRetries = 5
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

func mockUploadSignedClientCertificate(m *mtlskeystore.Mock, err error) {
	if err != nil {
		m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
			CertificateID: 12345,
			Version:       1,
			Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
				Certificate: "certificate-data",
				TrustChain:  ptr.To("trustchain-data"),
			},
		}).Return(err).Once()
		return
	}
	m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
		CertificateID: 12345,
		Version:       1,
		Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
			Certificate: "certificate-data",
			TrustChain:  ptr.To("trustchain-data"),
		},
	}).Return(nil).Once()

	m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
		CertificateID: 12345,
	}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				VersionGUID: "unique-guid",
				Status:      mtlskeystore.Deployed,
			},
		},
	}, nil).Twice()
}

func mockUploadSignedClientCertificateWithTimeout(m *mtlskeystore.Mock) {
	m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
		CertificateID: 12345,
		Version:       1,
		Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
			Certificate: "certificate-data",
			TrustChain:  ptr.To("trustchain-data"),
		},
	}).Return(nil).Once()

	// Simulate polling timeout
	m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
		CertificateID: 12345,
	}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				Status:      mtlskeystore.DeploymentPending,
				VersionGUID: "unique-guid",
			},
		},
	}, nil).Times(numberOfRetries + 1)
}

func mockUploadSignedClientCertificateUnexpectedStatus(m *mtlskeystore.Mock) {
	m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
		CertificateID: 12345,
		Version:       1,
		Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
			Certificate: "certificate-data",
			TrustChain:  ptr.To("trustchain-data"),
		},
	}).Return(nil).Once()
	m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
		CertificateID: 12345,
	}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{Version: 1, Status: "INVALID_STATUS", VersionGUID: "unique-guid"},
		},
	}, nil).Once()
}

func mockUploadSignedClientCertificatePollingError(m *mtlskeystore.Mock) {
	m.On("UploadSignedClientCertificate", testutils.MockContext, mtlskeystore.UploadSignedClientCertificateRequest{
		CertificateID: 12345,
		Version:       1,
		Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
			Certificate: "certificate-data",
			TrustChain:  ptr.To("trustchain-data"),
		},
	}).Return(nil).Once()
	m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
		CertificateID: 12345,
	}).Return(&mtlskeystore.GetClientCertificateVersionsResponse{
		Versions: []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				VersionGUID: "unique-guid",
				Status:      mtlskeystore.DeploymentPending,
			},
		},
	}, nil).Twice()
	m.On("GetClientCertificateVersions", testutils.MockContext, mtlskeystore.GetClientCertificateVersionsRequest{
		CertificateID: 12345,
	}).Return(nil, fmt.Errorf("API error during polling")).Once()
}
