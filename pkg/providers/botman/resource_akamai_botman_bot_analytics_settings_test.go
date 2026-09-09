package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceBotAnalyticsSettings(t *testing.T) {
	t.Parallel()
	t.Run("ResourceBotAnalyticsSettings", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.BotMan.On("UpdateBotAnalyticsSettings",
			testutils.MockContext,
			botman.UpdateBotAnalyticsSettingsRequest{
				ConfigID:    43253,
				Version:     15,
				JSONPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetBotAnalyticsSettings",
			testutils.MockContext,
			botman.GetBotAnalyticsSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		client.BotMan.On("UpdateBotAnalyticsSettings",
			testutils.MockContext,
			botman.UpdateBotAnalyticsSettingsRequest{
				ConfigID:    43253,
				Version:     15,
				JSONPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetBotAnalyticsSettings",
			testutils.MockContext,
			botman.GetBotAnalyticsSettingsRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotAnalyticsSettings/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_analytics_settings.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_botman_bot_analytics_settings.test", "bot_analytics_settings", expectedCreateJSON),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotAnalyticsSettings/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_analytics_settings.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_botman_bot_analytics_settings.test", "bot_analytics_settings", expectedUpdateJSON),
					),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
