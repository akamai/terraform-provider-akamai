provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_advanced_settings_pii_learning" "test" {
  config_id = 43253
}

