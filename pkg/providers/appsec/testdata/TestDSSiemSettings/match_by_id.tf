provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_siem_settings" "test" {
  config_id = 43253
}

