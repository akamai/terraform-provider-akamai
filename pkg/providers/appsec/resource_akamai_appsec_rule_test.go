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

		updateRuleResponse := appsec.UpdateRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRule/Rule.json"), &updateRuleResponse)

		getRuleResponse := appsec.GetRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRule/Rule.json"), &getRuleResponse)

		deleteRuleResponse := appsec.UpdateRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRule/Rule.json"), &deleteRuleResponse)

		getWAFModeResponse := appsec.GetWAFModeResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResWAFMode/WAFMode.json"), &getWAFModeResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetRule",
			mock.Anything,
			appsec.GetRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&getRuleResponse, nil)

		client.On("GetWAFMode",
			mock.Anything,
			appsec.GetWAFModeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getWAFModeResponse, nil)

		conditionExceptionJSON := loadFixtureBytes("testdata/TestResRule/ConditionException.json")
		client.On("UpdateRule",
			mock.Anything,
			appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", RuleID: 12345, JsonPayloadRaw: conditionExceptionJSON},
		).Return(&updateRuleResponse, nil)

		client.On("UpdateRule",
			mock.Anything,
			appsec.UpdateRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Action: "none"},
		).Return(&deleteRuleResponse, nil)

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

		updateConditionExceptionResponse := appsec.UpdateConditionExceptionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRule/RuleAseAuto.json"), &updateConditionExceptionResponse)

		getRuleResponse := appsec.GetRuleResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRule/RuleAseAuto.json"), &getRuleResponse)

		deleteConditionExceptionResponse := appsec.UpdateConditionExceptionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResRule/RuleAseAuto.json"), &deleteConditionExceptionResponse)

		getWAFModeResponse := appsec.GetWAFModeResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResWAFMode/WAFModeAseAuto.json"), &getWAFModeResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetRule",
			mock.Anything,
			appsec.GetRuleRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&getRuleResponse, nil)

		client.On("GetWAFMode",
			mock.Anything,
			appsec.GetWAFModeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getWAFModeResponse, nil)

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
			mock.Anything,
			appsec.UpdateConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345, Conditions: &conditions, Exception: &exception},
		).Return(&updateConditionExceptionResponse, nil)

		client.On("UpdateRuleConditionException",
			mock.Anything,
			appsec.UpdateConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RuleID: 12345},
		).Return(&updateConditionExceptionResponse, nil)

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
