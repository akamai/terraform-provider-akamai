provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_request_body" "test" {
  config_id                     = 43253
  request_body_inspection_limit = "16"
}