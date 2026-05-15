package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAAPSelectedHostnames_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by AAPSelectedHostnames ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfigurationAAP.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAAPSelectedHostnames/AAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetWAPSelectedHostnames",
			testutils.MockContext,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(&getWAPSelectedHostnamesResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSAAPSelectedHostnames/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_aap_selected_hostnames.test", "id", "43253:AAAA_81230"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("match by AAPSelectedHostnames ID - expect error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfigurationAAP.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getWAPSelectedHostnamesResponse := appsec.GetWAPSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAAPSelectedHostnames/AAPSelectedHostnames.json"), &getWAPSelectedHostnamesResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetWAPSelectedHostnames",
			testutils.MockContext,
			appsec.GetWAPSelectedHostnamesRequest{ConfigID: 43253, Version: 7, SecurityPolicyID: "AAAA_81230"},
		).Return(nil, fmt.Errorf("GetWAPSelectedHostnames failed"))

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSAAPSelectedHostnames/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_aap_selected_hostnames.test", "id", "43253:AAAA_81230"),
					),
					ExpectError: regexp.MustCompile(`GetWAPSelectedHostnames failed`),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("match by AAPSelectedHostnames ID - error retrieving hostnames", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAAPSelectedHostnames/AAPSelectedHostnames.json"), &getSelectedHostnamesResponse)
		require.NoError(t, err)

		getMatchTargetsResponse := appsec.GetMatchTargetsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSMatchTargets/MatchTargets.json"), &getMatchTargetsResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetSelectedHostnames",
			testutils.MockContext,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectedHostnamesResponse, nil)

		client.APPSEC.On("GetMatchTargets",
			testutils.MockContext,
			appsec.GetMatchTargetsRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&getMatchTargetsResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSAAPSelectedHostnames/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_aap_selected_hostnames.test", "id", "43253:AAAA_81230"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
