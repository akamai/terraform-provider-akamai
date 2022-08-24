provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_recategorized_akamai_defined_bot" "test" {
  config_id   = 43253
  bot_id      = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
  category_id = "c43b638c-8f9a-4ea3-b1bd-3c82c96fefbf"
}