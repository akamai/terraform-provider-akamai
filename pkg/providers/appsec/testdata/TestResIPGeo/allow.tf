provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_ip_geo" "test" {
  config_id                  = 43253
  security_policy_id         = "AAAA_81230"
  mode                       = "allow"
  block_action               = "deny"
  exception_ip_network_lists = ["68762_ADYEN", "69601_ADYENPRODWHITELIST"]
}

