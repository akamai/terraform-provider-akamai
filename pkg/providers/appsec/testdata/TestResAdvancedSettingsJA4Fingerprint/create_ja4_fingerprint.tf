provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_ja4_fingerprint" "test" {
  config_id    = 111111
  header_names = ["ja4-fingerprint-after"]
}