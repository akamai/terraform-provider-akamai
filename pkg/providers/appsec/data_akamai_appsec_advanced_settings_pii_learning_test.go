package appsec

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func checkOutputText(value string) error {
	matched, _ := regexp.MatchString("(?s).*ENABLE PII LEARNING.*true.*", value)
	if !matched {
		return errors.New("expected result not found")
	}
	return nil
}

func TestAkamaiAdvancedSettingsPIILearning_data_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPIILearning ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getPIILearningResponse := appsec.AdvancedSettingsPIILearningResponse{}
		piiLearningBytes := testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsPIILearning/PIILearning.json")
		var piiLearningJSON bytes.Buffer
		err = json.Compact(&piiLearningJSON, []byte(piiLearningBytes))
		require.NoError(t, err)
		err = json.Unmarshal(piiLearningBytes, &getPIILearningResponse)
		require.NoError(t, err)

		client.On("GetAdvancedSettingsPIILearning",
			testutils.MockContext,
			appsec.GetAdvancedSettingsPIILearningRequest{
				ConfigVersion: appsec.ConfigVersion{
					ConfigID: 43253,
					Version:  7,
				},
			},
		).Return(&getPIILearningResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsPIILearning/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_pii_learning.test", "id", "43253"),
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_pii_learning.test", "json", piiLearningJSON.String()),
							resource.TestCheckResourceAttrWith("data.akamai_appsec_advanced_settings_pii_learning.test", "output_text", checkOutputText),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiAdvancedSettingsPIILearning_data_missing_parameter(t *testing.T) {
	t.Run("match by AdvancedSettingsPIILearning ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getPIILearningResponse := appsec.AdvancedSettingsPIILearningResponse{}
		piiLearningBytes := testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsPIILearning/PIILearning.json")
		var piiLearningJSON bytes.Buffer
		err = json.Compact(&piiLearningJSON, []byte(piiLearningBytes))
		require.NoError(t, err)
		err = json.Unmarshal(piiLearningBytes, &getPIILearningResponse)
		require.NoError(t, err)

		client.On("GetAdvancedSettingsPIILearning",
			testutils.MockContext,
			appsec.GetAdvancedSettingsPIILearningRequest{
				ConfigVersion: appsec.ConfigVersion{
					ConfigID: 43253,
					Version:  7,
				},
			},
		).Return(nil, errors.New("API call failure"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsPIILearning/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_pii_learning.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`API call failure`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
