package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiFailoverHostnames_data_basic(t *testing.T) {
	t.Run("match by FailoverHostnames ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getFailoverHostnamesResponse := appsec.GetFailoverHostnamesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSFailoverHostnames/FailoverHostnames.json"), &getFailoverHostnamesResponse)
		require.NoError(t, err)

		client.On("GetFailoverHostnames",
			testutils.MockContext,
			appsec.GetFailoverHostnamesRequest{ConfigID: 43253},
		).Return(&getFailoverHostnamesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSFailoverHostnames/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_failover_hostnames.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
