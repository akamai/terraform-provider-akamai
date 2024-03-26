provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_custom_defined_bot" "test" {
  config_id = 43253
}