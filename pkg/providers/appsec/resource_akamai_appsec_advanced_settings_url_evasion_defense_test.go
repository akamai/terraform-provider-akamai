package appsec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type urlEvasionDefenseTestCase struct {
	init  func(*appsec.Mock)
	steps []resource.TestStep
}

type urlEvasionDefenseFixture struct {
	getConfigResponse   appsec.GetConfigurationResponse
	getResponseEnabled  appsec.GetAdvancedSettingsURLEvasionDefenseResponse
	getResponseDisabled appsec.GetAdvancedSettingsURLEvasionDefenseResponse
	updateResp          appsec.UpdateAdvancedSettingsURLEvasionDefenseResponse
}

func TestURLEvasionDefenseResource(t *testing.T) {
	t.Parallel()

	for name, test := range urlEvasionDefenseTests(t) {
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

func urlEvasionDefenseTests(t *testing.T) map[string]urlEvasionDefenseTestCase {
	t.Helper()

	fixture := newURLEvasionDefenseFixture()
	resourceName := "akamai_appsec_advanced_settings_url_evasion_defense.test"
	expectedAttrs := test.AttributeBatch{
		"config_id":                                 "43253",
		"status":                                    "enabled",
		"bypass_lists.#":                            "2",
		"rules.#":                                   "1",
		"rules.0.rule_id":                           "950001",
		"rules.0.action":                            "deny",
		"rules.0.name":                              "Rule 950001",
		"rules.0.description":                       "Test rule",
		"rules.0.condition_operator":                "OR",
		"rules.0.conditions.#":                      "9",
		"rules.0.conditions.0.type":                 "extensionMatch",
		"rules.0.conditions.0.extensions.#":         "2",
		"rules.0.conditions.0.positive_match":       "true",
		"rules.0.conditions.1.type":                 "filenameMatch",
		"rules.0.conditions.1.filenames.#":          "1",
		"rules.0.conditions.2.type":                 "hostMatch",
		"rules.0.conditions.2.hosts.#":              "1",
		"rules.0.conditions.3.type":                 "pathMatch",
		"rules.0.conditions.3.paths.#":              "1",
		"rules.0.conditions.4.type":                 "requestMethodMatch",
		"rules.0.conditions.4.methods.#":            "2",
		"rules.0.conditions.5.type":                 "ipMatch",
		"rules.0.conditions.5.ips.#":                "1",
		"rules.0.conditions.5.use_headers":          "true",
		"rules.0.conditions.6.type":                 "clientListMatch",
		"rules.0.conditions.6.client_lists.#":       "1",
		"rules.0.conditions.6.use_headers":          "true",
		"rules.0.conditions.7.type":                 "requestHeaderMatch",
		"rules.0.conditions.7.header":               "X-Test",
		"rules.0.conditions.7.value":                "abc*",
		"rules.0.conditions.7.value_case_sensitive": "true",
		"rules.0.conditions.7.value_wildcard":       "true",
		"rules.0.conditions.8.type":                 "uriQueryMatch",
		"rules.0.conditions.8.name":                 "param",
		"rules.0.conditions.8.value":                "val*",
		"rules.0.conditions.8.name_case_sensitive":  "true",
		"rules.0.conditions.8.value_case_sensitive": "true",
		"rules.0.conditions.8.value_wildcard":       "true",
	}

	createAndDestroyChecker := test.NewStateChecker(resourceName).
		CheckEqualBatch("", expectedAttrs)

	importChecker := newImportCheckerWithBatch(expectedAttrs)

	return map[string]urlEvasionDefenseTestCase{
		"schema validation - missing required status": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_missing_status.tf"),
				ExpectError: regexp.MustCompile("Missing required argument"),
			}},
		},
		"schema validation - invalid status": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_invalid_status.tf"),
				ExpectError: regexp.MustCompile("Attribute status value must be one of"),
			}},
		},
		"schema validation - invalid condition operator": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_invalid_condition_operator.tf"),
				ExpectError: regexp.MustCompile(`(?s)condition_operator.*must be one of`),
			}},
		},
		"schema validation - empty action": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_empty_action.tf"),
				ExpectError: regexp.MustCompile("Attribute action cannot be empty"),
			}},
		},
		"schema validation - empty bypass lists": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_empty_bypass_lists.tf"),
				ExpectError: regexp.MustCompile(`(?s)bypass_lists.*at least 1`),
			}},
		},
		"schema validation - positive match false": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_positive_match_false.tf"),
				ExpectError: regexp.MustCompile(`(?s)positive_match.*Value must be "true"`),
			}},
		},
		"create and destroy": {
			init: func(m *appsec.Mock) {
				mockURLEvasionDefenseCreateAndDestroy(m, fixture)
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense.tf"),
				Check:  createAndDestroyChecker.Build(),
			}},
		},
		"update": {
			init: func(m *appsec.Mock) {
				mockURLEvasionDefenseUpdate(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqual("status", "enabled").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_update.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqual("status", "disabled").Build(),
				},
			},
		},
		"update - 400 bad request": {
			init: func(m *appsec.Mock) {
				mockURLEvasionDefenseUpdateBadRequest(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense.tf"),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense_update.tf"),
					ExpectError: regexp.MustCompile(`(?s)URL Evasion Defense update failed.*Bad Request.*invalid payload`),
				},
			},
		},
		"import state": {
			init: func(m *appsec.Mock) {
				mockURLEvasionDefenseImport(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense.tf"),
					ResourceName:                         "akamai_appsec_advanced_settings_url_evasion_defense.test",
					ImportState:                          true,
					ImportStateId:                        "43253",
					ImportStateVerifyIdentifierAttribute: "config_id",
					ImportStateCheck:                     importChecker.Build(),
				},
			},
		},
		"import state - invalid id": {
			steps: []resource.TestStep{{
				ResourceName:  "akamai_appsec_advanced_settings_url_evasion_defense.test",
				Config:        testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsURLEvasionDefense/url_evasion_defense.tf"),
				ImportState:   true,
				ImportStateId: "not-a-config-id",
				ExpectError:   regexp.MustCompile(`ID 'not-a-config-id' incorrectly formatted: should be 'CONFIG_ID'`),
			}},
		},
	}
}

func newURLEvasionDefenseFixture() urlEvasionDefenseFixture {
	operator := string(appsec.AdvancedSettingsURLEvasionDefenseConditionOperatorOr)
	getResponseEnabled := appsec.GetAdvancedSettingsURLEvasionDefenseResponse{
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

	getResponseDisabled := getResponseEnabled
	getResponseDisabled.Status = string(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled)

	return urlEvasionDefenseFixture{
		getConfigResponse: appsec.GetConfigurationResponse{
			ID:                43253,
			LatestVersion:     7,
			StagingVersion:    6,
			ProductionVersion: 6,
		},
		getResponseEnabled:  getResponseEnabled,
		getResponseDisabled: getResponseDisabled,
		updateResp:          appsec.UpdateAdvancedSettingsURLEvasionDefenseResponse(getResponseEnabled),
	}
}

func mockURLEvasionDefenseCreateAndDestroy(m *appsec.Mock, fixture urlEvasionDefenseFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil).Times(4)

	m.On("GetConfigurationVersion", mock.Anything, appsec.GetConfigurationVersionRequest{ConfigID: 43253, Version: 7}).
		Return(nil, fmt.Errorf("not found")).Times(2)

	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{ConfigID: 43253, Version: 7}

	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).
		Return(&fixture.getResponseEnabled, nil).Once()

	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusEnabled))).
		Return(&fixture.updateResp, nil).Once()

	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).
		Return(&fixture.getResponseEnabled, nil).Times(3)

	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled))).
		Return(&fixture.updateResp, nil).Once()
}

func mockURLEvasionDefenseUpdate(m *appsec.Mock, fixture urlEvasionDefenseFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil).Times(8)

	m.On("GetConfigurationVersion", mock.Anything, appsec.GetConfigurationVersionRequest{ConfigID: 43253, Version: 7}).
		Return(nil, fmt.Errorf("not found")).Times(3)

	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{ConfigID: 43253, Version: 7}
	disabledUpdateResponse := ptr.To(appsec.UpdateAdvancedSettingsURLEvasionDefenseResponse(fixture.getResponseDisabled))

	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).Return(&fixture.getResponseEnabled, nil).Once()
	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusEnabled))).Return(&fixture.updateResp, nil).Once()
	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).Return(&fixture.getResponseEnabled, nil).Times(3)
	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled))).Return(disabledUpdateResponse, nil).Once()
	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).Return(&fixture.getResponseDisabled, nil).Once()
	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled))).Return(disabledUpdateResponse, nil).Once()
	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).Return(&fixture.getResponseDisabled, nil).Times(3)
}

func mockURLEvasionDefenseUpdateBadRequest(m *appsec.Mock, fixture urlEvasionDefenseFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil).Times(6)

	m.On("GetConfigurationVersion", mock.Anything, appsec.GetConfigurationVersionRequest{ConfigID: 43253, Version: 7}).
		Return(nil, fmt.Errorf("not found")).Times(3)

	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{ConfigID: 43253, Version: 7}

	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).
		Return(&fixture.getResponseEnabled, nil).Times(6)

	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusEnabled))).
		Return(&fixture.updateResp, nil).Once()

	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled))).
		Return(nil, &appsec.Error{StatusCode: 400, Title: "Bad Request", Detail: "invalid payload"}).Once()

	m.On("UpdateAdvancedSettingsURLEvasionDefense", mock.Anything, mock.MatchedBy(matchURLEvasionDefenseUpdateRequest(appsec.AdvancedSettingsURLEvasionDefenseStatusDisabled))).
		Return(ptr.To(appsec.UpdateAdvancedSettingsURLEvasionDefenseResponse(fixture.getResponseDisabled)), nil).Once()
}

func mockURLEvasionDefenseImport(m *appsec.Mock, fixture urlEvasionDefenseFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil)

	getRequest := appsec.GetAdvancedSettingsURLEvasionDefenseRequest{ConfigID: 43253, Version: 7}

	m.On("GetAdvancedSettingsURLEvasionDefense", mock.Anything, getRequest).
		Return(&fixture.getResponseEnabled, nil)
}

func matchURLEvasionDefenseUpdateRequest(status appsec.AdvancedSettingsURLEvasionDefenseStatus) func(appsec.UpdateAdvancedSettingsURLEvasionDefenseRequest) bool {
	return func(req appsec.UpdateAdvancedSettingsURLEvasionDefenseRequest) bool {
		return req.ConfigID == 43253 && req.Version == 7 && req.Body.Status == status && req.Body.LockVersion == 3
	}
}

func newImportCheckerWithBatch(batch test.AttributeBatch) test.ImportChecker {
	checker := test.NewImportChecker()
	for attr, val := range batch {
		checker = checker.CheckEqual(attr, val)
	}
	return checker
}
