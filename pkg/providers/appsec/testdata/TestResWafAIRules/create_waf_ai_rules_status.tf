provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_waf_ai_rules" "test" {
  config_id          = 111111
  security_policy_id = "2222_333333"
  ai_rule_status     = "ENABLED"
}
