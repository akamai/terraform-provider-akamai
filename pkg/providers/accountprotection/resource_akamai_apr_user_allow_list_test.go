package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceUserAllowList(t *testing.T) {
	t.Run("TestResourceUserAllowList", func(t *testing.T) {

		mockedAprClient := &apr.Mock{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedAprClient.On("UpsertUserAllowListID",
			testutils.MockContext,
			apr.UpsertUserAllowListIDRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedAprClient.On("GetUserAllowListID",
			testutils.MockContext,
			apr.GetUserAllowListIDRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		mockedAprClient.On("UpsertUserAllowListID",
			testutils.MockContext,
			apr.UpsertUserAllowListIDRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		mockedAprClient.On("GetUserAllowListID",
			testutils.MockContext,
			apr.GetUserAllowListIDRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)

		// Add the mock for DeleteUserAllowListID
		mockedAprClient.On("DeleteUserAllowListID",
			testutils.MockContext,
			apr.DeleteUserAllowListIDRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(nil).Once()

		useClient(mockedAprClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceAprUserAllowList/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_apr_user_allow_list.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_apr_user_allow_list.test", "user_allow_list", expectedCreateJSON)),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceAprUserAllowList/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_apr_user_allow_list.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_apr_user_allow_list.test", "user_allow_list", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedAprClient.AssertExpectations(t)
	})
}
