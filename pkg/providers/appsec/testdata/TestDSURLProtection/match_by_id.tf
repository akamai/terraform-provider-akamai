provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_url_protection_rule" "test" {
  config_id              = 43007
  url_protection_rule_id = 681
}

