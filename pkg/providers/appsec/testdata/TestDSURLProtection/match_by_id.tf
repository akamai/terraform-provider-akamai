provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_url_protection_policy" "test" {
  config_id                = 43007
  url_protection_policy_id = 681
}

