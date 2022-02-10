provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_custom_rule_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  custom_rule_id     = 60036362
  custom_rule_action = "none"
}

