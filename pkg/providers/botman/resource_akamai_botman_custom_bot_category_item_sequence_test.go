package botman

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCustomBotCategoryItemSequenceMissingInput(t *testing.T) {
	t.Run("ResourceCustomBotCategoryItemSequenceNoId", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategoryItemSequence/missing_config_id.tf"),
						ExpectError: regexp.MustCompile(`Error: Missing required argument`),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}

func TestResourceCustomBotCategoryItemSequenceError(t *testing.T) {
	t.Run("ResourceCustomBotCategoryItemSequenceError", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createCategoryIDs := botman.UUIDSequence{Sequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"}}
		mockedBotmanClient.On("UpdateCustomBotCategoryItemSequence",
			testutils.MockContext,
			botman.UpdateCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "fakecv20-eddb-4421-93d9-90954e509d5f",
				Sequence:   createCategoryIDs,
			},
		).Return(nil, &botman.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error fetching data",
			StatusCode: http.StatusInternalServerError,
		}).Once()

		useClient(mockedBotmanClient, func() {
			resource.Test(t, resource.TestCase{
				ErrorCheck: func(err error) error {
					print(err)
					return err
				},
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategoryItemSequence/create.tf"),
						ExpectError: regexp.MustCompile("Title: Internal Server Error; Type: internal_error; Detail: Error fetching data"),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}

func TestResourceCustomBotCategoryItemSequence(t *testing.T) {
	t.Run("ResourceCustomBotCategoryItemSequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createCategoryIDs := botman.UUIDSequence{Sequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"}}
		updateCategoryIDs := botman.UUIDSequence{Sequence: []string{createCategoryIDs.Sequence[1], createCategoryIDs.Sequence[0], createCategoryIDs.Sequence[2]}}
		createResponse := botman.UpdateCustomBotCategoryItemSequenceResponse(createCategoryIDs)
		readResponse := botman.GetCustomBotCategoryItemSequenceResponse(createCategoryIDs)
		updateResponse := botman.UpdateCustomBotCategoryItemSequenceResponse(updateCategoryIDs)
		readResponse2 := botman.GetCustomBotCategoryItemSequenceResponse(updateCategoryIDs)
		mockedBotmanClient.On("UpdateCustomBotCategoryItemSequence",
			testutils.MockContext,
			botman.UpdateCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "fakecv20-eddb-4421-93d9-90954e509d5f",
				Sequence:   createCategoryIDs,
			},
		).Return(&createResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategoryItemSequence",
			testutils.MockContext,
			botman.GetCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "fakecv20-eddb-4421-93d9-90954e509d5f",
			},
		).Return(&readResponse, nil).Times(3)

		mockedBotmanClient.On("UpdateCustomBotCategoryItemSequence",
			testutils.MockContext,
			botman.UpdateCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "fakecv20-eddb-4421-93d9-90954e509d5f",
				Sequence:   updateCategoryIDs,
			},
		).Return(&updateResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategoryItemSequence",
			testutils.MockContext,
			botman.GetCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				CategoryID: "fakecv20-eddb-4421-93d9-90954e509d5f",
				Version:    15,
			},
		).Return(&readResponse2, nil).Times(2)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategoryItemSequence/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "id", "43253:fakecv20-eddb-4421-93d9-90954e509d5f"),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.#", str.From(len(createCategoryIDs.Sequence))),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.0", createCategoryIDs.Sequence[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.1", createCategoryIDs.Sequence[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.2", createCategoryIDs.Sequence[2])),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategoryItemSequence/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "id", "43253:fakecv20-eddb-4421-93d9-90954e509d5f"),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.#", str.From(len(updateCategoryIDs.Sequence))),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.0", updateCategoryIDs.Sequence[0]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.1", updateCategoryIDs.Sequence[1]),
							resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.2", updateCategoryIDs.Sequence[2])),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
