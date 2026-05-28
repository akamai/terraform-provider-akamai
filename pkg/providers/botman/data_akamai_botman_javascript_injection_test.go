package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataJavascriptInjection(t *testing.T) {
	t.Parallel()
	t.Run("DataJavascriptInjection", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client)
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		client.BotMan.On("GetJavascriptInjection",
			testutils.MockContext,
			botman.GetJavascriptInjectionRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(response, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataJavascriptInjection/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_botman_javascript_injection.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
