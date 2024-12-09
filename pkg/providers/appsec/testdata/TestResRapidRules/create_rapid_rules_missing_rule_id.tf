provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_rapid_rules" "test" {
  config_id          = 111111
  security_policy_id = "2222_333333"
  default_action     = "alert"
  rule_definitions   = file("testdata/TestResRapidRules/RuleDefinitionsWithMissingRuleID.json")
}
