provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_clientside_protection_configuration" "test" {
  config_id = 43253
}
