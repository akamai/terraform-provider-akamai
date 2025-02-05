provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_rapid_rules" "test" {
  security_policy_id = "2222_333333"
  default_action     = "deny"
  rule_definitions   = file("testdata/TestResRapidRules/RuleDefinitions.json")
}
