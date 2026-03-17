provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_url_protection_policy_actions" "test" {
  config_id                = 43253
  security_policy_id       = "AAAA_81230"
  url_protection_policy_id = 999999
}
