package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiIPGeo_res_block(t *testing.T) {
	var (
		configVersion = func(configId int, client *appsec.Mock) appsec.GetConfigurationResponse {
			configResponse := appsec.GetConfigurationResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
			require.NoError(t, err)

			client.On("GetConfiguration",
				testutils.MockContext,
				appsec.GetConfigurationRequest{ConfigID: configId},
			).Return(&configResponse, nil)

			return configResponse
		}

		getIPGeoResponse = func(configId int, version int, policyId string, path string, client *appsec.Mock) appsec.GetIPGeoResponse {
			getIPGeoResponse := appsec.GetIPGeoResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &getIPGeoResponse)
			require.NoError(t, err)

			client.On("GetIPGeo",
				testutils.MockContext,
				appsec.GetIPGeoRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&getIPGeoResponse, nil)

			return getIPGeoResponse
		}

		updateIPGeoResponse = func(path string, request appsec.UpdateIPGeoRequest, client *appsec.Mock) appsec.UpdateIPGeoResponse {
			updateIPGeoResponse := appsec.UpdateIPGeoResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &updateIPGeoResponse)
			require.NoError(t, err)

			client.On("UpdateIPGeo",
				testutils.MockContext,
				request,
			).Return(&updateIPGeoResponse, nil)
			return updateIPGeoResponse
		}

		updateIPGeoProtectionResponseAllProtectionsFalse = func(configId int, version int, policyId string, path string, client *appsec.Mock) appsec.UpdateIPGeoProtectionResponse {
			updateIPGeoProtectionResponseAllProtectionsFalse := appsec.UpdateIPGeoProtectionResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &updateIPGeoProtectionResponseAllProtectionsFalse)
			require.NoError(t, err)

			client.On("UpdateIPGeoProtection",
				testutils.MockContext,
				appsec.UpdateIPGeoProtectionRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&updateIPGeoProtectionResponseAllProtectionsFalse, nil).Once()
			return updateIPGeoProtectionResponseAllProtectionsFalse
		}
	)
	t.Run("match by IPGeo ID", func(t *testing.T) {
		client := &appsec.Mock{}
		configVersion(43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
			ASNControls: &appsec.IPGeoASNControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"44811_ASNLIST2",
						"40721_ASNLIST1",
					},
				},
			},
			GeoControls: &appsec.IPGeoGeoControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"40731_BMROLLOUTGEO",
						"44831_ECSCGEOBLACKLIST",
					},
				},
			},
			IPControls: &appsec.IPGeoIPControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"49185_ADTWAFBYPASSLIST",
						"49181_ADTIPBLACKLIST",
					},
				},
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"69601_ADYENPRODWHITELIST",
						"68762_ADYEN",
					},
				},
			},
		}
		getIPGeoResponse(43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeo.json", client)
		updateIPGeoResponse("testdata/TestResIPGeo/IPGeo.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by Ukraine Geo ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configVersion(43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
			ASNControls: &appsec.IPGeoASNControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"44811_ASNLIST2",
						"40721_ASNLIST1",
					},
				},
			},
			GeoControls: &appsec.IPGeoGeoControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"40731_BMROLLOUTGEO",
						"44831_ECSCGEOBLACKLIST",
					},
				},
			},
			IPControls: &appsec.IPGeoIPControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"49185_ADTWAFBYPASSLIST",
						"49181_ADTIPBLACKLIST",
					},
				},
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"69601_ADYENPRODWHITELIST",
						"68762_ADYEN",
					},
				},
			},
			UkraineGeoControls: &appsec.UkraineGeoControl{
				Action: "alert",
			},
		}
		getIPGeoResponse(43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeo.json", client)
		updateIPGeoResponse("testdata/TestResIPGeo/IPGeo.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/ukraine_match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test1", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by Ukraine Geo ID field not included in config", func(t *testing.T) {
		client := &appsec.Mock{}

		configVersion(43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
			ASNControls: &appsec.IPGeoASNControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"44811_ASNLIST2",
						"40721_ASNLIST1",
					},
				},
			},
			GeoControls: &appsec.IPGeoGeoControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"40731_BMROLLOUTGEO",
						"44831_ECSCGEOBLACKLIST",
					},
				},
			},
			IPControls: &appsec.IPGeoIPControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"49185_ADTWAFBYPASSLIST",
						"49181_ADTIPBLACKLIST",
					},
				},
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"69601_ADYENPRODWHITELIST",
						"68762_ADYEN",
					},
				},
			},
		}
		getIPGeoResponse(43253, 7, "AAAA_81230", "testdata/TestResIPGeo/UkraineGeo.json", client)
		updateIPGeoResponse("testdata/TestResIPGeo/UkraineGeo.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by Ukraine Geo ID field included with same action", func(t *testing.T) {
		client := &appsec.Mock{}

		configVersion(43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
			ASNControls: &appsec.IPGeoASNControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"44811_ASNLIST2",
						"40721_ASNLIST1",
					},
				},
			},
			GeoControls: &appsec.IPGeoGeoControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"40731_BMROLLOUTGEO",
						"44831_ECSCGEOBLACKLIST",
					},
				},
			},
			IPControls: &appsec.IPGeoIPControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					Action: "deny",
					NetworkList: []string{
						"49185_ADTWAFBYPASSLIST",
						"49181_ADTIPBLACKLIST",
					},
				},
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"69601_ADYENPRODWHITELIST",
						"68762_ADYEN",
					},
				},
			},
			UkraineGeoControls: &appsec.UkraineGeoControl{
				Action: "alert",
			},
		}
		getIPGeoResponse(43253, 7, "AAAA_81230", "testdata/TestResIPGeo/UkraineGeo.json", client)
		updateIPGeoResponse("testdata/TestResIPGeo/UkraineGeo.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/ukraine_match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test1", "id", "43253:AAAA_81230"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/ukraine_match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test1", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by IPGeo Allow ID", func(t *testing.T) {

		client := &appsec.Mock{}

		configVersion(43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID:       43253,
			Version:        7,
			PolicyID:       "AAAA_81230",
			Block:          "blockAllTrafficExceptAllowedIPs",
			BlockAllAction: "deny",
			IPControls: &appsec.IPGeoIPControls{
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"69601_ADYENPRODWHITELIST",
						"68762_ADYEN",
					},
				},
			},
		}
		getIPGeoResponse(43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoAllow.json", client)
		updateIPGeoResponse("testdata/TestResIPGeo/IPGeoAllow.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/allow.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("block with empty lists", func(t *testing.T) {

		client := &appsec.Mock{}
		configVersion(43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
		}
		getIPGeoResponse(43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoBlockOnly.json", client)
		updateIPGeoResponse("testdata/TestResIPGeo/IPGeoBlockOnly.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/block_with_empty_lists.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("allow with empty lists", func(t *testing.T) {

		client := &appsec.Mock{}

		configVersion(43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID:       43253,
			Version:        7,
			PolicyID:       "AAAA_81230",
			BlockAllAction: "deny",
			Block:          "blockAllTrafficExceptAllowedIPs",
		}

		getIPGeoResponse(43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoAllowOnly.json", client)
		updateIPGeoResponse("testdata/TestResIPGeo/IPGeoAllowOnly.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/allow_with_empty_lists.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("empty string in Geo network lists input", func(t *testing.T) {
		client := &appsec.Mock{}
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/mode_allow_with_empty_geo_network_lists.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
						ExpectError: regexp.MustCompile("Error: empty or invalid string value for config parameter geo_controls"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("empty string input in all lists", func(t *testing.T) {
		client := &appsec.Mock{}
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResIPGeo/all_lists_with_empty_string_input.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
						ExpectError: regexp.MustCompile(`(?s)Error: empty or invalid string value for config parameter geo_controls.*Error: empty or invalid string value for config parameter ip_controls.*Error: empty or invalid string value for config parameter exception_ip_network_lists.*`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}
