provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_conditional_action" "test" {
  config_id = 43253
}