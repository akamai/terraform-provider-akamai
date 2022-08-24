provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_transactional_endpoint_protection" "test" {
  config_id = 43253
}