package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestJA4FingerprintData(t *testing.T) {
	t.Parallel()

	config := appsec.GetConfigurationResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
	require.NoError(t, err)

	getJA4FingerprintResponse := appsec.GetAdvancedSettingsJA4FingerprintResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSAdvancedSettingsJA4Fingerprint/JA4Fingerprint.json"), &getJA4FingerprintResponse)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker(ja4DataName).
		CheckEqual("config_id", "111111")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"Return JA4 Fingerprint settings": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationJA4Fingerprint(m, config, 3)
				mockGetJA4Fingerprint(m, getJA4FingerprintResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsJA4Fingerprint/match_by_id.tf"),
					Check: baseChecker.
						CheckEqual("json", "{\"headerNames\":[\"ja4-fingerprint\"]}").
						Build(),
				},
			},
		},
		"Error response from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailureJA4Fingerprint(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsJA4Fingerprint/match_by_id.tf"),
					ExpectError: regexp.MustCompile("Error getting configuration"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &appsec.Mock{}
			t.Parallel()
			if test.init != nil {
				test.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

var ja4DataName = "data.akamai_appsec_advanced_settings_ja4_fingerprint.test"
