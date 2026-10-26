package appsec

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestClientsideProtectionConfigurationResource(t *testing.T) {
	t.Parallel()

	for name, tc := range clientsideProtectionConfigurationTests(t) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.APPSEC)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    tc.steps,
			})
			client.APPSEC.AssertExpectations(t)
		})
	}
}

type clientsideProtectionConfigurationTestCase struct {
	init  func(*appsec.Mock)
	steps []resource.TestStep
}

func clientsideProtectionConfigurationTests(t *testing.T) map[string]clientsideProtectionConfigurationTestCase {
	t.Helper()

	fixture := newClientsideProtectionConfigurationFixture()
	resourceName := "akamai_appsec_clientside_protection_configuration.test"

	expectedAttrs := test.AttributeBatch{
		"config_id":                                    "43253",
		"security_policy_id":                           "test_policy_id",
		"edge_injection.load_script_async":             "true",
		"edge_test_parameters.disable_injection_key":   "disable-key",
		"edge_test_parameters.disable_injection_value": "disable-value",
		"edge_test_parameters.force_injection_key":     "force-key",
		"edge_test_parameters.force_injection_value":   "force-value",
		"injection_criteria.injection_path_includes":   "/include-path",
		"injection_criteria.injection_path_excludes":   "/exclude-path",
		"injection_criteria.injection_percent":         "100",
		"injection_criteria.injection_policy":          "head",
		"injection_criteria.script_tag_attributes.#":   "0",
	}

	expectedAttrsWithScriptTagAttrs := test.AttributeBatch{
		"config_id":          "43253",
		"security_policy_id": "test_policy_id",
		"injection_criteria.script_tag_attributes.#": "2",
		"injection_criteria.script_tag_attributes.0": `nonce="abc123"`,
		"injection_criteria.script_tag_attributes.1": `data-integrity="sha256"`,
	}

	createAndDestroyChecker := test.NewStateChecker(resourceName).
		CheckEqualBatch("", expectedAttrs)

	return map[string]clientsideProtectionConfigurationTestCase{
		"create and destroy": {
			init: func(m *appsec.Mock) {
				mockClientsideProtectionConfigurationCreateAndDestroy(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration.tf"),
					Check:  createAndDestroyChecker.Build(),
				},
			},
		},
		"update": {
			init: func(m *appsec.Mock) {
				mockClientsideProtectionConfigurationUpdate(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqual("edge_injection.load_script_async", "true").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration_update.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqual("edge_injection.load_script_async", "false").Build(),
				},
			},
		},
		"update - api error": {
			init: func(m *appsec.Mock) {
				mockClientsideProtectionConfigurationUpdateError(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration.tf"),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration_update.tf"),
					ExpectError: regexp.MustCompile("Error: update clientside protection configuration request failed"),
				},
			},
		},
		"create and destroy with script_tag_attributes": {
			init: func(m *appsec.Mock) {
				mockClientsideProtectionConfigurationCreateAndDestroyWithScriptTagAttrs(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration_with_script_tag_attributes.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqualBatch("", expectedAttrsWithScriptTagAttrs).Build(),
				},
			},
		},
		"import state": {
			init: func(m *appsec.Mock) {
				mockClientsideProtectionConfigurationImport(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration.tf"),
				},
				{
					ResourceName:                         resourceName,
					ImportState:                          true,
					ImportStateId:                        "43253:test_policy_id",
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "config_id",
				},
			},
		},
		"import state - invalid id": {
			init: func(m *appsec.Mock) {
				mockClientsideProtectionConfigurationCreateAndDestroy(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResClientsideProtectionConfiguration/clientside_protection_configuration.tf"),
				},
				{
					ResourceName:  resourceName,
					ImportState:   true,
					ImportStateId: "invalid-id",
					ExpectError:   regexp.MustCompile(`invalid import ID "invalid-id"`),
				},
			},
		},
	}
}

type clientsideProtectionConfigurationFixture struct {
	getConfigResponse             appsec.GetConfigurationResponse
	getResponse                   appsec.GetClientsideProtectionConfigurationResponse
	getResponseUpdate             appsec.GetClientsideProtectionConfigurationResponse
	getResponseWithScriptTagAttrs appsec.GetClientsideProtectionConfigurationResponse
	updateResponse                appsec.UpdateClientsideProtectionConfigurationResponse
	updateRespUpdated             appsec.UpdateClientsideProtectionConfigurationResponse
	updateRespWithScriptTagAttrs  appsec.UpdateClientsideProtectionConfigurationResponse
}

func newClientsideProtectionConfigurationFixture() clientsideProtectionConfigurationFixture {
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                43253,
		LatestVersion:     7,
		StagingVersion:    6,
		ProductionVersion: 6,
	}

	getResponse := appsec.GetClientsideProtectionConfigurationResponse{
		EdgeInjection: appsec.ClientsideProtectionEdgeInjection{
			ClientsideProtectionConfigID: 12345,
			LoadScriptAsync:              true,
		},
		EdgeTestParameters: appsec.ClientsideProtectionEdgeTestParameters{
			DisableInjectionKey:   "disable-key",
			DisableInjectionValue: "disable-value",
			ForceInjectionKey:     "force-key",
			ForceInjectionValue:   "force-value",
		},
		InjectionCriteria: appsec.ClientsideProtectionInjectionCriteria{
			InjectionPathIncludes: ptr.To("/include-path"),
			InjectionPathExcludes: ptr.To("/exclude-path"),
			InjectionPercent:      100,
			InjectionPolicy:       ptr.To("head"),
		},
	}

	getResponseUpdate := appsec.GetClientsideProtectionConfigurationResponse{
		EdgeInjection: appsec.ClientsideProtectionEdgeInjection{
			ClientsideProtectionConfigID: 12345,
			LoadScriptAsync:              false,
		},
		EdgeTestParameters: appsec.ClientsideProtectionEdgeTestParameters{
			DisableInjectionKey:   "disable-key",
			DisableInjectionValue: "disable-value",
			ForceInjectionKey:     "force-key",
			ForceInjectionValue:   "force-value",
		},
		InjectionCriteria: appsec.ClientsideProtectionInjectionCriteria{
			InjectionPathIncludes: ptr.To("/include-path"),
			InjectionPathExcludes: ptr.To("/exclude-path"),
			InjectionPercent:      50,
			InjectionPolicy:       ptr.To("head"),
		},
	}

	getResponseWithScriptTagAttrs := appsec.GetClientsideProtectionConfigurationResponse{
		EdgeInjection: appsec.ClientsideProtectionEdgeInjection{
			ClientsideProtectionConfigID: 12345,
			LoadScriptAsync:              true,
		},
		EdgeTestParameters: appsec.ClientsideProtectionEdgeTestParameters{
			DisableInjectionKey:   "disable-key",
			DisableInjectionValue: "disable-value",
			ForceInjectionKey:     "force-key",
			ForceInjectionValue:   "force-value",
		},
		InjectionCriteria: appsec.ClientsideProtectionInjectionCriteria{
			InjectionPathIncludes: ptr.To("/include-path"),
			InjectionPathExcludes: ptr.To("/exclude-path"),
			InjectionPercent:      100,
			InjectionPolicy:       ptr.To("head"),
			ScriptTagAttributes:   []string{`nonce="abc123"`, `data-integrity="sha256"`},
		},
	}

	updateResponse := appsec.UpdateClientsideProtectionConfigurationResponse(getResponse)
	updateRespUpdated := appsec.UpdateClientsideProtectionConfigurationResponse(getResponseUpdate)
	updateRespWithScriptTagAttrs := appsec.UpdateClientsideProtectionConfigurationResponse(getResponseWithScriptTagAttrs)

	return clientsideProtectionConfigurationFixture{
		getConfigResponse:             getConfigResponse,
		getResponse:                   getResponse,
		getResponseUpdate:             getResponseUpdate,
		getResponseWithScriptTagAttrs: getResponseWithScriptTagAttrs,
		updateResponse:                updateResponse,
		updateRespUpdated:             updateRespUpdated,
		updateRespWithScriptTagAttrs:  updateRespWithScriptTagAttrs,
	}
}

func mockClientsideProtectionConfigurationCreateAndDestroy(m *appsec.Mock, fixture clientsideProtectionConfigurationFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil)
	mockGetConfigurationVersionDefault(m)

	updateReq := buildCPCUpdateRequest(fixture, true, 100)
	m.On("UpdateClientsideProtectionConfiguration", mock.Anything, updateReq).
		Return(&fixture.updateResponse, nil).Once()

	m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
	}).Return(&fixture.getResponse, nil)
}

func mockClientsideProtectionConfigurationCreateAndDestroyWithScriptTagAttrs(m *appsec.Mock, fixture clientsideProtectionConfigurationFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil)
	mockGetConfigurationVersionDefault(m)

	updateReq := buildCPCUpdateRequest(fixture, true, 100, []string{`nonce="abc123"`, `data-integrity="sha256"`})
	m.On("UpdateClientsideProtectionConfiguration", mock.Anything, updateReq).
		Return(&fixture.updateRespWithScriptTagAttrs, nil).Once()

	m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
	}).Return(&fixture.getResponseWithScriptTagAttrs, nil)
}

func mockClientsideProtectionConfigurationUpdate(m *appsec.Mock, fixture clientsideProtectionConfigurationFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil)
	mockGetConfigurationVersionDefault(m)

	updateReq := buildCPCUpdateRequest(fixture, true, 100)
	m.On("UpdateClientsideProtectionConfiguration", mock.Anything, updateReq).
		Return(&fixture.updateResponse, nil).Once()

	m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
	}).Return(&fixture.getResponse, nil).Times(3)

	updateReqUpdated := buildCPCUpdateRequest(fixture, false, 50)
	m.On("UpdateClientsideProtectionConfiguration", mock.Anything, updateReqUpdated).
		Return(&fixture.updateRespUpdated, nil).Once()

	m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
	}).Return(&fixture.getResponseUpdate, nil)
}

func mockClientsideProtectionConfigurationUpdateError(m *appsec.Mock, fixture clientsideProtectionConfigurationFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil)
	mockGetConfigurationVersionDefault(m)

	updateReq := buildCPCUpdateRequest(fixture, true, 100)
	m.On("UpdateClientsideProtectionConfiguration", mock.Anything, updateReq).
		Return(&fixture.updateResponse, nil).Once()

	m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
	}).Return(&fixture.getResponse, nil)

	serverError := &appsec.Error{
		Type:       "internal_error",
		Title:      "Internal Server Error",
		Detail:     "update clientside protection configuration request failed",
		StatusCode: http.StatusInternalServerError,
	}
	updateReqUpdated := buildCPCUpdateRequest(fixture, false, 50)
	m.On("UpdateClientsideProtectionConfiguration", mock.Anything, updateReqUpdated).
		Return(nil, serverError).Once()
}

func mockClientsideProtectionConfigurationImport(m *appsec.Mock, fixture clientsideProtectionConfigurationFixture) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&fixture.getConfigResponse, nil)
	mockGetConfigurationVersionDefault(m)

	updateReq := buildCPCUpdateRequest(fixture, true, 100)
	m.On("UpdateClientsideProtectionConfiguration", mock.Anything, updateReq).
		Return(&fixture.updateResponse, nil).Once()

	m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
	}).Return(&fixture.getResponse, nil)
}

func buildCPCUpdateRequest(fixture clientsideProtectionConfigurationFixture, loadScriptAsync bool, injectionPercent int, scriptTagAttributes ...[]string) appsec.UpdateClientsideProtectionConfigurationRequest {
	var attrs []string
	if len(scriptTagAttributes) > 0 {
		attrs = scriptTagAttributes[0]
	}
	return appsec.UpdateClientsideProtectionConfigurationRequest{
		ConfigID:         43253,
		Version:          7,
		SecurityPolicyID: "test_policy_id",
		Body: appsec.UpdateClientsideProtectionConfigurationRequestBody{
			EdgeInjection: appsec.ClientsideProtectionEdgeInjectionRequest{
				LoadScriptAsync: loadScriptAsync,
			},
			EdgeTestParameters: fixture.getResponse.EdgeTestParameters,
			InjectionCriteria: appsec.ClientsideProtectionInjectionCriteria{
				InjectionPathIncludes: ptr.To("/include-path"),
				InjectionPathExcludes: ptr.To("/exclude-path"),
				InjectionPercent:      injectionPercent,
				InjectionPolicy:       ptr.To("head"),
				ScriptTagAttributes:   attrs,
			},
		},
	}
}
