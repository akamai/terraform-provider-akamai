package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCustomBotCategorySequence(t *testing.T) {
	t.Run("DataCustomBotCategorySequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.CustomBotCategorySequenceResponse{
			Sequence: []string{"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "d79285df-e399-43e8-bb0f-c0d980a88e4f", "afa309b8-4fd5-430e-a061-1c61df1d2ac2"},
		}
		mockedBotmanClient.On("GetCustomBotCategorySequence",
			testutils.MockContext,
			botman.GetCustomBotCategorySequenceRequest{ConfigID: 43253, Version: 15},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataCustomBotCategorySequence/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_sequence.test", "category_ids.#", "3"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_sequence.test", "category_ids.0", "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_sequence.test", "category_ids.1", "d79285df-e399-43e8-bb0f-c0d980a88e4f"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_sequence.test", "category_ids.2", "afa309b8-4fd5-430e-a061-1c61df1d2ac2")),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
