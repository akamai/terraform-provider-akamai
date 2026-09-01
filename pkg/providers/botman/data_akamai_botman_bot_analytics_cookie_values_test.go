package botman

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataBotAnalyticsCookieValue(t *testing.T) {
	t.Parallel()
	t.Run("DataBotAnalyticsCookieValues", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)

		response := map[string]interface{}{
			"values": []interface{}{
				map[string]interface{}{"testKey": "testValue1"},
				map[string]interface{}{"testKey": "testValue2"},
				map[string]interface{}{"testKey": "testValue3"},
				map[string]interface{}{"testKey": "testValue4"},
				map[string]interface{}{"testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"values": [
		{"testKey":"testValue1"},
		{"testKey":"testValue2"},
		{"testKey":"testValue3"},
		{"testKey":"testValue4"},
		{"testKey":"testValue5"}
	]
}`
		client.BotMan.On("GetBotAnalyticsCookieValues",
			testutils.MockContext,
		).Return(response, nil)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBotAnalyticsCookieValues/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_bot_analytics_cookie_values.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
