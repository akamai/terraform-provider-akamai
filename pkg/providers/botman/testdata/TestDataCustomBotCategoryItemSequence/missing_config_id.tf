provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_custom_bot_category_item_sequence" "test" {
  category_id = "fakecv20-eddb-4421-93d9-90954e509d5f"
}
