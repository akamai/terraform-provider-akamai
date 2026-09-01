package appsec

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/mock"
)

const securityPolicyProtectionsResourceNameTest = "akamai_appsec_security_policy_protections.test"
const securityPolicyProtectionsBasicFixture = "testdata/TestResSecurityPolicyProtections/basic.tf"

func TestAkamaiAppSecSecurityPolicyProtectionsResource(t *testing.T) {
	t.Parallel()
	priorProtections := &appsec.PolicyProtectionsResponse{
		ApplyAPIConstraints:            true,
		ApplyAccountProtectionControls: true,
		ApplyApplicationLayerControls:  true,
		ApplyBotmanControls:            false,
		ApplyMalwareControls:           true,
		ApplyNetworkLayerControls:      false,
		ApplyRateControls:              true,
		ApplyReputationControls:        false,
		ApplySlowPostControls:          true,
		ApplyURLProtectionControls:     true,
	}
	updatedProtections := &appsec.PolicyProtectionsResponse{
		ApplyAPIConstraints:            false,
		ApplyAccountProtectionControls: false,
		ApplyApplicationLayerControls:  false,
		ApplyBotmanControls:            true,
		ApplyMalwareControls:           false,
		ApplyNetworkLayerControls:      true,
		ApplyRateControls:              false,
		ApplyReputationControls:        true,
		ApplySlowPostControls:          false,
		ApplyURLProtectionControls:     false,
	}
	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"create protections - passed all required fields": {
			init: func(m *appsec.Mock) {
				// GetConfiguration called: 1 for Create (getModifiableConfigVersion),
				// 1 for readState after Create (getLatestConfigVersion),
				// 1 for Read during plan before Delete (getLatestConfigVersion),
				// 1 for Delete (getModifiableConfigVersion)
				mockGetConfigurationPolicyProtection(m, 4)
				// Mock for Create
				mockUpdatePolicyProtections(m, appsec.UpdatePolicyProtectionsRequest{
					ConfigID:                       12345,
					Version:                        1,
					PolicyID:                       "test_policy",
					ApplyAccountProtectionControls: true,
					ApplyAPIConstraints:            true,
					ApplyApplicationLayerControls:  true,
					ApplyBotmanControls:            false,
					ApplyMalwareControls:           true,
					ApplyNetworkLayerControls:      false,
					ApplyRateControls:              true,
					ApplyReputationControls:        false,
					ApplySlowPostControls:          true,
					ApplyURLProtectionControls:     true,
				}, 1)
				// Mock for Delete (all protections disabled)
				mockUpdatePolicyProtections(m, appsec.UpdatePolicyProtectionsRequest{
					ConfigID:                       12345,
					Version:                        1,
					PolicyID:                       "test_policy",
					ApplyAccountProtectionControls: false,
					ApplyAPIConstraints:            false,
					ApplyApplicationLayerControls:  false,
					ApplyBotmanControls:            false,
					ApplyMalwareControls:           false,
					ApplyNetworkLayerControls:      false,
					ApplyRateControls:              false,
					ApplyReputationControls:        false,
					ApplySlowPostControls:          false,
					ApplyURLProtectionControls:     false,
				}, 1)
				// GetPolicyProtections called: 1 for readState after Create, 1 for Read during plan before Delete
				mockGetPolicyProtections(m, 12345, 1, "test_policy", priorProtections, 2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					Check:  protectionsChecker(loadProtectionsExpected(t, "expected_basic_protections.json")),
				},
			},
		},
		"update protections": {
			init: func(m *appsec.Mock) {
				// GetConfiguration called: 1 for Create (getModifiableConfigVersion),
				// 1 for readState after Create (getLatestConfigVersion),
				// 1 for Read during plan before Delete (getLatestConfigVersion),
				// 1 for Delete (getModifiableConfigVersion)
				mockGetConfigurationPolicyProtection(m, 8)
				// Mock for Create
				mockUpdatePolicyProtections(m, appsec.UpdatePolicyProtectionsRequest{
					ConfigID:                       12345,
					Version:                        1,
					PolicyID:                       "test_policy",
					ApplyAccountProtectionControls: true,
					ApplyAPIConstraints:            true,
					ApplyApplicationLayerControls:  true,
					ApplyBotmanControls:            false,
					ApplyMalwareControls:           true,
					ApplyNetworkLayerControls:      false,
					ApplyRateControls:              true,
					ApplyReputationControls:        false,
					ApplySlowPostControls:          true,
					ApplyURLProtectionControls:     true,
				}, 1)
				// GetPolicyProtections called: 1 for readState after Create, 1 for Read during plan before Delete, 1 for readState before Update
				mockGetPolicyProtections(m, 12345, 1, "test_policy", priorProtections, 3)
				// mock for update
				mockUpdatePolicyProtections(m, appsec.UpdatePolicyProtectionsRequest{
					ConfigID:                       12345,
					Version:                        1,
					PolicyID:                       "test_policy",
					ApplyAccountProtectionControls: false,
					ApplyAPIConstraints:            false,
					ApplyApplicationLayerControls:  false,
					ApplyBotmanControls:            true,
					ApplyMalwareControls:           false,
					ApplyNetworkLayerControls:      true,
					ApplyRateControls:              false,
					ApplyReputationControls:        true,
					ApplySlowPostControls:          false,
					ApplyURLProtectionControls:     false,
				}, 1)
				// Mock for Delete (all protections disabled)
				mockUpdatePolicyProtections(m, appsec.UpdatePolicyProtectionsRequest{
					ConfigID:                       12345,
					Version:                        1,
					PolicyID:                       "test_policy",
					ApplyAccountProtectionControls: false,
					ApplyAPIConstraints:            false,
					ApplyApplicationLayerControls:  false,
					ApplyBotmanControls:            false,
					ApplyMalwareControls:           false,
					ApplyNetworkLayerControls:      false,
					ApplyRateControls:              false,
					ApplyReputationControls:        false,
					ApplySlowPostControls:          false,
					ApplyURLProtectionControls:     false,
				}, 1)
				// GetPolicyProtections called: 1 for readState after Update, 1 for Read during plan before Delete
				mockGetPolicyProtections(m, 12345, 1, "test_policy", updatedProtections, 2)
			},

			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					Check:  protectionsChecker(loadProtectionsExpected(t, "expected_basic_protections.json")),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyProtections/update.tf"),
					Check:  protectionsChecker(loadProtectionsExpected(t, "expected_updated_protections.json")),
				},
			},
		},
		"missing required security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyProtections/missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"empty security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyProtections/empty_security_policy_id.tf"),
					ExpectError: regexp.MustCompile("Error: Missing required argument"),
				},
			},
		},
		"missing required protections control": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyProtections/missing_protections.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"empty value for protections control": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyProtections/empty_protections.tf"),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute"),
				},
			},
		},
		"null value for protections control": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResSecurityPolicyProtections/null_protections.tf"),
					ExpectError: regexp.MustCompile("Error: Missing Configuration for Required Attribute"),
				},
			},
		},
		"create protections - Unable to read latest config version from API": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationPolicyProtectionFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					ExpectError: regexp.MustCompile("Failed to retrieve modifiable config version"),
				},
			},
		},
		"create protections - Unable to update policy protections": {
			init: func(m *appsec.Mock) {
				// GetConfiguration called once for getModifiableConfigVersion in updatePolicyProtections
				mockGetConfigurationPolicyProtection(m, 1)
				// updatePolicyProtections calls UpdatePolicyProtections directly — no GetPolicyProtections
				mockUpdatePolicyProtectionsFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					ExpectError: regexp.MustCompile("Failed to update security policy protections"),
				},
			},
		},
		"read protections - Unable to get policy protections": {
			init: func(m *appsec.Mock) {
				// GetConfiguration called: 1 for Create (getModifiableConfigVersion),
				// 1 for readState after Create (getLatestConfigVersion)
				mockGetConfigurationPolicyProtection(m, 2)
				// updatePolicyProtections succeeds — no GetPolicyProtections call inside it
				mockUpdatePolicyProtections(m, appsec.UpdatePolicyProtectionsRequest{
					ConfigID:                       12345,
					Version:                        1,
					PolicyID:                       "test_policy",
					ApplyAccountProtectionControls: true,
					ApplyAPIConstraints:            true,
					ApplyApplicationLayerControls:  true,
					ApplyBotmanControls:            false,
					ApplyMalwareControls:           true,
					ApplyNetworkLayerControls:      false,
					ApplyRateControls:              true,
					ApplyReputationControls:        false,
					ApplySlowPostControls:          true,
					ApplyURLProtectionControls:     true,
				}, 1)
				// GetPolicyProtections fails in readState after Create
				mockGetPolicyProtectionsFailure(m, "test_policy", 12345, 1, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					ExpectError: regexp.MustCompile("Failed to read security policy protections"),
				},
			},
		},
		"import protections": {
			init: func(m *appsec.Mock) {
				// GetConfiguration called: 1 for Read after ImportState (getLatestConfigVersion),
				// 1 for Delete (getModifiableConfigVersion)
				mockGetConfigurationPolicyProtection(m, 2)
				// GetPolicyProtections called once in Read (via readState) after ImportState
				mockGetPolicyProtections(m, 12345, 1, "test_policy", priorProtections, 1)
				// Mock for Delete (all protections disabled)
				mockUpdatePolicyProtections(m, appsec.UpdatePolicyProtectionsRequest{
					ConfigID:                       12345,
					Version:                        1,
					PolicyID:                       "test_policy",
					ApplyAccountProtectionControls: false,
					ApplyAPIConstraints:            false,
					ApplyApplicationLayerControls:  false,
					ApplyBotmanControls:            false,
					ApplyMalwareControls:           false,
					ApplyNetworkLayerControls:      false,
					ApplyRateControls:              false,
					ApplyReputationControls:        false,
					ApplySlowPostControls:          false,
					ApplyURLProtectionControls:     false,
				}, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					ImportState:        true,
					ImportStateId:      "12345:test_policy",
					ResourceName:       securityPolicyProtectionsResourceNameTest,
					ImportStatePersist: true,
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						if len(states) == 0 {
							return nil
						}
						s := states[0]
						if s.Attributes["config_id"] != "12345" {
							return fmt.Errorf("config_id mismatch: got %s", s.Attributes["config_id"])
						}
						if s.Attributes["security_policy_id"] != "test_policy" {
							return fmt.Errorf("security_policy_id mismatch: got %s", s.Attributes["security_policy_id"])
						}
						return nil
					},
				},
			},
		},
		"import protections - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					ImportState:        true,
					ImportStateId:      "12345",
					ResourceName:       securityPolicyProtectionsResourceNameTest,
					ExpectError:        regexp.MustCompile("incorrectly formatted"),
					ImportStatePersist: true,
				},
			},
		},
		"import protections - invalid config id": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					ImportState:        true,
					ImportStateId:      "abc:test_policy",
					ResourceName:       securityPolicyProtectionsResourceNameTest,
					ExpectError:        regexp.MustCompile("invalid configuration id 'abc'"),
					ImportStatePersist: true,
				},
			},
		},
		"import protections - invalid security policy id": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, securityPolicyProtectionsBasicFixture),
					ImportState:        true,
					ImportStateId:      "12345:",
					ResourceName:       securityPolicyProtectionsResourceNameTest,
					ExpectError:        regexp.MustCompile("invalid security policy id ''"),
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

// mockGetConfigurationPolicyProtection mocks the GetConfiguration API call for a given number of times.
func mockGetConfigurationPolicyProtection(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 12345}).
		Return(&appsec.GetConfigurationResponse{
			LatestVersion: 1,
		}, nil).Times(times)
}

func mockGetConfigurationPolicyProtectionFailure(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 12345}).
		Return(nil, &serverError).Times(times)
}

func mockUpdatePolicyProtections(client *appsec.Mock, expected appsec.UpdatePolicyProtectionsRequest, times int) {
	client.On("UpdatePolicyProtections", mock.Anything, mock.MatchedBy(func(req appsec.UpdatePolicyProtectionsRequest) bool {
		return reflect.DeepEqual(req, expected)
	})).Return(&appsec.PolicyProtectionsResponse{}, nil).Times(times)
}

func mockUpdatePolicyProtectionsFailure(client *appsec.Mock, times int) {
	client.On("UpdatePolicyProtections", mock.Anything, mock.Anything).
		Return(nil, &serverError).Times(times)
}

func mockGetPolicyProtections(client *appsec.Mock, configID int, version int, policyID string, resp *appsec.PolicyProtectionsResponse, times int) {
	client.On("GetPolicyProtections", mock.Anything, appsec.GetPolicyProtectionsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}).Return(resp, nil).Times(times)
}

func mockGetPolicyProtectionsFailure(client *appsec.Mock, policyID string, configID, version, times int) {
	client.On("GetPolicyProtections", mock.Anything, appsec.GetPolicyProtectionsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}).Return(nil, &serverError).Times(times)
}

func protectionsChecker(attrs map[string]string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	for k, v := range attrs {
		checks = append(checks, resource.TestCheckResourceAttr(securityPolicyProtectionsResourceNameTest, k, v))
	}
	return resource.ComposeTestCheckFunc(checks...)
}

// loadProtectionsExpected loads expected protection attributes from a JSON file.
func loadProtectionsExpected(t *testing.T, filename string) map[string]string {
	t.Helper()
	data, err := os.ReadFile(fmt.Sprintf("testdata/TestResSecurityPolicyProtections/%s", filename))
	if err != nil {
		t.Fatalf("failed to read expected protections file %s: %v", filename, err)
	}
	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal expected protections from %s: %v", filename, err)
	}
	return result
}
