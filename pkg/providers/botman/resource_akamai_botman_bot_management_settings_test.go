package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceBotManagementSettings(t *testing.T) {
	t.Run("ResourceBotManagementSettings", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := test.FixtureBytes("testdata/JsonPayload/create.json")
		mockedBotmanClient.On("UpdateBotManagementSetting",
			mock.Anything,
			botman.UpdateBotManagementSettingRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetBotManagementSetting",
			mock.Anything,
			botman.GetBotManagementSettingRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := test.FixtureBytes("testdata/JsonPayload/update.json")
		mockedBotmanClient.On("UpdateBotManagementSetting",
			mock.Anything,
			botman.UpdateBotManagementSettingRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      updateRequest,
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetBotManagementSetting",
			mock.Anything,
			botman.GetBotManagementSettingRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestResourceBotManagementSettings/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_bot_management_settings.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_botman_bot_management_settings.test", "bot_management_settings", expectedCreateJSON)),
					},
					{
						Config: test.Fixture("testdata/TestResourceBotManagementSettings/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_bot_management_settings.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_botman_bot_management_settings.test", "bot_management_settings", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
