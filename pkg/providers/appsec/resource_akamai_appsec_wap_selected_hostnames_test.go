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

func TestAccAkamaiWAPSelectedHostnames_res_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateWAPSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json")), &cu)

		cr := appsec.GetWAPSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json")), &cr)

		hns := appsec.GetWAPSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json")), &hns)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateWAPSelectedHostnames",
			mock.Anything,
			appsec.UpdateWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230",
				ProtectedHosts: []string{
					"rinaldi.sandbox.akamaideveloper.com",
				},
				EvaluatedHosts: []string{
					"sujala.sandbox.akamaideveloper.com",
				},
			},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAccAkamaiWAPSelectedHostnames_res_error_retrieving_hostnames(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateWAPSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json")), &cu)

		cr := appsec.GetWAPSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json")), &cr)

		hns := appsec.GetWAPSelectedHostnamesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json")), &hns)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetWAPSelectedHostnames failed"))

		client.On("UpdateWAPSelectedHostnames",
			mock.Anything,
			appsec.UpdateWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230",
				ProtectedHosts: []string{
					"rinaldi.sandbox.akamaideveloper.com",
				},
				EvaluatedHosts: []string{
					"sujala.sandbox.akamaideveloper.com",
				},
			},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						),
						ExpectError: regexp.MustCompile(`GetWAPSelectedHostnames failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
