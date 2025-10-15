package mtlskeystore

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlskeystore"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestClientCertificateUpload(t *testing.T) {
	t.Parallel()

	getDefaultUploadCertRequest := func() mtlskeystore.UploadSignedClientCertificateRequest {
		return mtlskeystore.UploadSignedClientCertificateRequest{
			CertificateID: 12345,
			Version:       1,
			Body: mtlskeystore.UploadSignedClientCertificateRequestBody{
				Certificate: "certificate-data",
				TrustChain:  ptr.To("trustchain-data"),
			},
			AcknowledgeAllWarnings: ptr.To(false),
		}
	}

	getDefaultReturnedVersions := func() []mtlskeystore.ClientCertificateVersion {
		return []mtlskeystore.ClientCertificateVersion{
			{
				Version:     1,
				VersionGUID: "unique-guid",
				Status:      string(mtlskeystore.CertificateVersionStatusDeployed),
			},
		}
	}

	getCustomReturnedVersions := func(version int64, versionGUID string, status mtlskeystore.CertificateVersionStatus) []mtlskeystore.ClientCertificateVersion {
		return []mtlskeystore.ClientCertificateVersion{
			{
				Version:     version,
				VersionGUID: versionGUID,
				Status:      string(status),
			},
		}
	}

	mockUploadSignedClientCertificate := func(m *mtlskeystore.Mock, request mtlskeystore.UploadSignedClientCertificateRequest, err error) {
		m.On("UploadSignedClientCertificate", testutils.MockContext, request).Return(err).Once()
	}

	mockListClientCertificateVersions := func(m *mtlskeystore.Mock, versions []mtlskeystore.ClientCertificateVersion, err error) *mock.Call {
		if err != nil {
			return m.On("ListClientCertificateVersions", testutils.MockContext, mock.Anything).Return(nil, err).Once()
		}
		return m.On("ListClientCertificateVersions", testutils.MockContext, mtlskeystore.ListClientCertificateVersionsRequest{
			CertificateID: 12345,
		}).Return(&mtlskeystore.ListClientCertificateVersionsResponse{
			Versions: versions,
		}, nil)
	}

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
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path with cert provided externally by another resource": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/external_cert.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path with auto_acknowledge_warnings": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				request := getDefaultUploadCertRequest()
				request.AcknowledgeAllWarnings = ptr.To(true)
				mockUploadSignedClientCertificate(m, request, nil)

				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/with_autoacknowledge.tf"),
					Check:  commonStateChecker.CheckEqual("auto_acknowledge_warnings", "true").Build(),
				},
			},
		},
		"happy path - without trust chain": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				request := getDefaultUploadCertRequest()
				request.Body.TrustChain = nil
				mockUploadSignedClientCertificate(m, request, nil)

				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/no_trust_chain.tf"),
					Check: commonStateChecker.
						CheckMissing("trust_chain").Build(),
				},
			},
		},
		"happy path with polling": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				// Upload returns nil (success)
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				// First poll returns pending
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusDeploymentPending), nil).Once()
				// Second poll returns deployed
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
				// Read
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"happy path without polling": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/no_wait_for_deployment.tf"),
					Check: commonStateChecker.
						CheckEqual("wait_for_deployment", "false").Build(),
				},
			},
		},
		"version dropped outside of terraform": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
				// On second call, version 1 is not returned
				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusDeployed), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:              commonStateChecker.Build(),
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"certificate dropped outside of terraform": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
				// On second call, no cert is found
				mockListClientCertificateVersions(m, nil, mtlskeystore.ErrClientCertificateNotFound).Once()
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:              commonStateChecker.Build(),
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"error API response from UploadSignedClientCertificate": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), fmt.Errorf("API error"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					ExpectError: regexp.MustCompile("API error"),
				},
			},
		},
		"error - timeout during deployment": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusDeploymentPending), nil) //Times(infinity)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					ExpectError: regexp.MustCompile("timeout waiting for client certificate deployment:"),
				},
			},
		},
		"error - unexpected status": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				m.On("ListClientCertificateVersions", testutils.MockContext, mtlskeystore.ListClientCertificateVersionsRequest{
					CertificateID: 12345,
				}).Return(&mtlskeystore.ListClientCertificateVersionsResponse{
					Versions: []mtlskeystore.ClientCertificateVersion{
						{Version: 1, Status: "INVALID_STATUS", VersionGUID: "unique-guid"},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					ExpectError: regexp.MustCompile("unexpected client certificate version status"),
				},
			},
		},
		"error - API error during polling": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusDeploymentPending), nil).Twice()
				mockListClientCertificateVersions(m, nil, fmt.Errorf("API error during polling"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					ExpectError: regexp.MustCompile("error retrieving client certificate versions: API error"),
				},
			},
		},
		"error - no version returned during polling": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)

				// Return a pending version on the first call and first poll
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusDeploymentPending), nil).Twice()
				// On the next call, return no versions
				mockListClientCertificateVersions(m, []mtlskeystore.ClientCertificateVersion{}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					ExpectError: regexp.MustCompile("no client certificate versions found"),
				},
			},
		},
		"update - happy path with new version": {
			init: func(m *mtlskeystore.Mock) {
				// Initial create
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()

				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
				// Update with new version
				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				updateRequest := getDefaultUploadCertRequest()
				updateRequest.Body.Certificate = "certificate-data-updated"
				updateRequest.Version = 2
				mockUploadSignedClientCertificate(m, updateRequest, nil)

				updatedVersions := getDefaultReturnedVersions()
				updatedVersions[0].Version = 2
				updatedVersions[0].VersionGUID = "unique-guid-updated"
				mockListClientCertificateVersions(m, updatedVersions, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_2.tf"),
					Check: commonStateChecker.
						CheckEqual("version_number", "2").
						CheckEqual("version_guid", "unique-guid-updated").
						CheckEqual("signed_certificate", "certificate-data-updated").Build(),
				},
			},
		},
		"update - happy path with polling": {
			init: func(m *mtlskeystore.Mock) {
				// Initial create
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
				// Update with new version
				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				updateRequest := getDefaultUploadCertRequest()
				updateRequest.Version = 2
				updateRequest.Body.Certificate = "certificate-data-updated"
				mockUploadSignedClientCertificate(m, updateRequest, nil)
				// Polling: first call returns pending, second returns deployed
				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusDeploymentPending), nil).Once()
				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusDeployed), nil).Once()

				// Final read
				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusDeployed), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_2.tf"),
					Check: commonStateChecker.
						CheckEqual("version_number", "2").
						CheckEqual("version_guid", "unique-guid-updated").
						CheckEqual("signed_certificate", "certificate-data-updated").Build(),
				},
			},
		},
		"update - happy path with custom timeouts": {
			init: func(m *mtlskeystore.Mock) {
				// Initial create
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()

				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
				// Update with new version
				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				updateRequest := getDefaultUploadCertRequest()
				updateRequest.Body.Certificate = "certificate-data-updated"
				updateRequest.Version = 2
				mockUploadSignedClientCertificate(m, updateRequest, nil)

				updatedVersions := getDefaultReturnedVersions()
				updatedVersions[0].Version = 2
				updatedVersions[0].VersionGUID = "unique-guid-updated"
				mockListClientCertificateVersions(m, updatedVersions, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/create_timeout.tf"),
					Check: commonStateChecker.
						CheckEqual("version_number", "1").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_2_update_timeout.tf"),
					Check: commonStateChecker.
						CheckEqual("version_number", "2").
						CheckEqual("version_guid", "unique-guid-updated").
						CheckEqual("signed_certificate", "certificate-data-updated").
						CheckEqual("timeouts.update", "21m").Build(),
				},
			},
		},
		"error update certificate - same version number": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1_different_cert.tf"),
					ExpectError: regexp.MustCompile("Only updates with a different version_number are supported"),
				},
			},
		},
		"error update trust chain - same version number": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1_different_chain.tf"),
					ExpectError: regexp.MustCompile("Only updates with a different version_number are supported"),
				},
			},
		},
		"error update wait_for_deployment - same version number": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/no_wait_for_deployment.tf"),
					ExpectError: regexp.MustCompile("Only updates with a different version_number are supported"),
				},
			},
		},
		"error update auto_acknowledge_warnings - same version number": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				updateRequest := getDefaultUploadCertRequest()
				updateRequest.AcknowledgeAllWarnings = ptr.To(true)
				mockUploadSignedClientCertificate(m, updateRequest, nil)

				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/with_autoacknowledge.tf"),
					Check: commonStateChecker.
						CheckEqual("auto_acknowledge_warnings", "true").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/without_autoacknowledge.tf"),
					ExpectError: regexp.MustCompile("Only updates with a different version_number are supported"),
				},
			},
		},
		"error update timeout - same version number": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
				// Read before 2nd update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/create_timeout.tf"),
					ExpectError: regexp.MustCompile("Only updates with a different version_number are supported"),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/update_timeout.tf"),
					ExpectError: regexp.MustCompile("Only updates with a different version_number are supported"),
				},
			},
		},
		"error update version number - same signed certificate": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_2_without_cert_updated.tf"),
					ExpectError: regexp.MustCompile("No change in signed certificate"),
				},
			},
		},
		"error update - different certificate id": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/different_cert_id.tf"),
					ExpectError: regexp.MustCompile("updating field `client_certificate_id` is not possible"),
				},
			},
		},
		"update - API error on version checking": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()

				mockListClientCertificateVersions(m, nil, fmt.Errorf("API error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_2.tf"),
					ExpectError: regexp.MustCompile("API error"),
				},
			},
		},
		"update - API error on upload": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()

				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				updateRequest := getDefaultUploadCertRequest()
				updateRequest.Body.Certificate = "certificate-data-updated"
				updateRequest.Version = 2
				mockUploadSignedClientCertificate(m, updateRequest, fmt.Errorf("update upload error"))
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_2.tf"),
					ExpectError: regexp.MustCompile("update upload error"),
				},
			},
		},
		"update - API error on version retrieval": {
			init: func(m *mtlskeystore.Mock) {
				mockListClientCertificateVersions(m, getCustomReturnedVersions(1, "unique-guid", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				mockUploadSignedClientCertificate(m, getDefaultUploadCertRequest(), nil)
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Twice()
				// Read before update
				mockListClientCertificateVersions(m, getDefaultReturnedVersions(), nil).Once()

				mockListClientCertificateVersions(m, getCustomReturnedVersions(2, "unique-guid-updated", mtlskeystore.CertificateVersionStatusAwaitingSigned), nil).Once()
				updateRequest := getDefaultUploadCertRequest()
				updateRequest.Body.Certificate = "certificate-data-updated"
				updateRequest.Version = 2
				mockUploadSignedClientCertificate(m, updateRequest, nil)
				mockListClientCertificateVersions(m, nil, fmt.Errorf("update get versions error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_1.tf"),
					Check:  commonStateChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResCertificateUpload/version_2.tf"),
					ExpectError: regexp.MustCompile("update get versions error"),
				},
			},
		},
		"error incorrect config - missing client_certificate_id": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
  version_number        = 1
  signed_certificate    = "certificate-data"
}
`,
					ExpectError: regexp.MustCompile("The argument \"client_certificate_id\" is required"),
				},
			},
		},
		"error incorrect config - missing version_number": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
  client_certificate_id = 12345
  signed_certificate    = "certificate-data"
}
`,
					ExpectError: regexp.MustCompile("The argument \"version_number\" is required"),
				},
			},
		},
		"error incorrect config - missing signed_certificate": {
			steps: []resource.TestStep{
				{
					Config: `
provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlskeystore_client_certificate_upload" "client_certificate_upload" {
  client_certificate_id = 12345
  version_number        = 1
}
`,
					ExpectError: regexp.MustCompile("The argument \"signed_certificate\" is required"),
				},
			},
		},
	}
	pollingInterval = 1 * time.Millisecond
	defaultTimeout = 50 * time.Millisecond
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlskeystore.Mock{}
			if tc.init != nil {
				tc.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ExternalProviders: map[string]resource.ExternalProvider{
						"random": {
							Source:            "registry.terraform.io/hashicorp/random",
							VersionConstraint: "3.1.0",
						},
					},
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
