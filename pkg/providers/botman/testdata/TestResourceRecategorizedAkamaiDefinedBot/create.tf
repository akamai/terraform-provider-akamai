provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_recategorized_akamai_defined_bot" "test" {
  config_id   = 43253
  bot_id      = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
  category_id = "87fb601b-4d30-4e0d-a74f-dc77e2b1bb74"
}