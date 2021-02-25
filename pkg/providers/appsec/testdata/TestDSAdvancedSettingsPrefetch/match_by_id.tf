provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_advanced_settings_prefetch" "test" {
  config_id = 43253
    version = 7
}
