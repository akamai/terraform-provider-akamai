provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_waf_ruleset" "test" {
  config_id          = 111111
  security_policy_id = "2222_333333"

  rules = [
    {
      rule_id     = 999999
      rule_action = "alert"
    }
  ]
}
