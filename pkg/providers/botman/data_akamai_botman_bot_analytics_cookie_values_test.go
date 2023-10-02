package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataBotAnalyticsCookieValue(t *testing.T) {
	t.Run("DataBotAnalyticsCookieValues", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}

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
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataBotAnalyticsCookieValues/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_analytics_cookie_values.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
