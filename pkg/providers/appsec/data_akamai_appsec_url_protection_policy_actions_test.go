package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestDataURLProtectionPolicyActions(t *testing.T) {
	t.Parallel()

	getURLProtectionPolicyActionsResponse := appsec.GetURLProtectionPolicyActionsResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSURLProtectionPolicyActions/URLProtectionPolicyActions.json"), &getURLProtectionPolicyActionsResponse)
	require.NoError(t, err)

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return URL protection policy actions": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationURLProtectionPolicyActions(m, 3)
				mockGetURLProtectionPolicyActionsDS(m, getURLProtectionPolicyActionsResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPolicyActions/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policy_actions.test", "max_rate_threshold_action", "alert"),
						resource.TestCheckResourceAttr("data.akamai_appsec_url_protection_policy_actions.test", "load_shedding_action", "none"),
					),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPolicyActions/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPolicyActions/missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument url_protection_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPolicyActions/missing_url_protection_policy_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "url_protection_policy_id" is required`),
				},
			},
		},
		"error when url_protection_policy_id not found": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationURLProtectionPolicyActions(m, 1)
				m.On("GetURLProtectionPolicyActions", testutils.MockContext, appsec.GetURLProtectionPolicyActionsRequest{
					ConfigID:              43253,
					ConfigVersion:         7,
					PolicyID:              "AAAA_81230",
					URLProtectionPolicyID: 999999,
				}).Return(nil, &appsec.Error{StatusCode: 404, Title: "Not Found"}).Times(1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSURLProtectionPolicyActions/not_found.tf"),
					ExpectError: regexp.MustCompile(`Read URL Protection Policy Actions failed`),
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

func mockGetURLProtectionPolicyActionsDS(m *appsec.Mock, response appsec.GetURLProtectionPolicyActionsResponse, times int) {
	m.On("GetURLProtectionPolicyActions", testutils.MockContext, appsec.GetURLProtectionPolicyActionsRequest{
		ConfigID:              43253,
		ConfigVersion:         7,
		PolicyID:              "AAAA_81230",
		URLProtectionPolicyID: 135355,
	}).Return(&response, nil).Times(times)
}

func mockGetConfigurationURLProtectionPolicyActions(m *appsec.Mock, times int) {
	m.On("GetConfiguration", testutils.MockContext, appsec.GetConfigurationRequest{ConfigID: 43253}).
		Return(&appsec.GetConfigurationResponse{
			FileType:          "RBAC",
			ID:                43253,
			LatestVersion:     7,
			Name:              "Test Config",
			ProductionVersion: 7,
			StagingVersion:    7,
			TargetProduct:     "KSD",
		}, nil).Times(times)
}
