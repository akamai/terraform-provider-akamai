package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCustomBotCategorySequence(t *testing.T) {
	t.Run("ResourceCustomBotCategorySequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createCategoryIDs := []string{"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "d79285df-e399-43e8-bb0f-c0d980a88e4f", "afa309b8-4fd5-430e-a061-1c61df1d2ac2"}
		createResponse := botman.CustomBotCategorySequenceResponse{Sequence: createCategoryIDs}
		mockedBotmanClient.On("UpdateCustomBotCategorySequence",
			testutils.MockContext,
			botman.UpdateCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: createCategoryIDs,
			},
		).Return(&createResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategorySequence",
			testutils.MockContext,
			botman.GetCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(&createResponse, nil).Times(3)

		updateCategoryIDs := []string{createCategoryIDs[1], createCategoryIDs[2], createCategoryIDs[0]}
		updateResponse := botman.CustomBotCategorySequenceResponse{Sequence: updateCategoryIDs}
		mockedBotmanClient.On("UpdateCustomBotCategorySequence",
			testutils.MockContext,
			botman.UpdateCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: updateCategoryIDs,
			},
		).Return(&updateResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategorySequence",
			testutils.MockContext,
			botman.GetCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(&updateResponse, nil).Times(2)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategorySequence/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.#", str.From(len(createCategoryIDs))),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.0", createCategoryIDs[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.1", createCategoryIDs[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.2", createCategoryIDs[2])),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategorySequence/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.#", str.From(len(updateCategoryIDs))),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.0", updateCategoryIDs[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.1", updateCategoryIDs[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.2", updateCategoryIDs[2])),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
