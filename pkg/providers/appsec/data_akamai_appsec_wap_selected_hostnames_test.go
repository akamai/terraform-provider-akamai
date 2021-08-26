package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiWAPSelectedHostnames_data_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfigurationWAP.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		cv := appsec.GetWAPSelectedHostnamesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSWAPSelectedHostnames/WAPSelectedHostnames.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
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

func TestAccAkamaiWAPSelectedHostnames_data_error_retrieving_hostnames(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfigurationWAP.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		cv := appsec.GetWAPSelectedHostnamesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSWAPSelectedHostnames/WAPSelectedHostnames.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetSelectedHostnames failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
						ExpectError: regexp.MustCompile(`GetSelectedHostnames failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAccAkamaiWAPSelectedHostnames_NonWAP_data_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		expectedGetSelectedHostnamesResponseJSON := compactJSON(loadFixtureBytes("testdata/TestDSSelectedHostnames/SelectedHostnames.json"))
		json.Unmarshal([]byte(expectedGetSelectedHostnamesResponseJSON), &getSelectedHostnamesResponse)

		getMatchTargetsResponse := appsec.GetMatchTargetsResponse{}
		expectedGetMatchTargetsResponseJSON := compactJSON(loadFixtureBytes("testdata/TestDSMatchTargets/MatchTargets.json"))
		json.Unmarshal([]byte(expectedGetMatchTargetsResponseJSON), &getMatchTargetsResponse)

		client.On("GetSelectedHostnames",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectedHostnamesResponse, nil)

		client.On("GetMatchTargets",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetMatchTargetsRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&getMatchTargetsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAPSelectedHostnames/match_by_id.tf"),
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
