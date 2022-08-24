provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_recategorized_akamai_defined_bot" "test" {
  config_id = 43253
}