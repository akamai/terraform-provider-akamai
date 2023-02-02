package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiConfiguration_data_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationsResponse := appsec.GetConfigurationsResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)
		require.NoError(t, err)

		client.On("GetConfigurations",
			mock.Anything,
			appsec.GetConfigurationsRequest{},
		).Return(&getConfigurationsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_configuration.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiConfiguration_data_nonexistentConfig(t *testing.T) {
	t.Run("nonexistent configuration", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationsResponse := appsec.GetConfigurationsResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)
		require.NoError(t, err)

		client.On("GetConfigurations",
			mock.Anything,
			appsec.GetConfigurationsRequest{},
		).Return(nil, fmt.Errorf("configuration 'Nonexistent' not found"))
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSConfiguration/nonexistent_config.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_configuration.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`configuration 'Nonexistent' not found`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
