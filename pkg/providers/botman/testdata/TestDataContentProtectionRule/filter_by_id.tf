provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_content_protection_rule" "test" {
  config_id                  = 43253
  security_policy_id         = "AAAA_81230"
  content_protection_rule_id = "fake3f89-e179-4892-89cf-d5e623ba9dc7"
}
