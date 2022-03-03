package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSiemSettings_res_basic(t *testing.T) {
	t.Run("match by SiemSettings ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateSiemSettingsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSiemSettings/SiemSettings.json")), &cu)

		cr := appsec.GetSiemSettingsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSiemSettings/SiemSettings.json")), &cr)

		crd := appsec.RemoveSiemSettingsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSiemSettings/SiemSettings.json")), &crd)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSiemSettings",
			mock.Anything,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateSiemSettings",
			mock.Anything,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, EnabledBotmanSiemEvents: true, SiemDefinitionID: 1, FirewallPolicyIds: []string{"12345"}},
		).Return(&cu, nil)

		client.On("RemoveSiemSettings",
			mock.Anything,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: false, EnabledBotmanSiemEvents: false, FirewallPolicyIds: []string(nil)},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSiemSettings/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
