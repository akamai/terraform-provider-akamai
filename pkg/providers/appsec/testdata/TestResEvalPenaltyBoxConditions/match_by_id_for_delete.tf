provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_eval_penalty_box_conditions" "delete_condition" {
  config_id              = 43253
  security_policy_id     = "AAAA"
  penalty_box_conditions = file("testdata/TestResEvalPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json")
}