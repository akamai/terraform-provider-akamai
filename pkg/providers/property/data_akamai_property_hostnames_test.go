package property

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataPropertyHostnames(t *testing.T) {
	t.Skip()
	t.Run("list hostnames", func(t *testing.T) {
		client := &papi.Mock{}
		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
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
		}, nil)
		client.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
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
		}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames.tf"),
					Check:  buildAggregatedHostnamesTest(hostnameItems, "prp_test1", "grp_test", "ctr_test", "prp_test", 1),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("list hostnames without group prefix", func(t *testing.T) {
		client := &papi.Mock{}
		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
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
		}, nil)
		client.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
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
		}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
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
		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
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
		}, nil)
		client.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
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
		}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
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
		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetLatestVersion", mock.Anything, papi.GetLatestVersionRequest{
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
		}, nil)
		client.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
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
		}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
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
		hostnames := papi.HostnameResponseItems{Items: buildPropertyHostnames()}
		hostnameItems := flattenHostnames(hostnames.Items)

		client.On("GetPropertyVersion", mock.Anything, papi.GetPropertyVersionRequest{
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
		}, nil)
		client.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
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
		}, nil)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
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

		client.On("GetPropertyVersion", mock.Anything, papi.GetPropertyVersionRequest{
			ContractID:      "ctr_test",
			GroupID:         "grp_test",
			PropertyID:      "prp_test",
			PropertyVersion: 5,
		}).Return(nil, fmt.Errorf("error fetching property version"))

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnames/property_hostnames_with_version.tf"),
					ExpectError: regexp.MustCompile(`error fetching property version`),
				}},
			})
		})

		client.AssertExpectations(t)
	})
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
