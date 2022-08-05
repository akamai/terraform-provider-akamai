provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_rate_policy_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  rate_policy_id     = 135355
  ipv4_action        = "none"
  ipv6_action        = "none"
}

