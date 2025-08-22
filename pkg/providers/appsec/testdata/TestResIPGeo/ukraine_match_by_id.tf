provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_ip_geo" "test1" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  mode               = "block"
  asn_controls {
    action            = "deny"
    asn_network_lists = ["40721_ASNLIST1", "44811_ASNLIST2"]
  }
  geo_controls {
    action            = "deny"
    geo_network_lists = ["40731_BMROLLOUTGEO", "44831_ECSCGEOBLACKLIST"]
  }
  ip_controls {
    action           = "deny"
    ip_network_lists = ["49181_ADTIPBLACKLIST", "49185_ADTWAFBYPASSLIST"]
  }
  exception_ip_network_lists = ["68762_ADYEN", "69601_ADYENPRODWHITELIST"]
  ukraine_geo_control_action = "alert"
}

