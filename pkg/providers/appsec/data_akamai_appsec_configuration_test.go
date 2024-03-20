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

func TestAkamaiConfiguration_data_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationsResponse := appsec.GetConfigurationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)
		require.NoError(t, err)

		client.On("GetConfigurations",
			mock.Anything,
			appsec.GetConfigurationsRequest{},
		).Return(&getConfigurationsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSConfiguration/match_by_id.tf"),
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
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSConfiguration/Configuration.json"), &getConfigurationsResponse)
		require.NoError(t, err)

		client.On("GetConfigurations",
			mock.Anything,
			appsec.GetConfigurationsRequest{},
		).Return(nil, fmt.Errorf("configuration 'Nonexistent' not found"))
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		client.AssertExpectations(t)
	})
}
