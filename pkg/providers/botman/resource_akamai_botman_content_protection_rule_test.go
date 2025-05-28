package botman

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceContentProtectionRule(t *testing.T) {
	t.Run("ResourceContentProtectionRule", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"contentProtectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedBotmanClient.On("CreateContentProtectionRule",
			testutils.MockContext,
			botman.CreateContentProtectionRuleRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetContentProtectionRule",
			testutils.MockContext,
			botman.GetContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"contentProtectionRuleId": "fake3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"contentProtectionRuleId":"fake3f89-e179-4892-89cf-d5e623ba9dc7","testKey":"updated_testValue3"}`
		mockedBotmanClient.On("UpdateContentProtectionRule",
			testutils.MockContext,
			botman.UpdateContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:             json.RawMessage(updateRequest),
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetContentProtectionRule",
			testutils.MockContext,
			botman.GetContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		mockedBotmanClient.On("RemoveContentProtectionRule",
			testutils.MockContext,
			botman.RemoveContentProtectionRuleRequest{
				ConfigID:                43253,
				Version:                 15,
				SecurityPolicyID:        "AAAA_81230",
				ContentProtectionRuleID: "fake3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRule/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_content_protection_rule.test", "id", "43253:AAAA_81230:fake3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_content_protection_rule.test", "content_protection_rule", expectedCreateJSON)),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRule/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_content_protection_rule.test", "id", "43253:AAAA_81230:fake3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_content_protection_rule.test", "content_protection_rule", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})

	t.Run("ResourceContentProtectionRule missing required fields", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}
		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRule/missing_config_id.tf"),
						ExpectError: regexp.MustCompile(`Error: Missing required argument`),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRule/missing_policy_id.tf"),
						ExpectError: regexp.MustCompile(`Error: Missing required argument`),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})

	t.Run("ResourceContentProtectionRule error", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedBotmanClient.On("CreateContentProtectionRule",
			testutils.MockContext,
			botman.CreateContentProtectionRuleRequest{
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

		useClient(mockedBotmanClient, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRule/create.tf"),
						ExpectError: regexp.MustCompile("Title: Internal Server Error; Type: internal_error; Detail: Error fetching data"),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
