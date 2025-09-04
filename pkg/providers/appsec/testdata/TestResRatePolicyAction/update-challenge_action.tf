provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_rate_policy_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  rate_policy_id     = 135355
  ipv4_action        = "challenge_1234"
  ipv6_action        = "challenge_5678"
}