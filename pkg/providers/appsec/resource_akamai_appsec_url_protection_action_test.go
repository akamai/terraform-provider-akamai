package appsec

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestURLProtectionActionResource(t *testing.T) {
	t.Parallel()

	actionAfterCreate := appsec.GetURLProtectionPolicyActionsResponse{
		MaxRateThresholdAction: "alert",
		LoadSheddingAction:     "none",
	}

	actionAfterUpdateWithILS := appsec.GetURLProtectionPolicyActionsResponse{
		MaxRateThresholdAction: "deny",
		LoadSheddingAction:     "alert",
	}

	actionAfterUpdateWithoutILS := appsec.GetURLProtectionPolicyActionsResponse{
		MaxRateThresholdAction: "deny",
		LoadSheddingAction:     "none", // API returns "none" when ILS is off
	}

	urlProtectionPolicyWithILS := appsec.GetURLProtectionPolicyResponse{
		IntelligentLoadShedding: true,
	}

	urlProtectionPolicyWithoutILS := appsec.GetURLProtectionPolicyResponse{
		IntelligentLoadShedding: false,
	}

	baseChecker := test.NewStateChecker(urlProtectionActionResourceReferenceName).
		CheckEqual("config_id", "43253").
		CheckEqual("security_policy_id", "AAAA_81230").
		CheckEqual("url_protection_policy_id", "135355")

	createActionInit := func(m *appsec.Mock) {
		mockGetURLProtectionConfiguration(m, 7) // Called 2 times for ValidateConfig, 1 for create, 1 for Validate config, 1 for Read, 1 for ValidateConfig and 1 for delete
		mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 4)
		mockUpdateURLProtectionPolicyActions(m, "alert", "none", 1)
		mockGetURLProtectionPolicyActions(m, actionAfterCreate, 2) // Called after create and during refresh
		mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)
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
		"check schema - missing required attribute url_protection_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/missing_url_protection_policy_id.tf"),
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
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 4)
				// When ILS is off, API always gets "none" no matter what user sets
				mockUpdateURLProtectionPolicyActions(m, "deny", "none", 1)
				mockGetURLProtectionPolicyActions(m, actionAfterUpdateWithoutILS, 2) // Called after create and during refresh
				mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)           // Called during delete
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
		// This test verifies the fix where ValidateConfig was changed from getModifiableConfigVersion
		// to getLatestConfigVersion. When LatestVersion == ProductionVersion == StagingVersion (all = 1),
		// getModifiableConfigVersion would clone the config on every Create and Update call.
		// ValidateConfig must NOT clone — it must only read the latest version.
		// The two different version-sets for GetURLProtectionPolicy (v1) vs
		// UpdateURLProtectionPolicyActions / GetURLProtectionPolicyActions (v2) prove the split:
		//   - ValidateConfig path → getLatestConfigVersion → uses v1, no clone
		//   - Create/Delete path → getModifiableConfigVersion → detects active version → clones v1→v2, uses v2
		"create url protection action - ValidateConfig uses getLatestConfigVersion, Create/Delete clones active version": {
			init: func(m *appsec.Mock) {
				// GetConfiguration always returns LatestVersion = ProductionVersion = StagingVersion = 1.
				// 4 calls come from ValidateConfig (via getLatestConfigVersion, which just reads LatestVersion).
				// 1 call comes from Create (via getModifiableConfigVersion, which detects active version → clone).
				// 1 call comes from Read (via getLatestConfigVersion).
				// 1 call comes from Delete (via getModifiableConfigVersion → clone).
				mockGetNonModifiableConfiguration(m, 7)

				// ValidateConfig always uses version 1 because getLatestConfigVersion reads LatestVersion=1
				// from the locked mock. 4 calls: 2 pre-create (plan+apply) + 2 post-create (plan+apply no-op).
				// getModifiableConfigVersion does NOT call GetURLProtectionPolicy — only ValidateConfig does.
				mockGetURLProtectionPolicyAtVersion1(m, urlProtectionPolicyWithoutILS, 4)

				// getModifiableConfigVersion detects LatestVersion==ProductionVersion, triggering a clone.
				// Create and Delete each clone independently (cache not shared across operations).
				mockCreateConfigurationVersionClone(m, 2)

				// Call GetURLProtectionPolicyActions at version 2 (the cloned version) after the update, and read them back.
				mockGetURLProtectionPolicyActionsAtVersion2(m, actionAfterUpdateWithoutILS, 1)

				// Create writes actions to version 2 (the cloned version) and reads them back.
				mockUpdateURLProtectionPolicyActionsAtVersion2(m, "deny", "none", 1)

				// Read (refresh) uses version 1: getLatestConfigVersion returns LatestVersion=1 from locked mock.
				mockGetURLProtectionPolicyActionsAtVersion1(m, actionAfterUpdateWithoutILS, 1)

				// Delete resets actions at version 2 (the second clone).
				mockUpdateURLProtectionPolicyActionsAtVersion2(m, "none", "none", 1)
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
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithILS, 4)
				mockUpdateURLProtectionPolicyActions(m, "deny", "alert", 1)
				mockGetURLProtectionPolicyActions(m, actionAfterUpdateWithILS, 2) // Called after create and during refresh
				mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)        // Called during delete
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
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithILS, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_with_no_load_shedding_when_ils_enabled.tf"),
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
		"create url protection action - Failed to get URL protection policy": {
			init: func(m *appsec.Mock) {
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicyFailure(m)
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
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)
				mockUpdateURLProtectionPolicyActionsFailure(m, "alert", "none")
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
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)
				mockUpdateURLProtectionPolicyActions(m, "alert", "none", 1)
				mockGetURLProtectionPolicyActionsFailure(m)
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
				// Create phase - ValidateConfig
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)

				// Create phase - Create
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "alert", "none", 1)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)

				// Create phase - ValidateConfig (post-create)
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)

				// Create phase - Read (refresh)
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)

				// Update phase - ValidateConfig
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithILS, 2) // ILS now enabled

				// Update phase - Update
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicyActions(m, actionAfterUpdateWithILS, 1)

				// Update phase - ValidateConfig (post-update)
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithILS, 1)

				// Update phase - Read (refresh)
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicyActions(m, actionAfterUpdateWithILS, 1)

				// Delete phase
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)
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
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)

				//create
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "alert", "none", 1)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 1)

				// read
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 1)

				// Update: GetConfiguration should fail here (no retry in unit tests)
				mockGetURLProtectionConfigurationFailure(m, 1)

				// Cleanup after failed update - Terraform tries to delete
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)
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
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)

				//create
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "alert", "none", 1)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)

				// read
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 2)

				// validate
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithILS, 1)

				// Cleanup
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_no_ils.tf"),
					Check:  createActionChecker,
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionAction/create_with_no_load_shedding_when_ils_enabled.tf"), // Changes action but no load_shedding_action
					ExpectError: regexp.MustCompile("load_shedding_action is required when intelligent load shedding is enabled"),
				},
			},
		},
		"read url protection action - resource not found removes from state": {
			init: func(m *appsec.Mock) {
				//validate
				mockGetURLProtectionConfiguration(m, 2)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 2)

				// Create
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "alert", "none", 1)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)

				// validate
				mockGetURLProtectionConfiguration(m, 4)
				mockGetURLProtectionPolicy(m, urlProtectionPolicyWithoutILS, 4)

				// read
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)

				// Read - not found (resource removed from state)
				mockGetURLProtectionConfiguration(m, 1)
				mockGetURLProtectionPolicyActionsNotFound(m)

				// Cleanup - even though resource is not in state, test framework tries to delete
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)
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
				mockGetURLProtectionPolicyActions(m, actionAfterCreate, 1)
				// Delete after test
				mockGetURLProtectionConfiguration(m, 1)
				mockUpdateURLProtectionPolicyActions(m, "none", "none", 1)
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
						CheckEqual("url_protection_policy_id", "135355").
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
			t.Parallel()
			client := edgegrid.NewTestClient()

			if test.init != nil {
				test.init(client.APPSEC)
			}

			mockGetConfigurationVersionDefault(client.APPSEC)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    test.steps,
			})
			client.APPSEC.AssertExpectations(t)
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

func mockGetURLProtectionPolicy(client *appsec.Mock, resp appsec.GetURLProtectionPolicyResponse, times int) {
	client.On("GetURLProtectionPolicy", mock.Anything, appsec.GetURLProtectionPolicyRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		URLProtectionPolicyID: 135355,
	}).Return(&resp, nil).Times(times)
}

func mockGetURLProtectionPolicyFailure(client *appsec.Mock) {
	client.On("GetURLProtectionPolicy", mock.Anything, appsec.GetURLProtectionPolicyRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		URLProtectionPolicyID: 135355,
	}).Return(nil, &urlProtectionActionServerError).Once()
}

func mockUpdateURLProtectionPolicyActions(client *appsec.Mock, maxRateAction, loadSheddingAction string, times int) {
	req := appsec.UpdateURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
		Body: appsec.URLProtectionPolicyActions{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		},
	}
	client.On("UpdateURLProtectionPolicyActions", mock.Anything, req).
		Return(&appsec.UpdateURLProtectionPolicyActionsResponse{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		}, nil).Times(times)
}

func mockUpdateURLProtectionPolicyActionsFailure(client *appsec.Mock, maxRateAction, loadSheddingAction string) {
	req := appsec.UpdateURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
		Body: appsec.URLProtectionPolicyActions{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		},
	}
	client.On("UpdateURLProtectionPolicyActions", mock.Anything, req).
		Return(nil, &urlProtectionActionServerError).Once()
}

func mockGetURLProtectionPolicyActions(client *appsec.Mock, resp appsec.GetURLProtectionPolicyActionsResponse, times int) {
	client.On("GetURLProtectionPolicyActions", mock.Anything, appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
	}).Return(&resp, nil).Times(times)
}

func mockGetURLProtectionPolicyActionsFailure(client *appsec.Mock) {
	client.On("GetURLProtectionPolicyActions", mock.Anything, appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
	}).Return(nil, &urlProtectionActionServerError).Once()
}

func mockGetURLProtectionPolicyActionsNotFound(client *appsec.Mock) {
	notFoundError := appsec.Error{
		Type:       "not_found",
		Title:      "Not Found",
		Detail:     "incorrect URL Protection Policy ID",
		StatusCode: http.StatusNotFound,
	}
	client.On("GetURLProtectionPolicyActions", mock.Anything, appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
	}).Return(nil, &notFoundError).Once()
}

// mockGetURLProtectionConfigurationLocked mocks GetConfiguration returning a config where
// LatestVersion = ProductionVersion = StagingVersion = 1. getModifiableConfigVersion (used by
// Create/Update/Delete) detects the version is active and immediately triggers a clone; meanwhile
// ValidateConfig's getLatestConfigVersion just reads and returns version 1 with no clone.
func mockGetNonModifiableConfiguration(client *appsec.Mock, times int) {
	// Each call gets its own fresh struct to prevent getModifiableConfigVersion from
	// mutating the shared pointer (it sets configuration.LatestVersion = clonedVersion
	// before caching, which would corrupt subsequent GetConfiguration mock returns).
	for i := 0; i < times; i++ {
		client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
			Return(&appsec.GetConfigurationResponse{
				FileType:          "RBAC",
				ID:                43253,
				LatestVersion:     1,
				Name:              "Akamai Tools",
				ProductionVersion: 1,
				StagingVersion:    1,
				TargetProduct:     "KSD",
			}, nil).Once()
	}
}

// mockCreateConfigurationVersionClone mocks the clone of version 1 into version 2.
func mockCreateConfigurationVersionClone(client *appsec.Mock, times int) {
	client.On("CreateConfigurationVersionClone", mock.Anything, appsec.CreateConfigurationVersionCloneRequest{
		ConfigID:          43253,
		CreateFromVersion: 1,
	}).Return(&appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: 43253,
		Version:  2,
	}, nil).Times(times)
}

// mockGetURLProtectionPolicyAtVersion1 mocks GetURLProtectionPolicy at config version 1.
// Used for ValidateConfig calls when getLatestConfigVersion returns version 1.
func mockGetURLProtectionPolicyAtVersion1(client *appsec.Mock, resp appsec.GetURLProtectionPolicyResponse, times int) {
	client.On("GetURLProtectionPolicy", mock.Anything, appsec.GetURLProtectionPolicyRequest{
		ConfigID:              43253,
		ConfigVersion:         1,
		URLProtectionPolicyID: 135355,
	}).Return(&resp, nil).Times(times)
}

// mockUpdateURLProtectionPolicyActionsAtVersion2 mocks UpdateURLProtectionPolicyActions at
// config version 2 — the cloned version produced by getModifiableConfigVersion.
func mockUpdateURLProtectionPolicyActionsAtVersion2(client *appsec.Mock, maxRateAction, loadSheddingAction string, times int) {
	req := appsec.UpdateURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         2,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
		Body: appsec.URLProtectionPolicyActions{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		},
	}
	client.On("UpdateURLProtectionPolicyActions", mock.Anything, req).
		Return(&appsec.UpdateURLProtectionPolicyActionsResponse{
			MaxRateThresholdAction: maxRateAction,
			LoadSheddingAction:     loadSheddingAction,
		}, nil).Times(times)
}

// mockGetURLProtectionPolicyActionsAtVersion2 mocks GetURLProtectionPolicyActions at
// config version 2 — the cloned version used by Create/Update.
func mockGetURLProtectionPolicyActionsAtVersion2(client *appsec.Mock, resp appsec.GetURLProtectionPolicyActionsResponse, times int) {
	client.On("GetURLProtectionPolicyActions", mock.Anything, appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         2,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
	}).Return(&resp, nil).Times(times)
}

// mockGetURLProtectionPolicyActionsAtVersion1 mocks GetURLProtectionPolicyActions at
// config version 1 — used by Read when getLatestConfigVersion returns version 1.
func mockGetURLProtectionPolicyActionsAtVersion1(client *appsec.Mock, resp appsec.GetURLProtectionPolicyActionsResponse, times int) {
	client.On("GetURLProtectionPolicyActions", mock.Anything, appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         1,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
	}).Return(&resp, nil).Times(times)
}

var urlProtectionActionResourceReferenceName = "akamai_appsec_url_protection_action.test"

var urlProtectionActionServerError = appsec.Error{
	Type:       "internal_error",
	Title:      "Internal Server Error",
	Detail:     "Error updating URL protection action",
	StatusCode: http.StatusInternalServerError,
}
