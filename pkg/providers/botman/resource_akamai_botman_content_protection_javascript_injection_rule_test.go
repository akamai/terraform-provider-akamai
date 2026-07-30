package botman

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceContentProtectionJavaScriptInjectionRule(t *testing.T) {
	t.Parallel()
	t.Run("ResourceContentProtectionJavaScriptInjectionRule", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createResponse := map[string]interface{}{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.BotMan.On("CreateContentProtectionJavaScriptInjectionRule",
			testutils.MockContext,
			botman.CreateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetContentProtectionJavaScriptInjectionRule",
			testutils.MockContext,
			botman.GetContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"contentProtectionJavaScriptInjectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"contentProtectionJavaScriptInjectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7","testKey":"updated_testValue3"}`
		client.BotMan.On("UpdateContentProtectionJavaScriptInjectionRule",
			testutils.MockContext,
			botman.UpdateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(updateRequest),
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetContentProtectionJavaScriptInjectionRule",
			testutils.MockContext,
			botman.GetContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		client.BotMan.On("RemoveContentProtectionJavaScriptInjectionRule",
			testutils.MockContext,
			botman.RemoveContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				ContentProtectionJavaScriptInjectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionJavaScriptInjectionRule/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_content_protection_javascript_injection_rule.test", "id", "43253:AAAA_81230:fake3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_javascript_injection_rule.test", "content_protection_javascript_injection_rule", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionJavaScriptInjectionRule/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_content_protection_javascript_injection_rule.test", "id", "43253:AAAA_81230:fake3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_javascript_injection_rule.test", "content_protection_javascript_injection_rule", expectedUpdateJSON)),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("ResourceContentProtectionJavaScriptInjectionRule missing required fields", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionJavaScriptInjectionRule/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionJavaScriptInjectionRule/missing_policy_id.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("ResourceContentProtectionJavaScriptInjectionRule error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.BotMan.On("CreateContentProtectionJavaScriptInjectionRule",
			testutils.MockContext,
			botman.CreateContentProtectionJavaScriptInjectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
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
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionJavaScriptInjectionRule/create.tf"),
					ExpectError: regexp.MustCompile("Title: Internal Server Error; Type: internal_error; Detail: Error fetching data"),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
