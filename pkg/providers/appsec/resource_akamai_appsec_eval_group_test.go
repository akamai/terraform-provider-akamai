package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalGroup_res_basic(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &mockappsec{}

		updResp := appsec.UpdateAttackGroupResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResEvalGroup/AttackGroup.json"))
		json.Unmarshal([]byte(expectJSU), &updResp)

		getResp := appsec.GetAttackGroupResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResEvalGroup/AttackGroup.json"))
		json.Unmarshal([]byte(expectJS), &getResp)

		delResp := appsec.UpdateAttackGroupResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResEvalGroup/AttackGroup.json"))
		json.Unmarshal([]byte(expectJSD), &delResp)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalGroup",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&getResp, nil)

		data := `{"conditions":[],"exception":{"headerCookieOrParamValues":["abc"],"specificHeaderCookieOrParamPrefix": {"prefix": "a*","selector": "REQUEST_COOKIES"}}}` + "\n"
		client.On("UpdateEvalGroup",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: []byte(data)},
		).Return(&updResp, nil)

		client.On("UpdateEvalGroup",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "none"},
		).Return(&delResp, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalGroup/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_group.test", "id", "43253:AAAA_81230:SQL"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAccAkamaiEvalGroup_res_error_updating_eval_group(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		data := `{"conditions":[],"exception":{"headerCookieOrParamValues":["abc"],"specificHeaderCookieOrParamPrefix": {"prefix": "a*","selector": "REQUEST_COOKIES"}}}` + "\n"
		client.On("UpdateEvalGroup",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: []byte(data)},
		).Return(nil, fmt.Errorf("UpdateEvalGroup request failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalGroup/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_group.test", "id", "43253:AAAA_81230:SQL"),
						),
						ExpectError:        regexp.MustCompile(`UpdateEvalGroup request failed`),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
