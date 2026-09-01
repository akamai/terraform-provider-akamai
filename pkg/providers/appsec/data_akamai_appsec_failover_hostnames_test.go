package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiFailoverHostnames_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by FailoverHostnames ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getFailoverHostnamesResponse := appsec.GetFailoverHostnamesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSFailoverHostnames/FailoverHostnames.json"), &getFailoverHostnamesResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetFailoverHostnames",
			testutils.MockContext,
			appsec.GetFailoverHostnamesRequest{ConfigID: 43253},
		).Return(&getFailoverHostnamesResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSFailoverHostnames/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_failover_hostnames.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
