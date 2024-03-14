provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_penalty_box_conditions" "test" {
  config_id              = 43253
  security_policy_id     = "AAAA_81230"
  penalty_box_conditions = file("testdata/TestResPenaltyBoxConditions/PenaltyBoxConditions.json")
}

