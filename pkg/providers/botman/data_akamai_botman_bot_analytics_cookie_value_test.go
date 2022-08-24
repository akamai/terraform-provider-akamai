package botman

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataBotAnalyticsCookieValue(t *testing.T) {
	t.Run("DataBotAnalyticsCookieValue", func(t *testing.T) {
		mockedBotmanClient := &mockbotman{}

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
		mockedBotmanClient.On("GetBotAnalyticsCookieValues",
			mock.Anything,
		).Return(response, nil)
		useClient(mockedBotmanClient, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataBotAnalyticsCookieValue/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_analytics_cookie_value.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
