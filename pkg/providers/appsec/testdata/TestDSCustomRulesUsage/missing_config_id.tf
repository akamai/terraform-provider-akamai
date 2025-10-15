provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_custom_rules_usage" "test" {
  rule_ids = [12345]
}
