provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_ip_geo" "test" {
  config_id                  = 43253
  security_policy_id         = "AAAA_12345_TEST"
  mode                       = "allow"
  geo_network_lists          = [""]
  ip_network_lists           = ["", "TEST_BYPASS_LIST"]
  exception_ip_network_lists = ["", "TEST_WHITE_LIST"]
  asn_network_lists          = [""]
}