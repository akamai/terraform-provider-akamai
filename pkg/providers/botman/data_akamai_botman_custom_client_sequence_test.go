package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataCustomClientSequence(t *testing.T) {
	t.Run("DataCustomClientSequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.CustomClientSequenceResponse{
			Sequence: []string{"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "d79285df-e399-43e8-bb0f-c0d980a88e4f", "afa309b8-4fd5-430e-a061-1c61df1d2ac2"},
		}
		mockedBotmanClient.On("GetCustomClientSequence",
			mock.Anything,
			botman.GetCustomClientSequenceRequest{ConfigID: 43253, Version: 15},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataCustomClientSequence/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_custom_client_sequence.test", "custom_client_ids.#", "3"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_client_sequence.test", "custom_client_ids.0", "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_client_sequence.test", "custom_client_ids.1", "d79285df-e399-43e8-bb0f-c0d980a88e4f"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_client_sequence.test", "custom_client_ids.2", "afa309b8-4fd5-430e-a061-1c61df1d2ac2")),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
