package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomDeny_res_basic(t *testing.T) {
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		client := &mockappsec{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		createResponse := appsec.CreateCustomDenyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyCreateResponse.json"), &createResponse)
		require.NoError(t, err)
		createRequestJSON := loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheTrue.json")
		client.On("CreateCustomDeny",
			mock.Anything,
			appsec.CreateCustomDenyRequest{ConfigID: 43253, Version: 7, JsonPayloadRaw: createRequestJSON},
		).Return(&createResponse, nil)

		getResponse := appsec.GetCustomDenyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyGetResponse.json"), &getResponse)
		require.NoError(t, err)
		client.On("GetCustomDeny",
			mock.Anything,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponse, nil).Times(3)

		updateRequestJSON := loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheFalse.json")
		updateResponse := appsec.UpdateCustomDenyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyUpdateResponse.json"), &updateResponse)
		require.NoError(t, err)
		client.On("UpdateCustomDeny",
			mock.Anything,
			appsec.UpdateCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918", JsonPayloadRaw: updateRequestJSON},
		).Return(&updateResponse, nil)

		getResponseAfterUpdate := appsec.GetCustomDenyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyGetResponseAfterUpdate.json"), &getResponseAfterUpdate)
		require.NoError(t, err)
		client.On("GetCustomDeny",
			mock.Anything,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponseAfterUpdate, nil).Twice()

		removeResponse := appsec.RemoveCustomDenyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"), &removeResponse)
		require.NoError(t, err)
		client.On("RemoveCustomDeny",
			mock.Anything,
			appsec.RemoveCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/update_by_id.tf"),
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
