package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAAPSelectedHostnames_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by AAPSelectedHostnames ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateWAPSelectedHostnamesResponse := appsec.UpdateWAPSelectedHostnamesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAAPSelectedHostnames/AAPSelectedHostnames.json"), &updateWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAAPSelectedHostnames/AAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		updatedSelectedHostnamesForUpdateResponse := appsec.UpdateWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAAPSelectedHostnames/AAPUpdatedSelectedHostnames.json"), &updatedSelectedHostnamesForUpdateResponse)
		require.NoError(t, err)

		updatedSelectedHostnamesForGetResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAAPSelectedHostnames/AAPUpdatedSelectedHostnames.json"), &updatedSelectedHostnamesForGetResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetWAPSelectedHostnames",
			testutils.MockContext,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&getWAPSelectedHostnamesResponse, nil).Times(3)

		client.APPSEC.On("GetWAPSelectedHostnames",
			testutils.MockContext,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&updatedSelectedHostnamesForGetResponse, nil)

		client.APPSEC.On("UpdateWAPSelectedHostnames",
			testutils.MockContext,
			appsec.UpdateWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230",
				ProtectedHosts: []string{
					"rinaldi.sandbox.akamaideveloper.com",
				},
				EvaluatedHosts: []string{
					"sujala.sandbox.akamaideveloper.com",
				},
			},
		).Return(&updateWAPSelectedHostnamesResponse, nil).Once()

		client.APPSEC.On("UpdateWAPSelectedHostnames",
			testutils.MockContext,
			appsec.UpdateWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230",
				ProtectedHosts: []string{
					"test.sandbox.akamaideveloper.com",
				},
				EvaluatedHosts: []string{
					"test.evaluated.sandbox.akamaideveloper.com",
				},
			},
		).Return(&updatedSelectedHostnamesForUpdateResponse, nil).Once()

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAAPSelectedHostnames/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_aap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_appsec_aap_selected_hostnames.test", "protected_hosts.0", "rinaldi.sandbox.akamaideveloper.com"),
						resource.TestCheckResourceAttr("akamai_appsec_aap_selected_hostnames.test", "evaluated_hosts.0", "sujala.sandbox.akamaideveloper.com"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAAPSelectedHostnames/update_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_aap_selected_hostnames.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_appsec_aap_selected_hostnames.test", "protected_hosts.0", "test.sandbox.akamaideveloper.com"),
						resource.TestCheckResourceAttr("akamai_appsec_aap_selected_hostnames.test", "evaluated_hosts.0", "test.evaluated.sandbox.akamaideveloper.com"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("match by AAPSelectedHostnames ID - error retrieving hostnames", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateWAPSelectedHostnamesResponse := appsec.UpdateWAPSelectedHostnamesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAAPSelectedHostnames/AAPSelectedHostnames.json"), &updateWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAAPSelectedHostnames/AAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetWAPSelectedHostnames",
			testutils.MockContext,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetWAPSelectedHostnames failed"))

		client.APPSEC.On("UpdateWAPSelectedHostnames",
			testutils.MockContext,
			appsec.UpdateWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230",
				ProtectedHosts: []string{
					"rinaldi.sandbox.akamaideveloper.com",
				},
				EvaluatedHosts: []string{
					"sujala.sandbox.akamaideveloper.com",
				},
			},
		).Return(&updateWAPSelectedHostnamesResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAAPSelectedHostnames/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_aap_selected_hostnames.test", "id", "43253:AAAA_81230"),
					),
					ExpectError: regexp.MustCompile(`GetWAPSelectedHostnames failed`),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})
}
