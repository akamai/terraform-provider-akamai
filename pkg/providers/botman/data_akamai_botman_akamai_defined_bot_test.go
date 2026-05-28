package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataAkamaiDefinedBot(t *testing.T) {
	t.Parallel()
	t.Run("DataAkamaiDefinedBot", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client)

		response := botman.GetAkamaiDefinedBotListResponse{
			Bots: []map[string]interface{}{
				{"botId": "b85e3eaa-d334-466d-857e-33308ce416be", "botName": "Test name 1", "testKey": "testValue1"},
				{"botId": "69acad64-7459-4c1d-9bad-672600150127", "botName": "Test name 2", "testKey": "testValue2"},
				{"botId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName": "Test name 3", "testKey": "testValue3"},
				{"botId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "botName": "Test name 4", "testKey": "testValue4"},
				{"botId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "botName": "Test name 5", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"bots":[
		{"botId":"b85e3eaa-d334-466d-857e-33308ce416be", "botName": "Test name 1", "testKey":"testValue1"},
		{"botId":"69acad64-7459-4c1d-9bad-672600150127", "botName": "Test name 2", "testKey":"testValue2"},
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName": "Test name 3", "testKey":"testValue3"},
		{"botId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "botName": "Test name 4", "testKey":"testValue4"},
		{"botId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "botName": "Test name 5", "testKey":"testValue5"}
	]
}`
		client.BotMan.On("GetAkamaiDefinedBotList",
			testutils.MockContext,
			botman.GetAkamaiDefinedBotListRequest{},
		).Return(&response, nil)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataAkamaiDefinedBot/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_akamai_defined_bot.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
	t.Run("DataAkamaiDefinedBot filter by name", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client)

		response := botman.GetAkamaiDefinedBotListResponse{
			Bots: []map[string]interface{}{
				{"botId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName": "Test name 3", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"bots":[
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "botName": "Test name 3", "testKey":"testValue3"}
	]
}`
		client.BotMan.On("GetAkamaiDefinedBotList",
			testutils.MockContext,
			botman.GetAkamaiDefinedBotListRequest{BotName: "Test name 3"},
		).Return(&response, nil)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataAkamaiDefinedBot/filter_by_name.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_akamai_defined_bot.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
