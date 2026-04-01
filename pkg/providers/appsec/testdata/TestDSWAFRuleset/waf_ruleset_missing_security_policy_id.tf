provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_waf_ruleset" "test" {
  config_id = 111111
}
