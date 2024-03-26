provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_custom_client" "test" {
  config_id = 43253
}