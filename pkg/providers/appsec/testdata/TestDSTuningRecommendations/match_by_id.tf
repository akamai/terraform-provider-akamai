provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_tuning_recommendations" "recommendations" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}

data "akamai_appsec_tuning_recommendations" "group_recommendations" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  attack_group       = "XSS"
  ruleset_type       = "evaluation"
}
