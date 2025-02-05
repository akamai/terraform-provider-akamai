package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceJavascriptInjection(t *testing.T) {
	t.Run("ResourceJavascriptInjection", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedBotmanClient.On("UpdateJavascriptInjection",
			testutils.MockContext,
			botman.UpdateJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetJavascriptInjection",
			testutils.MockContext,
			botman.GetJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		mockedBotmanClient.On("UpdateJavascriptInjection",
			testutils.MockContext,
			botman.UpdateJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      updateRequest,
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetJavascriptInjection",
			testutils.MockContext,
			botman.GetJavascriptInjectionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceJavascriptInjection/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_javascript_injection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_botman_javascript_injection.test", "javascript_injection", expectedCreateJSON)),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceJavascriptInjection/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_javascript_injection.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_botman_javascript_injection.test", "javascript_injection", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
