provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_advanced_settings_logging" "test" {
  config_id = 43253
}

