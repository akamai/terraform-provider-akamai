provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_siem_settings" "siem_settings" {
  config_id = 43253
    version = 7
}


