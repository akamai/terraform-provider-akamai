package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataUrlProtectionRules(t *testing.T) {

	listURLProtectionRules := appsec.ListURLProtectionRulesResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSURLProtection/URLProtectionRules.json"), &listURLProtectionRules)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker("data.akamai_appsec_url_protection_rules.test")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return a list of url protection rules": {
			init: func(m *appsec.Mock) {
				mockGetConfig(m, 3)
				mockListURLProtectionRules(m, listURLProtectionRules, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/list_rules_match_by_id.tf"),
					Check: baseChecker.
						CheckEqual("url_protection_rules.#", "4").
						Build(),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/config_missing_match_by_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"error response from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/list_rules_match_by_id.tf"),
					ExpectError: regexp.MustCompile("Error: invalid config version"),
				},
			},
		},
		"error response from ListURLProtectionRules api": {
			init: func(m *appsec.Mock) {
				mockGetConfig(m, 1)
				mockListURLProtectionRulesFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/list_rules_match_by_id.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'ListURLProtectionRules'"),
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

func mockListURLProtectionRules(m *appsec.Mock, response appsec.ListURLProtectionRulesResponse, times int) {
	m.On("ListURLProtectionRules", mock.Anything, appsec.ListURLProtectionRulesRequest{ConfigID: 43007, ConfigVersion: 40}).
		Return(&response, nil).Times(times)
}

func mockListURLProtectionRulesFailure(m *appsec.Mock, times int) {
	m.On("ListURLProtectionRules", mock.Anything, appsec.ListURLProtectionRulesRequest{ConfigID: 43007, ConfigVersion: 40}).
		Return(nil, &serverError).Times(times)
}
