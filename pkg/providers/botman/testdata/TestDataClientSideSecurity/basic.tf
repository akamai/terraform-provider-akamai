provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_client_side_security" "test" {
  config_id = 43253
}