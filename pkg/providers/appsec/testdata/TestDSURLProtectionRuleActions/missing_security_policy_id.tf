provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_url_protection_rule_actions" "test" {
  config_id         = 43253
  url_protection_id = 135355
}
