package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomDeny_res_basic(t *testing.T) {
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		createResponse := appsec.CreateCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyCreateResponse.json"), &createResponse)
		require.NoError(t, err)
		createRequestJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheTrue.json")
		client.On("CreateCustomDeny",
			testutils.MockContext,
			appsec.CreateCustomDenyRequest{ConfigID: 43253, Version: 7, JsonPayloadRaw: createRequestJSON},
		).Return(&createResponse, nil)

		getResponse := appsec.GetCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyGetResponse.json"), &getResponse)
		require.NoError(t, err)
		client.On("GetCustomDeny",
			testutils.MockContext,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponse, nil).Times(3)

		updateRequestJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheFalse.json")
		updateResponse := appsec.UpdateCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyUpdateResponse.json"), &updateResponse)
		require.NoError(t, err)
		client.On("UpdateCustomDeny",
			testutils.MockContext,
			appsec.UpdateCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918", JsonPayloadRaw: updateRequestJSON},
		).Return(&updateResponse, nil)

		getResponseAfterUpdate := appsec.GetCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyGetResponseAfterUpdate.json"), &getResponseAfterUpdate)
		require.NoError(t, err)
		client.On("GetCustomDeny",
			testutils.MockContext,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponseAfterUpdate, nil).Twice()

		removeResponse := appsec.RemoveCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDeny.json"), &removeResponse)
		require.NoError(t, err)
		client.On("RemoveCustomDeny",
			testutils.MockContext,
			appsec.RemoveCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResCustomDeny/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
