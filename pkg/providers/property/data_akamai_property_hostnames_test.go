package property

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataPropertyHostnames(t *testing.T) {
	t.Run("list hostnames", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, nil).Times(3)

		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
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
		client.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
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

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
					Check:  buildAggregatedHostnamesTest(hostnameItems, "prp_test1", "grp_test", "ctr_test", "prp_test", 1),
				}},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("list hostnames of type HOSTNAME_BUCKET", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, ptr.To("HOSTNAME_BUCKET")).Times(3)

		mockListActivePropertyHostnames(client, 0, &papi.ListActivePropertyHostnamesResponse{
			ContractID: "ctr_test",
			GroupID:    "grp_test",
			PropertyID: "prp_test",
			AccountID:  "act_test",
			Hostnames: papi.HostnamesResponseItems{
				Items: buildHostnameItems(12),
			},
		}, nil).Times(3)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
					Check: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
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
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames of type HOSTNAME_BUCKET with results on several pages", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, ptr.To("HOSTNAME_BUCKET")).Times(3)

		hostnames := buildHostnameItems(listActivePropertyHostnamesResultsPerPage + 3)

		mockListActivePropertyHostnames(client, 0, &papi.ListActivePropertyHostnamesResponse{
			ContractID: "ctr_test",
			GroupID:    "grp_test",
			PropertyID: "prp_test",
			AccountID:  "act_test",
			Hostnames: papi.HostnamesResponseItems{
				Items: hostnames[:listActivePropertyHostnamesResultsPerPage],
			},
		}, nil).Times(3)
		mockListActivePropertyHostnames(client, listActivePropertyHostnamesResultsPerPage, &papi.ListActivePropertyHostnamesResponse{
			ContractID: "ctr_test",
			GroupID:    "grp_test",
			PropertyID: "prp_test",
			AccountID:  "act_test",
			Hostnames: papi.HostnamesResponseItems{
				Items: hostnames[listActivePropertyHostnamesResultsPerPage:],
			},
		}, nil).Times(3)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
					Check: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
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
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames - filter_pending_default_certs set to true", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, ptr.To("HOSTNAME_BUCKET")).Times(3)

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
				ProductionEdgeHostnameId: "ehn13",
				StagingCertType:          papi.CertTypeCPSManaged,
				StagingCnameTo:           "cnamet13",
				StagingEdgeHostnameId:    "ehn13",
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
				ProductionEdgeHostnameId: "ehn14",
				StagingCertType:          papi.CertTypeDefault,
				StagingCnameTo:           "cnamet14",
				StagingEdgeHostnameId:    "ehn14",
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
				ProductionEdgeHostnameId: "",
				StagingCertType:          papi.CertTypeDefault,
				StagingCnameTo:           "cnamet15",
				StagingEdgeHostnameId:    "ehn15",
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
				ProductionEdgeHostnameId: "ehn16",
				StagingCertType:          papi.CertTypeDefault,
				StagingCnameTo:           "cnamet16",
				StagingEdgeHostnameId:    "ehn16",
			},
		}...)

		mockListActivePropertyHostnames(client, 0, &papi.ListActivePropertyHostnamesResponse{
			ContractID: "ctr_test",
			GroupID:    "grp_test",
			PropertyID: "prp_test",
			AccountID:  "act_test",
			Hostnames: papi.HostnamesResponseItems{
				Items: hostnames,
			},
		}, nil).Times(3)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_filter.tf"),
					Check: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
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
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames with status `DEPLOYED` - filter_pending_default_certs set to true", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, ptr.To("HOSTNAME_BUCKET")).Times(3)

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
				ProductionEdgeHostnameId: "ehn1",
				StagingCertType:          papi.CertTypeDefault,
				StagingCnameTo:           "cnamet1",
				StagingEdgeHostnameId:    "ehn1",
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
				ProductionEdgeHostnameId: "ehn2",
				StagingCertType:          papi.CertTypeDefault,
				StagingCnameTo:           "cnamet2",
				StagingEdgeHostnameId:    "ehn2",
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
				ProductionEdgeHostnameId: "ehn3",
				StagingCertType:          papi.CertTypeDefault,
				StagingCnameTo:           "cnamet3",
				StagingEdgeHostnameId:    "ehn3",
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
				ProductionEdgeHostnameId: "ehn4",
				StagingCertType:          papi.CertTypeDefault,
				StagingCnameTo:           "cnamet4",
				StagingEdgeHostnameId:    "ehn4",
			},
		}

		mockListActivePropertyHostnames(client, 0, &papi.ListActivePropertyHostnamesResponse{
			ContractID: "ctr_test",
			GroupID:    "grp_test",
			PropertyID: "prp_test",
			AccountID:  "act_test",
			Hostnames: papi.HostnamesResponseItems{
				Items: hostnames,
			},
		}, nil).Times(3)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_filter.tf"),
					Check: test.NewStateChecker("data.akamai_property_hostnames.akaprophosts").
						CheckEqual("id", "prp_test").
						CheckEqual("group_id", "grp_test").
						CheckEqual("contract_id", "ctr_test").
						CheckEqual("property_id", "prp_test").
						CheckMissing("version").
						CheckEqual("hostname_bucket.#", "0").
						Build(),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames without group prefix", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, nil).Times(3)

		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
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
		client.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
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

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_no_group_prefix.tf"),
					Check:  buildAggregatedHostnamesTest(hostnameItems, "prp_test1", "test", "ctr_test", "prp_test", 1),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames without contract prefix", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, nil).Times(3)

		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
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
		client.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
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

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_no_contract_prefix.tf"),
					Check:  buildAggregatedHostnamesTest(hostnameItems, "prp_test1", "grp_test", "test", "prp_test", 1),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames without property prefix", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, nil).Times(3)

		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", testutils.MockContext, papi.GetLatestVersionRequest{
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
		client.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
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

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_no_property_prefix.tf"),
					Check:  buildAggregatedHostnamesTest(hostnameItems, "prp_test1", "grp_test", "ctr_test", "prp_test", 1),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("specify property version to fetch", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, nil).Times(3)

		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
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
		client.On("GetPropertyVersionHostnames", testutils.MockContext, papi.GetPropertyVersionHostnamesRequest{
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

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
					Check:  buildAggregatedHostnamesTest(hostnameItems, "prp_test5", "grp_test", "ctr_test", "prp_test", 5),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("specify property version to fetch with error", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, nil).Once()

		client.On("GetPropertyVersion", testutils.MockContext, papi.GetPropertyVersionRequest{
			ContractID:      "ctr_test",
			GroupID:         "grp_test",
			PropertyID:      "prp_test",
			PropertyVersion: 5,
		}).Return(nil, fmt.Errorf("error fetching property version")).Once()

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
					ExpectError: regexp.MustCompile(`error fetching property version`),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames of type HOSTNAME_BUCKET fails", func(t *testing.T) {
		client := &papi.Mock{}

		mockGetPropertyWithPropertyType(client, ptr.To("HOSTNAME_BUCKET")).Once()

		mockListActivePropertyHostnames(client, 0, nil, fmt.Errorf("error fetching list hostnames")).Once()

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
					ExpectError: regexp.MustCompile(`error fetching list hostnames`),
				}},
			})
		})

		client.AssertExpectations(t)
	})
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

func buildPropertyHostnames() []papi.Hostname {
	hostnames := make([]papi.Hostname, 10)
	for i := 0; i < 10; i++ {
		hostnames[i] = papi.Hostname{
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

func buildHostnameItems(itemsNo int) []papi.HostnameItem {
	hostnames := make([]papi.HostnameItem, 0, itemsNo)
	for i := range itemsNo {
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
			ProductionEdgeHostnameId: fmt.Sprintf("ehn%v", i),
			StagingCertType:          papi.CertTypeDefault,
			StagingCnameTo:           fmt.Sprintf("cnamet%v", i),
			StagingEdgeHostnameId:    fmt.Sprintf("ehn%v", i),
		})
	}
	return hostnames
}

func buildAggregatedHostnamesTest(hostnames []map[string]interface{}, id, groupID, contractID, propertyID string, version int) resource.TestCheckFunc {
	testVar := make([]resource.TestCheckFunc, 0)
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", "id", id))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", "group_id", groupID))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", "contract_id", contractID))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", "property_id", propertyID))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", "version", strconv.Itoa(version)))
	testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", "hostnames.#", fmt.Sprintf("%v", len(hostnames))))
	for ind, hostname := range hostnames {
		for mapKey, mapVal := range hostname {
			if mapKey != "cert_status" {
				value := fmt.Sprintf("%v", mapVal)
				key := fmt.Sprintf("hostnames.%v.%v", ind, mapKey)
				testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", key, value))
			} else {
				certStatuses := mapVal.([]map[string]interface{})
				for cInd, cert := range certStatuses {
					for cKey, cVal := range cert {
						value := fmt.Sprintf("%v", cVal)
						key := fmt.Sprintf("hostnames.%v.%v.%v.%v", ind, mapKey, cInd, cKey)
						testVar = append(testVar, resource.TestCheckResourceAttr("data.akamai_property_hostnames.akaprophosts", key, value))
					}
				}
			}
		}
	}
	return resource.ComposeAggregateTestCheckFunc(testVar...)
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
