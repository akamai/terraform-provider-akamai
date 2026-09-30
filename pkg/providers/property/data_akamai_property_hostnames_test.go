package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	tst "github.com/akamai/terraform-provider-akamai/v11/internal/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataPropertyHostnames(t *testing.T) {
	t.Parallel()
	commonStateChecker := newHostnamesStateChecker(flattenHostnames(buildPropertyHostnames()))

	tests := map[string]struct {
		init        func(*edgegrid.TestClient)
		config      string
		checks      resource.TestCheckFunc
		expectError *regexp.Regexp
	}{
		"list hostnames with CCM Certificates, CCM Certificate Status, MTLS, TLS Configuration details": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnamesWithCCM()}

				client.PAPI.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
					ContractID:  "ctr_test",
					GroupID:     "grp_test",
					PropertyID:  "prp_test",
					ActivatedOn: "",
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 1,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   1,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 1,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
			checks: func() resource.TestCheckFunc {
				checker := test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
					CheckEqual("id", "prp_test1").
					CheckEqual("group_id", "grp_test").
					CheckEqual("contract_id", "ctr_test").
					CheckEqual("property_id", "prp_test").
					CheckEqual("version", "1").
					CheckEqual("hostnames.#", "3")

				for i := 0; i < 3; i++ {
					isEven := i%2 == 0
					checker = checker.
						CheckEqual(fmt.Sprintf("hostnames.%d.cname_type", i), "EDGE_HOSTNAME").
						CheckEqual(fmt.Sprintf("hostnames.%d.edge_hostname_id", i), fmt.Sprintf("ehn_%d", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.cname_from", i), fmt.Sprintf("cnamef%d.example.com", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.cname_to", i), fmt.Sprintf("cnamet%d.example.com.edgekey.net", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.cert_provisioning_type", i), "CCM").
						CheckEqual(fmt.Sprintf("hostnames.%d.cert_status.0.hostname", i), fmt.Sprintf("cnamef%d", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.cert_status.0.target", i), fmt.Sprintf("cnamet%d", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.cert_status.0.staging_status", i), "PENDING").
						CheckEqual(fmt.Sprintf("hostnames.%d.cert_status.0.production_status", i), "PENDING").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_certificates.0.rsa_cert_id", i), fmt.Sprintf("rsa_cert_%d", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_certificates.0.ecdsa_cert_id", i), fmt.Sprintf("ecdsa_cert_%d", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_status.0.ecdsa_staging_status", i), "ACTIVE").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_status.0.ecdsa_production_status", i), "ACTIVE").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_status.0.rsa_staging_status", i), "PENDING").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_status.0.rsa_production_status", i), "PENDING").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_id", i), fmt.Sprintf("%d", 164877+i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_link", i), fmt.Sprintf("/ccm/v2/lineages/%d", 164877+i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.#", i), "4").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.0.key_type", i), "RSA").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.0.network", i), "STAGING").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.0.status", i), "ACTIVE").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.1.key_type", i), "RSA").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.1.network", i), "PRODUCTION").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.1.status", i), "ACTIVE").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.2.key_type", i), "ECDSA").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.2.network", i), "STAGING").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.2.status", i), "ACTIVE").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.3.key_type", i), "ECDSA").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.3.network", i), "PRODUCTION").
						CheckEqual(fmt.Sprintf("hostnames.%d.ccm_cert_statuses.3.status", i), "ACTIVE").
						CheckEqual(fmt.Sprintf("hostnames.%d.mtls.0.ca_set_id", i), fmt.Sprintf("ca_set_%d", i)).
						CheckEqual(fmt.Sprintf("hostnames.%d.mtls.0.check_client_ocsp", i), fmt.Sprintf("%v", isEven)).
						CheckEqual(fmt.Sprintf("hostnames.%d.mtls.0.send_ca_set_client", i), fmt.Sprintf("%v", !isEven)).
						CheckEqual(fmt.Sprintf("hostnames.%d.tls_configuration.0.cipher_profile", i), "ak-akamai-default-2022q1").
						CheckEqual(fmt.Sprintf("hostnames.%d.tls_configuration.0.disallowed_tls_versions.#", i), "2").
						CheckEqual(fmt.Sprintf("hostnames.%d.tls_configuration.0.disallowed_tls_versions.0", i), "TLSv1").
						CheckEqual(fmt.Sprintf("hostnames.%d.tls_configuration.0.disallowed_tls_versions.1", i), "TLSv1_1").
						CheckEqual(fmt.Sprintf("hostnames.%d.tls_configuration.0.staple_server_ocsp_response", i), fmt.Sprintf("%v", isEven)).
						CheckEqual(fmt.Sprintf("hostnames.%d.tls_configuration.0.fips_mode", i), fmt.Sprintf("%v", !isEven))
				}
				return checker.Build()
			}(),
		},
		"list hostnames": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)
				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}

				client.PAPI.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
					ContractID:  "ctr_test",
					GroupID:     "grp_test",
					PropertyID:  "prp_test",
					ActivatedOn: "",
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 1,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   1,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 1,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
			checks: commonStateChecker.Build(),
		},
		"list hostnames of type HOSTNAME_BUCKET": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, ptr.To("HOSTNAME_BUCKET")).Times(3)

				mockListActivePropertyHostnames(client.PAPI, 0, &papi.ListActivePropertyHostnamesResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					PropertyID: "prp_test",
					AccountID:  "act_test",
					Hostnames: papi.HostnamesResponseItems{
						Items: buildHostnameItems(12),
					},
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
			checks: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
				CheckEqual("id", "prp_test").
				CheckEqual("group_id", "grp_test").
				CheckEqual("contract_id", "ctr_test").
				CheckEqual("property_id", "prp_test").
				CheckMissing("version").
				CheckEqual("hostname_bucket.#", "12").
				CheckEqual("hostname_bucket.0.cname_from", "cnamef0").
				CheckEqual("hostname_bucket.0.cname_type", "EDGE_HOSTNAME").
				CheckEqual("hostname_bucket.0.staging_edge_hostname_id", "ehn0").
				CheckEqual("hostname_bucket.0.staging_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.0.staging_cname_to", "cnamet0").
				CheckEqual("hostname_bucket.0.production_edge_hostname_id", "ehn0").
				CheckEqual("hostname_bucket.0.production_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.0.production_cname_to", "cnamet0").
				CheckEqual("hostname_bucket.0.cert_status.0.hostname", "cnamef0").
				CheckEqual("hostname_bucket.0.cert_status.0.target", "cnamet0").
				CheckEqual("hostname_bucket.0.cert_status.0.staging_status", "PENDING").
				CheckEqual("hostname_bucket.0.cert_status.0.production_status", "PENDING").
				Build(),
		},
		"list hostnames of type HOSTNAME_BUCKET with results on several pages": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, ptr.To("HOSTNAME_BUCKET")).Times(3)

				hostnames := buildHostnameItems(listActivePropertyHostnamesResultsPerPage + 3)

				mockListActivePropertyHostnames(client.PAPI, 0, &papi.ListActivePropertyHostnamesResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					PropertyID: "prp_test",
					AccountID:  "act_test",
					Hostnames: papi.HostnamesResponseItems{
						Items: hostnames[:listActivePropertyHostnamesResultsPerPage],
					},
				}, nil).Times(3)
				mockListActivePropertyHostnames(client.PAPI, listActivePropertyHostnamesResultsPerPage, &papi.ListActivePropertyHostnamesResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					PropertyID: "prp_test",
					AccountID:  "act_test",
					Hostnames: papi.HostnamesResponseItems{
						Items: hostnames[listActivePropertyHostnamesResultsPerPage:],
					},
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
			checks: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
				CheckEqual("id", "prp_test").
				CheckEqual("group_id", "grp_test").
				CheckEqual("contract_id", "ctr_test").
				CheckEqual("property_id", "prp_test").
				CheckMissing("version").
				CheckEqual("hostname_bucket.#", "53").
				CheckEqual("hostname_bucket.0.cname_from", "cnamef0").
				CheckEqual("hostname_bucket.0.cname_type", "EDGE_HOSTNAME").
				CheckEqual("hostname_bucket.0.staging_edge_hostname_id", "ehn0").
				CheckEqual("hostname_bucket.0.staging_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.0.staging_cname_to", "cnamet0").
				CheckEqual("hostname_bucket.0.production_edge_hostname_id", "ehn0").
				CheckEqual("hostname_bucket.0.production_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.0.production_cname_to", "cnamet0").
				CheckEqual("hostname_bucket.0.cert_status.0.hostname", "cnamef0").
				CheckEqual("hostname_bucket.0.cert_status.0.target", "cnamet0").
				CheckEqual("hostname_bucket.0.cert_status.0.staging_status", "PENDING").
				CheckEqual("hostname_bucket.0.cert_status.0.production_status", "PENDING").
				Build(),
		},
		"list hostnames - filter_pending_default_certs set to true": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, ptr.To("HOSTNAME_BUCKET")).Times(3)

				hostnames := buildHostnameItems(12)

				hostnames = append(hostnames, []papi.HostnameItem{
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef13",
								Target:   "cnamet13",
							},
							Staging:    []papi.StatusItem{{Status: "PENDING"}},
							Production: []papi.StatusItem{{Status: "PENDING"}},
						},
						CnameFrom:                "cnamef13",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       papi.CertTypeCPSManaged,
						ProductionCnameTo:        "cnamet13",
						ProductionEdgeHostnameID: "ehn13",
						StagingCertType:          papi.CertTypeCPSManaged,
						StagingCnameTo:           "cnamet13",
						StagingEdgeHostnameID:    "ehn13",
					},
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef14",
								Target:   "cnamet14",
							},
							Staging:    []papi.StatusItem{{Status: "DEPLOYED"}},
							Production: []papi.StatusItem{{Status: "DEPLOYED"}},
						},
						CnameFrom:                "cnamef14",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       papi.CertTypeDefault,
						ProductionCnameTo:        "cnamet14",
						ProductionEdgeHostnameID: "ehn14",
						StagingCertType:          papi.CertTypeDefault,
						StagingCnameTo:           "cnamet14",
						StagingEdgeHostnameID:    "ehn14",
					},
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef15",
								Target:   "cnamet15",
							},
							Staging:    []papi.StatusItem{{Status: "DEPLOYED"}},
							Production: []papi.StatusItem{},
						},
						CnameFrom:                "cnamef15",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       "",
						ProductionCnameTo:        "",
						ProductionEdgeHostnameID: "",
						StagingCertType:          papi.CertTypeDefault,
						StagingCnameTo:           "cnamet15",
						StagingEdgeHostnameID:    "ehn15",
					},
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef16",
								Target:   "cnamet16",
							},
							Staging:    []papi.StatusItem{{Status: "DEPLOYED"}},
							Production: []papi.StatusItem{{Status: "PENDING"}},
						},
						CnameFrom:                "cnamef16",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       papi.CertTypeDefault,
						ProductionCnameTo:        "cnamet16",
						ProductionEdgeHostnameID: "ehn16",
						StagingCertType:          papi.CertTypeDefault,
						StagingCnameTo:           "cnamet16",
						StagingEdgeHostnameID:    "ehn16",
					},
				}...)

				mockListActivePropertyHostnames(client.PAPI, 0, &papi.ListActivePropertyHostnamesResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					PropertyID: "prp_test",
					AccountID:  "act_test",
					Hostnames: papi.HostnamesResponseItems{
						Items: hostnames,
					},
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_filter.tf"),
			checks: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
				CheckEqual("id", "prp_test").
				CheckEqual("group_id", "grp_test").
				CheckEqual("contract_id", "ctr_test").
				CheckEqual("property_id", "prp_test").
				CheckMissing("version").
				CheckEqual("hostname_bucket.#", "13").
				CheckEqual("hostname_bucket.0.cname_from", "cnamef0").
				CheckEqual("hostname_bucket.0.cname_type", "EDGE_HOSTNAME").
				CheckEqual("hostname_bucket.0.staging_edge_hostname_id", "ehn0").
				CheckEqual("hostname_bucket.0.staging_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.0.staging_cname_to", "cnamet0").
				CheckEqual("hostname_bucket.0.production_edge_hostname_id", "ehn0").
				CheckEqual("hostname_bucket.0.production_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.0.production_cname_to", "cnamet0").
				CheckEqual("hostname_bucket.0.cert_status.0.hostname", "cnamef0").
				CheckEqual("hostname_bucket.0.cert_status.0.target", "cnamet0").
				CheckEqual("hostname_bucket.0.cert_status.0.staging_status", "PENDING").
				CheckEqual("hostname_bucket.0.cert_status.0.production_status", "PENDING").
				CheckEqual("hostname_bucket.12.cname_from", "cnamef16").
				CheckEqual("hostname_bucket.12.cname_type", "EDGE_HOSTNAME").
				CheckEqual("hostname_bucket.12.staging_edge_hostname_id", "ehn16").
				CheckEqual("hostname_bucket.12.staging_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.12.staging_cname_to", "cnamet16").
				CheckEqual("hostname_bucket.12.production_edge_hostname_id", "ehn16").
				CheckEqual("hostname_bucket.12.production_cert_type", "DEFAULT").
				CheckEqual("hostname_bucket.12.production_cname_to", "cnamet16").
				CheckEqual("hostname_bucket.12.cert_status.0.hostname", "cnamef16").
				CheckEqual("hostname_bucket.12.cert_status.0.target", "cnamet16").
				CheckEqual("hostname_bucket.12.cert_status.0.staging_status", "DEPLOYED").
				CheckEqual("hostname_bucket.12.cert_status.0.production_status", "PENDING").
				Build(),
		},
		"list hostnames with status `DEPLOYED` - filter_pending_default_certs set to true": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, ptr.To("HOSTNAME_BUCKET")).Times(3)

				hostnames := []papi.HostnameItem{
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef1",
								Target:   "cnamet1",
							},
							Staging:    []papi.StatusItem{{Status: "DEPLOYED"}},
							Production: []papi.StatusItem{{Status: "DEPLOYED"}},
						},
						CnameFrom:                "cnamef1",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       papi.CertTypeDefault,
						ProductionCnameTo:        "cnamet1",
						ProductionEdgeHostnameID: "ehn1",
						StagingCertType:          papi.CertTypeDefault,
						StagingCnameTo:           "cnamet1",
						StagingEdgeHostnameID:    "ehn1",
					},
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef2",
								Target:   "cnamet2",
							},
							Staging:    []papi.StatusItem{{Status: "DEPLOYED"}},
							Production: []papi.StatusItem{{Status: "DEPLOYED"}},
						},
						CnameFrom:                "cnamef2",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       papi.CertTypeDefault,
						ProductionCnameTo:        "cnamet2",
						ProductionEdgeHostnameID: "ehn2",
						StagingCertType:          papi.CertTypeDefault,
						StagingCnameTo:           "cnamet2",
						StagingEdgeHostnameID:    "ehn2",
					},
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef3",
								Target:   "cnamet3",
							},
							Staging:    []papi.StatusItem{{Status: "DEPLOYED"}},
							Production: []papi.StatusItem{{Status: "DEPLOYED"}},
						},
						CnameFrom:                "cnamef3",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       papi.CertTypeDefault,
						ProductionCnameTo:        "cnamet3",
						ProductionEdgeHostnameID: "ehn3",
						StagingCertType:          papi.CertTypeDefault,
						StagingCnameTo:           "cnamet3",
						StagingEdgeHostnameID:    "ehn3",
					},
					{
						CertStatus: &papi.CertStatusItem{
							ValidationCname: papi.ValidationCname{
								Hostname: "cnamef16",
								Target:   "cnamet16",
							},
							Staging:    []papi.StatusItem{{Status: "DEPLOYED"}},
							Production: []papi.StatusItem{{Status: "DEPLOYED"}},
						},
						CnameFrom:                "cnamef4",
						CnameType:                papi.HostnameCnameTypeEdgeHostname,
						ProductionCertType:       papi.CertTypeDefault,
						ProductionCnameTo:        "cnamet4",
						ProductionEdgeHostnameID: "ehn4",
						StagingCertType:          papi.CertTypeDefault,
						StagingCnameTo:           "cnamet4",
						StagingEdgeHostnameID:    "ehn4",
					},
				}

				mockListActivePropertyHostnames(client.PAPI, 0, &papi.ListActivePropertyHostnamesResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					PropertyID: "prp_test",
					AccountID:  "act_test",
					Hostnames: papi.HostnamesResponseItems{
						Items: hostnames,
					},
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_filter.tf"),
			checks: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
				CheckEqual("id", "prp_test").
				CheckEqual("group_id", "grp_test").
				CheckEqual("contract_id", "ctr_test").
				CheckEqual("property_id", "prp_test").
				CheckMissing("version").
				CheckEqual("hostname_bucket.#", "0").
				Build(),
		},
		"list hostnames without group prefix": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}

				client.PAPI.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
					ContractID:  "ctr_test",
					GroupID:     "grp_test",
					PropertyID:  "prp_test",
					ActivatedOn: "",
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 1,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   1,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 1,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_no_group_prefix.tf"),
			checks: commonStateChecker.
				CheckEqual("group_id", "test").
				Build(),
		},
		"list hostnames without contract prefix": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}

				client.PAPI.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
					ContractID:  "ctr_test",
					GroupID:     "grp_test",
					PropertyID:  "prp_test",
					ActivatedOn: "",
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 1,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   1,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 1,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_no_contract_prefix.tf"),
			checks: commonStateChecker.
				CheckEqual("contract_id", "test").
				Build(),
		},
		"list hostnames without property prefix": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}

				client.PAPI.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
					ContractID:  "ctr_test",
					GroupID:     "grp_test",
					PropertyID:  "prp_test",
					ActivatedOn: "",
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 1,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   1,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 1,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_no_property_prefix.tf"),
			checks: commonStateChecker.Build(),
		},
		"specify property version to fetch": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}

				client.PAPI.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 5,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   5,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
			checks: commonStateChecker.
				CheckEqual("version", "5").
				CheckEqual("id", "prp_test5").
				Build(),
		},
		"property with full domain ownership verification data": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				dov := papi.DomainOwnershipVerification{
					ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-12-31T23:59:59Z"),
					Status:                   "PENDING",
					ValidationCname: &papi.ValidationCname{
						Hostname: "example.com",
						Target:   "example.akamai.net",
					},
					ValidationHTTP: &papi.ValidationHTTP{
						FileContentMethod: papi.FileContentMethod{
							Body: "test-body",
							URL:  "test-url",
						},
						RedirectMethod: papi.RedirectMethod{
							HTTPRedirectFrom: "redirect-from",
							HTTPRedirectTo:   "redirect-to",
						},
					},
					ValidationTXT: &papi.ValidationTXT{
						ChallengeToken: "test-token",
						Hostname:       "test-example.com",
					},
				}
				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnamesWithDOV(&dov)}

				client.PAPI.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 5,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   5,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
			checks: newHostnamesStateChecker(flattenHostnames(buildPropertyHostnamesWithDOV(&papi.DomainOwnershipVerification{
				ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-12-31T23:59:59Z"),
				Status:                   "PENDING",
				ValidationCname: &papi.ValidationCname{
					Hostname: "example.com",
					Target:   "example.akamai.net",
				},
				ValidationHTTP: &papi.ValidationHTTP{
					FileContentMethod: papi.FileContentMethod{
						Body: "test-body",
						URL:  "test-url",
					},
					RedirectMethod: papi.RedirectMethod{
						HTTPRedirectFrom: "redirect-from",
						HTTPRedirectTo:   "redirect-to",
					},
				},
				ValidationTXT: &papi.ValidationTXT{
					ChallengeToken: "test-token",
					Hostname:       "test-example.com",
				},
			}))).
				CheckEqual("id", "prp_test5").
				CheckEqual("version", "5").
				Build(),
		},
		"property with domain ownership verification only status": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				dov := papi.DomainOwnershipVerification{
					Status: "VALIDATED",
				}
				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnamesWithDOV(&dov)}

				client.PAPI.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 5,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   5,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
			checks: newHostnamesStateChecker(flattenHostnames(buildPropertyHostnamesWithDOV(&papi.DomainOwnershipVerification{
				Status: "VALIDATED",
			}))).
				CheckEqual("id", "prp_test5").
				CheckEqual("version", "5").
				Build(),
		},
		"property with domain ownership verification partial data": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)
				dov := papi.DomainOwnershipVerification{
					ChallengeTokenExpiryDate: tst.NewTimeFromStringPtr(t, "2024-12-31T23:59:59Z"),
					Status:                   "PENDING",
					ValidationCname: &papi.ValidationCname{
						Hostname: "example.com",
						Target:   "example.akamai.net",
					},
					ValidationTXT: &papi.ValidationTXT{
						ChallengeToken: "test-token",
						Hostname:       "test-example.com",
					},
				}
				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnamesWithDOV(&dov)}

				client.PAPI.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 5,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   5,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
			checks: newHostnamesStateChecker(flattenHostnames(buildPropertyHostnamesWithDOV(&papi.DomainOwnershipVerification{
				Status: "PENDING",
				ValidationCname: &papi.ValidationCname{
					Hostname: "example.com",
					Target:   "example.akamai.net",
				},
				ValidationTXT: &papi.ValidationTXT{
					ChallengeToken: "test-token",
					Hostname:       "test-example.com",
				},
			}))).
				CheckEqual("id", "prp_test5").
				CheckEqual("version", "5").
				Build(),
		},
		"property with authorization data": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Times(3)

				hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnamesWithAuthorization()}

				client.PAPI.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
					ContractID:  "ctr_test",
					GroupID:     "grp_test",
					PropertyID:  "prp_test",
					ActivatedOn: "",
				}).Return(&papi.GetPropertyVersionsResponse{
					ContractID: "ctr_test",
					GroupID:    "grp_test",
					Version: papi.PropertyVersionGetItem{
						PropertyVersion: 1,
					},
				}, nil).Times(3)
				client.PAPI.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
					PropertyID:        "prp_test",
					PropertyVersion:   1,
					ContractID:        "ctr_test",
					GroupID:           "grp_test",
					ValidateHostnames: false,
					IncludeCertStatus: true,
				}).Return(&papi.GetPropertyVersionHostnamesResponse{
					AccountID:       "act_test",
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 1,
					Etag:            "etag",
					Hostnames:       hostnames,
				}, nil).Times(3)
			},
			config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
			checks: newHostnamesStateChecker(flattenHostnames(buildPropertyHostnamesWithAuthorization())).Build(),
		},
		"error - specify property version to fetch with error": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, nil).Once()

				client.PAPI.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
					ContractID:      "ctr_test",
					GroupID:         "grp_test",
					PropertyID:      "prp_test",
					PropertyVersion: 5,
				}).Return(nil, fmt.Errorf("error fetching property version")).Once()
			},
			config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
			expectError: regexp.MustCompile(`error fetching property version`),
		},
		"error - list hostnames of type HOSTNAME_BUCKET fails": {
			init: func(client *edgegrid.TestClient) {
				mockGetPropertyWithPropertyType(client.PAPI, ptr.To("HOSTNAME_BUCKET")).Once()

				mockListActivePropertyHostnames(client.PAPI, 0, nil, fmt.Errorf("error fetching list hostnames")).Once()
			},
			config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
			expectError: regexp.MustCompile(`error fetching list hostnames`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{{
					Config:      tc.config,
					Check:       tc.checks,
					ExpectError: tc.expectError,
				}},
			})

			client.PAPI.AssertExpectations(t)
		})
	}
}

func mockListActivePropertyHostnames(client *papi.Mock, offset int, resp *papi.ListActivePropertyHostnamesResponse, err error) *mock.Call {
	return client.On("ListActivePropertyHostnames", testutils.MockContext, papi.ListActivePropertyHostnamesRequest{
		ContractID:        "ctr_test",
		GroupID:           "grp_test",
		PropertyID:        "prp_test",
		IncludeCertStatus: true,
		Limit:             50,
		Offset:            offset,
	}).Return(resp, err)
}

func buildPropertyHostnames() []papi.HostnameResponseItem {
	hostnames := make([]papi.HostnameResponseItem, 10)
	for i := range 10 {
		hostnames[i] = papi.HostnameResponseItem{
			CnameType:            "EDGE_HOSTNAME",
			EdgeHostnameID:       fmt.Sprintf("ehn%v", i),
			CnameFrom:            fmt.Sprintf("cnamef%v", i),
			CnameTo:              fmt.Sprintf("cnamet%v", i),
			CertProvisioningType: "DEFAULT",
			CertStatus: papi.CertStatusItem{
				ValidationCname: papi.ValidationCname{
					Hostname: fmt.Sprintf("cnamef%v", i),
					Target:   fmt.Sprintf("cnamet%v", i),
				},
				Staging: []papi.StatusItem{{
					Status: "PENDING",
				}},
				Production: []papi.StatusItem{{
					Status: "PENDING",
				},
				},
			},
		}
	}
	return hostnames
}

func buildPropertyHostnamesWithDOV(dov *papi.DomainOwnershipVerification) []papi.HostnameResponseItem {
	hostnames := make([]papi.HostnameResponseItem, 10)
	for i := range 10 {
		hostnames[i] = papi.HostnameResponseItem{
			CnameType:            "EDGE_HOSTNAME",
			EdgeHostnameID:       fmt.Sprintf("ehn%v", i),
			CnameFrom:            fmt.Sprintf("cnamef%v", i),
			CnameTo:              fmt.Sprintf("cnamet%v", i),
			CertProvisioningType: "DEFAULT",
			CertStatus: papi.CertStatusItem{
				ValidationCname: papi.ValidationCname{
					Hostname: fmt.Sprintf("cnamef%v", i),
					Target:   fmt.Sprintf("cnamet%v", i),
				},
				Staging: []papi.StatusItem{{
					Status: "PENDING",
				}},
				Production: []papi.StatusItem{{
					Status: "PENDING",
				},
				},
			},
			DomainOwnershipVerification: dov,
		}
	}
	return hostnames
}

func buildPropertyHostnamesWithAuthorization() []papi.HostnameResponseItem {
	hostnames := make([]papi.HostnameResponseItem, 10)
	for i := range 10 {
		validUntil := tst.NewTimeFromStringPtr(&testing.T{}, "2024-12-31T23:59:59Z")
		hostnames[i] = papi.HostnameResponseItem{
			CnameType:            "EDGE_HOSTNAME",
			EdgeHostnameID:       fmt.Sprintf("ehn%v", i),
			CnameFrom:            fmt.Sprintf("cnamef%v", i),
			CnameTo:              fmt.Sprintf("cnamet%v", i),
			CertProvisioningType: "DEFAULT",
			CertStatus: papi.CertStatusItem{
				ValidationCname: papi.ValidationCname{
					Hostname: fmt.Sprintf("cnamef%v", i),
					Target:   fmt.Sprintf("cnamet%v", i),
				},
				Staging: []papi.StatusItem{{
					Status: "PENDING",
				}},
				Production: []papi.StatusItem{{
					Status: "PENDING",
				}},
				Authorization: &papi.Authorization{
					Status:     "VALID",
					ValidUntil: validUntil,
					DNS01: &papi.DNSAuthorization{
						Value: fmt.Sprintf("dns-token-%d", i),
						Result: papi.AuthorizationResult{
							Message:   "DNS challenge generated",
							Source:    "CPS",
							Timestamp: *validUntil,
						},
					},
					HTTP01: &papi.HTTPAuthorization{
						Body: fmt.Sprintf("http-body-%d", i),
						URL:  fmt.Sprintf("http://example.com/.well-known/acme-challenge-%d", i),
						Result: papi.AuthorizationResult{
							Message:   "HTTP challenge generated",
							Source:    "CA",
							Timestamp: *validUntil,
						},
					},
				},
			},
		}
	}
	return hostnames
}

func buildPropertyHostnamesWithCCM() []papi.HostnameResponseItem {
	hostnames := make([]papi.HostnameResponseItem, 3)
	for i := 0; i < 3; i++ {
		// Alternate boolean values to test both true and false scenarios
		isEven := i%2 == 0

		hostnames[i] = papi.HostnameResponseItem{
			CnameType:            "EDGE_HOSTNAME",
			EdgeHostnameID:       fmt.Sprintf("ehn_%d", i),
			CnameFrom:            fmt.Sprintf("cnamef%d.example.com", i),
			CnameTo:              fmt.Sprintf("cnamet%d.example.com.edgekey.net", i),
			CertProvisioningType: "CCM",
			CertStatus: papi.CertStatusItem{
				ValidationCname: papi.ValidationCname{
					Hostname: fmt.Sprintf("cnamef%v", i),
					Target:   fmt.Sprintf("cnamet%v", i),
				},
				Staging: []papi.StatusItem{{
					Status: "PENDING",
				}},
				Production: []papi.StatusItem{{
					Status: "PENDING",
				},
				},
			},
			CCMCertificates: &papi.CCMCertificatesResp{ //nolint:staticcheck
				CCMCertificates: papi.CCMCertificates{ //nolint:staticcheck
					ECDSACertID: fmt.Sprintf("ecdsa_cert_%d", i),
					RSACertID:   fmt.Sprintf("rsa_cert_%d", i),
				},
			},
			CCMCertStatus: &papi.CCMCertStatus{
				ECDSAStagingStatus:    "ACTIVE",
				ECDSAProductionStatus: "ACTIVE",
				RSAStagingStatus:      "PENDING",
				RSAProductionStatus:   "PENDING",
			},
			CCMCertID:   ptr.To(fmt.Sprintf("%d", 164877+i)),
			CCMCertLink: ptr.To(fmt.Sprintf("/ccm/v2/lineages/%d", 164877+i)),
			CCMCertStatuses: []papi.CCMCertStatusItem{
				{KeyType: "RSA", Network: "STAGING", Status: "ACTIVE"},
				{KeyType: "RSA", Network: "PRODUCTION", Status: "ACTIVE"},
				{KeyType: "ECDSA", Network: "STAGING", Status: "ACTIVE"},
				{KeyType: "ECDSA", Network: "PRODUCTION", Status: "ACTIVE"},
			},
			MTLS: &papi.MTLSResp{
				CASetLink: fmt.Sprintf("/ccm/v3/ca-sets/ca_set_%d", i),
				MTLS: papi.MTLS{
					CASetID:         fmt.Sprintf("ca_set_%d", i),
					CheckClientOCSP: isEven,
					SendCASetClient: !isEven,
				},
			},
			TLSConfiguration: &papi.TLSConfiguration{
				CipherProfile:            "ak-akamai-default-2022q1",
				DisallowedTLSVersions:    []string{"TLSv1", "TLSv1_1"},
				StapleServerOcspResponse: isEven,
				FIPSMode:                 !isEven,
			},
		}
	}
	return hostnames
}

func buildHostnameItems(itemsNo int) []papi.HostnameItem {
	hostnames := make([]papi.HostnameItem, 0, itemsNo)
	for i := 0; i < itemsNo; i++ {
		hostnames = append(hostnames, papi.HostnameItem{
			CertStatus: &papi.CertStatusItem{
				ValidationCname: papi.ValidationCname{
					Hostname: fmt.Sprintf("cnamef%v", i),
					Target:   fmt.Sprintf("cnamet%v", i),
				},
				Staging:    []papi.StatusItem{{Status: "PENDING"}},
				Production: []papi.StatusItem{{Status: "PENDING"}},
			},
			CnameFrom:                fmt.Sprintf("cnamef%v", i),
			CnameType:                papi.HostnameCnameTypeEdgeHostname,
			ProductionCertType:       papi.CertTypeDefault,
			ProductionCnameTo:        fmt.Sprintf("cnamet%v", i),
			ProductionEdgeHostnameID: fmt.Sprintf("ehn%v", i),
			StagingCertType:          papi.CertTypeDefault,
			StagingCnameTo:           fmt.Sprintf("cnamet%v", i),
			StagingEdgeHostnameID:    fmt.Sprintf("ehn%v", i),
		})
	}
	return hostnames
}

func mockGetPropertyWithPropertyType(client *papi.Mock, propertyType *string) *mock.Call {
	return client.On("GetProperty", testutils.MockContext, papi.GetPropertyRequest{
		ContractID: "ctr_test",
		GroupID:    "grp_test",
		PropertyID: "prp_test",
	}).Return(&papi.GetPropertyResponse{
		Property: &papi.Property{
			AccountID:    "act_test",
			ContractID:   "ctr_test",
			GroupID:      "grp_test",
			PropertyID:   "prp_test",
			PropertyType: propertyType,
		},
	}, nil)
}

func checkAuthorizationField(checker test.StateChecker, ind, cInd int, mapKey, cKey string, authList []map[string]any) test.StateChecker {
	if len(authList) == 0 {
		return checker
	}
	auth := authList[0]
	for aKey, aVal := range auth {
		switch aKey {
		case "dns01":
			if dns01List, ok := aVal.([]map[string]any); ok && len(dns01List) > 0 {
				dns01 := dns01List[0]
				for dKey, dVal := range dns01 {
					if dKey == "result" {
						if resultList, ok := dVal.([]map[string]any); ok && len(resultList) > 0 {
							result := resultList[0]
							for rKey, rVal := range result {
								value := fmt.Sprintf("%v", rVal)
								key := fmt.Sprintf("hostnames.%d.%s.%d.%s.0.%s.0.%s.0.%s", ind, mapKey, cInd, cKey, aKey, dKey, rKey)
								checker = checker.CheckEqual(key, value)
							}
						}
					} else {
						value := fmt.Sprintf("%v", dVal)
						key := fmt.Sprintf("hostnames.%d.%s.%d.%s.0.%s.0.%s", ind, mapKey, cInd, cKey, aKey, dKey)
						checker = checker.CheckEqual(key, value)
					}
				}
			}
		case "http01":
			if http01List, ok := aVal.([]map[string]any); ok && len(http01List) > 0 {
				http01 := http01List[0]
				for hKey, hVal := range http01 {
					if hKey == "result" {
						if resultList, ok := hVal.([]map[string]any); ok && len(resultList) > 0 {
							result := resultList[0]
							for rKey, rVal := range result {
								value := fmt.Sprintf("%v", rVal)
								key := fmt.Sprintf("hostnames.%d.%s.%d.%s.0.%s.0.%s.0.%s", ind, mapKey, cInd, cKey, aKey, hKey, rKey)
								checker = checker.CheckEqual(key, value)
							}
						}
					} else {
						value := fmt.Sprintf("%v", hVal)
						key := fmt.Sprintf("hostnames.%d.%s.%d.%s.0.%s.0.%s", ind, mapKey, cInd, cKey, aKey, hKey)
						checker = checker.CheckEqual(key, value)
					}
				}
			}
		default:
			value := fmt.Sprintf("%v", aVal)
			key := fmt.Sprintf("hostnames.%d.%s.%d.%s.0.%s", ind, mapKey, cInd, cKey, aKey)
			checker = checker.CheckEqual(key, value)
		}
	}
	return checker
}

func checkCertStatusField(checker test.StateChecker, ind int, mapKey string, mapVal any) test.StateChecker {
	certStatuses := mapVal.([]map[string]interface{})
	for cInd, cert := range certStatuses {
		for cKey, cVal := range cert {
			if cKey == "authorization" {
				if authList, ok := cVal.([]map[string]any); ok {
					checker = checkAuthorizationField(checker, ind, cInd, mapKey, cKey, authList)
				}
				continue
			}
			value := fmt.Sprintf("%v", cVal)
			key := fmt.Sprintf("hostnames.%v.%v.%v.%v", ind, mapKey, cInd, cKey)
			checker = checker.CheckEqual(key, value)
		}
	}
	return checker
}

func checkDomainOwnershipVerificationField(checker test.StateChecker, ind int, mapKey string, mapVal any) test.StateChecker {
	if mapVal == nil {
		return checker.CheckMissing(mapKey)
	}
	dovVal := mapVal.([]map[string]any)
	for dovKey, dovMapVal := range dovVal[0] {
		switch dovKey {
		case "validation_http":
			validationHTTP := dovMapVal.([]map[string]any)
			for vKey, vMapVal := range validationHTTP[0] {
				nestedAttrValidationHTTP := vMapVal.([]map[string]any)
				for nvKey, nvMapVal := range nestedAttrValidationHTTP[0] {
					value := fmt.Sprintf("%v", nvMapVal)
					key := fmt.Sprintf("hostnames.%v.%v.0.%v.0.%v.0.%v", ind, mapKey, dovKey, vKey, nvKey)
					checker = checker.CheckEqual(key, value)
				}
			}
		case "validation_txt":
			validationTXT := dovMapVal.([]map[string]any)
			for tKey, tMapVal := range validationTXT[0] {
				value := fmt.Sprintf("%v", tMapVal)
				key := fmt.Sprintf("hostnames.%v.%v.0.%v.0.%v", ind, mapKey, dovKey, tKey)
				checker = checker.CheckEqual(key, value)
			}
		case "validation_cname":
			validationCname := dovMapVal.([]map[string]any)
			for cKey, cMapVal := range validationCname[0] {
				value := fmt.Sprintf("%v", cMapVal)
				key := fmt.Sprintf("hostnames.%v.%v.0.%v.0.%v", ind, mapKey, dovKey, cKey)
				checker = checker.CheckEqual(key, value)
			}
		case "challenge_token_expiry_date", "status":
			if value := dovMapVal.(string); value != "" {
				key := fmt.Sprintf("hostnames.%v.%v.0.%v", ind, mapKey, dovKey)
				checker = checker.CheckEqual(key, value)
			}
		}
	}
	return checker
}

func checkCCMCertStatusesField(checker test.StateChecker, ind int, mapKey string, mapVal any) test.StateChecker {
	items, ok := mapVal.([]map[string]string)
	if !ok {
		return checker
	}
	baseKey := fmt.Sprintf("hostnames.%d.%s", ind, mapKey)
	checker = checker.CheckEqual(fmt.Sprintf("%s.#", baseKey), fmt.Sprintf("%d", len(items)))
	for itemIndex, itemMap := range items {
		for itemKey, itemValue := range itemMap {
			key := fmt.Sprintf("%s.%d.%s", baseKey, itemIndex, itemKey)
			checker = checker.CheckEqual(key, itemValue)
		}
	}
	return checker
}

func checkCCMDetailField(checker test.StateChecker, ind int, mapKey string, mapVal any) test.StateChecker {
	items, ok := mapVal.([]map[string]interface{})
	if !ok {
		return checker
	}
	for itemIndex, itemMap := range items {
		for itemKey, itemValue := range itemMap {
			if v, ok := itemValue.([]string); ok {
				baseKey := fmt.Sprintf("hostnames.%d.%s.%d.%s", ind, mapKey, itemIndex, itemKey)
				checker = checker.CheckEqual(fmt.Sprintf("%s.#", baseKey), fmt.Sprintf("%d", len(v)))
				for i, strVal := range v {
					elementKey := fmt.Sprintf("%s.%d", baseKey, i)
					checker = checker.CheckEqual(elementKey, strVal)
				}
				continue
			}
			value := fmt.Sprintf("%v", itemValue)
			key := fmt.Sprintf("hostnames.%d.%s.%d.%s", ind, mapKey, itemIndex, itemKey)
			checker = checker.CheckEqual(key, value)
		}
	}
	return checker
}

func newHostnamesStateChecker(hostnames []map[string]any) test.StateChecker {
	checker := test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
		CheckEqual("id", "prp_test1").
		CheckEqual("group_id", "grp_test").
		CheckEqual("contract_id", "ctr_test").
		CheckEqual("property_id", "prp_test").
		CheckEqual("version", "1").
		CheckEqual("hostnames.#", "10")

	for ind, hostname := range hostnames {
		for mapKey, mapVal := range hostname {
			switch mapKey {
			case "cert_status":
				checker = checkCertStatusField(checker, ind, mapKey, mapVal)
			case "domain_ownership_verification":
				checker = checkDomainOwnershipVerificationField(checker, ind, mapKey, mapVal)
			case "ccm_cert_statuses":
				checker = checkCCMCertStatusesField(checker, ind, mapKey, mapVal)
			case "ccm_cert_status", "ccm_certificates", "mtls", "tls_configuration":
				checker = checkCCMDetailField(checker, ind, mapKey, mapVal)
			default:
				value := fmt.Sprintf("%v", mapVal)
				key := fmt.Sprintf("hostnames.%v.%v", ind, mapKey)
				checker = checker.CheckEqual(key, value)
			}
		}
	}
	return checker
}
