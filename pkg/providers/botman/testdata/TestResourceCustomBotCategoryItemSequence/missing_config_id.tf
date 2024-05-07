provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_bot_category_item_sequence" "test" {
  category_id = "fakecv20-eddb-4421-93d9-90954e509d5f"
  bot_ids     = ["fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"]
}
