package appsec

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestURLProtectionActionResource(t *testing.T) {
	t.Parallel()

	actionAfterCreate := appsec.GetURLProtectionRuleActionsResponse{
		MaxRateThresholdAction: "alert",
		LoadSheddingAction:     "",
	}

	actionAfterUpdateWithILS := appsec.GetURLProtectionRuleActionsResponse{
		MaxRateThresholdAction: "deny",
		LoadSheddingAction:     "alert",
	}

	actionAfterUpdateWithoutILS := appsec.GetURLProtectionRuleActionsResponse{
		MaxRateThresholdAction: "deny",
		LoadSheddingAction:     "none", // API returns "none" when ILS is off
	}

	urlProtectionRuleWithILS := appsec.GetURLProtectionRuleResponse{
		IntelligentLoadShedding: true,
	}

	urlProtectionRuleWithoutILS := appsec.GetURLProtectionRuleResponse{
		IntelligentLoadShedding: false,
	}

	baseChecker := test.NewStateChecker(urlProtectionActionResourceReferenceName).
		CheckEqual("config_id", "43253").
		CheckEqual("security_policy_id", "AAAA_81230").
		CheckEqual("url_protection_rule_id", "135355")

	createActionInit := func(m *appsec.Mock) {
		mockGetURLProtectionConfiguration(m, 7) // Called 2 times for ValidateConfig, 1 for create, 1 for Validate config, 1 for Read, 1 for ValidateConfig and 1 for delete
		mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 4)
		mockUpdateURLProtectionRuleActions(m, "alert", "none", 1)
		mockGetURLProtectionRuleActions(m, actionAfterCreate, 2) // Called after create and during refresh
		mockUpdateURLProtectionRuleActions(m, "none", "none", 1)
	}

	createActionChecker := baseChecker.
		CheckEqual("max_rate_threshold_action", "alert").
		Build()

	var tests = map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"check schema - missing required attribute config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/missing_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema - missing required attribute security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema - missing required attribute url_protection_rule_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/missing_url_protection_rule_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema - missing required attribute max_rate_threshold_action": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/missing_max_rate_threshold_action.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"validate config - invalid max_rate_threshold_action value": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/invalid_action.tf"),
					ExpectError: regexp.MustCompile(`Attribute max_rate_threshold_action (value must be one of|must start with)`),
				},
			},
		},
		"validate config - invalid load_shedding_action value": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/invalid_load_shedding_action.tf"),
					ExpectError: regexp.MustCompile(`Attribute load_shedding_action (value must be one of|must start with)`),
				},
			},
		},
		"create url protection action - required only": {
			init: createActionInit,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					Check:  createActionChecker,
				},
			},
		},
		"create url protection action - with load_shedding_action when ILS disabled": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfiguration(m, 7) // Called 2 times for ValidateConfig, 1 for update, 1 for Validate config, 1 for Read, 1 for ValidateConfig and 1 for delete
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 4)
				// When ILS is off, API always gets "none" no matter what user sets
				mockUpdateURLProtectionRuleActions(m, "deny", "none", 1)
				mockGetURLProtectionRuleActions(m, actionAfterUpdateWithoutILS, 2) // Called after create and during refresh
				mockUpdateURLProtectionRuleActions(m, "none", "none", 1)           // Called during delete
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/update_no_ils.tf"),
					Check: baseChecker.
						CheckEqual("max_rate_threshold_action", "deny").
						Build(),
				},
			},
		},
		"create url protection action - with load_shedding_action when ILS enabled": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfiguration(m, 7) // Called 2 times for ValidateConfig, 1 for update, 1 for Validate config, 1 for Read, 1 for ValidateConfig and 1 for delete
				mockGetURLProtectionRule(m, urlProtectionRuleWithILS, 4)
				mockUpdateURLProtectionRuleActions(m, "deny", "alert", 1)
				mockGetURLProtectionRuleActions(m, actionAfterUpdateWithILS, 2) // Called after create and during refresh
				mockUpdateURLProtectionRuleActions(m, "none", "none", 1)        // Called during delete
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/update.tf"),
					Check: baseChecker.
						CheckEqual("max_rate_threshold_action", "deny").
						Build(),
				},
			},
		},
		"create url protection action - missing load_shedding_action when ILS enabled": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRule(m, urlProtectionRuleWithILS, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					ExpectError: regexp.MustCompile("load_shedding_action is required when intelligent load shedding is enabled"),
				},
			},
		},
		"create url protection action - Unable to read latest config version from API": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfigurationFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create.tf"),
					ExpectError: regexp.MustCompile("Unable to read latest config version from API"),
				},
			},
		},
		"create url protection action - Failed to get URL protection rule": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRuleFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create.tf"),
					ExpectError: regexp.MustCompile("Failed to prepare URL protection action"),
				},
			},
		},
		"create url protection action - Failed to create": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfiguration(m, 3)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 2)
				mockUpdateURLProtectionRuleActionsFailure(m, "alert", "none")
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					ExpectError: regexp.MustCompile("Failed to create URL protection action"),
				},
			},
		},
		"create url protection action - Failed to read after creation": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfiguration(m, 3)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 2)
				mockUpdateURLProtectionRuleActions(m, "alert", "none", 1)
				mockGetURLProtectionRuleActionsFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					ExpectError: regexp.MustCompile("Failed to read URL protection action after creation"),
				},
			},
		},
		"update url protection action": {
			init: func(m *appsec.Mock) {
				// Create phase
				mockGetURLProtectionConfiguration(m, 5)                     // 6 GetConfiguration calls during create
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 4) // 4 GetURLProtectionRule calls during create
				mockUpdateURLProtectionRuleActions(m, "alert", "none", 1)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 2)

				// Update phase
				mockGetURLProtectionConfiguration(m, 9)                  // 9 GetConfiguration calls during update
				mockGetURLProtectionRule(m, urlProtectionRuleWithILS, 4) // 4 GetURLProtectionRule calls during update
				mockUpdateURLProtectionRuleActions(m, "deny", "alert", 1)
				mockGetURLProtectionRuleActions(m, actionAfterUpdateWithILS, 3)
				mockUpdateURLProtectionRuleActions(m, "none", "none", 1)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					Check:  createActionChecker,
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/update.tf"),
					Check: baseChecker.
						CheckEqual("max_rate_threshold_action", "deny").
						CheckEqual("load_shedding_action", "alert").
						Build(),
				},
			},
		},
		"update url protection action - Unable to read latest config version": {
			init: func(m *appsec.Mock) {
				//validate
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 2)

				//create
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionRuleActions(m, "alert", "none", 1)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 1)

				// read
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 1)

				// Update: GetConfiguration should fail here (no retry in unit tests)
				mockGetURLProtectionConfigurationFailure(m, 1)

				// Cleanup after failed update - Terraform tries to delete
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionRuleActions(m, "none", "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					Check:  createActionChecker,
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/update_no_ils.tf"),
					ExpectError: regexp.MustCompile("Unable to read config version"),
				},
			},
		},
		"update url protection action - missing load_shedding_action when ILS enabled": {
			init: func(m *appsec.Mock) {
				//validate
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 2)

				//create
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionRuleActions(m, "alert", "none", 1)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 2)

				// read
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 2)

				// validate
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRule(m, urlProtectionRuleWithILS, 1)

				// Cleanup
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionRuleActions(m, "none", "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					Check:  createActionChecker,
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/update_no_ils.tf"), // Changes action but no load_shedding_action
					ExpectError: regexp.MustCompile("load_shedding_action is required when intelligent load shedding is enabled"),
				},
			},
		},
		"read url protection action - resource not found removes from state": {
			init: func(m *appsec.Mock) {
				//validate
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 2)

				// Create
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionRuleActions(m, "alert", "none", 1)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 4)
				mockGetURLProtectionRule(m, urlProtectionRuleWithoutILS, 4)

				// read
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 1)

				// Read - not found (resource removed from state)
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRuleActionsNotFound(m)

				// Cleanup - even though resource is not in state, test framework tries to delete
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionRuleActions(m, "none", "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					Check:  createActionChecker,
				},
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					PlanOnly:           true,
					ExpectNonEmptyPlan: true, // Resource will be recreated
				},
			},
		},
		"import url protection action": {
			init: func(m *appsec.Mock) {
				// Import and post-import refresh reads
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionRuleActions(m, actionAfterCreate, 1)
				// Delete after test
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionRuleActions(m, "none", "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create.tf"),
					ImportState:   true,
					ImportStateId: "43253:AAAA_81230:135355",
					ResourceName:  urlProtectionActionResourceReferenceName,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("config_id", "43253").
						CheckEqual("security_policy_id", "AAAA_81230").
						CheckEqual("url_protection_rule_id", "135355").
						CheckEqual("max_rate_threshold_action", "alert").
						Build(),
					ImportStatePersist: true,
				},
			},
		},
		"import url protection action - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create.tf"),
					ImportState:        true,
					ImportStateId:      "12345",
					ResourceName:       urlProtectionActionResourceReferenceName,
					ExpectError:        regexp.MustCompile("ID '12345' incorrectly formatted"),
					ImportStatePersist: true,
				},
			},
		},
		"import url protection action - invalid config id value": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create.tf"),
					ImportState:        true,
					ImportStateId:      "abc:AAAA_81230:135355",
					ResourceName:       urlProtectionActionResourceReferenceName,
					ExpectError:        regexp.MustCompile("invalid configuration id 'abc'"),
					ImportStatePersist: true,
				},
			},
		},
		"import url protection action - invalid url protection id value": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create.tf"),
					ImportState:        true,
					ImportStateId:      "43253:AAAA_81230:xyz",
					ResourceName:       urlProtectionActionResourceReferenceName,
					ExpectError:        regexp.MustCompile("invalid url protection id"),
					ImportStatePersist: true,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
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

// Mock helper functions

func mockGetURLProtectionConfiguration(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&appsec.GetConfigurationResponse{
			FileType:          "RBAC",
			ID:                43253,
			LatestVersion:     7,
			Name:              "Akamai Tools",
			ProductionVersion: 1,
			StagingVersion:    1,
			TargetProduct:     "KSD",
		}, nil).Times(times)
}

func mockGetURLProtectionConfigurationFailure(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(nil, &urlProtectionActionServerError).Times(times)
}

func mockGetURLProtectionRule(client *appsec.Mock, resp appsec.GetURLProtectionRuleResponse, times int) {
	client.On("GetURLProtectionRule", mock.Anything, appsec.GetURLProtectionRuleRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		URLProtectionRuleID: 135355,
	}).Return(&resp, nil).Times(times)
}

func mockGetURLProtectionRuleFailure(client *appsec.Mock) {
	client.On("GetURLProtectionRule", mock.Anything, appsec.GetURLProtectionRuleRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		URLProtectionRuleID: 135355,
	}).Return(nil, &urlProtectionActionServerError).Once()
}

func mockUpdateURLProtectionRuleActions(client *appsec.Mock, maxRateAction, loadSheddingAction string, times int) {
	req := appsec.UpdateURLProtectionRuleActionsRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		PolicyID:            "AAAA_81230",
		URLProtectionRuleID: 135355,
		Body: appsec.URLProtectionRuleActions{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		},
	}
	client.On("UpdateURLProtectionRuleActions", mock.Anything, req).
		Return(&appsec.UpdateURLProtectionRuleActionsResponse{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		}, nil).Times(times)
}

func mockUpdateURLProtectionRuleActionsFailure(client *appsec.Mock, maxRateAction, loadSheddingAction string) {
	req := appsec.UpdateURLProtectionRuleActionsRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		PolicyID:            "AAAA_81230",
		URLProtectionRuleID: 135355,
		Body: appsec.URLProtectionRuleActions{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		},
	}
	client.On("UpdateURLProtectionRuleActions", mock.Anything, req).
		Return(nil, &urlProtectionActionServerError).Once()
}

func mockGetURLProtectionRuleActions(client *appsec.Mock, resp appsec.GetURLProtectionRuleActionsResponse, times int) {
	client.On("GetURLProtectionRuleActions", mock.Anything, appsec.GetURLProtectionRuleActionsRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		PolicyID:            "AAAA_81230",
		URLProtectionRuleID: 135355,
	}).Return(&resp, nil).Times(times)
}

func mockGetURLProtectionRuleActionsFailure(client *appsec.Mock) {
	client.On("GetURLProtectionRuleActions", mock.Anything, appsec.GetURLProtectionRuleActionsRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		PolicyID:            "AAAA_81230",
		URLProtectionRuleID: 135355,
	}).Return(nil, &urlProtectionActionServerError).Once()
}

func mockGetURLProtectionRuleActionsNotFound(client *appsec.Mock) {
	notFoundError := appsec.Error{
		Type:       "not_found",
		Title:      "Not Found",
		Detail:     "incorrect URL Protection Rule ID",
		StatusCode: http.StatusNotFound,
	}
	client.On("GetURLProtectionRuleActions", mock.Anything, appsec.GetURLProtectionRuleActionsRequest{
		ConfigID:            43253,
		ConfigVersion:       7,
		PolicyID:            "AAAA_81230",
		URLProtectionRuleID: 135355,
	}).Return(nil, &notFoundError).Once()
}

var urlProtectionActionResourceReferenceName = "akamai_appsec_url_protection_action.test"

var urlProtectionActionServerError = appsec.Error{
	Type:       "internal_error",
	Title:      "Internal Server Error",
	Detail:     "Error updating URL protection action",
	StatusCode: http.StatusInternalServerError,
}
