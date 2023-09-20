package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiConfiguration_res_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		createConfigResponse := appsec.CreateConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationCreate.json"), &createConfigResponse)
		require.NoError(t, err)

		readConfigResponse := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &readConfigResponse)
		require.NoError(t, err)

		deleteConfigResponse := appsec.RemoveConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &deleteConfigResponse)
		require.NoError(t, err)

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationVersions.json"), &getConfigurationVersionsResponse)
		require.NoError(t, err)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSelectedHostname/SelectedHostname.json"), &getSelectedHostnamesResponse)
		require.NoError(t, err)

		client.On("GetSelectedHostnames",
			mock.Anything,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectedHostnamesResponse, nil)

		client.On("CreateConfiguration",
			mock.Anything,
			appsec.CreateConfigurationRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
		).Return(&createConfigResponse, nil)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&readConfigResponse, nil)

		client.On("RemoveConfiguration",
			mock.Anything,
			appsec.RemoveConfigurationRequest{ConfigID: 43253},
		).Return(&deleteConfigResponse, nil)

		client.On("GetConfigurationVersions",
			mock.Anything,
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAkamaiConfiguration_res_error_updating_configuration(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		createConfigResponse := appsec.CreateConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationCreate.json"), &createConfigResponse)
		require.NoError(t, err)

		readConfigResponse := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &readConfigResponse)
		require.NoError(t, err)

		deleteConfigResponse := appsec.RemoveConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &deleteConfigResponse)
		require.NoError(t, err)

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationVersions.json"), &getConfigurationVersionsResponse)
		require.NoError(t, err)

		hns := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSelectedHostname/SelectedHostname.json"), &hns)
		require.NoError(t, err)

		client.On("GetSelectedHostnames",
			mock.Anything,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&hns, nil)

		client.On("CreateConfiguration",
			mock.Anything,
			appsec.CreateConfigurationRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
		).Return(&createConfigResponse, nil)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&readConfigResponse, nil)

		client.On("UpdateConfiguration",
			mock.Anything,
			appsec.UpdateConfigurationRequest{ConfigID: 43253, Name: "Akamai Tools", Description: "Akamai Tools"},
		).Return(nil, fmt.Errorf("UpdateConfiguration failed"))

		client.On("RemoveConfiguration",
			mock.Anything,
			appsec.RemoveConfigurationRequest{ConfigID: 43253},
		).Return(&deleteConfigResponse, nil)

		client.On("GetConfigurationVersions",
			mock.Anything,
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/modify_contract.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`UpdateConfiguration failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
