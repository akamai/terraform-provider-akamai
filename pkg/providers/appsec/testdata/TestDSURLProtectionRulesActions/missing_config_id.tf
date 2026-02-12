provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_url_protection_rules_actions" "test" {
  security_policy_id = "AAAA_81230"
}
