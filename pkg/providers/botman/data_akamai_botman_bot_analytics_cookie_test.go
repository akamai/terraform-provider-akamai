package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataBotAnalyticsCookie(t *testing.T) {
	t.Parallel()
	t.Run("DataBotAnalyticsCookie", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		client.BotMan.On("GetBotAnalyticsCookie",
			testutils.MockContext,
			botman.GetBotAnalyticsCookieRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBotAnalyticsCookie/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_bot_analytics_cookie.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
