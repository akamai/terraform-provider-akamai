package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalProtectHost_res_basic(t *testing.T) {
	t.Run("match by EvalProtectHost ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateEvalProtectHostResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResEvalProtectHost/EvalProtectHost.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetEvalProtectHostsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResEvalProtectHost/EvalProtectHost.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalProtectHosts",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalProtectHostsRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateEvalProtectHost",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateEvalProtectHostRequest{ConfigID: 43253, Version: 7, Hostnames: []string{"example.com"}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalProtectHost/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_protect_host.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
