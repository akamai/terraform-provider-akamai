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

func TestDataClientsideProtectionConfiguration(t *testing.T) {
	t.Parallel()

	baseConfig := testutils.LoadFixtureString(t, "testdata/TestDSClientsideProtectionConfiguration/clientside_protection_configuration.tf")
	missingConfigID := testutils.LoadFixtureString(t, "testdata/TestDSClientsideProtectionConfiguration/clientside_protection_configuration_missing_config_id.tf")
	missingSecurityPolicyID := testutils.LoadFixtureString(t, "testdata/TestDSClientsideProtectionConfiguration/clientside_protection_configuration_missing_security_policy_id.tf")

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                43253,
		LatestVersion:     7,
		StagingVersion:    6,
		ProductionVersion: 6,
	}

	cpcResponse := appsec.GetClientsideProtectionConfigurationResponse{
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

	serverError := appsec.Error{
		Type:       "internal_error",
		Title:      "Internal Server Error",
		Detail:     "Error fetching CPC configuration",
		StatusCode: http.StatusInternalServerError,
	}

	checker := test.NewStateChecker("data.akamai_appsec_clientside_protection_configuration.test").
		CheckEqual("config_id", "43253").
		CheckEqual("security_policy_id", "test_policy_id")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
					Return(&getConfigResponse, nil)
				m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
					ConfigID:         43253,
					Version:          7,
					SecurityPolicyID: "test_policy_id",
				}).Return(&cpcResponse, nil)
			},
			steps: []resource.TestStep{
				{
					Config: baseConfig,
					Check:  checker.Build(),
				},
			},
		},
		"happy path with script_tag_attributes": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
					Return(&getConfigResponse, nil)
				m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
					ConfigID:         43253,
					Version:          7,
					SecurityPolicyID: "test_policy_id",
				}).Return(&appsec.GetClientsideProtectionConfigurationResponse{
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
				}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: baseConfig,
					Check: test.NewStateChecker("data.akamai_appsec_clientside_protection_configuration.test").
						CheckEqual("injection_criteria.script_tag_attributes.#", "2").
						CheckEqual("injection_criteria.script_tag_attributes.0", `nonce="abc123"`).
						CheckEqual("injection_criteria.script_tag_attributes.1", `data-integrity="sha256"`).
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
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      missingSecurityPolicyID,
					ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required, but no definition was found`),
				},
			},
		},
		"error response from GetConfiguration": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
					Return(nil, &serverError)
			},
			steps: []resource.TestStep{
				{
					Config:      baseConfig,
					ExpectError: regexp.MustCompile("Title: Internal Server Error"),
				},
			},
		},
		"error response from GetClientsideProtectionConfiguration": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43253}).
					Return(&getConfigResponse, nil)
				m.On("GetClientsideProtectionConfiguration", mock.Anything, appsec.GetClientsideProtectionConfigurationRequest{
					ConfigID:         43253,
					Version:          7,
					SecurityPolicyID: "test_policy_id",
				}).Return(nil, &serverError)
			},
			steps: []resource.TestStep{
				{
					Config:      baseConfig,
					ExpectError: regexp.MustCompile("Error: calling 'GetClientsideProtectionConfiguration'"),
				},
			},
		},
	}

	for name, tc := range tests {
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
