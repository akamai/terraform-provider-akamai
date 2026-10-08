provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_script_behavior_rule" "test" {
  security_policy_id = "ab_86352"
  rule_id            = 3000601
}
