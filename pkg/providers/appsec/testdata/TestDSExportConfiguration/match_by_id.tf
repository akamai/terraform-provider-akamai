provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_export_configuration" "test" {
  config_id = 43253
  version   = 7
  search    = ["ruleActions", "customRules", "rulesets", "reputationProfiles", "ratePolicies", "matchTargets"]
}

data "akamai_appsec_export_configuration" "evalGroups" {
  config_id = 43253
  version   = 7
  search    = ["EvalGroup.tf"]
}

