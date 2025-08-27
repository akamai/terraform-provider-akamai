package appsec

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestJA4FingerprintResource(t *testing.T) {
	t.Parallel()

	config := appsec.GetConfigurationResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
	require.NoError(t, err)

	advancedSettingsJA4FingerprintPriorState := appsec.GetAdvancedSettingsJA4FingerprintResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsJA4Fingerprint/JA4FingerprintBefore.json"), &advancedSettingsJA4FingerprintPriorState)
	require.NoError(t, err)

	advancedSettingsJA4FingerprintUpdatedState := appsec.GetAdvancedSettingsJA4FingerprintResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsJA4Fingerprint/JA4FingerprintAfter.json"), &advancedSettingsJA4FingerprintUpdatedState)
	require.NoError(t, err)

	advancedSettingsJA4FingerprintRemoveState := appsec.RemoveAdvancedSettingsJA4FingerprintResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsJA4Fingerprint/JA4FingerprintRemove.json"), &advancedSettingsJA4FingerprintRemoveState)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker(ja4ResourceName).
		CheckEqual("id", "111111").
		CheckEqual("config_id", "111111")

	var tests = map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"Check schema - missing required attribute config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsJA4Fingerprint/create_ja4_fingerprint_no_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"Create JA4 Fingerprint - Get configuration error": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailureJA4Fingerprint(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsJA4Fingerprint/create_ja4_fingerprint.tf"),
					ExpectError: regexp.MustCompile("Error getting configuration"),
				},
			},
		},
		"Create JA4 Fingerprint": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationJA4Fingerprint(m, config, 4)
				mockGetJA4Fingerprint(m, advancedSettingsJA4FingerprintUpdatedState, 2)
				mockUpdateJA4Fingerprint(m, []string{"ja4-fingerprint-after"}, 1)
				mockRemoveJA4Fingerprint(m, advancedSettingsJA4FingerprintRemoveState, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsJA4Fingerprint/create_ja4_fingerprint.tf"),
					Check: baseChecker.
						CheckEqual("header_names.0", "ja4-fingerprint-after").
						Build(),
				},
			},
		},
		"Import JA4 Fingerprint": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationJA4Fingerprint(m, config, 2)
				mockGetJA4Fingerprint(m, advancedSettingsJA4FingerprintPriorState, 1)
				mockRemoveJA4Fingerprint(m, advancedSettingsJA4FingerprintRemoveState, 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsJA4Fingerprint/ja4_fingerprint.tf"),
					ImportState:   true,
					ImportStateId: "111111",
					ResourceName:  ja4ResourceName,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("id", "111111").
						CheckEqual("config_id", "111111").
						Build(),
					ImportStatePersist: true,
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &appsec.Mock{}

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

var ja4ResourceName = "akamai_appsec_advanced_settings_ja4_fingerprint.test"

func mockGetConfigurationJA4Fingerprint(client *appsec.Mock, resp appsec.GetConfigurationResponse, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 111111}).
		Return(&resp, nil).Times(times)
}

func mockGetConfigurationFailureJA4Fingerprint(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 111111}).
		Return(nil, &appsec.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error getting configuration",
			StatusCode: http.StatusInternalServerError,
		}).Times(times)
}

func mockGetJA4Fingerprint(client *appsec.Mock, resp appsec.GetAdvancedSettingsJA4FingerprintResponse, times int) {
	client.On("GetAdvancedSettingsJA4Fingerprint", mock.Anything, appsec.GetAdvancedSettingsJA4FingerprintRequest{
		ConfigID: 111111,
		Version:  7,
	}).Return(&resp, nil).Times(times)
}

func mockUpdateJA4Fingerprint(client *appsec.Mock, headerNames []string, times int) {
	client.On("UpdateAdvancedSettingsJA4Fingerprint", mock.Anything, appsec.UpdateAdvancedSettingsJA4FingerprintRequest{
		ConfigID:    111111,
		Version:     7,
		HeaderNames: headerNames,
	}).Return(&appsec.UpdateAdvancedSettingsJA4FingerprintResponse{
		HeaderNames: headerNames,
	}, nil).Times(times)
}

func mockRemoveJA4Fingerprint(client *appsec.Mock, resp appsec.RemoveAdvancedSettingsJA4FingerprintResponse, times int) {
	client.On("RemoveAdvancedSettingsJA4Fingerprint", mock.Anything, appsec.RemoveAdvancedSettingsJA4FingerprintRequest{
		ConfigID: 111111,
		Version:  7,
	}).Return(&resp, nil).Times(times)
}
