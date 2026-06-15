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

func TestAkamaiConfiguration_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by Configuration ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getConfigurationsResponse := appsec.GetConfigurationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetConfigurations",
			testutils.MockContext,
			appsec.GetConfigurationsRequest{},
		).Return(&getConfigurationsResponse, nil)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfiguration/SelectedHostnames.json"), &getSelectedHostnamesResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetSelectedHostnames",
			testutils.MockContext,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 15},
		).Return(&getSelectedHostnamesResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSConfiguration/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_configuration.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})
}

func TestAkamaiConfiguration_data_hostnames(t *testing.T) {
	t.Parallel()
	t.Run("match by selected hostnames", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getConfigurationsResponse := appsec.GetConfigurationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetConfigurations",
			testutils.MockContext,
			appsec.GetConfigurationsRequest{},
		).Return(&getConfigurationsResponse, nil)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfiguration/SelectedHostnames.json"), &getSelectedHostnamesResponse)
		require.NoError(t, err)

		expectedOutputText := "\n+------------------------------------------------------------------------------------------------------+\n| Configurations                                                                                       |\n+-----------+--------------+----------------+---------------------------+------------------------------+\n| CONFIG_ID | NAME         | LATEST_VERSION | VERSION_ACTIVE_IN_STAGING | VERSION_ACTIVE_IN_PRODUCTION |\n+-----------+--------------+----------------+---------------------------+------------------------------+\n| 43253     | Akamai Tools | 15             | 0                         | 0                            |\n+-----------+--------------+----------------+---------------------------+------------------------------+\n"
		client.APPSEC.On("GetSelectedHostnames",
			testutils.MockContext,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 15},
		).Return(&getSelectedHostnamesResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSConfiguration/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_configuration.test", "host_names.#", "2"),
						resource.TestCheckResourceAttr("data.akamai_appsec_configuration.test", "output_text", expectedOutputText),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})
}

func TestAkamaiConfiguration_data_nonexistentConfig(t *testing.T) {
	t.Parallel()
	t.Run("nonexistent configuration", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getConfigurationsResponse := appsec.GetConfigurationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetConfigurations",
			testutils.MockContext,
			appsec.GetConfigurationsRequest{},
		).Return(nil, fmt.Errorf("configuration 'Nonexistent' not found"))
		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSConfiguration/nonexistent_config.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_configuration.test", "id", "43253"),
					),
					ExpectError: regexp.MustCompile(`configuration 'Nonexistent' not found`),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})
}
