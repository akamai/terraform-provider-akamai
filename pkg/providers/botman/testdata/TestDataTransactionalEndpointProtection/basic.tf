provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_transactional_endpoint_protection" "test" {
  config_id = 43253
}