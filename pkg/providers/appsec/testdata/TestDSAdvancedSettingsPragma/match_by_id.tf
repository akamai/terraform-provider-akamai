provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_advanced_settings_pragma_header" "test" {
  config_id = 43253
}

