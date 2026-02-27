provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_url_protection_action" "test" {
  config_id                 = 43253
  security_policy_id        = "AAAA_81230"
  url_protection_rule_id    = 135355
  max_rate_threshold_action = "invalid_action"
}
