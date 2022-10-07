provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_ip_geo" "test" {
  config_id                  = 43253
  security_policy_id         = "AAAA_81230"
  mode                       = "allow"
  geo_network_lists          = ["40731_BMROLLOUTGEO", "44831_ECSCGEOBLACKLIST"]
  ip_network_lists           = ["49181_ADTIPBLACKLIST", "49185_ADTWAFBYPASSLIST"]
  exception_ip_network_lists = ["69601_ADYENPRODWHITELIST", "68762_ADYEN"]
}

