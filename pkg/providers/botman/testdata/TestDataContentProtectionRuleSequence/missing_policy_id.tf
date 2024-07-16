provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_content_protection_rule_sequence" "test" {
  config_id = 43253
}
