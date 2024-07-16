package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiWAPSelectedHostnames_res_basic(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateWAPSelectedHostnamesResponse := appsec.UpdateWAPSelectedHostnamesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json"), &updateWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		updatedSelectedHostnamesForUpdateResponse := appsec.UpdateWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAPSelectedHostnames/WAPUpdatedSelectedHostnames.json"), &updatedSelectedHostnamesForUpdateResponse)
		require.NoError(t, err)

		updatedSelectedHostnamesForGetResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAPSelectedHostnames/WAPUpdatedSelectedHostnames.json"), &updatedSelectedHostnamesForGetResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&getWAPSelectedHostnamesResponse, nil).Times(3)

		client.On("GetWAPSelectedHostnames",
			mock.Anything,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&updatedSelectedHostnamesForGetResponse, nil)

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
		).Return(&updateWAPSelectedHostnamesResponse, nil).Once()

		client.On("UpdateWAPSelectedHostnames",
			mock.Anything,
			appsec.UpdateWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230",
				ProtectedHosts: []string{
					"test.sandbox.akamaideveloper.com",
				},
				EvaluatedHosts: []string{
					"test.evaluated.sandbox.akamaideveloper.com",
				},
			},
		).Return(&updatedSelectedHostnamesForUpdateResponse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResWAPSelectedHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "protected_hosts.0", "rinaldi.sandbox.akamaideveloper.com"),
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "evaluated_hosts.0", "sujala.sandbox.akamaideveloper.com"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResWAPSelectedHostnames/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "protected_hosts.0", "test.sandbox.akamaideveloper.com"),
							resource.TestCheckResourceAttr("akamai_appsec_wap_selected_hostnames.test", "evaluated_hosts.0", "test.evaluated.sandbox.akamaideveloper.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAkamaiWAPSelectedHostnames_res_error_retrieving_hostnames(t *testing.T) {
	t.Run("match by WAPSelectedHostnames ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateWAPSelectedHostnamesResponse := appsec.UpdateWAPSelectedHostnamesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json"), &updateWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAPSelectedHostnames/WAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

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
		).Return(&updateWAPSelectedHostnamesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResWAPSelectedHostnames/match_by_id.tf"),
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
