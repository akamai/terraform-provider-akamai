package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRules_data_basic(t *testing.T) {
	t.Run("match by Rules ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetRulesResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestDSRules/Rules.json")), &cv)

		configs := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &configs)

		client.On("GetRules",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRulesRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cv, nil)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configs, nil)

		wm := appsec.GetWAFModeResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResWAFMode/WAFMode.json")), &wm)

		client.On("GetWAFMode",
			mock.Anything,
			appsec.GetWAFModeRequest{
				ConfigID: 43253,
				Version:  7,
				PolicyID: "AAAA_81230",
			}).Return(&wm, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRules/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_rules.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
