package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroups_data_basic(t *testing.T) {
	t.Run("match by AttackGroups ID", func(t *testing.T) {
		client := &mockappsec{}

		getAttackGroupsResponse := appsec.GetAttackGroupsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSAttackGroups/AttackGroups.json"), &getAttackGroupsResponse)

		configs := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &configs)

		client.On("GetAttackGroups",
			mock.Anything,
			appsec.GetAttackGroupsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&getAttackGroupsResponse, nil)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configs, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAttackGroups/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_attack_groups.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
