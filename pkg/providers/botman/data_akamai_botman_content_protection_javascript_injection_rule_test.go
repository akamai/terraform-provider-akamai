package botman

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataContentProtectionJavaScriptInjectionRule(t *testing.T) {
	t.Parallel()
	t.Run("DataContentProtectionJavaScriptInjectionRule", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		response := botman.GetContentProtectionJavaScriptInjectionRuleListResponse{
			ContentProtectionJavaScriptInjectionRules: []map[string]interface{}{
				{"contentProtectionJavaScriptInjectionRuleId": "fake3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
				{"contentProtectionJavaScriptInjectionRuleId": "fakead64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
				{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				{"contentProtectionJavaScriptInjectionRuleId": "fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
				{"contentProtectionJavaScriptInjectionRuleId": "faked85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"contentProtectionJavaScriptInjectionRules":[
		{"contentProtectionJavaScriptInjectionRuleId":"fake3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"contentProtectionJavaScriptInjectionRuleId":"fakead64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"contentProtectionJavaScriptInjectionRuleId":"fake4ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"contentProtectionJavaScriptInjectionRuleId":"faked85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`
		client.BotMan.On("GetContentProtectionJavaScriptInjectionRuleList",
			testutils.MockContext,
			botman.GetContentProtectionJavaScriptInjectionRuleListRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(&response, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionJavaScriptInjectionRule/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_content_protection_javascript_injection_rule.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("DataContentProtectionJavaScriptInjectionRule filter by id", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		response := botman.GetContentProtectionJavaScriptInjectionRuleListResponse{
			ContentProtectionJavaScriptInjectionRules: []map[string]interface{}{
				{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"contentProtectionJavaScriptInjectionRules":[
		{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
	]
}`
		client.BotMan.On("GetContentProtectionJavaScriptInjectionRuleList",
			testutils.MockContext,
			botman.GetContentProtectionJavaScriptInjectionRuleListRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230", ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7"},
		).Return(&response, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionJavaScriptInjectionRule/filter_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_content_protection_javascript_injection_rule.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("DataContentProtectionJavaScriptInjectionRule error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		client.BotMan.On("GetContentProtectionJavaScriptInjectionRuleList",
			testutils.MockContext,
			botman.GetContentProtectionJavaScriptInjectionRuleListRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(nil, &botman.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error fetching data",
			StatusCode: http.StatusInternalServerError,
		}).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionJavaScriptInjectionRule/basic.tf"),
					ExpectError: regexp.MustCompile("Title: Internal Server Error; Type: internal_error; Detail: Error fetching data"),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("DataContentProtectionJavaScriptInjectionRule missing required fields", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionJavaScriptInjectionRule/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionJavaScriptInjectionRule/missing_policy_id.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
