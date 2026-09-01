package botman

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataBotAnalyticsSettings(t *testing.T) {
	t.Parallel()
	t.Run("DataBotAnalyticsSettings", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		client.BotMan.On("GetBotAnalyticsSettings",
			testutils.MockContext,
			botman.GetBotAnalyticsSettingsRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataBotAnalyticsSettings/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_bot_analytics_settings.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})

	t.Run("DataBotAnalyticsSettings wrong configID", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		client.APPSEC.On("GetConfiguration", testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 99999},
		).Return(nil, &appsec.Error{
			Type:       "not_found",
			Title:      "Not Found",
			Detail:     "Config ID not found",
			StatusCode: http.StatusNotFound,
		}).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBotAnalyticsSettings/wrong_config_id.tf"),
					ExpectError: regexp.MustCompile("Title: Not Found; Type: not_found; Detail: Config ID not found"),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("DataBotAnalyticsSettings missing config_id", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataBotAnalyticsSettings/missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`Error: Missing required argument`),
				},
			},
		})
	})
}
