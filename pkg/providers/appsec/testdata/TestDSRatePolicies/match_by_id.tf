provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_rate_policies" "test" {
  config_id = 43253
}

