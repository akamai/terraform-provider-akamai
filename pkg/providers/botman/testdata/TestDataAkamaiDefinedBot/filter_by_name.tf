provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_akamai_defined_bot" "test" {
  bot_name = "Test name 3"
}