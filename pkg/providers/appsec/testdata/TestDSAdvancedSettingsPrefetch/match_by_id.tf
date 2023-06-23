provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_advanced_settings_prefetch" "test" {
  config_id = 43253
}

