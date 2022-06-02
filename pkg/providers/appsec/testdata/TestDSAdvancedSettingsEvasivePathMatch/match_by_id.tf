provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_advanced_settings_evasive_path_match" "test" {
  config_id = 43253
}

