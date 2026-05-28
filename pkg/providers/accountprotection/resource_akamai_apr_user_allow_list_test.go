package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceUserAllowList(t *testing.T) {
	t.Parallel()
	t.Run("TestResourceUserAllowList", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.AccountProtection.On("UpsertUserAllowListID",
			testutils.MockContext,
			apr.UpsertUserAllowListIDRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		client.AccountProtection.On("GetUserAllowListID",
			testutils.MockContext,
			apr.GetUserAllowListIDRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		client.AccountProtection.On("UpsertUserAllowListID",
			testutils.MockContext,
			apr.UpsertUserAllowListIDRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		client.AccountProtection.On("GetUserAllowListID",
			testutils.MockContext,
			apr.GetUserAllowListIDRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)

		// Add the mock for DeleteUserAllowListID
		client.AccountProtection.On("DeleteUserAllowListID",
			testutils.MockContext,
			apr.DeleteUserAllowListIDRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(nil).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
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

		client.AccountProtection.AssertExpectations(t)
	})
}
