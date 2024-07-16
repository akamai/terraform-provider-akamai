provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_content_protection_javascript_injection_rule" "test" {
  security_policy_id = "AAAA_81230"
}
