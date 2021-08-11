package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRule_res_basic(t *testing.T) {
	t.Run("match by Rule ID", func(t *testing.T) {
		client := &mockappsec{}

		updResp := appsec.UpdateRuleResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRule/Rule.json"))
		json.Unmarshal([]byte(expectJSU), &updResp)

		getResp := appsec.GetRuleResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRule/Rule.json"))
		json.Unmarshal([]byte(expectJS), &getResp)

		delResp := appsec.UpdateRuleResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResRule/Rule.json"))
		json.Unmarshal([]byte(expectJSD), &delResp)

		wm := appsec.GetWAFModeResponse{}
		expectWM := compactJSON(loadFixtureBytes("testdata/TestResWAFMode/WAFMode.json"))
		json.Unmarshal([]byte(expectWM), &wm)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&getResp, nil)

		client.On("GetWAFMode",
			mock.Anything,
			appsec.GetWAFModeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&wm, nil)

		conditionExceptionJSON := loadFixtureBytes("testdata/TestResRule/ConditionException.json")
		client.On("UpdateRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", RuleID: 12345, JsonPayloadRaw: conditionExceptionJSON},
		).Return(&updResp, nil)

		/*client.On("UpdateRule",
		mock.Anything, // ctx is irrelevant for this test
		appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345}).Return(&updateRuleResponse, nil)*/

		client.On("UpdateRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Action: "none"},
		).Return(&delResp, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRule/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rule.test", "id", "43253:AAAA_81230:12345"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAccAkamaiRule_res_AseAuto(t *testing.T) {
	t.Run("match by Rule ID", func(t *testing.T) {
		client := &mockappsec{}

		updResp := appsec.UpdateConditionExceptionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRule/RuleAseAuto.json"))
		json.Unmarshal([]byte(expectJSU), &updResp)

		getResp := appsec.GetRuleResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRule/RuleAseAuto.json"))
		json.Unmarshal([]byte(expectJS), &getResp)

		delResp := appsec.UpdateConditionExceptionResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResRule/RuleAseAuto.json"))
		json.Unmarshal([]byte(expectJSD), &delResp)

		wm := appsec.GetWAFModeResponse{}
		expectWM := compactJSON(loadFixtureBytes("testdata/TestResWAFMode/WAFModeAseAuto.json"))
		json.Unmarshal([]byte(expectWM), &wm)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetRule",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&getResp, nil)

		client.On("GetWAFMode",
			mock.Anything,
			appsec.GetWAFModeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&wm, nil)

		conditions := appsec.RuleConditions{
			{
				Type:          "extensionMatch",
				Extensions:    []string{"test"},
				PositiveMatch: true,
			}}

		exception := appsec.RuleException{
			HeaderCookieOrParamValues: []string{"test"},
		}

		client.On("UpdateRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Conditions: &conditions, Exception: &exception},
		).Return(&updResp, nil)

		/*client.On("UpdateRule",
		mock.Anything, // ctx is irrelevant for this test
		appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345}).Return(&updateRuleResponse, nil)*/

		client.On("UpdateRuleConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&updResp, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRule/match_by_id_aseauto.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rule.test", "id", "43253:AAAA_81230:12345"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
