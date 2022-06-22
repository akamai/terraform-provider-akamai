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

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)

		getActivationsResponse := appsec.GetActivationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResActivations/Activations.json"), &getActivationsResponse)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResActivations/Activations.json"), &createActivationsResponse)

		client.On("GetActivations",
			mock.Anything,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil)

		client.On("CreateActivations",
			mock.Anything,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil)

		client.On("RemoveActivations",
			mock.Anything,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil)

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
