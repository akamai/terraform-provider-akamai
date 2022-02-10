provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_configuration_version" "test" {
  config_id = 43253
  version   = 7

}

