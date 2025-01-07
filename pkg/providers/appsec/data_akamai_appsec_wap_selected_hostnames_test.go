package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiWAPSelectedHostnames_data_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfigurationWAP.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSWAPSelectedHostnames/WAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		client.On("GetWAPSelectedHostnames",
			testutils.MockContext,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&getWAPSelectedHostnamesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiWAPSelectedHostnames_data_error_retrieving_hostnames(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfigurationWAP.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSWAPSelectedHostnames/WAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		client.On("GetWAPSelectedHostnames",
			testutils.MockContext,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetWAPSelectedHostnames failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
						ExpectError: regexp.MustCompile(`GetWAPSelectedHostnames failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiWAPSelectedHostnames_NonWAP_data_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSSelectedHostnames/SelectedHostnames.json"), &getSelectedHostnamesResponse)
		require.NoError(t, err)

		getMatchTargetsResponse := appsec.GetMatchTargetsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSMatchTargets/MatchTargets.json"), &getMatchTargetsResponse)
		require.NoError(t, err)

		client.On("GetSelectedHostnames",
			testutils.MockContext,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectedHostnamesResponse, nil)

		client.On("GetMatchTargets",
			testutils.MockContext,
			appsec.GetMatchTargetsRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&getMatchTargetsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
