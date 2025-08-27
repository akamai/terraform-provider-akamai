package botman

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataContentProtectionRuleSequence(t *testing.T) {
	t.Run("DataContentProtectionRuleSequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetContentProtectionRuleSequenceResponse{
			ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"},
		}
		mockedBotmanClient.On("GetContentProtectionRuleSequence",
			testutils.MockContext,
			botman.GetContentProtectionRuleSequenceRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionRuleSequence/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.#", "3"),
							resource.TestCheckResourceAttr("data.akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.0", "fake3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("data.akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.1", "fake85df-e399-43e8-bb0f-c0d980a88e4f"),
							resource.TestCheckResourceAttr("data.akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.2", "fake09b8-4fd5-430e-a061-1c61df1d2ac2")),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})

	t.Run("DataContentProtectionRuleSequenceError", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}
		mockedBotmanClient.On("GetContentProtectionRuleSequence",
			testutils.MockContext,
			botman.GetContentProtectionRuleSequenceRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
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
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionRuleSequence/basic.tf"),
						ExpectError: regexp.MustCompile("Title: Internal Server Error; Type: internal_error; Detail: Error fetching data"),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})

	t.Run("DataContentProtectionRuleSequenceMissingRequiredFields", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}
		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionRuleSequence/missing_config_id.tf"),
						ExpectError: regexp.MustCompile(`Error: Missing required argument`),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataContentProtectionRuleSequence/missing_policy_id.tf"),
						ExpectError: regexp.MustCompile(`Error: Missing required argument`),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
