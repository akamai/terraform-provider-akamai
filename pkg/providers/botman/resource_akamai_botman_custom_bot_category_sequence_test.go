package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/test"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceCustomBotCategorySequence(t *testing.T) {
	t.Run("ResourceCustomBotCategorySequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createCategoryIds := []string{"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "d79285df-e399-43e8-bb0f-c0d980a88e4f", "afa309b8-4fd5-430e-a061-1c61df1d2ac2"}
		createResponse := botman.CustomBotCategorySequenceResponse{Sequence: createCategoryIds}
		mockedBotmanClient.On("UpdateCustomBotCategorySequence",
			mock.Anything,
			botman.UpdateCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: createCategoryIds,
			},
		).Return(&createResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategorySequence",
			mock.Anything,
			botman.GetCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(&createResponse, nil).Times(3)

		updateCategoryIds := []string{createCategoryIds[1], createCategoryIds[2], createCategoryIds[0]}
		updateResponse := botman.CustomBotCategorySequenceResponse{Sequence: updateCategoryIds}
		mockedBotmanClient.On("UpdateCustomBotCategorySequence",
			mock.Anything,
			botman.UpdateCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: updateCategoryIds,
			},
		).Return(&updateResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategorySequence",
			mock.Anything,
			botman.GetCustomBotCategorySequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(&updateResponse, nil).Times(2)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestResourceCustomBotCategorySequence/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.#", tools.ConvertToString(len(createCategoryIds))),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.0", createCategoryIds[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.1", createCategoryIds[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.2", createCategoryIds[2])),
					},
					{
						Config: test.Fixture("testdata/TestResourceCustomBotCategorySequence/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.#", tools.ConvertToString(len(updateCategoryIds))),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.0", updateCategoryIds[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.1", updateCategoryIds[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_sequence.test", "category_ids.2", updateCategoryIds[2])),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
