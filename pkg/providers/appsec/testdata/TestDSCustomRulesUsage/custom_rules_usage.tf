provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_custom_rules_usage" "test" {
  config_id = 111111
  rule_ids  = [12345, 67890]
}
