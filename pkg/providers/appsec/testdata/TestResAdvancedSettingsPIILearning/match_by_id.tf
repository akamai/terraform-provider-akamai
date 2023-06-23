provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_pii_learning" "test" {
  config_id           = 43253
  enable_pii_learning = true
}

