package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomDeny_data_basic(t *testing.T) {
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetCustomDenyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSCustomDeny/CustomDeny.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_custom_deny.test", "custom_deny_id", "deny_custom_54994"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
