package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataBotAnalyticsCookie(t *testing.T) {
	t.Run("DataBotAnalyticsCookie", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		mockedBotmanClient.On("GetBotAnalyticsCookie",
			mock.Anything,
			botman.GetBotAnalyticsCookieRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataBotAnalyticsCookie/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_analytics_cookie.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
