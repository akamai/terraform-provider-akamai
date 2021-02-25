package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiFailoverHostnames_data_basic(t *testing.T) {
	t.Run("match by FailoverHostnames ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetFailoverHostnamesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSFailoverHostnames/FailoverHostnames.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetFailoverHostnames",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetFailoverHostnamesRequest{ConfigID: 43253},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSFailoverHostnames/match_by_id.tf"),
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
