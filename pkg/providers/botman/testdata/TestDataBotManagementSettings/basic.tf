provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_bot_management_settings" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}