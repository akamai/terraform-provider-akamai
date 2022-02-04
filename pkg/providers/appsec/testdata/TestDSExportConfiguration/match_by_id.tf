provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_export_configuration" "test" {
  config_id = 43253
  version   = 7
  search    = ["ruleActions", "customRules", "rulesets", "reputationProfiles", "ratePolicies", "matchTargets"]
}

