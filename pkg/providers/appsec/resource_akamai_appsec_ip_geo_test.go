package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiIPGeo_res_block(t *testing.T) {
	t.Run("match by IPGeo ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &getConfigurationResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigurationResponse, nil)

		updateIPGeoResponse := appsec.UpdateIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeo.json"), &updateIPGeoResponse)
		require.NoError(t, err)
		client.On("UpdateIPGeo",
			mock.Anything,
			appsec.UpdateIPGeoRequest{
				ConfigID: 43253,
				Version:  7,
				PolicyID: "AAAA_81230",
				Block:    "blockSpecificIPGeo",
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
			},
		).Return(&updateIPGeoResponse, nil)

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeo.json"), &getIPGeoResponse)
		require.NoError(t, err)
		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		updateIPGeoProtectionResponseAllProtectionsFalse := appsec.UpdateIPGeoProtectionResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeoProtection/PolicyProtections.json"), &updateIPGeoProtectionResponseAllProtectionsFalse)
		require.NoError(t, err)
		client.On("UpdateIPGeoProtection",
			mock.Anything,
			appsec.UpdateIPGeoProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateIPGeoProtectionResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResIPGeo/match_by_id.tf"),
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

func TestAkamaiIPGeo_res_allow(t *testing.T) {
	t.Run("match by IPGeo ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &getConfigurationResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigurationResponse, nil)

		updateIPGeoResponse := appsec.UpdateIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeoAllow.json"), &updateIPGeoResponse)
		require.NoError(t, err)
		client.On("UpdateIPGeo",
			mock.Anything,
			appsec.UpdateIPGeoRequest{
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
			},
		).Return(&updateIPGeoResponse, nil)

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeoAllow.json"), &getIPGeoResponse)
		require.NoError(t, err)
		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		updateIPGeoProtectionResponseAllProtectionsFalse := appsec.UpdateIPGeoProtectionResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeoProtection/PolicyProtections.json"), &updateIPGeoProtectionResponseAllProtectionsFalse)
		require.NoError(t, err)
		client.On("UpdateIPGeoProtection",
			mock.Anything,
			appsec.UpdateIPGeoProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateIPGeoProtectionResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResIPGeo/allow.tf"),
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

func TestAkamaiIPGeo_res_block_with_empty_lists(t *testing.T) {
	t.Run("block with empty lists", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &getConfigurationResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigurationResponse, nil)

		updateIPGeoResponse := appsec.UpdateIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeoBlockOnly.json"), &updateIPGeoResponse)
		require.NoError(t, err)
		client.On("UpdateIPGeo",
			mock.Anything,
			appsec.UpdateIPGeoRequest{
				ConfigID: 43253,
				Version:  7,
				PolicyID: "AAAA_81230",
				Block:    "blockSpecificIPGeo",
			},
		).Return(&updateIPGeoResponse, nil)

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeoBlockOnly.json"), &getIPGeoResponse)
		require.NoError(t, err)
		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		updateIPGeoProtectionResponseAllProtectionsFalse := appsec.UpdateIPGeoProtectionResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeoProtection/PolicyProtections.json"), &updateIPGeoProtectionResponseAllProtectionsFalse)
		require.NoError(t, err)
		client.On("UpdateIPGeoProtection",
			mock.Anything,
			appsec.UpdateIPGeoProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateIPGeoProtectionResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResIPGeo/block_with_empty_lists.tf"),
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

func TestAkamaiIPGeo_res_allow_with_empty_lists(t *testing.T) {
	t.Run("allow with empty lists", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &getConfigurationResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigurationResponse, nil)

		updateIPGeoResponse := appsec.UpdateIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeoAllowOnly.json"), &updateIPGeoResponse)
		require.NoError(t, err)
		client.On("UpdateIPGeo",
			mock.Anything,
			appsec.UpdateIPGeoRequest{
				ConfigID: 43253,
				Version:  7,
				PolicyID: "AAAA_81230",
				Block:    "blockAllTrafficExceptAllowedIPs",
			},
		).Return(&updateIPGeoResponse, nil)

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeo/IPGeoAllowOnly.json"), &getIPGeoResponse)
		require.NoError(t, err)
		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		updateIPGeoProtectionResponseAllProtectionsFalse := appsec.UpdateIPGeoProtectionResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResIPGeoProtection/PolicyProtections.json"), &updateIPGeoProtectionResponseAllProtectionsFalse)
		require.NoError(t, err)
		client.On("UpdateIPGeoProtection",
			mock.Anything,
			appsec.UpdateIPGeoProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateIPGeoProtectionResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResIPGeo/allow_with_empty_lists.tf"),
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
