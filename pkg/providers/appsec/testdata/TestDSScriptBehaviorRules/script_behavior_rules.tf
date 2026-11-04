provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_script_behavior_rules" "test" {
  config_id          = 40576
  security_policy_id = "ab_86352"
}