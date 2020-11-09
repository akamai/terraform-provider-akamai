package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroupConditionExceptions_data_basic(t *testing.T) {
	t.Run("match by AttackGroupConditionExceptions ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetAttackGroupConditionExceptionsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSAttackGroupConditionExceptions/AttackGroupConditionExceptions.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetAttackGroupConditionExceptions",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupConditionExceptionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAttackGroupConditionExceptions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_aag_rules.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
