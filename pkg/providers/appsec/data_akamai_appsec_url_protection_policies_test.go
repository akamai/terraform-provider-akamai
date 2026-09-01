package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataUrlProtectionPolicies(t *testing.T) {
	t.Parallel()

	listURLProtectionPolicies := appsec.ListURLProtectionPoliciesResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSURLProtection/URLProtectionPolicies.json"), &listURLProtectionPolicies)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker("data.akamai_appsec_url_protection_policies.test")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return a list of url protection policies": {
			init: func(m *appsec.Mock) {
				mockGetConfig(m, 3)
				mockListURLProtectionPolicies(m, listURLProtectionPolicies, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/list_policies_match_by_id.tf"),
					Check: baseChecker.
						CheckEqual("url_protection_policies.#", "4").
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
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/list_policies_match_by_id.tf"),
					ExpectError: regexp.MustCompile("Error: invalid config version"),
				},
			},
		},
		"error response from ListURLProtectionPolicies api": {
			init: func(m *appsec.Mock) {
				mockGetConfig(m, 1)
				mockListURLProtectionPoliciesFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtection/list_policies_match_by_id.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'ListURLProtectionPolicies'"),
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

func mockListURLProtectionPolicies(m *appsec.Mock, response appsec.ListURLProtectionPoliciesResponse, times int) {
	m.On("ListURLProtectionPolicies", mock.Anything, appsec.ListURLProtectionPoliciesRequest{ConfigID: 43007, ConfigVersion: 40}).
		Return(&response, nil).Times(times)
}

func mockListURLProtectionPoliciesFailure(m *appsec.Mock, times int) {
	m.On("ListURLProtectionPolicies", mock.Anything, appsec.ListURLProtectionPoliciesRequest{ConfigID: 43007, ConfigVersion: 40}).
		Return(nil, &serverError).Times(times)
}
