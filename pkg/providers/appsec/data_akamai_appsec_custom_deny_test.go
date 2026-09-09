package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomDeny_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getCustomDenyListResponse := appsec.GetCustomDenyListResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSCustomDeny/CustomDenyList.json"), &getCustomDenyListResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetCustomDenyList",
			testutils.MockContext,
			appsec.GetCustomDenyListRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_54994"},
		).Return(&getCustomDenyListResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSCustomDeny/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_custom_deny.test", "custom_deny_id", "deny_custom_54994"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
