package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCustomClientSequence(t *testing.T) {
	t.Parallel()
	t.Run("ResourceCustomClientSequence", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createCustomClientIDs := []string{"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "d79285df-e399-43e8-bb0f-c0d980a88e4f", "afa309b8-4fd5-430e-a061-1c61df1d2ac2"}
		createResponse := botman.CustomClientSequenceResponse{Sequence: createCustomClientIDs}
		client.BotMan.On("UpdateCustomClientSequence",
			testutils.MockContext,
			botman.UpdateCustomClientSequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: createCustomClientIDs,
			},
		).Return(&createResponse, nil).Once()

		client.BotMan.On("GetCustomClientSequence",
			testutils.MockContext,
			botman.GetCustomClientSequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(&createResponse, nil).Times(3)

		updateCustomClientIDs := []string{createCustomClientIDs[1], createCustomClientIDs[2], createCustomClientIDs[0]}
		updateResponse := botman.CustomClientSequenceResponse{Sequence: updateCustomClientIDs}
		client.BotMan.On("UpdateCustomClientSequence",
			testutils.MockContext,
			botman.UpdateCustomClientSequenceRequest{
				ConfigID: 43253,
				Version:  15,
				Sequence: updateCustomClientIDs,
			},
		).Return(&updateResponse, nil).Once()

		client.BotMan.On("GetCustomClientSequence",
			testutils.MockContext,
			botman.GetCustomClientSequenceRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(&updateResponse, nil).Times(2)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomClientSequence/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.#", str.From(len(createCustomClientIDs))),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.0", createCustomClientIDs[0]),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.1", createCustomClientIDs[1]),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.2", createCustomClientIDs[2])),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomClientSequence/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.#", str.From(len(updateCustomClientIDs))),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.0", updateCustomClientIDs[0]),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.1", updateCustomClientIDs[1]),
						resource.TestCheckResourceAttr("akamai_botman_custom_client_sequence.test", "custom_client_ids.2", updateCustomClientIDs[2])),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
