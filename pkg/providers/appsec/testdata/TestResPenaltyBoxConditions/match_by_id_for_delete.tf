provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_penalty_box_conditions" "delete_condition" {
  config_id              = 43253
  security_policy_id     = "AAAA"
  penalty_box_conditions = file("testdata/TestResPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json")
}