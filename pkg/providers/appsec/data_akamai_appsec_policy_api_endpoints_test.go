package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiPolicyApiEndpoints_data_basic(t *testing.T) {
	t.Run("match by ApiEndpoints ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetPolicyApiEndpointsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSPolicyApiEndpoints/PolicyApiEndpoints.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetApiEndpoints",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetPolicyApiEndpointsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ID: 619183},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSPolicyApiEndpoints/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_policy_api_endpoints.test", "id", "619183"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
