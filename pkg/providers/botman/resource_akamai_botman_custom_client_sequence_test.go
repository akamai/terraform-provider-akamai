package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceCustomClientSequence(t *testing.T) {
	t.Run("ResourceCustomClientSequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createCustomClientIds := []string{"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "d79285df-e399-43e8-bb0f-c0d980a88e4f", "afa309b8-4fd5-430e-a061-1c61df1d2ac2"}
		createResponse := botman.CustomClientSequenceResponse{Sequence: createCustomClientIds}
		mockedBotmanClient.On("UpdateCustomClientSequence",
			mock.Anything,
			botman.UpdateCustomClientSequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: createCustomClientIds,
			},
		).Return(&createResponse, nil).Once()

		mockedBotmanClient.On("GetCustomClientSequence",
			mock.Anything,
			botman.GetCustomClientSequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(&createResponse, nil).Times(3)

		updateCustomClientIds := []string{createCustomClientIds[1], createCustomClientIds[2], createCustomClientIds[0]}
		updateResponse := botman.CustomClientSequenceResponse{Sequence: updateCustomClientIds}
		mockedBotmanClient.On("UpdateCustomClientSequence",
			mock.Anything,
			botman.UpdateCustomClientSequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: updateCustomClientIds,
			},
		).Return(&updateResponse, nil).Once()

		mockedBotmanClient.On("GetCustomClientSequence",
			mock.Anything,
			botman.GetCustomClientSequenceRequest{
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
						Config: test.Fixture("testdata/TestResourceCustomClientSequence/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.#", tools.ConvertToString(len(createCustomClientIds))),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.0", createCustomClientIds[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.1", createCustomClientIds[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.2", createCustomClientIds[2])),
					},
					{
						Config: test.Fixture("testdata/TestResourceCustomClientSequence/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.#", tools.ConvertToString(len(updateCustomClientIds))),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.0", updateCustomClientIds[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.1", updateCustomClientIds[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.2", updateCustomClientIds[2])),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
