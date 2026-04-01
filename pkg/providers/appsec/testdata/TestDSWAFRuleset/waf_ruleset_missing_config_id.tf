provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_waf_ruleset" "test" {
  security_policy_id = "2222_333333"
}
