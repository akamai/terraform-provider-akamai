provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_url_evasion_defense" "test" {
  config_id = 43253
  status    = "invalid"
}

