provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_content_protection_rule_sequence" "test" {
  config_id                   = 43253
  security_policy_id          = "AAAA_81230"
  content_protection_rule_ids = ["fake3f89-e179-4892-89cf-d5e623ba9dc7", "fake85df-e399-43e8-bb0f-c0d980a88e4f", "fake09b8-4fd5-430e-a061-1c61df1d2ac2"]
}
