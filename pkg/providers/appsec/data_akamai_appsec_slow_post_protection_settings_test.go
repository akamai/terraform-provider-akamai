package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSlowPostProtectionSettings_data_basic(t *testing.T) {
	t.Run("match by SlowPostProtectionSettings ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetSlowPostProtectionSettingsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSSlowPostProtectionSettings/SlowPostProtectionSettings.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetSlowPostProtectionSettings",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSlowPostProtectionSettingsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSSlowPostProtectionSettings/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_slow_post.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
