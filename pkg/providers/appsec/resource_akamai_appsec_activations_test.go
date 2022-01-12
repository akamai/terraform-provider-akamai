package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiActivations_res_basic(t *testing.T) {
	t.Run("match by Activations ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.RemoveActivationsResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResActivations/ActivationsDelete.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		ga := appsec.GetActivationsResponse{}
		expectJSR := compactJSON(loadFixtureBytes("testdata/TestResActivations/Activations.json"))
		json.Unmarshal([]byte(expectJSR), &ga)

		cr := appsec.CreateActivationsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResActivations/Activations.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetActivations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&ga, nil)

		client.On("CreateActivations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateActivationsRequest{Action: "ACTIVATE", Network: "STAGING", Note: "", NotificationEmails: []string{"user@example.com"}, ActivationConfigs: []struct {
				ConfigID      int "json:\"configId\""
				ConfigVersion int "json:\"configVersion\""
			}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&cr, nil)

		client.On("RemoveActivations",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveActivationsRequest{ActivationID: 547694, Action: "DEACTIVATE", Network: "STAGING", Note: "", NotificationEmails: []string{"user@example.com"}, ActivationConfigs: []struct {
				ConfigID      int "json:\"configId\""
				ConfigVersion int "json:\"configVersion\""
			}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
