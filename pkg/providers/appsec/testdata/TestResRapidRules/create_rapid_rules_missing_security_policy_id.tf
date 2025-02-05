provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_rapid_rules" "test" {
  config_id        = 111111
  default_action   = "deny"
  rule_definitions = file("testdata/TestResRapidRules/RuleDefinitions.json")
}
