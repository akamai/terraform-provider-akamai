provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_siem_settings" "test" {
  config_id = 43253
    version = 7
}


