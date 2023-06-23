provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_akamai_bot_category_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}