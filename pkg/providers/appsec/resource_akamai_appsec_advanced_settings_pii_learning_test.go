package appsec

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsPIILearning_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPIILearning", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)

		getResponse := appsec.AdvancedSettingsPIILearningResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsPIILearning/PIILearning.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := appsec.AdvancedSettingsPIILearningResponse{}
		piiLearningBytes := testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsPIILearning/PIILearning.json")
		var piiLearningJSON bytes.Buffer
		err = json.Compact(&piiLearningJSON, []byte(piiLearningBytes))
		require.NoError(t, err)
		err = json.Unmarshal(piiLearningBytes, &updateResponse)
		require.NoError(t, err)

		removeResponse := appsec.AdvancedSettingsPIILearningResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsPIILearning/PIILearning.json"), &removeResponse)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.On("UpdateAdvancedSettingsPIILearning",
			testutils.MockContext,
			appsec.UpdateAdvancedSettingsPIILearningRequest{
				ConfigVersion: appsec.ConfigVersion{
					ConfigID: 43253,
					Version:  7,
				},
				EnablePIILearning: true},
		).Return(&updateResponse, nil)

		client.On("GetAdvancedSettingsPIILearning",
			testutils.MockContext,
			appsec.GetAdvancedSettingsPIILearningRequest{
				ConfigVersion: appsec.ConfigVersion{
					ConfigID: 43253,
					Version:  7},
			},
		).Return(&getResponse, nil)

		client.On("UpdateAdvancedSettingsPIILearning",
			testutils.MockContext,
			appsec.UpdateAdvancedSettingsPIILearningRequest{
				ConfigVersion: appsec.ConfigVersion{
					ConfigID: 43253,
					Version:  7,
				},
			},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsPIILearning/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_pii_learning.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_pii_learning.test", "enable_pii_learning", "true"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiAdvancedSettingsPIILearning_res_api_call_failure(t *testing.T) {
	t.Run("match by AdvancedSettingsPIILearning", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)

		getResponse := appsec.AdvancedSettingsPIILearningResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsPIILearning/PIILearning.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := appsec.AdvancedSettingsPIILearningResponse{}
		piiLearningBytes := testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsPIILearning/PIILearning.json")
		var piiLearningJSON bytes.Buffer
		err = json.Compact(&piiLearningJSON, []byte(piiLearningBytes))
		require.NoError(t, err)
		err = json.Unmarshal(piiLearningBytes, &updateResponse)
		require.NoError(t, err)

		removeResponse := appsec.AdvancedSettingsPIILearningResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsPIILearning/PIILearning.json"), &removeResponse)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.On("UpdateAdvancedSettingsPIILearning",
			testutils.MockContext,
			appsec.UpdateAdvancedSettingsPIILearningRequest{
				ConfigVersion: appsec.ConfigVersion{
					ConfigID: 43253,
					Version:  7,
				},
				EnablePIILearning: true},
		).Return(nil, errors.New("API call failure"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsPIILearning/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_pii_learning.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`API call failure`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
