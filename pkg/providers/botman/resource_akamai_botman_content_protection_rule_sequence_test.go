package botman

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceContentProtectionRuleSequence(t *testing.T) {
	t.Parallel()
	t.Run("ResourceContentProtectionRuleSequence", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createContentProtectionRuleIDs := botman.ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"}}
		updateContentProtectionRuleIDs := botman.ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"}}
		createResponse := botman.UpdateContentProtectionRuleSequenceResponse(createContentProtectionRuleIDs)
		readResponse := botman.GetContentProtectionRuleSequenceResponse(createContentProtectionRuleIDs)
		updateResponse := botman.UpdateContentProtectionRuleSequenceResponse(updateContentProtectionRuleIDs)
		readResponseAfterUpdate := botman.GetContentProtectionRuleSequenceResponse(updateContentProtectionRuleIDs)
		client.BotMan.On("UpdateContentProtectionRuleSequence",
			testutils.MockContext,
			botman.UpdateContentProtectionRuleSequenceRequest{
				ConfigID:                      43253,
				Version:                       15,
				SecurityPolicyID:              "AAAA_81230",
				ContentProtectionRuleSequence: createContentProtectionRuleIDs,
			},
		).Return(&createResponse, nil).Once()

		client.BotMan.On("GetContentProtectionRuleSequence",
			testutils.MockContext,
			botman.GetContentProtectionRuleSequenceRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(&readResponse, nil).Times(3)

		client.BotMan.On("UpdateContentProtectionRuleSequence",
			testutils.MockContext,
			botman.UpdateContentProtectionRuleSequenceRequest{
				ConfigID:                      43253,
				Version:                       15,
				SecurityPolicyID:              "AAAA_81230",
				ContentProtectionRuleSequence: updateContentProtectionRuleIDs,
			},
		).Return(&updateResponse, nil).Once()

		client.BotMan.On("GetContentProtectionRuleSequence",
			testutils.MockContext,
			botman.GetContentProtectionRuleSequenceRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(&readResponseAfterUpdate, nil).Times(2)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRuleSequence/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.#", "3"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.0", "fake3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.1", "fake85df-e399-43e8-bb0f-c0d980a88e4f"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.2", "fake09b8-4fd5-430e-a061-1c61df1d2ac2")),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRuleSequence/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.#", "3"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.0", "fake85df-e399-43e8-bb0f-c0d980a88e4f"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.1", "fake3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_content_protection_rule_sequence.test", "content_protection_rule_ids.2", "fake09b8-4fd5-430e-a061-1c61df1d2ac2")),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("ResourceContentProtectionRuleSequence missing required fields", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRuleSequence/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRuleSequence/missing_policy_id.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("ResourceContentProtectionRuleSequence error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createContentProtectionRuleIDs := botman.ContentProtectionRuleUUIDSequence{ContentProtectionRuleSequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"}}
		client.BotMan.On("UpdateContentProtectionRuleSequence",
			testutils.MockContext,
			botman.UpdateContentProtectionRuleSequenceRequest{
				ConfigID:                      43253,
				Version:                       15,
				SecurityPolicyID:              "AAAA_81230",
				ContentProtectionRuleSequence: createContentProtectionRuleIDs,
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
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceContentProtectionRuleSequence/create.tf"),
					ExpectError: regexp.MustCompile("Title: Internal Server Error; Type: internal_error; Detail: Error fetching data"),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
