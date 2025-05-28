package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiRatePolicy_res_basic(t *testing.T) {
	client := &appsec.Mock{}

	t.Run("match by RatePolicy ID", func(t *testing.T) {
		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		createResponse := appsec.CreateRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/RatePolicy.json"), &createResponse)
		require.NoError(t, err)
		createRatePolicyJSON := testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/CreateRatePolicy.json")
		client.On("CreateRatePolicy",
			testutils.MockContext,
			appsec.CreateRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createRatePolicyJSON},
		).Return(&createResponse, nil)

		getResponse := appsec.GetRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/RatePolicy.json"), &getResponse)
		require.NoError(t, err)
		client.On("GetRatePolicy",
			testutils.MockContext,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponse, nil).Once()
		client.On("GetRatePolicy",
			testutils.MockContext,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponse, nil).Once()
		client.On("GetRatePolicy",
			testutils.MockContext,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponse, nil).Once()

		updateResponse := appsec.UpdateRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/RatePolicyUpdated.json"), &updateResponse)
		require.NoError(t, err)
		updateRatePolicyJSON := testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/UpdateRatePolicy.json")
		client.On("UpdateRatePolicy",
			testutils.MockContext,
			appsec.UpdateRatePolicyRequest{RatePolicyID: 134644, ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: updateRatePolicyJSON},
		).Return(&updateResponse, nil)

		getResponseAfterUpdate := appsec.GetRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/RatePolicyUpdated.json"), &getResponseAfterUpdate)
		require.NoError(t, err)
		client.On("GetRatePolicy",
			testutils.MockContext,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponseAfterUpdate, nil).Twice()

		removeResponse := appsec.RemoveRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/RatePolicyEmpty.json"), &removeResponse)
		require.NoError(t, err)
		client.On("RemoveRatePolicy",
			testutils.MockContext,
			appsec.RemoveRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRatePolicy/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "43253:134644"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRatePolicy/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "43253:134644"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("suppress counterType and penaltyBoxDuration attributes", func(t *testing.T) {
		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		createResponse := appsec.CreateRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/suppress/RatePolicyResponse.json"), &createResponse)
		require.NoError(t, err)
		createRatePolicyJSON := testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/suppress/CreateRatePolicy.json")
		client.On("CreateRatePolicy",
			testutils.MockContext,
			appsec.CreateRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createRatePolicyJSON},
		).Return(&createResponse, nil)

		getResponse := appsec.GetRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/suppress/RatePolicyResponse.json"), &getResponse)
		require.NoError(t, err)
		client.On("GetRatePolicy",
			testutils.MockContext,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponse, nil).Times(3)

		updateResponse := appsec.UpdateRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/suppress/RatePolicyUpdated.json"), &updateResponse)
		require.NoError(t, err)
		updateRatePolicyJSON := testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/suppress/UpdateRatePolicy.json")
		client.On("UpdateRatePolicy",
			testutils.MockContext,
			appsec.UpdateRatePolicyRequest{RatePolicyID: 134644, ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: updateRatePolicyJSON},
		).Return(&updateResponse, nil)

		getResponseAfterUpdate := appsec.GetRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/suppress/RatePolicyUpdated.json"), &getResponseAfterUpdate)
		require.NoError(t, err)
		client.On("GetRatePolicy",
			testutils.MockContext,
			appsec.GetRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&getResponseAfterUpdate, nil).Twice()

		removeResponse := appsec.RemoveRatePolicyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicy/suppress/RatePolicyEmpty.json"), &removeResponse)
		require.NoError(t, err)
		client.On("RemoveRatePolicy",
			testutils.MockContext,
			appsec.RemoveRatePolicyRequest{ConfigID: 43253, ConfigVersion: 7, RatePolicyID: 134644},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRatePolicy/suppress/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "43253:134644"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRatePolicy/suppress/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy.test", "id", "43253:134644"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
