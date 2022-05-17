provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_api_endpoints" "test" {
  config_id = 43253
}

