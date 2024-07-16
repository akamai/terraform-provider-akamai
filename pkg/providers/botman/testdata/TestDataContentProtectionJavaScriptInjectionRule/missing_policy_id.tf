provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_content_protection_javascript_injection_rule" "test" {
  config_id = 43253
}
