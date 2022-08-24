provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_bot_detection_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  detection_id       = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
}