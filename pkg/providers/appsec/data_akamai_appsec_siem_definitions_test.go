package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSiemDefinitions_data_basic(t *testing.T) {
	t.Run("match by SiemDefinitions ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetSiemDefinitionsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSSiemDefinitions/SiemDefinitions.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetSiemDefinitions",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSiemDefinitionsRequest{SiemDefinitionName: "SIEM Version 01"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSSiemDefinitions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_siem_definitions.test", "id", "1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
