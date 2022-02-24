package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiApiRequestConstraints_res_basic(t *testing.T) {
	t.Run("match by ApiRequestConstraints ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateApiRequestConstraintsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json")), &cu)

		cr := appsec.GetApiRequestConstraintsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json")), &cr)

		crd := appsec.RemoveApiRequestConstraintsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json")), &crd)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetApiRequestConstraints",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1},
		).Return(&cr, nil)

		client.On("UpdateApiRequestConstraints",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1, Action: "alert"},
		).Return(&cu, nil)

		client.On("RemoveApiRequestConstraints",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1, Action: "none"},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResApiRequestConstraints/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_api_request_constraints.test", "id", "43253:AAAA_81230:1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
