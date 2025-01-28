package botman

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCustomBotCategoryItemSequenceError(t *testing.T) {
	t.Run("DataCustomBotCategoryItemSequenceError", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		mockedBotmanClient.On("GetCustomBotCategoryItemSequence",
			testutils.MockContext,
			botman.GetCustomBotCategoryItemSequenceRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "fakecv20-eddb-4421-93d9-90954e509d5f",
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
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataCustomBotCategoryItemSequence/basic.tf"),
						ExpectError: regexp.MustCompile("Title: Internal Server Error; Type: internal_error; Detail: Error fetching data"),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}

func TestDataCustomBotCategoryItemSequenceMissingInput(t *testing.T) {
	t.Run("DataCustomBotCategoryItemSequenceNoId", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataCustomBotCategoryItemSequence/missing_config_id.tf"),
						ExpectError: regexp.MustCompile(`Error: Missing required argument`),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}

func TestDataCustomBotCategoryItemSequence(t *testing.T) {
	t.Run("DataCustomBotCategoryItemSequence", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetCustomBotCategoryItemSequenceResponse{
			Sequence: []string{"fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"},
		}
		mockedBotmanClient.On("GetCustomBotCategoryItemSequence",
			testutils.MockContext,
			botman.GetCustomBotCategoryItemSequenceRequest{ConfigID: 43253, Version: 15, CategoryID: "fakecv20-eddb-4421-93d9-90954e509d5f"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataCustomBotCategoryItemSequence/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.#", "3"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.0", "fake3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.1", "fake85df-e399-43e8-bb0f-c0d980a88e4f"),
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_item_sequence.test", "bot_ids.2", "fake09b8-4fd5-430e-a061-1c61df1d2ac2")),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
