provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_script_behavior_rule" "test" {
  config_id = 40576
  rule_id   = 3000601
}
