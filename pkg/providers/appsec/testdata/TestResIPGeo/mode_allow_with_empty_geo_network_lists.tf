provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_ip_geo" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_12345_TEST"
  mode               = "allow"
  geo_controls {
    action            = "deny"
    geo_network_lists = [""]
  }
  ip_controls {
    action           = "deny"
    ip_network_lists = ["TEST_IP_BLACKLIST", "TEST_WAF_BYPASS_LIST"]
  }
  exception_ip_network_lists = ["TEST_PROD_WHITELIST"]
}