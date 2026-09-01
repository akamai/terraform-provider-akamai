package appsec

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataAdvancedSettingsURLEvasionDefense(t *testing.T) {
	t.Parallel()

	baseConfig := testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsURLEvasionDefense/url_evasion_defense.tf")
	missingConfigID := testutils.LoadFixtureString(t, "testdata/TestDSAdvancedSettingsURLEvasionDefense/url_evasion_defense_missing_config_id.tf")

	operator := string(appsec.AdvancedSettingsURLEvasionDefenseConditionOperatorOr)
	response := appsec.GetAdvancedSettingsURLEvasionDefenseResponse{
		Status:      string(appsec.AdvancedSettingsURLEvasionDefenseStatusEnabled),
		LockVersion: 3,
		BypassLists: []string{"123_TEST", "456_TEST"},
		Rules: []appsec.AdvancedSettingsURLEvasionDefenseRuleResponse{
			{
				RuleID:            950001,
				ConditionOperator: &operator,
				Conditions: []appsec.AdvancedSettingsURLEvasionDefenseRuleCondition{
					{Type: "extensionMatch", Extensions: []string{"php", "jsp"}, PositiveMatch: true},
					{Type: "filenameMatch", Filenames: []string{"index.php"}, PositiveMatch: true},
					{Type: "hostMatch", Hosts: []string{"example.com"}, PositiveMatch: true},
					{Type: "pathMatch", Paths: []string{"/admin/*"}, PositiveMatch: true},
					{Type: "requestMethodMatch", Methods: []string{"GET", "POST"}, PositiveMatch: true},
					{Type: "ipMatch", IPs: []string{"1.2.3.4"}, UseHeaders: true, PositiveMatch: true},
					{Type: "clientListMatch", ClientLists: []string{"123_TEST"}, UseHeaders: true, PositiveMatch: true},
					{Type: "requestHeaderMatch", Header: "X-Test", Value: "abc*", ValueCase: true, ValueWildcard: true, PositiveMatch: true},
					{Type: "uriQueryMatch", Name: "param", Value: "val*", NameCase: true, CaseSensitive: true, Wildcard: true, PositiveMatch: true},
				},
				Action:      ptr.To("deny"),
				Name:        "Rule 950001",
				Description: "Test rule",
			},
		},
	}

	baseChecker := test.NewStateChecker("data.akamai_appsec_advanced_settings_url_evasion_defense.test")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return URL evasion defense": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationForURLDefenseDS(m, 43253, 7)
				m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, appsec.GetAdvancedSettingsURLEvasionDefenseRequest{ConfigID: 43253, Version: 7}).
					Return(&response, nil)
			},
			steps: []resource.TestStep{
				{
					Config: baseConfig,
					Check: baseChecker.
						CheckEqual("config_id", "43253").
						CheckEqual("status", "enabled").
						CheckEqual("bypass_lists.#", "2").
						CheckEqual("rules.#", "1").
						CheckEqual("rules.0.rule_id", "950001").
						CheckEqual("rules.0.action", "deny").
						CheckEqual("rules.0.name", "Rule 950001").
						CheckEqual("rules.0.description", "Test rule").
						CheckEqual("rules.0.condition_operator", "OR").
						CheckEqual("rules.0.conditions.#", "9").
						CheckEqual("rules.0.conditions.0.type", "extensionMatch").
						CheckEqual("rules.0.conditions.0.extensions.#", "2").
						CheckEqual("rules.0.conditions.0.positive_match", "true").
						CheckEqual("rules.0.conditions.1.type", "filenameMatch").
						CheckEqual("rules.0.conditions.1.filenames.#", "1").
						CheckEqual("rules.0.conditions.1.positive_match", "true").
						CheckEqual("rules.0.conditions.2.type", "hostMatch").
						CheckEqual("rules.0.conditions.2.hosts.#", "1").
						CheckEqual("rules.0.conditions.2.positive_match", "true").
						CheckEqual("rules.0.conditions.3.type", "pathMatch").
						CheckEqual("rules.0.conditions.3.paths.#", "1").
						CheckEqual("rules.0.conditions.3.positive_match", "true").
						CheckEqual("rules.0.conditions.4.type", "requestMethodMatch").
						CheckEqual("rules.0.conditions.4.methods.#", "2").
						CheckEqual("rules.0.conditions.4.positive_match", "true").
						CheckEqual("rules.0.conditions.5.type", "ipMatch").
						CheckEqual("rules.0.conditions.5.ips.#", "1").
						CheckEqual("rules.0.conditions.5.use_headers", "true").
						CheckEqual("rules.0.conditions.5.positive_match", "true").
						CheckEqual("rules.0.conditions.6.type", "clientListMatch").
						CheckEqual("rules.0.conditions.6.client_lists.#", "1").
						CheckEqual("rules.0.conditions.6.use_headers", "true").
						CheckEqual("rules.0.conditions.6.positive_match", "true").
						CheckEqual("rules.0.conditions.7.type", "requestHeaderMatch").
						CheckEqual("rules.0.conditions.7.header", "X-Test").
						CheckEqual("rules.0.conditions.7.value", "abc*").
						CheckEqual("rules.0.conditions.7.value_case_sensitive", "true").
						CheckEqual("rules.0.conditions.7.value_wildcard", "true").
						CheckEqual("rules.0.conditions.7.positive_match", "true").
						CheckEqual("rules.0.conditions.8.type", "uriQueryMatch").
						CheckEqual("rules.0.conditions.8.name", "param").
						CheckEqual("rules.0.conditions.8.value", "val*").
						CheckEqual("rules.0.conditions.8.name_case_sensitive", "true").
						CheckEqual("rules.0.conditions.8.value_case_sensitive", "true").
						CheckEqual("rules.0.conditions.8.value_wildcard", "true").
						CheckEqual("rules.0.conditions.8.positive_match", "true").
						Build(),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      missingConfigID,
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"error response from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailureForURLDefenseDS(m, 43253)
			},
			steps: []resource.TestStep{
				{
					Config:      baseConfig,
					ExpectError: regexp.MustCompile("Error: invalid config version"),
				},
			},
		},
		"error response from GetAdvancedSettingsURLEvasionDefense api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationForURLDefenseDS(m, 43253, 7)
				m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, appsec.GetAdvancedSettingsURLEvasionDefenseRequest{ConfigID: 43253, Version: 7}).
					Return(nil, &urlEvasionDefenseDSServerError)
			},
			steps: []resource.TestStep{
				{
					Config:      baseConfig,
					ExpectError: regexp.MustCompile("Error: calling 'GetAdvancedSettingsURLEvasionDefense'"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if test.init != nil {
				test.init(client.APPSEC)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    test.steps,
			})

			client.APPSEC.AssertExpectations(t)
		})
	}
}

func mockGetConfigurationForURLDefenseDS(m *appsec.Mock, configID, latestVersion int) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: configID}).
		Return(&appsec.GetConfigurationResponse{
			ID:                configID,
			LatestVersion:     latestVersion,
			StagingVersion:    6,
			ProductionVersion: 6,
		}, nil)
}

func mockGetConfigurationFailureForURLDefenseDS(m *appsec.Mock, configID int) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: configID}).
		Return(nil, &urlEvasionDefenseDSServerError)
}

var urlEvasionDefenseDSServerError = appsec.Error{
	Type:       "internal_error",
	Title:      "Internal Server Error",
	Detail:     "Error Fetching URL Evasion Defense",
	StatusCode: http.StatusInternalServerError,
}

func TestNewURLEvasionDefenseDataSource(t *testing.T) {
	t.Parallel()

	if NewURLEvasionDefenseDataSource() == nil {
		t.Fatal(fmt.Errorf("new URL Evasion Defense data source is nil"))
	}
}
