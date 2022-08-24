provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_akamai_bot_category" "test" {
  category_name = "Test name 3"
}