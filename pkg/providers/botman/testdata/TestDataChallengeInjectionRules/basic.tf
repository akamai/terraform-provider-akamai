provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_challenge_injection_rules" "test" {
  config_id = 43253
}