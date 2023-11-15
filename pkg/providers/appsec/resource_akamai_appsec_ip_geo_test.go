package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiIPGeo_res_block(t *testing.T) {
	var (
		configVersion = func(t *testing.T, configId int, client *appsec.Mock) appsec.GetConfigurationResponse {
			configResponse := appsec.GetConfigurationResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
			require.NoError(t, err)

			client.On("GetConfiguration",
				mock.Anything,
				appsec.GetConfigurationRequest{ConfigID: configId},
			).Return(&configResponse, nil)

			return configResponse
		}

		getIPGeoResponse = func(t *testing.T, configId int, version int, policyId string, path string, client *appsec.Mock) appsec.GetIPGeoResponse {
			getIPGeoResponse := appsec.GetIPGeoResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &getIPGeoResponse)
			require.NoError(t, err)

			client.On("GetIPGeo",
				mock.Anything,
				appsec.GetIPGeoRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&getIPGeoResponse, nil)

			return getIPGeoResponse
		}

		updateIPGeoResponse = func(t *testing.T, configId int, version int, policyId string, path string, request appsec.UpdateIPGeoRequest, client *appsec.Mock) appsec.UpdateIPGeoResponse {
			updateIPGeoResponse := appsec.UpdateIPGeoResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &updateIPGeoResponse)
			require.NoError(t, err)

			client.On("UpdateIPGeo",
				mock.Anything,
				request,
			).Return(&updateIPGeoResponse, nil)
			return updateIPGeoResponse
		}

		updateIPGeoProtectionResponseAllProtectionsFalse = func(t *testing.T, configId int, version int, policyId string, path string, client *appsec.Mock) appsec.UpdateIPGeoProtectionResponse {
			updateIPGeoProtectionResponseAllProtectionsFalse := appsec.UpdateIPGeoProtectionResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &updateIPGeoProtectionResponseAllProtectionsFalse)
			require.NoError(t, err)

			client.On("UpdateIPGeoProtection",
				mock.Anything,
				appsec.UpdateIPGeoProtectionRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&updateIPGeoProtectionResponseAllProtectionsFalse, nil).Once()
			return updateIPGeoProtectionResponseAllProtectionsFalse
		}
	)
	t.Run("match by IPGeo ID", func(t *testing.T) {
		client := &appsec.Mock{}
		configVersion(t, 43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
			ASNControls: &appsec.IPGeoASNControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"40721_ASNLIST1",
						"44811_ASNLIST2",
					},
				},
			},
			GeoControls: &appsec.IPGeoGeoControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"40731_BMROLLOUTGEO",
						"44831_ECSCGEOBLACKLIST",
					},
				},
			},
			IPControls: &appsec.IPGeoIPControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"49181_ADTIPBLACKLIST",
						"49185_ADTWAFBYPASSLIST",
					},
				},
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"68762_ADYEN",
						"69601_ADYENPRODWHITELIST",
					},
				},
			},
		}
		getIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeo.json", client)
		updateIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeo.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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

		configVersion(t, 43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
			ASNControls: &appsec.IPGeoASNControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"40721_ASNLIST1",
						"44811_ASNLIST2",
					},
				},
			},
			GeoControls: &appsec.IPGeoGeoControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"40731_BMROLLOUTGEO",
						"44831_ECSCGEOBLACKLIST",
					},
				},
			},
			IPControls: &appsec.IPGeoIPControls{
				BlockedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"49181_ADTIPBLACKLIST",
						"49185_ADTWAFBYPASSLIST",
					},
				},
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"68762_ADYEN",
						"69601_ADYENPRODWHITELIST",
					},
				},
			},
			UkraineGeoControls: &appsec.UkraineGeoControl{
				Action: "alert",
			},
		}
		getIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeo.json", client)
		updateIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeo.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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

	t.Run("match by IPGeo Allow ID", func(t *testing.T) {

		client := &appsec.Mock{}

		configVersion(t, 43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockAllTrafficExceptAllowedIPs",
			IPControls: &appsec.IPGeoIPControls{
				AllowedIPNetworkLists: &appsec.IPGeoNetworkLists{
					NetworkList: []string{
						"68762_ADYEN",
						"69601_ADYENPRODWHITELIST",
					},
				},
			},
		}
		getIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoAllow.json", client)
		updateIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoAllow.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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

		configVersion(t, 43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockSpecificIPGeo",
		}
		getIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoBlockOnly.json", client)
		updateIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoBlockOnly.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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

		configVersion(t, 43253, client)

		updateRequest := appsec.UpdateIPGeoRequest{
			ConfigID: 43253,
			Version:  7,
			PolicyID: "AAAA_81230",
			Block:    "blockAllTrafficExceptAllowedIPs",
		}

		getIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoAllowOnly.json", client)
		updateIPGeoResponse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeo/IPGeoAllowOnly.json", updateRequest, client)
		updateIPGeoProtectionResponseAllProtectionsFalse(t, 43253, 7, "AAAA_81230", "testdata/TestResIPGeoProtection/PolicyProtections.json", client)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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
}
